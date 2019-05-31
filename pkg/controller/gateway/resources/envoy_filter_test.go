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

package resources

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func TestCreateEnvoyFilterForEnvoyTCPRoutes(t *testing.T) {
	gateway := &v1alpha1.Gateway{
		Spec: v1alpha1.GatewaySpec{
			Type: v1alpha1.GatewayTypeEnvoy,
			Host: "test.com",
			TCPRoutes: []v1alpha1.TCPRoute{
				{
					Port:        8080,
					BackendHost: "mysql.com",
					BackendPort: 3306,
				},
			},
		},
	}
	filter := CreateEnvoyFilter(gateway)
	expected := createExpectedEnvoyFilter(gateway)

	if diff := cmp.Diff(expected, filter); diff != "" {
		t.Errorf("CreateEnvoyFilter (-expected, +actual)\n%v", diff)
	}
}

func TestCreateEnvoyFilterForEnvoyHttpRoutes(t *testing.T) {
	gateway := &v1alpha1.Gateway{
		Spec: v1alpha1.GatewaySpec{
			Type: v1alpha1.GatewayTypeEnvoy,
			Host: "test.com",
			HTTPRoutes: []v1alpha1.HTTPRoute{
				{
					Authenticate: true,
					Global:       true,
					Backend:      "mytestservice",
					Context:      "hello",
					Definitions: []v1alpha1.APIDefinition{
						{
							Path:   "sayHello",
							Method: "GET",
						},
					},
				},
			},
		},
	}
	filter := CreateEnvoyFilter(gateway)
	expected := createExpectedEnvoyFilter(gateway)

	if diff := cmp.Diff(expected, filter); diff != "" {
		t.Errorf("CreateEnvoyFilter (-expected, +actual)\n%v", diff)
	}
}

func TestCreateEnvoyFilterForEnvoyMicroGateway(t *testing.T) {
	gateway := &v1alpha1.Gateway{
		Spec: v1alpha1.GatewaySpec{
			Type: v1alpha1.GatewayTypeMicroGateway,
			Host: "test.com",
			HTTPRoutes: []v1alpha1.HTTPRoute{
				{
					Authenticate: true,
					Global:       true,
					Backend:      "mytestservice",
					Context:      "hello",
					Definitions: []v1alpha1.APIDefinition{
						{
							Path:   "sayHello",
							Method: "GET",
						},
					},
				},
			},
			Tls: v1alpha1.TlsConfig{
				Secret: "mytlssecret",
			},
		},
	}
	filter := CreateEnvoyFilter(gateway)
	expected := createExpectedEnvoyFilter(gateway)

	if diff := cmp.Diff(expected, filter); diff != "" {
		t.Errorf("CreateEnvoyFilter (-expected, +actual)\n%v", diff)
	}
}

func TestCreateEnvoyFilterForEnvoyWebCell(t *testing.T) {
	gateway := &v1alpha1.Gateway{
		Spec: v1alpha1.GatewaySpec{
			Type: v1alpha1.GatewayTypeMicroGateway,
			Host: "test.com",
			HTTPRoutes: []v1alpha1.HTTPRoute{
				{
					Authenticate: true,
					Global:       true,
					Backend:      "mytestservice",
					Context:      "hello",
					Definitions: []v1alpha1.APIDefinition{
						{
							Path:   "sayHello",
							Method: "GET",
						},
					},
				},
			},
			Tls: v1alpha1.TlsConfig{
				Secret: "mytlssecret",
			},
			OidcConfig: &v1alpha1.OidcConfig{
				ProviderUrl:  "https://accounts.google.com",
				ClientId:     "xxxxxxxxxxxxxxxxxxx",
				ClientSecret: "yyyyyyyyyyyyyyyyyyy",
				RedirectUrl:  "http://pet-store.com/_auth/callback",
				BaseUrl:      "http://pet-store.com/",
				SubjectClaim: "given_name",
				NonSecurePaths: []string{
					"/*",
				},
			},
		},
	}
	filter := CreateEnvoyFilter(gateway)
	expected := createExpectedEnvoyFilter(gateway)

	if diff := cmp.Diff(expected, filter); diff != "" {
		t.Errorf("CreateEnvoyFilter (-expected, +actual)\n%v", diff)
	}
}

func createExpectedEnvoyFilter(gateway *v1alpha1.Gateway) *v1alpha3.EnvoyFilter {
	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      EnvoyFilterName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.EnvoyFilterSpec{
			WorkloadLabels: createGatewayLabels(gateway),
			Filters: []v1alpha3.Filter{
				{
					InsertPosition: v1alpha3.InsertPosition{
						Index: filterInsertPositionFirst,
					},
					ListenerMatch: v1alpha3.ListenerMatch{
						ListenerType:     filterListenerTypeGateway,
						ListenerProtocol: HTTPProtocol,
					},
					FilterName: baseFilterName,
					FilterType: HTTPProtocol,
					FilterConfig: v1alpha3.FilterConfig{
						GRPCService: v1alpha3.GRPCService{
							GoogleGRPC: v1alpha3.GoogleGRPC{
								TargetUri:  "127.0.0.1:15800", // filter is attached as a sidecar
								StatPrefix: statPrefix,
							},
							Timeout: filterTimeout,
						},
					},
				},
			},
		},
	}
}
