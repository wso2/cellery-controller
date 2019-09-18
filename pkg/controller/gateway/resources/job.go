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
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/config"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//func MakeJob(gateway *v1alpha2.Gateway) *batchv1.Job {
//	return &batchv1.Job{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      JobName(gateway),
//			Namespace: gateway.Namespace,
//			Labels:    makeLabels(gateway),
//			OwnerReferences: []metav1.OwnerReference{
//				*controller.CreateGatewayOwnerRef(gateway),
//			},
//		},
//		Spec: batchv1.JobSpec{
//			Template: corev1.PodTemplateSpec{
//				ObjectMeta: metav1.ObjectMeta{
//					Labels:      makeLabels(gateway),
//					Annotations: makePodAnnotations(gateway),
//				},
//				Spec: makePodSpec(gateway,
//					addPorts(gateway),
//					addConfigMapVolumes(gateway),
//					withRestartPolicy(corev1.RestartPolicyOnFailure),
//				),
//			},
//		},
//	}
//}

func RequireJob(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.IngressExtensions.HasApiPublisher()
}

//func RequireJobUpdate(gateway *v1alpha2.Gateway, job *batchv1.Job) bool {
//	return gateway.Generation != gateway.Status.ObservedGeneration ||
//		job.Generation != gateway.Status.JobGeneration
//}

func StatusFromJob(gateway *v1alpha2.Gateway, job *batchv1.Job) {
	gateway.Status.Type = v1alpha2.GatewayTypeJob
	gateway.Status.AvailableReplicas = job.Status.Active
	gateway.Status.JobGeneration = job.Generation
	// fixme add correct status for jobs
	if job.Status.Active > 0 {
		gateway.Status.Status = v1alpha2.GatewayCurrentStatusReady
	} else {
		gateway.Status.Status = v1alpha2.GatewayCurrentStatusNotReady
	}
}

func MakeJob(gateway *v1alpha2.Gateway, cfg config.Interface) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      JobName(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: batchv1.JobSpec{
			//Selector: makeSelector(gateway),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      makeLabels(gateway),
					Annotations: makePodAnnotations(gateway),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*makeApiPublisherContainer(gateway, cfg),
					},
					Volumes: []corev1.Volume{
						{
							Name: configVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ApiPublisherConfigMap(gateway),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  apiConfigKey,
											Path: apiConfigFile,
										},
										{
											Key:  gatewayConfigKey,
											Path: gatewayConfigFile,
										},
									},
								},
							},
						},
					},
					RestartPolicy:corev1.RestartPolicyOnFailure,
				},
			},
		},
	}
}

func makeApiPublisherContainer(gateway *v1alpha2.Gateway, cfg config.Interface) *corev1.Container {
	return &corev1.Container{
		Name:  "api-publisher",
		Image: "docker.io/madusha7/gateway-job:latest",
		// Ports: []corev1.ContainerPort{
		// 	{
		// 		ContainerPort: tokenServiceContainerJWKSPort,
		// 	},
		// },
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      configVolumeName,
				MountPath: configMountPath,
				ReadOnly:  true,
			},
		},
	}
}
