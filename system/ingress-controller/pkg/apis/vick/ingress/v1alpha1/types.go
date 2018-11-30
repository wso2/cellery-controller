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
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Ingress struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	meta_v1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	meta_v1.ObjectMeta `json:"metadata,omitempty"`
	//
	// Spec is the Ingress resource spec
	Spec IngressSpec `json:"spec"`
}

type IngressSpec struct {
	Name string `json:"name"`
	Context string `json:"context"`
	Version string `json:"version"`
	IsDefaultVersion string `json:"isDefaultVersion,omitempty"`
	GatewayEnvironments string `json:"gatewayEnvironments,omitempty"`
	Visibility string `json:"visibility,omitempty"`
	ApiDefinition string `json:"apiDefinition,omitempty"`
	EndpointConfig string `json:"endpointConfig,omitempty"`
	AuthType string `json:"authType,omitempty"`
	SwaggerVersion string `json:"swaggerVersion,omitempty"`
	Labels []string `json:"labels,omitempty"`
	Transports []string `json:"transports,omitempty"`
	Tiers []string `json:"tiers,omitempty"`
	Endpoints []Endpoint `json:"endpoints"`
	Paths []Path `json:"paths"`
}

func (ingSpec IngressSpec) IsEqual(anIngSpec IngressSpec) bool {
	if (ingSpec.IsDefaultVersion == anIngSpec.IsDefaultVersion) && (ingSpec.GatewayEnvironments == anIngSpec.GatewayEnvironments) &&
		(ingSpec.Visibility == anIngSpec.Visibility) && (ingSpec.ApiDefinition == anIngSpec.ApiDefinition) &&
		(ingSpec.EndpointConfig == anIngSpec.EndpointConfig) && (ingSpec.AuthType == anIngSpec.AuthType) &&
		(ingSpec.SwaggerVersion == anIngSpec.SwaggerVersion) && areStringArraysEqual(ingSpec.Labels, anIngSpec.Labels) &&
		(areEndpointsEqual(ingSpec.Endpoints, anIngSpec.Endpoints)) && (areStringArraysEqual(ingSpec.Transports, anIngSpec.Transports)) &&
		areStringArraysEqual(ingSpec.Tiers, anIngSpec.Tiers) && arePathsEqual(ingSpec.Paths, anIngSpec.Paths) {
		return true
	}
	return false
}

func areStringArraysEqual (arr []string, otherArr []string) bool {
	if arr == nil && otherArr == nil {
		return true
	}

	if len(arr) != len(otherArr) {
		return false
	}
	var matchCount int = 0
	for _, label := range arr {
		for _, anotherLabel := range arr {
			if label == anotherLabel {
				matchCount++
			}
		}
	}
	if matchCount == len(arr) {
		return true
	}
	return false
}

type Endpoint struct {
	Protocol string `json:"protocol"`
	Type string `json:"type"`
	ServiceName string `json:"serviceName"`
	Port int `json:"port"`
}

func areEndpointsEqual (endpoints []Endpoint, otherEndpoints []Endpoint) bool {
	if endpoints == nil && otherEndpoints == nil {
		return true
	}

	if len(endpoints) != len(otherEndpoints) {
		return false
	}
	var matchCount int = 0
	for _, ep := range endpoints {
		for _, anotherEp := range otherEndpoints {
			if ep.isEqual(anotherEp) {
				matchCount++
			}
		}
	}
	if matchCount == len(endpoints) {
		return true
	}
	return false
}

func (ep Endpoint) isEqual (anotherEp Endpoint) bool {
	if (ep.Port == anotherEp.Port) && (ep.Type == anotherEp.Type) && (ep.ServiceName == anotherEp.ServiceName) &&
		(ep.Protocol == anotherEp.Protocol) {
		return true
	}
	return false
}

type Path struct {
	Context string `json:"context"`
	Operation string `json:"operation"`
	AuthType string `json:"authType,omitempty"`
	Tiers []string `json:"tiers,omitempty"`
	Parameters []Parameters `json:"parameters,omitempty"`
}

func arePathsEqual (paths []Path, otherPaths []Path) bool {
	if paths == nil && otherPaths == nil {
		return true
	}

	if len(paths) != len(otherPaths) {
		return false
	}
	var matchCount int = 0
	for _, path := range paths {
		for _, anotherPath := range otherPaths {
			if path.isEqual(anotherPath) {
				matchCount++
			}
		}
	}
	if matchCount == len(paths) {
		return true
	}
	return false
}

func (path Path) isEqual(aPath Path) bool {
	if (path.Context == aPath.Context) && (path.Operation == aPath.Operation) && areStringArraysEqual(path.Tiers, aPath.Tiers) &&
		areParamsEqual(path.Parameters, aPath.Parameters) && (path.AuthType == aPath.AuthType) {
		return true
	}
	return false
}

type Parameters struct {
	Name string `json:"name,omitempty"`
	In string `json:"in,omitempty"`
	Description string `json:"description,omitempty"`
	Required string `json:"required,omitempty"`
	Type string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
}

func areParamsEqual (paramsArr []Parameters, otherParamsArr []Parameters) bool {
	if paramsArr == nil && otherParamsArr == nil {
		return true
	}

	if len(paramsArr) != len(otherParamsArr) {
		return false
	}
	var matchCount int = 0
	for _, params := range paramsArr {
		for _, otherParams := range otherParamsArr {
			if params.isEqual(otherParams) {
				matchCount++
			}
		}
	}
	if matchCount == len(paramsArr) {
		return true
	}
	return false
}

func (params Parameters) isEqual (otherParams Parameters) bool {
	if (params.Type == otherParams.Type) && (params.Name == otherParams.Name) && (params.Description == otherParams.Description) &&
		(params.Format == otherParams.Format) && (params.In == otherParams.In) && (params.Required == otherParams.Required) {
		return true
	}
	return false
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type IngressList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items []Ingress `json:"items"`
}

