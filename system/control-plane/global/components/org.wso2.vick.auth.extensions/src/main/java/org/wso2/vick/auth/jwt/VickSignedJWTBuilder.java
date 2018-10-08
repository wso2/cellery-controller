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
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.UUID;

public class VickSignedJWTBuilder {

    private static final String TENANT_DOMAIN = MultitenantConstants.SUPER_TENANT_DOMAIN_NAME;
    private static final int TENANT_ID = MultitenantConstants.SUPER_TENANT_ID;

    private static final String MICRO_GATEWAY_DEFAULT_AUDIENCE_VALUE = "http://org.wso2.apimgt/gateway";
    private static final String VICK_STS_ISSUER_CONFIG = "Vick.STS.Issuer";
    private static final String DEFAULT_ISSUER_VALUE = "https://sts.vick.wso2.com";
    private static final String SCOPE_CLAIM = "scope";

    private JWSHeader.Builder headerBuilder = new JWSHeader.Builder(JWSAlgorithm.RS256);
    private JWTClaimsSet.Builder claimSetBuilder = new JWTClaimsSet.Builder();

    private long expiryInSeconds = 500L;
    private List<String> audience;

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
                .audience(audience);
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
            return audience;
        }
    }
}
