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
import "./Timeline.css";
import Constants from "./utils/constants";
import PropTypes from "prop-types";
import React from "react";
import ReactDOM from "react-dom";
import Span from "./utils/span";
import TracingUtils from "./utils/tracingUtils";
import classNames from "classnames";
import interact from "interactjs";
import vis from "vis";
import {withStyles} from "@material-ui/core";

const styles = () => ({
    spanLabelContainer: {
        width: 500,
        whiteSpace: "nowrap",
        overflow: "hidden",
        textOverflow: "ellipsis",
        boxSizing: "border-box",
        display: "inline-block"
    },
    serviceName: {
        fontWeight: 500,
        fontSize: "normal"
    },
    operationName: {
        color: "#7c7c7c",
        fontSize: "small"
    },
    kindBadge: {
        borderRadius: "8px",
        color: "white",
        padding: "2px 5px",
        marginLeft: "15px",
        fontSize: "12px",
        display: "inline-block"
    }
});

class Timeline extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            traceId: props.match.params.traceId,
            traceTree: null,
            spans: []
        };

        this.timelineNode = React.createRef();
        this.timeline = null;
        this.spanLabelWidth = 0;
        this.RESIZE_HANDLE_CLASS = "vis-resize-handle";

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
                traceTree: rootSpan,
                spans: TracingUtils.getOrderedList(rootSpan)
            });
        });
    }

    componentDidUpdate() {
        const {classes} = this.props;
        const self = this;
        const kindsData = {
            CLIENT: {
                name: "Client",
                color: "#9c27b0"
            },
            SERVER: {
                name: "Server",
                color: "#4caf50"
            },
            PRODUCER: {
                name: "Producer",
                color: "#03a9f4"
            },
            CONSUMER: {
                name: "Consumer",
                color: "#009688"
            }
        };

        // Finding the maximum tree height
        let treeHeight = 0;
        this.state.traceTree.walk((span) => {
            if (span.treeDepth > treeHeight) {
                treeHeight = span.treeDepth;
            }
        });
        treeHeight += 1;

        const options = {
            orientation: "top",
            showMajorLabels: true,
            editable: false,
            groupEditable: false,
            showCurrentTime: false,
            groupTemplate: (item) => {
                const newElement = document.createElement("div");
                const kindData = kindsData[item.kind];
                ReactDOM.render((
                    <div>
                        <div style={{
                            paddingLeft: `${(item.depth + (item.isLeaf ? 1 : 0)) * 15}px`,
                            minWidth: `${(treeHeight + 1) * 15 + 100}px`,
                            width: (self.spanLabelWidth > 0 ? self.spanLabelWidth : null)
                        }} className={classNames(classes.spanLabelContainer)}>
                            <span className={classNames(classes.serviceName)}>{`${item.serviceName} `}</span>
                            <span className={classNames(classes.operationName)}>{item.operationName}</span>
                        </div>
                        {(kindData
                            ? <div className={classNames(classes.kindBadge)}
                                style={{backgroundColor: kindData.color}}>{kindData.name}</div>
                            : null
                        )}
                    </div>
                ), newElement);
                return newElement;
            }
        };

        /*
         * Due to the way the vis timeline inserts the items (insert children at each node) the order of the spans
         * drawn gets messed up. To avoid this the tree should be drawn from leaves to the root node. To achieve this
         * and make sure the tree is right-side up a dual reverse (reversing span array and reversing group order
         * function) is used.
         */
        options.groupOrder = (a, b) => b.order - a.order;
        const spans = this.state.spans.reverse();

        // Populating timeline items and groups
        const groups = [];
        const items = [];
        for (let i = 0; i < spans.length; i++) {
            const span = spans[i];

            // Fetching all the children and grand-chilren of this node
            const nestedGroups = [];
            span.walk((currentSpan, data) => {
                if (currentSpan !== span) {
                    data.push(`${span.getUniqueId()}-span-group`);
                }
                return data;
            }, nestedGroups);

            // Adding span
            items.push({
                id: `${span.getUniqueId()}-span`,
                start: new Date(span.startTime),
                end: new Date(span.startTime + span.duration),
                content: `${span.duration} ms`,
                group: `${span.getUniqueId()}-span-group`
            });
            groups.push({
                id: `${span.getUniqueId()}-span-group`,
                order: i,
                nestedGroups: nestedGroups.length > 0 ? nestedGroups : null,
                depth: span.treeDepth,
                isLeaf: span.children.size === 0,
                serviceName: span.serviceName,
                operationName: span.operationName,
                kind: span.kind
            });
        }

        // Clear previous interactions (eg:- resizability) added to the timeline
        const selector = `.${classNames(classes.spanLabelContainer)}`;
        Timeline.clearInteractions(`.${this.RESIZE_HANDLE_CLASS}`);
        Timeline.clearInteractions(selector);

        // Creating / Updating the timeline
        if (!this.timeline) {
            this.timeline = new vis.Timeline(this.timelineNode.current);
        }
        this.timeline.setOptions(options);
        this.timeline.setGroups(new vis.DataSet(groups));
        this.timeline.setItems(new vis.DataSet(items));
        this.timeline.on("change", () => {
            if (self.spanLabelWidth > 0) {
                // Set the last known width upon timeline redraw
                document.querySelectorAll(`.${classNames(classes.spanLabelContainer)}`).forEach((node) => {
                    node.style.width = self.spanLabelWidth;
                });
            }
        });

        // Add resizable functionality
        this.addHorizontalResizability(selector);
    }

    /**
     * Add resizability to a set of items in the timeline.
     *
     * @param {string} selector The CSS selector of the items to which the resizability should be added
     */
    addHorizontalResizability(selector) {
        const self = this;
        const edges = {right: true};

        // Add the horizontal resize handle
        const newNode = document.createElement("div");
        newNode.classList.add(this.RESIZE_HANDLE_CLASS);
        const parent = document.querySelector(".vis-panel.vis-top");
        parent.insertBefore(newNode, parent.childNodes[0]);

        // Handling the resizing
        interact(selector).resizable({
            manualStart: true,
            edges: edges
        }).on("resizemove", (event) => {
            const targets = event.target;
            targets.forEach((target) => {
                // Update the element's style
                target.style.width = `${event.rect.width}px`;

                // Store the current width to be used when the timeline is redrawn
                self.spanLabelWidth = event.rect.width;

                // Trigger timeline redraw
                self.timeline.body.emitter.emit("_change");
            });
        });

        // Handling dragging of the resize handle
        interact(".vis-resize-handle").on("down", (event) => {
            event.interaction.start(
                {
                    name: "resize",
                    edges: edges
                },
                interact(selector),
                document.querySelectorAll(selector)
            );
        });
    }

    /**
     * Clear interactions added to the timeline.
     *
     * @param {string} selector The CSS selector of the items from which the interactions should be cleared
     */
    static clearInteractions(selector) {
        interact(selector).unset();
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
    }).isRequired,
    classes: PropTypes.any.isRequired
};

export default withStyles(styles, {withTheme: true})(Timeline);
