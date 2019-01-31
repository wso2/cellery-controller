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
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/celleryio/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/celleryio/mesh-controller/pkg/controller"
)

func CreateService(cell *v1alpha1.Cell, serviceTemplate v1alpha1.ServiceTemplateSpec) *v1alpha1.Service {
	serviceSpec := serviceTemplate.Spec.DeepCopy()
	serviceSpec.Container.Name = serviceTemplate.Name

	var serviceEnvVars []corev1.EnvVar
	for _, serviceTemplate := range cell.Spec.ServiceTemplates {
		envKey := strings.Replace(strings.ToUpper(serviceTemplate.Name), "-", "_", -1)
		envValue := ServiceName(cell, serviceTemplate) + "-service"
		serviceEnvVars = append(serviceEnvVars,
			corev1.EnvVar{
				Name:  envKey,
				Value: envValue,
			},
		)
	}
	serviceSpec.Container.Env = append(serviceEnvVars, serviceSpec.Container.Env...)

	return &v1alpha1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ServiceName(cell, serviceTemplate),
			Namespace:   cell.Namespace,
			Labels:      createLabels(cell),
			Annotations: createServiceAnnotations(cell),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCellOwnerRef(cell),
			},
		},
		Spec: *serviceSpec,
	}
}
