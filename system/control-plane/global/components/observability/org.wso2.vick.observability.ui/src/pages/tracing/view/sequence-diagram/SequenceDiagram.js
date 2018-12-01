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

class SequenceDiagram extends Component {
    constructor(props){
        super(props);
        this.state = {
            config :"",
            heading: "Cell - level Sequence"
        }
        this.testFoo2 = this.testFoo2.bind(this);
        this.testFoo3 = this.testFoo3.bind(this);
        this.seperateCells = this.seperateCells.bind(this);
        this.drawCells = this.drawCells.bind(this);
        this.checkDraw = this.checkDraw.bind(this);


    }
    render() {

        return (
            <div>
                <h3> {this.state.heading} </h3>
            <div className="mermaid" id="mems">
                {this.state.config}
            </div>
            </div>
        );
    }

    componentDidMount(){
        let _this = this;
        _this.testFoo2();

        interact('.messageText').on('tap', function (event) {
            _this.testFoo3(event.srcElement.innerHTML);
        });
        console.log(this.props.spans);
    }
    componentDidUpdate(){
        setTimeout(()=>{
            $('#mems').removeAttr("data-processed");
        },1000);
        mermaid.init($("#mems"));
    }

    testFoo3(spanData){
        let index = findSpanIndex(this.props.spans,spanData);
        let _this = this;
        this.setState({
           config: drawSpan(this.props.spans[index]),
            heading : "Span - level Sequence"
        });
    }

    testFoo2() {
        let _this =this;
        this.setState({
            config:_this.drawCells()
        });

        var cellArray = []
        for (let i=0;i<this.props.spans.length;i++){
            if (this.props.spans[i].componentType==="Micro-service"){
                cellArray.push(this.props.spans[i]);
            }
        }
      this.checkDraw();

    }

     seperateCells(spanArray){
        var cellArray =[];
        for(let i=0;i<spanArray.length;i++){
            if((spanArray[i].serviceName.includes("global"))){
                cellArray.push("global") ;
            }
            if(Boolean(spanArray[i].cell)) {
                let cellname = spanArray[i].cell.name;
                if (cellname === "stock-options") {
                    cellname = "stock options";
                }
                if (!cellArray.includes(cellname)) {
                    cellArray.push(cellname);
                }
            }
        }

        return cellArray;
    }

    drawCells(){
        let array = this.seperateCells(this.props.spans);
        let mText ="sequenceDiagram \n";
        for (let i=0;i<array.length;i++){
            mText+= "participant "+ array[i] + "\n";
        }
        mText += "activate global\n";


        return   mText + this.checkDraw();

    }

    checkDraw() {
        let tree = TracingUtils.buildTree(this.props.spans);
        var tmp = "";
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
                        tmp += parentCellName + "->>+" + span.cell.name + ":"+ span.getUniqueId()  + "\n";
                    }
                }
            }
        }, undefined, (span) => {
            if(Boolean(span.cell)){
                let parentCellName = "";
                if (span.parent.cell === null) {
                    parentCellName = "global";
                }
                else {
                    parentCellName = span.parent.cell.name;
                }
                if (span.cell.name !== parentCellName){
                    tmp += span.cell.name + "-->>-"+ parentCellName + ": Return \n";
                }
            }
        });
        tmp += "deactivate global" ;
        return tmp;
    }

}

var spanData ="sequenceDiagram \n";

function checkDraw(span, data) {
    if (span.parentId !== "undefined") {
        if (!span.parent.cell){
            data += span.parent.cell.name + " ->> " + span.cell.name;
        }
    }
}

function removeDash(cellName){
    if(cellName.includes("-")){
        return cellName.replace("-"," ");
    }
    else {
        return cellName;
    }
}

function drawSpan(span){
        if(!Boolean(span.children)){
            alert("No further drill downs");
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
            iterateChildSpan(children);
        }
        console.log(spanData);
    return spanData;
}
function iterateChildSpan(childrenArr) {
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
                spanData += childrenArr[i].parent.spanId+" ->> " + childrenArr[i].spanId+": "+ childrenArr[i].serviceName + "\n";
            }
        }

        drawSpan(childrenArr[i]);
    }
}
function findSpanIndex(data,val){
    var index = data.findIndex(function(item){
        return item.getUniqueId() === val
    });
    return index;
}

export default SequenceDiagram;
