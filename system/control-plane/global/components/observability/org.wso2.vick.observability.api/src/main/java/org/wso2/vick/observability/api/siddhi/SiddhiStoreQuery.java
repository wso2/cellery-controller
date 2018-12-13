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

import org.wso2.siddhi.core.event.Event;
import org.wso2.vick.observability.api.internal.ServiceHolder;

/**
 * Executable Siddhi Store Query.
 *
 * This can be created using the @link {@link SiddhiStoreQueryTemplates}
 */
public class SiddhiStoreQuery {

    private String query;

    private SiddhiStoreQuery(String query) {
        this.query = query;
    }

    /**
     * Execute the Siddhi Store query and get the results as Json Array of Json Objects.
     *
     * @return Siddhi Store Query Results
     */
    public Object[][] execute() {
        Event[] queryResults = ServiceHolder.getSiddhiStoreQueryManager().query(query);

        Object[][] results;
        if (queryResults != null) {
            int rowCount = queryResults.length;
            int columnCount = queryResults[0].getData().length;

            results = new Object[rowCount][columnCount];
            for (int i = 0; i < rowCount; i++) {
                Object[] resultRow = new Object[columnCount];
                for (int j = 0; j < columnCount; j++) {
                    resultRow[j] = queryResults[i].getData(j);
                }
                results[i] = resultRow;
            }
        } else {    // No matching results
            results = new Object[0][0];
        }
        return results;
    }

    /**
     * Siddhi Store Query Builder for building a query string.
     * This supports replacing values in the query.
     * Method chaining can be used with this builder.
     */
    public static class Builder {

        private String query;

        Builder(String query) {
            this.query = query;
        }

        /**
         * Replace a parameter in the Siddhi Store Query.
         *
         * @param key   The name of the parameter to replace
         * @param value The value by which the parameter should be replaced with
         * @return The Siddhi Store Query Builder for chaining
         */
        public Builder setArg(String key, Object value) {
            this.query = this.query.replaceAll("\\$\\{" + key + "}", value.toString());
            return this;
        }

        /**
         * Build the Siddhi Store Query string from this builder.
         *
         * @return The Siddhi Store Query
         */
        public SiddhiStoreQuery build() {
            return new SiddhiStoreQuery(query);
        }
    }
}
