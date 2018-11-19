
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

import {ConfigHolder} from "./configHolder";

describe("ConfigHolder", () => {
    describe("set()", () => {
        it("should add the values proviced against the keys", () => {
            const configHolder = new ConfigHolder();
            configHolder.set("key1", "value1");
            configHolder.set("key2", "value2");
            configHolder.set("key3", "value3");

            expect(Object.keys(configHolder.config)).toHaveLength(3);
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key2.value).toBe("value2");
            expect(configHolder.config.key3.value).toBe("value3");
        });

        it("should not fail upon passing undefined or null", () => {
            const configHolder = new ConfigHolder();
            configHolder.set("key1", null);
            configHolder.set("key2", undefined);
            configHolder.set(null, "value1");
            configHolder.set(null, null);
            configHolder.set(null, undefined);
            configHolder.set(undefined, "value3");
            configHolder.set(undefined, null);
            configHolder.set(undefined, undefined);

            expect(Object.keys(configHolder.config)).toHaveLength(2);
            expect(configHolder.config.key1.value).toBeNull();
            expect(configHolder.config.key2.value).toBeUndefined();
        });
    });

    describe("get()", () => {
        it("should return the value stored before for the specified key", () => {
            const configHolder = new ConfigHolder();
            configHolder.config.key1 = {
                value: "value1",
                listeners: []
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: []
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: []
            };

            expect(configHolder.get("key1", "defaultValue1")).toBe("value1");
            expect(configHolder.get("key2", "defaultValue2")).toBe("value2");
            expect(configHolder.get("key3", "defaultValue3")).toBe("value3");
        });

        it("should return the value stored before for the specified key when the default value is not provided", () => {
            const configHolder = new ConfigHolder();
            configHolder.config.key1 = {
                value: "value1",
                listeners: []
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: []
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: []
            };

            expect(configHolder.get("key1", "defaultValue1")).toBe("value1");
            expect(configHolder.get("key2", "defaultValue2")).toBe("value2");
            expect(configHolder.get("key3", "defaultValue3")).toBe("value3");
        });

        it("should return the default value provided if the provided key was does not exist", () => {
            const configHolder = new ConfigHolder();
            configHolder.config.key1 = {
                value: "value1",
                listeners: []
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: []
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: []
            };

            expect(configHolder.get("non-existent-key1", "defaultValue1")).toBe("defaultValue1");
            expect(configHolder.get("non-existent-key2", "defaultValue2")).toBe("defaultValue2");
            expect(configHolder.get("non-existent-key3", "defaultValue3")).toBe("defaultValue3");
        });

        it("should return null if the provided key was does not exist and default value was not provided", () => {
            const configHolder = new ConfigHolder();
            configHolder.config.key1 = {
                value: "value1",
                listeners: []
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: []
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: []
            };

            expect(configHolder.get("non-existent-key1")).toBeNull();
            expect(configHolder.get("non-existent-key2")).toBeNull();
            expect(configHolder.get("non-existent-key3")).toBeNull();
        });
    });

    describe("loadConfig()", () => {
        it("should return null if the provided key was does not exist and default value was not provided", () => {
            const configHolder = new ConfigHolder();
            configHolder.loadConfig();

            expect(Object.keys(configHolder.config).length > 0).toBe(true);
        });
    });
});
