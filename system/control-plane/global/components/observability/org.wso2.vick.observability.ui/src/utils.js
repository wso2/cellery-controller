/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * Common utilities.
 */
class Utils {


    /**
     * Parse the query param string and get an object of key value pairs.
     *
     * @param {string} queryParamString Query param string
     * @returns {Object} Query param object
     */
    static parseQueryParams(queryParamString) {
        const queryParameters = {};
        if (queryParameters) {
            let query = queryParamString;
            if (queryParamString.startsWith("?")) {
                query = queryParamString.substr(1);
            }

            const queries = query.split("&");
            for (let i = 0; i < queries.length; i++) {
                const queryPair = queries[i].split("=");
                const key = decodeURIComponent(queryPair[0]);
                queryParameters[key] = decodeURIComponent(queryPair[1]);
            }
        }
        return queryParameters;
    }

}

export default Utils;
