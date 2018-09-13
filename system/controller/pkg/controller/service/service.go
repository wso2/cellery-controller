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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	serviceLister    listers.ServiceLister
	cellLister       listers.CellLister
	deploymentLister appsv1listers.DeploymentLister
	k8sServiceLister corev1listers.ServiceLister
	kubeClient       kubernetes.Interface
	vickClient       vickclientset.Interface
}

func NewController(kubeClient kubernetes.Interface, vickClient vickclientset.Interface, k8sServiceInformer corev1informers.ServiceInformer,
	cellInformer vickinformers.CellInformer, serviceInformer vickinformers.ServiceInformer,
	deploymentInformer appsv1informers.DeploymentInformer) *controller.Controller {
	h := &serviceHandler{
		kubeClient:       kubeClient,
		vickClient:       vickClient,
		serviceLister:    serviceInformer.Lister(),
		cellLister:       cellInformer.Lister(),
		k8sServiceLister: k8sServiceInformer.Lister(),
		deploymentLister: deploymentInformer.Lister(),
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

	ownerCellOriginal, err := h.cellLister.Cells(serviceOriginal.Namespace).Get(serviceOriginal.Spec.Cell)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("cell '%s' in work queue no longer exists", serviceOriginal.Spec.Cell))
			return nil
		}
		return err
	}
	glog.Infof("Found cell %+v", ownerCellOriginal)

	service := serviceOriginal.DeepCopy()

	if !metav1.IsControlledBy(service, ownerCellOriginal) {
		service.OwnerReferences = append(service.OwnerReferences,
			*metav1.NewControllerRef(ownerCellOriginal, schema.GroupVersionKind{
				Group:   v1alpha1.SchemeGroupVersion.Group,
				Version: v1alpha1.SchemeGroupVersion.Version,
				Kind:    "Cell",
			}),
		)
	}

	if _, err := h.vickClient.VickV1alpha1().Services(service.Namespace).Update(service); err != nil {
		glog.Errorf("Error updating service %+v", service)
		return err
	}
	// Get the deployment with the name
	deployment, err := h.deploymentLister.Deployments(service.Namespace).Get(service.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = h.kubeClient.AppsV1().Deployments(service.Namespace).Create(resources.CreateAppDeployment(service))
	}

	k8sService, err := h.k8sServiceLister.Services(service.Namespace).Get(service.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		k8sService, err = h.kubeClient.CoreV1().Services(service.Namespace).Create(resources.CreateCoreService(service))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	glog.Infof("Deployment created %+v", deployment)
	glog.Infof("Service created %+v", k8sService)

	return nil
}
