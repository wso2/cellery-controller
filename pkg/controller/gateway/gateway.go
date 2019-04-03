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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"

	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	extensionsv1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	meshclientset "github.com/cellery-io/mesh-controller/pkg/client/clientset/versioned"
	istioinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/gateway/config"
	"github.com/cellery-io/mesh-controller/pkg/controller/gateway/resources"

	//appsv1informers "k8s.io/client-go/informers/apps/v1"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	meshinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/mesh/v1alpha1"
	listers "github.com/cellery-io/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	istionetworklisters "github.com/cellery-io/mesh-controller/pkg/client/listers/networking/v1alpha3"

	//corev1informers "k8s.io/client-go/informers/core/v1"
	corev1 "k8s.io/api/core/v1"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	extensionsv1beta1informers "k8s.io/client-go/informers/extensions/v1beta1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

type gatewayHandler struct {
	kubeClient           kubernetes.Interface
	meshClient           meshclientset.Interface
	deploymentLister     appsv1listers.DeploymentLister
	k8sServiceLister     corev1listers.ServiceLister
	clusterIngressLister extensionsv1beta1listers.IngressLister
	secretLister         corev1listers.SecretLister
	istioGatewayLister   istionetworklisters.GatewayLister
	istioDRLister        istionetworklisters.DestinationRuleLister
	istioVSLister        istionetworklisters.VirtualServiceLister
	envoyFilterLister    istionetworklisters.EnvoyFilterLister
	configMapLister      corev1listers.ConfigMapLister
	gatewayLister        listers.GatewayLister
	gatewayConfig        config.Gateway
	gatewaySecret        config.Secret
	logger               *zap.SugaredLogger
}

func NewController(
	kubeClient kubernetes.Interface,
	meshClient meshclientset.Interface,
	systemConfigMapInformer corev1informers.ConfigMapInformer,
	systemSecretInformer corev1informers.SecretInformer,
	deploymentInformer appsv1informers.DeploymentInformer,
	k8sServiceInformer corev1informers.ServiceInformer,
	clusterIngressInformer extensionsv1beta1informers.IngressInformer,
	secretInformer corev1informers.SecretInformer,
	istioGatewayInformer istioinformers.GatewayInformer,
	istioDRInformer istioinformers.DestinationRuleInformer,
	istioVSInformer istioinformers.VirtualServiceInformer,
	envoyFilterInformer istioinformers.EnvoyFilterInformer,
	configMapInformer corev1informers.ConfigMapInformer,
	gatewayInformer meshinformers.GatewayInformer,
	logger *zap.SugaredLogger,
) *controller.Controller {

	h := &gatewayHandler{
		kubeClient:           kubeClient,
		meshClient:           meshClient,
		deploymentLister:     deploymentInformer.Lister(),
		k8sServiceLister:     k8sServiceInformer.Lister(),
		clusterIngressLister: clusterIngressInformer.Lister(),
		secretLister:         secretInformer.Lister(),
		istioGatewayLister:   istioGatewayInformer.Lister(),
		istioDRLister:        istioDRInformer.Lister(),
		istioVSLister:        istioVSInformer.Lister(),
		envoyFilterLister:    envoyFilterInformer.Lister(),
		configMapLister:      configMapInformer.Lister(),
		gatewayLister:        gatewayInformer.Lister(),
		logger:               logger.Named("gateway-controller"),
	}
	c := controller.New(h, h.logger, "Gateway")

	h.logger.Info("Setting up event handlers")
	gatewayInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			h.logger.Debugw("Informer update", "old", old, "new", new)
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

	systemSecretInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: h.updateSecret,
		UpdateFunc: func(old, new interface{}) {
			h.updateSecret(new)
		},
	})
	return c
}

func (h *gatewayHandler) Handle(key string) error {
	h.logger.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		h.logger.Errorf("invalid resource key: %s", key)
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
	h.logger.Debugw("lister instance", key, gatewayOriginal)
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

	if gateway.Spec.Type != v1alpha1.GatewayTypeMicroGateway {
		if err := h.handleIstioVirtualService(gateway); err != nil {
			return err
		}

		if err := h.handleIstioGateway(gateway); err != nil {
			return err
		}
	}

	if gateway.Spec.OidcConfig != nil {
		if err := h.handleEnvoyFilter(gateway); err != nil {
			return err
		}
	}

	if len(gateway.Spec.Host) > 0 {
		if len(gateway.Spec.Tls.Key) > 0 && len(gateway.Spec.Tls.Cert) > 0 {
			if err := h.handleClusterIngressSecret(gateway); err != nil {
				return err
			}
		}
		if err := h.handleClusterIngress(gateway); err != nil {
			return err
		}
	}

	//if err := h.handleIstioVirtualServicesForIngress(gateway); err != nil {
	//	return err
	//}

	// if err := h.handleIstioDestinationRules(gateway); err != nil {
	// 	return err
	// }

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
			h.logger.Errorf("Failed to create Gateway ConfigMap %v", err)
			return err
		}
		h.logger.Debugw("Config map created", resources.GatewayConfigMapName(gateway), configMap)
	} else if err != nil {
		return err
	}

	return nil
}

func (h *gatewayHandler) handleDeployment(gateway *v1alpha1.Gateway) error {
	deployment, err := h.deploymentLister.Deployments(gateway.Namespace).Get(resources.GatewayDeploymentName(gateway))
	if errors.IsNotFound(err) {
		desiredDeployment, err := resources.CreateGatewayDeployment(gateway, h.gatewayConfig, h.gatewaySecret)
		if err != nil {
			h.logger.Errorf("Cannot build the Gateway Deployment %v", err)
			return err
		}
		deployment, err = h.kubeClient.AppsV1().Deployments(gateway.Namespace).Create(desiredDeployment)
		if err != nil {
			h.logger.Errorf("Failed to create Gateway Deployment %v", err)
			return err
		}
		h.logger.Debugw("Deployment created", resources.GatewayDeploymentName(gateway), deployment)
	} else if err != nil {
		return err
	}

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
			h.logger.Errorf("Failed to create Gateway service %v", err)
			return err
		}
		h.logger.Debugw("Service created", resources.GatewayK8sServiceName(gateway), k8sService)
	} else if err != nil {
		return err
	}
	gateway.Status.HostName = k8sService.Name
	return nil
}

func (h *gatewayHandler) handleIstioGateway(gateway *v1alpha1.Gateway) error {
	istioGateway, err := h.istioGatewayLister.Gateways(gateway.Namespace).Get(resources.IstioGatewayName(gateway))
	if errors.IsNotFound(err) {
		istioGateway, err = h.meshClient.NetworkingV1alpha3().Gateways(gateway.Namespace).Create(resources.CreateIstioGateway(gateway))
		if err != nil {
			h.logger.Errorf("Failed to create Gateway service %v", err)
			return err
		}
		h.logger.Debugw("Istio gateway created", resources.IstioGatewayName(gateway), istioGateway)
	} else if err != nil {
		return err
	}

	return nil
}

func (h *gatewayHandler) handleIstioVirtualService(gateway *v1alpha1.Gateway) error {
	istioVS, err := h.istioVSLister.VirtualServices(gateway.Namespace).Get(resources.IstioVSName(gateway))
	if errors.IsNotFound(err) {
		istioVS, err = h.meshClient.NetworkingV1alpha3().VirtualServices(gateway.Namespace).Create(resources.CreateIstioVirtualService(gateway))
		if err != nil {
			h.logger.Errorf("Failed to create Gateway service %v", err)
			return err
		}
		h.logger.Debugw("Istio virtual service created", resources.IstioVSName(gateway), istioVS)
	} else if err != nil {
		return err
	}

	return nil
}

func (h *gatewayHandler) handleEnvoyFilter(gateway *v1alpha1.Gateway) error {
	envoyFilter, err := h.envoyFilterLister.EnvoyFilters(gateway.Namespace).Get(resources.EnvoyFilterName(gateway))
	if errors.IsNotFound(err) {
		envoyFilter, err = h.meshClient.NetworkingV1alpha3().EnvoyFilters(gateway.Namespace).Create(resources.CreateEnvoyFilter(gateway))
		if err != nil {
			h.logger.Errorf("Failed to create EnvoyFilter %v", err)
			return err
		}
		h.logger.Debugw("EnvoyFilter created", resources.EnvoyFilterName(gateway), envoyFilter)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *gatewayHandler) handleClusterIngressSecret(gateway *v1alpha1.Gateway) error {
	secret, err := h.secretLister.Secrets(gateway.Namespace).Get(resources.ClusterIngressSecretName(gateway))
	if errors.IsNotFound(err) {
		desiredSecret, err := resources.CreateClusterIngressSecret(gateway, h.gatewaySecret)
		if err != nil {
			h.logger.Errorf("Cannot build the cell ingress Secret %v", err)
			return err
		}
		secret, err = h.kubeClient.CoreV1().Secrets(gateway.Namespace).Create(desiredSecret)
		if err != nil {
			h.logger.Errorf("Failed to create cell ingress Secret %v", err)
			return err
		}
		h.logger.Debugw("Ingress Secret created", resources.ClusterIngressSecretName(gateway), secret)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *gatewayHandler) handleClusterIngress(gateway *v1alpha1.Gateway) error {
	clusterIngress, err := h.clusterIngressLister.Ingresses(gateway.Namespace).Get(resources.ClusterIngressName(gateway))
	if errors.IsNotFound(err) {
		clusterIngress, err = h.kubeClient.ExtensionsV1beta1().Ingresses(gateway.Namespace).Create(resources.CreateClusterIngress(gateway))
		if err != nil {
			h.logger.Errorf("Failed to create Ingress %v", err)
			return err
		}
		h.logger.Debugw("Ingress created", resources.ClusterIngressName(gateway), clusterIngress)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *gatewayHandler) handleIstioVirtualServicesForIngress(gateway *v1alpha1.Gateway) error {
	var hasGlobalApis = false

	for _, apiRoute := range gateway.Spec.HTTPRoutes {
		if apiRoute.Global == true {
			hasGlobalApis = true
			break
		}
	}

	if hasGlobalApis == true {
		istioVS, err := h.istioVSLister.VirtualServices(gateway.Namespace).Get(resources.IstioIngressVirtualServiceName(gateway))
		if errors.IsNotFound(err) {
			istioVS, err = h.meshClient.NetworkingV1alpha3().VirtualServices(gateway.Namespace).Create(
				resources.CreateIstioVirtualServiceForIngress(gateway))
			if err != nil {
				h.logger.Errorf("Failed to create virtual service for ingress %v", err)
				return err
			}
		} else if err != nil {
			return err
		}

		h.logger.Debugf("Istio virtual service for ingress created %+v", istioVS)
	} else {
		h.logger.Debugf("Ingress virtual services not created since gateway %+v does not have global APIs", gateway.Name)
	}

	return nil
}

func (h *gatewayHandler) handleIstioDestinationRules(gateway *v1alpha1.Gateway) error {
	istioDestinationRule, err := h.istioDRLister.DestinationRules(gateway.Namespace).Get(resources.IstioDestinationRuleName(gateway))
	if errors.IsNotFound(err) {
		istioDestinationRule, err = h.meshClient.NetworkingV1alpha3().DestinationRules(gateway.Namespace).Create(
			resources.CreateIstioDestinationRule(gateway))
		if err != nil {
			h.logger.Errorf("Failed to create destination rule %v", err)
			return err
		}
	} else if err != nil {
		return err
	}
	h.logger.Debugf("Istio destination rule created %+v", istioDestinationRule)

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
		// return h.meshClient.MeshV1alpha1().Services(service.Namespace).UpdateStatus(latestService)
		return h.meshClient.MeshV1alpha1().Gateways(gateway.Namespace).Update(latestGateway)
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

	if configMap.Name != mesh.SystemConfigMapName {
		return
	}

	conf := config.Gateway{}

	if gatewayInitConfig, ok := configMap.Data["cell-gateway-config"]; ok {
		conf.InitConfig = gatewayInitConfig
	} else {
		h.logger.Errorf("Cell gateway config is missing.")
	}

	if gatewaySetupConfig, ok := configMap.Data["cell-gateway-setup-config"]; ok {
		conf.SetupConfig = gatewaySetupConfig
	} else {
		h.logger.Errorf("Cell gateway setup config is missing.")
	}

	if gatewayInitImage, ok := configMap.Data["cell-gateway-init-image"]; ok {
		conf.InitImage = gatewayInitImage
	} else {
		h.logger.Errorf("Cell gateway init image missing.")
	}

	if gatewayImage, ok := configMap.Data["cell-gateway-image"]; ok {
		conf.Image = gatewayImage
	} else {
		h.logger.Errorf("Cell gateway image missing.")
	}

	if oidcFilterImage, ok := configMap.Data["oidc-filter-image"]; ok {
		conf.OidcFilterImage = oidcFilterImage
	} else {
		h.logger.Errorf("Cell gateway oidc filter image missing.")
	}

	if skipTlsVerify, ok := configMap.Data["skip-tls-verification"]; ok {
		conf.SkipTlsVerify = skipTlsVerify
	} else {
		conf.SkipTlsVerify = "false"
	}

	h.gatewayConfig = conf
}

func (h *gatewayHandler) updateSecret(obj interface{}) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return
	}

	if secret.Name != mesh.SystemSecretName {
		return
	}

	s := config.Secret{}

	if keyBytes, ok := secret.Data["tls.key"]; ok {
		block, _ := pem.Decode(keyBytes)
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			h.logger.Errorf("error while parsing cellery-secret: %v", err)
			s.PrivateKey = nil
		}
		s.PrivateKey = key
	} else {
		h.logger.Errorf("Missing tls.key in the cellery secret.")
	}

	h.gatewaySecret = s
}
