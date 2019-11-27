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
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha1"
	"cellery.io/cellery-controller/pkg/controller"
)

func TestCreateIstioGateway(t *testing.T) {
	gateway := &v1alpha1.Gateway{
		Spec: v1alpha1.GatewaySpec{
			Type: v1alpha1.GatewayTypeEnvoy,
			Host: "test.com",
			HTTPRoutes: []v1alpha1.HTTPRoute{
				{
					Authenticate: true,
					Global:       true,
					Backend:      "mytestservice",
					Context:      "hello",
					Definitions: []v1alpha1.APIDefinition{
						{
							Path:   "sayHello",
							Method: "GET",
						},
					},
				},
			},
		},
	}

	istiogw := CreateIstioGateway(gateway)
	expected := &v1alpha3.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IstioGatewayName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.GatewaySpec{
			Servers:  getGatewayServers(gateway),
			Selector: createGatewayLabels(gateway),
		},
	}

	if diff := cmp.Diff(expected, istiogw); diff != "" {
		t.Errorf("CreateIstioGateway (-expected, +actual)\n%v", diff)
	}
}

func getGatewayServers(gateway *v1alpha1.Gateway) []*v1alpha3.Server {
	var gatewayServers []*v1alpha3.Server

	for _, tcpRoute := range gateway.Spec.TCPRoutes {
		gatewayServers = append(gatewayServers, &v1alpha3.Server{
			Hosts: []string{"*"},
			Port: &v1alpha3.Port{
				Number:   tcpRoute.Port,
				Protocol: "TCP",
				Name:     fmt.Sprintf("tcp-%d", tcpRoute.Port),
			},
		})
	}
	gatewayServers = append(gatewayServers, &v1alpha3.Server{
		Hosts: []string{"*"},
		Port: &v1alpha3.Port{
			Number:   80,
			Protocol: "HTTP",
			Name:     fmt.Sprintf("http-%d", 80),
		},
	})

	return gatewayServers
}
