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

package fake

import (
	"fmt"
	"reflect"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"

	v1alpha2 "github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	clients "github.com/cellery-io/mesh-controller/pkg/clients"
	informers "github.com/cellery-io/mesh-controller/pkg/informers"
)

type fake struct {
	informers.Interface
	cache map[reflect.Type]cache.Indexer
}

func New(clients clients.Interface, resync time.Duration, objs ...runtime.Object) *fake {
	f := &fake{
		Interface: informers.New(clients, resync),
		cache:     make(map[reflect.Type]cache.Indexer),
	}

	f.addIndexer(&v1alpha2.Component{}, f.Components().Informer().GetIndexer())
	f.addIndexer(&corev1.Service{}, f.Services().Informer().GetIndexer())
	f.addIndexer(&appsv1.Deployment{}, f.Deployments().Informer().GetIndexer())

	for _, obj := range objs {
		t := reflect.TypeOf(obj).Elem()
		indexer, ok := f.cache[t]
		if !ok {
			panic(fmt.Sprintf("Unrecognized type %T", obj))
		}
		indexer.Add(obj)
	}
	return f
}

func (f *fake) addIndexer(obj runtime.Object, indexer cache.Indexer) {
	f.cache[reflect.TypeOf(obj).Elem()] = indexer
}
