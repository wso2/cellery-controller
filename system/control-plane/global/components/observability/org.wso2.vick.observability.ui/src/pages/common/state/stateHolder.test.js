
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

/* eslint max-lines: ["off"] */

import StateHolder from "./stateHolder";

describe("ConfigHolder", () => {
    const loggedInUser = "User1";
    localStorage.setItem(StateHolder.USER, loggedInUser);

    const validateInitialState = (stateHolder) => {
        {
            const config = stateHolder.state[StateHolder.CONFIG];
            expect(config).not.toBeUndefined();
            expect(Object.keys(config)).toHaveLength(1);
            expect(config.value).not.toBeUndefined();
        }
        {
            const globalFilter = stateHolder.state[StateHolder.GLOBAL_FILTER];
            expect(globalFilter).not.toBeUndefined();
            expect(Object.keys(globalFilter)).toHaveLength(1);
            expect(globalFilter.value).not.toBeUndefined();
            expect(globalFilter.value.startTime).toBe("now - 24 hours");
            expect(globalFilter.value.endTime).toBe("now");
            expect(globalFilter.value.dateRangeNickname).toBe("Last 24 hours");
            expect(globalFilter.value.refreshInterval).toBe(30 * 1000);
        }
        {
            const loadingState = stateHolder.state[StateHolder.LOADING_STATE];
            expect(loadingState).not.toBeUndefined();
            expect(Object.keys(loadingState)).toHaveLength(1);
            expect(loadingState.value).not.toBeUndefined();
            expect(loadingState.value.loadingOverlayCount).toBe(0);
            expect(loadingState.value.message).toBeNull();
        }
        {
            const notificationState = stateHolder.state[StateHolder.NOTIFICATION_STATE];
            expect(notificationState).not.toBeUndefined();
            expect(Object.keys(notificationState)).toHaveLength(1);
            expect(notificationState.value).not.toBeUndefined();
            expect(notificationState.value.isOpen).toBe(false);
            expect(notificationState.value.message).toBeNull();
            expect(notificationState.value.notificationLevel).toBeNull();
        }
        {
            const user = stateHolder.state[StateHolder.USER];
            expect(user).not.toBeUndefined();
            expect(Object.keys(user)).toHaveLength(1);
            expect(user.value).toBe(loggedInUser);
        }
    };

    describe("set()", () => {
        it("should add the values provided against the keys when a new key is provided", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            stateHolder.set("key2", "value2");
            stateHolder.set("key3", "value3");
            stateHolder.set("key4", null);
            stateHolder.set("key5", undefined);

            expect(Object.keys(stateHolder.state)).toHaveLength(10);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("value2");
            expect(stateHolder.state.key2.listeners).toBeUndefined();
            expect(stateHolder.state.key3).not.toBeUndefined();
            expect(stateHolder.state.key3.value).toBe("value3");
            expect(stateHolder.state.key3.listeners).toBeUndefined();
            expect(stateHolder.state.key4).not.toBeUndefined();
            expect(stateHolder.state.key4.value).toBeNull();
            expect(stateHolder.state.key4.listeners).toBeUndefined();
            expect(stateHolder.state.key5).not.toBeUndefined();
            expect(stateHolder.state.key5.value).toBeUndefined();
            expect(stateHolder.state.key5.listeners).toBeUndefined();
            expect(key1Callback1).not.toHaveBeenCalled();
        });

        it("should not fail upon passing undefined, null or empty string as key", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            stateHolder.set(null, "value1");
            stateHolder.set(null, null);
            stateHolder.set(null, undefined);
            stateHolder.set(undefined, "value2");
            stateHolder.set(undefined, null);
            stateHolder.set(undefined, undefined);
            stateHolder.set("", "value3");
            stateHolder.set("", null);
            stateHolder.set("", undefined);

            expect(Object.keys(stateHolder.state)).toHaveLength(6);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(key1Callback1).not.toHaveBeenCalled();
        });

        it("should set the new value and notify the listeners if an existing key is provided", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            stateHolder.state.key1 = {
                value: null,
                listeners: [key1Callback1]
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            stateHolder.set("key1", "new-value1");
            stateHolder.set("key2", "new-value2");

            expect(Object.keys(stateHolder.state)).toHaveLength(8);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("new-value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("new-value2");
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(3);
            expect(stateHolder.state.key2.listeners[0]).toBe(key2Callback1);
            expect(stateHolder.state.key2.listeners[1]).toBe(key2Callback2);
            expect(stateHolder.state.key2.listeners[2]).toBe(key2Callback3);
            expect(stateHolder.state.key3).not.toBeUndefined();
            expect(stateHolder.state.key3.value).toBe("value3");
            expect(stateHolder.state.key3.listeners).not.toBeUndefined();
            expect(stateHolder.state.key3.listeners).toHaveLength(2);
            expect(stateHolder.state.key3.listeners[0]).toBe(key3Callback1);
            expect(stateHolder.state.key3.listeners[1]).toBe(key3Callback2);
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
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            stateHolder.set("key2", "value2");

            expect(Object.keys(stateHolder.state)).toHaveLength(8);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("value2");
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(3);
            expect(stateHolder.state.key2.listeners[0]).toBe(key2Callback1);
            expect(stateHolder.state.key2.listeners[1]).toBe(key2Callback2);
            expect(stateHolder.state.key2.listeners[2]).toBe(key2Callback3);
            expect(stateHolder.state.key3).not.toBeUndefined();
            expect(stateHolder.state.key3.value).toBe("value3");
            expect(stateHolder.state.key3.listeners).not.toBeUndefined();
            expect(stateHolder.state.key3.listeners).toHaveLength(2);
            expect(stateHolder.state.key3.listeners[0]).toBe(key3Callback1);
            expect(stateHolder.state.key3.listeners[1]).toBe(key3Callback2);
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
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key1Callback2 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1, key1Callback2]
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: [key3Callback1]
            };
            stateHolder.unset("key1");

            expect(Object.keys(stateHolder.state)).toHaveLength(8);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBeNull();
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(2);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key1.listeners[1]).toBe(key1Callback2);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("value2");
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(3);
            expect(stateHolder.state.key2.listeners[0]).toBe(key2Callback1);
            expect(stateHolder.state.key2.listeners[1]).toBe(key2Callback2);
            expect(stateHolder.state.key2.listeners[2]).toBe(key2Callback3);
            expect(stateHolder.state.key3).not.toBeUndefined();
            expect(stateHolder.state.key3.value).toBe("value3");
            expect(stateHolder.state.key3.listeners).not.toBeUndefined();
            expect(stateHolder.state.key3.listeners).toHaveLength(1);
            expect(stateHolder.state.key3.listeners[0]).toBe(key3Callback1);
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
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key1Callback2 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1, key1Callback2]
            };
            stateHolder.unset("");
            stateHolder.unset(null);
            stateHolder.unset(undefined);

            expect(Object.keys(stateHolder.state)).toHaveLength(6);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(2);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key1.listeners[1]).toBe(key1Callback2);
            expect(key1Callback1).not.toHaveBeenCalled();
            expect(key1Callback2).not.toHaveBeenCalled();
        });

        it("should not notify the listeners if the key or the value for the key does not exist", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key1Callback2 = jest.fn();
            stateHolder.state.key1 = {
                value: null,
                listeners: [key1Callback1, key1Callback2]
            };
            stateHolder.unset("key1");
            stateHolder.unset("non-existent-key1");

            expect(Object.keys(stateHolder.state)).toHaveLength(6);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBeNull();
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(2);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key1.listeners[1]).toBe(key1Callback2);
            expect(key1Callback1).not.toHaveBeenCalled();
            expect(key1Callback2).not.toHaveBeenCalled();
        });
    });

    describe("get()", () => {
        it("should return the value stored before for the specified key", () => {
            const stateHolder = new StateHolder();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: []
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: []
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: []
            };

            expect(stateHolder.get("key1", "defaultValue1")).toBe("value1");
            expect(stateHolder.get("key2", "defaultValue2")).toBe("value2");
            expect(stateHolder.get("key3", "defaultValue3")).toBe("value3");
        });

        it("should return the value stored before for the specified key when the default value is not provided", () => {
            const stateHolder = new StateHolder();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: []
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: []
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: []
            };

            expect(stateHolder.get("key1")).toBe("value1");
            expect(stateHolder.get("key2")).toBe("value2");
            expect(stateHolder.get("key3")).toBe("value3");
        });

        it("should return the default value provided if the provided key does not exist", () => {
            const stateHolder = new StateHolder();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: []
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: []
            };
            stateHolder.state.key3 = {
                value: null,
                listeners: []
            };

            expect(stateHolder.get("key3", "defaultValue1")).toBeNull();
            expect(stateHolder.get("non-existent-key1", "defaultValue1")).toBe("defaultValue1");
            expect(stateHolder.get("non-existent-key2", "")).toBe("");
            expect(stateHolder.get("non-existent-key3", null)).toBeNull();
            expect(stateHolder.get("non-existent-key4", undefined)).toBeNull();
        });

        it("should return null if the provided key does not exist and default value was not provided", () => {
            const stateHolder = new StateHolder();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: []
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: []
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: []
            };

            expect(stateHolder.get("non-existent-key1")).toBeNull();
            expect(stateHolder.get("non-existent-key2")).toBeNull();
            expect(stateHolder.get("non-existent-key3")).toBeNull();
        });
    });

    describe("addListener()", () => {
        it("should add the listener to the provided key if the key exists", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            const newCallback = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            stateHolder.addListener("key2", newCallback);

            expect(Object.keys(stateHolder.state)).toHaveLength(8);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("value2");
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(4);
            expect(stateHolder.state.key2.listeners[0]).toBe(key2Callback1);
            expect(stateHolder.state.key2.listeners[1]).toBe(key2Callback2);
            expect(stateHolder.state.key2.listeners[2]).toBe(key2Callback3);
            expect(stateHolder.state.key2.listeners[3]).toBe(newCallback);
            expect(stateHolder.state.key3).not.toBeUndefined();
            expect(stateHolder.state.key3.value).toBe("value3");
            expect(stateHolder.state.key3.listeners).not.toBeUndefined();
            expect(stateHolder.state.key3.listeners).toHaveLength(2);
            expect(stateHolder.state.key3.listeners[0]).toBe(key3Callback1);
            expect(stateHolder.state.key3.listeners[1]).toBe(key3Callback2);
        });

        it("should add the listener to the provided key if it does not exist with value as null", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const newCallback = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            stateHolder.addListener("key2", newCallback);

            expect(Object.keys(stateHolder.state)).toHaveLength(7);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBeUndefined();
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(1);
            expect(stateHolder.state.key2.listeners[0]).toBe(newCallback);
        });
    });

    describe("removeListener()", () => {
        it("should remove the listener if it exists", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            stateHolder.removeListener("key2", key2Callback2);

            expect(Object.keys(stateHolder.state)).toHaveLength(8);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("value2");
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(2);
            expect(stateHolder.state.key2.listeners[0]).toBe(key2Callback1);
            expect(stateHolder.state.key2.listeners[1]).toBe(key2Callback3);
            expect(stateHolder.state.key3).not.toBeUndefined();
            expect(stateHolder.state.key3.value).toBe("value3");
            expect(stateHolder.state.key3.listeners).not.toBeUndefined();
            expect(stateHolder.state.key3.listeners).toHaveLength(2);
            expect(stateHolder.state.key3.listeners[0]).toBe(key3Callback1);
            expect(stateHolder.state.key3.listeners[1]).toBe(key3Callback2);

            stateHolder.set("key2", "newValue");

            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("newValue");
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(2);
            expect(stateHolder.state.key2.listeners[0]).toBe(key2Callback1);
            expect(stateHolder.state.key2.listeners[1]).toBe(key2Callback3);
            expect(key2Callback1).toHaveBeenCalled();
            expect(key2Callback2).not.toHaveBeenCalled();
            expect(key2Callback3).toHaveBeenCalled();
        });

        it("should not make any changes if the key does not exist", () => {
            const stateHolder = new StateHolder();
            const key1Callback1 = jest.fn();
            const key2Callback1 = jest.fn();
            const key2Callback2 = jest.fn();
            const key2Callback3 = jest.fn();
            const key3Callback1 = jest.fn();
            const key3Callback2 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1",
                listeners: [key1Callback1]
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: [key2Callback1, key2Callback2, key2Callback3]
            };
            stateHolder.state.key3 = {
                value: "value3",
                listeners: [key3Callback1, key3Callback2]
            };
            stateHolder.removeListener("non-existent-key1", key2Callback2);
            stateHolder.removeListener("non-existent-key2", key2Callback2);
            stateHolder.removeListener("non-existent-key3", key2Callback2);

            expect(Object.keys(stateHolder.state)).toHaveLength(8);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).not.toBeUndefined();
            expect(stateHolder.state.key1.listeners).toHaveLength(1);
            expect(stateHolder.state.key1.listeners[0]).toBe(key1Callback1);
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("value2");
            expect(stateHolder.state.key2.listeners).not.toBeUndefined();
            expect(stateHolder.state.key2.listeners).toHaveLength(3);
            expect(stateHolder.state.key2.listeners[0]).toBe(key2Callback1);
            expect(stateHolder.state.key2.listeners[1]).toBe(key2Callback2);
            expect(stateHolder.state.key2.listeners[2]).toBe(key2Callback3);
            expect(stateHolder.state.key3).not.toBeUndefined();
            expect(stateHolder.state.key3.value).toBe("value3");
            expect(stateHolder.state.key3.listeners).not.toBeUndefined();
            expect(stateHolder.state.key3.listeners).toHaveLength(2);
            expect(stateHolder.state.key3.listeners[0]).toBe(key3Callback1);
            expect(stateHolder.state.key3.listeners[1]).toBe(key3Callback2);
        });

        it("should not not fail and should not make any changes if the listeners list was not yet added", () => {
            const stateHolder = new StateHolder();
            const callback1 = jest.fn();
            const callback2 = jest.fn();
            stateHolder.state.key1 = {
                value: "value1"
            };
            stateHolder.state.key2 = {
                value: "value2",
                listeners: []
            };
            stateHolder.removeListener("key1", callback1);
            stateHolder.removeListener("key2", callback2);

            expect(Object.keys(stateHolder.state)).toHaveLength(7);
            validateInitialState(stateHolder);
            expect(stateHolder.state.key1).not.toBeUndefined();
            expect(stateHolder.state.key1.value).toBe("value1");
            expect(stateHolder.state.key1.listeners).toBeUndefined();
            expect(stateHolder.state.key2).not.toBeUndefined();
            expect(stateHolder.state.key2.value).toBe("value2");
            expect(stateHolder.state.key2.listeners).toEqual([]);
        });
    });
});
