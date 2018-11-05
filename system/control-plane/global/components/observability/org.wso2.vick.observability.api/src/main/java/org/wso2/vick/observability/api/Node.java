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

/**
 * This is the POJO to store the asset related details in the document.
 */
public class Node implements Comparable {
    private String name;
    private String tags;

    public Node(String name, String tags) {
        this.name = name;
        this.tags = tags;
    }

    public String getName() {
        return name;
    }

    public String getTags() {
        return tags;
    }

    public int compareTo(Object anotherNode) {
        if (anotherNode != null && anotherNode instanceof Node) {
            if (this.equals(anotherNode)) {
                return 0;
            }
            return name.compareTo(((Node) anotherNode).name);
        } else {
            return -1;
        }
    }

    public boolean equals(Object object) {
        return object != null && object instanceof Node && name.equalsIgnoreCase(((Node) object).name);
    }

    public int hashCode() {
        return name.hashCode();
    }
}
