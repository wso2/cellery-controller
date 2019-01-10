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

package org.wso2.vick.auth.cell.utils;

import org.apache.commons.codec.binary.Base64;
import org.wso2.vick.auth.cell.jwks.KeyResolver;
import org.wso2.vick.auth.cell.jwks.KeyResolverException;
import org.wso2.vick.auth.cell.jwks.SelfSignedKeyResolver;
import org.wso2.vick.auth.cell.sts.CellStsUtils;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;

import java.nio.charset.Charset;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.cert.Certificate;
import java.security.cert.CertificateEncodingException;

/**
 * Utilities used for certificate generation and parsing.
 */
public class CertificateUtils {

    private static KeyResolver keyResolver;

    public static String getThumbPrint(Certificate certificate) throws NoSuchAlgorithmException,
            CertificateEncodingException {

        MessageDigest digestValue = MessageDigest.getInstance("SHA-1");
        byte[] der = certificate.getEncoded();
        digestValue.update(der);
        byte[] digestInBytes = digestValue.digest();
        String publicCertThumbprint = hexify(digestInBytes);
        return new String((new Base64(0, (byte[]) null, true)).
                encode(publicCertThumbprint.getBytes(Charset.forName("UTF-8"))), Charset.forName("UTF-8"));
    }

    public static String hexify(byte[] bytes) {

        if (bytes == null) {
            String errorMsg = "Invalid byte array: 'NULL'";
            throw new IllegalArgumentException(errorMsg);
        } else {
            char[] hexDigits = new char[]{'0', '1', '2', '3', '4', '5', '6', '7',
                    '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'};
            StringBuilder buf = new StringBuilder(bytes.length * 2);

            for (int i = 0; i < bytes.length; ++i) {
                buf.append(hexDigits[(bytes[i] & 240) >> 4]);
                buf.append(hexDigits[bytes[i] & 15]);
            }

            return buf.toString();
        }
    }

    public static KeyResolver getKeyResolver() throws KeyResolverException {

        if (keyResolver == null) {
            try {
                keyResolver = new SelfSignedKeyResolver(CellStsUtils.getMyCellName());
            } catch (VickCellSTSException e) {
                throw new KeyResolverException("Error while retriving key resolver", e);
            }
        }
        return keyResolver;
    }
}
