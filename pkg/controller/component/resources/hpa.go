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

package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
)

func MakeHpa(component *v1alpha2.Component) *autoscalingv2beta1.HorizontalPodAutoscaler {
	return &autoscalingv2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HpaName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
			MinReplicas:    component.Spec.ScalingPolicy.Hpa.MinReplicas,
			MaxReplicas:    component.Spec.ScalingPolicy.Hpa.MaxReplicas,
			ScaleTargetRef: makeTargetRef(component),
			Metrics:        component.Spec.ScalingPolicy.Hpa.Metrics,
		},
	}
}

func makeTargetRef(component *v1alpha2.Component) autoscalingv2beta1.CrossVersionObjectReference {
	kind := "Deployment"
	name := DeploymentName(component)
	return autoscalingv2beta1.CrossVersionObjectReference{
		Kind:       kind,
		Name:       name,
		APIVersion: appsv1.SchemeGroupVersion.String(),
	}
}

func RequireHpa(component *v1alpha2.Component) bool {
	return (component.Spec.Type == v1alpha2.ComponentTypeDeployment || component.Spec.Type == v1alpha2.ComponentTypeStatefulSet) && component.Spec.ScalingPolicy.IsHpa()
}

func RequireHpaUpdate(component *v1alpha2.Component, hpa *autoscalingv2beta1.HorizontalPodAutoscaler) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		hpa.Generation != component.Status.HpaGeneration
}

func CopyHpa(source, destination *autoscalingv2beta1.HorizontalPodAutoscaler) {
	destination.Spec.ScaleTargetRef = source.Spec.ScaleTargetRef
	destination.Spec.MinReplicas = source.Spec.MinReplicas
	destination.Spec.MaxReplicas = source.Spec.MaxReplicas
	destination.Spec.Metrics = source.Spec.Metrics
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromHpa(component *v1alpha2.Component, hpa *autoscalingv2beta1.HorizontalPodAutoscaler) {
	component.Status.HpaGeneration = hpa.Generation
}
