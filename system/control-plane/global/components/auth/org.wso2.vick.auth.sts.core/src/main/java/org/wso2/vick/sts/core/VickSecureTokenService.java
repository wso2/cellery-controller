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
package org.wso2.vick.sts.core;

import com.nimbusds.jwt.SignedJWT;
import org.apache.commons.lang.StringUtils;
import org.wso2.carbon.identity.application.common.model.IdentityProvider;
import org.wso2.carbon.identity.oauth2.IdentityOAuth2Exception;
import org.wso2.carbon.idp.mgt.IdentityProviderManagementException;
import org.wso2.vick.auth.exception.VickAuthException;
import org.wso2.vick.auth.jwt.VickSignedJWTBuilder;
import org.wso2.vick.auth.util.Utils;

import java.text.ParseException;
import java.util.HashMap;
import java.util.Map;

/**
 * This class issues the JWT taken for the STS request in vick.
 */
public class VickSecureTokenService {

    public VickSTSResponse issueJWT(VickSTSRequest tokenRequest) throws VickStsException {

        // TODO we need to validate stuff before issuing the token...
        try {
            String subject = tokenRequest.getSource();
            Map<String, Object> claims = new HashMap<>();

            if (StringUtils.isNotBlank(tokenRequest.getUserContextJwt())) {
                // TODO: add logs here.
                // If a user context jwt is set this is a token requested to impersonate a user.
                SignedJWT userContextJwt = SignedJWT.parse(tokenRequest.getUserContextJwt());
                if (isUserContextJwtValid(userContextJwt)) {
                    subject = userContextJwt.getJWTClaimsSet().getSubject();
                    claims.putAll(Utils.getCustomClaims(userContextJwt));
                } else {
                    throw new VickStsException("Invalid user context JWT presented to obtain a STS token.");
                }
            }

            String jwt = new VickSignedJWTBuilder()
                    .subject(subject)
                    .scopes(tokenRequest.getScopes())
                    .audience(tokenRequest.getAudiences())
                    .claims(claims)
                    .build();

            VickSTSResponse vickSTSResponse = new VickSTSResponse();
            vickSTSResponse.setStsToken(jwt);

            return vickSTSResponse;
        } catch (VickAuthException e) {
            throw new VickStsException("Error issuing JWT.", e);
        } catch (ParseException e) {
            throw new VickStsException("Error while parsing the user context JWT", e);
        }
    }

    private boolean isUserContextJwtValid(SignedJWT userContextJwt) throws VickAuthException {
        // We can't blindly trust the user context JWT present. So we do a signature verification to see if it was
        // issued by the VICK sts
        try {
            IdentityProvider idp = Utils.getVickIdp();
            return Utils.validateSignature(userContextJwt, idp);
        } catch (IdentityProviderManagementException | IdentityOAuth2Exception e) {
            throw new VickAuthException("Error while validating user context jwt", e);
        }
    }

}
