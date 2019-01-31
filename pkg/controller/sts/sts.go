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

package sts

import (
	"fmt"
	"github.com/celleryio/mesh-controller/pkg/apis/mesh"
	"github.com/celleryio/mesh-controller/pkg/apis/mesh/v1alpha1"
	meshclientset "github.com/celleryio/mesh-controller/pkg/client/clientset/versioned"
	"github.com/celleryio/mesh-controller/pkg/controller"
	"github.com/celleryio/mesh-controller/pkg/controller/sts/config"
	"github.com/celleryio/mesh-controller/pkg/controller/sts/resources"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	//appsv1informers "k8s.io/client-go/informers/apps/v1"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	meshinformers "github.com/celleryio/mesh-controller/pkg/client/informers/externalversions/mesh/v1alpha1"
	listers "github.com/celleryio/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

type tokenServiceHandler struct {
	kubeClient         kubernetes.Interface
	meshClient         meshclientset.Interface
	deploymentLister   appsv1listers.DeploymentLister
	k8sServiceLister   corev1listers.ServiceLister
	configMapLister    corev1listers.ConfigMapLister
	tokenServiceLister listers.TokenServiceLister
	tokenServiceConfig config.TokenService
}

func NewController(
	kubeClient kubernetes.Interface,
	meshClient meshclientset.Interface,
	systemConfigMapInformer corev1informers.ConfigMapInformer,
	deploymentInformer appsv1informers.DeploymentInformer,
	k8sServiceInformer corev1informers.ServiceInformer,
	configMapInformer corev1informers.ConfigMapInformer,
	tokenServiceInformer meshinformers.TokenServiceInformer,
) *controller.Controller {

	h := &tokenServiceHandler{
		kubeClient:         kubeClient,
		meshClient:         meshClient,
		deploymentLister:   deploymentInformer.Lister(),
		k8sServiceLister:   k8sServiceInformer.Lister(),
		configMapLister:    configMapInformer.Lister(),
		tokenServiceLister: tokenServiceInformer.Lister(),
	}
	c := controller.New(h, "TokenService")

	glog.Info("Setting up event handlers")
	tokenServiceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
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

func (h *tokenServiceHandler) Handle(key string) error {
	glog.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		glog.Errorf("invalid resource key: %s", key)
		return nil
	}
	tokenServiceOriginal, err := h.tokenServiceLister.TokenServices(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("tokenService '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	glog.Infof("Found tokenService %+v", tokenServiceOriginal)
	tokenService := tokenServiceOriginal.DeepCopy()

	if err = h.handle(tokenService); err != nil {
		return err
	}
	return nil
}

func (h *tokenServiceHandler) handle(tokenService *v1alpha1.TokenService) error {

	configMap, err := h.configMapLister.ConfigMaps(tokenService.Namespace).Get(resources.TokenServiceConfigMapName(tokenService))
	if errors.IsNotFound(err) {
		configMap, err = h.kubeClient.CoreV1().ConfigMaps(tokenService.Namespace).Create(resources.CreateTokenServiceConfigMap(tokenService, h.tokenServiceConfig))
		if err != nil {
			glog.Errorf("Failed to create TokenService ConfigMap %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("TokenService config map created %+v", configMap)

	pocliyConfigMap, err := h.configMapLister.ConfigMaps(tokenService.Namespace).Get(resources.
		TokenServicePolicyConfigMapName(tokenService))
	if errors.IsNotFound(err) {
		pocliyConfigMap, err = h.kubeClient.CoreV1().ConfigMaps(tokenService.Namespace).Create(resources.CreateTokenServiceOPAConfigMap(tokenService, h.tokenServiceConfig))
		if err != nil {
			glog.Errorf("Failed to create TokenService OPA policy config map %v", err)
			return err
		}
	} else if err != nil {
		return err
	}

	glog.Infof("TokenService OPA polciy config map created %+v", pocliyConfigMap)

	deployment, err := h.deploymentLister.Deployments(tokenService.Namespace).Get(resources.TokenServiceDeploymentName(tokenService))
	if errors.IsNotFound(err) {
		deployment, err = h.kubeClient.AppsV1().Deployments(tokenService.Namespace).Create(resources.CreateTokenServiceDeployment(tokenService, h.tokenServiceConfig))
		if err != nil {
			glog.Errorf("Failed to create TokenService Deployment %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("TokenService deployment created %+v", deployment)

	k8sService, err := h.k8sServiceLister.Services(tokenService.Namespace).Get(resources.TokenServiceK8sServiceName(tokenService))
	if errors.IsNotFound(err) {
		k8sService, err = h.kubeClient.CoreV1().Services(tokenService.Namespace).Create(resources.CreateTokenServiceK8sService(tokenService))
		if err != nil {
			glog.Errorf("Failed to create TokenService service %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	glog.Infof("TokenService service created %+v", k8sService)

	return nil
}

func (h *tokenServiceHandler) updateConfig(obj interface{}) {
	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return
	}

	if configMap.Name != mesh.SystemConfigMapName {
		return
	}

	conf := config.TokenService{}

	if tokenServiceConfig, ok := configMap.Data["cell-sts-config"]; ok {
		conf.Config = tokenServiceConfig
	} else {
		glog.Errorf("Cell sts config is missing.")
	}

	if tokenServiceImage, ok := configMap.Data["cell-sts-image"]; ok {
		conf.Image = tokenServiceImage
	} else {
		glog.Errorf("Cell sts image missing.")
	}

	if opaImage, ok := configMap.Data["cell-sts-opa-image"]; ok {
		conf.OpaImage = opaImage
	} else {
		glog.Errorf("Cell sts OPA image missing.")
	}

	if opaPolicy, ok := configMap.Data["opa-default-policy"]; ok {
		conf.Policy = opaPolicy
	} else {
		glog.Errorf("opa default polciy is missing")
	}

	h.tokenServiceConfig = conf
}
