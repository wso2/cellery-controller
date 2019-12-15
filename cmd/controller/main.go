/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"flag"
	"log"
	"time"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"cellery.io/cellery-controller/pkg/clients"
	"cellery.io/cellery-controller/pkg/config"
	"cellery.io/cellery-controller/pkg/controller/cell"
	"cellery.io/cellery-controller/pkg/controller/component"
	"cellery.io/cellery-controller/pkg/controller/composite"
	"cellery.io/cellery-controller/pkg/controller/gateway"
	"cellery.io/cellery-controller/pkg/controller/sts"
	"cellery.io/cellery-controller/pkg/informers"
	"cellery.io/cellery-controller/pkg/logging"
	"cellery.io/cellery-controller/pkg/signals"
	"cellery.io/cellery-controller/pkg/version"
)

const (
	threadsPerController = 2
	componentName        = "Controller"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	logger, err := logging.NewLogger()
	if err != nil {
		log.Fatalf("Error building logger: %s", err.Error())
	}
	defer logger.Sync()

	logger.Info(version.String(componentName))

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		logger.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	clientset, err := clients.NewFromConfig(cfg)

	if err != nil {
		logger.Fatalf("Error building clients: %v", err)
	}

	// Create required informers
	informerset := informers.New(clientset, time.Second*60)

	// Create config watcher
	cw := config.NewWatcher(informerset, "cellery-config", "cellery-secret", "cellery-system", logger)

	// Create crd controllers
	gatewayController := gateway.NewController(
		clientset,
		informerset,
		cw,
		logger,
	)

	componentController := component.NewController(
		clientset,
		informerset,
		cw,
		logger,
	)

	tokenServiceController := sts.NewController(
		clientset,
		informerset,
		cw,
		logger,
	)

	cellController := cell.NewController(
		clientset,
		informerset,
		cw,
		logger,
	)

	compositeController := composite.NewController(
		clientset,
		informerset,
		cw,
		logger,
	)
	// Start informers and wait for caches to sync
	logger.Info("Starting informers...")
	err = informerset.Start(stopCh)
	if err != nil {
		logger.Fatalf("Error waiting for cache sync: %v", err)
	}

	err = cw.CheckResources()
	if err != nil {
		logger.Fatalf("Error checking config resources: %v", err)
	}

	//Start controllers
	logger.Info("Starting controllers...")
	go gatewayController.Run(threadsPerController, stopCh)
	go componentController.Run(threadsPerController, stopCh)
	go tokenServiceController.Run(threadsPerController, stopCh)
	go cellController.Run(threadsPerController, stopCh)
	go compositeController.Run(threadsPerController, stopCh)

	// Prevent exiting the main process
	<-stopCh
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
