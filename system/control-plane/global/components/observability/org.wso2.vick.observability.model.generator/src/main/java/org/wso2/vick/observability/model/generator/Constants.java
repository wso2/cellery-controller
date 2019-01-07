/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.observability.model.generator;

/**
 * This class holds the constants that are required for the model generator.
 */
public class Constants {
    public static final String SERVER_SPAN_KIND = "SERVER";
    public static final String CLIENT_SPAN_KIND = "CLIENT";
    public static final String EDGE_NAME_CONNECTOR = " ---> ";
    public static final String LINK_SEPARATOR = "##";
    public static final String CELL_SERVICE_NAME_SEPARATOR = ":";
    public static final String IGNORE_OPERATION_NAME = "async ext_authz egress";
}
