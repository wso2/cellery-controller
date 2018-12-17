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
import "./SequenceStyles.css";
import Button from "@material-ui/core/Button";
import PropTypes from "prop-types";
import React from "react";
import Span from "../../utils/span";
import TracingUtils from "../../utils/tracingUtils";
import interact from "interactjs";
import mermaid from "mermaid";

class SequenceDiagram extends React.Component {

    static GLOBAL = "global";

    constructor(props) {
        super(props);
        this.state = {
            config: "",
            heading: "Cell - Level Sequence",
            spanData: "sequenceDiagram \n",
            copyArr: [],
            clicked: false,
            cellName: null,
            clonedArray: []
        };

        this.mermaidDivRef = React.createRef();

        this.addCells = this.addCells.bind(this);
        this.addServices = this.addServices.bind(this);
        this.drawCells = this.drawCells.bind(this);
    }

    render() {
        return (

            <div>

                <h3> {this.state.heading} </h3>
                <Button color="primary" style={this.state.clicked ? {} : {display: "none"}} onClick={this.addCells}>
                    &lt;&lt; Back to Cell-level Diagram
                </Button>
                <div>{this.state.cellName}</div>
                <div className="mermaid" id="mermaidDiv" ref={this.mermaidDivRef}>
                    {this.state.config}
                </div>
            </div>
        );
    }

    componentDidMount() {
        this.addCells();
        interact(".messageText").on("tap", (event) => {
            if ((event.srcElement.innerHTML !== "Return" && event.srcElement.innerHTML !== "Calls")
                && (this.state.clicked !== true)) {
                const numb = event.srcElement.innerHTML.match(/\d+/g).map(Number);
                this.addServices(numb);
                this.setState({
                    clicked: true
                });
            }
        });
        this.cloneArray();
    }

    componentDidUpdate(prevProps, prevState) {
        const collections = this.mermaidDivRef.current.getElementsByClassName("messageText");
        for (let i = 0; i < collections.length; i++) {
            if (collections[i].innerHTML.includes("[")) {
                collections[i].classList.add("newMessageText");
            }
        }

        if (this.state.config !== prevState.config) {
            this.mermaidDivRef.current.removeAttribute("data-processed");
            mermaid.init(this.mermaidDivRef.current);

            for (let i = 0; i < collections.length; i++) {
                if (collections[i].innerHTML.match("\\s\\[([0-9]+)\\]+$")) {
                    collections[i].classList.add("newMessageText");
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
        data2 += `activate ${SequenceDiagram.removeDash(treeRoot.serviceName)}\n`;
        treeRoot.walk(
            (span) => {
                if (!span.isFromIstioSystemComponent() && !span.isFromVICKSystemComponent()) {
                    data2 += SequenceDiagram.updateDataText(span, parentName);
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
            heading: "Service - Level Sequence"
        });
    }

    /**
     * Updates the text data, which is used by the mermaid library to generate diagrams.
     *
     * @param {Span} span The span to be checked.
     * @param {String} parentName The parent cell name
     * @return {String} text The updated text
     */

    static updateDataText(span, parentName) {
        let text = "";
        if (!span.callingId && parentName === span.cell.name) {
            if (span.parent.serviceName !== span.serviceName) {
                text += `${SequenceDiagram.removeDash(span.parent.serviceName)}  ->>+`
                    + `${SequenceDiagram.removeDash(span.serviceName)}:${span.operationName}\n`;
            }
        }
        return text;
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
            clicked: false,
            heading: "Cell - Level Sequence"
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
                cellArray.push(SequenceDiagram.GLOBAL);
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
        dataText += `activate ${SequenceDiagram.GLOBAL}\n`;
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
            if (span.parentId !== "undefined") {
                if (span.parent.cell === null) {
                    parentCellName = SequenceDiagram.GLOBAL;
                } else {
                    parentCellName = span.parent.cell.name;
                }
                if (span.cell) {
                    parentCellName = SequenceDiagram.removeDash(parentCellName);
                    span.cell.name = SequenceDiagram.removeDash(span.cell.name);
                    if (parentCellName !== span.cell.name) {
                        span.callingId = callId;
                        dataText += `${parentCellName}->>+${span.cell.name}:${span.serviceName} [${callId}] \n`;
                        callId += 1;
                    }
                }
            }
        }, undefined, (span) => {
            if (span.cell) {
                let parentCellName = "";
                if (span.parent.cell === null) {
                    parentCellName = SequenceDiagram.GLOBAL;
                } else {
                    parentCellName = span.parent.cell.name;
                }
                if (span.cell.name !== parentCellName) {
                    dataText += `${span.cell.name}-->>-${parentCellName}: Return \n`;
                }
            }
        });
        dataText += `deactivate ${SequenceDiagram.GLOBAL}`;
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
    spans: PropTypes.arrayOf(
        PropTypes.instanceOf(Span).isRequired
    ).isRequired
};

export default SequenceDiagram;


