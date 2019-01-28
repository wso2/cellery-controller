/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.wso2.vick.observability.k8s.api.server.client;

import io.kubernetes.client.ApiClient;
import io.kubernetes.client.apis.CoreV1Api;
import io.kubernetes.client.models.V1Pod;
import io.kubernetes.client.models.V1PodList;
import io.kubernetes.client.util.Config;
import org.apache.log4j.Logger;
import org.wso2.siddhi.annotation.Example;
import org.wso2.siddhi.annotation.Extension;
import org.wso2.siddhi.core.config.SiddhiAppContext;
import org.wso2.siddhi.core.event.ComplexEventChunk;
import org.wso2.siddhi.core.event.stream.StreamEvent;
import org.wso2.siddhi.core.event.stream.StreamEventCloner;
import org.wso2.siddhi.core.event.stream.populater.ComplexEventPopulater;
import org.wso2.siddhi.core.executor.ExpressionExecutor;
import org.wso2.siddhi.core.query.processor.Processor;
import org.wso2.siddhi.core.query.processor.stream.StreamProcessor;
import org.wso2.siddhi.core.util.config.ConfigReader;
import org.wso2.siddhi.query.api.definition.AbstractDefinition;
import org.wso2.siddhi.query.api.definition.Attribute;
import org.wso2.siddhi.query.api.exception.SiddhiAppValidationException;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * This class implements the Stream Processor which can be used to call the K8s API Server and get data about Pods.
 */
@Extension(
        name = "getComponentPods",
        namespace = "k8sApiServerClient",
        description = "This is a client which calls the Kubernetes API server based on the received parameters and " +
                "adds the pod details received. This read the Service Account Token loaded into the pod and calls " +
                "the API Server using that.",
        examples = {
                @Example(
                        syntax = "k8sApiServerClient:getComponentPods()",
                        description = "This will fetch the currently running pods from the K8s API Servers"
                )
        }
)
public class GetComponentPodsStreamProcessor extends StreamProcessor {
    private static final Logger logger = Logger.getLogger(GetComponentPodsStreamProcessor.class.getName());

    private ApiClient k8sApiServerClient;

    @Override
    protected List<Attribute> init(AbstractDefinition inputDefinition,
                                   ExpressionExecutor[] attributeExpressionExecutors, ConfigReader configReader,
                                   SiddhiAppContext siddhiAppContext) {
        int attributeLength = attributeExpressionExecutors.length;
        if (attributeLength != 0) {
            throw new SiddhiAppValidationException("k8sApiServerClient expects exactly zero input parameters, but " +
                    attributeExpressionExecutors.length + " attributes found");
        }

        // Initializing the K8s API Client
        try {
            k8sApiServerClient = Config.defaultClient();
        } catch (IOException e) {
            String message = "Failed to initialize Kubernetes API Client";
            logger.error(message, e);
            throw new SiddhiAppValidationException(message);
        }

        List<Attribute> appendedAttributes = new ArrayList<>();
        appendedAttributes.add(new Attribute("cell", Attribute.Type.STRING));
        appendedAttributes.add(new Attribute("component", Attribute.Type.STRING));
        appendedAttributes.add(new Attribute("name", Attribute.Type.STRING));
        appendedAttributes.add(new Attribute("creationTimestamp", Attribute.Type.LONG));
        appendedAttributes.add(new Attribute("nodeName", Attribute.Type.STRING));
        return appendedAttributes;
    }

    @Override
    protected void process(ComplexEventChunk<StreamEvent> streamEventChunk, Processor nextProcessor,
                           StreamEventCloner streamEventCloner, ComplexEventPopulater complexEventPopulater) {
        ComplexEventChunk<StreamEvent> outputStreamEventChunk = new ComplexEventChunk<>(true);
        while (streamEventChunk.hasNext()) {
            StreamEvent incomingStreamEvent = streamEventChunk.next();
            addComponentPods(outputStreamEventChunk, incomingStreamEvent, Constants.COMPONENT_NAME_LABEL);
            addComponentPods(outputStreamEventChunk, incomingStreamEvent, Constants.GATEWAY_NAME_LABEL);
        }
        if (outputStreamEventChunk.getFirst() != null) {
            nextProcessor.process(outputStreamEventChunk);
        }
    }

    /**
     * Add component pods to output stream event chunk.
     * A new event will be cloned from the incoming stream event for each pod and added to the output event chunk.
     *
     * @param outputStreamEventChunk The output stream event chunk which will be sent to the next processor
     * @param incomingStreamEvent    The incoming stream event which will be cloned and used
     * @param componentNameLabel     The name of the label applied to store the component/gateway name
     */
    private void addComponentPods(ComplexEventChunk<StreamEvent> outputStreamEventChunk,
                                  StreamEvent incomingStreamEvent, String componentNameLabel) {
        // Calling the K8s API Servers to fetch component pods
        V1PodList componentPodList = null;
        try {
            CoreV1Api api = new CoreV1Api(k8sApiServerClient);
            componentPodList = api.listNamespacedPod(Constants.NAMESPACE, null, null,
                    Constants.RUNNING_STATUS_FIELD_SELECTOR, false,
                    Constants.CELL_NAME_LABEL + "," + componentNameLabel, null, null, null, false);
        } catch (Throwable e) {
            logger.error("Failed to fetch current pods for components", e);
        }

        if (componentPodList != null) {
            List<V1Pod> v1Pods = componentPodList.getItems();
            for (V1Pod v1Pod : v1Pods) {
                Object[] newData = new Object[5];
                newData[0] = v1Pod.getMetadata().getLabels().get(Constants.CELL_NAME_LABEL);
                newData[1] = getComponentName(v1Pod.getMetadata().getLabels().get(componentNameLabel));
                newData[2] = v1Pod.getMetadata().getName();
                newData[3] = v1Pod.getMetadata().getCreationTimestamp().getMillis();
                newData[4] = v1Pod.getSpec().getNodeName();

                StreamEvent streamEventCopy = streamEventCloner.copyStreamEvent(incomingStreamEvent);
                complexEventPopulater.populateComplexEvent(streamEventCopy, newData);
                outputStreamEventChunk.add(streamEventCopy);
            }
        }
    }

    /**
     * Get the actual component name from the fully qualified name.
     *
     * @param fullyQualifiedName The fully qualified name (eg:- "hr--hr")
     * @return The actual component name
     */
    private String getComponentName(String fullyQualifiedName) {
        String componentName = fullyQualifiedName;
        if (fullyQualifiedName.contains("--")) {
            componentName = fullyQualifiedName.split("--")[1];
        }
        return componentName;
    }

    @Override
    public void start() {
        // Do Nothing
    }

    @Override
    public void stop() {
        // Do Nothing
    }

    @Override
    public Map<String, Object> currentState() {
        // No State
        return null;
    }

    @Override
    public void restoreState(Map<String, Object> state) {
        // Do Nothing
    }
}
