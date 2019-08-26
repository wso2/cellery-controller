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

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/ptr"
)

var OwnerReferenceComponentFunc = func(name string) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion:         v1alpha2.SchemeGroupVersion.String(),
		Kind:               "Component",
		Name:               name,
		Controller:         ptr.Bool(true),
		BlockOwnerDeletion: ptr.Bool(true),
	}
}

type ComponentOption func(*v1alpha2.Component)

func Component(name, namespace string, opt ...ComponentOption) *v1alpha2.Component {
	c := &v1alpha2.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	for _, opt := range opt {
		opt(c)
	}

	return c
}

func ComponentWith(component *v1alpha2.Component, opt ...ComponentOption) *v1alpha2.Component {
	for _, opt := range opt {
		opt(component)
	}

	return component
}

func WithComponentType(cType v1alpha2.ComponentType) ComponentOption {
	return func(c *v1alpha2.Component) {
		c.Spec.Type = cType
	}
}

func WithComponentPodSpec(podSpec *corev1.PodSpec) ComponentOption {
	return func(c *v1alpha2.Component) {
		c.Spec.Template = *podSpec
	}
}

func WithComponentPortMaping(name string, protocol v1alpha2.Protocol, port int32, targetContainer string, targetPort int32) ComponentOption {
	return func(c *v1alpha2.Component) {
		c.Spec.Ports = append(c.Spec.Ports, v1alpha2.PortMapping{
			Name:            name,
			Protocol:        protocol,
			Port:            port,
			TargetContainer: targetContainer,
			TargetPort:      targetPort,
		})
	}
}

func WithComponentVolumeClaim(shared bool, pvc *corev1.PersistentVolumeClaim) ComponentOption {
	return func(c *v1alpha2.Component) {
		c.Spec.VolumeClaims = append(c.Spec.VolumeClaims, v1alpha2.VolumeClaim{
			Shared:   shared,
			Template: *pvc,
		})
	}
}

func WithComponentConfiguration(configMap *corev1.ConfigMap) ComponentOption {
	return func(c *v1alpha2.Component) {
		c.Spec.Configurations = append(c.Spec.Configurations, *configMap)
	}
}

func WithComponentSecret(secret *corev1.Secret) ComponentOption {
	return func(c *v1alpha2.Component) {
		c.Spec.Secrets = append(c.Spec.Secrets, *secret)
	}
}
