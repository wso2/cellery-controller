/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
)

func MakeService(component *v1alpha2.Component) *corev1.Service {

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports:    makeServicePorts(component),
			Selector: makeLabels(component),
		},
	}
}

func makeServicePorts(component *v1alpha2.Component) []corev1.ServicePort {
	var ports []corev1.ServicePort
	for _, p := range component.Spec.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:       fmt.Sprintf("%s-%s", strings.ToLower(string(p.Protocol)), p.Name),
			Protocol:   corev1.ProtocolTCP,
			Port:       p.Port,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: p.TargetPort},
		})
	}
	return ports
}

func RequireService(component *v1alpha2.Component) bool {
	return len(component.Spec.Ports) > 0 &&
		(component.Spec.Type == v1alpha2.ComponentTypeDeployment || component.Spec.Type == v1alpha2.ComponentTypeStatefulSet) &&
		!component.Spec.ScalingPolicy.IsKpa()
}

func RequireServiceUpdate(component *v1alpha2.Component, service *corev1.Service) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		service.Generation != component.Status.ServiceGeneration
}

func CopyService(source, destination *corev1.Service) {
	destination.Spec.Ports = source.Spec.Ports
	destination.Spec.Selector = source.Spec.Selector
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromService(component *v1alpha2.Component, service *corev1.Service) {
	component.Status.ServiceName = service.Name
	component.Status.ServiceGeneration = service.Generation
}
