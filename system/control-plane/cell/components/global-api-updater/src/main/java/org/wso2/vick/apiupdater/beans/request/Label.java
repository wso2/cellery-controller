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

import java.util.List;

/**
 * Class to represent label.
 */
public class Label {

    @JsonProperty(Constants.JsonParamNames.NAME)
    private String name;

    @JsonProperty(Constants.JsonParamNames.DESCRIPTION)
    private String description;

    @JsonProperty(Constants.JsonParamNames.ACCESS_URLS)
    private List accessUrls;

    public Label() {
        description = "cell label";
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public List getAccessUrls() {
        return accessUrls;
    }

    public void setAccessUrls(List accessUrls) {
        this.accessUrls = accessUrls;
    }
}