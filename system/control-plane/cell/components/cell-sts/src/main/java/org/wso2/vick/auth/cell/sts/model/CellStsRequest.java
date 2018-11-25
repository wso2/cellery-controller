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
package org.wso2.vick.auth.cell.sts.model;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

public class CellStsRequest {

    private String requestId;
    private RequestSource source;
    private RequestDestination destination;
    private RequestContext requestContext;
    private Map<String, String> requestHeaders = new HashMap<>();

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
        return Collections.unmodifiableMap(requestHeaders);
    }

    public String getRequestId() {
        return requestId;
    }

    private CellStsRequest() {
    }

    public static class CellStsRequestBuilder {

        private String requestId;
        private RequestSource source;
        private RequestDestination destination;
        private RequestContext requestContext;
        private Map<String, String> requestHeaders = new HashMap<>();


        public CellStsRequestBuilder setSource(RequestSource source) {
            this.source = source;
            return this;
        }

        public CellStsRequestBuilder setDestination(RequestDestination destination) {
            this.destination = destination;
            return this;
        }

        public CellStsRequestBuilder setRequestContext(RequestContext requestContext) {
            this.requestContext = requestContext;
            return this;
        }

        public CellStsRequestBuilder setRequestHeaders(Map<String, String> requestHeaders) {
            this.requestHeaders = requestHeaders;
            return this;
        }

        public CellStsRequestBuilder setRequestId(String requestId) {
            this.requestId = requestId;
            return this;
        }

        public CellStsRequest build() {
            CellStsRequest request = new CellStsRequest();
            request.requestId = requestId;
            request.source = source;
            request.destination = destination;
            request.requestContext = requestContext;
            request.requestHeaders = requestHeaders;

            return request;
        }
     }

}
