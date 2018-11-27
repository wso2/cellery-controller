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

import Button from "@material-ui/core/Button";
import {ColorGeneratorConstants} from "../common/color";
import FormControl from "@material-ui/core/FormControl/FormControl";
import Grid from "@material-ui/core/Grid/Grid";
import HttpUtils from "../common/utils/httpUtils";
import InputAdornment from "@material-ui/core/InputAdornment/InputAdornment";
import InputLabel from "@material-ui/core/InputLabel/InputLabel";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import NotificationUtils from "../common/utils/notificationUtils";
import Paper from "@material-ui/core/Paper/Paper";
import PropTypes from "prop-types";
import QueryUtils from "../common/utils/queryUtils";
import React from "react";
import SearchResult from "./SearchResult";
import Select from "@material-ui/core/Select/Select";
import Span from "./utils/span";
import TextField from "@material-ui/core/TextField/TextField";
import TopToolbar from "../common/TopToolbar";
import Typography from "@material-ui/core/Typography/Typography";
import {withConfig} from "../common/config";
import withStyles from "@material-ui/core/styles/withStyles";
import {ConfigConstants, ConfigHolder} from "../common/config/configHolder";

const styles = (theme) => ({
    container: {
        padding: theme.spacing.unit * 3
    },
    subheading: {
        marginBottom: theme.spacing.unit * 2
    },
    formControl: {
        marginBottom: theme.spacing.unit * 2
    },
    durationTextField: {
        marginTop: theme.spacing.unit * 2
    },
    startInputAdornment: {
        marginRight: theme.spacing.unit * 2,
        marginBottom: theme.spacing.unit * 2
    },
    searchForm: {
        marginBottom: Number(theme.spacing.unit)
    },
    resultContainer: {
        marginTop: theme.spacing.unit * 3
    }
});

/**
 * Trace Search.
 */
class Search extends React.Component {

    constructor(props) {
        super(props);
        const {location} = props;

        const queryParams = HttpUtils.parseQueryParams(location.search);
        this.state = Search.generateValidState({
            data: {
                cells: [],
                microservices: [],
                operations: []
            },
            filter: {
                cell: queryParams.cell ? queryParams.cell : Search.Constants.ALL_VALUE,
                microservice: queryParams.microservice ? queryParams.microservice : Search.Constants.ALL_VALUE,
                operation: queryParams.operation ? queryParams.operation : Search.Constants.ALL_VALUE,
                tags: queryParams.tags ? queryParams.operation : "",
                minDuration: queryParams.minDuration ? queryParams.operation : "",
                minDurationMultiplier: 1,
                maxDuration: queryParams.maxDuration ? queryParams.operation : "",
                maxDurationMultiplier: 1
            },
            metaData: {
                availableMicroservices: [],
                availableOperations: []
            },
            searchResults: []
        });

        this.loadCellData = this.loadCellData.bind(this);
        this.getChangeHandler = this.getChangeHandler.bind(this);
        this.search = this.search.bind(this);
    }

    componentDidMount() {
        const {location} = this.props;
        const queryParams = HttpUtils.parseQueryParams(location.search);
        let isQueryParamsEmpty = true;
        for (const key in queryParams) {
            if (queryParams.hasOwnProperty(key) && queryParams[key]) {
                isQueryParamsEmpty = false;
            }
        }
        if (!isQueryParamsEmpty) {
            this.search();
        }
    }

    render() {
        const {classes} = this.props;
        const {data, filter, metaData, searchResults} = this.state;

        const createMenuItemForSelect = (itemNames) => itemNames.map(
            (itemName) => (<MenuItem key={itemName} value={itemName}>{itemName}</MenuItem>)
        );

        return (
            <React.Fragment>
                <TopToolbar title={"Distributed Tracing"} onUpdate={this.loadCellData}/>
                <Paper className={classes.container}>
                    <Typography variant="h6" color="inherit" className={classes.subheading}>
                        Search Traces
                    </Typography>
                    <Grid container justify={"flex-start"} className={classes.searchForm}>
                        <Grid container justify={"flex-start"} spacing={24}>
                            <Grid item xs={3}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <InputLabel htmlFor="cell" shrink={true}>Cell</InputLabel>
                                    <Select value={filter.cell} onChange={this.getChangeHandler("cell")}
                                        inputProps={{name: "cell", id: "cell"}}>
                                        <MenuItem key={Search.Constants.ALL_VALUE} value={Search.Constants.ALL_VALUE}>
                                            {Search.Constants.ALL_VALUE}
                                        </MenuItem>
                                        {createMenuItemForSelect(data.cells)}
                                    </Select>
                                </FormControl>
                            </Grid>
                            <Grid item xs={3}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <InputLabel htmlFor="microservice" shrink={true}>Microservice</InputLabel>
                                    <Select value={filter.microservice} onChange={this.getChangeHandler("microservice")}
                                        inputProps={{name: "microservice", id: "microservice"}}>
                                        <MenuItem key={Search.Constants.ALL_VALUE} value={Search.Constants.ALL_VALUE}>
                                            {Search.Constants.ALL_VALUE}
                                        </MenuItem>
                                        {createMenuItemForSelect(metaData.availableMicroservices)}
                                    </Select>
                                </FormControl>
                            </Grid>
                            <Grid item xs={3}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <InputLabel htmlFor="operation" shrink={true}>Operation</InputLabel>
                                    <Select value={filter.operation} onChange={this.getChangeHandler("operation")}
                                        inputProps={{name: "operation", id: "operation"}}>
                                        <MenuItem key={Search.Constants.ALL_VALUE} value={Search.Constants.ALL_VALUE}>
                                            {Search.Constants.ALL_VALUE}
                                        </MenuItem>
                                        {createMenuItemForSelect(metaData.availableOperations)}
                                    </Select>
                                </FormControl>
                            </Grid>
                        </Grid>
                        <Grid container justify={"flex-start"} spacing={24}>
                            <Grid item xs={6}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <TextField label="Tags" id="tags" value={filter.tags} InputLabelProps={{shrink: true}}
                                        onChange={this.getChangeHandler("tags")} placeholder={"Eg: http.status_code=200"}/>
                                </FormControl>
                            </Grid>
                            <Grid item xs={3}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <InputLabel htmlFor="min-duration" shrink={true}>Duration</InputLabel>
                                    <TextField id="min-duration" value={filter.minDuration}
                                        className={classes.durationTextField}
                                        onChange={this.getChangeHandler("minDuration")} type="number"
                                        placeholder={"Eg: 10"}
                                        InputProps={{
                                            startAdornment: (
                                                <InputAdornment className={classes.startInputAdornment}
                                                    variant="filled" position="start">Min</InputAdornment>
                                            ),
                                            endAdornment: (
                                                <InputAdornment variant="filled" position="end">
                                                    <Select value={filter.minDurationMultiplier}
                                                        onChange={this.getChangeHandler("minDurationMultiplier")}
                                                        inputProps={{
                                                            name: "min-duration-multiplier",
                                                            id: "min-duration-multiplier"
                                                        }}>
                                                        <MenuItem value={1}>ms</MenuItem>
                                                        <MenuItem value={1000}>s</MenuItem>
                                                    </Select></InputAdornment>
                                            )
                                        }}/>
                                </FormControl>
                            </Grid>
                            <Grid item xs={3}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <TextField id="max-duration" value={filter.maxDuration}
                                        className={classes.durationTextField}
                                        onChange={this.getChangeHandler("maxDuration")} type="number"
                                        placeholder={"Eg: 1,000"}
                                        InputProps={{
                                            startAdornment: (
                                                <InputAdornment className={classes.startInputAdornment}
                                                    variant="filled" position="start">Max</InputAdornment>
                                            ),
                                            endAdornment: (
                                                <InputAdornment variant="filled" position="end">
                                                    <Select value={filter.maxDurationMultiplier}
                                                        onChange={this.getChangeHandler("maxDurationMultiplier")}
                                                        inputProps={{
                                                            name: "max-duration-multiplier",
                                                            id: "max-duration-multiplier"
                                                        }}>
                                                        <MenuItem value={1}>ms</MenuItem>
                                                        <MenuItem value={1000}>s</MenuItem>
                                                    </Select>
                                                </InputAdornment>
                                            )
                                        }}/>
                                </FormControl>
                            </Grid>
                        </Grid>
                    </Grid>
                    <Button variant="contained" color="primary" onClick={this.search}>Search</Button>
                    <div className={classes.resultContainer}>
                        <SearchResult data={searchResults}/>
                    </div>
                </Paper>
            </React.Fragment>
        );
    }

    /**
     * Search for traces.
     * Call the backend and search for traces.
     *
     * @param {boolean} showOverlay Show the overlay while loading
     */
    loadCellData(showOverlay) {
        const {config} = this.props;
        const self = this;

        if (showOverlay) {
            NotificationUtils.showLoadingOverlay("Loading Cell Information", config);
        }
        HttpUtils.callBackendAPI(
            {
                url: "/cells",
                method: "POST"
            },
            config
        ).then((data) => {
            const cells = [];
            const microservices = [];
            const operations = [];

            for (let i = 0; i < data.length; i++) {
                const span = new Span(data[i]);
                const cell = span.getCell();

                const cellName = (cell ? cell.name : null);
                const serviceName = span.serviceName;
                const operationName = span.operationName;

                if (cellName) {
                    if (!cells.includes(cellName)) {
                        cells.push(cellName);
                    }
                    if (!microservices.map((service) => service.name).includes(serviceName)) {
                        microservices.push({
                            name: serviceName,
                            cell: cellName
                        });
                    }
                    if (!operations.map((operation) => operation.name).includes(operationName)) {
                        operations.push({
                            name: operationName,
                            microservice: serviceName,
                            cell: cellName
                        });
                    }
                }
            }

            self.setState(Search.generateValidState({
                ...self.state,
                data: {
                    cells: cells,
                    microservices: microservices,
                    operations: operations
                }
            }));
            NotificationUtils.hideLoadingOverlay(config);
        }).catch(() => {
            NotificationUtils.hideLoadingOverlay(config);
        });
    }

    /**
     * Get the on change handler for a particular state filter attribute.
     *
     * @param {string} name The name of the filter
     * @returns {Function} The on change handler
     */
    getChangeHandler(name) {
        const self = this;
        return (event) => {
            self.setState(Search.generateValidState({
                ...this.state,
                filter: {
                    ...this.state.filter,
                    [name]: event.target.value
                }
            }));
        };
    }

    search() {
        const {
            cell, microservice, operation, tags, minDuration, minDurationMultiplier, maxDuration, maxDurationMultiplier
        } = this.state.filter;
        const {config} = this.props;
        const self = this;

        // Build search object
        const search = {};
        const addSearchParam = (key, value) => {
            if (value && value !== Search.Constants.ALL_VALUE) {
                search[key] = value;
            }
        };
        addSearchParam("cellName", cell);
        addSearchParam("serviceName", microservice);
        addSearchParam("operationName", operation);
        addSearchParam("tags", tags);
        addSearchParam("minDuration", minDuration * minDurationMultiplier);
        addSearchParam("maxDuration", maxDuration * maxDurationMultiplier);
        addSearchParam("queryStartTime",
            QueryUtils.parseTime(config.get(ConfigConstants.GLOBAL_FILTER).startTime).valueOf());
        addSearchParam("queryEndTime",
            QueryUtils.parseTime(config.get(ConfigConstants.GLOBAL_FILTER).endTime).valueOf());

        NotificationUtils.showLoadingOverlay("Searching for Traces", config);
        HttpUtils.callBackendAPI(
            {
                url: "/tracing/search",
                method: "POST",
                data: search
            },
            config
        ).then((data) => {
            const traces = {};
            for (let i = 0; i < data.length; i++) {
                const dataItem = data[i];
                if (!traces[dataItem.traceId]) {
                    traces[dataItem.traceId] = {};
                }
                if (!traces[dataItem.traceId][dataItem.cellName]) {
                    traces[dataItem.traceId][dataItem.cellName] = {};
                }
                if (!traces[dataItem.traceId][dataItem.cellName][dataItem.serviceName]) {
                    traces[dataItem.traceId][dataItem.cellName][dataItem.serviceName] = {};
                }
                const info = traces[dataItem.traceId][dataItem.cellName][dataItem.serviceName];
                info.count = dataItem.count;
                info.rootServiceName = dataItem.rootServiceName;
                info.rootOperationName = dataItem.rootOperationName;
                info.rootStartTime = dataItem.rootStartTime;
                info.rootDuration = dataItem.rootDuration;
            }
            const fillResult = (cellName, services, result) => {
                for (const serviceName in services) {
                    if (services.hasOwnProperty(serviceName)) {
                        const info = services[serviceName];

                        const span = new Span({
                            cellName: cellName,
                            serviceName: serviceName
                        });
                        const cell = span.getCell();

                        let cellNameKey;
                        if (span.isFromVICKSystemComponent()) {
                            cellNameKey = ColorGeneratorConstants.VICK;
                        } else if (span.isFromIstioSystemComponent()) {
                            cellNameKey = ColorGeneratorConstants.ISTIO;
                        } else {
                            cellNameKey = cell.name;
                        }

                        result.rootServiceName = info.rootServiceName;
                        result.rootOperationName = info.rootOperationName;
                        result.rootStartTime = info.rootStartTime;
                        result.rootDuration = info.rootDuration;
                        result.services.push({
                            cellName: cellNameKey,
                            serviceName: span.serviceName,
                            count: info.count
                        });
                    }
                }
            };
            const searchResults = [];
            for (const traceId in traces) {
                if (traces.hasOwnProperty(traceId)) {
                    const cells = traces[traceId];
                    const result = {
                        traceId: traceId,
                        services: []
                    };

                    for (const cellName in cells) {
                        if (cells.hasOwnProperty(cellName)) {
                            fillResult(cellName, cells[cellName], result);
                        }
                    }
                    searchResults.push(result);
                }
            }
            self.setState(Search.generateValidState({
                ...self.state,
                searchResults: searchResults
            }));
            NotificationUtils.hideLoadingOverlay(config);
        }).catch(() => {
            NotificationUtils.hideLoadingOverlay(config);
        });
    }

    /**
     * Current state from which the new valid state should be generated.
     *
     * @param {Object} state The current state
     * @returns {Object} The new valid state
     */
    static generateValidState(state) {
        const {data, filter, metaData} = state;

        // Finding the available microservices to be selected
        const selectedCells = (filter.cell === Search.Constants.ALL_VALUE ? data.cells : [filter.cell]);
        const availableMicroservices = data.microservices
            .filter((microservice) => selectedCells.includes(microservice.cell))
            .map((microservice) => microservice.name);

        const selectedMicroservice = (filter.microservice && availableMicroservices.includes(filter.microservice))
            ? filter.microservice
            : Search.Constants.ALL_VALUE;

        // Finding the available operations to be selected
        const selectedMicroservices = (selectedMicroservice === Search.Constants.ALL_VALUE
            ? availableMicroservices
            : [selectedMicroservice]);
        const availableOperations = data.operations
            .filter((operation) => selectedMicroservices.includes(operation.microservice))
            .map((operation) => operation.name);

        const selectedOperation = (filter.operation && availableOperations.includes(filter.operation))
            ? filter.operation
            : Search.Constants.ALL_VALUE;

        return {
            ...state,
            filter: {
                ...filter,
                microservice: selectedMicroservice,
                operation: selectedOperation
            },
            metaData: {
                ...metaData,
                availableMicroservices: availableMicroservices,
                availableOperations: availableOperations
            }
        };
    }

}

Search.Constants = {
    ALL_VALUE: "All"
};

Search.propTypes = {
    classes: PropTypes.object.isRequired,
    location: PropTypes.shape({
        search: PropTypes.string.isRequired
    }).isRequired,
    match: PropTypes.shape({
        url: PropTypes.string.isRequired
    }).isRequired,
    config: PropTypes.instanceOf(ConfigHolder).isRequired
};

export default withStyles(styles)(withConfig(Search));
