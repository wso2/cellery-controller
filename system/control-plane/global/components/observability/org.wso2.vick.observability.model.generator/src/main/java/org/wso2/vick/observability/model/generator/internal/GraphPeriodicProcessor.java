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
package org.wso2.vick.observability.model.generator.internal;

import com.google.common.graph.MutableNetwork;
import org.apache.log4j.Logger;
import org.wso2.vick.observability.model.generator.Node;
import org.wso2.vick.observability.model.generator.exception.GraphStoreException;

import java.util.Set;

/**
 * This is the Scheduled executor which periodically runs
 */
public class GraphPeriodicProcessor implements Runnable {
    private static final Logger log = Logger.getLogger(GraphPeriodicProcessor.class);
    private Object[] lastModel;


    @Override
    public void run() {
        try {
            MutableNetwork<Node, String> currentModel = ServiceHolder.getModelManager().getDependencyGraph();
            if (lastModel == null) {
                this.lastModel = ServiceHolder.getGraphStoreManager().loadGraph();
            }
            if (this.lastModel == null) {
                if (currentModel.nodes().size() != 0) {
                    lastModel = ServiceHolder.getGraphStoreManager().persistGraph(currentModel);
                }
            } else {
                Set<Node> currentNodes = currentModel.nodes();
                Set<String> currentEdges = currentModel.edges();
                Set<Node> lastNodes = (Set<Node>) this.lastModel[0];
                Set<String> lastEdges = (Set<String>) this.lastModel[1];
                if (currentEdges.size() == lastEdges.size() && currentNodes.size() == lastNodes.size()) {
                    if (isSameNodes(currentNodes, lastNodes) && isSameEdges(currentEdges, lastEdges)) {
                        return;
                    }
                }
                lastModel = ServiceHolder.getGraphStoreManager().persistGraph(currentModel);
            }
        } catch (GraphStoreException e) {
            log.error("Error occurred while handling the dependency graph persistence. ", e);
        }
    }

    private boolean isSameNodes(Set<Node> currentNodes, Set<Node> lastNodes) {
        boolean isSameModel = true;
        for (Node currentNode : currentNodes) {
            Node lastNode = null;
            //Find the node from the last persisted graph
            for (Node node : lastNodes) {
                if (node.equals(currentNode)) {
                    lastNode = node;
                    break;
                }
            }
            // Check the services within the node is also same.
            if (lastNode != null) {
                if (lastNode.getServices().size() == currentNode.getServices().size()) {
                    if (!lastNode.getServices().containsAll(currentNode.getServices())) {
                        isSameModel = false;
                        break;
                    }
                } else {
                    isSameModel = false;
                    break;
                }
            } else {
                isSameModel = false;
                break;
            }
        }
        return isSameModel;
    }

    private boolean isSameEdges(Set<String> currentEdges, Set<String> lastEdges) {
        return currentEdges.containsAll(lastEdges);
    }
}

