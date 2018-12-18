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

import javax.ws.rs.DefaultValue;
import javax.ws.rs.GET;
import javax.ws.rs.OPTIONS;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

/**
 * MSF4J service for fetching the aggregated request.
 */
@Path("/api/http-requests")
public class AggregatedRequestsAPI {

    @GET
    @Path("/cells")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getForCells(@QueryParam("queryStartTime") long queryStartTime,
                                @QueryParam("queryEndTime") long queryEndTime,
                                @DefaultValue("seconds") @QueryParam("timeGranularity") String timeGranularity) {
        Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_CELLS.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .setArg(SiddhiStoreQueryTemplates.Params.TIME_GRANULARITY, timeGranularity)
                .build()
                .execute();
        return Response.ok().entity(results).build();
    }

    @OPTIONS
    @Path("/cells")
    public Response getCellOverviewOptions() {
        return Response.ok().build();
    }
}
