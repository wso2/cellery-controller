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

import Button from "@material-ui/core/Button";
import ChipInput from "material-ui-chip-input";
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
import TopToolbar from "../common/toptoolbar";
import Typography from "@material-ui/core/Typography/Typography";
import withStyles from "@material-ui/core/styles/withStyles";
import withGlobalState, {StateHolder} from "../common/state";

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

class Search extends React.Component {

    static ALL_VALUE = "All";

    constructor(props) {
        super(props);
        const {location} = props;

        const queryParams = HttpUtils.parseQueryParams(location.search);
        this.state = {
            data: {
                cells: [],
                microservices: [],
                operations: []
            },
            filter: {
                cell: queryParams.cell ? queryParams.cell : Search.ALL_VALUE,
                microservice: queryParams.microservice ? queryParams.microservice : Search.ALL_VALUE,
                operation: queryParams.operation ? queryParams.operation : Search.ALL_VALUE,
                tags: queryParams.tags ? JSON.parse(queryParams.tags) : {},
                minDuration: queryParams.minDuration ? queryParams.minDuration : "",
                minDurationMultiplier: queryParams.minDurationMultiplier ? queryParams.minDurationMultiplier : 1,
                maxDuration: queryParams.maxDuration ? queryParams.maxDuration : "",
                maxDurationMultiplier: queryParams.maxDurationMultiplier ? queryParams.maxDurationMultiplier : 1
            },
            metaData: {
                availableMicroservices: [],
                availableOperations: []
            },
            hasSearchCompleted: false,
            searchResults: []
        };
    }

    componentDidMount = () => {
        const {location} = this.props;
        const queryParams = HttpUtils.parseQueryParams(location.search);
        let isQueryParamsEmpty = true;
        for (const key in queryParams) {
            if (queryParams.hasOwnProperty(key) && queryParams[key]) {
                isQueryParamsEmpty = false;
            }
        }

        if (!isQueryParamsEmpty) {
            this.search(true);
        }
    };

    render = () => {
        const {classes} = this.props;
        const {data, filter, metaData, hasSearchCompleted, searchResults} = this.state;

        const createMenuItemForSelect = (itemNames) => itemNames.map(
            (itemName) => (<MenuItem key={itemName} value={itemName}>{itemName}</MenuItem>)
        );

        const tagChips = [];
        for (const tagKey in filter.tags) {
            if (filter.tags.hasOwnProperty(tagKey)) {
                tagChips.push(`${tagKey}=${filter.tags[tagKey]}`);
            }
        }

        return (
            <React.Fragment>
                <TopToolbar title={"Distributed Tracing"} onUpdate={this.onGlobalRefresh}/>
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
                                        <MenuItem key={Search.ALL_VALUE} value={Search.ALL_VALUE}>
                                            {Search.ALL_VALUE}
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
                                        <MenuItem key={Search.ALL_VALUE} value={Search.ALL_VALUE}>
                                            {Search.ALL_VALUE}
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
                                        <MenuItem key={Search.ALL_VALUE} value={Search.ALL_VALUE}>
                                            {Search.ALL_VALUE}
                                        </MenuItem>
                                        {createMenuItemForSelect(metaData.availableOperations)}
                                    </Select>
                                </FormControl>
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid container justify={"flex-start"} spacing={24}>
                        <Grid item xs={6}>
                            <FormControl className={classes.formControl} fullWidth={true}>
                                <ChipInput label="Tags" InputLabelProps={{shrink: true}} value={tagChips}
                                    onChange={this.handleTagsChange} onDelete={this.handleTagsChange}
                                    placeholder={"Eg: http.status_code=200"}
                                />
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
                    <Button variant="contained" color="primary" onClick={this.onSearchButtonClick}>Search</Button>
                    {
                        hasSearchCompleted
                            ? (
                                <div className={classes.resultContainer}>
                                    <SearchResult data={searchResults}/>
                                </div>
                            )
                            : null
                    }
                </Paper>
            </React.Fragment>
        );
    };

    onSearchButtonClick = () => {
        const {history, match, location} = this.props;
        const {filter} = this.state;

        // Updating the URL to ensure that the user can come back to this page
        const searchString = HttpUtils.generateQueryParamString({
            ...filter,
            tags: JSON.stringify(filter.tags)
        });
        history.replace(match.url + searchString, {
            ...location.state
        });

        this.search(true);
    };

    onGlobalRefresh = (isUserAction, queryStartTime, queryEndTime) => {
        if (this.state.hasSearchCompleted) {
            this.search(isUserAction);
        }
        this.loadCellData(isUserAction && !this.state.hasSearchCompleted, queryStartTime, queryEndTime);
    };

    /**
     * Search for traces.
     * Call the backend and search for traces.
     *
     * @param {boolean} isUserAction Show the overlay while loading
     * @param {number} queryStartTime Start time of the global filter
     * @param {number} queryEndTime End time of the global filter
     */
    loadCellData = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState} = this.props;
        const self = this;
        const filter = {
            queryStartTime: queryStartTime.valueOf(),
            queryEndTime: queryEndTime.valueOf()
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Information", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/traces/metadata${HttpUtils.generateQueryParamString(filter)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const cells = [];
            const microservices = [];
            const operations = [];

            const cellData = data.map((dataItem) => ({
                cell: dataItem[0],
                serviceName: dataItem[1],
                operationName: dataItem[2]
            }));

            for (let i = 0; i < cellData.length; i++) {
                const span = new Span(cellData[i]);
                const cell = span.cell;

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

            self.setState((prevState) => ({
                ...prevState,
                data: {
                    cells: cells,
                    microservices: microservices,
                    operations: operations
                }
            }));
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load Cell Data",
                    StateHolder.NotificationLevels.ERROR,
                    globalState
                );
            }
        });
    };

    /**
     * Get the on change handler for a particular state filter attribute.
     *
     * @param {string} name The name of the filter
     * @returns {Function} The on change handler
     */
    getChangeHandler = (name) => (event) => {
        const value = event.target.value;
        this.setState((prevState) => ({
            ...prevState,
            filter: {
                ...prevState.filter,
                [name]: value
            }
        }));
    };

    /**
     * Handle the tags changing in the search.
     *
     * @param {Array.<string>} chips The chips in the tag search input
     */
    handleTagsChange = (chips) => {
        const parseChip = (chip) => {
            const chipContent = chip.split("=");
            return {
                key: chipContent[0].trim(),
                value: chipContent[1].trim()
            };
        };

        // Generating tags object
        let tags;
        if (typeof chips === "string") { // Delete tag
            tags = {...this.state.filter.tags};
            const tag = parseChip(chips);
            Reflect.deleteProperty(tags, tag.key);
        } else { // Tag change
            tags = {};
            for (let i = 0; i < chips.length; i++) {
                const tag = parseChip(chips[i]);
                tags[tag.key] = tag.value;
            }
        }

        this.setState((prevState) => ({
            ...prevState,
            filter: {
                ...prevState.filter,
                tags: tags
            }
        }));
    };

    search = (isUserAction) => {
        const {
            cell, microservice, operation, tags, minDuration, minDurationMultiplier, maxDuration, maxDurationMultiplier
        } = this.state.filter;
        const {globalState} = this.props;
        const self = this;

        // Build search object
        const search = {};
        const addSearchParam = (key, value) => {
            if (value && value !== Search.ALL_VALUE) {
                search[key] = value;
            }
        };
        addSearchParam("cell", cell);
        addSearchParam("serviceName", microservice);
        addSearchParam("operationName", operation);
        addSearchParam("tags", JSON.stringify(Object.keys(tags).length > 0 ? tags : {}));
        addSearchParam("minDuration", minDuration * minDurationMultiplier);
        addSearchParam("maxDuration", maxDuration * maxDurationMultiplier);
        addSearchParam("queryStartTime",
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime).valueOf());
        addSearchParam("queryEndTime",
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime).valueOf());

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
            const rootSpans = data.rootSpans
                .map((dataItem) => ({
                    traceId: dataItem[0],
                    rootServiceName: dataItem[1],
                    rootOperationName: dataItem[2],
                    rootStartTime: dataItem[3],
                    rootDuration: dataItem[4]
                }))
                .reduce((accumulator, dataItem) => {
                    accumulator[dataItem.traceId] = dataItem;
                    return accumulator;
                }, {});
            const searchResults = data.spanCounts
                .map((dataItem) => ({
                    traceId: dataItem[0],
                    cellNameKey: dataItem[1],
                    serviceName: dataItem[2],
                    count: dataItem[3]
                }))
                .reduce((accumulator, dataItem) => {
                    if (accumulator[dataItem.traceId]) {
                        if (!accumulator[dataItem.traceId].services) {
                            accumulator[dataItem.traceId].services = [];
                        }
                        accumulator[dataItem.traceId].services.push(dataItem);
                    }
                    return accumulator;
                }, rootSpans);

            const searchResultsArray = [];
            for (const traceId in searchResults) {
                if (searchResults.hasOwnProperty(traceId)) {
                    searchResultsArray.push(searchResults[traceId]);
                }
            }
            self.setState((prevState) => ({
                ...prevState,
                hasSearchCompleted: true,
                searchResults: searchResultsArray
            }));
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to search for Traces",
                    StateHolder.NotificationLevels.ERROR,
                    globalState
                );
            }
        });
    };

    static getDerivedStateFromProps = (props, state) => {
        const {data, filter, metaData} = state;

        // Finding the available microservices to be selected
        const selectedCells = (filter.cell === Search.ALL_VALUE ? data.cells : [filter.cell]);
        const availableMicroservices = data.microservices
            .filter((microservice) => selectedCells.includes(microservice.cell))
            .map((microservice) => microservice.name);

        const selectedMicroservice = (filter.microservice && availableMicroservices.includes(filter.microservice))
            ? filter.microservice
            : Search.ALL_VALUE;

        // Finding the available operations to be selected
        const selectedMicroservices = (selectedMicroservice === Search.ALL_VALUE
            ? availableMicroservices
            : [selectedMicroservice]);
        const availableOperations = data.operations
            .filter((operation) => selectedMicroservices.includes(operation.microservice))
            .map((operation) => operation.name);

        const selectedOperation = (filter.operation && availableOperations.includes(filter.operation))
            ? filter.operation
            : Search.ALL_VALUE;

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
    };

}

Search.propTypes = {
    classes: PropTypes.object.isRequired,
    history: PropTypes.shape({
        replace: PropTypes.func.isRequired
    }),
    location: PropTypes.shape({
        search: PropTypes.string.isRequired
    }).isRequired,
    match: PropTypes.shape({
        url: PropTypes.string.isRequired
    }).isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired
};

export default withStyles(styles)(withGlobalState(Search));
