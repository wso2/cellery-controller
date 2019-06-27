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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

const readinessInitialDelay = 15
const readinessTimeout = 5
const readinessPeriod = 20
const readinessFailureThreshold = 3
const readinessSuccessThreshold = 1

func CreateServiceDeployment(service *v1alpha1.Service) *appsv1.Deployment {
	podTemplateAnnotations := map[string]string{}
	podTemplateAnnotations[controller.IstioSidecarInjectAnnotation] = "true"
	//https://github.com/istio/istio/blob/master/install/kubernetes/helm/istio/templates/sidecar-injector-configmap.yaml

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceDeploymentName(service),
			Namespace: service.Namespace,
			Labels:    createLabels(service),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(service),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: service.Spec.Replicas,
			Selector: createSelector(service),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      createLabelsWithComponentFlag(createLabels(service)),
					Annotations: podTemplateAnnotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						buildContainerWithReadinessProbe(service),
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
	}
}

func UpdateContainersOfServiceDeployment(deployment *appsv1.Deployment, service *v1alpha1.Service) *appsv1.Deployment {
	newDeployment := deployment.DeepCopy()
	newDeployment.Spec.Template.Spec.Containers = []corev1.Container{
		buildContainerWithReadinessProbe(service),
	}
	return newDeployment
}

func buildContainerWithReadinessProbe(service *v1alpha1.Service) corev1.Container {
	if service.Spec.Container.Name == "" || len(service.Spec.Container.Ports) == 0 || service.Spec.Container.Image == "" {
		return service.Spec.Container
	}
	if service.Spec.Container.ReadinessProbe != nil {
		// readiness probe provided, return
		return service.Spec.Container
	}
	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: readinessInitialDelay,
		TimeoutSeconds:      readinessTimeout,
		PeriodSeconds:       readinessPeriod,
		FailureThreshold:    readinessFailureThreshold,
		SuccessThreshold:    readinessSuccessThreshold,
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: service.Spec.Container.Ports[0].ContainerPort,
				},
			},
		},
	}

	service.Spec.Container.ReadinessProbe = readinessProbe
	return service.Spec.Container
}
