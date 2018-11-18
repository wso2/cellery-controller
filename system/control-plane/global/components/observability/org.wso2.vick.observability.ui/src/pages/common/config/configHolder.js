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

import Utils from "../utils";

/**
 * Configuration holder.
 */
class ConfigHolder {

    /**
     * @type {Object}
     * @private
     */
    config = {};

    set(key, value) {
        if (key && value) {
            this.config[key] = value;
        }
    }

    get(key, defaultValue) {
        let value = (defaultValue ? defaultValue : null);
        if (this.config[key]) {
            value = this.config[key];
        }
        return value;
    }

    loadConfig() {
        const self = this;
        return new Promise((resolve) => {
            // TODO : Load configuration from server
            self.config = {
                user: Utils.getAuthenticatedUser(),
                backendURL: "wso2sp-worker.vick-system:9000",
                globalFilter: {
                    startTime: "now-24h",
                    endTime: "now",
                    refreshInterval: 30 * 1000 // 30 milli-seconds
                }
            };
            resolve(self.config);
        });
    }

}

const ConfigConstants = {
    GLOBAL_FILTER: "globalFilter",
    USER: "user",
    BACKEND_URL: "backendURL"
};

export {ConfigHolder, ConfigConstants};
