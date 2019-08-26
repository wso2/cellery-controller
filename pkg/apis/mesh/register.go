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

package mesh

const (
	GroupName = "mesh.cellery.io"

	// Cellery Labels
	// CellLabelKey             = GroupName + "/cell"
	// CompositeLabelKey        = GroupName + "/composite"
	// CellGatewayLabelKey      = GroupName + "/gateway"
	// CellTokenServiceLabelKey = GroupName + "/sts"
	// CellServiceLabelKey      = GroupName + "/service"
	// ComponentLabelKey        = GroupName + "/component"

	// tem fix for avoid pilot validation https://github.com/istio/istio/issues/14797
	// CellLabelKeySource            = GroupName + ".cell"
	// CompositeLabelKeySource       = GroupName + ".composite"
	// ComponentLabelKeySource       = GroupName + ".component"
	// CompositeTokenServiceLabelKey = GroupName + "/composite-sts"

	// Cellery Annotations
	// CellServicesAnnotationKey     = GroupName + "/cell-services"
	// CellDependenciesAnnotationKey = GroupName + "/cell-dependencies"
	// CellOriginalGatewaySvcKey     = GroupName + "/original-gw-svc"

	// Cellery System
	SystemNamespace     = "cellery-system"
	SystemConfigMapName = "cellery-config"
	SystemSecretName    = "cellery-secret"
)
