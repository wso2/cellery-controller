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

type Cell struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CellSpec   `json:"spec"`
	Status CellStatus `json:"status"`
}

type CellSpec struct {
	Gateway      Gateway      `json:"gateway"`
	Components   []Component  `json:"components"`
	TokenService TokenService `json:"sts"`
}

type CellStatus struct {
	ComponentCount       int                               `json:"componentCount"`
	ActiveComponentCount int                               `json:"activeComponentCount"`
	GatewayServiceName   string                            `json:"gatewayServiceName"`
	Status               CellCurrentStatus                 `json:"status"`
	GatewayStatus        GatewayCurrentStatus              `json:"gatewayStatus"`
	ComponentStatuses    map[string]ComponentCurrentStatus `json:"componentStatuses"`
	// Status                  string `json:"status"`
	ObservedGeneration      int64            `json:"observedGeneration,omitempty"`
	NetworkPolicyGeneration int64            `json:"networkPolicyGeneration,omitempty"`
	SecretGeneration        int64            `json:"secretGeneration,omitempty"`
	GatewayGeneration       int64            `json:"gatewayGeneration,omitempty"`
	TokenServiceGeneration  int64            `json:"tokenServiceGeneration,omitempty"`
	ComponentGenerations    map[string]int64 `json:"componentGenerations,omitempty"`
	// Current conditions of the cell.
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []CellCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type CellCurrentStatus string

const (
	CellCurrentStatusUnknown CellCurrentStatus = "Unknown"

	CellCurrentStatusReady CellCurrentStatus = "Ready"

	CellCurrentStatusNotReady CellCurrentStatus = "NotReady"
)

type CellCondition struct {
	Type   CellConditionType      `json:"type"`
	Status corev1.ConditionStatus `json:"status"`
}

type CellConditionType string

const (
	CellReady CellConditionType = "Ready"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CellList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Cell `json:"items"`
}
