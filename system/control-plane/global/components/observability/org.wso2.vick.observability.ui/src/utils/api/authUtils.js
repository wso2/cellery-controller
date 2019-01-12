/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import {StateHolder} from "../../components/common/state";

/**
 * Authentication/Authorization related utilities.
 */
class AuthUtils {

    /**
     * Sign in the user.
     *
     * @param {string} username The user to be signed in
     * @param {StateHolder} globalState The global state provided to the current component
     */
    static signIn = (username, globalState) => {
        // TODO: Implement User Login
        if (username) {
            localStorage.setItem(StateHolder.USER, username);
            globalState.set(StateHolder.USER, username);
        } else {
            throw Error(`Username provided cannot be "${username}"`);
        }
    };

    /**
     * Sign out the current user.
     * The provided global state will be updated accordingly as well.
     *
     * @param {StateHolder} globalState The global state provided to the current component
     */
    static signOut = (globalState) => {
        // TODO: Implement User Logout
        localStorage.removeItem(StateHolder.USER);
        globalState.unset(StateHolder.USER);
    };

    /**
     * Get the currently authenticated user.
     *
     * @returns {string} The current user
     */
    static getAuthenticatedUser = () => localStorage.getItem(StateHolder.USER);

}

export default AuthUtils;
