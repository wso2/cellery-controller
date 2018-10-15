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

package gateway

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	vickclientset "github.com/wso2/product-vick/system/controller/pkg/client/clientset/versioned"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	"github.com/wso2/product-vick/system/controller/pkg/controller/gateway/config"
	"github.com/wso2/product-vick/system/controller/pkg/controller/gateway/resources"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"reflect"

	//appsv1informers "k8s.io/client-go/informers/apps/v1"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	vickinformers "github.com/wso2/product-vick/system/controller/pkg/client/informers/externalversions/vick/v1alpha1"
	listers "github.com/wso2/product-vick/system/controller/pkg/client/listers/vick/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

type gatewayHandler struct {
	kubeClient       kubernetes.Interface
	vickClient       vickclientset.Interface
	deploymentLister appsv1listers.DeploymentLister
	k8sServiceLister corev1listers.ServiceLister
	configMapLister  corev1listers.ConfigMapLister
	gatewayLister    listers.GatewayLister
	gatewayConfig    config.Gateway
}

func NewController(
	kubeClient kubernetes.Interface,
	vickClient vickclientset.Interface,
	systemConfigMapInformer corev1informers.ConfigMapInformer,
	deploymentInformer appsv1informers.DeploymentInformer,
	k8sServiceInformer corev1informers.ServiceInformer,
	configMapInformer corev1informers.ConfigMapInformer,
	gatewayInformer vickinformers.GatewayInformer,
) *controller.Controller {

	h := &gatewayHandler{
		kubeClient:       kubeClient,
		vickClient:       vickClient,
		deploymentLister: deploymentInformer.Lister(),
		k8sServiceLister: k8sServiceInformer.Lister(),
		configMapLister:  configMapInformer.Lister(),
		gatewayLister:    gatewayInformer.Lister(),
	}
	c := controller.New(h, "Gateway")

	glog.Info("Setting up event handlers")
	gatewayInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			glog.Infof("Old %+v\nnew %+v", old, new)
			c.Enqueue(new)
		},
		DeleteFunc: c.Enqueue,
	})

	systemConfigMapInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: h.updateConfig,
		UpdateFunc: func(old, new interface{}) {
			h.updateConfig(new)
		},
	})
	return c
}

func (h *gatewayHandler) Handle(key string) error {
	glog.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		glog.Errorf("invalid resource key: %s", key)
		return nil
	}
	gatewayOriginal, err := h.gatewayLister.Gateways(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("gateway '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	glog.Infof("Found gateway %+v", gatewayOriginal)
	gateway := gatewayOriginal.DeepCopy()

	if err = h.handle(gateway); err != nil {
		return err
	}

	if _, err = h.updateStatus(gateway); err != nil {
		return err
	}

	return nil
}

func (h *gatewayHandler) handle(gateway *v1alpha1.Gateway) error {

	if err := h.handleConfigMap(gateway); err != nil {
		return err
	}

	if err := h.handleDeployment(gateway); err != nil {
		return err
	}

	if err := h.handleK8sService(gateway); err != nil {
		return err
	}

	h.updateOwnerCell(gateway)

	return nil
}

func (h *gatewayHandler) handleConfigMap(gateway *v1alpha1.Gateway) error {
	configMap, err := h.configMapLister.ConfigMaps(gateway.Namespace).Get(resources.GatewayConfigMapName(gateway))
	if errors.IsNotFound(err) {
		gatewayConfigMap, err := resources.CreateGatewayConfigMap(gateway, h.gatewayConfig)
		if err != nil {
			return err
		}
		configMap, err = h.kubeClient.CoreV1().ConfigMaps(gateway.Namespace).Create(gatewayConfigMap)
		if err != nil {
			glog.Errorf("Failed to create Gateway ConfigMap %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("Gateway config map created %+v", configMap)

	return nil
}

func (h *gatewayHandler) handleDeployment(gateway *v1alpha1.Gateway) error {
	deployment, err := h.deploymentLister.Deployments(gateway.Namespace).Get(resources.GatewayDeploymentName(gateway))
	if errors.IsNotFound(err) {
		deployment, err = h.kubeClient.AppsV1().Deployments(gateway.Namespace).Create(resources.CreateGatewayDeployment(gateway, h.gatewayConfig))
		if err != nil {
			glog.Errorf("Failed to create Gateway Deployment %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("Gateway deployment created %+v", deployment)

	if deployment.Status.AvailableReplicas > 0 {
		gateway.Status.Status = "Ready"
	} else {
		gateway.Status.Status = "NotReady"
	}

	return nil
}

func (h *gatewayHandler) handleK8sService(gateway *v1alpha1.Gateway) error {
	k8sService, err := h.k8sServiceLister.Services(gateway.Namespace).Get(resources.GatewayK8sServiceName(gateway))
	if errors.IsNotFound(err) {
		k8sService, err = h.kubeClient.CoreV1().Services(gateway.Namespace).Create(resources.CreateGatewayK8sService(gateway))
		if err != nil {
			glog.Errorf("Failed to create Gateway service %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("Gateway service created %+v", k8sService)

	gateway.Status.HostName = k8sService.Name

	return nil
}

func (h *gatewayHandler) updateStatus(gateway *v1alpha1.Gateway) (*v1alpha1.Gateway, error) {
	latestGateway, err := h.gatewayLister.Gateways(gateway.Namespace).Get(gateway.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(latestGateway.Status, gateway.Status) {
		latestGateway.Status = gateway.Status
		// https://github.com/kubernetes/kubernetes/issues/54672
		// return h.vickClient.VickV1alpha1().Services(service.Namespace).UpdateStatus(latestService)
		return h.vickClient.VickV1alpha1().Gateways(gateway.Namespace).Update(latestGateway)
	}
	return nil, nil
}

func (h *gatewayHandler) updateOwnerCell(gateway *v1alpha1.Gateway) {
	if gateway.Status.OwnerCell == "" {
		ownerRef := metav1.GetControllerOf(gateway)
		if ownerRef != nil {
			gateway.Status.OwnerCell = ownerRef.Name
		} else {
			gateway.Status.OwnerCell = "<none>"
		}
	}
}

func (h *gatewayHandler) updateConfig(obj interface{}) {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return
	}

	if configMap.Name != vick.SystemConfigMapName {
		return
	}

	conf := config.Gateway{}

	if gatewayInitConfig, ok := configMap.Data["cell-gateway-config"]; ok {
		conf.InitConfig = gatewayInitConfig
	} else {
		glog.Errorf("Cell gateway config is missing.")
	}

	if gatewaySetupConfig, ok := configMap.Data["cell-gateway-setup-config"]; ok {
		conf.SetupConfig = gatewaySetupConfig
	} else {
		glog.Errorf("Cell gateway setup config is missing.")
	}

	if gatewayInitImage, ok := configMap.Data["cell-gateway-init-image"]; ok {
		conf.InitImage = gatewayInitImage
	} else {
		glog.Errorf("Cell gateway init image missing.")
	}

	if gatewayImage, ok := configMap.Data["cell-gateway-image"]; ok {
		conf.Image = gatewayImage
	} else {
		glog.Errorf("Cell gateway image missing.")
	}

	h.gatewayConfig = conf
}
