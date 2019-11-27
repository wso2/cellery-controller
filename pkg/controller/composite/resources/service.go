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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"

	"cellery.io/cellery-controller/pkg/meta"

	"cellery.io/cellery-controller/pkg/controller"
)

func MakeOriginalComponentK8sService(composite *v1alpha2.Composite, compName string, targetPorts []int) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      K8sServiceName(compName),
			Namespace: composite.Namespace,
			Labels:    createLabelsForCurrentPodsWithPrevSvcName(composite, compName),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateCompositeOwnerRef(composite),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports:    *getAsContainerPorts(targetPorts),
			Selector: createLabelsForCurrentPodsWithPrevSvcName(composite, compName),
		},
	}
}

func getAsContainerPorts(ports []int) *[]corev1.ServicePort {
	var svcPorts []corev1.ServicePort
	for _, port := range ports {
		svcPorts = append(svcPorts, corev1.ServicePort{
			Name:       controller.HTTPServiceName,
			Protocol:   corev1.ProtocolTCP,
			Port:       80,
			TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: int32(port)},
		})
	}
	return &svcPorts
}

func createLabelsForCurrentPodsWithPrevSvcName(composite *v1alpha2.Composite, compName string) map[string]string {
	//labels := make(map[string]string, len(composite.ObjectMeta.Labels)+1)
	//for k, v := range composite.ObjectMeta.Labels {
	//	labels[k] = v
	//}
	//component := strings.Split(compName, "--")[1]
	labels := map[string]string{
		//meta.CellServiceLabelKey: fmt.Sprintf("%s--%s", composite.Name, component),
		meta.CompositeLabelKey: composite.Name,
	}

	return labels
}
