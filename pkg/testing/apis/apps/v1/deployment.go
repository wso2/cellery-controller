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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentOption func(*appsv1.Deployment)

func Deployment(name, namespace string, opt ...DeploymentOption) *appsv1.Deployment {
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range opt {
		opt(d)
	}

	return d
}

func DeploymentWith(deployment *appsv1.Deployment, opt ...DeploymentOption) *appsv1.Deployment {
	for _, opt := range opt {
		opt(deployment)
	}

	return deployment
}

func WithDeploymentLabels(m map[string]string) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Labels = m
	}
}

func WithDeploymentOwnerReference(ownerReference metav1.OwnerReference) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.OwnerReferences = append(d.OwnerReferences, ownerReference)
	}
}

func WithDeploymentReplicas(replicas int32) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Replicas = &replicas
	}
}

func WithDeploymentPodSpec(podSpec *corev1.PodSpec) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Spec = *podSpec
	}
}

func WithDeploymentSelector(m map[string]string) DeploymentOption {
	return func(d *appsv1.Deployment) {
		d.Spec.Selector = &metav1.LabelSelector{MatchLabels: m}
	}
}
