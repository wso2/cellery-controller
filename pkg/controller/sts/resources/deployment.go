/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/config"
	"cellery.io/cellery-controller/pkg/controller"
	"cellery.io/cellery-controller/pkg/ptr"
)

func MakeDeployment(tokenService *v1alpha2.TokenService, cfg config.Interface) *appsv1.Deployment {
	// cellName := tokenService.Labels[mesh.CellLabelKey]
	// if tokenService.Spec.Composite {
	// 	cellName = "composite"
	// }
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    makeLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.Int32(1),
			Selector: makeSelector(tokenService),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      makeLabels(tokenService),
					Annotations: makePodAnnotations(tokenService),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						*makeTokenServiceContainer(tokenService, cfg),
						*makeOpaContainer(tokenService, cfg),
						*makeJwksContainer(tokenService, cfg),
					},
					Volumes: []corev1.Volume{
						{
							Name: configVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ConfigMapName(tokenService),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  tokenServiceConfigKey,
											Path: tokenServiceConfigFile,
										},
										{
											Key:  unsecuredPathsConfigKey,
											Path: unsecuredPathsConfigFile,
										},
									},
								},
							},
						},
						{
							Name: policyVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: OpaPolicyConfigMapName(tokenService),
									},
								},
							},
						},
						{
							Name: keyPairVolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									// SecretName: func() string {
									// 	if tokenService.Spec.Composite {
									// 		return "composite-sts-secret"
									// 	}
									// 	return cellName + "--secret"
									// }(),
									SecretName: tokenService.Spec.SecretName,
									Items: []corev1.KeyToPath{
										{
											Key:  "key.pem",
											Path: "key.pem",
										},
										{
											Key:  "cert.pem",
											Path: "cert.pem",
										},
									},
								},
							},
						},
						{
							Name: caCertsVolumeName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									// SecretName: func() string {
									// 	if tokenService.Spec.Composite {
									// 		return "composite-sts-secret"
									// 	}
									// 	return cellName + "--secret"
									// }(),
									SecretName: tokenService.Spec.SecretName,
									Items: []corev1.KeyToPath{
										{
											Key:  "cellery-cert.pem",
											Path: "cellery-cert.pem",
										},
										{
											Key:  "cert-bundle.pem",
											Path: "cert-bundle.pem",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func makeTokenServiceContainer(tokenService *v1alpha2.TokenService, cfg config.Interface) *corev1.Container {

	return &corev1.Container{
		Name:  "sts",
		Image: cfg.StringValue(config.ConfigMapKeyTokenServiceImage),
		// Ports: []corev1.ContainerPort{
		// 	{
		// 		ContainerPort: tokenServiceContainerInboundPort,
		// 	},
		// 	{
		// 		ContainerPort: tokenServiceContainerOutboundPort,
		// 	},
		// },
		Env: []corev1.EnvVar{
			{
				Name:  "CELL_NAME",
				Value: tokenService.Spec.InstanceName,
			},
			{
				Name:  "CELL_NAMESPACE",
				Value: tokenService.Namespace,
			},
			{
				Name:  "VALIDATE_SERVER_CERT",
				Value: "true",
			},
			{
				Name:  "ENABLE_HOSTNAME_VERIFICATION",
				Value: "true",
			},
			{
				Name:  "CELL_IMAGE_NAME",
				Value: tokenService.Annotations["mesh.cellery.io/cell-image-name"],
			},
			{
				Name:  "CELL_IMAGE_VERSION",
				Value: tokenService.Annotations["mesh.cellery.io/cell-image-version"],
			},
			{
				Name:  "CELL_INSTANCE_NAME",
				Value: tokenService.Spec.InstanceName,
			},
			{
				Name:  "CELL_ORG_NAME",
				Value: tokenService.Annotations["mesh.cellery.io/cell-image-org"],
			},
		},
		ReadinessProbe: &corev1.Probe{
			Handler: corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{Type: intstr.Int, IntVal: tokenServiceContainerGatewayPort},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      configVolumeName,
				MountPath: configMountPath,
				ReadOnly:  true,
			},
			{
				Name:      policyVolumeName,
				MountPath: pocliyConfigMountPath,
				ReadOnly:  true,
			},
			{
				Name:      keyPairVolumeName,
				MountPath: keyPairMountPath,
				ReadOnly:  true,
			},
			{
				Name:      caCertsVolumeName,
				MountPath: caCertsMountPath,
				ReadOnly:  true,
			},
		},
	}
}

func makeOpaContainer(tokenService *v1alpha2.TokenService, cfg config.Interface) *corev1.Container {
	return &corev1.Container{
		Name:  "opa",
		Image: cfg.StringValue(config.ConfigMapKeyTokenServiceOpaImage),
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: opaServicePort,
				Name:          "http",
			},
		},
		Args: []string{
			"run",
			"--ignore=.*",
			"--server",
			"--watch",
			"/policies",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      policyVolumeName,
				MountPath: pocliyConfigMountPath,
				ReadOnly:  true,
			},
		},
	}
}

func makeJwksContainer(tokenService *v1alpha2.TokenService, cfg config.Interface) *corev1.Container {
	return &corev1.Container{
		Name:  "jwks-server",
		Image: cfg.StringValue(config.ConfigMapKeyTokenServiceJwksImage),
		Env: []corev1.EnvVar{
			{
				Name:  "jwksPort",
				Value: strconv.Itoa(tokenServiceContainerJWKSPort),
			},
		},
		// Ports: []corev1.ContainerPort{
		// 	{
		// 		ContainerPort: tokenServiceContainerJWKSPort,
		// 	},
		// },
		ReadinessProbe: &corev1.Probe{
			Handler: corev1.Handler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{Type: intstr.Int, IntVal: tokenServiceContainerJWKSPort},
				},
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      keyPairVolumeName,
				MountPath: keyPairMountPath,
				ReadOnly:  true,
			},
			{
				Name:      caCertsVolumeName,
				MountPath: caCertsMountPath,
				ReadOnly:  true,
			},
		},
	}
}

func RequireDeploymentUpdate(tokenService *v1alpha2.TokenService, deployment *appsv1.Deployment) bool {
	return tokenService.Generation != tokenService.Status.ObservedGeneration ||
		deployment.Generation != tokenService.Status.DeploymentGeneration
}

func CopyDeployment(source, destination *appsv1.Deployment) {
	destination.Spec.Template = source.Spec.Template
	destination.Spec.Selector = source.Spec.Selector
	destination.Spec.Replicas = source.Spec.Replicas
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromDeployment(tokenService *v1alpha2.TokenService, deployment *appsv1.Deployment) {
	tokenService.Status.DeploymentGeneration = deployment.Generation
	if deployment.Status.AvailableReplicas > 0 {
		tokenService.Status.Status = v1alpha2.TokenServiceCurrentStatusReady
	} else {
		tokenService.Status.Status = v1alpha2.TokenServiceCurrentStatusNotReady
	}
}
