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

import AppLayout from "./appLayout";
import Cells from "./cells";
import {ColorProvider} from "./common/color";
import CssBaseline from "@material-ui/core/CssBaseline/CssBaseline";
import ErrorBoundary from "./common/error/ErrorBoundary";
import NotFound from "./common/error/NotFound";
import Overview from "./overview";
import React from "react";
import SignIn from "./SignIn";
import SystemMetrics from "./systemMetrics";
import Tracing from "./tracing";
import {BrowserRouter, Route, Switch} from "react-router-dom";
import {MuiThemeProvider, createMuiTheme} from "@material-ui/core/styles";
import withGlobalState, {StateHolder, StateProvider} from "./common/state";
import * as PropTypes from "prop-types";

class StatelessProtectedPortal extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isAuthenticated: Boolean(props.globalState.get(StateHolder.USER))
        };

        props.globalState.addListener(StateHolder.USER, this.handleUserChange);
    }

    handleUserChange = (userKey, oldUser, newUser) => {
        this.setState({
            isAuthenticated: Boolean(newUser)
        });
    };

    render = () => {
        const {isAuthenticated} = this.state;
        return isAuthenticated
            ? (
                <AppLayout>
                    <ErrorBoundary>
                        <Switch>
                            <Route exact path="/" component={Overview}/>
                            <Route path="/cells" component={Cells}/>
                            <Route path="/tracing" component={Tracing}/>
                            <Route path="/system-metrics" component={SystemMetrics}/>
                            <Route path="/*" component={NotFound}/>
                        </Switch>
                    </ErrorBoundary>
                </AppLayout>
            )
            : <SignIn/>;
    };

}

StatelessProtectedPortal.propTypes = {
    globalState: PropTypes.instanceOf(StateHolder).isRequired
};

const ProtectedPortal = withGlobalState(StatelessProtectedPortal);

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
        <CssBaseline/>
        <ErrorBoundary>
            <BrowserRouter>
                <ColorProvider>
                    <StateProvider>
                        <ProtectedPortal/>
                    </StateProvider>
                </ColorProvider>
            </BrowserRouter>
        </ErrorBoundary>
    </MuiThemeProvider>
);

export default App;
