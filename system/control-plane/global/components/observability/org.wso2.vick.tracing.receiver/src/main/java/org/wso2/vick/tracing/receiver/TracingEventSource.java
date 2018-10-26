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
import com.sun.net.httpserver.HttpContext;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import org.apache.commons.io.IOUtils;
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
import org.wso2.vick.tracing.receiver.internal.Codec;
import org.wso2.vick.tracing.receiver.internal.ZipkinSpan;

import java.io.IOException;
import java.io.InputStream;
import java.net.InetSocketAddress;
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

    private HttpServer httpServer;
    private HttpServerListener httpServerListener;

    private String host;
    private String apiContext;
    private int port;

    @Override
    public void init(SourceEventListener sourceEventListener, OptionHolder optionHolder,
                     String[] requestedTransportPropertyNames, ConfigReader configReader,
                     SiddhiAppContext siddhiAppContext) {
        host = optionHolder.validateAndGetStaticValue(Constants.TRACING_RECEIVER_HOST_KEY,
                Constants.DEFAULT_TRACING_RECEIVER_IP);
        port = Integer.parseInt(optionHolder.validateAndGetStaticValue(Constants.TRACING_RECEIVER_PORT_KEY,
                Constants.DEFAULT_TRACING_RECEIVER_PORT));
        apiContext = optionHolder.validateAndGetStaticValue(Constants.TRACING_RECEIVER_API_CONTEXT_KEY,
                Constants.DEFAULT_TRACING_RECEIVER_API_CONTEXT);

        httpServerListener = new HttpServerListener(sourceEventListener);
    }

    @Override
    public void connect(ConnectionCallback connectionCallback) throws ConnectionUnavailableException {
        try {
            httpServer = HttpServer.create(new InetSocketAddress(host, port), 0);
        } catch (IOException e) {
            throw new ConnectionUnavailableException("Failed to instantiate HTTP Server");
        }
        HttpContext context = httpServer.createContext(apiContext);
        context.setHandler(httpServerListener);
        httpServer.start();
        if (logger.isDebugEnabled()) {
            logger.debug("Started HTTP Server started and receiving requests on http://" + host + ":" + port
                    + apiContext);
        }
    }

    @Override
    public void disconnect() {
        if (httpServer != null) {
            httpServer.stop(0);
            if (logger.isDebugEnabled()) {
                logger.debug("HTTP Server Shutdown");
            }
            httpServer = null;
        }
    }

    @Override
    public void destroy() {
        // Do nothing
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
    public static class HttpServerListener implements HttpHandler {

        private SourceEventListener sourceEventListener;

        HttpServerListener(SourceEventListener sourceEventListener) {
            this.sourceEventListener = sourceEventListener;
        }

        @Override
        public void handle(HttpExchange httpExchange) throws IOException {
            InputStream inputStream = httpExchange.getRequestBody();
            byte[] byteArray = IOUtils.toByteArray(inputStream);
            String contentType = httpExchange.getRequestHeaders().getFirst(Constants.HTTP_CONTENT_TYPE_HEADER);

            handleEventReceive(byteArray, contentType);
        }

        /**
         * Handle tracing data receive.
         *
         * @param byteArray   The byte array received
         * @param contentType The content type of the message received
         */
        private void handleEventReceive(byte[] byteArray, String contentType) {
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
    }
}
