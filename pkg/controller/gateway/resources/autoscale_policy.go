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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

const scaleTargetDeploymentKind = "Deployment"

func CreateAutoscalePolicy(gateway *v1alpha1.Gateway) *v1alpha1.AutoscalePolicy {

	return &v1alpha1.AutoscalePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayAutoscalePolicyName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: gateway.Spec.Autoscaling.Overridable,
			Policy: v1alpha1.Policy{
				ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
					Kind:       scaleTargetDeploymentKind,
					Name:       GatewayDeploymentName(gateway),
					APIVersion: appsv1.SchemeGroupVersion.String(),
				},
				MinReplicas: gateway.Spec.Autoscaling.Policy.MinReplicas,
				MaxReplicas: gateway.Spec.Autoscaling.Policy.MaxReplicas,
				Metrics:     gateway.Spec.Autoscaling.Policy.Metrics,
			},
		},
	}
}

func CreateDefaultAutoscalePolicy(gateway *v1alpha1.Gateway) *v1alpha1.AutoscalePolicy {

	var onePtr int32 = 1

	return &v1alpha1.AutoscalePolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AutoscalePolicy",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayAutoscalePolicyName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: true,
			Policy: v1alpha1.Policy{
				ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
					Kind:       scaleTargetDeploymentKind,
					Name:       GatewayDeploymentName(gateway),
					APIVersion: appsv1.SchemeGroupVersion.String(),
				},
				MinReplicas: &onePtr,
				MaxReplicas: 1,
				Metrics:     []autoscalingV2Beta1.MetricSpec{},
			},
		},
	}
}

func GatewayAutoscalePolicyName(gateway *v1alpha1.Gateway) string {
	return gateway.Name + "-autoscalepolicy"
}

func createLabels(gateway *v1alpha1.Gateway) map[string]string {
	labels := make(map[string]string, len(gateway.ObjectMeta.Labels)+2)

	labels[mesh.CellGatewayLabelKey] = gateway.Name
	labels[appLabelKey] = gateway.Name
	// order matters
	// todo: update the code if override is not possible
	for k, v := range gateway.ObjectMeta.Labels {
		labels[k] = v
	}
	return labels
}
