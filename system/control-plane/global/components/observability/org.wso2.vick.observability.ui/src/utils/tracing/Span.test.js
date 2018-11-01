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
import Span from "./Span";

describe("Span", () => {
    describe("isSiblingOf()", () => {
        it("should return true if the sibling span is provided", () => {
            const clientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const serverSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10010,
                duration: 50,
                tags: {}
            });

            expect(clientSpan.isSiblingOf(serverSpan)).toBe(true);
            expect(serverSpan.isSiblingOf(clientSpan)).toBe(true);
        });

        it("should return false if a span from another trace is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const traceIdMismatchesSpan = new Span({
                traceId: "trace-b-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.isSiblingOf(traceIdMismatchesSpan)).toBe(false);
        });

        it("should return false if a non equal span ID is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const spanIdMismatchedSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentSpanId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10010,
                duration: 50,
                tags: {}
            });

            expect(span.isSiblingOf(spanIdMismatchedSpan)).toBe(false);
        });

        it(`should return false if this span's or the provided span's kind is ${Constants.Span.Kind.CONSUMER}, `
                + `${Constants.Span.Kind.PRODUCER} or empty`, () => {
            const clientKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const serverKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const producerKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "producer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.PRODUCER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const consumerKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CONSUMER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const emptyKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "producer-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            // Spans belonging to kinds other than client and server
            const kindSpansList = [clientKindSpan, serverKindSpan, producerKindSpan, consumerKindSpan, emptyKindSpan];
            const nonSiblingSpansList = [producerKindSpan, consumerKindSpan, emptyKindSpan];
            for (let i = 0; i < nonSiblingSpansList.length; i++) {
                for (let j = 0; j < kindSpansList.length; j++) {
                    expect(nonSiblingSpansList[i].isSiblingOf(kindSpansList[j])).toBe(false);
                }
            }
        });

        it(`should return false if this span and the provided span is of the same type; `
                + `${Constants.Span.Kind.CLIENT}/${Constants.Span.Kind.SERVER}`, () => {
            const clientKindSpanA = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const clientKindSpanB = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const serverKindSpanA = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const serverKindSpanB = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            // Same kind siblings of type client or server
            expect(clientKindSpanA.isSiblingOf(clientKindSpanB)).toBe(false);
            expect(serverKindSpanA.isSiblingOf(serverKindSpanB)).toBe(false);
        });

        it("should return false if null/undefine is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.isSiblingOf(null)).toBe(false);
            expect(span.isSiblingOf(undefined)).toBe(false);
        });
    });

    describe("isParentOf()", () => {
        it("should return true if the direct parent span is provided", () => {
            const childSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.PRODUCER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const parentSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentSpanId: "span-c-id",
                serviceName: "parent-service",
                operationName: "get-resource",
                startTime: 10050,
                duration: 145,
                tags: {}
            });

            expect(parentSpan.isParentOf(childSpan)).toBe(true);
        });

        it("should return true if the sibling server span is provided", () => {
            const siblingServerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "sibling-span-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const siblingClientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "sibling-span-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10050,
                duration: 153,
                tags: {}
            });

            expect(siblingClientSpan.isParentOf(siblingServerSpan)).toBe(true);
        });

        it("should return false if the sibling client span is provided", () => {
            const siblingServerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "sibling-span-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const siblingClientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "sibling-span-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10050,
                duration: 153,
                tags: {}
            });

            expect(siblingServerSpan.isParentOf(siblingClientSpan)).toBe(false);
        });

        it("should return false if a span from another trace is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const mismatchedTraceIdSpan = new Span({
                traceId: "trace-b-id",
                spanId: "span-b-id",
                parentSpanId: "span-c-id",
                serviceName: "parent-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(mismatchedTraceIdSpan.isParentOf(span)).toBe(false);
        });

        it("should return false if a non-related span is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const nonRelatedSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-c-id",
                parentSpanId: "span-d-id",
                serviceName: "non-parent-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(nonRelatedSpan.isParentOf(span)).toBe(false);
        });

        it("should return false if null/undefined is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.isParentOf(null)).toBe(false);
            expect(span.isSiblingOf(undefined)).toBe(false);
        });
    });

    describe("addSpanReference()", () => {
        it("should add as child and return true if the child is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CONSUMER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const childSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-x-id",
                parentSpanId: "span-a-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.addSpanReference(childSpan)).toBe(true);
            expect(span.children.has(childSpan)).toBe(true);
            expect(span.parent).not.toBe(childSpan);
            expect(span.sibling).not.toBe(childSpan);
        });

        it("should add as parent and return true if the parent is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-x-id",
                parentSpanId: "span-a-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.PRODUCER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const parentSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.addSpanReference(parentSpan)).toBe(true);
            expect(span.children.has(parentSpan)).toBe(false);
            expect(span.parent).toBe(parentSpan);
            expect(span.sibling).not.toBe(parentSpan);
        });

        it("should add as child and sibling and return true if the sibling server span is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const siblingServerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.addSpanReference(siblingServerSpan)).toBe(true);
            expect(span.children.has(siblingServerSpan)).toBe(true);
            expect(span.parent).not.toBe(siblingServerSpan);
            expect(span.sibling).toBe(siblingServerSpan);
        });

        it("should add as parent and sibling and return true if the sibling client span is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const siblingClientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-b-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.addSpanReference(siblingClientSpan)).toBe(true);
            expect(span.children.has(siblingClientSpan)).toBe(false);
            expect(span.parent).toBe(siblingClientSpan);
            expect(span.sibling).toBe(siblingClientSpan);
        });

        it("should not add as child/parent/sibling and return false if non-related span is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-x-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const nonRelatedSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentSpanId: "span-y-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.addSpanReference(nonRelatedSpan)).toBe(false);
            expect(span.children.has(nonRelatedSpan)).toBe(false);
            expect(span.parent).not.toBe(nonRelatedSpan);
            expect(span.sibling).not.toBe(nonRelatedSpan);
        });

        it("should not add as child/parent/sibling and return false if null is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });

            expect(span.addSpanReference(null)).toBe(false);
            expect(span.children.size).toBe(0);
            expect(span.parent).toBeNull();
            expect(span.sibling).toBeNull();

            expect(span.addSpanReference(undefined)).toBe(false);
            expect(span.children.size).toBe(0);
            expect(span.parent).toBeNull();
            expect(span.sibling).toBeNull();
        });
    });

    describe("getUniqueId()", () => {
        it("should return a unique ID across traces", () => {
            const span = new Span({
                traceId: "trace-id",
                spanId: "span-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const siblingSpan = new Span({
                traceId: "trace-id",
                spanId: "span-server-id",
                parentSpanId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const parentSpan = new Span({
                traceId: "trace-id",
                spanId: "span-parent-id",
                parentSpanId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: {}
            });
            const differentTraceSpan = new Span({
                traceId: "trace-different-id",
                spanId: "span-id",
                parentSpanId: "span-parent-id",
                serviceName: "test-service",
                operationName: "set-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 200000,
                duration: 1203,
                tags: {}
            });

            expect(span.getUniqueId()).not.toBe(siblingSpan.getUniqueId());
            expect(span.getUniqueId()).not.toBe(parentSpan.getUniqueId());
            expect(span.getUniqueId()).not.toBe(differentTraceSpan.getUniqueId());
        });
    });
});
