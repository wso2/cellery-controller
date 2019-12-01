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

	"cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
)

func MakeVirtualService(gateway *v1alpha2.Gateway) *v1alpha3.VirtualService {

	var httpRoutes []*v1alpha3.HTTPRoute

	for _, httpRoute := range gateway.Spec.Ingress.HTTPRoutes {
		r := &v1alpha3.HTTPRoute{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Uri: &v1alpha3.StringMatch{
						// https://github.com/istio/istio/issues/8076
						Prefix: fmt.Sprintf("%s/", httpRoute.Context),
					},
				},
				{
					Uri: &v1alpha3.StringMatch{
						//Regex: fmt.Sprintf("\\/%s(\\?.*|\\/.*|\\#.*|\\s*)", apiRoute.Context),
						Prefix: httpRoute.Context,
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: httpRoute.Destination.Host,
						Port: &v1alpha3.PortSelector{
							Number: httpRoute.Destination.Port,
						},
					},
				},
			},
			Rewrite: &v1alpha3.HTTPRewrite{
				Uri: "/",
			},
			AppendHeaders: func() map[string]string {
				if httpRoute.ZeroScale {
					return map[string]string{
						"knative-serving-namespace": gateway.Namespace,
						"knative-serving-revision":  httpRoute.Destination.Host,
					}
				}
				return map[string]string{}
			}(),
		}
		httpRoutes = append(httpRoutes, r)
	}

	for _, grpcRoute := range gateway.Spec.Ingress.GRPCRoutes {
		r := &v1alpha3.HTTPRoute{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Port: grpcRoute.Port,
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: grpcRoute.Destination.Host,
						Port: &v1alpha3.PortSelector{
							Number: grpcRoute.Destination.Port,
						},
					},
				},
			},
			AppendHeaders: func() map[string]string {
				if grpcRoute.ZeroScale {
					return map[string]string{
						"knative-serving-namespace": gateway.Namespace,
						"knative-serving-revision":  grpcRoute.Destination.Host,
					}
				}
				return map[string]string{}
			}(),
		}
		httpRoutes = append(httpRoutes, r)
	}

	var tcpRoutes []*v1alpha3.TCPRoute

	for _, tcpRoute := range gateway.Spec.Ingress.TCPRoutes {
		tcpRoutes = append(tcpRoutes, &v1alpha3.TCPRoute{
			Match: []*v1alpha3.L4MatchAttributes{
				{
					Port: tcpRoute.Port,
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: tcpRoute.Destination.Host,
						Port: &v1alpha3.PortSelector{
							Number: tcpRoute.Destination.Port,
						},
					},
				},
			},
		})
	}

	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IstioVirtualServiceName(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Hosts:    []string{"*"},
			Gateways: []string{IstioGatewayName(gateway)},
			Http:     httpRoutes,
			Tcp:      tcpRoutes,
		},
	}
}

func RequireVirtualService(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.HasRoutes()
}

func RequireVirtualServiceUpdate(gateway *v1alpha2.Gateway, virtualService *v1alpha3.VirtualService) bool {
	return gateway.Generation != gateway.Status.ObservedGeneration ||
		virtualService.Generation != gateway.Status.VirtualServiceGeneration
}

func CopyVirtualService(source, destination *v1alpha3.VirtualService) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromVirtualService(gateway *v1alpha2.Gateway, virtualService *v1alpha3.VirtualService) {
	gateway.Status.VirtualServiceGeneration = virtualService.Generation
}

// func CreateIstioVirtualServiceForIngress(gateway *v1alpha1.Gateway) *v1alpha3.VirtualService {

// 	var routes []*v1alpha3.HTTPRoute

// 	for _, apiRoute := range gateway.Spec.HTTPRoutes {
// 		if apiRoute.Global == true {
// 			routes = append(routes, &v1alpha3.HTTPRoute{
// 				Match: []*v1alpha3.HTTPMatchRequest{
// 					{
// 						Uri: &v1alpha3.StringMatch{
// 							Prefix: fmt.Sprintf("/%s/", apiRoute.Context),
// 						},
// 					},
// 				},
// 				Route: []*v1alpha3.DestinationWeight{
// 					{
// 						Destination: &v1alpha3.Destination{
// 							Host: gateway.Status.HostName,
// 						},
// 					},
// 				},
// 			})
// 		}
// 	}

// 	return &v1alpha3.VirtualService{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      IstioIngressVirtualServiceName(gateway),
// 			Namespace: gateway.Namespace,
// 			Labels:    createGatewayLabels(gateway),
// 			OwnerReferences: []metav1.OwnerReference{
// 				*controller.CreateGatewayOwnerRef(gateway),
// 			},
// 		},
// 		Spec: v1alpha3.VirtualServiceSpec{
// 			Hosts:    []string{"*"},
// 			Gateways: []string{"ingress-gateway"},
// 			Http:     routes,
// 		},
// 	}
// }
