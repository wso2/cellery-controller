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

//
//import (
//	"fmt"
//
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//
//	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
//	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
//	listers "github.com/cellery-io/mesh-controller/pkg/generated/listers/mesh/v1alpha1"
//	"github.com/cellery-io/mesh-controller/pkg/controller"
//	controllercommons "github.com/cellery-io/mesh-controller/pkg/controller/commons"
//)
//
//func CreateCellVirtualService(composite *v1alpha1.Composite, compositeLister listers.CompositeLister, cellLister listers.CellLister) (*v1alpha3.VirtualService, error) {
//	hostNames, httpRoutes, tcpRoutes, err := buildInterCellRoutingInfo(composite, compositeLister, cellLister)
//	if err != nil {
//		return nil, err
//	}
//	if len(hostNames) == 0 || (len(httpRoutes) == 0 && len(tcpRoutes) == 0) {
//		// No virtual service needed
//		return nil, nil
//	}
//	return &v1alpha3.VirtualService{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      controllercommons.VirtualServiceName(composite.Name),
//			Namespace: composite.Namespace,
//			Labels:    createLabels(composite),
//			OwnerReferences: []metav1.OwnerReference{
//				*controller.CreateCellOwnerRef(composite),
//			},
//		},
//		Spec: v1alpha3.VirtualServiceSpec{
//			Hosts:    hostNames,
//			Gateways: []string{"mesh"},
//			Http:     httpRoutes,
//			// TCP is not supported atm
//			// TODO: support TCP
//			//Tcp:      tcpRoutes,
//		},
//	}, nil
//}
//
//func buildInterCellRoutingInfo(composite *v1alpha1.Composite, compositeLister listers.CompositeLister, cellLister listers.CellLister) ([]string, []*v1alpha3.HTTPRoute, []*v1alpha3.TCPRoute, error) {
//	var intercellHttpRoutes []*v1alpha3.HTTPRoute
//	var intercellTcpRoutes []*v1alpha3.TCPRoute
//	var hostNames []string
//	// get dependencies from cell annotations,
//	dependencies, err := controllercommons.ExtractDependencies(composite.Annotations)
//	if err != nil {
//		return nil, nil, nil, err
//	}
//	// for each dependency, create a route
//	for _, dependency := range dependencies {
//		dependencyInst := dependency[controllercommons.Instance]
//		if dependencyInst == "" {
//			return nil, nil, nil, fmt.Errorf("unable to extract dependency instance from annotations")
//		}
//		dependencyKind := dependency[controllercommons.Kind]
//		if dependencyKind == "" {
//			return nil, nil, nil, fmt.Errorf("unable to extract dependency kind from annotations")
//		}
//		if dependencyKind == controllercommons.CellKind {
//			depCell, err := cellLister.Cells(composite.Namespace).Get(dependencyInst)
//			if err != nil {
//				return nil, nil, nil, err
//			}
//			if len(depCell.Spec.GatewayTemplate.Spec.HTTPRoutes) > 0 {
//				hostNames = append(hostNames, controllercommons.BuildHostNameForCellDependency(dependencyInst))
//				// build http routes
//				intercellHttpRoutes = append(intercellHttpRoutes, controllercommons.BuildHttpRoutesForCellDependency(composite.Name, dependencyInst, false)...)
//			}
//		} else if dependencyKind == controllercommons.CompositeKind {
//			// retrieve the cell using the cell instance name
//			depComposite, err := compositeLister.Composites(composite.Namespace).Get(dependencyInst)
//			if err != nil {
//				return nil, nil, nil, err
//			}
//			if len(depComposite.Spec.ServiceTemplates) > 0 {
//				hostNames = append(hostNames, controllercommons.BuildHostNamesForCompositeDependency(dependencyInst, depComposite.Spec.ServiceTemplates)...)
//				intercellHttpRoutes = append(intercellHttpRoutes, controllercommons.BuildHttpRoutesForCompositeDependency(composite.Name, dependencyInst, depComposite.Spec.ServiceTemplates, false)...)
//			}
//		} else {
//			// unknown dependency kind
//			return nil, nil, nil, fmt.Errorf("unknown dependency kind '%s'", dependencyKind)
//		}
//	}
//	return hostNames, intercellHttpRoutes, intercellTcpRoutes, nil
//}
