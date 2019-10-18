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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func RequireHpa(gw *v1alpha2.Gateway) bool {
	return gw.IsHpa()
}

func MakeHpa(gw *v1alpha2.Gateway) *autoscalingv2beta1.HorizontalPodAutoscaler {
	return &autoscalingv2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HpaName(gw),
			Namespace: gw.Namespace,
			Labels:    makeLabels(gw),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(gw),
			},
		},
		Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
			MinReplicas:    gw.MinReplicas(),
			MaxReplicas:    gw.MaxReplicas(),
			ScaleTargetRef: makeTargetRef(gw),
			Metrics:        gw.Metrics(),
		},
	}
}

func makeTargetRef(gw *v1alpha2.Gateway) autoscalingv2beta1.CrossVersionObjectReference {
	kind := "Deployment"
	name := DeploymentName(gw)
	return autoscalingv2beta1.CrossVersionObjectReference{
		Kind:       kind,
		Name:       name,
		APIVersion: appsv1.SchemeGroupVersion.String(),
	}
}

func RequireHpaUpdate(gw *v1alpha2.Gateway, hpa *autoscalingv2beta1.HorizontalPodAutoscaler) bool {
	return gw.Generation != gw.Status.ObservedGeneration ||
		hpa.Generation != gw.Status.HpaGeneration
}

func CopyHpa(source, destination *autoscalingv2beta1.HorizontalPodAutoscaler) {
	destination.Spec.ScaleTargetRef = source.Spec.ScaleTargetRef
	destination.Spec.MinReplicas = source.Spec.MinReplicas
	destination.Spec.MaxReplicas = source.Spec.MaxReplicas
	destination.Spec.Metrics = source.Spec.Metrics
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromHpa(gw *v1alpha2.Gateway, hpa *autoscalingv2beta1.HorizontalPodAutoscaler) {
	gw.Status.HpaGeneration = hpa.Generation
}
