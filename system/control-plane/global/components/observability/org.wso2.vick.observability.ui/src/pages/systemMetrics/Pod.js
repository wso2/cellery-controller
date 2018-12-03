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
import Metrics from "./common/Metrics";
import Paper from "@material-ui/core/Paper";
import PropTypes from "prop-types";
import React from "react";
import Select from "@material-ui/core/Select";
import TopToolbar from "../common/toptoolbar";
import {withStyles} from "@material-ui/core/styles";

const styles = (theme) => ({
    root: {
        padding: theme.spacing.unit * 3,
        margin: theme.spacing.unit
    },
    filters: {
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

class Pod extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            selectedPod: "all"
        };
    }

    handleChange = (name) => (event) => {
        this.setState({
            [name]: event.target.value
        });
    };

    render = () => {
        const {classes} = this.props;
        const {selectedPod} = this.state;

        return (
            <React.Fragment>
                <TopToolbar title={"Pod Usage Metrics"}/>
                <Paper className={classes.root}>

                    <div className={classes.filters}>
                        <FormControl className={classes.formControl}>
                            <InputLabel htmlFor="pod">Pod</InputLabel>
                            <Select
                                native
                                value={selectedPod}
                                onChange={this.handleChange("pod")}
                                inputProps={{
                                    name: "pod",
                                    id: "pod"
                                }}
                            >
                                <option value="all">All</option>
                            </Select>
                        </FormControl>
                        <Button variant="outlined" size="small" color="primary" className={classes.button}>
                            Update
                        </Button>
                    </div>
                    <div className={classes.graphs}>
                        <Metrics/>
                    </div>
                </Paper>
            </React.Fragment>
        );
    }

}

Pod.propTypes = {
    classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Pod);
