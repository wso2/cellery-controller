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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func TestCreateOriginalComponentK8sService(t *testing.T) {
	composite := &v1alpha1.Composite{
		Spec: v1alpha1.CompositeSpec{
			ServiceTemplates: []v1alpha1.ServiceTemplateSpec{
				{
					Spec: v1alpha1.ServiceSpec{
						Type:        "Deployment",
						ServicePort: 8080,
						Replicas:    &intOne,
						Protocol:    "http",
						Container: corev1.Container{
							Name: "coNtaiNeR",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							Image: "izza/testcontainer",
						},
					},
				},
			},
		},
	}

	compName := "myhr--mycomponent"
	targetPorts := []int{5005}

	expected := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      K8sServiceName(compName),
			Namespace: composite.Namespace,
			Labels:    createLabelsForCurrentPodsWithPrevSvcName(composite, compName),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCompositeOwnerRef(composite),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports:    *getAsContainerPorts(targetPorts),
			Selector: createLabelsForCurrentPodsWithPrevSvcName(composite, compName),
		},
	}
	actual := CreateOriginalComponentK8sService(composite, compName, targetPorts)

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("CreateOriginalComponentK8sService (-want, +got)\n%v", diff)
	}
}
