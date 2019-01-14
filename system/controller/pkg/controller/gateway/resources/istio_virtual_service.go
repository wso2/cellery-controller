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
	"github.com/wso2/product-vick/system/controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateIstioVirtualService(gateway *v1alpha1.Gateway) *v1alpha3.VirtualService {

	var routes []*v1alpha3.HTTPRoute

	for _, apiRoute := range gateway.Spec.APIRoutes {
		routes = append(routes, &v1alpha3.HTTPRoute{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Uri: &v1alpha3.StringMatch{
						//Regex: fmt.Sprintf("\\/%s(\\?.*|\\/.*|\\#.*|\\s*)", apiRoute.Context),
						Prefix: fmt.Sprintf("/%s/", apiRoute.Context),
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: apiRoute.Backend,
					},
				},
			},
			Rewrite: &v1alpha3.HTTPRewrite{
				Uri: "/",
			},
		})
	}

	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IstioVSName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Hosts:    []string{"*"},
			Gateways: []string{IstioGatewayName(gateway)},
			Http:     routes,
		},
	}
}

func CreateIstioVirtualServiceForIngress(gateway *v1alpha1.Gateway) *v1alpha3.VirtualService {

	var routes []*v1alpha3.HTTPRoute

	for _, apiRoute := range gateway.Spec.APIRoutes {
		if apiRoute.Global == true {
			routes = append(routes, &v1alpha3.HTTPRoute{
				Match: []*v1alpha3.HTTPMatchRequest{
					{
						Uri: &v1alpha3.StringMatch{
							Prefix: fmt.Sprintf("/%s/", apiRoute.Context),
						},
					},
				},
				Route: []*v1alpha3.DestinationWeight{
					{
						Destination: &v1alpha3.Destination{
							Host: gateway.Status.HostName,
						},
					},
				},
			})
		}
	}

	return &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IstioIngressVirtualServiceName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Hosts:    []string{"*"},
			Gateways: []string{"ingress-gateway"},
			Http:     routes,
		},
	}
}

