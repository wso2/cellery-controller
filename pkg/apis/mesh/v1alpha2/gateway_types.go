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

type GatewaySpec struct {
	Ingress Ingress `json:"ingress,omitempty"`
}

type Ingress struct {
	IngressExtensions IngressExtensions `json:"extensions,omitempty"`
	HTTPRoutes        []HTTPRoute       `json:"http,omitempty"`
	GRPCRoutes        []GRPCRoute       `json:"grpc,omitempty"`
	TCPRoutes         []TCPRoute        `json:"tcp,omitempty"`
}

func (ing *Ingress) HasRoutes() bool {
	return len(ing.HTTPRoutes) > 0 || len(ing.GRPCRoutes) > 0 || len(ing.TCPRoutes) > 0
}

type IngressExtensions struct {
	ApiPublisher   *ApiPublisherConfig   `json:"apiPublisher,omitempty"`
	ClusterIngress *ClusterIngressConfig `json:"clusterIngress,omitempty"`
	OidcConfig     *OidcConfig           `json:"oidc,omitempty"`
}

func (ie *IngressExtensions) HasApiPublisher() bool {
	return ie.ApiPublisher != nil
}

func (ie *IngressExtensions) HasClusterIngress() bool {
	return ie.ClusterIngress != nil
}

func (ie *IngressExtensions) HasOidc() bool {
	return ie.OidcConfig != nil
}

type ApiPublisherConfig struct {
	Authenticate bool   `json:"authenticate"`
	Backend      string `json:"backend"`
	Context      string `json:"context"`
	Version      string `json:"version"`
}

func (ap *ApiPublisherConfig) HasVersion() bool {
	return len(ap.Version) > 0
}

type ClusterIngressConfig struct {
	Host string    `json:"host,omitempty"`
	Tls  TlsConfig `json:"tls,omitempty"`
}

func (ci *ClusterIngressConfig) HasSecret() bool {
	return len(ci.Tls.Secret) > 0
}

func (ci *ClusterIngressConfig) HasCertAndKey() bool {
	return len(ci.Tls.Key) > 0 && len(ci.Tls.Cert) > 0
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
	JwtIssuer      string   `json:"jwtIssuer"`
	JwtAudience    string   `json:"jwtAudience"`
	SecretName     string   `json:"secretName"`
	SecurePaths    []string `json:"securePaths,omitempty"`
	NonSecurePaths []string `json:"nonSecurePaths,omitempty"`
}

type Destination struct {
	Host string `json:"host,omitempty"`
	Port uint32 `json:"port,omitempty"`
}

type HTTPRoute struct {
	Context      string          `json:"context"`
	Version      string          `json:"version"`
	Definitions  []APIDefinition `json:"definitions"`
	Global       bool            `json:"global"`
	Authenticate bool            `json:"authenticate"`
	Port         uint32          `json:"port"`
	Destination  Destination     `json:"destination,omitempty"`
	ZeroScale    bool            `json:"zeroScale,omitempty"`
}

type APIDefinition struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type GRPCRoute struct {
	Port        uint32      `json:"port"`
	Destination Destination `json:"destination,omitempty"`
	ZeroScale   bool        `json:"zeroScale,omitempty"`
}

type TCPRoute struct {
	Port        uint32      `json:"port"`
	Destination Destination `json:"destination,omitempty"`
}

type GatewayStatus struct {
	PublisherStatus                PublisherCurrentStatus `json:"gatewayType"`
	ServiceName                    string                 `json:"serviceName"`
	Status                         GatewayCurrentStatus   `json:"status"`
	AvailableReplicas              int32                  `json:"availableReplicas"`
	ObservedGeneration             int64                  `json:"observedGeneration,omitempty"`
	DeploymentGeneration           int64                  `json:"deploymentGeneration,omitempty"`
	JobGeneration                  int64                  `json:"jobGeneration,omitempty"`
	ServiceGeneration              int64                  `json:"serviceGeneration,omitempty"`
	VirtualServiceGeneration       int64                  `json:"virtualServiceGeneration,omitempty"`
	IstioGatewayGeneration         int64                  `json:"istioGatewayGeneration,omitempty"`
	ClusterIngressGeneration       int64                  `json:"clusterIngressGeneration,omitempty"`
	ClusterIngressSecretGeneration int64                  `json:"clusterIngressSecretGeneration,omitempty"`
	OidcEnvoyFilterGeneration      int64                  `json:"oidcEnvoyFilterGeneration,omitempty"`
	ConfigMapGeneration            int64                  `json:"configMapGeneration,omitempty"`
}

func (gs *GatewayStatus) ResetServiceName() {
	gs.ServiceName = "<none>"
}

type GatewayCurrentStatus string

const (
	GatewayCurrentStatusUnknown GatewayCurrentStatus = "Unknown"

	GatewayCurrentStatusReady GatewayCurrentStatus = "Ready"

	GatewayCurrentStatusNotReady GatewayCurrentStatus = "NotReady"

	// GatewayCurrentStatusIdle GatewayCurrentStatus = "Idle"
)

type PublisherCurrentStatus string

const (
	PublisherCurrentStatusUnknown PublisherCurrentStatus = "Unknown"

	PublisherCurrentStatusRunning PublisherCurrentStatus = "Running"

	PublisherCurrentStatusSucceeded PublisherCurrentStatus = "Succeeded"

	PublisherCurrentStatusFailed PublisherCurrentStatus = "Failed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Gateway `json:"items"`
}
