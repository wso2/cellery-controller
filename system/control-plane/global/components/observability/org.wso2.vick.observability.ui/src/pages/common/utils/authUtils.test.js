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
import {ConfigConstants, ConfigHolder} from "../config/configHolder";

describe("AuthUtils", () => {
    afterEach(() => {
        localStorage.removeItem(ConfigConstants.USER);
    });

    describe("signIn()", () => {
        it("should set the username provided", () => {
            const config = new ConfigHolder();
            const spy = jest.spyOn(config, "set");
            AuthUtils.signIn("user1", config);

            expect(spy).toHaveBeenCalledTimes(1);
            expect(spy).toHaveBeenCalledWith(ConfigConstants.USER, "user1");
            expect(localStorage.getItem(ConfigConstants.USER)).toBe("user1");
        });

        it("should not set a username and should throw and error", () => {
            const config = new ConfigHolder();
            const spy = jest.spyOn(config, "set");

            expect(() => AuthUtils.signIn(null, config)).toThrow();
            expect(() => AuthUtils.signIn(undefined, config)).toThrow();
            expect(() => AuthUtils.signIn("", config)).toThrow();
            expect(spy).toHaveBeenCalledTimes(0);
            expect(spy).not.toHaveBeenCalled();
            expect(localStorage.getItem(ConfigConstants.USER)).toBeNull();
        });
    });

    describe("signOut()", () => {
        it("should unset the user in the configuration", () => {
            const config = new ConfigHolder();
            localStorage.setItem(ConfigConstants.USER, "user1");
            config.config[ConfigConstants.USER] = {
                value: "user1",
                listeners: []
            };
            const spy = jest.spyOn(config, "unset");
            AuthUtils.signOut(config);

            expect(spy).toHaveBeenCalledTimes(1);
            expect(spy).toHaveBeenCalledWith(ConfigConstants.USER);
            expect(localStorage.getItem(ConfigConstants.USER)).toBeNull();
        });
    });

    describe("getAuthenticatedUser()", () => {
        localStorage.setItem(ConfigConstants.USER, "user1");
        const user = AuthUtils.getAuthenticatedUser();

        expect(user).toBe("user1");
    });
});
