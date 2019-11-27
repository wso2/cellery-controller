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

	istioauthenticationv1alpha1 "cellery.io/cellery-controller/pkg/apis/istio/authentication/v1alpha1"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
)

// See https://istio.io/faq/security/#mysql-with-mtls
func MakeTlsPolicy(component *v1alpha2.Component) *istioauthenticationv1alpha1.Policy {
	return &istioauthenticationv1alpha1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TlsPolicyName(component),
			Namespace: component.Namespace,
			Labels:    makeLabels(component),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateComponentOwnerRef(component),
			},
		},
		Spec: istioauthenticationv1alpha1.PolicySpec{
			Targets: []*istioauthenticationv1alpha1.TargetSelector{
				{
					Name: ServiceName(component),
				},
			},
		},
	}
}

func RequireTlsPolicy(component *v1alpha2.Component) bool {
	b := false
	for _, v := range component.Spec.Ports {
		if v.Protocol == v1alpha2.ProtocolTCP {
			b = true
			break
		}
	}
	return b
}

func RequireTlsPolicyUpdate(component *v1alpha2.Component, policy *istioauthenticationv1alpha1.Policy) bool {
	return component.Generation != component.Status.ObservedGeneration ||
		policy.Generation != component.Status.TlsPolicyGeneration
}

func CopyTlsPolicy(source, destination *istioauthenticationv1alpha1.Policy) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromTlsPolicy(component *v1alpha2.Component, policy *istioauthenticationv1alpha1.Policy) {
	component.Status.TlsPolicyGeneration = policy.Generation
}
