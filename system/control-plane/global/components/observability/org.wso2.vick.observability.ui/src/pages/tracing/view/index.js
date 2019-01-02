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
import ErrorBoundary from "../../common/error/ErrorBoundary";
import HttpUtils from "../../common/utils/httpUtils";
import NotFound from "../../common/error/NotFound";
import NotificationUtils from "../../common/utils/notificationUtils";
import Paper from "@material-ui/core/Paper/Paper";
import React from "react";
import SequenceDiagram from "./sequenceDiagram/SequenceDiagram";
import Span from "../utils/span";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import Timeline from "./timeline";
import TopToolbar from "../../common/toptoolbar";
import TracingUtils from "../utils/tracingUtils";
import UnknownError from "../../common/error/UnknownError";
import withStyles from "@material-ui/core/styles/withStyles";
import withGlobalState, {StateHolder} from "../../common/state";
import * as PropTypes from "prop-types";

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
            "sequenceDiagram",
            "dependency-diagram"
        ];
        const queryParams = HttpUtils.parseQueryParams(props.location.search);
        const preSelectedTab = queryParams.tab ? this.tabs.indexOf(queryParams.tab) : null;

        this.state = {
            traceTree: null,
            spans: [],
            selectedTabIndex: (preSelectedTab ? preSelectedTab : 0),
            isLoading: false,
            errorMessage: null
        };

        this.traceViewRef = React.createRef();
    }

    componentDidMount = () => {
        const {globalState} = this.props;

        globalState.addListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
        this.loadTrace();
    };

    componentDidUpdate = () => {
        if (this.traceViewRef.current && this.traceViewRef.current.draw) {
            this.traceViewRef.current.draw();
        }
    };

    componentWillUnmount() {
        const {globalState} = this.props;

        globalState.removeListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
    }

    handleLoadingStateChange = (loadingStateKey, oldState, newState) => {
        this.setState({
            isLoading: newState.loadingOverlayCount > 0
        });
    };

    loadTrace = () => {
        const {globalState, match} = this.props;
        const traceId = match.params.traceId;
        const self = this;

        self.setState({
            traceTree: null,
            spans: []
        });
        NotificationUtils.showLoadingOverlay("Loading trace", globalState);
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

            try {
                const rootSpan = TracingUtils.buildTree(spans);
                TracingUtils.labelSpanTree(rootSpan);

                self.setState({
                    traceTree: rootSpan,
                    spans: TracingUtils.getOrderedList(rootSpan)
                });
            } catch (e) {
                NotificationUtils.showNotification(
                    "Unable to Render Invalid Trace", NotificationUtils.Levels.ERROR, globalState);
                self.setState({
                    errorMessage: e.message
                });
            }
            NotificationUtils.hideLoadingOverlay(globalState);
        }).catch(() => {
            NotificationUtils.hideLoadingOverlay(globalState);
            NotificationUtils.showNotification(
                `Failed to fetch Trace with ID ${traceId}`,
                NotificationUtils.Levels.ERROR,
                globalState
            );
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
        const {spans, selectedTabIndex, isLoading, errorMessage} = this.state;
        const selectedMicroservice = location.state ? location.state.selectedMicroservice : null;

        const traceId = match.params.traceId;

        const tabContent = [Timeline, SequenceDiagram, DependencyDiagram];
        const SelectedTabContent = tabContent[selectedTabIndex];

        let view;
        if (spans && spans.length) {
            view = (
                <Paper className={classes.container}>
                    <Tabs value={selectedTabIndex} indicatorColor="primary"
                        onChange={this.handleTabChange}>
                        <Tab label="Timeline"/>
                        <Tab label="Sequence Diagram"/>
                        <Tab label="Dependency Diagram"/>
                    </Tabs>
                    <ErrorBoundary title={"Unable to render Invalid Trace"}>
                        <SelectedTabContent spans={spans} innerRef={this.traceViewRef}
                            selectedMicroservice={selectedMicroservice}/>
                    </ErrorBoundary>
                </Paper>
            );
        } else if (errorMessage) {
            view = <UnknownError title={"Unable to Render Trace"} description={errorMessage}/>;
        } else {
            view = <NotFound title={`Trace with ID "${traceId}" Not Found`}/>;
        }

        return (
            isLoading
                ? null
                : (
                    <React.Fragment>
                        <TopToolbar title={"Distributed Tracing"}/>
                        {view}
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
