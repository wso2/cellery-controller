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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	networkv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/config"
	"cellery.io/cellery-controller/pkg/controller"
	"cellery.io/cellery-controller/pkg/crypto"
)

func MakeClusterIngress(gateway *v1alpha2.Gateway) *v1beta1.Ingress {

	var httpIngressPaths []networkv1.HTTPIngressPath
	extensions := gateway.Spec.Ingress.IngressExtensions
	httpIngressPaths = append(httpIngressPaths, networkv1.HTTPIngressPath{
		Path: "/",
		Backend: v1beta1.IngressBackend{
			ServiceName: ServiceName(gateway),
			ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
		},
	})

	// add callback endpoint to the ingress if oidc is enabled
	if extensions.HasOidc() {
		httpIngressPaths = append(httpIngressPaths, networkv1.HTTPIngressPath{
			Path: "/_auth",
			Backend: v1beta1.IngressBackend{
				ServiceName: ServiceName(gateway),
				ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 15810},
			},
		})
	}

	var tlsIngressHosts []v1beta1.IngressTLS

	if extensions.ClusterIngress.HasSecret() {
		tlsIngressHosts = append(tlsIngressHosts, networkv1.IngressTLS{
			Hosts:      []string{extensions.ClusterIngress.Host},
			SecretName: extensions.ClusterIngress.Tls.Secret,
		})
	} else if extensions.ClusterIngress.HasCertAndKey() {
		tlsIngressHosts = append(tlsIngressHosts, networkv1.IngressTLS{
			Hosts:      []string{extensions.ClusterIngress.Host},
			SecretName: ClusterIngressSecretName(gateway),
		})
	}

	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterIngressName(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []networkv1.IngressRule{
				{
					Host: extensions.ClusterIngress.Host,
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

func ClusterIngressName(gateway *v1alpha2.Gateway) string {
	return gateway.Name + "-ingress"
}

func RequireClusterIngress(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.IngressExtensions.HasClusterIngress()
}

func RequireClusterIngressUpdate(gateway *v1alpha2.Gateway, ingress *v1beta1.Ingress) bool {
	return gateway.Generation != gateway.Status.ObservedGeneration ||
		ingress.Generation != gateway.Status.ClusterIngressGeneration
}

func CopyClusterIngress(source, destination *v1beta1.Ingress) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromClusterIngress(gateway *v1alpha2.Gateway, ingress *v1beta1.Ingress) {
	gateway.Status.ClusterIngressGeneration = ingress.Generation
}

func MakeClusterIngressSecret(gateway *v1alpha2.Gateway, cfg config.Interface) (*corev1.Secret, error) {

	pvtKey, err := cfg.PrivateKey()
	if err != nil {
		return nil, err
	}
	ci := gateway.Spec.Ingress.IngressExtensions.ClusterIngress

	key, err := crypto.TryDecrypt(ci.Tls.Key, pvtKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt the tls.key: %v", err)
	}

	cert, err := crypto.TryDecrypt(ci.Tls.Cert, pvtKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt the tls.cert: %v", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterIngressSecretName(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.key": key,
			"tls.crt": cert,
		},
	}, nil
}

func ClusterIngressSecretName(gateway *v1alpha2.Gateway) string {
	return ClusterIngressName(gateway) + "-secret"
}

func RequireClusterIngressSecret(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.IngressExtensions.HasClusterIngress() &&
		gateway.Spec.Ingress.IngressExtensions.ClusterIngress.HasCertAndKey()
}

func RequireClusterIngressSecretUpdate(gateway *v1alpha2.Gateway, secret *corev1.Secret) bool {
	return gateway.Generation != gateway.Status.ObservedGeneration ||
		secret.Generation != gateway.Status.ClusterIngressSecretGeneration
}

func CopyClusterIngressSecret(source, destination *corev1.Secret) {
	destination.Data = source.Data
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromClusterIngressSecret(gateway *v1alpha2.Gateway, secret *corev1.Secret) {
	gateway.Status.ClusterIngressSecretGeneration = secret.Generation
}
