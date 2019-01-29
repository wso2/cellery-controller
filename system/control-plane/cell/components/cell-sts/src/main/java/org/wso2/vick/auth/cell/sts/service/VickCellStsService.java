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

import com.mashape.unirest.http.Unirest;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.PlainJWT;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.apache.http.conn.ssl.NoopHostnameVerifier;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.ssl.SSLContextBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.authorization.AuthorizationFailedException;
import org.wso2.vick.auth.cell.authorization.AuthorizationService;
import org.wso2.vick.auth.cell.sts.CellStsUtils;
import org.wso2.vick.auth.cell.sts.Constants;
import org.wso2.vick.auth.cell.sts.STSTokenGenerator;
import org.wso2.vick.auth.cell.sts.context.store.UserContextStore;
import org.wso2.vick.auth.cell.sts.exception.CellSTSRequestValidationFailedException;
import org.wso2.vick.auth.cell.sts.exception.TokenValidationFailureException;
import org.wso2.vick.auth.cell.sts.model.CellStsRequest;
import org.wso2.vick.auth.cell.sts.model.CellStsResponse;
import org.wso2.vick.auth.cell.sts.model.RequestDestination;
import org.wso2.vick.auth.cell.sts.model.config.CellStsConfiguration;
import org.wso2.vick.auth.cell.sts.validators.CellSTSRequestValidator;
import org.wso2.vick.auth.cell.sts.validators.DefaultCellSTSReqValidator;
import org.wso2.vick.auth.cell.sts.validators.SelfContainedTokenValidator;
import org.wso2.vick.auth.cell.sts.validators.TokenValidator;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.X509Certificate;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.HttpsURLConnection;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

public class VickCellStsService {

    private static final String VICK_AUTH_SUBJECT_CLAIMS_HEADER = "x-vick-auth-subject-claims";
    private static final String AUTHORIZATION_HEADER_NAME = "authorization";
    private static final String BEARER_HEADER_VALUE_PREFIX = "Bearer ";
    private static TokenValidator tokenValidator = new SelfContainedTokenValidator();
    private static CellSTSRequestValidator requestValidator = new DefaultCellSTSReqValidator(Collections.EMPTY_LIST);
    private static AuthorizationService authorizationService = new AuthorizationService();

    private static final Logger log = LoggerFactory.getLogger(VickCellStsService.class);

    private UserContextStore userContextStore;
    private UserContextStore localContextStore;

    private static CellStsConfiguration cellStsConfiguration;

    public VickCellStsService(CellStsConfiguration stsConfig,
                              UserContextStore contextStore, UserContextStore localContextStore)
            throws VickCellSTSException {

        this.userContextStore = contextStore;
        cellStsConfiguration = stsConfig;
        this.localContextStore = localContextStore;

        setHttpClientProperties();
    }

    public void handleInboundRequest(CellStsRequest cellStsRequest,
                                     CellStsResponse cellStsResponse) throws VickCellSTSException {

        // Extract the requestId
        String requestId = cellStsRequest.getRequestId();
        JWTClaimsSet jwtClaims;
        String jwt;

        try {
            boolean authenticationRequired = requestValidator.isAuthenticationRequired(cellStsRequest);
            if (!authenticationRequired) {
                return;
            }
            log.debug("Authentication is required for the request ID: {} ", requestId);
        } catch (CellSTSRequestValidationFailedException e) {
            throw new VickCellSTSException("Error while evaluating authentication requirement", e);
        }

        String callerCell = cellStsRequest.getSource().getCellName();
        log.debug("Caller cell : {}", callerCell);

        jwt = getUserContextJwt(cellStsRequest);
        log.debug("Incoming JWT : " + jwt);

        if (isRequestToMicroGateway(cellStsRequest)) {
            log.debug("Request to micro-gateway intercepted");
            jwtClaims = handleRequestToMicroGW(cellStsRequest, requestId, jwt);
        } else {
            jwtClaims = handleInternalRequest(cellStsRequest, requestId, jwt);
        }
          // TODO : Integrate OPA and enable authorization.
        try {
            authorizationService.authorize(cellStsRequest, jwt);
        } catch (AuthorizationFailedException e) {
            throw new VickCellSTSException("Authorization failure", e);
        }
        Map<String, String> headersToSet = new HashMap<>();

        headersToSet.put(Constants.VICK_AUTH_SUBJECT_HEADER, jwtClaims.getSubject());
        log.debug("Set {} to: {}", Constants.VICK_AUTH_SUBJECT_HEADER, jwtClaims.getSubject());

        headersToSet.put(VICK_AUTH_SUBJECT_CLAIMS_HEADER, new PlainJWT(jwtClaims).serialize());
        log.debug("Set {} to : {}", VICK_AUTH_SUBJECT_CLAIMS_HEADER, new PlainJWT(jwtClaims).serialize());

        cellStsResponse.addResponseHeaders(headersToSet);

    }

    private JWTClaimsSet handleInternalRequest(CellStsRequest cellStsRequest, String requestId, String jwt) throws
            VickCellSTSException {

        JWTClaimsSet jwtClaims;
        log.debug("Call from a workload to workload within cell {} ; Source workload {} ; Destination workload",
                cellStsRequest.getSource().getCellName(), cellStsRequest.getSource().getWorkload(),
                cellStsRequest.getDestination().getWorkload());

        try {
            if (localContextStore.get(requestId) == null) {
                log.debug("Initial entrace to cell from gateway. No cached token found.");
                validateInboundToken(cellStsRequest, jwt);
                localContextStore.put(requestId, jwt);
            } else {
                if (!StringUtils.equalsIgnoreCase(localContextStore.get(requestId), jwt)) {
                    throw new VickCellSTSException("Intra cell STS token is tampered.");
                }
            }
            jwtClaims = extractUserClaimsFromJwt(jwt);
        } catch (TokenValidationFailureException e) {
            throw new VickCellSTSException("Error while validating locally issued token.", e);
        }
        return jwtClaims;
    }

    private JWTClaimsSet handleRequestToMicroGW(CellStsRequest cellStsRequest, String requestId, String jwt) throws
            VickCellSTSException {

        JWTClaimsSet jwtClaims;
        log.debug("Incoming request to cell gateway {} from {}", CellStsUtils.getMyCellName(),
                cellStsRequest.getSource());
        try {
            log.debug("Validating incoming JWT {}", jwt);
            validateInboundToken(cellStsRequest, jwt);
            userContextStore.put(requestId, jwt);
            jwtClaims = extractUserClaimsFromJwt(jwt);

        } catch (TokenValidationFailureException e) {
            throw new VickCellSTSException("Error while validating JWT token", e);
        }
        return jwtClaims;
    }

    private void validateInboundToken(CellStsRequest cellStsRequest, String token) throws
            TokenValidationFailureException {

        tokenValidator.validateToken(token, cellStsRequest);
    }

    private String getUserContextJwt(CellStsRequest cellStsRequest) {

        String authzHeaderValue = getAuthorizationHeaderValue(cellStsRequest);
        return extractJwtFromAuthzHeader(authzHeaderValue);
    }

    public void handleOutboundRequest(CellStsRequest cellStsRequest,
                                      CellStsResponse cellStsResponse) throws VickCellSTSException {

        // First we check whether the destination of the intercepted call is within VICK
        RequestDestination destination = cellStsRequest.getDestination();
        if (destination.isExternalToVick()) {
            // If the intercepted call is to an external workload to VICK we cannot do anything in the Cell STS.
            log.info("Intercepted an outbound call to a workload:{} outside VICK. Passing the call through.", destination);
        } else {
            log.info("Intercepted an outbound call to a workload:{} within VICK. Injecting a STS token for " +
                    "authentication and user-context sharing from Cell STS.", destination);

            String stsToken = getStsToken(cellStsRequest);
            if (StringUtils.isEmpty(stsToken)) {
                throw new VickCellSTSException("No JWT token received from the STS endpoint: "
                        + cellStsConfiguration.getStsEndpoint());
            }
            log.debug("Attaching jwt to outbound request : {}", stsToken);
            // Set the authorization header
            if (cellStsRequest.getRequestHeaders().get(Constants.VICK_AUTH_SUBJECT_HEADER) != null) {
                log.info("Found user in outgoing request");
            }
            cellStsResponse.addResponseHeader(AUTHORIZATION_HEADER_NAME, BEARER_HEADER_VALUE_PREFIX + stsToken);
        }
    }

    private String getAuthorizationHeaderValue(CellStsRequest request) {

        return request.getRequestHeaders().get(AUTHORIZATION_HEADER_NAME);
    }

    private JWTClaimsSet extractUserClaimsFromJwt(String jwt) throws VickCellSTSException {

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
            String requestId = request.getRequestId();
            // This is the original JWT sent to the cell gateway.
            String jwt;

            if (isRequestFromMicroGateway(request)) {
                log.debug("Request with ID: {} from micro gateway to {} workload of cell {}", requestId, request
                        .getDestination().getWorkload(), request.getDestination().getCellName());
                jwt = userContextStore.get(requestId);
                return getTokenFromLocalSTS(jwt, CellStsUtils.getMyCellName());
            } else if (isIntraCellCall(request) && localContextStore.get(requestId) != null) {
                log.debug("Intra cell request with ID: {} from source workload {} to destination workload {} within " +
                        "cell {}", requestId, request.getSource().getWorkload(), request.getDestination().getWorkload());
                return localContextStore.get(requestId);
            } else if (!isIntraCellCall(request) && localContextStore.get(requestId) != null) {
                jwt = localContextStore.get(requestId);
                return getTokenFromLocalSTS(jwt, request.getDestination().getCellName());
            } else {
                log.debug("Request initiated within cell {} to {}", request.getSource().getCellName(), request
                        .getDestination().toString());
                return getTokenFromLocalSTS(CellStsUtils.getMyCellName());
            }
        } finally {
            // do nothing
        }
    }

    private boolean isIntraCellCall(CellStsRequest cellStsRequest) throws VickCellSTSException {

        String currentCell = CellStsUtils.getMyCellName();
        String destinationCell = cellStsRequest.getDestination().getCellName();

        return StringUtils.equals(currentCell, destinationCell);
    }

    private boolean isRequestFromMicroGateway(CellStsRequest cellStsRequest) throws VickCellSTSException {

        String workload = cellStsRequest.getSource().getWorkload();
        if (StringUtils.isNotEmpty(workload) && workload.startsWith(CellStsUtils.getMyCellName() +
                "--gateway-deployment-")) {
            return true;
        }
        return false;
    }

    private boolean isRequestToMicroGateway(CellStsRequest cellStsRequest) throws VickCellSTSException {

        String workload = cellStsRequest.getDestination().getWorkload();
        return (StringUtils.isNotEmpty(workload) && workload.startsWith(CellStsUtils.getMyCellName() +
                "--gateway-service"));
    }

    private String getTokenFromLocalSTS(String audience) throws VickCellSTSException {

        return STSTokenGenerator.generateToken(audience, null);
    }

    private String getTokenFromLocalSTS(String jwt, String audience) throws VickCellSTSException {

        String token = STSTokenGenerator.generateToken(jwt, audience,
                CellStsUtils.getIssuerName(CellStsUtils.getMyCellName()));
        return token;
    }

    private void setHttpClientProperties() throws VickCellSTSException {

        // Create a trust manager that does not validate certificate chains
        TrustManager[] trustAllCerts = new TrustManager[]{new X509TrustManager() {
            public java.security.cert.X509Certificate[] getAcceptedIssuers() {

                return null;
            }

            public void checkClientTrusted(X509Certificate[] certs, String authType) {
                // Do nothing
            }

            public void checkServerTrusted(X509Certificate[] certs, String authType) {
                // Do nothing
            }
        }
        };

        try {
            SSLContext sc = SSLContext.getInstance("SSL");
            sc.init(null, trustAllCerts, new java.security.SecureRandom());
            HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());
        } catch (KeyManagementException | NoSuchAlgorithmException e) {
            throw new VickCellSTSException("Error while initializing SSL context");
        }

        // Create all-trusting host name verifier
        HostnameVerifier allHostsValid = new HostnameVerifier() {
            public boolean verify(String hostname, SSLSession session) {

                return true;
            }
        };

        // Install the all-trusting host verifier
        HttpsURLConnection.setDefaultHostnameVerifier(allHostsValid);

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

    public static CellStsConfiguration getCellStsConfiguration() {

        return cellStsConfiguration;
    }
}
