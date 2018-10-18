/*
 *  Copyright (c) 2018 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 */

package org.wso2.vick.sts.core;

import org.json.simple.JSONObject;

/**
 * This is the POJO class that represents STS response in Vick.
 */
public class VickSTSResponse {

    private String stsToken;

    public String getStsToken() {

        return stsToken;
    }

    public void setStsToken(String stsToken) {

        this.stsToken = stsToken;
    }

    public String toJson() {

        JSONObject tokenResponse = new JSONObject();
        tokenResponse.put(VickSTSConstants.VickSTSResponse.STS_TOKEN, stsToken);
        return tokenResponse.toString();
    }
}
