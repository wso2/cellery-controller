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

	"cellery.io/cellery-controller/pkg/apis/istio/networking/v1alpha3"
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
)

func MakeOidcEnvoyFilter(gateway *v1alpha2.Gateway) *v1alpha3.EnvoyFilter {
	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      OidcEnvoyFilterName(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Spec: v1alpha3.EnvoyFilterSpec{
			WorkloadLabels: makeLabels(gateway),
			Filters: []v1alpha3.Filter{
				{
					InsertPosition: v1alpha3.InsertPosition{
						Index: filterInsertPositionFirst,
					},
					ListenerMatch: v1alpha3.ListenerMatch{
						ListenerType:     filterListenerTypeGateway,
						ListenerProtocol: HTTPProtocol,
					},
					FilterName: baseFilterName,
					FilterType: HTTPProtocol,
					FilterConfig: v1alpha3.FilterConfig{
						GRPCService: v1alpha3.GRPCService{
							GoogleGRPC: v1alpha3.GoogleGRPC{
								TargetUri:  "127.0.0.1:15800", // filter is attached as a sidecar
								StatPrefix: statPrefix,
							},
							Timeout: filterTimeout,
						},
					},
				},
			},
		},
	}
}

func OidcEnvoyFilterName(gateway *v1alpha2.Gateway) string {
	return gateway.Name + "-oidc"
}

func RequireOidcEnvoyFilter(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.IngressExtensions.HasOidc()
}

func RequireOidcEnvoyFilterUpdate(gateway *v1alpha2.Gateway, envoyFilter *v1alpha3.EnvoyFilter) bool {
	return gateway.Generation != gateway.Status.ObservedGeneration ||
		envoyFilter.Generation != gateway.Status.OidcEnvoyFilterGeneration
}

func CopyOidcEnvoyFilter(source, destination *v1alpha3.EnvoyFilter) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromOidcEnvoyFilter(gateway *v1alpha2.Gateway, envoyFilter *v1alpha3.EnvoyFilter) {
	gateway.Status.OidcEnvoyFilterGeneration = envoyFilter.Generation
}
