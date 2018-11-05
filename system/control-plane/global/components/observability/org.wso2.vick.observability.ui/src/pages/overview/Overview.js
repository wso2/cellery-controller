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

import React, {Component} from "react";
import PropTypes from 'prop-types';
import {withStyles} from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import DependencyGraph from "./DependencyGraph";

const graphConfig = {
    "automaticRearrangeAfterDropNode": false,
    "collapsible": false,
    "height": 800,
    "highlightDegree": 1,
    "highlightOpacity": 0.2,
    "linkHighlightBehavior": true,
    "maxZoom": 8,
    "minZoom": 0.1,
    "nodeHighlightBehavior": true,
    "panAndZoom": false,
    "staticGraph": false,
    "width": 1400,
    "d3": {
        "alphaTarget": 0.05,
        "gravity": -100,
        "linkLength": 100,
        "linkStrength": 1
    },
    "node": {
        "color": "#d3d3d3",
        "fontColor": "black",
        "fontSize": 14,
        "fontWeight": "normal",
        "highlightColor": "red",
        "highlightFontSize": 12,
        "highlightFontWeight": "bold",
        "highlightStrokeColor": "SAME",
        "highlightStrokeWidth": 1.5,
        "labelProperty": "name",
        "mouseCursor": "pointer",
        "opacity": 1,
        "renderLabel": true,
        "size": 600,
        "strokeColor": "green",
        "strokeWidth": 2,
        "svg": "red-cell.svg"
    },
    "link": {
        "color": "#d3d3d3",
        "opacity": 1,
        "semanticStrokeWidth": false,
        "strokeWidth": 4,
        "highlightColor": "black"
    }
};

const styles = {
    card: {
        minWidth: 275,
    },
    bullet: {
        display: 'inline-block',
        margin: '0 2px',
        transform: 'scale(0.8)',
    },
    title: {
        fontSize: 14,
    },
    pos: {
        marginBottom: 12,
    }
};

const cardCssStyle = {
    width: 300,
    height: 300,
    position: 'fixed',
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
                    },
                ]
            },
        };
        this.state = JSON.parse(JSON.stringify(this.defaultState));
        this.onMouseOverCell = this.onMouseOverCell.bind(this);
        this.onMouseOutCell = this.onMouseOutCell.bind(this);
    }

    onMouseOverCell(nodeId) {
        this.setState(prevState => ({
            summary: {
                ...prevState.summary,
                topic: 'Cell : ' + nodeId
            }
        }));
    };

    onMouseOutCell(nodeId) {
        this.setState({summary: JSON.parse(JSON.stringify(this.defaultState)).summary});
    }


    render() {
        const {classes} = this.props;
        const data = {
            nodes: [
                {id: 'Harry', "svg": "green-cell.svg", "onMouseOverNode": this.onMouseOverCell},
                {id: 'Sally', "svg": "yello-cell.svg", "onMouseOverNode": this.onMouseOverCell},
                {id: 'Alice', "onMouseOverNode": this.onMouseOverCell},
                {id: 'Sinthuja', "onMouseOverNode": this.onMouseOverCell}],
            links: [{source: 'Harry', target: 'Sally'}, {source: 'Harry', target: 'Alice'}]
        };
        return (
            <div>
                <DependencyGraph
                    id="graph-id"
                    data={data}
                    config={graphConfig}
                    onMouseOverNode={this.onMouseOverCell}
                    onMouseOutNode={this.onMouseOutCell}
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
                        {this.state.summary.content.map((element) =>
                            <Typography variant="subtitle1" gutterBottom>
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
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Overview);
