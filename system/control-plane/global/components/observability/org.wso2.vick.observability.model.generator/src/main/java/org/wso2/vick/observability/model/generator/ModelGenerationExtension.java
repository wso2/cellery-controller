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

import com.google.common.cache.Cache;
import com.google.common.cache.CacheBuilder;
import org.apache.log4j.Logger;
import org.wso2.siddhi.annotation.Example;
import org.wso2.siddhi.annotation.Extension;
import org.wso2.siddhi.core.config.SiddhiAppContext;
import org.wso2.siddhi.core.event.ComplexEventChunk;
import org.wso2.siddhi.core.event.stream.StreamEvent;
import org.wso2.siddhi.core.event.stream.StreamEventCloner;
import org.wso2.siddhi.core.event.stream.populater.ComplexEventPopulater;
import org.wso2.siddhi.core.exception.SiddhiAppCreationException;
import org.wso2.siddhi.core.executor.ExpressionExecutor;
import org.wso2.siddhi.core.query.processor.Processor;
import org.wso2.siddhi.core.query.processor.stream.StreamProcessor;
import org.wso2.siddhi.core.util.config.ConfigReader;
import org.wso2.siddhi.query.api.definition.AbstractDefinition;
import org.wso2.siddhi.query.api.definition.Attribute;
import org.wso2.vick.observability.model.generator.internal.ServiceHolder;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;

/**
 * This is the Siddhi extension which generates the dependency graph for the spans created.
 */
@Extension(
        name = "dependencyTree",
        namespace = "tracing",
        description = "This generates the dependency model of the spans",
        examples = @Example(description = "TBD"
                , syntax = "from inputStream#tracing:dependencyTree(componentName, spanId, parentId, serviceName," +
                " tags) \" +\n" +
                "                \"select * \n" +
                "                \"insert into outputStream;")
)
public class ModelGenerationExtension extends StreamProcessor {

    private static final Logger log = Logger.getLogger(ModelGenerationExtension.class);

    private ExpressionExecutor cellNameExecutor;
    private ExpressionExecutor serviceNameExecutor;
    private ExpressionExecutor operationNameExecutor;
    private ExpressionExecutor spanIdExecutor;
    private ExpressionExecutor parentIdExecutor;
    private ExpressionExecutor spanKindExecutor;
    private ExpressionExecutor tagExecutor;
    private final Map<String, List<SpanCacheInfo.NodeInfo>> pendingEdges = new ConcurrentHashMap<>();
    private final Cache<String, SpanCacheInfo> spanIdNodeCache = CacheBuilder.newBuilder().
            expireAfterAccess(60, TimeUnit.SECONDS).maximumSize(100000).build();
    private final Map<String, Node> nodeCache = new HashMap<>();


    @Override
    protected void process(ComplexEventChunk<StreamEvent> complexEventChunk, Processor processor,
                           StreamEventCloner streamEventCloner, ComplexEventPopulater complexEventPopulater) {
        while (complexEventChunk.hasNext()) {
            StreamEvent streamEvent = complexEventChunk.next();
            String cellName = (String) cellNameExecutor.execute(streamEvent);
            String serviceName = (String) serviceNameExecutor.execute(streamEvent);
            String operationName = (String) operationNameExecutor.execute(streamEvent);
            String spanId = (String) spanIdExecutor.execute(streamEvent);
            if (cellName != null && !cellName.trim().equalsIgnoreCase("")) {
                String spanKind = (String) spanKindExecutor.execute(streamEvent);
                String parentId = (String) parentIdExecutor.execute(streamEvent);
                String tags = (String) tagExecutor.execute(streamEvent);
                spanId = spanId.split(Constants.SPAN_ID_KIND_SEPARATOR)[0];
                Node node = getNode(cellName, tags);
                node.addService(serviceName);
                SpanCacheInfo spanCacheInfo = setSpanInfo(spanId, node, serviceName, operationName, spanKind);
                if (spanKind.equalsIgnoreCase(Constants.SERVER_SPAN_KIND) && spanCacheInfo.getClient() != null) {
                    ServiceHolder.getModelManager().moveLinks(spanCacheInfo.getClient().getNode(),
                            spanCacheInfo.getServer().getNode(),
                            serviceName + Constants.LINK_SEPARATOR + operationName);
                }

                ServiceHolder.getModelManager().addNode(node);
                if (parentId != null) {
                    SpanCacheInfo parentSpanCacheInfo = spanIdNodeCache.getIfPresent(parentId);
                    if (parentSpanCacheInfo != null) {
                        addLink(parentSpanCacheInfo, node, serviceName, operationName);
                    } else {
                        SpanCacheInfo.NodeInfo pendingNode;
                        if (spanKind.equalsIgnoreCase(Constants.LINK_SEPARATOR)) {
                            pendingNode = spanCacheInfo.getServer();
                        } else {
                            pendingNode = spanCacheInfo.getClient();
                        }
                        List<SpanCacheInfo.NodeInfo> waitingNodes = pendingEdges.putIfAbsent(parentId,
                                new ArrayList<>(Collections.singletonList(pendingNode)));
                        if (waitingNodes != null) {
                            waitingNodes.add(pendingNode);
                        }
                    }
                }
                List<SpanCacheInfo.NodeInfo> pendingChildNodes = this.pendingEdges.get(spanId);
                if (pendingChildNodes != null) {
                    for (SpanCacheInfo.NodeInfo child : pendingChildNodes) {
                        if (child != null) {
                            addLink(spanCacheInfo, child.getNode(), child.getService(), child.getOperationName());
                        }
                    }
                }
                this.pendingEdges.remove(spanId);
            }
        }
    }

    private void addLink(SpanCacheInfo parentSpan, Node childNode, String serviceName, String operationName) {
        SpanCacheInfo.NodeInfo parentNode;
        if (parentSpan.getServer() != null) {
            parentNode = parentSpan.getServer();
        } else {
            parentNode = parentSpan.getClient();
        }
        String linkName = parentNode.getService() + Constants.LINK_SEPARATOR + parentNode.getOperationName()
                + Constants.LINK_SEPARATOR + serviceName + Constants.LINK_SEPARATOR + operationName;
        ServiceHolder.getModelManager().addLink(parentNode.getNode(), childNode, linkName);
    }

    private Node getNode(String componentName, String tags) {
        Node node = nodeCache.get(componentName);
        if (node == null) {
            synchronized (nodeCache) {
                node = nodeCache.get(componentName);
                if (node == null) {
                    node = new Node(componentName, tags);
                    this.nodeCache.put(componentName, node);
                }
            }
        }
        return node;
    }

    private SpanCacheInfo setSpanInfo(String spanId, Node node, String serviceName, String operationName,
                                      String type) {
        SpanCacheInfo.Type nodeType;
        if (type.equalsIgnoreCase(Constants.SERVER_SPAN_KIND)) {
            nodeType = SpanCacheInfo.Type.SERVER;
        } else {
            nodeType = SpanCacheInfo.Type.CLIENT;
        }
        SpanCacheInfo spanInfo = spanIdNodeCache.getIfPresent(spanId);
        if (spanInfo == null) {
            synchronized (spanIdNodeCache) {
                spanInfo = spanIdNodeCache.getIfPresent(spanId);
                if (spanInfo == null) {
                    spanInfo = new SpanCacheInfo(spanId, new SpanCacheInfo.NodeInfo(node,
                            serviceName, operationName), nodeType);
                    spanIdNodeCache.put(spanId, spanInfo);
                    return spanInfo;
                }
            }
        }
        spanInfo.setNodeInfo(new SpanCacheInfo.NodeInfo(node, serviceName, operationName), nodeType);
        return spanInfo;
    }

    @Override
    protected List<Attribute> init(AbstractDefinition abstractDefinition, ExpressionExecutor[] expressionExecutors,
                                   ConfigReader configReader, SiddhiAppContext siddhiAppContext) {
        if (expressionExecutors.length != 7) {
            throw new SiddhiAppCreationException("Minimum number of attributes is six");
        } else {
            if (expressionExecutors[0].getReturnType() == Attribute.Type.STRING) {
                cellNameExecutor = expressionExecutors[0];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the component name "
                        + "field, but found a field with return type - " + expressionExecutors[0].getReturnType());
            }

            if (expressionExecutors[1].getReturnType() == Attribute.Type.STRING) {
                serviceNameExecutor = expressionExecutors[1];
            } else {
                throw new SiddhiAppCreationException("Expected a field with Long return type for the span id field," +
                        "but found a field with return type - " + expressionExecutors[1].getReturnType());
            }

            if (expressionExecutors[2].getReturnType() == Attribute.Type.STRING) {
                operationNameExecutor = expressionExecutors[2];
            } else {
                throw new SiddhiAppCreationException("Expected a field with Long return type for the parent id field,"
                        + "but found a field with return type - " + expressionExecutors[2].getReturnType());
            }

            if (expressionExecutors[3].getReturnType() == Attribute.Type.STRING) {
                spanIdExecutor = expressionExecutors[3];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the service name" +
                        " field, but found a field with return type - " + expressionExecutors[3].getReturnType());
            }

            if (expressionExecutors[4].getReturnType() == Attribute.Type.STRING) {
                parentIdExecutor = expressionExecutors[4];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the tags field," +
                        "but found a field with return type - " + expressionExecutors[4].getReturnType());
            }

            if (expressionExecutors[5].getReturnType() == Attribute.Type.STRING) {
                spanKindExecutor = expressionExecutors[5];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the " +
                        "spanKind field, but found a field with return type - "
                        + expressionExecutors[5].getReturnType());
            }

            if (expressionExecutors[6].getReturnType() == Attribute.Type.STRING) {
                tagExecutor = expressionExecutors[6];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the " +
                        "spanKind field, but found a field with return type - "
                        + expressionExecutors[5].getReturnType());
            }
        }
        return new ArrayList<>();
    }

    @Override
    public void start() {

    }

    @Override
    public void stop() {

    }

    @Override
    public Map<String, Object> currentState() {
        return null;
    }

    @Override
    public void restoreState(Map<String, Object> map) {

    }
}
