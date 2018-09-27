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
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	apiConfigVolumeName     = "api-config-volume"
	gatewayConfigVolumeName = "gateway-config-volume"
	gatewayConfigKey        = "cell-gateway-init-config"
	gatewayConfigFile       = "gw.json"
	configMountPath         = "/etc/config"
	apiConfigFile           = "api.json"
)

func CreateGatewayDeployment(gateway *v1alpha1.Gateway) *appsv1.Deployment {
	podTemplateAnnotations := map[string]string{}
	podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "false"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml
	one := int32(1)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayDeploymentName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &one,
			Selector: createGatewaySelector(gateway),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      createGatewayLabels(gateway),
					Annotations: podTemplateAnnotations,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "init-cell-gateway",
							Image: "nipunaprashan/microgateway_init",
							Args: []string{
								"cell", "cell", "admin", "admin", "https://10.100.1.217:9443/",
								"lib/platform/bre/security/ballerinaTruststore.p12", "ballerina",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      apiConfigVolumeName,
									MountPath: configMountPath + "/" + apiConfigFile,
									SubPath:   apiConfigFile,
									ReadOnly:  true,
								},
								{
									Name:      gatewayConfigVolumeName,
									MountPath: configMountPath + "/" + gatewayConfigFile,
									SubPath:   gatewayConfigFile,
									ReadOnly:  true,
								},
								{
									Name:      "targetdir",
									MountPath: "/target",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "cell-gateway",
							Image: "nipunaprashan/microgateway",
							Ports: []corev1.ContainerPort{{
								ContainerPort: gatewayContainerPort,
							}},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "targetdir",
									MountPath: "/target",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: apiConfigVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: GatewayConfigMapName(gateway),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  apiConfigKey,
											Path: apiConfigFile,
										},
									},
								},
							},
						},
						{
							Name: gatewayConfigVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: controller.VICKConfigMapName,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  gatewayConfigKey,
											Path: gatewayConfigFile,
										},
									},
								},
							},
						},
						{
							Name: "targetdir",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{
								},
							},
						},
					},
				},
			},
		},
	}
}
