/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import HealthIndicator from "../../common/HealthIndicator";
import HttpUtils from "../../common/utils/httpUtils";
import {Link} from "react-router-dom";
import NotificationUtils from "../../common/utils/notificationUtils";
import QueryUtils from "../../common/utils/queryUtils";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography/Typography";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core/styles";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    table: {
        width: "20%",
        marginTop: 25
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

class Details extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            health: -1,
            dependencyGraphData: [],
            isLoading: false
        };
    }

    componentDidMount = () => {
        const {globalState} = this.props;

        globalState.addListener(StateHolder.LOADING_STATE, this.handleLoadingStateChange);
        this.update(
            true,
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime).valueOf(),
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime).valueOf()
        );
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

    update = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell, microservice} = this.props;
        const self = this;

        const search = {
            queryStartTime: queryStartTime.valueOf(),
            queryEndTime: queryEndTime.valueOf(),
            destinationCell: cell,
            destinationMicroservice: microservice
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Microservice Info", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/microservices/metrics/${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const aggregatedData = data.map((datum) => ({
                isError: datum[1] === "5xx",
                count: datum[5]
            })).reduce((accumulator, currentValue) => {
                if (currentValue.isError) {
                    accumulator.errorsCount += currentValue.count;
                }
                accumulator.total += currentValue.count;
                return accumulator;
            }, {
                errorsCount: 0,
                total: 0
            });

            self.setState({
                health: 1 - (aggregatedData.count === 0 ? aggregatedData.errorsCount / aggregatedData.count : 0),
                dependencyGraphData: []
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load microservice information",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    render() {
        const {classes, cell} = this.props;
        const {health, isLoading} = this.state;
        return (
            isLoading
                ? null
                : (
                    <React.Fragment>
                        <Table className={classes.table}>
                            <TableBody>
                                <TableRow>
                                    <TableCell className={classes.tableCell}>
                                        <Typography color="textSecondary">
                                            Health
                                        </Typography>
                                    </TableCell>
                                    <TableCell className={classes.tableCell}>
                                        <HealthIndicator value={health}/>
                                    </TableCell>
                                </TableRow>
                                <TableRow>
                                    <TableCell className={classes.tableCell}>
                                        <Typography color="textSecondary">
                                            Cell
                                        </Typography>
                                    </TableCell>
                                    <TableCell className={classes.tableCell}>
                                        <Link to={`/cells/${cell}`}>{cell}</Link>
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
                )
        );
    }

}

Details.propTypes = {
    classes: PropTypes.object.isRequired,
    cell: PropTypes.string.isRequired,
    microservice: PropTypes.string.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired
};

export default withStyles(styles)(withGlobalState(Details));
