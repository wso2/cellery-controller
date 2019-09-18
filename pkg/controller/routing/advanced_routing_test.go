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

package commons

import (
	"fmt"
	"testing"

	"github.com/cellery-io/mesh-controller/pkg/apis/mesh/v1alpha2"
	"github.com/cellery-io/mesh-controller/pkg/meta"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cellery-io/mesh-controller/pkg/apis/istio/networking/v1alpha3"
)

func TestBuildHostNameForCellDependency(t *testing.T) {
	dependencyInst := "mydep"
	expected := dependencyInst + "--gateway-service"
	actual := BuildHostNameForCellDependency(dependencyInst)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("BuildHostNameForCellDependency (-expected, +actual)\n%v", diff)
	}
}

func TestBuildHostNamesForCompositeDependency(t *testing.T) {
	dependencyInst := "mydep"
	components := []v1alpha2.Component{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "mycomponent1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "mycomponent2",
			},
		},
	}
	expected := []string{
		"mydep--mycomponent1-service",
		"mydep--mycomponent2-service",
	}
	actual := BuildHostNamesForCompositeDependency(dependencyInst, components)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("BuildHostNamesForCompositeDependency (-expected, +actual)\n%v", diff)
	}
}

func TestBuildHttpRoutesForCellDependency(t *testing.T) {
	dependencyInst := "mydep"
	instName := "myinst"
	expected := []*v1alpha3.HTTPRoute{
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
					},
				},
			},
		},
	}
	actual := BuildHttpRoutesForCellDependency(instName, dependencyInst, false)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("BuildHttpRoutesForCellDependency (-expected, +actual)\n%v", diff)
	}
}

func TestBuildHttpRoutesForCellDependencyWithInstanceIdBasedRules(t *testing.T) {
	dependencyInst := "mydep"
	instName := "myinst"
	expected := []*v1alpha3.HTTPRoute{
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
					Headers: map[string]*v1alpha3.StringMatch{
						InstanceId: {
							Exact: "1",
						},
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
					},
				},
			},
		},
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
					Headers: map[string]*v1alpha3.StringMatch{
						InstanceId: {
							Exact: "2",
						},
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
					},
				},
			},
		},
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--gateway-service)(\\S*)$", dependencyInst),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: GatewayK8sServiceName(GatewayNameFromInstanceName(dependencyInst)),
					},
				},
			},
		},
	}
	actual := BuildHttpRoutesForCellDependency(instName, dependencyInst, true)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("BuildHttpRoutesForCellDependency (-expected, +actual)\n%v", diff)
	}
}

func TestBuildHttpRoutesForCompositeDependency(t *testing.T) {
	dependencyInst := "mydep"
	instName := "myinst"
	component := v1alpha2.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mycomponent1",
		},
	}
	expected := []*v1alpha3.HTTPRoute{
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--%s)(\\S*)$", dependencyInst, component.Name),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: CompositeK8sServiceNameFromInstance(dependencyInst, component),
					},
				},
			},
		},
	}
	actual := BuildHttpRoutesForCompositeDependency(instName, dependencyInst, []v1alpha2.Component{component}, false)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("BuildHttpRoutesForCompositeDependency (-expected, +actual)\n%v", diff)
	}
}

func TestBuildHttpRoutesForCompositeDependencyWithInstanceBasedRules(t *testing.T) {
	dependencyInst := "mydep"
	instName := "myinst"
	svcTemplate := v1alpha2.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mycomponent1",
		},
	}
	expected := []*v1alpha3.HTTPRoute{
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--%s)(\\S*)$", dependencyInst, svcTemplate.Name),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
					Headers: map[string]*v1alpha3.StringMatch{
						InstanceId: {
							Exact: "1",
						},
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: CompositeK8sServiceNameFromInstance(dependencyInst, svcTemplate),
					},
				},
			},
		},
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--%s)(\\S*)$", dependencyInst, svcTemplate.Name),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
					Headers: map[string]*v1alpha3.StringMatch{
						InstanceId: {
							Exact: "2",
						},
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: CompositeK8sServiceNameFromInstance(dependencyInst, svcTemplate),
					},
				},
			},
		},
		{
			Match: []*v1alpha3.HTTPMatchRequest{
				{
					Authority: &v1alpha3.StringMatch{
						Regex: fmt.Sprintf("^(%s)(--%s)(\\S*)$", dependencyInst, svcTemplate.Name),
					},
					SourceLabels: map[string]string{
						meta.CellLabelKeySource:      instName,
						meta.ComponentLabelKeySource: "true",
					},
				},
			},
			Route: []*v1alpha3.DestinationWeight{
				{
					Destination: &v1alpha3.Destination{
						Host: CompositeK8sServiceNameFromInstance(dependencyInst, svcTemplate),
					},
				},
			},
		},
	}
	actual := BuildHttpRoutesForCompositeDependency(instName, dependencyInst, []v1alpha2.Component{svcTemplate}, true)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("BuildHttpRoutesForCompositeDependency (-expected, +actual)\n%v", diff)
	}
}

func TestExtractDependencies(t *testing.T) {
	annotations := map[string]string{
		meta.CellDependenciesAnnotationKey: "[{\"org\":\"izza\",\"name\":\"emp-comp\",\"version\":\"0.0.4\",\"instance\":\"emp-comp-0-0-4-a1471a5b\",\"kind\":\"Composite\"},{\"org\":\"izza\",\"name\":\"stock-comp\",\"version\":\"0.0.4\",\"instance\":\"stock-comp-0-0-4-7af583f3\",\"kind\":\"Cell\"}]",
	}
	expected := []map[string]string{
		{
			"org":      "izza",
			"name":     "emp-comp",
			"version":  "0.0.4",
			"instance": "emp-comp-0-0-4-a1471a5b",
			"kind":     "Composite",
		},
		{
			"org":      "izza",
			"name":     "stock-comp",
			"version":  "0.0.4",
			"instance": "stock-comp-0-0-4-7af583f3",
			"kind":     "Cell",
		},
	}
	actual, err := ExtractDependencies(annotations)
	if err != nil {
		t.Errorf("Error while executing ExtractDependencies \n %v", err)
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("ExtractDependencies (-expected, +actual)\n%v", diff)
	}
}

func TestBuildVirtualServiceLastAppliedConfig(t *testing.T) {
	vs := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "inst1--vs",
			Namespace: "default",
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Hosts:    []string{"*"},
			Gateways: []string{"mesh"},
		},
	}
	expected := &v1alpha3.VirtualService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: v1alpha3.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "inst1--vs",
			Namespace: "default",
		},
		Spec: vs.Spec,
	}
	actual := BuildVirtualServiceLastAppliedConfig(vs)
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("BuildVirtualServiceLastAppliedConfig (-expected, +actual)\n%v", diff)
	}
}

func TestAnnotate(t *testing.T) {
	vs := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "inst1--vs",
			Namespace: "default",
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Hosts:    []string{"*"},
			Gateways: []string{"mesh"},
		},
	}
	annotations := make(map[string]string, len(vs.ObjectMeta.Annotations)+1)
	annotations["testAnn"] = "testAnnVal"
	expected := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "inst1--vs",
			Namespace:   "default",
			Annotations: annotations,
		},
		Spec: v1alpha3.VirtualServiceSpec{
			Hosts:    []string{"*"},
			Gateways: []string{"mesh"},
		},
	}
	Annotate(vs, "testAnn", "testAnnVal")
	if diff := cmp.Diff(expected, vs); diff != "" {
		t.Errorf("Annotate (-expected, +actual)\n%v", diff)
	}
}

func TestVirtualServiceName(t *testing.T) {
	actual := RoutingVirtualServiceName("inst1")
	expected := "inst1--vs"
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("RoutingVirtualServiceName (-expected, +actual)\n%v", diff)
	}
}

func TestGatewayK8sServiceName(t *testing.T) {
	actual := GatewayK8sServiceName("inst1")
	expected := "inst1-service"
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("GatewayK8sServiceName (-expected, +actual)\n%v", diff)
	}
}

func TestGatewayNameFromInstanceName(t *testing.T) {
	actual := GatewayNameFromInstanceName("inst1")
	expected := "inst1--gateway"
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("GatewayNameFromInstanceName (-expected, +actual)\n%v", diff)
	}
}

func TestCompositeK8sServiceNameFromInstance(t *testing.T) {
	component := v1alpha2.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mycomponent1",
		},
	}
	actual := CompositeK8sServiceNameFromInstance("inst1", component)
	expected := "inst1--mycomponent1-service"
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("CompositeK8sServiceNameFromInstance (-expected, +actual)\n%v", diff)
	}
}
