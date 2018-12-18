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

import FormControl from "@material-ui/core/FormControl";
import HttpUtils from "../../common/utils/httpUtils";
import InputLabel from "@material-ui/core/InputLabel";
import MetricsGraphs from "../MetricsGraphs";
import NotificationUtils from "../../common/utils/notificationUtils";
import PropTypes from "prop-types";
import QueryUtils from "../../common/utils/queryUtils";
import React from "react";
import Select from "@material-ui/core/Select";
import StateHolder from "../../common/state/stateHolder";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core/styles";

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

    static ALL_VALUE = "All";
    static INBOUND = "inbound";
    static OUTBOUND = "outbound";

    constructor(props) {
        super(props);

        this.state = {
            selectedType: Metrics.OUTBOUND,
            selectedCell: Metrics.ALL_VALUE,
            cells: [],
            cellData: [],
            timeRange: 0
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

        this.setState({
            timeRange: (endTime - startTime)
        });
    };

    getFilterChangeHandler = (name) => (event) => {
        const {globalState} = this.props;

        const newValue = event.target.value;
        this.setState({
            [name]: newValue
        });

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
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load cell information",
                    StateHolder.NotificationLevels.ERROR,
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
        if (selectedCell !== Metrics.ALL_VALUE) {
            if (selectedType === Metrics.INBOUND) {
                search.sourceCell = selectedCell;
            } else {
                search.destinationCell = selectedCell;
            }
        }
        if (selectedType === Metrics.INBOUND) {
            search.destinationCell = cell;
        } else {
            search.sourceCell = cell;
        }

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Info", globalState);
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
                totalResponseTimeSec: datum[2],
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
                    "Failed to load cell information",
                    StateHolder.NotificationLevels.ERROR,
                    globalState
                );
            }
        });
    };

    render = () => {
        const {classes} = this.props;
        const {selectedType, selectedCell, cells, cellData} = this.state;

        return (
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
                            <option value={Metrics.INBOUND}>Inbound</option>
                            <option value={Metrics.OUTBOUND}>Outbound</option>
                        </Select>
                    </FormControl>
                    <FormControl className={classes.formControl}>
                        <InputLabel htmlFor="selected-cell">
                            {selectedType === Metrics.INBOUND ? "Source" : "Target"} Cell
                        </InputLabel>
                        <Select value={selectedCell}
                            onChange={this.getFilterChangeHandler("selectedCell")}
                            inputProps={{
                                name: "selected-cell",
                                id: "selected-cell"
                            }}>
                            <option value={Metrics.ALL_VALUE}>{Metrics.ALL_VALUE}</option>
                            {
                                cells.map((cell) => (<option key={cell} value={cell}>{cell}</option>))
                            }
                        </Select>
                    </FormControl>
                </div>
                <div className={classes.graphs}>
                    <MetricsGraphs data={cellData}/>
                </div>
            </React.Fragment>
        );
    };

}

Metrics.propTypes = {
    classes: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    cell: PropTypes.string.isRequired
};

export default withStyles(styles)(withGlobalState(Metrics));
