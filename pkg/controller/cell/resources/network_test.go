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
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha1"
)

func TestCreateNetworkPolicy(t *testing.T) {
	tests := []struct {
		name string
		cell *v1alpha1.Cell
		want *networkv1.NetworkPolicy
	}{
		{
			name: "foo cell with single service",
			cell: &v1alpha1.Cell{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.CellSpec{
					ServiceTemplates: []v1alpha1.ServiceTemplateSpec{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "bar-service",
							},
						},
					},
				},
			},
			want: &networkv1.NetworkPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo--network",
					Labels: map[string]string{
						mesh.CellLabelKey:       "foo",
						mesh.CellLabelKeySource: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Cell",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: networkv1.NetworkPolicySpec{
					PodSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							mesh.CellLabelKey: "foo",
						},
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      mesh.CellServiceLabelKey,
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"foo--bar-service"},
							},
						},
					},
					PolicyTypes: []networkv1.PolicyType{
						networkv1.PolicyTypeIngress,
					},
					Ingress: []networkv1.NetworkPolicyIngressRule{
						{
							From: []networkv1.NetworkPolicyPeer{
								{
									PodSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											mesh.CellLabelKey:        "foo",
											mesh.CellGatewayLabelKey: "foo--gateway",
										},
									},
								},
								{
									PodSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											mesh.CellLabelKey: "foo",
										},
										MatchExpressions: []metav1.LabelSelectorRequirement{
											{
												Key:      mesh.CellServiceLabelKey,
												Operator: metav1.LabelSelectorOpIn,
												Values:   []string{"foo--bar-service"},
											},
										},
									},
								},
								{
									PodSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											mesh.TelepresenceLabelKey: "telepresence",
										},
									},
								},
								{
									NamespaceSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											"name": "knative-serving",
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
			name: "foo cell with multiple services",
			cell: &v1alpha1.Cell{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.CellSpec{
					GatewayTemplate: v1alpha1.GatewayTemplateSpec{
						Spec: v1alpha1.GatewaySpec{
							HTTPRoutes: []v1alpha1.HTTPRoute{
								{
									Context: "/context-1",
									Backend: "my-service",
									Global:  true,
									Definitions: []v1alpha1.APIDefinition{
										{
											Path:   "path1",
											Method: "GET",
										},
										{
											Path:   "path2",
											Method: "POST",
										},
									},
								},
							},
						},
					},
					ServiceTemplates: []v1alpha1.ServiceTemplateSpec{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "bar-service",
							},
							Spec: v1alpha1.ServiceSpec{
								Replicas:           &intOne,
								ServicePort:        80,
								ServiceAccountName: "admin",
								Container: corev1.Container{
									Name: "bar-container",
									Env: []corev1.EnvVar{
										{
											Name:  "my-key",
											Value: "my-value",
										},
									},
									Image: "sample/bar:v1",
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 8080,
										},
									},
								},
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "baz-service",
							},
							Spec: v1alpha1.ServiceSpec{
								Replicas:    &intOne,
								ServicePort: 80,
								Container: corev1.Container{
									Name:  "baz-container",
									Image: "sample/bar:v1",
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
			want: &networkv1.NetworkPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo--network",
					Labels: map[string]string{
						mesh.CellLabelKey:       "foo",
						mesh.CellLabelKeySource: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Cell",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: networkv1.NetworkPolicySpec{
					PodSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							mesh.CellLabelKey: "foo",
						},
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      mesh.CellServiceLabelKey,
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"foo--bar-service", "foo--baz-service"},
							},
						},
					},
					PolicyTypes: []networkv1.PolicyType{
						networkv1.PolicyTypeIngress,
					},
					Ingress: []networkv1.NetworkPolicyIngressRule{
						{
							From: []networkv1.NetworkPolicyPeer{
								{
									PodSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											mesh.CellLabelKey:        "foo",
											mesh.CellGatewayLabelKey: "foo--gateway",
										},
									},
								},
								{
									PodSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											mesh.CellLabelKey: "foo",
										},
										MatchExpressions: []metav1.LabelSelectorRequirement{
											{
												Key:      mesh.CellServiceLabelKey,
												Operator: metav1.LabelSelectorOpIn,
												Values:   []string{"foo--bar-service", "foo--baz-service"},
											},
										},
									},
								},
								{
									PodSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											mesh.TelepresenceLabelKey: "telepresence",
										},
									},
								},
								{
									NamespaceSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											"name": "knative-serving",
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
			got := CreateNetworkPolicy(test.cell)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateNetworkPolicy (-want, +got)\n%v", diff)
			}
		})
	}
}
