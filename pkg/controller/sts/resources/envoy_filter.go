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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

func CreateEnvoyFilter(tokenService *v1alpha1.TokenService) *v1alpha3.EnvoyFilter {

	var cellName string
	cellName, ok := tokenService.Labels[mesh.CellLabelKey]
	if !ok {
		cellName = tokenService.Name
	}

	var filters []v1alpha3.Filter
	switch tokenService.Spec.InterceptMode {
	case v1alpha1.InterceptModeInbound:
		filters = append(filters, buildInboundFilter(tokenService))
	case v1alpha1.InterceptModeOutbound:
		filters = append(filters, buildOutboundFilter(tokenService))
	case v1alpha1.InterceptModeAny:
		filters = append(filters, buildInboundFilter(tokenService))
		filters = append(filters, buildOutboundFilter(tokenService))
	}

	return &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      EnvoyFilterName(tokenService),
			Namespace: tokenService.Namespace,
			Labels:    createTokenServiceLabels(tokenService),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateTokenServiceOwnerRef(tokenService),
			},
		},
		Spec: v1alpha3.EnvoyFilterSpec{
			WorkloadLabels: map[string]string{
				mesh.CellLabelKey: cellName,
			},
			Filters: filters,
		},
	}
}

func buildInboundFilter(tokenService *v1alpha1.TokenService) v1alpha3.Filter {
	return v1alpha3.Filter{
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
					TargetUri:  TokenServiceK8sServiceName(tokenService) + ":8080",
					StatPrefix: statPrefix,
				},
				Timeout: filterTimeout,
			},
		},
	}
}

func buildOutboundFilter(tokenService *v1alpha1.TokenService) v1alpha3.Filter {
	return v1alpha3.Filter{
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
					TargetUri:  TokenServiceK8sServiceName(tokenService) + ":8081",
					StatPrefix: statPrefix,
				},
				Timeout: filterTimeout,
			},
		},
	}
}
