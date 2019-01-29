/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.vick.auth.cell.authorization.opa;

import com.google.gson.Gson;
import com.google.gson.internal.LinkedTreeMap;
import com.mashape.unirest.http.HttpResponse;
import com.mashape.unirest.http.JsonNode;
import com.mashape.unirest.http.Unirest;
import com.mashape.unirest.http.exceptions.UnirestException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.authorization.AuthorizationFailedException;
import org.wso2.vick.auth.cell.authorization.AuthorizationHandler;
import org.wso2.vick.auth.cell.authorization.AuthorizationUtils;
import org.wso2.vick.auth.cell.authorization.AuthorizeRequest;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;
import org.wso2.vick.auth.cell.utils.LambdaExceptionUtils;

import java.util.HashMap;
import java.util.Map;

/**
 * Calls local OPA server and validate the request.
 */
public class OPAAuthorizationHandler implements AuthorizationHandler {

    private static final Logger log = LoggerFactory.getLogger(OPAAuthorizationHandler.class);

    @Override
    public void authorize(AuthorizeRequest request) throws AuthorizationFailedException {
        log.debug("OPA authorization handler invoked for request id: {}", request.getRequestId());
        request.setAuthorizationContext(new OPAAuthorizationContext(request.getAuthorizationContext().getJwt()));
        Gson gson = new Gson();
        String requestString = gson.toJson(request);
        HttpResponse<JsonNode> apiResponse = null;
        log.info("Reqesut to OPA server : {}", requestString);
        try {
            apiResponse = Unirest.post(AuthorizationUtils.getOPAEndpoint()).body(requestString).asJson();
            log.info("Response from OPA server: {}" + apiResponse.getBody().toString());
            Object results = apiResponse.getBody().getObject().get("result");

            Map resultsMap = gson.fromJson(results.toString(), HashMap.class);
            resultsMap.forEach(LambdaExceptionUtils.rethrowBiConsumer((key, value) -> {
                if (Boolean.parseBoolean(((LinkedTreeMap) value).getOrDefault("deny", false).toString())) {
                    throw new AuthorizationFailedException("Error while authorizing request. Decision found : " +
                            value);
                }
            }));

            log.info("Authorization successfully completed for request: ", request.getRequestId());
        } catch (UnirestException | VickCellSTSException e) {
            throw new AuthorizationFailedException("Error while sending authorization request to OPA", e);
        }
    }
}
