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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TokenService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TokenServiceSpec   `json:"spec"`
	Status TokenServiceStatus `json:"status"`
}

type TokenServiceTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TokenServiceSpec `json:"spec,omitempty"`
}

type TokenServiceSpec struct {
	Selector       map[string]string `json:"selector,omitempty"`
	SecretName     string            `json:"secretName,omitempty"`
	InstanceName   string            `json:"instanceName,omitempty"`
	InterceptMode  InterceptMode     `json:"interceptMode,omitempty"`
	OpaPolicies    []OpaPolicy       `json:"opa,omitempty"`
	UnsecuredPaths []string          `json:"unsecuredPaths,omitempty"`
}

type OpaPolicy struct {
	Key    string `json:"key,omitempty"`
	Policy string `json:"regoPolicy,omitempty"`
}

type TokenServiceStatus struct {
	ObservedGeneration     int64 `json:"observedGeneration,omitempty"`
	DeploymentGeneration   int64 `json:"deploymentGeneration,omitempty"`
	ServiceGeneration      int64 `json:"serviceGeneration,omitempty"`
	ConfigMapGeneration    int64 `json:"configMapGeneration,omitempty"`
	OpaConfigMapGeneration int64 `json:"opaConfigMapGeneration,omitempty"`
	EnvoyFilterGeneration  int64 `json:"envoyFilterGeneration,omitempty"`
}

type InterceptMode string

const (
	// Intercept only the incoming traffic
	InterceptModeInbound InterceptMode = "Inbound"
	// Intercept only the outgoing traffic
	InterceptModeOutbound InterceptMode = "Outbound"
	// Intercept both incoming and outgoing traffic
	InterceptModeAny InterceptMode = "Any"
	// Do not intercept any traffic. This will disable the STS
	InterceptModeNone InterceptMode = "None"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TokenServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []TokenService `json:"items"`
}
