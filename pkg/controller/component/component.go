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
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	autoscalingv2beta2lister "k8s.io/client-go/listers/autoscaling/v2beta2"
	batchv1listers "k8s.io/client-go/listers/batch/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	istionetworkingv1alpha3 "github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/clients"
	"github.com/cellery-io/mesh-controller/pkg/config"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/component/resources"
	meshclient "github.com/cellery-io/mesh-controller/pkg/generated/clientset/versioned"
	v1alpha2listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/mesh/v1alpha2"
	istionetwork1alpha3listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/networking/v1alpha3"
	kservingv1alpha1listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/serving/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/informers"
	"github.com/cellery-io/mesh-controller/pkg/meta"
)

type reconciler struct {
	kubeClient                 kubeclient.Interface
	meshClient                 meshclient.Interface
	componentLister            v1alpha2listers.ComponentLister
	serviceLister              corev1listers.ServiceLister
	deploymentLister           appsv1listers.DeploymentLister
	statefulSetLister          appsv1listers.StatefulSetLister
	jobLister                  batchv1listers.JobLister
	hpaLister                  autoscalingv2beta2lister.HorizontalPodAutoscalerLister
	istioVirtualServiceLister  istionetwork1alpha3listers.VirtualServiceLister
	servingConfigurationLister kservingv1alpha1listers.ConfigurationLister
	logger                     *zap.SugaredLogger
	recorder                   record.EventRecorder
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
		componentLister:            informerset.Components().Lister(),
		serviceLister:              informerset.Services().Lister(),
		deploymentLister:           informerset.Deployments().Lister(),
		statefulSetLister:          informerset.StatefulSets().Lister(),
		jobLister:                  informerset.Jobs().Lister(),
		hpaLister:                  informerset.HorizontalPodAutoscalers().Lister(),
		istioVirtualServiceLister:  informerset.IstioVirtualServices().Lister(),
		servingConfigurationLister: informerset.KnativeServingConfigurations().Lister(),
		logger:                     logger.Named("component-controller"),
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
	component.SetDefaults()
	rErrs := &controller.ReconcileErrors{}

	rErrs.Add(r.reconcileService(component))
	rErrs.Add(r.reconcileDeployment(component))
	rErrs.Add(r.reconcileStatefulSet(component))
	rErrs.Add(r.reconcileJob(component))
	rErrs.Add(r.reconcileHpa(component))
	rErrs.Add(r.reconcileServingConfiguration(component))
	rErrs.Add(r.reconcileServingVirtualService(component))

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
			resources.CopyDeployment(desiredDeployment, existingDeployment)
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
			resources.CopyStatefulSet(desiredStatefulSet, existingStatefulSet)
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
			err = r.kubeClient.AutoscalingV2beta2().HorizontalPodAutoscalers(component.Namespace).Delete(hpaName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete HPA %q: %v", hpaName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		hpa, err = r.kubeClient.AutoscalingV2beta2().HorizontalPodAutoscalers(component.Namespace).Create(resources.MakeHpa(component))
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
		hpa, err = func(component *v1alpha2.Component, hpa *autoscalingv2beta2.HorizontalPodAutoscaler) (*autoscalingv2beta2.HorizontalPodAutoscaler, error) {
			if !resources.RequireHpaUpdate(component, hpa) {
				return hpa, nil
			}
			desiredHpa := resources.MakeHpa(component)
			existingHpa := hpa.DeepCopy()
			resources.CopyHpa(desiredHpa, existingHpa)
			return r.kubeClient.AutoscalingV2beta2().HorizontalPodAutoscalers(component.Namespace).Update(existingHpa)
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
	resources.StatusFromServingConfiguration(component, configuration)

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
