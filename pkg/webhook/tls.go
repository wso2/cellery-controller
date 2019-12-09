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
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cellery.io/cellery-controller/pkg/apis/mesh"
	"cellery.io/cellery-controller/pkg/crypto"
)

const (
	tlsKey  = "tls.key"
	tlsCert = "tls.crt"
	caCert  = "ca.crt"
)

func (s *server) configureTls() (*tls.Config, []byte, error) {
	secret, err := s.kubeClient.CoreV1().Secrets(s.options.Namespace).Get(s.options.ServerSecretName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		newSecret, err := s.generateSecret()
		if err != nil {
			return nil, nil, err
		}
		secret, err = s.kubeClient.CoreV1().Secrets(newSecret.Namespace).Create(newSecret)
		if apierrors.IsAlreadyExists(err) {
			secret, err = s.kubeClient.CoreV1().Secrets(s.options.Namespace).Get(s.options.ServerSecretName, metav1.GetOptions{})
			if err != nil {
				return nil, nil, err
			}
		} else if err != nil {
			return nil, nil, err
		}
	} else if err != nil {
		return nil, nil, err
	}

	var serverKeyPEM, serverCertPEM, rootCertPEM []byte
	var ok bool
	if serverKeyPEM, ok = secret.Data[tlsKey]; !ok {
		return nil, nil, fmt.Errorf("missing key %q in secret %q", tlsKey, s.options.ServerSecretName)
	}
	if serverCertPEM, ok = secret.Data[tlsCert]; !ok {
		return nil, nil, fmt.Errorf("missing key %q in secret %q", tlsCert, s.options.ServerSecretName)
	}
	if rootCertPEM, ok = secret.Data[caCert]; !ok {
		return nil, nil, fmt.Errorf("missing key %q in secret %q", caCert, s.options.ServerSecretName)
	}

	cert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		return nil, nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, rootCertPEM, nil
}

func (s *server) generateSecret() (*corev1.Secret, error) {
	rootSecret, err := s.kubeClient.CoreV1().Secrets(s.options.Namespace).Get(s.options.RootSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting root certificates: %v", err)
	}
	var rootKeyPEM, rootCertPEM []byte
	var ok bool
	if rootKeyPEM, ok = rootSecret.Data[tlsKey]; !ok {
		return nil, fmt.Errorf("missing key %q in secret %q", tlsKey, s.options.RootSecretName)
	}
	if rootCertPEM, ok = rootSecret.Data[tlsCert]; !ok {
		return nil, fmt.Errorf("missing key %q in secret %q", tlsCert, s.options.RootSecretName)
	}

	rootKey, err := crypto.ParsePrivateKey(rootKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("error while parsing root key %q from secret %q", tlsKey, s.options.RootSecretName)
	}
	rootCert, err := crypto.ParseCertificate(rootCertPEM)
	if err != nil {
		return nil, fmt.Errorf("error while parsing root cert %q from secret %q", tlsCert, s.options.RootSecretName)
	}

	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, fmt.Errorf("fail to generate rsa server private key: %v", err)
	}

	svcName := s.options.ServiceName
	svcNamespace := s.options.Namespace
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"cellery.io"},
		},
		DNSNames: []string{
			svcName,
			fmt.Sprintf("%s.%s", svcName, svcNamespace),
			fmt.Sprintf("%s.%s.svc", svcName, svcNamespace),
			fmt.Sprintf("%s.%s.svc.cluster.local", svcName, svcNamespace),
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	serverCertDER, err := x509.CreateCertificate(rand.Reader, &template, rootCert, serverKey.Public(), rootKey)
	if err != nil {
		return nil, fmt.Errorf("fail to sign server certificate: %v", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.options.ServerSecretName,
			Namespace: s.options.Namespace,
		},
		Type: mesh.GroupName + "/key-and-cert",
		Data: map[string][]byte{
			tlsKey:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)}),
			tlsCert: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCertDER}),
			caCert:  rootCertPEM,
		},
	}, nil
}
