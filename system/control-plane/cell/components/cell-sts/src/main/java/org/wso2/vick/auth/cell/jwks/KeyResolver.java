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

package org.wso2.vick.auth.cell.jwks;

import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.cert.X509Certificate;

/**
 * Resolves private keys, public keys and respective certificate of the system.
 */
public interface KeyResolver {

    /**
     * Returns the private key of the system.
     * @return Private Key.
     * @throws KeyResolverException KeyResolverException.
     */
    public PrivateKey getPrivateKey() throws KeyResolverException;

    /**
     * Returns the public key of the system.
     * @return Public Key.
     * @throws KeyResolverException KeyResolverException.
     */
    public PublicKey getPublicKey() throws KeyResolverException;

    /**
     * Returns the private key of the system.
     * @return Private Key.
     * @throws KeyResolverException KeyResolverException.
     */
    public X509Certificate getCertificate() throws KeyResolverException;

}
