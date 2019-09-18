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

package commons

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/meta"
)

const InstanceId = "x-instance-id"
const Instance = "instance"
const Kind = "kind"
const CompositeKind = "Composite"
const CellKind = "Cell"

func BuildHostNameForCellDependency(dependencyInst string) string {
	return GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst))
}

func BuildHostNamesForCompositeDependency(dependencyInst string, components []v1alpha2.Component) []string {
	var svcNames []string
	for _, component := range components {
		svcNames = append(svcNames, CompositeK8sServiceNameFromInstance(dependencyInst, component))
	}
	return svcNames
}

func BuildHttpRoutesForCellDependency(name string, dependencyInst string, isInstanceIdBasedRulesRequired bool) []*v1alpha3.HTTPRoute {
	var routes []*v1alpha3.HTTPRoute
	if isInstanceIdBasedRulesRequired {
		routes = append(routes, &v1alpha3.HTTPRoute{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
					},
					Headers: map[string]*v1alpha3.StringMatch{
						InstanceId: {
							Exact: "1",
						},
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      name,
						meta.ComponentLabelKeySource: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
					},
				},
			},
		})
		routes = append(routes, &v1alpha3.HTTPRoute{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
					},
					Headers: map[string]*v1alpha3.StringMatch{
						InstanceId: {
							Exact: "2",
						},
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      name,
						meta.ComponentLabelKeySource: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
					},
				},
			},
		})
	}
	routes = append(routes, &v1alpha3.HTTPRoute{
		Match: []*v1alpha3.HTTPMatchRequest{
			{
				Authority: &v1alpha3.StringMatch{
					Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
				},
				SourceLabels: map[string]string{
					meta.CellLabelKeySource:      name,
					meta.ComponentLabelKeySource: "true",
				},
			},
		},
		Route: []*v1alpha3.DestinationWeight{
			{
				Destination: &v1alpha3.Destination{
					Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
				},
			},
		},
	})
	return routes
}

func BuildHttpRoutesForCompositeDependency(name string, dependencyInst string, components []v1alpha2.Component, isInstanceIdBasedRulesRequired bool) []*v1alpha3.HTTPRoute {
	// three virtual services for each
	// TODO: create upon request from SDK side?
	var routes []*v1alpha3.HTTPRoute
	for _, component := range components {
		if isInstanceIdBasedRulesRequired {
			routes = append(routes, &v1alpha3.HTTPRoute{
				Match: []*v1alpha3.HTTPMatchRequest{
					{
						Authority: &v1alpha3.StringMatch{
							Regex: fmt.Sprintf("^(%s)(--%s)(\\S*)$", dependencyInst, component.Name),
						},
						Headers: map[string]*v1alpha3.StringMatch{
							InstanceId: {
								Exact: "1",
							},
						},
						SourceLabels: map[string]string{
							meta.CellLabelKeySource:      name,
							meta.ComponentLabelKeySource: "true",
						},
					},
				},
				Route: []*v1alpha3.DestinationWeight{
					{
						Destination: &v1alpha3.Destination{
							Host: CompositeK8sServiceNameFromInstance(dependencyInst, component),
						},
					},
				},
			})
			routes = append(routes, &v1alpha3.HTTPRoute{
				Match: []*v1alpha3.HTTPMatchRequest{
					{
						Authority: &v1alpha3.StringMatch{
							Regex: fmt.Sprintf("^(%s)(--%s)(\\S*)$", dependencyInst, component.Name),
						},
						Headers: map[string]*v1alpha3.StringMatch{
							InstanceId: {
								Exact: "2",
							},
						},
						SourceLabels: map[string]string{
							meta.CellLabelKeySource:      name,
							meta.ComponentLabelKeySource: "true",
						},
					},
				},
				Route: []*v1alpha3.DestinationWeight{
					{
						Destination: &v1alpha3.Destination{
							Host: CompositeK8sServiceNameFromInstance(dependencyInst, component),
						},
					},
				},
			})
		}
		routes = append(routes, &v1alpha3.HTTPRoute{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--%s)(\\S*)$", dependencyInst, component.Name),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      name,
						meta.ComponentLabelKeySource: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: CompositeK8sServiceNameFromInstance(dependencyInst, component),
					},
				},
			},
		})
	}
	return routes
}

func ExtractDependencies(annotations map[string]string) ([]map[string]string, error) {
	dependencies := annotations[meta.CellDependenciesAnnotationKey]
	var dependencyMap []map[string]string
	if dependencies == "" {
		// no dependencies
		return dependencyMap, nil
	}
	err := json.Unmarshal([]byte(dependencies), &dependencyMap)
	if err != nil {
		return dependencyMap, err
	}
	return dependencyMap, nil
}

func BuildVirtualServiceLastAppliedConfig(vs *v1alpha3.VirtualService) *v1alpha3.VirtualService {
	return &v1alpha3.VirtualService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: v1alpha3.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vs.Name,
			Namespace: vs.Namespace,
		},
		Spec: vs.Spec,
	}
}

func Annotate(vs *v1alpha3.VirtualService, name string, value string) {
	annotations := make(map[string]string, len(vs.ObjectMeta.Annotations)+1)
	annotations[name] = value
	for k, v := range vs.ObjectMeta.Annotations {
		annotations[k] = v
	}
	vs.Annotations = annotations
}

func VirtualServiceName(name string) string {
	return name + "--vs"
}

func CompositeK8sServiceNameFromInstance(instance string, component v1alpha2.Component) string {
	return instance + "--" + component.Name + "-service"
}

func GatewayK8sServiceName(gwName string) string {
	return gwName + "-service"
}

func GatewayNameFromInstanceName(instance string) string {
	return instance + "--gateway"
}
