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

// Taken from: https://github.com/istio/istio/blob/1.0.2/vendor/istio.io/api/networking/v1alpha3/virtual_service.pb.go

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VirtualServiceSpec `json:"spec"`
}

type VirtualServiceSpec struct {
	// REQUIRED. The destination hosts to which traffic is being sent. Could
	// be a DNS name with wildcard prefix or an IP address.  Depending on the
	// platform, short-names can also be used instead of a FQDN (i.e. has no
	// dots in the name). In such a scenario, the FQDN of the host would be
	// derived based on the underlying platform.
	//
	// **A host name can be defined by only one VirtualService**. A single
	// VirtualService can be used to describe traffic properties for multiple
	// HTTP and TCP ports.
	//
	// *Note for Kubernetes users*: When short names are used (e.g. "reviews"
	// instead of "reviews.default.svc.cluster.local"), Istio will interpret
	// the short name based on the namespace of the rule, not the service. A
	// rule in the "default" namespace containing a host "reviews will be
	// interpreted as "reviews.default.svc.cluster.local", irrespective of
	// the actual namespace associated with the reviews service. _To avoid
	// potential misconfigurations, it is recommended to always use fully
	// qualified domain names over short names._
	//
	// The hosts field applies to both HTTP and TCP services. Service inside
	// the mesh, i.e., those found in the service registry, must always be
	// referred to using their alphanumeric names. IP addresses are allowed
	// only for services defined via the Gateway.
	Hosts []string `json:"hosts,omitempty"`
	// The names of gateways and sidecars that should apply these routes. A
	// single VirtualService is used for sidecars inside the mesh as well as
	// for one or more gateways. The selection condition imposed by this
	// field can be overridden using the source field in the match conditions
	// of protocol-specific routes. The reserved word `mesh` is used to imply
	// all the sidecars in the mesh. When this field is omitted, the default
	// gateway (`mesh`) will be used, which would apply the rule to all
	// sidecars in the mesh. If a list of gateway names is provided, the
	// rules will apply only to the gateways. To apply the rules to both
	// gateways and sidecars, specify `mesh` as one of the gateway names.
	Gateways []string `json:"gateways,omitempty"`
	// An ordered list of route rules for HTTP traffic. HTTP routes will be
	// applied to platform service ports named 'http-*'/'http2-*'/'grpc-*', gateway
	// ports with protocol HTTP/HTTP2/GRPC/ TLS-terminated-HTTPS and service
	// entry ports using HTTP/HTTP2/GRPC protocols.  The first rule matching
	// an incoming request is used.
	Http []*HTTPRoute `json:"http,omitempty"`
	// An ordered list of route rule for non-terminated TLS & HTTPS
	// traffic. Routing is typically performed using the SNI value presented
	// by the ClientHello message. TLS routes will be applied to platform
	// service ports named 'https-*', 'tls-*', unterminated gateway ports using
	// HTTPS/TLS protocols (i.e. with "passthrough" TLS mode) and service
	// entry ports using HTTPS/TLS protocols.  The first rule matching an
	// incoming request is used.  NOTE: Traffic 'https-*' or 'tls-*' ports
	// without associated virtual service will be treated as opaque TCP
	// traffic.
	//Tls []*TLSRoute `protobuf:"bytes,5,rep,name=tls" json:"tls,omitempty"`
	// An ordered list of route rules for opaque TCP traffic. TCP routes will
	// be applied to any port that is not a HTTP or TLS port. The first rule
	// matching an incoming request is used.
	//Tcp []*TCPRoute `protobuf:"bytes,4,rep,name=tcp" json:"tcp,omitempty"`
}

type HTTPRoute struct {
	// Match conditions to be satisfied for the rule to be
	// activated. All conditions inside a single match block have AND
	// semantics, while the list of match blocks have OR semantics. The rule
	// is matched if any one of the match blocks succeed.
	Match []*HTTPMatchRequest `json:"match,omitempty"`
	// A http rule can either redirect or forward (default) traffic. The
	// forwarding target can be one of several versions of a service (see
	// glossary in beginning of document). Weights associated with the
	// service version determine the proportion of traffic it receives.
	Route []*DestinationWeight `json:"route,omitempty"`
	// Rewrite HTTP URIs and Authority headers. Rewrite cannot be used with
	// Redirect primitive. Rewrite will be performed before forwarding.
	Rewrite *HTTPRewrite `json:"rewrite,omitempty"`
}

type HTTPMatchRequest struct {
	// URI to match
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	Uri *StringMatch `json:"uri,omitempty"`
	// URI Scheme
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	Scheme *StringMatch `json:"scheme,omitempty"`
	// HTTP Method
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	Method *StringMatch `json:"method,omitempty"`
	// HTTP Authority
	// values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	Authority *StringMatch `json:"authority,omitempty"`
	// The header keys must be lowercase and use hyphen as the separator,
	// e.g. _x-request-id_.
	//
	// Header values are case-sensitive and formatted as follows:
	//
	// - `exact: "value"` for exact string match
	//
	// - `prefix: "value"` for prefix-based match
	//
	// - `regex: "value"` for ECMAscript style regex-based match
	//
	// **Note:** The keys `uri`, `scheme`, `method`, and `authority` will be ignored.
	Headers map[string]*StringMatch `json:"headers,omitempty"`
	// Specifies the ports on the host that is being addressed. Many services
	// only expose a single port or label ports with the protocols they support,
	// in these cases it is not required to explicitly select the port.
	Port uint32 `json:"port,omitempty"`
	// One or more labels that constrain the applicability of a rule to
	// workloads with the given labels. If the VirtualService has a list of
	// gateways specified at the top, it should include the reserved gateway
	// `mesh` in order for this field to be applicable.
	SourceLabels map[string]string `json:"sourceLabels,omitempty"`
	// Names of gateways where the rule should be applied to. Gateway names
	// at the top of the VirtualService (if any) are overridden. The gateway match is
	// independent of sourceLabels.
	Gateways []string `json:"gateways,omitempty"`
}

type StringMatch struct {
	Exact  string `json:"exact,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

type DestinationWeight struct {
	// REQUIRED. Destination uniquely identifies the instances of a service
	// to which the request/connection should be forwarded to.
	Destination *Destination `json:"destination,omitempty"`
	// REQUIRED. The proportion of traffic to be forwarded to the service
	// version. (0-100). Sum of weights across destinations SHOULD BE == 100.
	// If there is only destination in a rule, the weight value is assumed to
	// be 100.
	Weight int32 `json:"weight,omitempty"`
}

type HTTPRewrite struct {
	// rewrite the path (or the prefix) portion of the URI with this
	// value. If the original URI was matched based on prefix, the value
	// provided in this field will replace the corresponding matched prefix.
	Uri string `json:"uri,omitempty"`
	// rewrite the Authority/Host header with this value.
	Authority string `json:"authority,omitempty"`
}

type Destination struct {
	// REQUIRED. The name of a service from the service registry. Service
	// names are looked up from the platform's service registry (e.g.,
	// Kubernetes services, Consul services, etc.) and from the hosts
	// declared by [ServiceEntry](#ServiceEntry). Traffic forwarded to
	// destinations that are not found in either of the two, will be dropped.
	//
	// *Note for Kubernetes users*: When short names are used (e.g. "reviews"
	// instead of "reviews.default.svc.cluster.local"), Istio will interpret
	// the short name based on the namespace of the rule, not the service. A
	// rule in the "default" namespace containing a host "reviews will be
	// interpreted as "reviews.default.svc.cluster.local", irrespective of
	// the actual namespace associated with the reviews service. _To avoid
	// potential misconfigurations, it is recommended to always use fully
	// qualified domain names over short names._
	Host string `json:"host,omitempty"`
	// The name of a subset within the service. Applicable only to services
	// within the mesh. The subset must be defined in a corresponding
	// DestinationRule.
	Subset string `json:"subset,omitempty"`
	// Specifies the port on the host that is being addressed. If a service
	// exposes only a single port it is not required to explicitly select the
	// port.
	//Port *PortSelector `json:"port,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []VirtualService `json:"items"`
}
