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

import PropTypes from "prop-types";
import React from "react";

// Creating a context that can be accessed
const ConfigContext = React.createContext({});

/**
 * Config Provider to provide the configuration.
 *
 * @param {Object} props Props passed into the config provider
 * @returns {React.Component} Color Provider React Component
 * @constructor
 */
class ConfigProvider extends React.Component {

    render() {
        const {children, config} = this.props;

        return (
            <ConfigContext.Provider value={config}>
                {children}
            </ConfigContext.Provider>
        );
    }

}

ConfigProvider.propTypes = {
    children: PropTypes.any.isRequired,
    config: PropTypes.any.isRequired
};

/**
 * Higher Order Component for accessing the Color Generator.
 *
 * @param {React.Component} Component component which needs access to the configuration.
 * @returns {Object} The new HOC with access to the configuration.
 */
const withConfig = (Component) => class ColorGeneratorProvider extends React.Component {

    render() {
        return (
            <ConfigContext.Consumer>
                {(config) => <Component config={config} {...this.props}/>}
            </ConfigContext.Consumer>
        );
    }

};

export {withConfig, ConfigProvider};
