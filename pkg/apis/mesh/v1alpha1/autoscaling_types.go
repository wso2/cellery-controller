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

package v1alpha1

import (
	autoscalingV2beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AutoscalePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoscalePolicySpec `json:"spec"`
	Status string              `json:"status"`
}

type AutoscalePolicySpec struct {
	Overridable bool   `json:"overridable,omitempty"`
	Policy      Policy `json:"policy"`
}

type Policy struct {
	ScaleTargetRef autoscalingV2beta1.CrossVersionObjectReference `json:"scaleTargetRef,omitempty"`
	MinReplicas    *int32                                         `json:"minReplicas"`
	MaxReplicas    int32                                          `json:"maxReplicas"`
	Metrics        []autoscalingV2beta1.MetricSpec                `json:"metrics"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AutoscalePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AutoscalePolicy `json:"items"`
}
