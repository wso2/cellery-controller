/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Button from "@material-ui/core/Button/Button";
import Details from "./Details";
import Grey from "@material-ui/core/colors/grey";
import HttpUtils from "../../../utils/api/httpUtils";
import K8sObjects from "./K8sObjects";
import {Link} from "react-router-dom";
import Metrics from "./Metrics";
import Paper from "@material-ui/core/Paper";
import React from "react";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import Timeline from "@material-ui/icons/Timeline";
import TopToolbar from "../../common/toptoolbar";
import {withStyles} from "@material-ui/core/styles";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    root: {
        flexGrow: 1,
        backgroundColor: theme.palette.background.paper,
        padding: theme.spacing.unit * 3,
        paddingTop: 0,
        margin: theme.spacing.unit
    },
    tabBar: {
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        marginBottom: theme.spacing.unit * 2,
        borderBottomWidth: 1,
        borderBottomStyle: "solid",
        borderBottomColor: Grey[200]
    },
    viewTracesContent: {
        paddingLeft: theme.spacing.unit
    },
    traceButton: {
        fontSize: 12
    }
});

class Component extends React.Component {

    constructor(props) {
        super(props);

        this.tabs = [
            "details",
            "k8s-objects",
            "metrics"
        ];
        const queryParams = HttpUtils.parseQueryParams(props.location.search);
        const preSelectedTab = queryParams.tab ? this.tabs.indexOf(queryParams.tab) : null;

        this.state = {
            selectedTabIndex: (preSelectedTab ? preSelectedTab : 0)
        };

        this.tabContentRef = React.createRef();
    }

    handleTabChange = (event, value) => {
        const {history, location, match} = this.props;

        this.setState({
            selectedTabIndex: value
        });

        // Updating the Browser URL
        const queryParams = HttpUtils.generateQueryParamString({
            ...HttpUtils.parseQueryParams(location.search),
            tab: this.tabs[value]
        });
        history.replace(match.url + queryParams, {
            ...location.state
        });
    };

    handleOnUpdate = (isUserAction, startTime, endTime) => {
        if (this.tabContentRef.current && this.tabContentRef.current.update) {
            this.tabContentRef.current.update(isUserAction, startTime, endTime);
        }
    };

    onFilterUpdate = (newFilter) => {
        const {history, location, match} = this.props;

        // Updating the Browser URL
        const queryParams = HttpUtils.generateQueryParamString({
            ...HttpUtils.parseQueryParams(location.search),
            ...newFilter
        });
        history.replace(match.url + queryParams, {
            ...location.state
        });
    };

    render() {
        const {classes, location, match} = this.props;
        const {selectedTabIndex} = this.state;

        const cellName = match.params.cellName;
        const componentName = match.params.componentName;

        const tabContent = [Details, K8sObjects, Metrics];
        const SelectedTabContent = tabContent[selectedTabIndex];

        const queryParams = HttpUtils.parseQueryParams(location.search);

        const traceSearch = {
            cell: cellName,
            component: componentName
        };
        return (
            <React.Fragment>
                <TopToolbar title={`${componentName}`} subTitle="- Component" onUpdate={this.handleOnUpdate}/>
                <Paper className={classes.root}>
                    <div className={classes.tabBar}>
                        <Tabs value={selectedTabIndex} indicatorColor="primary"
                            onChange={this.handleTabChange} className={classes.tabs}>
                            <Tab label="Details"/>
                            <Tab label="K8s Objects"/>
                            <Tab label="Metrics"/>
                        </Tabs>
                        <Button className={classes.traceButton} component={Link}
                            to={`/tracing/search${HttpUtils.generateQueryParamString(traceSearch)}`}>
                            <Timeline/><span className={classes.viewTracesContent}>View Traces</span>
                        </Button>
                    </div>
                    <SelectedTabContent innerRef={this.tabContentRef} cell={cellName} component={componentName}
                        onFilterUpdate={this.onFilterUpdate} globalFilterOverrides={queryParams}/>
                </Paper>
            </React.Fragment>
        );
    }

}

Component.propTypes = {
    classes: PropTypes.object.isRequired,
    match: PropTypes.shape({
        params: PropTypes.shape({
            cellName: PropTypes.string.isRequired,
            componentName: PropTypes.string.isRequired
        }).isRequired
    }).isRequired,
    history: PropTypes.shape({
        replace: PropTypes.func.isRequired
    }),
    location: PropTypes.shape({
        search: PropTypes.string.isRequired
    }).isRequired
};

export default withStyles(styles)(Component);
