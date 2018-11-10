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
import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import PropTypes from "prop-types";
import React from "react";
import ReactDOM from "react-dom";
import Span from "../utils/span";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography";
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
    },
    resizeHandle: {
        transform: "translateX(-5px)",
        backgroundColor: "#919191",
        cursor: "ew-resize",
        position: "absolute",
        width: "5px",
        height: "100%",
        userSelect: "none"
    },
    spanDurationLabel: {
        maxHeight: null
    },
    spanDescriptionContent: {
        margin: "5px 12px 5px 7px"
    }
});

class Timeline extends React.Component {

    constructor(props) {
        super(props);

        this.timelineNode = React.createRef();
        this.timeline = null;
        this.spanLabelWidth = 400;
        this.selectedSpan = null;
    }

    componentDidMount() {
        const {classes, traceTree, spans} = this.props;
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
        let minLimit = Number.MAX_VALUE;
        let maxLimit = 0;
        traceTree.walk((span) => {
            if (span.treeDepth > treeHeight) {
                treeHeight = span.treeDepth;
            }
            if (span.startTime < minLimit) {
                minLimit = span.startTime;
            }
            if (span.startTime + span.duration > maxLimit) {
                maxLimit = span.startTime + span.duration;
            }
        });
        treeHeight += 1;
        const duration = (maxLimit - minLimit);
        minLimit -= duration * 0.05;
        maxLimit += duration * 0.12;

        const options = {
            orientation: "top",
            showMajorLabels: true,
            editable: false,
            selectable: false,
            groupEditable: false,
            showCurrentTime: false,
            min: new Date(minLimit),
            max: new Date(maxLimit),
            start: new Date(minLimit),
            end: new Date(maxLimit),
            order: (itemA, itemB) => itemA.order - itemB.order,
            groupTemplate: (item) => {
                const newElement = document.createElement("div");
                if (item && item.span.serviceName) {
                    const kindData = kindsData[item.span.kind];
                    const isLeaf = item.span.children.size;
                    ReactDOM.render((
                        <div>
                            <div style={{
                                paddingLeft: `${(item.span.treeDepth + (isLeaf === 0 ? 1 : 0)) * 15}px`,
                                minWidth: `${(treeHeight + 1) * 15 + 100}px`,
                                width: `${(self.spanLabelWidth > 0 ? self.spanLabelWidth : null)}px`
                            }} className={classNames(classes.spanLabelContainer)}>
                                <span className={classNames(classes.serviceName)}>{`${item.span.serviceName} `}</span>
                                <span className={classNames(classes.operationName)}>{item.span.operationName}</span>
                            </div>
                            {(kindData
                                ? <div className={classNames(classes.kindBadge)}
                                    style={{backgroundColor: kindData.color}}>{kindData.name}</div>
                                : null
                            )}
                        </div>
                    ), newElement);
                }
                return newElement;
            },
            template: (item, element) => {
                const newElement = document.createElement("div");
                let content = <span>{item.content}</span>;
                if (item.itemType === Timeline.Constants.ItemType.SPAN) {
                    content = <span>{item.span.duration} ms</span>;
                } else if (item.itemType === Timeline.Constants.ItemType.SPAN_DESCRIPTION) {
                    const rows = [];
                    for (const key in item.span.tags) {
                        if (item.span.tags.hasOwnProperty(key)) {
                            rows.push({
                                key: key,
                                value: item.span.tags[key]
                            });
                        }
                    }
                    if (rows.length > 0) {
                        content = (
                            <Card className={classes.spanDescriptionContent}>
                                <CardContent>
                                    <Typography color="textSecondary" gutterBottom>Tags</Typography>
                                    <Table>
                                        <TableBody>
                                            {
                                                rows.map((row, index) => (
                                                    <TableRow hover key={index}>
                                                        <TableCell component="th" scope="row">
                                                            <div>{row.key}</div>
                                                        </TableCell>
                                                        <TableCell>
                                                            <div>{row.value}</div>
                                                        </TableCell>
                                                    </TableRow>
                                                ))
                                            }
                                        </TableBody>
                                    </Table>
                                </CardContent>
                            </Card>
                        );
                    }
                }
                ReactDOM.render(content, newElement);
                element.setAttribute(Timeline.Constants.SPAN_ID_ATTRIBUTE_KEY, item.span.getUniqueId());
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
        const reversedSpans = spans.slice().reverse();

        // Clear previous interactions (eg:- resizability) added to the timeline
        const selector = `.${classNames(classes.spanLabelContainer)}`;
        Timeline.clearInteractions(`.${classNames(classes.resizeHandle)}`);
        Timeline.clearInteractions(selector);

        // Creating / Updating the timeline
        if (!this.timeline) {
            this.timeline = new vis.Timeline(this.timelineNode.current);
            this.timeline.on("changed", () => {
                // Adjust span description
                const timelineWindowWidth = document.querySelector(".vis-foreground").offsetWidth;
                const fitDescriptionToTimelineWindow = (node) => {
                    node.style.left = "0px";
                    node.style.width = `${timelineWindowWidth}px`;
                };
                document.querySelectorAll("div.vis-item-span-description")
                    .forEach(fitDescriptionToTimelineWindow);
                document.querySelectorAll(".vis-item-content")
                    .forEach(fitDescriptionToTimelineWindow);

                // Adjust span duration labels
                document.querySelectorAll("div.vis-item-span > div.vis-item-overflow").forEach((node) => {
                    node.style.transform = `translateX(${node.offsetWidth + 7}px)`;
                });
            });
            this.timeline.on("click", (event) => {
                if (event.what === "item" || event.what === "background") {
                    if (this.selectedSpan === event.group) {
                        this.selectedSpan = null;
                    } else {
                        this.selectedSpan = event.group;
                    }
                    this.updateTimelineItems(reversedSpans, {min: minLimit, max: maxLimit});
                }
            });
        }

        this.timeline.setOptions(options);
        this.updateTimelineItems(reversedSpans, {min: minLimit, max: maxLimit});

        // Add resizable functionality
        this.addHorizontalResizability(selector);
    }

    componentWillUnmount() {
        const {classes} = this.props;

        // Removing interactions added to the timeline
        Timeline.clearInteractions(`.${classNames(classes.resizeHandle)}`);
        Timeline.clearInteractions(`.${classNames(classes.spanLabelContainer)}`);

        // Destroying the timeline
        if (this.timeline) {
            this.timeline.destroy();
        }
    }

    /**
     * Update the items in the timeline.
     *
     * @param {Array.<Span>} spans The spans which should be displayed in the timeline
     * @param {{min: number, max: number}} limits of the timeline
     */
    updateTimelineItems(spans, limits) {
        // Populating timeline items and groups
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

            // Adding span
            items.push({
                id: `${span.getUniqueId()}-span`,
                itemType: Timeline.Constants.ItemType.SPAN,
                order: 1,
                start: new Date(span.startTime),
                end: new Date(span.startTime + span.duration),
                group: span.getUniqueId(),
                className: "vis-item-span",
                span: span
            });
            if (this.selectedSpan && this.selectedSpan === span.getUniqueId()) {
                items.push({
                    id: `${span.getUniqueId()}-span-description`,
                    itemType: Timeline.Constants.ItemType.SPAN_DESCRIPTION,
                    order: 2,
                    start: new Date(limits.min),
                    end: new Date(limits.max),
                    group: span.getUniqueId(),
                    className: "vis-item-span-description",
                    span: span
                });
            }
            groups.push({
                id: span.getUniqueId(),
                order: i * 2,
                nestedGroups: nestedGroups.length > 0 ? nestedGroups : null,
                span: span
            });
        }

        this.timeline.setGroups(new vis.DataSet(groups));
        this.timeline.setItems(new vis.DataSet(items));
    }

    /**
     * Add resizability to a set of items in the timeline.
     *
     * @param {string} selector The CSS selector of the items to which the resizability should be added
     */
    addHorizontalResizability(selector) {
        const self = this;
        const edges = {right: true};
        const {classes} = this.props;

        // Add the horizontal resize handle
        const newNode = document.createElement("div");
        newNode.classList.add(classNames(classes.resizeHandle));
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
        interact(`.${classNames(classes.resizeHandle)}`).on("down", (event) => {
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
    classes: PropTypes.any.isRequired,
    traceTree: PropTypes.instanceOf(Span),
    spans: PropTypes.arrayOf(PropTypes.instanceOf(Span))
};

Timeline.Constants = {
    SPAN_ID_ATTRIBUTE_KEY: "spanId",
    ItemType: {
        SPAN: "span",
        SPAN_DESCRIPTION: "description"
    }
};

export default withStyles(styles, {withTheme: true})(Timeline);
