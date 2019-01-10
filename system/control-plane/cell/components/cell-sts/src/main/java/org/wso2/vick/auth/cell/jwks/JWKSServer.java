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

import com.sun.net.httpserver.HttpContext;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.utils.CertificateUtils;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateException;
import java.text.ParseException;

/**
 * Minimal server which responds to JWKS requests. This will expose the STSes JWKS endpoint
 */
public class JWKSServer {

    private static final Logger log = LoggerFactory.getLogger(JWKSServer.class);
    private static KeyResolver keyResolver;
    private int port;

    public JWKSServer(int port) {

        this.port = port;
    }

    public void startServer() throws IOException, KeyResolverException {

        keyResolver = CertificateUtils.getKeyResolver();
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        HttpContext context = server.createContext("/");
        context.setHandler(JWKSServer::handleRequest);
        server.start();
        log.info("JWKS endpoint started in port : {}", port);
    }

    private static void handleRequest(HttpExchange exchange) throws IOException {

        String response = null;
        try {
            response = JWKSResponseBuilder.buildResponse(keyResolver.getPublicKey(), keyResolver.getCertificate());
        } catch (CertificateException | NoSuchAlgorithmException | ParseException | KeyResolverException e) {
            throw new IOException("Error while building response from JWKS endpoint");
        }
        exchange.sendResponseHeaders(200, response.getBytes().length);//response code and length
        OutputStream os = exchange.getResponseBody();
        os.write(response.getBytes());
        os.close();
    }
}
