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

import AuthUtils from "./authUtils";
import {StateHolder} from "../state";
import axios from "axios";

class HttpUtils {

    /**
     * Parse the query param string and get an object of key value pairs.
     *
     * @param {string} queryParamString Query param string
     * @returns {Object} Query param object
     */
    static parseQueryParams = (queryParamString) => {
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
    };

    /**
     * Generate a query param string from a query params object.
     *
     * @param {Object} queryParams Query params as an flat object
     * @returns {string} Query string
     */
    static generateQueryParamString = (queryParams) => {
        let queryString = "";
        if (queryParams) {
            for (const queryParamKey in queryParams) {
                if (queryParams.hasOwnProperty(queryParamKey)) {
                    const queryParamValue = queryParams[queryParamKey];

                    if (!queryParamValue) {
                        continue;
                    }

                    // Validating
                    if (typeof queryParamKey !== "string") {
                        throw Error(`Query param key need to be a string, instead found ${typeof queryParamKey}`);
                    }
                    if (typeof queryParamValue !== "string" && typeof queryParamValue !== "number"
                            && typeof queryParamValue !== "boolean") {
                        throw Error(`Query param value need to be a string, instead found ${typeof queryParamValue}`);
                    }

                    // Generating query string
                    queryString += queryString ? "&" : "?";
                    queryString += `${encodeURIComponent(queryParamKey)}=${encodeURIComponent(queryParamValue)}`;
                }
            }
        }
        return queryString;
    };

    /**
     * Call a deployed Siddhi App to fetch results.
     *
     * @param {Object} config Axios configuration object
     * @param {StateHolder} globalState The global state provided to the current component
     * @returns {Promise} A promise for the API call
     */
    static callObservabilityAPI = (config, globalState) => {
        config.url = `${globalState.get(StateHolder.CONFIG).observabilityAPIURL}${config.url}`;
        return HttpUtils.callAPI(config, globalState);
    };

    /**
     * Call the Siddhi backend API.
     *
     * @param {Object} config Axios configuration object
     * @param {StateHolder} globalState The global state provided to the current component
     * @returns {Promise} A promise for the API call
     */
    static callAPI = (config, globalState) => new Promise((resolve, reject) => {
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

        axios(config)
            .then((response) => {
                if (response.status >= 200 && response.status < 400) {
                    resolve(response.data);
                } else {
                    reject(response.data);
                }
            })
            .catch((error) => {
                if (error.response) {
                    const errorResponse = error.response;
                    if (errorResponse.status === 401) {
                        // Redirect to home page since the user is not authorised
                        AuthUtils.signOut(globalState);
                    }
                    reject(new Error(errorResponse.data));
                } else {
                    reject(error);
                }
            });
    });

}

export default HttpUtils;
