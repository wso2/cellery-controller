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

package cell

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	networkingv1listers "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/clients"
	"github.com/cellery-io/mesh-controller/pkg/config"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/cell/resources"
	meshclientset "github.com/cellery-io/mesh-controller/pkg/generated/clientset/versioned"
	v1alpha2listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/mesh/v1alpha2"
	istiov1alpha1listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/informers"
)

type reconciler struct {
	kubeClient                kubernetes.Interface
	meshClient                meshclientset.Interface
	secretLister              corev1listers.SecretLister
	networkPolicyLister       networkingv1listers.NetworkPolicyLister
	istioVirtualServiceLister istiov1alpha1listers.VirtualServiceLister
	istioEnvoyFilterLister    istiov1alpha1listers.EnvoyFilterLister
	cellLister                v1alpha2listers.CellLister
	gatewayLister             v1alpha2listers.GatewayLister
	tokenServiceLister        v1alpha2listers.TokenServiceLister
	componentLister           v1alpha2listers.ComponentLister
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
		cellLister:                informerset.Cells().Lister(),
		componentLister:           informerset.Components().Lister(),
		gatewayLister:             informerset.Gateways().Lister(),
		tokenServiceLister:        informerset.TokenServices().Lister(),
		networkPolicyLister:       informerset.NetworkPolicies().Lister(),
		secretLister:              informerset.Secrets().Lister(),
		istioEnvoyFilterLister:    informerset.IstioEnvoyFilters().Lister(),
		istioVirtualServiceLister: informerset.IstioVirtualServices().Lister(),
		cfg:                       cfg,
		logger:                    logger.Named("cell-controller"),
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(r.logger.Named("events").Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: r.kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "cell-controller"})
	r.recorder = recorder
	c := controller.New(r, r.logger, "Cell")

	r.logger.Info("Setting up event handlers")
	informerset.Cells().Informer().AddEventHandler(informers.HandleAll(c.Enqueue))

	informerset.Components().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Cell")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Gateways().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Cell")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.TokenServices().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Cell")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.NetworkPolicies().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Cell")),
		Handler:    informers.HandleAll(c.EnqueueControllerOf),
	})

	informerset.Secrets().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithOwnerGroupVersionKind(v1alpha2.SchemeGroupVersion.WithKind("Cell")),
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
	original, err := r.cellLister.Cells(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			r.logger.Errorf("cell '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	cell := original.DeepCopy()

	if err = r.reconcile(cell); err != nil {
		r.recorder.Eventf(cell, corev1.EventTypeWarning, "InternalError", "Failed to update cluster: %v", err)
		return err
	}

	if equality.Semantic.DeepEqual(original.Status, cell.Status) {
		return nil
	}

	if _, err = r.updateStatus(cell); err != nil {
		r.recorder.Eventf(cell, corev1.EventTypeWarning, "UpdateFailed", "Failed to update status: %v", err)
		return err
	}
	r.recorder.Eventf(cell, corev1.EventTypeNormal, "Updated", "Updated Cell status %q", cell.GetName())
	return nil
}

func (r *reconciler) reconcile(cell *v1alpha2.Cell) error {
	cell.SetDefaults()
	rErrs := &controller.ReconcileErrors{}

	rErrs.Add(r.reconcileNetworkPolicy(cell))
	rErrs.Add(r.reconcileSecret(cell))
	rErrs.Add(r.reconcileGateway(cell))
	rErrs.Add(r.reconcileTokenService(cell))

	for i, _ := range cell.Spec.Components {
		rErrs.Add(r.reconcileComponent(cell, &cell.Spec.Components[i]))
	}

	// if err := r.reconcileVirtualService(cell); err != nil {
	// 	return err
	// }
	if !rErrs.Empty() {
		return rErrs
	}

	activeCount := 0
	for _, v := range cell.Status.ComponentStatuses {
		if v == v1alpha2.ComponentCurrentStatusReady {
			activeCount++
		}
	}

	cell.Status.ActiveComponentCount = activeCount
	cell.Status.ComponentCount = len(cell.Spec.Components)
	if cell.Status.GatewayStatus == v1alpha2.GatewayCurrentStatusReady &&
		cell.Status.ActiveComponentCount == cell.Status.ComponentCount {
		cell.Status.Status = v1alpha2.CellCurrentStatusReady
		c := []v1alpha2.CellCondition{
			{
				Type:   v1alpha2.CellReady,
				Status: corev1.ConditionTrue,
			},
		}
		cell.Status.Conditions = c
	} else {
		cell.Status.Status = v1alpha2.CellCurrentStatusNotReady
		c := []v1alpha2.CellCondition{
			{
				Type:   v1alpha2.CellReady,
				Status: corev1.ConditionFalse,
			},
		}
		cell.Status.Conditions = c
	}
	cell.Status.ObservedGeneration = cell.Generation

	return nil
}

func (r *reconciler) reconcileNetworkPolicy(cell *v1alpha2.Cell) error {
	networkPolicyName := resources.NetworkPolicyName(cell)
	networkPolicy, err := r.networkPolicyLister.NetworkPolicies(cell.Namespace).Get(networkPolicyName)
	if errors.IsNotFound(err) {
		networkPolicy, err = r.kubeClient.NetworkingV1().NetworkPolicies(cell.Namespace).Create(resources.MakeNetworkPolicy(cell))
		if err != nil {
			r.logger.Errorf("Failed to create NetworkPolicy %q: %v", networkPolicyName, err)
			r.recorder.Eventf(cell, corev1.EventTypeWarning, "CreationFailed", "Failed to create NetworkPolicy %q: %v", networkPolicyName, err)
			return err
		}
		r.recorder.Eventf(cell, corev1.EventTypeNormal, "Created", "Created NetworkPolicy %q", networkPolicyName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve NetworkPolicy %q: %v", networkPolicyName, err)
		return err
	} else if !metav1.IsControlledBy(networkPolicy, cell) {
		return fmt.Errorf("cell: %q does not own the NetworkPolicy: %q", cell.Name, networkPolicyName)
	} else {
		networkPolicy, err = func(cell *v1alpha2.Cell, networkPolicy *networkv1.NetworkPolicy) (*networkv1.NetworkPolicy, error) {
			if !resources.RequireNetworkPolicyUpdate(cell, networkPolicy) {
				return networkPolicy, nil
			}
			desiredNetworkPolicy := resources.MakeNetworkPolicy(cell)
			existingNetworkPolicy := networkPolicy.DeepCopy()
			resources.CopyNetworkPolicy(desiredNetworkPolicy, existingNetworkPolicy)
			return r.kubeClient.NetworkingV1().NetworkPolicies(cell.Namespace).Update(existingNetworkPolicy)
		}(cell, networkPolicy)
		if err != nil {
			r.logger.Errorf("Failed to update NetworkPolicy %q: %v", networkPolicyName, err)
			return err
		}
	}
	resources.StatusFromNetworkPolicy(cell, networkPolicy)
	return nil
}

func (r *reconciler) reconcileSecret(cell *v1alpha2.Cell) error {
	secretName := resources.SecretName(cell)
	secret, err := r.secretLister.Secrets(cell.Namespace).Get(resources.SecretName(cell))

	if errors.IsNotFound(err) {
		secret, err = func(cell *v1alpha2.Cell) (*corev1.Secret, error) {
			desiredSecret, err := resources.MakeSecret(cell, r.cfg)
			if err != nil {
				return nil, err
			}
			return r.kubeClient.CoreV1().Secrets(cell.Namespace).Create(desiredSecret)
		}(cell)
		if err != nil {
			r.logger.Errorf("Failed to create Secret %q: %v", secretName, err)
			r.recorder.Eventf(cell, corev1.EventTypeWarning, "CreationFailed", "Failed to create Secret %q: %v", secretName, err)
			return err
		}
		r.recorder.Eventf(cell, corev1.EventTypeNormal, "Created", "Created Secret %q", secretName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Secret %q: %v", secretName, err)
		return err
	} else if !metav1.IsControlledBy(secret, cell) {
		return fmt.Errorf("cell: %q does not own the Secret: %q", cell.Name, secretName)
	}
	resources.StatusFromSecret(cell, secret)
	return nil
}

func (r *reconciler) reconcileGateway(cell *v1alpha2.Cell) error {
	gatewayName := resources.GatewayName(cell)
	gateway, err := r.gatewayLister.Gateways(cell.Namespace).Get(gatewayName)
	if errors.IsNotFound(err) {
		gateway, err = r.meshClient.MeshV1alpha2().Gateways(cell.Namespace).Create(resources.MakeGateway(cell))
		if err != nil {
			r.logger.Errorf("Failed to create Gateway %q: %v", gatewayName, err)
			r.recorder.Eventf(cell, corev1.EventTypeWarning, "CreationFailed", "Failed to create Gateway %q: %v", gatewayName, err)
			return err
		}
		r.recorder.Eventf(cell, corev1.EventTypeNormal, "Created", "Created Gateway %q", gatewayName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Gateway %q: %v", gatewayName, err)
		return err
	} else if !metav1.IsControlledBy(gateway, cell) {
		return fmt.Errorf("cell: %q does not own the Gateway: %q", cell.Name, gatewayName)
	} else {
		gateway, err = func(cell *v1alpha2.Cell, gateway *v1alpha2.Gateway) (*v1alpha2.Gateway, error) {
			if !resources.RequireGatewayUpdate(cell, gateway) {
				return gateway, nil
			}
			desiredGateway := resources.MakeGateway(cell)
			existingGateway := gateway.DeepCopy()
			resources.CopyGateway(desiredGateway, existingGateway)
			return r.meshClient.MeshV1alpha2().Gateways(cell.Namespace).Update(existingGateway)
		}(cell, gateway)
		if err != nil {
			r.logger.Errorf("Failed to update Gateway %q: %v", gatewayName, err)
			return err
		}
	}
	resources.StatusFromGateway(cell, gateway)
	return nil
}

// func (r *reconciler) reconcileGateway(cell *v1alpha2.Cell) error {
// 	gateway, err := r.gatewayLister.Gateways(cell.Namespace).Get(resources.GatewayName(cell))
// 	if errors.IsNotFound(err) {
// 		gateway = resources.MakeGateway(cell)
// 		lastAppliedConfig, err := json.Marshal(buildLastAppliedConfig(gateway))
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Gateway last applied config %v", err)
// 			return err
// 		}
// 		annotate(gateway, corev1.LastAppliedConfigAnnotation, string(lastAppliedConfig))
// 		gateway, err = r.meshClient.MeshV1alpha2().Gateways(cell.Namespace).Create(gateway)
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Gateway %v", err)
// 			return err
// 		}
// 		r.logger.Debugw("Gateway created", resources.GatewayName(cell), gateway)
// 	} else if err != nil {
// 		return err
// 	}

// 	cell.Status.GatewayHostname = gateway.Status.ServiceName
// 	cell.Status.GatewayStatus = string(gateway.Status.Status)
// 	return nil
// }

// func buildLastAppliedConfig(gw *v1alpha2.Gateway) *v1alpha2.Gateway {
// 	return &v1alpha2.Gateway{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Gateway",
// 			APIVersion: v1alpha2.SchemeGroupVersion.String(),
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      gw.Name,
// 			Namespace: gw.Namespace,
// 		},
// 		Spec: gw.Spec,
// 	}
// }

// func annotate(gw *v1alpha2.Gateway, name string, value string) {
// 	annotations := make(map[string]string, len(gw.ObjectMeta.Annotations)+1)
// 	annotations[name] = value
// 	for k, v := range gw.ObjectMeta.Annotations {
// 		annotations[k] = v
// 	}
// 	gw.Annotations = annotations
// }

func (r *reconciler) reconcileTokenService(cell *v1alpha2.Cell) error {
	tokenServiceName := resources.TokenServiceName(cell)
	tokenService, err := r.tokenServiceLister.TokenServices(cell.Namespace).Get(tokenServiceName)
	if errors.IsNotFound(err) {
		tokenService, err = r.meshClient.MeshV1alpha2().TokenServices(cell.Namespace).Create(resources.MakeTokenService(cell))
		if err != nil {
			r.logger.Errorf("Failed to create TokenService %q: %v", tokenServiceName, err)
			r.recorder.Eventf(cell, corev1.EventTypeWarning, "CreationFailed", "Failed to create TokenService %q: %v", tokenServiceName, err)
			return err
		}
		r.recorder.Eventf(cell, corev1.EventTypeNormal, "Created", "Created TokenService %q", tokenServiceName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve TokenService %q: %v", tokenServiceName, err)
		return err
	} else if !metav1.IsControlledBy(tokenService, cell) {
		return fmt.Errorf("cell: %q does not own the TokenService: %q", cell.Name, tokenServiceName)
	} else {
		tokenService, err = func(cell *v1alpha2.Cell, tokenService *v1alpha2.TokenService) (*v1alpha2.TokenService, error) {
			if !resources.RequireTokenServiceUpdate(cell, tokenService) {
				return tokenService, nil
			}
			desiredTokenService := resources.MakeTokenService(cell)
			existingTokenService := tokenService.DeepCopy()
			resources.CopyTokenService(desiredTokenService, existingTokenService)
			return r.meshClient.MeshV1alpha2().TokenServices(cell.Namespace).Update(existingTokenService)
		}(cell, tokenService)
		if err != nil {
			r.logger.Errorf("Failed to update TokenService %q: %v", tokenServiceName, err)
			return err
		}
	}
	resources.StatusFromTokenService(cell, tokenService)
	return nil
}

func (r *reconciler) reconcileComponent(cell *v1alpha2.Cell, componentTemplate *v1alpha2.Component) error {
	componentName := resources.ComponentName(cell, componentTemplate)
	component, err := r.componentLister.Components(cell.Namespace).Get(componentName)
	if errors.IsNotFound(err) {
		component, err = r.meshClient.MeshV1alpha2().Components(cell.Namespace).Create(resources.MakeComponent(cell, componentTemplate))
		if err != nil {
			r.logger.Errorf("Failed to create Component %q: %v", componentName, err)
			r.recorder.Eventf(cell, corev1.EventTypeWarning, "CreationFailed", "Failed to create Component %q: %v", componentName, err)
			return err
		}
		r.recorder.Eventf(cell, corev1.EventTypeNormal, "Created", "Created Component %q", componentName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Component %q: %v", componentName, err)
		return err
	} else if !metav1.IsControlledBy(component, cell) {
		return fmt.Errorf("cell: %q does not own the Component: %q", cell.Name, componentName)
	} else {
		component, err = func(cell *v1alpha2.Cell, component *v1alpha2.Component) (*v1alpha2.Component, error) {
			if !resources.RequireComponentUpdate(cell, component) {
				return component, nil
			}
			desiredComponent := resources.MakeComponent(cell, componentTemplate)
			existingComponent := component.DeepCopy()
			resources.CopyComponent(desiredComponent, existingComponent)
			return r.meshClient.MeshV1alpha2().Components(cell.Namespace).Update(existingComponent)
		}(cell, component)
		if err != nil {
			r.logger.Errorf("Failed to update Component %q: %v", componentName, err)
			return err
		}
	}
	resources.StatusFromComponent(cell, component)
	return nil
}

// func (r *reconciler) reconcileComponents(cell *v1alpha2.Cell) error {
// 	cell.Status.ServiceCount = 0
// 	for _, desiredComponent := range cell.Spec.Components {
// 		component, err := r.componentLister.Components(cell.Namespace).Get(resources.ComponentName(cell, desiredComponent))
// 		if errors.IsNotFound(err) {
// 			component, err = r.meshClient.MeshV1alpha2().Components(cell.Namespace).Create(resources.MakeComponent(cell, desiredComponent))
// 			if err != nil {
// 				r.logger.Errorf("Failed to create Service: %s : %v", desiredComponent.Name, err)
// 				return err
// 			}
// 			r.logger.Debugw("Service created", resources.ComponentName(cell, desiredComponent), component)
// 		} else if err != nil {
// 			return err
// 		}
// 		if component != nil {
// 			// service exists. if the new obj is not equal to old one, perform an update.
// 			newComponent := resources.MakeComponent(cell, desiredComponent)
// 			// set the previous service's `ResourceVersion` to the newService
// 			// Else the issue `metadata.resourceVersion: Invalid value: 0x0: must be specified for an update` will occur.
// 			newComponent.ResourceVersion = component.ResourceVersion
// 			if !isEqual(component, newComponent) {
// 				component, err = r.meshClient.MeshV1alpha2().Components(cell.Namespace).Update(newComponent)
// 				if err != nil {
// 					r.logger.Errorf("Failed to update Service: %s : %v", component.Name, err)
// 					return err
// 				}
// 				r.logger.Debugw("Service updated", resources.ComponentName(cell, desiredComponent), component)
// 			}
// 		}
// 		if component.Status.AvailableReplicas > 0 || component.Spec.ScalingPolicy.IsKpa() || component.Spec.Type == v1alpha2.ComponentTypeJob {
// 			cell.Status.ServiceCount++
// 		}
// 	}
// 	return nil
// }

// func (r *reconciler) reconcileVirtualService(cell *v1alpha2.Cell) error {
// 	cellVs, err := r.istioVirtualServiceLister.VirtualServices(cell.Namespace).Get(resources.CellVirtualServiceName(cell))
// 	if errors.IsNotFound(err) {
// 		cellVs, err = resources.CreateCellVirtualService(cell, r.cellLister)
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Cell VS object %v for instance %s", err, cell.Name)
// 			return err
// 		}
// 		if cellVs == nil {
// 			r.logger.Debugf("No VirtualService created for cell instance %s", cell.Name)
// 			return nil
// 		}
// 		lastAppliedConfig, err := json.Marshal(controller_commons.BuildVirtualServiceLastAppliedConfig(cellVs))
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Cell VS %v for instance %s", err, cell.Name)
// 			return err
// 		}
// 		controller_commons.Annotate(cellVs, corev1.LastAppliedConfigAnnotation, string(lastAppliedConfig))
// 		cellVs, err = r.meshClient.NetworkingV1alpha3().VirtualServices(cell.Namespace).Create(cellVs)
// 		if err != nil {
// 			r.logger.Errorf("Failed to create Cell VirtualService %v for instance %s", err, cell.Name)
// 			return err
// 		}
// 		r.logger.Debugw("Cell VirtualService created", resources.CellVirtualServiceName(cell), cellVs)
// 	} else if err != nil {
// 		return err
// 	}
// 	return nil
// }

func isEqual(oldService *v1alpha2.Component, newService *v1alpha2.Component) bool {
	// we only consider equality of the spec
	return reflect.DeepEqual(oldService.Spec, newService.Spec)
}

func (r *reconciler) updateStatus(desired *v1alpha2.Cell) (*v1alpha2.Cell, error) {
	cell, err := r.cellLister.Cells(desired.Namespace).Get(desired.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(cell.Status, desired.Status) {
		latest := cell.DeepCopy()
		latest.Status = desired.Status
		return r.meshClient.MeshV1alpha2().Cells(desired.Namespace).UpdateStatus(latest)
	}
	return desired, nil
}
