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

package gateway

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	istionetworkingv1alpha3 "github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/controller/gateway/resources"
	"github.com/cellery-io/mesh-controller/pkg/meta"
)

func (r *reconciler) reconcileApiPublisherConfigMap(gateway *v1alpha2.Gateway) error {
	configMapName := resources.ApiPublisherConfigMap(gateway)
	configMap, err := r.configMapLister.ConfigMaps(gateway.Namespace).Get(configMapName)

	if !resources.IsApiPublishingRequired(gateway) {
		if err == nil && metav1.IsControlledBy(configMap, gateway) {
			err = r.kubeClient.BatchV1().Jobs(gateway.Namespace).Delete(configMapName, meta.DeleteWithPropagationBackground())
			if err != nil {
				r.logger.Errorf("Failed to delete api publisher config map %q: %v", configMapName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		desiredConfigMap, err := resources.CreateGatewayConfigMap(gateway, r.cfg)
		configMap, err = r.kubeClient.CoreV1().ConfigMaps(gateway.Namespace).Create(desiredConfigMap)
		if err != nil {
			r.logger.Errorf("Failed to create api publisher ConfigMap %q: %v", configMapName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create api publisher ConfigMap %q: %v", configMapName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created Api Publisher ConfigMap %q",
			configMapName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve api publisher ConfigMap %q: %v", configMapName, err)
		return err
	} else if !metav1.IsControlledBy(configMap, gateway) {
		return fmt.Errorf("gateway: %q does not own the api publisher ConfigMap: %q", gateway.Name, configMapName)
	} else {
		configMap, err = func(gateway *v1alpha2.Gateway, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
			if !resources.RequireGatewayConfigMapUpdate(gateway, configMap) {
				return configMap, nil
			}
			desiredConfigMap, err := resources.CreateGatewayConfigMap(gateway, r.cfg)
			if err != nil {
				return nil, err
			}
			existingConfigMap := configMap.DeepCopy()
			resources.CopyGatewayConfigMap(desiredConfigMap, existingConfigMap)
			return r.kubeClient.CoreV1().ConfigMaps(gateway.Namespace).Update(existingConfigMap)
		}(gateway, configMap)
		if err != nil {
			r.logger.Errorf("Failed to update api publisher ConfigMap %q: %v", configMapName, err)
			return err
		}
	}
	resources.StatusFromConfigMap(gateway, configMap)
	return nil
}

func (r *reconciler) reconcileApiPublisherJob(gateway *v1alpha2.Gateway) error {
	jobName := resources.JobName(gateway)
	job, err := r.jobLister.Jobs(gateway.Namespace).Get(jobName)
	if !resources.RequireApiPublisherJob(gateway) {
		if err == nil && metav1.IsControlledBy(job, gateway) {
			err = r.kubeClient.BatchV1().Jobs(gateway.Namespace).Delete(jobName, meta.DeleteWithPropagationBackground())
			if err != nil {
				r.logger.Errorf("Failed to delete api publisher Job %q: %v", jobName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		job, err = r.kubeClient.BatchV1().Jobs(gateway.Namespace).Create(resources.MakeApiPublisherJob(gateway, r.cfg))
		if err != nil {
			r.logger.Errorf("Failed to create api publisher Job %q: %v", jobName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create api publisher Job %q: %v", jobName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created api publisher Job %q", jobName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve api publisher Job %q: %v", jobName, err)
		return err
	} else if !metav1.IsControlledBy(job, gateway) {
		return fmt.Errorf("component: %q does not own the api publisher Job: %q", gateway.Name, jobName)
	} else {
		if resources.RequireApiPublisherJobUpdate(gateway, job) {
			err = r.kubeClient.BatchV1().Jobs(gateway.Namespace).Delete(jobName, meta.DeleteWithPropagationBackground())
			if err != nil {
				r.logger.Errorf("Failed to delete api publisher Job %q: %v", jobName, err)
				return err
			}
		}
	}
	resources.StatusFromApiPublisherJob(gateway, job)
	return nil
}

func (r *reconciler) reconcileClusterIngress(gateway *v1alpha2.Gateway) error {
	ingressName := resources.ClusterIngressName(gateway)
	ingress, err := r.clusterIngressLister.Ingresses(gateway.Namespace).Get(ingressName)
	if !resources.RequireClusterIngress(gateway) {
		if err == nil && metav1.IsControlledBy(ingress, gateway) {
			err = r.kubeClient.ExtensionsV1beta1().Ingresses(gateway.Namespace).Delete(ingressName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete Ingress %q: %v", ingressName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		ingress, err = r.kubeClient.ExtensionsV1beta1().Ingresses(gateway.Namespace).Create(resources.MakeClusterIngress(gateway))
		if err != nil {
			r.logger.Errorf("Failed to create Ingress %q: %v", ingressName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create Ingress %q: %v", ingressName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created Ingress %q", ingressName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve Ingress %q: %v", ingressName, err)
		return err
	} else if !metav1.IsControlledBy(ingress, gateway) {
		return fmt.Errorf("gateway: %q does not own the Ingress: %q", gateway.Name, ingressName)
	} else {
		ingress, err = func(gateway *v1alpha2.Gateway, ingress *extensionsv1beta1.Ingress) (*extensionsv1beta1.Ingress, error) {
			if !resources.RequireClusterIngressUpdate(gateway, ingress) {
				return ingress, nil
			}
			desiredIngress := resources.MakeClusterIngress(gateway)
			if err != nil {
				return nil, err
			}
			existingIngress := ingress.DeepCopy()
			resources.CopyClusterIngress(desiredIngress, existingIngress)
			return r.kubeClient.ExtensionsV1beta1().Ingresses(gateway.Namespace).Update(existingIngress)
		}(gateway, ingress)
		if err != nil {
			r.logger.Errorf("Failed to update Ingress %q: %v", ingressName, err)
			return err
		}
	}
	resources.StatusFromClusterIngress(gateway, ingress)
	return nil
}

func (r *reconciler) reconcileClusterIngressSecret(gateway *v1alpha2.Gateway) error {
	secretName := resources.ClusterIngressSecretName(gateway)
	secret, err := r.secretLister.Secrets(gateway.Namespace).Get(secretName)
	if !resources.RequireClusterIngressSecret(gateway) {
		if err == nil && metav1.IsControlledBy(secret, gateway) {
			err = r.kubeClient.CoreV1().Secrets(gateway.Namespace).Delete(secretName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete ingress Secret %q: %v", secretName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		secret, err = func(gateway *v1alpha2.Gateway) (*corev1.Secret, error) {
			desiredSecret, err := resources.MakeClusterIngressSecret(gateway, r.cfg)
			if err != nil {
				return nil, err
			}
			return r.kubeClient.CoreV1().Secrets(gateway.Namespace).Create(desiredSecret)
		}(gateway)
		if err != nil {
			r.logger.Errorf("Failed to create ingress Secret %q: %v", secretName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create ingress Secret %q: %v", secretName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created ingress Secret %q", secretName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve ingress Secret %q: %v", secretName, err)
		return err
	} else if !metav1.IsControlledBy(secret, gateway) {
		return fmt.Errorf("gateway: %q does not own the ingress Secret: %q", gateway.Name, secretName)
	} else {
		secret, err = func(gateway *v1alpha2.Gateway, secret *corev1.Secret) (*corev1.Secret, error) {
			if !resources.RequireClusterIngressSecretUpdate(gateway, secret) {
				return secret, nil
			}
			desiredSecret, err := resources.MakeClusterIngressSecret(gateway, r.cfg)
			if err != nil {
				return nil, err
			}
			existingSecret := secret.DeepCopy()
			resources.CopyClusterIngressSecret(desiredSecret, existingSecret)
			return r.kubeClient.CoreV1().Secrets(gateway.Namespace).Update(existingSecret)
		}(gateway, secret)
		if err != nil {
			r.logger.Errorf("Failed to update ingress Secret %q: %v", secretName, err)
			return err
		}
	}
	resources.StatusFromClusterIngressSecret(gateway, secret)
	return nil
}

func (r *reconciler) reconcileOidcEnvoyFilter(gateway *v1alpha2.Gateway) error {
	envoyFilterName := resources.OidcEnvoyFilterName(gateway)
	envoyFilter, err := r.istioEnvoyFilterLister.EnvoyFilters(gateway.Namespace).Get(envoyFilterName)
	if !resources.RequireOidcEnvoyFilter(gateway) {
		if err == nil && metav1.IsControlledBy(envoyFilter, gateway) {
			err = r.meshClient.NetworkingV1alpha3().EnvoyFilters(gateway.Namespace).Delete(envoyFilterName, &metav1.DeleteOptions{})
			if err != nil {
				r.logger.Errorf("Failed to delete oidc EnvoyFilter %q: %v", envoyFilterName, err)
				return err
			}
		}
		return nil
	}

	if errors.IsNotFound(err) {
		envoyFilter, err = r.meshClient.NetworkingV1alpha3().EnvoyFilters(gateway.Namespace).Create(resources.MakeOidcEnvoyFilter(gateway))
		if err != nil {
			r.logger.Errorf("Failed to create oidc EnvoyFilter %q: %v", envoyFilterName, err)
			r.recorder.Eventf(gateway, corev1.EventTypeWarning, "CreationFailed", "Failed to create oidc EnvoyFilter %q: %v", envoyFilterName, err)
			return err
		}
		r.recorder.Eventf(gateway, corev1.EventTypeNormal, "Created", "Created oidc EnvoyFilter %q", envoyFilterName)
	} else if err != nil {
		r.logger.Errorf("Failed to retrieve oidc EnvoyFilter %q: %v", envoyFilterName, err)
		return err
	} else if !metav1.IsControlledBy(envoyFilter, gateway) {
		return fmt.Errorf("gateway: %q does not own the oidc EnvoyFilter: %q", gateway.Name, envoyFilterName)
	} else {
		envoyFilter, err = func(gateway *v1alpha2.Gateway, envoyFilter *istionetworkingv1alpha3.EnvoyFilter) (*istionetworkingv1alpha3.EnvoyFilter, error) {
			if !resources.RequireOidcEnvoyFilterUpdate(gateway, envoyFilter) {
				return envoyFilter, nil
			}
			desiredEnvoyFilter := resources.MakeOidcEnvoyFilter(gateway)
			if err != nil {
				return nil, err
			}
			existingEnvoyFilter := envoyFilter.DeepCopy()
			resources.CopyOidcEnvoyFilter(desiredEnvoyFilter, existingEnvoyFilter)
			return r.meshClient.NetworkingV1alpha3().EnvoyFilters(gateway.Namespace).Update(existingEnvoyFilter)
		}(gateway, envoyFilter)
		if err != nil {
			r.logger.Errorf("Failed to update oidc EnvoyFilter %q: %v", envoyFilterName, err)
			return err
		}
	}
	resources.StatusFromOidcEnvoyFilter(gateway, envoyFilter)
	return nil
}
