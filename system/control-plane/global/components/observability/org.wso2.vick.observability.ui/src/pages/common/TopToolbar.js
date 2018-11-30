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

import ArrowBack from "@material-ui/icons/ArrowBack";
import ArrowDropDown from "@material-ui/icons/ArrowDropDown";
import Button from "@material-ui/core/Button";
import CalendarToday from "@material-ui/icons/CalendarToday";
import {ConfigHolder} from "./config/configHolder";
import DateRangePicker from "./DateRangePicker";
import FormControl from "@material-ui/core/FormControl/FormControl";
import IconButton from "@material-ui/core/IconButton/IconButton";
import InputAdornment from "@material-ui/core/InputAdornment/InputAdornment";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import Popover from "@material-ui/core/Popover";
import PropTypes from "prop-types";
import QueryUtils from "./utils/queryUtils";
import React from "react";
import Refresh from "@material-ui/icons/Refresh";
import Select from "@material-ui/core/Select/Select";
import Toolbar from "@material-ui/core/Toolbar/Toolbar";
import Typography from "@material-ui/core/Typography/Typography";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core";
import {ConfigConstants, withConfig} from "./config";

const styles = (theme) => ({
    container: {
        position: "sticky",
        top: 64,
        backgroundColor: "#fafafa",
        marginBottom: theme.spacing.unit * 6,
        zIndex: 999,
        minHeight: 70
    },
    title: {
        flexGrow: 1,
        marginLeft: theme.spacing.unit,
        marginTop: theme.spacing.unit
    },
    dateRangeButton: {
        marginRight: theme.spacing.unit * 3
    },
    startInputAdornment: {
        marginRight: theme.spacing.unit * 2,
        marginBottom: theme.spacing.unit * 2
    },
    menuButton: {
        marginTop: theme.spacing.unit
    },
    dateRangeNicknameSelectedTime: {
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit
    }
});

/**
 * Top toolbar providing the global filter and page name.
 */
class TopToolbar extends React.Component {

    constructor(props) {
        super(props);
        const {config} = this.props;

        const globalFilter = config.get(ConfigConstants.GLOBAL_FILTER);
        if (globalFilter) {
            this.state = {
                startTime: globalFilter.startTime,
                endTime: globalFilter.endTime,
                dateRangeNickname: globalFilter.dateRangeNickname,
                refreshInterval: globalFilter.refreshInterval,
                dateRangeSelectorAnchorElement: undefined
            };
        } else {
            this.state = {
                startTime: "now - 24 hours",
                endTime: "now",
                dateRangeNickname: globalFilter.dateRangeNickname,
                refreshInterval: 30 * 1000,
                dateRangeSelectorAnchorElement: undefined
            };
            config.set(ConfigConstants.GLOBAL_FILTER, {...this.state});
        }

        this.refreshIntervalID = null;

        this.refreshManually = this.refreshManually.bind(this);
        this.setRefreshInterval = this.setRefreshInterval.bind(this);
        this.setTimePeriod = this.setTimePeriod.bind(this);
        this.startRefreshTask = this.startRefreshTask.bind(this);
        this.stopRefreshTask = this.stopRefreshTask.bind(this);
        this.refresh = this.refresh.bind(this);
        this.openDateRangeSelector = this.openDateRangeSelector.bind(this);
        this.closeDateRangeSelector = this.closeDateRangeSelector.bind(this);
    }

    componentDidMount() {
        this.refreshManually();
    }

    componentDidUpdate() {
        this.startRefreshTask();
    }

    componentWillUnmount() {
        this.stopRefreshTask();
    }

    render() {
        const {classes, title, location, history} = this.props;
        const {startTime, endTime, dateRangeNickname, refreshInterval, dateRangeSelectorAnchorElement} = this.state;

        const isDateRangeSelectorOpen = Boolean(dateRangeSelectorAnchorElement);
        return (
            <div className={classes.container}>
                <Toolbar disableGutters={true}>
                    {
                        location.state && location.state.hideBackButton
                            ? null
                            : (
                                <IconButton className={classes.menuButton} color="inherit" aria-label="Back"
                                    onClick={() => history.goBack()}>
                                    <ArrowBack/>
                                </IconButton>
                            )
                    }
                    <Typography variant="h5" color="inherit" className={classes.title}>
                        {title}
                    </Typography>
                    <Button aria-owns={isDateRangeSelectorOpen ? "date-range-picker-popper" : undefined}
                        className={classes.dateRangeButton} aria-haspopup="true" size="small" variant="text"
                        onClick={(event) => this.openDateRangeSelector(event.currentTarget)}>
                        {dateRangeNickname
                            ? dateRangeNickname
                            : (
                                <React.Fragment>
                                    <Typography color={"textSecondary"}>From</Typography>
                                    <Typography className={classes.dateRangeNicknameSelectedTime}>
                                        {startTime}
                                    </Typography>
                                    <Typography color={"textSecondary"}>To</Typography>
                                    <Typography className={classes.dateRangeNicknameSelectedTime}>
                                        {endTime}
                                    </Typography>
                                </React.Fragment>
                            )
                        }<ArrowDropDown/><CalendarToday/>
                    </Button>
                    <Popover id="date-range-picker-popper"
                        open={isDateRangeSelectorOpen}
                        anchorEl={dateRangeSelectorAnchorElement}
                        onClose={this.closeDateRangeSelector}
                        anchorOrigin={{
                            vertical: "bottom",
                            horizontal: "right"
                        }}
                        transformOrigin={{
                            vertical: "top",
                            horizontal: "right"
                        }}>
                        <DateRangePicker startTime={startTime} endTime={endTime} dateRangeNickname={dateRangeNickname}
                            onRangeChange={this.setTimePeriod}/>
                    </Popover>
                    <FormControl>
                        <Select value={refreshInterval} onChange={this.setRefreshInterval}
                            inputProps={{name: "refresh-interval", id: "refresh-interval"}}
                            startAdornment={(<InputAdornment className={classes.startInputAdornment}
                                variant="filled"
                                position="start">Refresh</InputAdornment>)}>
                            <MenuItem value={-1}>Off</MenuItem>
                            <MenuItem value={5 * 1000}>Every 5 sec</MenuItem>
                            <MenuItem value={10 * 1000}>Every 10 sec</MenuItem>
                            <MenuItem value={15 * 1000}>Every 15 sec</MenuItem>
                            <MenuItem value={30 * 1000}>Every 30 sec</MenuItem>
                            <MenuItem value={60 * 1000}>Every 1 min</MenuItem>
                            <MenuItem value={5 * 60 * 1000}>Every 5 min</MenuItem>
                        </Select>
                    </FormControl>
                    <IconButton aria-label="Refresh" onClick={this.refreshManually}>
                        <Refresh/>
                    </IconButton>
                </Toolbar>
            </div>
        );
    }

    /**
     * Open the date range selector with the provided element as the anchor.
     *
     * @param {HTMLElement} element The element which should act as the anchor of the date range popover
     */
    openDateRangeSelector(element) {
        this.setState({
            dateRangeSelectorAnchorElement: element
        });
    }

    /**
     * Close the date range selector currently open.
     */
    closeDateRangeSelector() {
        this.setState({
            dateRangeSelectorAnchorElement: undefined
        });
    }

    /**
     * Set the time period which should be considered for fetching data for the application.
     *
     * @param {string} startTime The new start time to be set
     * @param {string} endTime The new end time to be set
     * @param {string} dateRangeNickname The new date range nickname to be set
     */
    setTimePeriod(startTime, endTime, dateRangeNickname) {
        const {config} = this.props;

        config.set(ConfigConstants.GLOBAL_FILTER, {
            ...config.get(ConfigConstants.GLOBAL_FILTER),
            startTime: startTime,
            endTime: endTime,
            dateRangeNickname: dateRangeNickname
        });

        this.stopRefreshTask(); // Stop any existing refresh tasks (Will be restarted when the component is updated)
        this.setState({
            startTime: startTime,
            endTime: endTime,
            dateRangeNickname: dateRangeNickname
        });
        this.closeDateRangeSelector();
    }

    /**
     * Set the refresh interval to be used by the application.
     *
     * @param {Event} event The change event of the select input for the refresh interval
     */
    setRefreshInterval(event) {
        const {config} = this.props;
        const newRefreshInterval = event.target.value;

        config.set(ConfigConstants.GLOBAL_FILTER, {
            ...config.get(ConfigConstants.GLOBAL_FILTER),
            refreshInterval: newRefreshInterval
        });

        this.stopRefreshTask(); // Stop any existing refresh tasks (Will be restarted when the component is updated)
        this.setState({
            refreshInterval: newRefreshInterval
        });
    }

    /**
     * Call the update method and refresh the relevant parts of the page.
     * This will also reset the refresh task.
     */
    refreshManually() {
        this.stopRefreshTask();
        this.refresh(true);
        this.startRefreshTask();
    }

    /**
     * Start the refresh task which periodically calls the update method.
     */
    startRefreshTask() {
        const {refreshInterval} = this.state;
        const self = this;

        this.stopRefreshTask(); // Stop any existing refresh tasks
        if (refreshInterval && refreshInterval > 0) {
            this.refreshIntervalID = setInterval(() => self.refresh(false), refreshInterval);
        }
    }

    /**
     * Stop the refresh task which was started by #startRefreshTask().
     */
    stopRefreshTask() {
        if (this.refreshIntervalID) {
            clearInterval(this.refreshIntervalID);
            this.refreshIntervalID = null;
        }
    }

    /**
     * Refresh the components by calling the on Update.
     *
     * @param {boolean} isUserAction True if the refresh was initiated by user
     */
    refresh(isUserAction) {
        const {onUpdate} = this.props;
        const {startTime, endTime} = this.state;

        if (onUpdate) {
            onUpdate(
                isUserAction,
                QueryUtils.parseTime(startTime),
                QueryUtils.parseTime(endTime)
            );
        }
    }

}

TopToolbar.propTypes = {
    onUpdate: PropTypes.func,
    title: PropTypes.string.isRequired,
    classes: PropTypes.any.isRequired,
    config: PropTypes.instanceOf(ConfigHolder).isRequired,
    history: PropTypes.shape({
        goBack: PropTypes.func.isRequired
    }),
    location: PropTypes.shape({
        state: PropTypes.any
    }).isRequired
};

export default withStyles(styles, {withTheme: true})(withRouter(withConfig(TopToolbar)));
