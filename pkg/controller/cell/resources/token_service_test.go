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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

func TestCreateTokenService(t *testing.T) {
	tests := []struct {
		name string
		cell *v1alpha1.Cell
		want *v1alpha1.TokenService
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
			want: &v1alpha1.TokenService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo--sts",
					Labels: map[string]string{
						mesh.CellLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Cell",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateTokenService(test.cell)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateTokenService (-want, +got)\n%v", diff)
			}
		})
	}
}
