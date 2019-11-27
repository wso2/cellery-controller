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
	"encoding/json"
	"fmt"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/meta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha1"
	"cellery.io/cellery-controller/pkg/controller"
	routing "cellery.io/cellery-controller/pkg/controller/routing"
	listers "cellery.io/cellery-controller/pkg/generated/listers/mesh/v1alpha2"
)

func MakeRoutingVirtualService(cell *v1alpha2.Cell, cellLister listers.CellLister, compositeLister listers.CompositeLister) (*v1alpha3.VirtualService, error) {
	hostNames, httpRoutes, tcpRoutes, err := buildInterCellRoutingInfo(cell, cellLister, compositeLister)
	if err != nil {
		return nil, err
	}
	if len(hostNames) == 0 || (len(httpRoutes) == 0 && len(tcpRoutes) == 0) {
		// No virtual service needed
		return nil, nil
	}
	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      routing.RoutingVirtualServiceName(cell.Name),
			Namespace: cell.Namespace,
			Labels:    makeLabels(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Hosts:    hostNames,
			Gateways: []string{"mesh"},
			Http:     httpRoutes,
			// TCP is not supported atm
			// TODO: support TCP
			//Tcp:      tcpRoutes,
		},
	}, nil
}

func buildInterCellRoutingInfo(cell *v1alpha2.Cell, cellLister listers.CellLister, compositeLister listers.CompositeLister) ([]string, []*v1alpha3.HTTPRoute, []*v1alpha3.TCPRoute, error) {
	var intercellHttpRoutes []*v1alpha3.HTTPRoute
	var intercellTcpRoutes []*v1alpha3.TCPRoute
	var hostNames []string
	// get dependencies from cell annotations,
	dependencies, err := routing.ExtractDependencies(cell.Annotations)
	if err != nil {
		return nil, nil, nil, err
	}
	// if the source cell is a web cell, we need to create a few additional routing rules
	isWebCell := &cell.Spec.Gateway.Spec.Ingress.IngressExtensions != nil && cell.Spec.Gateway.Spec.Ingress.IngressExtensions.ClusterIngress != nil
	// for each dependency, create a route
	for _, dependency := range dependencies {
		dependencyInst := dependency[routing.Instance]
		if dependencyInst == "" {
			return nil, nil, nil, fmt.Errorf("unable to extract dependency instance from annotations")
		}
		dependencyKind := dependency[routing.Kind]
		if dependencyKind == "" {
			return nil, nil, nil, fmt.Errorf("unable to extract dependency kind from annotations")
		}
		if dependencyKind == routing.CellKind {
			depCell, err := cellLister.Cells(cell.Namespace).Get(dependencyInst)
			if err != nil {
				return nil, nil, nil, err
			}
			if len(depCell.Spec.Gateway.Spec.Ingress.HTTPRoutes) > 0 {
				hostNames = append(hostNames, routing.BuildHostNameForCellDependency(dependencyInst))
				// build http routes
				intercellHttpRoutes = append(intercellHttpRoutes, routing.BuildHttpRoutesForCellDependency(cell.Name, dependencyInst, isWebCell, CellSrcLabelBulder{})...)
			}
		} else if dependencyKind == routing.CompositeKind {
			depComposite, err := compositeLister.Composites(cell.Namespace).Get(dependencyInst)
			if err != nil {
				return nil, nil, nil, err
			}
			if len(depComposite.Spec.Components) > 0 {
				hostNames = append(hostNames, routing.BuildHostNamesForCompositeDependency(dependencyInst, depComposite.Spec.Components)...)
				intercellHttpRoutes = append(intercellHttpRoutes, routing.BuildHttpRoutesForCompositeDependency(cell.Name, dependencyInst, depComposite.Spec.Components, isWebCell, CellSrcLabelBulder{})...)
			}
		} else {
			// unknown dependency kind
			return nil, nil, nil, fmt.Errorf("unknown dependency kind '%s'", dependencyKind)
		}
		//if len(depCell.Spec.Gateway.Spec.Ingress.TCPRoutes) > 0 {
		// TCP is not supported atm
		// TODO: support TCP
		// hostNames = append(hostNames, buildHostName(dependencyInst))
		// build tcp routes
		// intercellTcpRoutes = append(intercellTcpRoutes, buildTcpRoutes(depCell, dependencyInst)...)
		//}
	}

	return hostNames, intercellHttpRoutes, intercellTcpRoutes, nil
}

func buildHostName(dependencyInst string) string {
	return GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst))
}

func buildHttpRoutes(cell *v1alpha1.Cell, dependencyInst string) []*v1alpha3.HTTPRoute {
	instanceIdMatch1Rule := &v1alpha3.HTTPRoute{
		Match: []*v1alpha3.HTTPMatchRequest{
			{
				Authority: &v1alpha3.StringMatch{
					Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
				},
				Headers: map[string]*v1alpha3.StringMatch{
					routing.InstanceId: {
						Exact: "1",
					},
				},
				SourceLabels: map[string]string{
					meta.CellLabelKeySource:      cell.Name,
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
	}
	instanceIdMatch2Rule := &v1alpha3.HTTPRoute{
		Match: []*v1alpha3.HTTPMatchRequest{
			{
				Authority: &v1alpha3.StringMatch{
					Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
				},
				Headers: map[string]*v1alpha3.StringMatch{
					routing.InstanceId: {
						Exact: "2",
					},
				},
				SourceLabels: map[string]string{
					meta.CellLabelKeySource:      cell.Name,
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
	}
	percentageBasedRule := &v1alpha3.HTTPRoute{
		Match: []*v1alpha3.HTTPMatchRequest{
			{
				Authority: &v1alpha3.StringMatch{
					Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
				},
				SourceLabels: map[string]string{
					meta.CellLabelKeySource:      cell.Name,
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
	}
	return []*v1alpha3.HTTPRoute{instanceIdMatch1Rule, instanceIdMatch2Rule, percentageBasedRule}
}

func buildTcpRoutes(cell *v1alpha1.Cell, dependencyInst string) []*v1alpha3.TCPRoute {
	var intercellRoutes []*v1alpha3.TCPRoute
	for _, cellTcpRoute := range cell.Spec.GatewayTemplate.Spec.TCPRoutes {
		route := v1alpha3.TCPRoute{
			Match: []*v1alpha3.L4MatchAttributes{
				{
					Port: cellTcpRoute.Port,
					SourceLabels: map[string]string{
						meta.CellLabelKey:      cell.Name,
						meta.ComponentLabelKey: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
						Port: &v1alpha3.PortSelector{
							Number: cellTcpRoute.Port,
						},
					},
				},
			},
		}
		intercellRoutes = append(intercellRoutes, &route)
	}
	return intercellRoutes
}

func extractDependencies(cell *v1alpha1.Cell) ([]map[string]string, error) {
	cellDependencies := cell.Annotations[meta.CellDependenciesAnnotationKey]
	var dependencies []map[string]string
	if cellDependencies == "" {
		// no dependencies
		return dependencies, nil
	}
	err := json.Unmarshal([]byte(cellDependencies), &dependencies)
	if err != nil {
		return dependencies, err
	}
	return dependencies, nil
}

func RequireRoutingVsUpdate(cell *v1alpha2.Cell, vs *v1alpha3.VirtualService) bool {
	return cell.Generation != cell.Status.ObservedGeneration ||
		vs.Generation != cell.Status.RoutingVsGeneration
}

func CopyRoutingVs(source, destination *v1alpha3.VirtualService) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromRoutingVs(cell *v1alpha2.Cell, vs *v1alpha3.VirtualService) {
	cell.Status.RoutingVsGeneration = vs.Generation
}

func BuildVirtualServiceiedConfig(vs *v1alpha3.VirtualService) *v1alpha3.VirtualService {
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

type CellSrcLabelBulder struct{}

func (labelBuilder CellSrcLabelBulder) Get(instance string) map[string]string {
	return map[string]string{
		meta.CellLabelKeySource:      instance,
		meta.ComponentLabelKeySource: "true",
	}
}
