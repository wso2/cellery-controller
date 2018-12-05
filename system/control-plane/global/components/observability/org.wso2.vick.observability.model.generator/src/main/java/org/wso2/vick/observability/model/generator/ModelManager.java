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

import com.google.common.graph.MutableNetwork;
import com.google.common.graph.NetworkBuilder;
import org.wso2.vick.observability.model.generator.exception.GraphStoreException;
import org.wso2.vick.observability.model.generator.exception.ModelException;
import org.wso2.vick.observability.model.generator.internal.ServiceHolder;

import java.util.Set;

import static org.wso2.vick.observability.model.generator.Constants.EDGE_NAME_CONNECTOR;

/**
 * This is the Manager, singleton class which performs the operations in the in memory dependency tree.
 */
public class ModelManager {
    private MutableNetwork<Node, String> dependencyGraph;

    public ModelManager() throws ModelException {
        try {
            Object[] graph = ServiceHolder.getGraphStoreManager().loadGraph();
            this.dependencyGraph = NetworkBuilder.directed()
                    .allowsParallelEdges(true)
                    .expectedNodeCount(100000)
                    .expectedEdgeCount(1000000)
                    .build();
            if (graph != null) {
                Set<Node> nodes = (Set<Node>) graph[0];
                Set<String> edges = (Set<String>) graph[1];
                addNodes(nodes);
                addEdges(edges);
            }
        } catch (GraphStoreException e) {
            throw new ModelException("Unable to load already persisted model.", e);
        }
    }

    private void addNodes(Set<Node> nodes) {
        for (Node node : nodes) {
            this.dependencyGraph.addNode(node);
        }
    }

    private void addEdges(Set<String> edges) throws ModelException {
        for (String edge : edges) {
            String[] elements = this.edgeNameElements(edge);
            Node parentNode = getNode(elements[0]);
            Node childNode = getNode(elements[1]);
            if (parentNode != null && childNode != null) {
                this.dependencyGraph.addEdge(parentNode, childNode, edge);
            } else {
                String msg = "";
                if (parentNode == null) {
                    msg += "Parent node doesn't exist in the graph for edgename :" + edge + ". ";
                }
                if (childNode == null) {
                    msg += "Client node doesn't exist in the graph for edgename :" + edge + ". ";
                }
                throw new ModelException(msg);
            }
        }
    }

    private Node getNode(String nodeName) {
        Set<Node> nodes = this.dependencyGraph.nodes();
        for (Node node : nodes) {
            if (node.getId().equalsIgnoreCase(nodeName)) {
                return node;
            }
        }
        return null;
    }


    public void addNode(Node node) {
        this.dependencyGraph.addNode(node);
    }

    public void addLink(Node parent, Node child, String serviceName) {
        try {
            this.dependencyGraph.addEdge(parent, child, generateEdgeName(parent.getId(), child.getId(), serviceName));
        } catch (Exception ignored) {
        }
    }

    public Set<Node> getNodes() {
        return this.dependencyGraph.nodes();
    }

    public Set<String> getLinks() {
        return this.dependencyGraph.edges();
    }

    public MutableNetwork<Node, String> getDependencyGraph() {
        return dependencyGraph;
    }

    private String generateEdgeName(String parentNodeId, String childNodeId, String serviceName) {
        return parentNodeId + EDGE_NAME_CONNECTOR + childNodeId + EDGE_NAME_CONNECTOR + serviceName;
    }

    public String[] edgeNameElements(String edgeName) {
        return edgeName.split(EDGE_NAME_CONNECTOR);
    }

    public void moveLinks(Node fromNode, Node targetNode, String newEdgePrefix) {
        Set<String> outEdges = this.dependencyGraph.outEdges(fromNode);
        for (String edgeName : outEdges) {
            Node outNode = this.dependencyGraph.incidentNodes(outEdges).target();
            String newEdgeName = newEdgePrefix + Constants.LINK_SEPARATOR + getEdgePostFix(edgeName);
            this.addLink(targetNode, outNode, newEdgeName);
            this.dependencyGraph.removeEdge(edgeName);
        }
    }

    private String getEdgePostFix(String edgeName) {
        int index = edgeName.indexOf(Constants.LINK_SEPARATOR) + Constants.LINK_SEPARATOR.length();
        index = edgeName.substring(index).indexOf(Constants.LINK_SEPARATOR) + Constants.LINK_SEPARATOR.length();
        return edgeName.substring(index);
    }
}
