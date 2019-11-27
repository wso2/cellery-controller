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

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
)

func MakeService(tokenService *v1alpha2.TokenService) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    makeLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       tokenServiceServicePortGatewayName,
					Protocol:   corev1.ProtocolTCP,
					Port:       tokenServiceServiceGatewayPort,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: tokenServiceContainerGatewayPort},
				},
				{
					Name:       tokenServiceServicePortInboundName,
					Protocol:   corev1.ProtocolTCP,
					Port:       tokenServiceServiceInboundPort,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: tokenServiceContainerInboundPort},
				},
				{
					Name:       tokenServiceServicePortOutboundName,
					Protocol:   corev1.ProtocolTCP,
					Port:       tokenServiceServiceOutboundPort,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: tokenServiceContainerOutboundPort},
				},
				{
					Name:       tokenServiceServicePortJWKSName,
					Protocol:   corev1.ProtocolTCP,
					Port:       tokenServiceContainerJWKSPort,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: tokenServiceContainerJWKSPort},
				},
			},
			Selector: makeLabels(tokenService),
		},
	}
}

func RequireServiceUpdate(tokenService *v1alpha2.TokenService, service *corev1.Service) bool {
	return tokenService.Generation != tokenService.Status.ObservedGeneration ||
		service.Generation != tokenService.Status.ServiceGeneration
}

func CopyService(source, destination *corev1.Service) {
	destination.Spec.Ports = source.Spec.Ports
	destination.Spec.Selector = source.Spec.Selector
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromService(tokenService *v1alpha2.TokenService, service *corev1.Service) {
	tokenService.Status.ServiceGeneration = service.Generation
}
