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
	meshinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/mesh/v1alpha1"
	listers "github.com/cellery-io/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/autoscale/resources"
)

type autoscalePolicyeHandler struct {
	kubeClient            kubernetes.Interface
	autoscalePolicyLister listers.AutoscalePolicyLister
	hpaLister             autoscalingV2beta1Lister.HorizontalPodAutoscalerLister
	logger                *zap.SugaredLogger
}

func NewController(
	kubeClient kubernetes.Interface,
	autoscalePolicyInformer meshinformers.AutoscalePolicyInformer,
	hpaInformer autoscalingV2beta1Informer.HorizontalPodAutoscalerInformer,
	logger *zap.SugaredLogger,
) *controller.Controller {
	h := &autoscalePolicyeHandler{
		kubeClient:            kubeClient,
		autoscalePolicyLister: autoscalePolicyInformer.Lister(),
		hpaLister:             hpaInformer.Lister(),
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
