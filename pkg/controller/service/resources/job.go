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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateServiceJob(service *v1alpha1.Service) *batchv1.Job {
	podTemplateAnnotations := map[string]string{}
	podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "true"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml

	boolTrue := true
	var intZero int32 = 0
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceJobName(service),
			Namespace: service.Namespace,
			Labels:    createLabels(service),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(service),
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &intZero,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      createLabelsWithComponentFlag(createLabels(service)),
					Annotations: podTemplateAnnotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						service.Spec.Container,
					},
					RestartPolicy:         corev1.RestartPolicyNever,
					ShareProcessNamespace: &boolTrue,
				},
			},
		},
	}
}
