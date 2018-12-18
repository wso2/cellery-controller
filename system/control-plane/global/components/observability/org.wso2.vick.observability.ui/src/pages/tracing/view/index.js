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

import DependencyDiagram from "./DependencyDiagram";
import HttpUtils from "../../common/utils/httpUtils";
import NotFound from "../../common/NotFound";
import NotificationUtils from "../../common/utils/notificationUtils";
import Paper from "@material-ui/core/Paper/Paper";
import PropTypes from "prop-types";
import React from "react";
import SequenceDiagram from "./sequence-diagram/SequenceDiagram";
import Span from "../utils/span";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import Timeline from "./timeline";
import TopToolbar from "../../common/toptoolbar";
import TracingUtils from "../utils/tracingUtils";
import withStyles from "@material-ui/core/styles/withStyles";
import withGlobalState, {StateHolder} from "../../common/state";

const styles = (theme) => ({
    container: {
        padding: theme.spacing.unit * 3
    }
});

class View extends React.Component {

    constructor(props) {
        super(props);

        this.tabs = [
            "timeline",
            "sequence-diagram",
            "dependency-diagram"
        ];
        const queryParams = HttpUtils.parseQueryParams(props.location.search);
        const preSelectedTab = queryParams.tab ? this.tabs.indexOf(queryParams.tab) : null;

        this.state = {
            spans: [],
            selectedTabIndex: (preSelectedTab ? preSelectedTab : 0),
            isLoading: true
        };

        this.traceViewRef = React.createRef();
    }

    componentDidMount = () => {
        this.loadTrace(true);
    };

    componentDidUpdate = () => {
        if (this.traceViewRef.current && this.traceViewRef.current.draw) {
            this.traceViewRef.current.draw();
        }
    };

    loadTrace = (isUserAction) => {
        const {globalState, match} = this.props;
        const traceId = match.params.traceId;
        const self = this;

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading trace", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/traces/${traceId}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const spans = data.map((dataItem) => new Span({
                traceId: dataItem[0],
                spanId: dataItem[1],
                parentId: dataItem[2],
                namespace: dataItem[3],
                cell: dataItem[4],
                serviceName: dataItem[5],
                pod: dataItem[6],
                operationName: dataItem[7],
                kind: dataItem[8],
                startTime: dataItem[9],
                duration: dataItem[10],
                tags: dataItem[11]
            }));

            const rootSpan = TracingUtils.buildTree(spans);
            TracingUtils.labelSpanTree(rootSpan);

            self.setState({
                traceTree: rootSpan,
                spans: TracingUtils.getOrderedList(rootSpan),
                isLoading: false
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            self.setState({
                isLoading: false
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    `Failed to fetch Trace with ID ${traceId}`,
                    StateHolder.NotificationLevels.ERROR,
                    globalState
                );
            }
        });
    };

    handleTabChange = (event, value) => {
        const {history, location, match} = this.props;

        this.setState({
            selectedTabIndex: value
        });

        // Updating the Browser URL
        const queryParamsString = HttpUtils.generateQueryParamString({
            tab: this.tabs[value]
        });
        history.replace(match.url + queryParamsString, {
            ...location.state
        });
    };

    render = () => {
        const {classes, location, match} = this.props;
        const {isLoading, spans, selectedTabIndex} = this.state;
        const selectedMicroservice = location.state.selectedMicroservice;

        const traceId = match.params.traceId;

        const tabContent = [Timeline, SequenceDiagram, DependencyDiagram];
        const SelectedTabContent = tabContent[selectedTabIndex];

        return (
            isLoading
                ? null
                : (
                    <React.Fragment>
                        <TopToolbar title={"Distributed Tracing"} onUpdate={this.loadTrace}/>
                        {
                            spans && spans.length === 0
                                ? (
                                    <NotFound content={`Trace with ID "${traceId}" Not Found`}/>
                                )
                                : (
                                    <Paper className={classes.container}>
                                        <Tabs value={selectedTabIndex} indicatorColor="primary"
                                            onChange={this.handleTabChange}>
                                            <Tab label="Timeline"/>
                                            <Tab label="Sequence Diagram"/>
                                            <Tab label="Dependency Diagram"/>
                                        </Tabs>
                                        <SelectedTabContent spans={spans} innerRef={this.traceViewRef}
                                            selectedMicroservice={selectedMicroservice}/>
                                    </Paper>
                                )
                        }
                    </React.Fragment>
                )
        );
    };

}

View.propTypes = {
    classes: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    match: PropTypes.shape({
        params: PropTypes.shape({
            traceId: PropTypes.string.isRequired
        }).isRequired
    }).isRequired,
    history: PropTypes.shape({
        replace: PropTypes.func.isRequired
    }),
    location: PropTypes.shape({
        search: PropTypes.string.isRequired,
        state: PropTypes.shape({
            selectedMicroservice: PropTypes.shape({
                cellName: PropTypes.string.isRequired,
                serviceName: PropTypes.string.isRequired
            }).isRequired
        }).isRequired
    }).isRequired
};

export default withStyles(styles)(withGlobalState(View));
