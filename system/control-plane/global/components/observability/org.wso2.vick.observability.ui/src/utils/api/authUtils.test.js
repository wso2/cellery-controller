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

import AuthUtils from "./authUtils";
import {StateHolder} from "../../components/common/state";

describe("AuthUtils", () => {
    afterEach(() => {
        localStorage.removeItem(StateHolder.USER);
    });

    describe("signIn()", () => {
        it("should set the username provided", () => {
            const stateHolder = new StateHolder();
            const spy = jest.spyOn(stateHolder, "set");
            AuthUtils.signIn("user1", stateHolder);

            expect(spy).toHaveBeenCalledTimes(1);
            expect(spy).toHaveBeenCalledWith(StateHolder.USER, "user1");
            expect(localStorage.getItem(StateHolder.USER)).toBe("user1");
        });

        it("should not set a username and should throw and error", () => {
            const stateHolder = new StateHolder();
            const spy = jest.spyOn(stateHolder, "set");

            expect(() => AuthUtils.signIn(null, stateHolder)).toThrow();
            expect(() => AuthUtils.signIn(undefined, stateHolder)).toThrow();
            expect(() => AuthUtils.signIn("", stateHolder)).toThrow();
            expect(spy).toHaveBeenCalledTimes(0);
            expect(spy).not.toHaveBeenCalled();
            expect(localStorage.getItem(StateHolder.USER)).toBeNull();
        });
    });

    describe("signOut()", () => {
        it("should unset the user in the state", () => {
            const stateHolder = new StateHolder();
            localStorage.setItem(StateHolder.USER, "user1");
            stateHolder.state[StateHolder.USER] = {
                value: "user1",
                listeners: []
            };
            const spy = jest.spyOn(stateHolder, "unset");
            AuthUtils.signOut(stateHolder);

            expect(spy).toHaveBeenCalledTimes(1);
            expect(spy).toHaveBeenCalledWith(StateHolder.USER);
            expect(localStorage.getItem(StateHolder.USER)).toBeNull();
        });
    });

    describe("getAuthenticatedUser()", () => {
        localStorage.setItem(StateHolder.USER, "user1");
        const user = AuthUtils.getAuthenticatedUser();

        expect(user).toBe("user1");
    });
});
