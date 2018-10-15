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

package org.wso2.vick.tracing.synapse.handler;

/**
 * Synapse handler related Constants.
 */
public class Constants {
    public static final String ZIPKIN_HOST = "zipkinHost";
    public static final String ZIPKIN_PORT = "zipkinPort";
    public static final String ZIPKIN_API_CONTEXT = "zipkinAPIContext";

    public static final String REQUEST_IN_SPAN = "VICK_TRACING_REQUEST_IN_SPAN";
    public static final String REQUEST_OUT_SPAN = "VICK_TRACING_REQUEST_OUT_SPAN";
    public static final String GLOBAL_GATEWAY_SERVICE_NAME = "global-gateway";
}
