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
	"testing"

	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	. "github.com/cellery-io/mesh-controller/pkg/testing/apis/apps/v1"
	. "github.com/cellery-io/mesh-controller/pkg/testing/apis/core/v1"
	. "github.com/cellery-io/mesh-controller/pkg/testing/apis/mesh/v1alpha2"
)

func TestMakeDeployment(t *testing.T) {
	tests := []struct {
		name      string
		component *v1alpha2.Component
		want      *appsv1.Deployment
	}{
		{
			name: "foo component without spec",
			component: ComponentWith(createComponent("foo", "foo-namespace"),
				func(component *v1alpha2.Component) {
					component.Spec = v1alpha2.ComponentSpec{}
				},
			),
			want: DeploymentWith(createDeployment("foo-deployment", "foo-namespace"),
				WithDeploymentOwnerReference(OwnerReferenceComponentFunc("foo")),
				WithDeploymentReplicas(1),
				func(deployment *appsv1.Deployment) {
					deployment.Spec.Template.Spec.Containers = nil
					deployment.Spec.Template.Spec.Volumes = nil
				},
			),
		},
		{
			name: "foo component with port mappings",
			component: ComponentWith(createComponent("foo", "foo-namespace"),
				func(component *v1alpha2.Component) {
					component.Spec.Template.Containers[0].VolumeMounts = nil
					component.Spec.Configurations = nil
					component.Spec.Secrets = nil
					component.Spec.VolumeClaims = nil
				},
			),
			want: DeploymentWith(createDeployment("foo-deployment", "foo-namespace"),
				func(deployment *appsv1.Deployment) {
					deployment.Spec.Template.Spec.Containers = nil
				},
				WithDeploymentOwnerReference(OwnerReferenceComponentFunc("foo")),
				WithDeploymentReplicas(1),
				WithDeploymentPodSpec(
					PodSpec(
						WithPodSpecContainer(
							Container(
								WithContainerName("container1"),
								WithContainerImage("busybox1:v1.2.3"),
								WithContainerPort(8080),
								WithContainerPort(9090),
								WithContainerPort(8001),
							),
						),
						WithPodSpecContainer(
							Container(
								WithContainerName("container2"),
								WithContainerImage("nginx:v1.2.3"),
								WithContainerPort(8002),
							),
						),
					),
				),
			),
		},
		{
			name: "foo component with volumes",
			component: ComponentWith(createComponent("foo", "foo-namespace"),
				func(component *v1alpha2.Component) {
					component.Spec.Ports = nil
				},
			),
			want: DeploymentWith(createDeployment("foo-deployment", "foo-namespace"),
				func(deployment *appsv1.Deployment) {
				},
				WithDeploymentOwnerReference(OwnerReferenceComponentFunc("foo")),
				WithDeploymentReplicas(1),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.component.SetDefaults()
			got := MakeDeployment(test.component)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("TestMakeService (-want, +got)\n%v", diff)
			}
		})
	}
}

func createComponent(name, namespace string) *v1alpha2.Component {
	return &v1alpha2.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha2.ComponentSpec{
			Type:     v1alpha2.ComponentTypeDeployment,
			Template: createPodTemplateSpec().Spec,
			Ports: []v1alpha2.PortMapping{
				{
					Name:       "port1",
					Protocol:   v1alpha2.ProtocolHTTP,
					Port:       80,
					TargetPort: 8080,
				},
				{
					Name:       "foo-rpc",
					Protocol:   v1alpha2.ProtocolGRPC,
					Port:       9090,
					TargetPort: 9090,
				},
				{
					Port:       15000,
					TargetPort: 8001,
				},
				{
					Port:            15001,
					TargetContainer: "container2",
					TargetPort:      8002,
				},
			},
			Configurations: []corev1.ConfigMap{
				*ConfigMap("config1", "foo-namespace",
					WithConfigMapData("key1", "value1"),
				),
			},
			Secrets: []corev1.Secret{
				*Secret("secret1", "foo-namespace",
					WithSecretData("key1", []byte("password")),
				),
			},
			VolumeClaims: []v1alpha2.VolumeClaim{
				{
					Shared:   false,
					Template: *PersistentVolumeClaim("pvc1", "foo-namespace"),
				},
			},
		},
	}
}

func createDeployment(name, namespace string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app":                       "foo",
				"mesh.cellery.io/component": "foo",
				"version":                   "v1.0.0",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"app":                       "foo",
				"mesh.cellery.io/component": "foo",
				"version":                   "v1.0.0",
			}},
			Template: *createPodTemplateSpec(),
		},
	}
}

func createPodTemplateSpec() *corev1.PodTemplateSpec {
	return &corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{"app": "foo", "mesh.cellery.io/component": "foo", "version": "v1.0.0"},
			Annotations: map[string]string{"sidecar.istio.io/inject": "true"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "container1",
					Image: "busybox1:v1.2.3",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config1",
							ReadOnly:  true,
							MountPath: "/etc/conf",
						},
						{
							Name:      "secret1",
							ReadOnly:  true,
							MountPath: "/etc/certs",
						},
						{
							Name:      "pvc1",
							ReadOnly:  true,
							MountPath: "/data",
						},
					},
				},
				{
					Name:  "container2",
					Image: "nginx:v1.2.3",
				},
			},
			Volumes: []corev1.Volume{
				*Volume("pvc1", WithVolumeSourcePersistentVolumeClaim("pvc1")),
				*Volume("config1", WithVolumeSourceConfigMap("config1")),
				*Volume("secret1", WithVolumeSourceSecret("secret1")),
			},
		},
	}
}
