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

func UnionMaps(maps ...map[string]string) map[string]string {
	size := 0
	for i, _ := range maps {
		size += len(maps[i])
	}

	out := make(map[string]string, size)

	for i, _ := range maps {
		for k, v := range maps[i] {
			out[k] = v
		}
	}

	return out
}
