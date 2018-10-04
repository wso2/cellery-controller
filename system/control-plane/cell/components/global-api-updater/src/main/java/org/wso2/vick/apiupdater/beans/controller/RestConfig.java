/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.vick.apiupdater.beans.controller;

import com.fasterxml.jackson.annotation.JsonProperty;
import org.json.simple.JSONObject;
import org.wso2.vick.apiupdater.utils.Constants;

/**
 * Rest configuration object.
 */
public class RestConfig {

    @JsonProperty(Constants.JsonParamNames.USERNAME)
    private String username;

    @JsonProperty(Constants.JsonParamNames.PASSWORD)
    private String password;

    @JsonProperty(Constants.JsonParamNames.API_VERSION)
    private String apiVersion;

    @JsonProperty(Constants.JsonParamNames.REGISTER_PAYLOAD)
    private JSONObject registerPayload;

    @JsonProperty(Constants.JsonParamNames.APIM_BASE_URL)
    private String apimBaseUrl;

    @JsonProperty(Constants.JsonParamNames.TOKEN_ENDPOINT)
    private String tokenEndpoint;

    @JsonProperty(Constants.JsonParamNames.TRUST_STORE)
    private JSONObject trustStore;

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public JSONObject getRegisterPayload() {
        return registerPayload;
    }

    public void setRegisterPayload(JSONObject registerPayload) {
        this.registerPayload = registerPayload;
    }

    public String getApimBaseUrl() {
        return apimBaseUrl;
    }

    public void setApimBaseUrl(String apimBaseUrl) {
        this.apimBaseUrl = apimBaseUrl;
    }

    public String getTokenEndpoint() {
        return tokenEndpoint;
    }

    public void setTokenEndpoint(String tokenEndpoint) {
        this.tokenEndpoint = tokenEndpoint;
    }

    public JSONObject getTrustStore() {
        return trustStore;
    }

    public void setTrustStore(JSONObject trustStore) {
        this.trustStore = trustStore;
    }
}
