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

import org.wso2.siddhi.core.SiddhiAppRuntime;
import org.wso2.siddhi.core.SiddhiManager;
import org.wso2.siddhi.core.event.Event;

/**
 * Manager for running Siddhi Store Queries.
 * This doesn't need to be accessed directly except for starting and stopping the service.
 * For executing use {@link SiddhiStoreQueryTemplates}.
 */
public class SiddhiStoreQueryManager {
    private static final String DISTRIBUTED_TRACING_TABLE_DEFINITION = "@Store(type=\"rdbms\", " +
                "datasource=\"VICK_OBSERVABILITY_DB\", field.length=\"tags:8000\")\n" +
            "@PrimaryKey(\"traceId\", \"spanId\", \"kind\")\n" +
            "define table DistributedTracingTable (traceId string, spanId string, parentId string, namespace string, " +
                "cell string, serviceName string, pod string, operationName string, kind string, startTime long, " +
                "duration long, tags string);";
    private static final String REQUEST_AGGREGATION_DEFINITION = "define stream ProcessedRequestsStream(" +
                "sourceNamespace string, sourceCell string, sourceVICKService string, destinationNamespace string, " +
                "destinationCell string, destinationVICKService string, requestPath string, requestMethod string, " +
                "httpResponseGroup string, responseTimeSec double, responseSizeBytes long);" +
            "@store(type=\"rdbms\", datasource=\"VICK_OBSERVABILITY_DB\")\n" +
            "@purge(enable=\"false\")\n" +
            "define aggregation RequestAggregation from ProcessedRequestsStream\n" +
                "select sourceNamespace, sourceCell, sourceVICKService, destinationNamespace, destinationCell, " +
                "destinationVICKService, requestPath, requestMethod, httpResponseGroup, " +
                "avg(responseTimeSec) as avgResponseTimeSec, avg(responseSizeBytes) as avgResponseSizeBytes, " +
                "count() as requestCount\n" +
            "group by sourceNamespace, sourceCell, sourceVICKService, destinationNamespace, destinationCell, " +
                "destinationVICKService, requestPath, requestMethod, httpResponseGroup\n" +
            "aggregate every sec...year;";
    private static final String SIDDHI_APP = DISTRIBUTED_TRACING_TABLE_DEFINITION + "\n" +
            REQUEST_AGGREGATION_DEFINITION;

    private SiddhiAppRuntime siddhiAppRuntime;

    public SiddhiStoreQueryManager() {
        SiddhiManager siddhiManager = new SiddhiManager();

        siddhiAppRuntime = siddhiManager.createSiddhiAppRuntime(SIDDHI_APP);
        siddhiAppRuntime.start();
    }

    /**
     * Run Siddhi Store Query and get the results.
     *
     * @param siddhiQuery Siddhi Store Query to run
     * @return The results of the Siddhi Store Query
     */
    Event[] query(String siddhiQuery) {
        return siddhiAppRuntime.query(siddhiQuery);
    }

    /**
     * Stop the Siddhi Store Query Manager.
     * This will stop the siddhi App run time to clear any resources allocated.
     */
    public void stop() {
        siddhiAppRuntime.shutdown();
    }
}
