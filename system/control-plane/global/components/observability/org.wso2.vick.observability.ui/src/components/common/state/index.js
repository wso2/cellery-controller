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

/* eslint react/prefer-stateless-function: ["off"] */

import CircularProgress from "@material-ui/core/CircularProgress/CircularProgress";
import Grid from "@material-ui/core/Grid/Grid";
import HttpUtils from "../../../utils/api/httpUtils";
import NotFound from "../error/NotFound";
import React from "react";
import StateHolder from "./stateHolder";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core";
import * as PropTypes from "prop-types";

// Creating a context that can be accessed
const StateContext = React.createContext({});

const styles = () => ({
    container: {
        minHeight: "100%",
        bottom: 0
    }
});

class UnStyledStateProvider extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isLoading: true,
            isConfigAvailable: false
        };

        this.mounted = false;

        const queryParams = HttpUtils.parseQueryParams(props.location.search);
        this.stateHolder = new StateHolder(queryParams);
    }

    componentDidMount = () => {
        const self = this;
        self.mounted = true;
        self.stateHolder.loadConfig()
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
    };

    componentWillUnmount = () => {
        this.mounted = false;
    };

    render = () => {
        const {children, classes} = this.props;
        const {isLoading, isConfigAvailable} = this.state;

        const content = (
            isConfigAvailable
                ? children
                : <NotFound title={"Failed to Load Cellery Observability Portal Configuration"}/>
        );
        return (
            <StateContext.Provider value={this.stateHolder}>
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
            </StateContext.Provider>
        );
    };

}

UnStyledStateProvider.propTypes = {
    children: PropTypes.any.isRequired,
    classes: PropTypes.any.isRequired,
    location: PropTypes.shape({
        search: PropTypes.string.isRequired
    }).isRequired
};

const StateProvider = withStyles(styles, {withTheme: true})(withRouter(UnStyledStateProvider));

/**
 * Higher Order Component for accessing the State Holder.
 *
 * @param {React.ComponentType} Component component which needs access to the state.
 * @returns {React.ComponentType} The new HOC with access to the state.
 */
const withGlobalState = (Component) => {
    class StateConsumer extends React.Component {

        render = () => {
            const {forwardedRef, ...otherProps} = this.props;

            return (
                <StateContext.Consumer>
                    {(state) => <Component globalState={state} ref={forwardedRef} {...otherProps}/>}
                </StateContext.Consumer>
            );
        };

    }

    return React.forwardRef((props, ref) => <StateConsumer {...props} forwardedRef={ref} />);
};

export default withGlobalState;
export {StateProvider, StateHolder};
