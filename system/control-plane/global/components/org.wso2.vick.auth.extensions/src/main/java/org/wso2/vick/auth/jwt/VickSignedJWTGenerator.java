package org.wso2.vick.auth.jwt;

import com.nimbusds.jwt.SignedJWT;
import edu.emory.mathcs.backport.java.util.Arrays;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.carbon.apimgt.api.APIManagementException;
import org.wso2.carbon.apimgt.keymgt.service.TokenValidationContext;
import org.wso2.carbon.apimgt.keymgt.token.JWTGenerator;
import org.wso2.vick.auth.exception.VickAuthException;
import org.wso2.vick.auth.util.Utils;

import java.text.ParseException;
import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * Generates a signed JWT with context information from API client authentication to be consumed by API backends.
 */
public class VickSignedJWTGenerator extends JWTGenerator {

    private static final Log log = LogFactory.getLog(VickSignedJWTGenerator.class);
    private static final String CONSUMER_KEY_CLAIM = "consumerKey";
    private static final long JWT_TOKEN_VALIDITY_IN_SECONDS = 300L;
    private static final String TRUSTED_IDP_CLAIMS = "trusted_idp_claims";

    @Override
    public String generateToken(TokenValidationContext validationContext) throws APIManagementException {

        VickSignedJWTBuilder jwtBuilder = new VickSignedJWTBuilder();
        try {
            return jwtBuilder.subject(getEndUserName(validationContext))
                    .expiryInSeconds(JWT_TOKEN_VALIDITY_IN_SECONDS)
                    .scopes(getScopes(validationContext))
                    .claim(CONSUMER_KEY_CLAIM, getConsumerKey(validationContext))
                    .claim(TRUSTED_IDP_CLAIMS, getClaimsFromSignedJWT(validationContext))
                    .build();
        } catch (VickAuthException e) {
            throw new APIManagementException("Error generating JWT for context: " + validationContext, e);
        }

    }

    private String getEndUserName(TokenValidationContext validationContext) {

        return validationContext.getValidationInfoDTO().getEndUserName();
    }

    private String getConsumerKey(TokenValidationContext validationContext) {

        return validationContext.getTokenInfo().getConsumerKey();
    }

    private List<String> getScopes(TokenValidationContext validationContext) {

        String[] scopes = validationContext.getTokenInfo().getScopes();
        return scopes != null ? Arrays.asList(scopes) : Collections.emptyList();
    }

    private Map<String, Object> getClaimsFromSignedJWT(TokenValidationContext validationContext) {

        // Get the signed JWT access token
        String accessToken = validationContext.getAccessToken();
        if (Utils.isSignedJWT(accessToken)) {
            try {
                SignedJWT signedJWT = SignedJWT.parse(accessToken);
                return signedJWT.getJWTClaimsSet().getClaims();
            } catch (ParseException e) {
                log.error("Error retrieving claims from the JWT Token.", e);
            }
        }

        return Collections.emptyMap();
    }
}
