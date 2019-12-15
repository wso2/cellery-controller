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

package config

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strconv"
	"sync"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"cellery.io/cellery-controller/pkg/crypto"
	"cellery.io/cellery-controller/pkg/informers"
)

type Interface interface {
	Value(key string) (string, bool)
	StringValue(key string) string
	BoolValue(key string) bool
	IntValue(key string) int64
	PrivateKey() (*rsa.PrivateKey, error)
	Certificate() (*x509.Certificate, error)
	CertificateBundle() []byte
}

type config struct {
	configMapName   string
	secretName      string
	namespace       string
	configMapLister corev1listers.ConfigMapLister
	secretLister    corev1listers.SecretLister

	rwlock      sync.RWMutex
	configData  map[string]string
	privateKey  *rsa.PrivateKey
	certificate *x509.Certificate
	certBundle  []byte

	logger *zap.SugaredLogger
}

func NewWatcher(inf informers.Interface, configMapName string, secretName string, namespace string, logger *zap.SugaredLogger) *config {
	cfg := &config{
		configMapName:   configMapName,
		secretName:      secretName,
		namespace:       namespace,
		configMapLister: inf.ConfigMaps().Lister(),
		secretLister:    inf.Secrets().Lister(),
		logger:          logger.Named("config-watcher"),
	}
	inf.ConfigMaps().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithNameAndNamespace(cfg.configMapName, cfg.namespace),
		Handler:    informers.HandleAddUpdate(cfg.updateConfigs),
	})

	inf.Secrets().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: informers.FilterWithNameAndNamespace(cfg.secretName, cfg.namespace),
		Handler:    informers.HandleAddUpdate(cfg.updateSecrets),
	})
	return cfg
}

func (c *config) CheckResources() error {
	if _, err := c.configMapLister.ConfigMaps(c.namespace).Get(c.configMapName); err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("configmap %q is not available in the %q namespace", c.configMapName, c.namespace)
		}
		return err
	}

	if _, err := c.secretLister.Secrets(c.namespace).Get(c.secretName); err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("secret %q is not available in the %q namespace", c.secretName, c.namespace)
		}
		return err
	}
	return nil
}

func (c *config) Value(key string) (string, bool) {
	c.rwlock.RLock()
	defer c.rwlock.RUnlock()
	v, ok := c.configData[key]
	return v, ok
}

func (c *config) StringValue(key string) string {
	def := ""
	if v, ok := c.Value(key); ok {
		return v
	}
	c.logger.Warnf("No configuration is found for key %q in %s/%s. Defaulting to %s",
		key, c.configMapName, c.namespace, def)
	return def
}

func (c *config) BoolValue(key string) bool {
	def := false
	if v, ok := c.Value(key); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			c.logger.Warnf("Error while parsing bool value from key %q: %v. Defaulting to %t", key, err, def)
			return def
		}
		return b
	}
	c.logger.Warnf("No configuration is found for key %q in %s/%s. Defaulting to %t",
		key, c.configMapName, c.namespace, def)
	return def
}

func (c *config) IntValue(key string) int64 {
	var def int64 = 0
	if v, ok := c.Value(key); ok {
		i64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.logger.Warnf("Error while parsing int value from key %q: %v. Defaulting to %d", key, err, def)
			return def
		}
		return i64
	}
	c.logger.Warnf("No configuration is found for key %q in %s/%s. Defaulting to %d",
		key, c.configMapName, c.namespace, def)
	return def
}

func (c *config) PrivateKey() (*rsa.PrivateKey, error) {
	c.rwlock.RLock()
	defer c.rwlock.RUnlock()
	if c.privateKey == nil {
		return nil, fmt.Errorf("no rsa private key is set with the key %q in secret %s/%s",
			SecretKeyPrivateKey, c.namespace, c.secretName)
	}
	return c.privateKey, nil
}

func (c *config) Certificate() (*x509.Certificate, error) {
	c.rwlock.RLock()
	defer c.rwlock.RUnlock()
	if c.certificate == nil {
		return nil, fmt.Errorf("no x509 certificate is set with the key %q in secret %s/%s",
			SecretKeyCertificate, c.namespace, c.secretName)
	}
	return c.certificate, nil
}

func (c *config) CertificateBundle() []byte {
	c.rwlock.RLock()
	defer c.rwlock.RUnlock()
	return c.certBundle
}

func (c *config) updateConfigs(obj interface{}) {
	c.rwlock.Lock()
	defer c.rwlock.Unlock()

	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return
	}
	c.configData = configMap.DeepCopy().Data
}

func (c *config) updateSecrets(obj interface{}) {
	c.rwlock.Lock()
	defer c.rwlock.Unlock()

	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return
	}

	if keyBytes, ok := secret.Data[SecretKeyPrivateKey]; ok {
		privateKey, err := crypto.ParsePrivateKey(keyBytes)
		if err != nil {
			c.logger.Errorf("Error while parsing %q from the secret %s/%s: %v", SecretKeyPrivateKey, c.namespace, c.secretName, err)
		} else {
			c.privateKey = privateKey
		}
	} else {
		c.logger.Errorf("Missing key %q in secret %s/%s", SecretKeyPrivateKey, c.namespace, c.secretName)
	}

	if certBytes, ok := secret.Data[SecretKeyCertificate]; ok {
		cert, err := crypto.ParseCertificate(certBytes)
		if err != nil {
			c.logger.Errorf("Error while parsing %q from the secret %s/%s: %v", SecretKeyCertificate, c.namespace, c.secretName, err)
		} else {
			c.certificate = cert
		}
	} else {
		c.logger.Errorf("Missing key %q in secret %s/%s", SecretKeyCertificate, c.namespace, c.secretName)
	}

	if certBundle, ok := secret.Data[SecretKeyCertificateBundle]; ok {
		c.certBundle = certBundle
	} else {
		c.logger.Errorf("Missing key %q in secret %s/%s", SecretKeyCertificateBundle, c.namespace, c.secretName)
	}

}
