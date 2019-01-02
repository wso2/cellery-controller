
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

import NotificationUtils from "./notificationUtils";
import StateHolder from "../state/stateHolder";

describe("NotificationUtils", () => {
    let state;

    beforeEach(() => {
        state = new StateHolder();
    });

    describe("showLoadingOverlay()", () => {
        it("should set the message and make the loading state true when called with the initial state", () => {
            NotificationUtils.showLoadingOverlay("Test Message 1", state);

            const loadingState = state.get(StateHolder.LOADING_STATE);
            expect(loadingState.loadingOverlayCount).toBe(1);
            expect(loadingState.message).toBe("Test Message 1");
        });

        it("should set the message and make the loading state true if the loading is currently hidden", () => {
            state.set(StateHolder.LOADING_STATE, {
                loadingOverlayCount: 0,
                message: null
            });
            NotificationUtils.showLoadingOverlay("Test Message 2", state);

            const loadingState = state.get(StateHolder.LOADING_STATE);
            expect(loadingState.loadingOverlayCount).toBe(1);
            expect(loadingState.message).toBe("Test Message 2");
        });

        it("should change to the new message and the loading state unchanged if called when overlay is shown", () => {
            state.set(StateHolder.LOADING_STATE, {
                loadingOverlayCount: 103,
                message: "Initial Test Message 3"
            });
            NotificationUtils.showLoadingOverlay("Test Message 3", state);

            const loadingState = state.get(StateHolder.LOADING_STATE);
            expect(loadingState.loadingOverlayCount).toBe(104);
            expect(loadingState.message).toBe("Test Message 3");
        });
    });

    describe("hideLoadingOverlay()", () => {
        it("should set the message to null and make the loading state false when applied on the initial state", () => {
            NotificationUtils.hideLoadingOverlay(state);

            const loadingState = state.get(StateHolder.LOADING_STATE);
            expect(loadingState.loadingOverlayCount).toBe(0);
            expect(loadingState.message).toBeNull();
        });

        it("should keep the message and loading state unchanged if the loading is currently hidden", () => {
            state.set(StateHolder.LOADING_STATE, {
                loadingOverlayCount: 0,
                message: null
            });
            NotificationUtils.hideLoadingOverlay(state);

            const loadingState = state.get(StateHolder.LOADING_STATE);
            expect(loadingState.loadingOverlayCount).toBe(0);
            expect(loadingState.message).toBeNull();
        });

        it("should set the message to null and loading state to false if called when overlay is shown", () => {
            state.set(StateHolder.LOADING_STATE, {
                loadingOverlayCount: 214,
                message: "Initial Test Message 3"
            });
            NotificationUtils.hideLoadingOverlay(state);

            const loadingState = state.get(StateHolder.LOADING_STATE);
            expect(loadingState.loadingOverlayCount).toBe(213);
            expect(loadingState.message).toBeNull();
        });
    });

    describe("isLoadingOverlayShown()", () => {
        it("should return true if the overlay is currently shown", () => {
            state.set(StateHolder.LOADING_STATE, {
                loadingOverlayCount: 3,
                message: "Initial Test Message 3"
            });

            expect(NotificationUtils.isLoadingOverlayShown(state)).toBe(true);
        });

        it("should return false if the overlay is currently not shown", () => {
            state.set(StateHolder.LOADING_STATE, {
                loadingOverlayCount: 0,
                message: 0
            });

            expect(NotificationUtils.isLoadingOverlayShown(state)).toBe(false);
        });
    });

    describe("showNotification()", () => {
        it("should set the message, notification level and state when applied on the initial state", () => {
            NotificationUtils.showNotification("Test Message 1", NotificationUtils.Levels.INFO, state);

            const notificationState = state.get(StateHolder.NOTIFICATION_STATE);
            expect(notificationState.isOpen).toBe(true);
            expect(notificationState.message).toBe("Test Message 1");
            expect(notificationState.notificationLevel).toBe(NotificationUtils.Levels.INFO);
        });

        it("should set the message, notification level and state if the notification is currently hidden", () => {
            state.set(StateHolder.NOTIFICATION_STATE, {
                isOpen: false,
                message: null,
                notificationLevel: null
            });
            NotificationUtils.showNotification("Test Message 2", NotificationUtils.Levels.WARNING, state);

            const notificationState = state.get(StateHolder.NOTIFICATION_STATE);
            expect(notificationState.isOpen).toBe(true);
            expect(notificationState.message).toBe("Test Message 2");
            expect(notificationState.notificationLevel).toBe(NotificationUtils.Levels.WARNING);
        });

        it("should change the message and notification level if called when notification is already shown", () => {
            state.set(StateHolder.LOADING_STATE, {
                isOpen: true,
                message: "Initial Test Message 3",
                notificationLevel: NotificationUtils.Levels.INFO
            });
            NotificationUtils.showNotification("Test Message 3", NotificationUtils.Levels.ERROR, state);

            const notificationState = state.get(StateHolder.NOTIFICATION_STATE);
            expect(notificationState.isOpen).toBe(true);
            expect(notificationState.message).toBe("Test Message 3");
            expect(notificationState.notificationLevel).toBe(NotificationUtils.Levels.ERROR);
        });
    });

    describe("closeNotification()", () => {
        it("should remove the message, notification level and state when applied on the initial state", () => {
            NotificationUtils.closeNotification(state);

            const notificationState = state.get(StateHolder.NOTIFICATION_STATE);
            expect(notificationState.isOpen).toBe(false);
            expect(notificationState.message).toBeNull();
            expect(notificationState.notificationLevel).toBeNull();
        });

        it("should remove the message, notification level and state if the notification is currently hidden", () => {
            state.set(StateHolder.NOTIFICATION_STATE, {
                isOpen: false,
                message: null,
                notificationLevel: null
            });
            NotificationUtils.closeNotification(state);

            const notificationState = state.get(StateHolder.NOTIFICATION_STATE);
            expect(notificationState.isOpen).toBe(false);
            expect(notificationState.message).toBeNull();
            expect(notificationState.notificationLevel).toBeNull();
        });

        it("should remove the message and notification level if called when notification is already shown", () => {
            state.set(StateHolder.LOADING_STATE, {
                isOpen: true,
                message: "Initial Test Message 3",
                notificationLevel: NotificationUtils.Levels.INFO
            });
            NotificationUtils.closeNotification(state);

            const notificationState = state.get(StateHolder.NOTIFICATION_STATE);
            expect(notificationState.isOpen).toBe(false);
            expect(notificationState.message).toBeNull();
            expect(notificationState.notificationLevel).toBeNull();
        });
    });
});
