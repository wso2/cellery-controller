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
	"fmt"
	"math/big"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller/composite/config"
)

func CreateKeyPairSecret(composite *v1alpha1.Composite, cellerySecret config.Secret) (*corev1.Secret, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, fmt.Errorf("fail to generate rsa private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:         "composite",
			Country:            []string{"LK"},
			Locality:           []string{"Colombo"},
			Organization:       []string{"WSO2"},
			OrganizationalUnit: []string{"WSO2"},
			Province:           []string{"West"},
		},
		DNSNames:              []string{fmt.Sprintf("%s--sts-service", "composite"), "composite"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	err = cellerySecret.Validate()
	if err != nil {
		return nil, fmt.Errorf("cellery system secret validation failed: %v", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, cellerySecret.Certificate, privateKey.Public(), cellerySecret.PrivateKey)

	if err != nil {
		return nil, fmt.Errorf("fail to sign cell certificate: %v", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecretName(composite),
			Namespace: composite.Namespace,
			Labels:    createLabels(composite),
		},
		Type: mesh.GroupName + "/key-and-cert",
		Data: map[string][]byte{
			"key.pem":          pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}),
			"cert.pem":         pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes}),
			"cellery-cert.pem": pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cellerySecret.Certificate.Raw}),
			"cert-bundle.pem":  cellerySecret.CertBundle,
		},
	}, nil
}
