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
import Divider from "@material-ui/core/Divider/Divider";
import FormControl from "@material-ui/core/FormControl/FormControl";
import IconButton from "@material-ui/core/IconButton/IconButton";
import InputAdornment from "@material-ui/core/InputAdornment/InputAdornment";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import PropTypes from "prop-types";
import React from "react";
import Select from "@material-ui/core/Select/Select";
import Toolbar from "@material-ui/core/Toolbar/Toolbar";
import Typography from "@material-ui/core/Typography/Typography";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core";
import {ConfigConstants, withConfig} from "./config";

const styles = (theme) => ({
    subToolbar: {
        paddingLeft: theme.spacing.unit,
        paddingRight: theme.spacing.unit
    },
    container: {
        marginBottom: theme.spacing.unit * 3
    },
    title: {
        flexGrow: 1
    },
    startInputAdornment: {
        marginBottom: theme.spacing.unit * 2
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
                refreshInterval: globalFilter.refreshInterval
            };
        } else {
            this.state = {
                startTime: "now - 24h",
                endTime: "now",
                refreshInterval: 30 * 1000
            };
            config.set(ConfigConstants.GLOBAL_FILTER, {...this.state});
        }

        this.refreshIntervalID = null;

        this.setRefreshInterval = this.setRefreshInterval.bind(this);
        this.setTimePeriod = this.setTimePeriod.bind(this);
        this.startRefreshTask = this.startRefreshTask.bind(this);
        this.stopRefreshTask = this.stopRefreshTask.bind(this);
    }

    componentDidMount() {
        this.startRefreshTask();
    }

    componentDidUpdate() {
        this.startRefreshTask();
    }

    componentWillUnmount() {
        this.stopRefreshTask();
    }

    render() {
        const {classes, title, location, history} = this.props;
        const {refreshInterval} = this.state;

        return (
            <div className={classes.container}>
                <Toolbar disableGutters={true} className={classes.subToolbar}>
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
                    <FormControl>
                        <Select value={refreshInterval} onChange={this.setRefreshInterval}
                            inputProps={{name: "refresh-interval", id: "refresh-interval"}}
                            startAdornment={(<InputAdornment className={classes.startInputAdornment}
                                variant="filled" position="start">Refresh</InputAdornment>)}>
                            <MenuItem value={-1}>Off</MenuItem>
                            <MenuItem value={5 * 1000}>Every 5 sec</MenuItem>
                            <MenuItem value={10 * 1000}>Every 10 sec</MenuItem>
                            <MenuItem value={15 * 1000}>Every 15 sec</MenuItem>
                            <MenuItem value={30 * 1000}>Every 30 sec</MenuItem>
                            <MenuItem value={60 * 1000}>Every 1 min</MenuItem>
                            <MenuItem value={5 * 60 * 1000}>Every 5 min</MenuItem>
                        </Select>
                    </FormControl>
                </Toolbar>
                <Divider/>
            </div>
        );
    }

    setTimePeriod() {
        const {config} = this.props;

        config.set(ConfigConstants.GLOBAL_FILTER, {
            ...config.get(ConfigConstants.GLOBAL_FILTER),
            startTime: 0,
            endTime: 100
        });

        this.stopRefreshTask(); // Stop any existing refresh tasks (Will be restarted when the component is updated)
        this.setState({
            startTime: config.globalFilter.startTime,
            endTime: config.globalFilter.endTime
        });
    }

    setRefreshInterval(event) {
        const {config} = this.props;
        const newRefreshInterval = event.target.value;

        config.set(ConfigConstants.GLOBAL_FILTER, {
            ...config.get(ConfigConstants.GLOBAL_FILTER),
            refreshInterval: newRefreshInterval
        });

        this.stopRefreshTask(); // Stop any existing refresh tasks (Will be restarted when the component is updated)
        this.setState({
            refreshInterval: config.get(ConfigConstants.GLOBAL_FILTER).refreshInterval
        });
    }

    startRefreshTask() {
        const {onUpdate} = this.props;
        const {startTime, endTime, refreshInterval} = this.state;

        this.stopRefreshTask(); // Stop any existing refresh tasks
        if (refreshInterval && refreshInterval > 0) {
            this.refreshIntervalID = setInterval(() => onUpdate(startTime, endTime), refreshInterval);
        }
    }

    stopRefreshTask() {
        if (this.refreshIntervalID) {
            clearInterval(this.refreshIntervalID);
            this.refreshIntervalID = null;
        }
    }

}

TopToolbar.propTypes = {
    onUpdate: PropTypes.func,
    title: PropTypes.string.isRequired,
    classes: PropTypes.any.isRequired,
    config: PropTypes.any.isRequired,
    history: PropTypes.shape({
        goBack: PropTypes.func.isRequired
    }),
    location: PropTypes.shape({
        state: PropTypes.any
    }).isRequired
};

export default withStyles(styles, {withTheme: true})(withRouter(withConfig(TopToolbar)));
