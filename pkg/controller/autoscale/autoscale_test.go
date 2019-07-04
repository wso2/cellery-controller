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
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	autoscalingV2Beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

var onePtr int32 = 1

func TestBuildAutoscalePolicyLastAppliedConfig(t *testing.T) {
	scalePolicy := &v1alpha1.AutoscalePolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AutoscalePolicy",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "testinst--svc-autoscalepolicy",
			Namespace:       "default",
			Labels:          nil,
			OwnerReferences: nil,
		},
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: true,
			Policy: v1alpha1.Policy{
				ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
					Kind:       "Deployment",
					Name:       "testint--svc-deployment",
					APIVersion: "apps/v1",
				},
				MinReplicas: "1",
				MaxReplicas: 1,
				Metrics:     nil,
			},
		},
	}
	lastAppliedConfig := BuildAutoscalePolicyLastAppliedConfig(scalePolicy)
	if diff := cmp.Diff(scalePolicy, lastAppliedConfig); diff != "" {
		t.Errorf("CreateAutoscalePolicy (-expected, +actual)\n%v", diff)
	}
}

func TestAnnotate(t *testing.T) {
	scalePolicy := &v1alpha1.AutoscalePolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AutoscalePolicy",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "testinst--svc-autoscalepolicy",
			Namespace:       "default",
			Labels:          nil,
			OwnerReferences: nil,
			Annotations: map[string]string{
				"cellery.org/simple-annotation": "simple-annotation-value",
			},
		},
		Spec: v1alpha1.AutoscalePolicySpec{
			Overridable: true,
			Policy: v1alpha1.Policy{
				ScaleTargetRef: autoscalingV2Beta1.CrossVersionObjectReference{
					Kind:       "Deployment",
					Name:       "testint--svc-deployment",
					APIVersion: "apps/v1",
				},
				MinReplicas: "1",
				MaxReplicas: 1,
				Metrics:     nil,
			},
		},
	}
	lastAppliedConfig, err := json.Marshal(BuildAutoscalePolicyLastAppliedConfig(scalePolicy))
	if err != nil {
		t.Errorf("Error %+v", err)

	}
	Annotate(scalePolicy, corev1.LastAppliedConfigAnnotation, string(lastAppliedConfig))
	if diff := cmp.Diff(scalePolicy.Annotations[corev1.LastAppliedConfigAnnotation], string(lastAppliedConfig)); diff != "" {
		t.Errorf("Annotate (-expected, +actual)\n%v", diff)
	}
}
