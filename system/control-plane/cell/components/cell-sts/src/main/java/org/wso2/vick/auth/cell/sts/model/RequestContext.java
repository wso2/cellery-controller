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
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.auth.cell.sts.model;

public class RequestContext {

    private String host;
    private String path;
    private String protocol;
    // HTTP method
    private String method;

    public String getHost() {
        return host;
    }

    public String getPath() {
        return path;
    }

    public String getProtocol() {
        return protocol;
    }

    public RequestContext setHost(String host) {
        this.host = host;
        return this;
    }

    public RequestContext setPath(String path) {
        this.path = path;
        return this;
    }

    public RequestContext setProtocol(String protocol) {
        this.protocol = protocol;
        return this;
    }

    public String getMethod() {
        return method;
    }

    public RequestContext setMethod(String method) {
        this.method = method;
        return this;
    }

    // TODO builder
}
