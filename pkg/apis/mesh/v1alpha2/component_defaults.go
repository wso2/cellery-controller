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

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"cellery.io/cellery-controller/pkg/ptr"
)

func (c *Component) Default() {
	c.Spec.Default()
	c.Status.Default()
}

func (cs *ComponentSpec) Default() {
	if cs.Type == "" {
		cs.Type = ComponentTypeDeployment
	}
	cs.ScalingPolicy.Default()

	for i := range cs.Template.Containers {
		for j := range cs.Template.Containers[i].Ports {
			if len(cs.Template.Containers[i].Ports[j].Protocol) == 0 {
				cs.Template.Containers[i].Ports[j].Protocol = corev1.ProtocolTCP
			}
		}
	}

	for i, _ := range cs.Ports {
		cs.Ports[i].Default()
	}
}

func (sp *ComponentScalingPolicy) Default() {
	if sp.Replicas == nil && sp.Hpa == nil && sp.Kpa == nil {
		sp.Replicas = ptr.Int32(1)
	}
	if sp.Hpa != nil && sp.Hpa.MinReplicas == nil {
		sp.Hpa.MinReplicas = ptr.Int32(1)
	}

	if sp.Kpa != nil && sp.Kpa.MinReplicas == nil {
		sp.Kpa.MinReplicas = ptr.Int32(0)
	}
}

func (pm *PortMapping) Default() {
	if pm.Protocol == "" {
		pm.Protocol = ProtocolTCP
	}
	if pm.Name == "" {
		pm.Name = fmt.Sprintf("%d-%d", pm.Port, pm.TargetPort)
	}
}

func (cstat *ComponentStatus) Default() {
	if cstat.Status == "" {
		cstat.Status = ComponentCurrentStatusUnknown
	}
	if cstat.PersistantVolumeClaimGenerations == nil {
		cstat.PersistantVolumeClaimGenerations = make(map[string]int64)
	}
	if cstat.ConfigMapGenerations == nil {
		cstat.ConfigMapGenerations = make(map[string]int64)
	}
	if cstat.SecretGenerations == nil {
		cstat.SecretGenerations = make(map[string]int64)
	}
}
