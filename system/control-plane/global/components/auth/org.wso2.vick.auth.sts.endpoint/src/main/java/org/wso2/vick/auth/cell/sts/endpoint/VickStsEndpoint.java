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

package org.wso2.vick.auth.cell.sts.endpoint;

import org.apache.commons.lang.StringUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.vick.sts.core.VickSTSConstants;
import org.wso2.vick.sts.core.VickSTSRequest;
import org.wso2.vick.sts.core.VickSTSResponse;
import org.wso2.vick.sts.core.VickSecureTokenService;
import org.wso2.vick.sts.core.VickStsException;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import javax.servlet.http.HttpServletRequest;
import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.Context;
import javax.ws.rs.core.MultivaluedMap;
import javax.ws.rs.core.Response;

/**
 * This is the REST service that is exposed to get the STS token.
 */
@Path("/sts")
public class VickStsEndpoint {

    private static final Log log = LogFactory.getLog(VickStsEndpoint.class);

    private VickSecureTokenService tokenService = new VickSecureTokenService();

    @POST
    @Path("/token")
    @Consumes("application/x-www-form-urlencoded")
    @Produces("application/json")
    public Response getStsToken(@Context HttpServletRequest request, MultivaluedMap<String, String> form) {

        VickSTSResponse stsResponse;
        try {
            VickSTSRequest vickSTSRequest = buildStsRequest(request, form);
            stsResponse = tokenService.issueJWT(vickSTSRequest);
        } catch (VickStsException e) {
            log.error("Error while issuing STS Token.", e);
            return Response.serverError().build();
        }

        // Build response.
        return Response.ok().entity(stsResponse.toJson()).build();
    }

    private VickSTSRequest buildStsRequest(HttpServletRequest request, MultivaluedMap<String, String> form) {

        VickSTSRequest stsRequest = new VickSTSRequest();
        stsRequest.setSource(form.getFirst(VickSTSConstants.VickSTSRequest.SUBJECT));
        stsRequest.setScopes(buildValueList(form.getFirst(VickSTSConstants.VickSTSRequest.SCOPE)));
        stsRequest.setAudiences(buildValueList(form.getFirst(VickSTSConstants.VickSTSRequest.AUDIENCE)));
        stsRequest.setUserContextJwt(form.getFirst(VickSTSConstants.VickSTSRequest.USER_CONTEXT_JWT));
        return stsRequest;
    }

    private List<String> buildValueList(String value) {

        if (StringUtils.isNotBlank(value)) {
            value = value.trim();
            return Arrays.asList(value.split("\\s"));
        } else {
            return Collections.emptyList();
        }
    }
}

