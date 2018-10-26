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
import zipkin2.codec.SpanBytesEncoder;
import zipkin2.reporter.AsyncReporter;
import zipkin2.reporter.urlconnection.URLConnectionSender;

import java.util.HashMap;
import java.util.Map;

/**
 * Tracing Handler to work with tracing headers and publish tracing data.
 */
public class TracingSynapseHandler extends AbstractSynapseHandler {

    private static final Logger logger = Logger.getLogger(TracingSynapseHandler.class.getName());

    private Tracer tracer;

    public TracingSynapseHandler() {
        String hostname = System.getenv(Constants.ZIPKIN_HOST);
        int port = Integer.parseInt(System.getenv(Constants.ZIPKIN_PORT));
        String apiContext = System.getenv(Constants.ZIPKIN_API_CONTEXT);

        // Instantiating the reporter
        String tracingReceiverEndpoint = "http://" + hostname + ":" + port + apiContext;
        URLConnectionSender sender = URLConnectionSender.create(tracingReceiverEndpoint).toBuilder()
                .compressionEnabled(Constants.TRACING_SENDER_COMPRESSION_ENABLED)
                .build();
        AsyncReporter<zipkin2.Span> reporter = AsyncReporter.builder(sender)
                .build(SpanBytesEncoder.JSON_V1);
        if (logger.isDebugEnabled()) {
            logger.debug("Initialized tracing sender to send to " + tracingReceiverEndpoint);
        }

        // Instantiating the tracer
        Tracing braveTracing = Tracing.newBuilder()
                .localServiceName(Constants.GLOBAL_GATEWAY_SERVICE_NAME)
                .spanReporter(reporter)
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

        // Settings tags
        org.apache.axis2.context.MessageContext axis2MessageContext = getAxis2MessageContext(messageContext);
        addTag(span, Constants.TAG_KEY_SPAN_KIND, Constants.SPAN_KIND_SERVER);
        addTag(span, Constants.TAG_KEY_HTTP_METHOD,
                axis2MessageContext.getProperty(Constants.AXIS2_MESSAGE_CONTEXT_PROPERTY_HTTP_METHOD));
        addTag(span, Constants.TAG_KEY_PEER_ADDRESS,
                axis2MessageContext.getProperty(Constants.AXIS2_MESSAGE_CONTEXT_PROPERTY_REMOTE_HOST));
        addTag(span, Constants.TAG_KEY_PROTOCOL,
                axis2MessageContext.getProperty(Constants.AXIS2_MESSAGE_CONTEXT_PROPERTY_REMOTE_HOST));
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

            // Settings tags
            addTag(span, Constants.TAG_KEY_SPAN_KIND, Constants.SPAN_KIND_CLIENT);
            addTag(span, Constants.TAG_KEY_HTTP_METHOD,
                    messageContext.getProperty(Constants.SYNAPSE_MESSAGE_CONTEXT_PROPERTY_HTTP_METHOD));
            addTag(span, Constants.TAG_KEY_HTTP_URL,
                    messageContext.getProperty(Constants.SYNAPSE_MESSAGE_CONTEXT_PROPERTY_ENDPOINT));
            addTag(span, Constants.TAG_KEY_PEER_ADDRESS,
                    messageContext.getProperty(Constants.SYNAPSE_MESSAGE_CONTEXT_PROPERTY_PEER_ADDRESS));
            addTag(span, Constants.TAG_KEY_PROTOCOL,
                    messageContext.getProperty(Constants.SYNAPSE_MESSAGE_CONTEXT_PROPERTY_TRANSPORT));

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
     * Finish the stored in message context.
     *
     * @param messageContext The synapse message context
     * @param type           The type of span to finish
     * @return True if sequence should continue
     */
    private boolean finishSpan(MessageContext messageContext, String type) {
        Span span = (Span) messageContext.getProperty(type);
        if (span != null) {
            org.apache.axis2.context.MessageContext axis2MessageContext = getAxis2MessageContext(messageContext);

            if (axis2MessageContext != null) {
                // Settings tags
                addTag(span, Constants.TAG_KEY_HTTP_STATUS_CODE,
                        axis2MessageContext.getProperty(Constants.AXIS2_MESSAGE_CONTEXT_PROPERTY_HTTP_STATUS_CODE));
            }

            span.finish();
            if (logger.isDebugEnabled()) {
                logger.debug("Finished span");
            }
        }
        return true;
    }

    /**
     * Extract headers map from synapse message context.
     *
     * @param synapseMessageContext Synapse message context
     * @return Headers map
     */
    private Map<String, String> extractHeadersFromSynapseContext(MessageContext synapseMessageContext) {
        Map<String, String> headersMap = null;
        org.apache.axis2.context.MessageContext axis2MessageContext = getAxis2MessageContext(synapseMessageContext);
        if (axis2MessageContext != null) {
            Object headers = axis2MessageContext.getProperty(org.apache.axis2.context.MessageContext.TRANSPORT_HEADERS);
            if (headers instanceof Map) {
                headersMap = (Map<String, String>) headers;
            }
        }
        if (headersMap == null) {
            headersMap = new HashMap<>(0);
        }
        return headersMap;
    }

    /**
     * Get the Axis2 message context from the synapse message context.
     *
     * @param synapseMessageContext Synapse message context
     * @return The relevant Axis2 message context
     */
    private org.apache.axis2.context.MessageContext getAxis2MessageContext(MessageContext synapseMessageContext) {
        org.apache.axis2.context.MessageContext axis2MessageContext = null;
        if (synapseMessageContext instanceof Axis2MessageContext) {
            axis2MessageContext = ((Axis2MessageContext) synapseMessageContext).getAxis2MessageContext();
        }
        return axis2MessageContext;
    }

    /**
     * Add a tag to a span.
     *
     * @param span     The span to which the tag should be added
     * @param tagKey   The key of the tag to be added
     * @param tagValue The value of the tag to be added
     */
    private void addTag(Span span, String tagKey, Object tagValue) {
        if (tagValue != null) {
            if (tagValue instanceof String) {
                span.setTag(tagKey, (String) tagValue);
            } else if (tagValue instanceof Number) {
                span.setTag(tagKey, (Number) tagValue);
            } else if (tagValue instanceof Boolean) {
                span.setTag(tagKey, (boolean) tagValue);
            }
        }
    }
}
