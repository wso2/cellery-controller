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
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	. "github.com/cellery-io/mesh-controller/pkg/meta"
)

func MakeNetworkPolicy(cell *v1alpha2.Cell) *networkv1.NetworkPolicy {

	cellName := cell.Name
	gatewayName := GatewayName(cell)
	var componentNames []string

	for _, component := range cell.Spec.Components {
		componentNames = append(componentNames, ComponentName(cell, &component))
	}

	return &networkv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NetworkPolicyName(cell),
			Namespace: cell.Namespace,
			Labels:    makeLabels(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: networkv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					CellLabelKey: cellName,
				},
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      ComponentLabelKey,
						Operator: metav1.LabelSelectorOpIn,
						Values:   componentNames,
					},
				},
			},
			PolicyTypes: []networkv1.PolicyType{
				networkv1.PolicyTypeIngress,
			},
			Ingress: []networkv1.NetworkPolicyIngressRule{
				{
					From: []networkv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									CellLabelKey:    cellName,
									GatewayLabelKey: gatewayName,
								},
							},
						},
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									CellLabelKey: cellName,
								},
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      ComponentLabelKey,
										Operator: metav1.LabelSelectorOpIn,
										Values:   componentNames,
									},
								},
							},
						},
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									mesh.TelepresenceLabelKey: "telepresence",
								},
							},
						},
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"name": "knative-serving",
								},
							},
						},
					},
				},
			},
		},
	}
}

func RequireNetworkPolicyUpdate(cell *v1alpha2.Cell, networkPolicy *networkv1.NetworkPolicy) bool {
	return cell.Generation != cell.Status.ObservedGeneration ||
		networkPolicy.Generation != cell.Status.NetworkPolicyGeneration
}

func CopyNetworkPolicy(source, destination *networkv1.NetworkPolicy) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromNetworkPolicy(cell *v1alpha2.Cell, networkPolicy *networkv1.NetworkPolicy) {
	cell.Status.NetworkPolicyGeneration = networkPolicy.Generation
}
