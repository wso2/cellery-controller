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

import CellDependencyView from "./CellDependencyView";
import ColorGenerator from "../../common/color/colorGenerator";
import HealthIndicator from "../../common/HealthIndicator";
import HttpUtils from "../../../utils/api/httpUtils";
import NotFound from "../../common/error/NotFound";
import NotificationUtils from "../../../utils/common/notificationUtils";
import QueryUtils from "../../../utils/common/queryUtils";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography/Typography";
import withColor from "../../common/color";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core/styles";
import * as PropTypes from "prop-types";

const styles = () => ({
    table: {
        width: "20%",
        marginTop: 25
    },
    tableCell: {
        borderBottom: "none"
    }
});

class Details extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            isDataAvailable: false,
            health: -1,
            dependencyGraphData: [],
            isLoading: false
        };
    }

    componentDidMount = () => {
        const {globalState} = this.props;

        this.update(
            true,
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime),
            QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime)
        );
    };

    update = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell} = this.props;
        const self = this;

        const search = {
            queryStartTime: queryStartTime.valueOf(),
            queryEndTime: queryEndTime.valueOf(),
            destinationCell: cell
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Cell Info", globalState);
            self.setState({
                isLoading: true
            });
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/metrics/${HttpUtils.generateQueryParamString(search)}`,
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
                health: 1 - (aggregatedData.total === 0 ? aggregatedData.errorsCount / aggregatedData.total : 0),
                isDataAvailable: aggregatedData.total > 0
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                self.setState({
                    isLoading: false
                });
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                self.setState({
                    isLoading: false
                });
                NotificationUtils.showNotification(
                    "Failed to load cell information",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    render = () => {
        const {classes, cell} = this.props;
        const {health, isLoading, isDataAvailable} = this.state;

        let view;
        if (isDataAvailable) {
            view = (
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
                    </TableBody>
                </Table>
            );
        } else {
            view = (
                <NotFound title={"Cell Not Found"} description={`The "${cell}" cell not found. This is possibly `
                    + "because no requests had been received/sent by this cell in the selected time period"}/>
            );
        }
        return (
            <React.Fragment>
                {isLoading ? null : view}
                {
                    isDataAvailable
                        ? <CellDependencyView cell={cell}/>
                        : null
                }
            </React.Fragment>
        );
    }

}

Details.propTypes = {
    classes: PropTypes.object.isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    cell: PropTypes.string.isRequired
};

export default withStyles(styles)(withColor(withGlobalState(Details)));

