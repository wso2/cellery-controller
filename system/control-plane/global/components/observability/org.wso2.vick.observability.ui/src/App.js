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
import SignIn from './pages/SignIn';
import {BrowserRouter, Route} from 'react-router-dom';

class App extends Component {
    render() {
        const user = this.props.username ? this.props.username : localStorage.getItem("username");
        return (
            <BrowserRouter>
                <Route path='/' render={user ? props => <AppLayout {...props} username={user}/> : props =>
                    <SignIn {...props} />}>
                </Route>
            </BrowserRouter>
        )
    }
}

export default App;
