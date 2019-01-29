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

package org.wso2.vick.auth.cell.authorization;

/**
 * Interface for authorization handling. Any authorization provider should override this interface.
 */
public interface AuthorizationHandler {

    /**
     * Check authorization on a particular request.
     * @param request Authorization request.
     * @throws AuthorizationFailedException AuthorizationFailedException.
     */
    public void authorize(AuthorizeRequest request) throws AuthorizationFailedException;
}
