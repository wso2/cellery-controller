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
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateHpa(scalePolicy *v1alpha1.AutoscalePolicy) *autoscalingV2Beta1.HorizontalPodAutoscaler {

	minReplicas := scalePolicy.Spec.MinReplicas()
	return &autoscalingV2Beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HpaName(scalePolicy),
			Namespace: scalePolicy.Namespace,
			Labels:    scalePolicy.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateAutoscalerOwnerRef(scalePolicy),
			},
		},
		Spec: autoscalingV2Beta1.HorizontalPodAutoscalerSpec{
			MinReplicas:    &minReplicas,
			MaxReplicas:    scalePolicy.Spec.Policy.MaxReplicas,
			ScaleTargetRef: scalePolicy.Spec.Policy.ScaleTargetRef,
			Metrics:        scalePolicy.Spec.Policy.Metrics,
		},
	}
}

func HpaName(policy *v1alpha1.AutoscalePolicy) string {
	return policy.Name + "-hpa"
}
