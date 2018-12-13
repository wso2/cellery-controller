/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.observability.api;

import javax.ws.rs.core.Response;

/**
 * This class holds the util methods that can be shared by different classes.
 */
public class Utils {
    private Utils() {
    }

    public static void addAllowOriginHeader(Response.ResponseBuilder responseBuilder, String value) {
        responseBuilder.header(Constants.ACCESS_CONTROL_ALLOW_ORIGIN, value);
    }

    public static void addAllowCredentialsHeader(Response.ResponseBuilder responseBuilder, String value) {
        responseBuilder.header(Constants.ACCESS_CONTROL_ALLOW_CREDENTIALS, value);
    }

    public static void addAllowMethodsHeader(Response.ResponseBuilder responseBuilder, String value) {
        responseBuilder.header(Constants.ACCESS_CONTROL_ALLOW_METHODS, value);
    }

    public static void addAllowHeaders(Response.ResponseBuilder responseBuilder, String value) {
        responseBuilder.header(Constants.ACCESS_CONTROL_ALLOW_HEADERS, value);
    }

    public static void addCorsResponseBuilder(Response.ResponseBuilder responseBuilder, String httpMethod) {
        Utils.addAllowOriginHeader(responseBuilder, Constants.ALL_ORIGIN);
        Utils.addAllowCredentialsHeader(responseBuilder, Boolean.TRUE.toString());
        Utils.addAllowMethodsHeader(responseBuilder, httpMethod);
        Utils.addAllowHeaders(responseBuilder, Constants.ORIGIN);
    }
}
