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

package resources

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

var boolTrue = true

func TestCreateGatewayK8sService(t *testing.T) {
	tests := []struct {
		name    string
		gateway *v1alpha1.Gateway
		want    *corev1.Service
	}{
		{
			name: "foo gateway without spec",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-service",
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "http2",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
						},
						{
							Name:       "https",
							Protocol:   corev1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 443},
						},
					},
					Selector: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
					},
				},
			},
		},
		{
			name: "foo gateway with micro-gateway type",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.GatewaySpec{
					Type: v1alpha1.GatewayTypeMicroGateway,
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-service",
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "http",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
						},
					},
					Selector: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
					},
				},
			},
		},
		{
			name: "foo gateway with envoy-gateway type and oidc config",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.GatewaySpec{
					Type: v1alpha1.GatewayTypeEnvoy,
					OidcConfig: &v1alpha1.OidcConfig{
						ProviderUrl:  "http://provider.com",
						ClientId:     "cid",
						ClientSecret: "secret",
						BaseUrl:      "http://example.com",
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-service",
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "http2",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
						},
						{
							Name:       "https",
							Protocol:   corev1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 443},
						},
						{
							Name:       "http-oidc-callback",
							Protocol:   corev1.ProtocolTCP,
							Port:       15810,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 15810},
						},
					},
					Selector: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
					},
				},
			},
		},
		{
			name: "foo gateway with tcp mode envoy-gateway type",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.GatewaySpec{
					Type: v1alpha1.GatewayTypeEnvoy,
					TCPRoutes: []v1alpha1.TCPRoute{
						{
							Port:        31220,
							BackendHost: "tcp-host",
							BackendPort: 8080,
						},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-service",
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
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "http2",
							Protocol:   corev1.ProtocolTCP,
							Port:       80,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
						},
						{
							Name:       "https",
							Protocol:   corev1.ProtocolTCP,
							Port:       443,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 443},
						},
						{
							Name:       "tcp-31220",
							Protocol:   corev1.ProtocolTCP,
							Port:       31220,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 31220},
						},
					},
					Selector: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateGatewayK8sService(test.gateway)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateGatewayK8sService (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestCreateOriginalGatewayK8sService(t *testing.T) {
	gateway := &v1alpha1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo-namespace",
			Name:      "foo",
		},
		Spec: v1alpha1.GatewaySpec{
			Type: v1alpha1.GatewayTypeMicroGateway,
		},
	}

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo-namespace",
			Name:      "original-foo-gw-svc",
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
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				},
			},
			Selector: map[string]string{
				mesh.CellGatewayLabelKey: "foo",
				appLabelKey:              "foo",
			},
		},
	}
	actual := MakeOriginalGatewayK8sService(gateway, "original-foo-gw-svc")

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("CreateGatewayK8sService (-expected, +actual)\n%v", diff)
	}
}
