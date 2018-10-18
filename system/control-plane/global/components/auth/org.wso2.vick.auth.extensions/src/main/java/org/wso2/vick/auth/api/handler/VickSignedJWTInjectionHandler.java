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

package org.wso2.vick.auth.api.handler;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.synapse.MessageContext;
import org.apache.synapse.core.axis2.Axis2MessageContext;
import org.apache.synapse.rest.AbstractHandler;

import java.util.Map;

/**
 *  Injects the signed jwt issued by the global STS into the authorization header to be forwarded to the API back ends.
 */
public class VickSignedJWTInjectionHandler extends AbstractHandler {

    private static final String AUTHORIZATION_HEADER_NAME = "Authorization";
    private static final String JWT_ASSERTION_HEADER = "X-JWT-Assertion";

    private Log log = LogFactory.getLog(VickSignedJWTInjectionHandler.class);

    public boolean handleRequest(MessageContext messageContext) {

        String vickJWT = getVickJWT(messageContext);
        if (log.isDebugEnabled()) {
            log.debug("JWT issued from VICK STS: " + vickJWT);
        }

        removeVickSTSHeader(messageContext);
        if (log.isDebugEnabled()) {
            log.debug("Removed JWT Assertion Header: " + JWT_ASSERTION_HEADER);
        }

        String bearerHeader = "Bearer " + vickJWT;
        setAuthorizationHeader(messageContext, bearerHeader);
        if (log.isDebugEnabled()) {
            log.debug("Set new Authorization Header Value to: " + bearerHeader);
        }

        return true;
    }

    public boolean handleResponse(MessageContext messageContext) {

        return true;
    }

    private String getVickJWT(MessageContext messageContext) {

        return (String) getTransportHeaders(messageContext).get(JWT_ASSERTION_HEADER);
    }

    private void setAuthorizationHeader(MessageContext messageContext, String value) {

        getTransportHeaders(messageContext).put(AUTHORIZATION_HEADER_NAME, value);
    }

    private void removeVickSTSHeader(MessageContext messageContext) {

        getTransportHeaders(messageContext).remove(JWT_ASSERTION_HEADER);
    }

    private Map getTransportHeaders(MessageContext messageContext) {

        return (Map) ((Axis2MessageContext) messageContext).getAxis2MessageContext().
                getProperty(org.apache.axis2.context.MessageContext.TRANSPORT_HEADERS);
    }
}
