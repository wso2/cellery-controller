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
	"reflect"

	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/api/errors"

	istionetworkingv1alpha3 "cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/clients"
	"cellery.io/cellery-controller/pkg/config"
	"cellery.io/cellery-controller/pkg/controller"
	"cellery.io/cellery-controller/pkg/controller/sts/resources"
	meshclientset "cellery.io/cellery-controller/pkg/generated/clientset/versioned"
	"cellery.io/cellery-controller/pkg/informers"

	//appsv1informers "k8s.io/client-go/informers/apps/v1"
	//corev1informers "k8s.io/client-go/informers/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	//corev1informers "k8s.io/client-go/informers/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/record"

	mesh1alpha2listers "cellery.io/cellery-controller/pkg/generated/listers/mesh/v1alpha2"
	istionetwork1alpha3listers "cellery.io/cellery-controller/pkg/generated/listers/networking/v1alpha3"
)

type reconciler struct {
	kubeClient             kubernetes.Interface
	meshClient             meshclientset.Interface
	deploymentLister       appsv1listers.DeploymentLister
	serviceLister          corev1listers.ServiceLister
	istioEnvoyFilterLister istionetwork1alpha3listers.EnvoyFilterLister
	configMapLister        corev1listers.ConfigMapLister
	tokenServiceLister     mesh1alpha2listers.TokenServiceLister
	cfg                    config.Interface
	logger                 *zap.SugaredLogger
	recorder               record.EventRecorder
}

func NewController(
	clientset clients.Interface,
	informerset informers.Interface,
	cfg config.Interface,
	logger *zap.SugaredLogger,
) *controller.Controller {

	r := &reconciler{
		kubeClient:             clientset.Kubernetes(),
		meshClient:             clientset.Mesh(),
		deploymentLister:       informerset.Deployments().Lister(),
		serviceLister:          informerset.Services().Lister(),
		istioEnvoyFilterLister: informerset.IstioEnvoyFilters().Lister(),
		configMapLister:        informerset.ConfigMaps().Lister(),
		tokenServiceLister:     informerset.TokenServices().Lister(),
		cfg:                    cfg,
		logger:                 logger.Named("tokenservice-controller"),
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(r.logger.Named("events").Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: r.kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "tokenservice-controller"})
	r.recorder = recorder
	c := controller.New(r, r.logger, "TokenService")

	r.logger.Info("Setting up event handlers")
	informerset.TokenServices().Informer().AddEventHandler(informers.HandleAll(c.Enqueue))

	informerset.Services().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("TokenService")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Deployments().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("TokenService")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.ConfigMaps().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("TokenService")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.IstioEnvoyFilters().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("TokenService")),
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
	original, err := r.tokenServiceLister.TokenServices(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			r.logger.Errorf("tokenService '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	tokenService := original.DeepCopy()

	if err = r.reconcile(tokenService); err != nil {
		r.recorder.Eventf(tokenService, corev1.EventTypeWarning, "InternalError", "Failed to update cluster: %v", err)
		return err
	}

	if equality.Semantic.DeepEqual(original.Status, tokenService.Status) {
		return nil
	}

	if _, err = r.updateStatus(tokenService); err != nil {
		r.recorder.Eventf(tokenService, corev1.EventTypeWarning, "UpdateFailed", "Failed to update status: %v", err)
		return err
	}
	r.recorder.Eventf(tokenService, corev1.EventTypeNormal, "Updated", "Updated TokenService status %q", tokenService.GetName())
	return nil
}

func (r *reconciler) reconcile(tokenService *v1alpha2.TokenService) error {
	tokenService.SetDefaults()
	rErrs := &controller.ReconcileErrors{}
	rErrs.Add(r.reconcileService(tokenService))
	rErrs.Add(r.reconcileConfigMap(tokenService))
	rErrs.Add(r.reconcileOpaConfigMap(tokenService))
	rErrs.Add(r.reconcileDeployment(tokenService))
	rErrs.Add(r.reconcileEnvoyFilter(tokenService))

	if !rErrs.Empty() {
		return rErrs
	}

	tokenService.Status.ObservedGeneration = tokenService.Generation
	return nil
}

func (r *reconciler) reconcileService(tokenService *v1alpha2.TokenService) error {
	serviceName := resources.ServiceName(tokenService)
	service, err := r.serviceLister.Services(tokenService.Namespace).Get(serviceName)
	if errors.IsNotFound(err) {
		service, err = r.kubeClient.CoreV1().Services(tokenService.Namespace).Create(resources.MakeService(tokenService))
		if err != nil {
			r.logger.Errorf("Failed to create Service %q: %v", serviceName, err)
			r.recorder.Eventf(tokenService, corev1.EventTypeWarning, "CreationFailed", "Failed to create Service %q: %v", serviceName, err)
			return err
		}
		r.recorder.Eventf(tokenService, corev1.EventTypeNormal, "Created", "Created Service %q", serviceName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Service %q: %v", serviceName, err)
		return err
	} else if !metav1.IsControlledBy(service, tokenService) {
		return fmt.Errorf("tokenService: %q does not own the Service: %q", tokenService.Name, serviceName)
	} else {
		service, err = func(tokenService *v1alpha2.TokenService, service *corev1.Service) (*corev1.Service, error) {
			if !resources.RequireServiceUpdate(tokenService, service) {
				return service, nil
			}
			desiredService := resources.MakeService(tokenService)
			existingService := service.DeepCopy()
			resources.CopyService(desiredService, existingService)
			return r.kubeClient.CoreV1().Services(tokenService.Namespace).Update(existingService)
		}(tokenService, service)
		if err != nil {
			r.logger.Errorf("Failed to update Service %q: %v", serviceName, err)
			return err
		}
	}
	resources.StatusFromService(tokenService, service)
	return nil
}

func (r *reconciler) reconcileConfigMap(tokenService *v1alpha2.TokenService) error {
	configMapName := resources.ConfigMapName(tokenService)
	configMap, err := r.configMapLister.ConfigMaps(tokenService.Namespace).Get(configMapName)
	if errors.IsNotFound(err) {
		configMap, err = r.kubeClient.CoreV1().ConfigMaps(tokenService.Namespace).Create(resources.MakeConfigMap(tokenService, r.cfg))
		if err != nil {
			r.logger.Errorf("Failed to create ConfigMap %q: %v", configMapName, err)
			r.recorder.Eventf(tokenService, corev1.EventTypeWarning, "CreationFailed", "Failed to create ConfigMap %q: %v", configMapName, err)
			return err
		}
		r.recorder.Eventf(tokenService, corev1.EventTypeNormal, "Created", "Created ConfigMap %q", configMapName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve ConfigMap %q: %v", configMapName, err)
		return err
	} else if !metav1.IsControlledBy(configMap, tokenService) {
		return fmt.Errorf("tokenService: %q does not own the ConfigMap: %q", tokenService.Name, configMapName)
	} else {
		configMap, err = func(tokenService *v1alpha2.TokenService, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
			if !resources.RequireConfigMapUpdate(tokenService, configMap) {
				return configMap, nil
			}
			desiredConfigMap := resources.MakeConfigMap(tokenService, r.cfg)
			existingConfigMap := configMap.DeepCopy()
			resources.CopyConfigMap(desiredConfigMap, existingConfigMap)
			return r.kubeClient.CoreV1().ConfigMaps(tokenService.Namespace).Update(existingConfigMap)
		}(tokenService, configMap)
		if err != nil {
			r.logger.Errorf("Failed to update ConfigMap %q: %v", configMapName, err)
			return err
		}
	}
	resources.StatusFromConfigMap(tokenService, configMap)
	return nil
}

func (r *reconciler) reconcileOpaConfigMap(tokenService *v1alpha2.TokenService) error {
	configMapName := resources.OpaPolicyConfigMapName(tokenService)
	configMap, err := r.configMapLister.ConfigMaps(tokenService.Namespace).Get(configMapName)
	if errors.IsNotFound(err) {
		configMap, err = r.kubeClient.CoreV1().ConfigMaps(tokenService.Namespace).Create(resources.MakeOpaConfigMap(tokenService, r.cfg))
		if err != nil {
			r.logger.Errorf("Failed to create OPA ConfigMap %q: %v", configMapName, err)
			r.recorder.Eventf(tokenService, corev1.EventTypeWarning, "CreationFailed", "Failed to create OPA ConfigMap %q: %v", configMapName, err)
			return err
		}
		r.recorder.Eventf(tokenService, corev1.EventTypeNormal, "Created", "Created OPA ConfigMap %q", configMapName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve OPA ConfigMap %q: %v", configMapName, err)
		return err
	} else if !metav1.IsControlledBy(configMap, tokenService) {
		return fmt.Errorf("tokenService: %q does not own the OPA ConfigMap: %q", tokenService.Name, configMapName)
	} else {
		configMap, err = func(tokenService *v1alpha2.TokenService, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
			if !resources.RequireOpaConfigMapUpdate(tokenService, configMap) {
				return configMap, nil
			}
			desiredConfigMap := resources.MakeOpaConfigMap(tokenService, r.cfg)
			existingConfigMap := configMap.DeepCopy()
			resources.CopyOpaConfigMap(desiredConfigMap, existingConfigMap)
			return r.kubeClient.CoreV1().ConfigMaps(tokenService.Namespace).Update(existingConfigMap)
		}(tokenService, configMap)
		if err != nil {
			r.logger.Errorf("Failed to update OPA ConfigMap %q: %v", configMapName, err)
			return err
		}
	}
	resources.StatusFromOpaConfigMap(tokenService, configMap)
	return nil
}

func (r *reconciler) reconcileDeployment(tokenService *v1alpha2.TokenService) error {
	deploymentName := resources.DeploymentName(tokenService)
	deployment, err := r.deploymentLister.Deployments(tokenService.Namespace).Get(deploymentName)
	if errors.IsNotFound(err) {
		deployment, err = r.kubeClient.AppsV1().Deployments(tokenService.Namespace).Create(resources.MakeDeployment(tokenService, r.cfg))
		if err != nil {
			r.logger.Errorf("Failed to create Deployment %q: %v", deploymentName, err)
			r.recorder.Eventf(tokenService, corev1.EventTypeWarning, "CreationFailed", "Failed to create Deployment %q: %v", deploymentName, err)
			return err
		}
		r.recorder.Eventf(tokenService, corev1.EventTypeNormal, "Created", "Created Deployment %q", deploymentName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Deployment %q: %v", deploymentName, err)
		return err
	} else if !metav1.IsControlledBy(deployment, tokenService) {
		return fmt.Errorf("tokenService: %q does not own the Deployment: %q", tokenService.Name, deploymentName)
	} else {
		deployment, err = func(tokenService *v1alpha2.TokenService, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
			if !resources.RequireDeploymentUpdate(tokenService, deployment) {
				return deployment, nil
			}
			desiredDeployment := resources.MakeDeployment(tokenService, r.cfg)
			existingDeployment := deployment.DeepCopy()
			resources.CopyDeployment(desiredDeployment, existingDeployment)
			return r.kubeClient.AppsV1().Deployments(tokenService.Namespace).Update(existingDeployment)
		}(tokenService, deployment)
		if err != nil {
			r.logger.Errorf("Failed to update Deployment %q: %v", deploymentName, err)
			return err
		}
	}
	resources.StatusFromDeployment(tokenService, deployment)
	return nil
}

func (r *reconciler) reconcileEnvoyFilter(tokenService *v1alpha2.TokenService) error {
	envoyFilterName := resources.EnvoyFilterName(tokenService)
	envoyFilter, err := r.istioEnvoyFilterLister.EnvoyFilters(tokenService.Namespace).Get(envoyFilterName)
	if !resources.RequireEnvoyFilter(tokenService) {
		if err == nil && metav1.IsControlledBy(envoyFilter, tokenService) {
			err = r.meshClient.NetworkingV1alpha3().EnvoyFilters(tokenService.Namespace).Delete(envoyFilterName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete EnvoyFilter %q: %v", envoyFilterName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		envoyFilter, err = r.meshClient.NetworkingV1alpha3().EnvoyFilters(tokenService.Namespace).Create(resources.MakeEnvoyFilter(tokenService))
		if err != nil {
			r.logger.Errorf("Failed to create EnvoyFilter %q: %v", envoyFilterName, err)
			r.recorder.Eventf(tokenService, corev1.EventTypeWarning, "CreationFailed", "Failed to create EnvoyFilter %q: %v", envoyFilterName, err)
			return err
		}
		r.recorder.Eventf(tokenService, corev1.EventTypeNormal, "Created", "Created EnvoyFilter %q", envoyFilterName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve EnvoyFilter %q: %v", envoyFilterName, err)
		return err
	} else if !metav1.IsControlledBy(envoyFilter, tokenService) {
		return fmt.Errorf("tokenService: %q does not own the EnvoyFilter: %q", tokenService.Name, envoyFilterName)
	} else {
		envoyFilter, err = func(tokenService *v1alpha2.TokenService, envoyFilter *istionetworkingv1alpha3.EnvoyFilter) (*istionetworkingv1alpha3.EnvoyFilter, error) {
			if !resources.RequireEnvoyFilterUpdate(tokenService, envoyFilter) {
				return envoyFilter, nil
			}
			desiredEnvoyFilter := resources.MakeEnvoyFilter(tokenService)
			if err != nil {
				return nil, err
			}
			existingEnvoyFilter := envoyFilter.DeepCopy()
			resources.CopyEnvoyFilter(desiredEnvoyFilter, existingEnvoyFilter)
			return r.meshClient.NetworkingV1alpha3().EnvoyFilters(tokenService.Namespace).Update(existingEnvoyFilter)
		}(tokenService, envoyFilter)
		if err != nil {
			r.logger.Errorf("Failed to update EnvoyFilter %q: %v", envoyFilterName, err)
			return err
		}
	}
	resources.StatusFromEnvoyFilter(tokenService, envoyFilter)
	return nil
}

func (r *reconciler) updateStatus(desired *v1alpha2.TokenService) (*v1alpha2.TokenService, error) {
	gateway, err := r.tokenServiceLister.TokenServices(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(gateway.Status, desired.Status) {
		latest := gateway.DeepCopy()
		latest.Status = desired.Status
		return r.meshClient.MeshV1alpha2().TokenServices(desired.Namespace).UpdateStatus(latest)
	}
	return desired, nil
}
