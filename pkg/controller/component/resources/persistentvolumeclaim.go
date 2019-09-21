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
	. "github.com/cellery-io/mesh-controller/pkg/meta"
)

func MakePersistentVolumeClaim(component *v1alpha2.Component, volumeClaim *v1alpha2.VolumeClaim) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PersistentVolumeClaimName(component, volumeClaim),
			Namespace: component.Namespace,
			Labels: UnionMaps(
				map[string]string{
					VolumeLabelKey: "pvc",
				},
				makeLabels(component),
				volumeClaim.Template.Labels,
			),
		},
		Spec: volumeClaim.Template.Spec,
	}
}

// func RequirePersistentVolumeClaimUpdate(component *v1alpha2.Component, service *corev1.Service) bool {
// return component.Generation != component.Status.ObservedGeneration ||
// service.Generation != component.Status.ConfigMapGenerations[configMap.Name]
// }

// func CopyPersistentVolumeClaim(source, destination *corev1.PersistentVolumeClaim) {
// destination.Spec.Ports = source.Spec.Ports
// destination.Spec.Selector = source.Spec.Selector
// destination.Labels = source.Labels
// destination.Annotations = source.Annotations
// }

func StatusFromPersistentVolumeClaim(component *v1alpha2.Component, pvc *corev1.PersistentVolumeClaim) {
	component.Status.PersistantVolumeClaimGenerations[pvc.Name] = pvc.Generation
}
