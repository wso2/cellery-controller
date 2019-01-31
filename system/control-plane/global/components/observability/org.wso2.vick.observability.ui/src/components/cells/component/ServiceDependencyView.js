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
import ErrorBoundary from "../../common/error/ErrorBoundary";
import HttpUtils from "../../../utils/api/httpUtils";
import InfoOutlined from "@material-ui/icons/InfoOutlined";
import NotificationUtils from "../../../utils/common/notificationUtils";
import QueryUtils from "../../../utils/common/queryUtils";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import Typography from "@material-ui/core/Typography/Typography";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core";
import withColor, {ColorGenerator} from "../../common/color";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    subtitle: {
        fontWeight: 400,
        fontSize: "1rem"
    },
    graph: {
        width: "100%",
        height: "100%"
    },
    diagram: {
        padding: theme.spacing.unit * 3
    },
    dependencies: {
        marginTop: theme.spacing.unit * 3
    },
    info: {
        display: "inline-flex"
    },
    infoIcon: {
        verticalAlign: "middle",
        display: "inline-flex",
        fontSize: 18,
        marginRight: 4
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
        fontColor: "#555",
        fontSize: 16,
        fontWeight: "normal",
        highlightFontSize: 16,
        highlightColor: "SAME",
        highlightFontWeight: "bold",
        highlightStrokeColor: "SAME",
        highlightStrokeWidth: 1,
        labelProperty: "name",
        mouseCursor: "pointer",
        opacity: 1,
        renderLabel: true,
        size: 600,
        strokeColor: "#777",
        strokeWidth: 0
    },
    link: {
        color: "#d3d3d3",
        opacity: 1,
        semanticStrokeWidth: false,
        strokeWidth: 4,
        highlightColor: "#777"
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
        const {globalState, cell, component} = this.props;
        const self = this;

        const search = {
            fromTime: queryStartTime.valueOf(),
            toTime: queryEndTime.valueOf()
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Component Dependency Graph", globalState);
        }
        let url = `/dependency-model/cells/${cell}/microservices/${component}`;
        url += `${HttpUtils.generateQueryParamString(search)}`;
        HttpUtils.callObservabilityAPI(
            {
                url: url,
                method: "GET"
            },
            globalState
        ).then((data) => {
            // Update node,edge data to show external cell dependencies
            const nodes = [];
            const edges = [];

            data.nodes.forEach((node) => {
                if (cell === node.id.split(":")[0]) {
                    nodes.push(node);
                } else if (node.id.split(":")[1] === "gateway") {
                    nodes.push({
                        ...node,
                        id: node.id.split(":")[0]
                    });
                }
            });

            data.edges.forEach((edge) => {
                if (cell === edge.source.split(":")[0] && cell === edge.target.split(":")[0]) {
                    edges.push(edge);
                } else if (cell === edge.source.split(":")[0] && edge.target.split(":")[1] === "gateway") {
                    edges.push({
                        ...edge,
                        target: edge.target.split(":")[0]
                    });
                }
            });

            self.setState({
                data: {
                    nodes: nodes,
                    links: edges
                }
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load component dependency view",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    viewGenerator = (nodeProps) => {
        const color = this.props.colorGenerator.getColor(nodeProps.id.split(":")[0]);
        const outlineColor = ColorGenerator.shadeColor(color, -0.08);
        const {cell, component} = this.props;
        const style = {};
        style.transform = "translate(1%, 10%) scale(2.2, 2.2)";
        if (cell === nodeProps.id.split(":")[0]) {
            return (
                <svg x="0px" y="0px" width="100%" height="100%" viewBox="0 0 240 240">
                    <circle cx="100" cy="100" r={80} fill={color}
                        strokeWidth="10" stroke={(component === nodeProps.id.split(":")[1]) ? "#444" : outlineColor}/>
                </svg>
            );
        }
        return (
            <svg x="0px" y="0px" width="100%" height="100%" viewBox="0 0 240 240">
                <polygon strokeWidth="4" fill={color} stroke={outlineColor} style={style} strokeLinejoin="round"
                    points="34.2,87.4 12.3,65.5 12.3,34.5 34.2,12.6 65.2,12.6 87.1,34.5 87.1,65.5 65.2,87.4"
                />
            </svg>
        );
    };

    onClickCell = (nodeId) => {
        // TODO: redirect to another cell view.
    };

    render = () => {
        const {classes, cell, component} = this.props;
        const dependedNodeCount = this.state.data.nodes.length;

        let view;
        if (dependedNodeCount > 1) {
            view = (
                <ErrorBoundary title={"Unable to Render"} description={"Unable to Render due to Invalid Data"}>
                    <DependencyGraph
                        id="component-dependency-graph"
                        data={this.state.data}
                        config={{
                            ...graphConfig,
                            node: {
                                ...graphConfig.node,
                                viewGenerator: this.viewGenerator
                            }
                        }}
                        onClickNode={this.onClickCell}
                    />
                </ErrorBoundary>
            );
        } else {
            view = (
                <div>
                    <InfoOutlined className={classes.infoIcon} color="action"/>
                    <Typography variant="subtitle2" color="textSecondary" className={classes.info}>
                        {`"${component}"`} component in {`"${cell}"`} cell does not depend on any other Component
                    </Typography>
                </div>
            );
        }
        return (
            <div className={classes.dependencies}>
                <Typography color="textSecondary" className={classes.subtitle}>
                    Dependencies
                </Typography>
                <div className={classes.diagram}>
                    {view}
                </div>
            </div>
        );
    }

}

ServiceDependencyView.propTypes = {
    classes: PropTypes.object.isRequired,
    cell: PropTypes.string.isRequired,
    component: PropTypes.string.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired
};

export default withStyles(styles, {withTheme: true})(withColor(withGlobalState(ServiceDependencyView)));
