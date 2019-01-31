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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/celleryio/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/celleryio/mesh-controller/pkg/controller"
)

func CreateServiceK8sService(service *v1alpha1.Service) *corev1.Service {
	containerPort := defaultServiceContainerPort

	if len(service.Spec.Container.Ports) > 0 {
		containerPort = service.Spec.Container.Ports[0].ContainerPort
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceK8sServiceName(service),
			Namespace: service.Namespace,
			Labels:    createLabels(service),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(service),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name:       controller.HTTPServiceName,
				Protocol:   corev1.ProtocolTCP,
				Port:       service.Spec.ServicePort,
				TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: containerPort},
			}},
			Selector: createLabels(service),
		},
	}
}
