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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/config"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func MakeConfigMap(tokenService *v1alpha2.TokenService, cfg config.Interface) *corev1.ConfigMap {

	unsecuredPathsStr := "[]"
	if len(tokenService.Spec.UnsecuredPaths) > 0 {
		unsecuredPaths, _ := json.Marshal(tokenService.Spec.UnsecuredPaths)
		unsecuredPathsStr = string(unsecuredPaths)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    makeLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Data: map[string]string{
			tokenServiceConfigKey:   cfg.StringValue(config.ConfigMapKeyTokenServiceConfig),
			unsecuredPathsConfigKey: unsecuredPathsStr,
		},
	}
}

func RequireConfigMapUpdate(tokenService *v1alpha2.TokenService, configMap *corev1.ConfigMap) bool {
	return tokenService.Generation != tokenService.Status.ObservedGeneration ||
		configMap.Generation != tokenService.Status.ConfigMapGeneration
}

func CopyConfigMap(source, destination *corev1.ConfigMap) {
	destination.Data = source.Data
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromConfigMap(tokenService *v1alpha2.TokenService, configMap *corev1.ConfigMap) {
	tokenService.Status.ConfigMapGeneration = configMap.Generation
}

func MakeOpaConfigMap(tokenService *v1alpha2.TokenService, cfg config.Interface) *corev1.ConfigMap {

	m := make(map[string]string)
	m["default.rego"] = cfg.StringValue(config.ConfigMapKeyTokenServiceDefaultOpaPolicy)

	for _, v := range tokenService.Spec.OpaPolicies {
		m[fmt.Sprintf("%s.rego", v.Key)] = v.Policy
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      OpaPolicyConfigMapName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    makeLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Data: m,
	}
}

func RequireOpaConfigMapUpdate(tokenService *v1alpha2.TokenService, configMap *corev1.ConfigMap) bool {
	return tokenService.Generation != tokenService.Status.ObservedGeneration ||
		configMap.Generation != tokenService.Status.OpaConfigMapGeneration
}

func CopyOpaConfigMap(source, destination *corev1.ConfigMap) {
	destination.Data = source.Data
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromOpaConfigMap(tokenService *v1alpha2.TokenService, configMap *corev1.ConfigMap) {
	tokenService.Status.OpaConfigMapGeneration = configMap.Generation
}
