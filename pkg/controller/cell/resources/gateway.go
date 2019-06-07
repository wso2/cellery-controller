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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateGateway(cell *v1alpha1.Cell) *v1alpha1.Gateway {
	gatewaySpec := cell.Spec.GatewayTemplate.Spec

	// Default to envoy gateway if not specified
	if len(gatewaySpec.Type) == 0 {
		gatewaySpec.Type = v1alpha1.GatewayTypeEnvoy
	}

	// Find the zero scale services to inform the gateway
	zeroScale := make(map[string]bool)
	for _, v := range cell.Spec.ServiceTemplates {
		if v.Spec.IsZeroScaled() {
			zeroScale[v.Name] = true
		} else {
			zeroScale[v.Name] = false
		}
	}

	for i, _ := range gatewaySpec.HTTPRoutes {
		if zeroScale[gatewaySpec.HTTPRoutes[i].Backend] {
			gatewaySpec.HTTPRoutes[i].ZeroScale = true
		}
		if gatewaySpec.Type == v1alpha1.GatewayTypeMicroGateway {
			gatewaySpec.HTTPRoutes[i].Backend = "http://" + cell.Name + "--" + gatewaySpec.HTTPRoutes[i].Backend + "-service"
		} else {
			gatewaySpec.HTTPRoutes[i].Backend = cell.Name + "--" + gatewaySpec.HTTPRoutes[i].Backend + "-service"
		}
	}

	for i, _ := range gatewaySpec.TCPRoutes {
		gatewaySpec.TCPRoutes[i].BackendHost = cell.Name + "--" + gatewaySpec.TCPRoutes[i].BackendHost + "-service"
	}

	return &v1alpha1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayName(cell),
			Namespace: cell.Namespace,
			Labels:    createLabels(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: gatewaySpec,
	}
}
