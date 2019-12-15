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
	"reflect"
	"strconv"

	"cellery.io/cellery-controller/pkg/meta"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	batchv1listers "k8s.io/client-go/listers/batch/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	extensionsv1beta1listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingv2beta1lister "k8s.io/client-go/listers/autoscaling/v2beta1"

	istionetworkingv1alpha3 "cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/clients"
	"cellery.io/cellery-controller/pkg/config"
	"cellery.io/cellery-controller/pkg/controller"
	"cellery.io/cellery-controller/pkg/controller/gateway/resources"
	meshclientset "cellery.io/cellery-controller/pkg/generated/clientset/versioned"
	mesh1alpha2listers "cellery.io/cellery-controller/pkg/generated/listers/mesh/v1alpha2"
	istionetwork1alpha3listers "cellery.io/cellery-controller/pkg/generated/listers/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/informers"
)

type reconciler struct {
	kubeClient                 kubernetes.Interface
	meshClient                 meshclientset.Interface
	deploymentLister           appsv1listers.DeploymentLister
	serviceLister              corev1listers.ServiceLister
	jobLister                  batchv1listers.JobLister
	clusterIngressLister       extensionsv1beta1listers.IngressLister
	secretLister               corev1listers.SecretLister
	istioGatewayLister         istionetwork1alpha3listers.GatewayLister
	istioDestinationRuleLister istionetwork1alpha3listers.DestinationRuleLister
	istioVirtualServiceLister  istionetwork1alpha3listers.VirtualServiceLister
	istioEnvoyFilterLister     istionetwork1alpha3listers.EnvoyFilterLister
	configMapLister            corev1listers.ConfigMapLister
	gatewayLister              mesh1alpha2listers.GatewayLister
	hpaLister                  autoscalingv2beta1lister.HorizontalPodAutoscalerLister

	cfg      config.Interface
	logger   *zap.SugaredLogger
	recorder record.EventRecorder
}

func NewController(
	clientset clients.Interface,
	informerset informers.Interface,
	cfg config.Interface,
	logger *zap.SugaredLogger,
) *controller.Controller {

	r := &reconciler{
		kubeClient:                 clientset.Kubernetes(),
		meshClient:                 clientset.Mesh(),
		deploymentLister:           informerset.Deployments().Lister(),
		serviceLister:              informerset.Services().Lister(),
		jobLister:                  informerset.Jobs().Lister(),
		clusterIngressLister:       informerset.Ingresses().Lister(),
		secretLister:               informerset.Secrets().Lister(),
		istioGatewayLister:         informerset.IstioGateways().Lister(),
		istioDestinationRuleLister: informerset.IstioDestinationRules().Lister(),
		istioVirtualServiceLister:  informerset.IstioVirtualServices().Lister(),
		istioEnvoyFilterLister:     informerset.IstioEnvoyFilters().Lister(),
		configMapLister:            informerset.ConfigMaps().Lister(),
		gatewayLister:              informerset.Gateways().Lister(),
		hpaLister:                  informerset.HorizontalPodAutoscalers().Lister(),
		cfg:                        cfg,
		logger:                     logger.Named("gateway-controller"),
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(r.logger.Named("events").Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: r.kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "gateway-controller"})
	r.recorder = recorder
	c := controller.New(r, r.logger, "Gateway")

	r.logger.Info("Setting up event handlers")
	informerset.Gateways().Informer().AddEventHandler(informers.HandleAll(c.Enqueue))

	informerset.Services().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Deployments().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.IstioVirtualServices().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.IstioGateways().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.IstioEnvoyFilters().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Ingresses().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Secrets().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.ConfigMaps().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Gateway")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	return c
}

func (r *reconciler) Reconcile(key string) error {
	r.logger.Infof("Reconcile called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		r.logger.Errorf("invalid resource key: %s", key)
		return nil
	}
	original, err := r.gatewayLister.Gateways(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			r.logger.Errorf("gateway '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	gateway := original.DeepCopy()

	if err = r.reconcile(gateway); err != nil {
		r.recorder.Eventf(gateway, corev1.EventTypeWarning, "InternalError", "Failed to update cluster: %v", err)
		return err
	}

	if equality.Semantic.DeepEqual(original.Status, gateway.Status) {
		return nil
	}

	if _, err = r.updateStatus(gateway); err != nil {
		r.recorder.Eventf(gateway, corev1.EventTypeWarning, "UpdateFailed", "Failed to update status: %v", err)
		return err
	}
	r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Updated", "Updated Gateway status %q", gateway.GetName())
	return nil
}

func (r *reconciler) reconcile(gateway *v1alpha2.Gateway) error {
	gateway.Default()
	rErrs := &controller.ReconcileErrors{}

	rErrs.Add(r.reconcileService(gateway))
	rErrs.Add(r.reconcileDeployment(gateway))
	rErrs.Add(r.reconcileIstioVirtualService(gateway))
	rErrs.Add(r.reconcileIstioGateway(gateway))
	rErrs.Add(r.reconcileHpa(gateway))

	// Extensions
	rErrs.Add(r.reconcileApiPublisherConfigMap(gateway))
	rErrs.Add(r.reconcileApiPublisherJob(gateway))

	rErrs.Add(r.reconcileClusterIngressSecret(gateway))
	rErrs.Add(r.reconcileClusterIngress(gateway))

	rErrs.Add(r.reconcileOidcEnvoyFilter(gateway))

	rErrs.Add(r.reconcileRoutingK8sService(gateway))
	// if gateway.Spec.Empty() {
	// 	gateway.Status.Status = "Ready"
	// 	gateway.Status.HostName = "N/A"
	// } else {

	// 	if err := r.reconcileConfigMap(gateway); err != nil {
	// 		return err
	// 	}

	// 	if err := r.reconcileDeployment(gateway); err != nil {
	// 		return err
	// 	}

	// 	if err := r.reconcileAutoscalePolicy(gateway); err != nil {
	// 		return err
	// 	}

	// 	if err := r.reconcileK8sService(gateway); err != nil {
	// 		return err
	// 	}

	// 	if gateway.Spec.Type != v1alpha1.GatewayTypeMicroGateway {
	// 		if err := r.reconcileIstioVirtualService(gateway); err != nil {
	// 			return err
	// 		}

	// 		if err := r.reconcileIstioGateway(gateway); err != nil {
	// 			return err
	// 		}
	// 	}

	// 	if gateway.Spec.OidcConfig != nil {
	// 		if err := r.reconcileEnvoyFilter(gateway); err != nil {
	// 			return err
	// 		}
	// 	}

	// 	if len(gateway.Spec.Host) > 0 {
	// 		if len(gateway.Spec.Tls.Key) > 0 && len(gateway.Spec.Tls.Cert) > 0 {
	// 			if err := r.reconcileClusterIngressSecret(gateway); err != nil {
	// 				return err
	// 			}
	// 		}
	// 		if err := r.reconcileClusterIngress(gateway); err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	//if err := r.reconcileIstioVirtualServicesForIngress(gateway); err != nil {
	//	return err
	//}

	// if err := r.reconcileIstioDestinationRules(gateway); err != nil {
	// 	return err
	// }
	if !rErrs.Empty() {
		return rErrs
	}

	gateway.Status.ObservedGeneration = gateway.Generation
	return nil
}

func (r *reconciler) reconcileService(gateway *v1alpha2.Gateway) error {
	serviceName := resources.ServiceName(gateway)
	service, err := r.serviceLister.Services(gateway.Namespace).Get(serviceName)
	if !resources.RequireService(gateway) {
		if err == nil && metav1.IsControlledBy(service, gateway) {
			err = r.kubeClient.CoreV1().Services(gateway.Namespace).Delete(serviceName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Service %q: %v", serviceName, err)
				return err
			}
		}
		gateway.Status.ResetServiceName()
		return nil
	}

	if errors.IsNotFound(err) {
		service, err = r.kubeClient.CoreV1().Services(gateway.Namespace).Create(resources.MakeService(gateway))
		if err != nil {
			r.logger.Errorf("Failed to create Service %q: %v", serviceName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create Service %q: %v", serviceName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created Service %q", serviceName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Service %q: %v", serviceName, err)
		return err
	} else if !metav1.IsControlledBy(service, gateway) {
		return fmt.Errorf("gateway: %q does not own the Service: %q", gateway.Name, serviceName)
	} else {
		service, err = func(gateway *v1alpha2.Gateway, service *corev1.Service) (*corev1.Service, error) {
			if !resources.RequireServiceUpdate(gateway, service) {
				return service, nil
			}
			desiredService := resources.MakeService(gateway)
			existingService := service.DeepCopy()
			resources.CopyService(desiredService, existingService)
			return r.kubeClient.CoreV1().Services(gateway.Namespace).Update(existingService)
		}(gateway, service)
		if err != nil {
			r.logger.Errorf("Failed to update Service %q: %v", serviceName, err)
			return err
		}
	}
	resources.StatusFromService(gateway, service)
	return nil
}

func (r *reconciler) reconcileDeployment(gateway *v1alpha2.Gateway) error {
	deploymentName := resources.DeploymentName(gateway)
	deployment, err := r.deploymentLister.Deployments(gateway.Namespace).Get(deploymentName)
	if !resources.RequireDeployment(gateway) {
		if err == nil && metav1.IsControlledBy(deployment, gateway) {
			err = r.kubeClient.AppsV1().Deployments(gateway.Namespace).Delete(deploymentName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Deployment %q: %v", deploymentName, err)
				return err
			}
		}
		gateway.Status.Status = v1alpha2.GatewayCurrentStatusReady
		return nil
	}

	if errors.IsNotFound(err) {
		deployment, err = func(gateway *v1alpha2.Gateway) (*appsv1.Deployment, error) {
			desiredDeployment, err := resources.MakeDeployment(gateway, r.cfg)
			if err != nil {
				return nil, err
			}
			return r.kubeClient.AppsV1().Deployments(gateway.Namespace).Create(desiredDeployment)
		}(gateway)
		if err != nil {
			r.logger.Errorf("Failed to create Deployment %q: %v", deploymentName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create Deployment %q: %v", deploymentName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created Deployment %q", deploymentName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Deployment %q: %v", deploymentName, err)
		return err
	} else if !metav1.IsControlledBy(deployment, gateway) {
		return fmt.Errorf("gateway: %q does not own the Deployment: %q", gateway.Name, deploymentName)
	} else {
		deployment, err = func(gateway *v1alpha2.Gateway, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
			if !resources.RequireDeploymentUpdate(gateway, deployment) {
				return deployment, nil
			}
			desiredDeployment, err := resources.MakeDeployment(gateway, r.cfg)
			if err != nil {
				return nil, err
			}
			existingDeployment := deployment.DeepCopy()
			resources.CopyDeployment(desiredDeployment, existingDeployment)
			return r.kubeClient.AppsV1().Deployments(gateway.Namespace).Update(existingDeployment)
		}(gateway, deployment)
		if err != nil {
			r.logger.Errorf("Failed to update Deployment %q: %v", deploymentName, err)
			return err
		}
	}
	resources.StatusFromDeployment(gateway, deployment)
	return nil
}

func (r *reconciler) reconcileIstioGateway(gateway *v1alpha2.Gateway) error {
	istioGatewayName := resources.IstioGatewayName(gateway)
	istioGateway, err := r.istioGatewayLister.Gateways(gateway.Namespace).Get(istioGatewayName)
	if !resources.RequireIstioGateway(gateway) {
		if err == nil && metav1.IsControlledBy(istioGateway, gateway) {
			err = r.meshClient.NetworkingV1alpha3().Gateways(gateway.Namespace).Delete(istioGatewayName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Istio Gateway %q: %v", istioGatewayName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		istioGateway, err = r.meshClient.NetworkingV1alpha3().Gateways(gateway.Namespace).Create(resources.MakeIstioGateway(gateway))
		if err != nil {
			r.logger.Errorf("Failed to create Istio Gateway %q: %v", istioGatewayName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create Istio Gateway %q: %v", istioGatewayName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created Istio Gateway %q", istioGatewayName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Istio Gateway %q: %v", istioGatewayName, err)
		return err
	} else if !metav1.IsControlledBy(istioGateway, gateway) {
		return fmt.Errorf("gateway: %q does not own the Istio Gateway: %q", gateway.Name, istioGatewayName)
	} else {
		istioGateway, err = func(gateway *v1alpha2.Gateway, istioGateway *istionetworkingv1alpha3.Gateway) (*istionetworkingv1alpha3.Gateway, error) {
			if !resources.RequireIstioGatewayUpdate(gateway, istioGateway) {
				return istioGateway, nil
			}
			desiredIstioGateway := resources.MakeIstioGateway(gateway)
			existingIstioGateway := istioGateway.DeepCopy()
			resources.CopyIstioGateway(desiredIstioGateway, existingIstioGateway)
			return r.meshClient.NetworkingV1alpha3().Gateways(gateway.Namespace).Update(existingIstioGateway)
		}(gateway, istioGateway)
		if err != nil {
			r.logger.Errorf("Failed to update Istio Gateway %q: %v", istioGatewayName, err)
			return err
		}
	}
	resources.StatusFromIstioGateway(gateway, istioGateway)
	return nil
}

func (r *reconciler) reconcileIstioVirtualService(gateway *v1alpha2.Gateway) error {
	virtualServiceName := resources.IstioVirtualServiceName(gateway)
	virtualService, err := r.istioVirtualServiceLister.VirtualServices(gateway.Namespace).Get(virtualServiceName)
	if !resources.RequireVirtualService(gateway) {
		if err == nil && metav1.IsControlledBy(virtualService, gateway) {
			err = r.meshClient.NetworkingV1alpha3().VirtualServices(gateway.Namespace).Delete(virtualServiceName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete VirtualService %q: %v", virtualServiceName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		virtualService, err = r.meshClient.NetworkingV1alpha3().VirtualServices(gateway.Namespace).Create(resources.MakeVirtualService(gateway))
		if err != nil {
			r.logger.Errorf("Failed to create VirtualService %q: %v", virtualServiceName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create VirtualService %q: %v", virtualServiceName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created VirtualService %q", virtualServiceName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve VirtualService %q: %v", virtualServiceName, err)
		return err
	} else if !metav1.IsControlledBy(virtualService, gateway) {
		return fmt.Errorf("gateway: %q does not own the VirtualService: %q", gateway.Name, virtualServiceName)
	} else {
		virtualService, err = func(gateway *v1alpha2.Gateway, virtualService *istionetworkingv1alpha3.VirtualService) (*istionetworkingv1alpha3.VirtualService, error) {
			if !resources.RequireVirtualServiceUpdate(gateway, virtualService) {
				return virtualService, nil
			}
			desiredVirtualService := resources.MakeVirtualService(gateway)
			existingVirtualService := virtualService.DeepCopy()
			resources.CopyVirtualService(desiredVirtualService, existingVirtualService)
			return r.meshClient.NetworkingV1alpha3().VirtualServices(gateway.Namespace).Update(existingVirtualService)
		}(gateway, virtualService)
		if err != nil {
			r.logger.Errorf("Failed to update VirtualService %q: %v", virtualServiceName, err)
			return err
		}
	}
	resources.StatusFromVirtualService(gateway, virtualService)
	return nil
}

// func (r *reconciler) reconcileConfigMap(gateway *v1alpha2.Gateway) error {
// 	configMap, err := r.configMapLister.ConfigMaps(gateway.Namespace).Get(resources.ApiPublisherConfigMap(gateway))
// 	if errors.IsNotFound(err) {
// 		gatewayConfigMap, err := resources.CreateGatewayConfigMap(gateway, r.gatewayConfig)
// 		if err != nil {
// 			return err
// 		}
// 		configMap, err = r.kubeClient.CoreV1().ConfigMaps(gateway.Namespace).Create(gatewayConfigMap)
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Gateway ConfigMap %v", err)
// 			return err
// 		}
// 		r.logger.Debugw("Config map created", resources.ApiPublisherConfigMap(gateway), configMap)
// 	} else if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *reconciler) reconcileAutoscalePolicy(gateway *v1alpha2.Gateway) error {
// 	// if autoscaling is disabled in cellery system wide configs, do nothing
// 	if !r.gatewayConfig.EnableAutoscaling {
// 		// if the Autoscaler has been previously created, delete it
// 		_, err := r.autoscalePolicyLister.AutoscalePolicies(gateway.Namespace).Get(resources.GatewayAutoscalePolicyName(gateway))
// 		if errors.IsNotFound(err) {
// 			// has not been created previously
// 			r.logger.Debugf("Autoscaling disabled, hence no autoscaler created for Gateway %s", gateway.Name)
// 			return nil
// 		} else if err != nil {
// 			return err
// 		}
// 		// created previously, delete it
// 		var delPropPtr = v1.DeletePropagationBackground
// 		err = r.meshClient.MeshV1alpha1().AutoscalePolicies(gateway.Namespace).Delete(resources.GatewayAutoscalePolicyName(gateway),
// 			&v1.DeleteOptions{
// 				PropagationPolicy: &delPropPtr,
// 			})
// 		if err != nil {
// 			return err
// 		}
// 		r.logger.Debugf("Terminated Autoscaler for Gateway %s", gateway.Name)

// 		return nil
// 	}

// 	autoscalePolicy, err := r.autoscalePolicyLister.AutoscalePolicies(gateway.Namespace).Get(resources.GatewayAutoscalePolicyName(gateway))
// 	if errors.IsNotFound(err) {
// 		if gateway.Spec.Autoscaling != nil {
// 			autoscalePolicy = resources.CreateAutoscalePolicy(gateway)
// 		} else {
// 			// if autoscaling is not defined, no autoscaler will be added by default
// 			// autoscalePolicy = resources.CreateDefaultAutoscalePolicy(gateway)
// 			return nil
// 		}
// 		lastAppliedConfig, err := json.Marshal(autoscale.BuildAutoscalePolicyLastAppliedConfig(autoscalePolicy))
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Autoscale policy %v", err)
// 			return err
// 		}
// 		autoscale.Annotate(autoscalePolicy, corev1.LastAppliedConfigAnnotation, string(lastAppliedConfig))
// 		autoscalePolicy, err = r.meshClient.MeshV1alpha1().AutoscalePolicies(gateway.Namespace).Create(autoscalePolicy)
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Autoscale policy %v", err)
// 			return err
// 		}
// 		r.logger.Infow("Autoscale policy created", resources.GatewayAutoscalePolicyName(gateway), autoscalePolicy)
// 		autoscalePolicy.Status = "Ready"
// 		r.updateAutoscalePolicyStatus(autoscalePolicy)

// 	} else if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *reconciler) updateAutoscalePolicyStatus(policy *v1alpha2.AutoscalePolicy) error {
// 	_, err := r.meshClient.MeshV1alpha1().AutoscalePolicies(policy.Namespace).Update(policy)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *reconciler) reconcileRoutingK8sService(gateway *v1alpha2.Gateway) error {
	// This is a workaround for an issue with switching traffic 100% to a new instance, and terminating the old one.
	// When the old instance is terminated, the associated gateway k8s service will be deleted as well. Since the
	// Istio Virtual Service uses that particular gateway k8s service name as a hostname, once its deleted the DNS
	// lookup will fail, which will lead to traffic routing to fail.
	// To overcome this issue, whenever traffic is switched to 100% to a new instance, the previous instance's gateway
	// k8s service name is written to an annotation of the new instance gateway. This annotation will be picked up
	// by this method and that particular service will be re-created if it does not exist.

	originalGwK8sSvcName := gateway.Annotations[meta.CellOriginalGatewaySvcKey]
	if originalGwK8sSvcName != "" {
		k8sService, err := r.serviceLister.Services(gateway.Namespace).Get(originalGwK8sSvcName)
		if errors.IsNotFound(err) {
			k8sService, err = r.kubeClient.CoreV1().Services(gateway.Namespace).Create(resources.MakeOriginalGatewayK8sService(gateway, originalGwK8sSvcName))
			if err != nil {
				r.logger.Errorf("Failed to create K8s service for original gateway %v", err)
				return err
			}
			r.logger.Debugw("K8s service for original gateway created", originalGwK8sSvcName, k8sService)
		} else if err != nil {
			return err
		}
		//gateway.Status.HostName = k8sService.Name
	}
	return nil
}

// func (r *reconciler) reconcileEnvoyFilter(gateway *v1alpha2.Gateway) error {
// 	envoyFilter, err := r.envoyFilterLister.EnvoyFilters(gateway.Namespace).Get(resources.EnvoyFilterName(gateway))
// 	if errors.IsNotFound(err) {
// 		envoyFilter, err = r.meshClient.NetworkingV1alpha3().EnvoyFilters(gateway.Namespace).Create(resources.CreateEnvoyFilter(gateway))
// 		if err != nil {
// 			r.logger.Errorf("Failed to create EnvoyFilter %v", err)
// 			return err
// 		}
// 		r.logger.Debugw("EnvoyFilter created", resources.EnvoyFilterName(gateway), envoyFilter)
// 	} else if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *reconciler) reconcileIstioVirtualServicesForIngress(gateway *v1alpha2.Gateway) error {
// 	var hasGlobalApis = false

// 	for _, apiRoute := range gateway.Spec.HTTPRoutes {
// 		if apiRoute.Global == true {
// 			hasGlobalApis = true
// 			break
// 		}
// 	}

// 	if hasGlobalApis == true {
// 		istioVS, err := r.istioVirtualServiceLister.VirtualServices(gateway.Namespace).Get(resources.IstioIngressVirtualServiceName(gateway))
// 		if errors.IsNotFound(err) {
// 			istioVS, err = r.meshClient.NetworkingV1alpha3().VirtualServices(gateway.Namespace).Create(
// 				resources.CreateIstioVirtualServiceForIngress(gateway))
// 			if err != nil {
// 				r.logger.Errorf("Failed to create virtual service for ingress %v", err)
// 				return err
// 			}
// 		} else if err != nil {
// 			return err
// 		}

// 		r.logger.Debugf("Istio virtual service for ingress created %+v", istioVS)
// 	} else {
// 		r.logger.Debugf("Ingress virtual services not created since gateway %+v does not have global APIs", gateway.Name)
// 	}

// 	return nil
// }

// func (r *reconciler) reconcileIstioDestinationRules(gateway *v1alpha2.Gateway) error {
// 	istioDestinationRule, err := r.istioDestinationRuleLister.DestinationRules(gateway.Namespace).Get(resources.IstioDestinationRuleName(gateway))
// 	if errors.IsNotFound(err) {
// 		istioDestinationRule, err = r.meshClient.NetworkingV1alpha3().DestinationRules(gateway.Namespace).Create(
// 			resources.CreateIstioDestinationRule(gateway))
// 		if err != nil {
// 			r.logger.Errorf("Failed to create destination rule %v", err)
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	}
// 	r.logger.Debugf("Istio destination rule created %+v", istioDestinationRule)

// 	return nil
// }

func (r *reconciler) updateStatus(desired *v1alpha2.Gateway) (*v1alpha2.Gateway, error) {
	gateway, err := r.gatewayLister.Gateways(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(gateway.Status, desired.Status) {
		latest := gateway.DeepCopy()
		latest.Status = desired.Status
		return r.meshClient.MeshV1alpha2().Gateways(desired.Namespace).UpdateStatus(latest)
	}
	return desired, nil
}

func (r *reconciler) reconcileHpa(gw *v1alpha2.Gateway) error {
	hpaName := resources.HpaName(gw)
	hpa, err := r.hpaLister.HorizontalPodAutoscalers(gw.Namespace).Get(hpaName)
	if !resources.RequireHpa(gw) {
		if err == nil && metav1.IsControlledBy(hpa, gw) {
			err = r.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(gw.Namespace).Delete(hpaName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete HPA %q: %v", hpaName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		hpa, err = r.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(gw.Namespace).Create(resources.MakeHpa(gw))
		if err != nil {
			r.logger.Errorf("Failed to create HPA %q: %v", hpaName, err)
			r.recorder.Eventf(gw, corev1.EventTypeWarning, "CreationFailed", "Failed to create HPA %q: %v", hpaName, err)
			return err
		}
		r.recorder.Eventf(gw, corev1.EventTypeNormal, "Created", "Created HPA %q", hpaName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve HPA %q: %v", hpaName, err)
		return err
	} else if !metav1.IsControlledBy(hpa, gw) {
		return fmt.Errorf("gw: %q does not own the HPA: %q", gw.Name, hpaName)
	} else {
		hpa, err = func(gw *v1alpha2.Gateway, hpa *autoscalingv2beta1.HorizontalPodAutoscaler) (*autoscalingv2beta1.HorizontalPodAutoscaler, error) {
			if !resources.RequireHpaUpdate(gw, hpa) {
				return hpa, nil
			}
			desiredHpa := resources.MakeHpa(gw)
			existingHpa := hpa.DeepCopy()
			resources.CopyHpa(desiredHpa, existingHpa)
			return r.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(gw.Namespace).Update(existingHpa)
		}(gw, hpa)
		if err != nil {
			r.logger.Errorf("Failed to update HPA %q: %v", hpaName, err)
			return err
		}
	}
	resources.StatusFromHpa(gw, hpa)
	return nil
}

func isAutoscalingEnabled(str string) bool {
	if enabled, err := strconv.ParseBool(str); err == nil {
		return enabled
	}
	return false
}
