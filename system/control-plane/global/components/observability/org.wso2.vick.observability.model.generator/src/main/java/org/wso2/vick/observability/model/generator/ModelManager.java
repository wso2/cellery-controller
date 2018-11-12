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

import java.util.Set;

/**
 * This is the Manager, singleton class which performs the operations in the in memory dependency tree.
 */

public class ModelManager {
    private static final ModelManager instance = new ModelManager();
    private static final String EDGE_NAME_CONNECTOR = " ---> ";

    private MutableNetwork<Node, String> dependencyGraph;

    private ModelManager() {
        this.dependencyGraph = NetworkBuilder.directed()
                .allowsParallelEdges(true)
                .expectedNodeCount(100000)
                .expectedEdgeCount(1000000)
                .build();
    }

    public static ModelManager getInstance() {
        return instance;
    }

    public void addNode(Node node) {
        this.dependencyGraph.addNode(node);
    }

    public void addLink(Node parent, Node child) {
        this.dependencyGraph.addEdge(parent, child, getLinkName(parent, child));
    }

    private String getLinkName(Node parent, Node child) {
        return parent.getId() + EDGE_NAME_CONNECTOR + child.getId();
    }

    public Set<Node> getNodes() {
        return this.dependencyGraph.nodes();
    }

    public Set<String> getLinks() {
        return this.dependencyGraph.edges();
    }

    public String[] getParentChildNodeNames(String edgeName) {
        return edgeName.split(EDGE_NAME_CONNECTOR);
    }
}
