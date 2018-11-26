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

import AppLayout from "./AppLayout";
import Cells from "./pages/cells";
import {ColorProvider} from "./pages/common/color";
import {ConfigHolder} from "./pages/common/config/configHolder";
import NotFound from "./pages/common/NotFound";
import Overview from "./pages/Overview";
import PropTypes from "prop-types";
import React from "react";
import SignIn from "./pages/SignIn";
import SystemMetrics from "./pages/systemMetrics";
import Tracing from "./pages/tracing";
import {BrowserRouter, Route, Switch} from "react-router-dom";
import {ConfigConstants, ConfigProvider, withConfig} from "./pages/common/config";
import {MuiThemeProvider, createMuiTheme} from "@material-ui/core/styles";

/**
 * Protected Portal of the App.
 *
 * @param {Object} props Props passed to the component
 * @returns {React.Component} Protected portal react component
 */
class UnConfiguredProtectedPortal extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isAuthenticated: Boolean(props.config.get(ConfigConstants.USER))
        };

        this.handleUserChange = this.handleUserChange.bind(this);

        props.config.addListener(ConfigConstants.USER, this.handleUserChange);
    }

    handleUserChange(userKey, oldUser, newUser) {
        this.setState({
            isAuthenticated: Boolean(newUser)
        });
    }

    render() {
        const {isAuthenticated} = this.state;
        return isAuthenticated
            ? (
                <AppLayout>
                    <Switch>
                        <Route exact path="/" component={Overview}/>
                        <Route path="/cells" component={Cells}/>
                        <Route path="/tracing" component={Tracing}/>
                        <Route path="/system-metrics" component={SystemMetrics}/>
                        <Route path="/*" component={NotFound}/>
                    </Switch>
                </AppLayout>
            )
            : <SignIn/>;
    }

}

UnConfiguredProtectedPortal.propTypes = {
    config: PropTypes.instanceOf(ConfigHolder).isRequired
};

const ProtectedPortal = withConfig(UnConfiguredProtectedPortal);

// Create the main theme of the App
const theme = createMuiTheme({
    typography: {
        useNextVariants: true
    }
});

/**
 * The Observability Main App.
 *
 * @returns {React.Component} App react component
 */
const App = () => (
    <MuiThemeProvider theme={theme}>
        <ColorProvider>
            <ConfigProvider>
                <BrowserRouter>
                    <ProtectedPortal/>
                </BrowserRouter>
            </ConfigProvider>
        </ColorProvider>
    </MuiThemeProvider>
);

export default App;
