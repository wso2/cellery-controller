/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"cellery.io/cellery-controller/pkg/clients"
	"cellery.io/cellery-controller/pkg/version"

	"cellery.io/cellery-controller/pkg/logging"
	"cellery.io/cellery-controller/pkg/signals"
	"cellery.io/cellery-controller/pkg/webhook"
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

	logger.Info(version.String())

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		logger.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	clientset, err := clients.NewFromConfig(cfg)

	if err != nil {
		logger.Fatalf("Error building clients: %v", err)
	}

	opt := webhook.ServerOptions{
		Namespace:             "cellery-system",
		ServerSecretName:      "webhook-certs",
		RootSecretName:        "cellery-secret",
		ServiceName:           "webhook",
		DeploymentName:        "webhook",
		MutatingWebhookName:   "defaulting.mesh.cellery.io",
		ValidatingWebhookName: "validating.mesh.cellery.io",
		Port:                  8443,
	}

	server := webhook.NewServer(clientset.Kubernetes(), opt, logger)

	if err = server.Run(stopCh); err != nil {
		logger.Fatalf("Failed to run the admission webhook: %v", err)
	}

}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
