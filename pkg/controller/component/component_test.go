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

package component

import (
	"fmt"

	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime"
	clienttesting "k8s.io/client-go/testing"

	fakeclients "github.com/cellery-io/mesh-controller/pkg/clients/fake"
	fakeinformers "github.com/cellery-io/mesh-controller/pkg/informers/fake"
	"github.com/cellery-io/mesh-controller/pkg/logging"
	. "github.com/cellery-io/mesh-controller/pkg/testing/apis/apps/v1"
	. "github.com/cellery-io/mesh-controller/pkg/testing/apis/core/v1"
	. "github.com/cellery-io/mesh-controller/pkg/testing/apis/mesh/v1alpha2"
)

var noResyncPeriodFunc = func() time.Duration { return 0 }

type Test struct {
	Name        string
	Key         string
	Objects     []runtime.Object
	WantCreates []runtime.Object
}

func TestReconcile(t *testing.T) {

	tests := []Test{
		{
			Name: "invalid key",
			Key:  "foo/bar/baz",
		},
		{
			Name: "non existing key",
			Key:  "foo/bar",
		},
		{
			Name: "create service and deployment",
			Key:  "foo/component",
			Objects: []runtime.Object{
				Component("component", "foo",
					WithComponentPodSpec(
						PodSpec(
							WithPodSpecContainer(
								Container(
									WithContainerImage("busybox:v1.2.3"),
									WithContainerEnvFromValue("env-key1", "env-value1"),
									WithContainerVolumeMounts("pvc1", true, "/data", ""),
									WithContainerVolumeMounts("config1", true, "/etc/conf", ""),
									WithContainerVolumeMounts("secret1", true, "/etc/certs", ""),
								),
							),
						),
					),
					WithComponentPortMaping("port1", "HTTP", 80, "", 8080),
					WithComponentPortMaping("foo-rpc", "GRPC", 9090, "", 9090),
					WithComponentPortMaping("", "", 15000, "", 8001),
					WithComponentVolumeClaim(
						false,
						PersistentVolumeClaim("pvc1", "foo-namespace"),
					),
					WithComponentConfiguration(
						ConfigMap("config1", "foo-namespace",
							WithConfigMapData("key1", "value1"),
						),
					),
					WithComponentSecret(
						Secret("secret1", "foo-namespace",
							WithSecretData("key1", []byte("password")),
						),
					),
				),
			},
			WantCreates: []runtime.Object{
				Service("component-service", "foo",
					WithServiceLabel("app", "component"),
					WithServiceOwnerReference(OwnerReferenceComponentFunc("component")),
					WithServicePort("http-port1", "TCP", 80, 8080),
					WithServicePort("grpc-foo-rpc", "TCP", 9090, 9090),
					WithServicePort("tcp-15000-8001", "TCP", 15000, 8001),
				),
				Deployment("component-deployment", "foo",
					WithDeploymentOwnerReference(OwnerReferenceComponentFunc("component")),
					WithDeploymentReplicas(1),
					WithDeploymentPodSpec(
						PodSpec(
							WithPodSpecContainer(
								Container(
									WithContainerImage("busybox:v1.2.3"),
									WithContainerEnvFromValue("env-key1", "env-value1"),
									WithContainerPort(8080),
									WithContainerPort(9090),
									WithContainerPort(8001),
									WithContainerVolumeMounts("pvc1", true, "/data", ""),
									WithContainerVolumeMounts("config1", true, "/etc/conf", ""),
									WithContainerVolumeMounts("secret1", true, "/etc/certs", ""),
								),
							),
							WithPodSpecVolume(
								Volume("pvc1", WithVolumeSourcePersistentVolumeClaim("pvc1")),
							),
							WithPodSpecVolume(
								Volume("config1", WithVolumeSourceConfigMap("config1")),
							),
							WithPodSpecVolume(
								Volume("secret1", WithVolumeSourceSecret("secret1")),
							),
						),
					),
				),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			clients := fakeclients.New()
			log, _ := logging.NewLogger()
			informers := fakeinformers.New(clients, noResyncPeriodFunc(), test.Objects...)

			clients.FakeKubernetesClient().PrependReactor("create", "*", func(action clienttesting.Action) (handled bool, ret runtime.Object, err error) {
				return false, nil, fmt.Errorf("dasda")
			})

			r := reconciler{
				kubeClient:       clients.Kubernetes(),
				meshClient:       clients.Mesh(),
				componentLister:  informers.Components().Lister(),
				serviceLister:    informers.Services().Lister(),
				deploymentLister: informers.Deployments().Lister(),
				logger:           log,
			}

			err := r.Reconcile(test.Key)
			if err != nil {
				t.Error(err)
			}

			clients.ProcessActions()

			if got, want := len(clients.GetCreateActions()), len(test.WantCreates); got > want {
				for _, extra := range clients.GetCreateActions()[want:] {
					t.Errorf("Additional object create: %#v", extra.GetObject())
				}
			}

			for i, want := range test.WantCreates {
				if i >= len(clients.GetCreateActions()) {
					t.Errorf("Missing object create: %#v", want)
					continue
				}
				if diff := cmp.Diff(want, clients.GetCreateActions()[i].GetObject()); diff != "" {
					t.Errorf("Unexpected create (-want, +got)\n%s", diff)
				}
			}
		})
	}
}

func TestNewController(t *testing.T) {
	clients := fakeclients.New()
	informers := fakeinformers.New(clients, noResyncPeriodFunc())

	log, _ := logging.NewLogger()
	c := NewController(clients, informers, log)

	if c == nil {
		t.Fatal("Expected NewController to return non-nil value")
	}
}
