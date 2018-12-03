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

import Details from "../common/Details";
import Grey from "@material-ui/core/colors/grey";
import K8sObjects from "./K8sObjects";
import Metrics from "../common/Metrics";
import Paper from "@material-ui/core/Paper";
import PropTypes from "prop-types";
import React from "react";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import TopToolbar from "../../common/toptoolbar";
import {withStyles} from "@material-ui/core/styles";

const styles = (theme) => ({
    root: {
        flexGrow: 1,
        backgroundColor: theme.palette.background.paper,
        padding: theme.spacing.unit * 3,
        paddingTop: 0,
        margin: theme.spacing.unit
    },
    tabs: {
        marginBottom: theme.spacing.unit * 2,
        borderBottomWidth: 1,
        borderBottomStyle: "solid",
        borderBottomColor: Grey[200]
    }
});

class MicroService extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            selectedMicroservice: 0
        };
    }

    handleChange = (event, value) => {
        this.setState({
            selectedMicroservice: value
        });
    };

    render() {
        const {classes} = this.props;
        const {selectedMicroservice} = this.state;

        const details = <Details/>;
        const k8sObjects = <K8sObjects/>;
        const metrics = <Metrics/>;
        const tabContent = [details, k8sObjects, metrics];

        return (
            <React.Fragment>
                <TopToolbar title={"Microservice Name"}/>
                <Paper className={classes.root}>

                    <Tabs
                        value={selectedMicroservice}
                        onChange={this.handleChange}
                        indicatorColor="primary"
                        textColor="primary"
                        className={classes.tabs}
                    >
                        <Tab label="DETAILS"/>
                        <Tab label="K8S OBJECTS"/>
                        <Tab label="METRICS"/>
                    </Tabs>
                    {tabContent[selectedMicroservice]}
                </Paper>
            </React.Fragment>
        );
    }

}


MicroService.propTypes = {
    classes: PropTypes.object.isRequired
};

export default withStyles(styles)(MicroService);
