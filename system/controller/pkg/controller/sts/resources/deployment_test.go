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
	"github.com/google/go-cmp/cmp"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	"github.com/wso2/product-vick/system/controller/pkg/controller/sts/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var intOne int32 = 1

func TestCreateTokenServiceDeployment(t *testing.T) {
	tests := []struct {
		name         string
		tokenService *v1alpha1.TokenService
		config       config.TokenService
		want         *appsv1.Deployment
	}{
		{
			name: "foo token service without spec",
			tokenService: &v1alpha1.TokenService{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
					Labels: map[string]string{
						vick.CellLabelKey: "my-cell",
					},
				},
			},
			config: config.TokenService{
				Image: "vick/cell-sts",
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						vick.CellTokenServiceLabelKey: "foo",
						vick.CellLabelKey:             "my-cell",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "TokenService",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &intOne,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							vick.CellTokenServiceLabelKey: "foo",
							vick.CellLabelKey:             "my-cell",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								vick.CellTokenServiceLabelKey: "foo",
								vick.CellLabelKey:             "my-cell",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "false",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "cell-sts",
									Image: "vick/cell-sts",
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 8080,
										},
										{
											ContainerPort: 8081,
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  envCellNameKey,
											Value: "my-cell",
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
												Name: "foo-config",
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
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateTokenServiceDeployment(test.tokenService, test.config)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateTokenServiceDeployment (-want, +got)\n%v", diff)
			}
		})
	}
}
