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

package controller

import (
	"fmt"
	"github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"time"

	log "github.com/golang/glog"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"github.com/wso2/product-vick/system/ingress-controller/handler"
)

// Controller struct defines how a controller should encapsulate
// logging, client connectivity, informing (list and watching)
// queueing, and handling of resource changes
type Controller struct {
	Clientset kubernetes.Interface
	Queue     workqueue.RateLimitingInterface
	Informer  cache.SharedIndexInformer
	Handler   handler.Handler
}

type IngressHolder struct {
	Operation string      // add, update or delete
	OldObj    *v1alpha1.Ingress // previous object for update operations only
	Obj       *v1alpha1.Ingress // current object for add and delete
}

const (
	AddOp = "add"
	UpdateOp = "update"
	DeleteOp = "delete"
)

// Run is the main path of execution for the controller loop
func (c *Controller) Run(stopCh <-chan struct{}) {
	// handle a panic with logging and exiting
	defer utilruntime.HandleCrash()
	// ignore new items in the queue but when all goroutines
	// have completed existing items then shutdown
	defer c.Queue.ShutDown()

	log.Infof("%s","Controller.Run: initiating")

	// run the informer to start listing and watching resources
	go c.Informer.Run(stopCh)

	// do the initial synchronization (one time) to populate resources
	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Error syncing cache"))
		return
	}
	log.Infof("Controller.Run: cache sync complete")
	// call handler init
	err := c.Handler.Init()
	if err != nil {
		log.Fatal(err)
	}

	// run the runWorker method every second with a stop channel
	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced allows us to satisfy the Controller interface
// by wiring up the informer's HasSynced method to it
func (c *Controller) HasSynced() bool {
	return c.Informer.HasSynced()
}

// runWorker executes the loop to process new items added to the queue
func (c *Controller) runWorker() {
	log.Infof("Controller.runWorker: starting")

	// invoke processNextItem to fetch and consume the next change
	// to a watched or listed resource
	for c.processNextItem() {
		log.Infof("Controller.runWorker: processing next item")
	}

	log.Infof("Controller.runWorker: completed")
}

// processNextItem retrieves each queued item and takes the
// necessary handler action based off of if the item was
// created or deleted
func (c *Controller) processNextItem() bool {
	log.Infof("Controller.processNextItem: start")

	// fetch the next item (blocking) from the queue to process or
	// if a shutdown is requested then return out of this to stop
	// processing
	key, quit := c.Queue.Get()

	// stop the worker loop from running as this indicates we
	// have sent a shutdown message that the queue has indicated
	// from the Get method
	if quit {
		return false
	}

	defer c.Queue.Done(key)

	// assert the string out of the key (format `namespace/name`)
	ingHolder, ok := key.(*IngressHolder)
	if !ok {
		log.Errorf("Item from the Queue is not of the type ingressHolder")
		return false
	}

	// check the related operation and call the relevant handlers
	if ingHolder.Operation == AddOp {
		// call handler's created method
		log.Infof("Controller.processNextItem: object created detected: %+v", ingHolder.Obj)
		c.Handler.ObjectCreated(ingHolder.Obj)
	} else if (ingHolder.Operation == UpdateOp) {
		// call handler's updated method
		log.Infof("Controller.processNextItem: object updated detected, new: %+v, old: %+v", ingHolder.Obj,
			ingHolder.OldObj)
		c.Handler.ObjectUpdated(ingHolder.OldObj, ingHolder.Obj)
	} else if (ingHolder.Operation == DeleteOp) {
		log.Infof("Controller.processNextItem: object deleted detected: %+v", ingHolder.Obj)
		c.Handler.ObjectDeleted(ingHolder.Obj)
	}

	// keep the worker loop running by returning true
	return true
}
