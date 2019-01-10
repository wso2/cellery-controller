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

import com.nimbusds.jose.jwk.source.RemoteJWKSet;
import com.nimbusds.jose.proc.SecurityContext;
import com.nimbusds.jose.util.DefaultResourceRetriever;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.vick.auth.cell.sts.exception.TokenValidationFailureException;

import java.net.MalformedURLException;
import java.net.URL;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Provides JWK sources for JWT validation.
 */
public class JWKSourceDataProvider {

    private static final int DEFAULT_HTTP_CONNECTION_TIMEOUT = 1000;
    private static final int DEFAULT_HTTP_READ_TIMEOUT = 1000;
    private static final Log log = LogFactory.getLog(JWKSourceDataProvider.class);

    private static JWKSourceDataProvider jwkSourceDataProvider = new JWKSourceDataProvider();
    private Map<String, RemoteJWKSet<SecurityContext>> jwkSourceMap = new ConcurrentHashMap<>();

    private JWKSourceDataProvider() {

    }

    /**
     * Returns an instance of JWK-Source data holder.
     *
     * @return JWKSSourceDataHolder.
     */
    public static JWKSourceDataProvider getInstance() {

        return jwkSourceDataProvider;
    }

    /**
     * Get cached JWKSet for the jwks_uri.
     *
     * @param jwksUri Identity provider's JWKS endpoint.
     * @return RemoteJWKSet.
     * @throws MalformedURLException for invalid URL.
     */
    public RemoteJWKSet<SecurityContext> getJWKSource(String jwksUri) throws MalformedURLException {

        RemoteJWKSet<SecurityContext> jwkSet = jwkSourceMap.get(jwksUri);

        if (jwkSet == null) {
            jwkSet = retrieveJWKSFromJWKSEndpoint(jwksUri);
            jwkSourceMap.put(jwksUri, jwkSet);
        }
        return jwkSet;
    }

    public Map<String, RemoteJWKSet<SecurityContext>> getJwkSourceMap() {

        return jwkSourceMap;
    }

    /**
     * Retrieve the new-keyset from the JWKS endpoint in case of signature validation failure.
     *
     * @param jwksUri Identity providers jwks_uri.
     * @throws TokenValidationFailureException for invalid/malformed URL.
     */
    public void refreshJWKSResource(String jwksUri) throws TokenValidationFailureException {

        try {
            jwkSourceMap.remove(jwksUri);
            RemoteJWKSet<SecurityContext> jwkSet = retrieveJWKSFromJWKSEndpoint(jwksUri);
            jwkSourceMap.put(jwksUri, jwkSet);
        } catch (MalformedURLException e) {
            throw new TokenValidationFailureException("Provided URI is malformed. jwks_uri: " + jwksUri, e);
        }
    }

    /**
     * Retrieve JWKS from jwks_uri.
     *
     * @param jwksUri Identity provider's jwks_uri.
     * @return RemoteJWKSet
     * @throws MalformedURLException for invalid URL.
     */
    private RemoteJWKSet<SecurityContext> retrieveJWKSFromJWKSEndpoint(String jwksUri) throws MalformedURLException {

        // Retrieve HTTP endpoint configurations.
        int connectionTimeout = DEFAULT_HTTP_CONNECTION_TIMEOUT;
        ;
        int readTimeout = DEFAULT_HTTP_READ_TIMEOUT;
        int sizeLimit = RemoteJWKSet.DEFAULT_HTTP_SIZE_LIMIT;

        DefaultResourceRetriever resourceRetriever = new DefaultResourceRetriever(
                connectionTimeout,
                readTimeout,
                sizeLimit);

        return new RemoteJWKSet<>(new URL(jwksUri), resourceRetriever);
    }
}
