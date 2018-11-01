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

import React, {Component} from 'react';
import AppLayout from './AppLayout';
import Cell from "./pages/Cell";
import MicroService from "./pages/MicroService";
import Overview from "./pages/Overview";
import SignIn from './pages/SignIn';
import {BrowserRouter, Redirect, Route, Switch} from "react-router-dom";
import Tracing from "./pages/tracing/Tracing";

class App extends Component {
    render() {
        const user = this.props.username ? this.props.username : localStorage.getItem("username");
        return (
            <BrowserRouter>
                {
                    user
                        ? (
                            <AppLayout username={user}>
                                <Switch>
                                    <Route exact path="/" component={Overview}/>
                                    <Route exact path="/cells" component={Cell}/>
                                    <Route exact path="/micro-services" component={MicroService}/>
                                    <Route path="/tracing" component={Tracing}/>
                                    <Redirect from="/*" to="/"/>
                                </Switch>
                            </AppLayout>
                        )
                        : <SignIn/>
                }
            </BrowserRouter>
        )
    }
}

export default App;
