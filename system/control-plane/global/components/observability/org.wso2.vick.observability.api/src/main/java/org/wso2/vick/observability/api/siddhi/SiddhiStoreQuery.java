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

import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import org.wso2.siddhi.core.event.Event;
import org.wso2.vick.observability.api.internal.ServiceHolder;

import java.util.List;

/**
 * Executable Siddhi Store Query.
 *
 * This can be created using the @link {@link SiddhiStoreQueryTemplates}
 */
public class SiddhiStoreQuery {

    private String query;
    private List<SiddhiStoreQueryTemplates.Attribute> selectAttributes;

    private SiddhiStoreQuery(String query, List<SiddhiStoreQueryTemplates.Attribute> selectAttributes) {
        this.query = query;
        this.selectAttributes = selectAttributes;
    }

    /**
     * Execute the Siddhi Store query and get the results as Json Array of Json Objects.
     *
     * @return Siddhi Store Query Results
     */
    public JsonArray execute() {
        Event[] queryResults = ServiceHolder.getSiddhiStoreQueryManager().query(query);

        JsonArray jsonResultArray = new JsonArray();
        for (Event jsonResultRow : queryResults) {
            JsonObject jsonObject = new JsonObject();
            for (int i = 0; i < selectAttributes.size(); i++) {
                SiddhiStoreQueryTemplates.Attribute attribute = selectAttributes.get(i);
                jsonObject.add(attribute.getName(), attribute.parseValue(jsonResultRow.getData(i)));
            }
            jsonResultArray.add(jsonObject);
        }
        return jsonResultArray;
    }

    /**
     * Siddhi Store Query Builder for building a query string.
     * This supports replacing values in the query.
     * Method chaining can be used with this builder.
     */
    public static class Builder {

        private String query;
        private List<SiddhiStoreQueryTemplates.Attribute> selectAttributes;

        Builder(String query, List<SiddhiStoreQueryTemplates.Attribute> selectAttributes) {
            this.query = query;
            this.selectAttributes = selectAttributes;
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
            return new SiddhiStoreQuery(query, selectAttributes);
        }
    }
}
