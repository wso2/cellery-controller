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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Taken from: https://github.com/istio/istio/blob/1.0.2/vendor/istio.io/api/networking/v1alpha3/envoy_filter.pb.go#L180

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EnvoyFilter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EnvoyFilterSpec `json:"spec"`
}

type EnvoyFilterSpec struct {
	WorkloadLabels map[string]string `json:"workloadLabels,omitempty"`
	Filters        []Filter         `json:"filters"`
}

type Filter struct {
	// Insert position in the filter chain. Defaults to FIRST
	InsertPosition InsertPosition `json:"insertPosition"`

	// Filter will be added to the listener only if the match conditions are true.
	// If not specified, the filters will be applied to all listeners.
	ListenerMatch ListenerMatch `json:"listenerMatch"`

	// REQUIRED: The type of filter to instantiate.
	FilterType string `json:"filterType"`

	// REQUIRED: The name of the filter to instantiate. The name must match a supported
	// filter _compiled into_ Envoy.
	FilterName string `json:"filterName"`

	// REQUIRED: Filter specific configuration which depends on the filter being
	// instantiated.
	FilterConfig FilterConfig `json:"filterConfig"`
}

// Indicates the relative index in the filter chain where the filter should be inserted.
type InsertPosition struct {
	// Position of this filter in the filter chain.
	Index string `json:"index,omitempty"`
	// If BEFORE or AFTER position is specified, specify the name of the
	// filter relative to which this filter should be inserted.
	RelativeTo string `json:"relativeTo,omitempty"`
}

// All conditions specified in the ListenerMatch must be met for the filter
// to be applied to a listener.
type ListenerMatch struct {
	// The service port/gateway port to which traffic is being
	// sent/received. If not specified, matches all listeners. Even though
	// inbound listeners are generated for the instance/pod ports, only
	// service ports should be used to match listeners.
	PortNumber uint32 `json:"portNumber,omitempty"`
	// Instead of using specific port numbers, a set of ports matching a
	// given port name prefix can be selected. E.g., "mongo" selects ports
	// named mongo-port, mongo, mongoDB, MONGO, etc. Matching is case
	// insensitive.
	PortNamePrefix string `json:"portNamePrefix,omitempty"`
	// Inbound vs outbound sidecar listener or gateway listener. If not specified,
	// matches all listeners.
	ListenerType string `json:"listenerType,omitempty"`
	// Selects a class of listeners for the same protocol. If not
	// specified, applies to listeners on all protocols. Use the protocol
	// selection to select all HTTP listeners (includes HTTP2/gRPC/HTTPS
	// where Envoy terminates TLS) or all TCP listeners (includes HTTPS
	// passthrough using SNI).
	ListenerProtocol string `json:"listenerProtocol,omitempty"`
	// One or more IP addresses to which the listener is bound. If
	// specified, should match at least one address in the list.
	Address []string `json:"address,omitempty"`
}

type FilterConfig struct {
	GRPCService GRPCService `json:"grpc_service"`
}

type GRPCService struct {
	GoogleGRPC GoogleGRPC `json:"google_grpc"`
	Timeout    string     `json:"timeout"`
}

type GoogleGRPC struct {
	TargetUri  string `json:"target_uri"`
	StatPrefix string `json:"stat_prefix"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EnvoyFilterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []EnvoyFilter `json:"items"`
}
