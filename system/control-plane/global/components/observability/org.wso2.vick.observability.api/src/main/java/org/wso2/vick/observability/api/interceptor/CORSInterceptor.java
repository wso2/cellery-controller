/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.wso2.vick.observability.api.interceptor;

import org.apache.commons.lang3.StringUtils;
import org.wso2.msf4j.Request;
import org.wso2.msf4j.Response;
import org.wso2.msf4j.interceptor.RequestInterceptor;
import org.wso2.vick.observability.api.Constants;

import javax.ws.rs.HttpMethod;
import javax.ws.rs.core.HttpHeaders;

/**
 * Used for applying header for allowing cross origin requests.
 * <p>
 * This alone is not enough since the browsers send an OPTIONS HTTP call to the endpoint first to check
 * which methods are allowed. Therefore an endpoint with HTTP Method OPTIONS should be added to all services.
 */
public class CORSInterceptor implements RequestInterceptor {
    @Override
    public boolean interceptRequest(Request request, Response response) {
        response.setHeader(Constants.ACCESS_CONTROL_ALLOW_METHODS, HttpMethod.GET + "," + HttpMethod.POST +
                "," + HttpMethod.PUT + "," + HttpMethod.DELETE);
        response.setHeader(Constants.ACCESS_CONTROL_MAX_AGE, Constants.MAX_AGE);
        response.setHeader(Constants.ACCESS_CONTROL_ALLOW_HEADERS, HttpHeaders.CONTENT_TYPE);

        if (StringUtils.isNotBlank(request.getHeader(Constants.ORIGIN))) {
            response.setHeader(Constants.ACCESS_CONTROL_ALLOW_ORIGIN, request.getHeader(Constants.ORIGIN));
        } else {
            response.setHeader(Constants.ACCESS_CONTROL_ALLOW_ORIGIN, Constants.ALL_ORIGIN);
        }
        return true;
    }
}
