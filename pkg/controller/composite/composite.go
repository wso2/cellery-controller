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

package composite

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	meshclientset "github.com/cellery-io/mesh-controller/pkg/client/clientset/versioned"
	meshinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/mesh/v1alpha1"
	istioinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/networking/v1alpha3"
	listers "github.com/cellery-io/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	istiov1alpha1listers "github.com/cellery-io/mesh-controller/pkg/client/listers/networking/v1alpha3"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	controller_commons "github.com/cellery-io/mesh-controller/pkg/controller/commons"
	"github.com/cellery-io/mesh-controller/pkg/controller/composite/config"
	"github.com/cellery-io/mesh-controller/pkg/controller/composite/resources"
)

type compositeHandler struct {
	kubeClient         kubernetes.Interface
	meshClient         meshclientset.Interface
	compositeLister    listers.CompositeLister
	serviceLister      listers.ServiceLister
	secretLister       corev1listers.SecretLister
	tokenServiceLister listers.TokenServiceLister
	cellerySecret      config.Secret
	virtualSvcLister   istiov1alpha1listers.VirtualServiceLister
	cellLister         listers.CellLister
	logger             *zap.SugaredLogger
}

func NewController(
	kubeClient kubernetes.Interface,
	meshClient meshclientset.Interface,
	compositeInformer meshinformers.CompositeInformer,
	serviceInformer meshinformers.ServiceInformer,
	tokenServiceInformer meshinformers.TokenServiceInformer,
	secretInformer corev1informers.SecretInformer,
	systemSecretInformer corev1informers.SecretInformer,
	virtualSvcInformer istioinformers.VirtualServiceInformer,
	cellInformer meshinformers.CellInformer,
	logger *zap.SugaredLogger,
) *controller.Controller {
	h := &compositeHandler{
		kubeClient:         kubeClient,
		meshClient:         meshClient,
		compositeLister:    compositeInformer.Lister(),
		serviceLister:      serviceInformer.Lister(),
		tokenServiceLister: tokenServiceInformer.Lister(),
		secretLister:       secretInformer.Lister(),
		virtualSvcLister:   virtualSvcInformer.Lister(),
		cellLister:         cellInformer.Lister(),
		logger:             logger.Named("composite-controller"),
	}
	c := controller.New(h, h.logger, "Composite")

	h.logger.Info("Setting up event handlers")
	compositeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
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

func (h *compositeHandler) Handle(key string) error {
	h.logger.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		h.logger.Errorf("invalid resource key: %s", key)
		return nil
	}
	compositeOriginal, err := h.compositeLister.Composites(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("cell '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	h.logger.Debugw("lister instance", key, compositeOriginal)
	composite := compositeOriginal.DeepCopy()

	if err = h.handle(composite); err != nil {
		return err
	}

	if _, err = h.updateStatus(composite); err != nil {
		return err
	}
	return nil
}

func (h *compositeHandler) handle(composite *v1alpha1.Composite) error {

	if err := h.handleServices(composite); err != nil {
		return err
	}

	if err := h.handleSecret(composite); err != nil {
		return err
	}

	if err := h.handleTokenService(composite); err != nil {
		return err
	}

	if err := h.handleVirtualService(composite); err != nil {
		return err
	}

	h.updateCellStatus(composite)
	return nil
}

func (h *compositeHandler) handleVirtualService(composite *v1alpha1.Composite) error {
	cellVs, err := h.virtualSvcLister.VirtualServices(composite.Namespace).Get(controller_commons.VirtualServiceName(composite.Name))
	if errors.IsNotFound(err) {
		cellVs, err = resources.CreateCellVirtualService(composite, h.compositeLister, h.cellLister)
		if err != nil {
			h.logger.Errorf("Failed to create Cell VS object %v for instance %s", err, composite.Name)
			return err
		}
		if cellVs == nil {
			h.logger.Debugf("No VirtualService created for composite instance %s", composite.Name)
			return nil
		}
		lastAppliedConfig, err := json.Marshal(controller_commons.BuildVirtualServiceLastAppliedConfig(cellVs))
		if err != nil {
			h.logger.Errorf("Failed to create Cell VS %v for instance %s", err, composite.Name)
			return err
		}
		controller_commons.Annotate(cellVs, corev1.LastAppliedConfigAnnotation, string(lastAppliedConfig))
		cellVs, err = h.meshClient.NetworkingV1alpha3().VirtualServices(composite.Namespace).Create(cellVs)
		if err != nil {
			h.logger.Errorf("Failed to create Cell VirtualService %v for instance %s", err, composite.Name)
			return err
		}
		h.logger.Debugw("Cell VirtualService created", controller_commons.VirtualServiceName(composite.Name), cellVs)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *compositeHandler) handleServices(composite *v1alpha1.Composite) error {
	servicesSpecs := composite.Spec.ServiceTemplates
	composite.Status.ServiceCount = 0
	for _, serviceSpec := range servicesSpecs {
		service, err := h.serviceLister.Services(composite.Namespace).Get(resources.ServiceName(composite, serviceSpec))
		if errors.IsNotFound(err) {
			service, err = h.meshClient.MeshV1alpha1().Services(composite.Namespace).Create(resources.CreateService(composite, serviceSpec))
			if err != nil {
				h.logger.Errorf("Failed to create Service: %s : %v", serviceSpec.Name, err)
				return err
			}
			h.logger.Debugw("Service created", resources.ServiceName(composite, serviceSpec), service)
		} else if err != nil {
			return err
		}
		if service != nil {
			// service exists. if the new obj is not equal to old one, perform an update.
			newService := resources.CreateService(composite, serviceSpec)
			// set the previous service's `ResourceVersion` to the newService
			// Else the issue `metadata.resourceVersion: Invalid value: 0x0: must be specified for an update` will occur.
			newService.ResourceVersion = service.ResourceVersion
			if !isEqual(service, newService) {
				service, err = h.meshClient.MeshV1alpha1().Services(composite.Namespace).Update(newService)
				if err != nil {
					h.logger.Errorf("Failed to update Service: %s : %v", service.Name, err)
					return err
				}
				h.logger.Debugw("Service updated", resources.ServiceName(composite, serviceSpec), service)
			}
		}
		if service.Status.AvailableReplicas > 0 || service.Spec.IsZeroScaled() || service.Spec.Type == v1alpha1.ServiceTypeJob {
			composite.Status.ServiceCount++
		}
	}
	return nil
}

func (h *compositeHandler) handleTokenService(composite *v1alpha1.Composite) error {
	tokenService, err := h.tokenServiceLister.TokenServices(composite.Namespace).Get(resources.TokenServiceName(composite))
	if errors.IsNotFound(err) {
		tokenService, err = h.meshClient.MeshV1alpha1().TokenServices(composite.Namespace).Create(resources.CreateTokenService(composite))
		if err != nil {
			h.logger.Errorf("Failed to create Composite TokenService %v", err)
			return err
		}
		h.logger.Debugw("Composite TokenService created", resources.TokenServiceName(composite), tokenService)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *compositeHandler) handleSecret(composite *v1alpha1.Composite) error {
	secret, err := h.secretLister.Secrets(composite.Namespace).Get(resources.SecretName(composite))
	if errors.IsNotFound(err) {
		desiredSecret, err := resources.CreateKeyPairSecret(composite, h.cellerySecret)
		if err != nil {
			h.logger.Errorf("Cannot build the Composite Secret %v", err)
			return err
		}
		secret, err = h.kubeClient.CoreV1().Secrets(composite.Namespace).Create(desiredSecret)
		if err != nil {
			h.logger.Errorf("Failed to create Composite Secret %v", err)
			return err
		}
		h.logger.Debugw("Secret created", resources.SecretName(composite), secret)
	} else if err != nil {
		return err
	}
	return nil
}

func isEqual(oldService *v1alpha1.Service, newService *v1alpha1.Service) bool {
	// we only consider equality of the spec
	return reflect.DeepEqual(oldService.Spec, newService.Spec)
}

func (h *compositeHandler) updateStatus(composite *v1alpha1.Composite) (*v1alpha1.Composite, error) {
	latestComposite, err := h.compositeLister.Composites(composite.Namespace).Get(composite.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(latestComposite.Status, composite.Status) {
		latestComposite.Status = composite.Status

		return h.meshClient.MeshV1alpha1().Composites(composite.Namespace).Update(latestComposite)
	}
	return composite, nil
}

func (h *compositeHandler) updateCellStatus(composite *v1alpha1.Composite) {
	if int(composite.Status.ServiceCount) == len(composite.Spec.ServiceTemplates) {
		composite.Status.Status = "Ready"
		c := []v1alpha1.CompositeCondition{
			{
				Type:   v1alpha1.CompositeReady,
				Status: corev1.ConditionTrue,
			},
		}
		composite.Status.Conditions = c
	} else {
		composite.Status.Status = "NotReady"
		c := []v1alpha1.CompositeCondition{
			{
				Type:   v1alpha1.CompositeReady,
				Status: corev1.ConditionFalse,
			},
		}
		composite.Status.Conditions = c
	}
}

func (h *compositeHandler) updateSecret(obj interface{}) {
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
