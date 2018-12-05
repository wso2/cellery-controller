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

import CheckCircleOutline from "@material-ui/icons/CheckCircleOutline";
import Green from "@material-ui/core/colors/green";
import PropTypes from "prop-types";
import React from "react";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography/Typography";
import {withStyles} from "@material-ui/core/styles";

const styles = (theme) => ({
    table: {
        width: "20%"
    },
    successIcon: {
        color: Green[600]
    },
    tableCell: {
        borderBottom: "none"
    },
    dependencies: {
        marginTop: theme.spacing.unit * 3
    },
    diagram: {
        padding: theme.spacing.unit * 3
    },
    subtitle: {
        fontWeight: 400,
        fontSize: "1rem"
    }
});

const Details = (props) => {
    const {classes} = props;
    return (
        <React.Fragment>
            <Table className={classes.table}>
                <TableBody>
                    <TableRow>
                        <TableCell className={classes.tableCell}>
                            <Typography color="textSecondary">
                                Namespace
                            </Typography>
                        </TableCell>
                        <TableCell className={classes.tableCell}>
                            <Typography>
                                Default
                            </Typography>
                        </TableCell>
                    </TableRow>
                    <TableRow>
                        <TableCell className={classes.tableCell}>
                            <Typography color="textSecondary">
                                Health
                            </Typography>
                        </TableCell>
                        <TableCell className={classes.tableCell}>
                            <CheckCircleOutline className={classes.successIcon}/>
                        </TableCell>
                    </TableRow>
                </TableBody>
            </Table>
            <div className={classes.dependencies}>
                <Typography color="textSecondary" className={classes.subtitle}>
                    Dependencies
                </Typography>
                <div className={classes.diagram}>
                    Dependency Diagram
                </div>
            </div>
        </React.Fragment>
    );
};

Details.propTypes = {
    classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Details);
