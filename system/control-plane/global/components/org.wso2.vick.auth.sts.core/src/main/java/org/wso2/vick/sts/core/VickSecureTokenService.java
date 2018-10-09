package org.wso2.vick.sts.core;

import org.wso2.vick.auth.exception.VickAuthException;
import org.wso2.vick.auth.jwt.VickSignedJWTBuilder;

public class VickSecureTokenService {

    public VickSTSResponse issueJWT(VickSTSRequest tokenRequest) throws VickStsException {

        // TODO we need to validate stuff before issuing the token...
        try {
            String jwt = new VickSignedJWTBuilder()
                    .subject(tokenRequest.getSource())
                    .scopes(tokenRequest.getScopes())
                    .audience(tokenRequest.getAudiences())
                    .build();

            VickSTSResponse vickSTSResponse = new VickSTSResponse();
            vickSTSResponse.setStsToken(jwt);

            return vickSTSResponse;
        } catch (VickAuthException e) {
            throw new VickStsException("Error issuing JWT", e);
        }
    }
}
