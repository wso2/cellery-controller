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

import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import DependencyGraph from "./DependencyGraph";
import Divider from "@material-ui/core/Divider";
import Drawer from "@material-ui/core/Drawer";
import Grey from "@material-ui/core/colors/grey";
import IconButton from "@material-ui/core/IconButton";
import MoreIcon from "@material-ui/icons/MoreHoriz";
import NotificationUtils from "../common/utils/notificationUtils";
import Paper from "@material-ui/core/Paper";
import PropTypes from "prop-types";
import React from "react";
import SidePanelContent from "./SidePanelContent";
import StateHolder from "../common/state/stateHolder";
import TopToolbar from "../common/toptoolbar";
import Typography from "@material-ui/core/Typography";
import axios from "axios";
import classNames from "classnames";
import withGlobalState from "../common/state";
import {withStyles} from "@material-ui/core/styles";
import withColor, {ColorGenerator} from "../common/color";

const graphConfig = {
    directed: true,
    automaticRearrangeAfterDropNode: false,
    collapsible: false,
    height: 800,
    highlightDegree: 1,
    highlightOpacity: 0.2,
    linkHighlightBehavior: false,
    maxZoom: 8,
    minZoom: 0.1,
    nodeHighlightBehavior: true,
    panAndZoom: false,
    staticGraph: false,
    width: 1400,
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
        top: 140,
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
        const color = this.props.colorGenerator.getColor(nodeProps.id);
        return <svg x="0px" y="0px"
                    width="50px" height="50px" viewBox="0 0 240 240">
            <polygon fill={color} points="224,179.5 119.5,239.5 15,179.5 15,59.5 119.5,-0.5 224,59.5 "/>
        </svg>;
    };

    onClickCell = (nodeId) => {
        let cell = null;
        this.state.data.nodes.forEach((element) => {
            if (element.id === nodeId) {
                cell = element;
            }
        });
        let serviceInfo = this.loadServicesInfo(cell.services);
        this.setState((prevState) => ({
            summary: {
                ...prevState.summary,
                topic: nodeId,
                content: [
                    {
                        key: "Total",
                        value: serviceInfo.length,
                    },
                    {
                        key: "Successful",
                        value: serviceInfo.length
                    },
                    {
                        key: "Failed",
                        value: 1
                    },
                    {
                        key: "Warning",
                        value: 1
                    }
                ]
            },
            listData: serviceInfo,
            reloadGraph: false,
            isOverallSummary: false
        }))
        ;
    };

    onClickGraph = () => {
        this.setState({
            summary: JSON.parse(JSON.stringify(this.defaultState)).summary,
            listData: this.loadCellInfo(this.state.data.nodes),
            reloadGraph: true,
            isOverallSummary: true
        });
    };

    handleDrawerOpen = () => {
        this.setState({open: true});
    };

    handleDrawerClose = () => {
        this.setState({open: false});
    };

    loadListInfo = (isUserAction) => {
        const {globalState} = this.props;
        const self = this;

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Info", globalState);
        }
        setTimeout(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }, 1000);
    };

    loadCellInfo = (nodes) => {
        console.log(nodes);
        const nodeInfo = [];
        nodes.forEach((node) => {
            nodeInfo.push([1, node.id, node.id]);
        });
        return nodeInfo;
    };

    loadServicesInfo = (services) => {
        const serviceInfo = [];
        services.forEach((service) => {
            serviceInfo.push([1, service, service]);
        });
        return serviceInfo;
    };

    constructor(props) {
        super(props);
        const colorGenerator = this.props.colorGenerator;
        graphConfig.node.viewGenerator = this.viewGenerator;
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
                    }
                ]
            },
            isOverallSummary: true,
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
        this.state = JSON.parse(JSON.stringify(this.defaultState));
        // TODO: Update the url to the WSO2-sp worker node.
        axios({
            method: "GET",
            url: "http://localhost:9123/dependency-model/cell-overview"
        }).then((response) => {
            const result = response.data;
            // TODO: Update values with real data
            const summaryContent = [
                {
                    key: "Total",
                    value: result.nodes.length
                },
                {
                    key: "Successful",
                    value: 2
                },
                {
                    key: "Failed",
                    value: 1
                },
                {
                    key: "Warning",
                    value: 0
                }
            ];
            const statusCodeContent = [
                {
                    key: "Total",
                    value: 0
                },
                {
                    key: "OK",
                    value: 70
                },
                {
                    key: "3xx",
                    value: 20
                },
                {
                    key: "4xx",
                    value: 5
                },
                {
                    key: "5xx",
                    value: 5
                }
            ];
            this.defaultState.summary.content = summaryContent;
            colorGenerator.addKeys(result.nodes);
            const cellList = this.loadCellInfo(result.nodes);
            this.setState((prevState) => ({
                data: {
                    nodes: result.nodes,
                    links: result.edges
                },
                summary: {
                    ...prevState.summary,
                    content: summaryContent
                },
                request: {
                    ...prevState.request,
                    statusCodes: statusCodeContent
                },
                listData: cellList

            }));
        }).catch((error) => {
            this.setState({error: error});
        });
    }


    render() {
        const {classes, theme} = this.props;
        const {open} = this.state;

        return (
            <React.Fragment>
                <TopToolbar title={"Overview"} onUpdate={this.loadListInfo}/>

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
                                {(this.state.isOverallSummary) ? "Overview" : "Cell Details"}
                            </Typography>
                        </div>
                        <Divider/>
                        <SidePanelContent
                            summary={this.state.summary}
                            request={this.state.request}
                            isOverview={this.state.isOverallSummary}
                            open={this.state.open}
                            listData={this.state.listData}>
                        </SidePanelContent>
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

