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

	"cellery.io/cellery-controller/pkg/apis/mesh"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha1"
	"cellery.io/cellery-controller/pkg/controller"
	"cellery.io/cellery-controller/pkg/controller/gateway/config"
)

var intOne int32 = 1

func TestCreateGatewayDeployment(t *testing.T) {
	tests := []struct {
		name    string
		gateway *v1alpha1.Gateway
		config  config.Gateway
		secret  config.Secret
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
			secret: config.Secret{},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
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
							appLabelKey:              "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								mesh.CellGatewayLabelKey: "foo",
								appLabelKey:              "foo",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "false",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "cell-gateway",
									Image: "docker.io/istio/proxyv2:1.2.2",
									Ports: []corev1.ContainerPort{
										{
											Protocol:      corev1.ProtocolTCP,
											ContainerPort: 80,
										},
										{
											ContainerPort: 443,
											Protocol:      corev1.ProtocolTCP,
										},
									},
									Args: []string{
										"proxy",
										"router",
										"--domain",
										"$(POD_NAMESPACE).svc.cluster.local",
										"--drainDuration",
										"45s",
										"--parentShutdownDuration",
										"1m0s",
										"--connectTimeout",
										"10s",
										"--serviceCluster",
										"foo.$(POD_NAMESPACE)",
										"--zipkinAddress",
										"zipkin.istio-system:9411",
										"--proxyAdminPort",
										"15000",
										"--statusPort",
										"15020",
										"--controlPlaneAuthPolicy",
										"NONE",
										"--discoveryAddress",
										"istio-pilot.istio-system:15010",
									},
									Env: []corev1.EnvVar{
										{
											Name:  "CELL_NAME",
											Value: "foo",
										},
										{
											Name: "NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: "v1",
													FieldPath:  "spec.nodeName",
												},
											},
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
											Name: "HOST_IP",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: "v1",
													FieldPath:  "status.hostIP",
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
										{
											Name: "ISTIO_META_CONFIG_NAMESPACE",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: "v1",
													FieldPath:  "metadata.namespace",
												},
											},
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
								{
									Name: "cell-certs",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											SecretName: "foo--secret",
										},
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
				Spec: v1alpha1.GatewaySpec{
					Type: v1alpha1.GatewayTypeMicroGateway,
				},
			},
			config: config.Gateway{
				InitConfig:  "",
				SetupConfig: "",
				InitImage:   "vick/init-cell-gateway",
				Image:       "vick/cell-gateway",
			},
			secret: config.Secret{},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
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
							appLabelKey:              "foo",
							"my-label-key":           "my-label-value",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								mesh.CellGatewayLabelKey: "foo",
								appLabelKey:              "foo",
								"my-label-key":           "my-label-value",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "true",
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
		{
			name: "foo gateway with oidc config",
			gateway: &v1alpha1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.GatewaySpec{
					Type: v1alpha1.GatewayTypeEnvoy,
					OidcConfig: &v1alpha1.OidcConfig{
						ProviderUrl:    "http://provider.com",
						ClientId:       "cid",
						ClientSecret:   "secret",
						BaseUrl:        "http://example.com",
						NonSecurePaths: []string{"/foo1", "/foo2"},
						SubjectClaim:   "claim",
						RedirectUrl:    "http://example.com",
						DcrUser:        "dcr-user",
						DcrPassword:    "dcr-pass",
						DcrUrl:         "http://dcr-url",
						SecurePaths:    []string{"/bar1", "/bar2"},
					},
				},
			},
			config: config.Gateway{
				OidcFilterImage: "oidc-image",
				SkipTlsVerify:   "false",
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						mesh.CellGatewayLabelKey: "foo",
						appLabelKey:              "foo",
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
							appLabelKey:              "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								mesh.CellGatewayLabelKey: "foo",
								appLabelKey:              "foo",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "false",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "cell-gateway",
									Image: "docker.io/istio/proxyv2:1.2.2",
									Ports: []corev1.ContainerPort{
										{
											Protocol:      corev1.ProtocolTCP,
											ContainerPort: 80,
										},
										{
											ContainerPort: 443,
											Protocol:      corev1.ProtocolTCP,
										},
									},
									Args: []string{
										"proxy",
										"router",
										"--domain",
										"$(POD_NAMESPACE).svc.cluster.local",
										"--drainDuration",
										"45s",
										"--parentShutdownDuration",
										"1m0s",
										"--connectTimeout",
										"10s",
										"--serviceCluster",
										"foo.$(POD_NAMESPACE)",
										"--zipkinAddress",
										"zipkin.istio-system:9411",
										"--proxyAdminPort",
										"15000",
										"--statusPort",
										"15020",
										"--controlPlaneAuthPolicy",
										"NONE",
										"--discoveryAddress",
										"istio-pilot.istio-system:15010",
									},
									Env: []corev1.EnvVar{
										{
											Name:  "CELL_NAME",
											Value: "foo",
										},
										{
											Name: "NODE_NAME",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: "v1",
													FieldPath:  "spec.nodeName",
												},
											},
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
											Name: "HOST_IP",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: "v1",
													FieldPath:  "status.hostIP",
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
										{
											Name: "ISTIO_META_CONFIG_NAMESPACE",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													APIVersion: "v1",
													FieldPath:  "metadata.namespace",
												},
											},
										},
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "istio-certs",
											MountPath: "/etc/certs",
										},
									},
								},
								{
									Name:  "envoy-oidc-filter",
									Image: "oidc-image",
									Env: []corev1.EnvVar{
										{
											Name:  "PROVIDER_URL",
											Value: "http://provider.com",
										},
										{
											Name:  "CLIENT_ID",
											Value: "cid",
										},
										{
											Name:  "CLIENT_SECRET",
											Value: "secret",
										},
										{
											Name:  "DCR_ENDPOINT",
											Value: "http://dcr-url",
										},
										{
											Name:  "DCR_USER",
											Value: "dcr-user",
										},
										{
											Name:  "DCR_PASSWORD",
											Value: "dcr-pass",
										},
										{
											Name:  "REDIRECT_URL",
											Value: "http://example.com",
										},
										{
											Name:  "APP_BASE_URL",
											Value: "http://example.com",
										},
										{
											Name:  "PRIVATE_KEY_FILE",
											Value: "/etc/certs/key.pem",
										},
										{
											Name:  "CERTIFICATE_FILE",
											Value: "/etc/certs/cert.pem",
										},
										{
											Name:  "JWT_ISSUER",
											Value: "foo--gateway",
										},
										{
											Name:  "JWT_AUDIENCE",
											Value: "foo",
										},
										{
											Name:  "SUBJECT_CLAIM",
											Value: "claim",
										},
										{
											Name:  "NON_SECURE_PATHS",
											Value: "/foo1,/foo2",
										},
										{
											Name:  "SECURE_PATHS",
											Value: "/bar1,/bar2",
										},
										{
											Name:  "SKIP_DISCOVERY_URL_CERT_VERIFY",
											Value: "false",
										},
									},
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 15800,
											Protocol:      corev1.ProtocolTCP,
										},
										{
											ContainerPort: 15810,
											Protocol:      corev1.ProtocolTCP,
										},
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "cell-certs",
											MountPath: "/etc/certs",
											ReadOnly:  true,
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
								{
									Name: "cell-certs",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											SecretName: "foo--secret",
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
			got, _ := CreateGatewayDeployment(test.gateway, test.config, test.secret)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("TestCreateGatewayDeployment (-want, +got)\n%v", diff)
			}
		})
	}
}
