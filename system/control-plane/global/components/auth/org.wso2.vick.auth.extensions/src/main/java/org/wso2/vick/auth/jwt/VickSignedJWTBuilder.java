/*
 *  Copyright (c) 2018 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 */

package org.wso2.vick.auth.jwt;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.JWSSigner;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.util.Base64URL;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.collections.CollectionUtils;
import org.apache.commons.lang.StringUtils;
import org.wso2.carbon.identity.core.util.IdentityUtil;
import org.wso2.carbon.identity.oauth2.IdentityOAuth2Exception;
import org.wso2.carbon.identity.oauth2.util.OAuth2Util;
import org.wso2.carbon.utils.multitenancy.MultitenantConstants;
import org.wso2.vick.auth.exception.VickAuthException;

import java.security.Key;
import java.security.interfaces.RSAPrivateKey;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

/**
 * The JWT token builder for VICK.
 */
public class VickSignedJWTBuilder {

    private static final String TENANT_DOMAIN = MultitenantConstants.SUPER_TENANT_DOMAIN_NAME;
    private static final int TENANT_ID = MultitenantConstants.SUPER_TENANT_ID;

    private static final String MICRO_GATEWAY_DEFAULT_AUDIENCE_VALUE = "http://org.wso2.apimgt/gateway";
    private static final String VICK_STS_ISSUER_CONFIG = "Vick.STS.Issuer";
    private static final String DEFAULT_ISSUER_VALUE = "https://sts.vick.wso2.com";
    private static final String SCOPE_CLAIM = "scope";
    private static final String KEY_TYPE_CLAIM = "keytype";
    private static final String PRODUCTION_KEY_TYPE = "PRODUCTION";

    private JWSHeader.Builder headerBuilder = new JWSHeader.Builder(JWSAlgorithm.RS256);
    private JWTClaimsSet.Builder claimSetBuilder = new JWTClaimsSet.Builder();

    // By default we set the expiry to be 20mins (So that it will be greater than GatewayCache timeout)
    private long expiryInSeconds = 1200L;
    private List<String> audience = new ArrayList<>();

    public VickSignedJWTBuilder subject(String subject) {

        claimSetBuilder.subject(subject);
        return this;
    }

    public VickSignedJWTBuilder claim(String name, Object value) {

        claimSetBuilder.claim(name, value);
        return this;
    }

    public VickSignedJWTBuilder claims(Map<String, Object> customClaims) {

        customClaims.forEach((x, y) -> claimSetBuilder.claim(x, y));
        return this;
    }

    public VickSignedJWTBuilder scopes(List<String> scopes) {

        return claim(SCOPE_CLAIM, scopes);
    }

    public VickSignedJWTBuilder expiryInSeconds(long expiryInSeconds) {

        this.expiryInSeconds = expiryInSeconds;
        return this;
    }

    public VickSignedJWTBuilder audience(List<String> audience) {

        this.audience = audience;
        return this;
    }

    public VickSignedJWTBuilder audience(String audience) {

        this.audience.add(audience);
        return this;
    }

    private JWSHeader buildJWSHeader() throws IdentityOAuth2Exception {

        String certThumbPrint = OAuth2Util.getThumbPrint(TENANT_DOMAIN, TENANT_ID);
        headerBuilder.keyID(certThumbPrint);
        headerBuilder.x509CertThumbprint(new Base64URL(certThumbPrint));
        return headerBuilder.build();
    }

    public String build() throws VickAuthException {

        // Build the JWT Header
        try {
            JWSHeader jwsHeader = buildJWSHeader();
            // Add mandatory claims
            addMandatoryClaims(claimSetBuilder);
            JWTClaimsSet claimsSet = this.claimSetBuilder.build();

            SignedJWT signedJWT = new SignedJWT(jwsHeader, claimsSet);
            JWSSigner signer = new RSASSASigner(getRSASigningKey());

            signedJWT.sign(signer);
            return signedJWT.serialize();
        } catch (IdentityOAuth2Exception | JOSEException e) {
            throw new VickAuthException("Error while generating the signed JWT.", e);
        }
    }

    private void addMandatoryClaims(JWTClaimsSet.Builder claimsSet) {

        Date issuedAt = new Date(System.currentTimeMillis());
        Date expiryTime = new Date(issuedAt.getTime() + expiryInSeconds * 1000);

        List<String> audience = getAudience(this.audience);

        claimsSet.jwtID(UUID.randomUUID().toString())
                .issuer(getIssuer())
                .issueTime(issuedAt)
                .expirationTime(expiryTime)
                .audience(audience)
                .claim(KEY_TYPE_CLAIM, PRODUCTION_KEY_TYPE);
    }

    private String getIssuer() {

        String issuer = IdentityUtil.getProperty(VICK_STS_ISSUER_CONFIG);
        if (StringUtils.isEmpty(issuer)) {
            issuer = DEFAULT_ISSUER_VALUE;
        }
        return issuer;
    }

    private RSAPrivateKey getRSASigningKey() throws IdentityOAuth2Exception {

        Key privateKey = OAuth2Util.getPrivateKey(TENANT_DOMAIN, TENANT_ID);
        return (RSAPrivateKey) privateKey;
    }

    private List<String> getAudience(List<String> audience) {

        if (CollectionUtils.isEmpty(audience)) {
            return Collections.singletonList(MICRO_GATEWAY_DEFAULT_AUDIENCE_VALUE);
        } else {
            return audience.stream()
                    .filter(StringUtils::isNotBlank)
                    .collect(Collectors.toList());
        }
    }
}
