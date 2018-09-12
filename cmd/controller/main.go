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
	"k8s.io/client-go/tools/cache"
	"time"

	"github.com/golang/glog"
	vickclientset "github.com/wso2/product-vick/pkg/client/clientset/versioned"
	vickinformers "github.com/wso2/product-vick/pkg/client/informers/externalversions"
	"github.com/wso2/product-vick/pkg/controller/cell"
	"github.com/wso2/product-vick/pkg/controller/service"
	"github.com/wso2/product-vick/pkg/signals"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	threadsPerController = 2
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	vickClient, err := vickclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building vick clientset: %s", err.Error())
	}

	// Create informers
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	vickInformerFactory := vickinformers.NewSharedInformerFactory(vickClient, time.Second*30)

	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	k8sServiceInformer := kubeInformerFactory.Core().V1().Services()
	serviceInformer := vickInformerFactory.Vick().V1alpha1().Services()
	cellInformer := vickInformerFactory.Vick().V1alpha1().Cells()

	// Create crd controllers
	cellController := cell.NewController(kubeClient, cellInformer)
	serviceController := service.NewController(kubeClient, vickClient, k8sServiceInformer, cellInformer, serviceInformer, deploymentInformer)

	// Start informers
	go kubeInformerFactory.Start(stopCh)
	go vickInformerFactory.Start(stopCh)

	// Wait for cache sync
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh,
		deploymentInformer.Informer().HasSynced,
		k8sServiceInformer.Informer().HasSynced,
		cellInformer.Informer().HasSynced,
		serviceInformer.Informer().HasSynced); !ok {
		glog.Fatal("failed to wait for caches to sync")
	}

	//Start controllers
	go cellController.Run(threadsPerController, stopCh)
	go serviceController.Run(threadsPerController, stopCh)

	// Prevent exiting the main process
	<-stopCh
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
