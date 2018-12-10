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
import org.wso2.vick.observability.model.generator.model.Edge;
import org.wso2.vick.observability.model.generator.model.Model;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

/**
 * This is the Manager, singleton class which performs the operations in the in memory dependency tree.
 */
public class ModelManager {
    private MutableNetwork<Node, String> dependencyGraph;

    public ModelManager() throws ModelException {
        try {
            Model model = ServiceHolder.getModelStoreManager().loadLastModel();
            this.dependencyGraph = NetworkBuilder.directed()
                    .allowsParallelEdges(true)
                    .expectedNodeCount(100000)
                    .expectedEdgeCount(1000000)
                    .build();
            if (model != null) {
                Set<Node> nodes = model.getNodes();
                Set<String> edges = Utils.getEdgesString(model.getEdges());
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
            String[] elements = Utils.edgeNameElements(edge);
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
            this.dependencyGraph.addEdge(parent, child, Utils.generateEdgeName(parent.getId(), child.getId(),
                    serviceName));
        } catch (Exception ignored) {
        }
    }

    public MutableNetwork<Node, String> getDependencyGraph() {
        return dependencyGraph;
    }

    public void moveLinks(Node fromNode, Node targetNode, String newEdgePrefix) {
        Set<String> outEdges = this.dependencyGraph.outEdges(fromNode);
        for (String edgeName : outEdges) {
            Node outNode = this.dependencyGraph.incidentNodes(outEdges).target();
            String newEdgeName = newEdgePrefix + Constants.LINK_SEPARATOR + Utils.getEdgePostFix(edgeName);
            this.addLink(targetNode, outNode, newEdgeName);
            this.dependencyGraph.removeEdge(edgeName);
        }
    }

    public Model getGraph(long fromTime, long toTime) throws GraphStoreException {
        if (fromTime == 0 && toTime == 0) {
            return new Model(this.dependencyGraph.nodes(), Utils.getEdges(this.dependencyGraph.edges()));
        } else {
            if (toTime == 0) {
                toTime = System.currentTimeMillis();
            }
            List<Model> models = ServiceHolder.getModelStoreManager().loadModel(fromTime, toTime);
            return getMergedModel(models);
        }
    }

    private Model getMergedModel(List<Model> models) {
        Set<Node> allNodes = new HashSet<>();
        Set<Edge> allEdges = new HashSet<>();
        for (Model model : models) {
            Set<Node> aModelNodes = model.getNodes();
            for (Node node : aModelNodes) {
                Node nodeFromAllNodes = Utils.getNode(allNodes, node);
                if (nodeFromAllNodes == null) {
                    allNodes.add(node);
                } else {
                    nodeFromAllNodes.getServices().addAll(node.getServices());
                }
            }
            allEdges.addAll(model.getEdges());
        }
        return new Model(allNodes, allEdges);
    }
}
