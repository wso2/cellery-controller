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

package org.wso2.vick.observability.api.siddhi;

/**
 * Siddhi Store Query Templates Enum class containing all the Siddhi Store Queries.
 * The Siddhi Store Query Builder can be accessed from the Siddhi Store Query Templates.
 */
public enum SiddhiStoreQueryTemplates {

    /*
     * Siddhi Store Queries Start Here
     */

    REQUEST_AGGREGATION_CELLS("from RequestAggregation\n" +
            "within ${" + Params.QUERY_START_TIME + "}, ${" + Params.QUERY_END_TIME + "}\n" +
            "per \"${" + Params.TIME_GRANULARITY + "}\"\n" +
            "select sourceCell, destinationCell, " +
            "sum(avgResponseTime * requestCount) / sum(requestCount) as avgResponseTime, " +
            "sum(requestCount) as requestCount\n" +
            "group by sourceCell, destinationCell"
    ),
    DISTRIBUTED_TRACING_METADATA("from DistributedTracingTable\n" +
            "on (${" + Params.QUERY_START_TIME + "}L == -1 or startTime >= ${" + Params.QUERY_START_TIME + "}L) " +
            "and (${" + Params.QUERY_END_TIME + "}L == -1 or startTime <= ${" + Params.QUERY_END_TIME + "}L)\n" +
            "select cell, serviceName, operationName\n" +
            "group by cell, serviceName, operationName"
    ),
    DISTRIBUTED_TRACING_SEARCH_GET_TRACE_IDS("from DistributedTracingTable\n" +
            "on (\"${" + Params.CELL + "}\" == \"\" or cell == \"${" + Params.CELL + "}\") " +
            "and (\"${" + Params.SERVICE_NAME + "}\" == \"\" or " +
            "serviceName == \"${" + Params.SERVICE_NAME + "}\") " +
            "and (\"${" + Params.OPERATION_NAME + "}\" == \"\" or " +
            "operationName == \"${" + Params.OPERATION_NAME + "}\") " +
            "and (${" + Params.MIN_DURATION + "}L == -1 or duration >= ${" + Params.MIN_DURATION + "}L) " +
            "and (${" + Params.MAX_DURATION + "}L == -1 or duration <= ${" + Params.MAX_DURATION + "}L) " +
            "and (${" + Params.QUERY_START_TIME + "}L == -1 or startTime >= ${" + Params.QUERY_START_TIME + "}L) " +
            "and (${" + Params.QUERY_END_TIME + "}L == -1 or startTime <= ${" + Params.QUERY_END_TIME + "}L)\n" +
            "select traceId\n" +
            "group by traceId"
    ),
    DISTRIBUTED_TRACING_SEARCH_GET_TRACE_IDS_WITH_TAGS("from DistributedTracingTable\n" +
            "on (\"${" + Params.CELL + "}\" == \"\" or cell == \"${" + Params.CELL + "}\") " +
            "and (\"${" + Params.SERVICE_NAME + "}\" == \"\" or " +
            "serviceName == \"${" + Params.SERVICE_NAME + "}\") " +
            "and (\"${" + Params.OPERATION_NAME + "}\" == \"\" or " +
            "operationName == \"${" + Params.OPERATION_NAME + "}\") " +
            "and (${" + Params.MIN_DURATION + "}L == -1 or duration >= ${" + Params.MIN_DURATION + "}L) " +
            "and (${" + Params.MAX_DURATION + "}L == -1 or duration <= ${" + Params.MAX_DURATION + "}L) " +
            "and (${" + Params.QUERY_START_TIME + "}L == -1 or startTime >= ${" + Params.QUERY_START_TIME + "}L) " +
            "and (${" + Params.QUERY_END_TIME + "}L == -1 or startTime <= ${" + Params.QUERY_END_TIME + "}L)\n" +
            "select traceId, tags"
    ),
    DISTRIBUTED_TRACING_SEARCH_GET_MULTIPLE_CELL_SERVICE_COUNTS("from DistributedTracingTable\n" +
            "on ${" + Params.CONDITION + "}\n" +
            "select traceId, cell, serviceName, count() as count\n" +
            "group by traceId, cell, serviceName\n" +
            "order by startTime desc"
    ),
    DISTRIBUTED_TRACING_SEARCH_GET_ROOT_SPAN_METADATA("from DistributedTracingTable\n" +
            "on traceId == spanId and (${" + Params.CONDITION + "})\n" +
            "select traceId, serviceName, operationName, startTime, duration"
    ),
    DISTRIBUTED_TRACING_GET_TRACE("from DistributedTracingTable\n" +
            "on traceId == \"${" + Params.TRACE_ID + "}\"\n" +
            "select traceId, spanId, parentId, namespace, cell, serviceName, pod, operationName, kind, startTime, " +
            "duration, tags"
    );

    /*
     * Siddhi Store Queries End Here
     */

    private String query;

    SiddhiStoreQueryTemplates(String query) {
        this.query = query;
    }

    /**
     * Parameters to be replaced in the Siddhi Queries.
     */
    public static class Params {
        public static final String QUERY_START_TIME = "queryStartTime";
        public static final String QUERY_END_TIME = "queryEndTime";
        public static final String TIME_GRANULARITY = "timeGranularity";
        public static final String CELL = "cell";
        public static final String SERVICE_NAME = "serviceName";
        public static final String OPERATION_NAME = "operationName";
        public static final String MIN_DURATION = "minDuration";
        public static final String MAX_DURATION = "maxDuration";
        public static final String TRACE_ID = "traceId";
        public static final String TAGS_JSON_ENCODED = "tagsJsonEncoded";
        public static final String CONDITION = "condition";     // Should be used with caution considering SQL injection
    }

    /**
     * Get the build for the query.
     *
     * @return The Siddhi Store Query Builder for the particular query
     */
    public SiddhiStoreQuery.Builder builder() {
        return new SiddhiStoreQuery.Builder(query);
    }
}
