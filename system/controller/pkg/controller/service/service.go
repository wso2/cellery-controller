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

package service

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	vickclientset "github.com/wso2/product-vick/system/controller/pkg/client/clientset/versioned"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	"github.com/wso2/product-vick/system/controller/pkg/controller/service/resources"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	vickinformers "github.com/wso2/product-vick/system/controller/pkg/client/informers/externalversions/vick/v1alpha1"
	listers "github.com/wso2/product-vick/system/controller/pkg/client/listers/vick/v1alpha1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

type serviceHandler struct {
	kubeClient       kubernetes.Interface
	vickClient       vickclientset.Interface
	deploymentLister appsv1listers.DeploymentLister
	k8sServiceLister corev1listers.ServiceLister
	serviceLister    listers.ServiceLister
}

func NewController(
	kubeClient kubernetes.Interface,
	vickClient vickclientset.Interface,
	deploymentInformer appsv1informers.DeploymentInformer,
	k8sServiceInformer corev1informers.ServiceInformer,
	serviceInformer vickinformers.ServiceInformer,
) *controller.Controller {
	h := &serviceHandler{
		kubeClient:       kubeClient,
		vickClient:       vickClient,
		deploymentLister: deploymentInformer.Lister(),
		k8sServiceLister: k8sServiceInformer.Lister(),
		serviceLister:    serviceInformer.Lister(),
	}
	c := controller.New(h, "Service")

	glog.Info("Setting up event handlers")
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			glog.Infof("Old %+v\nnew %+v", old, new)
			c.Enqueue(new)
		},
		DeleteFunc: c.Enqueue,
	})
	return c
}

func (h *serviceHandler) Handle(key string) error {
	glog.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		glog.Errorf("invalid resource key: %s", key)
		return nil
	}
	serviceOriginal, err := h.serviceLister.Services(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("service '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	glog.Infof("Found service %+v", serviceOriginal)
	service := serviceOriginal.DeepCopy()

	if err = h.handle(service); err != nil {
		return err
	}
	return nil

	if _, err = h.updateStatus(service); err != nil {
		return err
	}

	return nil
}

func (h *serviceHandler) handle(service *v1alpha1.Service) error {
	// Get the gateway deployment
	deployment, err := h.deploymentLister.Deployments(service.Namespace).Get(resources.ServiceDeploymentName(service))
	if errors.IsNotFound(err) {
		deployment, err = h.kubeClient.AppsV1().Deployments(service.Namespace).Create(resources.CreateServiceDeployment(service))
		if err != nil {
			glog.Errorf("Failed to create Service Deployment %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("Service deployment created %+v", deployment)

	k8sService, err := h.k8sServiceLister.Services(service.Namespace).Get(resources.ServiceK8sServiceName(service))
	if errors.IsNotFound(err) {
		k8sService, err = h.kubeClient.CoreV1().Services(service.Namespace).Create(resources.CreateServiceK8sService(service))
		if err != nil {
			glog.Errorf("Failed to create Service service %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("Service service created %+v", k8sService)
	service.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	service.Status.OwnerCell = ""
	return nil
}

func (h *serviceHandler) updateStatus(service *v1alpha1.Service) (*v1alpha1.Service, error) {
	return h.vickClient.VickV1alpha1().Services(service.Namespace).Update(service)
}
