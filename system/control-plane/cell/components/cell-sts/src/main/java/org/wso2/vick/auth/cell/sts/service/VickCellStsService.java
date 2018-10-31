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

import com.mashape.unirest.http.HttpResponse;
import com.mashape.unirest.http.JsonNode;
import com.mashape.unirest.http.Unirest;
import com.mashape.unirest.http.exceptions.UnirestException;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.PlainJWT;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.apache.http.conn.ssl.NoopHostnameVerifier;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.ssl.SSLContextBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.sts.context.store.UserContextStore;
import org.wso2.vick.auth.cell.sts.model.CellStsRequest;
import org.wso2.vick.auth.cell.sts.model.CellStsResponse;
import org.wso2.vick.auth.cell.sts.model.config.CellStsConfiguration;
import org.wso2.vick.sts.core.VickSTSConstants;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.util.HashMap;
import java.util.Map;

public class VickCellStsService {

    private static final String REQUEST_ID_HEADER = "x-request-id";
    private static final String VICK_AUTH_SUBJECT_HEADER = "x-vick-auth-subject";
    private static final String VICK_AUTH_SUBJECT_CLAIMS_HEADER = "x-vick-auth-subject-claims";
    private static final String AUTHORIZATION_HEADER_NAME = "authorization";
    private static final String BEARER_HEADER_VALUE_PREFIX = "Bearer ";
    private static final String STS_RESPONSE_TOKEN_PARAM = "token";

    private static final Logger log = LoggerFactory.getLogger(VickCellStsService.class);

    private UserContextStore userContextStore;
    private CellStsConfiguration cellStsConfiguration;


    public VickCellStsService(CellStsConfiguration stsConfig,
                              UserContextStore contextStore) throws VickCellSTSException {
        userContextStore = contextStore;
        cellStsConfiguration = stsConfig;

        setHttpClientProperties();
    }

    public void handleInboundRequest(CellStsRequest cellStsRequest,
                                     CellStsResponse cellStsResponse) throws VickCellSTSException {

        // Extract the requestId
        String requestId = getRequestId(cellStsRequest);

        JWTClaimsSet jwtClaims;
        if (userContextStore.containsKey(requestId)) {
            // User context is already available in the cell local context store. Load the user context from the store.
            log.debug("User context JWT found in context store. Loading user claims using context for requestId:{}",
                    requestId);
            jwtClaims = getUserClaimsFromContextStore(requestId);
        } else {
            // User context is not available in the cell local context store. This means we have intercepted a service
            // call from the Cell Gateway into a service. We need to extract the user claims from the JWT sent in
            // authorization header and store it in our user context store.
            log.debug("User context JWT not found in context store for requestId:{}. " +
                    "Extracting the user context JWT from the authorization header", requestId);
            String authzHeaderValue = getAuthorizationHeaderValue(cellStsRequest);
            jwtClaims = extractUserClaimsFromAuthzHeader(authzHeaderValue);

            // We store the JWT sent in the authorization header against the request Id
            userContextStore.put(requestId, authzHeaderValue);
            log.debug("User context JWT added to context store for requestId:{}", requestId);
        }

        Map<String, String> headersToSet = new HashMap<>();
        headersToSet.put(VICK_AUTH_SUBJECT_HEADER, jwtClaims.getSubject());
        headersToSet.put(VICK_AUTH_SUBJECT_CLAIMS_HEADER, new PlainJWT(jwtClaims).serialize());

        cellStsResponse.addResponseHeaders(headersToSet);
    }


    public void handleOutboundRequest(CellStsRequest cellStsRequest,
                                      CellStsResponse cellStsResponse) throws VickCellSTSException {

        String authzHeaderInRequest = getAuthorizationHeaderValue(cellStsRequest);
        if (StringUtils.isEmpty(authzHeaderInRequest)) {
            log.debug("Authorization Header is missing in the outbound call. Injecting a JWT from at cell STS.");

            String stsToken = getStsToken(cellStsRequest);
            if (StringUtils.isEmpty(stsToken)) {
                throw new VickCellSTSException("No JWT token received from the STS endpoint: "
                        + cellStsConfiguration.getStsEndpoint());
            }
            // Set the authorization header
            cellStsResponse.addResponseHeader(AUTHORIZATION_HEADER_NAME, BEARER_HEADER_VALUE_PREFIX + stsToken);
        } else {
            log.info("Authorization Header is present in the request. Continuing without injecting a new JWT at " +
                    "cell STS.");
        }
    }

    private String getAuthorizationHeaderValue(CellStsRequest request) {

        return request.getRequestHeaders().get(AUTHORIZATION_HEADER_NAME);
    }

    private String getRequestId(CellStsRequest cellStsRequest) throws VickCellSTSException {

        String id = cellStsRequest.getRequestHeaders().get(REQUEST_ID_HEADER);
        if (StringUtils.isBlank(id)) {
            throw new VickCellSTSException("Request Id cannot be found in the header: " + REQUEST_ID_HEADER);
        }
        return id;
    }

    private JWTClaimsSet extractUserClaimsFromAuthzHeader(String authzHeaderValue) throws VickCellSTSException {

        String jwt = extractJwtFromAuthzHeader(authzHeaderValue);
        if (StringUtils.isBlank(jwt)) {
            throw new VickCellSTSException("Cannot extract user context JWT from Authorization header.");
        }

        return getJWTClaims(jwt);
    }

    private JWTClaimsSet getUserClaimsFromContextStore(String requestId) throws VickCellSTSException {

        String jwt = userContextStore.get(requestId);
        return getJWTClaims(jwt);
    }

    private String extractJwtFromAuthzHeader(String authzHeader) {
        if (StringUtils.isBlank(authzHeader)) {
            return null;
        }

        String[] split = authzHeader.split("\\s+");
        return split.length > 1 ? split[1] : null;
    }


    private JWTClaimsSet getJWTClaims(String jwt) throws VickCellSTSException {
        try {
            return SignedJWT.parse(jwt).getJWTClaimsSet();
        } catch (java.text.ParseException e) {
            throw new VickCellSTSException("Error while parsing the Signed JWT in authorization header.", e);
        }
    }

    private String getStsToken(CellStsRequest request) throws VickCellSTSException {

        try {
            // Check for a stored user context
            String requestId = getRequestId(request);
            // This is the original JWT sent to the cell gateway.
            String jwt = userContextStore.get(requestId);

            if (StringUtils.isNotBlank(jwt)) {
                log.debug("JWT token to inject into outbound call was retrieved from the user context store.");
                return jwt;
            } else {
                // Most probably this is the case where a request has been initiated from within a Cell without a
                // user interaction.
                log.debug("User Context JWT cannot be retrieved from the context store." +
                        "Calling the Global STS endpoint {} to get a JWT.", cellStsConfiguration.getStsEndpoint());
                return getTokenFromGlobalSTS();
            }
        } catch (UnirestException e) {
            throw new VickCellSTSException("Error while obtaining the STS token.", e);
        }
    }

    private String getTokenFromGlobalSTS() throws UnirestException {

        String stsEndpointUrl = cellStsConfiguration.getStsEndpoint();
        String username = cellStsConfiguration.getUsername();
        String password = cellStsConfiguration.getPassword();
        String cellName = cellStsConfiguration.getCellName();

        HttpResponse<JsonNode> apiResponse =
                Unirest.post(stsEndpointUrl)
                        .basicAuth(username, password)
                        .field(VickSTSConstants.VickSTSRequest.SUBJECT, cellName)
                        .asJson();

        log.debug("Response from the STS:\nstatus:{}\nbody:{}",
                apiResponse.getStatus(), apiResponse.getBody().toString());

        if (apiResponse.getStatus() == 200) {
            Object stsTokenValue = apiResponse.getBody().getObject().get(STS_RESPONSE_TOKEN_PARAM);
            return stsTokenValue != null ? stsTokenValue.toString() : null;
        } else {
            log.error("Error from STS endpoint. statusCode= " + apiResponse.getStatus() + ", " +
                    "statusMessage=" + apiResponse.getStatusText());
            return null;
        }
    }

    private void setHttpClientProperties() throws VickCellSTSException {
        try {
            // TODO add the correct certs for hostname verification..
            Unirest.setHttpClient(HttpClients.custom()
                    .setSSLContext(new SSLContextBuilder().loadTrustMaterial(null, (x509Certificates, s) -> true).build())
                    .setSSLHostnameVerifier(NoopHostnameVerifier.INSTANCE)
                    .disableRedirectHandling()
                    .build());
        } catch (NoSuchAlgorithmException | KeyManagementException | KeyStoreException e) {
            throw new VickCellSTSException("Error initializing the http client.", e);
        }
    }
}
