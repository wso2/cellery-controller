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

import CircularProgress from "@material-ui/core/CircularProgress/CircularProgress";
import Grid from "@material-ui/core/Grid/Grid";
import PropTypes from "prop-types";
import React from "react";
import {withStyles} from "@material-ui/core";
import {ConfigConstants, ConfigHolder} from "./configHolder";

// Creating a context that can be accessed
const ConfigContext = React.createContext({});

const styles = () => ({
    container: {
        minHeight: "100%",
        bottom: 0
    }
});

/**
 * Config Provider to provide the configuration.
 *
 * @param {Object} props Props passed into the config provider
 * @returns {React.Component} Color Provider React Component
 * @constructor
 */
class UnStyledConfigProvider extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isLoading: true,
            isConfigAvailable: false
        };

        this.mounted = false;
        this.config = new ConfigHolder();
    }

    componentDidMount() {
        const self = this;
        self.mounted = true;
        self.config.loadConfig()
            .then(() => {
                if (self.mounted) {
                    self.setState({
                        isLoading: false,
                        isConfigAvailable: true
                    });
                }
            })
            .catch(() => {
                if (self.mounted) {
                    self.setState({
                        isLoading: false
                    });
                }
            });
    }

    componentWillUnmount() {
        this.mounted = false;
    }

    render() {
        const {children, classes} = this.props;
        const {isLoading, isConfigAvailable} = this.state;

        const content = (isConfigAvailable ? children : <div>Error</div>);
        return (
            <ConfigContext.Provider value={this.config}>
                {
                    isLoading
                        ? (
                            <Grid container justify="center" alignItems="center"
                                className={classes.container}>
                                <Grid item>
                                    <CircularProgress size={60}/>
                                </Grid>
                            </Grid>
                        )
                        : content
                }
            </ConfigContext.Provider>
        );
    }

}

UnStyledConfigProvider.propTypes = {
    children: PropTypes.any.isRequired,
    classes: PropTypes.any.isRequired
};

const ConfigProvider = withStyles(styles, {withTheme: true})(UnStyledConfigProvider);

/**
 * Higher Order Component for accessing the Color Generator.
 *
 * @param {React.Component | Function} Component component which needs access to the configuration.
 * @returns {React.Component | Function} The new HOC with access to the configuration.
 */
const withConfig = (Component) => class ConfigConsumer extends React.Component {

    render() {
        return (
            <ConfigContext.Consumer>
                {(config) => <Component config={config} {...this.props}/>}
            </ConfigContext.Consumer>
        );
    }

};

export {withConfig, ConfigProvider, ConfigConstants};
