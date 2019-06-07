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
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateZeroScaleService(service *v1alpha1.Service) *servingv1alpha1.Configuration {
	// Container name is set by the knative controller
	service.Spec.Container.Name = ""
	return &servingv1alpha1.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceServingConfigurationName(service),
			Namespace: service.Namespace,
			Labels:    createLabels(service),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(service),
			},
		},
		Spec: servingv1alpha1.ConfigurationSpec{
			Template: &servingv1alpha1.RevisionTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: ServiceServingRevisionName(service),
					Annotations: map[string]string{
						"autoscaling.knative.dev/maxScale": func() string {
							if service.Spec.Autoscaling == nil {
								return "0"
							}
							return strconv.Itoa(int(service.Spec.Autoscaling.Policy.MaxReplicas))
						}(),
					},
					Labels: createLabels(service),
				},
				Spec: servingv1alpha1.RevisionSpec{
					RevisionSpec: servingv1beta1.RevisionSpec{
						PodSpec: servingv1beta1.PodSpec{
							Containers: []corev1.Container{
								service.Spec.Container,
							},
						},
					},
				},
			},
		},
	}
}

func CreateZeroScaleVirtualService(service *v1alpha1.Service) *v1alpha3.VirtualService {
	// Container name is set by the knative controller
	service.Spec.Container.Name = ""
	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceServingVirtualServiceName(service),
			Namespace: service.Namespace,
			Labels:    createLabels(service),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateServiceOwnerRef(service),
			},
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Gateways: []string{"mesh"},
			Hosts:    []string{ServiceServingRevisionName(service)},
			Http: []*v1alpha3.HTTPRoute{
				{
					AppendHeaders: map[string]string{
						"knative-serving-namespace": "default",
						"knative-serving-revision":  ServiceServingRevisionName(service),
					},
					Match: []*v1alpha3.HTTPMatchRequest{
						{
							Authority: &v1alpha3.StringMatch{
								Regex: fmt.Sprintf("^%s(?::\\d{1,5})?$", ServiceServingRevisionName(service)),
							},
							SourceLabels: map[string]string{
								mesh.CellLabelKey: service.Labels[mesh.CellLabelKey],
							},
						},
						{
							Authority: &v1alpha3.StringMatch{
								Regex: fmt.Sprintf("^%s\\.default(?::\\d{1,5})?$", ServiceServingRevisionName(service)),
							},
							SourceLabels: map[string]string{
								mesh.CellLabelKey: service.Labels[mesh.CellLabelKey],
							},
						},
						{
							Authority: &v1alpha3.StringMatch{
								Regex: fmt.Sprintf("^%s\\.default\\.svc\\.cluster\\.local(?::\\d{1,5})?$", ServiceServingRevisionName(service)),
							},
							SourceLabels: map[string]string{
								mesh.CellLabelKey: service.Labels[mesh.CellLabelKey],
							},
						},
					},
					Route: []*v1alpha3.DestinationWeight{
						{
							Destination: &v1alpha3.Destination{
								Host: ServiceServingRevisionName(service),
							},
						},
					},
				},
			},
		},
	}
}
