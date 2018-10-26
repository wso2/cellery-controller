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
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.PlainJWT;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.sts.context.store.UserContextStore;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.ExternalAuth;

import java.util.HashMap;
import java.util.Map;

/**
 * Intercepts inbounds calls to pods within the Cell.
 */
public class VickCellInboundInterceptorService extends VickCellInterceptorService {

    private static final String VICK_AUTH_SUBJECT_HEADER = "x-vick-auth-subject";
    private static final String VICK_AUTH_SUBJECT_CLAIMS_HEADER = "x-vick-auth-subject-claims";

    private Logger log = LoggerFactory.getLogger(VickCellInboundInterceptorService.class);

    public VickCellInboundInterceptorService(UserContextStore userContextStore) throws VickCellSTSException {
        super(userContextStore);
    }

    @Override
    protected ExternalAuth.CheckResponse handleRequest(ExternalAuth.CheckRequest request) throws VickCellSTSException {
        log.info("Intercepting Sidecar Inbound call to destination:{}", getDestination(request));

        // Extract the requestId
        String requestId = getRequestId(request);

        JWTClaimsSet jwtClaims;
        if (userContextStore.containsKey(requestId)) {
            // We have intercepted intra cell communication here. So we load the user attributes from the cell local
            // context store.
            log.debug("User context JWT found in context store. Loading user claims using context.");
            String jwt = userContextStore.get(requestId);
            jwtClaims = getJWTClaims(jwt);
        } else {
            // We have intercepted a service call from the Cell Gateway into a service. We need to extract the user
            // claims from the JWT sent in authorization header and store it in our user context store.
            String authzHeader = getAuthorizationHeaderValue(request);
            String jwt = extractJwtFromAuthzHeader(authzHeader);
            if (StringUtils.isBlank(jwt)) {
                throw new VickCellSTSException("Cannot extract user context JWT from Authorization header.");
            }

            jwtClaims = getJWTClaims(jwt);
            // Add the JWT to the user context store
            userContextStore.put(requestId, jwt);
            log.debug("User context JWT added to context store.");
        }

        Map<String, String> headersToSet = new HashMap<>();
        headersToSet.put(VICK_AUTH_SUBJECT_HEADER, jwtClaims.getSubject());
        headersToSet.put(VICK_AUTH_SUBJECT_CLAIMS_HEADER, new PlainJWT(jwtClaims).serialize());

        return ExternalAuth.CheckResponse.newBuilder()
                .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
                .setOkResponse(buildOkHttpResponse(headersToSet))
                .build();
    }

    private String extractJwtFromAuthzHeader(String authzHeader) {
        if (StringUtils.isBlank(authzHeader)) {
            return null;
        }

        String[] split = authzHeader.split("\\s+");
        return split.length > 1 ? split[1] : null;
    }


    protected JWTClaimsSet getJWTClaims(String jwt) throws VickCellSTSException {
        try {
            return SignedJWT.parse(jwt).getJWTClaimsSet();
        } catch (java.text.ParseException e) {
            throw new VickCellSTSException("Error while parsing the Signed JWT in authorization header.", e);
        }
    }

}
