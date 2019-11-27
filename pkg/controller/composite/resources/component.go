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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"cellery.io/cellery-controller/pkg/controller"
	. "cellery.io/cellery-controller/pkg/meta"
)

func MakeComponent(composite *v1alpha2.Composite, component *v1alpha2.Component) *v1alpha2.Component {
	return &v1alpha2.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ComponentName(composite, component),
			Namespace: composite.Namespace,
			Labels: UnionMaps(
				makeLabels(composite),
				map[string]string{
					ObservabilityComponentLabelKey: component.Name,
					AppLabelKey:                    fmt.Sprintf("%s--composite", ComponentName(composite, component)),
				},
			),
			Annotations: makeAnnotations(composite),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCompositeOwnerRef(composite),
			},
		},
		Spec: component.Spec,
	}
}

func RequireComponentUpdate(composite *v1alpha2.Composite, component *v1alpha2.Component) bool {
	return composite.Generation != composite.Status.ObservedGeneration ||
		component.Generation != composite.Status.ComponentGenerations[component.Name]
}

func CopyComponent(source, destination *v1alpha2.Component) {
	destination.Spec = source.Spec
	destination.Labels = source.Labels
	destination.Annotations = source.Annotations
}

func StatusFromComponent(composite *v1alpha2.Composite, component *v1alpha2.Component) {
	composite.Status.ComponentGenerations[component.Name] = component.Generation
	composite.Status.ComponentStatuses[component.Name] = component.Status.Status
}
