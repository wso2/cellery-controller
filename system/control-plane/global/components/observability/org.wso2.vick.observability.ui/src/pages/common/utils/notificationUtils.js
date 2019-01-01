/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
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

import {StateHolder} from "../state";

class NotificationUtils {

    static Levels = {
        INFO: "INFO",
        WARNING: "WARNING",
        ERROR: "ERROR"
    };

    /**
     * Show the loading overlay.
     *
     * @param {string} message The message to be shown in the loading overlay
     * @param {StateHolder} globalState The global state provided to the current component
     */
    static showLoadingOverlay = (message, globalState) => {
        const prevState = globalState.get(StateHolder.LOADING_STATE);
        globalState.set(StateHolder.LOADING_STATE, {
            loadingOverlayCount: prevState.loadingOverlayCount + 1,
            message: message
        });
    };

    /**
     * Hide the loading overlay.
     *
     * @param {StateHolder} globalState The global state provided to the current component
     */
    static hideLoadingOverlay = (globalState) => {
        const prevState = globalState.get(StateHolder.LOADING_STATE);
        globalState.set(StateHolder.LOADING_STATE, {
            loadingOverlayCount: prevState.loadingOverlayCount === 0 ? 0 : prevState.loadingOverlayCount - 1,
            message: null
        });
    };

    /**
     * Check if the loading overlay is currently shown.
     *
     * @param {StateHolder} globalState The global state provided to the current component
     * @returns {boolean} True if the loading state is currently shown
     */
    static isLoadingOverlayShown = (globalState) => globalState.get(StateHolder.LOADING_STATE).loadingOverlayCount > 0;

    /**
     * Show a notification to the user.
     *
     * @param {string} message The message to be shown in the notification
     * @param {string} level The notification level
     * @param {StateHolder} globalState The global state provided to the current component
     */
    static showNotification = (message, level, globalState) => {
        globalState.set(StateHolder.NOTIFICATION_STATE, {
            isOpen: true,
            message: message,
            notificationLevel: level
        });
    };

    /**
     * Close the notification shown to the user.
     *
     * @param {StateHolder} globalState The global state provided to the current component
     */
    static closeNotification = (globalState) => {
        globalState.set(StateHolder.NOTIFICATION_STATE, {
            isOpen: false,
            message: null,
            notificationLevel: null
        });
    };

}

export default NotificationUtils;
