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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	. "github.com/cellery-io/mesh-controller/pkg/meta"
)

func makeLabels(component *v1alpha2.Component) map[string]string {
	return UnionMaps(
		map[string]string{
			AppLabelKey:                       component.Name,
			VersionLabelKey:                   "v1.0.0",
			ObservabilityComponentLabelKey:    component.Name,
			ObservabilityWorkloadTypeLabelKey: string(component.Spec.Type),
		},
		component.Labels,
		map[string]string{
			ComponentLabelKey: component.Name,
		},
	)
}

func makeSelector(component *v1alpha2.Component) *metav1.LabelSelector {
	return &metav1.LabelSelector{MatchLabels: makeLabels(component)}
}

func makePodAnnotations(component *v1alpha2.Component) map[string]string {
	return UnionMaps(
		map[string]string{
			IstioSidecarInjectAnnotationKey: "true",
		},
		component.Labels,
	)
}

func makeServingDeploymentSelector(component *v1alpha2.Component) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		"serving.knative.dev/configuration": ServingConfigurationName(component),
		"serving.knative.dev/revision":      ServingRevisionName(component),
	})
}

func ServiceName(component *v1alpha2.Component) string {
	return component.Name + "-service"
}

func DeploymentName(component *v1alpha2.Component) string {
	return component.Name + "-deployment"
}

func JobName(component *v1alpha2.Component) string {
	return component.Name + "-job"
}

func StatefulSetName(component *v1alpha2.Component) string {
	return component.Name + "-statefulset"
}

func PersistentVolumeClaimName(component *v1alpha2.Component, volumeClaim *v1alpha2.VolumeClaim) string {
	return component.Name + "-" + volumeClaim.Template.Name + "-pvc"
}

func ConfigMapName(component *v1alpha2.Component, configMap *corev1.ConfigMap) string {
	return component.Name + "-" + configMap.Name + "-config"
}

func SecretName(component *v1alpha2.Component, secret *corev1.Secret) string {
	return component.Name + "-" + secret.Name + "-secret"
}

func HpaName(component *v1alpha2.Component) string {
	return component.Name + "-hpa"
}

func ServingConfigurationName(component *v1alpha2.Component) string {
	return component.Name
}

func ServingVirtualServiceName(component *v1alpha2.Component) string {
	return component.Name + "-mesh"
}

func ServingRevisionName(component *v1alpha2.Component) string {
	return ServingConfigurationName(component) + "-service"
}
