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

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Composite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompositeSpec   `json:"spec"`
	Status CompositeStatus `json:"status"`
}

type CompositeSpec struct {
	Components []Component `json:"components"`
}

type CompositeStatus struct {
	ComponentCount         int                               `json:"componentCount"`
	ActiveComponentCount   int                               `json:"activeComponentCount"`
	Status                 CompositeCurrentStatus            `json:"status"`
	ComponentStatuses      map[string]ComponentCurrentStatus `json:"componentStatuses"`
	ObservedGeneration     int64                             `json:"observedGeneration,omitempty"`
	SecretGeneration       int64                             `json:"secretGeneration,omitempty"`
	TokenServiceGeneration int64                             `json:"tokenServiceGeneration,omitempty"`
	ComponentGenerations   map[string]int64                  `json:"componentGenerations,omitempty"`
	// Current conditions of the composite.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []CompositeCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type CompositeCurrentStatus string

const (
	CompositeCurrentStatusUnknown CompositeCurrentStatus = "Unknown"

	CompositeCurrentStatusReady CompositeCurrentStatus = "Ready"

	CompositeCurrentStatusNotReady CompositeCurrentStatus = "NotReady"
)

type CompositeCondition struct {
	Type   CompositeConditionType `json:"type"`
	Status corev1.ConditionStatus `json:"status"`
}

type CompositeConditionType string

const (
	CompositeReady CompositeConditionType = "Ready"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CompositeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Composite `json:"items"`
}
