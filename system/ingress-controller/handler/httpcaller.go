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
	"github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"net/http"
)

type HttpCaller interface {
	DoHttpCall() (interface{}, error)
}

type DefaultHttpCaller struct {
	client *http.Client
	req *http.Request
}

func (caller *DefaultHttpCaller) DoHttpCall() (interface{}, error)  {
	return caller.client.Do(caller.req)
}

func registerClientHttpCaller(ingRuleCreator *IngressRuleCreator, client *http.Client) (*DefaultHttpCaller, error) {
	// create headers
	authHeader := getBasicAuthHeader(ingRuleCreator.IngressConfig.Username, ingRuleCreator.IngressConfig.Password)
	// create request
	requestBody, err := json.Marshal(ingRuleCreator.IngressConfig.RegisterPayload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + clientRegistrationPath +
		ingRuleCreator.IngressConfig.ApiVersion + clientRegisterPath, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Basic " + authHeader)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return &DefaultHttpCaller{client:client, req:request}, nil
}

func generateTokenHttpCaller(ingRuleCreator *IngressRuleCreator, client *http.Client,
								clientId string, clientSecret string) (*DefaultHttpCaller, error) {

	// get basic auth header for registered client, with clientId and clientSecret
	authHeader := getBasicAuthHeader(clientId, clientSecret)
	payload := "grant_type=password&username=" + ingRuleCreator.IngressConfig.Username + "&password=" +
		ingRuleCreator.IngressConfig.Password + "&scope=" + tokenScopes
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.TokenEp,
		bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Basic " + authHeader)
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	return &DefaultHttpCaller{client:client, req:request}, nil
}

func createApiHttpCaller(ingRuleCreator *IngressRuleCreator, client *http.Client, ingSpec *v1alpha1.IngressSpec,
																token string) (*DefaultHttpCaller, error) {

	err := validateAndSetDefaults(ingSpec)
	if err != nil {
		return nil, err
	}
	createApiPayload, err := createApiCreationPayload(ingSpec)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath, bytes.NewBuffer([]byte(createApiPayload)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return &DefaultHttpCaller{client:client, req:request}, nil
}

func publishApiHttpCaller(ingRuleCreator *IngressRuleCreator, client *http.Client, apiId string,
	token string) (*DefaultHttpCaller, error) {

	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + lifecylePath + "?apiId=" + apiId + "&action=Publish", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return &DefaultHttpCaller{client:client, req:request}, nil
}

func updateApiHttpCaller(ingRuleCreator *IngressRuleCreator, ingSpec *v1alpha1.IngressSpec, client *http.Client,
										apiId string, token string) (*DefaultHttpCaller, error) {

	err := validateAndSetDefaults(ingSpec)
	if err != nil {
		return nil, err
	}
	createApiPayload, err := createApiCreationPayload(ingSpec)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("PUT", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "/" + apiId, bytes.NewBuffer([]byte(createApiPayload)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return &DefaultHttpCaller{client:client, req:request}, nil
}

func getApiIdHttpCaller(ingRuleCreator *IngressRuleCreator, context string,
												client *http.Client, token string) (*DefaultHttpCaller, error) {

	request, err := http.NewRequest("GET", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "?query=context:" + context, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return &DefaultHttpCaller{client:client, req:request}, nil
}

func deleteApiHttpCaller(ingRuleCreator *IngressRuleCreator, client *http.Client, apiId string,
																token string) (*DefaultHttpCaller, error) {

	request, err := http.NewRequest("DELETE", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "/" + apiId, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	return &DefaultHttpCaller{client:client, req:request}, nil
}

//type HttpReponseDecoder interface {
//	Decode(res *http.Response) (interface{}, error)
//}
//
//type RegisterClientHttpResDecoder struct {
//}

//func (decoder *RegisterClientHttpResDecoder) Decode(res *http.Response) (interface{}, error)  {
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return &RegisterClientResults{}, err
//	}
//	var regResponse registrationResponse
//	err = json.Unmarshal(body, &regResponse)
//	if err != nil {
//		return &RegisterClientResults{}, err
//	}
//	if len(regResponse.ClientId) == 0 || len(regResponse.ClientSecret) == 0 {
//		return &RegisterClientResults{}, errors.New("Client registration response does not contain id/secret")
//	}
//	return &RegisterClientResults{regResponse.ClientId, regResponse.ClientSecret}, nil
//}

//type GenTokenResponseDecoder struct {
//}
//
//func (decoder *GenTokenResponseDecoder) Decode(res *http.Response) (interface{}, error)  {
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return "", err
//	}
//	var tokResponse tokenResponse
//	err = json.Unmarshal(body, &tokResponse)
//	if err != nil {
//		return "", err
//	}
//	// check if there is an error
//	if len(tokResponse.Error) != 0 {
//		return "", errors.New("Error in retrieving token: " + tokResponse.Error + ", "+ tokResponse.ErrorDesc)
//	}
//	if len(tokResponse.AccessToken) == 0  {
//		return "", errors.New("Token response does not contain access token")
//	}
//	// the token scopes should match the expected scope set
//	glog.Infof("Scopes of the access token: %s", tokenScopes)
//	return tokResponse.AccessToken, nil
//}
//
//type ApiCreateResponseDecoder struct {
//	ApiName string
//	ApiContext string
//	ApiVersion string
//}
//
//func (decoder *ApiCreateResponseDecoder) Decode(res *http.Response) (interface{}, error)  {
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return "", err
//	}
//	glog.Infof("Api creation response body: %+v\n", body)
//
//	if res.StatusCode == 201 {
//		// api created successfully
//		glog.Infof("Api with name: %s, context: %s, version: %s created successfully", decoder.ApiName,
//			decoder.ApiContext, decoder.ApiVersion)
//		var creatApiRes createApiSuccessResponse
//		err = json.Unmarshal(body, &creatApiRes)
//		if err != nil {
//			return "", err
//		}
//		return creatApiRes.Id, nil
//
//	} else {
//		glog.Infof("Api creation failed for api: %s, context: %s, version: %s", decoder.ApiName,
//			decoder.ApiContext, decoder.ApiVersion)
//		// marshal the error response to get the error message
//		var createApiErrorRes createApiErrorResponse
//		err = json.Unmarshal(body, &createApiErrorRes)
//		if err != nil {
//			return "", err
//		}
//		glog.Info(createApiErrorRes)
//		return "", errors.New("Api creation failed for api: " + decoder.ApiName + ", context: " + decoder.ApiContext +
//			", version: " + decoder.ApiVersion + ", error code: " + createApiErrorRes.Code + ", error msg: " +
//			createApiErrorRes.Message + ", error desc: " + createApiErrorRes.Description)
//	}
//}
//
//type ApiUpdateResponseDecoder struct {
//	ApiName string
//	ApiContext string
//	ApiVersion string
//}
//
//func (decoder *ApiUpdateResponseDecoder) Decode(res *http.Response) (interface{}, error)  {
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return "", err
//	}
//
//	if res.StatusCode == 200 {
//		// api created successfully
//		glog.Infof("Api with name: %s, context: %s, version: %s updated successfully", decoder.ApiName,
//			decoder.ApiContext, decoder.ApiVersion)
//		var updateApiSuccessRes createApiSuccessResponse
//		err = json.Unmarshal(body, &updateApiSuccessRes)
//		if err != nil {
//			return "", err
//		}
//		return updateApiSuccessRes.Id, nil
//	} else {
//		glog.Infof("Api update failed for api: %s, context: %s, version: %s", decoder.ApiName,
//			decoder.ApiContext, decoder.ApiVersion)
//		// marshal the error response to get the error message
//		var publishApiErrorRes updateApiErrorResponse
//		err = json.Unmarshal(body, &publishApiErrorRes)
//		if err != nil {
//			return "", err
//		}
//		glog.Info(publishApiErrorRes)
//		return "", errors.New("Api update failed for api: " + decoder.ApiName + ", context: " + decoder.ApiContext +
//			", version: " + decoder.ApiVersion + ", error code: " + strconv.Itoa(publishApiErrorRes.Code) + ", error msg: " +
//			publishApiErrorRes.Message + ", error desc: " + publishApiErrorRes.Description)
//
//	}
//}
//
//type ApiPublishResponseDecoder struct {
//	ApiId string
//}
//
//func (decoder *ApiPublishResponseDecoder) Decode(res *http.Response) (interface{}, error)  {
//	if res.StatusCode == 200 {
//		glog.Infof("Api with id %s published successfully", decoder.ApiId)
//	} else {
//		// read response
//		body, err := ioutil.ReadAll(res.Body)
//		if err != nil {
//			return nil, err
//		}
//		// marshal the error response to get the error message
//		var publishApiErrorRes publisheApiErrorResponse
//		err = json.Unmarshal(body, &publishApiErrorRes)
//		if err != nil {
//			return nil, err
//		}
//		return nil, errors.New("Api update failed for api " + decoder.ApiId + ", error code: " +
//			strconv.Itoa(publishApiErrorRes.Code) + ", error msg: " + publishApiErrorRes.Message + ", error desc: " +
//			publishApiErrorRes.Description)
//	}
//	return nil, nil
//}
//
