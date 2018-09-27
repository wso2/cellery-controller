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

package org.wso2.vick.apiupdater.beans.request;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import org.wso2.vick.apiupdater.utils.Constants;

import java.util.List;

/**
 * Class t represent method information.
 */
@JsonInclude(JsonInclude.Include.NON_NULL)
public class Method {

    @JsonProperty(Constants.JsonParamNames.PARAMETERS)
    private List<Parameter> parameters;

    @JsonProperty(Constants.JsonParamNames.X_AUTH_TYPE)
    private String xAuthType;

    public Method() {
        this.xAuthType = "None";
    }

    public List<Parameter> getParameters() {
        return parameters;
    }

    public void setParameters(List<Parameter> parameters) {
        this.parameters = parameters;
    }

    public String getxAuthType() {
        return xAuthType;
    }

    public void setxAuthType(String xAuthType) {
        this.xAuthType = xAuthType;
    }
}
