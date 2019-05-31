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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller/cell/config"
)

func TestCreateKeyPairSecret(t *testing.T) {
	cell := &v1alpha1.Cell{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo-namespace",
			Name:      "foo",
		},
		Spec: v1alpha1.CellSpec{
			ServiceTemplates: []v1alpha1.ServiceTemplateSpec{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "bar-service",
					},
				},
			},
		},
	}

	sec := config.Secret{
		PrivateKey: func() *rsa.PrivateKey {
			key, _ := rsa.GenerateKey(rand.Reader, 2048)
			return key
		}(),
		Certificate: func() *x509.Certificate {
			return &x509.Certificate{}
		}(),
	}

	if _, err := CreateKeyPairSecret(cell, sec); err != nil {
		t.Errorf("Secret creation failed with error: %v", err)
	}
}
