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

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
	. "cellery.io/cellery-controller/pkg/meta"
)

func MakeTokenService(cell *v1alpha2.Cell) *v1alpha2.TokenService {

	tSpec := cell.Spec.TokenService.Spec

	// if tSpec.InterceptMode != v1alpha2.InterceptModeNone {
	if cell.Spec.Gateway.Spec.Ingress.IngressExtensions.HasClusterIngress() {
		tSpec.InterceptMode = v1alpha2.InterceptModeOutbound
	}
	// else {
	// 		tSpec.InterceptMode = v1alpha1.InterceptModeAny
	// 	}
	// }
	tSpec.SecretName = SecretName(cell)
	tSpec.InstanceName = cell.Name
	tSpec.Selector = map[string]string{
		CellLabelKey: cell.Name,
	}

	return &v1alpha2.TokenService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TokenServiceName(cell),
			Namespace: cell.Namespace,
			Labels:    makeLabels(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: tSpec,
	}
}

func RequireTokenServiceUpdate(cell *v1alpha2.Cell, tokenService *v1alpha2.TokenService) bool {
	return cell.Generation != cell.Status.ObservedGeneration ||
		tokenService.Generation != cell.Status.TokenServiceGeneration
}

func CopyTokenService(source, destination *v1alpha2.TokenService) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromTokenService(cell *v1alpha2.Cell, tokenService *v1alpha2.TokenService) {
	cell.Status.TokenServiceStatus = tokenService.Status.Status
	cell.Status.TokenServiceGeneration = tokenService.Generation
}
