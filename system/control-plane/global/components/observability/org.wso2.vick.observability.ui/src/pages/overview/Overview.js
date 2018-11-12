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

import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import DependencyGraph from "./DependencyGraph";
import PropTypes from "prop-types";
import Typography from "@material-ui/core/Typography";
import {withStyles} from "@material-ui/core/styles";
import React, {Component} from "react";

const graphConfig = {
    automaticRearrangeAfterDropNode: false,
    collapsible: false,
    height: 800,
    highlightDegree: 1,
    highlightOpacity: 0.2,
    linkHighlightBehavior: true,
    maxZoom: 8,
    minZoom: 0.1,
    nodeHighlightBehavior: true,
    panAndZoom: false,
    staticGraph: false,
    width: 1400,
    d3: {
        alphaTarget: 0.05,
        gravity: -2000,
        linkLength: 100,
        linkStrength: 1
    },
    node: {
        color: "#d3d3d3",
        fontColor: "black",
        fontSize: 14,
        fontWeight: "normal",
        highlightColor: "red",
        highlightFontSize: 12,
        highlightFontWeight: "bold",
        highlightStrokeColor: "SAME",
        highlightStrokeWidth: 1.5,
        labelProperty: "name",
        mouseCursor: "pointer",
        opacity: 1,
        renderLabel: true,
        size: 600,
        strokeColor: "green",
        strokeWidth: 2,
        svg: "red-cell.svg"
    },
    link: {
        color: "#d3d3d3",
        opacity: 1,
        semanticStrokeWidth: false,
        strokeWidth: 4,
        highlightColor: "black"
    }
};

const styles = {
    card: {
        minWidth: 275
    },
    title: {
        fontSize: 14
    },
    pos: {
        marginBottom: 12
    }
};

const cardCssStyle = {
    width: 300,
    height: 300,
    position: "fixed",
    bottom: 0,
    right: 0,
    top: 70
};

class Overview extends Component {

    constructor(props) {
        super(props);
        this.defaultState = {
            summary: {
                topic: "VICK Deployment",
                content: [
                    {
                        key: "Total cells",
                        value: "11"
                    },
                    {
                        key: "Successful cells",
                        value: "5"
                    },
                    {
                        key: "Failed cells",
                        value: "5"
                    },
                    {
                        key: "Warning cells",
                        value: "1"
                    }
                ]
            },
            data: {
                nodes: null,
                links: null,
            },
            error: null,
            reloadGraph: true
        };
        this.state = JSON.parse(JSON.stringify(this.defaultState));
        fetch("http://localhost:9123/dependencyModel/graph")
            .then(res => res.json())
            .then(
                (result) => {
                    console.log("result came...");
                    console.log(result);
                    this.setState({
                        data: {
                            nodes: result.nodes,
                            links: result.edges
                        }
                    });

                    // TODO: testing...
                    // let nodeName = "Harry";
                    // const results = {
                    //     nodes: [
                    //         {id: nodeName, svg: "green-cell.svg", onMouseOverNode: this.onMouseOverCell},
                    //         {id: "Sally", svg: "yello-cell.svg", onMouseOverNode: this.onMouseOverCell},
                    //         {id: "Alice", onMouseOverNode: this.onMouseOverCell},
                    //         {id: "Sinthuja", onMouseOverNode: this.onMouseOverCell}
                    //     ],
                    //     links: [{source: nodeName, target: "Sally"}, {source: nodeName, target: "Alice"}]
                    // };
                    // this.setState({
                    //     data: {
                    //         nodes: results.nodes,
                    //         links: results.links
                    //     }
                    // });
                },
                (error) => {
                    this.setState({error: error});
                }
            );
        console.log(this.state.data);
        this.onMouseOverCell = this.onMouseOverCell.bind(this);
        this.onMouseOutCell = this.onMouseOutCell.bind(this);
    }

    onMouseOverCell(nodeId) {
        this.setState((prevState) => ({
            summary: {
                ...prevState.summary,
                topic: `Cell : ${nodeId}`
            },
            reloadGraph: false
        }));
    }

    onMouseOutCell(nodeId) {
        this.setState({
            summary: JSON.parse(JSON.stringify(this.defaultState)).summary,
            reloadGraph: true
        });
    }




    render() {
        const {classes} = this.props;
        // let nodeName = "Harry";
        // const data = {
        //     nodes: [
        //         {id: nodeName, svg: "green-cell.svg", onMouseOverNode: this.onMouseOverCell},
        //         {id: "Sally", svg: "yello-cell.svg", onMouseOverNode: this.onMouseOverCell},
        //         {id: "Alice", onMouseOverNode: this.onMouseOverCell},
        //         {id: "Sinthuja", onMouseOverNode: this.onMouseOverCell}
        //     ],
        //     links: [{source: nodeName, target: "Sally"}, {source: nodeName, target: "Alice"}]
        // };
        return (
            <div>
                <DependencyGraph
                    id="graph-id"
                    data={this.state.data}
                    config={graphConfig}
                    onMouseOverNode={this.onMouseOverCell}
                    onMouseOutNode={this.onMouseOutCell}
                    reloadGraph={this.state.reloadGraph}
                />
                <Card className={classes.card} style={cardCssStyle} transformOrigin={{
                    vertical: "top",
                    horizontal: "right"
                }}>
                    <CardContent>
                        <Typography className={classes.title} color="textSecondary" gutterBottom>
                            Summary
                        </Typography>
                        <Typography variant="h5" component="h2">
                            {this.state.summary.topic}
                        </Typography>
                        <br/>
                        {this.state.summary.content.map((element) => <Typography variant="subtitle1" key={element.key}
                                                                                 gutterBottom>
                                {element.key} : {element.value}
                            </Typography>
                        )}
                    </CardContent>
                </Card>
            </div>
        );
    }

}

Overview.propTypes = {
    classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Overview);
