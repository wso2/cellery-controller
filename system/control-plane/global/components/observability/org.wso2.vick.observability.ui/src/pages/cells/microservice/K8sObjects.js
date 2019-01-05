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

import BarChartIcon from "@material-ui/icons/BarChart";
import Green from "@material-ui/core/colors/green";
import IconButton from "@material-ui/core/IconButton";
import {Link} from "react-router-dom";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography/Typography";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core/styles";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    container: {
        marginTop: theme.spacing.unit * 3
    },
    successIcon: {
        color: Green[600]
    },
    table: {
        overflowX: "auto",
        marginBottom: theme.spacing.unit * 3
    },
    tableCell: {
        borderBottom: "none"
    },
    cellWidth10: {
        width: "10%"
    },
    cellWidth20: {
        width: "20%"
    },
    subtitle: {
        fontWeight: 400,
        fontSize: "1rem"
    }
});

class K8sObjects extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isLoading: false
        };
    }

    componentDidMount = () => {
        const {globalState} = this.props;
        globalState.addListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
    };

    componentWillUnmount = () => {
        const {globalState} = this.props;
        globalState.removeListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
    };

    handleLoadingStateChange = (loadingStateKey, oldState, newState) => {
        this.setState({
            isLoading: newState.loadingOverlayCount > 0
        });
    };

    render = () => {
        const {classes} = this.props;
        const {isLoading} = this.state;
        return (
            isLoading
                ? null
                : (
                    <React.Fragment>
                        <div className={classes.container}>
                            <Typography color="inherit" className={classes.subtitle}>Service</Typography>
                            <Table className={classes.table}>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Service</TableCell>
                                        <TableCell>Type</TableCell>
                                        <TableCell className={classes.cellWidth20}>Creation Time</TableCell>
                                        <TableCell className={classes.cellWidth20}>Age</TableCell>
                                        <TableCell className={classes.cellWidth10}/>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    <TableRow key="">
                                        <TableCell className={classes.tableCell} component="th" scope="row"/>
                                        <TableCell className={classes.tableCell}/>
                                        <TableCell className={classes.tableCell}/>
                                        <TableCell className={classes.tableCell}/>
                                        <TableCell className={classes.tableCell}/>
                                    </TableRow>
                                </TableBody>
                            </Table>
                            <Typography color="inherit" className={classes.subtitle}>Workload</Typography>
                            <Table className={classes.table}>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Workload</TableCell>
                                        <TableCell className={classes.cellWidth20}>Creation Time</TableCell>
                                        <TableCell className={classes.cellWidth20}>Age</TableCell>
                                        <TableCell className={classes.cellWidth10}/>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    <TableRow key="">
                                        <TableCell className={classes.tableCell} component="th" scope="row"/>
                                        <TableCell className={classes.tableCell}/>
                                        <TableCell className={classes.tableCell}/>
                                        <TableCell className={classes.tableCell}/>
                                    </TableRow>
                                </TableBody>
                            </Table>
                            <Typography color="inherit" className={classes.subtitle}>Pods</Typography>
                            <Table className={classes.table}>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Pod</TableCell>
                                        <TableCell>Restarts</TableCell>
                                        <TableCell className={classes.cellWidth20}>Age</TableCell>
                                        <TableCell>Metrics</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    <TableRow key="">
                                        <TableCell component="th" scope="row"/>
                                        <TableCell className={classes.cellWidth20}/>
                                        <TableCell className={classes.cellWidth20}/>
                                        <TableCell className={classes.cellWidth10}>
                                            <IconButton className={classes.button} size="small" color="inherit"
                                                component={Link} to="/system-metrics/pod-usage">
                                                <BarChartIcon/>
                                            </IconButton>
                                        </TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </div>
                    </React.Fragment>
                )
        );
    }

}

K8sObjects.propTypes = {
    classes: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired
};

export default withStyles(styles)(withGlobalState(K8sObjects));
