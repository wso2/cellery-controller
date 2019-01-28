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

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.rpc.Code;
import com.google.rpc.Status;
import io.grpc.stub.StreamObserver;
import org.apache.commons.lang.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.wso2.vick.auth.cell.sts.CellStsUtils;
import org.wso2.vick.auth.cell.sts.generated.envoy.core.Base;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.AttributeContextOuterClass;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.AuthorizationGrpc;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.ExternalAuth;
import org.wso2.vick.auth.cell.sts.generated.istio.mixer.v1.AttributesOuterClass;
import org.wso2.vick.auth.cell.sts.model.CellStsRequest;
import org.wso2.vick.auth.cell.sts.model.CellStsResponse;
import org.wso2.vick.auth.cell.sts.model.RequestContext;
import org.wso2.vick.auth.cell.sts.model.RequestDestination;
import org.wso2.vick.auth.cell.sts.model.RequestSource;

import java.util.Base64;
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

    private static final String REQUEST_ID = "request.id";
    private static final String CELL_NAME = "cell.name";
    private static final String REQUEST_ID_HEADER = "x-request-id";
    private static final String DESTINATION_HEADER = ":authority";
    private static final String CELL_NAME_ENV_VARIABLE = "CELL_NAME";
    private static final String ISTIO_ATTRIBUTES_HEADER = "x-istio-attributes";
    private static final String ISTIO_INGRESS_PREFIX = "istio-ingressgateway";

    protected VickCellStsService cellStsService;

    public VickCellInterceptorService(VickCellStsService cellStsService) throws VickCellSTSException {

        this.cellStsService = cellStsService;
    }

    @Override
    public final void check(ExternalAuth.CheckRequest requestFromProxy,
                            StreamObserver<ExternalAuth.CheckResponse> responseObserver) {

        ExternalAuth.CheckResponse responseToProxy;
        try {
            // Add request ID for log correlation.
            MDC.put(REQUEST_ID, getRequestId(requestFromProxy));
            // Add cell name to log entries
            MDC.put(CELL_NAME, getMyCellName());

            String destination = getDestination(requestFromProxy);
            log.debug("Request from Istio-Proxy (destination:{}):\n{}", destination, requestFromProxy);

            // Build Cell STS request from the Envoy Proxy Check Request
            CellStsRequest cellStsRequest = buildCellStsRequest(requestFromProxy);

            AttributesOuterClass.Attributes attributesFromRequest = getAttributesFromRequest(requestFromProxy);
            Map<String, AttributesOuterClass.Attributes.AttributeValue> attributesMap;
            if (attributesFromRequest != null) {
                attributesMap = attributesFromRequest.getAttributesMap();
            } else {
                attributesMap = Collections.emptyMap();
            }

            log.debug("Request Attributes: \n" + attributesMap);
            CellStsResponse cellStsResponse = new CellStsResponse();

            // Let the request be handled by inbound/outbound interceptors
            handleRequest(cellStsRequest, cellStsResponse);

            // Build the response to envoy proxy response from Cell STS Response
            responseToProxy = ExternalAuth.CheckResponse.newBuilder()
                    .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
                    .setOkResponse(buildOkHttpResponseWithHeaders(cellStsResponse.getResponseHeaders()))
                    .build();

            log.debug("Response to istio-proxy (destination:{}):\n{}", destination, responseToProxy);
        } catch (VickCellSTSException e) {
            log.error("Error while handling request from istio-proxy to (destination:{})",
                    getDestination(requestFromProxy), e);
            responseToProxy = buildErrorResponse();
        } finally {
            MDC.clear();
        }

        responseObserver.onNext(responseToProxy);
        responseObserver.onCompleted();
    }

    protected abstract void handleRequest(CellStsRequest cellStsRequest,
                                          CellStsResponse cellStsResponse) throws VickCellSTSException;


    private ExternalAuth.CheckResponse buildErrorResponse() {
        return ExternalAuth.CheckResponse.newBuilder()
                .setStatus(Status.newBuilder().setCode(Code.PERMISSION_DENIED_VALUE).build())
                .build();
    }

    private ExternalAuth.OkHttpResponse buildOkHttpResponseWithHeaders(Map<String, String> headers) {

        ExternalAuth.OkHttpResponse.Builder builder = ExternalAuth.OkHttpResponse.newBuilder();
        headers.forEach((key, value) -> builder.addHeaders(buildHeader(key, value)));
        return builder.build();
    }

    private Base.HeaderValueOption buildHeader(String headerName, String headerValue) {

        return Base.HeaderValueOption.newBuilder()
                .setHeader(Base.HeaderValue.newBuilder().setKey(headerName).setValue(headerValue))
                .build();
    }

    private String getRequestId(ExternalAuth.CheckRequest request) throws VickCellSTSException {

        String id = request.getAttributes().getRequest().getHttp().getHeadersMap().get(REQUEST_ID_HEADER);
        if (StringUtils.isBlank(id)) {
            throw new VickCellSTSException("Request Id cannot be found in the header: " + REQUEST_ID_HEADER);
        }
        return id;
    }

    private String getDestination(ExternalAuth.CheckRequest request) {

        String destination = "";

        AttributesOuterClass.Attributes attributesFromRequest = getAttributesFromRequest(request);
        if (attributesFromRequest != null) {
            AttributesOuterClass.Attributes.AttributeValue attributeValue =
                    attributesFromRequest.getAttributesMap().get("destination.service.name");
            if (attributeValue != null) {
                destination = attributeValue.getStringValue();
            }
        }
        if (StringUtils.isEmpty(destination)) {
            log.debug("Destination cannot be found in request attributes.");
            destination = request.getAttributes().getRequest().getHttp().getHeadersMap().get(DESTINATION_HEADER);
            if (StringUtils.isBlank(destination)) {
                destination = request.getAttributes().getRequest().getHttp().getHost();
                log.debug("Destination is picked from host value in the request.");
            }
        }
        return destination;
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

    private CellStsRequest buildCellStsRequest(ExternalAuth.CheckRequest requestFromProxy) throws VickCellSTSException {

        return new CellStsRequest.CellStsRequestBuilder()
                .setRequestId(getRequestId(requestFromProxy))
                .setRequestHeaders(requestFromProxy.getAttributes().getRequest().getHttp().getHeaders())
                .setSource(buildRequestSource(requestFromProxy))
                .setDestination(buildRequestDestination(requestFromProxy))
                .setRequestContext(buildRequestContext(requestFromProxy))
                .build();
    }

    private AttributesOuterClass.Attributes getAttributesFromRequest(ExternalAuth.CheckRequest requestFromProxy) {

        String mixerAttributesHeaderValue = getMixerAttributesHeader(requestFromProxy);

        if (StringUtils.isNotBlank(mixerAttributesHeaderValue)) {
            byte[] decodeHeader = Base64.getDecoder().decode(mixerAttributesHeaderValue.getBytes());
            try {
                return AttributesOuterClass.Attributes.parseFrom(decodeHeader);
            } catch (InvalidProtocolBufferException e) {
                log.error("Error while trying to decode mixer attributes from '{}' header", ISTIO_ATTRIBUTES_HEADER);
            }
        }
        return null;
    }

    private String getMixerAttributesHeader(ExternalAuth.CheckRequest requestFromProxy) {
        return requestFromProxy.getAttributes().getRequest().getHttp().getHeadersMap().get(ISTIO_ATTRIBUTES_HEADER);
    }

    private RequestContext buildRequestContext(ExternalAuth.CheckRequest checkRequest) {

        AttributeContextOuterClass.AttributeContext.HttpRequest httpRequest =
                checkRequest.getAttributes().getRequest().getHttp();

        return new RequestContext()
                .setHost(httpRequest.getHost())
                .setProtocol(httpRequest.getProtocol())
                .setMethod(httpRequest.getMethod())
                .setPath(httpRequest.getPath());
    }

    private RequestSource buildRequestSource(ExternalAuth.CheckRequest checkRequest) {
        AttributesOuterClass.Attributes attributesFromRequest = getAttributesFromRequest(checkRequest);
        RequestSource.RequestSourceBuilder requestSourceBuilder = new RequestSource.RequestSourceBuilder();
        if (attributesFromRequest != null) {
            AttributesOuterClass.Attributes.AttributeValue sourceId =
                    attributesFromRequest.getAttributesMap().get("source.uid");
            if (sourceId != null) {
                // "source.uid" -> "kubernetes://hr--hr-deployment-596946948d-vvgln.default"
                String sourceUid = sourceId.getStringValue();
                String sourceWorkloadName = sourceUid.replace("kubernetes://", "");

                requestSourceBuilder.setWorkload(sourceWorkloadName)
                        .setCellName(extractCellNameFromWorkloadName(sourceWorkloadName));
            }
        }
        return requestSourceBuilder.build();
    }

    private String extractCellNameFromWorkloadName(String workloadName) {

        // When requests reaches cells through istio ingress (after mTLS), source cell is not avialable.
        if (StringUtils.isNotEmpty(workloadName) && workloadName.startsWith(ISTIO_INGRESS_PREFIX)) {
            return null;
        }
        // Workload name is in the format hr--hr-deployment-596946948d-vvgln.default where the value before --
        // is the cell name.
        return workloadName.split("--")[0];
    }

    private RequestDestination buildRequestDestination(ExternalAuth.CheckRequest checkRequest) {
        RequestDestination.RequestDestinationBuilder destinationBuilder =
                new RequestDestination.RequestDestinationBuilder();

        String destinationWorkloadName = getDestination(checkRequest);
        if (StringUtils.isNotBlank(destinationWorkloadName)) {
            destinationBuilder.setWorkload(destinationWorkloadName);
        }

        if (CellStsUtils.isWorkloadExternalToVick(destinationWorkloadName)) {
            destinationBuilder.setExternalToVick(true);
        } else {
            destinationBuilder.setCellName(extractCellNameFromWorkloadName(destinationWorkloadName));
        }

        return destinationBuilder.build();
    }
}