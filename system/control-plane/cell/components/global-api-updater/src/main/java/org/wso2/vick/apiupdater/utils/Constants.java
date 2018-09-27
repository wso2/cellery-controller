/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.vick.apiupdater.utils;

/**
 * This class represents the constants.
 */
public class Constants {

    /**
     * Json param names Constants.
     */
    public static class JsonParamNames {

        public static final String CONTEXT = "context";
        public static final String DEFINITION = "definitions";
        public static final String BACKEND = "backend";
        public static final String GLOBAL = "global";
        public static final String PATH = "path";
        public static final String METHOD = "method";
        public static final String CELL = "cell";
        public static final String VERSION = "version";
        public static final String APIS = "apis";
        public static final String USERNAME = "username";
        public static final String PASSWORD = "password";
        public static final String REGISTER_OBJECT = "registerObj";
        public static final String APIM_BASE_URL = "apimBaseUrl";
        public static final String REGISTER_PATH = "registerPath";
        public static final String TOKEN_ENDPOINT = "tokenEndpoint";
        public static final String API_CREATE_PATH = "apiCreatePath";
        public static final String API_PUBLISH_PATH = "apiPublishPath";
        public static final String ADD_LABEL_PATH = "addLabelPath";
        public static final String NAME = "name";
        public static final String IS_DEFAULT_VERSION = "isDefaultVersion";
        public static final String TRANSPORT = "transport";
        public static final String TIERS = "tiers";
        public static final String GATEWAY_ENVIRONMENTS = "gatewayEnvironments";
        public static final String VISIBILITY = "visibility";
        public static final String LABELS = "labels";
        public static final String API_DEFINITION = "apiDefinition";
        public static final String ENDPOINT_CONFIG = "endpointConfig";
        public static final String PRODUCTION_ENDPOINTS = "production_endpoints";
        public static final String ENDPOINT_TYPE = "endpoint_type";
        public static final String ACCESS_URLS = "accessUrls";
        public static final String PARAMETERS = "parameters";
        public static final String X_AUTH_TYPE = "x-auth-type";
        public static final String REQUIRED = "required";
        public static final String IN = "in";
        public static final String GET = "get";
        public static final String POST = "post";
        public static final String PUT = "put";
        public static final String DELETE = "delete";
        public static final String PATHS = "paths";
        public static final String URL = "url";
    }

    /**
     * Json param names Constants.
     */
    public static class Utils {
        // Util constants
        public static final String CONTENT_TYPE_APPLICATION_JSON = "application/json";
        public static final String CONTENT_TYPE_APPLICATION_URL_ENCODED = "application/x-www-form-urlencoded";
        public static final String PRODUCTION_AND_SANDBOX = "Production and Sandbox";
        public static final String EMPTY_STRING = "";

        // Rest API constants
        public static final String HTTP_RESPONSE_TYPE_ACCEPT = "Accept";
        public static final String HTTP_CONTENT_TYPE = "Content-type";
        public static final String HTTP_REQ_HEADER_AUTHZ = "Authorization";
        public static final String BEARER = "Bearer ";
        public static final String BASIC = "Basic ";
        public static final String BODY = "body";
        public static final String ALLOW_QUERY_PATTERN = "/*";

        public static final String TOKEN_PAYLOAD =
                "grant_type=password&username=$USER&password=$PASS&scope=apim:api_create apim:api_publish apim:label_manage";

        // Config map file paths
        public static final String CELL_CONFIGURATION_FILE_PATH = "/etc/config/api.json";
        public static final String REST_CONFIGURATION_FILE_PATH = "/etc/config/gw.json";

        // Token constants
        public static final String CLIENT_ID = "clientId";
        public static final String CLIENT_SECRET = "clientSecret";
        public static final String ACCESS_TOKEN = "access_token";
        public static final String ID = "id";

        // Error constants
        public static final String DUPLICATE_LABEL_ERROR = "Error while adding new Label for";
        public static final String DUPLICATE_API_ERROR = "A duplicate API already exists for";
        public static final String DIFFERENT_CONTEXT_ERROR = "already exists with different context";
    }
}
