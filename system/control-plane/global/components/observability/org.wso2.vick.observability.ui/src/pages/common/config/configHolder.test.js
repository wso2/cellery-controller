
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

/* eslint max-lines: ["off"] */

import ConfigHolder from "./configHolder";

describe("ConfigHolder", () => {
    describe("set()", () => {
        it("should add the values provided against the keys when a new key is provided", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            configHolder.set("key2", "value2");
            configHolder.set("key3", "value3");
            configHolder.set("key4", null);
            configHolder.set("key5", undefined);

            expect(Object.keys(configHolder.config)).toHaveLength(5);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBe("value2");
            expect(configHolder.config.key2.listeners).toBeUndefined();
            expect(configHolder.config.key3).not.toBeUndefined();
            expect(configHolder.config.key3.value).toBe("value3");
            expect(configHolder.config.key3.listeners).toBeUndefined();
            expect(configHolder.config.key4).not.toBeUndefined();
            expect(configHolder.config.key4.value).toBeNull();
            expect(configHolder.config.key4.listeners).toBeUndefined();
            expect(configHolder.config.key5).not.toBeUndefined();
            expect(configHolder.config.key5.value).toBeUndefined();
            expect(configHolder.config.key5.listeners).toBeUndefined();
            expect(key1Callback1).not.toHaveBeenCalled();
        });

        it("should not fail upon passing undefined, null or empty string as key", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            configHolder.set(null, "value1");
            configHolder.set(null, null);
            configHolder.set(null, undefined);
            configHolder.set(undefined, "value2");
            configHolder.set(undefined, null);
            configHolder.set(undefined, undefined);
            configHolder.set("", "value3");
            configHolder.set("", null);
            configHolder.set("", undefined);

            expect(Object.keys(configHolder.config)).toHaveLength(1);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(key1Callback1).not.toHaveBeenCalled();
        });

        it("should set the new value and notify the listeners if an existing key is provided", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            configHolder.config.key1 = {
                value: null,
                listeners: [key1Callback1]
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            configHolder.set("key1", "new-value1");
            configHolder.set("key2", "new-value2");

            expect(Object.keys(configHolder.config)).toHaveLength(3);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("new-value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBe("new-value2");
            expect(configHolder.config.key2.listeners).not.toBeUndefined();
            expect(configHolder.config.key2.listeners).toHaveLength(3);
            expect(configHolder.config.key2.listeners[0]).toBe(key2Callback1);
            expect(configHolder.config.key2.listeners[1]).toBe(key2Callback2);
            expect(configHolder.config.key2.listeners[2]).toBe(key2Callback3);
            expect(configHolder.config.key3).not.toBeUndefined();
            expect(configHolder.config.key3.value).toBe("value3");
            expect(configHolder.config.key3.listeners).not.toBeUndefined();
            expect(configHolder.config.key3.listeners).toHaveLength(2);
            expect(configHolder.config.key3.listeners[0]).toBe(key3Callback1);
            expect(configHolder.config.key3.listeners[1]).toBe(key3Callback2);
            expect(key1Callback1).toHaveBeenCalledTimes(1);
            expect(key1Callback1).toHaveBeenCalledWith("key1", null, "new-value1");
            expect(key2Callback1).toHaveBeenCalledTimes(1);
            expect(key2Callback1).toHaveBeenCalledWith("key2", "value2", "new-value2");
            expect(key2Callback2).toHaveBeenCalledTimes(1);
            expect(key2Callback2).toHaveBeenCalledWith("key2", "value2", "new-value2");
            expect(key2Callback3).toHaveBeenCalledTimes(1);
            expect(key2Callback3).toHaveBeenCalledWith("key2", "value2", "new-value2");
            expect(key3Callback1).not.toHaveBeenCalled();
            expect(key3Callback2).not.toHaveBeenCalled();
        });

        it("should not call the listeners if the provided new value is equal to the existing value", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            configHolder.set("key2", "value2");

            expect(Object.keys(configHolder.config)).toHaveLength(3);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBe("value2");
            expect(configHolder.config.key2.listeners).not.toBeUndefined();
            expect(configHolder.config.key2.listeners).toHaveLength(3);
            expect(configHolder.config.key2.listeners[0]).toBe(key2Callback1);
            expect(configHolder.config.key2.listeners[1]).toBe(key2Callback2);
            expect(configHolder.config.key2.listeners[2]).toBe(key2Callback3);
            expect(configHolder.config.key3).not.toBeUndefined();
            expect(configHolder.config.key3.value).toBe("value3");
            expect(configHolder.config.key3.listeners).not.toBeUndefined();
            expect(configHolder.config.key3.listeners).toHaveLength(2);
            expect(configHolder.config.key3.listeners[0]).toBe(key3Callback1);
            expect(configHolder.config.key3.listeners[1]).toBe(key3Callback2);
            expect(key1Callback1).not.toHaveBeenCalled();
            expect(key2Callback1).not.toHaveBeenCalled();
            expect(key2Callback2).not.toHaveBeenCalled();
            expect(key2Callback3).not.toHaveBeenCalled();
            expect(key3Callback1).not.toHaveBeenCalled();
            expect(key3Callback2).not.toHaveBeenCalled();
        });
    });

    describe("unset()", () => {
        it("should set the value to null", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key1Callback2 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1, key1Callback2]
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: [key3Callback1]
            };
            configHolder.unset("key1");

            expect(Object.keys(configHolder.config)).toHaveLength(3);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBeNull();
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(2);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key1.listeners[1]).toBe(key1Callback2);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBe("value2");
            expect(configHolder.config.key2.listeners).not.toBeUndefined();
            expect(configHolder.config.key2.listeners).toHaveLength(3);
            expect(configHolder.config.key2.listeners[0]).toBe(key2Callback1);
            expect(configHolder.config.key2.listeners[1]).toBe(key2Callback2);
            expect(configHolder.config.key2.listeners[2]).toBe(key2Callback3);
            expect(configHolder.config.key3).not.toBeUndefined();
            expect(configHolder.config.key3.value).toBe("value3");
            expect(configHolder.config.key3.listeners).not.toBeUndefined();
            expect(configHolder.config.key3.listeners).toHaveLength(1);
            expect(configHolder.config.key3.listeners[0]).toBe(key3Callback1);
            expect(key1Callback1).toHaveBeenCalledTimes(1);
            expect(key1Callback1).toHaveBeenCalledWith("key1", "value1", null);
            expect(key1Callback2).toHaveBeenCalledTimes(1);
            expect(key1Callback2).toHaveBeenCalledWith("key1", "value1", null);
            expect(key2Callback1).not.toHaveBeenCalled();
            expect(key2Callback2).not.toHaveBeenCalled();
            expect(key2Callback3).not.toHaveBeenCalled();
            expect(key3Callback1).not.toHaveBeenCalled();
        });

        it("should not do anything and not fail when the provided key is null, undefined or empty string", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key1Callback2 = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1, key1Callback2]
            };
            configHolder.unset("");
            configHolder.unset(null);
            configHolder.unset(undefined);

            expect(Object.keys(configHolder.config)).toHaveLength(1);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(2);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key1.listeners[1]).toBe(key1Callback2);
            expect(key1Callback1).not.toHaveBeenCalled();
            expect(key1Callback2).not.toHaveBeenCalled();
        });

        it("should not notify the listeners if the key or the value for the key does not exist", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key1Callback2 = jest.fn();
            configHolder.config.key1 = {
                value: null,
                listeners: [key1Callback1, key1Callback2]
            };
            configHolder.unset("key1");
            configHolder.unset("non-existent-key1");

            expect(Object.keys(configHolder.config)).toHaveLength(1);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBeNull();
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(2);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key1.listeners[1]).toBe(key1Callback2);
            expect(key1Callback1).not.toHaveBeenCalled();
            expect(key1Callback2).not.toHaveBeenCalled();
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

            expect(configHolder.get("key1")).toBe("value1");
            expect(configHolder.get("key2")).toBe("value2");
            expect(configHolder.get("key3")).toBe("value3");
        });

        it("should return the default value provided if the provided key does not exist", () => {
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
                value: null,
                listeners: []
            };

            expect(configHolder.get("key3", "defaultValue1")).toBeNull();
            expect(configHolder.get("non-existent-key1", "defaultValue1")).toBe("defaultValue1");
            expect(configHolder.get("non-existent-key2", "")).toBe("");
            expect(configHolder.get("non-existent-key3", null)).toBeNull();
            expect(configHolder.get("non-existent-key4", undefined)).toBeNull();
        });

        it("should return null if the provided key does not exist and default value was not provided", () => {
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

    describe("addListener()", () => {
        it("should add the listener to the provided key if the key exists", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            const newCallback = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            configHolder.addListener("key2", newCallback);

            expect(Object.keys(configHolder.config)).toHaveLength(3);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBe("value2");
            expect(configHolder.config.key2.listeners).not.toBeUndefined();
            expect(configHolder.config.key2.listeners).toHaveLength(4);
            expect(configHolder.config.key2.listeners[0]).toBe(key2Callback1);
            expect(configHolder.config.key2.listeners[1]).toBe(key2Callback2);
            expect(configHolder.config.key2.listeners[2]).toBe(key2Callback3);
            expect(configHolder.config.key2.listeners[3]).toBe(newCallback);
            expect(configHolder.config.key3).not.toBeUndefined();
            expect(configHolder.config.key3.value).toBe("value3");
            expect(configHolder.config.key3.listeners).not.toBeUndefined();
            expect(configHolder.config.key3.listeners).toHaveLength(2);
            expect(configHolder.config.key3.listeners[0]).toBe(key3Callback1);
            expect(configHolder.config.key3.listeners[1]).toBe(key3Callback2);
        });

        it("should add the listener to the provided key if it does not exist with value as null", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const newCallback = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            configHolder.addListener("key2", newCallback);

            expect(Object.keys(configHolder.config)).toHaveLength(2);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBeUndefined();
            expect(configHolder.config.key2.listeners).not.toBeUndefined();
            expect(configHolder.config.key2.listeners).toHaveLength(1);
            expect(configHolder.config.key2.listeners[0]).toBe(newCallback);
        });
    });

    describe("removeListener()", () => {
        it("should remove the listener if it exists", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            configHolder.removeListener("key2", key2Callback2);

            expect(Object.keys(configHolder.config)).toHaveLength(3);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBe("value2");
            expect(configHolder.config.key2.listeners).not.toBeUndefined();
            expect(configHolder.config.key2.listeners).toHaveLength(2);
            expect(configHolder.config.key2.listeners[0]).toBe(key2Callback1);
            expect(configHolder.config.key2.listeners[1]).toBe(key2Callback3);
            expect(configHolder.config.key3).not.toBeUndefined();
            expect(configHolder.config.key3.value).toBe("value3");
            expect(configHolder.config.key3.listeners).not.toBeUndefined();
            expect(configHolder.config.key3.listeners).toHaveLength(2);
            expect(configHolder.config.key3.listeners[0]).toBe(key3Callback1);
            expect(configHolder.config.key3.listeners[1]).toBe(key3Callback2);
        });

        it("should not make any changes if the key does not exist", () => {
            const configHolder = new ConfigHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            configHolder.config.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            configHolder.config.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            configHolder.config.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            configHolder.removeListener("non-existent-key1", key2Callback2);
            configHolder.removeListener("non-existent-key2", key2Callback2);
            configHolder.removeListener("non-existent-key3", key2Callback2);

            expect(Object.keys(configHolder.config)).toHaveLength(3);
            expect(configHolder.config.key1).not.toBeUndefined();
            expect(configHolder.config.key1.value).toBe("value1");
            expect(configHolder.config.key1.listeners).not.toBeUndefined();
            expect(configHolder.config.key1.listeners).toHaveLength(1);
            expect(configHolder.config.key1.listeners[0]).toBe(key1Callback1);
            expect(configHolder.config.key2).not.toBeUndefined();
            expect(configHolder.config.key2.value).toBe("value2");
            expect(configHolder.config.key2.listeners).not.toBeUndefined();
            expect(configHolder.config.key2.listeners).toHaveLength(3);
            expect(configHolder.config.key2.listeners[0]).toBe(key2Callback1);
            expect(configHolder.config.key2.listeners[1]).toBe(key2Callback2);
            expect(configHolder.config.key2.listeners[2]).toBe(key2Callback3);
            expect(configHolder.config.key3).not.toBeUndefined();
            expect(configHolder.config.key3.value).toBe("value3");
            expect(configHolder.config.key3.listeners).not.toBeUndefined();
            expect(configHolder.config.key3.listeners).toHaveLength(2);
            expect(configHolder.config.key3.listeners[0]).toBe(key3Callback1);
            expect(configHolder.config.key3.listeners[1]).toBe(key3Callback2);
        });
    });

    describe("loadConfig()", () => {
        it("should return null if the provided key was does not exist and default value was not provided", async () => {
            const configHolder = new ConfigHolder();

            await expect(configHolder.loadConfig()).resolves.toBeUndefined();
            expect(Object.keys(configHolder.config).length > 0).toBe(true);
        });
    });
});
