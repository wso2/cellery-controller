/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	. "github.com/cellery-io/mesh-controller/pkg/meta"
)

func MakeGateway(cell *v1alpha2.Cell) *v1alpha2.Gateway {
	gatewaySpec := cell.Spec.Gateway.Spec

	// Find the zero scale services to inform the gateway
	zeroScale := make(map[string]bool)
	for _, v := range cell.Spec.Components {
		if v.Spec.ScalingPolicy.IsKpa() {
			zeroScale[v.Name] = true
		} else {
			zeroScale[v.Name] = false
		}
	}

	for i, _ := range gatewaySpec.Ingress.HTTPRoutes {
		if zeroScale[gatewaySpec.Ingress.HTTPRoutes[i].Destination.Host] {
			gatewaySpec.Ingress.HTTPRoutes[i].ZeroScale = true
		}
		// if gatewaySpec.Type == v1alpha1.GatewayTypeMicroGateway {
		// 	gatewaySpec.HTTPRoutes[i].Backend = "http://" + cell.Name + "--" + gatewaySpec.HTTPRoutes[i].Backend + "-service"
		// } else {
		gatewaySpec.Ingress.HTTPRoutes[i].Destination.Host = cell.Name + "--" + gatewaySpec.Ingress.HTTPRoutes[i].Destination.Host + "-service"
		// }
	}

	for i, _ := range gatewaySpec.Ingress.TCPRoutes {
		gatewaySpec.Ingress.TCPRoutes[i].Destination.Host = cell.Name + "--" + gatewaySpec.Ingress.TCPRoutes[i].Destination.Host + "-service"
	}

	for i, _ := range gatewaySpec.Ingress.GRPCRoutes {
		if zeroScale[gatewaySpec.Ingress.GRPCRoutes[i].Destination.Host] {
			gatewaySpec.Ingress.GRPCRoutes[i].ZeroScale = true
		}
		gatewaySpec.Ingress.GRPCRoutes[i].Destination.Host = cell.Name + "--" + gatewaySpec.Ingress.GRPCRoutes[i].Destination.Host + "-service"
	}

	if gatewaySpec.Ingress.IngressExtensions.HasOidc() {
		gatewaySpec.Ingress.IngressExtensions.OidcConfig.SecretName = SecretName(cell)
		gatewaySpec.Ingress.IngressExtensions.OidcConfig.JwtIssuer = GatewayName(cell)
		gatewaySpec.Ingress.IngressExtensions.OidcConfig.JwtAudience = cell.Name
	}

	return &v1alpha2.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayName(cell),
			Namespace: cell.Namespace,
			Labels: UnionMaps(
				makeLabels(cell),
				map[string]string{
					ObservabilityGatewayLabelKey: "gateway",
				},
			),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: gatewaySpec,
	}
}

func RequireGatewayUpdate(cell *v1alpha2.Cell, gateway *v1alpha2.Gateway) bool {
	return cell.Generation != cell.Status.ObservedGeneration ||
		gateway.Generation != cell.Status.GatewayGeneration
}

func CopyGateway(source, destination *v1alpha2.Gateway) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromGateway(cell *v1alpha2.Cell, gateway *v1alpha2.Gateway) {
	cell.Status.GatewayServiceName = gateway.Status.ServiceName
	cell.Status.GatewayStatus = gateway.Status.Status
	cell.Status.GatewayGeneration = gateway.Generation
}
