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

func TestCreateTokenServiceK8sService(t *testing.T) {
	tests := []struct {
		name         string
		tokenService *v1alpha1.TokenService
		want         *corev1.Service
	}{
		{
			name: "foo token service without spec",
			tokenService: &v1alpha1.TokenService{
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
						mesh.CellTokenServiceLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "TokenService",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name:       "grpc-inbound",
							Protocol:   corev1.ProtocolTCP,
							Port:       8080,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
						},
						{
							Name:       "grpc-outbound",
							Protocol:   corev1.ProtocolTCP,
							Port:       8081,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8081},
						},
						{
							Name:       "http-jwks",
							Protocol:   corev1.ProtocolTCP,
							Port:       8090,
							TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8090},
						},
					},
					Selector: map[string]string{
						mesh.CellTokenServiceLabelKey: "foo",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateTokenServiceK8sService(test.tokenService)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateTokenServiceK8sService (-want, +got)\n%v", diff)
			}
		})
	}
}
