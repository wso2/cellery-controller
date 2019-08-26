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

package meta

import "github.com/cellery-io/mesh-controller/pkg/apis/mesh"

const (
	AppLabelKey     = "app"
	VersionLabelKey = "version"

	// Cellery Labels
	CellLabelKey         = mesh.GroupName + "/cell"
	CompositeLabelKey    = mesh.GroupName + "/composite"
	ComponentLabelKey    = mesh.GroupName + "/component"
	GatewayLabelKey      = mesh.GroupName + "/gateway"
	TokenServiceLabelKey = mesh.GroupName + "/token-service"

	// Cellery observability labels
	ObservabilityGroupPrefix          = "observability."
	ObservabilityInstanceLabelKey     = ObservabilityGroupPrefix + mesh.GroupName + "/instance"
	ObservabilityInstanceKindLabelKey = ObservabilityGroupPrefix + mesh.GroupName + "/instance-kind"
	ObservabilityWorkloadTypeLabelKey = ObservabilityGroupPrefix + mesh.GroupName + "/workload-type"
	ObservabilityComponentLabelKey    = ObservabilityGroupPrefix + mesh.GroupName + "/component"
	ObservabilityGatewayLabelKey      = ObservabilityGroupPrefix + mesh.GroupName + "/gateway"

	// tmp fix for avoid pilot validation https://github.com/istio/istio/issues/14797
	CellLabelKeySource            = mesh.GroupName + ".cell"
	CompositeLabelKeySource       = mesh.GroupName + ".composite"
	ComponentLabelKeySource       = mesh.GroupName + ".component"
	CompositeTokenServiceLabelKey = mesh.GroupName + "/composite-sts"
)
