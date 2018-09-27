/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.vick.apiupdater.utils;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;
import java.io.IOException;
import java.net.HttpURLConnection;
import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.cert.X509Certificate;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.http.HttpEntity;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.util.EntityUtils;
import org.wso2.vick.apiupdater.exceptions.APIException;

/**
 * Utility methods for HTTP request processors
 */
public class RequestProcessor {

    private static final Log log = LogFactory.getLog(RequestProcessor.class);
    private CloseableHttpClient httpClient;

    public RequestProcessor() throws APIException {
        try {
            if (log.isDebugEnabled()) {
                log.debug("Ignoring SSL verification...");
            }
            SSLContext sslContext = SSLContext.getInstance("SSL");

            sslContext.init(null, new TrustManager[] { new X509TrustManager() {
                public X509Certificate[] getAcceptedIssuers() {
                    return null;
                }

                public void checkClientTrusted(X509Certificate[] certs, String authType) {
                }

                public void checkServerTrusted(X509Certificate[] certs, String authType) {
                }
            } }, new SecureRandom());

            SSLConnectionSocketFactory sslsocketFactory =
                    new SSLConnectionSocketFactory(sslContext, new String[] { "TLSv1.2" }, null,
                            new HostnameVerifier() {
                                @Override
                                public boolean verify(String s, SSLSession sslSession) {
                                    return true;
                                }
                            });

            httpClient = HttpClients.custom().setSSLSocketFactory(sslsocketFactory).build();
        } catch (NoSuchAlgorithmException | KeyManagementException e) {
            String errorMessage = "Error occurred while ignoring ssl certificates to allow http connections";
            log.error(errorMessage, e);
            throw new APIException(errorMessage, e);
        }
    }

    /**
     * Execute http get request.
     *
     * @param url         url
     * @param contentType content type
     * @param acceptType  accept type
     * @param authHeader  authorization header
     * @return Closable http response
     * @throws APIException Api exception when an error occurred
     */
    public CloseableHttpResponse doGet(String url, String contentType, String acceptType, String authHeader)
            throws APIException {

        CloseableHttpResponse response;
        try {
            HttpGet httpGet = new HttpGet(url);
            httpGet.setHeader(Constants.Utils.HTTP_CONTENT_TYPE, contentType);
            httpGet.setHeader(Constants.Utils.HTTP_RESPONSE_TYPE_ACCEPT, acceptType);
            httpGet.setHeader(Constants.Utils.HTTP_REQ_HEADER_AUTHZ, authHeader);

            response = httpClient.execute(httpGet);
            closeClientConnection();
        } catch (IOException e) {
            String errorMessage = "Error occurred while executing the http Get connection.";
            log.error(errorMessage, e);
            throw new APIException(errorMessage, e);
        }
        return response;
    }

    /**
     * Execute http post request
     *
     * @param url         url
     * @param contentType content type
     * @param acceptType  accept type
     * @param authHeader  authorization header
     * @param payload     post payload
     * @return Closable http response
     * @throws APIException Api exception when an error occurred
     */
    public String doPost(String url, String contentType, String acceptType, String authHeader, String payload)
            throws APIException {
        String returnObj = null;
        try {
            StringEntity payloadEntity = new StringEntity(payload);
            HttpPost httpPost = new HttpPost(url);
            httpPost.setHeader(Constants.Utils.HTTP_CONTENT_TYPE, contentType);
            httpPost.setHeader(Constants.Utils.HTTP_RESPONSE_TYPE_ACCEPT, acceptType);
            httpPost.setHeader(Constants.Utils.HTTP_REQ_HEADER_AUTHZ, authHeader);
            httpPost.setEntity(payloadEntity);

            CloseableHttpResponse response = httpClient.execute(httpPost);
            HttpEntity entity = response.getEntity();
            String responseStr = EntityUtils.toString(entity);
            int statusCode = response.getStatusLine().getStatusCode();

            if (responseValidate(statusCode, responseStr)) {
                returnObj = responseStr;
            }
            closeClientConnection();
        } catch (IOException e) {
            String errorMessage = "Error occurred while executing the http Post connection.";
            log.error(errorMessage, e);
            throw new APIException(errorMessage, e);
        }
        return returnObj;
    }

    /**
     * Close http client connection
     *
     * @throws IOException throws an IO exception if an error occurred while closing the connection.
     */
    private void closeClientConnection() throws IOException {
        if (httpClient != null) {
            try {
                httpClient.close();
            } catch (IOException e) {
                log.error("Error while closing the http client connection", e);
                throw e;
            }
        }
    }

    /**
     * Validate the http response.
     * @param statusCode status code
     * @return boolean to validate response
     */
    private boolean responseValidate(int statusCode, String response) throws IOException {
        switch (statusCode) {
            case HttpURLConnection.HTTP_OK:
                return true;
            case HttpURLConnection.HTTP_CREATED:
                return true;
            case HttpURLConnection.HTTP_ACCEPTED:
                return true;
            case HttpURLConnection.HTTP_BAD_REQUEST:
                if (response != null && !Constants.Utils.EMPTY_STRING.equals(response)) {
                    if (response.contains(Constants.Utils.DIFFERENT_CONTEXT_ERROR)) {
                        // skip the error when trying to add the same api with different context.
                        return true;
                    }
                } else {
                    return false;
                }
            case HttpURLConnection.HTTP_UNAUTHORIZED:
                return false;
            case HttpURLConnection.HTTP_NOT_FOUND:
                return false;
            case HttpURLConnection.HTTP_CONFLICT:
                if (response != null && !Constants.Utils.EMPTY_STRING.equals(response)) {
                    if (response.contains(Constants.Utils.DUPLICATE_API_ERROR)) {
                        // skip the error when trying to add the same api.
                        return true;
                    }
                } else {
                    return false;
                }
            case HttpURLConnection.HTTP_INTERNAL_ERROR:
                if (response != null && !Constants.Utils.EMPTY_STRING.equals(response)) {
                    if (response.contains(Constants.Utils.DUPLICATE_LABEL_ERROR)) {
                        // skip the error when trying to add the same label.
                        return true;
                    }
                } else {
                    return false;
                }
            default:
                return false;
        }
    }
}
