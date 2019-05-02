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
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

func TestCreateEnvoyFilter(t *testing.T) {
	tests := []struct {
		name         string
		tokenService *v1alpha1.TokenService
		want         *v1alpha3.EnvoyFilter
	}{
		{
			name: "foo token service with any intercept mode",
			tokenService: &v1alpha1.TokenService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
					Labels: map[string]string{
						"mesh.cellery.io/cell": "foo-cell",
					},
				},
				Spec: v1alpha1.TokenServiceSpec{
					InterceptMode:  v1alpha1.InterceptModeAny,
					UnsecuredPaths: []string{"/path1", "/path2"},
					OpaPolicies: []v1alpha1.OpaPolicy{
						{
							Key:    "p1",
							Policy: "policy1",
						},
					},
				},
			},
			want: &v1alpha3.EnvoyFilter{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-envoyfilter",
					Labels: map[string]string{
						"mesh.cellery.io/cell": "foo-cell",
						"mesh.cellery.io/sts":  "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "TokenService",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: v1alpha3.EnvoyFilterSpec{
					WorkloadLabels: map[string]string{
						mesh.CellLabelKey: "foo-cell",
					},
					Filters: []v1alpha3.Filter{
						{
							InsertPosition: v1alpha3.InsertPosition{
								Index: "FIRST",
							},
							ListenerMatch: v1alpha3.ListenerMatch{
								ListenerType:     "SIDECAR_INBOUND",
								ListenerProtocol: "HTTP",
							},
							FilterType: "HTTP",
							FilterName: "envoy.ext_authz",
							FilterConfig: v1alpha3.FilterConfig{
								GRPCService: v1alpha3.GRPCService{
									GoogleGRPC: v1alpha3.GoogleGRPC{
										TargetUri:  "foo-service:8080",
										StatPrefix: "ext_authz",
									},
									Timeout: "10s",
								},
							},
						},
						{
							InsertPosition: v1alpha3.InsertPosition{
								Index: "LAST",
							},
							ListenerMatch: v1alpha3.ListenerMatch{
								ListenerType:     "SIDECAR_OUTBOUND",
								ListenerProtocol: "HTTP",
							},
							FilterType: "HTTP",
							FilterName: "envoy.ext_authz",
							FilterConfig: v1alpha3.FilterConfig{
								GRPCService: v1alpha3.GRPCService{
									GoogleGRPC: v1alpha3.GoogleGRPC{
										TargetUri:  "foo-service:8081",
										StatPrefix: "ext_authz",
									},
									Timeout: "10s",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "foo token service with outbound intercept mode",
			tokenService: &v1alpha1.TokenService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
					Labels: map[string]string{
						"mesh.cellery.io/cell": "foo-cell",
					},
				},
				Spec: v1alpha1.TokenServiceSpec{
					InterceptMode:  v1alpha1.InterceptModeOutbound,
					UnsecuredPaths: []string{"/path1", "/path2"},
					OpaPolicies: []v1alpha1.OpaPolicy{
						{
							Key:    "p1",
							Policy: "policy1",
						},
					},
				},
			},
			want: &v1alpha3.EnvoyFilter{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-envoyfilter",
					Labels: map[string]string{
						"mesh.cellery.io/cell": "foo-cell",
						"mesh.cellery.io/sts":  "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "TokenService",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: v1alpha3.EnvoyFilterSpec{
					WorkloadLabels: map[string]string{
						mesh.CellLabelKey: "foo-cell",
					},
					Filters: []v1alpha3.Filter{
						{
							InsertPosition: v1alpha3.InsertPosition{
								Index: "LAST",
							},
							ListenerMatch: v1alpha3.ListenerMatch{
								ListenerType:     "SIDECAR_OUTBOUND",
								ListenerProtocol: "HTTP",
							},
							FilterType: "HTTP",
							FilterName: "envoy.ext_authz",
							FilterConfig: v1alpha3.FilterConfig{
								GRPCService: v1alpha3.GRPCService{
									GoogleGRPC: v1alpha3.GoogleGRPC{
										TargetUri:  "foo-service:8081",
										StatPrefix: "ext_authz",
									},
									Timeout: "10s",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateEnvoyFilter(test.tokenService)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateEnvoyFilter (-want, +got)\n%v", diff)
			}
		})
	}
}
