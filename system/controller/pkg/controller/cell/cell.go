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

package cell

import (
	"fmt"
	"github.com/wso2/product-vick/system/controller/pkg/controller/cell/resources"

	"github.com/golang/glog"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	//appsv1informers "k8s.io/client-go/informers/apps/v1"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	//corev1informers "k8s.io/client-go/informers/core/v1"
	vickinformers "github.com/wso2/product-vick/system/controller/pkg/client/informers/externalversions/vick/v1alpha1"
	networkv1informers "k8s.io/client-go/informers/networking/v1"

	listers "github.com/wso2/product-vick/system/controller/pkg/client/listers/vick/v1alpha1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	networkv1listers "k8s.io/client-go/listers/networking/v1"
)

type cellHandler struct {
	cellLister          listers.CellLister
	deploymentLister    appsv1listers.DeploymentLister
	networkPilicyLister networkv1listers.NetworkPolicyLister
	k8sServiceLister    corev1listers.ServiceLister
	kubeClient          kubernetes.Interface
}

func NewController(kubeClient kubernetes.Interface, cellInformer vickinformers.CellInformer, networkPolicyInformer networkv1informers.NetworkPolicyInformer) *controller.Controller {
	h := &cellHandler{
		kubeClient:          kubeClient,
		cellLister:          cellInformer.Lister(),
		networkPilicyLister: networkPolicyInformer.Lister(),
	}
	c := controller.New(h, "Cell")

	glog.Info("Setting up event handlers")
	cellInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			glog.Infof("Old %+v\nnew %+v", old, new)
			c.Enqueue(new)
		},
		DeleteFunc: c.Enqueue,
	})
	return c
}

func (h *cellHandler) Handle(key string) error {
	glog.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		glog.Errorf("invalid resource key: %s", key)
		return nil
	}
	cell, err := h.cellLister.Cells(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("cell '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	glog.Infof("Found cell %+v", cell)

	networkPoliy, err := h.networkPilicyLister.NetworkPolicies(cell.Namespace).Get(resources.NetworkPolicyName(cell))
	if errors.IsNotFound(err) {
		networkPoliy, err = h.kubeClient.NetworkingV1().NetworkPolicies(cell.Namespace).Create(resources.CreateNetworkPolicy(cell))
	}

	if err != nil {
		return err
	}

	glog.Infof("NetworkPolicy created %+v", networkPoliy)
	return nil
}
