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

import org.apache.log4j.Logger;
import com.google.gson.JsonObject;
import org.wso2.vick.observability.api.siddhi.SiddhiStoreQueryTemplates;

import java.util.HashSet;
import java.util.Set;
import javax.ws.rs.DefaultValue;
import javax.ws.rs.GET;
import javax.ws.rs.OPTIONS;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

/**
 * MSF4J service for fetching the aggregated request.
 */
@Path("/api/http-requests")
public class AggregatedRequestsAPI {
    private static final Logger log = Logger.getLogger(AggregatedRequestsAPI.class);


    @GET
    @Path("/cells")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getCellStats(@QueryParam("queryStartTime") long queryStartTime,
                                @QueryParam("queryEndTime") long queryEndTime,
                                @DefaultValue("seconds") @QueryParam("timeGranularity") String timeGranularity) {
        try {
            Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_CELLS.builder()
                    .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                    .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                    .setArg(SiddhiStoreQueryTemplates.Params.TIME_GRANULARITY, timeGranularity)
                    .build()
                    .execute();
            return Response.ok().entity(results).build();
        } catch (Throwable throwable) {
            log.error("Unable to get the aggregated results for cells. ", throwable);
            return Response.serverError().entity(throwable).build();
        }
    }

    @GET
    @Path("/cells/{cellName}/services")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getServicesStats(@QueryParam("queryStartTime") long queryStartTime,
                                @QueryParam("queryEndTime") long queryEndTime,
                                @DefaultValue("seconds") @QueryParam("timeGranularity") String timeGranularity,
                                @PathParam("cellName") String cellName) {
        try {
            Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_SERVICES_OF_CELL.builder()
                    .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                    .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                    .setArg(SiddhiStoreQueryTemplates.Params.TIME_GRANULARITY, timeGranularity)
                    .setArg(SiddhiStoreQueryTemplates.Params.CELL, cellName)
                    .build()
                    .execute();
            return Response.ok().entity(results).build();
        } catch (Throwable throwable) {
            log.error("Unable to get the aggregated results for cells. ", throwable);
            return Response.serverError().entity(throwable).build();
        }
    public Response getAggregatedRequestsForCells(@QueryParam("queryStartTime") long queryStartTime,
                                                  @QueryParam("queryEndTime") long queryEndTime,
                                                  @DefaultValue("seconds")
                                                      @QueryParam("timeGranularity") String timeGranularity) {
        Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_CELLS.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .setArg(SiddhiStoreQueryTemplates.Params.TIME_GRANULARITY, timeGranularity)
                .build()
                .execute();
        return Response.ok().entity(results).build();
    }

    @GET
    @Path("/cells/metrics")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getMetricsForCells(@QueryParam("queryStartTime") long queryStartTime,
                                       @QueryParam("queryEndTime") long queryEndTime,
                                       @DefaultValue("") @QueryParam("sourceCell") String sourceCell,
                                       @DefaultValue("") @QueryParam("destinationCell") String destinationCell,
                                       @DefaultValue("seconds") @QueryParam("timeGranularity") String timeGranularity) {
        Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_CELLS_METRICS.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .setArg(SiddhiStoreQueryTemplates.Params.TIME_GRANULARITY, timeGranularity)
                .setArg(SiddhiStoreQueryTemplates.Params.SOURCE_CELL, sourceCell)
                .setArg(SiddhiStoreQueryTemplates.Params.DESTINATION_CELL, destinationCell)
                .build()
                .execute();
        return Response.ok().entity(results).build();
    }

    @GET
    @Path("/cells/metadata")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getMetadataForCells(@QueryParam("queryStartTime") long queryStartTime,
                                        @QueryParam("queryEndTime") long queryEndTime) {
        Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_CELLS_METADATA.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .build()
                .execute();

        Set<String> cells = new HashSet<>();
        for (Object[] result : results) {
            cells.add((String) result[0]);
            cells.add((String) result[1]);
        }

        return Response.ok().entity(cells).build();
    }

    @GET
    @Path("/cells/{cellName}/microservices")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getForMicroservices(@PathParam("cellName") String cellName,
                                        @QueryParam("queryStartTime") long queryStartTime,
                                        @QueryParam("queryEndTime") long queryEndTime,
                                        @DefaultValue("seconds")
                                            @QueryParam("timeGranularity") String timeGranularity) {
        Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_CELL_MICROSERVICES.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .setArg(SiddhiStoreQueryTemplates.Params.TIME_GRANULARITY, timeGranularity)
                .setArg(SiddhiStoreQueryTemplates.Params.CELL, cellName)
                .build()
                .execute();
        return Response.ok().entity(results).build();
    }

    @Path("/*")
    public Response getOptions() {
        return Response.ok().build();
    }

    @GET
    @Path("/cells/microservices/metrics")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getMetricsForMicroservices(@QueryParam("queryStartTime") long queryStartTime,
                                               @QueryParam("queryEndTime") long queryEndTime,
                                               @DefaultValue("") @QueryParam("sourceCell") String sourceCell,
                                               @DefaultValue("")
                                                   @QueryParam("sourceMicroservice") String sourceMicroservice,
                                               @DefaultValue("") @QueryParam("destinationCell") String destinationCell,
                                               @DefaultValue("")
                                                   @QueryParam("destinationMicroservice")String destinationMicroservice,
                                               @DefaultValue("seconds")
                                                   @QueryParam("timeGranularity") String timeGranularity) {
        Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_MICROSERVICES_METRICS.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .setArg(SiddhiStoreQueryTemplates.Params.TIME_GRANULARITY, timeGranularity)
                .setArg(SiddhiStoreQueryTemplates.Params.SOURCE_CELL, sourceCell)
                .setArg(SiddhiStoreQueryTemplates.Params.SOURCE_MICROSERVICE, sourceMicroservice)
                .setArg(SiddhiStoreQueryTemplates.Params.DESTINATION_CELL, destinationCell)
                .setArg(SiddhiStoreQueryTemplates.Params.DESTINATION_MICROSERVICE, destinationMicroservice)
                .build()
                .execute();
        return Response.ok().entity(results).build();
    }

    @GET
    @Path("/cells/microservices/metadata")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getMetadataForMicroservices(@QueryParam("queryStartTime") long queryStartTime,
                                                @QueryParam("queryEndTime") long queryEndTime) {
        Object[][] results = SiddhiStoreQueryTemplates.REQUEST_AGGREGATION_MICROSERVICES_METADATA.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .build()
                .execute();

        Set<JsonObject> microservices = new HashSet<>();
        for (Object[] result : results) {
            for (int i = 0; i < 2; i++) {
                JsonObject microservice = new JsonObject();
                microservice.addProperty("cell", (String) result[i * 2]);
                microservice.addProperty("name", (String) result[i * 2 + 1]);
                microservices.add(microservice);
            }
        }

        return Response.ok().entity(microservices).build();
    }
}
