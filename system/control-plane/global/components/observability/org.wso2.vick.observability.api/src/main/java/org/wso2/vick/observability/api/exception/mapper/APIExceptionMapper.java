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

package org.wso2.vick.observability.api.exception.mapper;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.google.gson.JsonPrimitive;
import org.wso2.vick.observability.api.exception.InvalidParamException;

import javax.ws.rs.core.HttpHeaders;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.ext.ExceptionMapper;

/**
 * Exception Mapper for mapping Server Error Exceptions.
 */
public class APIExceptionMapper implements ExceptionMapper {
    private Gson gson = new Gson();

    @Override
    public Response toResponse(Throwable throwable) {
        Response.Status status;
        String message;

        if (throwable instanceof InvalidParamException) {
            status = Response.Status.PRECONDITION_FAILED;
            message = throwable.getMessage();
        } else {
            status = Response.Status.INTERNAL_SERVER_ERROR;
            message = "Unknown Error Occurred";
        }

        JsonObject errorResponseJsonObject = new JsonObject();
        errorResponseJsonObject.add("message", new JsonPrimitive(message));

        return Response.status(status)
                .header(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_JSON)
                .entity(gson.toJson(errorResponseJsonObject))
                .build();
    }
}
