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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/gateway/config"
	"github.com/cellery-io/mesh-controller/pkg/crypto"
)

func CreateClusterIngressSecret(gateway *v1alpha1.Gateway, secret config.Secret) (*corev1.Secret, error) {

	key, err := crypto.TryDecrypt(gateway.Spec.Tls.Key, secret.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt the tls.key: %v", err)
	}

	cert, err := crypto.TryDecrypt(gateway.Spec.Tls.Cert, secret.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt the tls.cert: %v", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterIngressSecretName(gateway),
			Namespace: gateway.Namespace,
			Labels:    createGatewayLabels(gateway),
			OwnerReferences: []metav1.OwnerReference{
				*controller.CreateGatewayOwnerRef(gateway),
			},
		},
		Type: corev1.SecretTypeTLS,
		Data: map[string][]byte{
			"tls.key": key,
			"tls.crt": cert,
		},
	}, nil
}
