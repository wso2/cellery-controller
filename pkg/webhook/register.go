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

package webhook

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/ptr"
)

func (s *server) registerWebhooks(caCertPEM []byte) error {
	deployment, err := s.kubeClient.AppsV1().Deployments(s.options.Namespace).Get(s.options.DeploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch the webhook deployment: %v", err)
	}
	ownerRef := metav1.NewControllerRef(deployment, appsv1.SchemeGroupVersion.WithKind("Deployment"))

	if err := s.registerMutatingWebhook(ownerRef, caCertPEM); err != nil {
		return fmt.Errorf("mutating webhook registration failed: %v", err)
	}

	if err := s.registerValidatingWebhook(ownerRef, caCertPEM); err != nil {
		return fmt.Errorf("validating webhook registration failed: %v", err)
	}
	return nil
}

func (s *server) registerMutatingWebhook(ownerRef *metav1.OwnerReference, caCertPEM []byte) error {
	_, err := s.kubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Get(s.options.MutatingWebhookName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = s.kubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Create(s.makeMutatingWebhookConfiguration(ownerRef, caCertPEM))
		if apierrors.IsAlreadyExists(err) {
			_, err = s.kubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Get(s.options.MutatingWebhookName, metav1.GetOptions{})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (s *server) registerValidatingWebhook(ownerRef *metav1.OwnerReference, caCertPEM []byte) error {
	_, err := s.kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Get(s.options.ValidatingWebhookName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = s.kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Create(s.makeValidatingWebhookConfiguration(ownerRef, caCertPEM))
		if apierrors.IsAlreadyExists(err) {
			_, err = s.kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Get(s.options.ValidatingWebhookName, metav1.GetOptions{})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (s *server) makeMutatingWebhookConfiguration(ownerRef *metav1.OwnerReference, caCertPEM []byte) *admissionregistrationv1beta1.MutatingWebhookConfiguration {
	var resources []schema.GroupVersionResource
	for gvk := range s.defaulters {
		plural := fmt.Sprintf("%ss", strings.ToLower(gvk.Kind))
		resources = append(resources, gvk.GroupVersion().WithResource(plural))
	}
	return &admissionregistrationv1beta1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:            s.options.MutatingWebhookName,
			OwnerReferences: []metav1.OwnerReference{*ownerRef},
		},
		Webhooks: []admissionregistrationv1beta1.MutatingWebhook{
			{
				Name:                    s.options.MutatingWebhookName,
				Rules:                   makeWebhookRules(resources),
				ClientConfig:            makeClientConfig(s.options.ServiceName, s.options.Namespace, pathMutate, caCertPEM),
				FailurePolicy:           failurePolicy(),
				AdmissionReviewVersions: admissionReviewVersions(),
			},
		},
	}
}

func (s *server) makeValidatingWebhookConfiguration(ownerRef *metav1.OwnerReference, caCertPEM []byte) *admissionregistrationv1beta1.ValidatingWebhookConfiguration {
	var resources []schema.GroupVersionResource
	for gvk := range s.validators {
		plural := fmt.Sprintf("%ss", strings.ToLower(gvk.Kind))
		resources = append(resources, gvk.GroupVersion().WithResource(plural))
	}
	return &admissionregistrationv1beta1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:            s.options.ValidatingWebhookName,
			OwnerReferences: []metav1.OwnerReference{*ownerRef},
		},
		Webhooks: []admissionregistrationv1beta1.ValidatingWebhook{
			{
				Name:                    s.options.ValidatingWebhookName,
				Rules:                   makeWebhookRules(resources),
				ClientConfig:            makeClientConfig(s.options.ServiceName, s.options.Namespace, pathValidate, caCertPEM),
				FailurePolicy:           failurePolicy(),
				AdmissionReviewVersions: admissionReviewVersions(),
			},
		},
	}
}

func makeWebhookRules(resources []schema.GroupVersionResource) []admissionregistrationv1beta1.RuleWithOperations {
	scope := admissionregistrationv1beta1.NamespacedScope
	var rules []admissionregistrationv1beta1.RuleWithOperations
	for _, gvr := range resources {
		rules = append(rules, admissionregistrationv1beta1.RuleWithOperations{
			Operations: []admissionregistrationv1beta1.OperationType{
				admissionregistrationv1beta1.Create,
				admissionregistrationv1beta1.Update,
			},
			Rule: admissionregistrationv1beta1.Rule{
				APIGroups:   []string{gvr.Group},
				APIVersions: []string{gvr.Version},
				Resources:   []string{gvr.Resource},
				Scope:       &scope,
			},
		})
	}
	return rules
}

func makeClientConfig(svcName, namespace, path string, caCertPEM []byte) admissionregistrationv1beta1.WebhookClientConfig {
	return admissionregistrationv1beta1.WebhookClientConfig{
		Service: &admissionregistrationv1beta1.ServiceReference{
			Namespace: namespace,
			Name:      svcName,
			Path:      ptr.String(path),
		},
		CABundle: caCertPEM,
	}
}

func admissionReviewVersions() []string {
	return []string{"v1beta1"}
}

func failurePolicy() *admissionregistrationv1beta1.FailurePolicyType {
	failurePolicy := admissionregistrationv1beta1.Fail
	return &failurePolicy
}
