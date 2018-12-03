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

import AuthUtils from "../utils/authUtils";

/**
 * Configuration holder.
 */
class StateHolder {

    static USER = "user";
    static LOADING_STATE = "loadingState";
    static NOTIFICATION_STATE = "notificationState";
    static CONFIG = "config";
    static GLOBAL_FILTER = "globalFilter";

    static NotificationLevels = {
        INFO: "INFO",
        WARNING: "WARNING",
        ERROR: "ERROR"
    };

    /**
     * @type {Object}
     * @private
     */
    state = {};

    constructor() {
        const rawState = {
            user: AuthUtils.getAuthenticatedUser(),
            loadingState: {
                isLoading: false,
                message: null
            },
            notificationState: {
                isOpen: false,
                message: null,
                notificationLevel: null
            },
            config: {},
            globalFilter: {
                startTime: "now - 1 year",
                endTime: "now",
                dateRangeNickname: "Last 1 year",
                refreshInterval: 30 * 1000 // 30 milli-seconds
            }
        };

        const initialState = {};
        for (const stateKey in rawState) {
            if (rawState.hasOwnProperty(stateKey)) {
                initialState[stateKey] = {
                    value: rawState[stateKey],
                    listeners: []
                };
            }
        }

        this.state = initialState;
    }

    /**
     * Set the value for a particular key.
     *
     * @param {string} key The key for which the value should be added
     * @param {Object} value The new value that should be set
     */
    set = (key, value) => {
        if (key) {
            if (!this.state[key]) {
                this.state[key] = {
                    value: null
                };
            }
            this.notify(key, value);
            this.state[key].value = value;
        }
    };

    /**
     * Unset the value for a particular key.
     *
     * @param {string} key The key for which the value should be removed
     */
    unset = (key) => {
        if (key && this.state[key]) {
            this.notify(key, null);
            this.state[key].value = null;
        }
    };

    /**
     * Get the value for a particular key.
     *
     * @param {string} key The key for which the value should be retrieved
     * @param {Object} defaultValue The default value which should be returned if the value does not exist
     * @returns {Object} The value for the key provided
     */
    get = (key, defaultValue = null) => {
        let value = defaultValue;
        if (this.state[key]) {
            value = this.state[key].value;
        }
        return value;
    };

    /**
     * Add a listener for a particular state key.
     *
     * @param {string} key The state key for which the listener should be added
     * @param {Function} callback The callback function which should be called upon update
     */
    addListener = (key, callback) => {
        if (!this.state[key]) {
            this.state[key] = {};
        }
        if (!this.state[key].listeners) {
            this.state[key].listeners = [];
        }
        this.state[key].listeners.push(callback);
    };

    /**
     * Remove a listener previously added.
     *
     * @param {string} key The key from which the listener should be added
     * @param {Function} callback The callback which should be removed
     */
    removeListener = (key, callback) => {
        if (this.state[key]) {
            const listeners = this.state[key].listeners;
            if (listeners) {
                const removeIndex = listeners.indexOf(callback);
                listeners.splice(removeIndex, 1);
            }
        }
    };

    /**
     * Notify the listeners about a state change.
     *
     * @param {string} key The key of which the listeners should be notified
     * @param {Object} newValue The new value of the key
     * @private
     */
    notify = (key, newValue) => {
        const oldValue = this.state[key].value;
        const listeners = this.state[key].listeners;
        if (oldValue !== newValue && listeners) {
            listeners.forEach((listener) => listener(key, oldValue, newValue));
        }
    };

    /**
     * Load the state that should be used.
     *
     * @returns {Promise<Object>} Promise which resolves when the state is loaded or rejects
     */
    loadConfig = () => {
        const self = this;
        return new Promise((resolve) => {
            // TODO : Load configuration from server
            const loadedConfiguration = {
                backendURL: "http://localhost:8080/api"
            };
            self.state.config = {
                value: loadedConfiguration,
                listeners: []
            };
            resolve();
        });
    };

}

export default StateHolder;
