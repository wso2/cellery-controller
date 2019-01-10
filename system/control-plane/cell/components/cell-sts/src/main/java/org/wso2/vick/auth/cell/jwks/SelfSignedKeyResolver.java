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
 * Resolves keys using self signed keys.
 */
public class SelfSignedKeyResolver implements KeyResolver {

    SelfSignedCertGenerator selfSignedCertGenerator;

    public SelfSignedKeyResolver(String commonName) throws KeyResolverException {

        try {
            selfSignedCertGenerator = new SelfSignedCertGenerator(commonName);
        } catch (Exception e) {
            throw new KeyResolverException("Error while resolving keys", e);
        }
    }

    @Override
    public PrivateKey getPrivateKey() {

        return selfSignedCertGenerator.getPrivateKey();
    }

    @Override
    public PublicKey getPublicKey() {

        return selfSignedCertGenerator.getPublicKey();
    }

    @Override
    public X509Certificate getCertificate() {

        return selfSignedCertGenerator.getCertificate();
    }
}
