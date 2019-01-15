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
import Constants from "../../../utils/constants";
import DurationInput from "./DurationInput";
import FormControl from "@material-ui/core/FormControl/FormControl";
import Grid from "@material-ui/core/Grid/Grid";
import HttpUtils from "../../../utils/api/httpUtils";
import InputLabel from "@material-ui/core/InputLabel/InputLabel";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import NotFound from "../../common/error/NotFound";
import NotificationUtils from "../../../utils/common/notificationUtils";
import Paper from "@material-ui/core/Paper/Paper";
import React from "react";
import Select from "@material-ui/core/Select/Select";
import Span from "../../../utils/tracing/span";
import TagsInput from "./TagsInput";
import TopToolbar from "../../common/toptoolbar";
import TracesList from "./TracesList";
import Typography from "@material-ui/core/Typography/Typography";
import withStyles from "@material-ui/core/styles/withStyles";
import withGlobalState, {StateHolder} from "../../common/state";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    container: {
        flexGrow: 1,
        padding: theme.spacing.unit * 3,
        margin: theme.spacing.unit
    },
    subheading: {
        marginBottom: theme.spacing.unit * 2
    },
    formControl: {
        marginBottom: theme.spacing.unit * 2
    },
    searchForm: {
        marginBottom: Number(theme.spacing.unit)
    },
    resultContainer: {
        marginTop: theme.spacing.unit * 3
    }
});

class TraceSearch extends React.Component {

    constructor(props) {
        super(props);
        const {location} = props;

        const queryParams = HttpUtils.parseQueryParams(location.search);
        this.state = {
            data: {
                cells: [],
                components: [],
                operations: []
            },
            filter: {
                cell: queryParams.cell ? queryParams.cell : Constants.Dashboard.ALL_VALUE,
                component: queryParams.component ? queryParams.component : Constants.Dashboard.ALL_VALUE,
                operation: queryParams.operation ? queryParams.operation : Constants.Dashboard.ALL_VALUE,
                tags: queryParams.tags ? JSON.parse(queryParams.tags) : {},
                minDuration: queryParams.minDuration
                    ? parseInt(queryParams.minDuration, 10)
                    : undefined,
                minDurationMultiplier: queryParams.minDurationMultiplier
                    ? parseInt(queryParams.minDurationMultiplier, 10)
                    : 1,
                maxDuration: queryParams.maxDuration
                    ? parseInt(queryParams.maxDuration, 10)
                    : undefined,
                maxDurationMultiplier: queryParams.maxDurationMultiplier
                    ? parseInt(queryParams.maxDurationMultiplier, 10)
                    : 1
            },
            metaData: {
                availableComponents: [],
                availableOperations: []
            },
            isLoading: false,
            hasSearchCompleted: false
        };

        this.tracesListRef = React.createRef();
    }

    render = () => {
        const {classes, location} = this.props;
        const {data, filter, metaData, isLoading} = this.state;

        /*
         * Checking if the search should be run for the right after rendering
         * If the query params are present, it indicates that the search should be run.
         */
        const queryParams = HttpUtils.parseQueryParams(location.search);
        let isQueryParamsEmpty = true;
        for (const key in queryParams) {
            if (queryParams.hasOwnProperty(key) && queryParams[key]) {
                isQueryParamsEmpty = false;
            }
        }

        const createMenuItemsForSelect = (itemNames) => itemNames.map(
            (itemName) => (<MenuItem key={itemName} value={itemName}>{itemName}</MenuItem>)
        );

        return (
            <React.Fragment>
                <TopToolbar title={"Distributed Tracing"} onUpdate={this.onGlobalRefresh}/>
                {
                    isLoading
                        ? null
                        : (
                            <Paper className={classes.container}>
                                <Typography variant="h6" color="inherit" className={classes.subheading}>
                                    Search Traces
                                </Typography>
                                <Grid container justify={"flex-start"} className={classes.searchForm}>
                                    <Grid container justify={"flex-start"} spacing={24}>
                                        <Grid item xs={3}>
                                            <FormControl className={classes.formControl} fullWidth={true}>
                                                <InputLabel htmlFor="cell" shrink={true}>Cell</InputLabel>
                                                <Select value={filter.cell} inputProps={{name: "cell", id: "cell"}}
                                                    onChange={this.getChangeHandlerForString("cell")}>
                                                    <MenuItem key={Constants.Dashboard.ALL_VALUE}
                                                        value={Constants.Dashboard.ALL_VALUE}>
                                                        {Constants.Dashboard.ALL_VALUE}
                                                    </MenuItem>
                                                    {createMenuItemsForSelect(data.cells)}
                                                </Select>
                                            </FormControl>
                                        </Grid>
                                        <Grid item xs={3}>
                                            <FormControl className={classes.formControl} fullWidth={true}>
                                                <InputLabel htmlFor="component" shrink={true}>Component</InputLabel>
                                                <Select value={filter.component}
                                                    onChange={this.getChangeHandlerForString("component")}
                                                    inputProps={{name: "component", id: "component"}}>
                                                    <MenuItem key={Constants.Dashboard.ALL_VALUE}
                                                        value={Constants.Dashboard.ALL_VALUE}>
                                                        {Constants.Dashboard.ALL_VALUE}
                                                    </MenuItem>
                                                    {createMenuItemsForSelect(metaData.availableComponents)}
                                                </Select>
                                            </FormControl>
                                        </Grid>
                                        <Grid item xs={3}>
                                            <FormControl className={classes.formControl} fullWidth={true}>
                                                <InputLabel htmlFor="operation" shrink={true}>Operation</InputLabel>
                                                <Select value={filter.operation}
                                                    onChange={this.getChangeHandlerForString("operation")}
                                                    inputProps={{name: "operation", id: "operation"}}>
                                                    <MenuItem key={Constants.Dashboard.ALL_VALUE}
                                                        value={Constants.Dashboard.ALL_VALUE}>
                                                        {Constants.Dashboard.ALL_VALUE}
                                                    </MenuItem>
                                                    {createMenuItemsForSelect(metaData.availableOperations)}
                                                </Select>
                                            </FormControl>
                                        </Grid>
                                    </Grid>
                                </Grid>
                                <Grid container justify={"flex-start"} spacing={24} className={classes.searchForm}>
                                    <Grid item xs={6}>
                                        <FormControl className={classes.formControl} fullWidth={true}>
                                            <TagsInput onTagsUpdate={this.handleTagsUpdate} defaultTags={filter.tags}/>
                                        </FormControl>
                                    </Grid>
                                    <Grid item xs={3}>
                                        <FormControl className={classes.formControl} fullWidth={true}>
                                            <InputLabel htmlFor="min-duration-input" shrink={true}>Duration</InputLabel>
                                            <DurationInput onDurationUpdate={this.handleMinDurationUpdate} label={"Min"}
                                                durationInputId={"min-duration-input"}
                                                defaultDuration={filter.minDuration}
                                                defaultDurationMultiplier={filter.minDurationMultiplier}/>
                                        </FormControl>
                                    </Grid>
                                    <Grid item xs={3}>
                                        <FormControl className={classes.formControl} fullWidth={true}>
                                            <DurationInput onDurationUpdate={this.handleMaxDurationUpdate} label={"Max"}
                                                durationInputId={"max-duration-input"}
                                                defaultDuration={filter.maxDuration}
                                                defaultDurationMultiplier={filter.maxDurationMultiplier}/>
                                        </FormControl>
                                    </Grid>
                                </Grid>
                                <Button variant="contained" color="primary" onClick={this.onSearchButtonClick}
                                    disabled={data.cells.length === 0}>
                                    Search
                                </Button>
                                {
                                    data.cells.length > 0
                                        ? (
                                            <div className={classes.resultContainer}>
                                                <TracesList innerRef={this.tracesListRef} filter={filter}
                                                    onTraceClick={this.onTraceClick}
                                                    loadTracesOnMount={!isQueryParamsEmpty}/>
                                            </div>
                                        )
                                        : (
                                            <NotFound title={"No Traces Available"}
                                                description={"No Traces are available in the Selected Time Range. "
                                                    + "This is because no requests were sent/received during this "
                                                    + "time period."}/>
                                        )
                                }
                            </Paper>
                        )
                }
            </React.Fragment>
        );
    };

    onTraceClick = (traceId, selectedCellName, selectedComponent) => {
        this.props.history.push({
            pathname: `./id/${traceId}`,
            state: {
                selectedComponent: {
                    cellName: selectedCellName,
                    serviceName: selectedComponent
                }
            }
        });
    };

    onSearchButtonClick = () => {
        const {history, match, location} = this.props;
        const {filter} = this.state;

        // Updating the URL to ensure that the user can come back to this page
        const searchString = HttpUtils.generateQueryParamString({
            ...HttpUtils.parseQueryParams(location.search),
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
            self.setState({
                isLoading: true
            });
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/traces/metadata${HttpUtils.generateQueryParamString(filter)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const cells = [];
            const components = [];
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
                    if (!components.map((service) => service.name).includes(serviceName)) {
                        components.push({
                            name: serviceName,
                            cell: cellName
                        });
                    }
                    if (!operations.map((operation) => operation.name).includes(operationName)) {
                        operations.push({
                            name: operationName,
                            component: serviceName,
                            cell: cellName
                        });
                    }
                }
            }

            self.setState((prevState) => ({
                ...prevState,
                data: {
                    cells: cells,
                    components: components,
                    operations: operations
                }
            }));
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
                    "Failed to load Cell Data",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    /**
     * Get the on change handler for a particular state filter attribute of type string.
     *
     * @param {string} name The name of the filter
     * @returns {Function} The on change handler
     */
    getChangeHandlerForString = (name) => (event) => {
        const value = event.target.value;
        this.setState((prevState) => ({
            ...prevState,
            filter: {
                ...prevState.filter,
                [name]: value
            }
        }));
    };

    handleMinDurationUpdate = ({duration, durationMultiplier}) => {
        this.setState((prevState) => ({
            filter: {
                ...prevState.filter,
                minDuration: duration,
                minDurationMultiplier: durationMultiplier
            }
        }));
    };

    handleMaxDurationUpdate = ({duration, durationMultiplier}) => {
        this.setState((prevState) => ({
            filter: {
                ...prevState.filter,
                maxDuration: duration,
                maxDurationMultiplier: durationMultiplier
            }
        }));
    };

    /**
     * Handle a tags object changes.
     *
     * @param {Object} newTags The new tags object
     */
    handleTagsUpdate = (newTags) => {
        this.setState((prevState) => ({
            filter: {
                ...prevState.filter,
                tags: newTags
            }
        }));
    };

    search = (isUserAction) => {
        if (this.tracesListRef.current && this.tracesListRef.current.loadTraces) {
            this.tracesListRef.current.loadTraces(isUserAction);
        }
    };

    static getDerivedStateFromProps = (props, state) => {
        const {data, filter, metaData} = state;

        // Finding the available components to be selected
        const selectedCells = (filter.cell === Constants.Dashboard.ALL_VALUE ? data.cells : [filter.cell]);
        const availableComponents = data.components
            .filter((component) => selectedCells.includes(component.cell))
            .map((component) => component.name);

        const selectedComponent = data.cells.length === 0 || (filter.component
            && availableComponents.includes(filter.component))
            ? filter.component
            : Constants.Dashboard.ALL_VALUE;

        // Finding the available operations to be selected
        const selectedComponents = (selectedComponent === Constants.Dashboard.ALL_VALUE
            ? availableComponents
            : [selectedComponent]);
        const availableOperations = data.operations
            .filter((operation) => selectedComponents.includes(operation.component))
            .map((operation) => operation.name);

        const selectedOperation = data.cells.length === 0 || (filter.operation
            && availableOperations.includes(filter.operation))
            ? filter.operation
            : Constants.Dashboard.ALL_VALUE;

        return {
            ...state,
            filter: {
                ...filter,
                component: selectedComponent,
                operation: selectedOperation
            },
            metaData: {
                ...metaData,
                availableComponents: availableComponents,
                availableOperations: availableOperations
            }
        };
    };

}

TraceSearch.propTypes = {
    classes: PropTypes.object.isRequired,
    history: PropTypes.shape({
        replace: PropTypes.func.isRequired,
        push: PropTypes.func.isRequired
    }).isRequired,
    location: PropTypes.shape({
        search: PropTypes.string.isRequired
    }).isRequired,
    match: PropTypes.shape({
        url: PropTypes.string.isRequired
    }).isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired
};

export default withStyles(styles)(withGlobalState(TraceSearch));
