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

import Details from "./common/Details";
import Grey from "@material-ui/core/colors/grey";
import Metrics from "./common/Metrics";
import Microservices from "./common/Table";
import Paper from "@material-ui/core/Paper";
import PropTypes from "prop-types";
import React from "react";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import TopToolbar from "../common/TopToolbar";
import {withStyles} from "@material-ui/core/styles";

const styles = (theme) => ({
    root: {
        flexGrow: 1,
        backgroundColor: theme.palette.background.paper,
        padding: theme.spacing.unit * 3,
        paddingTop: 0,
        margin: Number(theme.spacing.unit)
    },
    tabs: {
        marginBottom: theme.spacing.unit * 2,
        borderBottomWidth: 1,
        borderBottomStyle: "solid",
        borderBottomColor: Grey[200]
    }
});

class Cell extends React.Component {

    state = {
        value: 0
    };

    handleChange = (event, value) => {
        this.setState({value: value});
    };

    render() {
        const {classes} = this.props;
        const details = <Details></Details>;
        const microservices = <Microservices></Microservices>;
        const metrics = <Metrics isHidden={true}></Metrics>;
        const tabContent = [details, microservices, metrics];

        return (
            <React.Fragment>
                <TopToolbar title={"Cell Name"} onUpdate={this.loadCellData}/>
                <Paper className={classes.root}>
                    <Tabs
                        value={this.state.value}
                        onChange={this.handleChange}
                        indicatorColor="primary"
                        textColor="primary"
                        className={classes.tabs}
                    >
                        <Tab label="DETAILS"/>
                        <Tab label="MICROSERVICES"/>
                        <Tab label="METRICS"/>
                    </Tabs>
                    {tabContent[this.state.value]}
                </Paper>
            </React.Fragment>
        );
    }

}


Cell.propTypes = {
    classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Cell);
