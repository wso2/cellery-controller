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

import Constants from "../utils/constants";
import PropTypes from "prop-types";
import React from "react";
import SequenceDiagram from "./SequenceDiagram";
import Span from "../utils/span";
import Tab from "@material-ui/core/Tab";
import Tabs from "@material-ui/core/Tabs";
import Timeline from "./Timeline";
import TracingUtils from "../utils/tracingUtils";

class View extends React.Component {

    constructor(props) {
        super(props);
        const traceId = props.match.params.traceId;

        this.state = {
            traceTree: null,
            spans: [],
            selectedTabIndex: 0,
            isLoading: true
        };

        this.handleTabChange = this.handleTabChange.bind(this);

        // TODO : Remove this section and query the backend and retrieve spans.
        const self = this;
        setTimeout(() => {
            const globalGatewayServerSpan = new Span({
                traceId: traceId,
                spanId: "span-a-id",
                parentSpanId: traceId,
                serviceName: Constants.VICK.System.GLOBAL_GATEWAY_NAME,
                operationName: "get-hr-info",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000000,
                duration: 3160000,
                tags: {
                    "span.kind": Constants.Span.Kind.SERVER,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const globalGatewayClientSpan = new Span({
                traceId: traceId,
                spanId: "span-b-id",
                parentSpanId: "span-a-id",
                serviceName: Constants.VICK.System.GLOBAL_GATEWAY_NAME,
                operationName: "call-hr-cell",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10010000,
                duration: 3110000,
                tags: {
                    "span.kind": Constants.Span.Kind.CLIENT,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const hrCellGatewayServerSpan = new Span({
                traceId: traceId,
                spanId: "span-b-id",
                parentSpanId: "span-a-id",
                serviceName: "src:0.0.0.hr_1_0_0_employee",
                operationName: "call-hr-cell",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10020000,
                duration: 3090000,
                tags: {
                    "span.kind": Constants.Span.Kind.SERVER,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const hrCellGatewayClientSpan = new Span({
                traceId: traceId,
                spanId: "span-c-id",
                parentSpanId: "span-b-id",
                serviceName: "src:0.0.0.hr_1_0_0_employee",
                operationName: "get-employee-data",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10030000,
                duration: 3060000,
                tags: {
                    "span.kind": Constants.Span.Kind.CLIENT,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const employeeServiceServerSpan = new Span({
                traceId: traceId,
                spanId: "span-c-id",
                parentSpanId: "span-b-id",
                serviceName: "employee-service",
                operationName: "get-employee-data",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10040000,
                duration: 3040000,
                tags: {
                    "span.kind": Constants.Span.Kind.SERVER,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const employeeServiceToIstioMixerClientSpan = new Span({
                traceId: traceId,
                spanId: "span-d-id",
                parentSpanId: "span-c-id",
                serviceName: "employee-service",
                operationName: "is-authorized",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10050000,
                duration: 990000,
                tags: {
                    "span.kind": Constants.Span.Kind.CLIENT,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const istioMixerServerSpan = new Span({
                traceId: traceId,
                spanId: "span-d-id",
                parentSpanId: "span-c-id",
                serviceName: Constants.VICK.System.ISTIO_MIXER_NAME,
                operationName: "is-authorized",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10060000,
                duration: 940000,
                tags: {
                    "span.kind": Constants.Span.Kind.SERVER,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const istioMixerWorkerSpan = new Span({
                traceId: traceId,
                spanId: "span-e-id",
                parentSpanId: "span-d-id",
                serviceName: Constants.VICK.System.ISTIO_MIXER_NAME,
                operationName: "authorization",
                startTime: 10070000,
                duration: 890000,
                tags: {
                    "cache.hit": true
                }
            });
            const employeeServiceToReviewsServiceClientSpan = new Span({
                traceId: traceId,
                spanId: "span-f-id",
                parentSpanId: "span-c-id",
                serviceName: "employee-service",
                operationName: "get-reviews",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 11050000,
                duration: 990000,
                tags: {
                    "span.kind": Constants.Span.Kind.CLIENT,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const reviewsServiceServerSpan = new Span({
                traceId: traceId,
                spanId: "span-f-id",
                parentSpanId: "span-c-id",
                serviceName: "reviews-service",
                operationName: "get-reviews",
                kind: Constants.Span.Kind.SERVER,
                startTime: 11100000,
                duration: 890000,
                tags: {
                    "span.kind": Constants.Span.Kind.SERVER,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const employeeServiceToStockOptionsCellClientSpan = new Span({
                traceId: traceId,
                spanId: "span-g-id",
                parentSpanId: "span-c-id",
                serviceName: "employee-service",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 12060000,
                duration: 990000,
                tags: {
                    "span.kind": Constants.Span.Kind.CLIENT,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const stockOptionsCellGatewayServerSpan = new Span({
                traceId: traceId,
                spanId: "span-g-id",
                parentSpanId: "span-c-id",
                serviceName: "src:0.0.0.stock_options_1_0_0_employee",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.SERVER,
                startTime: 12100000,
                duration: 890000,
                tags: {
                    "span.kind": Constants.Span.Kind.SERVER,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const stockOptionsCellGatewayClientSpan = new Span({
                traceId: traceId,
                spanId: "span-h-id",
                parentSpanId: "span-g-id",
                serviceName: "src:0.0.0.stock_options_1_0_0_employee",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 12150000,
                duration: 790000,
                tags: {
                    "span.kind": Constants.Span.Kind.CLIENT,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const stockOptionsServiceServerSpan = new Span({
                traceId: traceId,
                spanId: "span-h-id",
                parentSpanId: "span-g-id",
                serviceName: "stock-options-service",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.SERVER,
                startTime: 12200000,
                duration: 690000,
                tags: {
                    "span.kind": Constants.Span.Kind.SERVER,
                    "http.status.code": 200,
                    "http.method": "GET"
                }
            });
            const spans = [
                globalGatewayServerSpan, globalGatewayClientSpan, hrCellGatewayServerSpan, hrCellGatewayClientSpan,
                employeeServiceServerSpan, employeeServiceToIstioMixerClientSpan, istioMixerServerSpan,
                istioMixerWorkerSpan, employeeServiceToReviewsServiceClientSpan, reviewsServiceServerSpan,
                employeeServiceToStockOptionsCellClientSpan, stockOptionsCellGatewayServerSpan,
                stockOptionsCellGatewayClientSpan, stockOptionsServiceServerSpan
            ];
            const rootSpan = TracingUtils.buildTree(spans);
            TracingUtils.labelSpanTree(rootSpan);

            self.setState({
                traceTree: rootSpan,
                spans: TracingUtils.getOrderedList(rootSpan),
                isLoading: false
            });
        });
    }

    handleTabChange(event, value) {
        this.setState({
            selectedTabIndex: value
        });
    }

    render() {
        const {traceTree, spans, selectedTabIndex} = this.state;

        const timeline = <Timeline traceTree={traceTree} spans={spans}/>;
        const sequenceDiagram = <SequenceDiagram/>;
        const tabContent = [timeline, sequenceDiagram];

        return (
            this.state.isLoading
                ? null
                : (
                    <div>
                        <Tabs value={selectedTabIndex} indicatorColor="primary"
                            onChange={this.handleTabChange}>
                            <Tab label="Timeline"/>
                            <Tab label="Sequence Diagram"/>
                        </Tabs>
                        {tabContent[selectedTabIndex]}
                    </div>
                )
        );
    }

}

View.propTypes = {
    match: PropTypes.shape({
        params: PropTypes.shape({
            traceId: PropTypes.string.isRequired
        }).isRequired
    }).isRequired
};

export default View;
