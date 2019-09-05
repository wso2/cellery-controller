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

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

func createLabels(composite *v1alpha1.Composite) map[string]string {
	labels := make(map[string]string, len(composite.ObjectMeta.Labels)+2)
	labels[mesh.CompositeLabelKey] = composite.Name
	labels[mesh.CompositeLabelKeySource] = composite.Name

	for k, v := range composite.ObjectMeta.Labels {
		labels[k] = v
	}
	return labels
}

func addTokenServiceLabels(labels map[string]string) map[string]string {
	newLabels := make(map[string]string, len(labels)+1)
	labels[mesh.CompositeTokenServiceLabelKey] = "true"
	for k, v := range labels {
		newLabels[k] = v
	}
	return newLabels
}

func createServiceAnnotations(composite *v1alpha1.Composite) map[string]string {
	annotations := make(map[string]string, len(composite.ObjectMeta.Annotations))

	for k, v := range composite.ObjectMeta.Annotations {
		annotations[k] = v
	}
	return annotations
}

func NetworkPolicyName(cell *v1alpha1.Cell) string {
	return cell.Name + "--network"
}

func GatewayName(cell *v1alpha1.Cell) string {
	return cell.Name + "--gateway"
}

func GatewayNameFromInstanceName(instance string) string {
	return instance + "--gateway"
}

func TokenServiceName(composite *v1alpha1.Composite) string {
	return "composite--sts"
}

func EnvoyFilterName(cell *v1alpha1.Cell) string {
	return cell.Name + "--envoyfilter"
}

func ServiceName(composite *v1alpha1.Composite, serviceTemplate v1alpha1.ServiceTemplateSpec) string {
	return composite.Name + "--" + serviceTemplate.Name
}

func SecretName(composite *v1alpha1.Composite) string {
	return "composite-sts-secret"
}

func K8sServiceName(name string) string {
	return fmt.Sprintf("%s-service", name)
}
