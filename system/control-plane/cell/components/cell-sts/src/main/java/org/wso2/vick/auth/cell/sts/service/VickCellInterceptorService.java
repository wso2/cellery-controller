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
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.auth.cell.sts.service;

import com.google.rpc.Code;
import com.google.rpc.Status;
import io.grpc.stub.StreamObserver;
import org.apache.commons.lang.StringUtils;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.wso2.vick.auth.cell.sts.context.store.UserContextStore;
import org.wso2.vick.auth.cell.sts.generated.envoy.core.Base;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.AuthorizationGrpc;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.ExternalAuth;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Collections;
import java.util.Map;

/**
 * Intercepts inbound/outbound calls among sidecars within and out of the cells.
 * <p>
 * Inbound calls are intercepted to inject user attributes are headers to be consumed by services within the cell.
 * Outbound calls are intercepted to inject authorization token required for authentication.
 */
public abstract class VickCellInterceptorService extends AuthorizationGrpc.AuthorizationImplBase {

    private static final Logger log = LoggerFactory.getLogger(VickCellInterceptorService.class);
    private static final String AUTHORIZATION_HEADER_NAME = "authorization";
    private static final String CELL_NAME_ENV_VARIABLE = "CELL_NAME";
    private static final String STS_CONFIG_PATH_ENV_VARIABLE = "CONF_PATH";
    private static final String CONFIG_FILE_PATH = "/etc/config/sts.json";

    private static final String CONFIG_STS_ENDPOINT = "endpoint";
    private static final String CONFIG_AUTH_USERNAME = "username";
    private static final String CONFIG_AUTH_PASSWORD = "password";
    private static final String BEARER_HEADER_VALUE_PREFIX = "Bearer ";

    private static final String REQUEST_ID = "request.id";
    private static final String CELL_NAME = "cell.name";
    private static final String REQUEST_ID_HEADER = "x-request-id";
    private static final String DESTINATION_HEADER = ":authority";

    protected String stsEndpointUrl;
    protected String userName;
    protected String password;
    protected String cellName;

    protected UserContextStore userContextStore;

    public VickCellInterceptorService(UserContextStore userContextStore) throws VickCellSTSException {

        setUpConfigurationParams();
        this.userContextStore = userContextStore;
    }

    private void setUpConfigurationParams() throws VickCellSTSException {

        try {
            String configFilePath = getConfigFilePath();
            String content = new String(Files.readAllBytes(Paths.get(configFilePath)));
            JSONObject config = (JSONObject) new JSONParser().parse(content);
            stsEndpointUrl = (String) config.get(CONFIG_STS_ENDPOINT);
            userName = (String) config.get(CONFIG_AUTH_USERNAME);
            password = (String) config.get(CONFIG_AUTH_PASSWORD);
            cellName = getMyCellName();

            log.info("Global STS Endpoint: " + stsEndpointUrl);
            log.info("Cell Name: " + cellName);
        } catch (ParseException | IOException e) {
            throw new VickCellSTSException("Error while setting up STS configurations", e);
        }
    }

    private String getConfigFilePath() {

        String configPath = System.getenv(STS_CONFIG_PATH_ENV_VARIABLE);
        return StringUtils.isNotBlank(configPath) ? configPath : CONFIG_FILE_PATH;
    }

    @Override
    public void check(ExternalAuth.CheckRequest request, StreamObserver<ExternalAuth.CheckResponse> responseObserver) {

        ExternalAuth.CheckResponse response;
        try {
            // Add request ID for log correlation.
            MDC.put(REQUEST_ID, getRequestId(request));
            // Add cell name to log entries
            MDC.put(CELL_NAME, cellName);

            String destination = getDestination(request);
            log.debug("Request from Istio-Proxy (destination:{}):\n{}", destination, request);
            response = handleRequest(request);
            log.debug("Response to istio-proxy (destination:{}):\n{}", destination, response);
        } catch (VickCellSTSException e) {
            log.error("Error while handling request from istio-proxy to (destination:{})", getDestination(request), e);
            response = buildErrorResponse();
        } finally {
            MDC.clear();
        }

        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    protected abstract ExternalAuth.CheckResponse handleRequest(ExternalAuth.CheckRequest request) throws VickCellSTSException;

    protected ExternalAuth.CheckResponse buildErrorResponse() {
        return ExternalAuth.CheckResponse.newBuilder()
                .setStatus(Status.newBuilder().setCode(Code.PERMISSION_DENIED_VALUE).build())
                .build();
    }

    protected ExternalAuth.OkHttpResponse buildOkHttpResponse(String stsToken) {

        return buildOkHttpResponseWithHeaders(
                Collections.singletonMap(AUTHORIZATION_HEADER_NAME, BEARER_HEADER_VALUE_PREFIX + stsToken));
    }

    protected ExternalAuth.OkHttpResponse buildOkHttpResponseWithHeaders(Map<String, String> headers) {

        ExternalAuth.OkHttpResponse.Builder builder = ExternalAuth.OkHttpResponse.newBuilder();
        headers.forEach((key, value) -> builder.addHeaders(buildHeader(key, value)));
        return builder.build();
    }


    private String getMyCellName() throws VickCellSTSException {

        // For now we pick the cell name from the environment variable. In future we need to figure out a way to derive
        // values from the authz request.
        String cellName = System.getenv(CELL_NAME_ENV_VARIABLE);
        if (StringUtils.isBlank(cellName)) {
            throw new VickCellSTSException("Environment variable '" + CELL_NAME_ENV_VARIABLE + "' is empty.");
        }
        return cellName;
    }

    private Base.HeaderValueOption buildHeader(String headerName, String headerValue) {

        return Base.HeaderValueOption.newBuilder()
                .setHeader(Base.HeaderValue.newBuilder().setKey(headerName).setValue(headerValue))
                .build();
    }

    protected String getAuthorizationHeaderValue(ExternalAuth.CheckRequest request) {

        return request.getAttributes().getRequest().getHttp().getHeaders().get(AUTHORIZATION_HEADER_NAME);
    }

    protected String getRequestId(ExternalAuth.CheckRequest request) throws VickCellSTSException {

        String id = request.getAttributes().getRequest().getHttp().getHeaders().get(REQUEST_ID_HEADER);
        if (StringUtils.isBlank(id)) {
            throw new VickCellSTSException("Request Id cannot be found in the header: " + REQUEST_ID_HEADER);
        }
        return id;
    }

    protected String getDestination(ExternalAuth.CheckRequest request) {

        String destination = request.getAttributes().getRequest().getHttp().getHeaders().get(DESTINATION_HEADER);
        if (StringUtils.isBlank(destination)) {
            destination = getHost(request);
            log.debug("Destination is picked from host value in the request.");
        }
        return destination;
    }

    private String getHost(ExternalAuth.CheckRequest request) {

        return request.getAttributes().getRequest().getHttp().getHost();
    }

}
