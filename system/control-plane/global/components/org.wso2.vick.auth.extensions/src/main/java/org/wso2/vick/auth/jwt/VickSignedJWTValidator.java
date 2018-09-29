package org.wso2.vick.auth.jwt;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSVerifier;
import com.nimbusds.jose.crypto.RSASSAVerifier;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.carbon.base.MultitenantConstants;
import org.wso2.carbon.identity.application.common.model.FederatedAuthenticatorConfig;
import org.wso2.carbon.identity.application.common.model.IdentityProvider;
import org.wso2.carbon.identity.application.common.util.IdentityApplicationConstants;
import org.wso2.carbon.identity.application.common.util.IdentityApplicationManagementUtil;
import org.wso2.carbon.identity.oauth.common.exception.InvalidOAuthClientException;
import org.wso2.carbon.identity.oauth.config.OAuthServerConfiguration;
import org.wso2.carbon.identity.oauth2.IdentityOAuth2Exception;
import org.wso2.carbon.identity.oauth2.util.OAuth2Util;
import org.wso2.carbon.identity.oauth2.validators.OAuth2JWTTokenValidator;
import org.wso2.carbon.identity.oauth2.validators.OAuth2TokenValidationMessageContext;
import org.wso2.carbon.idp.mgt.IdentityProviderManagementException;
import org.wso2.carbon.idp.mgt.IdentityProviderManager;
import org.wso2.vick.auth.util.Utils;

import java.security.PublicKey;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;
import java.text.ParseException;
import java.util.Date;
import java.util.List;

/**
 * Validates signed JWTs issued by trusted IDPs.
 */
public class VickSignedJWTValidator extends OAuth2JWTTokenValidator {

    private static final Log log = LogFactory.getLog(VickSignedJWTValidator.class);
    private static final String CONSUMER_KEY = "consumerKey";

    @Override
    public boolean validateAccessToken(OAuth2TokenValidationMessageContext validationContext) throws IdentityOAuth2Exception {

        // validate mandatory attributes
        String accessToken = getAccessTokenIdentifier(validationContext);
        try {
            SignedJWT signedJWT = SignedJWT.parse(accessToken);
            boolean signedJWTValid = isSignedJWTValid(signedJWT);
            if (signedJWTValid) {
                JWTClaimsSet claimsSet = signedJWT.getJWTClaimsSet();

                // These two properties are set to avoid token lookup from the database in the case of signed JWTs
                // issued by external IDPs.
                validationContext.addProperty(OAuth2Util.REMOTE_ACCESS_TOKEN, Boolean.TRUE.toString());
                validationContext.addProperty(OAuth2Util.JWT_ACCESS_TOKEN, Boolean.TRUE.toString());

                validationContext.addProperty(OAuth2Util.IAT,
                        String.valueOf(getTimeInSeconds(claimsSet.getIssueTime())));
                validationContext.addProperty(OAuth2Util.EXP,
                        String.valueOf(getTimeInSeconds(claimsSet.getExpirationTime())));
                validationContext.addProperty(OAuth2Util.CLIENT_ID, claimsSet.getClaim(CONSUMER_KEY));
                validationContext.addProperty(OAuth2Util.SUB, claimsSet.getSubject());
                validationContext.addProperty(OAuth2Util.SCOPE, claimsSet.getClaim(OAuth2Util.SCOPE));
                validationContext.addProperty(OAuth2Util.ISS, claimsSet.getIssuer());
                validationContext.addProperty(OAuth2Util.JTI, claimsSet.getJWTID());
            }

            return signedJWTValid;
        } catch (ParseException e) {
            throw new IdentityOAuth2Exception("Error validating signed jwt.", e);
        }
    }

    private long getTimeInSeconds(Date date) {

        return date.getTime() / 1000L;
    }

    @Override
    public boolean validateScope(OAuth2TokenValidationMessageContext messageContext) throws IdentityOAuth2Exception {

        if (Utils.isSignedJWT(getAccessTokenIdentifier(messageContext))) {
            return true;
        }
        return super.validateScope(messageContext);
    }

    private String getAccessTokenIdentifier(OAuth2TokenValidationMessageContext messageContext) {

        return messageContext.getRequestDTO().getAccessToken().getIdentifier();
    }

    private boolean isSignedJWTValid(SignedJWT signedJWT) throws IdentityOAuth2Exception {

        try {
            JWTClaimsSet claimsSet = signedJWT.getJWTClaimsSet();

            if (claimsSet == null) {
                throw new IdentityOAuth2Exception("Claim values are empty in the validated JWT.");
            } else {
                validateMandatoryJWTClaims(claimsSet);
                validateConsumerKey(claimsSet);
                validateExpiryTime(claimsSet);
                validateNotBeforeTime(claimsSet);
                validateAudience(claimsSet);

                IdentityProvider trustedIdp = getTrustedIdp(claimsSet);
                return validateSignature(signedJWT, trustedIdp);
            }
        } catch (ParseException ex) {
            throw new IdentityOAuth2Exception("Error while validating JWT.", ex);
        }
    }

    private void validateConsumerKey(JWTClaimsSet claimsSet) throws IdentityOAuth2Exception {

        String consumerKey = (String) claimsSet.getClaim(CONSUMER_KEY);
        if (StringUtils.isNotBlank(consumerKey)) {
            try {
                OAuth2Util.getAppInformationByClientId(consumerKey);
            } catch (IdentityOAuth2Exception | InvalidOAuthClientException e) {
                throw new IdentityOAuth2Exception("Invalid consumerKey. Cannot find a registered app for consumerKey: "
                        + consumerKey);
            }
        } else {
            throw new IdentityOAuth2Exception("Mandatory claim 'consumerKey' is missing in the signedJWT.");
        }
    }

    private void validateAudience(JWTClaimsSet claimsSet) throws IdentityOAuth2Exception {
        // We do not validate audience at the moment....
    }

    private void validateMandatoryJWTClaims(JWTClaimsSet claimsSet) throws IdentityOAuth2Exception {

        String subject = claimsSet.getSubject();
        List<String> audience = claimsSet.getAudience();
        String jti = claimsSet.getJWTID();
        if (StringUtils.isEmpty(claimsSet.getIssuer()) || StringUtils.isEmpty(subject) ||
                claimsSet.getExpirationTime() == null || audience == null || jti == null) {
            throw new IdentityOAuth2Exception("Mandatory fields(Issuer, Subject, Expiration time, jtl or Audience) " +
                    "are empty in the given Token.");
        }
    }

    private void validateExpiryTime(JWTClaimsSet claimsSet) throws IdentityOAuth2Exception {

        long timeStampSkewMillis = OAuthServerConfiguration.getInstance().getTimeStampSkewInSeconds() * 1000;
        long expirationTimeInMillis = claimsSet.getExpirationTime().getTime();
        long currentTimeInMillis = System.currentTimeMillis();
        if ((currentTimeInMillis + timeStampSkewMillis) > expirationTimeInMillis) {
            if (log.isDebugEnabled()) {
                log.debug("Token is expired." +
                        ", Expiration Time(ms) : " + expirationTimeInMillis +
                        ", TimeStamp Skew : " + timeStampSkewMillis +
                        ", Current Time : " + currentTimeInMillis + ". Token Rejected and validation terminated.");
            }
            throw new IdentityOAuth2Exception("Token is expired.");
        }

        if (log.isDebugEnabled()) {
            log.debug("Expiration Time(exp) of Token was validated successfully.");
        }
    }

    private void validateNotBeforeTime(JWTClaimsSet claimsSet) throws IdentityOAuth2Exception {

        Date notBeforeTime = claimsSet.getNotBeforeTime();
        if (notBeforeTime != null) {
            long timeStampSkewMillis = OAuthServerConfiguration.getInstance().getTimeStampSkewInSeconds() * 1000;
            long notBeforeTimeMillis = notBeforeTime.getTime();
            long currentTimeInMillis = System.currentTimeMillis();
            if (currentTimeInMillis + timeStampSkewMillis < notBeforeTimeMillis) {
                if (log.isDebugEnabled()) {
                    log.debug("Token is used before Not_Before_Time." +
                            ", Not Before Time(ms) : " + notBeforeTimeMillis +
                            ", TimeStamp Skew : " + timeStampSkewMillis +
                            ", Current Time : " + currentTimeInMillis + ". Token Rejected and validation terminated.");
                }
                throw new IdentityOAuth2Exception("Token is used before Not_Before_Time.");
            }
            if (log.isDebugEnabled()) {
                log.debug("Not Before Time(nbf) of Token was validated successfully.");
            }
        }
    }

    private boolean validateSignature(SignedJWT signedJWT, IdentityProvider idp) throws IdentityOAuth2Exception {

        JWSVerifier verifier;
        X509Certificate x509Certificate = getCertToValidateJwt(signedJWT, idp);

        String signatureAlgorithm = signedJWT.getHeader().getAlgorithm().getName();
        validateSignatureAlgorithm(signatureAlgorithm);

        PublicKey publicKey = x509Certificate.getPublicKey();
        if (publicKey instanceof RSAPublicKey) {
            verifier = new RSASSAVerifier((RSAPublicKey) publicKey);
        } else {
            throw new IdentityOAuth2Exception("Public key is not an RSA public key.");
        }

        try {
            return signedJWT.verify(verifier);
        } catch (JOSEException e) {
            throw new IdentityOAuth2Exception("Error while validating signature of jwt.", e);
        }

    }

    private void validateSignatureAlgorithm(String signatureAlgorithm) throws IdentityOAuth2Exception {

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

    private X509Certificate getCertToValidateJwt(SignedJWT signedJWT,
                                                 IdentityProvider idp) throws IdentityOAuth2Exception {

        X509Certificate x509Certificate = resolveSignerCertificate(signedJWT.getHeader(), idp);
        if (x509Certificate == null) {
            throw new IdentityOAuth2Exception("Unable to locate certificate for Identity Provider: "
                    + idp.getDisplayName());
        }
        return x509Certificate;
    }

    private IdentityProvider getTrustedIdp(JWTClaimsSet claimsSet) throws IdentityOAuth2Exception {

        String jwtIssuer = claimsSet.getIssuer();
        String tenantDomain = getTenantDomain(claimsSet);

        IdentityProvider identityProvider;
        try {
            identityProvider = IdentityProviderManager.getInstance().getIdPByName(jwtIssuer, tenantDomain);
            if (identityProvider != null) {
                // if no IDPs were found for a given name, the IdentityProviderManager returns a dummy IDP with the
                // name "default". We need to handle this case.
                if (StringUtils.equalsIgnoreCase(identityProvider.getIdentityProviderName(), "default")) {
                    // Check whether this jwt was issued by our local idp
                    identityProvider = getLocalIdpForIssuer(jwtIssuer, tenantDomain);
                }
            }

            if (identityProvider == null) {
                throw new IdentityOAuth2Exception("No trusted IDP registered with the issuer: " + jwtIssuer
                        + " in tenantDomain: " + tenantDomain);
            } else {
                return identityProvider;
            }
        } catch (IdentityProviderManagementException e) {
            throw new IdentityOAuth2Exception("Error while retrieving trusted IDP information for issuer: " + jwtIssuer
                    + " in tenantDomain: " + tenantDomain);
        }
    }

    private IdentityProvider getLocalIdpForIssuer(String jwtIssuer,
                                                  String tenantDomain) throws IdentityOAuth2Exception {

        String residentIdpIssuer = null;
        IdentityProvider residentIdentityProvider;
        try {
            residentIdentityProvider = IdentityProviderManager.getInstance().getResidentIdP(tenantDomain);
        } catch (IdentityProviderManagementException e) {
            throw new IdentityOAuth2Exception("Error retrieving resident IDP information for issuer: " + jwtIssuer +
                    " of tenantDomain: " + tenantDomain, e);
        }

        FederatedAuthenticatorConfig[] fedAuthnConfigs = residentIdentityProvider.getFederatedAuthenticatorConfigs();
        FederatedAuthenticatorConfig oauthAuthenticatorConfig =
                IdentityApplicationManagementUtil.getFederatedAuthenticator(fedAuthnConfigs,
                        IdentityApplicationConstants.Authenticator.OIDC.NAME);
        if (oauthAuthenticatorConfig != null) {
            residentIdpIssuer = IdentityApplicationManagementUtil.getProperty(oauthAuthenticatorConfig.getProperties(),
                    Utils.OPENID_IDP_ENTITY_ID).getValue();
        }
        return StringUtils.equalsIgnoreCase(residentIdpIssuer, jwtIssuer) ? residentIdentityProvider : null;
    }

    private String getTenantDomain(JWTClaimsSet jwtClaimsSet) {

        return MultitenantConstants.SUPER_TENANT_DOMAIN_NAME;
    }

}
