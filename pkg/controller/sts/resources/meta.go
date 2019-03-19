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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

func createTokenServiceLabels(tokenService *v1alpha1.TokenService) map[string]string {
	labels := make(map[string]string, len(tokenService.ObjectMeta.Labels)+1)
	labels[mesh.CellTokenServiceLabelKey] = tokenService.Name

	for k, v := range tokenService.ObjectMeta.Labels {
		labels[k] = v
	}
	return labels
}

func createTokenServiceSelector(tokenService *v1alpha1.TokenService) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: createTokenServiceLabels(tokenService)}
}

func TokenServiceConfigMapName(tokenService *v1alpha1.TokenService) string {
	return tokenService.Name + "-config"
}

func TokenServicePolicyConfigMapName(tokenService *v1alpha1.TokenService) string {
	return tokenService.Name + "-policy"
}

func TokenServiceDeploymentName(tokenService *v1alpha1.TokenService) string {
	return tokenService.Name + "-deployment"
}

func TokenServiceK8sServiceName(tokenService *v1alpha1.TokenService) string {
	return tokenService.Name + "-service"
}

func EnvoyFilterName(tokenService *v1alpha1.TokenService) string {
	return tokenService.Name + "-envoyfilter"
}
