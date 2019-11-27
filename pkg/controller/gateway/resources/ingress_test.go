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
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"cellery.io/cellery-controller/pkg/apis/mesh"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha1"
)

func TestCreateClusterIngress(t *testing.T) {
	tests := []struct {
		name    string
		gateway *v1alpha1.Gateway
		want    *v1beta1.Ingress
	}{
		{
			name: "foo gateway without spec",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.GatewaySpec{
					Host: "my-host.com",
				},
			},
			want: &v1beta1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-ingress",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Gateway",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
						{
							Host: "my-host.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: v1beta1.IngressBackend{
												ServiceName: "foo-service",
												ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "foo gateway with oidc config",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.GatewaySpec{
					Type: v1alpha1.GatewayTypeMicroGateway,
					Host: "my-host.com",
					OidcConfig: &v1alpha1.OidcConfig{
						ProviderUrl:    "http://provider.com",
						ClientId:       "cid",
						ClientSecret:   "secret",
						BaseUrl:        "http://example.com",
						NonSecurePaths: []string{"/foo1", "/foo2"},
						SubjectClaim:   "claim",
						RedirectUrl:    "http://example.com",
						DcrUser:        "dcr-user",
						DcrPassword:    "dcr-pass",
						DcrUrl:         "http://dcr-url",
						SecurePaths:    []string{"/bar1", "/bar2"},
					},
				},
			},
			want: &v1beta1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-ingress",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Gateway",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: v1beta1.IngressSpec{
					Rules: []v1beta1.IngressRule{
						{
							Host: "my-host.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: v1beta1.IngressBackend{
												ServiceName: "foo-service",
												ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
											},
										},
										{
											Path: "/_auth",
											Backend: v1beta1.IngressBackend{
												ServiceName: "foo-service",
												ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 15810},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "foo gateway with tls secret",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.GatewaySpec{
					Type: v1alpha1.GatewayTypeMicroGateway,
					Host: "my-host.com",
					Tls: v1alpha1.TlsConfig{
						Secret: "my-secret",
					},
				},
			},
			want: &v1beta1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-ingress",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Gateway",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: v1beta1.IngressSpec{
					TLS: []v1beta1.IngressTLS{
						{
							Hosts:      []string{"my-host.com"},
							SecretName: "my-secret",
						},
					},
					Rules: []v1beta1.IngressRule{
						{
							Host: "my-host.com",
							IngressRuleValue: v1beta1.IngressRuleValue{
								HTTP: &v1beta1.HTTPIngressRuleValue{
									Paths: []v1beta1.HTTPIngressPath{
										{
											Path: "/",
											Backend: v1beta1.IngressBackend{
												ServiceName: "foo-service",
												ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
											},
										},
									},
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
			got := CreateClusterIngress(test.gateway)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateGatewayK8sService (-want, +got)\n%v", diff)
			}
		})
	}
}
