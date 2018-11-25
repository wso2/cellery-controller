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
 * This class holds the information that is required to regenerate the graph links from the cache elements.
 */
public class SpanCacheInfo {
    private String spanId;
    private NodeInfo client;
    private NodeInfo server;

    SpanCacheInfo(String spanId, NodeInfo nodeInfo, Type type) {
        this.spanId = spanId;
        this.setNodeInfo(nodeInfo, type);
    }

    public void setNodeInfo(NodeInfo nodeInfo, Type type) {
        if (type == Type.CLIENT) {
            this.client = nodeInfo;
        } else {
            this.server = nodeInfo;
        }
    }

    public String getSpanId() {
        return spanId;
    }

    public NodeInfo getClient() {
        return client;
    }

    public NodeInfo getServer() {
        return server;
    }

    /**
     * This class wraps the Node information that participates on the span.
     *
     */
    public static class NodeInfo {
        private Node node;
        private String service;
        private String operationName;

        NodeInfo(Node node, String service, String operationName) {
            this.node = node;
            this.service = service;
            this.operationName = operationName;
        }

        public Node getNode() {
            return node;
        }

        public String getService() {
            return service;
        }

        public String getOperationName() {
            return operationName;
        }
    }

    /**
     * This enum defines the types of the Node within a span.
     */
    public enum Type {
        CLIENT, SERVER
    }

}
