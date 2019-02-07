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
	"testing"

	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/gateway/config"
)

var intOne int32 = 1

func TestCreateGatewayDeployment(t *testing.T) {
	tests := []struct {
		name    string
		gateway *v1alpha1.Gateway
		config  config.Gateway
		want    *appsv1.Deployment
	}{
		{
			name: "foo gateway without spec",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
			},
			config: config.Gateway{},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Gateway",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &intOne,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							mesh.CellGatewayLabelKey: "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								mesh.CellGatewayLabelKey: "foo",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "false",
							},
						},
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name: "cell-gateway-init",
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
									Name: "cell-gateway",
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 8080,
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "CELL_NAME",
											Value: "foo",
										},
									},
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
												Name: "foo-config",
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
												Name: "foo-config",
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
			},
		},
		{
			name: "foo gateway with spec and config",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
					Labels: map[string]string{
						"my-label-key": "my-label-value",
					},
				},
			},
			config: config.Gateway{
				InitConfig:  "",
				SetupConfig: "",
				InitImage:   "vick/init-cell-gateway",
				Image:       "vick/cell-gateway",
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						"my-label-key":           "my-label-value",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Gateway",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &intOne,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							mesh.CellGatewayLabelKey: "foo",
							"my-label-key":           "my-label-value",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								mesh.CellGatewayLabelKey: "foo",
								"my-label-key":           "my-label-value",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "false",
							},
						},
						Spec: corev1.PodSpec{
							InitContainers: []corev1.Container{
								{
									Name:  "cell-gateway-init",
									Image: "vick/init-cell-gateway",
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
									Image: "vick/cell-gateway",
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 8080,
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "CELL_NAME",
											Value: "foo",
										},
									},
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
												Name: "foo-config",
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
												Name: "foo-config",
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
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateGatewayDeployment(test.gateway, test.config)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("TestCreateGatewayDeployment (-want, +got)\n%v", diff)
			}
		})
	}
}
