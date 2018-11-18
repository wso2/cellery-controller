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
import Cell from "./pages/Cell";
import {ColorProvider} from "./pages/common/color";
import Microservice from "./pages/Microservice";
import Overview from "./pages/Overview";
import PropTypes from "prop-types";
import React from "react";
import SignIn from "./pages/SignIn";
import Tracing from "./pages/tracing";
import {BrowserRouter, Redirect, Route, Switch} from "react-router-dom";
import {ConfigConstants, ConfigHolder} from "./pages/common/config/configHolder";
import {ConfigProvider, withConfig} from "./pages/common/config";
import {MuiThemeProvider, createMuiTheme} from "@material-ui/core/styles";

/**
 * Protected Portal of the App.
 *
 * @param {Object} props Props passed to the component
 * @returns {React.Component} Protected portal react component
 */
let ProtectedPortal = (props) => (
    props.config.get(ConfigConstants.USER)
        ? (
            <AppLayout>
                <Switch>
                    <Route exact path="/" component={Overview}/>
                    <Route exact path="/cells" component={Cell}/>
                    <Route exact path="/microservices" component={Microservice}/>
                    <Route path="/tracing" component={Tracing}/>
                    <Redirect from="/*" to="/"/>
                </Switch>
            </AppLayout>
        )
        : <SignIn/>
);

ProtectedPortal.propTypes = {
    config: PropTypes.instanceOf(ConfigHolder).isRequired
};

ProtectedPortal = withConfig(ProtectedPortal);

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
