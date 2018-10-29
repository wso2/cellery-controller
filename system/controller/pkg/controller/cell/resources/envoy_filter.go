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
	"github.com/wso2/product-vick/system/controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick"
	"github.com/wso2/product-vick/system/controller/pkg/apis/vick/v1alpha1"
	"github.com/wso2/product-vick/system/controller/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateEnvoyFilter(cell *v1alpha1.Cell) *v1alpha3.EnvoyFilter {
	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      EnvoyFilterName(cell),
			Namespace: cell.Namespace,
			Labels:    createLabels(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: v1alpha3.EnvoyFilterSpec{
			WorkloadLabels: map[string]string{
				vick.CellLabelKey: cell.Name,
			},
			Filters: []v1alpha3.Filter{
				{
					InsertPosition: v1alpha3.InsertPosition{
						Index: filterInsertPositionFirst,
					},
					ListenerMatch: v1alpha3.ListenerMatch{
						ListenerType:     filterListenerTypeInbound,
						ListenerProtocol: HTTPProtocol,
					},
					FilterName: baseFilterName,
					FilterType: HTTPProtocol,
					FilterConfig: v1alpha3.FilterConfig{
						GRPCService: v1alpha3.GRPCService{
							GoogleGRPC: v1alpha3.GoogleGRPC{
								TargetUri:  TokenServiceName(cell) + "-service:8080",
								StatPrefix: startPrefix,
							},
							Timeout: filterTimeout,
						},
					},
				},
				{
					InsertPosition: v1alpha3.InsertPosition{
						Index: filterInsertPositionLast,
					},
					ListenerMatch: v1alpha3.ListenerMatch{
						ListenerType:     filterListenerTypeOutbound,
						ListenerProtocol: HTTPProtocol,
					},
					FilterName: baseFilterName,
					FilterType: HTTPProtocol,
					FilterConfig: v1alpha3.FilterConfig{
						GRPCService: v1alpha3.GRPCService{
							GoogleGRPC: v1alpha3.GoogleGRPC{
								TargetUri:  TokenServiceName(cell) + "-service:8081",
								StatPrefix: startPrefix,
							},
							Timeout: filterTimeout,
						},
					},
				},
			},
		},
	}
}
