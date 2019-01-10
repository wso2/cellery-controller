/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.auth.cell.sts.model.config;

import org.wso2.vick.auth.cell.sts.CellStsUtils;

import java.util.HashMap;
import java.util.Map;

public class CellStsConfiguration {

    private String stsEndpoint;

    private String username;

    private String password;

    private String cellName;

    private String globalJWKEndpoint;

    /**
     * Get the global JWKS endpoint
     * @return Global JWKS endpoint
     */
    public String getGlobalJWKEndpoint() {

        return globalJWKEndpoint;
    }

    /**
     * Set global JWKS endpoint.
     * @param globalJWKEndpoint.
     * @return CellStsConfiguration.
     */
    public CellStsConfiguration setGlobalJWKEndpoint(String globalJWKEndpoint) {

        this.globalJWKEndpoint = globalJWKEndpoint;
        return this;
    }

    public String getStsEndpoint() {
        return stsEndpoint;
    }

    public CellStsConfiguration setStsEndpoint(String stsEndpoint) {
        this.stsEndpoint = stsEndpoint;
        return this;
    }

    public String getUsername() {
        return username;
    }

    public CellStsConfiguration setUsername(String username) {
        this.username = username;
        return this;
    }

    public String getPassword() {
        return password;
    }

    public CellStsConfiguration setPassword(String password) {
        this.password = password;
        return this;
    }

    public String getCellName() {
        return cellName;
    }

    public CellStsConfiguration setCellName(String cellName) {
        this.cellName = cellName;
        return this;
    }

    @Override
    public String toString() {
        Map<String, String> configJson = new HashMap<>();
        configJson.put("Global STS Endpoint", stsEndpoint);
        configJson.put("Cell Name", cellName);

        return CellStsUtils.getPrettyPrintJson(configJson);
    }
}
