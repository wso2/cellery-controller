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

import Constants from "../../../utils/constants";
import FormControl from "@material-ui/core/FormControl";
import HttpUtils from "../../../utils/api/httpUtils";
import InputLabel from "@material-ui/core/InputLabel";
import MetricsGraphs from "../metricsGraphs";
import NotificationUtils from "../../../utils/common/notificationUtils";
import QueryUtils from "../../../utils/common/queryUtils";
import React from "react";
import Select from "@material-ui/core/Select";
import StateHolder from "../../common/state/stateHolder";
import Typography from "@material-ui/core/Typography/Typography";
import withGlobalState from "../../common/state/index";
import {withStyles} from "@material-ui/core/styles";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    filters: {
        marginTop: theme.spacing.unit * 4,
        marginBottom: theme.spacing.unit * 4
    },
    formControl: {
        marginRight: theme.spacing.unit * 4,
        minWidth: 150
    },
    graphs: {
        marginBottom: theme.spacing.unit * 4
    },
    button: {
        marginTop: theme.spacing.unit * 2
    }
});

class Metrics extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            selectedType: props.globalFilterOverrides && props.globalFilterOverrides.selectedType
                ? props.globalFilterOverrides.selectedType
                : Constants.Dashboard.INBOUND,
            selectedCell: props.globalFilterOverrides && props.globalFilterOverrides.selectedCell
                ? props.globalFilterOverrides.selectedCell
                : Constants.Dashboard.ALL_VALUE,
            cells: [],
            cellData: [],
            isLoading: false
        };
    }

    componentDidMount = () => {
        const {globalState} = this.props;

        this.update(
            true,
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime),
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime)
        );
    };

    update = (isUserAction, startTime, endTime, selectedTypeOverride, selectedCellOverride) => {
        const {selectedType, selectedCell} = this.state;
        const queryStartTime = startTime.valueOf();
        const queryEndTime = endTime.valueOf();

        this.loadMetrics(
            isUserAction, queryStartTime, queryEndTime,
            selectedTypeOverride ? selectedTypeOverride : selectedType,
            selectedCellOverride ? selectedCellOverride : selectedCell
        );
        this.loadCellMetadata(isUserAction, queryStartTime, queryEndTime);
    };

    getFilterChangeHandler = (name) => (event) => {
        const {globalState, onFilterUpdate} = this.props;
        const {selectedType, selectedCell} = this.state;

        const newValue = event.target.value;
        this.setState({
            [name]: newValue
        });

        if (onFilterUpdate) {
            onFilterUpdate({
                selectedType: selectedType,
                selectedCell: selectedCell,
                [name]: newValue
            });
        }

        this.update(
            true,
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime),
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime),
            name === "selectedType" ? newValue : null,
            name === "selectedCell" ? newValue : null,
        );
    };

    loadCellMetadata = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell} = this.props;
        const self = this;

        const search = {
            queryStartTime: queryStartTime,
            queryEndTime: queryEndTime
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Info", globalState);
            self.setState({
                isLoading: true
            });
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/metadata${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            self.setState({
                cells: data.filter((datum) => Boolean(datum) && datum !== cell)
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                self.setState({
                    isLoading: false
                });
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                self.setState({
                    isLoading: false
                });
                NotificationUtils.showNotification(
                    "Failed to load cell information",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    loadMetrics = (isUserAction, queryStartTime, queryEndTime, selectedType, selectedCell) => {
        const {globalState, cell} = this.props;
        const self = this;

        // Creating the search params
        const search = {
            queryStartTime: queryStartTime,
            queryEndTime: queryEndTime
        };
        if (selectedCell !== Constants.Dashboard.ALL_VALUE) {
            if (selectedType === Constants.Dashboard.INBOUND) {
                search.sourceCell = selectedCell;
            } else {
                search.destinationCell = selectedCell;
            }
        }
        if (selectedType === Constants.Dashboard.INBOUND) {
            search.destinationCell = cell;
        } else {
            search.sourceCell = cell;
        }

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Metrics", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/metrics${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const cellData = data.map((datum) => ({
                timestamp: datum[0],
                httpResponseGroup: datum[1],
                totalResponseTimeMilliSec: datum[2],
                totalRequestSizeBytes: datum[3],
                totalResponseSizeBytes: datum[4],
                requestCount: datum[5]
            }));

            self.setState({
                cellData: cellData
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load cell metrics",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    render = () => {
        const {classes, cell} = this.props;
        const {selectedType, selectedCell, cells, cellData, isLoading} = this.state;

        const targetSourcePrefix = selectedType === Constants.Dashboard.INBOUND ? "Source" : "Target";

        return (
            isLoading
                ? null
                : (
                    <React.Fragment>
                        <div className={classes.filters}>
                            <FormControl className={classes.formControl}>
                                <InputLabel htmlFor="selected-type">Type</InputLabel>
                                <Select value={selectedType}
                                    onChange={this.getFilterChangeHandler("selectedType")}
                                    inputProps={{
                                        name: "selected-type",
                                        id: "selected-type"
                                    }}>
                                    <option value={Constants.Dashboard.INBOUND}>Inbound</option>
                                    <option value={Constants.Dashboard.OUTBOUND}>Outbound</option>
                                </Select>
                            </FormControl>
                            <FormControl className={classes.formControl}>
                                <InputLabel htmlFor="selected-cell">{targetSourcePrefix} Cell</InputLabel>
                                <Select value={selectedCell}
                                    onChange={this.getFilterChangeHandler("selectedCell")}
                                    inputProps={{
                                        name: "selected-cell",
                                        id: "selected-cell"
                                    }}>
                                    <option value={Constants.Dashboard.ALL_VALUE}>
                                        {Constants.Dashboard.ALL_VALUE}
                                    </option>
                                    {
                                        cells.map((cell) => (<option key={cell} value={cell}>{cell}</option>))
                                    }
                                </Select>
                            </FormControl>
                        </div>
                        <div className={classes.graphs}>
                            {
                                cellData.length > 0
                                    ? (
                                        <MetricsGraphs cell={cell} data={cellData}
                                            direction={selectedType === Constants.Dashboard.INBOUND ? "In" : "Out"}/>
                                    )
                                    : (
                                        <Typography>
                                            {
                                                selectedType === Constants.Dashboard.INBOUND
                                                    ? `No Requests from the selected cell to "${cell}" cell`
                                                    : `No Requests from "${cell}" cell to the selected cell`
                                            }
                                        </Typography>
                                    )
                            }
                        </div>
                    </React.Fragment>
                )
        );
    };

}

Metrics.propTypes = {
    classes: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    cell: PropTypes.string.isRequired,
    onFilterUpdate: PropTypes.func.isRequired,
    globalFilterOverrides: PropTypes.shape({
        selectedType: PropTypes.string,
        selectedCell: PropTypes.string
    })
};

export default withStyles(styles)(withGlobalState(Metrics));
