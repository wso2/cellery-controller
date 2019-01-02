/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

/**
 * Common utilities for the API.
 */
public class Utils {

    /**
     * Generate a Siddhi match condition to match any value from a array of values for a particular attribute.
     *
     * Eg:-
     *     Input  - traceId, ["id01", "id02", "id03"]
     *     Output - traceId == "id01" or traceId == "id02" or traceId == "id03"
     *
     * @param attributeName The name of the attribute
     * @param values The array of values from which at least one should match
     * @return The match condition which would match any value from the provided array
     */
    public static String generateSiddhiMatchConditionForAnyValues(String attributeName, String[] values) {
        StringBuilder traceIdMatchConditionBuilder = new StringBuilder();
        for (int i = 0; i < values.length; i++) {
            if (i != 0) {
                traceIdMatchConditionBuilder.append(" or ");
            }
            traceIdMatchConditionBuilder.append(attributeName)
                    .append(" == \"")
                    .append(values[i])
                    .append("\"");
        }
        return traceIdMatchConditionBuilder.toString();
    }
}
