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
	corev1 "k8s.io/api/core/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	. "cellery.io/cellery-controller/pkg/testing/apis/core/v1"
	. "cellery.io/cellery-controller/pkg/testing/apis/mesh/v1alpha2"
)

func TestMakeService(t *testing.T) {
	tests := []struct {
		name      string
		component *v1alpha2.Component
		want      *corev1.Service
	}{
		{
			name:      "foo component without spec",
			component: Component("foo", "foo-namespace"),
			want: Service("foo-service", "foo-namespace",
				WithServiceLabels(map[string]string{
					"app":                       "foo",
					"mesh.cellery.io/component": "foo",
					"version":                   "v1.0.0",
				}),
				WithServiceOwnerReference(OwnerReferenceComponentFunc("foo")),
				WithServiceSelector(map[string]string{
					"app":                       "foo",
					"mesh.cellery.io/component": "foo",
					"version":                   "v1.0.0",
				}),
			),
		},
		{
			name: "foo component with port mappings",
			component: Component("foo", "foo-namespace",
				WithComponentPortMaping("port1", "HTTP", 80, "", 8080),
				WithComponentPortMaping("foo-rpc", "GRPC", 9090, "", 9090),
				WithComponentPortMaping("", "", 15000, "", 8001),
			),
			want: Service("foo-service", "foo-namespace",
				WithServiceLabels(map[string]string{
					"app":                       "foo",
					"mesh.cellery.io/component": "foo",
					"version":                   "v1.0.0",
				}),
				WithServiceOwnerReference(OwnerReferenceComponentFunc("foo")),
				WithServiceSelector(map[string]string{
					"app":                       "foo",
					"mesh.cellery.io/component": "foo",
					"version":                   "v1.0.0",
				}),
				WithServicePort("http-port1", "TCP", 80, 8080),
				WithServicePort("grpc-foo-rpc", "TCP", 9090, 9090),
				WithServicePort("tcp-15000-8001", "TCP", 15000, 8001),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.component.SetDefaults()
			got := MakeService(test.component)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("TestMakeService (-want, +got)\n%v", diff)
			}
		})
	}
}
