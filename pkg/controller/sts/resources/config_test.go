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

	"cellery.io/cellery-controller/pkg/apis/mesh"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha1"
	"cellery.io/cellery-controller/pkg/controller/sts/config"
)

func TestCreateTokenServiceConfigMap(t *testing.T) {
	tests := []struct {
		name         string
		tokenService *v1alpha1.TokenService
		config       config.TokenService
		want         *corev1.ConfigMap
	}{
		{
			name: "foo token service without spec",
			tokenService: &v1alpha1.TokenService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
			},
			config: config.TokenService{
				Config: "{my-key:my-value}",
			},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-config",
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
				Data: map[string]string{
					"sts-config":      "{my-key:my-value}",
					"unsecured-paths": "[]",
				},
			},
		},
		{
			name: "foo token service with unsecured path",
			tokenService: &v1alpha1.TokenService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.TokenServiceSpec{
					UnsecuredPaths: []string{"/path1", "/path2"},
				},
			},
			config: config.TokenService{
				Config: "{my-key:my-value}",
			},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-config",
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
				Data: map[string]string{
					"sts-config":      "{my-key:my-value}",
					"unsecured-paths": "[\"/path1\",\"/path2\"]",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateTokenServiceConfigMap(test.tokenService, test.config)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateTokenServiceConfigMap (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestCreateTokenServiceOPAConfigMap(t *testing.T) {
	tests := []struct {
		name         string
		tokenService *v1alpha1.TokenService
		config       config.TokenService
		want         *corev1.ConfigMap
	}{
		{
			name: "foo token service with a policy spec",
			tokenService: &v1alpha1.TokenService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.TokenServiceSpec{
					OpaPolicies: []v1alpha1.OpaPolicy{
						{
							Key:    "policy-key",
							Policy: "policy rego",
						},
					},
				},
			},
			config: config.TokenService{
				Policy: "default policy",
			},
			want: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-policy",
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
				Data: map[string]string{
					"default.rego":    "default policy",
					"policy-key.rego": "policy rego",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateTokenServiceOPAConfigMap(test.tokenService, test.config)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("TestCreateTokenServiceOPAConfigMap (-want, +got)\n%v", diff)
			}
		})
	}
}
