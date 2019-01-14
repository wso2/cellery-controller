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

import ErrorBoundary from "./error/ErrorBoundary";
import {Graph} from "react-d3-graph";
import React from "react";
import UnknownError from "./error/UnknownError";
import * as PropTypes from "prop-types";

class DependencyGraph extends React.Component {

    static DEFAULT_GRAPH_CONFIG = {
        directed: true,
        automaticRearrangeAfterDropNode: false,
        collapsible: false,
        highlightDegree: 1,
        highlightOpacity: 0.2,
        linkHighlightBehavior: false,
        maxZoom: 8,
        minZoom: 0.1,
        nodeHighlightBehavior: true,
        panAndZoom: false,
        staticGraph: false,
        height: 580,
        width: 1050,
        d3: {
            alphaTarget: 0.05,
            gravity: -1500,
            linkLength: 150,
            linkStrength: 1
        },
        node: {
            color: "#d3d3d3",
            fontColor: "#555",
            fontSize: 16,
            fontWeight: "normal",
            highlightColor: "red",
            highlightFontSize: 16,
            highlightFontWeight: "bold",
            highlightStrokeColor: "SAME",
            highlightStrokeWidth: 1,
            labelProperty: "name",
            mouseCursor: "pointer",
            opacity: 1,
            renderLabel: true,
            size: 600,
            strokeColor: "green",
            strokeWidth: 2
        },
        link: {
            color: "#d3d3d3",
            opacity: 1,
            semanticStrokeWidth: false,
            strokeWidth: 4,
            highlightColor: "#777"
        }
    };


    render = () => {
        const {data, config, ...otherProps} = this.props;
        // Finding distinct links
        const links = [];
        if (data.links) {
            data.links.forEach((link) => {
                const linkMatches = links.find(
                    (existingEdge) => existingEdge.source === link.source && existingEdge.target === link.target);
                if (!linkMatches) {
                    links.push({
                        source: link.source,
                        target: link.target,
                        edgeString: link.edgeString
                    });
                }
            });
        }

        let view;
        if (data.nodes && data.nodes.length > 0) {
            view = (
                <ErrorBoundary title={"Unable to Render"} description={"Unable to Render due to Invalid Data"}>
                    <Graph {...otherProps} data={{...data, links: links}} config={{
                        ...DependencyGraph.DEFAULT_GRAPH_CONFIG,
                        ...config,
                        d3: {
                            ...DependencyGraph.DEFAULT_GRAPH_CONFIG.d3,
                            ...config.d3
                        },
                        node: {
                            ...DependencyGraph.DEFAULT_GRAPH_CONFIG.node,
                            ...config.node
                        },
                        link: {
                            ...DependencyGraph.DEFAULT_GRAPH_CONFIG.link,
                            ...config.link
                        }
                    }}/>
                </ErrorBoundary>
            );
        } else {
            view = (
                <UnknownError title={"No Data Available"} description={"No Data Available to render Dependency Graph"}/>
            );
        }
        return view;
    };

}

DependencyGraph.propTypes = {
    data: PropTypes.object.isRequired,
    reloadGraph: PropTypes.bool
};

export default DependencyGraph;
