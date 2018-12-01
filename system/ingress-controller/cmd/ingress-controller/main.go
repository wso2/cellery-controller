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

package main

import (
	"github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"github.com/wso2/product-vick/system/ingress-controller/controller"
	"github.com/wso2/product-vick/system/ingress-controller/handler"
	"os"
	"os/signal"
	"syscall"

	log "github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	ingressclientset "github.com/wso2/product-vick/system/ingress-controller/pkg/client/clientset/versioned"
	ingressinformer_v1 "github.com/wso2/product-vick/system/ingress-controller/pkg/client/informers/externalversions/ingress/v1alpha1"
)

// retrieve the Kubernetes cluster client from outside of the cluster
func getKubernetesClient() (kubernetes.Interface, ingressclientset.Interface) {
	// construct the path to resolve to `~/.kube/config`
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"

	// create the config from the path
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	// generate the client based off of the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	vickingressClient, err := ingressclientset.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}

	log.Infof("Successfully constructed k8s client")
	return client, vickingressClient
}

// main code path
func main() {
	// get the Kubernetes client for connectivity
	client, globalCpIngressClient := getKubernetesClient()

	// retrieve our custom resource informer which was generated from
	// the code generator and pass it the custom resource client, specifying
	// we should be looking through all namespaces for listing and watching
	informer := ingressinformer_v1.NewIngressInformer(
		globalCpIngressClient,
		metav1.NamespaceAll,
		0,
		cache.Indexers{},
	)

	// create a new queue so that when the informer gets a resource that is either
	// a result of listing or watching, we can add an idenfitying key to the queue
	// so that it can be handled in the handler
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// add event handlers to handle the three types of events for resources:
	//  - adding new resources
	//  - updating existing resources
	//  - deleting resources
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// check the object type
			ingress, ok := obj.(*v1alpha1.Ingress)
			if ok {
				ingObj := &controller.IngressHolder{Operation: controller.AddOp, OldObj:nil, Obj:ingress.DeepCopy()}
				queue.Add(ingObj)
				log.Infof("Add vick ingress: %+v", ingObj.Obj)
			} else {
				log.Errorf("Added object type is not Ingress")
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldIngress, oldOk := oldObj.(*v1alpha1.Ingress)
			if !oldOk {
				log.Errorf("Old object type is not Ingress")
			}
			newIngress, newOk := newObj.(*v1alpha1.Ingress)
			if !newOk {
				log.Errorf("New object type is not Ingress")
			}
			if oldOk && newOk {
				ingObj := &controller.IngressHolder{Operation: controller.UpdateOp, OldObj:oldIngress.DeepCopy(), Obj:newIngress.DeepCopy()}
				queue.Add(ingObj)
				log.Infof("Updated vick ingress, new:%+v, old:%+v", ingObj.Obj, ingObj.OldObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// check the object type
			ingress, ok := obj.(*v1alpha1.Ingress)
			if ok {
				ingObj := &controller.IngressHolder{Operation: controller.DeleteOp, OldObj:nil, Obj:ingress.DeepCopy()}
				queue.Add(ingObj)
				log.Infof("Delete vick ingress: %+v", ingObj.Obj)
			} else {
				log.Errorf("Deleted object type is not Ingress")
			}
		},
	})

	// construct the Controller object which has all of the necessary components to
	// handle logging, connections, informing (listing and watching), the queue,
	// and the handler
	controller := controller.Controller{
		client,
		queue,
		informer,
		&handler.IngressHandler{},
	}

	// use a channel to synchronize the finalization for a graceful shutdown
	stopCh := make(chan struct{})
	defer close(stopCh)

	// run the controller loop to process items
	go controller.Run(stopCh)

	// use a channel to handle OS signals to terminate and gracefully shut
	// down processing
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)
	signal.Notify(sigTerm, syscall.SIGINT)
	<-sigTerm
}
