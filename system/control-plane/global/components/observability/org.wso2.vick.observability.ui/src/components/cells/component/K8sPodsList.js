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

import BarChart from "@material-ui/icons/BarChart";
import Constants from "../../../utils/constants";
import DataTable from "../../common/DataTable";
import HttpUtils from "../../../utils/api/httpUtils";
import IconButton from "@material-ui/core/IconButton/IconButton";
import Menu from "@material-ui/core/Menu/Menu";
import MenuItem from "@material-ui/core/MenuItem/MenuItem";
import NotFound from "../../common/error/NotFound";
import NotificationUtils from "../../../utils/common/notificationUtils";
import QueryUtils from "../../../utils/common/queryUtils";
import React from "react";
import StateHolder from "../../common/state/stateHolder";
import moment from "moment";
import withGlobalState from "../../common/state";
import {withStyles} from "@material-ui/core";
import * as PropTypes from "prop-types";

class K8sPodsList extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            podInfo: [],
            isLoading: false,
            metricsPopperElement: null,
            openMetricsPopoverPod: null
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
        this.loadPodInfo(isUserAction, startTime, endTime);
    };

    loadPodInfo = (isUserAction, queryStartTime, queryEndTime) => {
        const {globalState, cell, component} = this.props;
        const self = this;

        const search = {
            queryStartTime: queryStartTime.valueOf(),
            queryEndTime: queryEndTime.valueOf(),
            cell: cell,
            component: component
        };

        if (isUserAction) {
            NotificationUtils.showLoadingOverlay("Loading Pod Info", globalState);
            self.setState({
                isLoading: true
            });
        }
        HttpUtils.callObservabilityAPI(
            {
                url: `/k8s/pods${HttpUtils.generateQueryParamString(search)}`,
                method: "GET"
            },
            globalState
        ).then((data) => {
            const podInfo = data.map((dataItem) => ({
                name: dataItem[2],
                creationTimestamp: dataItem[3],
                lastKnownAliveTimestamp: dataItem[4],
                nodeName: dataItem[5]
            }));

            self.setState({
                podInfo: podInfo
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
                    "Failed to load pod information",
                    NotificationUtils.Levels.ERROR,
                    globalState
                );
            }
        });
    };

    metricsButtonRenderer = (value) => {
        const self = this;
        const {globalState} = self.props;
        const {metricsPopperElement, openMetricsPopoverPod} = self.state;

        const podMetricsLinkTemplate = globalState.get(StateHolder.CONFIG).templates.kubernetesMetricsLinks.pod;
        const nodeMetricsLinkTemplate = globalState.get(StateHolder.CONFIG).templates.kubernetesMetricsLinks.node;

        const handleMetricsMenuOpen = (event) => {
            self.setState({
                metricsPopperElement: event.currentTarget,
                openMetricsPopoverPod: value.podName
            });
        };
        const handleMetricsMenuClose = () => {
            self.setState({
                metricsPopperElement: null,
                openMetricsPopoverPod: null
            });
        };
        const replaceLinkTemplate = (linkTemplate) => {
            const queryStartTIme = QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).startTime).valueOf();
            const queryEndTIme = QueryUtils.parseTime(globalState.get(StateHolder.GLOBAL_FILTER).endTime).valueOf();

            return linkTemplate.replace(/\${pod}/g, value.podName)
                .replace(/\${node}/g, value.nodeName)
                .replace(/\${fromTime}/g, queryStartTIme)
                .replace(/\${toTime}/g, queryEndTIme);
        };
        const onPodMetricsMenuItemClick = () => {
            handleMetricsMenuClose();
            window.open(replaceLinkTemplate(podMetricsLinkTemplate));
        };
        const onNodeMetricsMenuItemClick = () => {
            handleMetricsMenuClose();
            window.open(replaceLinkTemplate(nodeMetricsLinkTemplate));
        };

        const anchorEl = openMetricsPopoverPod === value.podName ? metricsPopperElement : null;
        return (
            <React.Fragment>
                <IconButton size="small" color="inherit" onClick={handleMetricsMenuOpen}
                    disabled={!podMetricsLinkTemplate && !nodeMetricsLinkTemplate}>
                    <BarChart/>
                </IconButton>
                {
                    podMetricsLinkTemplate || nodeMetricsLinkTemplate
                        ? (
                            <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMetricsMenuClose}>
                                {
                                    podMetricsLinkTemplate
                                        ? <MenuItem onClick={onPodMetricsMenuItemClick}>Pod Metrics</MenuItem>
                                        : null
                                }
                                {
                                    nodeMetricsLinkTemplate
                                        ? <MenuItem onClick={onNodeMetricsMenuItemClick}>Node Metrics</MenuItem>
                                        : null
                                }
                            </Menu>
                        )
                        : null
                }
            </React.Fragment>
        );
    };

    render = () => {
        const {component} = this.props;
        const {podInfo, isLoading} = this.state;

        const columns = [
            {
                name: "Pod"
            },
            {
                name: "Created Timestamp",
                options: {
                    customBodyRender: (value) => moment(value).format(Constants.Pattern.DATE_TIME)
                }
            },
            {
                name: "Last Known Alive Timestamp",
                options: {
                    customBodyRender: (value) => moment(value).format(Constants.Pattern.DATE_TIME)
                }
            },
            {
                name: "Node"
            },
            {
                name: "",
                options: {
                    customBodyRender: this.metricsButtonRenderer
                }
            }
        ];
        const options = {
            filter: false
        };

        let listView;
        if (podInfo.length > 0) {
            listView = <DataTable columns={columns} options={options} data={podInfo.map((podDatum) => [
                podDatum.name,
                podDatum.creationTimestamp,
                podDatum.lastKnownAliveTimestamp,
                podDatum.nodeName,
                {
                    podName: podDatum.name,
                    nodeName: podDatum.nodeName
                }
            ])}/>;
        } else {
            listView = (
                <NotFound title={"No Pods Found"} description={`No Pods found for "${component}" component `
                    + "in the selected time range"}/>
            );
        }

        return (isLoading ? null : listView);
    };

}

K8sPodsList.propTypes = {
    globalState: PropTypes.instanceOf(StateHolder).isRequired,
    cell: PropTypes.string.isRequired,
    component: PropTypes.string.isRequired
};

export default withStyles({})(withGlobalState(K8sPodsList));
