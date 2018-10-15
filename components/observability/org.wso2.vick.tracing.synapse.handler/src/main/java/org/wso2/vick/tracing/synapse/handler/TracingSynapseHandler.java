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
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.vick.tracing.synapse.handler;

import brave.Tracing;
import brave.opentracing.BraveTracer;
import brave.propagation.B3Propagation;
import io.opentracing.Span;
import io.opentracing.SpanContext;
import io.opentracing.Tracer;
import io.opentracing.propagation.Format;
import io.opentracing.propagation.TextMapExtractAdapter;
import io.opentracing.propagation.TextMapInjectAdapter;
import org.apache.log4j.Logger;
import org.apache.synapse.AbstractSynapseHandler;
import org.apache.synapse.MessageContext;
import org.apache.synapse.core.axis2.Axis2MessageContext;
import zipkin2.reporter.AsyncReporter;
import zipkin2.reporter.Sender;
import zipkin2.reporter.okhttp3.OkHttpSender;

import java.util.HashMap;
import java.util.Map;

/**
 * Tracing Handler to work with tracing headers and publish tracing data.
 */
public class TracingSynapseHandler extends AbstractSynapseHandler {

    private static final Logger logger = Logger.getLogger(TracingSynapseHandler.class.getName());

    private Tracer tracer;

    public TracingSynapseHandler() {
        Map properties = getProperties();
        String hostname = (String) properties.get(Constants.ZIPKIN_HOST);
        int port = Integer.parseInt((String) properties.get(Constants.ZIPKIN_PORT));
        String apiContext = (String) properties.get(Constants.ZIPKIN_API_CONTEXT);

        // Instantiating the tracer
        String tracingReceiverEndpoint = "http://" + hostname + ":" + port + apiContext;
        Sender sender = OkHttpSender.create(tracingReceiverEndpoint);
        if (logger.isDebugEnabled()) {
            logger.debug("Initialized tracing sender to send to " + tracingReceiverEndpoint);
        }
        Tracing braveTracing = Tracing.newBuilder()
                .localServiceName(Constants.GLOBAL_GATEWAY_SERVICE_NAME)
                .spanReporter(AsyncReporter.create(sender))
                .build();
        tracer = BraveTracer.newBuilder(braveTracing)
                .textMapPropagation(Format.Builtin.HTTP_HEADERS, B3Propagation.B3_STRING)
                .build();

        logger.info("Initialized VICK Tracer");
    }

    @Override
    public boolean handleRequestInFlow(MessageContext messageContext) {
        // Extracting the B3 headers from the incoming headers
        Map<String, String> headersMap = extractHeadersFromSynapseContext(messageContext);
        SpanContext parentSpanContext = tracer.extract(
                Format.Builtin.HTTP_HEADERS,
                new TextMapExtractAdapter(headersMap)
        );

        // Building the request in span
        String spanName = messageContext.getTo().getAddress();
        Span span = tracer.buildSpan(spanName)
                .asChildOf(parentSpanContext)
                .start();
        if (logger.isDebugEnabled()) {
            logger.debug("Started span: " + spanName);
        }
        messageContext.setProperty(Constants.REQUEST_IN_SPAN, span);
        return true;
    }

    @Override
    public boolean handleRequestOutFlow(MessageContext messageContext) {
        Object requestInSpan = messageContext.getProperty(Constants.REQUEST_IN_SPAN);
        if (requestInSpan instanceof Span) {
            Span parentSpan = (Span) requestInSpan;

            // Building the request out span
            String spanName = messageContext.getTo().getAddress();
            Span span = tracer.buildSpan(spanName)
                    .asChildOf(parentSpan)
                    .start();
            if (logger.isDebugEnabled()) {
                logger.debug("Started span: " + spanName);
            }
            messageContext.setProperty(Constants.REQUEST_OUT_SPAN, span);

            // Injecting B3 headers into the outgoing headers
            Map<String, String> headersMap = extractHeadersFromSynapseContext(messageContext);
            tracer.inject(
                    span.context(),
                    Format.Builtin.HTTP_HEADERS,
                    new TextMapInjectAdapter(headersMap)
            );
        }
        return true;
    }

    @Override
    public boolean handleResponseInFlow(MessageContext messageContext) {
        return finishSpan(messageContext, Constants.REQUEST_OUT_SPAN);
    }

    @Override
    public boolean handleResponseOutFlow(MessageContext messageContext) {
        return finishSpan(messageContext, Constants.REQUEST_IN_SPAN);
    }

    /**
     * Finish an existing span and remove the span from the store
     *
     * @param messageContext The synapse message context
     * @return Whether the mediation flow should continue
     */
    private boolean finishSpan(MessageContext messageContext, String type) {
        Span span = (Span) messageContext.getProperty(type);
        if (span != null) {
            span.finish();
            if (logger.isDebugEnabled()) {
                logger.debug("Finished span");
            }
        }
        return true;
    }

    /**
     * Extract headers map from synapse message context
     *
     * @param synapseMessageContext Synapse message context
     * @return Headers map
     */
    private Map<String, String> extractHeadersFromSynapseContext(MessageContext synapseMessageContext) {
        Map<String, String> headersMap = null;
        if (synapseMessageContext instanceof Axis2MessageContext) {
            Axis2MessageContext axis2MessageContext = ((Axis2MessageContext) synapseMessageContext);
            Object headers = axis2MessageContext.getAxis2MessageContext()
                    .getProperty(org.apache.axis2.context.MessageContext.TRANSPORT_HEADERS);
            if (headers instanceof Map) {
                headersMap = (Map<String, String>) headers;
            }
        }
        if (headersMap == null) {
            headersMap = new HashMap<>(0);
        }
        return headersMap;
    }
}
