package org.wso2.vick.auth.km;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.carbon.apimgt.api.APIManagementException;
import org.wso2.carbon.apimgt.api.model.AccessTokenInfo;
import org.wso2.carbon.apimgt.api.model.KeyManagerConfiguration;
import org.wso2.carbon.apimgt.impl.AMDefaultKeyManagerImpl;
import org.wso2.carbon.apimgt.impl.APIConstants;
import org.wso2.carbon.identity.oauth2.OAuth2TokenValidationService;
import org.wso2.carbon.identity.oauth2.dto.OAuth2IntrospectionResponseDTO;
import org.wso2.carbon.identity.oauth2.dto.OAuth2TokenValidationRequestDTO;
import org.wso2.carbon.identity.oauth2.util.OAuth2Util;
import org.wso2.vick.auth.util.Utils;

import java.util.Arrays;

/**
 * Allows signed JWTs issued by trusted external IDPs to be used for API Authentication.
 */
public class VickExtendedKeyManagerImpl extends AMDefaultKeyManagerImpl {

    private static final Log log = LogFactory.getLog(VickExtendedKeyManagerImpl.class);
    private static final String JWT_TOKEN_TYPE = "jwt";
    private static final String BEARER_TOKEN_TYPE = "bearer";

    public AccessTokenInfo getTokenMetaData(String accessToken) throws APIManagementException {

        OAuth2TokenValidationRequestDTO tokenValidationRequest = buildTokenValidationRequest(accessToken);
        OAuth2IntrospectionResponseDTO introspectionResponse = introspectToken(tokenValidationRequest);

        AccessTokenInfo tokenInfo = new AccessTokenInfo();
        if (isTokenInvalid(introspectionResponse)) {
            tokenInfo.setTokenValid(false);
            tokenInfo.setErrorcode(APIConstants.KeyValidationStatus.API_AUTH_INVALID_CREDENTIALS);
        } else {
            tokenInfo.setTokenValid(true);
            tokenInfo.setEndUserName(introspectionResponse.getSub());
            tokenInfo.setConsumerKey(introspectionResponse.getClientId());
            // TODO : check on this..
            tokenInfo.setIssuedTime(System.currentTimeMillis());
            tokenInfo.setScope(buildScopes(introspectionResponse));

            // Convert Expiry Time to milliseconds.
            tokenInfo.setValidityPeriod(getExpiryPeriodInMillis(introspectionResponse));

            // If token has am_application_scope, consider the token as an Application token.
            handleScopes(introspectionResponse, tokenInfo);
        }
        return tokenInfo;
    }

    private String[] buildScopes(OAuth2IntrospectionResponseDTO introspectionResponse) {

        return OAuth2Util.buildScopeArray(introspectionResponse.getScope());
    }

    private OAuth2IntrospectionResponseDTO introspectToken(OAuth2TokenValidationRequestDTO requestDTO) {

        OAuth2TokenValidationService tokenValidationService = new OAuth2TokenValidationService();
        return tokenValidationService.buildIntrospectionResponse(requestDTO);
    }

    private boolean isTokenInvalid(OAuth2IntrospectionResponseDTO responseDTO) {

        return !responseDTO.isActive();
    }

    private long getExpiryPeriodInMillis(OAuth2IntrospectionResponseDTO introspectionResponse) {

        return System.currentTimeMillis() - (introspectionResponse.getExp() * 1000L);
    }

    @Override
    public void loadConfiguration(KeyManagerConfiguration configuration) throws APIManagementException {
        // This is a workaround to force APIM to use default values.
        configuration = null;
        super.loadConfiguration(configuration);
    }

    private void handleScopes(OAuth2IntrospectionResponseDTO responseDTO, AccessTokenInfo tokenInfo) {

        String[] scopes = OAuth2Util.buildScopeArray(responseDTO.getScope());
        String applicationTokenScope = getConfigurationElementValue(APIConstants.APPLICATION_TOKEN_SCOPE);
        if (scopes != null && applicationTokenScope != null && !applicationTokenScope.isEmpty()) {
            if (Arrays.asList(scopes).contains(applicationTokenScope)) {
                tokenInfo.setApplicationToken(true);
            }
        }
    }

    private OAuth2TokenValidationRequestDTO buildTokenValidationRequest(String accessToken) {

        OAuth2TokenValidationRequestDTO requestDTO = new OAuth2TokenValidationRequestDTO();

        OAuth2TokenValidationRequestDTO.OAuth2AccessToken token = requestDTO.new OAuth2AccessToken();
        token.setIdentifier(accessToken);

        if (Utils.isSignedJWT(accessToken)) {
            token.setTokenType(JWT_TOKEN_TYPE);
        } else {
            token.setTokenType(BEARER_TOKEN_TYPE);
        }

        requestDTO.setAccessToken(token);
        return requestDTO;
    }
}
