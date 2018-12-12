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

package org.wso2.vick.observability.api;

import org.wso2.vick.observability.api.siddhi.SiddhiStoreQueryTemplates;

import javax.ws.rs.GET;
import javax.ws.rs.OPTIONS;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

/**
 * MSF4J service for fetching distributed tracing data.
 */
@Path("/api/tracing")
public class DistributedTracingAPI {

    @GET
    @Path("/metadata")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getMetadata() {
        Object[][] results = SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_METADATA.builder()
                .build()
                .execute();
        return Response.ok().entity(results).build();
    }

    @OPTIONS
    @Path("/metadata")
    public Response getMetadataOptions() {
        return Response.ok().build();
    }

    @GET
    @Path("/search")
    @Produces(MediaType.APPLICATION_JSON)
    public Response search(@QueryParam("cell") String cell,
                           @QueryParam("serviceName") String serviceName,
                           @QueryParam("operationName") String operationName,
                           @QueryParam("minDuration") String minDuration,
                           @QueryParam("maxDuration") String maxDuration,
                           @QueryParam("queryStartTime") String queryStartTime,
                           @QueryParam("queryEndTime") String queryEndTime) {
        Object[][] results = SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_SEARCH_GET_TRACE_IDS.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.CELL, cell)
                .setArg(SiddhiStoreQueryTemplates.Params.SERVICE_NAME, serviceName)
                .setArg(SiddhiStoreQueryTemplates.Params.OPERATION_NAME, operationName)
                .setArg(SiddhiStoreQueryTemplates.Params.MIN_DURATION, minDuration)
                .setArg(SiddhiStoreQueryTemplates.Params.MAX_DURATION, maxDuration)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .build()
                .execute();

        String[] selectedTraceIds = new String[results.length];
        for (int i = 0; i < results.length; i++) {
            selectedTraceIds[i] = (String) results[i][0];
        }

        return Response.ok().entity(selectedTraceIds).build();
    }

    @OPTIONS
    @Path("/search")
    public Response searchOptions() {
        return Response.ok().build();
    }
}
