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

package component

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	autoscalingv2beta1lister "k8s.io/client-go/listers/autoscaling/v2beta1"
	batchv1listers "k8s.io/client-go/listers/batch/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	istioauthenticationv1alpha1 "cellery.io/cellery-controller/pkg/apis/istio/authentication/v1alpha1"
	istionetworkingv1alpha3 "cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/clients"
	"cellery.io/cellery-controller/pkg/config"
	"cellery.io/cellery-controller/pkg/controller"
	"cellery.io/cellery-controller/pkg/controller/component/resources"
	meshclient "cellery.io/cellery-controller/pkg/generated/clientset/versioned"
	istioauthenticationv1alpha1listers "cellery.io/cellery-controller/pkg/generated/listers/authentication/v1alpha1"
	v1alpha2listers "cellery.io/cellery-controller/pkg/generated/listers/mesh/v1alpha2"
	istionetworkv1alpha3listers "cellery.io/cellery-controller/pkg/generated/listers/networking/v1alpha3"
	kservingv1alpha1listers "cellery.io/cellery-controller/pkg/generated/listers/serving/v1alpha1"
	"cellery.io/cellery-controller/pkg/informers"
	"cellery.io/cellery-controller/pkg/meta"
)

type reconciler struct {
	kubeClient                  kubeclient.Interface
	meshClient                  meshclient.Interface
	componentLister             v1alpha2listers.ComponentLister
	serviceLister               corev1listers.ServiceLister
	deploymentLister            appsv1listers.DeploymentLister
	statefulSetLister           appsv1listers.StatefulSetLister
	persistentVolumeClaimLister corev1listers.PersistentVolumeClaimLister
	configMapLister             corev1listers.ConfigMapLister
	secretLister                corev1listers.SecretLister
	jobLister                   batchv1listers.JobLister
	hpaLister                   autoscalingv2beta1lister.HorizontalPodAutoscalerLister
	istioVirtualServiceLister   istionetworkv1alpha3listers.VirtualServiceLister
	istioPolicyLister           istioauthenticationv1alpha1listers.PolicyLister
	servingConfigurationLister  kservingv1alpha1listers.ConfigurationLister
	cfg                         config.Interface
	logger                      *zap.SugaredLogger
	recorder                    record.EventRecorder
}

func NewController(
	clientset clients.Interface,
	informerset informers.Interface,
	cfg config.Interface,
	logger *zap.SugaredLogger,
) *controller.Controller {
	r := &reconciler{
		kubeClient:                  clientset.Kubernetes(),
		meshClient:                  clientset.Mesh(),
		componentLister:             informerset.Components().Lister(),
		serviceLister:               informerset.Services().Lister(),
		deploymentLister:            informerset.Deployments().Lister(),
		statefulSetLister:           informerset.StatefulSets().Lister(),
		persistentVolumeClaimLister: informerset.PersistentVolumeClaims().Lister(),
		configMapLister:             informerset.ConfigMaps().Lister(),
		secretLister:                informerset.Secrets().Lister(),
		jobLister:                   informerset.Jobs().Lister(),
		hpaLister:                   informerset.HorizontalPodAutoscalers().Lister(),
		istioVirtualServiceLister:   informerset.IstioVirtualServices().Lister(),
		istioPolicyLister:           informerset.IstioPolicy().Lister(),
		servingConfigurationLister:  informerset.KnativeServingConfigurations().Lister(),
		cfg:                         cfg,
		logger:                      logger.Named("component-controller"),
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(r.logger.Named("events").Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: r.kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "component-controller"})
	r.recorder = recorder
	c := controller.New(r, r.logger, "Component")

	r.logger.Info("Setting up event handlers")
	informerset.Components().Informer().AddEventHandler(informers.HandleAll(c.Enqueue))

	informerset.Services().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.ConfigMaps().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Secrets().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.PersistentVolumeClaims().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Deployments().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.StatefulSets().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Jobs().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.HorizontalPodAutoscalers().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.IstioVirtualServices().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.IstioPolicy().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.KnativeServingConfigurations().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Component")),
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
	original, err := r.componentLister.Components(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			r.logger.Errorf("component '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	component := original.DeepCopy()

	if err = r.reconcile(component); err != nil {
		r.recorder.Eventf(component, corev1.EventTypeWarning, "InternalError", "Failed to update cluster: %v", err)
		return err
	}

	if equality.Semantic.DeepEqual(original.Status, component.Status) {
		return nil
	}

	if _, err = r.updateStatus(component); err != nil {
		r.recorder.Eventf(component, corev1.EventTypeWarning, "UpdateFailed", "Failed to update status: %v", err)
		return err
	}
	r.recorder.Eventf(component, corev1.EventTypeNormal, "Updated", "Updated Component status %q", component.GetName())
	return nil
}

func (r *reconciler) reconcile(component *v1alpha2.Component) error {
	component.Default()
	rErrs := &controller.ReconcileErrors{}

	rErrs.Add(r.reconcileService(component))
	rErrs.Add(r.reconcileDeployment(component))
	rErrs.Add(r.reconcileStatefulSet(component))
	rErrs.Add(r.reconcileJob(component))
	rErrs.Add(r.reconcileHpa(component))
	rErrs.Add(r.reconcileServingConfiguration(component))
	rErrs.Add(r.reconcileServingVirtualService(component))
	rErrs.Add(r.reconcileTlsPolicy(component))

	for i, _ := range component.Spec.VolumeClaims {
		rErrs.Add(r.reconcilePersistentVolumeClaim(component, &component.Spec.VolumeClaims[i]))
	}

	for i, _ := range component.Spec.Configurations {
		rErrs.Add(r.reconcileConfiguration(component, &component.Spec.Configurations[i]))
	}

	for i, _ := range component.Spec.Secrets {
		rErrs.Add(r.reconcileSecret(component, &component.Spec.Secrets[i]))
	}

	if !rErrs.Empty() {
		return rErrs
	}

	component.Status.ObservedGeneration = component.Generation
	return nil
}

func (r *reconciler) reconcileService(component *v1alpha2.Component) error {
	serviceName := resources.ServiceName(component)
	service, err := r.serviceLister.Services(component.Namespace).Get(serviceName)
	if !resources.RequireService(component) {
		if err == nil && metav1.IsControlledBy(service, component) {
			err = r.kubeClient.CoreV1().Services(component.Namespace).Delete(serviceName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Service %q: %v", serviceName, err)
				return err
			}
		}
		component.Status.ResetServiceName()
		return nil
	}

	if errors.IsNotFound(err) {
		service, err = r.kubeClient.CoreV1().Services(component.Namespace).Create(resources.MakeService(component))
		if err != nil {
			r.logger.Errorf("Failed to create Service %q: %v", serviceName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create Service %q: %v", serviceName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created Service %q", serviceName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Service %q: %v", serviceName, err)
		return err
	} else if !metav1.IsControlledBy(service, component) {
		return fmt.Errorf("component: %q does not own the Service: %q", component.Name, serviceName)
	} else {
		service, err = func(component *v1alpha2.Component, service *corev1.Service) (*corev1.Service, error) {
			if !resources.RequireServiceUpdate(component, service) {
				return service, nil
			}
			desiredService := resources.MakeService(component)
			existingService := service.DeepCopy()
			resources.CopyService(desiredService, existingService)
			return r.kubeClient.CoreV1().Services(component.Namespace).Update(existingService)
		}(component, service)
		if err != nil {
			r.logger.Errorf("Failed to update Service %q: %v", serviceName, err)
			return err
		}
	}
	resources.StatusFromService(component, service)
	return nil
}

func (r *reconciler) reconcileDeployment(component *v1alpha2.Component) error {
	deploymentName := resources.DeploymentName(component)
	deployment, err := r.deploymentLister.Deployments(component.Namespace).Get(deploymentName)
	if !resources.RequireDeployment(component) {
		if err == nil && metav1.IsControlledBy(deployment, component) {
			err = r.kubeClient.AppsV1().Deployments(component.Namespace).Delete(deploymentName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Deployment %q: %v", deploymentName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		deployment, err = r.kubeClient.AppsV1().Deployments(component.Namespace).Create(resources.MakeDeployment(component))
		if err != nil {
			r.logger.Errorf("Failed to create Deployment %q: %v", deploymentName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create Deployment %q: %v", deploymentName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created Deployment %q", deploymentName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Deployment %q: %v", deploymentName, err)
		return err
	} else if !metav1.IsControlledBy(deployment, component) {
		return fmt.Errorf("component: %q does not own the Deployment: %q", component.Name, deploymentName)
	} else {
		deployment, err = func(component *v1alpha2.Component, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
			if !resources.RequireDeploymentUpdate(component, deployment) {
				return deployment, nil
			}
			desiredDeployment := resources.MakeDeployment(component)
			existingDeployment := deployment.DeepCopy()
			resources.CopyDeployment(desiredDeployment, existingDeployment, component)
			return r.kubeClient.AppsV1().Deployments(component.Namespace).Update(existingDeployment)
		}(component, deployment)
		if err != nil {
			r.logger.Errorf("Failed to update Deployment %q: %v", deploymentName, err)
			return err
		}
	}
	resources.StatusFromDeployment(component, deployment)
	return nil
}

func (r *reconciler) reconcileStatefulSet(component *v1alpha2.Component) error {
	statefulSetName := resources.StatefulSetName(component)
	statefulSet, err := r.statefulSetLister.StatefulSets(component.Namespace).Get(statefulSetName)
	if !resources.RequireStatefulSet(component) {
		if err == nil && metav1.IsControlledBy(statefulSet, component) {
			err = r.kubeClient.AppsV1().StatefulSets(component.Namespace).Delete(statefulSetName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete StatefulSet %q: %v", statefulSetName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		statefulSet, err = r.kubeClient.AppsV1().StatefulSets(component.Namespace).Create(resources.MakeStatefulSet(component))
		if err != nil {
			r.logger.Errorf("Failed to create StatefulSet %q: %v", statefulSetName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create StatefulSet %q: %v", statefulSetName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created StatefulSet %q", statefulSetName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve StatefulSet %q: %v", statefulSetName, err)
		return err
	} else if !metav1.IsControlledBy(statefulSet, component) {
		return fmt.Errorf("component: %q does not own the StatefulSet: %q", component.Name, statefulSetName)
	} else {
		statefulSet, err = func(component *v1alpha2.Component, statefulSet *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
			if !resources.RequireStatefulSetUpdate(component, statefulSet) {
				return statefulSet, nil
			}
			desiredStatefulSet := resources.MakeStatefulSet(component)
			existingStatefulSet := statefulSet.DeepCopy()
			resources.CopyStatefulSet(desiredStatefulSet, existingStatefulSet, component)
			return r.kubeClient.AppsV1().StatefulSets(component.Namespace).Update(existingStatefulSet)
		}(component, statefulSet)
		if err != nil {
			r.logger.Errorf("Failed to update StatefulSet %q: %v", statefulSetName, err)
			return err
		}
	}
	resources.StatusFromStatefulSet(component, statefulSet)
	return nil
}

func (r *reconciler) reconcileJob(component *v1alpha2.Component) error {
	jobName := resources.JobName(component)
	job, err := r.jobLister.Jobs(component.Namespace).Get(jobName)
	if !resources.RequireJob(component) {
		if err == nil && metav1.IsControlledBy(job, component) {
			err = r.kubeClient.BatchV1().Jobs(component.Namespace).Delete(jobName, meta.DeleteWithPropagationBackground())
			if err != nil {
				r.logger.Errorf("Failed to delete Job %q: %v", jobName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		job, err = r.kubeClient.BatchV1().Jobs(component.Namespace).Create(resources.MakeJob(component))
		if err != nil {
			r.logger.Errorf("Failed to create Job %q: %v", jobName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create Job %q: %v", jobName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created Job %q", jobName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Job %q: %v", jobName, err)
		return err
	} else if !metav1.IsControlledBy(job, component) {
		return fmt.Errorf("component: %q does not own the Job: %q", component.Name, jobName)
	} else {
		if resources.RequireJobUpdate(component, job) {
			err = r.kubeClient.BatchV1().Jobs(component.Namespace).Delete(jobName, meta.DeleteWithPropagationBackground())
			if err != nil {
				r.logger.Errorf("Failed to delete Job %q: %v", jobName, err)
				return err
			}
		}

	}
	resources.StatusFromJob(component, job)
	return nil
}

func (r *reconciler) reconcileHpa(component *v1alpha2.Component) error {
	hpaName := resources.HpaName(component)
	hpa, err := r.hpaLister.HorizontalPodAutoscalers(component.Namespace).Get(hpaName)

	if !resources.RequireHpa(component) {
		if err == nil && metav1.IsControlledBy(hpa, component) {
			err = r.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(component.Namespace).Delete(hpaName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete HPA %q: %v", hpaName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		hpa, err = r.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(component.Namespace).Create(resources.MakeHpa(component))
		if err != nil {
			r.logger.Errorf("Failed to create HPA %q: %v", hpaName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create HPA %q: %v", hpaName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created HPA %q", hpaName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve HPA %q: %v", hpaName, err)
		return err
	} else if !metav1.IsControlledBy(hpa, component) {
		return fmt.Errorf("component: %q does not own the HPA: %q", component.Name, hpaName)
	} else {
		hpa, err = func(component *v1alpha2.Component, hpa *autoscalingv2beta1.HorizontalPodAutoscaler) (*autoscalingv2beta1.HorizontalPodAutoscaler, error) {
			if !resources.RequireHpaUpdate(component, hpa) {
				return hpa, nil
			}
			desiredHpa := resources.MakeHpa(component)
			existingHpa := hpa.DeepCopy()
			resources.CopyHpa(desiredHpa, existingHpa)
			return r.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(component.Namespace).Update(existingHpa)
		}(component, hpa)
		if err != nil {
			r.logger.Errorf("Failed to update HPA %q: %v", hpaName, err)
			return err
		}
	}
	resources.StatusFromHpa(component, hpa)
	return nil
}

func (r *reconciler) reconcileServingConfiguration(component *v1alpha2.Component) error {
	configurationName := resources.ServingConfigurationName(component)
	configuration, err := r.servingConfigurationLister.Configurations(component.Namespace).Get(configurationName)

	if !resources.RequireKnativeServing(component) {
		if err == nil && metav1.IsControlledBy(configuration, component) {
			err = r.meshClient.ServingV1alpha1().Configurations(component.Namespace).Delete(configurationName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Serving Configuration %q: %v", configurationName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		configuration, err = r.meshClient.ServingV1alpha1().Configurations(component.Namespace).Create(resources.MakeServingConfiguration(component))
		if err != nil {
			r.logger.Errorf("Failed to create Serving Configuration %q: %v", configurationName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create Serving Configuration %q: %v", configurationName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created Serving Configuration %q", configurationName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Serving Configuration %q: %v", configurationName, err)
		return err
	} else if !metav1.IsControlledBy(configuration, component) {
		return fmt.Errorf("component: %q does not own the Serving Configuration: %q", component.Name, configurationName)
	} else {
		if resources.RequireServingConfigurationUpdate(component, configuration) {
			err = r.meshClient.ServingV1alpha1().Configurations(component.Namespace).Delete(configurationName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Serving Configuration %q: %v", configurationName, err)
				return err
			}
		}
	}
	resources.StatusFromServingConfiguration(component, configuration, r.deploymentLister.List)

	component.Status.ServiceName = configuration.Name
	return nil
}

func (r *reconciler) reconcileServingVirtualService(component *v1alpha2.Component) error {
	virtualServiceName := resources.ServingVirtualServiceName(component)
	virtualService, err := r.istioVirtualServiceLister.VirtualServices(component.Namespace).Get(virtualServiceName)
	if !resources.RequireKnativeServing(component) {
		if err == nil && metav1.IsControlledBy(virtualService, component) {
			err = r.meshClient.NetworkingV1alpha3().VirtualServices(component.Namespace).Delete(virtualServiceName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Serving VirtualService %q: %v", virtualServiceName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		virtualService, err = r.meshClient.NetworkingV1alpha3().VirtualServices(component.Namespace).Create(resources.MakeServingVirtualService(component))
		if err != nil {
			r.logger.Errorf("Failed to create Serving VirtualService %q: %v", virtualServiceName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create Serving VirtualService %q: %v", virtualServiceName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created Serving VirtualService %q", virtualServiceName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Serving VirtualService %q: %v", virtualServiceName, err)
		return err
	} else if !metav1.IsControlledBy(virtualService, component) {
		return fmt.Errorf("component: %q does not own the Serving VirtualService: %q", component.Name, virtualServiceName)
	} else {
		virtualService, err = func(component *v1alpha2.Component, virtualService *istionetworkingv1alpha3.VirtualService) (*istionetworkingv1alpha3.VirtualService, error) {
			if !resources.RequireServingVirtualServiceUpdate(component, virtualService) {
				return virtualService, nil
			}
			desiredVirtualService := resources.MakeServingVirtualService(component)
			existingVirtualService := virtualService.DeepCopy()
			resources.CopyServingVirtualService(desiredVirtualService, existingVirtualService)
			return r.meshClient.NetworkingV1alpha3().VirtualServices(component.Namespace).Update(existingVirtualService)
		}(component, virtualService)
		if err != nil {
			r.logger.Errorf("Failed to update Serving VirtualService %q: %v", virtualServiceName, err)
			return err
		}
	}
	resources.StatusFromServingVirtualService(component, virtualService)
	return nil
}

func (r *reconciler) reconcileTlsPolicy(component *v1alpha2.Component) error {
	policyName := resources.TlsPolicyName(component)
	policy, err := r.istioPolicyLister.Policies(component.Namespace).Get(policyName)
	if !resources.RequireTlsPolicy(component) {
		if err == nil && metav1.IsControlledBy(policy, component) {
			err = r.meshClient.AuthenticationV1alpha1().Policies(component.Namespace).Delete(policyName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Tls Policy %q: %v", policyName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		policy, err = r.meshClient.AuthenticationV1alpha1().Policies(component.Namespace).Create(resources.MakeTlsPolicy(component))
		if err != nil {
			r.logger.Errorf("Failed to create Tls Policy %q: %v", policyName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create Tls Policy %q: %v", policyName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created Tls Policy %q", policyName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Tls Policy %q: %v", policyName, err)
		return err
	} else if !metav1.IsControlledBy(policy, component) {
		return fmt.Errorf("component: %q does not own the Tls Policy: %q", component.Name, policyName)
	} else {
		policy, err = func(component *v1alpha2.Component, policy *istioauthenticationv1alpha1.Policy) (*istioauthenticationv1alpha1.Policy, error) {
			if !resources.RequireTlsPolicyUpdate(component, policy) {
				return policy, nil
			}
			desiredPolicy := resources.MakeTlsPolicy(component)
			existingPolicy := policy.DeepCopy()
			resources.CopyTlsPolicy(desiredPolicy, existingPolicy)
			return r.meshClient.AuthenticationV1alpha1().Policies(component.Namespace).Update(existingPolicy)
		}(component, policy)
		if err != nil {
			r.logger.Errorf("Failed to update Tls Policy %q: %v", policyName, err)
			return err
		}
	}
	resources.StatusFromTlsPolicy(component, policy)
	return nil
}

func (r *reconciler) reconcilePersistentVolumeClaim(component *v1alpha2.Component, volumeClaim *v1alpha2.VolumeClaim) error {
	persistentVolumeClaimName := resources.PersistentVolumeClaimName(component, volumeClaim)
	persistentVolumeClaim, err := r.persistentVolumeClaimLister.PersistentVolumeClaims(component.Namespace).Get(persistentVolumeClaimName)
	if errors.IsNotFound(err) {
		persistentVolumeClaim, err = r.kubeClient.CoreV1().PersistentVolumeClaims(component.Namespace).Create(resources.MakePersistentVolumeClaim(component, volumeClaim))
		if err != nil {
			r.logger.Errorf("Failed to create PersistentVolumeClaim %q: %v", persistentVolumeClaimName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create PersistentVolumeClaim %q: %v", persistentVolumeClaimName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created PersistentVolumeClaim %q", persistentVolumeClaimName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve PersistentVolumeClaim %q: %v", persistentVolumeClaimName, err)
		return err
	}
	resources.StatusFromPersistentVolumeClaim(component, persistentVolumeClaim)
	return nil
}

func (r *reconciler) reconcileConfiguration(component *v1alpha2.Component, configMapTemplate *corev1.ConfigMap) error {
	configMapName := resources.ConfigMapName(component, configMapTemplate)
	configMap, err := r.configMapLister.ConfigMaps(component.Namespace).Get(configMapName)
	if errors.IsNotFound(err) {
		configMap, err = r.kubeClient.CoreV1().ConfigMaps(component.Namespace).Create(resources.MakeConfigMap(component, configMapTemplate))
		if err != nil {
			r.logger.Errorf("Failed to create ConfigMap %q: %v", configMapName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create ConfigMap %q: %v", configMapName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created ConfigMap %q", configMapName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve ConfigMap %q: %v", configMapName, err)
		return err
	} else if !metav1.IsControlledBy(configMap, component) {
		return fmt.Errorf("component: %q does not own the ConfigMap: %q", component.Name, configMapName)
	} else {
		configMap, err = func(component *v1alpha2.Component, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
			if !resources.RequireConfigMapUpdate(component, configMap) {
				return configMap, nil
			}
			desiredConfigMap := resources.MakeConfigMap(component, configMapTemplate)
			existingConfigMap := configMap.DeepCopy()
			resources.CopyConfigMap(desiredConfigMap, existingConfigMap)
			return r.kubeClient.CoreV1().ConfigMaps(component.Namespace).Update(existingConfigMap)
		}(component, configMap)
		if err != nil {
			r.logger.Errorf("Failed to update ConfigMap %q: %v", configMapName, err)
			return err
		}
	}
	resources.StatusFromConfigMap(component, configMap)
	return nil
}

func (r *reconciler) reconcileSecret(component *v1alpha2.Component, secretTemplate *corev1.Secret) error {
	secretName := resources.SecretName(component, secretTemplate)
	secret, err := r.secretLister.Secrets(component.Namespace).Get(secretName)
	if errors.IsNotFound(err) {
		secret, err = func(component *v1alpha2.Component, secretTemplate *corev1.Secret) (*corev1.Secret, error) {
			desiredSecret, err := resources.MakeSecret(component, secretTemplate, r.cfg)
			if err != nil {
				return nil, err
			}
			return r.kubeClient.CoreV1().Secrets(component.Namespace).Create(desiredSecret)
		}(component, secretTemplate)
		if err != nil {
			r.logger.Errorf("Failed to create Secret %q: %v", secretName, err)
			r.recorder.Eventf(component, corev1.EventTypeWarning, "CreationFailed", "Failed to create Secret %q: %v", secretName, err)
			return err
		}
		r.recorder.Eventf(component, corev1.EventTypeNormal, "Created", "Created ConfigMap %q", secretName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Secret %q: %v", secretName, err)
		return err
	} else if !metav1.IsControlledBy(secret, component) {
		return fmt.Errorf("component: %q does not own the Secret: %q", component.Name, secretName)
	} else {
		secret, err = func(component *v1alpha2.Component, secret *corev1.Secret) (*corev1.Secret, error) {
			if !resources.RequireSecretUpdate(component, secret) {
				return secret, nil
			}
			desiredSecret, err := resources.MakeSecret(component, secretTemplate, r.cfg)
			if err != nil {
				return nil, err
			}
			existingSecret := secret.DeepCopy()
			resources.CopySecret(desiredSecret, existingSecret)
			return r.kubeClient.CoreV1().Secrets(component.Namespace).Update(existingSecret)
		}(component, secret)
		if err != nil {
			r.logger.Errorf("Failed to update Secret %q: %v", secretName, err)
			return err
		}
	}
	resources.StatusFromSecret(component, secret)
	return nil
}

func (r *reconciler) updateStatus(desired *v1alpha2.Component) (*v1alpha2.Component, error) {
	component, err := r.componentLister.Components(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(component.Status, desired.Status) {
		latest := component.DeepCopy()
		latest.Status = desired.Status
		return r.meshClient.MeshV1alpha2().Components(desired.Namespace).UpdateStatus(latest)
	}
	return desired, nil
}
