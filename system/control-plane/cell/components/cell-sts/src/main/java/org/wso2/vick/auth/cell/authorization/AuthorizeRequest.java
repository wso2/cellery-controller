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

package org.wso2.vick.auth.cell.authorization;

import org.wso2.vick.auth.cell.sts.model.CellStsRequest;
import org.wso2.vick.auth.cell.sts.model.RequestContext;
import org.wso2.vick.auth.cell.sts.model.RequestDestination;
import org.wso2.vick.auth.cell.sts.model.RequestSource;

import java.util.Map;

/**
 * Authorization request.
 */
public class AuthorizeRequest {

    private String requestId;
    private RequestSource source;
    private RequestDestination destination;
    private RequestContext requestContext;
    private Map<String, String> requestHeaders;
    private AuthorizationContext authorizationContext;

    public void setRequestId(String requestId) {

        this.requestId = requestId;
    }

    public void setSource(RequestSource source) {

        this.source = source;
    }

    public void setDestination(RequestDestination destination) {

        this.destination = destination;
    }

    public void setRequestContext(RequestContext requestContext) {

        this.requestContext = requestContext;
    }

    public void setRequestHeaders(Map<String, String> requestHeaders) {

        this.requestHeaders = requestHeaders;
    }

    public AuthorizationContext getAuthorizationContext() {

        return authorizationContext;
    }

    public void setAuthorizationContext(AuthorizationContext authorizationContext) {

        this.authorizationContext = authorizationContext;
    }

    public String getRequestId() {

        return requestId;
    }

    public RequestSource getSource() {

        return source;
    }

    public RequestDestination getDestination() {

        return destination;
    }

    public RequestContext getRequestContext() {

        return requestContext;
    }

    public Map<String, String> getRequestHeaders() {

        return requestHeaders;
    }

    public AuthorizationContext getUser() {

        return authorizationContext;
    }

    public AuthorizeRequest(CellStsRequest request, AuthorizationContext authorizationContext) {

        this.requestId = request.getRequestId();
        this.source = request.getSource();
        this.destination = request.getDestination();
        this.requestContext = request.getRequestContext();
        this.requestHeaders = request.getRequestHeaders();
        this.authorizationContext = authorizationContext;
    }

}
