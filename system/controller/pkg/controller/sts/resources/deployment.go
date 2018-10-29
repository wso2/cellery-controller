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
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	"github.com/wso2/product-vick/system/controller/pkg/controller/sts/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTokenServiceDeployment(tokenService *v1alpha1.TokenService, tokenServiceConfig config.TokenService) *appsv1.Deployment {
	podTemplateAnnotations := map[string]string{}
	podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "false"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml
	one := int32(1)
	cellName := tokenService.Labels[vick.CellLabelKey]
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenServiceDeploymentName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    createTokenServiceLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &one,
			Selector: createTokenServiceSelector(tokenService),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      createTokenServiceLabels(tokenService),
					Annotations: podTemplateAnnotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "cell-sts",
							Image: tokenServiceConfig.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: tokenServiceContainerInboundPort,
								},
								{
									ContainerPort: tokenServiceContainerOutboundPort,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  envCellNameKey,
									Value: cellName,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      configVolumeName,
									MountPath: configMountPath,
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: configVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: TokenServiceConfigMapName(tokenService),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  tokenServiceConfigKey,
											Path: tokenServiceConfigFile,
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
