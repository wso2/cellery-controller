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

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"strconv"
)

const (
	cellGwConfigUsernameElem = "username"
	cellGwConfigPasswordElem = "password"
	cellGwConfigApiVersionElem = "apiVersion"
	cellGwConfigRegisterPayloadElem = "registerPayload"
	cellGwConfigRegisterPayloadClientNameElem = "clientName"
	cellGwConfigRegisterPayloadOwnerElem = "owner"
	cellGwConfigRegisterPayloadGrantTypeElem = "grantType"
	cellGwConfigRegisterPayloadSaasAppElem = "saasApp"
	cellGwConfigGlobalApimBaseUrlElem = "apimBaseUrl"
	cellGwConfigGlobalApimtokenEpElem = "tokenEndpoint"
	cellGwConfigTrustStoreElem = "trustStore"
	cellGwConfigTrustStoreLocationElem = "location"
	cellGwConfigTrustStorePasswordElem = "password"
)

type IngressConfig struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ApiVersion      string `json:"apiVersion"`
	BaseUrl         string `json:"baseUrl"`
	TokenEp         string `json:"tokenEp"`
	RegisterPayload RegisterPayloadConfig `json:"registerPayload"`
	Truststore      TruststoreConfig `json:"truststore"`
}

type RegisterPayloadConfig struct {
	ClientName string `json:"clientName"`
	Owner      string `json:"owner"`
	GrantType  string `json:"grantType"`
	SaasApp    bool `json:"saasApp"`
}

type TruststoreConfig struct {
	Location string `json:"location"`
	Password string `json:"password"`
}

func GetIngressConfigs(path string) (*IngressConfig, error) {
	ingConfig, err := readGlobalAPIMConfigs(path)
	if err != nil {
		return nil, err
	}
	return validateAndParse(ingConfig)
}

func readGlobalAPIMConfigs(path string) (map[string]interface{}, error) {
	configFileByteArr, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// parse the file
	var parsedConfig map[string]interface{}
	err = json.Unmarshal([]byte(configFileByteArr), &parsedConfig)
	if err != nil {
		return nil, err
	}
	glog.Infof("Successfully read the config file at %s\n", path)
	return parsedConfig, nil
}

func validateAndParse(config map[string]interface{}) (*IngressConfig, error) {
	username := config[cellGwConfigUsernameElem]
	if username == nil {
		return nil, handleConfigNotFound(cellGwConfigUsernameElem)
	}
	password := config[cellGwConfigPasswordElem]
	if password == nil {
		return nil, handleConfigNotFound(cellGwConfigPasswordElem)
	}
	apiVersion := config[cellGwConfigApiVersionElem]
	if apiVersion == nil {
		return nil, handleConfigNotFound(cellGwConfigApiVersionElem)
	}
	registerPayload := config[cellGwConfigRegisterPayloadElem]
	isMap, registerPayloadMapStruct := isMapType(registerPayload)
	if !isMap {
		return nil, handleConfigNotFound(cellGwConfigRegisterPayloadElem)
	}
	client := registerPayloadMapStruct[cellGwConfigRegisterPayloadClientNameElem]
	if client == nil {
		return nil, handleConfigNotFound(cellGwConfigRegisterPayloadClientNameElem)
	}
	owner := registerPayloadMapStruct[cellGwConfigRegisterPayloadOwnerElem]
	if owner == nil {
		return nil, handleConfigNotFound(cellGwConfigRegisterPayloadOwnerElem)
	}
	grantType := registerPayloadMapStruct[cellGwConfigRegisterPayloadGrantTypeElem]
	if grantType == nil {
		return nil, handleConfigNotFound(cellGwConfigRegisterPayloadGrantTypeElem)
	}
	saasApp := registerPayloadMapStruct[cellGwConfigRegisterPayloadSaasAppElem]
	if saasApp == nil {
		return nil, handleConfigNotFound(cellGwConfigRegisterPayloadSaasAppElem)
	}
	saasAppBool, err := getBoolValue(fmt.Sprint(saasApp))
	if err != nil {
		return nil, handleIncorrectConfig(cellGwConfigRegisterPayloadSaasAppElem)
	}
	apimBaseUrl := config[cellGwConfigGlobalApimBaseUrlElem]
	if apimBaseUrl == nil {
		return nil, handleConfigNotFound(cellGwConfigGlobalApimBaseUrlElem)
	}
	tokenEndpoint := config[cellGwConfigGlobalApimtokenEpElem]
	if tokenEndpoint == nil {
		return nil, handleConfigNotFound(cellGwConfigGlobalApimtokenEpElem)
	}
	trustStore := config[cellGwConfigTrustStoreElem]
	isMap, trustStoreMapStruct := isMapType(trustStore)
	if !isMap {
		return nil, handleConfigNotFound(cellGwConfigTrustStoreElem)
	}
	location := trustStoreMapStruct[cellGwConfigTrustStoreLocationElem]
	if location == nil {
		return nil, handleConfigNotFound(cellGwConfigTrustStoreLocationElem)
	}
	truststorePassword := trustStoreMapStruct[cellGwConfigTrustStorePasswordElem]
	if truststorePassword == nil {
		return nil, handleConfigNotFound(cellGwConfigTrustStorePasswordElem)
	}
	glog.Infof("Successfully validated and parsed the configurations")

	// create the beans
	return &IngressConfig{
		Username:   fmt.Sprint(username),
		Password:   fmt.Sprint(password),
		ApiVersion: fmt.Sprint(apiVersion),
		BaseUrl:    fmt.Sprint(apimBaseUrl),
		TokenEp:    fmt.Sprint(tokenEndpoint),
		RegisterPayload: RegisterPayloadConfig {
			ClientName: fmt.Sprint(client),
			Owner:      fmt.Sprint(owner),
			GrantType:  fmt.Sprint(grantType),
			SaasApp:    saasAppBool,
		},
		Truststore: TruststoreConfig {
			Location: fmt.Sprint(location),
			Password: fmt.Sprint(password),
		},
	}, nil
}

func isMapType (obj interface{}) (bool, map[string]interface{}) {
	configMap, isMap := obj.(map[string]interface{})
	return isMap, configMap
}

func getBoolValue (str string) (bool, error) {
	return strconv.ParseBool(str)
}

func handleConfigNotFound (configElem string) error {
	return errors.New("Unable to extract parameter '" + configElem + "'")
}

func handleIncorrectConfig (configElem string) error {
	return errors.New("Incorrect value for '" + configElem + "'")
}

