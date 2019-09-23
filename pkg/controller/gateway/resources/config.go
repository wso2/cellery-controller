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
	"encoding/json"
	"fmt"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/meta"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/config"
	"github.com/cellery-io/mesh-controller/pkg/controller"
)

type apiConfig struct {
	Cell          string               `json:"cell"`
	Version       string               `json:"version"`
	Hostname      string               `json:"hostname"`
	HTTPRoutes    []v1alpha2.HTTPRoute `json:"apis"`
	GlobalContext string               `json:"globalContext"`
}

func RequireConfigMap(gateway *v1alpha2.Gateway) bool {
	return gateway.Spec.Ingress.IngressExtensions.HasApiPublisher()
}

func CreateGatewayConfigMap(gateway *v1alpha2.Gateway, cfg config.Interface) (*corev1.ConfigMap, error) {
	globalContext := ""
	version := "0.1"
	cellName, ok := gateway.Labels[meta.CellLabelKey]
	if !ok {
		cellName = gateway.Name
	}
	if gateway.Spec.Ingress.IngressExtensions.HasApiPublisher() {
		globalContext = gateway.Spec.Ingress.IngressExtensions.ApiPublisher.Context
		if gateway.Spec.Ingress.IngressExtensions.ApiPublisher.HasVersion() {
			version = gateway.Spec.Ingress.IngressExtensions.ApiPublisher.Version
		}
	}
	api := &apiConfig{
		Cell:          cellName,
		Version:       version,
		Hostname:      GatewayFullK8sServiceName(gateway),
		HTTPRoutes:    gateway.Spec.Ingress.HTTPRoutes,
		GlobalContext: globalContext,
	}
	apiConfigJsonBytes, err := json.Marshal(api)
	if err != nil {
		return nil, fmt.Errorf("cannot create apiConfig json for the ConfigMap %q: %v",
			ApiPublisherConfigMap(gateway), err)
	}
	apiConfigJson := string(apiConfigJsonBytes)
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ApiPublisherConfigMap(gateway),
			Namespace: gateway.Namespace,
			Labels:    makeLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Data: map[string]string{
			apiConfigKey:     apiConfigJson,
			apiPublisherConfigKey: cfg.StringValue(config.ConfigMapKeyApiPublisherConfig),
		},
	}, nil
}

func StatusFromConfigMap(gateway *v1alpha2.Gateway, configMap *corev1.ConfigMap) {
	gateway.Status.ConfigMapGeneration = configMap.Generation
}
