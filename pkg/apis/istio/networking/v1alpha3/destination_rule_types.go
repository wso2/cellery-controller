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

// Taken from: https://github.com/istio/istio/blob/1.0.2/vendor/istio.io/api/networking/v1alpha3/destination_rule.pb.go

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DestinationRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DestinationRuleSpec `json:"spec"`
}

type DestinationRuleSpec struct {
	// REQUIRED. The name of a service from the service registry. Service
	// names are looked up from the platform's service registry (e.g.,
	// Kubernetes services, Consul services, etc.) and from the hosts
	// declared by [ServiceEntries](#ServiceEntry). Rules defined for
	// services that do not exist in the service registry will be ignored.
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
	// Note that the host field applies to both HTTP and TCP services.
	Host string `json:"host,omitempty"`
	// Traffic policies to apply (load balancing policy, connection pool
	// sizes, outlier detection).
	TrafficPolicy *TrafficPolicy `json:"trafficPolicy,omitempty"`
	// One or more named sets that represent individual versions of a
	// service. Traffic policies can be overridden at subset level.
	//Subsets []*Subset `json:"subsets,omitempty"`
}

// Traffic policies to apply for a specific destination, across all
// destination ports. See DestinationRule for examples.
type TrafficPolicy struct {
	// Settings controlling the load balancer algorithms.
	LoadBalancer *LoadBalancerSettings `json:"loadBalancer,omitempty"`
	// Settings controlling the volume of connections to an upstream service
	//ConnectionPool *ConnectionPoolSettings `json:"connectionPool,omitempty"`
	// Settings controlling eviction of unhealthy hosts from the load balancing pool
	//OutlierDetection *OutlierDetection `json:"outlierDetection,omitempty"`
	// TLS related settings for connections to the upstream service.
	//Tls *TLSSettings `json:"tls,omitempty"`
	// Traffic policies specific to individual ports. Note that port level
	// settings will override the destination-level settings. Traffic
	// settings specified at the destination-level will not be inherited when
	// overridden by port-level settings, i.e. default values will be applied
	// to fields omitted in port-level traffic policies.
	PortLevelSettings []*TrafficPolicy_PortTrafficPolicy `json:"portLevelSettings,omitempty"`
}

type LoadBalancerSettings struct {
	// Upstream load balancing policy.
	//
	// Types that are valid to be assigned to LbPolicy:
	//	*LoadBalancerSettings_Simple
	//	*LoadBalancerSettings_ConsistentHash
	Simple string `json:"simple,omitempty"`
}

// Traffic policies that apply to specific ports of the service
type TrafficPolicy_PortTrafficPolicy struct {
	// Specifies the port name or number of a port on the destination service
	// on which this policy is being applied.
	//
	// Names must comply with DNS label syntax (rfc1035) and therefore cannot
	// collide with numbers. If there are multiple ports on a service with
	// the same protocol the names should be of the form <protocol-name>-<DNS
	// label>.
	Port *PortSelector `json:"port,omitempty"`
	// Settings controlling the load balancer algorithms.
	//LoadBalancer *LoadBalancerSettings `protobuf:"bytes,2,opt,name=load_balancer,json=loadBalancer" json:"load_balancer,omitempty"`
	// Settings controlling the volume of connections to an upstream service
	//ConnectionPool *ConnectionPoolSettings `protobuf:"bytes,3,opt,name=connection_pool,json=connectionPool" json:"connection_pool,omitempty"`
	// Settings controlling eviction of unhealthy hosts from the load balancing pool
	//OutlierDetection *OutlierDetection `protobuf:"bytes,4,opt,name=outlier_detection,json=outlierDetection" json:"outlier_detection,omitempty"`
	// TLS related settings for connections to the upstream service.
	Tls *TLSSettings `json:"tls,omitempty"`
}

type TLSSettings struct {
	// REQUIRED: Indicates whether connections to this port should be secured
	// using TLS. The value of this field determines how TLS is enforced.
	Mode string `json:"mode,omitempty"`
	// REQUIRED if mode is `MUTUAL`. The path to the file holding the
	// client-side TLS certificate to use.
	// Should be empty if mode is `ISTIO_MUTUAL`.
	ClientCertificate string `json:"clientCertificate,omitempty"`
	// REQUIRED if mode is `MUTUAL`. The path to the file holding the
	// client's private key.
	// Should be empty if mode is `ISTIO_MUTUAL`.
	PrivateKey string `json:"privateKey,omitempty"`
	// OPTIONAL: The path to the file containing certificate authority
	// certificates to use in verifying a presented server certificate. If
	// omitted, the proxy will not verify the server's certificate.
	// Should be empty if mode is `ISTIO_MUTUAL`.
	CaCertificates string `json:"caCertificates,omitempty"`
	// A list of alternate names to verify the subject identity in the
	// certificate. If specified, the proxy will verify that the server
	// certificate's subject alt name matches one of the specified values.
	// Should be empty if mode is `ISTIO_MUTUAL`.
	SubjectAltNames []string `json:"subjectAltNames,omitempty"`
	// SNI string to present to the server during TLS handshake.
	// Should be empty if mode is `ISTIO_MUTUAL`.
	Sni string `json:"sni,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DestinationRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []DestinationRule `json:"items"`
}
