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

package informers

import (
	"fmt"
	"reflect"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	appsv1 "k8s.io/client-go/informers/apps/v1"
	autoscalingv2beta1 "k8s.io/client-go/informers/autoscaling/v2beta1"
	batchv1 "k8s.io/client-go/informers/batch/v1"
	corev1 "k8s.io/client-go/informers/core/v1"
	extensionsv1beta1 "k8s.io/client-go/informers/extensions/v1beta1"
	networkingv1 "k8s.io/client-go/informers/networking/v1"

	"github.com/cellery-io/mesh-controller/pkg/clients"
	meshinformers "github.com/cellery-io/mesh-controller/pkg/generated/informers/externalversions"
	istioauthenticationv1alpha1 "github.com/cellery-io/mesh-controller/pkg/generated/informers/externalversions/authentication/v1alpha1"
	meshv1alpha2 "github.com/cellery-io/mesh-controller/pkg/generated/informers/externalversions/mesh/v1alpha2"
	istionetworkv1alpha3 "github.com/cellery-io/mesh-controller/pkg/generated/informers/externalversions/networking/v1alpha3"
	knativeservingv1alpha1 "github.com/cellery-io/mesh-controller/pkg/generated/informers/externalversions/serving/v1alpha1"
)

type Interface interface {
	// K8s informers
	ConfigMaps() corev1.ConfigMapInformer
	Deployments() appsv1.DeploymentInformer
	HorizontalPodAutoscalers() autoscalingv2beta1.HorizontalPodAutoscalerInformer
	Jobs() batchv1.JobInformer
	NetworkPolicies() networkingv1.NetworkPolicyInformer
	PersistentVolumeClaims() corev1.PersistentVolumeClaimInformer
	Secrets() corev1.SecretInformer
	Services() corev1.ServiceInformer
	StatefulSets() appsv1.StatefulSetInformer
	Ingresses() extensionsv1beta1.IngressInformer

	// Istio informers
	IstioDestinationRules() istionetworkv1alpha3.DestinationRuleInformer
	IstioEnvoyFilters() istionetworkv1alpha3.EnvoyFilterInformer
	IstioGateways() istionetworkv1alpha3.GatewayInformer
	IstioVirtualServices() istionetworkv1alpha3.VirtualServiceInformer
	IstioPolicy() istioauthenticationv1alpha1.PolicyInformer

	// Knative serving informers
	KnativeServingConfigurations() knativeservingv1alpha1.ConfigurationInformer

	// Cellery mesh informers
	Cells() meshv1alpha2.CellInformer
	Components() meshv1alpha2.ComponentInformer
	Composites() meshv1alpha2.CompositeInformer
	Gateways() meshv1alpha2.GatewayInformer
	TokenServices() meshv1alpha2.TokenServiceInformer
}

type informers struct {
	kubeInformerFactory kubeinformers.SharedInformerFactory
	meshInformerFactory meshinformers.SharedInformerFactory
}

func New(clients clients.Interface, resync time.Duration) *informers {
	return &informers{
		kubeInformerFactory: kubeinformers.NewSharedInformerFactory(clients.Kubernetes(), resync),
		meshInformerFactory: meshinformers.NewSharedInformerFactory(clients.Mesh(), resync),
	}
}

func (i *informers) Start(stopCh <-chan struct{}) error {
	i.kubeInformerFactory.Start(stopCh)
	i.meshInformerFactory.Start(stopCh)
	return func(results ...map[reflect.Type]bool) error {
		for i, _ := range results {
			for t, ok := range results[i] {
				if !ok {
					return fmt.Errorf("failed to wait for cache with type %s", t)
				}
			}
		}
		return nil
	}(i.kubeInformerFactory.WaitForCacheSync(stopCh), i.meshInformerFactory.WaitForCacheSync(stopCh))
}

func (i *informers) ConfigMaps() corev1.ConfigMapInformer {
	return i.kubeInformerFactory.Core().V1().ConfigMaps()
}

func (i *informers) Deployments() appsv1.DeploymentInformer {
	return i.kubeInformerFactory.Apps().V1().Deployments()
}

func (i *informers) HorizontalPodAutoscalers() autoscalingv2beta1.HorizontalPodAutoscalerInformer {
	return i.kubeInformerFactory.Autoscaling().V2beta1().HorizontalPodAutoscalers()
}

func (i *informers) Jobs() batchv1.JobInformer {
	return i.kubeInformerFactory.Batch().V1().Jobs()
}

func (i *informers) NetworkPolicies() networkingv1.NetworkPolicyInformer {
	return i.kubeInformerFactory.Networking().V1().NetworkPolicies()
}

func (i *informers) PersistentVolumeClaims() corev1.PersistentVolumeClaimInformer {
	return i.kubeInformerFactory.Core().V1().PersistentVolumeClaims()
}

func (i *informers) Secrets() corev1.SecretInformer {
	return i.kubeInformerFactory.Core().V1().Secrets()
}

func (i *informers) Services() corev1.ServiceInformer {
	return i.kubeInformerFactory.Core().V1().Services()
}

func (i *informers) StatefulSets() appsv1.StatefulSetInformer {
	return i.kubeInformerFactory.Apps().V1().StatefulSets()
}

func (i *informers) Ingresses() extensionsv1beta1.IngressInformer {
	return i.kubeInformerFactory.Extensions().V1beta1().Ingresses()
}

func (i *informers) IstioDestinationRules() istionetworkv1alpha3.DestinationRuleInformer {
	return i.meshInformerFactory.Networking().V1alpha3().DestinationRules()
}

func (i *informers) IstioEnvoyFilters() istionetworkv1alpha3.EnvoyFilterInformer {
	return i.meshInformerFactory.Networking().V1alpha3().EnvoyFilters()
}

func (i *informers) IstioGateways() istionetworkv1alpha3.GatewayInformer {
	return i.meshInformerFactory.Networking().V1alpha3().Gateways()
}

func (i *informers) IstioVirtualServices() istionetworkv1alpha3.VirtualServiceInformer {
	return i.meshInformerFactory.Networking().V1alpha3().VirtualServices()
}

func (i *informers) IstioPolicy() istioauthenticationv1alpha1.PolicyInformer {
	return i.meshInformerFactory.Authentication().V1alpha1().Policies()
}

func (i *informers) KnativeServingConfigurations() knativeservingv1alpha1.ConfigurationInformer {
	return i.meshInformerFactory.Serving().V1alpha1().Configurations()
}

func (i *informers) Cells() meshv1alpha2.CellInformer {
	return i.meshInformerFactory.Mesh().V1alpha2().Cells()
}

func (i *informers) Components() meshv1alpha2.ComponentInformer {
	return i.meshInformerFactory.Mesh().V1alpha2().Components()
}

func (i *informers) Composites() meshv1alpha2.CompositeInformer {
	return i.meshInformerFactory.Mesh().V1alpha2().Composites()
}

func (i *informers) Gateways() meshv1alpha2.GatewayInformer {
	return i.meshInformerFactory.Mesh().V1alpha2().Gateways()
}

func (i *informers) TokenServices() meshv1alpha2.TokenServiceInformer {
	return i.meshInformerFactory.Mesh().V1alpha2().TokenServices()
}
