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
package org.wso2.vick.observability.model.generator.model;

import org.wso2.vick.observability.model.generator.Utils;

/**
 * This is the POJO of the Model edge information that is being returned to by the MSF4J services - DependencyModelAPI
 */
public class Edge {
    private String source;
    private String target;
    private String edgeString;

    public Edge(String edgeString) {
        this.edgeString = edgeString;
        String[] edgeElements = Utils.edgeNameElements(edgeString);
        this.source = edgeElements[0];
        this.target = edgeElements[1];
    }

    public String getSource() {
        return source;
    }

    public String getTarget() {
        return target;
    }

    public String getEdgeString() {
        return edgeString;
    }

    public int compareTo(Object anotherNode) {
        if (anotherNode != null && anotherNode instanceof Edge) {
            if (this.equals(anotherNode)) {
                return 0;
            }
            return edgeString.compareTo(((Edge) anotherNode).edgeString);
        } else {
            return -1;
        }
    }

    public boolean equals(Object object) {
        return object != null && object instanceof Edge && edgeString.equalsIgnoreCase(((Edge) object).edgeString);
    }
}
