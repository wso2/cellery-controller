/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import DependencyGraph from "../../common/DependencyGraph";
import HttpUtils from "../../common/utils/httpUtils";
import NotificationUtils from "../../common/utils/notificationUtils";
import PropTypes from "prop-types";
import QueryUtils from "../../common/utils/queryUtils";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core";
import withColor, {ColorGenerator} from "../../common/color";

const styles = () => ({
    graph: {
        width: "100%",
        height: "100%"
    }
});

const graphConfig = {
    directed: true,
    automaticRearrangeAfterDropNode: false,
    collapsible: false,
    height: 500,
    highlightDegree: 1,
    highlightOpacity: 0.2,
    linkHighlightBehavior: false,
    maxZoom: 8,
    minZoom: 1,
    nodeHighlightBehavior: true,
    panAndZoom: false,
    staticGraph: false,
    width: 1000,
    d3: {
        alphaTarget: 0.05,
        gravity: -700,
        linkLength: 150,
        linkStrength: 1
    },
    node: {
        color: "#d3d3d3",
        fontColor: "black",
        fontSize: 18,
        fontWeight: "normal",
        highlightFontSize: 18,
        highlightFontWeight: "bold",
        highlightStrokeColor: "SAME",
        highlightStrokeWidth: 1.5,
        labelProperty: "name",
        mouseCursor: "pointer",
        opacity: 1,
        renderLabel: true,
        size: 600,
        strokeColor: "black",
        strokeWidth: 2
    },
    link: {
        color: "#d3d3d3",
        opacity: 1,
        semanticStrokeWidth: false,
        strokeWidth: 4,
        highlightColor: "black"
    }
};

class ServiceDependencyView extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: {
                nodes: [],
                links: []
            }
        };
    }

    componentDidMount = () => {
        const {globalState} = this.props;

        this.update(
            true,
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime).valueOf(),
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime).valueOf()
        );
    };

    update = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell, service} = this.props;
        const self = this;

        const search = {
            fromTime: queryStartTime.valueOf(),
            toTime: queryEndTime.valueOf()
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Dependency Graph", globalState);
        }
        let url = `/dependency-model/cells/${cell}/microservices/${service}`;
        url += `${HttpUtils.generateQueryParamString(search)}`;
        HttpUtils.callObservabilityAPI(
            {
                url: url,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const nodes = [];
            data.nodes.forEach((node) => {
                nodes.push(
                    {
                        id: node.id,
                        color: self.props.colorGenerator.getColor(node.id.split(":")[0])
                    });
            });

            /*
             * Data.nodes.map((node) => {
             *     nodes.push(
             *         {
             *             id: node.id,
             *             color: self.props.colorGenerator.getColor(node.id.split(":")[0])
             *         });
             * });
             */
            self.setState({
                data: {
                    nodes: nodes,
                    links: data.edges
                }
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load service dependency view",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    onClickCell = (nodeId) => {
        // TODO: redirect to another cell view.
    };

    render = () => (
        <React.Fragment>
            <DependencyGraph
                id="service-dependency-graph"
                data={this.state.data}
                config={graphConfig}
                reloadGraph={true}
                onClickNode={this.onClickCell}
            />
        </React.Fragment>
    );

}

ServiceDependencyView.propTypes = {
    classes: PropTypes.object.isRequired,
    cell: PropTypes.string.isRequired,
    service: PropTypes.string.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired
};

export default withStyles(styles, {withTheme: true})(withColor(withGlobalState(ServiceDependencyView)));
