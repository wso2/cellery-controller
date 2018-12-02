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

import Constants from "../common/constants";
import Grid from "@material-ui/core/Grid/Grid";
import Paper from "@material-ui/core/Paper/Paper";
import PropTypes from "prop-types";
import React from "react";
import TablePagination from "@material-ui/core/TablePagination/TablePagination";
import Typography from "@material-ui/core/Typography/Typography";
import moment from "moment";
import {withRouter} from "react-router-dom";
import withStyles from "@material-ui/core/styles/withStyles";
import withColor, {ColorGenerator} from "../common/color";

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
    serviceName: {
        fontWeight: 500,
        fontSize: "normal"
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
    traceSubHeader: {
        backgroundColor: "#929292",
        color: "#ffffff",
        textAlign: "right",
        fontWeight: 400,
        fontSize: "small",
        padding: theme.spacing.unit
    },
    traceContent: {
        padding: theme.spacing.unit
    },
    serviceTag: {
        borderStyle: "solid",
        borderWidth: "thin",
        borderColor: "#000000",
        margin: theme.spacing.unit,
        display: "inline-block"
    },
    serviceTagColor: {
        height: "100%",
        width: theme.spacing.unit * 2,
        display: "table-cell"
    },
    serviceTagContent: {
        padding: theme.spacing.unit,
        display: "table-cell",
        fontSize: "small"
    }
});

class SearchResult extends React.Component {

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
     * @param {Event} event Event for the click event
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
        const {classes, data, colorGenerator} = this.props;
        const {rowsPerPage, page} = this.state;

        const cellNames = [];
        data.forEach(
            (result) => result.services.forEach(
                (service) => {
                    if (!cellNames.includes(service.cellNameKey)) {
                        cellNames.push(service.cellNameKey);
                    }
                }
            )
        );
        colorGenerator.addKeys(cellNames);

        return (
            data.length > 0
                ? (
                    <React.Fragment>
                        <Typography variant="h6" color="inherit" className={classes.subheading}>
                            Traces
                        </Typography>
                        {
                            data.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                                .map((result) => (
                                    <Paper key={result.traceId} className={classes.trace}
                                        onClick={(event) => this.loadTracePage(event, result.traceId)}>
                                        <Grid container className={classes.traceHeader}>
                                            <Grid item xs={8}>
                                                <span className={classes.serviceName}>
                                                    {`${result.rootServiceName} `}
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
                                            </Grid>
                                        </Grid>
                                        <div className={classes.traceSubHeader}>{result.rootDuration / 1000} s</div>
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
                                                        <div key={service.serviceName} className={classes.serviceTag}
                                                            onClick={
                                                                (event) => this.loadTracePage(event, result.traceId,
                                                                    service.cellNameKey, service.serviceName)}>
                                                            <div className={classes.serviceTagColor} style={{
                                                                backgroundColor: colorGenerator
                                                                    .getColor(service.cellNameKey)
                                                            }}/>
                                                            <div className={classes.serviceTagContent}>
                                                                {service.serviceName} ({service.count})
                                                            </div>
                                                        </div>
                                                    ))
                                            }
                                        </div>
                                    </Paper>
                                ))
                        }
                        <TablePagination
                            component="div"
                            count={data.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            labelRowsPerPage={"Traces Per Page"}
                            backIconButtonProps={{"aria-label": "Previous Page"}}
                            nextIconButtonProps={{"aria-label": "Next Page"}}
                            onChangePage={this.handleChangePage}
                            onChangeRowsPerPage={this.handleChangeRowsPerPage}/>
                    </React.Fragment>
                )
                : (
                    <div>No Results</div>
                )
        );
    };

}

SearchResult.propTypes = {
    classes: PropTypes.any.isRequired,
    history: PropTypes.shape({
        push: PropTypes.func.isRequired
    }).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator),
    data: PropTypes.arrayOf(PropTypes.shape({
        traceId: PropTypes.string.isRequired,
        rootServiceName: PropTypes.string.isRequired,
        rootOperationName: PropTypes.string.isRequired,
        rootStartTime: PropTypes.number.isRequired,
        rootDuration: PropTypes.number.isRequired,
        services: PropTypes.arrayOf(PropTypes.shape({
            cellNameKey: PropTypes.string.isRequired,
            serviceName: PropTypes.string.isRequired,
            count: PropTypes.number.isRequired
        })).isRequired
    })).isRequired
};

export default withStyles(styles)(withRouter(withColor(SearchResult)));
