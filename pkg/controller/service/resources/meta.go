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

func createLabels(service *v1alpha1.Service) map[string]string {
	labels := make(map[string]string, len(service.ObjectMeta.Labels)+2)

	labels[mesh.CellServiceLabelKey] = service.Name
	labels[appLabelKey] = service.Name
	// order matters
	// todo: update the code if override is not possible
	for k, v := range service.ObjectMeta.Labels {
		labels[k] = v
	}
	return labels
}

func createLabelsWithComponentFlag(labels map[string]string) map[string]string {
	newLabels := make(map[string]string, len(labels)+1)
	newLabels[mesh.ComponentLabelKey] = "true"
	for k, v := range labels {
		newLabels[k] = v
	}
	return newLabels
}

func createSelector(service *v1alpha1.Service) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: createLabels(service)}
}

func ServiceDeploymentName(service *v1alpha1.Service) string {
	return service.Name + "-deployment"
}

func ServiceK8sServiceName(service *v1alpha1.Service) string {
	return service.Name + "-service"
}

func ServiceServingConfigurationName(service *v1alpha1.Service) string {
	return service.Name + "-service"
}

func ServiceServingVirtualServiceName(service *v1alpha1.Service) string {
	return service.Name + "-mesh"
}

func ServiceServingRevisionName(service *v1alpha1.Service) string {
	return ServiceServingConfigurationName(service) + "-rev"
}
