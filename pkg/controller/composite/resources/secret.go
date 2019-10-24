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
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/config"
)

func MakeSecret(composite *v1alpha2.Composite, cfg config.Interface) (*corev1.Secret, error) {

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
		DNSNames:              []string{fmt.Sprintf("%s--sts-service", "composite"), "composite", fmt.Sprintf("%s--sts-service.%s", "composite", mesh.SystemNamespace)},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	keySystem, err := cfg.PrivateKey()

	if err != nil {
		return nil, err
	}
	certSystem, err := cfg.Certificate()
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, certSystem, privateKey.Public(), keySystem)

	if err != nil {
		return nil, fmt.Errorf("fail to sign cell certificate: %v", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecretName(composite),
			Namespace: mesh.SystemNamespace,
			Labels:    makeLabels(composite),
		},
		Type: mesh.GroupName + "/key-and-cert",
		Data: map[string][]byte{
			"key.pem":          pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}),
			"cert.pem":         pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes}),
			"cellery-cert.pem": pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certSystem.Raw}),
			"cert-bundle.pem":  cfg.CertificateBundle(),
		},
	}, nil
}

func StatusFromSecret(composite *v1alpha2.Composite, secret *corev1.Secret) {
	composite.Status.SecretGeneration = secret.Generation
}
