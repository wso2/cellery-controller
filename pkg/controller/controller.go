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
	"strings"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	meshscheme "cellery.io/cellery-controller/pkg/generated/clientset/versioned/scheme"
)

type Reconciler interface {
	Reconcile(key string) error
}

type Controller struct {
	reconciler Reconciler
	name       string
	workqueue  workqueue.RateLimitingInterface
	logger     *zap.SugaredLogger
}

func New(r Reconciler, logger *zap.SugaredLogger, workQueueName string) *Controller {
	return &Controller{
		reconciler: r,
		name:       workQueueName,
		workqueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), workQueueName),
		logger:     logger,
	}
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	c.logger.Infof("Starting %s controller", c.name)

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	// wait until we're told to stop
	<-stopCh
	c.logger.Infof("Shutting down the %s controller", c.name)
}

func (c *Controller) Enqueue(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		c.logger.Error(err)
		return
	}
	c.EnqueueKey(key)
}

func (c *Controller) EnqueueControllerOf(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			c.logger.Errorf("error decoding object, invalid type")
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			c.logger.Errorf("error decoding object tombstone, invalid type")
			return
		}
		c.logger.Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	if owner := metav1.GetControllerOf(object); owner != nil {
		c.EnqueueKey(object.GetNamespace() + "/" + owner.Name)
	}
}

func (c *Controller) EnqueueKey(key string) {
	c.workqueue.AddRateLimited(key)
	c.logger.Debugf("Adding key %q to queue (depth: %d)", key, c.workqueue.Len())
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}
	c.logger.Debugf("Processing %q from queue (depth: %d)", obj, c.workqueue.Len())

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			c.logger.Errorf("expected string in workqueue but got %#v", obj)
			return nil
		}
		t := time.Now()
		// Run the reconciler, passing it the namespace/name string of the resource.
		if err := c.reconciler.Reconcile(key); err != nil {
			c.logger.Infow("Reconcile failed", "key", key, "time", time.Since(t))
			return fmt.Errorf("error reconciling '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		c.logger.Infow("Reconcile succeeded", "key", key, "time", time.Since(t))
		return nil
	}(obj)

	if err != nil {
		c.logger.Error(err)
		return true
	}

	return true
}

type ReconcileErrors struct {
	errors []error
}

func (re *ReconcileErrors) Add(err error) {
	if err != nil {
		re.errors = append(re.errors, err)
	}
}

func (re *ReconcileErrors) Empty() bool {
	return len(re.errors) == 0
}

func (re *ReconcileErrors) Error() string {
	errLen := len(re.errors)
	if errLen == 0 {
		return ""
	} else if errLen == 1 {
		return re.errors[0].Error()
	}
	var errStrings []string
	for i, _ := range re.errors {
		errStrings = append(errStrings, fmt.Sprintf("\t- %s", re.errors[i].Error()))
	}
	return fmt.Sprintf("multiple errors:\n%s", strings.Join(errStrings, "\n"))
}

func init() {
	// Add custom resource types to the default Kubernetes Scheme
	// so Events can be logged.
	utilruntime.Must(meshscheme.AddToScheme(scheme.Scheme))
}
