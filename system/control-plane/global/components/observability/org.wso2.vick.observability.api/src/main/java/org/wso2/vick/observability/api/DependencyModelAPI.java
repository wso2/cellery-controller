/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.observability.api;

import org.wso2.vick.observability.api.model.Graph;
import org.wso2.vick.observability.api.model.GraphEdge;
import org.wso2.vick.observability.model.generator.ModelManager;

import java.util.ArrayList;
import java.util.List;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.Response;

/**
 * This is MSF4J service for the dependency model to plot the UI graph.
 */
@Path("/dependency-model")
public class DependencyModelAPI {

    @GET
    @Path("/cell-overview")
    @Produces("application/json")
    public Response getCellOverview() {
        List<GraphEdge> graphEdges = new ArrayList<>();

        for (String edges : ModelManager.getInstance().getLinks()) {
            String[] edgeNameElements = ModelManager.getInstance().
                    edgeNameElements(edges);
            GraphEdge edge = new GraphEdge(edgeNameElements[0], edgeNameElements[1]);
            graphEdges.add(edge);
        }
        return Response.ok().header("Access-Control-Allow-Origin", "*")
                .header("Access-Control-Allow-Credentials", "true")
                .header("Access-Control-Allow-Methods", "POST, GET, PUT, UPDATE, DELETE, OPTIONS, HEAD")
                .header("Access-Control-Allow-Headers", "Content-Type, Accept, X-Requested-With")
                .entity(new Graph(ModelManager.getInstance().getNodes(), graphEdges))
                .build();
    }

}
