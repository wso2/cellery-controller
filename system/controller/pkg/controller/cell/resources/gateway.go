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
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateGateway(cell *v1alpha1.Cell) *v1alpha1.Gateway {
	gatewaySpec := cell.Spec.GatewayTemplate.Spec

	for i, _ := range gatewaySpec.APIRoutes {
		gatewaySpec.APIRoutes[i].Backend = "http://" + cell.Name + "--" + gatewaySpec.APIRoutes[i].Backend + "-service"
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
