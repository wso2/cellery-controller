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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

func CreateTokenService(composite *v1alpha1.Composite) *v1alpha1.TokenService {
	return &v1alpha1.TokenService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenServiceName(composite),
			Namespace: composite.Namespace,
			Labels:    createLabels(composite),
		},
		Spec: v1alpha1.TokenServiceSpec{
			InterceptMode: v1alpha1.InterceptModeAny,
			Composite:     true,
		},
	}
}
