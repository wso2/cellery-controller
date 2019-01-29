/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package org.wso2.vick.auth.cell.authorization.opa;

import com.nimbusds.jwt.JWT;
import com.nimbusds.jwt.JWTParser;
import org.wso2.vick.auth.cell.authorization.AuthorizationContext;
import org.wso2.vick.auth.cell.authorization.AuthorizationFailedException;

import java.text.ParseException;

public class OPAAuthorizationContext extends AuthorizationContext {

    private String jwtContent;

    public OPAAuthorizationContext(String jwt) throws AuthorizationFailedException {

        super(jwt);
        try {
            JWT parsedJWT = JWTParser.parse(jwt);
            jwtContent = parsedJWT.getJWTClaimsSet().toJSONObject().toJSONString();
        } catch (ParseException e) {
            throw new AuthorizationFailedException("Error while parsing JWT", e);
        }

    }
}
