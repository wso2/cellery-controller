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

/**
 * Represent Path definition information.
 */
@JsonInclude(JsonInclude.Include.NON_NULL)
public class PathDefinition {

    @JsonProperty(Constants.JsonParamNames.GET)
    private Method get;

    @JsonProperty(Constants.JsonParamNames.POST)
    private Method post;

    @JsonProperty(Constants.JsonParamNames.PUT)
    private Method put;

    @JsonProperty(Constants.JsonParamNames.DELETE)
    private Method delete;

    public Method getGet() {
        return get;
    }

    public void setGet(Method get) {
        this.get = get;
    }

    public Method getPost() {
        return post;
    }

    public void setPost(Method post) {
        this.post = post;
    }

    public Method getPut() {
        return put;
    }

    public void setPut(Method put) {
        this.put = put;
    }

    public Method getDelete() {
        return delete;
    }

    public void setDelete(Method delete) {
        this.delete = delete;
    }
}
