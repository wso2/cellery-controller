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
import withGlobalState from "../../common/state";
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
            selectedType: Constants.Dashboard.INBOUND,
            selectedCell: Constants.Dashboard.ALL_VALUE,
            selectedComponent: Constants.Dashboard.ALL_VALUE,
            components: [],
            metadata: {
                availableCells: [],
                availableComponents: [] // Filtered based on the selected cell
            },
            componentData: [],
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

    update = (isUserAction, startTime, endTime, selectedTypeOverride, selectedCellOverride,
        selectedComponentOverride) => {
        const {selectedType, selectedCell, selectedComponent} = this.state;
        const queryStartTime = startTime.valueOf();
        const queryEndTime = endTime.valueOf();

        this.loadMetrics(
            isUserAction, queryStartTime, queryEndTime,
            selectedTypeOverride ? selectedTypeOverride : selectedType,
            selectedCellOverride ? selectedCellOverride : selectedCell,
            selectedComponentOverride ? selectedComponentOverride : selectedComponent
        );
        this.loadComponentMetadata(isUserAction, queryStartTime, queryEndTime);
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
            name === "selectedComponent" ? newValue : null
        );
    };

    loadComponentMetadata = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell, component} = this.props;
        const self = this;

        const search = {
            queryStartTime: queryStartTime,
            queryEndTime: queryEndTime
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Component Info", globalState);
            self.setState({
                isLoading: true
            });
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/microservices/metadata${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            self.setState({
                components: data
                    .filter((datum) => (datum.cell !== cell || datum.component !== component))
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
                    "Failed to load component information",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    loadMetrics = (isUserAction, queryStartTime, queryEndTime, selectedType, selectedCell, selectedComponent) => {
        const {globalState, cell, component} = this.props;
        const self = this;

        // Creating the search params
        const search = {
            queryStartTime: queryStartTime,
            queryEndTime: queryEndTime
        };
        if (selectedCell !== Constants.Dashboard.ALL_VALUE) {
            if (selectedType === Constants.Dashboard.INBOUND) {
                search.sourceCell = selectedCell;
                search.sourceComponent = selectedComponent;
            } else {
                search.destinationCell = selectedCell;
                search.destinationComponent = selectedComponent;
            }
        }
        if (selectedType === Constants.Dashboard.INBOUND) {
            search.destinationCell = cell;
            search.destinationComponent = component;
        } else {
            search.sourceCell = cell;
            search.sourceComponent = component;
        }

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Component Metrics", globalState);
            self.setState({
                isLoading: true
            });
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/microservices/metrics${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const componentData = data.map((datum) => ({
                timestamp: datum[0],
                httpResponseGroup: datum[1],
                totalResponseTimeMilliSec: datum[2],
                totalRequestSizeBytes: datum[3],
                totalResponseSizeBytes: datum[4],
                requestCount: datum[5]
            }));

            self.setState({
                componentData: componentData
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
                    "Failed to load component metrics",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    static getDerivedStateFromProps = (props, state) => {
        const {cell, component} = props;
        const {components, selectedCell, selectedComponent} = state;

        const availableCells = [];
        components.forEach((componentDatum) => {
            // Validating whether at-least one relevant component in this cell exists
            let hasRelevantComponent = true;
            if (Boolean(componentDatum.cell) && componentDatum.cell === cell) {
                const relevantComponents = components.find(
                    (datum) => cell === datum.cell && datum.name !== component);
                hasRelevantComponent = Boolean(relevantComponents);
            }

            if (hasRelevantComponent && Boolean(componentDatum.cell)
                    && !availableCells.includes(componentDatum.cell)) {
                availableCells.push(componentDatum.cell);
            }
        });

        const availableComponents = [];
        components.forEach((componentDatum) => {
            if (Boolean(componentDatum.name) && Boolean(componentDatum.cell)
                && (selectedCell === Constants.Dashboard.ALL_VALUE || componentDatum.cell === selectedCell)
                && (componentDatum.cell !== cell || componentDatum.name !== component)
                && !availableComponents.includes(componentDatum.name)) {
                availableComponents.push(componentDatum.name);
            }
        });

        const selectedComponentToShow = availableComponents.includes(selectedComponent)
            ? selectedComponent
            : Constants.Dashboard.ALL_VALUE;

        return {
            ...state,
            selectedComponent: selectedComponentToShow,
            metadata: {
                availableCells: availableCells,
                availableComponents: availableComponents
            }
        };
    };

    render = () => {
        const {classes, cell, component} = this.props;
        const {selectedType, selectedCell, selectedComponent, componentData, metadata, isLoading} = this.state;

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
                                    <option value={Constants.Dashboard.INBOUND}>{Constants.Dashboard.INBOUND}</option>
                                    <option value={Constants.Dashboard.OUTBOUND}>{Constants.Dashboard.OUTBOUND}</option>
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
                                        metadata.availableCells.map(
                                            (cell) => (<option key={cell} value={cell}>{cell}</option>))
                                    }
                                </Select>
                            </FormControl>
                            <FormControl className={classes.formControl}>
                                <InputLabel htmlFor="selected-component">
                                    {targetSourcePrefix} Component
                                </InputLabel>
                                <Select value={selectedComponent}
                                    onChange={this.getFilterChangeHandler("selectedComponent")}
                                    inputProps={{
                                        name: "selected-component",
                                        id: "selected-component"
                                    }}>
                                    <option value={Constants.Dashboard.ALL_VALUE}>
                                        {Constants.Dashboard.ALL_VALUE}
                                    </option>
                                    {
                                        metadata.availableComponents.map((component) => (
                                            <option key={component} value={component}>{component}</option>
                                        ))
                                    }
                                </Select>
                            </FormControl>
                        </div>
                        <div className={classes.graphs}>
                            {
                                componentData.length > 0
                                    ? (
                                        <MetricsGraphs cell={cell} component={component} data={componentData}
                                            direction={selectedType === Constants.Dashboard.INBOUND ? "In" : "Out"}/>
                                    )
                                    : (
                                        <Typography>
                                            {
                                                selectedType === Constants.Dashboard.INBOUND
                                                    ? "No Requests from the selected component "
                                                        + `to the "${cell}" cell's "${component}" component`
                                                    : `No Requests from the "${cell}" cell's "${component}" `
                                                        + "component to the selected component"
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
    component: PropTypes.string.isRequired
};

export default withStyles(styles)(withGlobalState(Metrics));
