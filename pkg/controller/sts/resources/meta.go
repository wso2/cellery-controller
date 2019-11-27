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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	. "cellery.io/cellery-controller/pkg/meta"
)

func makeLabels(tokenService *v1alpha2.TokenService) map[string]string {
	return UnionMaps(
		tokenService.Labels,
		map[string]string{
			TokenServiceLabelKey: tokenService.Name,
		},
	)
}

func makeSelector(tokenService *v1alpha2.TokenService) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: makeLabels(tokenService)}
}

func makePodAnnotations(tokenService *v1alpha2.TokenService) map[string]string {
	return UnionMaps(
		map[string]string{
			IstioSidecarInjectAnnotationKey: "false",
		},
		tokenService.Labels,
	)
}

func ServiceName(tokenService *v1alpha2.TokenService) string {
	return tokenService.Name + "-service"
}

func DeploymentName(tokenService *v1alpha2.TokenService) string {
	return tokenService.Name + "-deployment"
}

func ConfigMapName(tokenService *v1alpha2.TokenService) string {
	return tokenService.Name + "-config"
}

func OpaPolicyConfigMapName(tokenService *v1alpha2.TokenService) string {
	return tokenService.Name + "-policy"
}

func EnvoyFilterName(tokenService *v1alpha2.TokenService) string {
	return tokenService.Name + "-envoyfilter"
}
