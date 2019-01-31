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
	"github.com/celleryio/mesh-controller/pkg/apis/mesh"
	"github.com/celleryio/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/celleryio/mesh-controller/pkg/controller"
	"github.com/celleryio/mesh-controller/pkg/controller/gateway/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateGatewayDeployment(gateway *v1alpha1.Gateway, gatewayConfig config.Gateway) *appsv1.Deployment {
	podTemplateAnnotations := map[string]string{}
	podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "true"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml
	var cellName string
	cellName, ok := gateway.Labels[mesh.CellLabelKey]
	if !ok {
		cellName = gateway.Name
	}

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
							Name:  "cell-gateway-init",
							Image: gatewayConfig.InitImage,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      configVolumeName,
									MountPath: configMountPath,
									ReadOnly:  true,
								},
								{
									Name:      setupConfigVolumeName,
									MountPath: setupConfigMountPath,
									ReadOnly:  true,
								},
								{
									Name:      gatewayBuildVolumeName,
									MountPath: gatewayBuildMountPath,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "cell-gateway",
							Image: gatewayConfig.Image,
							Env: []corev1.EnvVar{
								{
									Name:  "CELL_NAME",
									Value: cellName,
								},
							},
							Ports: []corev1.ContainerPort{{
								ContainerPort: gatewayContainerPort,
							}},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      gatewayBuildVolumeName,
									MountPath: gatewayBuildMountPath,
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
										Name: GatewayConfigMapName(gateway),
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
						{
							Name: setupConfigVolumeName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: GatewayConfigMapName(gateway),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  gatewaySetupConfigKey,
											Path: gatewaySetupConfigFile,
										},
									},
								},
							},
						},
						{
							Name: gatewayBuildVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
}

func CreateGatewayDeploymentEnvoy(gateway *v1alpha1.Gateway, gatewayConfig config.Gateway) *appsv1.Deployment {
	podTemplateAnnotations := map[string]string{}
	podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "false"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml
	var cellName string
	cellName, ok := gateway.Labels[mesh.CellLabelKey]
	if !ok {
		cellName = gateway.Name
	}

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
					Containers: []corev1.Container{
						{
							Name:  "cell-gateway",
							Image: "gcr.io/istio-release/proxyv2:1.0.2",
							Env: []corev1.EnvVar{
								{
									Name:  "CELL_NAME",
									Value: cellName,
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.name",
										},
									},
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.namespace",
										},
									},
								},
								{
									Name: "INSTANCE_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "status.podIP",
										},
									},
								},
								{
									Name: "ISTIO_META_POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "metadata.name",
										},
									},
								},
							},
							Args: []string{
								"proxy",
								"router",
								"-v",
								"2",
								"--discoveryRefreshDelay",
								"1s",
								"--drainDuration",
								"45s",
								"--parentShutdownDuration",
								"1m0s",
								"--connectTimeout",
								"10s",
								"--serviceCluster",
								GatewayDeploymentName(gateway),
								"--zipkinAddress",
								"zipkin.istio-system:9411",
								"--statsdUdpAddress",
								"istio-statsd-prom-bridge.istio-system:9125",
								"--proxyAdminPort",
								"15000",
								"--controlPlaneAuthPolicy",
								"NONE",
								"--discoveryAddress",
								"istio-pilot.istio-system:8080",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									ContainerPort: 443,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "istio-certs",
									MountPath: "/etc/certs",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "istio-certs",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "istio.default",
								},
							},
						},
					},
				},
			},
		},
	}
}
