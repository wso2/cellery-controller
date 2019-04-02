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
	"testing"

	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

var intFifty int32 = 50

func TestCreateAutoscalePolicy(t *testing.T) {
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cell-svc",
			Namespace: "default",
		},
		Spec: v1alpha1.ServiceSpec{
			Replicas: &intOne,
			Autoscaling: &v1alpha1.AutoscalePolicySpec{
				Overridable: true,
				Policy: v1alpha1.Policy{
					MinReplicas: &intOne,
					MaxReplicas: 5,
					Metrics: []autoscalingV2Beta1.MetricSpec{
						autoscalingV2Beta1.MetricSpec{
							Type: autoscalingV2Beta1.MetricSourceType("Resournce"),
							Resource: &autoscalingV2Beta1.ResourceMetricSource{
								Name:                     "cpu",
								TargetAverageUtilization: &intFifty,
							},
						},
					},
				},
			},
		},
	}
	policy := CreateAutoscalePolicy(svc)

	expected := &v1alpha1.AutoscalePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceAutoscalePolicyName(svc),
			Namespace: svc.Namespace,
			Labels:    createLabels(svc),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(svc),
			},
		},
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: svc.Spec.Autoscaling.Overridable,
			Policy: v1alpha1.Policy{
				ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
					Kind:       scaleTargetDeploymentKind,
					Name:       ServiceDeploymentName(svc),
					APIVersion: appsv1.SchemeGroupVersion.String(),
				},
				MinReplicas: svc.Spec.Autoscaling.Policy.MinReplicas,
				MaxReplicas: svc.Spec.Autoscaling.Policy.MaxReplicas,
				Metrics:     svc.Spec.Autoscaling.Policy.Metrics,
			},
		},
	}
	if diff := cmp.Diff(expected, policy); diff != "" {
		t.Errorf("CreateAutoscalePolicy (-expected, +actual)\n%v", diff)
	}
}

func TestCreateDefaultAutoscalePolicy(t *testing.T) {
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cell-svc",
			Namespace: "default",
		},
		Spec: v1alpha1.ServiceSpec{
			Replicas: &intOne,
		},
	}
	policy := CreateDefaultAutoscalePolicy(svc)

	expected := &v1alpha1.AutoscalePolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AutoscalePolicy",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceAutoscalePolicyName(svc),
			Namespace: svc.Namespace,
			Labels:    createLabels(svc),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(svc),
			},
		},
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: true,
			Policy: v1alpha1.Policy{
				ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
					Kind:       scaleTargetDeploymentKind,
					Name:       ServiceDeploymentName(svc),
					APIVersion: appsv1.SchemeGroupVersion.String(),
				},
				MinReplicas: svc.Spec.Replicas,
				MaxReplicas: *svc.Spec.Replicas,
				Metrics:     []autoscalingV2Beta1.MetricSpec{},
			},
		},
	}

	if diff := cmp.Diff(expected, policy); diff != "" {
		t.Errorf("CreateAutoscalePolicy (-expected, +actual)\n%v", diff)
	}
}

func TestServiceAutoscalePolicyName(t *testing.T) {
	svc := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cell-svc",
			Namespace: "default",
		},
		Spec: v1alpha1.ServiceSpec{
			Replicas: &intOne,
			Autoscaling: &v1alpha1.AutoscalePolicySpec{
				Overridable: true,
				Policy: v1alpha1.Policy{
					MinReplicas: &intOne,
					MaxReplicas: 5,
					Metrics: []autoscalingV2Beta1.MetricSpec{
						autoscalingV2Beta1.MetricSpec{
							Type: autoscalingV2Beta1.MetricSourceType("Resournce"),
							Resource: &autoscalingV2Beta1.ResourceMetricSource{
								Name:                     "cpu",
								TargetAverageUtilization: &intFifty,
							},
						},
					},
				},
			},
		},
	}
	name := ServiceAutoscalePolicyName(svc)
	expected := svc.Name + "-autoscalepolicy"
	if name != expected {
		t.Errorf("ServiceAutoscalePolicyName incorrect, got: %v, expected: %v", name, expected)
	}
}
