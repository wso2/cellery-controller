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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceOption func(*corev1.Service)

func Service(name, namespace string, opt ...ServiceOption) *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range opt {
		opt(s)
	}

	return s
}

func WithServiceLabels(m map[string]string) ServiceOption {
	return func(s *corev1.Service) {
		s.Labels = m
	}
}

func WithServiceOwnerReference(ownerReference metav1.OwnerReference) ServiceOption {
	return func(s *corev1.Service) {
		s.OwnerReferences = append(s.OwnerReferences, ownerReference)
	}
}

func WithServicePort(name string, protocol corev1.Protocol, port int32, targetPort int32) ServiceOption {
	return func(s *corev1.Service) {
		s.Spec.Ports = append(s.Spec.Ports, corev1.ServicePort{
			Name:       name,
			Protocol:   protocol,
			Port:       port,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: targetPort},
		})
	}
}

func WithServiceSelector(m map[string]string) ServiceOption {
	return func(s *corev1.Service) {
		s.Spec.Selector = m
	}
}
