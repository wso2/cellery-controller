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

package org.wso2.vick.tracing.receiver;

import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import io.netty.buffer.ByteBuf;
import org.apache.log4j.Logger;
import org.apache.thrift.TException;
import org.wso2.siddhi.annotation.Example;
import org.wso2.siddhi.annotation.Extension;
import org.wso2.siddhi.annotation.Parameter;
import org.wso2.siddhi.annotation.util.DataType;
import org.wso2.siddhi.core.config.SiddhiAppContext;
import org.wso2.siddhi.core.exception.ConnectionUnavailableException;
import org.wso2.siddhi.core.stream.input.source.Source;
import org.wso2.siddhi.core.stream.input.source.SourceEventListener;
import org.wso2.siddhi.core.util.config.ConfigReader;
import org.wso2.siddhi.core.util.transport.OptionHolder;
import org.wso2.transport.http.netty.config.ListenerConfiguration;
import org.wso2.transport.http.netty.contract.HttpConnectorListener;
import org.wso2.transport.http.netty.contract.HttpWsConnectorFactory;
import org.wso2.transport.http.netty.contract.ServerConnector;
import org.wso2.transport.http.netty.contract.ServerConnectorFuture;
import org.wso2.transport.http.netty.contractimpl.DefaultHttpWsConnectorFactory;
import org.wso2.transport.http.netty.listener.ServerBootstrapConfiguration;
import org.wso2.transport.http.netty.message.HTTPCarbonMessage;
import org.wso2.vick.tracing.receiver.internal.Codec;
import org.wso2.vick.tracing.receiver.internal.ZipkinSpan;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

/**
 * This class implements the event source, where the received telemetry attributes can be injected to streams.
 */
@Extension(
        name = "tracing-receiver",
        namespace = "source",
        description = "This is the tracing Receiver for VICK. This accepts Zipkin encoded tracing data. " +
                "The event source outputs a map of attributes. Therefore a key-value mapper needs to be used.",
        parameters = {
                @Parameter(
                        name = "ip",
                        type = DataType.STRING,
                        description = "IP to which the server connector should be bound to",
                        optional = true,
                        defaultValue = "0.0.0.0"
                ),
                @Parameter(
                        name = "port",
                        type = DataType.INT,
                        description = "Port on which the server connector should listen on",
                        optional = true,
                        defaultValue = "9411"
                )
        },
        examples = {
                @Example(
                        syntax = "@source(type='tracing-receiver', @map(type='keyvalue', " +
                                "fail.on.missing.attribute='false'))\n" +
                                "define stream ZipkinStreamIn (traceId string, id string, parentId string, " +
                                "name string, serviceName string, kind string, timestamp long, duration long, " +
                                "tags string)",
                        description = "This produced events when Zipkin tracing data is received on amy interface on " +
                                "port 9411. The stream definition of the event source is fixed since it depends on " +
                                "the Zipkin format"
                )
        }
)
public class TracingEventSource extends Source {

    private static final Logger logger = Logger.getLogger(TracingEventSource.class.getName());
    private static final Gson gson = new Gson();
    private static final String HTTP_SERVER_ID = "TRACING_HTTP_SERVER";

    private HttpWsConnectorFactory httpWsConnectorFactory;
    private ServerConnector serverConnector;
    private HttpServerListener httpServerListener;

    private String ip;
    private int port;

    @Override
    public void init(SourceEventListener sourceEventListener, OptionHolder optionHolder,
                     String[] requestedTransportPropertyNames, ConfigReader configReader,
                     SiddhiAppContext siddhiAppContext) {
        ip = optionHolder.validateAndGetStaticValue(Constants.TRACING_RECEIVER_IP_KEY,
                Constants.DEFAULT_TRACING_RECEIVER_IP);
        port = Integer.parseInt(optionHolder.validateAndGetStaticValue(Constants.TRACING_RECEIVER_PORT_KEY,
                Constants.DEFAULT_TRACING_RECEIVER_PORT));

        httpServerListener = new HttpServerListener(sourceEventListener);
        httpWsConnectorFactory = new DefaultHttpWsConnectorFactory();
    }

    @Override
    public void connect(ConnectionCallback connectionCallback) throws ConnectionUnavailableException {
        ServerBootstrapConfiguration serverBootstrapConfiguration
                = new ServerBootstrapConfiguration(new HashMap<>());
        ListenerConfiguration listenerConfiguration
                = new ListenerConfiguration(HTTP_SERVER_ID, ip, port);

        serverConnector = httpWsConnectorFactory
                .createServerConnector(serverBootstrapConfiguration, listenerConfiguration);

        ServerConnectorFuture serverConnectorFuture = serverConnector.start();
        if (logger.isDebugEnabled()) {
            logger.debug("Successfully started HTTP Server Connector " + serverConnector.getConnectorID());
        }
        serverConnectorFuture.setHttpConnectorListener(httpServerListener);
        if (logger.isDebugEnabled()) {
            logger.debug("Registered event listener");
        }
    }

    @Override
    public void disconnect() {
        if (serverConnector != null) {
            serverConnector.stop();
            if (logger.isDebugEnabled()) {
                logger.debug("Successfully stopped HTTP Server Connector " + serverConnector.getConnectorID());
            }
            serverConnector = null;
        }
    }

    @Override
    public void destroy() {
        if (httpWsConnectorFactory != null) {
            try {
                httpWsConnectorFactory.shutdown();
                if (logger.isDebugEnabled()) {
                    logger.debug("Successfully shutdown HTTP Connector Factory");
                }
                httpWsConnectorFactory = null;
            } catch (InterruptedException e) {
                logger.error("Failed to shutdown HTTP Connector Factory", e);
            }
        }
    }

    @Override
    public void pause() {
        // Do nothing
    }

    @Override
    public void resume() {
        // Do nothing
    }

    @Override
    public Map<String, Object> currentState() {
        return null;    // Do nothing
    }

    @Override
    public void restoreState(Map<String, Object> state) {
        // Do nothing
    }

    @Override
    public Class[] getOutputEventClasses() {
        return new Class[]{Map.class};
    }

    /**
     * HTTP Server Listener used for listening to HTTP Request.
     * This notifies the source event listener.
     */
    public static class HttpServerListener implements HttpConnectorListener {

        private SourceEventListener sourceEventListener;

        HttpServerListener(SourceEventListener sourceEventListener) {
            this.sourceEventListener = sourceEventListener;
        }

        @Override
        public void onMessage(HTTPCarbonMessage httpMessage) {
            ByteBuf inputByteBuf = httpMessage.getBlockingEntityCollector().getMessageBody();
            byte[] byteArray = inputByteBuf.array();
            String contentType = httpMessage.getHeader(Constants.HTTP_CONTENT_TYPE_HEADER);

            if (logger.isDebugEnabled()) {
                logger.debug("Received message of type " + contentType);
            }

            // Decoding the Zipkin spans
            List<ZipkinSpan> spans = null;
            if (Objects.equals(contentType, Constants.HTTP_APPLICATION_THRIFT_CONTENT_TYPE)) {
                try {
                    // Decoding Thrift encoded Zipkin spans
                    spans = Codec.decodeThriftData(byteArray);
                    if (logger.isDebugEnabled()) {
                        logger.debug("Decoded " + spans.size() + " Thrift encoded Zipkin Spans");
                    }
                } catch (TException e) {
                    logger.error("Failed to decode Thrift tracing data", e);
                }
            } else {
                // Decoding Zipkin spans (encoding is automatically detected)
                spans = Codec.decodeData(byteArray);
                if (logger.isDebugEnabled()) {
                    logger.debug("Decoded " + spans.size() + " Thrift encoded Zipkin Spans");
                }
            }

            if (spans != null) {
                // Sending the received Zipkin spans to the source event listener
                for (ZipkinSpan span : spans) {
                    Map<String, Object> attributes = new HashMap<>();

                    JsonArray tagsJsonArray = new JsonArray();
                    for (String tagKey : span.getTags().keySet()) {
                        JsonObject tagJsonObject = new JsonObject();
                        tagJsonObject.addProperty(tagKey, span.getTags().get(tagKey));
                        tagsJsonArray.add(tagJsonObject);
                    }
                    JsonObject kindTagJsonObject = new JsonObject();
                    kindTagJsonObject.addProperty(Constants.SPAN_KIND_TAG_KEY, span.getKind());
                    tagsJsonArray.add(kindTagJsonObject);

                    attributes.put(Constants.TRACE_ID, span.getTraceId());
                    attributes.put(Constants.SPAN_ID, span.getId());
                    attributes.put(Constants.PARENT_ID, span.getParentId());
                    attributes.put(Constants.NAME, span.getName());
                    attributes.put(Constants.SERVICE_NAME, span.getServiceName());
                    attributes.put(Constants.KIND, span.getKind());
                    attributes.put(Constants.TIMESTAMP, span.getTimestamp() / 1000);
                    attributes.put(Constants.DURATION, span.getDuration() / 1000);
                    attributes.put(Constants.TAGS, gson.toJson(tagsJsonArray));

                    sourceEventListener.onEvent(attributes, new String[0]);
                    if (logger.isDebugEnabled()) {
                        logger.debug("Sent span " + span.getTraceId() + "-" + span.getId()
                                + " to event source listener");
                    }
                }
            }
        }

        @Override
        public void onError(Throwable throwable) {
            logger.error("Failed to handler incoming Zipkin Spans", throwable);
        }
    }
}
