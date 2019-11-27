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

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
)

type podSpecOption func(*corev1.PodSpec)

func makePodSpec(component *v1alpha2.Component, opt ...podSpecOption) corev1.PodSpec {
	podSpec := component.Spec.Template
	for _, opt := range opt {
		opt(&podSpec)
	}
	return podSpec
}

func addPorts(component *v1alpha2.Component) podSpecOption {
	return func(podSpec *corev1.PodSpec) {
		if len(component.Spec.Template.Containers) == 0 {
			return
		}
		m := make(map[int][]corev1.ContainerPort)

		for _, pm := range component.Spec.Ports {
			if pm.TargetContainer == "" {
				m[0] = append(m[0], corev1.ContainerPort{
					ContainerPort: pm.TargetPort,
				})
				continue
			}
			for i, c := range component.Spec.Template.Containers {
				if pm.TargetContainer == c.Name {
					m[i] = append(m[i], corev1.ContainerPort{
						ContainerPort: pm.TargetPort,
						Protocol:      corev1.ProtocolTCP,
					})
					continue
				}
			}
		}

		for k, v := range m {
			podSpec.Containers[k].Ports = v
		}
	}
}

func addPersistentVolumeClaimVolumes(component *v1alpha2.Component) podSpecOption {
	return func(podSpec *corev1.PodSpec) {
		var volumes []corev1.Volume
		for _, vc := range component.Spec.VolumeClaims {
			volumes = append(volumes, corev1.Volume{
				Name: vc.Template.Name,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: PersistentVolumeClaimName(component, &vc),
					},
				},
			})
		}
		podSpec.Volumes = append(podSpec.Volumes, volumes...)
	}
}

func addConfigMapVolumes(component *v1alpha2.Component) podSpecOption {
	return func(podSpec *corev1.PodSpec) {
		var volumes []corev1.Volume
		for _, c := range component.Spec.Configurations {
			volumes = append(volumes, corev1.Volume{
				Name: c.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: ConfigMapName(component, &c),
						},
					},
				},
			})
		}
		podSpec.Volumes = append(podSpec.Volumes, volumes...)
	}
}

func addSecretVolumes(component *v1alpha2.Component) podSpecOption {
	return func(podSpec *corev1.PodSpec) {
		var volumes []corev1.Volume
		for _, s := range component.Spec.Secrets {
			volumes = append(volumes, corev1.Volume{
				Name: s.Name,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: SecretName(component, &s),
					},
				},
			})
		}
		podSpec.Volumes = append(podSpec.Volumes, volumes...)
	}
}
