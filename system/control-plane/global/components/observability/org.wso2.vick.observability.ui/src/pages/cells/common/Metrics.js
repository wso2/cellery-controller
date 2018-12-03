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
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import PropTypes from "prop-types";
import React from "react";
import Select from "@material-ui/core/Select";
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

    constructor(props) {
        super(props);

        this.state = {
            type: "inbound",
            cell: "all",
            microservice: "all"
        };
    }

    handleChange = (name) => (event) => {
        this.setState({
            [name]: event.target.value
        });
    };

    render = () => {
        const {classes, isHidden} = this.props;
        return (
            <React.Fragment>
                <div className={classes.filters}>
                    <FormControl className={classes.formControl}>
                        <InputLabel htmlFor="metrics-type">Type</InputLabel>
                        <Select
                            native
                            value={this.state.type}
                            onChange={this.handleChange("type")}
                            inputProps={{
                                name: "type",
                                id: "metrics-type"
                            }}
                        >
                            <option value="inbound">Inbound</option>
                            <option value="outbound">Outbound</option>
                        </Select>
                    </FormControl>
                    <FormControl className={classes.formControl}>
                        <InputLabel htmlFor="cell">Cell</InputLabel>
                        <Select
                            native
                            value={this.state.cell}
                            onChange={this.handleChange("type")}
                            inputProps={{
                                name: "cell",
                                id: "cell"
                            }}
                        >
                            <option value="all">All</option>
                        </Select>
                    </FormControl>
                    {
                        isHidden
                            ? null
                            : (
                                <FormControl className={classes.formControl}>
                                    <InputLabel htmlFor="microservice">Microservice</InputLabel>
                                    <Select
                                        native
                                        value={this.state.age}
                                        onChange={this.handleChange("microservice")}
                                        inputProps={{
                                            name: "microservice",
                                            id: "microservice"
                                        }}
                                    >
                                        <option value="all">All</option>
                                    </Select>
                                </FormControl>
                            )
                    }
                    <Button variant="outlined" size="small" color="primary" className={classes.button}>
                        Update
                    </Button>
                </div>
                <div className={classes.graphs}>
                    Graphs
                </div>
            </React.Fragment>
        );
    };

}

Metrics.propTypes = {
    classes: PropTypes.object.isRequired,
    isHidden: PropTypes.bool.isRequired
};

export default withStyles(styles)(Metrics);
