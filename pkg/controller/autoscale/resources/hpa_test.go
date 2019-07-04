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
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

var intOne int32 = 1
var intFifty int32 = 50

func TestCreateHpa(t *testing.T) {
	scalePolicy := &v1alpha1.AutoscalePolicy{
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: true,
			Policy: v1alpha1.Policy{
				MinReplicas: "1",
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
	}

	hpa := CreateHpa(scalePolicy)

	min := scalePolicy.Spec.MinReplicas()
	expected := &autoscalingV2Beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      HpaName(scalePolicy),
			Namespace: scalePolicy.Namespace,
			Labels:    scalePolicy.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateAutoscalerOwnerRef(scalePolicy),
			},
		},
		Spec: autoscalingV2Beta1.HorizontalPodAutoscalerSpec{
			MinReplicas:    &min,
			MaxReplicas:    scalePolicy.Spec.Policy.MaxReplicas,
			ScaleTargetRef: scalePolicy.Spec.Policy.ScaleTargetRef,
			Metrics:        scalePolicy.Spec.Policy.Metrics,
		},
	}

	if diff := cmp.Diff(expected, hpa); diff != "" {
		t.Errorf("CreateHpa (-expected, +actual)\n%v", diff)
	}
}

func TestHpaName(t *testing.T) {
	scalePolicy := &v1alpha1.AutoscalePolicy{
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: true,
			Policy: v1alpha1.Policy{
				MinReplicas: "1",
				MaxReplicas: 5,
				Metrics: []autoscalingV2Beta1.MetricSpec{
					autoscalingV2Beta1.MetricSpec{
						Type: autoscalingV2Beta1.MetricSourceType("Resource"),
						Resource: &autoscalingV2Beta1.ResourceMetricSource{
							Name:                     "cpu",
							TargetAverageUtilization: &intFifty,
						},
					},
				},
			},
		},
	}
	name := HpaName(scalePolicy)
	expected := scalePolicy.Name + "-hpa"
	if name != expected {
		t.Errorf("HpaName incorrect, got: %v, expected: %v", name, expected)
	}
}
