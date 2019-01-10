/*
 *  Copyright (c) 2019 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.vick.auth.cell.sts;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.JWSSigner;
import com.nimbusds.jose.crypto.RSASSASigner;
import com.nimbusds.jose.util.Base64URL;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.wso2.vick.auth.cell.jwks.KeyResolverException;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;
import org.wso2.vick.auth.cell.utils.CertificateUtils;

import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateEncodingException;
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
public class STSJWTBuilder {

    private static final String MICRO_GATEWAY_DEFAULT_AUDIENCE_VALUE = "http://org.wso2.apimgt/gateway";
    private static final String VICK_STS_ISSUER_CONFIG = "Vick.STS.Issuer";
    private static final String SCOPE_CLAIM = "scope";
    private static final String KEY_TYPE_CLAIM = "keytype";
    private static final String PRODUCTION_KEY_TYPE = "PRODUCTION";

    private JWSHeader.Builder headerBuilder = new JWSHeader.Builder(JWSAlgorithm.RS256);
    private JWTClaimsSet.Builder claimSetBuilder = new JWTClaimsSet.Builder();

    // By default we set this to 20m
    private long expiryInSeconds = 1200L;
    private List<String> audience = new ArrayList<>();
    private static String issuer = "https://sts.vick.wso2.com";

    public STSJWTBuilder subject(String subject) {

        claimSetBuilder.subject(subject);
        return this;
    }

    public STSJWTBuilder issuer(String issuer) {

        if (issuer != null) {
            this.issuer = issuer;
        }
        return this;
    }

    public STSJWTBuilder claim(String name, Object value) {

        claimSetBuilder.claim(name, value);
        return this;
    }

    public STSJWTBuilder claims(Map<String, Object> customClaims) {

        customClaims.forEach((x, y) -> claimSetBuilder.claim(x, y));
        return this;
    }

    public STSJWTBuilder scopes(List<String> scopes) {

        return claim(SCOPE_CLAIM, scopes);
    }

    public STSJWTBuilder expiryInSeconds(long expiryInSeconds) {

        this.expiryInSeconds = expiryInSeconds;
        return this;
    }

    public STSJWTBuilder audience(List<String> audience) {

        this.audience = audience;
        return this;
    }

    public STSJWTBuilder audience(String audience) {

        this.audience.add(audience);
        return this;
    }

    public String build() throws VickCellSTSException {

        JWSHeader jwsHeader = null;
        try {
            jwsHeader = buildJWSHeader();
        } catch (KeyResolverException | CertificateEncodingException | NoSuchAlgorithmException e) {
            throw new VickCellSTSException("Error while building JWS header", e);
        }
        // Add mandatory claims
        addMandatoryClaims(claimSetBuilder);
        JWTClaimsSet claimsSet = this.claimSetBuilder.build();

        SignedJWT signedJWT = new SignedJWT(jwsHeader, claimsSet);

        try {
            JWSSigner signer = new RSASSASigner(CertificateUtils.getKeyResolver().getPrivateKey());
            signedJWT.sign(signer);
        } catch (JOSEException | KeyResolverException e) {
            throw new VickCellSTSException("Error while signing JWT", e);
        }
        return signedJWT.serialize();
    }

    private JWSHeader buildJWSHeader() throws KeyResolverException, CertificateEncodingException,
            NoSuchAlgorithmException {

        String certThumbPrint = null;
        certThumbPrint = CertificateUtils.getThumbPrint(CertificateUtils.getKeyResolver().getCertificate());
        headerBuilder.keyID(certThumbPrint);
        headerBuilder.x509CertThumbprint(new Base64URL(certThumbPrint));
        return headerBuilder.build();
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

        return issuer;
    }

    private List<String> getAudience(List<String> audience) {

        if (audience == null || audience.isEmpty()) {
            return Collections.singletonList(MICRO_GATEWAY_DEFAULT_AUDIENCE_VALUE);
        } else {
            return audience.stream()
                    .filter(StringUtils::isNotBlank)
                    .collect(Collectors.toList());
        }
    }
}
