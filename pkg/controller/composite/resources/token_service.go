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

	"cellery.io/cellery-controller/pkg/apis/mesh"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	. "cellery.io/cellery-controller/pkg/meta"
)

func MakeTokenService(composite *v1alpha2.Composite) *v1alpha2.TokenService {
	return &v1alpha2.TokenService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenServiceName(composite),
			Namespace: mesh.SystemNamespace,
			Labels:    makeLabels(composite),
		},
		Spec: v1alpha2.TokenServiceSpec{
			InterceptMode: v1alpha2.InterceptModeAny,
			SecretName:    SecretName(composite),
			InstanceName:  "composite",
			Selector: map[string]string{
				CompositeTokenServiceLabelKey: "true",
			},
		},
	}
}

func StatusFromTokenService(composite *v1alpha2.Composite, tokenService *v1alpha2.TokenService) {
	composite.Status.TokenServiceStatus = tokenService.Status.Status
	composite.Status.TokenServiceGeneration = tokenService.Generation
}
