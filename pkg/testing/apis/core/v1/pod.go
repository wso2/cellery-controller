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

type PodSpecOption func(*corev1.PodSpec)

func PodSpec(opt ...PodSpecOption) *corev1.PodSpec {
	ps := &corev1.PodSpec{}
	for _, opt := range opt {
		opt(ps)
	}

	return ps
}

func WithPodSpecContainer(container *corev1.Container) PodSpecOption {
	return func(ps *corev1.PodSpec) {
		ps.Containers = append(ps.Containers, *container)
	}
}

func WithPodSpecVolume(volume *corev1.Volume) PodSpecOption {
	return func(ps *corev1.PodSpec) {
		ps.Volumes = append(ps.Volumes, *volume)
	}
}
