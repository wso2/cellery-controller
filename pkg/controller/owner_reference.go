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

package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
)

func CreateCellOwnerRef(obj metav1.Object) *metav1.OwnerReference {
	return metav1.NewControllerRef(obj, schema.GroupVersionKind{
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Version: v1alpha1.SchemeGroupVersion.Version,
		Kind:    "Cell",
	})
}

func CreateCompositeOwnerRef(obj metav1.Object) *metav1.OwnerReference {
	return metav1.NewControllerRef(obj, schema.GroupVersionKind{
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Version: v1alpha1.SchemeGroupVersion.Version,
		Kind:    "Composite",
	})
}

func CreateGatewayOwnerRef(obj metav1.Object) *metav1.OwnerReference {
	return metav1.NewControllerRef(obj, schema.GroupVersionKind{
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Version: v1alpha1.SchemeGroupVersion.Version,
		Kind:    "Gateway",
	})
}

func CreateTokenServiceOwnerRef(obj metav1.Object) *metav1.OwnerReference {
	return metav1.NewControllerRef(obj, schema.GroupVersionKind{
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Version: v1alpha1.SchemeGroupVersion.Version,
		Kind:    "TokenService",
	})
}

func CreateServiceOwnerRef(obj metav1.Object) *metav1.OwnerReference {
	return metav1.NewControllerRef(obj, schema.GroupVersionKind{
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Version: v1alpha1.SchemeGroupVersion.Version,
		Kind:    "Service",
	})
}

func CreateAutoscalerOwnerRef(obj metav1.Object) *metav1.OwnerReference {
	return metav1.NewControllerRef(obj, schema.GroupVersionKind{
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Version: v1alpha1.SchemeGroupVersion.Version,
		Kind:    "Autoscaler",
	})
}
