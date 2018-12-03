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
        }
        this.testFoo2 = this.testFoo2.bind(this);
        this.testFoo3 = this.testFoo3.bind(this);
        this.seperateCells = this.seperateCells.bind(this);
        this.drawCells = this.drawCells.bind(this);
        this.combine = this.combine.bind(this);
        this.drawSpan = this .drawSpan.bind(this);
        this.iterateChildSpan = this.iterateChildSpan.bind(this);

        var spanData = "";

    }
    render() {

        return (
            <div>
                <h3> {this.state.heading} </h3>
                <div className="mermaid" id="mems">
                {this.state.config}
            </div>
                <div id="hiddenDiv" style={this.state.clicked ? {} : {display:'none'}} onClick={this.testFoo2}> &lt;&lt; Back to Cell-level Diagram</div>
            </div>
        );
    }

    componentDidMount(){
        let _this = this;
        _this.testFoo2();

        interact('.messageText').on('tap', function (event) {
            if ((event.srcElement.innerHTML !== "Return") && (_this.state.clicked !== true)){
                _this.testFoo3(event.srcElement.innerHTML);
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
              $('#mems').removeAttr("data-processed");
          mermaid.init($("#mems"));
      }

    }

    testFoo3(span){
        let index = findSpanIndex(this.props.spans,span);
       // let _this = this;
        let data2 = this.combine(this.props.spans[index]);


            console.log(data2);


        this.setState({
           config: data2,
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


        _this.setState({
            clicked:false
        });

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


        return   mText + this.drawCells2();

    }

    drawCells2() {
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

    drawSpan(span, tmp){
        if(span.children === undefined){
            console.log("dsdsfdsf")
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
            this.iterateChildSpan(children, tmp);
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

                        var jsonObj = {
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

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var spanData ="sequenceDiagram \n";
var copyArr =[];

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

// function drawSpan(span){
//         if(span.children === null){
//           console.log("dsdsfdsf")
//         }
//         else {
//             const children = [];
//             console.log(span);
//             const childrenIterator = span.children.values();
//             let currentChild = childrenIterator.next();
//             while (!currentChild.done) {
//                 children.push(currentChild.value);
//                 currentChild = childrenIterator.next();
//             }
//             iterateChildSpan(children);
//         }
//
//     return spanData;
// }
// function iterateChildSpan(childrenArr) {
//     for(let i=0;i<childrenArr.length;i++){
//
//         let parentCellName;
//         if (Boolean(childrenArr[i].cell)) {
//             if (childrenArr[i].parent.cell === null) {
//                 parentCellName = "global";
//             } else {
//                 parentCellName = childrenArr[i].parent.cell.name;
//             }
//             if (childrenArr[i].cell.name !== parentCellName) {
//                 console.log("hit here");
//                 break;
//             }
//             else {
//                 if (childrenArr[i].parent.kind !== Constants.Span.Kind.CLIENT) {
//                     spanData += childrenArr[i].parent.spanId + " ->>+ " + childrenArr[i].spanId + ": " + childrenArr[i].serviceName + "\n";
//                     var jsonObj = {
//                         "from": childrenArr[i].spanId,
//                         "to": childrenArr[i].parent.spanId
//                     };
//                     copyArr.push(jsonObj);
//                 }
//             }
//         }
//
//         drawSpan(childrenArr[i]);
//     }
//
// }
function findSpanIndex(data,val){
    var index = data.findIndex(function(item){
        return item.getUniqueId() === val
    });
    return index;
}

export default SequenceDiagram;
