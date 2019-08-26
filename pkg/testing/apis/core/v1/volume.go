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

package v1

import (
	corev1 "k8s.io/api/core/v1"
)

type VolumeOption func(*corev1.Volume)

func Volume(name string, opt ...VolumeOption) *corev1.Volume {
	v := &corev1.Volume{
		Name: name,
	}
	for _, opt := range opt {
		opt(v)
	}
	return v
}

func WithVolumeSourcePersistentVolumeClaim(name string) VolumeOption {
	return func(v *corev1.Volume) {
		v.VolumeSource = corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: name,
			},
		}
	}
}

func WithVolumeSourceConfigMap(name string) VolumeOption {
	return func(v *corev1.Volume) {
		v.VolumeSource = corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
			},
		}
	}
}

func WithVolumeSourceSecret(name string) VolumeOption {
	return func(v *corev1.Volume) {
		v.VolumeSource = corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: name,
			},
		}
	}
}
