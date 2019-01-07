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

import "react-datetime/css/react-datetime.css";
import Button from "@material-ui/core/Button/Button";
import Checkbox from "@material-ui/core/Checkbox/Checkbox";
import Collapse from "@material-ui/core/Collapse";
import Constants from "../constants";
import Datetime from "react-datetime";
import Divider from "@material-ui/core/Divider";
import FormControl from "@material-ui/core/FormControl/FormControl";
import FormControlLabel from "@material-ui/core/FormControlLabel/FormControlLabel";
import Grid from "@material-ui/core/Grid/Grid";
import InputLabel from "@material-ui/core/InputLabel/InputLabel";
import PropTypes from "prop-types";
import QueryUtils from "../utils/queryUtils";
import React from "react";
import TextField from "@material-ui/core/TextField/TextField";
import Typography from "@material-ui/core/Typography/Typography";
import classNames from "classnames";
import {withStyles} from "@material-ui/core";

const styles = (theme) => ({
    dateRangePopOver: {
        width: 700,
        padding: theme.spacing.unit * 2
    },
    dateRangesTitleDivider: {
        marginBottom: theme.spacing.unit * 2
    },
    customRangeContainer: {
        padding: theme.spacing.unit * 2,
        borderRightStyle: "solid",
        borderRightColor: "#d0d0d0",
        borderRightWidth: 2
    },
    customDateRangeInputLabel: {
        paddingTop: theme.spacing.unit,
        paddingLeft: theme.spacing.unit,
        paddingBottom: theme.spacing.unit
    },
    customRangeApplyButton: {
        marginTop: theme.spacing.unit * 2
    },
    isRangeToNowCheckbox: {
        marginLeft: theme.spacing.unit * 0.5
    },
    defaultRangesContainer: {
        padding: theme.spacing.unit * 2
    },
    defaultRange: {
        cursor: "pointer",
        padding: theme.spacing.unit * 0.5
    },
    selectedDefaultRange: {
        fontWeight: 500,
        color: theme.palette.primary.main
    },
    formControl: {
        marginBottom: theme.spacing.unit * 2
    },
    collapsible: {
        marginBottom: 15
    }
});

class DateRangePicker extends React.Component {

    DEFAULT_RANGES = {
        LAST_MIN: {
            name: "Last mins",
            from: "now - 1 minute"
        },
        LAST_5_MINS: {
            name: "Last 5 mins",
            from: "now - 5 minutes"
        },
        LAST_10_MINS: {
            name: "Last 10 mins",
            from: "now - 10 minutes"
        },
        LAST_30_MINS: {
            name: "Last 30 mins",
            from: "now - 30 minutes"
        },
        LAST_1_HOUR: {
            name: "Last 1 hour",
            from: "now - 1 hour"
        },
        LAST_3_HOURS: {
            name: "Last 2 hours",
            from: "now - 2 hours"
        },
        LAST_6_HOURS: {
            name: "Last 6 hours",
            from: "now - 6 hours"
        },
        LAST_12_HOURS: {
            name: "Last 12 hours",
            from: "now - 12 hours"
        },
        LAST_24_HOURS: {
            name: "Last 24 hours",
            from: "now - 24 hours"
        },
        LAST_7_DAYS: {
            name: "Last 7 days",
            from: "now - 7 days"
        },
        LAST_30_DAYS: {
            name: "Last 30 days",
            from: "now - 30 days"
        }
    };

    DATE_RANGE_FROM = "DATE_RANGE_FROM";
    DATE_RANGE_TO = "DATE_RANGE_TO";

    constructor(props) {
        super(props);

        this.state = {
            isRangeFromCalendarOpen: false,
            isRangeToCalendarOpen: false,
            isRangeToNow: props.endTime.trim() === "now",
            startTime: props.startTime,
            endTime: props.endTime,
            dateRangeNickname: props.dateRangeNickname,
            lastFocussedDateRangeInput: undefined
        };
    }

    render = () => {
        const {classes} = this.props;
        const {
            isRangeFromCalendarOpen, isRangeToCalendarOpen, isRangeToNow, startTime, endTime, dateRangeNickname,
            lastFocussedDateRangeInput
        } = this.state;

        // Parsing the start time for the Date Time Picker
        let parsedStartTime;
        try {
            parsedStartTime = QueryUtils.parseTime(startTime);
        } catch (e) {
            parsedStartTime = undefined;
        }

        // Parsing the start time for the Date Time Picker
        let parsedEndTime;
        try {
            parsedEndTime = QueryUtils.parseTime(endTime);
        } catch (e) {
            parsedEndTime = undefined;
        }

        // Validating the provided time range
        let fromErrorMessage;
        let toErrorMessage;
        if (!parsedStartTime || !parsedEndTime) {
            if (!parsedStartTime) {
                fromErrorMessage = "Invalid time";
            }
            if (!parsedEndTime) {
                toErrorMessage = "Invalid time";
            }
        } else if (startTime.includes("now") && !endTime.includes("now")) {
            const errorMessage = "Date range cannot be selected from a relative time to a absolute time";
            if (lastFocussedDateRangeInput === this.DATE_RANGE_FROM) {
                fromErrorMessage = errorMessage;
            } else if (lastFocussedDateRangeInput === this.DATE_RANGE_TO) {
                toErrorMessage = errorMessage;
            }
        } else if (parsedStartTime.isAfter(parsedEndTime)) {
            const errorMessage = "Date range cannot be from a later time to a earlier time";
            if (lastFocussedDateRangeInput === this.DATE_RANGE_FROM) {
                fromErrorMessage = errorMessage;
            } else if (lastFocussedDateRangeInput === this.DATE_RANGE_TO) {
                toErrorMessage = errorMessage;
            }
        }

        return (
            <div className={classes.dateRangePopOver}>
                <Grid container>
                    <Grid item xs={8} className={classes.customRangeContainer}>
                        <Typography variant={"subtitle2"}>Custom Range</Typography>
                        <Divider className={classes.dateRangesTitleDivider}/>
                        <Grid container>
                            <Grid item xs={2} className={classes.customDateRangeInputLabel}>
                                <InputLabel>From</InputLabel>
                            </Grid>
                            <Grid item xs={7}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <TextField error={Boolean(fromErrorMessage)} value={startTime}
                                        onFocus={this.onFromDateRangeInputFocus}
                                        onChange={this.getCustomDateRangeInputChangeHandler("startTime")}
                                        helperText={fromErrorMessage}/>
                                </FormControl>
                                <Collapse in={isRangeFromCalendarOpen} className={classes.collapsible}>
                                    <Datetime input={false} value={parsedStartTime}
                                        onChange={this.getCustomCalendarChangeHandler("startTime")}/>
                                </Collapse>
                            </Grid>
                        </Grid>
                        <Grid container>
                            <Grid item xs={2} className={classes.customDateRangeInputLabel}>
                                <InputLabel>To</InputLabel>
                            </Grid>
                            <Grid item xs={7}>
                                <FormControl className={classes.formControl} fullWidth={true}>
                                    <TextField error={Boolean(toErrorMessage)} value={endTime}
                                        onFocus={this.onToDateRangeInputFocus}
                                        onChange={this.getCustomDateRangeInputChangeHandler("endTime")}
                                        helperText={toErrorMessage}/>
                                </FormControl>
                                <Collapse in={isRangeToCalendarOpen}>
                                    <Datetime input={false} value={parsedEndTime}
                                        onChange={this.getCustomCalendarChangeHandler("endTime")}/>
                                </Collapse>
                            </Grid>
                            <Grid item xs={3}>
                                <FormControlLabel className={classes.isRangeToNowCheckbox}
                                    control={
                                        <Checkbox
                                            checked={isRangeToNow}
                                            onChange={this.onIsRangeToNowCheckBoxClick}
                                            value="isRangeToNow"
                                            color="default"
                                        />
                                    }
                                    label="now"
                                />
                            </Grid>
                        </Grid>
                        <Grid container>
                            <Grid item xs={2}/>
                            <Grid item xs={10}>
                                <Button variant="outlined" size="small" color="primary"
                                    disabled={Boolean(fromErrorMessage || toErrorMessage)}
                                    className={classes.customRangeApplyButton} onClick={this.onCustomRangeApply}>
                                    Apply
                                </Button>
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item xs={4} className={classes.defaultRangesContainer}>
                        <Typography variant={"subtitle2"}>Default Ranges</Typography>
                        <Divider className={classes.dateRangesTitleDivider}/>
                        {
                            Object.keys(this.DEFAULT_RANGES).map((defaultRangeKey) => {
                                const dateRangeName = this.DEFAULT_RANGES[defaultRangeKey].name;
                                return (
                                    <Typography key={defaultRangeKey}
                                        className={classNames({
                                            [classes.defaultRange]: true,
                                            [classes.selectedDefaultRange]: dateRangeName === dateRangeNickname
                                        })}
                                        onClick={this.getDefaultRangeClickEventHandler(defaultRangeKey)}>
                                        {dateRangeName}
                                    </Typography>
                                );
                            })
                        }
                    </Grid>
                </Grid>
            </div>
        );
    };

    /**
     * Get a change handler for a date range input.
     *
     * @param {string} type One of startTime or endTime
     * @returns {Function} The range change handler
     */
    getCustomDateRangeInputChangeHandler = (type) => {
        const self = this;
        return (event) => {
            const newDate = event.target.value;
            self.setState((prevState) => ({
                [type]: newDate,
                isRangeToNow: type === "endTime" ? newDate.trim() === "now" : prevState.isRangeToNow
            }));
        };
    };

    /**
     * Get a change handler for a calendar.
     *
     * @param {string} type One of startTime or endTime
     * @returns {Function} The calendar date change handler
     */
    getCustomCalendarChangeHandler = (type) => {
        const self = this;
        return (date) => {
            self.setState((prevState) => ({
                [type]: date.format(Constants.Pattern.DATE_TIME),
                isRangeToNow: type === "endTime" ? false : prevState.isRangeToNow
            }));
        };
    };

    /**
     * Get an event handler for default range select.
     *
     * @param {string} defaultRange The selected default range
     * @returns {Function} The event handler for the default range select
     */
    getDefaultRangeClickEventHandler = (defaultRange) => {
        const self = this;
        return () => {
            const startTime = this.DEFAULT_RANGES[defaultRange].from;
            const endTime = "now";
            const dateRangeNickname = this.DEFAULT_RANGES[defaultRange].name;

            self.setState({
                isRangeToNow: true,
                isRangeFromCalendarOpen: false,
                isRangeToCalendarOpen: false,
                startTime: startTime,
                endTime: endTime,
                dateRangeNickname: dateRangeNickname,
                lastFocussedDateRangeInput: undefined
            });
            this.applyDateRange(startTime, endTime, dateRangeNickname);
        };
    };

    /**
     * Handle the custom range apply button click.
     */
    onCustomRangeApply = () => {
        this.applyDateRange();
    };

    /**
     * Handle the user clicking on the range is to now checkbox.
     *
     * @param {SyntheticEvent} event The checkbox click event
     */
    onIsRangeToNowCheckBoxClick = (event) => {
        const isChecked = event.target.checked;
        this.setState((prevState) => ({
            isRangeToNow: isChecked,
            endTime: isChecked ? "now" : prevState.endTime
        }));
    };

    /**
     * Open the date range from calendar.
     */
    onFromDateRangeInputFocus = () => {
        this.setState({
            isRangeFromCalendarOpen: true,
            isRangeToCalendarOpen: false,
            lastFocussedDateRangeInput: this.DATE_RANGE_FROM
        });
    };

    /**
     * Open the date range to calendar.
     */
    onToDateRangeInputFocus = () => {
        this.setState({
            isRangeFromCalendarOpen: false,
            isRangeToCalendarOpen: true,
            lastFocussedDateRangeInput: this.DATE_RANGE_TO
        });
    };

    /**
     * Handle on date range apply event.
     *
     * @param {string} newStartTime The new start time to be applied
     * @param {string} newEndTime The new end time to be applied
     * @param {string} newDateRangeNickname The new date range nickname to be applied
     */
    applyDateRange = (newStartTime = "", newEndTime = "", newDateRangeNickname = "") => {
        const {onRangeChange} = this.props;
        const {startTime, endTime} = this.state;

        onRangeChange(
            newStartTime ? newStartTime : startTime,
            newEndTime ? newEndTime : endTime,
            newDateRangeNickname ? newDateRangeNickname : null
        );
    };

}

DateRangePicker.propTypes = {
    classes: PropTypes.any.isRequired,
    startTime: PropTypes.string.isRequired,
    endTime: PropTypes.string.isRequired,
    dateRangeNickname: PropTypes.string,
    onRangeChange: PropTypes.func.isRequired
};

export default withStyles(styles, {withTheme: true})(DateRangePicker);
