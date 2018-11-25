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
import FormControl from "@material-ui/core/FormControl/FormControl";
import Grid from "@material-ui/core/Grid/Grid";
import HttpUtils from "../common/utils/httpUtils";
import InputAdornment from "@material-ui/core/InputAdornment/InputAdornment";
import InputLabel from "@material-ui/core/InputLabel/InputLabel";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import PropTypes from "prop-types";
import React from "react";
import Select from "@material-ui/core/Select/Select";
import TextField from "@material-ui/core/TextField/TextField";
import TopToolbar from "../common/TopToolbar";
import Typography from "@material-ui/core/Typography/Typography";
import {withRouter} from "react-router-dom/";
import withStyles from "@material-ui/core/styles/withStyles";

const styles = (theme) => ({
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
        marginBottom: theme.spacing.unit * 2
    }
});

const TraceResult = () => (
    <div>Result</div>
);

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
            traces: []
        });

        this.loadCellData = this.loadCellData.bind(this);
        this.getChangeHandler = this.getChangeHandler.bind(this);
        this.search = this.search.bind(this);
        this.loadTracePage = this.loadTracePage.bind(this);
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
        const {data, filter, metaData, traces} = this.state;

        const createMenuItemForSelect = (itemNames) => itemNames.map(
            (itemName) => (<MenuItem key={itemName} value={itemName}>{itemName}</MenuItem>)
        );

        return (
            <React.Fragment>
                <TopToolbar title={"Distributed Tracing"} onUpdate={this.loadCellData}/>
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
                {
                    traces.length > 0
                        ? traces.map(
                            (trace) => <TraceResult key={trace.traceId} trace={trace} onClick={this.loadTracePage}/>)
                        : null
                }
            </React.Fragment>
        );
    }

    /**
     * Search for traces.
     * Call the backend and search for traces.
     */
    loadCellData() {
        // TODO : Load values from the server
        const self = this;
        setTimeout(() => {
            self.setState(Search.generateValidState({
                ...self.state,
                data: {
                    cells: ["cellA", "cellB", "cellC"],
                    microservices: [
                        {name: "cellA-microserviceA", cell: "cellA"},
                        {name: "cellA-microserviceB", cell: "cellA"},
                        {name: "cellA-microserviceC", cell: "cellA"},
                        {name: "cellB-microserviceA", cell: "cellB"},
                        {name: "cellB-microserviceB", cell: "cellB"},
                        {name: "cellB-microserviceC", cell: "cellB"},
                        {name: "cellC-microserviceA", cell: "cellC"},
                        {name: "cellC-microserviceB", cell: "cellC"},
                        {name: "cellC-microserviceC", cell: "cellC"}
                    ],
                    operations: [
                        {name: "cellA-microserviceA-operationA", microservice: "cellA-microserviceA", cell: "cellA"},
                        {name: "cellA-microserviceA-operationB", microservice: "cellA-microserviceA", cell: "cellA"},
                        {name: "cellA-microserviceA-operationC", microservice: "cellA-microserviceA", cell: "cellA"},
                        {name: "cellA-microserviceB-operationA", microservice: "cellA-microserviceB", cell: "cellA"},
                        {name: "cellA-microserviceB-operationB", microservice: "cellA-microserviceB", cell: "cellA"},
                        {name: "cellA-microserviceB-operationC", microservice: "cellA-microserviceB", cell: "cellA"},
                        {name: "cellA-microserviceC-operationA", microservice: "cellA-microserviceC", cell: "cellA"},
                        {name: "cellA-microserviceC-operationB", microservice: "cellA-microserviceC", cell: "cellA"},
                        {name: "cellA-microserviceC-operationC", microservice: "cellA-microserviceC", cell: "cellA"},
                        {name: "cellB-microserviceA-operationA", microservice: "cellB-microserviceA", cell: "cellB"},
                        {name: "cellB-microserviceA-operationB", microservice: "cellB-microserviceA", cell: "cellB"},
                        {name: "cellB-microserviceA-operationC", microservice: "cellB-microserviceA", cell: "cellB"},
                        {name: "cellB-microserviceB-operationA", microservice: "cellB-microserviceB", cell: "cellB"},
                        {name: "cellB-microserviceB-operationB", microservice: "cellB-microserviceB", cell: "cellB"},
                        {name: "cellB-microserviceB-operationC", microservice: "cellB-microserviceB", cell: "cellB"},
                        {name: "cellB-microserviceC-operationA", microservice: "cellB-microserviceC", cell: "cellB"},
                        {name: "cellB-microserviceC-operationB", microservice: "cellB-microserviceC", cell: "cellB"},
                        {name: "cellB-microserviceC-operationC", microservice: "cellB-microserviceC", cell: "cellB"},
                        {name: "cellC-microserviceA-operationA", microservice: "cellC-microserviceA", cell: "cellC"},
                        {name: "cellC-microserviceA-operationB", microservice: "cellC-microserviceA", cell: "cellC"},
                        {name: "cellC-microserviceA-operationC", microservice: "cellC-microserviceA", cell: "cellC"},
                        {name: "cellC-microserviceB-operationA", microservice: "cellC-microserviceB", cell: "cellC"},
                        {name: "cellC-microserviceB-operationB", microservice: "cellC-microserviceB", cell: "cellC"},
                        {name: "cellC-microserviceB-operationC", microservice: "cellC-microserviceB", cell: "cellC"},
                        {name: "cellC-microserviceC-operationA", microservice: "cellC-microserviceC", cell: "cellC"},
                        {name: "cellC-microserviceC-operationB", microservice: "cellC-microserviceC", cell: "cellC"},
                        {name: "cellC-microserviceC-operationC", microservice: "cellC-microserviceC", cell: "cellC"}
                    ]
                }
            }));
        }, 1000);
    }

    /**
     * Get the on change handler for a particular state filter attribute.
     *
     * @param {string} name The name of the filter
     * @returns {Function} The on change handler
     */
    getChangeHandler(name) {
        return (event) => {
            const newState = Search.generateValidState({
                ...this.state,
                filter: {
                    ...this.state.filter,
                    [name]: event.target.value
                }
            });
            this.setState(newState);
        };
    }

    search() {
        const {
            cell, microservice, operation, tags, minDuration, minDurationMultiplier, maxDuration, maxDurationMultiplier
        } = this.state.filter;

        // Build search object
        const search = {};
        const addSearchParam = (key, value) => {
            if (value) {
                search[key] = value;
            }
        };
        addSearchParam("cell", cell);
        addSearchParam("microservice", microservice);
        addSearchParam("operation", operation);
        addSearchParam("tags", tags);
        addSearchParam("minDuration", minDuration * minDurationMultiplier);
        addSearchParam("maxDuration", maxDuration * maxDurationMultiplier);

        // TODO : Search in the backend
    }

    /**
     * Load the trace page.
     *
     * @param {string} traceId The trace ID of the selected trace
     * @param {string} microservice The microservice name if a microservice was selected
     */
    loadTracePage(traceId, microservice) {
        this.props.history.push({
            pathname: `./id/${traceId}`,
            state: {
                highlightedMicroservice: microservice
            }
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
    history: PropTypes.shape({
        push: PropTypes.func.isRequired
    }).isRequired,
    location: PropTypes.shape({
        search: PropTypes.string.isRequired
    }).isRequired,
    match: PropTypes.shape({
        url: PropTypes.string.isRequired
    }).isRequired
};

export default withStyles(styles)(withRouter(Search));
