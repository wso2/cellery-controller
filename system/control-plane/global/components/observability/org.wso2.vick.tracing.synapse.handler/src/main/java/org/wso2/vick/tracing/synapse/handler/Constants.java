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
    public static final String ZIPKIN_HOST = "ZIPKIN_HOST";
    public static final String ZIPKIN_PORT = "ZIPKIN_PORT";
    public static final String ZIPKIN_API_CONTEXT = "ZIPKIN_API_CONTEXT";

    public static final String TRACING_CORRELATION_ID = "TRACING_CORRELATION_ID";
    public static final String GLOBAL_GATEWAY_SERVICE_NAME = "global-gateway";
    public static final boolean TRACING_SENDER_COMPRESSION_ENABLED = false;
    public static final String B3_GLOBAL_GATEWAY_CORRELATION_ID_HEADER = "X-B3-GlobalGatewayCorrelationID";

    // Tag keys
    public static final String TAG_KEY_HTTP_METHOD = "http.method";
    public static final String TAG_KEY_HTTP_URL = "http.url";
    public static final String TAG_KEY_PROTOCOL = "protocol";
    public static final String TAG_KEY_PEER_ADDRESS = "peer.address";
    public static final String TAG_KEY_SPAN_KIND = "span.kind";
    public static final String TAG_KEY_HTTP_STATUS_CODE = "http.status_code";

    // Tag values
    public static final String SPAN_KIND_CLIENT = "client";
    public static final String SPAN_KIND_SERVER = "server";

    // Synapse message context properties
    public static final String SYNAPSE_MESSAGE_CONTEXT_PROPERTY_HTTP_METHOD = "REST_METHOD";
    public static final String SYNAPSE_MESSAGE_CONTEXT_PROPERTY_ENDPOINT = "ENDPOINT_ADDRESS";
    public static final String SYNAPSE_MESSAGE_CONTEXT_PROPERTY_PEER_ADDRESS = "api.ut.hostName";
    public static final String SYNAPSE_MESSAGE_CONTEXT_PROPERTY_TRANSPORT = "TRANSPORT_IN_NAME";

    // Axis2 message context properties
    public static final String AXIS2_MESSAGE_CONTEXT_PROPERTY_HTTP_STATUS_CODE = "HTTP_SC";
    public static final String AXIS2_MESSAGE_CONTEXT_PROPERTY_HTTP_METHOD = "HTTP_METHOD";
    public static final String AXIS2_MESSAGE_CONTEXT_PROPERTY_REMOTE_HOST = "REMOTE_HOST";
}
