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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"

	v1alpha3 "github.com/celleryio/mesh-controller/pkg/apis/istio/networking/v1alpha3"
)

// FakeEnvoyFilters implements EnvoyFilterInterface
type FakeEnvoyFilters struct {
	Fake *FakeNetworkingV1alpha3
	ns   string
}

var envoyfiltersResource = schema.GroupVersionResource{Group: "networking", Version: "v1alpha3", Resource: "envoyfilters"}

var envoyfiltersKind = schema.GroupVersionKind{Group: "networking", Version: "v1alpha3", Kind: "EnvoyFilter"}

// Get takes name of the envoyFilter, and returns the corresponding envoyFilter object, and an error if there is any.
func (c *FakeEnvoyFilters) Get(name string, options v1.GetOptions) (result *v1alpha3.EnvoyFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(envoyfiltersResource, c.ns, name), &v1alpha3.EnvoyFilter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha3.EnvoyFilter), err
}

// List takes label and field selectors, and returns the list of EnvoyFilters that match those selectors.
func (c *FakeEnvoyFilters) List(opts v1.ListOptions) (result *v1alpha3.EnvoyFilterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(envoyfiltersResource, envoyfiltersKind, c.ns, opts), &v1alpha3.EnvoyFilterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha3.EnvoyFilterList{ListMeta: obj.(*v1alpha3.EnvoyFilterList).ListMeta}
	for _, item := range obj.(*v1alpha3.EnvoyFilterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested envoyFilters.
func (c *FakeEnvoyFilters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(envoyfiltersResource, c.ns, opts))

}

// Create takes the representation of a envoyFilter and creates it.  Returns the server's representation of the envoyFilter, and an error, if there is any.
func (c *FakeEnvoyFilters) Create(envoyFilter *v1alpha3.EnvoyFilter) (result *v1alpha3.EnvoyFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(envoyfiltersResource, c.ns, envoyFilter), &v1alpha3.EnvoyFilter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha3.EnvoyFilter), err
}

// Update takes the representation of a envoyFilter and updates it. Returns the server's representation of the envoyFilter, and an error, if there is any.
func (c *FakeEnvoyFilters) Update(envoyFilter *v1alpha3.EnvoyFilter) (result *v1alpha3.EnvoyFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(envoyfiltersResource, c.ns, envoyFilter), &v1alpha3.EnvoyFilter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha3.EnvoyFilter), err
}

// Delete takes name of the envoyFilter and deletes it. Returns an error if one occurs.
func (c *FakeEnvoyFilters) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(envoyfiltersResource, c.ns, name), &v1alpha3.EnvoyFilter{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEnvoyFilters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(envoyfiltersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha3.EnvoyFilterList{})
	return err
}

// Patch applies the patch and returns the patched envoyFilter.
func (c *FakeEnvoyFilters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.EnvoyFilter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(envoyfiltersResource, c.ns, name, data, subresources...), &v1alpha3.EnvoyFilter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha3.EnvoyFilter), err
}
