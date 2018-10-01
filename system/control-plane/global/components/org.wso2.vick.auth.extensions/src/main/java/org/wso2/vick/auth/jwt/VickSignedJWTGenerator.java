package org.wso2.vick.auth.jwt;

import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.carbon.apimgt.api.APIManagementException;
import org.wso2.carbon.apimgt.keymgt.service.TokenValidationContext;
import org.wso2.carbon.apimgt.keymgt.token.JWTGenerator;
import org.wso2.carbon.identity.application.common.model.FederatedAuthenticatorConfig;
import org.wso2.carbon.identity.application.common.model.IdentityProvider;
import org.wso2.carbon.identity.application.common.util.IdentityApplicationConstants;
import org.wso2.carbon.identity.application.common.util.IdentityApplicationManagementUtil;
import org.wso2.carbon.identity.oauth.common.exception.InvalidOAuthClientException;
import org.wso2.carbon.identity.oauth.dao.OAuthAppDO;
import org.wso2.carbon.identity.oauth2.IdentityOAuth2Exception;
import org.wso2.carbon.identity.oauth2.util.OAuth2Util;
import org.wso2.carbon.idp.mgt.IdentityProviderManagementException;
import org.wso2.carbon.idp.mgt.IdentityProviderManager;
import org.wso2.vick.auth.util.Utils;

import java.nio.charset.Charset;
import java.text.ParseException;
import java.util.Arrays;
import java.util.Base64;
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.UUID;

/**
 * Generates a signed JWT with context information from API client authentication to be consumed by API backends.
 */
public class VickSignedJWTGenerator extends JWTGenerator {

    private static final Log log = LogFactory.getLog(VickSignedJWTGenerator.class);
    private static final String CONSUMER_KEY_CLAIM = "consumerKey";
    private static final String SCOPE_CLAIM = "scope";
    private static final long JWT_TOKEN_VALIDITY_IN_SECONDS = 300;
    private static final String TRUSTED_IDP_CLAIMS = "trusted_idp_claims";

    @Override
    public String generateToken(TokenValidationContext validationContext) throws APIManagementException {

        String base64UrlEncodedHeader = getBase64UrlEncodeJWTHeader(validationContext);
        String base64UrlEncodedBody = getBase64UrlEncodedJWTBody(validationContext);
        String base64UrlEncodedAssertion =
                getBase64UrlEncodedSignature(validationContext, base64UrlEncodedHeader, base64UrlEncodedBody);

        return base64UrlEncodedHeader + '.' + base64UrlEncodedBody + '.' + base64UrlEncodedAssertion;
    }

    private String getBase64UrlEncodedSignature(TokenValidationContext validationContext,
                                                String base64UrlEncodedHeader,
                                                String base64UrlEncodedBody) throws APIManagementException {

        String assertion = base64UrlEncodedHeader + '.' + base64UrlEncodedBody;
        byte[] signedAssertion = signJWT(assertion, validationContext.getTokenInfo().getEndUserName());
        if (log.isDebugEnabled()) {
            log.debug("Signed assertion value : " + new String(signedAssertion, Charset.defaultCharset()));
        }
        return Base64.getUrlEncoder().encodeToString(signedAssertion);
    }

    private String getBase64UrlEncodedJWTBody(TokenValidationContext validationContext) throws APIManagementException {

        String jwtBody = buildBody(validationContext);
        String base64UrlEncodedBody = "";
        if (jwtBody != null) {
            base64UrlEncodedBody = Base64.getUrlEncoder().encodeToString(jwtBody.getBytes());
        }
        return base64UrlEncodedBody;
    }

    private String getBase64UrlEncodeJWTHeader(TokenValidationContext validationContext) throws APIManagementException {

        String jwtHeader = buildHeader(validationContext);
        String base64UrlEncodedHeader = "";
        if (jwtHeader != null) {
            base64UrlEncodedHeader = base64UrlEncode(jwtHeader);
        }
        return base64UrlEncodedHeader;
    }

    private String base64UrlEncode(String jwtHeader) {

        return Base64.getUrlEncoder().encodeToString(jwtHeader.getBytes(Charset.defaultCharset()));
    }

    @Override
    public String buildBody(TokenValidationContext validationContext) throws APIManagementException {

        String endUserName = validationContext.getValidationInfoDTO().getEndUserName();
        Date issuedTime = new Date(System.currentTimeMillis());
        Date expiryTime = new Date(issuedTime.getTime() + JWT_TOKEN_VALIDITY_IN_SECONDS * 1000);

        JWTClaimsSet.Builder builder = new JWTClaimsSet.Builder();
        builder.subject(endUserName)
                .jwtID(UUID.randomUUID().toString())
                .issuer(getIssuer(validationContext))
                .audience(getAudienceValue(validationContext))
                .issueTime(issuedTime)
                .expirationTime(expiryTime)
                .claim(SCOPE_CLAIM, getScopes(validationContext))
                .claim(CONSUMER_KEY_CLAIM, getConsumerKey(validationContext))
                .claim(TRUSTED_IDP_CLAIMS, getClaimsFromSignedJWT(validationContext));

        return builder.build().toJSONObject().toJSONString();
    }

    private String getIssuer(TokenValidationContext validationContext) throws APIManagementException {

        String consumerKey = validationContext.getTokenInfo().getConsumerKey();
        String appTenantDomain;
        try {
            appTenantDomain = OAuth2Util.getTenantDomainOfOauthApp(consumerKey);
            return getIssuer(appTenantDomain);
        } catch (IdentityOAuth2Exception | InvalidOAuthClientException e) {
            throw new APIManagementException("Error while getting issuer value for JWT Token.", e);
        }
    }

    private String getIssuer(String tenantDomain) throws IdentityOAuth2Exception {

        IdentityProvider identityProvider = getResidentIdp(tenantDomain);
        FederatedAuthenticatorConfig[] fedAuthnConfigs = identityProvider.getFederatedAuthenticatorConfigs();
        // Get OIDC authenticator
        FederatedAuthenticatorConfig oidcAuthenticatorConfig =
                IdentityApplicationManagementUtil.getFederatedAuthenticator(fedAuthnConfigs,
                        IdentityApplicationConstants.Authenticator.OIDC.NAME);
        return IdentityApplicationManagementUtil.getProperty(oidcAuthenticatorConfig.getProperties(),
                Utils.OPENID_IDP_ENTITY_ID).getValue();
    }

    private IdentityProvider getResidentIdp(String tenantDomain) throws IdentityOAuth2Exception {

        try {
            return IdentityProviderManager.getInstance().getResidentIdP(tenantDomain);
        } catch (IdentityProviderManagementException e) {
            final String ERROR_GET_RESIDENT_IDP = "Error while getting Resident Identity Provider of '%s' tenant.";
            String errorMsg = String.format(ERROR_GET_RESIDENT_IDP, tenantDomain);
            throw new IdentityOAuth2Exception(errorMsg, e);
        }
    }

    private String getConsumerKey(TokenValidationContext validationContext) {

        return validationContext.getTokenInfo().getConsumerKey();
    }

    private String[] getScopes(TokenValidationContext validationContext) {

        return validationContext.getTokenInfo().getScopes();
    }

    private Map<String, Object> getClaimsFromSignedJWT(TokenValidationContext validationContext) {

        // Get the signed JWT access token
        String accessToken = validationContext.getAccessToken();
        if (isJWT(accessToken)) {
            try {
                SignedJWT signedJWT = SignedJWT.parse(accessToken);
                return signedJWT.getJWTClaimsSet().getClaims();
            } catch (ParseException e) {
                log.error("Error retrieving claims from the JWT Token.", e);
            }
        }

        return Collections.emptyMap();
    }

    private boolean isJWT(String accessToken) {
        // Signed JWT token contains 3 base64 encoded components separated by periods.
        return StringUtils.countMatches(accessToken, ".") == 2;
    }

    private List<String> getAudienceValue(TokenValidationContext validationContext) {

        String consumerKey = validationContext.getTokenInfo().getConsumerKey();
        try {
            OAuthAppDO app = OAuth2Util.getAppInformationByClientId(consumerKey);
            if (app.getAudiences() != null) {
                return Arrays.asList(app.getAudiences());
            }
        } catch (IdentityOAuth2Exception | InvalidOAuthClientException e) {
            log.error("Error while retrieving audience information for app with consumerKey: " + consumerKey);
        }
        return Collections.emptyList();
    }

}
