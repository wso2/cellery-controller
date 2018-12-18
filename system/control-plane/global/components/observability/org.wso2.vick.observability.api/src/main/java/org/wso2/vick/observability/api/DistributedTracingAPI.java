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

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.gson.JsonSyntaxException;
import org.wso2.vick.observability.api.exception.InvalidParamException;
import org.wso2.vick.observability.api.siddhi.SiddhiStoreQueryTemplates;

import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
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
 * MSF4J service for fetching distributed tracing data.
 */
@Path("/api/traces")
public class DistributedTracingAPI {

    @GET
    @Path("/metadata")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getMetadata(@DefaultValue("-1") @QueryParam("queryStartTime") long queryStartTime,
                                @DefaultValue("-1") @QueryParam("queryEndTime") long queryEndTime) {
        Object[][] results = SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_METADATA.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
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
    public Response search(@DefaultValue("") @QueryParam("cell") String cell,
                           @DefaultValue("") @QueryParam("serviceName") String serviceName,
                           @DefaultValue("") @QueryParam("operationName") String operationName,
                           @DefaultValue("-1") @QueryParam("minDuration") long minDuration,
                           @DefaultValue("-1") @QueryParam("maxDuration") long maxDuration,
                           @DefaultValue("-1") @QueryParam("queryStartTime") long queryStartTime,
                           @DefaultValue("-1") @QueryParam("queryEndTime") long queryEndTime,
                           @DefaultValue("{}") @QueryParam("tags") String jsonEncodedTags) {
        Map<String, String> queryTags = new HashMap<>();
        SiddhiStoreQueryTemplates siddhiStoreQueryTemplates;
        if ("{}".equals(jsonEncodedTags)) {
            siddhiStoreQueryTemplates = SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_SEARCH_GET_TRACE_IDS;
        } else {
            siddhiStoreQueryTemplates = SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_SEARCH_GET_TRACE_IDS_WITH_TAGS;

            // Parsing the provided JSON encoded tags
            try {
                JsonElement jsonElement = new JsonParser().parse(jsonEncodedTags);
                if (jsonElement.isJsonObject()) {
                    JsonObject queryTagsJsonObject = jsonElement.getAsJsonObject();
                    for (Map.Entry<String, JsonElement> queryTagsJsonObjectEntry : queryTagsJsonObject.entrySet()) {
                        if (queryTagsJsonObjectEntry.getValue().isJsonPrimitive()) {
                            queryTags.put(
                                    queryTagsJsonObjectEntry.getKey(),
                                    queryTagsJsonObjectEntry.getValue().getAsString()
                            );
                        } else {
                            throw new InvalidParamException("tags", "JSON encoded object");
                        }
                    }
                } else {
                    throw new InvalidParamException("tags", "JSON encoded object");
                }
            } catch (JsonSyntaxException e) {
                throw new InvalidParamException("tags", "JSON encoded object", e);
            }
        }

        Object[][] traceIdResults = siddhiStoreQueryTemplates.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.CELL, cell)
                .setArg(SiddhiStoreQueryTemplates.Params.SERVICE_NAME, serviceName)
                .setArg(SiddhiStoreQueryTemplates.Params.OPERATION_NAME, operationName)
                .setArg(SiddhiStoreQueryTemplates.Params.MIN_DURATION, minDuration)
                .setArg(SiddhiStoreQueryTemplates.Params.MAX_DURATION, maxDuration)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_START_TIME, queryStartTime)
                .setArg(SiddhiStoreQueryTemplates.Params.QUERY_END_TIME, queryEndTime)
                .setArg(SiddhiStoreQueryTemplates.Params.TAGS_JSON_ENCODED, jsonEncodedTags)
                .build()
                .execute();

        // Filtering based on tags
        Set<String> traceIds = new HashSet<>();
        for (Object[] traceIdResult : traceIdResults) {
            String traceId = (String) traceIdResult[0];

            if ("{}".equals(jsonEncodedTags)) {
                traceIds.add(traceId);
            } else if (!traceIds.contains(traceId)) {   // To consider a traceId a single matching span is enough
                boolean isMatch = false;
                JsonElement parsedJsonElement = new JsonParser().parse((String) traceIdResult[1]);
                if (parsedJsonElement.isJsonObject()) {
                    JsonObject traceTags = parsedJsonElement.getAsJsonObject();
                    for (Map.Entry<String, String> queryTagEntry : queryTags.entrySet()) {
                        String tagKey = queryTagEntry.getKey();
                        String tagValue = queryTagEntry.getValue();

                        JsonElement traceTagValueJsonElement = traceTags.get(tagKey);
                        if (traceTagValueJsonElement != null && traceTagValueJsonElement.isJsonPrimitive()
                                && tagValue.equals(traceTagValueJsonElement.getAsString())) {
                            isMatch = true;
                            break;
                        }
                    }
                }
                if (isMatch) {
                    traceIds.add(traceId);
                }
            }
        }

        Object[][] spanCountResults;
        Object[][] rootSpanResults;
        if (traceIds.size() > 0) {
            // Calculating a match condition for the traceIds
            StringBuilder traceIdMatchConditionBuilder = new StringBuilder();
            String[] traceIdArray = traceIds.toArray(new String [0]);
            for (int i = 0; i < traceIdArray.length; i++) {
                if (i != 0) {
                    traceIdMatchConditionBuilder.append(" or ");
                }
                traceIdMatchConditionBuilder.append("traceId == \"")
                        .append(traceIdArray[i])
                        .append("\"");
            }
            String traceIdMatchCondition = traceIdMatchConditionBuilder.toString();

            spanCountResults =
                    SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_SEARCH_GET_MULTIPLE_CELL_SERVICE_COUNTS.builder()
                            .setArg(SiddhiStoreQueryTemplates.Params.CONDITION, traceIdMatchCondition)
                            .build()
                            .execute();

            rootSpanResults = SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_SEARCH_GET_ROOT_SPAN_METADATA
                    .builder()
                    .setArg(SiddhiStoreQueryTemplates.Params.CONDITION, traceIdMatchCondition)
                    .build()
                    .execute();
        } else {
            spanCountResults = new Object[0][0];
            rootSpanResults = new Object[0][0];
        }

        Map<String, Object[][]> resultsMap = new HashMap<>(2);
        resultsMap.put("spanCounts", spanCountResults);
        resultsMap.put("rootSpans", rootSpanResults);

        return Response.ok().entity(resultsMap).build();
    }

    @OPTIONS
    @Path("/search")
    public Response searchOptions() {
        return Response.ok().build();
    }


    @GET
    @Path("/{traceId}")
    @Produces(MediaType.APPLICATION_JSON)
    public Response getTraceByTraceId(@PathParam("traceId") String traceId) {
        Object[][] results = SiddhiStoreQueryTemplates.DISTRIBUTED_TRACING_GET_TRACE.builder()
                .setArg(SiddhiStoreQueryTemplates.Params.TRACE_ID, traceId)
                .build()
                .execute();
        return Response.ok().entity(results).build();
    }

    @OPTIONS
    @Path("/{traceId}")
    public Response getTraceByTraceId() {
        return Response.ok().build();
    }
}
