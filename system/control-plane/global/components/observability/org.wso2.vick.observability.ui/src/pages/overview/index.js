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

import Button from "@material-ui/core/Button";
import Card from "@material-ui/core/Card";
import CardActions from "@material-ui/core/CardActions";
import CardContent from "@material-ui/core/CardContent";
import DependencyGraph from "./DependencyGraph";
import PropTypes from "prop-types";
import React from "react";
import Typography from "@material-ui/core/Typography";
import axios from "axios";
import {withRouter} from "react-router-dom";
import {withStyles} from "@material-ui/core/styles";
import withColor, {ColorGenerator} from "../common/color";

const graphConfig = {
    directed: true,
    automaticRearrangeAfterDropNode: false,
    collapsible: false,
    height: 800,
    highlightDegree: 1,
    highlightOpacity: 0.2,
    linkHighlightBehavior: false,
    maxZoom: 8,
    minZoom: 0.1,
    nodeHighlightBehavior: true,
    panAndZoom: false,
    staticGraph: false,
    width: 1400,
    d3: {
        alphaTarget: 0.05,
        gravity: -1500,
        linkLength: 150,
        linkStrength: 1
    },
    node: {
        color: "#d3d3d3",
        fontColor: "black",
        fontSize: 18,
        fontWeight: "normal",
        highlightColor: "red",
        highlightFontSize: 18,
        highlightFontWeight: "bold",
        highlightStrokeColor: "SAME",
        highlightStrokeWidth: 1.5,
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
        highlightColor: "black"
    }
};

const styles = {
    card: {
        minWidth: 275,
        width: 300,
        height: 500,
        position: "fixed",
        bottom: 0,
        right: 0,
        top: 70
    },
    title: {
        fontSize: 14
    },
    moreDetailsButton: {
        left: 170,
        bottom: 0
    }
};

class Overview extends React.Component {

    constructor(props) {
        super(props);
        const colorGenerator = this.props.colorGenerator;
        graphConfig.node.viewGenerator = this.viewGenerator;
        this.defaultState = {
            summary: {
                topic: "VICK Deployment",
                content: [
                    {
                        key: "Total cells",
                        value: ""
                    },
                    {
                        key: "Successful cells",
                        value: ""
                    },
                    {
                        key: "Failed cells",
                        value: ""
                    },
                    {
                        key: "Warning cells",
                        value: ""
                    }
                ]
            },
            isOverallSummary: true,
            data: {
                nodes: null,
                links: null
            },
            error: null,
            reloadGraph: true
        };
        this.state = JSON.parse(JSON.stringify(this.defaultState));
        // TODO: Update the url to the WSO2-sp worker node.
        axios({
            method: "GET",
            url: "http://localhost:9123/dependency-model/cell-overview"
        }).then((response) => {
            const result = response.data;
            const summaryContent = [
                {
                    key: "Total cells",
                    value: result.nodes.length
                },
                {
                    key: "Successful cells",
                    value: result.nodes.length
                },
                {
                    key: "Failed cells",
                    value: "0"
                },
                {
                    key: "Warning cells",
                    value: "0"
                }
            ];
            this.defaultState.summary.content = summaryContent;
            colorGenerator.addKeys(result.nodes);
            this.setState((prevState) => ({
                data: {
                    nodes: result.nodes,
                    links: result.edges
                },
                summary: {
                    ...prevState.summary,
                    content: summaryContent
                }
            }));
        }).catch((error) => {
            this.setState({error: error});
        });
    }

    viewGenerator = (nodeProps) => {
        const color = this.props.colorGenerator.getColor(nodeProps.id);
        return <svg x="0px" y="0px"
            width="50px" height="50px" viewBox="0 0 240 240">
            <polygon fill={color} points="224,179.5 119.5,239.5 15,179.5 15,59.5 119.5,-0.5 224,59.5 "/>
        </svg>;
    };

    onClickCell = (nodeId) => {
        const outbound = new Set();
        const inbound = new Set();
        this.state.data.links.forEach((element) => {
            if (element.source === nodeId) {
                outbound.add(element.target);
            } else if (element.target === nodeId) {
                inbound.add(element.source);
            }
        });
        const services = new Set();
        this.state.data.nodes.forEach((element) => {
            if (element.id === nodeId) {
                element.services.forEach((service) => {
                    services.add(service);
                });
            }
        });
        this.setState((prevState) => ({
            summary: {
                ...prevState.summary,
                topic: `Cell : ${nodeId}`,
                content: [
                    {
                        key: "Outbound Cells",
                        setValue: this.populateArray(outbound)
                    },
                    {
                        key: "Inbound Cells",
                        setValue: this.populateArray(inbound)
                    },
                    {
                        key: "Micro Services",
                        setValue: this.populateArray(services)
                    }

                ]
            },
            reloadGraph: false,
            isOverallSummary: false
        }));
    };

    populateArray = (setElements) => {
        const arrayElements = [];
        setElements.forEach((setElement) => {
            arrayElements.push(setElement);
        });
        return arrayElements;
    };

    onClickGraph = () => {
        this.setState({
            summary: JSON.parse(JSON.stringify(this.defaultState)).summary,
            reloadGraph: true,
            isOverallSummary: true
        });
    };


    render = () => {
        const {classes} = this.props;
        return (
            <div>
                <DependencyGraph
                    id="graph-id"
                    data={this.state.data}
                    config={graphConfig}
                    reloadGraph={this.state.reloadGraph}
                    onClickNode={this.onClickCell}
                    onClickGraph={this.onClickGraph}
                />
                <Card className={classes.card} transformOrigin={{
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
                            {(element.setValue && element.setValue.length > 0)
                                && (
                                    <ul>
                                        {element.setValue.map((setValueElement) => <li
                                            key={setValueElement}>{setValueElement}</li>)}
                                    </ul>)
                            }
                            {element.value || (element.setValue && element.setValue.length > 0) ? "" : "None"}
                        </Typography>
                        )}

                    </CardContent>
                    {(!this.state.isOverallSummary) && (
                        <CardActions>
                            <Button size="small" variant="contained" color="primary"
                                className={classes.moreDetailsButton}>
                                More Details
                            </Button>
                        </CardActions>
                    )}
                </Card>
            </div>
        );
    };

}

Overview.propTypes = {
    classes: PropTypes.object.isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator)
};

export default withStyles(styles)(withRouter(withColor(Overview)));
