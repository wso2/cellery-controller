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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateGatewayK8sService(gateway *v1alpha1.Gateway) *corev1.Service {
	if gateway.Spec.Type == v1alpha1.GatewayTypeMicroGateway {
		return createMicroGatewayK8sService(gateway)
	} else {
		return createEnvoyGatewayK8sService(gateway)
	}
}

func createMicroGatewayK8sService(gateway *v1alpha1.Gateway) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayK8sServiceName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       controller.HTTPServiceName,
					Protocol:   corev1.ProtocolTCP,
					Port:       gatewayServicePort,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: gatewayContainerPort},
				},
			},
			Selector: createGatewayLabels(gateway),
		},
	}
}

func createEnvoyGatewayK8sService(gateway *v1alpha1.Gateway) *corev1.Service {
	var servicePorts []corev1.ServicePort

	servicePorts = append(servicePorts, corev1.ServicePort{
		Name:       "http2",
		Protocol:   corev1.ProtocolTCP,
		Port:       80,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
	})

	servicePorts = append(servicePorts, corev1.ServicePort{
		Name:       "https",
		Protocol:   corev1.ProtocolTCP,
		Port:       443,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 443},
	})

	if gateway.Spec.OidcConfig != nil {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       "http-oidc-callback",
			Protocol:   corev1.ProtocolTCP,
			Port:       15810,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 15810},
		})
	}

	for _, tcpRoute := range gateway.Spec.TCPRoutes {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       fmt.Sprintf("tcp-%d", tcpRoute.Port),
			Protocol:   corev1.ProtocolTCP,
			Port:       int32(tcpRoute.Port),
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: int32(tcpRoute.Port)},
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayK8sServiceName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports:    servicePorts,
			Selector: createGatewayLabels(gateway),
		},
	}
}
