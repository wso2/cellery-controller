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

package org.wso2.vick.tracing.receiver;

/**
 * Tracing Receiver Constants.
 */
public class Constants {
    public static final String TRACE_ID = "traceId";
    public static final String SPAN_ID = "id";
    public static final String PARENT_ID = "parentId";
    public static final String NAME = "name";
    public static final String SERVICE_NAME = "serviceName";
    public static final String KIND = "kind";
    public static final String TIMESTAMP = "timestamp";
    public static final String DURATION = "duration";
    public static final String TAGS = "tags";

    public static final String HTTP_CONTENT_TYPE_HEADER = "Content-Type";
    public static final String HTTP_APPLICATION_THRIFT_CONTENT_TYPE = "application/x-thrift";
}
