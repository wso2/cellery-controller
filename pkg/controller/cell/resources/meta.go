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
	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	. "cellery.io/cellery-controller/pkg/meta"
)

// func createLabels(cell *v1alpha2.Cell) map[string]string {
// 	labels := make(map[string]string, len(cell.ObjectMeta.Labels)+2)
// 	labels[mesh.CellLabelKey] = cell.Name
// 	labels[mesh.CellLabelKeySource] = cell.Name

// 	for k, v := range cell.ObjectMeta.Labels {
// 		labels[k] = v
// 	}
// 	return labels
// }

// func createServiceAnnotations(cell *v1alpha2.Cell) map[string]string {
// 	annotations := make(map[string]string, len(cell.ObjectMeta.Annotations))

// 	for k, v := range cell.ObjectMeta.Annotations {
// 		annotations[k] = v
// 	}
// 	return annotations
// }

func makeLabels(cell *v1alpha2.Cell) map[string]string {
	return UnionMaps(
		cell.Labels,
		map[string]string{
			CellLabelKey:                      cell.Name,
			CellLabelKeySource:                cell.Name,
			ObservabilityInstanceKindLabelKey: "Cell",
			ObservabilityInstanceLabelKey:     cell.Name,
		},
	)
}

func makeAnnotations(cell *v1alpha2.Cell) map[string]string {
	return UnionMaps(
		map[string]string{
			IstioSidecarInjectAnnotationKey: "false",
		},
		cell.Annotations,
	)
}

func NetworkPolicyName(cell *v1alpha2.Cell) string {
	return cell.Name + "--network"
}

func GatewayName(cell *v1alpha2.Cell) string {
	return cell.Name + "--gateway"
}

func GatewayNameFromInstanceName(instance string) string {
	return instance + "--gateway"
}

func GatewayK8sServiceName(gwName string) string {
	return gwName + "-service"
}

func TokenServiceName(cell *v1alpha2.Cell) string {
	return cell.Name + "--sts"
}

func EnvoyFilterName(cell *v1alpha2.Cell) string {
	return cell.Name + "--envoyfilter"
}

func ComponentName(cell *v1alpha2.Cell, component *v1alpha2.Component) string {
	return cell.Name + "--" + component.Name
}

func SecretName(cell *v1alpha2.Cell) string {
	return cell.Name + "--secret"
}
