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

package v1alpha2

import "cellery.io/cellery-controller/pkg/ptr"

func (g *Gateway) SetDefaults() {
	g.Spec.SetDefaults()
	g.Status.SetDefaults()
}

func (gs *GatewaySpec) SetDefaults() {
	if !gs.Ingress.IngressExtensions.HasApiPublisher() {
		for _, h := range gs.Ingress.HTTPRoutes {
			if h.Global {
				gs.Ingress.IngressExtensions.ApiPublisher = &ApiPublisherConfig{}
				return
			}
		}
	}

	// if cs.Type == "" {
	// 	cs.Type = ComponentTypeDeployment
	// }
	gs.ScalingPolicy.SetDefaults()
	// for i, _ := range cs.Ports {
	// 	cs.Ports[i].SetDefaults()
	// }
}

func (sp *GwScalingPolicy) SetDefaults() {
	if sp.Hpa != nil {
		if sp.Hpa.MinReplicas == nil {
			// default min replicas = 1
			sp.Hpa.MinReplicas = ptr.Int32(1)
		}
	}
	if sp.Hpa == nil && sp.Replicas == nil {
		sp.Replicas = ptr.Int32(1)
	}
}

// func (pm *PortMapping) SetDefaults() {
// 	if pm.Protocol == "" {
// 		pm.Protocol = ProtocolTCP
// 	}
// 	if pm.Name == "" {
// 		pm.Name = fmt.Sprintf("%d-%d", pm.Port, pm.TargetPort)
// 	}
// }

func (gs *GatewayStatus) SetDefaults() {
	// if cstat.Status == "" {
	// 	cstat.Status = ComponentCurrentStatusUnknown
	// }
}
