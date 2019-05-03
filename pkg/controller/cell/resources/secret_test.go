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
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
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

	if CreateKeyPairSecret(cell) == nil {
		t.Errorf("Secret for cell with spec: %v not created properly", cell.Spec)
	}
}

func createPrivateKeyAndCert(cell *v1alpha1.Cell) (*pem.Block, *pem.Block) {

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{cell.Name},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.Public(), privateKey)

	return &pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
		},
		&pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
}
