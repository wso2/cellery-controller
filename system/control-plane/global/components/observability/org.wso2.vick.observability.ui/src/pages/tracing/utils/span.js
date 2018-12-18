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

import Constants from "../../common/constants";

/**
 * Single span in a Trace.
 */
class Span {


    /**
     * Span constructor.
     *
     * @param {Object} spanData Span data object
     */
    constructor(spanData) {
        this.traceId = spanData.traceId;
        this.spanId = spanData.spanId;
        this.parentId = spanData.parentId;
        this.serviceName = spanData.serviceName;
        this.operationName = spanData.operationName;
        this.kind = (spanData.kind ? spanData.kind.toUpperCase() : null);
        this.startTime = spanData.startTime ? spanData.startTime : 0;
        this.duration = spanData.duration ? spanData.duration : 0;
        this.tags = spanData.tags ? JSON.parse(spanData.tags) : {};

        // Adding additional tags
        const addTagFromSpanData = (tagKey, spanDataKey) => {
            if (spanData[spanDataKey]) {
                this.tags[tagKey] = spanData[[spanDataKey]];
            }
        };
        addTagFromSpanData("k8s.pod", "pod");
        addTagFromSpanData("k8s.namespace", "namespace");

        /** @type {string} **/
        this.componentType = "";

        /** @type {{name: string, version: string}} **/
        this.cell = (spanData.cell ? {name: spanData.cell, version: null} : null);

        /** @type {Span} **/
        this.parent = null;

        /** @type {Set.<Span>} **/
        this.children = new Set();

        /** @type {Span} **/
        this.sibling = null;

        this.treeDepth = null;
    }

    /**
     * Check if another span is a sibling of this span.
     *
     * @param {Span} span The span to check if it is a sibling
     * @returns {boolean} True if this is a sibling of the other span
     */
    isSiblingOf = (span) => Boolean(span) && this.traceId === span.traceId && this.spanId === span.spanId
        && ((this.kind === Constants.Span.Kind.CLIENT && span.kind === Constants.Span.Kind.SERVER)
            || (this.kind === Constants.Span.Kind.SERVER && span.kind === Constants.Span.Kind.CLIENT));

    /**
     * Check if this is the parent of another span.
     *
     * @param {Span} span The span to check if it is a child
     * @returns {boolean} True if this is the parent of the other span
     */
    isParentOf = (span) => {
        let isParentOfSpan = false;
        if (Boolean(span) && this.traceId === span.traceId) {
            if (this.spanId === span.spanId && this.kind === Constants.Span.Kind.CLIENT
                    && span.kind === Constants.Span.Kind.SERVER) { // Siblings
                isParentOfSpan = true;
            } else if (this.spanId === span.parentId) {
                isParentOfSpan = true;
                if (this.hasSibling()) {
                    isParentOfSpan = isParentOfSpan && this.kind === Constants.Span.Kind.SERVER;
                }
                if (span.hasSibling()) {
                    isParentOfSpan = isParentOfSpan && span.kind === Constants.Span.Kind.CLIENT;
                }
            }
        }
        return isParentOfSpan;
    };

    /**
     * Check if this span has a sibling.
     *
     * @returns {boolean} True if this span has a sibling
     */
    hasSibling = () => this.kind === Constants.Span.Kind.CLIENT || this.kind === Constants.Span.Kind.SERVER;

    /**
     * Add a reference to another span in this span.
     * Only child, parent and sibling spans are added as references.
     *
     * @param {Span} span The to which the reference should be added
     * @returns {boolean} True if the span was added as a reference
     */
    addSpanReference = (span) => {
        let spanAdded = false;
        if (this.isParentOf(span)) {
            this.children.add(span);
            spanAdded = true;
        } else if (Boolean(span) && span.isParentOf(this)) {
            this.parent = span;
            spanAdded = true;
        }
        if (this.isSiblingOf(span)) {
            this.sibling = span;
            spanAdded = true;
        }
        return spanAdded;
    };

    /**
     * Reset all references to spans.
     */
    resetSpanReferences = () => {
        this.children.clear();
        this.parent = null;
        this.sibling = null;
        this.treeDepth = 0;
    };

    /**
     * Walk down the trace tree starting from this span in DFS manner.
     * When a node has multiple children they will be traveled in the order of their start time.
     *
     * @param {function} nodeCallBack The callback to be called in each node.
     *                                The function should return the data that should be passed down to the children.
     * @param {Object} data The initial data to be passed down the trace tree
     * @param {function} postTraverseCallBack The callback to be called after traversing a node.
     * @param {function} shouldTerminate
     */

    walk(nodeCallBack, data = {}, postTraverseCallBack = null, shouldTerminate = null) {
        if (shouldTerminate && shouldTerminate(this)) {
            return;
        }


        let newData;
        if (nodeCallBack) {
            newData = nodeCallBack(this, data);
        }

        // Get the list of children of this node
        const children = [];
        const childrenIterator = this.children.values();
        let currentChild = childrenIterator.next();
        while (!currentChild.done) {
            children.push(currentChild.value);
            currentChild = childrenIterator.next();
        }

        // Sorting by start time
        children.sort((a, b) => a.startTime - b.startTime);

        // Traversing down the tree structure
        for (let i = 0; i < children.length; i++) {
            children[i].walk(nodeCallBack, newData, postTraverseCallBack, shouldTerminate);
        }
        if (postTraverseCallBack) {
            postTraverseCallBack(this);
        }
    }

    /**
     * Get a unique ID to represent this span.
     *
     * @returns {string} the unique ID to represent this span
     */
    getUniqueId = () => `${this.traceId}--${this.spanId}${this.kind ? `--${this.kind}` : ""}`;

    /**
     * Check whether a span belongs to the cell gateway.
     *
     * @returns {boolean} True if the component to which the span belongs to is a cell gateway
     */
    isFromCellGateway = () => Constants.Cell.GATEWAY_NAME_PATTERN.test(this.serviceName);

    /**
     * Check whether a span belongs to the Istio System.
     *
     * @returns {boolean} True if the component to which the span belongs to is a system component
     */
    isFromIstioSystemComponent = () => Constants.System.ISTIO_MIXER_NAME_PATTERN.test(this.serviceName);

    /**
     * Check whether a span belongs to the side car.
     *
     * @returns {boolean} True if the span is from side car
     */
    isFromSideCar = () => this.tags.component === "proxy";

    /**
     * Check whether a span belongs to the VICK System.
     *
     * @returns {boolean} True if the component to which the span belongs to is a system component
     */
    isFromVICKSystemComponent = () => (
        this.isFromCellGateway() || Constants.System.GLOBAL_GATEWAY_NAME_PATTERN.test(this.serviceName)
            || Constants.System.SIDECAR_AUTH_FILTER_OPERATION_NAME_PATTERN.test(this.operationName)
    );

    /**
     * Check whether an error occurred during this span.
     *
     * @returns {boolean} True if an error had occurred in this span
     */
    hasError = () => this.tags.error === "true";

    /**
     * Create a shallow clone.
     * This will create a clone without the span references.
     *
     * @returns {Span} The cloned span
     */
    shallowClone = () => {
        const span = new Span({
            traceId: this.traceId,
            spanId: this.spanId,
            parentId: this.parentId,
            serviceName: this.serviceName,
            operationName: this.operationName,
            kind: this.kind,
            startTime: this.startTime,
            duration: this.duration
        });
        span.tags = {...this.tags};
        span.componentType = this.componentType;
        span.cell = {...this.cell};
        return span;
    };

}

export default Span;
