/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.vick.observability.k8s.api.server.client;

/**
 * This contains the constants for the K8S API Server Client Extensions
 */
public class Constants {
    public static final String NAMESPACE = "default";
    public static final String CELL_NAME_LABEL = "vick.wso2.com/cell";
    public static final String COMPONENT_NAME_LABEL = "vick.wso2.com/service";
    public static final String GATEWAY_NAME_LABEL = "vick.wso2.com/gateway";
    public static final String RUNNING_STATUS_FIELD_SELECTOR = "status.phase=Running";
}
