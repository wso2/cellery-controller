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

import ColorGenerator from "./colorGenerator";
import PropTypes from "prop-types";
import React from "react";

// Creating a context that can be accessed
const ColorContext = React.createContext(null);

/**
 * Color Provider to provide the color generator.
 *
 * @param {Object} props Props passed into the color provider
 * @returns {React.Component} Color Provider React Component
 * @constructor
 */
class ColorProvider extends React.Component {

    render() {
        this.colorGenerator = new ColorGenerator();

        return (
            <ColorContext.Provider value={this.colorGenerator}>
                {this.props.children}
            </ColorContext.Provider>
        );
    }

}

ColorProvider.propTypes = {
    children: PropTypes.any.isRequired
};

/**
 * Higher Order Component for accessing the Color Generator.
 *
 * @param {React.Component} Component component which needs access to the color generator.
 * @returns {Object} The new HOC with access to the color generator.
 */
const withColor = (Component) => class ColorGeneratorProvider extends React.Component {

    render() {
        return (
            <ColorContext.Consumer>
                {(colorGenerator) => <Component colorGenerator={colorGenerator} {...this.props}/>}
            </ColorContext.Consumer>
        );
    }

};

export {withColor, ColorProvider};
