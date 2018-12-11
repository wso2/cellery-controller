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

import com.google.gson.JsonElement;
import com.google.gson.JsonNull;
import com.google.gson.JsonPrimitive;

import java.util.ArrayList;
import java.util.List;

/**
 * Siddhi Store Query Templates Enum class containing all the Siddhi Store Queries.
 * The Siddhi Store Query Builder can be accessed from the Siddhi Store Query Templates.
 */
public enum SiddhiStoreQueryTemplates {

    /*
     * Siddhi Store Queries Start Here
     */

    CELL_LEVEL_REQUEST_AGGREGATION("from RequestAggregation\n" +
            "within ${" + Params.QUERY_START_TIME + "}, ${" + Params.QUERY_END_TIME + "}" +
            "per ${" + Params.TIME_GRANULARITY + "}" +
            "select sourceCell, destinationCell, " +
                "sum(avgResponseTime * requestCount) / sum(requestCount) as avgResponseTime, " +
                "sum(requestCount) as requestCount\n" +
            "group by sourceCell, destinationCell",
            "sourceCell string, destinationCell string, avgResponseTime long, requestCount long\n"
    );

    /*
     * Siddhi Store Queries End Here
     */

    private String query;
    private List<Attribute> selectAttributes;

    SiddhiStoreQueryTemplates(String query, String selectStatement) {
        this.query = query;

        this.selectAttributes = new ArrayList<>();
        String[] attributesWithType = selectStatement.trim().split("\\s*,\\s*");
        for (String attributeWithType : attributesWithType) {
            String[] attributeAndType = attributeWithType.split("\\s+");
            this.selectAttributes.add(new Attribute(attributeAndType[0], attributeAndType[1]));
        }
    }

    /**
     * Parameters to be replaced in the Siddhi Queries.
     */
    public static class Params {
        public static final String QUERY_START_TIME = "queryStartTime";
        public static final String QUERY_END_TIME = "queryEndTime";
        public static final String TIME_GRANULARITY = "timeGranularity";
    }

    /**
     * Class for storing Siddhi Store Query Template Attribute.
     * This also helps in parsing the value and getting a JsonElement for the value.
     */
    public static class Attribute {
        private String name;
        private String type;

        Attribute(String name, String type) {
            this.name = name;
            this.type = type;
        }

        /**
         * Get the name of this attribute.
         *
         * @return The name of the attribute
         */
        String getName() {
            return name;
        }

        /**
         * Parse the Siddhi Attribute value and get a Json Element.
         *
         * @param value The object to parse
         * @return The Json Element for the attribute value
         */
        JsonElement parseValue(Object value) {
            JsonElement jsonElement;
            switch (this.type.toUpperCase()) {
                case "STRING":
                    jsonElement = new JsonPrimitive((String) value);
                    break;
                case "INT":
                    jsonElement = new JsonPrimitive((int) value);
                    break;
                case "LONG":
                    jsonElement = new JsonPrimitive((long) value);
                    break;
                case "DOUBLE":
                    jsonElement = new JsonPrimitive((double) value);
                    break;
                case "FLOAT":
                    jsonElement = new JsonPrimitive((float) value);
                    break;
                case "BOOL":
                    jsonElement = new JsonPrimitive((boolean) value);
                    break;
                default:
                    jsonElement = JsonNull.INSTANCE;
            }
            return jsonElement;
        }
    }

    /**
     * Get the build for the query.
     *
     * @return The Siddhi Store Query Builder for the particular query
     */
    public SiddhiStoreQuery.Builder builder() {
        return new SiddhiStoreQuery.Builder(query, selectAttributes);
    }
}
