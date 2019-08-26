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
	"crypto/md5"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AddObjectHash(obj metav1.Object) {
	obj.SetAnnotations(
		UnionMaps(obj.GetAnnotations(),
			map[string]string{
				LastAppliedHashAnnotationKey: Hash(obj),
			},
		),
	)
}

func Hash(objs ...interface{}) string {
	bytes := []byte{}
	for _, obj := range objs {
		jsonBytes, _ := json.Marshal(obj)
		bytes = append(bytes, jsonBytes...)
	}
	return fmt.Sprintf("%x", md5.Sum(bytes))
}

func HashEqual(obj1 metav1.Object, obj2 metav1.Object) bool {
	h1, ok := obj1.GetAnnotations()[LastAppliedHashAnnotationKey]
	if !ok {
		return false
	}
	h2, ok := obj2.GetAnnotations()[LastAppliedHashAnnotationKey]
	if !ok {
		return false
	}
	return h1 == h2
}
