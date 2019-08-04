/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec   `json:"spec"`
	Status ServiceStatus `json:"status"`
}

// ServiceTemplateSpec describes the data a service should have when created from a template
type ServiceTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceSpec `json:"spec,omitempty"`
}

type ServiceSpec struct {
	Replicas           *int32               `json:"replicas"`
	ServicePort        int32                `json:"servicePort"`
	ServiceAccountName string               `json:"serviceAccountName"`
	Protocol           string               `json:"protocol"`
	Container          corev1.Container     `json:"container"`
	Autoscaling        *AutoscalePolicySpec `json:"autoscaling,omitempty"`
	Type               ServiceType          `json:"type,omitempty"`
}

type ServiceStatus struct {
	OwnerCell         string `json:"ownerCell"`
	AvailableReplicas int32  `json:"availableReplicas"`
	HostName          string `json:"hostname"`
}

func (s *ServiceSpec) IsZeroScaled() bool {
	if s.Autoscaling == nil {
		return false
	}
	min := s.Autoscaling.Policy.MinReplicas
	return *min == 0
}

type ServiceType string

const (
	// ServiceTypeDeployment is the default type which run as services.
	ServiceTypeDeployment ServiceType = "Deployment"

	// ServiceTypeJob is a job which run into completion.
	ServiceTypeJob ServiceType = "Job"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Service `json:"items"`
}
