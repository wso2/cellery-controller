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

	"k8s.io/apimachinery/pkg/runtime"
	kubefake "k8s.io/client-go/kubernetes/fake"

	clientgotesting "k8s.io/client-go/testing"

	"cellery.io/cellery-controller/pkg/clients"
	meshfake "cellery.io/cellery-controller/pkg/generated/clientset/versioned/fake"
)

type fake struct {
	clients.Interface
	fakeKubeClient *kubefake.Clientset
	fakeMeshClient *meshfake.Clientset
	createActions  []clientgotesting.CreateAction
	updatesActions []clientgotesting.UpdateAction
}

func New(objs ...runtime.Object) *fake {
	ks := kubefake.NewSimpleClientset()
	ms := meshfake.NewSimpleClientset()
	return &fake{
		Interface:      clients.NewFromClients(ks, ms),
		fakeKubeClient: ks,
		fakeMeshClient: ms,
	}
}

func (f *fake) FakeKubernetesClient() *kubefake.Clientset {
	return f.fakeKubeClient
}

func (f *fake) FakeMeshClient() *meshfake.Clientset {
	return f.fakeMeshClient
}

func (f *fake) ProcessActions() {
	var actions []clientgotesting.Action
	actions = append(actions, f.fakeMeshClient.Actions()...)
	actions = append(actions, f.fakeKubeClient.Actions()...)

	for _, action := range actions {
		switch action.(type) {
		case clientgotesting.CreateAction:
			f.createActions = append(f.createActions, action.(clientgotesting.CreateAction))
		case clientgotesting.UpdateAction:
			f.updatesActions = append(f.updatesActions, action.(clientgotesting.UpdateAction))
		default:
			panic(fmt.Sprintf("unknown action type %+v", action))
		}
	}
}

func (f *fake) ClearActions() {
	f.fakeKubeClient.ClearActions()
	f.fakeMeshClient.ClearActions()
}

func (f *fake) GetCreateActions() []clientgotesting.CreateAction {
	return f.createActions
}

func (f *fake) GetUpdateActions() []clientgotesting.UpdateAction {
	return f.updatesActions
}
