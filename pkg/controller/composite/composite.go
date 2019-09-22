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

package composite

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cellery-io/mesh-controller/pkg/meta"

	"github.com/cellery-io/mesh-controller/pkg/clients"
	"github.com/cellery-io/mesh-controller/pkg/informers"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/config"
	"github.com/cellery-io/mesh-controller/pkg/controller"

	//controller_commons "github.com/cellery-io/mesh-controller/pkg/controller/commons"
	"github.com/cellery-io/mesh-controller/pkg/controller/composite/resources"
	routing "github.com/cellery-io/mesh-controller/pkg/controller/routing"
	meshclientset "github.com/cellery-io/mesh-controller/pkg/generated/clientset/versioned"
	v1alpha2listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/mesh/v1alpha2"
	istionetwork1alpha3listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/networking/v1alpha3"
)

type reconciler struct {
	kubeClient                kubernetes.Interface
	meshClient                meshclientset.Interface
	compositeLister           v1alpha2listers.CompositeLister
	componentLister           v1alpha2listers.ComponentLister
	serviceLister             corev1listers.ServiceLister
	secretLister              corev1listers.SecretLister
	tokenServiceLister        v1alpha2listers.TokenServiceLister
	istioVirtualServiceLister istionetwork1alpha3listers.VirtualServiceLister
	cellLister                v1alpha2listers.CellLister
	cfg                       config.Interface
	logger                    *zap.SugaredLogger
	recorder                  record.EventRecorder
}

func NewController(
	clientset clients.Interface,
	informerset informers.Interface,
	cfg config.Interface,
	logger *zap.SugaredLogger,
) *controller.Controller {
	r := &reconciler{
		kubeClient:                clientset.Kubernetes(),
		meshClient:                clientset.Mesh(),
		compositeLister:           informerset.Composites().Lister(),
		componentLister:           informerset.Components().Lister(),
		serviceLister:             informerset.Services().Lister(),
		tokenServiceLister:        informerset.TokenServices().Lister(),
		secretLister:              informerset.Secrets().Lister(),
		istioVirtualServiceLister: informerset.IstioVirtualServices().Lister(),
		cellLister:                informerset.Cells().Lister(),
		cfg:                       cfg,
		logger:                    logger.Named("composite-controller"),
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(r.logger.Named("events").Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: r.kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "composite-controller"})
	r.recorder = recorder
	c := controller.New(r, r.logger, "Composite")

	r.logger.Info("Setting up event handlers")
	informerset.Composites().Informer().AddEventHandler(informers.HandleAll(c.Enqueue))

	informerset.Components().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Composite")),
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
	original, err := r.compositeLister.Composites(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			r.logger.Errorf("composite '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	composite := original.DeepCopy()

	if err = r.reconcile(composite); err != nil {
		r.recorder.Eventf(composite, corev1.EventTypeWarning, "InternalError", "Failed to update cluster: %v", err)
		return err
	}

	if equality.Semantic.DeepEqual(original.Status, composite.Status) {
		return nil
	}

	if _, err = r.updateStatus(composite); err != nil {
		r.recorder.Eventf(composite, corev1.EventTypeWarning, "UpdateFailed", "Failed to update status: %v", err)
		return err
	}
	r.recorder.Eventf(composite, corev1.EventTypeNormal, "Updated", "Updated Composite status %q", composite.GetName())
	return nil
}

func (r *reconciler) reconcile(composite *v1alpha2.Composite) error {
	composite.SetDefaults()
	rErrs := &controller.ReconcileErrors{}

	rErrs.Add(r.reconcileSecret(composite))
	rErrs.Add(r.reconcileTokenService(composite))

	for i, _ := range composite.Spec.Components {
		rErrs.Add(r.reconcileComponent(composite, &composite.Spec.Components[i]))
	}

	rErrs.Add(r.reconcileVirtualService(composite))

	rErrs.Add(r.reconcileRoutingK8sService(composite))

	if !rErrs.Empty() {
		return rErrs
	}

	activeCount := 0
	for _, v := range composite.Status.ComponentStatuses {
		if v == v1alpha2.ComponentCurrentStatusReady || v == v1alpha2.ComponentCurrentStatusIdle {
			activeCount++
		}
	}

	composite.Status.ActiveComponentCount = activeCount
	composite.Status.ComponentCount = len(composite.Spec.Components)
	if composite.Status.ActiveComponentCount == composite.Status.ComponentCount {
		composite.Status.Status = v1alpha2.CompositeCurrentStatusReady
		c := []v1alpha2.CompositeCondition{
			{
				Type:   v1alpha2.CompositeReady,
				Status: corev1.ConditionTrue,
			},
		}
		composite.Status.Conditions = c
	} else {
		composite.Status.Status = v1alpha2.CompositeCurrentStatusNotReady
		c := []v1alpha2.CompositeCondition{
			{
				Type:   v1alpha2.CompositeReady,
				Status: corev1.ConditionFalse,
			},
		}
		composite.Status.Conditions = c
	}
	composite.Status.ObservedGeneration = composite.Generation

	return nil
}

func (r *reconciler) reconcileSecret(composite *v1alpha2.Composite) error {
	secretName := resources.SecretName(composite)
	secret, err := r.secretLister.Secrets(composite.Namespace).Get(resources.SecretName(composite))

	if errors.IsNotFound(err) {
		secret, err = func(composite *v1alpha2.Composite) (*corev1.Secret, error) {
			desiredSecret, err := resources.MakeSecret(composite, r.cfg)
			if err != nil {
				return nil, err
			}
			return r.kubeClient.CoreV1().Secrets(composite.Namespace).Create(desiredSecret)
		}(composite)
		if err != nil {
			r.logger.Errorf("Failed to create Secret %q: %v", secretName, err)
			r.recorder.Eventf(composite, corev1.EventTypeWarning, "CreationFailed", "Failed to create Secret %q: %v", secretName, err)
			return err
		}
		r.recorder.Eventf(composite, corev1.EventTypeNormal, "Created", "Created Secret %q", secretName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Secret %q: %v", secretName, err)
		return err
	}
	resources.StatusFromSecret(composite, secret)
	return nil
}

func (r *reconciler) reconcileTokenService(composite *v1alpha2.Composite) error {
	tokenServiceName := resources.TokenServiceName(composite)
	tokenService, err := r.tokenServiceLister.TokenServices(composite.Namespace).Get(tokenServiceName)
	if errors.IsNotFound(err) {
		tokenService, err = r.meshClient.MeshV1alpha2().TokenServices(composite.Namespace).Create(resources.MakeTokenService(composite))
		if err != nil {
			r.logger.Errorf("Failed to create TokenService %q: %v", tokenServiceName, err)
			r.recorder.Eventf(composite, corev1.EventTypeWarning, "CreationFailed", "Failed to create TokenService %q: %v", tokenServiceName, err)
			return err
		}
		r.recorder.Eventf(composite, corev1.EventTypeNormal, "Created", "Created TokenService %q", tokenServiceName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve TokenService %q: %v", tokenServiceName, err)
		return err
	}
	resources.StatusFromTokenService(composite, tokenService)
	return nil
}

func (r *reconciler) reconcileComponent(composite *v1alpha2.Composite, componentTemplate *v1alpha2.Component) error {
	componentName := resources.ComponentName(composite, componentTemplate)
	component, err := r.componentLister.Components(composite.Namespace).Get(componentName)
	if errors.IsNotFound(err) {
		component, err = r.meshClient.MeshV1alpha2().Components(composite.Namespace).Create(resources.MakeComponent(composite, componentTemplate))
		if err != nil {
			r.logger.Errorf("Failed to create Component %q: %v", componentName, err)
			r.recorder.Eventf(composite, corev1.EventTypeWarning, "CreationFailed", "Failed to create Component %q: %v", componentName, err)
			return err
		}
		r.recorder.Eventf(composite, corev1.EventTypeNormal, "Created", "Created Component %q", componentName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Component %q: %v", componentName, err)
		return err
	} else if !metav1.IsControlledBy(component, composite) {
		return fmt.Errorf("composite: %q does not own the Component: %q", composite.Name, componentName)
	} else {
		component, err = func(composite *v1alpha2.Composite, component *v1alpha2.Component) (*v1alpha2.Component, error) {
			if !resources.RequireComponentUpdate(composite, component) {
				return component, nil
			}
			desiredComponent := resources.MakeComponent(composite, componentTemplate)
			existingComponent := component.DeepCopy()
			resources.CopyComponent(desiredComponent, existingComponent)
			return r.meshClient.MeshV1alpha2().Components(composite.Namespace).Update(existingComponent)
		}(composite, component)
		if err != nil {
			r.logger.Errorf("Failed to update Component %q: %v", componentName, err)
			return err
		}
	}
	resources.StatusFromComponent(composite, component)
	return nil
}

func (r *reconciler) reconcileVirtualService(composite *v1alpha2.Composite) error {
	name := routing.RoutingVirtualServiceName(composite.Name)
	routingVs, err := r.istioVirtualServiceLister.VirtualServices(composite.Namespace).Get(name)
	if errors.IsNotFound(err) {
		routingVs, err = resources.MakeRoutingVirtualService(composite, r.compositeLister, r.cellLister)
		if err != nil {
			r.logger.Errorf("Failed to create Composite VS object %v for instance %s", err, composite.Name)
			r.recorder.Eventf(composite, corev1.EventTypeWarning, "CreationFailed", "Failed to create Virtual Service %q: %v", name, err)
			return err
		}
		if routingVs == nil {
			r.logger.Debugf("No VirtualService created for composite instance %s", composite.Name)
			return nil
		}
		lastAppliedConfig, err := json.Marshal(routing.BuildVirtualServiceLastAppliedConfig(routingVs))
		if err != nil {
			r.logger.Errorf("Failed to create routing VS %v for instance %s", err, composite.Name)
			r.recorder.Eventf(composite, corev1.EventTypeWarning, "CreationFailed", "Failed to create Virtual Service %q: %v", name, err)
			return err
		}
		routing.Annotate(routingVs, corev1.LastAppliedConfigAnnotation, string(lastAppliedConfig))
		routingVs, err = r.meshClient.NetworkingV1alpha3().VirtualServices(composite.Namespace).Create(routingVs)
		if err != nil {
			r.logger.Errorf("Failed to create routing VS %v for instance %s", err, composite.Name)
			r.recorder.Eventf(composite, corev1.EventTypeWarning, "CreationFailed", "Failed to create Virtual Service %q: %v", name, err)
			return err
		}
		r.logger.Debugw("Routing VirtualService created", name, routingVs)
		r.recorder.Eventf(composite, corev1.EventTypeNormal, "Created", "Created Virtual Service %q", name)
	} else if err != nil {
		return err
	} else if !metav1.IsControlledBy(routingVs, composite) {
		return fmt.Errorf("Composite: %q does not own the VS: %q", composite.Name, routingVs)
	} else {
		// TODO: find a better solution
		//routingVs, err = func(composite *v1alpha2.Composite, routingVs *v1alpha3.VirtualService) (*v1alpha3.VirtualService, error) {
		//	if !resources.RequireRoutingVsUpdate(composite, routingVs) {
		//		return routingVs, nil
		//	}
		//	desiredVs, err := resources.MakeRoutingVirtualService(composite, r.compositeLister, r.cellLister)
		//	if err != nil {
		//		r.logger.Errorf("Failed to obtain desired VS %q: %v", desiredVs, err)
		//		return nil, err
		//	}
		//	existingVs := routingVs.DeepCopy()
		//	resources.CopyRoutingVs(desiredVs, existingVs)
		//	return r.meshClient.NetworkingV1alpha3().VirtualServices(composite.Namespace).Update(existingVs)
		//}(composite, routingVs)
		//if err != nil {
		//	r.logger.Errorf("Failed to update VS %q: %v", name, err)
		//	return err
		//}
	}

	resources.StatusFromRoutingVs(composite, routingVs)
	return nil
}

func isEqual(oldService *v1alpha2.Component, newService *v1alpha2.Component) bool {
	// we only consider equality of the spec
	return reflect.DeepEqual(oldService.Spec, newService.Spec)
}

func (r *reconciler) updateStatus(desired *v1alpha2.Composite) (*v1alpha2.Composite, error) {
	composite, err := r.compositeLister.Composites(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(composite.Status, desired.Status) {
		latest := composite.DeepCopy()
		latest.Status = desired.Status
		return r.meshClient.MeshV1alpha2().Composites(desired.Namespace).UpdateStatus(latest)
	}
	return desired, nil
}

type origComponentData struct {
	ComponentName  string `json:"componentName"`
	ContainerPorts []int  `json:"containerPorts"`
}

func (r *reconciler) reconcileRoutingK8sService(composite *v1alpha2.Composite) error {
	// This is a workaround for an issue with switching traffic 100% to a new instance, and terminating the old one.
	// When the old composite instance is terminated, the associated k8s service will be deleted as well. Since the
	// Istio Virtual Service uses that particular gateway k8s service name as a hostname, once its deleted the DNS
	// lookup will fail, which will lead to traffic routing to fail.
	// To overcome this issue, whenever traffic is switched to 100% to a new instance, the previous instance's components'
	// k8s service names are written to an annotation of the new composite instance. This annotation will be picked up
	// by this method and that particular service will be re-created if it does not exist.

	originalComponentSvcKey := composite.Annotations[meta.CompositeOriginalComponentSvcKey]
	if originalComponentSvcKey != "" {
		var origCompData []origComponentData
		err := json.Unmarshal([]byte(originalComponentSvcKey), &origCompData)
		if err != nil {
			return err
		}
		for _, data := range origCompData {
			k8sService, err := r.serviceLister.Services(composite.Namespace).Get(resources.K8sServiceName(data.ComponentName))
			if errors.IsNotFound(err) {
				k8sService, err = r.kubeClient.CoreV1().Services(composite.Namespace).Create(
					resources.MakeOriginalComponentK8sService(composite, data.ComponentName, data.ContainerPorts))
				if err != nil {
					r.logger.Errorf("Failed to create K8s service for component %v", err)
					return err
				}
				r.logger.Debugw("K8s service for component created", data.ComponentName, k8sService)
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}
