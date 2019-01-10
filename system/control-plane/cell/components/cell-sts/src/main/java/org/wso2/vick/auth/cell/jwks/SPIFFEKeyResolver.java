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
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;

/**
 * Resolves keys derived in the format of SPIFFE
 */
public class SPIFFEKeyResolver extends StaticKeyResolver {

    // Sample SPIFFE public certificate string
    private static String publicKeyString = "-----BEGIN CERTIFICATE-----\n" +
            "MIIDJjCCAg6gAwIBAgIRANmawRKKvf2ZSIXfzIXpj9cwDQYJKoZIhvcNAQELBQAw\n" +
            "HDEaMBgGA1UEChMRazhzLmNsdXN0ZXIubG9jYWwwHhcNMTgxMTIzMTAyNDU2WhcN\n" +
            "MTkwMjIxMTAyNDU2WjALMQkwBwYDVQQKEwAwggEiMA0GCSqGSIb3DQEBAQUAA4IB\n" +
            "DwAwggEKAoIBAQCdWHtkvQop9FJrPQ6MgIkdxBBSSyuT38VXxkfp7bx8bU/HEImg\n" +
            "UUpuH60MXJSAWh3CQ2lbjcwCupt94QSDu3bL9qfhekbRcrduwY6lZWD+nx9Yy3gL\n" +
            "8dukUozFJHAUYIpwa4Sjjrt28zWxMQnQ6/tGqUPXzD2NYZyQmQM0ZngUq/34ypoz\n" +
            "CUrWaPCsHnYYSKfHPb0xcoCSg3vZpNlsyqX8qaisCEFSVI4psuB6ktGBJ2bVfhQS\n" +
            "3RlvncpEd+6gaKrTS3i9b4hdvhPHOkwpvyF9jkNu4QWTRmq1/7j+ptLvNatqNEIW\n" +
            "IPO0McQcnMAvwb1IuuQQKmbgmR36d/4m9yuNAgMBAAGjdDByMA4GA1UdDwEB/wQE\n" +
            "AwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw\n" +
            "ADAzBgNVHREELDAqhihzcGlmZmU6Ly9jbHVzdGVyLmxvY2FsL25zL2Zvby9zYS9k\n" +
            "ZWZhdWx0MA0GCSqGSIb3DQEBCwUAA4IBAQBD24+bdQ6n43soOZsjqf1OA3/3WNiE\n" +
            "IL+GANg1jDDjYJw+gjl2RNa+EJNoHGdJ9x4X3U9dFkHzSF7/sa4Wr5wa5Ztkn9dN\n" +
            "PS7srNjplfoeIDW7y9uqlept1hU/8LL9J4PKzHWQtRYrbwmyJ1cql3rubrGSRjSu\n" +
            "nO7xDZGl6H33UhWnFT3zJYPJ7+mY6jxgyhYxoNylC1pQFloWGmjf3Q33LYXcawJN\n" +
            "9qqHgRbvo9RTMHrPqk/9mDl5KLhQr//d3K+rBfO9A7MjAaKk+/1SQHCd9DwLyAlL\n" +
            "gZwKV3qnWX9WAevQJPHiQ0KIBAUwygm8iUlp+6x7fpCIrH/sKZCFpzNa\n" +
            "-----END CERTIFICATE-----\n";

    @Override
    public PrivateKey getPrivateKey() {

        // No private key.
        return null;
    }

    @Override
    public PublicKey getPublicKey() throws KeyResolverException {

        try {
            return buildCertificate(publicKeyString).getPublicKey();
        } catch (CertificateException e) {
            throw new KeyResolverException("Error while retrieving public key", e);
        }
    }

    @Override
    public X509Certificate getCertificate() throws KeyResolverException {

        try {
            return (X509Certificate) buildCertificate(publicKeyString);
        } catch (CertificateException e) {
            throw new KeyResolverException("Error while building certificate", e);
        }
    }

}
