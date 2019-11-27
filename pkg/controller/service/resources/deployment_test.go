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

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha1"
	"cellery.io/cellery-controller/pkg/controller"
)

var intOne int32 = 1

func TestCreateServiceDeployment(t *testing.T) {
	tests := []struct {
		name    string
		service *v1alpha1.Service
		want    *appsv1.Deployment
	}{
		{
			name: "foo service without spec",
			service: &v1alpha1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Service",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							appLabelKey:              "foo",
							mesh.CellServiceLabelKey: "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								appLabelKey:                  "foo",
								mesh.CellServiceLabelKey:     "foo",
								mesh.ComponentLabelKey:       "true",
								mesh.ComponentLabelKeySource: "true",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "true",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{},
							},
						},
					},
					Strategy: appsv1.DeploymentStrategy{
						Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
						RollingUpdate: &appsv1.RollingUpdateDeployment{
							MaxSurge: &intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 1,
							},
							MaxUnavailable: &intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 0,
							},
						},
					},
				},
			},
		},
		{
			name: "foo service with spec",
			service: &v1alpha1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo",
				},
				Spec: v1alpha1.ServiceSpec{
					Replicas:           &intOne,
					ServicePort:        80,
					ServiceAccountName: "admin",
					Container: corev1.Container{
						Name:  "foo-container",
						Image: "example.io/app/foo",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 8080,
							},
						},
					},
				},
			},
			want: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo-namespace",
					Name:      "foo-deployment",
					Labels: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
					},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion:         v1alpha1.SchemeGroupVersion.String(),
						Kind:               "Service",
						Name:               "foo",
						Controller:         &boolTrue,
						BlockOwnerDeletion: &boolTrue,
					}},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &intOne,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							appLabelKey:              "foo",
							mesh.CellServiceLabelKey: "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								appLabelKey:                  "foo",
								mesh.CellServiceLabelKey:     "foo",
								mesh.ComponentLabelKey:       "true",
								mesh.ComponentLabelKeySource: "true",
							},
							Annotations: map[string]string{
								controller.IstioSidecarInjectAnnotation: "true",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "foo-container",
									Image: "example.io/app/foo",
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: 8080,
										},
									},
									ReadinessProbe: &corev1.Probe{
										InitialDelaySeconds: readinessInitialDelay,
										TimeoutSeconds:      readinessTimeout,
										PeriodSeconds:       readinessPeriod,
										FailureThreshold:    readinessFailureThreshold,
										SuccessThreshold:    readinessSuccessThreshold,
										Handler: corev1.Handler{
											TCPSocket: &corev1.TCPSocketAction{
												Port: intstr.IntOrString{
													Type:   intstr.Int,
													IntVal: 8080,
												},
											},
										},
									},
								},
							},
						},
					},
					Strategy: appsv1.DeploymentStrategy{
						Type: appsv1.DeploymentStrategyType(appsv1.RollingUpdateDaemonSetStrategyType),
						RollingUpdate: &appsv1.RollingUpdateDeployment{
							MaxSurge: &intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 1,
							},
							MaxUnavailable: &intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 0,
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CreateServiceDeployment(test.service)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CreateServiceDeployment (-want, +got)\n%v", diff)
			}
		})
	}
}

func TestUpdateContainersOfServiceDeployment(t *testing.T) {

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo-namespace",
			Name:      "foo-deployment",
			Labels: map[string]string{
				appLabelKey:              "foo",
				mesh.CellServiceLabelKey: "foo",
			},
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         v1alpha1.SchemeGroupVersion.String(),
				Kind:               "Service",
				Name:               "foo",
				Controller:         &boolTrue,
				BlockOwnerDeletion: &boolTrue,
			}},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &intOne,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					appLabelKey:              "foo",
					mesh.CellServiceLabelKey: "foo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						appLabelKey:              "foo",
						mesh.CellServiceLabelKey: "foo",
					},
					Annotations: map[string]string{
						controller.IstioSidecarInjectAnnotation: "true",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "foo-container",
							Image: "cellerysamples.io/foo:1.0.0",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	svcWithNewContainer := &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo-namespace",
			Name:      "foo",
		},
		Spec: v1alpha1.ServiceSpec{
			Replicas:           &intOne,
			ServicePort:        80,
			ServiceAccountName: "admin",
			Container: corev1.Container{
				Name: "foo-container",
				Env: []corev1.EnvVar{
					{
						Name:  "TestEnv",
						Value: "TestVal",
					},
				},
				Image: "cellerysamples.io/foo:1.0.1",
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: 8080,
					},
				},
			},
		},
	}

	updatedDeployment := UpdateContainersOfServiceDeployment(deployment, svcWithNewContainer)
	expected := svcWithNewContainer.Spec.Container
	actual := updatedDeployment.Spec.Template.Spec.Containers[0]
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("UpdateContainersOfServiceDeployment (-expected, +actual)\n%v", diff)
	}
}
