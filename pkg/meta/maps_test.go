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

package meta

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnionMaps(t *testing.T) {

	tests := []struct {
		name string
		maps []map[string]string
		want map[string]string
	}{{
		name: "nil map",
		want: map[string]string{},
	}, {
		name: "empty map",
		maps: []map[string]string{
			map[string]string{},
		},
		want: map[string]string{},
	}, {
		name: "single map",
		maps: []map[string]string{
			map[string]string{"a": "b", "c": "d"},
		},
		want: map[string]string{"a": "b", "c": "d"},
	}, {
		name: "two maps without shared key",
		maps: []map[string]string{
			map[string]string{"a": "b", "c": "d"},
			map[string]string{"w": "x", "y": "z"},
		},
		want: map[string]string{"a": "b", "c": "d", "w": "x", "y": "z"},
	}, {
		name: "two maps with shared key",
		maps: []map[string]string{
			map[string]string{"a": "b", "c": "d"},
			map[string]string{"c": "f", "y": "z"},
		},
		want: map[string]string{"a": "b", "c": "f", "y": "z"},
	}, {
		name: "three maps with shared key",
		maps: []map[string]string{
			map[string]string{"a": "b", "c": "d"},
			map[string]string{"m": "n", "a": "p"},
			map[string]string{"c": "f", "y": "z"},
		},
		want: map[string]string{"a": "p", "c": "f", "m": "n", "y": "z"},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := UnionMaps(test.maps...)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("MakeLabels (-want, +got) = %v", diff)
			}
		})
	}
}
