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
	"fmt"
	"testing"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	servingv1alpha1 "github.com/cellery-io/mesh-controller/pkg/apis/knative/serving/v1alpha1"
	servingv1beta1 "github.com/cellery-io/mesh-controller/pkg/apis/knative/serving/v1beta1"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

var zero int32 = 0

func TestCreateZeroScaleDeployment(t *testing.T) {
	tests := []struct {
		name    string
		service *v1alpha1.Service
		want    *servingv1alpha1.Configuration
	}{
		{
			name: "foo zero scale service with spec",
			service: &v1alpha1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.ServiceSpec{
					Autoscaling: &v1alpha1.AutoscalePolicySpec{
						Policy: v1alpha1.Policy{
							MinReplicas: "0",
							MaxReplicas: 10,
						},
					},
					ServicePort:        80,
					ServiceAccountName: "admin",
					Container: corev1.Container{
						Name:  "foo-container",
						Image: "example.io/app/foo",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 8080,
							},
						},
					},
				},
			},
			want: &servingv1alpha1.Configuration{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-service",
					Labels: map[string]string{
						"app":                     "foo",
						"mesh.cellery.io/service": "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Service",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: servingv1alpha1.ConfigurationSpec{
					Template: &servingv1alpha1.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: "foo-service-rev",
							Annotations: map[string]string{
								"autoscaling.knative.dev/maxScale": "10",
							},
							Labels: map[string]string{
								"app":                        "foo",
								"mesh.cellery.io/service":    "foo",
								mesh.ComponentLabelKey:       "true",
								mesh.ComponentLabelKeySource: "true",
							},
						},
						Spec: servingv1alpha1.RevisionSpec{
							RevisionSpec: servingv1beta1.RevisionSpec{
								PodSpec: servingv1beta1.PodSpec{
									Containers: []corev1.Container{
										{
											Image: "example.io/app/foo",
											Ports: []corev1.ContainerPort{
												{
													ContainerPort: 8080,
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateZeroScaleService(test.service)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateZeroScaleService (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestCreateZeroScaleVirtualService(t *testing.T) {
	tests := []struct {
		name    string
		service *v1alpha1.Service
		want    *v1alpha3.VirtualService
	}{
		{
			name: "foo zero scale service with spec",
			service: &v1alpha1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
					Labels: map[string]string{
						"mesh.cellery.io/cell": "foo-cell",
						"mesh.cellery.io.cell": "foo-cell",
					},
				},
				Spec: v1alpha1.ServiceSpec{
					Autoscaling: &v1alpha1.AutoscalePolicySpec{
						Policy: v1alpha1.Policy{
							MinReplicas: "0",
							MaxReplicas: 10,
						},
					},
					ServicePort:        80,
					ServiceAccountName: "admin",
					Container: corev1.Container{
						Name:  "foo-container",
						Image: "example.io/app/foo",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 8080,
							},
						},
					},
				},
			},
			want: &v1alpha3.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-mesh",
					Labels: map[string]string{
						"app":                     "foo",
						"mesh.cellery.io/service": "foo",
						"mesh.cellery.io/cell":    "foo-cell",
						"mesh.cellery.io.cell":    "foo-cell",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Service",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: v1alpha3.VirtualServiceSpec{
					Gateways: []string{"mesh"},
					Hosts:    []string{"foo-service-rev"},
					Http: []*v1alpha3.HTTPRoute{
						{
							AppendHeaders: map[string]string{
								"knative-serving-namespace": "default",
								"knative-serving-revision":  "foo-service-rev",
							},
							Match: []*v1alpha3.HTTPMatchRequest{
								{
									Authority: &v1alpha3.StringMatch{
										Regex: fmt.Sprintf("^%s(?::\\d{1,5})?$", "foo-service-rev"),
									},
									SourceLabels: map[string]string{
										"mesh.cellery.io.cell": "foo-cell",
									},
								},
								{
									Authority: &v1alpha3.StringMatch{
										Regex: fmt.Sprintf("^%s\\.default(?::\\d{1,5})?$", "foo-service-rev"),
									},
									SourceLabels: map[string]string{
										"mesh.cellery.io.cell": "foo-cell",
									},
								},
								{
									Authority: &v1alpha3.StringMatch{
										Regex: fmt.Sprintf("^%s\\.default\\.svc\\.cluster\\.local(?::\\d{1,5})?$", "foo-service-rev"),
									},
									SourceLabels: map[string]string{
										"mesh.cellery.io.cell": "foo-cell",
									},
								},
							},
							Route: []*v1alpha3.DestinationWeight{
								{
									Destination: &v1alpha3.Destination{
										Host: "foo-service-rev",
										Port: &v1alpha3.PortSelector{
											Number: 80,
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
			got := CreateZeroScaleVirtualService(test.service)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateZeroScaleVirtualService (-want, +got)\n%v", diff)
			}
		})
	}
}
