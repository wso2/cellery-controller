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

func TestCreateServiceK8sService(t *testing.T) {
	tests := []struct {
		name    string
		service *v1alpha1.Service
		want    *corev1.Service
	}{
		{
			name: "foo service without spec",
			service: &v1alpha1.Service{
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
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Service",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       8080,
						TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
					}},
					Selector: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
					},
				},
			},
		},
		{
			name: "foo service with spec",
			service: &v1alpha1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
					Labels: map[string]string{
						"my-label-key": "my-label-value",
					},
				},
				Spec: v1alpha1.ServiceSpec{
					ServicePort: 9443,
					Container: corev1.Container{
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 8082,
							},
						},
					},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-service",
					Labels: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
						"my-label-key":           "my-label-value",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Service",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       9443,
						TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8082},
					}},
					Selector: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
						"my-label-key":           "my-label-value",
					},
				},
			},
		},
		{
			name: "foo service with spec without container port",
			service: &v1alpha1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.ServiceSpec{
					ServicePort: 9443,
					Container:   corev1.Container{},
				},
			},
			want: &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-service",
					Labels: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Service",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Name:       "http",
						Protocol:   corev1.ProtocolTCP,
						Port:       9443,
						TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
					}},
					Selector: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateServiceK8sService(test.service)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateServiceK8sService (-want, +got)\n%v", diff)
			}
		})
	}
}
