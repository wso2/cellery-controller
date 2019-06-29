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

package autoscale

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	autoscalingV2beta1Informer "k8s.io/client-go/informers/autoscaling/v2beta1"
	"k8s.io/client-go/kubernetes"
	autoscalingV2beta1Lister "k8s.io/client-go/listers/autoscaling/v2beta1"
	"k8s.io/client-go/tools/cache"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	meshclientset "github.com/cellery-io/mesh-controller/pkg/client/clientset/versioned"
	meshinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/mesh/v1alpha1"
	listers "github.com/cellery-io/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/autoscale/resources"
)

type autoscalePolicyeHandler struct {
	kubeClient            kubernetes.Interface
	meshClient            meshclientset.Interface
	autoscalePolicyLister listers.AutoscalePolicyLister
	hpaLister             autoscalingV2beta1Lister.HorizontalPodAutoscalerLister
	serviceLister         listers.ServiceLister
	gatewayLister         listers.GatewayLister
	logger                *zap.SugaredLogger
}

func NewController(
	kubeClient kubernetes.Interface,
	meshClient meshclientset.Interface,
	autoscalePolicyInformer meshinformers.AutoscalePolicyInformer,
	hpaInformer autoscalingV2beta1Informer.HorizontalPodAutoscalerInformer,
	serviceInformer meshinformers.ServiceInformer,
	gatewayInformer meshinformers.GatewayInformer,
	logger *zap.SugaredLogger,
) *controller.Controller {
	h := &autoscalePolicyeHandler{
		kubeClient:            kubeClient,
		meshClient:            meshClient,
		autoscalePolicyLister: autoscalePolicyInformer.Lister(),
		hpaLister:             hpaInformer.Lister(),
		serviceLister:         serviceInformer.Lister(),
		gatewayLister:         gatewayInformer.Lister(),
		logger:                logger.Named("autoscale-policy-controller"),
	}
	c := controller.New(h, h.logger, "AutoscalePolicy")

	h.logger.Info("Setting up event handlers")
	autoscalePolicyInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			h.logger.Infow("Informer update", "old", old, "new", new)
			c.Enqueue(new)
		},
		DeleteFunc: c.Enqueue,
	})
	return c
}

func (h *autoscalePolicyeHandler) Handle(key string) error {
	h.logger.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		h.logger.Errorf("invalid resource key: %s", key)
		return nil
	}
	autoscalePolicyOrig, err := h.autoscalePolicyLister.AutoscalePolicies(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("Autoscale policy '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	h.logger.Debugw("lister instance", key, autoscalePolicyOrig)
	autoscalePolicy := autoscalePolicyOrig.DeepCopy()

	if err = h.handle(autoscalePolicy); err != nil {
		return err
	}

	return nil
}

func (h *autoscalePolicyeHandler) handle(autoscalePolicy *v1alpha1.AutoscalePolicy) error {
	// if the autoscale policy does not have a parent reference, add it
	if len(autoscalePolicy.OwnerReferences) == 0 {
		err := buildParentReference(autoscalePolicy, h.serviceLister, h.gatewayLister)
		if err != nil {
			h.logger.Errorf("Failed to create owner reference for autoscale policy %s, error: %v", autoscalePolicy.Name, err)
			return err
		}
		autoscalePolicy.Status = "Ready"
		_, err = h.meshClient.MeshV1alpha1().AutoscalePolicies(autoscalePolicy.Namespace).Update(autoscalePolicy)
		if err != nil {
			h.logger.Errorf("Failed to update Autoscale policy %s after owner reference creation %v", autoscalePolicy.Name, err)
			return err
		}
	}
	hpa, err := h.hpaLister.HorizontalPodAutoscalers(autoscalePolicy.Namespace).Get(resources.HpaName(autoscalePolicy))
	if errors.IsNotFound(err) {
		hpa, err = h.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(autoscalePolicy.Namespace).
			Create(resources.CreateHpa(autoscalePolicy))
		if err != nil {
			h.logger.Errorf("Failed to create HPA %v", err)
			return err
		}
		h.logger.Infow("HPA created", resources.HpaName(autoscalePolicy), hpa)
	} else if err != nil {
		return err
	}
	if hpa != nil {
		// HPA already exists, if its not equal to the current object, update it.
		newHpa := resources.CreateHpa(autoscalePolicy)
		if !isEqual(hpa, newHpa) {
			updatedHpa, err := h.kubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(hpa.Namespace).Update(newHpa)
			if err != nil {
				h.logger.Errorf("Failed to update HPA %v", err)
				return err
			}
			h.logger.Debugw("HPA updated", resources.HpaName(autoscalePolicy), updatedHpa)
		}
	}
	return nil
}

func isEqual(oldHpa *autoscalingV2Beta1.HorizontalPodAutoscaler, newHpa *autoscalingV2Beta1.HorizontalPodAutoscaler) bool {
	return reflect.DeepEqual(oldHpa.Spec.MinReplicas, newHpa.Spec.MinReplicas) &&
		reflect.DeepEqual(oldHpa.Spec.MaxReplicas, newHpa.Spec.MaxReplicas) &&
		reflect.DeepEqual(oldHpa.Spec.Metrics, newHpa.Spec.Metrics)
}

func Annotate(autoscalePolicy *v1alpha1.AutoscalePolicy, name string, value string) {
	annotations := make(map[string]string, len(autoscalePolicy.ObjectMeta.Annotations)+1)
	annotations[name] = value
	for k, v := range autoscalePolicy.ObjectMeta.Annotations {
		annotations[k] = v
	}
	autoscalePolicy.Annotations = annotations
}

func BuildAutoscalePolicyLastAppliedConfig(autoscalePolicy *v1alpha1.AutoscalePolicy) *v1alpha1.AutoscalePolicy {
	return &v1alpha1.AutoscalePolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AutoscalePolicy",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      autoscalePolicy.Name,
			Namespace: autoscalePolicy.Namespace,
		},
		Spec: autoscalePolicy.Spec,
	}
}

func buildParentReference(autoscalePolicy *v1alpha1.AutoscalePolicy, serviceLister listers.ServiceLister, gatewayLister listers.GatewayLister) error {
	// extract the service name and lookup the service
	parts := strings.Split(autoscalePolicy.Spec.Policy.ScaleTargetRef.Name, "-deployment")
	if len(parts) < 2 {
		return fmt.Errorf("Unable to extract the service name scale target reference of autoscale policy %s", autoscalePolicy.Name)
	}
	if csvc, err := serviceLister.Services(autoscalePolicy.Namespace).Get(parts[0]); err == nil {
		autoscalePolicy.OwnerReferences = []metav1.OwnerReference{*controller.CreateServiceOwnerRef(csvc)}
	} else {
		if errors.IsNotFound(err) {
			// could be a autoscale policy for a gateway
			if gw, gWerr := gatewayLister.Gateways(autoscalePolicy.Namespace).Get(parts[0]); gWerr == nil {
				autoscalePolicy.OwnerReferences = []metav1.OwnerReference{*controller.CreateGatewayOwnerRef(gw)}
			} else {
				return fmt.Errorf("Unable to retrieve the gateway %s for to create parent reference for autoscale policy %s", parts[0], autoscalePolicy.Name)
			}
		} else {
			return fmt.Errorf("Unable to retrieve the service %s for to create parent reference for autoscale policy %s", parts[0], autoscalePolicy.Name)
		}
	}
	return nil
}
