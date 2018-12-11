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
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetIngressConfigsWrongFilePath(t *testing.T) {
	_, err := GetIngressConfigs("conf.json")
	if err == nil {
		t.Fatal("expected error 'no such file or directory' not recieved")
	}
}

func TestGetIngressConfigsCorruptedFile(t *testing.T) {
	_, err := GetIngressConfigs("corrupted_config.json")
	if err == nil {
		t.Fatal("expected error 'malformed config file' not recieved")
	}
}

func TestGetNonCompleteIngressConfigs(t *testing.T) {
	_, err := GetIngressConfigs("incomplete_config.json")
	if err == nil {
		t.Fatal("expected error 'Unable to extract parameter' not recieved")
	}
}

func TestGetIngressConfigs(t *testing.T) {
	expectedConfig := &IngressConfig{
		Username:   fmt.Sprint("admin"),
		Password:   fmt.Sprint("admin"),
		ApiVersion: fmt.Sprint("v0.14"),
		BaseUrl:    fmt.Sprint("https://localhost:9443"),
		TokenEp:    fmt.Sprint("https://localhost:8243/token"),
		RegisterPayload: RegisterPayloadConfig {
			ClientName: fmt.Sprint("rest_api_publisher"),
			Owner:      fmt.Sprint("admin"),
			GrantType:  fmt.Sprint("password refresh_token"),
			SaasApp:    true,
		},
		Truststore: TruststoreConfig {
			Location: fmt.Sprint("lib/platform/bre/security/ballerinaTruststore.p12"),
			Password: fmt.Sprint("ballerina"),
		},
	}

	actualConfig, err := GetIngressConfigs("config.json")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(expectedConfig, actualConfig); diff != "" {
		t.Errorf("diff exists (-expectedConfig, +actualConfig) \n%v", diff)
	}
}