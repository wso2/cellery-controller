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
	"github.com/cellery-io/mesh-controller/pkg/ptr"
)

func MakeStatefulSet(component *v1alpha2.Component) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      StatefulSetName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    ptr.Int32(component.Spec.ScalingPolicy.MinReplicas()),
			ServiceName: ServiceName(component),
			Selector:    makeSelector(component),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      makeLabels(component),
					Annotations: makePodAnnotations(component),
				},
				Spec: makePodSpec(component,
					addPorts(component),
					addConfigMapVolumes(component),
					addSecretVolumes(component),
				),
			},
			VolumeClaimTemplates: makeVolumeClaimTemplates(component),
		},
	}
}

func makeVolumeClaimTemplates(component *v1alpha2.Component) []corev1.PersistentVolumeClaim {
	var pvcs []corev1.PersistentVolumeClaim
	for _, v := range component.Spec.VolumeClaims {
		pvcs = append(pvcs, v.Template)
	}
	return pvcs
}

func RequireStatefulSet(component *v1alpha2.Component) bool {
	return component.Spec.Type == v1alpha2.ComponentTypeStatefulSet
}

func RequireStatefulSetUpdate(component *v1alpha2.Component, statefulSet *appsv1.StatefulSet) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		statefulSet.Generation != component.Status.StatefulSetGeneration
}

func CopyStatefulSet(source, destination *appsv1.StatefulSet) {
	destination.Spec.Template = source.Spec.Template
	destination.Spec.Replicas = source.Spec.Replicas
	destination.Spec.Selector = source.Spec.Selector
	destination.Spec.ServiceName = source.Spec.ServiceName
	destination.Spec.VolumeClaimTemplates = source.Spec.VolumeClaimTemplates
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromStatefulSet(component *v1alpha2.Component, statefulSet *appsv1.StatefulSet) {
	component.Status.Type = v1alpha2.ComponentTypeStatefulSet
	component.Status.AvailableReplicas = statefulSet.Status.ReadyReplicas
	component.Status.StatefulSetGeneration = statefulSet.Generation
	if statefulSet.Status.CurrentReplicas > 0 && statefulSet.Status.ReadyReplicas > 0 {
		component.Status.Status = v1alpha2.ComponentCurrentStatusReady
	} else {
		component.Status.Status = v1alpha2.ComponentCurrentStatusNotReady
	}
}
