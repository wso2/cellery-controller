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

import Constants from "./Constants";

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
        if (!(spanData && spanData.traceId && spanData.spanId)) {
            throw Error("Invalid span: Expected Trace Id and Span Id but not found");
        }

        this.traceId = spanData.traceId;
        this.spanId = spanData.spanId;
        this.parentSpanId = spanData.parentSpanId;
        this.serviceName = spanData.serviceName;
        this.operationName = spanData.operationName;
        this.kind = (spanData.kind ? spanData.kind.toUpperCase() : null);
        this.startTime = spanData.startTime ? spanData.startTime : 0;
        this.duration = spanData.duration ? spanData.duration : 0;
        this.tags = spanData.tags ? spanData.tags : {};

        /** @type {Span} **/
        this.parent = null;

        /** @type {Set.<Span>} **/
        this.children = new Set();

        /** @type {Span} **/
        this.sibling = null;

        /** @type {boolean} **/
        this.isSystemComponent = false;

        /** @type {{name: string, version: string}} **/
        this.cell = null;
    }

    /**
     * Check if another span is a sibling of this span.
     *
     * @param {Span} span The span to check if it is a sibling
     * @returns {boolean} True if this is a sibling of the other span
     */
    isSiblingOf(span) {
        return Boolean(span) && this.traceId === span.traceId && this.spanId === span.spanId
            && ((this.kind === Constants.Span.Kind.CLIENT && span.kind === Constants.Span.Kind.SERVER)
            || (this.kind === Constants.Span.Kind.SERVER && span.kind === Constants.Span.Kind.CLIENT));
    }

    /**
     * Check if this is the parent of another span.
     *
     * @param {Span} span The span to check if it is a child
     * @returns {boolean} True if this is the parent of the other span
     */
    isParentOf(span) {
        let isParentOfSpan = false;
        if (Boolean(span) && this.traceId === span.traceId) {
            if (this.spanId === span.spanId && this.kind === Constants.Span.Kind.CLIENT
                    && span.kind === Constants.Span.Kind.SERVER) { // Siblings
                isParentOfSpan = true;
            } else if (this.spanId === span.parentSpanId && this.kind !== Constants.Span.Kind.CLIENT
                    && span.kind !== Constants.Span.Kind.SERVER) {
                isParentOfSpan = true;
            }
        }
        return isParentOfSpan;
    }

    /**
     * Add a reference to another span in this span.
     * Only child, parent and sibling spans are added as references.
     *
     * @param {Span} span The to which the reference should be added
     * @returns {boolean} True if the span was added as a reference
     */
    addSpanReference(span) {
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
    }

    /**
     * Get a unique ID to represent this span.
     *
     * @returns {string} the unique ID to represent this span
     */
    getUniqueId() {
        return `${this.traceId}--${this.spanId}--${this.kind}`;
    }

}

export default Span;
