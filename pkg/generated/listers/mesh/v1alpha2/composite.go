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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "cellery.io/cellery-controller/pkg/apis/mesh/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CompositeLister helps list Composites.
type CompositeLister interface {
	// List lists all Composites in the indexer.
	List(selector labels.Selector) (ret []*v1alpha2.Composite, err error)
	// Composites returns an object that can list and get Composites.
	Composites(namespace string) CompositeNamespaceLister
	CompositeListerExpansion
}

// compositeLister implements the CompositeLister interface.
type compositeLister struct {
	indexer cache.Indexer
}

// NewCompositeLister returns a new CompositeLister.
func NewCompositeLister(indexer cache.Indexer) CompositeLister {
	return &compositeLister{indexer: indexer}
}

// List lists all Composites in the indexer.
func (s *compositeLister) List(selector labels.Selector) (ret []*v1alpha2.Composite, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.Composite))
	})
	return ret, err
}

// Composites returns an object that can list and get Composites.
func (s *compositeLister) Composites(namespace string) CompositeNamespaceLister {
	return compositeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CompositeNamespaceLister helps list and get Composites.
type CompositeNamespaceLister interface {
	// List lists all Composites in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha2.Composite, err error)
	// Get retrieves the Composite from the indexer for a given namespace and name.
	Get(name string) (*v1alpha2.Composite, error)
	CompositeNamespaceListerExpansion
}

// compositeNamespaceLister implements the CompositeNamespaceLister
// interface.
type compositeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Composites in the indexer for a given namespace.
func (s compositeNamespaceLister) List(selector labels.Selector) (ret []*v1alpha2.Composite, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha2.Composite))
	})
	return ret, err
}

// Get retrieves the Composite from the indexer for a given namespace and name.
func (s compositeNamespaceLister) Get(name string) (*v1alpha2.Composite, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha2.Resource("composite"), name)
	}
	return obj.(*v1alpha2.Composite), nil
}
