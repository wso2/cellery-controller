/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.wso2.vick.observability.api;

import org.apache.log4j.Logger;
import org.wso2.vick.observability.api.internal.ServiceHolder;
import org.wso2.vick.observability.model.generator.exception.GraphStoreException;
import org.wso2.vick.observability.model.generator.model.Model;

import javax.ws.rs.DefaultValue;
import javax.ws.rs.GET;
import javax.ws.rs.OPTIONS;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.QueryParam;
import javax.ws.rs.core.Response;

/**
 * MSF4J service for fetching the dependency models.
 */
@Path("/api/dependency-model")
public class DependencyModelAPI {
    private static final Logger log = Logger.getLogger(DependencyModelAPI.class);

    @GET
    @Path("/cells")
    @Produces("application/json")
    public Response getCellOverview(@DefaultValue("0") @QueryParam("fromTime") Long fromTime,
                                    @DefaultValue("0") @QueryParam("toTime") Long toTime) {
        try {
            Model model = ServiceHolder.getModelManager().getGraph(fromTime, toTime);
            return Response.ok().entity(model).build();
        } catch (GraphStoreException e) {
            log.error("Error occured while retrieving the dependency API", e);
            return Response.serverError().entity(e).build();
        }
    }

    @OPTIONS
    @Path("/*")
    public Response getOptions() {
        return Response.ok().build();
    }
}
