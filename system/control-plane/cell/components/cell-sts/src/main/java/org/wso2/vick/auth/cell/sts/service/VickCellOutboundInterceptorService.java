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
import com.mashape.unirest.http.HttpResponse;
import com.mashape.unirest.http.JsonNode;
import com.mashape.unirest.http.Unirest;
import com.mashape.unirest.http.exceptions.UnirestException;
import org.apache.commons.lang.StringUtils;
import org.apache.http.conn.ssl.NoopHostnameVerifier;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.ssl.SSLContextBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.sts.context.store.UserContextStore;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.ExternalAuth;
import org.wso2.vick.sts.core.VickSTSConstants;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;

/**
 * Intercepts Outbound calls from pods within the Cell.
 */
public class VickCellOutboundInterceptorService extends VickCellInterceptorService {

    private static final String STS_RESPONSE_TOKEN_PARAM = "token";
    private Logger log = LoggerFactory.getLogger(VickCellOutboundInterceptorService.class);

    public VickCellOutboundInterceptorService(UserContextStore userContextStore) throws VickCellSTSException {
        super(userContextStore);
        setHttpClientProperties();
    }

    @Override
    protected ExternalAuth.CheckResponse handleRequest(ExternalAuth.CheckRequest request) throws VickCellSTSException {
        log.info("Intercepting Sidecar Outbound call to destination:{}", getDestination(request));

        String authzHeaderInRequest = getAuthorizationHeaderValue(request);
        ExternalAuth.CheckResponse response;

        if (StringUtils.isEmpty(authzHeaderInRequest)) {
            log.debug("Authorization Header is missing in the outbound call. Injecting a JWT from at cell STS.");

            String stsToken = getStsToken(request);
            if (StringUtils.isEmpty(stsToken)) {
                log.error("No JWT token received from the STS endpoint: " + stsEndpointUrl);
            }
            response = ExternalAuth.CheckResponse.newBuilder()
                    .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
                    .setOkResponse(buildOkHttpResponse(stsToken))
                    .build();
        } else {
            log.info("Authorization Header is present in the request. Continuing without injecting a new JWT at " +
                    "cell STS.");
            response = ExternalAuth.CheckResponse.newBuilder()
                    .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
                    .build();
        }

        return response;
    }

    private String getStsToken(ExternalAuth.CheckRequest request) throws VickCellSTSException {

        try {
            // Check for a stored user context
            String requestId = getRequestId(request);
            // This is the original JWT sent to the cell gateway.
            String jwt = userContextStore.get(requestId);

            if (StringUtils.isNotBlank(jwt)) {
                log.debug("JWT token to inject into outbound call was retrieved from the user context store.");
                return jwt;
            } else {
                // TODO : decide whether we should simply fail the request here.
                log.debug("User Context JWT cannot be retrieved from the context store." +
                        "Calling the Global STS endpoint {} to get a JWT.", stsEndpointUrl);
                HttpResponse<JsonNode> apiResponse =
                        Unirest.post(stsEndpointUrl)
                                .basicAuth(userName, password)
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
                }
            }
        } catch (UnirestException e) {
            log.error("Error while obtaining the STS token.", e);
        }

        return null;
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
