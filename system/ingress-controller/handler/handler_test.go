package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wso2/product-vick/system/ingress-controller/config"
	"github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"io"
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
)

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