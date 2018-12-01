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

import "./DependencyDiagram.css";
import {ColorGenerator} from "../../common/color/colorGenerator";
import Constants from "../utils/constants";
import {Graph} from "react-d3-graph";
import PropTypes from "prop-types";
import React from "react";
import Span from "../utils/span";
import TracingUtils from "../utils/tracingUtils";
import {withStyles} from "@material-ui/core";
import {ColorGeneratorConstants, withColor} from "../../common/color/index";

const styles = () => ({
    graph: {
        width: "100%",
        height: "100%"
    }
});

class DependencyDiagram extends React.Component {

    render() {
        const {classes, spans, colorGenerator} = this.props;
        const rootSpan = TracingUtils.buildTree(spans);

        const nodeIdList = [];
        const nodes = [];
        const links = [];
        const addNodeIfNotPresent = (span) => {
            if (!nodeIdList.includes(span.serviceName)) {
                // Finding the proper color for this item
                let colorKey = span.cell ? span.cell.name : null;
                if (!colorKey) {
                    if (span.componentType === Constants.Span.ComponentType.VICK) {
                        colorKey = ColorGeneratorConstants.VICK;
                    } else if (span.componentType === Constants.Span.ComponentType.ISTIO) {
                        colorKey = ColorGeneratorConstants.ISTIO;
                    } else {
                        colorKey = span.componentType;
                    }
                }
                const color = colorGenerator.getColor(colorKey);

                nodeIdList.push(span.serviceName);
                nodes.push({
                    id: span.serviceName,
                    color: color,
                    size: span.duration
                });
            }
        };
        const addLink = (sourceSpan, destinationSpan) => {
            const link = {
                source: sourceSpan.serviceName,
                target: destinationSpan.serviceName
            };
            if (sourceSpan.hasError() || destinationSpan.hasError()) {
                link.color = colorGenerator.getColor(ColorGeneratorConstants.ERROR);
            }
            links.push(link);
        };
        rootSpan.walk((span, data) => {
            let linkSource = data;
            if (linkSource && span.kind === Constants.Span.Kind.SERVER) { // Ending link traversing
                addNodeIfNotPresent(span);
                addLink(linkSource, span);
                linkSource = null;
            } else if (!linkSource && span.kind === Constants.Span.Kind.CLIENT) { // Starting link traversing
                addNodeIfNotPresent(span);
                linkSource = span;
            }
            return linkSource;
        }, null);

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
                                highlightFontSize: 15
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
    }

}

DependencyDiagram.propTypes = {
    classes: PropTypes.object.isRequired,
    spans: PropTypes.arrayOf(PropTypes.instanceOf(Span)).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator).isRequired
};

export default withStyles(styles, {withTheme: true})(withColor(DependencyDiagram));
