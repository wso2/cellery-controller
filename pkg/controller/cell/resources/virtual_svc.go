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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	listers "github.com/cellery-io/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

const instance = "instance"

func CreateCellVirtualService(cell *v1alpha1.Cell, cellLister listers.CellLister) (*v1alpha3.VirtualService, error) {
	// http
	httpHostNames, httpRoutes, err := buildInterCellHttpRoutingInfo(cell, cellLister)
	if err != nil {
		return nil, err
	}
	// tcp
	tcpHostNames, tcpRoutes, err := buildInterCellTcpRoutes(cell, cellLister)
	if err != nil {
		return nil, err
	}
	if len(httpHostNames) == 0 && len(tcpHostNames) == 0 {
		// implies no dependencies, hence no need of a virtual service
		return nil, nil
	}
	if len(httpRoutes) == 0 && len(tcpRoutes) == 0 {
		// implies no dependencies, hence no need of a virtual service
		return nil, nil
	}
	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CellVirtualServiceName(cell),
			Namespace: cell.Namespace,
			Labels:    createLabels(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: v1alpha3.VirtualServiceSpec{
			// atm only HTTP is considered
			// TODO: fix this to work with TCP
			//Hosts:    append(httpHostNames, tcpHostNames...),
			Hosts:    httpHostNames,
			Gateways: []string{"mesh"},
			Http:     httpRoutes,
			// atm if there are any tcp routes exposed from dependency instance are not used
			// TODO: fix this to work with TCP
			//Tcp:      tcpRoutes,
		},
	}, nil
}

func buildInterCellHttpRoutingInfo(cell *v1alpha1.Cell, cellLister listers.CellLister) ([]string, []*v1alpha3.HTTPRoute, error) {
	var intercellRoutes []*v1alpha3.HTTPRoute
	var httpHostNames []string
	// get dependencies from cell annotations
	dependencies, err := extractDependencies(cell)
	if err != nil {
		return nil, nil, err
	}
	// for each dependency, create a route
	for _, dependency := range dependencies {
		dependencyInst := dependency[instance]
		if dependencyInst == "" {
			return nil, nil, fmt.Errorf("unable to extract dependency instance from annotations")
		}
		// retrieve the cell using the cell instance name
		depCell, err := cellLister.Cells(cell.Namespace).Get(dependencyInst)
		if err != nil {
			return nil, nil, err
		}
		// if there are no HTTP components exposed from gateway, return
		if len(depCell.Spec.GatewayTemplate.Spec.HTTPRoutes) == 0 {
			return httpHostNames, intercellRoutes, nil
		}
		// build http host names
		gwSvcName := GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst))
		httpHostNames = append(httpHostNames, gwSvcName)

		route := v1alpha3.HTTPRoute{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
					},
					SourceLabels: map[string]string{
						mesh.CellLabelKey:      cell.Name,
						mesh.ComponentLabelKey: "true",
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
		intercellRoutes = append(intercellRoutes, &route)
	}

	return httpHostNames, intercellRoutes, nil
}

func buildInterCellTcpRoutes(cell *v1alpha1.Cell, cellLister listers.CellLister) ([]string, []*v1alpha3.TCPRoute, error) {
	var intercellRoutes []*v1alpha3.TCPRoute
	var tcpHostNames []string
	// get dependencies from cell annotations
	dependencies, err := extractDependencies(cell)
	if err != nil {
		return nil, nil, err
	}
	// for each dependency, create a TCP route if the dependency is exposing TCP services from the gateway.
	for _, dependency := range dependencies {
		dependencyInst := dependency[instance]
		if dependencyInst == "" {
			return nil, nil, fmt.Errorf("unable to extract dependency instance from annotations")
		}
		// retrieve the cell using the cell instance name
		depCell, err := cellLister.Cells(cell.Namespace).Get(dependencyInst)
		if err != nil {
			return nil, nil, err
		}
		gwSvcName := GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst))
		tcpHostNames = append(tcpHostNames, gwSvcName)
		// for each TCP service, create a rule
		for _, tcpService := range *getTcpServices(depCell) {
			route := v1alpha3.TCPRoute{
				Match: []*v1alpha3.L4MatchAttributes{
					{
						Port: tcpService.Port,
						SourceLabels: map[string]string{
							mesh.CellLabelKey:      cell.Name,
							mesh.ComponentLabelKey: "true",
						},
					},
				},
				Route: []*v1alpha3.DestinationWeight{
					{
						Destination: &v1alpha3.Destination{
							Host: GatewayNameFromInstanceName(dependencyInst),
							Port: &v1alpha3.PortSelector{
								Number: tcpService.Port,
							},
						},
					},
				},
			}
			intercellRoutes = append(intercellRoutes, &route)
		}
	}

	return tcpHostNames, intercellRoutes, nil
}

func getTcpServices(cell *v1alpha1.Cell) *[]v1alpha1.TCPRoute {
	return &cell.Spec.GatewayTemplate.Spec.TCPRoutes
}

func extractDependencies(cell *v1alpha1.Cell) ([]map[string]string, error) {
	cellDependencies := cell.Annotations[mesh.CellDependenciesAnnotationKey]
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
