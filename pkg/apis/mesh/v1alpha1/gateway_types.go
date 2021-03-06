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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec"`
	Status GatewayStatus `json:"status"`
}

type GatewayTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GatewaySpec `json:"spec,omitempty"`
}

type GatewaySpec struct {
	Type        GatewayType          `json:"type,omitempty"`
	Host        string               `json:"host,omitempty"`
	Tls         TlsConfig            `json:"tls,omitempty"`
	OidcConfig  *OidcConfig          `json:"oidc,omitempty"`
	HTTPRoutes  []HTTPRoute          `json:"http,omitempty"`
	GRPCRoutes  []GRPCRoute          `json:"grpc,omitempty"`
	TCPRoutes   []TCPRoute           `json:"tcp,omitempty"`
	Autoscaling *AutoscalePolicySpec `json:"autoscaling,omitempty"`
}

func (gw *GatewaySpec) Empty() bool {
	if len(gw.TCPRoutes) == 0 && len(gw.HTTPRoutes) == 0 && len(gw.GRPCRoutes) == 0 {
		return true
	}
	return false
}

type TlsConfig struct {
	Secret string `json:"secret,omitempty"`
	Key    string `json:"key,omitempty"`
	Cert   string `json:"cert,omitempty"`
}

type OidcConfig struct {
	ProviderUrl    string   `json:"providerUrl"`
	ClientId       string   `json:"clientId"`
	ClientSecret   string   `json:"clientSecret"`
	DcrUrl         string   `json:"dcrUrl"`
	DcrUser        string   `json:"dcrUser"`
	DcrPassword    string   `json:"dcrPassword"`
	RedirectUrl    string   `json:"redirectUrl"`
	BaseUrl        string   `json:"baseUrl"`
	SubjectClaim   string   `json:"subjectClaim"`
	SecurePaths    []string `json:"securePaths,omitempty"`
	NonSecurePaths []string `json:"nonSecurePaths,omitempty"`
}

type HTTPRoute struct {
	Context      string          `json:"context"`
	Definitions  []APIDefinition `json:"definitions"`
	Backend      string          `json:"backend"`
	Global       bool            `json:"global"`
	Authenticate bool            `json:"authenticate"`
	ZeroScale    bool            `json:"zeroScale,omitempty"`
}

type APIDefinition struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type GRPCRoute struct {
	Port        uint32 `json:"port"`
	BackendHost string `json:"backendHost"`
	BackendPort uint32 `json:"backendPort"`
	ZeroScale   bool   `json:"zeroScale,omitempty"`
}

type TCPRoute struct {
	Port        uint32 `json:"port"`
	BackendHost string `json:"backendHost"`
	BackendPort uint32 `json:"backendPort"`
}

type GatewayStatus struct {
	OwnerCell string `json:"ownerCell"`
	HostName  string `json:"hostname"`
	Status    string `json:"status"`
}

type GatewayType string

const (
	// GatewayTypeEnvoy uses envoy proxy as the gateway.
	GatewayTypeEnvoy GatewayType = "Envoy"

	// GatewayTypeMicroGateway uses WSO2 micro-gateway as the gateway.
	GatewayTypeMicroGateway GatewayType = "MicroGateway"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Gateway `json:"items"`
}
