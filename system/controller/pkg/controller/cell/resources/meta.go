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
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
)

func createLabels(cell *v1alpha1.Cell) map[string]string {
	labels := make(map[string]string, len(cell.ObjectMeta.Labels)+1)

	labels[vick.CellNameLabelKey] = cell.Name

	for k, v := range cell.ObjectMeta.Labels {
		labels[k] = v
	}
	return labels
}

func NetworkPolicyName(cell *v1alpha1.Cell) string {
	return cell.Name + "-network"
}
