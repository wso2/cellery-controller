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
	"k8s.io/api/extensions/v1beta1"
	networkv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateClusterIngress(gateway *v1alpha1.Gateway) *v1beta1.Ingress {

	var httpIngressPaths []networkv1.HTTPIngressPath

	httpIngressPaths = append(httpIngressPaths, networkv1.HTTPIngressPath{
		Path: "/",
		Backend: v1beta1.IngressBackend{
			ServiceName: GatewayK8sServiceName(gateway),
			ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
		},
	})

	// add callback endpoint to the ingress if oidc is enabled
	if gateway.Spec.OidcConfig != nil {
		httpIngressPaths = append(httpIngressPaths, networkv1.HTTPIngressPath{
			Path: "/_auth",
			Backend: v1beta1.IngressBackend{
				ServiceName: GatewayK8sServiceName(gateway),
				ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 15810},
			},
		})
	}

	var tlsIngressHosts []v1beta1.IngressTLS

	if len(gateway.Spec.Tls.Secret) > 0 {
		tlsIngressHosts = append(tlsIngressHosts, networkv1.IngressTLS{
			Hosts:      []string{gateway.Spec.Host},
			SecretName: gateway.Spec.Tls.Secret,
		})
	} else if len(gateway.Spec.Tls.Key) > 0 && len(gateway.Spec.Tls.Cert) > 0 {
		tlsIngressHosts = append(tlsIngressHosts, networkv1.IngressTLS{
			Hosts:      []string{gateway.Spec.Host},
			SecretName: ClusterIngressSecretName(gateway),
		})
	}

	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterIngressName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: gateway.Spec.Host,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: httpIngressPaths,
						},
					},
				},
			},
			TLS: tlsIngressHosts,
		},
	}
}
