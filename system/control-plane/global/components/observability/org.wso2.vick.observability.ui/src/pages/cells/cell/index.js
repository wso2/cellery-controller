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

import Details from "./Details";
import Grey from "@material-ui/core/colors/grey";
import HttpUtils from "../../common/utils/httpUtils";
import Metrics from "./Metrics";
import MicroserviceList from "./MicroserviceList";
import Paper from "@material-ui/core/Paper/Paper";
import PropTypes from "prop-types";
import React from "react";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import TopToolbar from "../../common/toptoolbar";
import {withStyles} from "@material-ui/core/styles";

const styles = (theme) => ({
    root: {
        flexGrow: 1,
        backgroundColor: theme.palette.background.paper,
        padding: theme.spacing.unit * 3,
        paddingTop: 0,
        margin: Number(theme.spacing.unit)
    },
    tabs: {
        marginBottom: theme.spacing.unit * 2,
        borderBottomWidth: 1,
        borderBottomStyle: "solid",
        borderBottomColor: Grey[200]
    }
});

class Cell extends React.Component {

    constructor(props) {
        super(props);

        this.tabs = [
            "details",
            "microservices",
            "metrics"
        ];
        const queryParams = HttpUtils.parseQueryParams(props.location.search);
        const preSelectedTab = queryParams.tab ? this.tabs.indexOf(queryParams.tab) : null;

        this.state = {
            selectedTabIndex: (preSelectedTab ? preSelectedTab : 0)
        };
    }

    handleTabChange = (event, value) => {
        const {history, location, match} = this.props;

        this.setState({
            selectedTabIndex: value
        });

        const queryParams = HttpUtils.generateQueryParamString({
            tab: this.tabs[value]
        });
        history.replace(match.url + queryParams, {
            ...location.state
        });
    };

    render = () => {
        const {classes, match} = this.props;
        const {selectedTabIndex} = this.state;

        const cellName = match.params.cellName;

        const details = <Details/>;
        const microservices = <MicroserviceList/>;
        const metrics = <Metrics/>;
        const tabContent = [details, microservices, metrics];

        return (
            <React.Fragment>
                <TopToolbar title={`Cell: ${cellName}`} onUpdate={this.loadCellData}/>
                <Paper className={classes.root}>
                    <Tabs value={selectedTabIndex} indicatorColor="primary"
                        onChange={this.handleTabChange} className={classes.tabs}>
                        <Tab label="Details"/>
                        <Tab label="Microservices"/>
                        <Tab label="Metrics"/>
                    </Tabs>
                    {tabContent[selectedTabIndex]}
                </Paper>
            </React.Fragment>
        );
    };

}

Cell.propTypes = {
    classes: PropTypes.object.isRequired,
    match: PropTypes.shape({
        params: PropTypes.shape({
            cellName: PropTypes.string.isRequired
        }).isRequired
    }).isRequired,
    history: PropTypes.shape({
        replace: PropTypes.func.isRequired
    }),
    location: PropTypes.shape({
        search: PropTypes.string.isRequired
    }).isRequired
};

export default withStyles(styles)(Cell);
