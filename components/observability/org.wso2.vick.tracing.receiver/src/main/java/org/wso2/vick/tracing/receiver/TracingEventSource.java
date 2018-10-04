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
import org.apache.thrift.TException;
import org.wso2.siddhi.annotation.Example;
import org.wso2.siddhi.annotation.Extension;
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
import java.util.logging.Logger;

/**
 * This class implements the event source, where the received telemetry attributes can be injected to streams.
 */
@Extension(
        name = "tracing-receiver",
        namespace = "source",
        description = "Tracing Receiver for VICK",
        examples = {
                @Example(
                        syntax = "this is synatax",
                        description = "some desc"
                )
        }
)
public class TracingEventSource extends Source {

    private static final Logger logger = Logger.getLogger(TracingEventSource.class.getName());
    private static final Gson gson = new Gson();

    private HttpServer server;
    private SourceEventListener sourceEventListener;

    @Override
    public void init(SourceEventListener sourceEventListener, OptionHolder optionHolder,
                     String[] requestedTransportPropertyNames, ConfigReader configReader,
                     SiddhiAppContext siddhiAppContext) {
        this.sourceEventListener = sourceEventListener;
    }

    @Override
    public void connect(ConnectionCallback connectionCallback) throws ConnectionUnavailableException {
        try {
            server = HttpServer.create(new InetSocketAddress("0.0.0.0", 9411), 0);
        } catch (IOException e) {
            throw new ConnectionUnavailableException("Failed to initialize HTTP Server", e);
        }
        HttpContext context = server.createContext("/api/v1/spans");
        context.setHandler(new HttpServerHandler(sourceEventListener));     // Handles tracing data receiving
        server.start();
    }

    @Override
    public void disconnect() {
        if (server != null) {
            server.stop(0);
            server = null;
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
     * Http Server handler for requests to api/v1/spans context.
     */
    public static class HttpServerHandler implements HttpHandler {

        private SourceEventListener sourceEventListener;

        HttpServerHandler(SourceEventListener sourceEventListener) {
            this.sourceEventListener = sourceEventListener;
        }

        @Override
        public void handle(HttpExchange httpExchange) throws IOException {
            InputStream inputStream = httpExchange.getRequestBody();
            byte[] byteArray = IOUtils.toByteArray(inputStream);

            String contentType = httpExchange.getRequestHeaders().getFirst(Constants.HTTP_CONTENT_TYPE_HEADER);

            List<ZipkinSpan> spans;
            if (Objects.equals(contentType, Constants.HTTP_APPLICATION_THRIFT_CONTENT_TYPE)) {
                try {
                    spans = Codec.decodeThriftData(byteArray);
                } catch (TException e) {
                    throw new IOException("Failed to decode thrift tracing data", e);
                }
            } else {
                spans = Codec.decodeData(byteArray);
            }

            for (ZipkinSpan span : spans) {
                Map<String, Object> attributes = new HashMap<>();

                JsonArray tagsJsonArray = new JsonArray();
                for (String tagKey : span.getTags().keySet()) {
                    JsonObject tagJsonObject = new JsonObject();
                    tagJsonObject.addProperty(tagKey, span.getTags().get(tagKey));
                    tagsJsonArray.add(tagJsonObject);
                }

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
            }
        }
    }
}
