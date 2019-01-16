/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import InputAdornment from "@material-ui/core/InputAdornment/InputAdornment";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import React from "react";
import Select from "@material-ui/core/Select/Select";
import TextField from "@material-ui/core/TextField/TextField";
import withStyles from "@material-ui/core/styles/withStyles";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    durationTextField: {
        marginTop: theme.spacing.unit * 2
    },
    startInputAdornment: {
        marginRight: theme.spacing.unit * 2,
        marginBottom: theme.spacing.unit * 2
    }
});

class DurationInput extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            duration: props.defaultDuration,
            durationMultiplier: props.defaultDurationMultiplier
        };
    }

    render() {
        const {classes, label, durationInputId, durationMultiplierInputId} = this.props;
        const {duration, durationMultiplier} = this.state;

        return (
            <TextField id={durationInputId} value={duration ? duration : ""} type={"number"}
                className={classes.durationTextField} placeholder={"Eg: 1,000"}
                onChange={this.getChangeHandlerForNumber("duration")}
                InputProps={{
                    startAdornment: (
                        <InputAdornment className={classes.startInputAdornment} variant="filled" position="start">
                            {label}
                        </InputAdornment>
                    ),
                    endAdornment: (
                        <InputAdornment variant="filled" position="end">
                            <Select value={durationMultiplier}
                                onChange={this.getChangeHandlerForNumber("durationMultiplier")}
                                inputProps={{
                                    name: durationMultiplierInputId,
                                    id: durationMultiplierInputId
                                }}>
                                <MenuItem value={1}>ms</MenuItem>
                                <MenuItem value={1000}>s</MenuItem>
                            </Select>
                        </InputAdornment>
                    )
                }}/>
        );
    }

    /**
     * Get the on change handler for a particular state filter attribute of type number.
     *
     * @param {string} name The name of the filter
     * @returns {Function} The on change handler
     */
    getChangeHandlerForNumber = (name) => (event) => {
        const {onDurationUpdate} = this.props;
        const value = event.target.value === "" ? undefined : parseFloat(event.target.value);
        if (value === undefined || !isNaN(value)) {
            this.setState((prevState) => {
                const newState = {
                    ...prevState,
                    [name]: value
                };
                if (onDurationUpdate) {
                    onDurationUpdate(newState);
                }
                return newState;
            });
        }
    };

}

DurationInput.propTypes = {
    classes: PropTypes.object.isRequired,
    durationInputId: PropTypes.string,
    durationMultiplierInputId: PropTypes.string,
    label: PropTypes.string.isRequired,
    defaultDuration: PropTypes.number,
    defaultDurationMultiplier: PropTypes.number.isRequired,
    onDurationUpdate: PropTypes.func.isRequired
};

export default withStyles(styles)(DurationInput);
