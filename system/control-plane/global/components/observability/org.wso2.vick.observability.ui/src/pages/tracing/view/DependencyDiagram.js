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

import "./DependencyDiagram.css";
import Constants from "../../common/constants";
import {Graph} from "react-d3-graph";
import PropTypes from "prop-types";
import React from "react";
import Span from "../utils/span";
import TracingUtils from "../utils/tracingUtils";
import {withStyles} from "@material-ui/core";
import withColor, {ColorGenerator} from "../../common/color";

const styles = () => ({
    graph: {
        width: "100%",
        height: "100%"
    }
});

class DependencyDiagram extends React.Component {

    static MIN_RADIUS = 60;
    static MAX_RADIUS = 120;

    render = () => {
        const {classes, spans, colorGenerator} = this.props;
        const rootSpan = TracingUtils.getTreeRoot(spans);

        const nodeIdList = [];
        const nodes = [];
        const links = [];
        const getUniqueNodeId = (span) => (
            `${span.cell && span.cell.name ? `${span.cell.name}:` : ""}${span.serviceName}`
        );
        const addNodeIfNotPresent = (span) => {
            if (!nodeIdList.includes(getUniqueNodeId(span))) {
                // Finding the proper color for this item
                let colorKey = span.cell ? span.cell.name : null;
                if (!colorKey) {
                    if (span.componentType === Constants.ComponentType.VICK) {
                        colorKey = ColorGenerator.VICK;
                    } else if (span.componentType === Constants.ComponentType.ISTIO) {
                        colorKey = ColorGenerator.ISTIO;
                    } else {
                        colorKey = span.componentType;
                    }
                }
                const color = colorGenerator.getColor(colorKey);

                nodeIdList.push(getUniqueNodeId(span));
                nodes.push({
                    id: getUniqueNodeId(span),
                    color: color,
                    size: 350,
                    span: span
                });
            }
        };
        const addLink = (sourceSpan, destinationSpan) => {
            const link = {
                source: getUniqueNodeId(sourceSpan),
                target: getUniqueNodeId(destinationSpan)
            };
            if (sourceSpan.hasError() || destinationSpan.hasError()) {
                link.color = colorGenerator.getColor(ColorGenerator.ERROR);
            }
            links.push(link);
        };
        rootSpan.walk((span, data) => {
            let linkSource = data;
            if (!Constants.System.SIDECAR_AUTH_FILTER_OPERATION_NAME_PATTERN.test(span.operationName)
                && !Constants.System.ISTIO_MIXER_NAME_PATTERN.test(span.serviceName)) {
                if (linkSource && span.kind === Constants.Span.Kind.SERVER) { // Ending link traversing
                    addNodeIfNotPresent(span);
                    addLink(linkSource, span);
                    linkSource = null;
                } else if (!linkSource && span.kind === Constants.Span.Kind.CLIENT) { // Starting link traversing
                    addNodeIfNotPresent(span);
                    linkSource = span;
                }
            }
            return linkSource;
        }, null);

        let minDuration = Number.MAX_SAFE_INTEGER;
        let maxDuration = 0;
        for (const node of nodes) {
            if (node.span.duration < minDuration) {
                minDuration = node.span.duration;
            }
            if (node.span.duration > maxDuration) {
                maxDuration = node.span.duration;
            }
        }

        return (
            nodes.length > 0 && links.length > 0
                ? (
                    <Graph id={"trace-dependency-graph"} className={classes.graph}
                        data={{
                            nodes: nodes,
                            links: links
                        }}
                        config={{
                            directed: true,
                            nodeHighlightBehavior: true,
                            highlightOpacity: 0.2,
                            node: {
                                fontSize: 15,
                                highlightFontSize: 15,
                                viewGenerator: (node) => {
                                    const radius = (((node.span.duration - minDuration)
                                        * (DependencyDiagram.MAX_RADIUS - DependencyDiagram.MIN_RADIUS))
                                        / (maxDuration - minDuration)) + DependencyDiagram.MIN_RADIUS;

                                    let nodeSVGContent;
                                    const circle = <circle cx="120" cy="120" r={radius} fill={node.color}/>;
                                    if (node.span.hasError()) {
                                        const errorColor = colorGenerator.getColor(ColorGenerator.ERROR);

                                        const iconTranslation = radius * (Math.PI / 4);
                                        const xTranslation = 150;
                                        const yTranslation = 120 - iconTranslation - 30;
                                        nodeSVGContent = (
                                            <g>
                                                <g>
                                                    <g>
                                                        {circle}
                                                    </g>
                                                </g>
                                                <g transform={
                                                    `translate(${xTranslation}, ${yTranslation}) scale(0.4, 0.4)`
                                                }>
                                                    <path stroke="#fff" strokeWidth="10" fill={errorColor}
                                                        d="M120.5,9.6C59.1,9.6,9,59.8,9,121.3S59.1,233,120.5,
                                                             233S232,182.8,232,121.3S181.9,9.6,120.5,9.6z"/>
                                                    <path fill="#ffffff"
                                                        d="M105.4,164.5h29.9v29.9h-29.9V164.5z M105.4,
                                                        44.2h29.9v90.1h-29.9V44.2z"/>
                                                </g>
                                            </g>
                                        );
                                    } else {
                                        nodeSVGContent = circle;
                                    }
                                    return (
                                        <svg x="0" y="0" width="100%" height="100%" viewBox="0 0 240 240">
                                            {nodeSVGContent}
                                        </svg>
                                    );
                                }
                            },
                            link: {
                                strokeWidth: 3,
                                highlightColor: "#444"
                            },
                            d3: {
                                gravity: -200
                            }
                        }}
                    />
                )
                : null
        );
    };

}

DependencyDiagram.propTypes = {
    classes: PropTypes.object.isRequired,
    spans: PropTypes.arrayOf(PropTypes.instanceOf(Span)).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired
};

export default withStyles(styles, {withTheme: true})(withColor(DependencyDiagram));
