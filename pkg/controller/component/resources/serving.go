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
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	servingv1alpha1 "github.com/cellery-io/mesh-controller/pkg/apis/knative/serving/v1alpha1"
	servingv1beta1 "github.com/cellery-io/mesh-controller/pkg/apis/knative/serving/v1beta1"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func MakeServingConfiguration(component *v1alpha2.Component) *servingv1alpha1.Configuration {

	var container corev1.Container

	if len(component.Spec.Template.Containers) > 0 {
		container = component.Spec.Template.Containers[0]
	}
	// Reset the user specified container ports
	// https://github.com/knative/serving/blob/master/docs/runtime-contract.md
	container.Ports = []corev1.ContainerPort{}
	if pm := findServingContainerPortMapping(component, container.Name); pm != nil {
		container.Ports = append(container.Ports, corev1.ContainerPort{
			ContainerPort: pm.TargetPort,
			Name: func() string {
				if pm.Protocol == v1alpha2.ProtocolGRPC {
					return "h2c"
				}
				return "http1"
			}(),
		})
	}

	// Container name is set by the knative controller
	container.Name = ""

	kpaOpts := component.Spec.ScalingPolicy.Kpa

	return &servingv1alpha1.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServingConfigurationName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: servingv1alpha1.ConfigurationSpec{
			Template: &servingv1alpha1.RevisionTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: ServingRevisionName(component),
					Annotations: map[string]string{
						"autoscaling.knative.dev/maxScale": strconv.Itoa(int(kpaOpts.MaxReplicas)),
					},
					// Labels: createLabelsWithComponentFlag(createLabels(component)),
					Labels: makeLabels(component),
				},
				Spec: servingv1alpha1.RevisionSpec{
					RevisionSpec: servingv1beta1.RevisionSpec{
						PodSpec: servingv1beta1.PodSpec{
							Containers: []corev1.Container{
								container,
							},
						},
						ContainerConcurrency: func() servingv1beta1.RevisionContainerConcurrencyType {
							c := kpaOpts.Concurrency
							if c < 1 {
								return 0
							}
							return servingv1beta1.RevisionContainerConcurrencyType(c)
						}(),
					},
				},
			},
		},
	}
}

func MakeServingVirtualService(component *v1alpha2.Component) *v1alpha3.VirtualService {

	var container corev1.Container

	if len(component.Spec.Template.Containers) > 0 {
		container = component.Spec.Template.Containers[0]
	}

	selector := component.Spec.ScalingPolicy.Kpa.Selector

	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServingVirtualServiceName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Gateways: []string{"mesh"},
			Hosts:    []string{ServingRevisionName(component)},
			Http: []*v1alpha3.HTTPRoute{
				{
					AppendHeaders: map[string]string{
						"knative-serving-namespace": "default",
						"knative-serving-revision":  ServingRevisionName(component),
					},
					Match: []*v1alpha3.HTTPMatchRequest{
						{
							Authority: &v1alpha3.StringMatch{
								Regex: fmt.Sprintf("^%s(?::\\d{1,5})?$", ServingRevisionName(component)),
							},
							SourceLabels: selector,
						},
						{
							Authority: &v1alpha3.StringMatch{
								Regex: fmt.Sprintf("^%s\\.default(?::\\d{1,5})?$", ServingRevisionName(component)),
							},
							SourceLabels: selector,
						},
						{
							Authority: &v1alpha3.StringMatch{
								Regex: fmt.Sprintf("^%s\\.default\\.svc\\.cluster\\.local(?::\\d{1,5})?$", ServingRevisionName(component)),
							},
							SourceLabels: selector,
						},
					},
					Route: []*v1alpha3.DestinationWeight{
						{
							Destination: &v1alpha3.Destination{
								Host: ServingRevisionName(component),
								Port: &v1alpha3.PortSelector{
									Number: func() uint32 {
										if pm := findServingContainerPortMapping(component, container.Name); pm != nil {
											if pm.Protocol == v1alpha2.ProtocolGRPC {
												return 81
											}
										}
										return 80
									}(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func RequireKnativeServing(component *v1alpha2.Component) bool {
	return component.Spec.Type == v1alpha2.ComponentTypeDeployment && component.Spec.ScalingPolicy.IsKpa()
}

func RequireServingConfigurationUpdate(component *v1alpha2.Component, configuration *servingv1alpha1.Configuration) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		configuration.Generation != component.Status.ServingConfigurationGeneration
}

func StatusFromServingConfiguration(component *v1alpha2.Component, configuration *servingv1alpha1.Configuration) {
	component.Status.Type = v1alpha2.ComponentTypeDeployment

	component.Status.ServingConfigurationGeneration = configuration.Generation
}

func RequireServingVirtualServiceUpdate(component *v1alpha2.Component, virtualService *v1alpha3.VirtualService) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		virtualService.Generation != component.Status.ServingVirtualServiceGeneration
}

func CopyServingVirtualService(source, destination *v1alpha3.VirtualService) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromServingVirtualService(component *v1alpha2.Component, virtualService *v1alpha3.VirtualService) {
	component.Status.ServingVirtualServiceGeneration = virtualService.Generation
}

func findServingContainerPortMapping(component *v1alpha2.Component, targetContainer string) *v1alpha2.PortMapping {
	for _, pm := range component.Spec.Ports {
		if pm.TargetContainer == "" || pm.TargetContainer == targetContainer {
			return &pm
		}
	}
	return nil
}
