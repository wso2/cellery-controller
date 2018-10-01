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
import org.wso2.vick.apiupdater.utils.Constants;

/**
 * Class to represent API.
 */
public class API {

    @JsonProperty(Constants.JsonParamNames.CONTEXT)
    private String context;

    @JsonProperty(Constants.JsonParamNames.DEFINITION)
    private ApiDefinition[] definitions;

    @JsonProperty(Constants.JsonParamNames.BACKEND)
    private String backend;

    @JsonProperty(Constants.JsonParamNames.GLOBAL)
    private boolean global;

    public String getContext() {
        return context;
    }

    public void setContext(String context) {
        this.context = context;
    }

    public ApiDefinition[] getDefinitions() {
        return definitions;
    }

    public void setDefinitions(ApiDefinition[] definitions) {
        this.definitions = definitions;
    }

    public String getBackend() {
        return backend;
    }

    public void setBackend(String backend) {
        this.backend = backend;
    }

    public boolean isGlobal() {
        return global;
    }

    public void setGlobal(boolean global) {
        this.global = global;
    }
}
