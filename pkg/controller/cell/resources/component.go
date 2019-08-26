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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	. "github.com/cellery-io/mesh-controller/pkg/meta"
)

func MakeComponent(cell *v1alpha2.Cell, component *v1alpha2.Component) *v1alpha2.Component {

	if component.Spec.ScalingPolicy.IsKpa() {
		component.Spec.ScalingPolicy.Kpa.Selector = map[string]string{
			CellLabelKeySource: cell.Name,
		}
	}

	return &v1alpha2.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentName(cell, component),
			Namespace: cell.Namespace,
			Labels: UnionMaps(
				makeLabels(cell),
				map[string]string{
					ObservabilityComponentLabelKey: component.Name,
					AppLabelKey:                    fmt.Sprintf("%s--cell", ComponentName(cell, component)),
				},
			),
			Annotations: makeAnnotations(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: component.Spec,
	}
}

func RequireComponentUpdate(cell *v1alpha2.Cell, component *v1alpha2.Component) bool {
	return cell.Generation != cell.Status.ObservedGeneration ||
		component.Generation != cell.Status.ComponentGenerations[component.Name]
}

func CopyComponent(source, destination *v1alpha2.Component) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromComponent(cell *v1alpha2.Cell, component *v1alpha2.Component) {
	cell.Status.ComponentGenerations[component.Name] = component.Generation
	cell.Status.ComponentStatuses[component.Name] = component.Status.Status
}
