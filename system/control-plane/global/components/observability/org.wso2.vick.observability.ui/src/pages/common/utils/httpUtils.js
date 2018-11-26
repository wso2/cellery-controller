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

import AuthUtils from "./authUtils";
import {ConfigConstants} from "../config";
import axios from "axios";

class HttpUtils {

    /**
     * Parse the query param string and get an object of key value pairs.
     *
     * @param {string} queryParamString Query param string
     * @returns {Object} Query param object
     */
    static parseQueryParams(queryParamString) {
        const queryParameters = {};
        if (queryParamString) {
            let query = queryParamString;
            if (queryParamString.startsWith("?")) {
                query = queryParamString.substr(1);
            }

            if (query) {
                const queries = query.split("&");
                for (let i = 0; i < queries.length; i++) {
                    const queryPair = queries[i].split("=");
                    const key = decodeURIComponent(queryPair[0]);

                    if (key) {
                        queryParameters[key] = (queryPair.length === 2 && queryPair[1])
                            ? decodeURIComponent(queryPair[1])
                            : true;
                    }
                }
            }
        }
        return queryParameters;
    }

    /**
     * Call the Siddhi backend API.
     *
     * @param {Object} config Axios configuration object
     * @param {ConfigHolder} globalConfig The global configuration provided to the current component
     * @returns {Promise} A promise for the API call
     */
    static callBackendAPI(config, globalConfig) {
        return new Promise((resolve, reject) => {
            if (!config.headers) {
                config.headers = {};
            }
            if (!config.headers.Accept) {
                config.headers.Accept = "application/json";
            }
            if (!config.headers["Content-Type"]) {
                config.headers["Content-Type"] = "application/json";
            }
            if (!config.data && (config.method === "POST" || config.method === "PUT" || config.method === "PATCH")) {
                config.data = {};
            }
            config.url = `${globalConfig.get(ConfigConstants.BACKEND_URL)}${config.url}`;

            axios(config)
                .then((response) => {
                    if (response.status >= 200 && response.status < 400) {
                        if (response.data.map) {
                            resolve(response.data.map((dataItem) => dataItem.event));
                        } else {
                            resolve(response.data);
                        }
                    } else {
                        reject(response.data);
                    }
                })
                .catch((error) => {
                    if (error.response) {
                        const errorResponse = error.response;
                        if (errorResponse.status === 401) {
                            // Redirect to home page since the user is not authorised
                            AuthUtils.signOut(globalConfig);
                        }
                        reject(new Error(errorResponse.data));
                    } else {
                        reject(error);
                    }
                });
        });
    }

}

export default HttpUtils;
