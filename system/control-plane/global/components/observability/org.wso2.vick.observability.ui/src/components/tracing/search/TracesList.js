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

import AccessTime from "@material-ui/icons/AccessTime";
import Constants from "../../../utils/constants";
import Grid from "@material-ui/core/Grid/Grid";
import HttpUtils from "../../../utils/api/httpUtils";
import NotificationUtils from "../../../utils/common/notificationUtils";
import Paper from "@material-ui/core/Paper/Paper";
import QueryUtils from "../../../utils/common/queryUtils";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import TablePagination from "@material-ui/core/TablePagination/TablePagination";
import Typography from "@material-ui/core/Typography/Typography";
import moment from "moment";
import withGlobalState from "../../common/state";
import withStyles from "@material-ui/core/styles/withStyles";
import withColor, {ColorGenerator} from "../../common/color";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    trace: {
        cursor: "pointer",
        marginTop: theme.spacing.unit * 2,
        marginRight: 0,
        marginBottom: theme.spacing.unit * 2,
        marginLeft: 0
    },
    traceHeader: {
        backgroundColor: "#dfdfdf",
        padding: theme.spacing.unit
    },
    traceHeaderRight: {
        fontWeight: 400,
        fontSize: "small",
        textAlign: "right"
    },
    cellName: {
        fontWeight: 500,
        fontSize: "normal"
    },
    serviceName: {
        fontWeight: 500,
        fontSize: "normal",
        paddingRight: theme.spacing.unit
    },
    operationName: {
        color: "#616161",
        fontSize: "small"
    },
    rootStartTime: {
        fontWeight: 300,
        color: "#616161",
        marginLeft: theme.spacing.unit,
        fontSize: "normal"
    },
    duration: {
        color: "#444",
        fontStyle: "italic",
        padding: Number(theme.spacing.unit) / 2
    },
    durationIcon: {
        verticalAlign: "-webkit-baseline-middle",
        paddingLeft: theme.spacing.unit * 2,
        color: "#666"
    },
    tagCellName: {
        color: "#666",
        paddingLeft: Number(theme.spacing.unit) / 2
    },
    tagServiceName: {
        color: "#222"
    },
    traceContent: {
        padding: theme.spacing.unit
    },
    serviceTag: {
        borderStyle: "solid",
        borderWidth: "thin",
        borderColor: "#c9c9c9",
        margin: theme.spacing.unit,
        display: "inline-block"
    },
    serviceTagColor: {
        height: "100%",
        width: theme.spacing.unit,
        display: "table-cell"
    },
    serviceTagContent: {
        padding: theme.spacing.unit,
        display: "table-cell",
        fontSize: 12
    }
});

class TracesList extends React.PureComponent {

    constructor(props) {
        super(props);

        this.state = {
            rowsPerPage: 5,
            page: 0,
            hasSearchCompleted: false,
            isLoading: false,
            searchResults: {
                rootSpans: [],
                spanCounts: []
            }
        };
    }

    componentDidMount = () => {
        const {globalState, loadTracesOnMount} = this.props;

        globalState.addListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
        if (loadTracesOnMount) {
            this.loadTraces(true);
        }
    };

    componentWillUnmount = () => {
        const {globalState} = this.props;
        globalState.removeListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
    };

    handleLoadingStateChange = (loadingStateKey, oldState, newState) => {
        this.setState({
            isLoading: newState.loadingOverlayCount > 0
        });
    };

    handleChangeRowsPerPage = (event) => {
        const rowsPerPage = event.target.value;
        this.setState({
            rowsPerPage: rowsPerPage
        });
    };

    handleChangePage = (event, page) => {
        this.setState({
            page: page
        });
    };

    /**
     * Load the trace page.
     *
     * @param {MouseEvent} event Event for the click event
     * @param {string} traceId The trace ID of the selected trace
     * @param {string} cellName The name of the cell the component belongs to if a component was selected
     * @param {string} component The component name if a component was selected
     */
    loadTracePage = (event, traceId, cellName = "", component = "") => {
        event.stopPropagation();
        this.props.onTraceClick(traceId, cellName, component);
    };

    /**
     * Get the suitable color for Component.
     *
     * @param {Object} component The name of the Component
     * @returns {string} The suitable color for the component
     */
    getColorForComponent = (component) => {
        const {colorGenerator} = this.props;
        let colorKey = ColorGenerator.UNKNOWN;
        if (component.cellName) {
            colorKey = component.cellName;
        } else if (Constants.System.GLOBAL_GATEWAY_NAME_PATTERN.test(component.serviceName)) {
            colorKey = ColorGenerator.VICK;
        } else if (Constants.System.ISTIO_MIXER_NAME_PATTERN.test(component.serviceName)) {
            colorKey = ColorGenerator.ISTIO;
        } else if (component.serviceName) {
            colorKey = component.serviceName;
        }
        return colorGenerator.getColor(colorKey);
    };

    loadTraces = (isUserAction) => {
        const self = this;
        const {globalState, filter, globalFilterOverrides} = self.props;
        const {
            cell, component, operation, tags, minDuration, minDurationMultiplier, maxDuration, maxDurationMultiplier
        } = filter;

        // Build search object
        const search = {};
        const addSearchParam = (key, value) => {
            if (value && value !== Constants.Dashboard.ALL_VALUE) {
                search[key] = value;
            }
        };
        addSearchParam("cell", cell);
        addSearchParam("serviceName", component);
        addSearchParam("operationName", operation);
        addSearchParam("tags", JSON.stringify(tags && Object.keys(tags).length > 0 ? tags : {}));
        addSearchParam("minDuration", minDuration * minDurationMultiplier);
        addSearchParam("maxDuration", maxDuration * maxDurationMultiplier);
        addSearchParam("queryStartTime", globalFilterOverrides && globalFilterOverrides.queryStartTime
            ? globalFilterOverrides.queryStartTime.valueOf()
            : QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime).valueOf());
        addSearchParam("queryEndTime", globalFilterOverrides && globalFilterOverrides.queryEndTime
            ? globalFilterOverrides.queryEndTime.valueOf()
            : QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime).valueOf());

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Searching for Traces", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/traces/search${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            self.setState((prevState) => ({
                ...prevState,
                hasSearchCompleted: true,
                searchResults: {
                    rootSpans: data.rootSpans.map((dataItem) => ({
                        traceId: dataItem[0],
                        rootCellName: dataItem[1],
                        rootServiceName: dataItem[2],
                        rootOperationName: dataItem[3],
                        rootStartTime: dataItem[4],
                        rootDuration: dataItem[5]
                    })),
                    spanCounts: data.spanCounts.map((dataItem) => ({
                        traceId: dataItem[0],
                        cellName: dataItem[1],
                        serviceName: dataItem[2],
                        count: dataItem[3]
                    }))
                }
            }));
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to search for Traces",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    render = () => {
        const {classes, hideTitle} = this.props;
        const {rowsPerPage, page, hasSearchCompleted, isLoading, searchResults} = this.state;

        // Merging the span counts and root span information
        const rootSpans = searchResults.rootSpans.reduce((accumulator, dataItem) => {
            accumulator[dataItem.traceId] = {...dataItem};
            return accumulator;
        }, {});
        const processedSearchResults = searchResults.spanCounts.reduce((accumulator, dataItem) => {
            if (accumulator[dataItem.traceId]) {
                if (!accumulator[dataItem.traceId].services) {
                    accumulator[dataItem.traceId].services = [];
                }
                accumulator[dataItem.traceId].services.push({...dataItem});
            }
            return accumulator;
        }, rootSpans);

        const searchResultsArray = [];
        for (const traceId in processedSearchResults) {
            if (processedSearchResults.hasOwnProperty(traceId)) {
                searchResultsArray.push(processedSearchResults[traceId]);
            }
        }

        let view;
        if (hasSearchCompleted && !isLoading) {
            view = (
                searchResultsArray.length > 0
                    ? (
                        <React.Fragment>
                            {
                                hideTitle
                                    ? null
                                    : (
                                        <Typography variant="h6" color="inherit" className={classes.subheading}>
                                            Traces
                                        </Typography>
                                    )
                            }
                            {
                                searchResultsArray.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                                    .map((result) => (
                                        <Paper key={result.traceId} className={classes.trace}
                                            onClick={(event) => this.loadTracePage(event, result.traceId)}>
                                            <Grid container className={classes.traceHeader}>
                                                <Grid item xs={8}>
                                                    {
                                                        result.rootCellName
                                                            ? (
                                                                <span className={classes.cellName}>
                                                                    {`${result.rootCellName}:`}
                                                                </span>
                                                            )
                                                            : null
                                                    }
                                                    <span className={classes.serviceName}>
                                                        {result.rootServiceName}
                                                    </span>
                                                    <span className={classes.operationName}>
                                                        {result.rootOperationName}
                                                    </span>
                                                </Grid>
                                                <Grid item xs={4} className={classes.traceHeaderRight}>
                                                    <span>
                                                        {
                                                            result.services.reduce(
                                                                (accumulator, currentValue) => accumulator
                                                                + currentValue.count,
                                                                0
                                                            )
                                                        } Spans
                                                    </span>
                                                    <span className={classes.rootStartTime}>
                                                        {
                                                            moment(result.rootStartTime)
                                                                .format(Constants.Pattern.DATE_TIME)
                                                        }
                                                    </span>
                                                    <span className={classes.durationIcon}>
                                                        <AccessTime varient="inherit" fontSize="small"/>
                                                    </span>
                                                    <span className={classes.duration}>
                                                        {result.rootDuration / 1000} s
                                                    </span>
                                                </Grid>
                                            </Grid>
                                            <div className={classes.traceContent}>
                                                {
                                                    result.services
                                                        .sort((a, b) => {
                                                            if (a.serviceName < b.serviceName) {
                                                                return -1;
                                                            }
                                                            if (a.serviceName > b.serviceName) {
                                                                return 1;
                                                            }
                                                            return 0;
                                                        })
                                                        .map((service) => (
                                                            <div key={`${service.cellName}-${service.serviceName}`}
                                                                className={classes.serviceTag}
                                                                onClick={
                                                                    (event) => this.loadTracePage(event, result.traceId,
                                                                        service.cellName, service.serviceName)}>
                                                                <div className={classes.serviceTagColor} style={{
                                                                    backgroundColor: this.getColorForComponent(service)
                                                                }}/>
                                                                <div className={classes.serviceTagContent}>
                                                                    <span className={classes.tagCellName}>
                                                                        {service.cellName
                                                                            ? `${service.cellName}: `
                                                                            : null} </span>
                                                                    <span className={classes.tagServiceName}>
                                                                        {service.serviceName} ({service.count})
                                                                    </span>
                                                                </div>
                                                            </div>
                                                        ))
                                                }
                                            </div>
                                        </Paper>
                                    ))
                            }
                            <TablePagination component="div" count={searchResultsArray.length} rowsPerPage={rowsPerPage}
                                backIconButtonProps={{"aria-label": "Previous Page"}} page={page}
                                labelRowsPerPage={"Traces Per Page"} onChangePage={this.handleChangePage}
                                nextIconButtonProps={{"aria-label": "Next Page"}}
                                onChangeRowsPerPage={this.handleChangeRowsPerPage}/>
                        </React.Fragment>
                    )
                    : (
                        <div>No Traces Found</div>
                    )
            );
        } else {
            view = null;
        }
        return view;
    };

}

TracesList.propTypes = {
    classes: PropTypes.any.isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator),
    globalState: PropTypes.instanceOf(StateHolder),
    hideTitle: PropTypes.bool,
    loadTracesOnMount: PropTypes.bool,
    onTraceClick: PropTypes.func.isRequired,
    filter: PropTypes.shape({
        cell: PropTypes.string,
        component: PropTypes.string,
        operation: PropTypes.string,
        tags: PropTypes.object,
        minDuration: PropTypes.number,
        minDurationMultiplier: PropTypes.number,
        maxDuration: PropTypes.number,
        maxDurationMultiplier: PropTypes.number
    }).isRequired,
    globalFilterOverrides: PropTypes.shape({
        queryStartTime: PropTypes.any.isRequired,
        queryEndTime: PropTypes.any.isRequired
    })
};

export default withStyles(styles)(withColor(withGlobalState(TracesList)));
