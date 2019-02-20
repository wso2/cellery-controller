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
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

const scaleTargetDeploymentKind = "Deployment"

func CreateHpa(service *v1alpha1.Service) *autoscalingV2Beta1.HorizontalPodAutoscaler {

	return &autoscalingV2Beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceHpaName(service),
			Namespace: service.Namespace,
			Labels:    createLabels(service),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(service),
			},
		},
		Spec: autoscalingV2Beta1.HorizontalPodAutoscalerSpec{
			MinReplicas: service.Spec.Autoscaling.Policy.MinReplicas,
			MaxReplicas: service.Spec.Autoscaling.Policy.MaxReplicas,
			ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
				Kind:       scaleTargetDeploymentKind,
				Name:       ServiceDeploymentName(service),
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
			Metrics: service.Spec.Autoscaling.Policy.Metrics,
		},
	}
}

func CreateDefaultHpa(service *v1alpha1.Service) *autoscalingV2Beta1.HorizontalPodAutoscaler {
	return &autoscalingV2Beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceHpaName(service),
			Namespace: service.Namespace,
			Labels:    createLabels(service),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(service),
			},
		},
		Spec: autoscalingV2Beta1.HorizontalPodAutoscalerSpec{
			MinReplicas: service.Spec.Replicas,
			MaxReplicas: *service.Spec.Replicas,
			ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
				Kind:       scaleTargetDeploymentKind,
				Name:       ServiceDeploymentName(service),
				APIVersion: appsv1.SchemeGroupVersion.String(),
			},
		},
	}
}
