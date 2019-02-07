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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/gateway/config"
)

type apiConfig struct {
	Cell      string              `json:"cell"`
	Version   string              `json:"version"`
	Hostname  string              `json:"hostname"`
	APIRoutes []v1alpha1.APIRoute `json:"apis"`
}

func CreateGatewayConfigMap(gateway *v1alpha1.Gateway, gatewayConfig config.Gateway) (*corev1.ConfigMap, error) {
	var cellName string
	cellName, ok := gateway.Labels[mesh.CellLabelKey]
	if !ok {
		cellName = gateway.Name
	}

	api := &apiConfig{
		Cell:      cellName,
		Version:   "1.0.0",
		Hostname:  GatewayFullK8sServiceName(gateway),
		APIRoutes: gateway.Spec.APIRoutes,
	}
	apiConfigJsonBytes, err := json.Marshal(api)
	if err != nil {
		return nil, fmt.Errorf("cannot create apiConfig json for the ConfigMap %q: %v",
			GatewayConfigMapName(gateway), err)
	}
	apiConfigJson := string(apiConfigJsonBytes)
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GatewayConfigMapName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Data: map[string]string{
			apiConfigKey:          apiConfigJson,
			gatewayConfigKey:      gatewayConfig.InitConfig,
			gatewaySetupConfigKey: gatewayConfig.SetupConfig,
		},
	}, nil
}
