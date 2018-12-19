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

import Constants from "../../common/constants";
import DataTable from "../../common/DataTable";
import HealthIndicator from "../../common/HealthIndicator";
import HttpUtils from "../../common/utils/httpUtils";
import {Link} from "react-router-dom";
import NotificationUtils from "../../common/utils/notificationUtils";
import Paper from "@material-ui/core/Paper/Paper";
import PropTypes from "prop-types";
import QueryUtils from "../../common/utils/queryUtils";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core";

const styles = (theme) => ({
    root: {
        margin: theme.spacing.unit
    }
});

class MicroserviceList extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            microserviceInfo: []
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

    update = (isUserAction, startTime, endTime) => {
        this.loadMicroserviceInfo(isUserAction, startTime, endTime);
    };

    loadMicroserviceInfo = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell} = this.props;
        const self = this;

        const search = {
            queryStartTime: queryStartTime.valueOf(),
            queryEndTime: queryEndTime.valueOf()
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Microservice Info", globalState);
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/http-requests/cells/${cell}/microservices/${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const microserviceInfo = data.map((dataItem) => ({
                sourceCell: dataItem[0],
                sourceMicroservice: dataItem[1],
                destinationCell: dataItem[2],
                destinationMicroservice: dataItem[3],
                httpResponseGroup: dataItem[4],
                totalResponseTimeMilliSec: dataItem[5],
                requestCount: dataItem[6]
            }));

            self.setState({
                microserviceInfo: microserviceInfo
            });
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
            }
        }).catch(() => {
            if (isUserAction) {
                NotificationUtils.hideLoadingOverlay(globalState);
                NotificationUtils.showNotification(
                    "Failed to load microservice information",
                    StateHolder.NotificationLevels.ERROR,
                    globalState
                );
            }
        });
    };

    render = () => {
        const {classes, cell} = this.props;
        const {microserviceInfo} = this.state;
        const columns = [
            {
                name: "Health",
                options: {
                    customBodyRender: (value) => <HealthIndicator value={value}/>
                }
            },
            {
                name: "Microservice",
                options: {
                    customBodyRender: (value) => <Link to={`/cells/${cell}/microservices/${value}`}>{value}</Link>
                }
            },
            {
                name: "Inbound Error Rate",
                options: {
                    customBodyRender: (value) => `${value * 100} %`
                }
            },
            {
                name: "Outbound Error Rate",
                options: {
                    customBodyRender: (value) => `${value * 100} %`
                }
            },
            {
                name: "Average Response Time (ms)",
                options: {
                    customBodyRender: (value) => (Math.round(value))
                }
            },
            {
                name: "Average Inbound Request Count (requests/s)"
            }
        ];

        // Processing data to find the required values
        const dataTableMap = {};
        const initializeDataTableMapEntryIfNotPresent = (microservice) => {
            if (!dataTableMap[microservice]) {
                dataTableMap[microservice] = {
                    inboundErrorCount: 0,
                    outboundErrorCount: 0,
                    requestCount: 0,
                    totalResponseTimeMilliSec: 0
                };
            }
        };
        const isMicroserviceRelevant = (consideredCell, microservice) => (
            !Constants.System.GLOBAL_GATEWAY_NAME_PATTERN.test(microservice)
                && !Constants.Cell.GATEWAY_NAME_PATTERN.test(microservice) && consideredCell === cell
        );
        for (const microserviceDatum of microserviceInfo) {
            if (isMicroserviceRelevant(microserviceDatum.sourceCell, microserviceDatum.sourceMicroservice)) {
                initializeDataTableMapEntryIfNotPresent(microserviceDatum.sourceMicroservice);

                if (microserviceDatum.httpResponseGroup === "5xx") {
                    dataTableMap[microserviceDatum.sourceMicroservice].outboundErrorCount
                        += microserviceDatum.requestCount;
                }
            }
            if (isMicroserviceRelevant(microserviceDatum.destinationCell, microserviceDatum.destinationMicroservice)) {
                initializeDataTableMapEntryIfNotPresent(microserviceDatum.destinationMicroservice);

                if (microserviceDatum.httpResponseGroup === "5xx") {
                    dataTableMap[microserviceDatum.destinationMicroservice].inboundErrorCount
                        += microserviceDatum.requestCount;
                }
                dataTableMap[microserviceDatum.destinationMicroservice].requestCount += microserviceDatum.requestCount;
                dataTableMap[microserviceDatum.destinationMicroservice].totalResponseTimeMilliSec
                    += microserviceDatum.totalResponseTimeMilliSec;
            }
        }

        // Transforming the objects into 2D array accepted by the Table library
        const tableData = [];
        for (const microservice in dataTableMap) {
            if (dataTableMap.hasOwnProperty(microservice) && Boolean(microservice)) {
                const microserviceData = dataTableMap[microservice];
                tableData.push([
                    microserviceData.requestCount === 0
                        ? -1
                        : 1 - microserviceData.inboundErrorCount / microserviceData.requestCount,
                    microservice,
                    microserviceData.requestCount === 0
                        ? 0
                        : microserviceData.inboundErrorCount / microserviceData.requestCount,
                    microserviceData.requestCount === 0
                        ? 0
                        : microserviceData.outboundErrorCount / microserviceData.requestCount,
                    microserviceData.requestCount === 0
                        ? 0
                        : microserviceData.totalResponseTimeMilliSec / microserviceData.requestCount,
                    microserviceData.requestCount
                ]);
            }
        }

        return (
            <Paper className={classes.root}>
                <DataTable columns={columns} data={tableData}/>
            </Paper>
        );
    };

}

MicroserviceList.propTypes = {
    classes: PropTypes.object.isRequired,
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    cell: PropTypes.string.isRequired
};

export default withStyles(styles)(withGlobalState(MicroserviceList));
