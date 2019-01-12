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

import "./SequenceDiagram.css";
import ChevronRight from "@material-ui/icons/ChevronRight";
import Constants from "../../../utils/constants";
import React from "react";
import Span from "../../../utils/tracing/span";
import TracingUtils from "../../../utils/tracing/tracingUtils";
import Typography from "@material-ui/core/Typography";
import interact from "interactjs";
import mermaid from "mermaid";
import {withStyles} from "@material-ui/core/styles";
import withColor, {ColorGenerator} from "../../common/color";
import * as PropTypes from "prop-types";

const styles = (theme) => ({
    newMessageText: {
        fill: "#4c4cb3 !important",
        cursor: "pointer"
    },
    subtitle: {
        fontWeight: 400,
        fontSize: "1.0rem"
    },
    mermaid: {
        padding: 10
    },
    navigation: {
        paddingTop: 20,
        paddingBottom: 5,
        marginLeft: 50
    },
    navList: {
        display: "block",
        float: "left"
    }
});

class SequenceDiagram extends React.Component {

    static GLOBAL = "global";
    static GLOBAL_GATEWAY = "global gateway";

    constructor(props) {
        super(props);
        this.state = {
            config: "",
            spanData: "sequenceDiagram \n",
            copyArr: [],
            clicked: false,
            cellName: null,
            clonedArray: [],
            cellClicked: "global",
            callIdClicked: ""
        };

        this.mermaidDivRef = React.createRef();

        this.addCells = this.addCells.bind(this);
        this.addServices = this.addServices.bind(this);
        this.drawCells = this.drawCells.bind(this);
    }

    render() {
        const {classes} = this.props;
        return (
            <div>
                <div className={classes.navigation}>
                    <span className={classes.navList}>
                        <Typography color="textSecondary" className={classes.subtitle} onClick={this.addCells}
                            style={this.state.clicked
                                ? {color: "#3e51b5", cursor: "pointer", textDecoration: "underline"}
                                : {}}>
                       Cells
                        </Typography>
                    </span>
                    <span className={classes.navList}>
                        <ChevronRight color="action" style={this.state.clicked ? {} : {display: "none"}}/>
                    </span>
                    <span className={classes.navList}>
                        <Typography color="textSecondary" className={classes.subtitle}
                            style={this.state.clicked ? {} : {display: "none"}}>
                            {this.state.cellClicked} cell [{this.state.callIdClicked}] - Services
                        </Typography>
                    </span>

                </div>
                <br/>
                <div className={classes.mermaid} ref={this.mermaidDivRef}>
                    {this.state.config}
                </div>
            </div>

        );
    }

    componentDidMount() {
        this.addCells();
        interact(".messageText").on("tap", (event) => {
            if ((event.srcElement.innerHTML !== "Return")
                && (this.state.clicked !== true)) {
                const numb = event.srcElement.innerHTML.match(/\d+/g).map(Number);
                this.addServices(numb);
            }
        });
        this.cloneArray();
    }

    componentDidUpdate(prevProps, prevState) {
        if (this.state.config !== prevState.config) {
            const collectionsMessage = this.mermaidDivRef.current.getElementsByClassName("messageText");
            this.mermaidDivRef.current.removeAttribute("data-processed");
            mermaid.init(this.mermaidDivRef.current);

            if (!this.state.clicked) {
                this.setMessageLinkStyle(collectionsMessage);
            }
        }

        if (this.state.config !== prevState.config || this.props.colorGenerator !== prevProps.colorGenerator) {
            const collectionsActor = this.mermaidDivRef.current.getElementsByClassName("actor");
            this.addActorColor(this.state.clicked, collectionsActor);
        }
    }

    /**
     * Sets the style for message links that are clickable (cell -level diagram).
     *
     * @param {Element []} messageElementArray The array of message link elements.
     */

    setMessageLinkStyle(messageElementArray) {
        const {classes} = this.props;
        for (let i = 0; i < messageElementArray.length; i++) {
            if (messageElementArray[i].innerHTML.match("\\s\\[([0-9]+)\\]+$")) {
                messageElementArray[i].classList.add(classes.newMessageText);
            }
        }
    }


    /**
     * Adds the relevant cell color, which is consistent throughout the dashboard, to the actors.
     *
     * @param {boolean} serviceClicked The variable to check if the user is in service-level diagram.
     * @param {Element[]} elementArray The array of elements with the class name `actor`.
     */
    addActorColor(serviceClicked, elementArray) {
        const {colorGenerator} = this.props;
        let color;
        let cellName;
        let actorStyle;
        if (serviceClicked) {
            cellName = this.state.cellClicked;
            if (cellName === SequenceDiagram.GLOBAL) {
                color = colorGenerator.getColor(ColorGenerator.VICK);
            } else {
                color = colorGenerator.getColor(SequenceDiagram.addDash(cellName));
            }
            actorStyle = `
                stroke: ${color};
                stroke-width: 3;
                fill: #ffffff`;

            for (let i = 0; i < elementArray.length; i += 2) {
                elementArray[i].style = actorStyle;
            }
        } else {
            // For loop with iteration by factor 2 to skip SVG `rect` element and get the text in each actor.
            for (let i = 1; i < elementArray.length; i += 2) {
                if (elementArray[i].firstElementChild !== null) {
                    const cellName = SequenceDiagram.addDash(elementArray[i].firstElementChild.innerHTML);
                    if (cellName === SequenceDiagram.addDash(SequenceDiagram.GLOBAL_GATEWAY)) {
                        color = colorGenerator.getColor(ColorGenerator.VICK);
                    } else {
                        color = colorGenerator.getColor(cellName);
                    }
                    actorStyle = `
                stroke: ${color};
                stroke-width: 3;
                fill: #ffffff`;
                    // Index of i-1 is given to set the style to the respective SVG `rect` element.
                    elementArray[i - 1].style = actorStyle;
                }
            }
        }
    }

    /**
     * Create a copy of the original span list
     */

    cloneArray() {
        this.setState({
            clonedArray: this.props.spans
        });
    }


    /**
     * Adds the service calls made for a particular cell to the diagram.
     *
     * @param {number[]} callId The span's call Id of the particular cell call.
     */

    addServices(callId) {
        let data2 = "sequenceDiagram \n";
        const treeRoot = this.state.clonedArray[SequenceDiagram.findSpanIndexCall(this.state.clonedArray, callId)];
        const parentName = treeRoot.cell.name;
        this.setState({
            cellClicked: parentName,
            callIdClicked: callId
        });
        data2 += `activate ${SequenceDiagram.removeDash(treeRoot.serviceName)}\n`;
        let j = 0;
        treeRoot.walk(
            (span) => {
                if (!span.isFromIstioSystemComponent() && !span.isFromVICKSystemComponent()) {
                    if (!span.callingId && parentName === span.cell.name) {
                        if (span.parent.serviceName !== span.serviceName) {
                            j += 1;
                            data2 += `${`${SequenceDiagram.removeDash(span.parent.serviceName)}  ->>+`
                                + `${SequenceDiagram.removeDash(span.serviceName)}:`}${span.operationName}`
                                + `- [${callId}.${j}] \n`;
                        }
                    }
                }
            }, null,
            (span) => {
                if (!span.isFromIstioSystemComponent() && !span.isFromVICKSystemComponent()) {
                    data2 += SequenceDiagram.updateTextDatawithReturn(span, parentName);
                }
            },
            (span) => (!span.isFromIstioSystemComponent() && !span.isFromVICKSystemComponent()
                && !span.callingId && parentName !== span.parent.cell.name)
        );
        data2 += `deactivate ${SequenceDiagram.removeDash(treeRoot.serviceName)}\n`;
        this.setState({
            config: data2,
            clicked: true
        });
    }

    /**
     * Updates the text data, which is used by the mermaid library to generate diagrams, with return drawn.
     *
     * @param {span} span The span array.
     * @param {String} parentName The parent cell name
     * @return {String} text The updated text
     */

    static updateTextDatawithReturn(span, parentName) {
        let text = "";
        if (!span.callingId && parentName === span.cell.name) {
            if (span.parent.serviceName !== span.serviceName) {
                text += `${SequenceDiagram.removeDash(span.serviceName)}-->>- `
                    + `${SequenceDiagram.removeDash(span.parent.serviceName)}: Return \n`;
            }
        }
        return text;
    }

    /**
     * Adds the cell calls made for a particular trace to the diagram..
     */
    addCells() {
        this.setState({
            config: this.drawCells()
        });
        const cellArray = [];
        for (let i = 0; i < this.props.spans.length; i++) {
            if (this.props.spans[i].componentType === "Micro-service") {
                cellArray.push(this.props.spans[i]);
            }
        }
        this.setState({
            clicked: false
        });
    }

    /**
     * Gets all the cells that has been involved in the particular trace.
     *
     * @param {Array} spanArray The array containing the list of all spans.
     * @return {Array} cellArray The array containing all the cells in the trace.
     */

    static separateCells(spanArray) {
        const cellArray = [];
        for (let i = 0; i < spanArray.length; i++) {
            if ((spanArray[i].serviceName.includes(SequenceDiagram.GLOBAL))) {
                if (!cellArray.includes(SequenceDiagram.GLOBAL_GATEWAY)) {
                    cellArray.push(SequenceDiagram.GLOBAL_GATEWAY);
                }
            }
            if (spanArray[i].cell) {
                const cellName = SequenceDiagram.removeDash(spanArray[i].cell.name);
                if (!cellArray.includes(cellName)) {
                    cellArray.push(cellName);
                }
            }
        }
        return cellArray;
    }

    /**
     * Include all the cells in the trace as actors in the sequence diagram..
     *
     * @return {String} dataText The text data as string which is converted to the diagram by the mermaid library.
     */

    drawCells() {
        const array = SequenceDiagram.separateCells(this.props.spans);
        let dataText = "sequenceDiagram \n";
        for (let i = 0; i < array.length; i++) {
            dataText += `participant ${array[i]}\n`;
        }
        dataText += `activate ${SequenceDiagram.GLOBAL_GATEWAY}\n`;
        return dataText + this.addCellConnections();
    }

    /**
     * Connects all the cell communications in the diagram.
     *
     * @returns {string} dataText The text data of string type that is converted by the mermaid
     *                             library to depict the cell connections.
     */
    addCellConnections() {
        let callId = 1;
        const tree = TracingUtils.getTreeRoot(this.props.spans);
        let dataText = "";
        tree.walk((span, data) => {
            let parentCellName;
            let childCellName;
            if (span.parent !== null) {
                if (span.parent.cell === null) {
                    parentCellName = SequenceDiagram.GLOBAL_GATEWAY;
                } else {
                    parentCellName = span.parent.cell.name;
                }
                if (span.cell) {
                    parentCellName = SequenceDiagram.removeDash(parentCellName);
                    childCellName = SequenceDiagram.removeDash(span.cell.name);
                    if (parentCellName !== childCellName
                        && !span.operationName.match(Constants.System.SIDECAR_AUTH_FILTER_OPERATION_NAME_PATTERN)) {
                        span.callingId = callId;
                        dataText += `${parentCellName}->>+${childCellName}: call ${span.cell.name}-cell [${callId}] \n`;
                        callId += 1;
                    }
                }
            }
        }, undefined, (span) => {
            if (span.cell) {
                let parentCellName = "";
                if (span.parent.cell === null) {
                    parentCellName = SequenceDiagram.GLOBAL_GATEWAY;
                } else {
                    parentCellName = span.parent.cell.name;
                }
                if (span.cell.name !== parentCellName
                    && !span.operationName.match(Constants.System.SIDECAR_AUTH_FILTER_OPERATION_NAME_PATTERN)) {
                    dataText += `${SequenceDiagram.removeDash(span.cell.name)}-->>-`
                        + `${SequenceDiagram.removeDash(parentCellName)}: Return \n`;
                }
            }
        });
        dataText += `deactivate ${SequenceDiagram.GLOBAL_GATEWAY}`;
        return dataText;
    }

    /**
     * Removes dash symbol from cell/service names as the library doesn't support dashes in the actors name.
     *
     * @param {string} name The cell/service name that needs to be checked for dashes.
     * @returns {string} name The cell/service name after removing the dashes.
     */
    static removeDash(name) {
        if (name.includes("-")) {
            return name.replace(/-/g, " ");
        }
        return name;
    }

    static addDash(name) {
        if (name.includes(" ")) {
            return name.replace(" ", "-");
        }
        return name;
    }

    /**
     * Gets the index of the span object from an array by checking the span's unique id.
     *
     * @param {Array} data The array from which the index should be found.
     * @param {number[]} value The call Id of the span object.
     *
     */

    static findSpanIndexCall(data, value) {
        let isFound = false;
        return data.findIndex((item) => {
            if (item.callingId) {
                isFound = item.callingId === value[0];
            }
            return isFound;
        });
    }

}

SequenceDiagram.propTypes = {
    classes: PropTypes.any.isRequired,
    spans: PropTypes.arrayOf(
        PropTypes.instanceOf(Span).isRequired
    ).isRequired,
    colorGenerator: PropTypes.instanceOf(ColorGenerator)
};

export default withStyles(styles, {withTheme: true})(withColor(SequenceDiagram));


