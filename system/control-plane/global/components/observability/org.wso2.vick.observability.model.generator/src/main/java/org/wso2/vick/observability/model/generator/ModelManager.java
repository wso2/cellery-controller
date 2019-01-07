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

import com.google.common.graph.EndpointPair;
import com.google.common.graph.MutableNetwork;
import com.google.common.graph.NetworkBuilder;
import org.wso2.vick.observability.model.generator.exception.GraphStoreException;
import org.wso2.vick.observability.model.generator.exception.ModelException;
import org.wso2.vick.observability.model.generator.internal.ServiceHolder;
import org.wso2.vick.observability.model.generator.model.Edge;
import org.wso2.vick.observability.model.generator.model.Model;

import java.util.ArrayList;
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
            if (!parent.equals(child)) {
                this.dependencyGraph.addEdge(parent, child, Utils.generateEdgeName(parent.getId(), child.getId(),
                        serviceName));
            } else {
                parent.addEdge(serviceName);
            }
        } catch (Exception ignored) {
        }
    }

    public MutableNetwork<Node, String> getDependencyGraph() {
        return dependencyGraph;
    }

    public void moveLinks(Node targetNode, String newEdgePrefix, List<String> edgesToRemove, boolean moveOnlyTarget) {
        if (edgesToRemove != null) {
            for (String edgeName : edgesToRemove) {
                try {
                    EndpointPair<Node> endpointPair = this.dependencyGraph.incidentNodes(edgeName);
                    if (endpointPair != null) {
                        Node outNode = endpointPair.target();
                        if (!moveOnlyTarget || outNode.equals(targetNode)) {
                            String newEdgeName = newEdgePrefix + Constants.LINK_SEPARATOR
                                    + Utils.getEdgePostFix(edgeName);
                            this.addLink(targetNode, outNode, newEdgeName);
                            this.dependencyGraph.removeEdge(edgeName);
                        }
                    }
                } catch (Exception ignored) {
                }
            }
        }
    }

    public Model getGraph(long fromTime, long toTime) throws GraphStoreException {
        if (!ServiceHolder.getPeriodicProcessor().isStarted()) {
            ServiceHolder.getPeriodicProcessor().run();
        }
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

    public Model getDependencyModel(long fromTime, long toTime, String cellName) throws GraphStoreException {
        Model graph = getGraph(fromTime, toTime);
        List<Model> models = new ArrayList<>();
        Set<String> processedNodes = new HashSet<>();
        populateDependencyModel(graph, cellName, models, processedNodes);
        return getMergedModel(models);
    }

    private void populateDependencyModel(Model graph, String cellName, List<Model> models, Set<String> processedNodes) {
        if (!processedNodes.contains(cellName)) {
            Model dependencyGraph = getDependencyModel(graph, cellName);
            processedNodes.add(cellName);
            if (dependencyGraph.getNodes().size() > 1) {
                models.add(dependencyGraph);
                for (Node node : dependencyGraph.getNodes()) {
                    populateDependencyModel(graph, node.getId(), models, processedNodes);
                }
            } else {
                if (models.isEmpty()) {
                    models.add(dependencyGraph);
                }
            }
        }
    }


    private Model getDependencyModel(Model model, String cell) {
        Set<Node> dependedNodes = new HashSet<>();
        Set<Edge> dependedEdges = new HashSet<>();
        for (Edge edge : model.getEdges()) {
            if (edge.getSource().equalsIgnoreCase(cell)) {
                dependedEdges.add(edge);
                dependedNodes.add(Utils.getNode(model.getNodes(), new Node(edge.getTarget())));
            }
        }
        dependedNodes.add(Utils.getNode(model.getNodes(), new Node(cell)));
        return new Model(dependedNodes, dependedEdges);
    }

    public Model getDependencyModel(long fromTime, long toTime, String cellName, String serviceName)
            throws GraphStoreException {
        Model graph = getGraph(fromTime, toTime);
        Node cell = Utils.getNode(graph.getNodes(), new Node(cellName));
        Model serviceModel = null;
        if (cell != null) {
            String qualifiedServiceName = Utils.getQualifiedServiceName(cellName, serviceName);
            serviceModel = getDependencyModel(cell, serviceName);
            for (Edge edge : graph.getEdges()) {
                String edgeServiceName = Utils.getEdgeServiceName(edge.getEdgeString());
                String[] services = Utils.getServices(edgeServiceName);
                if (services[0].equalsIgnoreCase(serviceName) && !services[0].equalsIgnoreCase(services[1])) {
                    Node dependedCell = Utils.getNode(graph.getNodes(),
                            new Node(edge.getTarget()));
                    if (dependedCell != null) {
                        Model dependentCellModel = getDependencyModel(dependedCell, services[1]);
                        serviceModel.mergeModel(dependentCellModel,
                                new Edge(Utils.generateEdgeName(qualifiedServiceName,
                                        Utils.getQualifiedServiceName(dependedCell.getId(), services[1]), "")));
                    }
                }
            }
        }
        if (serviceModel == null) {
            serviceModel = new Model(new HashSet<>(), new HashSet<>());
        }
        return serviceModel;
    }

    private Model getDependencyModel(Node cell, String serviceName) {
        Set<Node> nodes = new HashSet<>();
        Set<Edge> edges = new HashSet<>();
        String qualifiedServiceName = Utils.getQualifiedServiceName(cell.getId(), serviceName);
        nodes.add(new Node(qualifiedServiceName));
        for (String edge : cell.getEdges()) {
            String[] sourceTarget = edge.split(Constants.LINK_SEPARATOR);
            if (sourceTarget[0].trim().equalsIgnoreCase(serviceName) &&
                    !sourceTarget[0].trim().equalsIgnoreCase(sourceTarget[1])) {
                String qualifiedTargetServiceName = Utils.getQualifiedServiceName(cell.getId(),
                        sourceTarget[1].trim());
                nodes.add(new Node(qualifiedTargetServiceName));
                edges.add(new Edge(Utils.generateEdgeName(qualifiedServiceName, qualifiedTargetServiceName,
                        "")));
                Model nextLevelModel = getDependencyModel(cell, sourceTarget[1]);
                if (nextLevelModel.getNodes().size() > 1) {
                    nodes.addAll(nextLevelModel.getNodes());
                    edges.addAll(nextLevelModel.getEdges());
                }
            }
        }
        return new Model(nodes, edges);
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
