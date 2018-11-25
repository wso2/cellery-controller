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

package org.wso2.vick.auth.util;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSVerifier;
import com.nimbusds.jose.crypto.RSASSAVerifier;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.carbon.identity.application.common.model.IdentityProvider;
import org.wso2.carbon.identity.application.common.util.IdentityApplicationManagementUtil;
import org.wso2.carbon.identity.oauth2.IdentityOAuth2Exception;
import org.wso2.carbon.idp.mgt.IdentityProviderManagementException;
import org.wso2.carbon.idp.mgt.IdentityProviderManager;
import org.wso2.carbon.utils.multitenancy.MultitenantConstants;

import java.security.PublicKey;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;
import java.text.ParseException;
import java.util.Collections;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * This is the Utils class for token validations.
 */
public class Utils {

    private static final Log log = LogFactory.getLog(Utils.class);

    public static final String OPENID_IDP_ENTITY_ID = "IdPEntityId";
    private static final Set<String> FILTERED_CLAIMS;

    static {
        Set<String> n = new HashSet<>();
        n.add("iss");
        n.add("aud");
        n.add("exp");
        n.add("nbf");
        n.add("iat");
        n.add("jti");
        n.add("scope");
        n.add("keytype");

        FILTERED_CLAIMS = Collections.unmodifiableSet(n);
    }

    private Utils() {

    }

    public static IdentityProvider getVickIdp() throws IdentityProviderManagementException {
        return IdentityProviderManager.getInstance().getResidentIdP(MultitenantConstants.SUPER_TENANT_DOMAIN_NAME);
    }

    public static boolean isSignedJWT(String jwtToTest) {
        // Signed JWT token contains 3 base64 encoded components separated by periods.
        return StringUtils.countMatches(jwtToTest, ".") == 2;
    }

    public static Map<String, Object> getCustomClaims(SignedJWT signedJWT) throws ParseException {
        return signedJWT.getJWTClaimsSet().getClaims().entrySet()
                .stream()
                .filter(x -> !FILTERED_CLAIMS.contains(x.getKey()))
                .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    }

    public static boolean validateSignature(SignedJWT signedJWT,
                                            IdentityProvider idp) throws IdentityOAuth2Exception {

        try {
            X509Certificate x509Certificate = getCertToValidateJwt(idp);
            String signatureAlgorithm = signedJWT.getHeader().getAlgorithm().getName();
            validateSignatureAlgorithm(signatureAlgorithm);

            PublicKey publicKey = x509Certificate.getPublicKey();
            JWSVerifier verifier;
            if (publicKey instanceof RSAPublicKey) {
                verifier = new RSASSAVerifier((RSAPublicKey) publicKey);
            } else {
                throw new IdentityOAuth2Exception("Public key is not an RSA public key.");
            }
            return signedJWT.verify(verifier);
        } catch (JOSEException | CertificateException e) {
            throw new IdentityOAuth2Exception("Error while validating signature of jwt.", e);
        }

    }

    private static void validateSignatureAlgorithm(String signatureAlgorithm) throws IdentityOAuth2Exception {

        if (StringUtils.isEmpty(signatureAlgorithm)) {
            throw new IdentityOAuth2Exception("Algorithm must not be null.");

        } else {
            if (log.isDebugEnabled()) {
                log.debug("Signature Algorithm found in the Token Header: " + signatureAlgorithm);
            }
            if (!StringUtils.startsWithIgnoreCase(signatureAlgorithm, "RS")) {
                throw new IdentityOAuth2Exception("Signature validation for algorithm: " + signatureAlgorithm
                        + " is not supported.");
            }
        }
    }

    private static X509Certificate getCertToValidateJwt(IdentityProvider idp) throws IdentityOAuth2Exception,
            CertificateException {

        X509Certificate x509Certificate = (X509Certificate) IdentityApplicationManagementUtil
                .decodeCertificate(idp.getCertificate());
        if (x509Certificate == null) {
            throw new IdentityOAuth2Exception("Unable to locate certificate for Identity Provider: "
                    + idp.getDisplayName());
        }
        return x509Certificate;
    }
}
