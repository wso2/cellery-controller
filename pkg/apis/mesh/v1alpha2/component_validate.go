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

package v1alpha2

import (
	"fmt"

	//apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (c *Component) Validate() error {
	//var allErrs field.ErrorList
	//allErrs = append(allErrs, c.Spec.Validate(field.NewPath("spec"))...)
	//return apierrors.NewInvalid(c.GroupVersionKind().GroupKind(), c.Name, allErrs)
	return nil
}

func (cs *ComponentSpec) Validate(fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	switch cs.Type {
	case ComponentTypeDeployment, ComponentTypeJob, ComponentTypeStatefulSet:
	default:
		allErrs = append(allErrs, field.Invalid(fldPath.Child("type"), cs.Type,
			fmt.Sprintf("must be one of '%s', '%s', '%s'", ComponentTypeDeployment, ComponentTypeJob, ComponentTypeStatefulSet)))
	}
	return allErrs
}
