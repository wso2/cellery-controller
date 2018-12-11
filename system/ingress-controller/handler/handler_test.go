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

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/wso2/product-vick/system/ingress-controller/config"
	"github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	clientId = "abc123"
	clientSecret = "def678"
	token = "abc123def678"
	apiId = "123rty567"
	apiName = "TestApi1"
	apiContext = "/test"
	apiVersion = "1.0"
	postOperation = "POST"
	protocol = "https"
	epType = "production"
	serviceName = "testSvc"
	port = 9090
)

// IngressRuleCreator tests

var ingConfig, _ = config.GetIngressConfigs("gw.json")
var ingRuleCreator = NewIngressRuleCreator(ingConfig)

type dummyReader struct {
	io.Reader
}

func (dummyReader) Close() error {
	return nil
}

type dummyRegisterCaller struct {}

func (caller *dummyRegisterCaller) DoHttpCall() (interface{}, error)  {
	regResp := registrationResponse{ClientId:clientId, ClientSecret:clientSecret}
	regRespBytes, err := json.Marshal(regResp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &http.Response{Body:dummyReader{bytes.NewBuffer(regRespBytes)} }, nil
}

func TestIngressRuleCreator_RegisterClient(t *testing.T) {
	id, secret, err := ingRuleCreator.RegisterClient(&dummyRegisterCaller{})
	if err != nil {
		t.Fatal(err)
	}
	if id != clientId {
		t.Fatal("Client id expected: " + clientId + ", actual: " + id)
	}
	if secret != clientSecret {
		t.Fatal("Client secret expected: " + clientSecret + ", actual: " + secret)
	}
}

type dummyGenTokenCaller struct {}

func (caller *dummyGenTokenCaller) DoHttpCall() (interface{}, error)  {
	tokResp := tokenResponse{AccessToken:token, Scope:tokenScopes}
	tokRespBytes, err := json.Marshal(tokResp)
	if err != nil {
		return nil, err
	}
	return &http.Response{Body:dummyReader{bytes.NewBuffer(tokRespBytes)} }, nil
}

func TestIngressRuleCreator_GenerateAccessToken(t *testing.T) {
	tok, err := ingRuleCreator.GenerateAccessToken(&dummyGenTokenCaller{})
	if err != nil {
		t.Fatal(err)
	}
	if tok != token {
		t.Fatal("Token expected: " + token + ", actual: " + tok)
	}
}

type dummyCreateApiCaller struct {}

func (caller *dummyCreateApiCaller) DoHttpCall() (interface{}, error)  {
	apiCreateSuccessResp := createApiSuccessResponse{Id:apiId}
	tokRespBytes, err := json.Marshal(apiCreateSuccessResp)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer(tokRespBytes)} }
	resp.StatusCode = 201
	return resp, nil
}

func TestIngressRuleCreator_CreateApi(t *testing.T) {
	id, err := ingRuleCreator.CreateApi(&dummyCreateApiCaller{},
	&v1alpha1.IngressSpec{Name:apiName, Context:apiContext, Version:apiVersion})
	if err != nil {
		t.Fatal(err)
	}
	if id != apiId {
		t.Fatal("Api id expected: " + apiId + ", actual: " + id)
	}
}

type dummyCreateApiFailuerCaller struct {}

func (caller *dummyCreateApiFailuerCaller) DoHttpCall() (interface{}, error)  {
	apiCreateFailureResp := createApiErrorResponse{Code:"10100", Message:"Failed", Description:"Api creation failure"}
	apiCreateFailureRespBytes, err := json.Marshal(apiCreateFailureResp)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer(apiCreateFailureRespBytes)} }
	resp.StatusCode = 500
	return resp, nil
}

func TestIngressRuleCreator_CreateApiFailure(t *testing.T) {
	_, err := ingRuleCreator.CreateApi(&dummyCreateApiFailuerCaller{},
		&v1alpha1.IngressSpec{Name:apiName, Context:apiContext, Version:apiVersion})
	t.Logf("Expected error: %+v\n", err)
	if err == nil {
		t.Fatal("Expected error not recieved")
	}
}

type dummyPublishApiCaller struct {}

func (caller *dummyPublishApiCaller) DoHttpCall() (interface{}, error)  {
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer([]byte(apiId))} }
	resp.StatusCode = 200
	return resp, nil
}

func TestIngressRuleCreator_PublishApi(t *testing.T) {
	err := ingRuleCreator.PublishApi(&dummyPublishApiCaller{}, apiId)
	if err != nil {
		t.Fatal(err)
	}
}

type dummyPublishApiFailureCaller struct {}

func (caller *dummyPublishApiFailureCaller) DoHttpCall() (interface{}, error)  {
	apiPublishFailedResp := publisheApiErrorResponse{Code:100300, Message:"Failed", Description:"API Publisher failed"}
	apiPublishFailedRespBytes, err := json.Marshal(apiPublishFailedResp)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer(apiPublishFailedRespBytes)} }
	resp.StatusCode = 500
	return resp, nil
}

func TestIngressRuleCreator_PublishApiFailure(t *testing.T) {
	err := ingRuleCreator.PublishApi(&dummyPublishApiFailureCaller{}, apiId)
	t.Logf("Expected error: %+v\n", err)
	if err == nil {
		t.Fatal("Expected error not recieved")
	}
}

type dummyUpdateCaller struct {}

func (caller *dummyUpdateCaller) DoHttpCall() (interface{}, error)  {
	apiUpdateResp := createApiSuccessResponse{Id: apiId, Name:apiName, Context:apiContext, Version:apiVersion}
	apiUpdateRespBytes, err := json.Marshal(apiUpdateResp)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer(apiUpdateRespBytes)} }
	resp.StatusCode = 200
	return resp, nil
}

func TestIngressRuleCreator_UpdateApi(t *testing.T) {
	id, err := ingRuleCreator.UpdateApi(&dummyUpdateCaller{},
		&v1alpha1.IngressSpec{Name:apiName, Context:apiContext, Version:apiVersion})
	if err != nil {
		t.Fatal(err)
	}
	if id != apiId {
		t.Fatal("Api id expected: " + apiId + ", actual: " + id)
	}
}

type dummyUpdateFailureCaller struct {}

func (caller *dummyUpdateFailureCaller) DoHttpCall() (interface{}, error)  {
	apiUpdateResp := updateApiErrorResponse{Code:100201, Message:"Failed", Description:"API update failed"}
	apiUpdateRespBytes, err := json.Marshal(apiUpdateResp)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer(apiUpdateRespBytes)} }
	resp.StatusCode = 500
	return resp, nil
}

func TestIngressRuleCreator_UpdateApiFailure(t *testing.T) {
	_, err := ingRuleCreator.UpdateApi(&dummyUpdateFailureCaller{},
		&v1alpha1.IngressSpec{Name:apiName, Context:apiContext, Version:apiVersion})
	t.Logf("Expected error: %+v\n", err)
	if err == nil {
		t.Fatal("Expected error not recieved")
	}
}

type dummyGetIdCaller struct {}

func (caller *dummyGetIdCaller) DoHttpCall() (interface{}, error) {
	v := version{Id:apiId, Version:apiVersion}
	getApiIdResp := searchApiResult{List:[]version{v}}
	getApiIdRespBytes, err := json.Marshal(getApiIdResp)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer(getApiIdRespBytes)} }
	resp.StatusCode = 200
	return resp, nil
}

func TestIngressRuleCreator_GetApiId(t *testing.T) {
	id, err := ingRuleCreator.GetApiId(&dummyGetIdCaller{}, apiContext, apiVersion)
	if err != nil {
		t.Fatal(err)
	}
	if id != apiId {
		t.Fatal("Api id expected: " + apiId + ", actual: " + id)
	}
}

type dummyApiDeleteCaller struct {}

func (caller *dummyApiDeleteCaller) DoHttpCall() (interface{}, error) {
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer([]byte(apiId))} }
	resp.StatusCode = 200
	return resp, nil
}

func TestIngressRuleCreator_DeleteApi(t *testing.T) {
	err := ingRuleCreator.DeleteApi(&dummyApiDeleteCaller{}, apiId)
	if err != nil {
		t.Fatal(err)
	}
}

type dummyApiDeleteFailureCaller struct {}

func (caller *dummyApiDeleteFailureCaller) DoHttpCall() (interface{}, error) {
	resp := &http.Response{Body:dummyReader{bytes.NewBuffer([]byte(apiId))} }
	resp.StatusCode = 404
	return resp, nil
}

func TestIngressRuleCreator_DeleteApiFailure(t *testing.T) {
	err := ingRuleCreator.DeleteApi(&dummyApiDeleteFailureCaller{}, apiId)
	t.Logf("Expected error: %+v\n", err)
	if err == nil {
		t.Fatal("Expected error not recieved")
	}
}

// HttpCaller tests
func TestRegisterClientHttpCaller(t *testing.T) {
	authHeader := getBasicAuthHeader(ingRuleCreator.IngressConfig.Username, ingRuleCreator.IngressConfig.Password)
	// create request
	requestBody, err := json.Marshal(ingRuleCreator.IngressConfig.RegisterPayload)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + clientRegistrationPath +
		ingRuleCreator.IngressConfig.ApiVersion + clientRegisterPath, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Basic " + authHeader)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	httpCaller, err := registerClientHttpCaller(ingRuleCreator, httpClient)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	compare(t, request, httpCaller, requestBody, "registerClientHttpCaller")
}

func TestGenerateTokenHttpCaller(t *testing.T) {
	authHeader := getBasicAuthHeader(clientId, clientSecret)
	payload := "grant_type=password&username=" + ingRuleCreator.IngressConfig.Username + "&password=" +
		ingRuleCreator.IngressConfig.Password + "&scope=" + tokenScopes
	requestBody := []byte(payload)
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.TokenEp,
		bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Basic " + authHeader)
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	httpCaller, err := generateTokenHttpCaller(ingRuleCreator, httpClient, clientId, clientSecret)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	compare(t, request, httpCaller, requestBody, "generateTokenHttpCaller")
}

func TestCreateApiHttpCaller(t *testing.T) {

	paths := []v1alpha1.Path{{Context:apiContext, Operation:postOperation}}
	endpoints := []v1alpha1.Endpoint{{Protocol:protocol, Type:epType, ServiceName:serviceName, Port:port}}
	ingSpec := &v1alpha1.IngressSpec{Name:apiName, Context:apiContext, Version:apiVersion, Paths:paths, Endpoints:endpoints}
	err := validateAndSetDefaults(ingSpec)
	if err != nil {
		t.Fatal(err)
	}
	createApiPayload, err := createApiCreationPayload(ingSpec)
	if err != nil {
		t.Fatal(err)
	}
	requestBody := []byte(createApiPayload)
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	httpCaller, err := createApiHttpCaller(ingRuleCreator, httpClient, ingSpec, token)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	compare(t, request, httpCaller, requestBody, "createApiHttpCaller")
}

func TestPublishApiHttpCaller(t *testing.T) {

	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + lifecylePath + "?apiId=" + apiId + "&action=Publish", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	httpCaller, err := publishApiHttpCaller(ingRuleCreator, httpClient, apiId, token)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	compare(t, request, httpCaller, nil, "publishApiHttpCaller")
}

func TestUpdateApiHttpCaller(t *testing.T) {

	paths := []v1alpha1.Path{{Context:apiContext, Operation:postOperation}}
	endpoints := []v1alpha1.Endpoint{{Protocol:protocol, Type:epType, ServiceName:serviceName, Port:port}}
	ingSpec := &v1alpha1.IngressSpec{Name:apiName, Context:apiContext, Version:apiVersion, Paths:paths, Endpoints:endpoints}
	err := validateAndSetDefaults(ingSpec)
	if err != nil {
		t.Fatal(err)
	}
	createApiPayload, err := createApiCreationPayload(ingSpec)
	if err != nil {
		t.Fatal(err)
	}
	requestBody := []byte(createApiPayload)
	request, err := http.NewRequest("PUT", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "/" + apiId, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	httpCaller, err := updateApiHttpCaller(ingRuleCreator, ingSpec, httpClient, apiId, token)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	compare(t, request, httpCaller, requestBody, "updateApiHttpCaller")
}

func TestGetApiIdHttpCaller (t *testing.T) {

	request, err := http.NewRequest("GET", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "?query=context:" + apiContext, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	httpCaller, err := getApiIdHttpCaller(ingRuleCreator, apiContext, httpClient, token)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	compare(t, request, httpCaller, nil, "getApiIdHttpCaller")
}

func TestDeleteApiHttpCaller (t *testing.T) {
	request, err := http.NewRequest("DELETE", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "/" + apiId, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	httpCaller, err := deleteApiHttpCaller(ingRuleCreator, httpClient,apiId, token)
	if err != nil {
		t.Fatal(err)
	}
	// verify
	compare(t, request, httpCaller, nil, "deleteApiHttpCaller")
}

func compare (t *testing.T, request *http.Request, httpCaller *DefaultHttpCaller, requestBody []byte, method string) {
	if diff := cmp.Diff(request.Header, httpCaller.req.Header); diff != "" {
		t.Errorf("diff exists (-request.Header, +httpCaller.req.Header) " + method + " method \n%v", diff)
	}
	if requestBody != nil {
		body, err := ioutil.ReadAll(httpCaller.req.Body)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(requestBody, body); diff != "" {
			t.Errorf("diff exists (-requestBody, +body) "+method+" method \n%v", diff)
		}
	}
	if diff := cmp.Diff(request.Method, httpCaller.req.Method); diff != "" {
		t.Errorf("diff exists (-request.Method, +httpCaller.req.Method) " + method + " method \n%v", diff)
	}
	if diff := cmp.Diff(request.URL, httpCaller.req.URL); diff != "" {
		t.Errorf("diff exists (-request.URL, +httpCaller.req.URL) " + method + " method \n%v", diff)
	}
}