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

package v1alpha2

import (
	"fmt"
	"math"
	"net"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	apimachineryvalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

//
// Original Source: https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/core/validation/validation.go
//
const isNegativeErrorMsg string = apimachineryvalidation.IsNegativeErrorMsg

//const isInvalidQuotaResource string = `must be a standard resource for quota`
//const fieldImmutableErrorMsg string = apimachineryvalidation.FieldImmutableErrorMsg
const isNotIntegerErrorMsg string = `must be an integer`

//const isNotPositiveErrorMsg string = `must be greater than zero`

//var pdPartitionErrorMsg string = validation.InclusiveRangeError(1, 255)
var fileModeErrorMsg = "must be a number between 0 and 0777 (octal), both inclusive"

//// BannedOwners is a black list of object that are not allowed to be owners.
//var BannedOwners = apimachineryvalidation.BannedOwners
//
//var iscsiInitiatorIqnRegex = regexp.MustCompile(`iqn\.\d{4}-\d{2}\.([[:alnum:]-.]+)(:[^,;*&$|\s]+)$`)
//var iscsiInitiatorEuiRegex = regexp.MustCompile(`^eui.[[:alnum:]]{16}$`)
//var iscsiInitiatorNaaRegex = regexp.MustCompile(`^naa.[[:alnum:]]{32}$`)
//
//// ValidateHasLabel requires that metav1.ObjectMeta has a Label with key and expectedValue
//func ValidateHasLabel(meta metav1.ObjectMeta, fldPath *field.Path, key, expectedValue string) field.ErrorList {
//	allErrs := field.ErrorList{}
//	actualValue, found := meta.Labels[key]
//	if !found {
//		allErrs = append(allErrs, field.Required(fldPath.Child("labels").Key(key),
//			fmt.Sprintf("must be '%s'", expectedValue)))
//		return allErrs
//	}
//	if actualValue != expectedValue {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("labels").Key(key), meta.Labels,
//			fmt.Sprintf("must be '%s'", expectedValue)))
//	}
//	return allErrs
//}
//
//// ValidateAnnotations validates that a set of annotations are correctly defined.
//func ValidateAnnotations(annotations map[string]string, fldPath *field.Path) field.ErrorList {
//	return apimachineryvalidation.ValidateAnnotations(annotations, fldPath)
//}
//
func ValidateDNS1123Label(value string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, msg := range validation.IsDNS1123Label(value) {
		allErrs = append(allErrs, field.Invalid(fldPath, value, msg))
	}
	return allErrs
}

// ValidateDNS1123Subdomain validates that a name is a proper DNS subdomain.
func ValidateDNS1123Subdomain(value string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, msg := range validation.IsDNS1123Subdomain(value) {
		allErrs = append(allErrs, field.Invalid(fldPath, value, msg))
	}
	return allErrs
}

// ValidateNameFunc validates that the provided name is valid for a given resource type.
// Not all resources have the same validation rules for names. Prefix is true
// if the name will have a value appended to it.  If the name is not valid,
// this returns a list of descriptions of individual characteristics of the
// value that were not valid.  Otherwise this returns an empty list or nil.
type ValidateNameFunc apimachineryvalidation.ValidateNameFunc

//
//// ValidatePodName can be used to check whether the given pod name is valid.
//// Prefix indicates this name will be used as part of generation, in which case
//// trailing dashes are allowed.
//var ValidatePodName = apimachineryvalidation.NameIsDNSSubdomain
//
//// ValidateReplicationControllerName can be used to check whether the given replication
//// controller name is valid.
//// Prefix indicates this name will be used as part of generation, in which case
//// trailing dashes are allowed.
//var ValidateReplicationControllerName = apimachineryvalidation.NameIsDNSSubdomain
//
//// ValidateServiceName can be used to check whether the given service name is valid.
//// Prefix indicates this name will be used as part of generation, in which case
//// trailing dashes are allowed.
//var ValidateServiceName = apimachineryvalidation.NameIsDNS1035Label

// ValidateNodeName can be used to check whether the given node name is valid.
// Prefix indicates this name will be used as part of generation, in which case
// trailing dashes are allowed.
var ValidateNodeName = apimachineryvalidation.NameIsDNSSubdomain

//// ValidateNamespaceName can be used to check whether the given namespace name is valid.
//// Prefix indicates this name will be used as part of generation, in which case
//// trailing dashes are allowed.
//var ValidateNamespaceName = apimachineryvalidation.ValidateNamespaceName
//
//// ValidateLimitRangeName can be used to check whether the given limit range name is valid.
//// Prefix indicates this name will be used as part of generation, in which case
//// trailing dashes are allowed.
//var ValidateLimitRangeName = apimachineryvalidation.NameIsDNSSubdomain
//
//// ValidateResourceQuotaName can be used to check whether the given
//// resource quota name is valid.
//// Prefix indicates this name will be used as part of generation, in which case
//// trailing dashes are allowed.
//var ValidateResourceQuotaName = apimachineryvalidation.NameIsDNSSubdomain
//
// ValidateSecretName can be used to check whether the given secret name is valid.
// Prefix indicates this name will be used as part of generation, in which case
// trailing dashes are allowed.
var ValidateSecretName = apimachineryvalidation.NameIsDNSSubdomain

// ValidateServiceAccountName can be used to check whether the given service account name is valid.
// Prefix indicates this name will be used as part of generation, in which case
// trailing dashes are allowed.
var ValidateServiceAccountName = apimachineryvalidation.ValidateServiceAccountName

//// ValidateEndpointsName can be used to check whether the given endpoints name is valid.
//// Prefix indicates this name will be used as part of generation, in which case
//// trailing dashes are allowed.
//var ValidateEndpointsName = apimachineryvalidation.NameIsDNSSubdomain
//
//// ValidateClusterName can be used to check whether the given cluster name is valid.
//var ValidateClusterName = apimachineryvalidation.ValidateClusterName
//
//// ValidateClassName can be used to check whether the given class name is valid.
//// It is defined here to avoid import cycle between pkg/apis/storage/validation
//// (where it should be) and this file.
//var ValidateClassName = apimachineryvalidation.NameIsDNSSubdomain

// ValidatePiorityClassName can be used to check whether the given priority
// class name is valid.
var ValidatePriorityClassName = apimachineryvalidation.NameIsDNSSubdomain

// ValidateRuntimeClassName can be used to check whether the given RuntimeClass name is valid.
// Prefix indicates this name will be used as part of generation, in which case
// trailing dashes are allowed.
func ValidateRuntimeClassName(name string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	for _, msg := range apimachineryvalidation.NameIsDNSSubdomain(name, false) {
		allErrs = append(allErrs, field.Invalid(fldPath, name, msg))
	}
	return allErrs
}

// Validates that given value is not negative.
func ValidateNonnegativeField(value int64, fldPath *field.Path) field.ErrorList {
	return apimachineryvalidation.ValidateNonnegativeField(value, fldPath)
}

// Validates that a Quantity is not negative
func ValidateNonnegativeQuantity(value resource.Quantity, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if value.Cmp(resource.Quantity{}) < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath, value.String(), isNegativeErrorMsg))
	}
	return allErrs
}

//// Validates that a Quantity is positive
//func ValidatePositiveQuantityValue(value resource.Quantity, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if value.Cmp(resource.Quantity{}) <= 0 {
//		allErrs = append(allErrs, field.Invalid(fldPath, value.String(), isNotPositiveErrorMsg))
//	}
//	return allErrs
//}
//
//func ValidateImmutableField(newVal, oldVal interface{}, fldPath *field.Path) field.ErrorList {
//	return apimachineryvalidation.ValidateImmutableField(newVal, oldVal, fldPath)
//}
//
//func ValidateImmutableAnnotation(newVal string, oldVal string, annotation string, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	if oldVal != newVal {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("annotations", annotation), newVal, fieldImmutableErrorMsg))
//	}
//	return allErrs
//}
//
// ValidateObjectMeta validates an object's metadata on creation. It expects that name generation has already
// been performed.
// It doesn't return an error for rootscoped resources with namespace, because namespace should already be cleared before.
// TODO: Remove calls to this method scattered in validations of specific resources, e.g., ValidatePodUpdate.
func ValidateObjectMeta(meta *metav1.ObjectMeta, requiresNamespace bool, nameFn ValidateNameFunc, fldPath *field.Path) field.ErrorList {
	allErrs := apimachineryvalidation.ValidateObjectMeta(meta, requiresNamespace, apimachineryvalidation.ValidateNameFunc(nameFn), fldPath)
	// run additional checks for the finalizer name
	for i := range meta.Finalizers {
		allErrs = append(allErrs, validateKubeFinalizerName(string(meta.Finalizers[i]), fldPath.Child("finalizers").Index(i))...)
	}
	return allErrs
}

//// ValidateObjectMetaUpdate validates an object's metadata when updated
//func ValidateObjectMetaUpdate(newMeta, oldMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
//	allErrs := apimachineryvalidation.ValidateObjectMetaUpdate(newMeta, oldMeta, fldPath)
//	// run additional checks for the finalizer name
//	for i := range newMeta.Finalizers {
//		allErrs = append(allErrs, validateKubeFinalizerName(string(newMeta.Finalizers[i]), fldPath.Child("finalizers").Index(i))...)
//	}
//
//	return allErrs
//}
//
//
func IsMatchedVolume(name string, volumes map[string]corev1.VolumeSource) bool {
	if _, ok := volumes[name]; ok {
		return true
	}
	return false
}

func isMatchedDevice(name string, volumes map[string]corev1.VolumeSource) (bool, bool) {
	if source, ok := volumes[name]; ok {
		if source.PersistentVolumeClaim != nil {
			return true, true
		}
		return true, false
	}
	return false, false
}

func mountNameAlreadyExists(name string, devices map[string]string) bool {
	if _, ok := devices[name]; ok {
		return true
	}
	return false
}

func mountPathAlreadyExists(mountPath string, devices map[string]string) bool {
	for _, devPath := range devices {
		if mountPath == devPath {
			return true
		}
	}

	return false
}

func deviceNameAlreadyExists(name string, mounts map[string]string) bool {
	if _, ok := mounts[name]; ok {
		return true
	}
	return false
}

func devicePathAlreadyExists(devicePath string, mounts map[string]string) bool {
	for _, mountPath := range mounts {
		if mountPath == devicePath {
			return true
		}
	}

	return false
}

func validateHostPathVolumeSource(hostPath *corev1.HostPathVolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(hostPath.Path) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
		return allErrs
	}

	allErrs = append(allErrs, validatePathNoBacksteps(hostPath.Path, fldPath.Child("path"))...)
	allErrs = append(allErrs, validateHostPathType(hostPath.Type, fldPath.Child("type"))...)
	return allErrs
}

//func validateGitRepoVolumeSource(gitRepo *corev1.GitRepoVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(gitRepo.Repository) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("repository"), ""))
//	}
//
//	pathErrs := validateLocalDescendingPath(gitRepo.Directory, fldPath.Child("directory"))
//	allErrs = append(allErrs, pathErrs...)
//	return allErrs
//}
//
//func validateISCSIVolumeSource(iscsi *corev1.ISCSIVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(iscsi.TargetPortal) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("targetPortal"), ""))
//	}
//	if len(iscsi.IQN) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("iqn"), ""))
//	} else {
//		if !strings.HasPrefix(iscsi.IQN, "iqn") && !strings.HasPrefix(iscsi.IQN, "eui") && !strings.HasPrefix(iscsi.IQN, "naa") {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format starting with iqn, eui, or naa"))
//		} else if strings.HasPrefix(iscsi.IQN, "iqn") && !iscsiInitiatorIqnRegex.MatchString(iscsi.IQN) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format"))
//		} else if strings.HasPrefix(iscsi.IQN, "eui") && !iscsiInitiatorEuiRegex.MatchString(iscsi.IQN) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format"))
//		} else if strings.HasPrefix(iscsi.IQN, "naa") && !iscsiInitiatorNaaRegex.MatchString(iscsi.IQN) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format"))
//		}
//	}
//	if iscsi.Lun < 0 || iscsi.Lun > 255 {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("lun"), iscsi.Lun, validation.InclusiveRangeError(0, 255)))
//	}
//	if (iscsi.DiscoveryCHAPAuth || iscsi.SessionCHAPAuth) && iscsi.SecretRef == nil {
//		allErrs = append(allErrs, field.Required(fldPath.Child("secretRef"), ""))
//	}
//	if iscsi.InitiatorName != nil {
//		initiator := *iscsi.InitiatorName
//		if !strings.HasPrefix(initiator, "iqn") && !strings.HasPrefix(initiator, "eui") && !strings.HasPrefix(initiator, "naa") {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format starting with iqn, eui, or naa"))
//		}
//		if strings.HasPrefix(initiator, "iqn") && !iscsiInitiatorIqnRegex.MatchString(initiator) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format"))
//		} else if strings.HasPrefix(initiator, "eui") && !iscsiInitiatorEuiRegex.MatchString(initiator) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format"))
//		} else if strings.HasPrefix(initiator, "naa") && !iscsiInitiatorNaaRegex.MatchString(initiator) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format"))
//		}
//	}
//	return allErrs
//}
//
//func validateISCSIPersistentVolumeSource(iscsi *corev1.ISCSIPersistentVolumeSource, pvName string, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(iscsi.TargetPortal) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("targetPortal"), ""))
//	}
//	if iscsi.InitiatorName != nil && len(pvName+":"+iscsi.TargetPortal) > 64 {
//		tooLongErr := "Total length of <volume name>:<iscsi.targetPortal> must be under 64 characters if iscsi.initiatorName is specified."
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("targetportal"), iscsi.TargetPortal, tooLongErr))
//	}
//	if len(iscsi.IQN) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("iqn"), ""))
//	} else {
//		if !strings.HasPrefix(iscsi.IQN, "iqn") && !strings.HasPrefix(iscsi.IQN, "eui") && !strings.HasPrefix(iscsi.IQN, "naa") {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format"))
//		} else if strings.HasPrefix(iscsi.IQN, "iqn") && !iscsiInitiatorIqnRegex.MatchString(iscsi.IQN) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format"))
//		} else if strings.HasPrefix(iscsi.IQN, "eui") && !iscsiInitiatorEuiRegex.MatchString(iscsi.IQN) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format"))
//		} else if strings.HasPrefix(iscsi.IQN, "naa") && !iscsiInitiatorNaaRegex.MatchString(iscsi.IQN) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("iqn"), iscsi.IQN, "must be valid format"))
//		}
//	}
//	if iscsi.Lun < 0 || iscsi.Lun > 255 {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("lun"), iscsi.Lun, validation.InclusiveRangeError(0, 255)))
//	}
//	if (iscsi.DiscoveryCHAPAuth || iscsi.SessionCHAPAuth) && iscsi.SecretRef == nil {
//		allErrs = append(allErrs, field.Required(fldPath.Child("secretRef"), ""))
//	}
//	if iscsi.SecretRef != nil {
//		if len(iscsi.SecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretRef", "name"), ""))
//		}
//	}
//	if iscsi.InitiatorName != nil {
//		initiator := *iscsi.InitiatorName
//		if !strings.HasPrefix(initiator, "iqn") && !strings.HasPrefix(initiator, "eui") && !strings.HasPrefix(initiator, "naa") {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format"))
//		}
//		if strings.HasPrefix(initiator, "iqn") && !iscsiInitiatorIqnRegex.MatchString(initiator) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format"))
//		} else if strings.HasPrefix(initiator, "eui") && !iscsiInitiatorEuiRegex.MatchString(initiator) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format"))
//		} else if strings.HasPrefix(initiator, "naa") && !iscsiInitiatorNaaRegex.MatchString(initiator) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("initiatorname"), initiator, "must be valid format"))
//		}
//	}
//	return allErrs
//}
//
//func validateFCVolumeSource(fc *corev1.FCVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(fc.TargetWWNs) < 1 && len(fc.WWIDs) < 1 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("targetWWNs"), "must specify either targetWWNs or wwids, but not both"))
//	}
//
//	if len(fc.TargetWWNs) != 0 && len(fc.WWIDs) != 0 {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("targetWWNs"), fc.TargetWWNs, "targetWWNs and wwids can not be specified simultaneously"))
//	}
//
//	if len(fc.TargetWWNs) != 0 {
//		if fc.Lun == nil {
//			allErrs = append(allErrs, field.Required(fldPath.Child("lun"), "lun is required if targetWWNs is specified"))
//		} else {
//			if *fc.Lun < 0 || *fc.Lun > 255 {
//				allErrs = append(allErrs, field.Invalid(fldPath.Child("lun"), fc.Lun, validation.InclusiveRangeError(0, 255)))
//			}
//		}
//	}
//	return allErrs
//}
//
//func validateGCEPersistentDiskVolumeSource(pd *corev1.GCEPersistentDiskVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(pd.PDName) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("pdName"), ""))
//	}
//	if pd.Partition < 0 || pd.Partition > 255 {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("partition"), pd.Partition, pdPartitionErrorMsg))
//	}
//	return allErrs
//}
//
//func validateAWSElasticBlockStoreVolumeSource(PD *corev1.AWSElasticBlockStoreVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(PD.VolumeID) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeID"), ""))
//	}
//	if PD.Partition < 0 || PD.Partition > 255 {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("partition"), PD.Partition, pdPartitionErrorMsg))
//	}
//	return allErrs
//}

func validateSecretVolumeSource(secretSource *corev1.SecretVolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(secretSource.SecretName) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("secretName"), ""))
	}

	secretMode := secretSource.DefaultMode
	if secretMode != nil && (*secretMode > 0777 || *secretMode < 0) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("defaultMode"), *secretMode, fileModeErrorMsg))
	}

	itemsPath := fldPath.Child("items")
	for i, kp := range secretSource.Items {
		itemPath := itemsPath.Index(i)
		allErrs = append(allErrs, validateKeyToPath(&kp, itemPath)...)
	}
	return allErrs
}

func validateConfigMapVolumeSource(configMapSource *corev1.ConfigMapVolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(configMapSource.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	}

	configMapMode := configMapSource.DefaultMode
	if configMapMode != nil && (*configMapMode > 0777 || *configMapMode < 0) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("defaultMode"), *configMapMode, fileModeErrorMsg))
	}

	itemsPath := fldPath.Child("items")
	for i, kp := range configMapSource.Items {
		itemPath := itemsPath.Index(i)
		allErrs = append(allErrs, validateKeyToPath(&kp, itemPath)...)
	}
	return allErrs
}

func validateKeyToPath(kp *corev1.KeyToPath, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(kp.Key) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), ""))
	}
	if len(kp.Path) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
	}
	allErrs = append(allErrs, validateLocalNonReservedPath(kp.Path, fldPath.Child("path"))...)
	if kp.Mode != nil && (*kp.Mode > 0777 || *kp.Mode < 0) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("mode"), *kp.Mode, fileModeErrorMsg))
	}

	return allErrs
}

func validatePersistentClaimVolumeSource(claim *corev1.PersistentVolumeClaimVolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(claim.ClaimName) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("claimName"), ""))
	}
	return allErrs
}

//func validateNFSVolumeSource(nfs *corev1.NFSVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(nfs.Server) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("server"), ""))
//	}
//	if len(nfs.Path) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
//	}
//	if !path.IsAbs(nfs.Path) {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("path"), nfs.Path, "must be an absolute path"))
//	}
//	return allErrs
//}
//
//func validateQuobyteVolumeSource(quobyte *corev1.QuobyteVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(quobyte.Registry) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("registry"), "must be a host:port pair or multiple pairs separated by commas"))
//	} else if len(quobyte.Tenant) >= 65 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("tenant"), "must be a UUID and may not exceed a length of 64 characters"))
//	} else {
//		for _, hostPortPair := range strings.Split(quobyte.Registry, ",") {
//			if _, _, err := net.SplitHostPort(hostPortPair); err != nil {
//				allErrs = append(allErrs, field.Invalid(fldPath.Child("registry"), quobyte.Registry, "must be a host:port pair or multiple pairs separated by commas"))
//			}
//		}
//	}
//
//	if len(quobyte.Volume) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volume"), ""))
//	}
//	return allErrs
//}
//
//func validateGlusterfsVolumeSource(glusterfs *corev1.GlusterfsVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(glusterfs.EndpointsName) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("endpoints"), ""))
//	}
//	if len(glusterfs.Path) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
//	}
//	return allErrs
//}
//func validateGlusterfsPersistentVolumeSource(glusterfs *corev1.GlusterfsPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(glusterfs.EndpointsName) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("endpoints"), ""))
//	}
//	if len(glusterfs.Path) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
//	}
//	if glusterfs.EndpointsNamespace != nil {
//		endpointNs := glusterfs.EndpointsNamespace
//		if *endpointNs == "" {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("endpointsNamespace"), *endpointNs, "if the endpointnamespace is set, it must be a valid namespace name"))
//		} else {
//			for _, msg := range ValidateNamespaceName(*endpointNs, false) {
//				allErrs = append(allErrs, field.Invalid(fldPath.Child("endpointsNamespace"), *endpointNs, msg))
//			}
//		}
//	}
//	return allErrs
//}
//
//func validateFlockerVolumeSource(flocker *corev1.FlockerVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(flocker.DatasetName) == 0 && len(flocker.DatasetUUID) == 0 {
//		//TODO: consider adding a RequiredOneOf() error for this and similar cases
//		allErrs = append(allErrs, field.Required(fldPath, "one of datasetName and datasetUUID is required"))
//	}
//	if len(flocker.DatasetName) != 0 && len(flocker.DatasetUUID) != 0 {
//		allErrs = append(allErrs, field.Invalid(fldPath, "resource", "datasetName and datasetUUID can not be specified simultaneously"))
//	}
//	if strings.Contains(flocker.DatasetName, "/") {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("datasetName"), flocker.DatasetName, "must not contain '/'"))
//	}
//	return allErrs
//}
//
//var validVolumeDownwardAPIFieldPathExpressions = sets.NewString(
//	"metadata.name",
//	"metadata.namespace",
//	"metadata.labels",
//	"metadata.annotations",
//	"metadata.uid")
//
//func validateDownwardAPIVolumeFile(file *corev1.DownwardAPIVolumeFile, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	if len(file.Path) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
//	}
//	allErrs = append(allErrs, validateLocalNonReservedPath(file.Path, fldPath.Child("path"))...)
//	if file.FieldRef != nil {
//		allErrs = append(allErrs, validateObjectFieldSelector(file.FieldRef, &validVolumeDownwardAPIFieldPathExpressions, fldPath.Child("fieldRef"))...)
//		if file.ResourceFieldRef != nil {
//			allErrs = append(allErrs, field.Invalid(fldPath, "resource", "fieldRef and resourceFieldRef can not be specified simultaneously"))
//		}
//	} else if file.ResourceFieldRef != nil {
//		allErrs = append(allErrs, validateContainerResourceFieldSelector(file.ResourceFieldRef, &validContainerResourceFieldPathExpressions, fldPath.Child("resourceFieldRef"), true)...)
//	} else {
//		allErrs = append(allErrs, field.Required(fldPath, "one of fieldRef and resourceFieldRef is required"))
//	}
//	if file.Mode != nil && (*file.Mode > 0777 || *file.Mode < 0) {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("mode"), *file.Mode, fileModeErrorMsg))
//	}
//
//	return allErrs
//}
//
//func validateDownwardAPIVolumeSource(downwardAPIVolume *corev1.DownwardAPIVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	downwardAPIMode := downwardAPIVolume.DefaultMode
//	if downwardAPIMode != nil && (*downwardAPIMode > 0777 || *downwardAPIMode < 0) {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("defaultMode"), *downwardAPIMode, fileModeErrorMsg))
//	}
//
//	for _, file := range downwardAPIVolume.Items {
//		allErrs = append(allErrs, validateDownwardAPIVolumeFile(&file, fldPath)...)
//	}
//	return allErrs
//}

var supportedHostPathTypes = sets.NewString(
	string(corev1.HostPathUnset),
	string(corev1.HostPathDirectoryOrCreate),
	string(corev1.HostPathDirectory),
	string(corev1.HostPathFileOrCreate),
	string(corev1.HostPathFile),
	string(corev1.HostPathSocket),
	string(corev1.HostPathCharDev),
	string(corev1.HostPathBlockDev))

func validateHostPathType(hostPathType *corev1.HostPathType, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if hostPathType != nil && !supportedHostPathTypes.Has(string(*hostPathType)) {
		allErrs = append(allErrs, field.NotSupported(fldPath, hostPathType, supportedHostPathTypes.List()))
	}

	return allErrs
}

// This validate will make sure targetPath:
// 1. is not abs path
// 2. does not have any element which is ".."
func validateLocalDescendingPath(targetPath string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if path.IsAbs(targetPath) {
		allErrs = append(allErrs, field.Invalid(fldPath, targetPath, "must be a relative path"))
	}

	allErrs = append(allErrs, validatePathNoBacksteps(targetPath, fldPath)...)

	return allErrs
}

// validatePathNoBacksteps makes sure the targetPath does not have any `..` path elements when split
//
// This assumes the OS of the apiserver and the nodes are the same. The same check should be done
// on the node to ensure there are no backsteps.
func validatePathNoBacksteps(targetPath string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	parts := strings.Split(filepath.ToSlash(targetPath), "/")
	for _, item := range parts {
		if item == ".." {
			allErrs = append(allErrs, field.Invalid(fldPath, targetPath, "must not contain '..'"))
			break // even for `../../..`, one error is sufficient to make the point
		}
	}
	return allErrs
}

// validateMountPropagation verifies that MountPropagation field is valid and
// allowed for given container.
func validateMountPropagation(mountPropagation *corev1.MountPropagationMode, container *corev1.Container, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if mountPropagation == nil {
		return allErrs
	}

	supportedMountPropagations := sets.NewString(string(corev1.MountPropagationBidirectional), string(corev1.MountPropagationHostToContainer), string(corev1.MountPropagationNone))
	if !supportedMountPropagations.Has(string(*mountPropagation)) {
		allErrs = append(allErrs, field.NotSupported(fldPath, *mountPropagation, supportedMountPropagations.List()))
	}

	if container == nil {
		// The container is not available yet, e.g. during validation of
		// PodPreset. Stop validation now, Pod validation will refuse final
		// Pods with Bidirectional propagation in non-privileged containers.
		return allErrs
	}

	privileged := container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged
	if *mountPropagation == corev1.MountPropagationBidirectional && !privileged {
		allErrs = append(allErrs, field.Forbidden(fldPath, "Bidirectional mount propagation is available only to privileged containers"))
	}
	return allErrs
}

// This validate will make sure targetPath:
// 1. is not abs path
// 2. does not contain any '..' elements
// 3. does not start with '..'
func validateLocalNonReservedPath(targetPath string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateLocalDescendingPath(targetPath, fldPath)...)
	// Don't report this error if the check for .. elements already caught it.
	if strings.HasPrefix(targetPath, "..") && !strings.HasPrefix(targetPath, "../") {
		allErrs = append(allErrs, field.Invalid(fldPath, targetPath, "must not start with '..'"))
	}
	return allErrs
}

//func validateRBDVolumeSource(rbd *corev1.RBDVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(rbd.CephMonitors) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("monitors"), ""))
//	}
//	if len(rbd.RBDImage) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("image"), ""))
//	}
//	return allErrs
//}
//
//func validateRBDPersistentVolumeSource(rbd *corev1.RBDPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(rbd.CephMonitors) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("monitors"), ""))
//	}
//	if len(rbd.RBDImage) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("image"), ""))
//	}
//	return allErrs
//}
//
//func validateCinderVolumeSource(cd *corev1.CinderVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(cd.VolumeID) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeID"), ""))
//	}
//	if cd.SecretRef != nil {
//		if len(cd.SecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretRef", "name"), ""))
//		}
//	}
//	return allErrs
//}
//
//func validateCinderPersistentVolumeSource(cd *corev1.CinderPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(cd.VolumeID) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeID"), ""))
//	}
//	if cd.SecretRef != nil {
//		if len(cd.SecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretRef", "name"), ""))
//		}
//		if len(cd.SecretRef.Namespace) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretRef", "namespace"), ""))
//		}
//	}
//	return allErrs
//}
//
//func validateCephFSVolumeSource(cephfs *corev1.CephFSVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(cephfs.Monitors) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("monitors"), ""))
//	}
//	return allErrs
//}
//
//func validateCephFSPersistentVolumeSource(cephfs *corev1.CephFSPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(cephfs.Monitors) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("monitors"), ""))
//	}
//	return allErrs
//}
//
//func validateFlexVolumeSource(fv *corev1.FlexVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(fv.Driver) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("driver"), ""))
//	}
//
//	// Make sure user-specified options don't use kubernetes namespaces
//	for k := range fv.Options {
//		namespace := k
//		if parts := strings.SplitN(k, "/", 2); len(parts) == 2 {
//			namespace = parts[0]
//		}
//		normalized := "." + strings.ToLower(namespace)
//		if strings.HasSuffix(normalized, ".kubernetes.io") || strings.HasSuffix(normalized, ".k8s.io") {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("options").Key(k), k, "kubernetes.io and k8s.io namespaces are reserved"))
//		}
//	}
//
//	return allErrs
//}
//
//func validateFlexPersistentVolumeSource(fv *corev1.FlexPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(fv.Driver) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("driver"), ""))
//	}
//
//	// Make sure user-specified options don't use kubernetes namespaces
//	for k := range fv.Options {
//		namespace := k
//		if parts := strings.SplitN(k, "/", 2); len(parts) == 2 {
//			namespace = parts[0]
//		}
//		normalized := "." + strings.ToLower(namespace)
//		if strings.HasSuffix(normalized, ".kubernetes.io") || strings.HasSuffix(normalized, ".k8s.io") {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("options").Key(k), k, "kubernetes.io and k8s.io namespaces are reserved"))
//		}
//	}
//
//	return allErrs
//}
//
//func validateAzureFile(azure *corev1.AzureFileVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if azure.SecretName == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("secretName"), ""))
//	}
//	if azure.ShareName == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("shareName"), ""))
//	}
//	return allErrs
//}
//
//func validateAzureFilePV(azure *corev1.AzureFilePersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if azure.SecretName == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("secretName"), ""))
//	}
//	if azure.ShareName == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("shareName"), ""))
//	}
//	if azure.SecretNamespace != nil {
//		if len(*azure.SecretNamespace) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretNamespace"), ""))
//		}
//	}
//	return allErrs
//}
//
//func validateAzureDisk(azure *corev1.AzureDiskVolumeSource, fldPath *field.Path) field.ErrorList {
//	var supportedCachingModes = sets.NewString(string(corev1.AzureDataDiskCachingNone), string(corev1.AzureDataDiskCachingReadOnly), string(corev1.AzureDataDiskCachingReadWrite))
//	var supportedDiskKinds = sets.NewString(string(corev1.AzureSharedBlobDisk), string(corev1.AzureDedicatedBlobDisk), string(corev1.AzureManagedDisk))
//
//	diskURISupportedManaged := []string{"/subscriptions/{sub-id}/resourcegroups/{group-name}/providers/microsoft.compute/disks/{disk-id}"}
//	diskURISupportedblob := []string{"https://{account-name}.blob.corev1.windows.net/{container-name}/{disk-name}.vhd"}
//
//	allErrs := field.ErrorList{}
//	if azure.DiskName == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("diskName"), ""))
//	}
//
//	if azure.DataDiskURI == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("diskURI"), ""))
//	}
//
//	if azure.CachingMode != nil && !supportedCachingModes.Has(string(*azure.CachingMode)) {
//		allErrs = append(allErrs, field.NotSupported(fldPath.Child("cachingMode"), *azure.CachingMode, supportedCachingModes.List()))
//	}
//
//	if azure.Kind != nil && !supportedDiskKinds.Has(string(*azure.Kind)) {
//		allErrs = append(allErrs, field.NotSupported(fldPath.Child("kind"), *azure.Kind, supportedDiskKinds.List()))
//	}
//
//	// validate that DiskUri is the correct format
//	if azure.Kind != nil && *azure.Kind == corev1.AzureManagedDisk && strings.Index(azure.DataDiskURI, "/subscriptions/") != 0 {
//		allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskURI"), azure.DataDiskURI, diskURISupportedManaged))
//	}
//
//	if azure.Kind != nil && *azure.Kind != corev1.AzureManagedDisk && strings.Index(azure.DataDiskURI, "https://") != 0 {
//		allErrs = append(allErrs, field.NotSupported(fldPath.Child("diskURI"), azure.DataDiskURI, diskURISupportedblob))
//	}
//
//	return allErrs
//}
//
//func validateVsphereVolumeSource(cd *corev1.VsphereVirtualDiskVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(cd.VolumePath) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumePath"), ""))
//	}
//	return allErrs
//}
//
//func validatePhotonPersistentDiskVolumeSource(cd *corev1.PhotonPersistentDiskVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(cd.PdID) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("pdID"), ""))
//	}
//	return allErrs
//}
//
//func validatePortworxVolumeSource(pwx *corev1.PortworxVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(pwx.VolumeID) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeID"), ""))
//	}
//	return allErrs
//}
//
//func validateScaleIOVolumeSource(sio *corev1.ScaleIOVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if sio.Gateway == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("gateway"), ""))
//	}
//	if sio.System == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("system"), ""))
//	}
//	if sio.VolumeName == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeName"), ""))
//	}
//	return allErrs
//}
//
//func validateScaleIOPersistentVolumeSource(sio *corev1.ScaleIOPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if sio.Gateway == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("gateway"), ""))
//	}
//	if sio.System == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("system"), ""))
//	}
//	if sio.VolumeName == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeName"), ""))
//	}
//	return allErrs
//}
//
//func validateLocalVolumeSource(ls *corev1.LocalVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if ls.Path == "" {
//		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
//		return allErrs
//	}
//
//	allErrs = append(allErrs, validatePathNoBacksteps(ls.Path, fldPath.Child("path"))...)
//	return allErrs
//}
//
//func validateStorageOSVolumeSource(storageos *corev1.StorageOSVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(storageos.VolumeName) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeName"), ""))
//	} else {
//		allErrs = append(allErrs, ValidateDNS1123Label(storageos.VolumeName, fldPath.Child("volumeName"))...)
//	}
//	if len(storageos.VolumeNamespace) > 0 {
//		allErrs = append(allErrs, ValidateDNS1123Label(storageos.VolumeNamespace, fldPath.Child("volumeNamespace"))...)
//	}
//	if storageos.SecretRef != nil {
//		if len(storageos.SecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretRef", "name"), ""))
//		}
//	}
//	return allErrs
//}
//
//func validateStorageOSPersistentVolumeSource(storageos *corev1.StorageOSPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(storageos.VolumeName) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeName"), ""))
//	} else {
//		allErrs = append(allErrs, ValidateDNS1123Label(storageos.VolumeName, fldPath.Child("volumeName"))...)
//	}
//	if len(storageos.VolumeNamespace) > 0 {
//		allErrs = append(allErrs, ValidateDNS1123Label(storageos.VolumeNamespace, fldPath.Child("volumeNamespace"))...)
//	}
//	if storageos.SecretRef != nil {
//		if len(storageos.SecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretRef", "name"), ""))
//		}
//		if len(storageos.SecretRef.Namespace) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("secretRef", "namespace"), ""))
//		}
//	}
//	return allErrs
//}
//
//func ValidateCSIDriverName(driverName string, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	if len(driverName) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath, ""))
//	}
//
//	if len(driverName) > 63 {
//		allErrs = append(allErrs, field.TooLong(fldPath, driverName, 63))
//	}
//
//	for _, msg := range validation.IsDNS1123Subdomain(strings.ToLower(driverName)) {
//		allErrs = append(allErrs, field.Invalid(fldPath, driverName, msg))
//	}
//
//	return allErrs
//}
//
//func validateCSIPersistentVolumeSource(csi *corev1.CSIPersistentVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	allErrs = append(allErrs, ValidateCSIDriverName(csi.Driver, fldPath.Child("driver"))...)
//
//	if len(csi.VolumeHandle) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("volumeHandle"), ""))
//	}
//
//	if csi.ControllerPublishSecretRef != nil {
//		if len(csi.ControllerPublishSecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("controllerPublishSecretRef", "name"), ""))
//		} else {
//			allErrs = append(allErrs, ValidateDNS1123Label(csi.ControllerPublishSecretRef.Name, fldPath.Child("name"))...)
//		}
//		if len(csi.ControllerPublishSecretRef.Namespace) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("controllerPublishSecretRef", "namespace"), ""))
//		} else {
//			allErrs = append(allErrs, ValidateDNS1123Label(csi.ControllerPublishSecretRef.Namespace, fldPath.Child("namespace"))...)
//		}
//	}
//
//	if csi.ControllerExpandSecretRef != nil {
//		if len(csi.ControllerExpandSecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("controllerExpandSecretRef", "name"), ""))
//		} else {
//			allErrs = append(allErrs, ValidateDNS1123Label(csi.ControllerExpandSecretRef.Name, fldPath.Child("name"))...)
//		}
//		if len(csi.ControllerExpandSecretRef.Namespace) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("controllerExpandSecretRef", "namespace"), ""))
//		} else {
//			allErrs = append(allErrs, ValidateDNS1123Label(csi.ControllerExpandSecretRef.Namespace, fldPath.Child("namespace"))...)
//		}
//	}
//
//	if csi.NodePublishSecretRef != nil {
//		if len(csi.NodePublishSecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("nodePublishSecretRef ", "name"), ""))
//		} else {
//			allErrs = append(allErrs, ValidateDNS1123Label(csi.NodePublishSecretRef.Name, fldPath.Child("name"))...)
//		}
//		if len(csi.NodePublishSecretRef.Namespace) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("nodePublishSecretRef ", "namespace"), ""))
//		} else {
//			allErrs = append(allErrs, ValidateDNS1123Label(csi.NodePublishSecretRef.Namespace, fldPath.Child("namespace"))...)
//		}
//	}
//
//	return allErrs
//}
//
//func validateCSIVolumeSource(csi *corev1.CSIVolumeSource, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	allErrs = append(allErrs, ValidateCSIDriverName(csi.Driver, fldPath.Child("driver"))...)
//
//	if csi.NodePublishSecretRef != nil {
//		if len(csi.NodePublishSecretRef.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("nodePublishSecretRef ", "name"), ""))
//		} else {
//			for _, msg := range ValidateSecretName(csi.NodePublishSecretRef.Name, false) {
//				allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), csi.NodePublishSecretRef.Name, msg))
//			}
//		}
//	}
//
//	return allErrs
//}
//
//// ValidatePersistentVolumeName checks that a name is appropriate for a
//// PersistentVolumeName object.
//var ValidatePersistentVolumeName = apimachineryvalidation.NameIsDNSSubdomain
//
//var supportedAccessModes = sets.NewString(string(corev1.ReadWriteOnce), string(corev1.ReadOnlyMany), string(corev1.ReadWriteMany))
//
//var supportedReclaimPolicy = sets.NewString(string(corev1.PersistentVolumeReclaimDelete), string(corev1.PersistentVolumeReclaimRecycle), string(corev1.PersistentVolumeReclaimRetain))
//
//var supportedVolumeModes = sets.NewString(string(corev1.PersistentVolumeBlock), string(corev1.PersistentVolumeFilesystem))
//
//var supportedDataSourceAPIGroupKinds = map[schema.GroupKind]bool{
//	{Group: "snapshot.storage.k8s.io", Kind: "VolumeSnapshot"}: true,
//	{Group: "", Kind: "PersistentVolumeClaim"}:                 true,
//}
//
//func ValidatePersistentVolumeSpec(pvSpec *corev1.PersistentVolumeSpec, pvName string, validateInlinePersistentVolumeSpec bool, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	if validateInlinePersistentVolumeSpec {
//		if pvSpec.ClaimRef != nil {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("claimRef"), "may not be specified in the context of inline volumes"))
//		}
//		if len(pvSpec.Capacity) != 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("capacity"), "may not be specified in the context of inline volumes"))
//		}
//		if pvSpec.CSI == nil {
//			allErrs = append(allErrs, field.Required(fldPath.Child("csi"), "has to be specified in the context of inline volumes"))
//		}
//	}
//
//	if len(pvSpec.AccessModes) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("accessModes"), ""))
//	}
//	for _, mode := range pvSpec.AccessModes {
//		if !supportedAccessModes.Has(string(mode)) {
//			allErrs = append(allErrs, field.NotSupported(fldPath.Child("accessModes"), mode, supportedAccessModes.List()))
//		}
//	}
//
//	if !validateInlinePersistentVolumeSpec {
//		if len(pvSpec.Capacity) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("capacity"), ""))
//		}
//
//		if _, ok := pvSpec.Capacity[corev1.ResourceStorage]; !ok || len(pvSpec.Capacity) > 1 {
//			allErrs = append(allErrs, field.NotSupported(fldPath.Child("capacity"), pvSpec.Capacity, []string{string(corev1.ResourceStorage)}))
//		}
//		capPath := fldPath.Child("capacity")
//		for r, qty := range pvSpec.Capacity {
//			allErrs = append(allErrs, validateBasicResource(qty, capPath.Key(string(r)))...)
//			allErrs = append(allErrs, ValidatePositiveQuantityValue(qty, capPath.Key(string(r)))...)
//		}
//	}
//
//	if len(string(pvSpec.PersistentVolumeReclaimPolicy)) > 0 {
//		if validateInlinePersistentVolumeSpec {
//			if pvSpec.PersistentVolumeReclaimPolicy != corev1.PersistentVolumeReclaimRetain {
//				allErrs = append(allErrs, field.Forbidden(fldPath.Child("persistentVolumeReclaimPolicy"), "may only be "+string(corev1.PersistentVolumeReclaimRetain)+" in the context of inline volumes"))
//			}
//		} else {
//			if !supportedReclaimPolicy.Has(string(pvSpec.PersistentVolumeReclaimPolicy)) {
//				allErrs = append(allErrs, field.NotSupported(fldPath.Child("persistentVolumeReclaimPolicy"), pvSpec.PersistentVolumeReclaimPolicy, supportedReclaimPolicy.List()))
//			}
//		}
//	}
//
//	var nodeAffinitySpecified bool
//	var errs field.ErrorList
//	if pvSpec.NodeAffinity != nil {
//		if validateInlinePersistentVolumeSpec {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("nodeAffinity"), "may not be specified in the context of inline volumes"))
//		} else {
//			nodeAffinitySpecified, errs = validateVolumeNodeAffinity(pvSpec.NodeAffinity, fldPath.Child("nodeAffinity"))
//			allErrs = append(allErrs, errs...)
//		}
//	}
//	numVolumes := 0
//	if pvSpec.HostPath != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("hostPath"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateHostPathVolumeSource(pvSpec.HostPath, fldPath.Child("hostPath"))...)
//		}
//	}
//	if pvSpec.GCEPersistentDisk != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("gcePersistentDisk"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateGCEPersistentDiskVolumeSource(pvSpec.GCEPersistentDisk, fldPath.Child("persistentDisk"))...)
//		}
//	}
//	if pvSpec.AWSElasticBlockStore != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("awsElasticBlockStore"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateAWSElasticBlockStoreVolumeSource(pvSpec.AWSElasticBlockStore, fldPath.Child("awsElasticBlockStore"))...)
//		}
//	}
//	if pvSpec.Glusterfs != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("glusterfs"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateGlusterfsPersistentVolumeSource(pvSpec.Glusterfs, fldPath.Child("glusterfs"))...)
//		}
//	}
//	if pvSpec.Flocker != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("flocker"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateFlockerVolumeSource(pvSpec.Flocker, fldPath.Child("flocker"))...)
//		}
//	}
//	if pvSpec.NFS != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("nfs"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateNFSVolumeSource(pvSpec.NFS, fldPath.Child("nfs"))...)
//		}
//	}
//	if pvSpec.RBD != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("rbd"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateRBDPersistentVolumeSource(pvSpec.RBD, fldPath.Child("rbd"))...)
//		}
//	}
//	if pvSpec.Quobyte != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("quobyte"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateQuobyteVolumeSource(pvSpec.Quobyte, fldPath.Child("quobyte"))...)
//		}
//	}
//	if pvSpec.CephFS != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("cephFS"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateCephFSPersistentVolumeSource(pvSpec.CephFS, fldPath.Child("cephfs"))...)
//		}
//	}
//	if pvSpec.ISCSI != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("iscsi"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateISCSIPersistentVolumeSource(pvSpec.ISCSI, pvName, fldPath.Child("iscsi"))...)
//		}
//	}
//	if pvSpec.Cinder != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("cinder"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateCinderPersistentVolumeSource(pvSpec.Cinder, fldPath.Child("cinder"))...)
//		}
//	}
//	if pvSpec.FC != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("fc"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateFCVolumeSource(pvSpec.FC, fldPath.Child("fc"))...)
//		}
//	}
//	if pvSpec.FlexVolume != nil {
//		numVolumes++
//		allErrs = append(allErrs, validateFlexPersistentVolumeSource(pvSpec.FlexVolume, fldPath.Child("flexVolume"))...)
//	}
//	if pvSpec.AzureFile != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("azureFile"), "may not specify more than 1 volume type"))
//
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateAzureFilePV(pvSpec.AzureFile, fldPath.Child("azureFile"))...)
//		}
//	}
//
//	if pvSpec.VsphereVolume != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("vsphereVolume"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateVsphereVolumeSource(pvSpec.VsphereVolume, fldPath.Child("vsphereVolume"))...)
//		}
//	}
//	if pvSpec.PhotonPersistentDisk != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("photonPersistentDisk"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validatePhotonPersistentDiskVolumeSource(pvSpec.PhotonPersistentDisk, fldPath.Child("photonPersistentDisk"))...)
//		}
//	}
//	if pvSpec.PortworxVolume != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("portworxVolume"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validatePortworxVolumeSource(pvSpec.PortworxVolume, fldPath.Child("portworxVolume"))...)
//		}
//	}
//	if pvSpec.AzureDisk != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("azureDisk"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateAzureDisk(pvSpec.AzureDisk, fldPath.Child("azureDisk"))...)
//		}
//	}
//	if pvSpec.ScaleIO != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("scaleIO"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateScaleIOPersistentVolumeSource(pvSpec.ScaleIO, fldPath.Child("scaleIO"))...)
//		}
//	}
//	if pvSpec.Local != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("local"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateLocalVolumeSource(pvSpec.Local, fldPath.Child("local"))...)
//			// NodeAffinity is required
//			if !nodeAffinitySpecified {
//				allErrs = append(allErrs, field.Required(fldPath.Child("nodeAffinity"), "Local volume requires node affinity"))
//			}
//		}
//	}
//	if pvSpec.StorageOS != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("storageos"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateStorageOSPersistentVolumeSource(pvSpec.StorageOS, fldPath.Child("storageos"))...)
//		}
//	}
//
//	if pvSpec.CSI != nil {
//		if numVolumes > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("csi"), "may not specify more than 1 volume type"))
//		} else {
//			numVolumes++
//			allErrs = append(allErrs, validateCSIPersistentVolumeSource(pvSpec.CSI, fldPath.Child("csi"))...)
//		}
//	}
//
//	if numVolumes == 0 {
//		allErrs = append(allErrs, field.Required(fldPath, "must specify a volume type"))
//	}
//
//	// do not allow hostPath mounts of '/' to have a 'recycle' reclaim policy
//	if pvSpec.HostPath != nil && path.Clean(pvSpec.HostPath.Path) == "/" && pvSpec.PersistentVolumeReclaimPolicy == corev1.PersistentVolumeReclaimRecycle {
//		allErrs = append(allErrs, field.Forbidden(fldPath.Child("persistentVolumeReclaimPolicy"), "may not be 'recycle' for a hostPath mount of '/'"))
//	}
//
//	if len(pvSpec.StorageClassName) > 0 {
//		if validateInlinePersistentVolumeSpec {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("storageClassName"), "may not be specified in the context of inline volumes"))
//		} else {
//			for _, msg := range ValidateClassName(pvSpec.StorageClassName, false) {
//				allErrs = append(allErrs, field.Invalid(fldPath.Child("storageClassName"), pvSpec.StorageClassName, msg))
//			}
//		}
//	}
//	if pvSpec.VolumeMode != nil {
//		if validateInlinePersistentVolumeSpec {
//			if *pvSpec.VolumeMode != corev1.PersistentVolumeFilesystem {
//				allErrs = append(allErrs, field.Forbidden(fldPath.Child("volumeMode"), "may not specify volumeMode other than "+string(corev1.PersistentVolumeFilesystem)+" in the context of inline volumes"))
//			}
//		} else {
//			if !supportedVolumeModes.Has(string(*pvSpec.VolumeMode)) {
//				allErrs = append(allErrs, field.NotSupported(fldPath.Child("volumeMode"), *pvSpec.VolumeMode, supportedVolumeModes.List()))
//			}
//		}
//	}
//	return allErrs
//}
//
//func ValidatePersistentVolume(pv *corev1.PersistentVolume) field.ErrorList {
//	metaPath := field.NewPath("metadata")
//	allErrs := ValidateObjectMeta(&pv.ObjectMeta, false, ValidatePersistentVolumeName, metaPath)
//	allErrs = append(allErrs, ValidatePersistentVolumeSpec(&pv.Spec, pv.ObjectMeta.Name, false, field.NewPath("spec"))...)
//	return allErrs
//}
//
//// ValidatePersistentVolumeUpdate tests to see if the update is legal for an end user to make.
//// newPv is updated with fields that cannot be changed.
//func ValidatePersistentVolumeUpdate(newPv, oldPv *corev1.PersistentVolume) field.ErrorList {
//	allErrs := ValidatePersistentVolume(newPv)
//
//	// if oldPV does not have ControllerExpandSecretRef then allow it to be set
//	if (oldPv.Spec.CSI != nil && oldPv.Spec.CSI.ControllerExpandSecretRef == nil) &&
//		(newPv.Spec.CSI != nil && newPv.Spec.CSI.ControllerExpandSecretRef != nil) {
//		newPv = newPv.DeepCopy()
//		newPv.Spec.CSI.ControllerExpandSecretRef = nil
//	}
//
//	// PersistentVolumeSource should be immutable after creation.
//	if !apiequality.Semantic.DeepEqual(newPv.Spec.PersistentVolumeSource, oldPv.Spec.PersistentVolumeSource) {
//		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "persistentvolumesource"), "is immutable after creation"))
//	}
//	allErrs = append(allErrs, ValidateImmutableField(newPv.Spec.VolumeMode, oldPv.Spec.VolumeMode, field.NewPath("volumeMode"))...)
//
//	// Allow setting NodeAffinity if oldPv NodeAffinity was not set
//	if oldPv.Spec.NodeAffinity != nil {
//		allErrs = append(allErrs, ValidateImmutableField(newPv.Spec.NodeAffinity, oldPv.Spec.NodeAffinity, field.NewPath("nodeAffinity"))...)
//	}
//
//	return allErrs
//}
//
//// ValidatePersistentVolumeStatusUpdate tests to see if the status update is legal for an end user to make.
//// newPv is updated with fields that cannot be changed.
//func ValidatePersistentVolumeStatusUpdate(newPv, oldPv *corev1.PersistentVolume) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newPv.ObjectMeta, &oldPv.ObjectMeta, field.NewPath("metadata"))
//	if len(newPv.ResourceVersion) == 0 {
//		allErrs = append(allErrs, field.Required(field.NewPath("resourceVersion"), ""))
//	}
//	newPv.Spec = oldPv.Spec
//	return allErrs
//}
//
//// ValidatePersistentVolumeClaim validates a PersistentVolumeClaim
//func ValidatePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) field.ErrorList {
//	allErrs := ValidateObjectMeta(&pvc.ObjectMeta, true, ValidatePersistentVolumeName, field.NewPath("metadata"))
//	allErrs = append(allErrs, ValidatePersistentVolumeClaimSpec(&pvc.Spec, field.NewPath("spec"))...)
//	return allErrs
//}
//
//// ValidatePersistentVolumeClaimSpec validates a PersistentVolumeClaimSpec
//func ValidatePersistentVolumeClaimSpec(spec *corev1.PersistentVolumeClaimSpec, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(spec.AccessModes) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("accessModes"), "at least 1 access mode is required"))
//	}
//	if spec.Selector != nil {
//		allErrs = append(allErrs, metav1validation.ValidateLabelSelector(spec.Selector, fldPath.Child("selector"))...)
//	}
//	for _, mode := range spec.AccessModes {
//		if mode != corev1.ReadWriteOnce && mode != corev1.ReadOnlyMany && mode != corev1.ReadWriteMany {
//			allErrs = append(allErrs, field.NotSupported(fldPath.Child("accessModes"), mode, supportedAccessModes.List()))
//		}
//	}
//	storageValue, ok := spec.Resources.Requests[corev1.ResourceStorage]
//	if !ok {
//		allErrs = append(allErrs, field.Required(fldPath.Child("resources").Key(string(corev1.ResourceStorage)), ""))
//	} else {
//		allErrs = append(allErrs, ValidateResourceQuantityValue(string(corev1.ResourceStorage), storageValue, fldPath.Child("resources").Key(string(corev1.ResourceStorage)))...)
//		allErrs = append(allErrs, ValidatePositiveQuantityValue(storageValue, fldPath.Child("resources").Key(string(corev1.ResourceStorage)))...)
//	}
//
//	if spec.StorageClassName != nil && len(*spec.StorageClassName) > 0 {
//		for _, msg := range ValidateClassName(*spec.StorageClassName, false) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("storageClassName"), *spec.StorageClassName, msg))
//		}
//	}
//	if spec.VolumeMode != nil && !supportedVolumeModes.Has(string(*spec.VolumeMode)) {
//		allErrs = append(allErrs, field.NotSupported(fldPath.Child("volumeMode"), *spec.VolumeMode, supportedVolumeModes.List()))
//	}
//
//	if spec.DataSource != nil {
//		if len(spec.DataSource.Name) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("dataSource", "name"), ""))
//		}
//
//		groupKind := schema.GroupKind{Group: "", Kind: spec.DataSource.Kind}
//		if spec.DataSource.APIGroup != nil {
//			groupKind.Group = string(*spec.DataSource.APIGroup)
//		}
//		groupKindList := make([]string, 0, len(supportedDataSourceAPIGroupKinds))
//		for grp := range supportedDataSourceAPIGroupKinds {
//			groupKindList = append(groupKindList, grp.String())
//		}
//		if !supportedDataSourceAPIGroupKinds[groupKind] {
//			allErrs = append(allErrs, field.NotSupported(fldPath.Child("dataSource"), groupKind.String(), groupKindList))
//		}
//	}
//
//	return allErrs
//}
//
var supportedPortProtocols = sets.NewString(string(corev1.ProtocolTCP), string(corev1.ProtocolUDP), string(corev1.ProtocolSCTP))

//
func validateContainerPorts(ports []corev1.ContainerPort, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allNames := sets.String{}
	for i, port := range ports {
		idxPath := fldPath.Index(i)
		if len(port.Name) > 0 {
			if msgs := validation.IsValidPortName(port.Name); len(msgs) != 0 {
				for i = range msgs {
					allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), port.Name, msgs[i]))
				}
			} else if allNames.Has(port.Name) {
				allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), port.Name))
			} else {
				allNames.Insert(port.Name)
			}
		}
		if port.ContainerPort == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("containerPort"), ""))
		} else {
			for _, msg := range validation.IsValidPortNum(int(port.ContainerPort)) {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("containerPort"), port.ContainerPort, msg))
			}
		}
		if port.HostPort != 0 {
			for _, msg := range validation.IsValidPortNum(int(port.HostPort)) {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("hostPort"), port.HostPort, msg))
			}
		}
		if len(port.Protocol) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("protocol"), ""))
		} else if !supportedPortProtocols.Has(string(port.Protocol)) {
			allErrs = append(allErrs, field.NotSupported(idxPath.Child("protocol"), port.Protocol, supportedPortProtocols.List()))
		}
	}
	return allErrs
}

// ValidateEnv validates env vars
func ValidateEnv(vars []corev1.EnvVar, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, ev := range vars {
		idxPath := fldPath.Index(i)
		if len(ev.Name) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("name"), ""))
		} else {
			for _, msg := range validation.IsEnvVarName(ev.Name) {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), ev.Name, msg))
			}
		}
		allErrs = append(allErrs, validateEnvVarValueFrom(ev, idxPath.Child("valueFrom"))...)
	}
	return allErrs
}

func validateEnvVarValueFrom(ev corev1.EnvVar, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ev.ValueFrom == nil {
		return allErrs
	}

	numSources := 0

	if ev.ValueFrom.ConfigMapKeyRef != nil {
		numSources++
		allErrs = append(allErrs, validateConfigMapKeySelector(ev.ValueFrom.ConfigMapKeyRef, fldPath.Child("configMapKeyRef"))...)
	}
	if ev.ValueFrom.SecretKeyRef != nil {
		numSources++
		allErrs = append(allErrs, validateSecretKeySelector(ev.ValueFrom.SecretKeyRef, fldPath.Child("secretKeyRef"))...)
	}

	if numSources == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath, "", "must specify one of: `fieldRef`, `resourceFieldRef`, `configMapKeyRef` or `secretKeyRef`"))
	} else if len(ev.Value) != 0 {
		if numSources != 0 {
			allErrs = append(allErrs, field.Invalid(fldPath, "", "may not be specified when `value` is not empty"))
		}
	} else if numSources > 1 {
		allErrs = append(allErrs, field.Invalid(fldPath, "", "may not have more than one field specified at a time"))
	}

	return allErrs
}

//
//func validateContainerResourceFieldSelector(fs *corev1.ResourceFieldSelector, expressions *sets.String, fldPath *field.Path, volume bool) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	if volume && len(fs.ContainerName) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("containerName"), ""))
//	} else if len(fs.Resource) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("resource"), ""))
//	} else if !expressions.Has(fs.Resource) {
//		allErrs = append(allErrs, field.NotSupported(fldPath.Child("resource"), fs.Resource, expressions.List()))
//	}
//	allErrs = append(allErrs, validateContainerResourceDivisor(fs.Resource, fs.Divisor, fldPath)...)
//	return allErrs
//}

func ValidateEnvFrom(vars []corev1.EnvFromSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, ev := range vars {
		idxPath := fldPath.Index(i)
		if len(ev.Prefix) > 0 {
			for _, msg := range validation.IsEnvVarName(ev.Prefix) {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("prefix"), ev.Prefix, msg))
			}
		}

		numSources := 0
		if ev.ConfigMapRef != nil {
			numSources++
			allErrs = append(allErrs, validateConfigMapEnvSource(ev.ConfigMapRef, idxPath.Child("configMapRef"))...)
		}
		if ev.SecretRef != nil {
			numSources++
			allErrs = append(allErrs, validateSecretEnvSource(ev.SecretRef, idxPath.Child("secretRef"))...)
		}

		if numSources == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath, "", "must specify one of: `configMapRef` or `secretRef`"))
		} else if numSources > 1 {
			allErrs = append(allErrs, field.Invalid(fldPath, "", "may not have more than one field specified at a time"))
		}
	}
	return allErrs
}

func validateConfigMapEnvSource(configMapSource *corev1.ConfigMapEnvSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(configMapSource.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	} else {
		for _, msg := range ValidateConfigMapName(configMapSource.Name, true) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), configMapSource.Name, msg))
		}
	}
	return allErrs
}

func validateSecretEnvSource(secretSource *corev1.SecretEnvSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(secretSource.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	} else {
		for _, msg := range ValidateSecretName(secretSource.Name, true) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), secretSource.Name, msg))
		}
	}
	return allErrs
}

//
//var validContainerResourceDivisorForCPU = sets.NewString("1m", "1")
//var validContainerResourceDivisorForMemory = sets.NewString("1", "1k", "1M", "1G", "1T", "1P", "1E", "1Ki", "1Mi", "1Gi", "1Ti", "1Pi", "1Ei")
//var validContainerResourceDivisorForEphemeralStorage = sets.NewString("1", "1k", "1M", "1G", "1T", "1P", "1E", "1Ki", "1Mi", "1Gi", "1Ti", "1Pi", "1Ei")
//
//func validateContainerResourceDivisor(rName string, divisor resource.Quantity, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	unsetDivisor := resource.Quantity{}
//	if unsetDivisor.Cmp(divisor) == 0 {
//		return allErrs
//	}
//	switch rName {
//	case "limits.cpu", "requests.cpu":
//		if !validContainerResourceDivisorForCPU.Has(divisor.String()) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("divisor"), rName, "only divisor's values 1m and 1 are supported with the cpu resource"))
//		}
//	case "limits.memory", "requests.memory":
//		if !validContainerResourceDivisorForMemory.Has(divisor.String()) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("divisor"), rName, "only divisor's values 1, 1k, 1M, 1G, 1T, 1P, 1E, 1Ki, 1Mi, 1Gi, 1Ti, 1Pi, 1Ei are supported with the memory resource"))
//		}
//	case "limits.ephemeral-storage", "requests.ephemeral-storage":
//		if !validContainerResourceDivisorForEphemeralStorage.Has(divisor.String()) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("divisor"), rName, "only divisor's values 1, 1k, 1M, 1G, 1T, 1P, 1E, 1Ki, 1Mi, 1Gi, 1Ti, 1Pi, 1Ei are supported with the local ephemeral storage resource"))
//		}
//	}
//	return allErrs
//}

func validateConfigMapKeySelector(s *corev1.ConfigMapKeySelector, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	nameFn := ValidateNameFunc(ValidateSecretName)
	for _, msg := range nameFn(s.Name, false) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), s.Name, msg))
	}
	if len(s.Key) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), ""))
	} else {
		for _, msg := range validation.IsConfigMapKey(s.Key) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("key"), s.Key, msg))
		}
	}

	return allErrs
}

func validateSecretKeySelector(s *corev1.SecretKeySelector, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	nameFn := ValidateNameFunc(ValidateSecretName)
	for _, msg := range nameFn(s.Name, false) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), s.Name, msg))
	}
	if len(s.Key) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), ""))
	} else {
		for _, msg := range validation.IsConfigMapKey(s.Key) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("key"), s.Key, msg))
		}
	}

	return allErrs
}

func GetVolumeMountMap(mounts []corev1.VolumeMount) map[string]string {
	volmounts := make(map[string]string)

	for _, mnt := range mounts {
		volmounts[mnt.Name] = mnt.MountPath
	}

	return volmounts
}

//
func GetVolumeDeviceMap(devices []corev1.VolumeDevice) map[string]string {
	voldevices := make(map[string]string)

	for _, dev := range devices {
		voldevices[dev.Name] = dev.DevicePath
	}

	return voldevices
}

func ValidateVolumeMounts(mounts []corev1.VolumeMount, voldevices map[string]string, volumes map[string]corev1.VolumeSource, container *corev1.Container, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	mountpoints := sets.NewString()

	for i, mnt := range mounts {
		idxPath := fldPath.Index(i)
		if len(mnt.Name) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("name"), ""))
		}
		//if !IsMatchedVolume(mnt.Name, volumes) {
		//	allErrs = append(allErrs, field.NotFound(idxPath.Child("name"), mnt.Name))
		//}
		if len(mnt.MountPath) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("mountPath"), ""))
		}
		if mountpoints.Has(mnt.MountPath) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("mountPath"), mnt.MountPath, "must be unique"))
		}
		mountpoints.Insert(mnt.MountPath)

		// check for overlap with VolumeDevice
		if mountNameAlreadyExists(mnt.Name, voldevices) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), mnt.Name, "must not already exist in volumeDevices"))
		}
		if mountPathAlreadyExists(mnt.MountPath, voldevices) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("mountPath"), mnt.MountPath, "must not already exist as a path in volumeDevices"))
		}

		if len(mnt.SubPath) > 0 {
			allErrs = append(allErrs, validateLocalDescendingPath(mnt.SubPath, fldPath.Child("subPath"))...)
		}

		if len(mnt.SubPathExpr) > 0 {
			if len(mnt.SubPath) > 0 {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("subPathExpr"), mnt.SubPathExpr, "subPathExpr and subPath are mutually exclusive"))
			}

			allErrs = append(allErrs, validateLocalDescendingPath(mnt.SubPathExpr, fldPath.Child("subPathExpr"))...)
		}

		if mnt.MountPropagation != nil {
			allErrs = append(allErrs, validateMountPropagation(mnt.MountPropagation, container, fldPath.Child("mountPropagation"))...)
		}
	}
	return allErrs
}

func ValidateVolumeDevices(devices []corev1.VolumeDevice, volmounts map[string]string, volumes map[string]corev1.VolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	devicepath := sets.NewString()
	devicename := sets.NewString()

	for i, dev := range devices {
		idxPath := fldPath.Index(i)
		devName := dev.Name
		devPath := dev.DevicePath
		didMatch, isPVC := isMatchedDevice(devName, volumes)
		if len(devName) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("name"), ""))
		}
		if devicename.Has(devName) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), devName, "must be unique"))
		}
		// Must be PersistentVolumeClaim volume source
		if didMatch && !isPVC {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), devName, "can only use volume source type of PersistentVolumeClaim for block mode"))
		}
		if !didMatch {
			allErrs = append(allErrs, field.NotFound(idxPath.Child("name"), devName))
		}
		if len(devPath) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("devicePath"), ""))
		}
		if devicepath.Has(devPath) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("devicePath"), devPath, "must be unique"))
		}
		if len(devPath) > 0 && len(validatePathNoBacksteps(devPath, fldPath.Child("devicePath"))) > 0 {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("devicePath"), devPath, "can not contain backsteps ('..')"))
		} else {
			devicepath.Insert(devPath)
		}
		// check for overlap with VolumeMount
		if deviceNameAlreadyExists(devName, volmounts) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), devName, "must not already exist in volumeMounts"))
		}
		if devicePathAlreadyExists(devPath, volmounts) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("devicePath"), devPath, "must not already exist as a path in volumeMounts"))
		}
		if len(devName) > 0 {
			devicename.Insert(devName)
		}
	}
	return allErrs
}

func validateProbe(probe *corev1.Probe, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if probe == nil {
		return allErrs
	}
	allErrs = append(allErrs, validateHandler(&probe.Handler, fldPath)...)

	allErrs = append(allErrs, ValidateNonnegativeField(int64(probe.InitialDelaySeconds), fldPath.Child("initialDelaySeconds"))...)
	allErrs = append(allErrs, ValidateNonnegativeField(int64(probe.TimeoutSeconds), fldPath.Child("timeoutSeconds"))...)
	allErrs = append(allErrs, ValidateNonnegativeField(int64(probe.PeriodSeconds), fldPath.Child("periodSeconds"))...)
	allErrs = append(allErrs, ValidateNonnegativeField(int64(probe.SuccessThreshold), fldPath.Child("successThreshold"))...)
	allErrs = append(allErrs, ValidateNonnegativeField(int64(probe.FailureThreshold), fldPath.Child("failureThreshold"))...)
	return allErrs
}

//
//func validateClientIPAffinityConfig(config *corev1.SessionAffinityConfig, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if config == nil {
//		allErrs = append(allErrs, field.Required(fldPath, fmt.Sprintf("when session affinity type is %s", corev1.ServiceAffinityClientIP)))
//		return allErrs
//	}
//	if config.ClientIP == nil {
//		allErrs = append(allErrs, field.Required(fldPath.Child("clientIP"), fmt.Sprintf("when session affinity type is %s", corev1.ServiceAffinityClientIP)))
//		return allErrs
//	}
//	if config.ClientIP.TimeoutSeconds == nil {
//		allErrs = append(allErrs, field.Required(fldPath.Child("clientIP").Child("timeoutSeconds"), fmt.Sprintf("when session affinity type is %s", corev1.ServiceAffinityClientIP)))
//		return allErrs
//	}
//	allErrs = append(allErrs, validateAffinityTimeout(config.ClientIP.TimeoutSeconds, fldPath.Child("clientIP").Child("timeoutSeconds"))...)
//
//	return allErrs
//}
//
//func validateAffinityTimeout(timeout *int32, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if *timeout <= 0 || *timeout > corev1.MaxClientIPServiceAffinitySeconds {
//		allErrs = append(allErrs, field.Invalid(fldPath, timeout, fmt.Sprintf("must be greater than 0 and less than %d", corev1.MaxClientIPServiceAffinitySeconds)))
//	}
//	return allErrs
//}

// AccumulateUniqueHostPorts extracts each HostPort of each Container,
// accumulating the results and returning an error if any ports conflict.
func AccumulateUniqueHostPorts(containers []corev1.Container, accumulator *sets.String, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for ci, ctr := range containers {
		idxPath := fldPath.Index(ci)
		portsPath := idxPath.Child("ports")
		for pi := range ctr.Ports {
			idxPath := portsPath.Index(pi)
			port := ctr.Ports[pi].HostPort
			if port == 0 {
				continue
			}
			str := fmt.Sprintf("%s/%s/%d", ctr.Ports[pi].Protocol, ctr.Ports[pi].HostIP, port)
			if accumulator.Has(str) {
				allErrs = append(allErrs, field.Duplicate(idxPath.Child("hostPort"), str))
			} else {
				accumulator.Insert(str)
			}
		}
	}
	return allErrs
}

// checkHostPortConflicts checks for colliding Port.HostPort values across
// a slice of containers.
func checkHostPortConflicts(containers []corev1.Container, fldPath *field.Path) field.ErrorList {
	allPorts := sets.String{}
	return AccumulateUniqueHostPorts(containers, &allPorts, fldPath)
}

func validateExecAction(exec *corev1.ExecAction, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	if len(exec.Command) == 0 {
		allErrors = append(allErrors, field.Required(fldPath.Child("command"), ""))
	}
	return allErrors
}

var supportedHTTPSchemes = sets.NewString(string(corev1.URISchemeHTTP), string(corev1.URISchemeHTTPS))

func validateHTTPGetAction(http *corev1.HTTPGetAction, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	if len(http.Path) == 0 {
		allErrors = append(allErrors, field.Required(fldPath.Child("path"), ""))
	}
	allErrors = append(allErrors, ValidatePortNumOrName(http.Port, fldPath.Child("port"))...)
	if !supportedHTTPSchemes.Has(string(http.Scheme)) {
		allErrors = append(allErrors, field.NotSupported(fldPath.Child("scheme"), http.Scheme, supportedHTTPSchemes.List()))
	}
	for _, header := range http.HTTPHeaders {
		for _, msg := range validation.IsHTTPHeaderName(header.Name) {
			allErrors = append(allErrors, field.Invalid(fldPath.Child("httpHeaders"), header.Name, msg))
		}
	}
	return allErrors
}

func ValidatePortNumOrName(port intstr.IntOrString, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if port.Type == intstr.Int {
		for _, msg := range validation.IsValidPortNum(port.IntValue()) {
			allErrs = append(allErrs, field.Invalid(fldPath, port.IntValue(), msg))
		}
	} else if port.Type == intstr.String {
		for _, msg := range validation.IsValidPortName(port.StrVal) {
			allErrs = append(allErrs, field.Invalid(fldPath, port.StrVal, msg))
		}
	} else {
		allErrs = append(allErrs, field.InternalError(fldPath, fmt.Errorf("unknown type: %v", port.Type)))
	}
	return allErrs
}

func validateTCPSocketAction(tcp *corev1.TCPSocketAction, fldPath *field.Path) field.ErrorList {
	return ValidatePortNumOrName(tcp.Port, fldPath.Child("port"))
}

func validateHandler(handler *corev1.Handler, fldPath *field.Path) field.ErrorList {
	numHandlers := 0
	allErrors := field.ErrorList{}
	if handler.Exec != nil {
		if numHandlers > 0 {
			allErrors = append(allErrors, field.Forbidden(fldPath.Child("exec"), "may not specify more than 1 handler type"))
		} else {
			numHandlers++
			allErrors = append(allErrors, validateExecAction(handler.Exec, fldPath.Child("exec"))...)
		}
	}
	if handler.HTTPGet != nil {
		if numHandlers > 0 {
			allErrors = append(allErrors, field.Forbidden(fldPath.Child("httpGet"), "may not specify more than 1 handler type"))
		} else {
			numHandlers++
			allErrors = append(allErrors, validateHTTPGetAction(handler.HTTPGet, fldPath.Child("httpGet"))...)
		}
	}
	if handler.TCPSocket != nil {
		if numHandlers > 0 {
			allErrors = append(allErrors, field.Forbidden(fldPath.Child("tcpSocket"), "may not specify more than 1 handler type"))
		} else {
			numHandlers++
			allErrors = append(allErrors, validateTCPSocketAction(handler.TCPSocket, fldPath.Child("tcpSocket"))...)
		}
	}
	if numHandlers == 0 {
		allErrors = append(allErrors, field.Required(fldPath, "must specify a handler type"))
	}
	return allErrors
}

func validateLifecycle(lifecycle *corev1.Lifecycle, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if lifecycle.PostStart != nil {
		allErrs = append(allErrs, validateHandler(lifecycle.PostStart, fldPath.Child("postStart"))...)
	}
	if lifecycle.PreStop != nil {
		allErrs = append(allErrs, validateHandler(lifecycle.PreStop, fldPath.Child("preStop"))...)
	}
	return allErrs
}

var supportedPullPolicies = sets.NewString(string(corev1.PullAlways), string(corev1.PullIfNotPresent), string(corev1.PullNever))

func validatePullPolicy(policy corev1.PullPolicy, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}

	switch policy {
	case corev1.PullAlways, corev1.PullIfNotPresent, corev1.PullNever:
		break
	case "":
		allErrors = append(allErrors, field.Required(fldPath, ""))
	default:
		allErrors = append(allErrors, field.NotSupported(fldPath, policy, supportedPullPolicies.List()))
	}

	return allErrors
}

func validateContainers(containers []corev1.Container, isInitContainers bool, volumes map[string]corev1.VolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(containers) == 0 {
		return append(allErrs, field.Required(fldPath, ""))
	}

	allNames := sets.String{}
	for i, ctr := range containers {
		idxPath := fldPath.Index(i)
		namePath := idxPath.Child("name")
		volMounts := GetVolumeMountMap(ctr.VolumeMounts)
		volDevices := GetVolumeDeviceMap(ctr.VolumeDevices)

		if len(ctr.Name) == 0 {
			allErrs = append(allErrs, field.Required(namePath, ""))
		} else {
			allErrs = append(allErrs, ValidateDNS1123Label(ctr.Name, namePath)...)
		}
		if allNames.Has(ctr.Name) {
			allErrs = append(allErrs, field.Duplicate(namePath, ctr.Name))
		} else {
			allNames.Insert(ctr.Name)
		}
		// TODO: do not validate leading and trailing whitespace to preserve backward compatibility.
		// for example: https://github.com/openshift/origin/issues/14659 image = " " is special token in pod template
		// others may have done similar
		if len(ctr.Image) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("image"), ""))
		}
		if ctr.Lifecycle != nil {
			allErrs = append(allErrs, validateLifecycle(ctr.Lifecycle, idxPath.Child("lifecycle"))...)
		}
		allErrs = append(allErrs, validateProbe(ctr.LivenessProbe, idxPath.Child("livenessProbe"))...)
		// Liveness-specific validation
		if ctr.LivenessProbe != nil && ctr.LivenessProbe.SuccessThreshold != 1 {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("livenessProbe", "successThreshold"), ctr.LivenessProbe.SuccessThreshold, "must be 1"))
		}

		//switch ctr.TerminationMessagePolicy {
		//case corev1.TerminationMessageReadFile, corev1.TerminationMessageFallbackToLogsOnError:
		//case "":
		//	allErrs = append(allErrs, field.Required(idxPath.Child("terminationMessagePolicy"), "must be 'File' or 'FallbackToLogsOnError'"))
		//default:
		//	allErrs = append(allErrs, field.Invalid(idxPath.Child("terminationMessagePolicy"), ctr.TerminationMessagePolicy, "must be 'File' or 'FallbackToLogsOnError'"))
		//}

		allErrs = append(allErrs, validateProbe(ctr.ReadinessProbe, idxPath.Child("readinessProbe"))...)
		allErrs = append(allErrs, validateContainerPorts(ctr.Ports, idxPath.Child("ports"))...)
		allErrs = append(allErrs, ValidateEnv(ctr.Env, idxPath.Child("env"))...)
		allErrs = append(allErrs, ValidateEnvFrom(ctr.EnvFrom, idxPath.Child("envFrom"))...)
		allErrs = append(allErrs, ValidateVolumeMounts(ctr.VolumeMounts, volDevices, volumes, &ctr, idxPath.Child("volumeMounts"))...)
		allErrs = append(allErrs, ValidateVolumeDevices(ctr.VolumeDevices, volMounts, volumes, idxPath.Child("volumeDevices"))...)
		//allErrs = append(allErrs, validatePullPolicy(ctr.ImagePullPolicy, idxPath.Child("imagePullPolicy"))...)
		allErrs = append(allErrs, ValidateResourceRequirements(&ctr.Resources, idxPath.Child("resources"))...)
		allErrs = append(allErrs, ValidateSecurityContext(ctr.SecurityContext, idxPath.Child("securityContext"))...)
	}

	if isInitContainers {
		// check initContainers one by one since they are running in sequential order.
		for _, initContainer := range containers {
			allErrs = append(allErrs, checkHostPortConflicts([]corev1.Container{initContainer}, fldPath)...)
		}
	} else {
		// Check for colliding ports across all containers.
		allErrs = append(allErrs, checkHostPortConflicts(containers, fldPath)...)
	}

	return allErrs
}

func validateInitContainers(containers, otherContainers []corev1.Container, deviceVolumes map[string]corev1.VolumeSource, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if len(containers) > 0 {
		allErrs = append(allErrs, validateContainers(containers, true, deviceVolumes, fldPath)...)
	}

	allNames := sets.String{}
	for _, ctr := range otherContainers {
		allNames.Insert(ctr.Name)
	}
	for i, ctr := range containers {
		idxPath := fldPath.Index(i)
		if allNames.Has(ctr.Name) {
			allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), ctr.Name))
		}
		if len(ctr.Name) > 0 {
			allNames.Insert(ctr.Name)
		}
		if ctr.Lifecycle != nil {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("lifecycle"), ctr.Lifecycle, "must not be set for init containers"))
		}
		if ctr.LivenessProbe != nil {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("livenessProbe"), ctr.LivenessProbe, "must not be set for init containers"))
		}
		if ctr.ReadinessProbe != nil {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("readinessProbe"), ctr.ReadinessProbe, "must not be set for init containers"))
		}
	}
	return allErrs
}

func validateRestartPolicy(restartPolicy *corev1.RestartPolicy, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	switch *restartPolicy {
	case corev1.RestartPolicyAlways, corev1.RestartPolicyOnFailure, corev1.RestartPolicyNever:
		break
	case "":
		allErrors = append(allErrors, field.Required(fldPath, ""))
	default:
		validValues := []string{string(corev1.RestartPolicyAlways), string(corev1.RestartPolicyOnFailure), string(corev1.RestartPolicyNever)}
		allErrors = append(allErrors, field.NotSupported(fldPath, *restartPolicy, validValues))
	}

	return allErrors
}

func ValidatePreemptionPolicy(preemptionPolicy *corev1.PreemptionPolicy, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	switch *preemptionPolicy {
	case corev1.PreemptLowerPriority, corev1.PreemptNever:
	case "":
		allErrors = append(allErrors, field.Required(fldPath, ""))
	default:
		validValues := []string{string(corev1.PreemptLowerPriority), string(corev1.PreemptNever)}
		allErrors = append(allErrors, field.NotSupported(fldPath, preemptionPolicy, validValues))
	}
	return allErrors
}

func validateDNSPolicy(dnsPolicy *corev1.DNSPolicy, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	switch *dnsPolicy {
	case corev1.DNSClusterFirstWithHostNet, corev1.DNSClusterFirst, corev1.DNSDefault, corev1.DNSNone:
	case "":
		allErrors = append(allErrors, field.Required(fldPath, ""))
	default:
		validValues := []string{string(corev1.DNSClusterFirstWithHostNet), string(corev1.DNSClusterFirst), string(corev1.DNSDefault), string(corev1.DNSNone)}
		allErrors = append(allErrors, field.NotSupported(fldPath, dnsPolicy, validValues))
	}
	return allErrors
}

//const (
//	// Limits on various DNS parameters. These are derived from
//	// restrictions in Linux libc name resolution handling.
//	// Max number of DNS name servers.
//	MaxDNSNameservers = 3
//	// Max number of domains in search path.
//	MaxDNSSearchPaths = 6
//	// Max number of characters in search path.
//	MaxDNSSearchListChars = 256
//)
//
//func validateReadinessGates(readinessGates []corev1.PodReadinessGate, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	for i, value := range readinessGates {
//		for _, msg := range validation.IsQualifiedName(string(value.ConditionType)) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Index(i).Child("conditionType"), string(value.ConditionType), msg))
//		}
//	}
//	return allErrs
//}

//func validateHostNetwork(hostNetwork bool, containers []corev1.Container, fldPath *field.Path) field.ErrorList {
//	allErrors := field.ErrorList{}
//	if hostNetwork {
//		for i, container := range containers {
//			portsPath := fldPath.Index(i).Child("ports")
//			for i, port := range container.Ports {
//				idxPath := portsPath.Index(i)
//				if port.HostPort != port.ContainerPort {
//					allErrors = append(allErrors, field.Invalid(idxPath.Child("containerPort"), port.ContainerPort, "must match `hostPort` when `hostNetwork` is true"))
//				}
//			}
//		}
//	}
//	return allErrors
//}

// validateImagePullSecrets checks to make sure the pull secrets are well
// formed.  Right now, we only expect name to be set (it's the only field).  If
// this ever changes and someone decides to set those fields, we'd like to
// know.
func validateImagePullSecrets(imagePullSecrets []corev1.LocalObjectReference, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	for i, currPullSecret := range imagePullSecrets {
		idxPath := fldPath.Index(i)
		strippedRef := corev1.LocalObjectReference{Name: currPullSecret.Name}
		if !reflect.DeepEqual(strippedRef, currPullSecret) {
			allErrors = append(allErrors, field.Invalid(idxPath, currPullSecret, "only name may be set"))
		}
	}
	return allErrors
}

func validateTaintEffect(effect *corev1.TaintEffect, allowEmpty bool, fldPath *field.Path) field.ErrorList {
	if !allowEmpty && len(*effect) == 0 {
		return field.ErrorList{field.Required(fldPath, "")}
	}

	allErrors := field.ErrorList{}
	switch *effect {
	// TODO: Replace next line with subsequent commented-out line when implement TaintEffectNoScheduleNoAdmit.
	case corev1.TaintEffectNoSchedule, corev1.TaintEffectPreferNoSchedule, corev1.TaintEffectNoExecute:
		// case corev1.TaintEffectNoSchedule, corev1.TaintEffectPreferNoSchedule, corev1.TaintEffectNoScheduleNoAdmit, corev1.TaintEffectNoExecute:
	default:
		validValues := []string{
			string(corev1.TaintEffectNoSchedule),
			string(corev1.TaintEffectPreferNoSchedule),
			string(corev1.TaintEffectNoExecute),
			// TODO: Uncomment this block when implement TaintEffectNoScheduleNoAdmit.
			// string(corev1.TaintEffectNoScheduleNoAdmit),
		}
		allErrors = append(allErrors, field.NotSupported(fldPath, *effect, validValues))
	}
	return allErrors
}

//// validateOnlyAddedTolerations validates updated pod tolerations.
//func validateOnlyAddedTolerations(newTolerations []corev1.Toleration, oldTolerations []corev1.Toleration, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	for _, old := range oldTolerations {
//		found := false
//		old.TolerationSeconds = nil
//		for _, new := range newTolerations {
//			new.TolerationSeconds = nil
//			if reflect.DeepEqual(old, new) {
//				found = true
//				break
//			}
//		}
//		if !found {
//			allErrs = append(allErrs, field.Forbidden(fldPath, "existing toleration can not be modified except its tolerationSeconds"))
//			return allErrs
//		}
//	}
//
//	allErrs = append(allErrs, ValidateTolerations(newTolerations, fldPath)...)
//	return allErrs
//}

func ValidateHostAliases(hostAliases []corev1.HostAlias, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, hostAlias := range hostAliases {
		if ip := net.ParseIP(hostAlias.IP); ip == nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ip"), hostAlias.IP, "must be valid IP address"))
		}
		for _, hostname := range hostAlias.Hostnames {
			allErrs = append(allErrs, ValidateDNS1123Subdomain(hostname, fldPath.Child("hostnames"))...)
		}
	}
	return allErrs
}

// ValidateTolerations tests if given tolerations have valid data.
func ValidateTolerations(tolerations []corev1.Toleration, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	for i, toleration := range tolerations {
		idxPath := fldPath.Index(i)
		// validate the toleration key
		if len(toleration.Key) > 0 {
			allErrors = append(allErrors, metav1validation.ValidateLabelName(toleration.Key, idxPath.Child("key"))...)
		}

		// empty toleration key with Exists operator and empty value means match all taints
		if len(toleration.Key) == 0 && toleration.Operator != corev1.TolerationOpExists {
			allErrors = append(allErrors, field.Invalid(idxPath.Child("operator"), toleration.Operator,
				"operator must be Exists when `key` is empty, which means \"match all values and all keys\""))
		}

		if toleration.TolerationSeconds != nil && toleration.Effect != corev1.TaintEffectNoExecute {
			allErrors = append(allErrors, field.Invalid(idxPath.Child("effect"), toleration.Effect,
				"effect must be 'NoExecute' when `tolerationSeconds` is set"))
		}

		// validate toleration operator and value
		switch toleration.Operator {
		// empty operator means Equal
		case corev1.TolerationOpEqual, "":
			if errs := validation.IsValidLabelValue(toleration.Value); len(errs) != 0 {
				allErrors = append(allErrors, field.Invalid(idxPath.Child("operator"), toleration.Value, strings.Join(errs, ";")))
			}
		case corev1.TolerationOpExists:
			if len(toleration.Value) > 0 {
				allErrors = append(allErrors, field.Invalid(idxPath.Child("operator"), toleration, "value must be empty when `operator` is 'Exists'"))
			}
		default:
			validValues := []string{string(corev1.TolerationOpEqual), string(corev1.TolerationOpExists)}
			allErrors = append(allErrors, field.NotSupported(idxPath.Child("operator"), toleration.Operator, validValues))
		}

		// validate toleration effect, empty toleration effect means match all taint effects
		if len(toleration.Effect) > 0 {
			allErrors = append(allErrors, validateTaintEffect(&toleration.Effect, true, idxPath.Child("effect"))...)
		}
	}
	return allErrors
}

//// ValidatePodSpec tests that the specified PodSpec has valid data.
//// This includes checking formatting and uniqueness.  It also canonicalizes the
//// structure by setting default values and implementing any backwards-compatibility
//// tricks.
//
//// ValidateNodeSelectorRequirement tests that the specified NodeSelectorRequirement fields has valid data
//func ValidateNodeSelectorRequirement(rq corev1.NodeSelectorRequirement, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	switch rq.Operator {
//	case corev1.NodeSelectorOpIn, corev1.NodeSelectorOpNotIn:
//		if len(rq.Values) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("values"), "must be specified when `operator` is 'In' or 'NotIn'"))
//		}
//	case corev1.NodeSelectorOpExists, corev1.NodeSelectorOpDoesNotExist:
//		if len(rq.Values) > 0 {
//			allErrs = append(allErrs, field.Forbidden(fldPath.Child("values"), "may not be specified when `operator` is 'Exists' or 'DoesNotExist'"))
//		}
//
//	case corev1.NodeSelectorOpGt, corev1.NodeSelectorOpLt:
//		if len(rq.Values) != 1 {
//			allErrs = append(allErrs, field.Required(fldPath.Child("values"), "must be specified single value when `operator` is 'Lt' or 'Gt'"))
//		}
//	default:
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("operator"), rq.Operator, "not a valid selector operator"))
//	}
//
//	allErrs = append(allErrs, metav1validation.ValidateLabelName(rq.Key, fldPath.Child("key"))...)
//
//	return allErrs
//}
//
//func ValidateSeccompProfile(p string, fldPath *field.Path) field.ErrorList {
//	if p == corev1.SeccompProfileRuntimeDefault || p == corev1.DeprecatedSeccompProfileDockerDefault {
//		return nil
//	}
//	if p == "unconfined" {
//		return nil
//	}
//	if strings.HasPrefix(p, "localhost/") {
//		return validateLocalDescendingPath(strings.TrimPrefix(p, "localhost/"), fldPath)
//	}
//	return field.ErrorList{field.Invalid(fldPath, p, "must be a valid seccomp profile")}
//}
//
//func ValidateSeccompPodAnnotations(annotations map[string]string, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if p, exists := annotations[corev1.SeccompPodAnnotationKey]; exists {
//		allErrs = append(allErrs, ValidateSeccompProfile(p, fldPath.Child(corev1.SeccompPodAnnotationKey))...)
//	}
//	for k, p := range annotations {
//		if strings.HasPrefix(k, corev1.SeccompContainerAnnotationKeyPrefix) {
//			allErrs = append(allErrs, ValidateSeccompProfile(p, fldPath.Child(k))...)
//		}
//	}
//
//	return allErrs
//}

const (
	// a sysctl segment regex, concatenated with dots to form a sysctl name
	SysctlSegmentFmt string = "[a-z0-9]([-_a-z0-9]*[a-z0-9])?"

	// a sysctl name regex
	SysctlFmt string = "(" + SysctlSegmentFmt + "\\.)*" + SysctlSegmentFmt

	// the maximal length of a sysctl name
	SysctlMaxLength int = 253
)

var sysctlRegexp = regexp.MustCompile("^" + SysctlFmt + "$")

// IsValidSysctlName checks that the given string is a valid sysctl name,
// i.e. matches SysctlFmt.
func IsValidSysctlName(name string) bool {
	if len(name) > SysctlMaxLength {
		return false
	}
	return sysctlRegexp.MatchString(name)
}

func validateSysctls(sysctls []corev1.Sysctl, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	names := make(map[string]struct{})
	for i, s := range sysctls {
		if len(s.Name) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("name"), ""))
		} else if !IsValidSysctlName(s.Name) {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i).Child("name"), s.Name, fmt.Sprintf("must have at most %d characters and match regex %s", SysctlMaxLength, SysctlFmt)))
		} else if _, ok := names[s.Name]; ok {
			allErrs = append(allErrs, field.Duplicate(fldPath.Index(i).Child("name"), s.Name))
		}
		names[s.Name] = struct{}{}
	}
	return allErrs
}

// ValidatePodSecurityContext test that the specified PodSecurityContext has valid data.
func ValidatePodSecurityContext(securityContext *corev1.PodSecurityContext, spec *corev1.PodSpec, specPath, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if securityContext != nil {
		if securityContext.FSGroup != nil {
			for _, msg := range validation.IsValidGroupID(*securityContext.FSGroup) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("fsGroup"), *(securityContext.FSGroup), msg))
			}
		}
		if securityContext.RunAsUser != nil {
			for _, msg := range validation.IsValidUserID(*securityContext.RunAsUser) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("runAsUser"), *(securityContext.RunAsUser), msg))
			}
		}
		if securityContext.RunAsGroup != nil {
			for _, msg := range validation.IsValidGroupID(*securityContext.RunAsGroup) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("runAsGroup"), *(securityContext.RunAsGroup), msg))
			}
		}
		for g, gid := range securityContext.SupplementalGroups {
			for _, msg := range validation.IsValidGroupID(gid) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("supplementalGroups").Index(g), gid, msg))
			}
		}

		if len(securityContext.Sysctls) != 0 {
			allErrs = append(allErrs, validateSysctls(securityContext.Sysctls, fldPath.Child("sysctls"))...)
		}

		allErrs = append(allErrs, validateWindowsSecurityContextOptions(securityContext.WindowsOptions, fldPath.Child("windowsOptions"))...)
	}

	return allErrs
}

//func ValidateContainerUpdates(newContainers, oldContainers []corev1.Container, fldPath *field.Path) (allErrs field.ErrorList, stop bool) {
//	allErrs = field.ErrorList{}
//	if len(newContainers) != len(oldContainers) {
//		//TODO: Pinpoint the specific container that causes the invalid error after we have strategic merge diff
//		allErrs = append(allErrs, field.Forbidden(fldPath, "pod updates may not add or remove containers"))
//		return allErrs, true
//	}
//
//	// validate updated container images
//	for i, ctr := range newContainers {
//		if len(ctr.Image) == 0 {
//			allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("image"), ""))
//		}
//		// this is only called from ValidatePodUpdate so its safe to check leading/trailing whitespace.
//		if len(strings.TrimSpace(ctr.Image)) != len(ctr.Image) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Index(i).Child("image"), ctr.Image, "must not have leading or trailing whitespace"))
//		}
//	}
//	return allErrs, false
//}
//
// Validate compute resource typename.
// Refer to docs/design/resources.md for more details.
func validateResourceName(value string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, msg := range validation.IsQualifiedName(value) {
		allErrs = append(allErrs, field.Invalid(fldPath, value, msg))
	}
	if len(allErrs) != 0 {
		return allErrs
	}

	if len(strings.Split(value, "/")) == 1 {
		if !IsStandardResourceName(value) {
			return append(allErrs, field.Invalid(fldPath, value, "must be a standard resource type or fully qualified"))
		}
	}

	return allErrs
}

// Validate container resource name
// Refer to docs/design/resources.md for more details.
func validateContainerResourceName(value string, fldPath *field.Path) field.ErrorList {
	allErrs := validateResourceName(value, fldPath)

	if len(strings.Split(value, "/")) == 1 {
		if !IsStandardContainerResourceName(value) {
			return append(allErrs, field.Invalid(fldPath, value, "must be a standard resource for containers"))
		}
	} else if !IsNativeResource(corev1.ResourceName(value)) {
		if !IsExtendedResourceName(corev1.ResourceName(value)) {
			return append(allErrs, field.Invalid(fldPath, value, "doesn't follow extended resource name standard"))
		}
	}
	return allErrs
}

func ValidatePodSpec(spec *corev1.PodSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	vols, vErrs := ValidateVolumes(spec.Volumes, fldPath.Child("volumes"))
	allErrs = append(allErrs, vErrs...)
	allErrs = append(allErrs, validateContainers(spec.Containers, false, vols, fldPath.Child("containers"))...)
	allErrs = append(allErrs, validateInitContainers(spec.InitContainers, spec.Containers, vols, fldPath.Child("initContainers"))...)
	//allErrs = append(allErrs, validateRestartPolicy(&spec.RestartPolicy, fldPath.Child("restartPolicy"))...)
	//allErrs = append(allErrs, validateDNSPolicy(&spec.DNSPolicy, fldPath.Child("dnsPolicy"))...)
	allErrs = append(allErrs, metav1validation.ValidateLabels(spec.NodeSelector, fldPath.Child("nodeSelector"))...)
	allErrs = append(allErrs, ValidatePodSecurityContext(spec.SecurityContext, spec, fldPath, fldPath.Child("securityContext"))...)
	allErrs = append(allErrs, validateImagePullSecrets(spec.ImagePullSecrets, fldPath.Child("imagePullSecrets"))...)
	if len(spec.ServiceAccountName) > 0 {
		for _, msg := range ValidateServiceAccountName(spec.ServiceAccountName, false) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("serviceAccountName"), spec.ServiceAccountName, msg))
		}
	}

	if len(spec.NodeName) > 0 {
		for _, msg := range ValidateNodeName(spec.NodeName, false) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("nodeName"), spec.NodeName, msg))
		}
	}

	if spec.ActiveDeadlineSeconds != nil {
		value := *spec.ActiveDeadlineSeconds
		if value < 1 || value > math.MaxInt32 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("activeDeadlineSeconds"), value, validation.InclusiveRangeError(1, math.MaxInt32)))
		}
	}

	if len(spec.Hostname) > 0 {
		allErrs = append(allErrs, ValidateDNS1123Label(spec.Hostname, fldPath.Child("hostname"))...)
	}

	if len(spec.Subdomain) > 0 {
		allErrs = append(allErrs, ValidateDNS1123Label(spec.Subdomain, fldPath.Child("subdomain"))...)
	}

	if len(spec.Tolerations) > 0 {
		allErrs = append(allErrs, ValidateTolerations(spec.Tolerations, fldPath.Child("tolerations"))...)
	}

	if len(spec.HostAliases) > 0 {
		allErrs = append(allErrs, ValidateHostAliases(spec.HostAliases, fldPath.Child("hostAliases"))...)
	}

	if len(spec.PriorityClassName) > 0 {
		for _, msg := range ValidatePriorityClassName(spec.PriorityClassName, false) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("priorityClassName"), spec.PriorityClassName, msg))
		}
	}

	if spec.RuntimeClassName != nil {
		allErrs = append(allErrs, ValidateRuntimeClassName(*spec.RuntimeClassName, fldPath.Child("runtimeClassName"))...)
	}

	if spec.PreemptionPolicy != nil {
		allErrs = append(allErrs, ValidatePreemptionPolicy(spec.PreemptionPolicy, fldPath.Child("preemptionPolicy"))...)
	}

	return allErrs
}

func ValidateVolumes(volumes []corev1.Volume, fldPath *field.Path) (map[string]corev1.VolumeSource, field.ErrorList) {
	allErrs := field.ErrorList{}

	allNames := sets.String{}
	vols := make(map[string]corev1.VolumeSource)
	for i, vol := range volumes {
		idxPath := fldPath.Index(i)
		namePath := idxPath.Child("name")
		el := validateVolumeSource(&vol.VolumeSource, idxPath, vol.Name)
		if len(vol.Name) == 0 {
			el = append(el, field.Required(namePath, ""))
		} else {
			el = append(el, ValidateDNS1123Label(vol.Name, namePath)...)
		}
		if allNames.Has(vol.Name) {
			el = append(el, field.Duplicate(namePath, vol.Name))
		}
		if len(el) == 0 {
			allNames.Insert(vol.Name)
			vols[vol.Name] = vol.VolumeSource
		} else {
			allErrs = append(allErrs, el...)
		}

	}
	return vols, allErrs
}

func validateVolumeSource(source *corev1.VolumeSource, fldPath *field.Path, volName string) field.ErrorList {
	numVolumes := 0
	allErrs := field.ErrorList{}
	if source.EmptyDir != nil {
		numVolumes++
		if source.EmptyDir.SizeLimit != nil && source.EmptyDir.SizeLimit.Cmp(resource.Quantity{}) < 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("emptyDir").Child("sizeLimit"), "SizeLimit field must be a valid resource quantity"))
		}
	}
	if source.HostPath != nil {
		if numVolumes > 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("hostPath"), "may not specify more than 1 volume type"))
		} else {
			numVolumes++
			allErrs = append(allErrs, validateHostPathVolumeSource(source.HostPath, fldPath.Child("hostPath"))...)
		}
	}
	if source.Secret != nil {
		if numVolumes > 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("secret"), "may not specify more than 1 volume type"))
		} else {
			numVolumes++
			allErrs = append(allErrs, validateSecretVolumeSource(source.Secret, fldPath.Child("secret"))...)
		}
	}
	if source.PersistentVolumeClaim != nil {
		if numVolumes > 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("persistentVolumeClaim"), "may not specify more than 1 volume type"))
		} else {
			numVolumes++
			allErrs = append(allErrs, validatePersistentClaimVolumeSource(source.PersistentVolumeClaim, fldPath.Child("persistentVolumeClaim"))...)
		}
	}
	if source.ConfigMap != nil {
		if numVolumes > 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("configMap"), "may not specify more than 1 volume type"))
		} else {
			numVolumes++
			allErrs = append(allErrs, validateConfigMapVolumeSource(source.ConfigMap, fldPath.Child("configMap"))...)
		}
	}

	if numVolumes == 0 {
		allErrs = append(allErrs, field.Required(fldPath, "must specify a volume type"))
	}
	return allErrs
}

//// Validate resource names that can go in a resource quota
//// Refer to docs/design/resources.md for more details.
//func ValidateResourceQuotaResourceName(value string, fldPath *field.Path) field.ErrorList {
//	allErrs := validateResourceName(value, fldPath)
//
//	if len(strings.Split(value, "/")) == 1 {
//		if !helper.IsStandardQuotaResourceName(value) {
//			return append(allErrs, field.Invalid(fldPath, value, isInvalidQuotaResource))
//		}
//	}
//	return allErrs
//}
//
//// Validate limit range types
//func validateLimitRangeTypeName(value string, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	for _, msg := range validation.IsQualifiedName(value) {
//		allErrs = append(allErrs, field.Invalid(fldPath, value, msg))
//	}
//	if len(allErrs) != 0 {
//		return allErrs
//	}
//
//	if len(strings.Split(value, "/")) == 1 {
//		if !helper.IsStandardLimitRangeType(value) {
//			return append(allErrs, field.Invalid(fldPath, value, "must be a standard limit type or fully qualified"))
//		}
//	}
//
//	return allErrs
//}
//
//// Validate limit range resource name
//// limit types (other than Pod/Container) could contain storage not just cpu or memory
//func validateLimitRangeResourceName(limitType corev1.LimitType, value string, fldPath *field.Path) field.ErrorList {
//	switch limitType {
//	case corev1.LimitTypePod, corev1.LimitTypeContainer:
//		return validateContainerResourceName(value, fldPath)
//	default:
//		return validateResourceName(value, fldPath)
//	}
//}
//
//// ValidateLimitRange tests if required fields in the LimitRange are set.
//func ValidateLimitRange(limitRange *corev1.LimitRange) field.ErrorList {
//	allErrs := ValidateObjectMeta(&limitRange.ObjectMeta, true, ValidateLimitRangeName, field.NewPath("metadata"))
//
//	// ensure resource names are properly qualified per docs/design/resources.md
//	limitTypeSet := map[corev1.LimitType]bool{}
//	fldPath := field.NewPath("spec", "limits")
//	for i := range limitRange.Spec.Limits {
//		idxPath := fldPath.Index(i)
//		limit := &limitRange.Spec.Limits[i]
//		allErrs = append(allErrs, validateLimitRangeTypeName(string(limit.Type), idxPath.Child("type"))...)
//
//		_, found := limitTypeSet[limit.Type]
//		if found {
//			allErrs = append(allErrs, field.Duplicate(idxPath.Child("type"), limit.Type))
//		}
//		limitTypeSet[limit.Type] = true
//
//		keys := sets.String{}
//		min := map[string]resource.Quantity{}
//		max := map[string]resource.Quantity{}
//		defaults := map[string]resource.Quantity{}
//		defaultRequests := map[string]resource.Quantity{}
//		maxLimitRequestRatios := map[string]resource.Quantity{}
//
//		for k, q := range limit.Max {
//			allErrs = append(allErrs, validateLimitRangeResourceName(limit.Type, string(k), idxPath.Child("max").Key(string(k)))...)
//			keys.Insert(string(k))
//			max[string(k)] = q
//		}
//		for k, q := range limit.Min {
//			allErrs = append(allErrs, validateLimitRangeResourceName(limit.Type, string(k), idxPath.Child("min").Key(string(k)))...)
//			keys.Insert(string(k))
//			min[string(k)] = q
//		}
//
//		if limit.Type == corev1.LimitTypePod {
//			if len(limit.Default) > 0 {
//				allErrs = append(allErrs, field.Forbidden(idxPath.Child("default"), "may not be specified when `type` is 'Pod'"))
//			}
//			if len(limit.DefaultRequest) > 0 {
//				allErrs = append(allErrs, field.Forbidden(idxPath.Child("defaultRequest"), "may not be specified when `type` is 'Pod'"))
//			}
//		} else {
//			for k, q := range limit.Default {
//				allErrs = append(allErrs, validateLimitRangeResourceName(limit.Type, string(k), idxPath.Child("default").Key(string(k)))...)
//				keys.Insert(string(k))
//				defaults[string(k)] = q
//			}
//			for k, q := range limit.DefaultRequest {
//				allErrs = append(allErrs, validateLimitRangeResourceName(limit.Type, string(k), idxPath.Child("defaultRequest").Key(string(k)))...)
//				keys.Insert(string(k))
//				defaultRequests[string(k)] = q
//			}
//		}
//
//		if limit.Type == corev1.LimitTypePersistentVolumeClaim {
//			_, minQuantityFound := limit.Min[corev1.ResourceStorage]
//			_, maxQuantityFound := limit.Max[corev1.ResourceStorage]
//			if !minQuantityFound && !maxQuantityFound {
//				allErrs = append(allErrs, field.Required(idxPath.Child("limits"), "either minimum or maximum storage value is required, but neither was provided"))
//			}
//		}
//
//		for k, q := range limit.MaxLimitRequestRatio {
//			allErrs = append(allErrs, validateLimitRangeResourceName(limit.Type, string(k), idxPath.Child("maxLimitRequestRatio").Key(string(k)))...)
//			keys.Insert(string(k))
//			maxLimitRequestRatios[string(k)] = q
//		}
//
//		for k := range keys {
//			minQuantity, minQuantityFound := min[k]
//			maxQuantity, maxQuantityFound := max[k]
//			defaultQuantity, defaultQuantityFound := defaults[k]
//			defaultRequestQuantity, defaultRequestQuantityFound := defaultRequests[k]
//			maxRatio, maxRatioFound := maxLimitRequestRatios[k]
//
//			if minQuantityFound && maxQuantityFound && minQuantity.Cmp(maxQuantity) > 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("min").Key(string(k)), minQuantity, fmt.Sprintf("min value %s is greater than max value %s", minQuantity.String(), maxQuantity.String())))
//			}
//
//			if defaultRequestQuantityFound && minQuantityFound && minQuantity.Cmp(defaultRequestQuantity) > 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("defaultRequest").Key(string(k)), defaultRequestQuantity, fmt.Sprintf("min value %s is greater than default request value %s", minQuantity.String(), defaultRequestQuantity.String())))
//			}
//
//			if defaultRequestQuantityFound && maxQuantityFound && defaultRequestQuantity.Cmp(maxQuantity) > 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("defaultRequest").Key(string(k)), defaultRequestQuantity, fmt.Sprintf("default request value %s is greater than max value %s", defaultRequestQuantity.String(), maxQuantity.String())))
//			}
//
//			if defaultRequestQuantityFound && defaultQuantityFound && defaultRequestQuantity.Cmp(defaultQuantity) > 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("defaultRequest").Key(string(k)), defaultRequestQuantity, fmt.Sprintf("default request value %s is greater than default limit value %s", defaultRequestQuantity.String(), defaultQuantity.String())))
//			}
//
//			if defaultQuantityFound && minQuantityFound && minQuantity.Cmp(defaultQuantity) > 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("default").Key(string(k)), minQuantity, fmt.Sprintf("min value %s is greater than default value %s", minQuantity.String(), defaultQuantity.String())))
//			}
//
//			if defaultQuantityFound && maxQuantityFound && defaultQuantity.Cmp(maxQuantity) > 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("default").Key(string(k)), maxQuantity, fmt.Sprintf("default value %s is greater than max value %s", defaultQuantity.String(), maxQuantity.String())))
//			}
//			if maxRatioFound && maxRatio.Cmp(*resource.NewQuantity(1, resource.DecimalSI)) < 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("maxLimitRequestRatio").Key(string(k)), maxRatio, fmt.Sprintf("ratio %s is less than 1", maxRatio.String())))
//			}
//			if maxRatioFound && minQuantityFound && maxQuantityFound {
//				maxRatioValue := float64(maxRatio.Value())
//				minQuantityValue := minQuantity.Value()
//				maxQuantityValue := maxQuantity.Value()
//				if maxRatio.Value() < resource.MaxMilliValue && minQuantityValue < resource.MaxMilliValue && maxQuantityValue < resource.MaxMilliValue {
//					maxRatioValue = float64(maxRatio.MilliValue()) / 1000
//					minQuantityValue = minQuantity.MilliValue()
//					maxQuantityValue = maxQuantity.MilliValue()
//				}
//				maxRatioLimit := float64(maxQuantityValue) / float64(minQuantityValue)
//				if maxRatioValue > maxRatioLimit {
//					allErrs = append(allErrs, field.Invalid(idxPath.Child("maxLimitRequestRatio").Key(string(k)), maxRatio, fmt.Sprintf("ratio %s is greater than max/min = %f", maxRatio.String(), maxRatioLimit)))
//				}
//			}
//
//			// for GPU, hugepages and other resources that are not allowed to overcommit,
//			// the default value and defaultRequest value must match if both are specified
//			if !helper.IsOvercommitAllowed(corev1.ResourceName(k)) && defaultQuantityFound && defaultRequestQuantityFound && defaultQuantity.Cmp(defaultRequestQuantity) != 0 {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("defaultRequest").Key(string(k)), defaultRequestQuantity, fmt.Sprintf("default value %s must equal to defaultRequest value %s in %s", defaultQuantity.String(), defaultRequestQuantity.String(), k)))
//			}
//		}
//	}
//
//	return allErrs
//}
//
//// ValidateServiceAccount tests if required fields in the ServiceAccount are set.
//func ValidateServiceAccount(serviceAccount *corev1.ServiceAccount) field.ErrorList {
//	allErrs := ValidateObjectMeta(&serviceAccount.ObjectMeta, true, ValidateServiceAccountName, field.NewPath("metadata"))
//	return allErrs
//}
//
//// ValidateServiceAccountUpdate tests if required fields in the ServiceAccount are set.
//func ValidateServiceAccountUpdate(newServiceAccount, oldServiceAccount *corev1.ServiceAccount) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newServiceAccount.ObjectMeta, &oldServiceAccount.ObjectMeta, field.NewPath("metadata"))
//	allErrs = append(allErrs, ValidateServiceAccount(newServiceAccount)...)
//	return allErrs
//}
//
//// ValidateSecret tests if required fields in the Secret are set.
//func ValidateSecret(secret *corev1.Secret) field.ErrorList {
//	allErrs := ValidateObjectMeta(&secret.ObjectMeta, true, ValidateSecretName, field.NewPath("metadata"))
//
//	dataPath := field.NewPath("data")
//	totalSize := 0
//	for key, value := range secret.Data {
//		for _, msg := range validation.IsConfigMapKey(key) {
//			allErrs = append(allErrs, field.Invalid(dataPath.Key(key), key, msg))
//		}
//		totalSize += len(value)
//	}
//	if totalSize > corev1.MaxSecretSize {
//		allErrs = append(allErrs, field.TooLong(dataPath, "", corev1.MaxSecretSize))
//	}
//
//	switch secret.Type {
//	case corev1.SecretTypeServiceAccountToken:
//		// Only require Annotations[kubernetes.io/service-account.name]
//		// Additional fields (like Annotations[kubernetes.io/service-account.uid] and Data[token]) might be contributed later by a controller loop
//		if value := secret.Annotations[corev1.ServiceAccountNameKey]; len(value) == 0 {
//			allErrs = append(allErrs, field.Required(field.NewPath("metadata", "annotations").Key(corev1.ServiceAccountNameKey), ""))
//		}
//	case corev1.SecretTypeOpaque, "":
//	// no-op
//	case corev1.SecretTypeDockercfg:
//		dockercfgBytes, exists := secret.Data[corev1.DockerConfigKey]
//		if !exists {
//			allErrs = append(allErrs, field.Required(dataPath.Key(corev1.DockerConfigKey), ""))
//			break
//		}
//
//		// make sure that the content is well-formed json.
//		if err := json.Unmarshal(dockercfgBytes, &map[string]interface{}{}); err != nil {
//			allErrs = append(allErrs, field.Invalid(dataPath.Key(corev1.DockerConfigKey), "<secret contents redacted>", err.Error()))
//		}
//	case corev1.SecretTypeDockerConfigJson:
//		dockerConfigJsonBytes, exists := secret.Data[corev1.DockerConfigJsonKey]
//		if !exists {
//			allErrs = append(allErrs, field.Required(dataPath.Key(corev1.DockerConfigJsonKey), ""))
//			break
//		}
//
//		// make sure that the content is well-formed json.
//		if err := json.Unmarshal(dockerConfigJsonBytes, &map[string]interface{}{}); err != nil {
//			allErrs = append(allErrs, field.Invalid(dataPath.Key(corev1.DockerConfigJsonKey), "<secret contents redacted>", err.Error()))
//		}
//	case corev1.SecretTypeBasicAuth:
//		_, usernameFieldExists := secret.Data[corev1.BasicAuthUsernameKey]
//		_, passwordFieldExists := secret.Data[corev1.BasicAuthPasswordKey]
//
//		// username or password might be empty, but the field must be present
//		if !usernameFieldExists && !passwordFieldExists {
//			allErrs = append(allErrs, field.Required(field.NewPath("data[%s]").Key(corev1.BasicAuthUsernameKey), ""))
//			allErrs = append(allErrs, field.Required(field.NewPath("data[%s]").Key(corev1.BasicAuthPasswordKey), ""))
//			break
//		}
//	case corev1.SecretTypeSSHAuth:
//		if len(secret.Data[corev1.SSHAuthPrivateKey]) == 0 {
//			allErrs = append(allErrs, field.Required(field.NewPath("data[%s]").Key(corev1.SSHAuthPrivateKey), ""))
//			break
//		}
//
//	case corev1.SecretTypeTLS:
//		if _, exists := secret.Data[corev1.TLSCertKey]; !exists {
//			allErrs = append(allErrs, field.Required(dataPath.Key(corev1.TLSCertKey), ""))
//		}
//		if _, exists := secret.Data[corev1.TLSPrivateKeyKey]; !exists {
//			allErrs = append(allErrs, field.Required(dataPath.Key(corev1.TLSPrivateKeyKey), ""))
//		}
//	// TODO: Verify that the key matches the cert.
//	default:
//		// no-op
//	}
//
//	return allErrs
//}
//
//// ValidateSecretUpdate tests if required fields in the Secret are set.
//func ValidateSecretUpdate(newSecret, oldSecret *corev1.Secret) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newSecret.ObjectMeta, &oldSecret.ObjectMeta, field.NewPath("metadata"))
//
//	if len(newSecret.Type) == 0 {
//		newSecret.Type = oldSecret.Type
//	}
//
//	allErrs = append(allErrs, ValidateImmutableField(newSecret.Type, oldSecret.Type, field.NewPath("type"))...)
//
//	allErrs = append(allErrs, ValidateSecret(newSecret)...)
//	return allErrs
//}
//
// ValidateConfigMapName can be used to check whether the given ConfigMap name is valid.
// Prefix indicates this name will be used as part of generation, in which case
// trailing dashes are allowed.
var ValidateConfigMapName = apimachineryvalidation.NameIsDNSSubdomain

// ValidateConfigMap tests whether required fields in the ConfigMap are set.
func ValidateConfigMap(cfg *corev1.ConfigMap) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateObjectMeta(&cfg.ObjectMeta, true, ValidateConfigMapName, field.NewPath("metadata"))...)

	totalSize := 0

	for key, value := range cfg.Data {
		for _, msg := range validation.IsConfigMapKey(key) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("data").Key(key), key, msg))
		}
		// check if we have a duplicate key in the other bag
		if _, isValue := cfg.BinaryData[key]; isValue {
			msg := "duplicate of key present in binaryData"
			allErrs = append(allErrs, field.Invalid(field.NewPath("data").Key(key), key, msg))
		}
		totalSize += len(value)
	}
	for key, value := range cfg.BinaryData {
		for _, msg := range validation.IsConfigMapKey(key) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("binaryData").Key(key), key, msg))
		}
		totalSize += len(value)
	}
	if totalSize > corev1.MaxSecretSize {
		// pass back "" to indicate that the error refers to the whole object.
		allErrs = append(allErrs, field.TooLong(field.NewPath(""), cfg, corev1.MaxSecretSize))
	}

	return allErrs
}

//// ValidateConfigMapUpdate tests if required fields in the ConfigMap are set.
//func ValidateConfigMapUpdate(newCfg, oldCfg *corev1.ConfigMap) field.ErrorList {
//	allErrs := field.ErrorList{}
//	allErrs = append(allErrs, ValidateObjectMetaUpdate(&newCfg.ObjectMeta, &oldCfg.ObjectMeta, field.NewPath("metadata"))...)
//	allErrs = append(allErrs, ValidateConfigMap(newCfg)...)
//
//	return allErrs
//}
//
//func validateBasicResource(quantity resource.Quantity, fldPath *field.Path) field.ErrorList {
//	if quantity.Value() < 0 {
//		return field.ErrorList{field.Invalid(fldPath, quantity.Value(), "must be a valid resource quantity")}
//	}
//	return field.ErrorList{}
//}
//
// Validates resource requirement spec.
func ValidateResourceRequirements(requirements *corev1.ResourceRequirements, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	limPath := fldPath.Child("limits")
	reqPath := fldPath.Child("requests")
	limContainsCPUOrMemory := false
	reqContainsCPUOrMemory := false
	limContainsHugePages := false
	reqContainsHugePages := false
	supportedQoSComputeResources := sets.NewString(string(corev1.ResourceCPU), string(corev1.ResourceMemory))
	for resourceName, quantity := range requirements.Limits {

		fldPath := limPath.Key(string(resourceName))
		// Validate resource name.
		allErrs = append(allErrs, validateContainerResourceName(string(resourceName), fldPath)...)

		// Validate resource quantity.
		allErrs = append(allErrs, ValidateResourceQuantityValue(string(resourceName), quantity, fldPath)...)

		if IsHugePageResourceName(resourceName) {
			limContainsHugePages = true
		}

		if supportedQoSComputeResources.Has(string(resourceName)) {
			limContainsCPUOrMemory = true
		}
	}
	for resourceName, quantity := range requirements.Requests {
		fldPath := reqPath.Key(string(resourceName))
		// Validate resource name.
		allErrs = append(allErrs, validateContainerResourceName(string(resourceName), fldPath)...)
		// Validate resource quantity.
		allErrs = append(allErrs, ValidateResourceQuantityValue(string(resourceName), quantity, fldPath)...)

		// Check that request <= limit.
		limitQuantity, exists := requirements.Limits[resourceName]
		if exists {
			// For non overcommitable resources, not only requests can't exceed limits, they also can't be lower, i.e. must be equal.
			if quantity.Cmp(limitQuantity) != 0 && !IsOvercommitAllowed(resourceName) {
				allErrs = append(allErrs, field.Invalid(reqPath, quantity.String(), fmt.Sprintf("must be equal to %s limit", resourceName)))
			} else if quantity.Cmp(limitQuantity) > 0 {
				allErrs = append(allErrs, field.Invalid(reqPath, quantity.String(), fmt.Sprintf("must be less than or equal to %s limit", resourceName)))
			}
		} else if !IsOvercommitAllowed(resourceName) {
			allErrs = append(allErrs, field.Required(limPath, "Limit must be set for non overcommitable resources"))
		}
		if IsHugePageResourceName(resourceName) {
			reqContainsHugePages = true
		}
		if supportedQoSComputeResources.Has(string(resourceName)) {
			reqContainsCPUOrMemory = true
		}

	}
	if !limContainsCPUOrMemory && !reqContainsCPUOrMemory && (reqContainsHugePages || limContainsHugePages) {
		allErrs = append(allErrs, field.Forbidden(fldPath, fmt.Sprintf("HugePages require cpu or memory")))
	}

	return allErrs
}

//// validateResourceQuotaScopes ensures that each enumerated hard resource constraint is valid for set of scopes
//func validateResourceQuotaScopes(resourceQuotaSpec *corev1.ResourceQuotaSpec, fld *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if len(resourceQuotaSpec.Scopes) == 0 {
//		return allErrs
//	}
//	hardLimits := sets.NewString()
//	for k := range resourceQuotaSpec.Hard {
//		hardLimits.Insert(string(k))
//	}
//	fldPath := fld.Child("scopes")
//	scopeSet := sets.NewString()
//	for _, scope := range resourceQuotaSpec.Scopes {
//		if !helper.IsStandardResourceQuotaScope(string(scope)) {
//			allErrs = append(allErrs, field.Invalid(fldPath, resourceQuotaSpec.Scopes, "unsupported scope"))
//		}
//		for _, k := range hardLimits.List() {
//			if helper.IsStandardQuotaResourceName(k) && !helper.IsResourceQuotaScopeValidForResource(scope, k) {
//				allErrs = append(allErrs, field.Invalid(fldPath, resourceQuotaSpec.Scopes, "unsupported scope applied to resource"))
//			}
//		}
//		scopeSet.Insert(string(scope))
//	}
//	invalidScopePairs := []sets.String{
//		sets.NewString(string(corev1.ResourceQuotaScopeBestEffort), string(corev1.ResourceQuotaScopeNotBestEffort)),
//		sets.NewString(string(corev1.ResourceQuotaScopeTerminating), string(corev1.ResourceQuotaScopeNotTerminating)),
//	}
//	for _, invalidScopePair := range invalidScopePairs {
//		if scopeSet.HasAll(invalidScopePair.List()...) {
//			allErrs = append(allErrs, field.Invalid(fldPath, resourceQuotaSpec.Scopes, "conflicting scopes"))
//		}
//	}
//	return allErrs
//}
//
//// validateScopedResourceSelectorRequirement tests that the match expressions has valid data
//func validateScopedResourceSelectorRequirement(resourceQuotaSpec *corev1.ResourceQuotaSpec, fld *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	hardLimits := sets.NewString()
//	for k := range resourceQuotaSpec.Hard {
//		hardLimits.Insert(string(k))
//	}
//	fldPath := fld.Child("matchExpressions")
//	scopeSet := sets.NewString()
//	for _, req := range resourceQuotaSpec.ScopeSelector.MatchExpressions {
//		if !helper.IsStandardResourceQuotaScope(string(req.ScopeName)) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("scopeName"), req.ScopeName, "unsupported scope"))
//		}
//		for _, k := range hardLimits.List() {
//			if helper.IsStandardQuotaResourceName(k) && !helper.IsResourceQuotaScopeValidForResource(req.ScopeName, k) {
//				allErrs = append(allErrs, field.Invalid(fldPath, resourceQuotaSpec.ScopeSelector, "unsupported scope applied to resource"))
//			}
//		}
//		switch req.ScopeName {
//		case corev1.ResourceQuotaScopeBestEffort, corev1.ResourceQuotaScopeNotBestEffort, corev1.ResourceQuotaScopeTerminating, corev1.ResourceQuotaScopeNotTerminating:
//			if req.Operator != corev1.ScopeSelectorOpExists {
//				allErrs = append(allErrs, field.Invalid(fldPath.Child("operator"), req.Operator,
//					"must be 'Exist' only operator when scope is any of ResourceQuotaScopeTerminating, ResourceQuotaScopeNotTerminating, ResourceQuotaScopeBestEffort and ResourceQuotaScopeNotBestEffort"))
//			}
//		}
//
//		switch req.Operator {
//		case corev1.ScopeSelectorOpIn, corev1.ScopeSelectorOpNotIn:
//			if len(req.Values) == 0 {
//				allErrs = append(allErrs, field.Required(fldPath.Child("values"),
//					"must be at least one value when `operator` is 'In' or 'NotIn' for scope selector"))
//			}
//		case corev1.ScopeSelectorOpExists, corev1.ScopeSelectorOpDoesNotExist:
//			if len(req.Values) != 0 {
//				allErrs = append(allErrs, field.Invalid(fldPath.Child("values"), req.Values,
//					"must be no value when `operator` is 'Exist' or 'DoesNotExist' for scope selector"))
//			}
//		default:
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("operator"), req.Operator, "not a valid selector operator"))
//		}
//		scopeSet.Insert(string(req.ScopeName))
//	}
//	invalidScopePairs := []sets.String{
//		sets.NewString(string(corev1.ResourceQuotaScopeBestEffort), string(corev1.ResourceQuotaScopeNotBestEffort)),
//		sets.NewString(string(corev1.ResourceQuotaScopeTerminating), string(corev1.ResourceQuotaScopeNotTerminating)),
//	}
//	for _, invalidScopePair := range invalidScopePairs {
//		if scopeSet.HasAll(invalidScopePair.List()...) {
//			allErrs = append(allErrs, field.Invalid(fldPath, resourceQuotaSpec.Scopes, "conflicting scopes"))
//		}
//	}
//
//	return allErrs
//}
//
//// validateScopeSelector tests that the specified scope selector has valid data
//func validateScopeSelector(resourceQuotaSpec *corev1.ResourceQuotaSpec, fld *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if resourceQuotaSpec.ScopeSelector == nil {
//		return allErrs
//	}
//	allErrs = append(allErrs, validateScopedResourceSelectorRequirement(resourceQuotaSpec, fld.Child("scopeSelector"))...)
//	return allErrs
//}
//
//// ValidateResourceQuota tests if required fields in the ResourceQuota are set.
//func ValidateResourceQuota(resourceQuota *corev1.ResourceQuota) field.ErrorList {
//	allErrs := ValidateObjectMeta(&resourceQuota.ObjectMeta, true, ValidateResourceQuotaName, field.NewPath("metadata"))
//
//	allErrs = append(allErrs, ValidateResourceQuotaSpec(&resourceQuota.Spec, field.NewPath("spec"))...)
//	allErrs = append(allErrs, ValidateResourceQuotaStatus(&resourceQuota.Status, field.NewPath("status"))...)
//
//	return allErrs
//}
//
//func ValidateResourceQuotaStatus(status *corev1.ResourceQuotaStatus, fld *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	fldPath := fld.Child("hard")
//	for k, v := range status.Hard {
//		resPath := fldPath.Key(string(k))
//		allErrs = append(allErrs, ValidateResourceQuotaResourceName(string(k), resPath)...)
//		allErrs = append(allErrs, ValidateResourceQuantityValue(string(k), v, resPath)...)
//	}
//	fldPath = fld.Child("used")
//	for k, v := range status.Used {
//		resPath := fldPath.Key(string(k))
//		allErrs = append(allErrs, ValidateResourceQuotaResourceName(string(k), resPath)...)
//		allErrs = append(allErrs, ValidateResourceQuantityValue(string(k), v, resPath)...)
//	}
//
//	return allErrs
//}
//
//func ValidateResourceQuotaSpec(resourceQuotaSpec *corev1.ResourceQuotaSpec, fld *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//
//	fldPath := fld.Child("hard")
//	for k, v := range resourceQuotaSpec.Hard {
//		resPath := fldPath.Key(string(k))
//		allErrs = append(allErrs, ValidateResourceQuotaResourceName(string(k), resPath)...)
//		allErrs = append(allErrs, ValidateResourceQuantityValue(string(k), v, resPath)...)
//	}
//	allErrs = append(allErrs, validateResourceQuotaScopes(resourceQuotaSpec, fld)...)
//	allErrs = append(allErrs, validateScopeSelector(resourceQuotaSpec, fld)...)
//
//	return allErrs
//}
//
// ValidateResourceQuantityValue enforces that specified quantity is valid for specified resource
func ValidateResourceQuantityValue(resource string, value resource.Quantity, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateNonnegativeQuantity(value, fldPath)...)
	if IsIntegerResourceName(resource) {
		if value.MilliValue()%int64(1000) != int64(0) {
			allErrs = append(allErrs, field.Invalid(fldPath, value, isNotIntegerErrorMsg))
		}
	}
	return allErrs
}

//// ValidateResourceQuotaUpdate tests to see if the update is legal for an end user to make.
//// newResourceQuota is updated with fields that cannot be changed.
//func ValidateResourceQuotaUpdate(newResourceQuota, oldResourceQuota *corev1.ResourceQuota) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newResourceQuota.ObjectMeta, &oldResourceQuota.ObjectMeta, field.NewPath("metadata"))
//	allErrs = append(allErrs, ValidateResourceQuotaSpec(&newResourceQuota.Spec, field.NewPath("spec"))...)
//
//	// ensure scopes cannot change, and that resources are still valid for scope
//	fldPath := field.NewPath("spec", "scopes")
//	oldScopes := sets.NewString()
//	newScopes := sets.NewString()
//	for _, scope := range newResourceQuota.Spec.Scopes {
//		newScopes.Insert(string(scope))
//	}
//	for _, scope := range oldResourceQuota.Spec.Scopes {
//		oldScopes.Insert(string(scope))
//	}
//	if !oldScopes.Equal(newScopes) {
//		allErrs = append(allErrs, field.Invalid(fldPath, newResourceQuota.Spec.Scopes, fieldImmutableErrorMsg))
//	}
//
//	newResourceQuota.Status = oldResourceQuota.Status
//	return allErrs
//}
//
//// ValidateResourceQuotaStatusUpdate tests to see if the status update is legal for an end user to make.
//// newResourceQuota is updated with fields that cannot be changed.
//func ValidateResourceQuotaStatusUpdate(newResourceQuota, oldResourceQuota *corev1.ResourceQuota) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newResourceQuota.ObjectMeta, &oldResourceQuota.ObjectMeta, field.NewPath("metadata"))
//	if len(newResourceQuota.ResourceVersion) == 0 {
//		allErrs = append(allErrs, field.Required(field.NewPath("resourceVersion"), ""))
//	}
//	fldPath := field.NewPath("status", "hard")
//	for k, v := range newResourceQuota.Status.Hard {
//		resPath := fldPath.Key(string(k))
//		allErrs = append(allErrs, ValidateResourceQuotaResourceName(string(k), resPath)...)
//		allErrs = append(allErrs, ValidateResourceQuantityValue(string(k), v, resPath)...)
//	}
//	fldPath = field.NewPath("status", "used")
//	for k, v := range newResourceQuota.Status.Used {
//		resPath := fldPath.Key(string(k))
//		allErrs = append(allErrs, ValidateResourceQuotaResourceName(string(k), resPath)...)
//		allErrs = append(allErrs, ValidateResourceQuantityValue(string(k), v, resPath)...)
//	}
//	newResourceQuota.Spec = oldResourceQuota.Spec
//	return allErrs
//}
//
//// ValidateNamespace tests if required fields are set.
//func ValidateNamespace(namespace *corev1.Namespace) field.ErrorList {
//	allErrs := ValidateObjectMeta(&namespace.ObjectMeta, false, ValidateNamespaceName, field.NewPath("metadata"))
//	for i := range namespace.Spec.Finalizers {
//		allErrs = append(allErrs, validateFinalizerName(string(namespace.Spec.Finalizers[i]), field.NewPath("spec", "finalizers"))...)
//	}
//	return allErrs
//}
//
//// Validate finalizer names
//func validateFinalizerName(stringValue string, fldPath *field.Path) field.ErrorList {
//	allErrs := apimachineryvalidation.ValidateFinalizerName(stringValue, fldPath)
//	allErrs = append(allErrs, validateKubeFinalizerName(stringValue, fldPath)...)
//	return allErrs
//}

// validateKubeFinalizerName checks for "standard" names of legacy finalizer
func validateKubeFinalizerName(stringValue string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(strings.Split(stringValue, "/")) == 1 {
		if !IsStandardFinalizerName(stringValue) {
			return append(allErrs, field.Invalid(fldPath, stringValue, "name is neither a standard finalizer name nor is it fully qualified"))
		}
	}

	return allErrs
}

//// ValidateNamespaceUpdate tests to make sure a namespace update can be applied.
//// newNamespace is updated with fields that cannot be changed
//func ValidateNamespaceUpdate(newNamespace *corev1.Namespace, oldNamespace *corev1.Namespace) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newNamespace.ObjectMeta, &oldNamespace.ObjectMeta, field.NewPath("metadata"))
//	newNamespace.Spec.Finalizers = oldNamespace.Spec.Finalizers
//	newNamespace.Status = oldNamespace.Status
//	return allErrs
//}
//
//// ValidateNamespaceStatusUpdate tests to see if the update is legal for an end user to make. newNamespace is updated with fields
//// that cannot be changed.
//func ValidateNamespaceStatusUpdate(newNamespace, oldNamespace *corev1.Namespace) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newNamespace.ObjectMeta, &oldNamespace.ObjectMeta, field.NewPath("metadata"))
//	newNamespace.Spec = oldNamespace.Spec
//	if newNamespace.DeletionTimestamp.IsZero() {
//		if newNamespace.Status.Phase != corev1.NamespaceActive {
//			allErrs = append(allErrs, field.Invalid(field.NewPath("status", "Phase"), newNamespace.Status.Phase, "may only be 'Active' if `deletionTimestamp` is empty"))
//		}
//	} else {
//		if newNamespace.Status.Phase != corev1.NamespaceTerminating {
//			allErrs = append(allErrs, field.Invalid(field.NewPath("status", "Phase"), newNamespace.Status.Phase, "may only be 'Terminating' if `deletionTimestamp` is not empty"))
//		}
//	}
//	return allErrs
//}
//
//// ValidateNamespaceFinalizeUpdate tests to see if the update is legal for an end user to make.
//// newNamespace is updated with fields that cannot be changed.
//func ValidateNamespaceFinalizeUpdate(newNamespace, oldNamespace *corev1.Namespace) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newNamespace.ObjectMeta, &oldNamespace.ObjectMeta, field.NewPath("metadata"))
//
//	fldPath := field.NewPath("spec", "finalizers")
//	for i := range newNamespace.Spec.Finalizers {
//		idxPath := fldPath.Index(i)
//		allErrs = append(allErrs, validateFinalizerName(string(newNamespace.Spec.Finalizers[i]), idxPath)...)
//	}
//	newNamespace.Status = oldNamespace.Status
//	return allErrs
//}
//
//// ValidateEndpoints tests if required fields are set.
//func ValidateEndpoints(endpoints *corev1.Endpoints) field.ErrorList {
//	allErrs := ValidateObjectMeta(&endpoints.ObjectMeta, true, ValidateEndpointsName, field.NewPath("metadata"))
//	allErrs = append(allErrs, ValidateEndpointsSpecificAnnotations(endpoints.Annotations, field.NewPath("annotations"))...)
//	allErrs = append(allErrs, validateEndpointSubsets(endpoints.Subsets, field.NewPath("subsets"))...)
//	return allErrs
//}
//
//func validateEndpointSubsets(subsets []corev1.EndpointSubset, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	for i := range subsets {
//		ss := &subsets[i]
//		idxPath := fldPath.Index(i)
//
//		// EndpointSubsets must include endpoint address. For headless service, we allow its endpoints not to have ports.
//		if len(ss.Addresses) == 0 && len(ss.NotReadyAddresses) == 0 {
//			//TODO: consider adding a RequiredOneOf() error for this and similar cases
//			allErrs = append(allErrs, field.Required(idxPath, "must specify `addresses` or `notReadyAddresses`"))
//		}
//		for addr := range ss.Addresses {
//			allErrs = append(allErrs, validateEndpointAddress(&ss.Addresses[addr], idxPath.Child("addresses").Index(addr))...)
//		}
//		for addr := range ss.NotReadyAddresses {
//			allErrs = append(allErrs, validateEndpointAddress(&ss.NotReadyAddresses[addr], idxPath.Child("notReadyAddresses").Index(addr))...)
//		}
//		for port := range ss.Ports {
//			allErrs = append(allErrs, validateEndpointPort(&ss.Ports[port], len(ss.Ports) > 1, idxPath.Child("ports").Index(port))...)
//		}
//	}
//
//	return allErrs
//}
//
//func validateEndpointAddress(address *corev1.EndpointAddress, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	for _, msg := range validation.IsValidIP(address.IP) {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("ip"), address.IP, msg))
//	}
//	if len(address.Hostname) > 0 {
//		allErrs = append(allErrs, ValidateDNS1123Label(address.Hostname, fldPath.Child("hostname"))...)
//	}
//	// During endpoint update, verify that NodeName is a DNS subdomain and transition rules allow the update
//	if address.NodeName != nil {
//		for _, msg := range ValidateNodeName(*address.NodeName, false) {
//			allErrs = append(allErrs, field.Invalid(fldPath.Child("nodeName"), *address.NodeName, msg))
//		}
//	}
//	allErrs = append(allErrs, validateNonSpecialIP(address.IP, fldPath.Child("ip"))...)
//	return allErrs
//}
//
//func validateNonSpecialIP(ipAddress string, fldPath *field.Path) field.ErrorList {
//	// We disallow some IPs as endpoints or external-ips.  Specifically,
//	// unspecified and loopback addresses are nonsensical and link-local
//	// addresses tend to be used for node-centric purposes (e.g. metadata
//	// service).
//	allErrs := field.ErrorList{}
//	ip := net.ParseIP(ipAddress)
//	if ip == nil {
//		allErrs = append(allErrs, field.Invalid(fldPath, ipAddress, "must be a valid IP address"))
//		return allErrs
//	}
//	if ip.IsUnspecified() {
//		allErrs = append(allErrs, field.Invalid(fldPath, ipAddress, "may not be unspecified (0.0.0.0)"))
//	}
//	if ip.IsLoopback() {
//		allErrs = append(allErrs, field.Invalid(fldPath, ipAddress, "may not be in the loopback range (127.0.0.0/8)"))
//	}
//	if ip.IsLinkLocalUnicast() {
//		allErrs = append(allErrs, field.Invalid(fldPath, ipAddress, "may not be in the link-local range (169.254.0.0/16)"))
//	}
//	if ip.IsLinkLocalMulticast() {
//		allErrs = append(allErrs, field.Invalid(fldPath, ipAddress, "may not be in the link-local multicast range (224.0.0.0/24)"))
//	}
//	return allErrs
//}
//
//func validateEndpointPort(port *corev1.EndpointPort, requireName bool, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if requireName && len(port.Name) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
//	} else if len(port.Name) != 0 {
//		allErrs = append(allErrs, ValidateDNS1123Label(port.Name, fldPath.Child("name"))...)
//	}
//	for _, msg := range validation.IsValidPortNum(int(port.Port)) {
//		allErrs = append(allErrs, field.Invalid(fldPath.Child("port"), port.Port, msg))
//	}
//	if len(port.Protocol) == 0 {
//		allErrs = append(allErrs, field.Required(fldPath.Child("protocol"), ""))
//	} else if !supportedPortProtocols.Has(string(port.Protocol)) {
//		allErrs = append(allErrs, field.NotSupported(fldPath.Child("protocol"), port.Protocol, supportedPortProtocols.List()))
//	}
//	return allErrs
//}
//
//// ValidateEndpointsUpdate tests to make sure an endpoints update can be applied.
//// NodeName changes are allowed during update to accommodate the case where nodeIP or PodCIDR is reused.
//// An existing endpoint ip will have a different nodeName if this happens.
//func ValidateEndpointsUpdate(newEndpoints, oldEndpoints *corev1.Endpoints) field.ErrorList {
//	allErrs := ValidateObjectMetaUpdate(&newEndpoints.ObjectMeta, &oldEndpoints.ObjectMeta, field.NewPath("metadata"))
//	allErrs = append(allErrs, validateEndpointSubsets(newEndpoints.Subsets, field.NewPath("subsets"))...)
//	allErrs = append(allErrs, ValidateEndpointsSpecificAnnotations(newEndpoints.Annotations, field.NewPath("annotations"))...)
//	return allErrs
//}
//
// ValidateSecurityContext ensures the security context contains valid settings
func ValidateSecurityContext(sc *corev1.SecurityContext, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	//this should only be true for testing since SecurityContext is defaulted by the core
	if sc == nil {
		return allErrs
	}

	if sc.RunAsUser != nil {
		for _, msg := range validation.IsValidUserID(*sc.RunAsUser) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("runAsUser"), *sc.RunAsUser, msg))
		}
	}

	if sc.RunAsGroup != nil {
		for _, msg := range validation.IsValidGroupID(*sc.RunAsGroup) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("runAsGroup"), *sc.RunAsGroup, msg))
		}
	}

	if sc.ProcMount != nil {
		if err := ValidateProcMountType(fldPath.Child("procMount"), *sc.ProcMount); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if sc.AllowPrivilegeEscalation != nil && !*sc.AllowPrivilegeEscalation {
		if sc.Privileged != nil && *sc.Privileged {
			allErrs = append(allErrs, field.Invalid(fldPath, sc, "cannot set `allowPrivilegeEscalation` to false and `privileged` to true"))
		}

		if sc.Capabilities != nil {
			for _, cap := range sc.Capabilities.Add {
				if string(cap) == "CAP_SYS_ADMIN" {
					allErrs = append(allErrs, field.Invalid(fldPath, sc, "cannot set `allowPrivilegeEscalation` to false and `capabilities.Add` CAP_SYS_ADMIN"))
				}
			}
		}
	}

	allErrs = append(allErrs, validateWindowsSecurityContextOptions(sc.WindowsOptions, fldPath.Child("windowsOptions"))...)

	return allErrs
}

// maxGMSACredentialSpecLength is the max length, in bytes, for the actual contents
// of a GMSA cred spec. In general, those shouldn't be more than a few hundred bytes,
// so we want to give plenty of room here while still providing an upper bound.
const (
	maxGMSACredentialSpecLengthInKiB = 64
	maxGMSACredentialSpecLength      = maxGMSACredentialSpecLengthInKiB * 1024
)

func validateWindowsSecurityContextOptions(windowsOptions *corev1.WindowsSecurityContextOptions, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if windowsOptions == nil {
		return allErrs
	}

	if windowsOptions.GMSACredentialSpecName != nil {
		// gmsaCredentialSpecName must be the name of a custom resource
		for _, msg := range validation.IsDNS1123Subdomain(*windowsOptions.GMSACredentialSpecName) {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("gmsaCredentialSpecName"), windowsOptions.GMSACredentialSpecName, msg))
		}
	}

	if windowsOptions.GMSACredentialSpec != nil {
		if l := len(*windowsOptions.GMSACredentialSpec); l == 0 {
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("gmsaCredentialSpec"), windowsOptions.GMSACredentialSpec, "gmsaCredentialSpec cannot be an empty string"))
		} else if l > maxGMSACredentialSpecLength {
			errMsg := fmt.Sprintf("gmsaCredentialSpec size must be under %d KiB", maxGMSACredentialSpecLengthInKiB)
			allErrs = append(allErrs, field.Invalid(fieldPath.Child("gmsaCredentialSpec"), windowsOptions.GMSACredentialSpec, errMsg))
		}
	}

	return allErrs
}

//func ValidatePodLogOptions(opts *corev1.PodLogOptions) field.ErrorList {
//	allErrs := field.ErrorList{}
//	if opts.TailLines != nil && *opts.TailLines < 0 {
//		allErrs = append(allErrs, field.Invalid(field.NewPath("tailLines"), *opts.TailLines, isNegativeErrorMsg))
//	}
//	if opts.LimitBytes != nil && *opts.LimitBytes < 1 {
//		allErrs = append(allErrs, field.Invalid(field.NewPath("limitBytes"), *opts.LimitBytes, "must be greater than 0"))
//	}
//	switch {
//	case opts.SinceSeconds != nil && opts.SinceTime != nil:
//		allErrs = append(allErrs, field.Forbidden(field.NewPath(""), "at most one of `sinceTime` or `sinceSeconds` may be specified"))
//	case opts.SinceSeconds != nil:
//		if *opts.SinceSeconds < 1 {
//			allErrs = append(allErrs, field.Invalid(field.NewPath("sinceSeconds"), *opts.SinceSeconds, "must be greater than 0"))
//		}
//	}
//	return allErrs
//}
//
//// ValidateLoadBalancerStatus validates required fields on a LoadBalancerStatus
//func ValidateLoadBalancerStatus(status *corev1.LoadBalancerStatus, fldPath *field.Path) field.ErrorList {
//	allErrs := field.ErrorList{}
//	for i, ingress := range status.Ingress {
//		idxPath := fldPath.Child("ingress").Index(i)
//		if len(ingress.IP) > 0 {
//			if isIP := (net.ParseIP(ingress.IP) != nil); !isIP {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("ip"), ingress.IP, "must be a valid IP address"))
//			}
//		}
//		if len(ingress.Hostname) > 0 {
//			for _, msg := range validation.IsDNS1123Subdomain(ingress.Hostname) {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("hostname"), ingress.Hostname, msg))
//			}
//			if isIP := (net.ParseIP(ingress.Hostname) != nil); isIP {
//				allErrs = append(allErrs, field.Invalid(idxPath.Child("hostname"), ingress.Hostname, "must be a DNS name, not an IP address"))
//			}
//		}
//	}
//	return allErrs
//}
//
//// validateVolumeNodeAffinity tests that the PersistentVolume.NodeAffinity has valid data
//// returns:
//// - true if volumeNodeAffinity is set
//// - errorList if there are validation errors
//func validateVolumeNodeAffinity(nodeAffinity *corev1.VolumeNodeAffinity, fldPath *field.Path) (bool, field.ErrorList) {
//	allErrs := field.ErrorList{}
//
//	if nodeAffinity == nil {
//		return false, allErrs
//	}
//
//	if nodeAffinity.Required != nil {
//		allErrs = append(allErrs, ValidateNodeSelector(nodeAffinity.Required, fldPath.Child("required"))...)
//	} else {
//		allErrs = append(allErrs, field.Required(fldPath.Child("required"), "must specify required node constraints"))
//	}
//
//	return true, allErrs
//}
//
//// ValidateCIDR validates whether a CIDR matches the conventions expected by net.ParseCIDR
//func ValidateCIDR(cidr string) (*net.IPNet, error) {
//	_, net, err := net.ParseCIDR(cidr)
//	if err != nil {
//		return nil, err
//	}
//	return net, nil
//}
//
//func IsDecremented(update, old *int32) bool {
//	if update == nil && old != nil {
//		return true
//	}
//	if update == nil || old == nil {
//		return false
//	}
//	return *update < *old
//}

// ValidateProcMountType tests that the argument is a valid ProcMountType.
func ValidateProcMountType(fldPath *field.Path, procMountType corev1.ProcMountType) *field.Error {
	switch procMountType {
	case corev1.DefaultProcMount, corev1.UnmaskedProcMount:
		return nil
	default:
		return field.NotSupported(fldPath, procMountType, []string{string(corev1.DefaultProcMount), string(corev1.UnmaskedProcMount)})
	}
}

//// ConvertDownwardAPIFieldLabel converts the specified downward API field label
//// and its value in the pod of the specified version to the internal version,
//// and returns the converted label and value. This function returns an error if
//// the conversion fails.
//func ConvertDownwardAPIFieldLabel(version, label, value string) (string, string, error) {
//	if version != "v1" {
//		return "", "", fmt.Errorf("unsupported pod version: %s", version)
//	}
//
//	if path, _, ok := fieldpath.SplitMaybeSubscriptedPath(label); ok {
//		switch path {
//		case "metadata.annotations", "metadata.labels":
//			return label, value, nil
//		default:
//			return "", "", fmt.Errorf("field label does not support subscript: %s", label)
//		}
//	}
//
//	switch label {
//	case "metadata.annotations",
//		"metadata.labels",
//		"metadata.name",
//		"metadata.namespace",
//		"metadata.uid",
//		"spec.nodeName",
//		"spec.restartPolicy",
//		"spec.serviceAccountName",
//		"spec.schedulerName",
//		"status.phase",
//		"status.hostIP",
//		"status.podIP":
//		return label, value, nil
//	// This is for backwards compatibility with old v1 clients which send spec.host
//	case "spec.host":
//		return "spec.nodeName", value, nil
//	default:
//		return "", "", fmt.Errorf("field label not supported: %s", label)
//	}
//}
