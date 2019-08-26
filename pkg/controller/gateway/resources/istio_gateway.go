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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func MakeIstioGateway(gateway *v1alpha2.Gateway) *v1alpha3.Gateway {

	var gatewayServers []*v1alpha3.Server

	for _, httpRoute := range gateway.Spec.Ingress.HTTPRoutes {
		gatewayServers = append(gatewayServers, &v1alpha3.Server{
			Hosts: []string{"*"},
			Port: &v1alpha3.Port{
				Number:   httpRoute.Port,
				Protocol: "HTTP",
				Name:     fmt.Sprintf("tcp-%d", httpRoute.Port),
			},
		})
	}

	for _, grpcRoute := range gateway.Spec.Ingress.GRPCRoutes {
		gatewayServers = append(gatewayServers, &v1alpha3.Server{
			Hosts: []string{"*"},
			Port: &v1alpha3.Port{
				Number:   grpcRoute.Port,
				Protocol: "GRPC",
				Name:     fmt.Sprintf("grpc-%d", grpcRoute.Port),
			},
		})
	}

	for _, tcpRoute := range gateway.Spec.Ingress.TCPRoutes {
		gatewayServers = append(gatewayServers, &v1alpha3.Server{
			Hosts: []string{"*"},
			Port: &v1alpha3.Port{
				Number:   tcpRoute.Port,
				Protocol: "TCP",
				Name:     fmt.Sprintf("tcp-%d", tcpRoute.Port),
			},
		})
	}
	// gatewayServers = append(gatewayServers, &v1alpha3.Server{
	// 	Hosts: []string{"*"},
	// 	Port: &v1alpha3.Port{
	// 		Number:   80,
	// 		Protocol: "HTTP",
	// 		Name:     fmt.Sprintf("http-%d", 80),
	// 	},
	// })

	return &v1alpha3.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IstioGatewayName(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.GatewaySpec{
			Servers:  gatewayServers,
			Selector: makeLabels(gateway),
		},
	}
}

func RequireIstioGateway(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.HasRoutes()
}

func RequireIstioGatewayUpdate(gateway *v1alpha2.Gateway, istioGateway *v1alpha3.Gateway) bool {
	return gateway.Generation != gateway.Status.ObservedGeneration ||
		istioGateway.Generation != gateway.Status.IstioGatewayGeneration
}

func CopyIstioGateway(source, destination *v1alpha3.Gateway) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromIstioGateway(gateway *v1alpha2.Gateway, istioGateway *v1alpha3.Gateway) {
	gateway.Status.IstioGatewayGeneration = istioGateway.Generation
}
