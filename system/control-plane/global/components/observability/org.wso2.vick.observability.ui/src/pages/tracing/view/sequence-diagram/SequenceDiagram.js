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

import interact from "interactjs";
import React, { Component } from 'react';
import  ReactDOM from "react-dom";
import mermaid, {mermaidAPI} from 'mermaid';
import $ from 'jquery';
import TracingUtils from "../../utils/tracingUtils";
import Span from "../../utils/span";
import Constants from "../../utils/constants";
import './SequenceStyles.css'


class SequenceDiagram extends Component {
    constructor(props){
        super(props);
        this.state = {
            config :"",
            heading: "Cell - level Sequence",
            spanData:"sequenceDiagram \n",
            copyArr: [],
            clicked: false
        };
        this.addCells = this.addCells.bind(this);
        this.addServices = this.addServices.bind(this);
        this.seperateCells = this.seperateCells.bind(this);
        this.drawCells = this.drawCells.bind(this);
        this.combine = this.combine.bind(this);
        this.drawSpan = this .drawSpan.bind(this);
        this.iterateChildSpan = this.iterateChildSpan.bind(this);

    }
    render() {
        return (

            <div>

                <h3> {this.state.heading} </h3>
                <div id="hiddenDiv" style={this.state.clicked ? {} : {display:'none'}} onClick={this.addCells}> &lt;&lt; Back to Cell-level Diagram</div>
                <div className="mermaid" id="mems">
                {this.state.config}
            </div>
            </div>
        );
    }

    componentDidMount(){
        let _this = this;
        _this.addCells();

        interact('.messageText').on('tap', function (event) {
            if ((event.srcElement.innerHTML !== "Return") && (_this.state.clicked !== true)){
                let numb = event.srcElement.innerHTML.match(/\d+/g).map(Number);
                _this.addServices(numb);
                _this.setState({
                    clicked:true
                });

            }

        });
        console.log(this.props.spans);
        console.log(this.props.clicked);
    }
    componentDidUpdate(prevProps, prevState) {
      if (this.state.config !== prevState.config) {
          document.getElementById("mems").removeAttribute("data-processed");
          mermaid.init($("#mems"));
      }

    }

    /**
     * Adds the service calls made for a particular cell to the diagram.
     *
     * @param {number[]} callId The span's call Id of the particular cell call.
     */

    addServices(callId){
       let spanArray = this.tryGetSpan(callId);
        console.log(spanArray);
        let cellArray =[];
        for(let i=0;i<spanArray.length;i++){
                let cellname = removeDash(spanArray[i].serviceName);
                if (!cellArray.includes(cellname)) {
                    cellArray.push(cellname);
                }
        }

        let data2 ="sequenceDiagram \n";
        for (let i=1;i<spanArray.length;i++){
            data2+= removeDash(spanArray[i-1].serviceName) +" ->> "+ removeDash(spanArray[i].serviceName)+ ": Calls \n";
        }
        console.log(data2);
        this.setState({
            config: data2,
            heading : "Span - level Sequence"
        });
    }

    /**
     * Adds the cell calls made for a particular trace to the diagram..
     */
    addCells() {
        let _this =this;
        this.setState({
            config:_this.drawCells()
        });
        let cellArray = []
        for (let i=0;i<this.props.spans.length;i++){
            if (this.props.spans[i].componentType==="Micro-service"){
                cellArray.push(this.props.spans[i]);
            }
        }
        _this.setState({
            clicked:false,
            heading: "Cell - level Sequence"
        });
    }

    /**
     * Gets all the cells that has been involved in the particular trace.
     *
     * @param {Array} spanArray The array containing the list of all spans.
     * @return {Array} cellArray The array containing all the cells in the trace.
     */

     seperateCells(spanArray){
        let cellArray =[];
        for(let i=0;i<spanArray.length;i++){
            if((spanArray[i].serviceName.includes("global"))){
                cellArray.push("global") ;
            }
            if(Boolean(spanArray[i].cell)) {
                let cellname = removeDash(spanArray[i].cell.name);
                if (!cellArray.includes(cellname)) {
                    cellArray.push(cellname);
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

    drawCells(){
        let array = this.seperateCells(this.props.spans);
        let dataText ="sequenceDiagram \n";
        for (let i=0;i<array.length;i++){
            dataText+= "participant "+ array[i] + "\n";
        }
        dataText += "activate global\n";
        return   dataText + this.addCellConnections();
    }

    /**
     * Connects all the cell communications in the diagram.
     *
     * @return {String} dataText The text data of string type that is converted by the mermaid library to depict the cell connections.
     */
    addCellConnections() {
        let callId = 1;
        let serviceListArr = []
        let tree = TracingUtils.buildTree(this.props.spans);
        let dataText = "";
        tree.walk((span, data) => {
            let parentCellName;
            if (span.parentId !== "undefined") {
                if (span.parent.cell === null) {
                    parentCellName = "global";
                }
                else {
                    parentCellName = span.parent.cell.name;
                }
                if (Boolean(span.cell)) {
                    parentCellName = removeDash(parentCellName);
                    span.cell.name = removeDash(span.cell.name);
                    if (parentCellName !== span.cell.name) {
                        span.callingId = callId;
                        dataText += parentCellName + "->>+" + span.cell.name + ":"+ span.serviceName+" ["+ callId + "] \n";
                        callId +=1;
                    }
                    else {
                        serviceListArr.push(span);
                    }
                }
            }
        }, undefined, (span) => {
            if (Boolean(span.cell)) {
                let parentCellName = "";
                if (span.parent.cell === null) {
                    parentCellName = "global";
                }
                else {
                    parentCellName = span.parent.cell.name;
                }
                if (span.cell.name !== parentCellName) {
                    dataText += span.cell.name + "-->>-" + parentCellName + ": Return \n";
                }
            }
        });
        dataText += "deactivate global";
        return dataText;
    }

    /**
     * Gets all the cells that has been involved in the particular trace.
     *
     * @param {Array} spanArray The array containing the list of all spans.
     */
    drawSpan(span, strData){
        if(span.children === undefined){
            console.log("Test");
        }
        else {
            const children = [];
            console.log(span);
            const childrenIterator = span.children.values();
            let currentChild = childrenIterator.next();
            while (!currentChild.done) {
                children.push(currentChild.value);
                currentChild = childrenIterator.next();
            }
            this.iterateChildSpan(children, strData);
        }

    }

    iterateChildSpan(childrenArr,str) {
        for(let i=0;i<childrenArr.length;i++){

            let parentCellName;
            if (Boolean(childrenArr[i].cell)) {
                if (childrenArr[i].parent.cell === null) {
                    parentCellName = "global";
                } else {
                    parentCellName = childrenArr[i].parent.cell.name;
                }
                if (childrenArr[i].cell.name !== parentCellName) {
                    console.log("hit here");
                    break;
                }
                else {
                    if (childrenArr[i].parent.kind !== Constants.Span.Kind.CLIENT) {
                        str.data += childrenArr[i].parent.spanId + " ->>+ " + childrenArr[i].spanId + ": " + childrenArr[i].serviceName + "\n";
                        let jsonObj = {
                            "from": childrenArr[i].spanId,
                            "to": childrenArr[i].parent.spanId
                        };
                        str.copyArr1.push(jsonObj);
                    }
                }
            }

            this.drawSpan(childrenArr[i],str);
        }

    }

    combine(span){
        let tmp = { data : "sequenceDiagram \n",copyArr1:[]};
        this.drawSpan(span,tmp);
        console.log(tmp);

        tmp.copyArr1 = tmp.copyArr1.reverse();
        for (let k = 0; k < tmp.copyArr1.length; k++) {
            tmp.data += tmp.copyArr1[k].from +  " -->>- "+ tmp.copyArr1[k].to + ": Finish Span \n";
        }
        console.log(tmp.data);
        return tmp.data;
    }

    /**
     * Get the list of spans of the services processed inside a cell.
     *
     * @param {number[]} spanId The unique Id of the span which made the call to the current cell.
     * @param {Array} arr Iterate through this array to get the spans.
     * @return {Array} spanList The array containing all the spans of the services processed inside a cell.
     */
    getServices(spanId,arr){
        let spanList =[];
        console.log(spanId);
        let span  = arr[findSpanIndexCall(arr,spanId)];
        console.log(findSpanIndexCall(arr,spanId));
        console.log(span);
        let cellName = span.getCell().name;
        console.log(cellName);
        for (let i = findSpanIndexCall(arr,spanId);i<arr.length;i++){
            if(!arr[i].isFromIstioSystemComponent() && !arr[i].isFromVICKSystemComponent()) {
                if (cellName === arr[i].getCell().name) {
                    spanList.push(arr[i]);
                }
                else {
                    break;
                }
            }
        }
        return spanList;
    }


    tryGetSpan(spanCallId){
        let serviceListArr = [];
        let treeRoot = this.props.spans[findSpanIndexCall(this.props.spans,spanCallId)];
        let parentName = treeRoot.cell.name;
        treeRoot.walk(
            (span) => {
                if(!Boolean(span.callingId) && parentName === span.parent.cell.name){
                    console.log(span.cell.name);
                    serviceListArr.push(span);
                }
            }, null,
            null,
            (span) => (!Boolean(span.callingId) && parentName !== span.parent.cell.name)
        );

        return serviceListArr;
    }

}
export default SequenceDiagram;

/**
 * Removes dash symbol from cell/service names as the library doesn't support dashes in the actors name.
 *
 * @param {String} name The cell/service name that needs to be checked for dashes.
 * @return {String} name The cell/service name after removing the dashes.
 */
function removeDash(name){
    if(name.includes("-")){
        return name.replace(/-/g," ");
    }
    else {
        return name;
    }
}

/**
 * Gets the index of the span object from an array by checking the span's unique id.
 *
 * @param {Array} data The array from which the index should be found.
 * @return {number} index The index of the span.
 */

function findSpanIndex(data,value){
    let index = data.findIndex(function(item){
        return item.getUniqueId() === value
    });
    return index;
}

function findSpanIndexCall(data,value){
    let index = data.findIndex(function(item){
        if (Boolean(item.callingId)) {
            return item.callingId === value[0]
        }
    });

    return index;
}

