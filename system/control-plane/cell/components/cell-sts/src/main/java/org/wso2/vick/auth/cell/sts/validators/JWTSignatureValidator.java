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
 */

package org.wso2.vick.auth.cell.sts.validators;

import com.nimbusds.jwt.JWT;
import org.wso2.vick.auth.cell.sts.exception.TokenValidationFailureException;

import java.util.Map;

/**
 * Validates JWT Signature.
 */
public interface JWTSignatureValidator {

    /**
     * Validates a JWT string using jwks_uri.
     *
     * @param jwt       JWT string.
     * @param jwksUri   Identity provider's jwks_uri.
     * @param algorithm JWT signature algorithm.
     * @param opts      Optional values.
     * @return whether the provided token is valid.
     * @throws TokenValidationFailureException
     */
    public boolean validateSignature(String jwt, String jwksUri, String algorithm, Map<String, Object> opts) throws
            TokenValidationFailureException;

    /**
     * Validates a JWT token using jwks_uri.
     *
     * @param jwt       JWT.
     * @param jwksUri   Identity provider's jwks_uri.
     * @param algorithm JWT signature algorithm.
     * @param opts      Optional values.
     * @return whether the provided token is valid.
     * @throws TokenValidationFailureException
     */
    public boolean validateSignature(JWT jwt, String jwksUri, String algorithm, Map<String, Object> opts) throws
            TokenValidationFailureException;
}
