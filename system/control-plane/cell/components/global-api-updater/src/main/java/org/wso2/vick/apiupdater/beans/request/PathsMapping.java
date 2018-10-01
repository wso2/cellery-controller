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

import com.fasterxml.jackson.annotation.JsonProperty;
import org.wso2.vick.apiupdater.utils.Constants;

import java.util.HashMap;
import java.util.Map;

/**
 * Class to represent paths map.
 */
public class PathsMapping {

    @JsonProperty(Constants.JsonParamNames.PATHS)
    private Map<String, PathDefinition> paths = new HashMap<>();

    @JsonProperty(Constants.JsonParamNames.SWAGGER)
    private String swagger;

    public PathsMapping() {
        swagger = Constants.Utils.SWAGGER_VERSION;
    }

    public Map<String, PathDefinition> getPaths() {
        return paths;
    }

    public String getSwagger() {
        return swagger;
    }

    public void setSwagger(String swagger) {
        this.swagger = swagger;
    }

    public void addPathDefinition(String key, PathDefinition def) {
        this.paths.put(key, def);
    }
}
