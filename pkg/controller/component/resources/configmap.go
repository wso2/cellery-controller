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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func MakeConfigMap(component *v1alpha2.Component, configMap *corev1.ConfigMap) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName(component, configMap),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Data:       configMap.Data,
		BinaryData: configMap.BinaryData,
	}
}

func RequireConfigMapUpdate(component *v1alpha2.Component, configMap *corev1.ConfigMap) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		configMap.Generation != component.Status.ConfigMapGenerations[configMap.Name]
}

func CopyConfigMap(source, destination *corev1.ConfigMap) {
	destination.Data = source.Data
	destination.BinaryData = source.BinaryData
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromConfigMap(component *v1alpha2.Component, configMap *corev1.ConfigMap) {
	component.Status.ConfigMapGenerations[configMap.Name] = configMap.Generation
}
