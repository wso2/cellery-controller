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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Taken from: https://github.com/istio/istio/blob/1.0.2/vendor/istio.io/api/networking/v1alpha3/gateway.pb.go

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GatewaySpec `json:"spec"`
}

type GatewaySpec struct {
	// REQUIRED: A list of server specifications.
	Servers []*Server `json:"servers,omitempty"`
	// REQUIRED: One or more labels that indicate a specific set of pods/VMs
	// on which this gateway configuration should be applied.
	// The scope of label search is platform dependent.
	// On Kubernetes, for example, the scope includes pods running in
	// all reachable namespaces.
	Selector map[string]string `json:"selector,omitempty"`
}

type Server struct {
	// REQUIRED: The Port on which the proxy should listen for incoming
	// connections
	Port *Port `json:"port,omitempty"`
	// REQUIRED. A list of hosts exposed by this gateway. At least one
	// host is required. While typically applicable to
	// HTTP services, it can also be used for TCP services using TLS with
	// SNI. May contain a wildcard prefix for the bottom-level component of
	// a domain name. For example `*.foo.com` matches `bar.foo.com`
	// and `*.com` matches `bar.foo.com`, `example.com`, and so on.
	//
	// **Note**: A `VirtualService` that is bound to a gateway must have one
	// or more hosts that match the hosts specified in a server. The match
	// could be an exact match or a suffix match with the server's hosts. For
	// example, if the server's hosts specifies "*.example.com",
	// VirtualServices with hosts dev.example.com, prod.example.com will
	// match. However, VirtualServices with hosts example.com or
	// newexample.com will not match.
	Hosts []string `json:"hosts,omitempty"`
	// Set of TLS related options that govern the server's behavior. Use
	// these options to control if all http requests should be redirected to
	// https, and the TLS modes to use.
	Tls *Server_TLSOptions `json:"tls,omitempty"`
}

// Port describes the properties of a specific port of a service.
type Port struct {
	// REQUIRED: A valid non-negative integer port number.
	Number uint32 `json:"number,omitempty"`
	// REQUIRED: The protocol exposed on the port.
	// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP|TLS.
	// TLS is used to indicate secure connections to non HTTP services.
	Protocol string `json:"protocol,omitempty"`
	// Label assigned to the port.
	Name string `json:"name,omitempty"`
}

type Server_TLSOptions struct {
	// If set to true, the load balancer will send a 301 redirect for all
	// http connections, asking the clients to use HTTPS.
	HttpsRedirect bool `json:"httpsRedirect,omitempty"`
	// Optional: Indicates whether connections to this port should be
	// secured using TLS. The value of this field determines how TLS is
	// enforced.
	Mode string `json:"mode,omitempty"`
	// REQUIRED if mode is `SIMPLE` or `MUTUAL`. The path to the file
	// holding the server-side TLS certificate to use.
	ServerCertificate string `json:"serverCertificate,omitempty"`
	// REQUIRED if mode is `SIMPLE` or `MUTUAL`. The path to the file
	// holding the server's private key.
	PrivateKey string `json:"privateKey,omitempty"`
	// REQUIRED if mode is `MUTUAL`. The path to a file containing
	// certificate authority certificates to use in verifying a presented
	// client side certificate.
	CaCertificates string `json:"caCertificates,omitempty"`
	// A list of alternate names to verify the subject identity in the
	// certificate presented by the client.
	SubjectAltNames []string `json:"subjectAltNames,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Gateway `json:"items"`
}
