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

package clients

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	meshclient "github.com/cellery-io/mesh-controller/pkg/generated/clientset/versioned"
)

type Interface interface {
	Kubernetes() kubeclient.Interface
	Mesh() meshclient.Interface
}

type clients struct {
	kubeClient kubeclient.Interface
	meshClient meshclient.Interface
}

func NewFromConfig(cfg *rest.Config) (*clients, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("error building kubernetes clientset: %s", err.Error())
	}

	meshClient, err := meshclient.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("error building mesh clientset: %s", err.Error())
	}
	return &clients{
		kubeClient: kubeClient,
		meshClient: meshClient,
	}, nil
}

func NewFromClients(kubeClient kubeclient.Interface, meshClient meshclient.Interface) *clients {
	return &clients{
		kubeClient: kubeClient,
		meshClient: meshClient,
	}
}

func (c *clients) Kubernetes() kubeclient.Interface {
	return c.kubeClient
}

func (c *clients) Mesh() meshclient.Interface {
	return c.meshClient
}
