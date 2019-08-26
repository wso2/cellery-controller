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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/meta"
	"github.com/cellery-io/mesh-controller/pkg/ptr"
)

const readinessInitialDelay = 15
const readinessTimeout = 5
const readinessPeriod = 20
const readinessFailureThreshold = 3
const readinessSuccessThreshold = 1

func MakeDeployment(component *v1alpha2.Component) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.Int32(component.Spec.ScalingPolicy.MinReplicas()),
			Selector: makeSelector(component),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      makeLabels(component),
					Annotations: makePodAnnotations(component),
				},
				Spec: makePodSpec(component,
					addPorts(component),
					addPersistentVolumeClaimVolumes(component),
					addConfigMapVolumes(component),
					addSecretVolumes(component),
				),
			},
		},
	}
	meta.AddObjectHash(deployment)
	return deployment
}

func RequireDeployment(component *v1alpha2.Component) bool {
	return component.Spec.Type == v1alpha2.ComponentTypeDeployment && !component.Spec.ScalingPolicy.IsKpa()
}

func RequireDeploymentUpdate(component *v1alpha2.Component, deployment *appsv1.Deployment) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		deployment.Generation != component.Status.DeploymentGeneration
}

func CopyDeployment(source, destination *appsv1.Deployment) {
	destination.Spec.Template = source.Spec.Template
	destination.Spec.Selector = source.Spec.Selector
	destination.Spec.Replicas = source.Spec.Replicas
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromDeployment(component *v1alpha2.Component, deployment *appsv1.Deployment) {
	component.Status.Type = v1alpha2.ComponentTypeDeployment
	component.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	component.Status.DeploymentGeneration = deployment.Generation
	if deployment.Status.AvailableReplicas > 0 {
		component.Status.Status = v1alpha2.ComponentCurrentStatusReady
	} else {
		component.Status.Status = v1alpha2.ComponentCurrentStatusNotReady
	}
}
