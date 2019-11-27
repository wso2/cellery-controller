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

	"cellery.io/cellery-controller/pkg/meta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
	routing "cellery.io/cellery-controller/pkg/controller/routing"
	listers "cellery.io/cellery-controller/pkg/generated/listers/mesh/v1alpha2"
)

func MakeRoutingVirtualService(composite *v1alpha2.Composite, compositeLister listers.CompositeLister, cellLister listers.CellLister) (*v1alpha3.VirtualService, error) {
	hostNames, httpRoutes, tcpRoutes, err := buildInterCellRoutingInfo(composite, compositeLister, cellLister)
	if err != nil {
		return nil, err
	}
	if len(hostNames) == 0 || (len(httpRoutes) == 0 && len(tcpRoutes) == 0) {
		// No virtual service needed
		return nil, nil
	}
	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      routing.RoutingVirtualServiceName(composite.Name),
			Namespace: composite.Namespace,
			Labels:    makeLabels(composite),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(composite),
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

func buildInterCellRoutingInfo(composite *v1alpha2.Composite, compositeLister listers.CompositeLister, cellLister listers.CellLister) ([]string, []*v1alpha3.HTTPRoute, []*v1alpha3.TCPRoute, error) {
	var intercellHttpRoutes []*v1alpha3.HTTPRoute
	var intercellTcpRoutes []*v1alpha3.TCPRoute
	var hostNames []string
	// get dependencies from cell annotations,
	dependencies, err := routing.ExtractDependencies(composite.Annotations)
	if err != nil {
		return nil, nil, nil, err
	}
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
			depCell, err := cellLister.Cells(composite.Namespace).Get(dependencyInst)
			if err != nil {
				return nil, nil, nil, err
			}
			if len(depCell.Spec.Gateway.Spec.Ingress.HTTPRoutes) > 0 {
				hostNames = append(hostNames, routing.BuildHostNameForCellDependency(dependencyInst))
				// build http routes
				intercellHttpRoutes = append(intercellHttpRoutes, routing.BuildHttpRoutesForCellDependency(composite.Name, dependencyInst, false, CompositeSrcLabelBulder{})...)
			}
		} else if dependencyKind == routing.CompositeKind {
			// retrieve the cell using the cell instance name
			depComposite, err := compositeLister.Composites(composite.Namespace).Get(dependencyInst)
			if err != nil {
				return nil, nil, nil, err
			}
			if len(depComposite.Spec.Components) > 0 {
				hostNames = append(hostNames, routing.BuildHostNamesForCompositeDependency(dependencyInst, depComposite.Spec.Components)...)
				intercellHttpRoutes = append(intercellHttpRoutes, routing.BuildHttpRoutesForCompositeDependency(composite.Name, dependencyInst, depComposite.Spec.Components, false, CompositeSrcLabelBulder{})...)
			}
		} else {
			// unknown dependency kind
			return nil, nil, nil, fmt.Errorf("unknown dependency kind '%s'", dependencyKind)
		}
	}
	return hostNames, intercellHttpRoutes, intercellTcpRoutes, nil
}

func RequireRoutingVsUpdate(composite *v1alpha2.Composite, vs *v1alpha3.VirtualService) bool {
	return composite.Generation != composite.Status.ObservedGeneration ||
		vs.Generation != composite.Status.RoutingVsGeneration
}

func CopyRoutingVs(source, destination *v1alpha3.VirtualService) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromRoutingVs(composite *v1alpha2.Composite, vs *v1alpha3.VirtualService) {
	composite.Status.RoutingVsGeneration = vs.Generation
}

type CompositeSrcLabelBulder struct{}

func (labelBuilder CompositeSrcLabelBulder) Get(instance string) map[string]string {
	return map[string]string{
		meta.CompositeLabelKeySource: instance,
		meta.ComponentLabelKeySource: "true",
	}
}
