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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/ptr"
)

func MakeJob(component *v1alpha2.Component) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      JobName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: ptr.Int32(0),
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
					withRestartPolicy(corev1.RestartPolicyOnFailure),
					withSharedProcessNamespace(),
				),
			},
		},
	}
}

func RequireJob(component *v1alpha2.Component) bool {
	return component.Spec.Type == v1alpha2.ComponentTypeJob
}

func RequireJobUpdate(component *v1alpha2.Component, job *batchv1.Job) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		job.Generation != component.Status.JobGeneration
}

func StatusFromJob(component *v1alpha2.Component, job *batchv1.Job) {
	component.Status.Type = v1alpha2.ComponentTypeJob
	component.Status.AvailableReplicas = job.Status.Active
	component.Status.JobGeneration = job.Generation
	// fixme add correct status for jobs
	if job.Status.Active > 0 {
		component.Status.Status = v1alpha2.ComponentCurrentStatusReady
	} else {
		component.Status.Status = v1alpha2.ComponentCurrentStatusNotReady
	}
}
