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
	autoscalingV2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec"`
	Status ComponentStatus `json:"status"`
}

type ComponentSpec struct {
	Type           ComponentType          `json:"type,omitempty"`
	ScalingPolicy  ComponentScalingPolicy `json:"scalingPolicy,omitempty"`
	Template       corev1.PodSpec         `json:"template,omitempty"`
	Ports          []PortMapping          `json:"ports,omitempty"`
	VolumeClaims   []VolumeClaim          `json:"volumeClaims,omitempty"`
	Configurations []corev1.ConfigMap     `json:"configurations,omitempty"`
	Secrets        []corev1.Secret        `json:"secrets,omitempty"`
}

type ComponentScalingPolicy struct {
	Replicas *int32                   `json:"replicas,omitempty"`
	Hpa      *HorizontalPodAutoscaler `json:"hpa,omitempty"`
	Kpa      *KnativePodAutoscaler    `json:"kpa,omitempty"`
}

func (sp *ComponentScalingPolicy) MinReplicas() int32 {
	if sp.Hpa != nil {
		if sp.Hpa.MinReplicas != nil {
			return *sp.Hpa.MinReplicas
		}
		return 1
	}
	if sp.Replicas != nil {
		return *sp.Replicas
	}
	return 1
}

func (sp *ComponentScalingPolicy) IsHpa() bool {
	return sp != nil && sp.Hpa != nil
}

func (sp *ComponentScalingPolicy) IsKpa() bool {
	return !sp.IsHpa() && sp.Kpa != nil
}

type ReplicaRange struct {
	MinReplicas *int32 `json:"minReplicas,omitempty"`
	MaxReplicas int32  `json:"maxReplicas,omitempty"`
}

type HorizontalPodAutoscaler struct {
	Overridable  *bool `json:"overridable"`
	ReplicaRange `json:",inline"`
	Metrics      []autoscalingV2beta2.MetricSpec `json:"metrics,omitempty"`
}

type KnativePodAutoscaler struct {
	ReplicaRange `json:",inline"`
	Concurrency  int64             `json:"concurrency"`
	Selector     map[string]string `json:"selector,omitempty"`
}

type PortMapping struct {
	Name            string   `json:"name"`
	Protocol        Protocol `json:"protocol"`
	Port            int32    `json:"port"`
	TargetContainer string   `json:"targetContainer"`
	TargetPort      int32    `json:"targetPort"`
}

type VolumeClaim struct {
	Shared   bool                         `json:"shared"`
	Template corev1.PersistentVolumeClaim `json:"template"`
}

type ComponentStatus struct {
	Type                             ComponentType          `json:"componentType"`
	Status                           ComponentCurrentStatus `json:"status"`
	ServiceName                      string                 `json:"serviceName"`
	AvailableReplicas                int32                  `json:"availableReplicas"`
	ObservedGeneration               int64                  `json:"observedGeneration,omitempty"`
	DeploymentGeneration             int64                  `json:"deploymentGeneration,omitempty"`
	StatefulSetGeneration            int64                  `json:"statefulSetGeneration,omitempty"`
	JobGeneration                    int64                  `json:"jobGeneration,omitempty"`
	ServiceGeneration                int64                  `json:"serviceGeneration,omitempty"`
	HpaGeneration                    int64                  `json:"hpaGeneration,omitempty"`
	ServingConfigurationGeneration   int64                  `json:"servingConfigurationGeneration,omitempty"`
	ServingVirtualServiceGeneration  int64                  `json:"servingVirtualServiceGeneration,omitempty"`
	TlsPolicyGeneration              int64                  `json:"TlsPolicyGeneration,omitempty"`
	PersistantVolumeClaimGenerations map[string]int64       `json:"persistantVolumeClaimGenerations,omitempty"`
	ConfigMapGenerations             map[string]int64       `json:"configMapGenerations,omitempty"`
	SecretGenerations                map[string]int64       `json:"secretGenerations,omitempty"`
}

func (cstat *ComponentStatus) SetType(t ComponentType) {
	switch t {
	case ComponentTypeDeployment:
		cstat.Type = ComponentTypeDeployment
	case ComponentTypeStatefulSet:
		cstat.Type = ComponentTypeStatefulSet
	case ComponentTypeJob:
		cstat.Type = ComponentTypeJob
	default:
		cstat.Type = "Unknown"
	}
}

func (cs *ComponentStatus) ResetServiceName() {
	cs.ServiceName = "<none>"
}

type Protocol string

const (
	// ProtocolTCP is the TCP protocol.
	ProtocolTCP Protocol = "TCP"

	// ProtocolHTTP is the HTTP protocol.
	ProtocolHTTP Protocol = "HTTP"

	// ProtocolGRPC is the GRPC protocol.
	ProtocolGRPC Protocol = "GRPC"
)

type ComponentType string

const (
	// ServiceTypeDeployment is the default type which run as services.
	ComponentTypeDeployment ComponentType = "Deployment"

	// ServiceTypeStatefulSet is the default type which runs for the stateful services.
	ComponentTypeStatefulSet ComponentType = "StatefulSet"

	// ServiceTypeJob is a job which run into completion.
	ComponentTypeJob ComponentType = "Job"
)

type ComponentCurrentStatus string

const (
	ComponentCurrentStatusUnknown ComponentCurrentStatus = "Unknown"

	ComponentCurrentStatusReady ComponentCurrentStatus = "Ready"

	ComponentCurrentStatusNotReady ComponentCurrentStatus = "NotReady"

	ComponentCurrentStatusIdle ComponentCurrentStatus = "Idle"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Component `json:"items"`
}
