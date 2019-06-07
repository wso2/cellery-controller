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

package cell

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	corev1informers "k8s.io/client-go/informers/core/v1"
	networkv1informers "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	networkv1listers "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	meshclientset "github.com/cellery-io/mesh-controller/pkg/client/clientset/versioned"
	meshinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/mesh/v1alpha1"
	istioinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/networking/v1alpha3"
	listers "github.com/cellery-io/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	istiov1alpha1listers "github.com/cellery-io/mesh-controller/pkg/client/listers/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/cell/config"
	"github.com/cellery-io/mesh-controller/pkg/controller/cell/resources"
)

type cellHandler struct {
	kubeClient          kubernetes.Interface
	meshClient          meshclientset.Interface
	networkPilicyLister networkv1listers.NetworkPolicyLister
	secretLister        corev1listers.SecretLister
	cellLister          listers.CellLister
	gatewayLister       listers.GatewayLister
	tokenServiceLister  listers.TokenServiceLister
	serviceLister       listers.ServiceLister
	envoyFilterLister   istiov1alpha1listers.EnvoyFilterLister
	cellerySecret       config.Secret
	logger              *zap.SugaredLogger
}

func NewController(
	kubeClient kubernetes.Interface,
	meshClient meshclientset.Interface,
	cellInformer meshinformers.CellInformer,
	gatewayInformer meshinformers.GatewayInformer,
	tokenServiceInformer meshinformers.TokenServiceInformer,
	serviceInformer meshinformers.ServiceInformer,
	networkPolicyInformer networkv1informers.NetworkPolicyInformer,
	secretInformer corev1informers.SecretInformer,
	envoyFilterInformer istioinformers.EnvoyFilterInformer,
	systemSecretInformer corev1informers.SecretInformer,
	logger *zap.SugaredLogger,
) *controller.Controller {
	h := &cellHandler{
		kubeClient:          kubeClient,
		meshClient:          meshClient,
		cellLister:          cellInformer.Lister(),
		serviceLister:       serviceInformer.Lister(),
		gatewayLister:       gatewayInformer.Lister(),
		tokenServiceLister:  tokenServiceInformer.Lister(),
		networkPilicyLister: networkPolicyInformer.Lister(),
		secretLister:        secretInformer.Lister(),
		envoyFilterLister:   envoyFilterInformer.Lister(),
		logger:              logger.Named("cell-controller"),
	}
	c := controller.New(h, h.logger, "Cell")

	h.logger.Info("Setting up event handlers")
	cellInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			h.logger.Debugw("Informer update", "old", old, "new", new)
			c.Enqueue(new)
		},
		DeleteFunc: c.Enqueue,
	})

	systemSecretInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: h.updateSecret,
		UpdateFunc: func(old, new interface{}) {
			h.updateSecret(new)
		},
	})

	return c
}

func (h *cellHandler) Handle(key string) error {
	h.logger.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		h.logger.Errorf("invalid resource key: %s", key)
		return nil
	}
	cellOriginal, err := h.cellLister.Cells(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("cell '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	h.logger.Debugw("lister instance", key, cellOriginal)
	cell := cellOriginal.DeepCopy()

	if err = h.handle(cell); err != nil {
		return err
	}

	if _, err = h.updateStatus(cell); err != nil {
		return err
	}
	return nil
}

func (h *cellHandler) handle(cell *v1alpha1.Cell) error {

	if err := h.handleNetworkPolicy(cell); err != nil {
		return err
	}

	if err := h.handleSecret(cell); err != nil {
		return err
	}

	if err := h.handleGateway(cell); err != nil {
		return err
	}

	if err := h.handleTokenService(cell); err != nil {
		return err
	}

	if err := h.handleServices(cell); err != nil {
		return err
	}

	h.updateCellStatus(cell)
	return nil
}

func (h *cellHandler) handleNetworkPolicy(cell *v1alpha1.Cell) error {
	networkPolicy, err := h.networkPilicyLister.NetworkPolicies(cell.Namespace).Get(resources.NetworkPolicyName(cell))
	if errors.IsNotFound(err) {
		networkPolicy, err = h.kubeClient.NetworkingV1().NetworkPolicies(cell.Namespace).Create(resources.CreateNetworkPolicy(cell))
		if err != nil {
			h.logger.Errorf("Failed to create NetworkPolicy %v", err)
			return err
		}
		h.logger.Debugw("NetworkPolicy created", resources.NetworkPolicyName(cell), networkPolicy)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *cellHandler) handleSecret(cell *v1alpha1.Cell) error {
	secret, err := h.secretLister.Secrets(cell.Namespace).Get(resources.SecretName(cell))
	if errors.IsNotFound(err) {
		desiredSecret, err := resources.CreateKeyPairSecret(cell, h.cellerySecret)
		if err != nil {
			h.logger.Errorf("Cannot build the Cell Secret %v", err)
			return err
		}
		secret, err = h.kubeClient.CoreV1().Secrets(cell.Namespace).Create(desiredSecret)
		if err != nil {
			h.logger.Errorf("Failed to create cell Secret %v", err)
			return err
		}
		h.logger.Debugw("Secret created", resources.SecretName(cell), secret)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *cellHandler) handleGateway(cell *v1alpha1.Cell) error {
	gateway, err := h.gatewayLister.Gateways(cell.Namespace).Get(resources.GatewayName(cell))
	if errors.IsNotFound(err) {
		gateway, err = h.meshClient.MeshV1alpha1().Gateways(cell.Namespace).Create(resources.CreateGateway(cell))
		if err != nil {
			h.logger.Errorf("Failed to create Gateway %v", err)
			return err
		}
		h.logger.Debugw("Gateway created", resources.GatewayName(cell), gateway)
	} else if err != nil {
		return err
	}

	cell.Status.GatewayHostname = gateway.Status.HostName
	cell.Status.GatewayStatus = gateway.Status.Status
	return nil
}

func (h *cellHandler) handleTokenService(cell *v1alpha1.Cell) error {
	tokenService, err := h.tokenServiceLister.TokenServices(cell.Namespace).Get(resources.TokenServiceName(cell))
	if errors.IsNotFound(err) {
		tokenService, err = h.meshClient.MeshV1alpha1().TokenServices(cell.Namespace).Create(resources.CreateTokenService(cell))
		if err != nil {
			h.logger.Errorf("Failed to create TokenService %v", err)
			return err
		}
		h.logger.Debugw("TokenService created", resources.TokenServiceName(cell), tokenService)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *cellHandler) handleEnvoyFilter(cell *v1alpha1.Cell) error {
	envoyFilter, err := h.envoyFilterLister.EnvoyFilters(cell.Namespace).Get(resources.EnvoyFilterName(cell))
	if errors.IsNotFound(err) {
		envoyFilter, err = h.meshClient.NetworkingV1alpha3().EnvoyFilters(cell.Namespace).Create(resources.CreateEnvoyFilter(cell))
		if err != nil {
			h.logger.Errorf("Failed to create EnvoyFilter %v", err)
			return err
		}
		h.logger.Debugw("EnvoyFilter created", resources.EnvoyFilterName(cell), envoyFilter)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *cellHandler) handleServices(cell *v1alpha1.Cell) error {
	servicesSpecs := cell.Spec.ServiceTemplates
	cell.Status.ServiceCount = 0
	for _, serviceSpec := range servicesSpecs {
		service, err := h.serviceLister.Services(cell.Namespace).Get(resources.ServiceName(cell, serviceSpec))
		if errors.IsNotFound(err) {
			service, err = h.meshClient.MeshV1alpha1().Services(cell.Namespace).Create(resources.CreateService(cell, serviceSpec))
			if err != nil {
				h.logger.Errorf("Failed to create Service: %s : %v", serviceSpec.Name, err)
				return err
			}
			h.logger.Debugw("Service created", resources.ServiceName(cell, serviceSpec), service)
		} else if err != nil {
			return err
		}
		if service.Status.AvailableReplicas > 0 || service.Spec.IsZeroScaled() {
			cell.Status.ServiceCount++
		}
	}
	return nil
}

func (h *cellHandler) updateStatus(cell *v1alpha1.Cell) (*v1alpha1.Cell, error) {
	latestCell, err := h.cellLister.Cells(cell.Namespace).Get(cell.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(latestCell.Status, cell.Status) {
		latestCell.Status = cell.Status

		return h.meshClient.MeshV1alpha1().Cells(cell.Namespace).Update(latestCell)
	}
	return cell, nil
}

func (h *cellHandler) updateCellStatus(cell *v1alpha1.Cell) {
	if cell.Status.GatewayStatus == "Ready" && int(cell.Status.ServiceCount) == len(cell.Spec.ServiceTemplates) {
		cell.Status.Status = "Ready"
		c := []v1alpha1.CellCondition{
			{
				Type:   v1alpha1.CellReady,
				Status: corev1.ConditionTrue,
			},
		}
		cell.Status.Conditions = c
	} else {
		cell.Status.Status = "NotReady"
		c := []v1alpha1.CellCondition{
			{
				Type:   v1alpha1.CellReady,
				Status: corev1.ConditionFalse,
			},
		}
		cell.Status.Conditions = c
	}
}

func (h *cellHandler) updateSecret(obj interface{}) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return
	}

	if secret.Name != mesh.SystemSecretName {
		return
	}

	s := config.Secret{}

	if keyBytes, ok := secret.Data["tls.key"]; ok {
		block, _ := pem.Decode(keyBytes)
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			h.logger.Errorf("error while parsing cellery-secret tls.key: %v", err)
			s.PrivateKey = nil
		}
		key, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			h.logger.Errorf("error while parsing cellery-secret tls.key: non rsa private key")
			s.PrivateKey = nil
		}
		s.PrivateKey = key
	} else {
		h.logger.Errorf("Missing tls.key in the cellery secret.")
	}

	if certBytes, ok := secret.Data["tls.crt"]; ok {
		block, _ := pem.Decode(certBytes)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			h.logger.Errorf("error while parsing cellery-secret tls.crt: %v", err)
			s.PrivateKey = nil
		}
		s.Certificate = cert
	} else {
		h.logger.Errorf("Missing tls.cert in the cellery secret.")
	}

	if certBundle, ok := secret.Data["cert-bundle.pem"]; ok {
		s.CertBundle = certBundle
	} else {
		h.logger.Errorf("Missing cert-bundle.pem in the cellery secret.")
	}

	h.cellerySecret = s
}
