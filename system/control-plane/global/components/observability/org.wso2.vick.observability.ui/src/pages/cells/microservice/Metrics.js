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
import Typography from "@material-ui/core/Typography/Typography";
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
    static INBOUND = "Inbound";
    static OUTBOUND = "Outbound";

    constructor(props) {
        super(props);

        this.state = {
            selectedType: Metrics.INBOUND,
            selectedCell: Metrics.ALL_VALUE,
            selectedMicroservice: Metrics.ALL_VALUE,
            microservices: [],
            metadata: {
                availableCells: [],
                availableMicroservices: [] // Filtered based on the selected cell
            },
            microserviceData: []
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

    update = (isUserAction, startTime, endTime, selectedTypeOverride, selectedCellOverride,
        selectedMicroserviceOverride) => {
        const {selectedType, selectedCell, selectedMicroservice} = this.state;
        const queryStartTime = startTime.valueOf();
        const queryEndTime = endTime.valueOf();

        this.loadMetrics(
            isUserAction, queryStartTime, queryEndTime,
            selectedTypeOverride ? selectedTypeOverride : selectedType,
            selectedCellOverride ? selectedCellOverride : selectedCell,
            selectedMicroserviceOverride ? selectedMicroserviceOverride : selectedMicroservice
        );
        this.loadMicroserviceMetadata(isUserAction, queryStartTime, queryEndTime);
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
            name === "selectedMicroservice" ? newValue : null
        );
    };

    loadMicroserviceMetadata = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell, microservice} = this.props;
        const self = this;

        const search = {
            queryStartTime: queryStartTime,
            queryEndTime: queryEndTime
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Microservice Info", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/microservices/metadata${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            self.setState({
                microservices: data
                    .filter((datum) => (datum.cell !== cell || datum.microservice !== microservice))
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load microservice information",
                    StateHolder.NotificationLevels.ERROR,
                    globalState
                );
            }
        });
    };

    loadMetrics = (isUserAction, queryStartTime, queryEndTime, selectedType, selectedCell, selectedMicroservice) => {
        const {globalState, cell, microservice} = this.props;
        const self = this;

        // Creating the search params
        const search = {
            queryStartTime: queryStartTime,
            queryEndTime: queryEndTime
        };
        if (selectedCell !== Metrics.ALL_VALUE) {
            if (selectedType === Metrics.INBOUND) {
                search.sourceCell = selectedCell;
                search.sourceMicroservice = selectedMicroservice;
            } else {
                search.destinationCell = selectedCell;
                search.destinationMicroservice = selectedMicroservice;
            }
        }
        if (selectedType === Metrics.INBOUND) {
            search.destinationCell = cell;
            search.destinationMicroservice = microservice;
        } else {
            search.sourceCell = cell;
            search.sourceMicroservice = microservice;
        }

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Microservice Metrics", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/microservices/metrics${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const microserviceData = data.map((datum) => ({
                timestamp: datum[0],
                httpResponseGroup: datum[1],
                totalResponseTimeMilliSec: datum[2],
                totalRequestSizeBytes: datum[3],
                totalResponseSizeBytes: datum[4],
                requestCount: datum[5]
            }));

            self.setState({
                microserviceData: microserviceData
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load microservice metrics",
                    StateHolder.NotificationLevels.ERROR,
                    globalState
                );
            }
        });
    };

    static getDerivedStateFromProps = (props, state) => {
        const {cell, microservice} = props;
        const {microservices, selectedCell, selectedMicroservice} = state;

        const availableCells = [];
        microservices.forEach((microserviceDatum) => {
            // Validating whether at-least one relevant microservice in this cell exists
            let hasRelevantMicroservice = true;
            if (Boolean(microserviceDatum.cell) && microserviceDatum.cell === cell) {
                const relevantMicroservices = microservices.find(
                    (datum) => cell === datum.cell && datum.name !== microservice);
                hasRelevantMicroservice = Boolean(relevantMicroservices);
            }

            if (hasRelevantMicroservice && Boolean(microserviceDatum.cell)
                    && !availableCells.includes(microserviceDatum.cell)) {
                availableCells.push(microserviceDatum.cell);
            }
        });

        const availableMicroservices = [];
        microservices.forEach((microserviceDatum) => {
            if (Boolean(microserviceDatum.name) && Boolean(microserviceDatum.cell)
                && (selectedCell === Metrics.ALL_VALUE || microserviceDatum.cell === selectedCell)
                && (microserviceDatum.cell !== cell || microserviceDatum.name !== microservice)
                && !availableMicroservices.includes(microserviceDatum.name)) {
                availableMicroservices.push(microserviceDatum.name);
            }
        });

        const selectedMicroserviceToShow = availableMicroservices.includes(selectedMicroservice)
            ? selectedMicroservice
            : Metrics.ALL_VALUE;

        return {
            ...state,
            selectedMicroservice: selectedMicroserviceToShow,
            metadata: {
                availableCells: availableCells,
                availableMicroservices: availableMicroservices
            }
        };
    };

    render = () => {
        const {classes, cell, microservice} = this.props;
        const {selectedType, selectedCell, selectedMicroservice, microserviceData, metadata} = this.state;

        const targetSourcePrefix = selectedType === Metrics.INBOUND ? "Source" : "Target";

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
                            <option value={Metrics.INBOUND}>{Metrics.INBOUND}</option>
                            <option value={Metrics.OUTBOUND}>{Metrics.OUTBOUND}</option>
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
                            <option value={Metrics.ALL_VALUE}>{Metrics.ALL_VALUE}</option>
                            {
                                metadata.availableCells.map((cell) => (<option key={cell} value={cell}>{cell}</option>))
                            }
                        </Select>
                    </FormControl>
                    <FormControl className={classes.formControl}>
                        <InputLabel htmlFor="selected-microservice">{targetSourcePrefix} Microservice</InputLabel>
                        <Select value={selectedMicroservice}
                            onChange={this.getFilterChangeHandler("selectedMicroservice")}
                            inputProps={{
                                name: "selected-microservice",
                                id: "selected-microservice"
                            }}>
                            <option value={Metrics.ALL_VALUE}>{Metrics.ALL_VALUE}</option>
                            {
                                metadata.availableMicroservices.map((microservice) => (
                                    <option key={microservice} value={microservice}>{microservice}</option>
                                ))
                            }
                        </Select>
                    </FormControl>
                </div>
                <div className={classes.graphs}>
                    {
                        microserviceData.length > 0
                            ? (
                                <MetricsGraphs data={microserviceData}/>
                            )
                            : (
                                <Typography>
                                    {
                                        selectedType === Metrics.INBOUND
                                            ? "No Requests from the selected microservice "
                                                + `to the "${cell}" cell's "${microservice}" microservice`
                                            : `No Requests from the "${cell}" cell's "${microservice}" microservice `
                                                + "to the selected microservice"
                                    }
                                </Typography>
                            )
                    }
                </div>
            </React.Fragment>
        );
    };

}

Metrics.propTypes = {
    classes: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    cell: PropTypes.string.isRequired,
    microservice: PropTypes.string.isRequired
};

export default withStyles(styles)(withGlobalState(Metrics));
