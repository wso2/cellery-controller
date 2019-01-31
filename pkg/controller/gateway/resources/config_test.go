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

	"github.com/celleryio/mesh-controller/pkg/apis/mesh"
	"github.com/celleryio/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/celleryio/mesh-controller/pkg/controller/gateway/config"
)

func TestCreateGatewayConfigMap(t *testing.T) {
	tests := []struct {
		name    string
		gateway *v1alpha1.Gateway
		config  config.Gateway
		want    *corev1.ConfigMap
	}{
		{
			name: "foo gateway without spec",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
			},
			config: config.Gateway{},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-config",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Gateway",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Data: map[string]string{
					apiConfigKey:          `{"cell":"foo","version":"1.0.0","hostname":"foo-service.foo-namespace","apis":null}`,
					gatewayConfigKey:      "",
					gatewaySetupConfigKey: "",
				},
			},
		},
		{
			name: "foo gateway with spec and config",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
					Labels: map[string]string{
						"my-label-key": "my-label-value",
					},
				},
				Spec: v1alpha1.GatewaySpec{
					APIRoutes: []v1alpha1.APIRoute{
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
			config: config.Gateway{
				InitConfig:  "{key:value}",
				SetupConfig: "{key:value}",
				InitImage:   "vick/init-cell-gateway",
				Image:       "vick/cell-gateway",
			},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-config",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						"my-label-key":           "my-label-value",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Gateway",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Data: map[string]string{
					apiConfigKey:          `{"cell":"foo","version":"1.0.0","hostname":"foo-service.foo-namespace","apis":[{"context":"/context-1","definitions":[{"path":"path1","method":"GET"},{"path":"path2","method":"POST"}],"backend":"my-service","global":true}]}`,
					gatewayConfigKey:      "{key:value}",
					gatewaySetupConfigKey: "{key:value}",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CreateGatewayConfigMap(test.gateway, test.config)
			if err != nil {
				t.Errorf("Error while creating the config map: %v", err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateGatewayConfigMap (-want, +got)\n%v", diff)
			}
		})
	}
}
