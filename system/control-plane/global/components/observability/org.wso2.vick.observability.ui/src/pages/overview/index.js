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

import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import Constants from "../common/constants";
import DependencyGraph from "./DependencyGraph";
import Divider from "@material-ui/core/Divider";
import Drawer from "@material-ui/core/Drawer";
import Grey from "@material-ui/core/colors/grey";
import HttpUtils from "../common/utils/httpUtils";
import IconButton from "@material-ui/core/IconButton";
import MoreIcon from "@material-ui/icons/MoreHoriz";
import NotificationUtils from "../common/utils/notificationUtils";
import Paper from "@material-ui/core/Paper";
import QueryUtils from "../common/utils/queryUtils";
import React from "react";
import SidePanelContent from "./SidePanelContent";
import StateHolder from "../common/state/stateHolder";
import TopToolbar from "../common/toptoolbar";
import Typography from "@material-ui/core/Typography";
import classNames from "classnames";
import withGlobalState from "../common/state";
import {withStyles} from "@material-ui/core/styles";
import withColor, {ColorGenerator} from "../common/color";
import * as PropTypes from "prop-types";

const graphConfig = {
    directed: true,
    automaticRearrangeAfterDropNode: false,
    collapsible: false,
    highlightDegree: 1,
    highlightOpacity: 0.2,
    linkHighlightBehavior: false,
    maxZoom: 8,
    minZoom: 0.1,
    nodeHighlightBehavior: true,
    panAndZoom: false,
    staticGraph: false,
    height: 700,
    width: 1050,
    d3: {
        alphaTarget: 0.05,
        gravity: -1500,
        linkLength: 150,
        linkStrength: 1
    },
    node: {
        color: "#d3d3d3",
        fontColor: "black",
        fontSize: 18,
        fontWeight: "normal",
        highlightColor: "red",
        highlightFontSize: 18,
        highlightFontWeight: "bold",
        highlightStrokeColor: "SAME",
        highlightStrokeWidth: 1.5,
        labelProperty: "name",
        mouseCursor: "pointer",
        opacity: 1,
        renderLabel: true,
        size: 600,
        strokeColor: "green",
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

const drawerWidth = 300;

const styles = (theme) => ({
    root: {
        display: "flex"
    },
    moreDetails: {
        position: "absolute",
        right: 25,
        transition: theme.transitions.create(["margin", "width"], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen
        })
    },
    moreDetailsShift: {
        transition: theme.transitions.create(["margin", "width"], {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen
        })
    },
    menuButton: {
        marginTop: 8,
        marginRight: 8
    },
    hide: {
        display: "none"
    },
    drawer: {
        width: drawerWidth,
        flexShrink: 0
    },
    drawerPaper: {
        top: 135,
        width: drawerWidth,
        borderTopWidth: 1,
        borderTopStyle: "solid",
        borderTopColor: Grey[200]
    },
    drawerHeader: {
        display: "flex",
        alignItems: "center",
        padding: 5,
        justifyContent: "flex-start",
        textTransform: "uppercase",
        minHeight: "fit-content"
    },
    content: {
        flexGrow: 1,
        padding: theme.spacing.unit * 3,
        transition: theme.transitions.create("margin", {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen
        }),
        marginLeft: Number(theme.spacing.unit),
        marginRight: -drawerWidth
    },
    contentShift: {
        transition: theme.transitions.create("margin", {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen
        }),
        marginRight: 0
    },
    sideBarHeading: {
        letterSpacing: 1,
        fontSize: 12,
        marginLeft: 4
    }
});

class Overview extends React.Component {

    viewGenerator = (nodeProps) => {
        const {colorGenerator} = this.props;
        const {selectedCell} = this.state;

        const nodeId = nodeProps.id;
        const color = colorGenerator.getColor(nodeId);
        const state = this.getCellState(nodeId);

        const style = {};
        if (selectedCell && selectedCell.id === nodeId) {
            style.stroke = "#383058";
            style.strokeWidth = 15;
        }
        style.transform = "translate(15%, 10%) scale(0.8, 0.8)";
        const cellIcon = <polygon fill={color} points="208,179.5 103.5,239.5 -1,179.5 -1,59.5 103.5,-0.5 208,59.5"
            style={style}/>;

        let cellView;
        if (state === Constants.Status.Success) {
            cellView = (
                <svg x="0px" y="0px" width="100%" height="100%" viewBox="0 0 240 240">
                    {cellIcon}
                </svg>
            );
        } else if (state === Constants.Status.Warning) {
            cellView = (
                <svg x="0px" y="0px" width="100%" height="100%" viewBox="0 0 240 240">
                    <g>
                        <g>
                            {cellIcon}
                        </g>
                    </g>
                    <g>
                        <path d="M146.5,6.1c-23.6,0-42.9,19.3-42.9,42.9s19.3,42.9,42.9,42.9s42.9-19.3,42.9-42.9S170.1,
                              6.1,146.5,6.1z" stroke="#fff" strokeWidth="3" fill="#ff9800"/>
                        <path fill="#ffffff" d="M144.6,56.9h7.9v7.9h-7.9V56.9z M144.6,25.2h7.9V49h-7.9V25.2z"/>
                    </g>
                </svg>
            );
        } else {
            cellView = (
                <svg x="0px" y="0px" width="100%" height="100%" viewBox="0 0 240 240">
                    <g>
                        <g>
                            {cellIcon}
                        </g>
                    </g>
                    <g>
                        <path d="M146.5,6.1c-23.6,0-42.9,19.3-42.9,42.9s19.3,42.9,42.9,42.9s42.9-19.3,42.9-42.9S170.1,
                              6.1, 146.5,6.1z" stroke="#fff" strokeWidth="3" fill="#F44336"/>
                        <path fill="#ffffff" d="M144.6,56.9h7.9v7.9h-7.9V56.9z M144.6,25.2h7.9V49h-7.9V25.2z"/>
                    </g>
                </svg>
            );
        }
        return cellView;
    };


    getCellState = (nodeId) => {
        const healthInfo = this.defaultState.healthInfo.find((element) => element.nodeId === nodeId);
        return healthInfo.status;
    };

    onClickCell = (nodeId) => {
        const fromTime = this.state.currentTimeRange.fromTime;
        const toTime = this.state.currentTimeRange.toTime;

        const search = {
            queryStartTime: fromTime.valueOf(),
            queryEndTime: toTime.valueOf(),
            timeGranularity: QueryUtils.getTimeGranularity(fromTime, toTime)
        };
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/${nodeId}/microservices${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            this.props.globalState
        ).then((response) => {
            const cell = this.state.data.nodes.find((element) => element.id === nodeId);
            const serviceHealth = this.getServiceHealth(cell.services, response);
            const serviceHealthCount = this.getHealthCount(serviceHealth);
            const statusCodeContent = this.getStatusCodeContent(nodeId, this.defaultState.request.cellStats);
            const serviceInfo = this.loadServicesInfo(cell.services, serviceHealth);
            this.setState((prevState) => ({
                summary: {
                    ...prevState.summary,
                    topic: nodeId,
                    content: [
                        {
                            key: "Total",
                            value: serviceInfo.length
                        },
                        {
                            key: "Successful",
                            value: serviceHealthCount.success
                        },
                        {
                            key: "Failed",
                            value: serviceHealthCount.error
                        },
                        {
                            key: "Warning",
                            value: serviceHealthCount.warning
                        }
                    ]
                },
                data: {...prevState.data},
                listData: serviceInfo,
                reloadGraph: true,
                selectedCell: cell,
                request: {
                    ...prevState.request,
                    statusCodes: statusCodeContent
                }
            }));
        }).catch((error) => {
            this.setState({error: error});
        });
    };

    getServiceHealth = (services, responseCodeStats) => {
        const {globalState} = this.props;
        const config = globalState.get(StateHolder.CONFIG);
        const healthInfo = [];
        services.forEach((service) => {
            const total = this.getTotalServiceRequests(service, responseCodeStats, "*");
            if (total === 0) {
                healthInfo.push({nodeId: service, status: Constants.Status.Success, percentage: 1});
            } else {
                const error = this.getTotalServiceRequests(service, responseCodeStats, "5xx");
                const successPercentage = 1 - (error / total);

                if (successPercentage > config.percentageRangeMinValue.warningThreshold) {
                    healthInfo.push({nodeId: service, status: Constants.Status.Success, percentage: successPercentage});
                } else if (successPercentage > config.percentageRangeMinValue.errorThreshold) {
                    healthInfo.push({nodeId: service, status: Constants.Status.Warning, percentage: successPercentage});
                } else {
                    healthInfo.push({nodeId: service, status: Constants.Status.Error, percentage: successPercentage});
                }
            }
        });
        return healthInfo;
    };

    onClickGraph = () => {
        const defaultState = JSON.parse(JSON.stringify(this.defaultState));
        this.setState((prevState) => ({
            ...prevState,
            summary: defaultState.summary,
            listData: this.loadCellInfo(defaultState.data.nodes),
            reloadGraph: true,
            selectedCell: null,
            request: defaultState.request
        }));
    };

    handleDrawerOpen = () => {
        this.setState({open: true});
    };

    handleDrawerClose = () => {
        this.setState({open: false});
    };

    loadCellInfo = (nodes) => {
        const nodeInfo = [];
        nodes.forEach((node) => {
            const healthInfo = this.defaultState.healthInfo.find((element) => element.nodeId === node.id);
            nodeInfo.push([healthInfo.percentage, node.id, node.id]);
        });
        return nodeInfo;
    };

    loadServicesInfo = (services, healthInfo) => {
        const serviceInfo = [];
        services.forEach((service) => {
            const healthElement = healthInfo.find((element) => element.nodeId === service);
            serviceInfo.push([healthElement.percentage, service, service]);
        });
        return serviceInfo;
    };

    constructor(props) {
        super(props);
        this.initializeDefault();
        this.state = JSON.parse(JSON.stringify(this.defaultState));
        graphConfig.node.viewGenerator = this.viewGenerator;
    }

    initializeDefault = () => {
        this.defaultState = {
            summary: {
                topic: "VICK Deployment",
                content: [
                    {
                        key: "Total",
                        value: 0
                    },
                    {
                        key: "Successful",
                        value: 0
                    },
                    {
                        key: "Failed",
                        value: 0
                    },
                    {
                        key: "Warning",
                        value: 0
                    }
                ]
            },
            request: {
                statusCodes: [
                    {
                        key: "Total",
                        value: 0
                    },
                    {
                        key: "OK",
                        value: 0
                    },
                    {
                        key: "3xx",
                        value: 0
                    },
                    {
                        key: "4xx",
                        value: 0
                    },
                    {
                        key: "5xx",
                        value: 0
                    },
                    {
                        key: "Unknown",
                        value: 0
                    }
                ],
                cellStats: []
            },
            healthInfo: [],
            selectedCell: null,
            data: {
                nodes: null,
                links: null
            },
            error: null,
            reloadGraph: true,
            open: true,
            listData: [],
            page: 0,
            rowsPerPage: 5
        };
    };

    callOverviewInfo = (fromTime, toTime) => {
        const colorGenerator = this.props.colorGenerator;

        const search = {};
        if (fromTime && toTime) {
            search.fromTime = fromTime.valueOf();
            search.toTime = toTime.valueOf();
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/dependency-model/cells${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            this.props.globalState
        ).then((result) => {
            const {nodes, edges} = result;

            this.defaultState.healthInfo = this.getCellHealth(nodes);
            const healthCount = this.getHealthCount(this.defaultState.healthInfo);
            const summaryContent = [
                {
                    key: "Total",
                    value: nodes.length
                },
                {
                    key: "Successful",
                    value: healthCount.success
                },
                {
                    key: "Failed",
                    value: healthCount.error
                },
                {
                    key: "Warning",
                    value: healthCount.warning
                }
            ];
            this.defaultState.summary.content = summaryContent;
            this.defaultState.data.nodes = nodes;
            this.defaultState.data.links = edges;
            colorGenerator.addKeys(nodes);
            const cellList = this.loadCellInfo(nodes);
            this.setState((prevState) => ({
                data: {
                    nodes: nodes,
                    links: edges
                },
                summary: {
                    ...prevState.summary,
                    content: summaryContent
                },
                listData: cellList
            }));
        }).catch((error) => {
            this.setState({error: error});
        });
    };

    getHealthCount = (healthInfo) => {
        let successCount = 0;
        let warningCount = 0;
        let errorCount = 0;
        healthInfo.forEach((info) => {
            if (info.status === Constants.Status.Success) {
                successCount += 1;
            } else if (info.status === Constants.Status.Warning) {
                warningCount += 1;
            } else {
                errorCount += 1;
            }
        });
        return {success: successCount, warning: warningCount, error: errorCount};
    };

    getCellHealth = (nodes) => {
        const {globalState} = this.props;
        const config = globalState.get(StateHolder.CONFIG);
        const healthInfo = [];
        nodes.forEach((node) => {
            const total = this.getTotalRequests(node.id, this.defaultState.request.cellStats, "*");
            if (total === 0) {
                healthInfo.push({nodeId: node.id, status: Constants.Status.Success, percentage: 1});
            } else {
                const error = this.getTotalRequests(node.id, this.defaultState.request.cellStats, "5xx");
                const successPercentage = 1 - (error / total);
                if (successPercentage > config.percentageRangeMinValue.warningThreshold) {
                    healthInfo.push({nodeId: node.id, status: Constants.Status.Success, percentage: successPercentage});
                } else if (successPercentage > config.percentageRangeMinValue.errorThreshold) {
                    healthInfo.push({nodeId: node.id, status: Constants.Status.Warning, percentage: successPercentage});
                } else {
                    healthInfo.push({nodeId: node.id, status: Constants.Status.Error, percentage: successPercentage});
                }
            }
        });
        return healthInfo;
    };

    callRequestStats = (fromTime, toTime) => {
        const search = {
            queryStartTime: fromTime.valueOf(),
            queryEndTime: toTime.valueOf(),
            timeGranularity: QueryUtils.getTimeGranularity(fromTime, toTime)
        };
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            this.props.globalState
        ).then((response) => {
            const statusCodeContent = this.getStatusCodeContent(null, response);
            this.defaultState.request.statusCodes = statusCodeContent;
            this.defaultState.request.cellStats = response;
            this.setState((prevState) => ({
                stats: {
                    cellStats: response
                },
                request: {
                    ...prevState.request,
                    statusCodes: statusCodeContent
                }
            }));
            this.callOverviewInfo(fromTime, toTime);
        }).catch((error) => {
            this.setState({error: error});
        });
    };

    getStatusCodeContent = (cell, data) => {
        const total = this.getTotalRequests(cell, data, "*");
        const response2xx = this.getTotalRequests(cell, data, "2xx");
        const response3xx = this.getTotalRequests(cell, data, "3xx");
        const response4xx = this.getTotalRequests(cell, data, "4xx");
        const response5xx = this.getTotalRequests(cell, data, "5xx");
        const responseUnknown = total - (response2xx + response3xx + response4xx + response5xx);
        return [
            {
                key: "Total",
                value: total
            },
            {
                key: "OK",
                count: response2xx,
                value: this.getPercentage(response2xx, total)
            },
            {
                key: "3xx",
                count: response3xx,
                value: this.getPercentage(response3xx, total)
            },
            {
                key: "4xx",
                count: response4xx,
                value: this.getPercentage(response4xx, total)
            },
            {
                key: "5xx",
                count: response5xx,
                value: this.getPercentage(response5xx, total)
            },
            {
                key: "Unknown",
                count: responseUnknown,
                value: this.getPercentage(responseUnknown, total)
            }
        ];
    };

    getPercentage = (responseCode, total) => {
        if (total !== 0) {
            return Math.round((responseCode / total) * 100);
        }
        return 0;
    };

    getTotalRequests = (cell, stats, responseCode) => {
        let total = 0;
        stats.forEach((stat) => {
            if (stat[1] !== "") {
                if (!cell || cell === stat[1]) {
                    if (responseCode === "*") {
                        total += stat[4];
                    } else if (responseCode === stat[2]) {
                        total += stat[4];
                    }
                }
            }
        });
        return total;
    };

    getTotalServiceRequests = (cell, stats, responseCode) => {
        let total = 0;
        stats.forEach((stat) => {
            if (stat[2] !== "") {
                if (!cell || cell === stat[2]) {
                    if (responseCode === "*") {
                        total += stat[6];
                    } else if (responseCode === stat[4]) {
                        total += stat[6];
                    }
                }
            }
        });
        return total;
    };

    loadOverviewOnTimeUpdate = (isUserAction, startTime, endTime) => {
        const {globalState} = this.props;
        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Overview", globalState);
            this.setState((prevState) => ({
                ...prevState,
                currentTimeRange: {
                    fromTime: startTime,
                    toTime: endTime
                }
            }));
            this.callRequestStats(startTime, endTime);
        }
        setTimeout(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }, 1000);
    };

    render() {
        const {classes, theme} = this.props;
        const {open, selectedCell} = this.state;

        return (
            <React.Fragment>
                <TopToolbar title={"Overview"} onUpdate={this.loadOverviewOnTimeUpdate}/>

                <div className={classes.root}>
                    <Paper
                        className={classNames(classes.content, {
                            [classes.contentShift]: open
                        })}
                    >
                        <DependencyGraph
                            id="graph-id"
                            data={this.state.data}
                            config={graphConfig}
                            reloadGraph={this.state.reloadGraph}
                            onClickNode={this.onClickCell}
                            onClickGraph={this.onClickGraph}
                        />
                    </Paper>
                    <div className={classNames(classes.moreDetails, {
                        [classes.moreDetailsShift]: open
                    })}>
                        <IconButton
                            color="inherit"
                            aria-label="Open drawer"
                            onClick={this.handleDrawerOpen}
                            className={classNames(classes.menuButton, open && classes.hide)}
                        >
                            <MoreIcon/>
                        </IconButton>
                    </div>

                    <Drawer
                        className={classes.drawer}
                        variant="persistent"
                        anchor="right"
                        open={open}
                        classes={{
                            paper: classes.drawerPaper
                        }}
                    >
                        <div className={classes.drawerHeader}>
                            <IconButton onClick={this.handleDrawerClose}>
                                {theme.direction === "rtl" ? <ChevronLeftIcon/> : <ChevronRightIcon/>}
                            </IconButton>
                            <Typography color="textSecondary" className={classes.sideBarHeading}>
                                {selectedCell ? "Cell Details" : "Overview"}
                            </Typography>
                        </div>
                        <Divider/>
                        <SidePanelContent
                            summary={this.state.summary}
                            request={this.state.request}
                            selectedCell={selectedCell}
                            open={this.state.open}
                            listData={this.state.listData}/>
                    </Drawer>
                </div>
            </React.Fragment>
        );
    }

}

Overview.propTypes = {
    classes: PropTypes.object.isRequired,
    theme: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator)
};

export default withStyles(styles, {withTheme: true})(withGlobalState(withColor(Overview)));

