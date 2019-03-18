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
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/sts/config"
)

func CreateTokenServiceConfigMap(tokenService *v1alpha1.TokenService, tokenServiceConfig config.TokenService) *corev1.ConfigMap {

	unsecuredPathsStr := "[]"
	if len(tokenService.Spec.UnsecuredPaths) > 0 {
		unsecuredPaths, _ := json.Marshal(tokenService.Spec.UnsecuredPaths)
		unsecuredPathsStr = string(unsecuredPaths)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenServiceConfigMapName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    createTokenServiceLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Data: map[string]string{
			tokenServiceConfigKey:   tokenServiceConfig.Config,
			unsecuredPathsConfigKey: unsecuredPathsStr,
		},
	}
}

func CreateTokenServiceOPAConfigMap(tokenService *v1alpha1.TokenService, tokenServiceConfig config.TokenService) *corev1.ConfigMap {

	m := make(map[string]string)
	m["default.rego"] = tokenServiceConfig.Policy

	for _, v := range tokenService.Spec.OpaPolicies {
		m[fmt.Sprintf("%s.rego", v.Key)] = v.Policy
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenServicePolicyConfigMapName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    createTokenServiceLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Data: m,
	}
}
