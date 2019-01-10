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
import Paper from "@material-ui/core/Paper/Paper";
import React from "react";
import TablePagination from "@material-ui/core/TablePagination/TablePagination";
import Typography from "@material-ui/core/Typography/Typography";
import moment from "moment";
import {withRouter} from "react-router-dom";
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
            page: 0
        };
    }

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
     * @param {string} cellName The name of the cell the microservice belongs to if a microservice was selected
     * @param {string} microservice The microservice name if a microservice was selected
     */
    loadTracePage = (event, traceId, cellName = "", microservice = "") => {
        event.stopPropagation();
        this.props.history.push({
            pathname: `./id/${traceId}`,
            state: {
                selectedMicroservice: {
                    cellName: cellName,
                    serviceName: microservice
                }
            }
        });
    };

    render = () => {
        const {classes, searchResults, colorGenerator} = this.props;
        const {rowsPerPage, page} = this.state;

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

        return (
            searchResultsArray.length > 0
                ? (
                    <React.Fragment>
                        <Typography variant="h6" color="inherit" className={classes.subheading}>
                            Traces
                        </Typography>
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
                                                    {moment(result.rootStartTime).format(Constants.Pattern.DATE_TIME)}
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
                                                                backgroundColor: colorGenerator
                                                                    .getColor(service.cellName)
                                                            }}/>
                                                            <div className={classes.serviceTagContent}>
                                                                <span className={classes.tagCellName}>
                                                                    {service.cellName
                                                                        ? `${service.cellName}: `
                                                                        : null} </span>
                                                                <span className={classes.tagServiceName}>
                                                                    {service.serviceName} ({service.count})</span>
                                                            </div>
                                                        </div>
                                                    ))
                                            }
                                        </div>
                                    </Paper>
                                ))
                        }
                        <TablePagination component="div" count={searchResultsArray.length} rowsPerPage={rowsPerPage}
                            backIconButtonProps={{"aria-label": "Previous Page"}} labelRowsPerPage={"Traces Per Page"}
                            nextIconButtonProps={{"aria-label": "Next Page"}} onChangePage={this.handleChangePage}
                            onChangeRowsPerPage={this.handleChangeRowsPerPage} page={page}/>
                    </React.Fragment>
                )
                : (
                    <div>No Traces Found</div>
                )
        );
    };

}

TracesList.propTypes = {
    classes: PropTypes.any.isRequired,
    history: PropTypes.shape({
        push: PropTypes.func.isRequired
    }).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator),
    searchResults: PropTypes.shape({
        rootSpans: PropTypes.arrayOf(PropTypes.shape({
            traceId: PropTypes.string.isRequired,
            rootCellName: PropTypes.string.isRequired,
            rootServiceName: PropTypes.string.isRequired,
            rootOperationName: PropTypes.string.isRequired,
            rootStartTime: PropTypes.number.isRequired,
            rootDuration: PropTypes.number.isRequired
        })).isRequired,
        spanCounts: PropTypes.arrayOf(PropTypes.shape({
            traceId: PropTypes.string.isRequired,
            cellName: PropTypes.string.isRequired,
            serviceName: PropTypes.string.isRequired,
            count: PropTypes.number.isRequired
        })).isRequired
    }).isRequired
};

export default withStyles(styles)(withRouter(withColor(TracesList)));
