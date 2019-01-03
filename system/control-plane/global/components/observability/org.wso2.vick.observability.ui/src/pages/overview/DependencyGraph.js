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

import ErrorBoundary from "../common/error/ErrorBoundary";
import {Graph} from "react-d3-graph";
import React from "react";
import * as PropTypes from "prop-types";

class DependencyGraph extends React.Component {

    shouldComponentUpdate = (nextProps) => nextProps.reloadGraph;

    render = () => {
        const {data, ...otherProps} = this.props;

        // Finding distinct links
        const links = [];
        if (data.links) {
            data.links.forEach((link) => {
                const linkMatches = links.find(
                    (existingEdge) => existingEdge.source === link.source && existingEdge.target === link.target);
                if (!linkMatches) {
                    links.push({
                        source: link.source,
                        target: link.target
                    });
                }
            });
        }

        let view;
        if (data.nodes && data.nodes.length > 0) {
            view = (
                <ErrorBoundary title={"Unable to Render"} description={"Unable to Render due to Invalid Data"}>
                    <Graph {...otherProps} data={{...data, links: links}}/>
                </ErrorBoundary>
            );
        } else {
            view = <div>No Data Available</div>;
        }
        return view;
    };

}

DependencyGraph.propTypes = {
    data: PropTypes.object.isRequired,
    reloadGraph: PropTypes.bool
};

export default DependencyGraph;
