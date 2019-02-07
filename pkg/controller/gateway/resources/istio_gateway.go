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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateIstioGateway(gateway *v1alpha1.Gateway) *v1alpha3.Gateway {
	return &v1alpha3.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      IstioGatewayName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.GatewaySpec{
			Servers: []*v1alpha3.Server{
				{
					Hosts: []string{"*"},
					Port: &v1alpha3.Port{
						Number:   80,
						Protocol: "HTTP2",
						Name:     "http2",
					},
				},
				{
					Hosts: []string{"*"},
					Port: &v1alpha3.Port{
						Number:   443,
						Protocol: "HTTPS",
						Name:     "https",
					},
					Tls: &v1alpha3.Server_TLSOptions{
						Mode:              "MUTUAL",
						ServerCertificate: "/etc/certs/cert-chain.pem",
						PrivateKey:        "/etc/certs/key.pem",
						CaCertificates:    "/etc/certs/root-cert.pem",
					},
				},
			},
			Selector: createGatewayLabels(gateway),
		},
	}
}
