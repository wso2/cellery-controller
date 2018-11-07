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

import "vis/dist/vis-timeline-graph2d.min.css";
import Constants from "./utils/constants";
import PropTypes from "prop-types";
import React from "react";
import Span from "./utils/span";
import TracingUtils from "./utils/tracingUtils";
import vis from "vis";

class Timeline extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            traceId: props.match.params.traceId,
            spanTreeRootNode: null,
            spans: []
        };

        this.timelineNode = React.createRef();
        this.timeline = null;

        // TODO : Remove this section and query the backend and retrieve spans.
        const self = this;
        setTimeout(() => {
            const globalGatewayServerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-a-id",
                parentSpanId: this.state.traceId,
                serviceName: Constants.VICK.System.GLOBAL_GATEWAY_NAME,
                operationName: "get-hr-info",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000000,
                duration: 3160000,
                tags: {}
            });
            const globalGatewayClientSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-b-id",
                parentSpanId: "span-a-id",
                serviceName: Constants.VICK.System.GLOBAL_GATEWAY_NAME,
                operationName: "call-hr-cell",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10010000,
                duration: 3110000,
                tags: {}
            });
            const hrCellGatewayServerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-b-id",
                parentSpanId: "span-a-id",
                serviceName: "src:0.0.0.hr_1_0_0_employee",
                operationName: "call-hr-cell",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10020000,
                duration: 3090000,
                tags: {}
            });
            const hrCellGatewayClientSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-c-id",
                parentSpanId: "span-b-id",
                serviceName: "src:0.0.0.hr_1_0_0_employee",
                operationName: "get-employee-data",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10030000,
                duration: 3060000,
                tags: {}
            });
            const employeeServiceServerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-c-id",
                parentSpanId: "span-b-id",
                serviceName: "employee-service",
                operationName: "get-employee-data",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10040000,
                duration: 3040000,
                tags: {}
            });
            const employeeServiceToIstioMixerClientSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-d-id",
                parentSpanId: "span-c-id",
                serviceName: "employee-service",
                operationName: "is-authorized",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10050000,
                duration: 990000,
                tags: {}
            });
            const istioMixerServerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-d-id",
                parentSpanId: "span-c-id",
                serviceName: Constants.VICK.System.ISTIO_MIXER_NAME,
                operationName: "is-authorized",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10060000,
                duration: 940000,
                tags: {}
            });
            const istioMixerWorkerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-e-id",
                parentSpanId: "span-d-id",
                serviceName: Constants.VICK.System.ISTIO_MIXER_NAME,
                operationName: "authorization",
                startTime: 10070000,
                duration: 890000,
                tags: {}
            });
            const employeeServiceToReviewsServiceClientSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-f-id",
                parentSpanId: "span-c-id",
                serviceName: "employee-service",
                operationName: "get-reviews",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 11050000,
                duration: 990000,
                tags: {}
            });
            const reviewsServiceServerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-f-id",
                parentSpanId: "span-c-id",
                serviceName: "reviews-service",
                operationName: "get-reviews",
                kind: Constants.Span.Kind.SERVER,
                startTime: 11100000,
                duration: 890000,
                tags: {}
            });
            const employeeServiceToStockOptionsCellClientSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-g-id",
                parentSpanId: "span-c-id",
                serviceName: "employee-service",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 12060000,
                duration: 990000,
                tags: {}
            });
            const stockOptionsCellGatewayServerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-g-id",
                parentSpanId: "span-c-id",
                serviceName: "src:0.0.0.stock_options_1_0_0_employee",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.SERVER,
                startTime: 12100000,
                duration: 890000,
                tags: {}
            });
            const stockOptionsCellGatewayClientSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-h-id",
                parentSpanId: "span-g-id",
                serviceName: "src:0.0.0.stock_options_1_0_0_employee",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 12150000,
                duration: 790000,
                tags: {}
            });
            const stockOptionsServiceServerSpan = new Span({
                traceId: this.state.traceId,
                spanId: "span-h-id",
                parentSpanId: "span-g-id",
                serviceName: "stock-options-service",
                operationName: "get-employee-stock-options",
                kind: Constants.Span.Kind.SERVER,
                startTime: 12200000,
                duration: 690000,
                tags: {}
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
                spanTreeRootNode: rootSpan,
                spans: TracingUtils.getOrderedList(rootSpan)
            });
        });
    }

    componentDidUpdate() {
        const options = {
            width: "100%",
            height: "700px",
            orientation: "top",
            showMajorLabels: true,
            editable: false,
            groupEditable: false
        };

        /*
         * Due to the way the vis timeline inserts the items (insert children at each node) the order of the spans
         * drawn gets messed up. To avoid this the tree should be drawn from leaves to the root node. To achieve this
         * and make sure the tree is right-side up a dual reverse (reversing span array and reversing group order
         * function) is used.
         */
        options.groupOrder = (a, b) => b.order - a.order;
        const spans = this.state.spans.reverse();

        const groups = [];
        const items = [];
        for (let i = 0; i < spans.length; i++) {
            const span = spans[i];

            // Fetching all the children and grand-chilren of this node
            const nestedGroups = [];
            span.walk((currentSpan, data) => {
                if (currentSpan !== span) {
                    data.push(currentSpan.getUniqueId());
                }
                return data;
            }, nestedGroups);

            items.push({
                id: `${span.getUniqueId()}-span`,
                start: new Date(span.startTime),
                end: new Date(span.startTime + span.duration),
                content: `${span.duration} ms`,
                group: span.getUniqueId()
            });
            groups.push({
                id: span.getUniqueId(),
                order: i,
                content: `${span.serviceName} <small style="color: #7c7c7c;">${span.operationName}</small>`,
                nestedGroups: nestedGroups.length > 0 ? nestedGroups : null
            });
        }

        // Creating / Updating the timeline
        if (!this.timeline) {
            this.timeline = new vis.Timeline(this.timelineNode.current);
        }
        this.timeline.setOptions(options);
        this.timeline.setGroups(new vis.DataSet(groups));
        this.timeline.setItems(new vis.DataSet(items));
    }

    render() {
        return <div ref={this.timelineNode}/>;
    }

}

Timeline.propTypes = {
    match: PropTypes.shape({
        params: PropTypes.shape({
            traceId: PropTypes.string.isRequired
        }).isRequired
    }).isRequired
};

export default Timeline;
