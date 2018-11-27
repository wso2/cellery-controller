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


import React, { Component } from 'react';
import  ReactDOM from "react-dom";
import mermaid, {mermaidAPI} from 'mermaid';
import $ from 'jquery';
import TracingUtils from "../../utils/tracingUtils";

class SequenceDiagram extends Component {
    constructor(props){
        super(props);
        this.state = {
            config :""
        }
        this.testFoo2 = this.testFoo2.bind(this);
        this.seperateCells = this.seperateCells.bind(this);
        this.drawCells = this.drawCells.bind(this);

    }
    render() {

        return (
            <div className="mermaid" id="mems">
                {this.state.config}
            </div>
        );
    }

    componentDidMount(){
        let _this = this;
        _this.testFoo2();

    }
    componentDidUpdate(){
        mermaid.init(undefined, $("#mems"));
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
        console.log(this.props.spans);
    }

     seperateCells(spanArray){
        var cellArray =[];
        for(let i=0;i<spanArray.length;i++){
            let  cellname = spanArray[i].serviceName.split("-")[0];
            if(!cellArray.includes(cellname)){
                cellArray.push(cellname);
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

        return mText;
    }
}

export default SequenceDiagram;
