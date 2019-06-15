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

package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	v1 "k8s.io/api/apps/v1"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	appsv1informers "k8s.io/client-go/informers/apps/v1"
	autoscalingV2beta1Informer "k8s.io/client-go/informers/autoscaling/v2beta1"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	autoscalingV2beta1Lister "k8s.io/client-go/listers/autoscaling/v2beta1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh"
	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha1"
	meshclientset "github.com/cellery-io/mesh-controller/pkg/client/clientset/versioned"
	meshinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/mesh/v1alpha1"
	istioinformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/networking/v1alpha3"
	knativeservinginformers "github.com/cellery-io/mesh-controller/pkg/client/informers/externalversions/serving/v1alpha1"
	listers "github.com/cellery-io/mesh-controller/pkg/client/listers/mesh/v1alpha1"
	istionetworklisters "github.com/cellery-io/mesh-controller/pkg/client/listers/networking/v1alpha3"
	knativeservinglisters "github.com/cellery-io/mesh-controller/pkg/client/listers/serving/v1alpha1"
	"github.com/cellery-io/mesh-controller/pkg/controller"
	"github.com/cellery-io/mesh-controller/pkg/controller/autoscale"
	"github.com/cellery-io/mesh-controller/pkg/controller/service/config"
	"github.com/cellery-io/mesh-controller/pkg/controller/service/resources"
)

type serviceHandler struct {
	kubeClient                 kubernetes.Interface
	meshClient                 meshclientset.Interface
	deploymentLister           appsv1listers.DeploymentLister
	hpaLister                  autoscalingV2beta1Lister.HorizontalPodAutoscalerLister
	autoscalePolicyLister      listers.AutoscalePolicyLister
	servingConfigurationLister knativeservinglisters.ConfigurationLister
	istioVSLister              istionetworklisters.VirtualServiceLister
	k8sServiceLister           corev1listers.ServiceLister
	serviceLister              listers.ServiceLister
	configMapLister            corev1listers.ConfigMapLister
	serviceConfig              config.Service
	logger                     *zap.SugaredLogger
}

func NewController(
	kubeClient kubernetes.Interface,
	meshClient meshclientset.Interface,
	deploymentInformer appsv1informers.DeploymentInformer,
	hpaInformer autoscalingV2beta1Informer.HorizontalPodAutoscalerInformer,
	autoscalePolicyInformer meshinformers.AutoscalePolicyInformer,
	servingConfigurationInformer knativeservinginformers.ConfigurationInformer,
	istioVSInformer istioinformers.VirtualServiceInformer,
	k8sServiceInformer corev1informers.ServiceInformer,
	serviceInformer meshinformers.ServiceInformer,
	configMapInformer corev1informers.ConfigMapInformer,
	logger *zap.SugaredLogger,
) *controller.Controller {
	h := &serviceHandler{
		kubeClient:                 kubeClient,
		meshClient:                 meshClient,
		deploymentLister:           deploymentInformer.Lister(),
		hpaLister:                  hpaInformer.Lister(),
		autoscalePolicyLister:      autoscalePolicyInformer.Lister(),
		servingConfigurationLister: servingConfigurationInformer.Lister(),
		istioVSLister:              istioVSInformer.Lister(),
		k8sServiceLister:           k8sServiceInformer.Lister(),
		serviceLister:              serviceInformer.Lister(),
		configMapLister:            configMapInformer.Lister(),
		logger:                     logger.Named("service-controller"),
	}
	c := controller.New(h, h.logger, "Service")

	h.logger.Info("Setting up event handlers")
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.Enqueue,
		UpdateFunc: func(old, new interface{}) {
			h.logger.Debugw("Informer update", "old", old, "new", new)
			c.Enqueue(new)
		},
		DeleteFunc: c.Enqueue,
	})

	configMapInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: h.updateConfig,
		UpdateFunc: func(old, new interface{}) {
			h.updateConfig(new)
		},
	})

	return c
}

func (h *serviceHandler) Handle(key string) error {
	h.logger.Infof("Handle called with %s", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		h.logger.Errorf("invalid resource key: %s", key)
		return nil
	}
	serviceOriginal, err := h.serviceLister.Services(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("service '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}
	h.logger.Debugw("lister instance", key, serviceOriginal)
	service := serviceOriginal.DeepCopy()

	if err = h.handle(service); err != nil {
		return err
	}

	if _, err = h.updateStatus(service); err != nil {
		return err
	}

	return nil
}

func (h *serviceHandler) handle(service *v1alpha1.Service) error {

	if service.Spec.IsZeroScaled() {
		if err := h.handleZeroScaleDeployment(service); err != nil {
			return err
		}
		if err := h.handleZeroScaleVirtualService(service); err != nil {
			return err
		}
	} else {
		if err := h.handleDeployment(service); err != nil {
			return err
		}

		if err := h.handleAutoscalePolicy(service); err != nil {
			return err
		}

		if err := h.handleK8sService(service); err != nil {
			return err
		}
	}

	h.updateOwnerCell(service)
	return nil
}

func (h *serviceHandler) handleDeployment(service *v1alpha1.Service) error {
	deployment, err := h.deploymentLister.Deployments(service.Namespace).Get(resources.ServiceDeploymentName(service))
	if errors.IsNotFound(err) {
		deployment, err = h.kubeClient.AppsV1().Deployments(service.Namespace).Create(resources.CreateServiceDeployment(service))
		if err != nil {
			h.logger.Errorf("Failed to create Service Deployment %v", err)
			return err
		}
		h.logger.Debugw("Deployment created", resources.ServiceDeploymentName(service), deployment)
	} else if err != nil {
		return err
	}
	if deployment != nil {
		// deployment already exists, update if there is a change.
		newDeployment := resources.UpdateContainersOfServiceDeployment(deployment, service)
		// set the previous deployment's `ResourceVersion` to the newDeployment.
		// Else the issue `metadata.resourceVersion: Invalid value: 0x0: must be specified for an update` will occur.
		newDeployment.ResourceVersion = deployment.ResourceVersion
		if !isEqual(deployment, newDeployment) {
			deployment, err = h.kubeClient.AppsV1().Deployments(service.Namespace).Update(newDeployment)
			if err != nil {
				h.logger.Errorf("Failed to update Service Deployment %v", err)
				return err
			}
			h.logger.Debugw("Deployment updated", resources.ServiceDeploymentName(service), deployment)
		}
	}

	service.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	return nil
}

func isEqual(oldDeployment *v1.Deployment, newDeployment *v1.Deployment) bool {
	// only the container specs of the two deployments are considered for any differences
	return reflect.DeepEqual(oldDeployment.Spec.Template.Spec.Containers, newDeployment.Spec.Template.Spec.Containers)
}

func (h *serviceHandler) handleAutoscalePolicy(service *v1alpha1.Service) error {
	// if autoscaling is disabled in cellery system wide configs, do nothing
	if !h.serviceConfig.EnableAutoscaling {
		// if the Autoscaler has been previously created, delete it
		_, err := h.autoscalePolicyLister.AutoscalePolicies(service.Namespace).Get(resources.ServiceAutoscalePolicyName(service))
		if errors.IsNotFound(err) {
			// has not been created previously
			h.logger.Debugf("Autoscaling disabled, hence no autoscaler created for Service %s", service.Name)
			return nil
		} else if err != nil {
			return err
		}
		// created previously, delete it
		var delPropPtr = metav1.DeletePropagationBackground
		err = h.meshClient.MeshV1alpha1().AutoscalePolicies(service.Namespace).Delete(resources.ServiceAutoscalePolicyName(service),
			&metav1.DeleteOptions{
				PropagationPolicy: &delPropPtr,
			})
		if err != nil {
			return err
		}
		h.logger.Debugf("Terminated Autoscaler for Service %s", service.Name)

		return nil
	}

	autoscalePolicy, err := h.autoscalePolicyLister.AutoscalePolicies(service.Namespace).Get(resources.ServiceAutoscalePolicyName(service))
	if errors.IsNotFound(err) {
		if service.Spec.Autoscaling != nil {
			autoscalePolicy = resources.CreateAutoscalePolicy(service)
		} else {
			autoscalePolicy = resources.CreateDefaultAutoscalePolicy(service)
		}
		lastAppliedConfig, err := json.Marshal(autoscale.BuildAutoscalePolicyLastAppliedConfig(autoscalePolicy))
		if err != nil {
			h.logger.Errorf("Failed to create Autoscale policy %v", err)
			return err
		}
		autoscale.Annotate(autoscalePolicy, corev1.LastAppliedConfigAnnotation, string(lastAppliedConfig))
		autoscalePolicy, err = h.meshClient.MeshV1alpha1().AutoscalePolicies(service.Namespace).Create(autoscalePolicy)
		if err != nil {
			h.logger.Errorf("Failed to create Autoscale policy %v", err)
			return err
		}
		h.logger.Infow("Autoscale policy created", resources.ServiceAutoscalePolicyName(service), autoscalePolicy)
		autoscalePolicy.Status = "Ready"
		h.updateAutoscalePolicyStatus(autoscalePolicy)

	} else if err != nil {
		return err
	}
	return nil
}

func (h *serviceHandler) handleZeroScaleDeployment(service *v1alpha1.Service) error {
	zeroScaleDeployment, err := h.servingConfigurationLister.Configurations(service.Namespace).Get(resources.ServiceServingConfigurationName(service))
	if errors.IsNotFound(err) {
		zeroScaleDeployment, err = h.meshClient.ServingV1alpha1().Configurations(service.Namespace).Create(resources.CreateZeroScaleService(service))
		if err != nil {
			h.logger.Errorf("Failed to create Zero Scale Service Deployment %v", err)
			return err
		}
		h.logger.Debugw("Zero Scale Deployment created", resources.ServiceServingConfigurationName(service), zeroScaleDeployment)
	} else if err != nil {
		return err
	}

	service.Status.HostName = zeroScaleDeployment.Status.LatestReadyRevisionName
	return nil
}

func (h *serviceHandler) handleZeroScaleVirtualService(service *v1alpha1.Service) error {
	zeroScaleVirtualService, err := h.istioVSLister.VirtualServices(service.Namespace).Get(resources.ServiceServingVirtualServiceName(service))
	if errors.IsNotFound(err) {
		zeroScaleVirtualService, err = h.meshClient.NetworkingV1alpha3().VirtualServices(service.Namespace).Create(resources.CreateZeroScaleVirtualService(service))
		if err != nil {
			h.logger.Errorf("Failed to create Zero Scale VirtualService %v", err)
			return err
		}
		h.logger.Debugw("Zero Scale VirtualService created", resources.ServiceServingVirtualServiceName(service), zeroScaleVirtualService)
	} else if err != nil {
		return err
	}
	return nil
}

func (h *serviceHandler) updateAutoscalePolicyStatus(policy *v1alpha1.AutoscalePolicy) error {
	_, err := h.meshClient.MeshV1alpha1().AutoscalePolicies(policy.Namespace).Update(policy)
	if err != nil {
		return err
	}
	return nil
}

func (h *serviceHandler) handleK8sService(service *v1alpha1.Service) error {
	k8sService, err := h.k8sServiceLister.Services(service.Namespace).Get(resources.ServiceK8sServiceName(service))
	if errors.IsNotFound(err) {
		k8sService, err = h.kubeClient.CoreV1().Services(service.Namespace).Create(resources.CreateServiceK8sService(service))
		if err != nil {
			h.logger.Errorf("Failed to create Service service %v", err)
			return err
		}
		h.logger.Debugw("Service created", resources.ServiceK8sServiceName(service), k8sService)
	} else if err != nil {
		return err
	}

	service.Status.HostName = k8sService.Name

	return nil
}

func (h *serviceHandler) updateStatus(service *v1alpha1.Service) (*v1alpha1.Service, error) {
	latestService, err := h.serviceLister.Services(service.Namespace).Get(service.Name)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual(latestService.Status, service.Status) {
		latestService.Status = service.Status
		// https://github.com/kubernetes/kubernetes/issues/54672
		// return h.meshClient.MeshV1alpha1().Services(service.Namespace).UpdateStatus(latestService)
		return h.meshClient.MeshV1alpha1().Services(service.Namespace).Update(latestService)
	}
	return nil, nil
}

func (h *serviceHandler) updateOwnerCell(service *v1alpha1.Service) {
	if service.Status.OwnerCell == "" {
		ownerRef := metav1.GetControllerOf(service)
		if ownerRef != nil {
			service.Status.OwnerCell = ownerRef.Name
		} else {
			service.Status.OwnerCell = "<none>"
		}
	}
}

func (h *serviceHandler) updateConfig(obj interface{}) {

	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return
	}

	if configMap.Name != mesh.SystemConfigMapName {
		return
	}

	conf := config.Service{}

	if enableAutoscaling, ok := configMap.Data["enable-autoscaling"]; ok {
		conf.EnableAutoscaling = isAutoscalingEnabled(enableAutoscaling)
	} else {
		conf.EnableAutoscaling = false
	}

	h.serviceConfig = conf
}

func isAutoscalingEnabled(str string) bool {
	if enabled, err := strconv.ParseBool(str); err == nil {
		return enabled
	}
	return false
}
