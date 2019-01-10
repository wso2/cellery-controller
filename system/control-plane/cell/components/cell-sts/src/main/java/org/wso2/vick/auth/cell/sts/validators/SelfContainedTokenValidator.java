/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
 */

package org.wso2.vick.auth.cell.sts.validators;

import com.nimbusds.jwt.JWT;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.sts.CellStsUtils;
import org.wso2.vick.auth.cell.sts.exception.TokenValidationFailureException;
import org.wso2.vick.auth.cell.sts.model.CellStsRequest;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;
import org.wso2.vick.auth.cell.sts.service.VickCellStsService;

import java.text.ParseException;
import java.util.Date;
import java.util.Optional;

/**
 * Validate self contained acess tokens.
 */
public class SelfContainedTokenValidator implements TokenValidator {

    private JWTSignatureValidator jwtValidator = new JWKSBasedJWTValidator();
    private static final Logger log = LoggerFactory.getLogger(SelfContainedTokenValidator.class);
    private static String globalIssuer = "https://sts.vick.wso2.com";

    /**
     * Validates a self contained access token.
     *
     * @param token          Incoming token. JWT to be validated.
     * @param cellStsRequest Request which reaches cell STS.
     * @throws TokenValidationFailureException TokenValidationFailureException.
     */
    @Override
    public void validateToken(String token, CellStsRequest cellStsRequest) throws TokenValidationFailureException {

        if (StringUtils.isEmpty(token)) {
            throw new TokenValidationFailureException("No token found in the request.");
        }
        try {
            log.debug("Validating token: {}" + token);
            SignedJWT parsedJWT = SignedJWT.parse(token);
            JWTClaimsSet jwtClaimsSet = parsedJWT.getJWTClaimsSet();
            validateIssuer(jwtClaimsSet, cellStsRequest);
            validateAudience(jwtClaimsSet, cellStsRequest);
            ValidateExpiry(jwtClaimsSet);
            validateSignature(parsedJWT, cellStsRequest);
        } catch (ParseException e) {
            throw new TokenValidationFailureException("Error while parsing JWT: " + token, e);
        }
    }

    private void ValidateExpiry(JWTClaimsSet jwtClaimsSet) throws TokenValidationFailureException {

        if (jwtClaimsSet.getExpirationTime().before(new Date(System.currentTimeMillis()))) {
            throw new TokenValidationFailureException("Token has expired. Expiry time: " + jwtClaimsSet
                    .getExpirationTime());
        }
        log.debug("Token life time is valid, expiry time: {}" + jwtClaimsSet.getExpirationTime());
    }

    private void validateAudience(JWTClaimsSet jwtClaimsSet, CellStsRequest cellStsRequest) throws
            TokenValidationFailureException {

        if (jwtClaimsSet.getAudience().isEmpty()) {
            throw new TokenValidationFailureException("No audiences found in the token");
        }

        try {
            String cellAudience = CellStsUtils.getMyCellName();
            Optional<String> audienceMatch = jwtClaimsSet.getAudience().stream().filter(audience ->
                    audience.equalsIgnoreCase(cellAudience)).findAny();
            if (!audienceMatch.isPresent()) {
                throw new TokenValidationFailureException("Error while validating audience. Expected audience :" +
                        cellAudience);
            }
            log.debug("Audience validation successful");
        } catch (VickCellSTSException e) {
            throw new TokenValidationFailureException("Cannot infer cell name", e);
        }

    }

    private void validateIssuer(JWTClaimsSet claimsSet, CellStsRequest request) throws TokenValidationFailureException {

        String issuer = globalIssuer;
        if (StringUtils.isNotEmpty(request.getSource().getCellName())) {
            issuer = CellStsUtils.getIssuerName(request.getSource().getCellName());
        }
        String issuerInToken = claimsSet.getIssuer();
        if (StringUtils.isEmpty(issuerInToken)) {
            throw new TokenValidationFailureException("No issuer found in the JWT");
        }
        if (!StringUtils.equalsIgnoreCase(issuerInToken, issuer)) {
            throw new TokenValidationFailureException("Issuer validation failed. Expected issuer : " + issuer + ". " +
                    "Received issuer: " + issuerInToken);
        }
        log.debug("Issuer validated successfully. Issuer : {}" + issuer);
    }

    private void validateSignature(JWT jwt, CellStsRequest cellStsRequest) throws TokenValidationFailureException {

        String jwkEndpoint = VickCellStsService.getCellStsConfiguration().getGlobalJWKEndpoint();

        if (StringUtils.isNotEmpty(cellStsRequest.getSource().getCellName())) {
            int port = resolvePort(cellStsRequest.getSource().getCellName());
            jwkEndpoint = "http://" + CellStsUtils.getIssuerName(cellStsRequest.getSource().getCellName()) + ":" + port;
        }

        log.debug("Calling JWKS endpoint: " + jwkEndpoint);
        try {
            log.debug("Validating signature of the token");
            jwtValidator.validateSignature(jwt, jwkEndpoint, jwt.getHeader().getAlgorithm().getName(), null);
        } catch (TokenValidationFailureException e) {
            throw new TokenValidationFailureException("Error while validating signature of the token", e);
        }
        log.debug("Token signature validated successfully");
    }

    private int resolvePort(String cellName) {

        int port = 8090;
        // Keep this commented code for the easiness of testing.
//        switch (cellName) {
//            case "hr":
//                port = 8090;
//                break;
//            case "employee":
//                port = 8091;
//                break;
//            case "stock-options":
//                port = 8092;
//                break;
//        }
        return port;
    }
}
