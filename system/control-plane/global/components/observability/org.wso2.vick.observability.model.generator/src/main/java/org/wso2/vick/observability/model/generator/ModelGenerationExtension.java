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

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
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

    private ExpressionExecutor componentNameExecutor;
    private ExpressionExecutor spanIdExecutor;
    private ExpressionExecutor parentIdExecutor;
    private ExpressionExecutor serviceNameExecutor;
    private ExpressionExecutor tagExecutor;
    private Map<String, List<Node>> pendingEdges = new HashMap<>();
    private Cache<String, Node> spanIdNodeCache = CacheBuilder.newBuilder().
            expireAfterAccess(60, TimeUnit.SECONDS).maximumSize(100000).build();

    @Override
    protected void process(ComplexEventChunk<StreamEvent> complexEventChunk, Processor processor,
                           StreamEventCloner streamEventCloner, ComplexEventPopulater complexEventPopulater) {
        while (complexEventChunk.hasNext()) {
            StreamEvent streamEvent = complexEventChunk.next();
            String componentName = (String) componentNameExecutor.execute(streamEvent);
            if (!componentName.trim().startsWith("istio")) {
                String spanId = (String) spanIdExecutor.execute(streamEvent);
                String parentId = (String) parentIdExecutor.execute(streamEvent);
                String serviceName = (String) serviceNameExecutor.execute(streamEvent);
                String tags = (String) tagExecutor.execute(streamEvent);
                log.info(spanId);
                spanId = spanId.split("-")[0].trim();
                Node node = new Node(componentName, serviceName, tags);
                spanIdNodeCache.put(spanId, node);

                ModelManager.getInstance().addNode(node);
                if (parentId != null) {
                    Node parentNode = spanIdNodeCache.getIfPresent(parentId);
                    if (parentNode != null) {
                        ModelManager.getInstance().addLink(parentNode, node, serviceName);
                    } else {
                        List<Node> waitingNodes = pendingEdges.putIfAbsent(parentId,
                                new ArrayList<>(Collections.singletonList(node)));
                        if (waitingNodes != null) {
                            waitingNodes.add(node);
                        }
                    }
                }
                List<Node> pendingChildNodes = this.pendingEdges.get(spanId);
                if (pendingChildNodes != null) {
                    for (Node child : pendingChildNodes) {
                        Node parentNode = spanIdNodeCache.getIfPresent(spanId);
                        if (parentNode != null) {
                            ModelManager.getInstance().addLink(parentNode, child, parentNode.getServiceName());
                        }
                    }
                }
                this.pendingEdges.remove(spanId);
            }
        }
    }

    private String getNodeName(String componentName, String serviceName) {
        return componentName + " - " + serviceName;
    }

    @Override
    protected List<Attribute> init(AbstractDefinition abstractDefinition, ExpressionExecutor[] expressionExecutors,
                                   ConfigReader configReader, SiddhiAppContext siddhiAppContext) {
        if (expressionExecutors.length != 5) {
            throw new SiddhiAppCreationException("Minimum number of attributes is 4");
        } else {
            if (expressionExecutors[0].getReturnType() == Attribute.Type.STRING) {
                componentNameExecutor = expressionExecutors[0];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the component name "
                        + "field, but found a field with return type - " + expressionExecutors[0].getReturnType());
            }

            if (expressionExecutors[1].getReturnType() == Attribute.Type.STRING) {
                spanIdExecutor = expressionExecutors[1];
            } else {
                throw new SiddhiAppCreationException("Expected a field with Long return type for the span id field," +
                        "but found a field with return type - " + expressionExecutors[1].getReturnType());
            }

            if (expressionExecutors[2].getReturnType() == Attribute.Type.STRING) {
                parentIdExecutor = expressionExecutors[2];
            } else {
                throw new SiddhiAppCreationException("Expected a field with Long return type for the parent id field,"
                        + "but found a field with return type - " + expressionExecutors[2].getReturnType());
            }

            if (expressionExecutors[3].getReturnType() == Attribute.Type.STRING) {
                serviceNameExecutor = expressionExecutors[3];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the service name" +
                        " field, but found a field with return type - " + expressionExecutors[3].getReturnType());
            }

            if (expressionExecutors[4].getReturnType() == Attribute.Type.STRING) {
                tagExecutor = expressionExecutors[4];
            } else {
                throw new SiddhiAppCreationException("Expected a field with String return type for the tags field," +
                        "but found a field with return type - " + expressionExecutors[4].getReturnType());
            }
        }
        return new ArrayList<Attribute>();
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
