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

import NotFound from "../common/NotFound";
import PropTypes from "prop-types";
import React from "react";
import Search from "./Search";
import View from "./view";
import {Redirect, Route, Switch, withRouter} from "react-router-dom";

const Tracing = ({match, location}) => (
    <Switch>
        <Route exact path={`${match.path}/search`} component={Search}/>
        <Route exact path={`${match.path}/id/:traceId`} component={View}/>
        <Redirect exact from={`${match.url}/`} to={{pathname: `${match.url}/search`, state: location.state}}/>
        <Route path={`${match.url}/*`} component={NotFound}/>
    </Switch>
);

Tracing.propTypes = {
    match: PropTypes.object.isRequired,
    location: PropTypes.shape({
        state: PropTypes.object
    }).isRequired
};

export default withRouter(Tracing);
