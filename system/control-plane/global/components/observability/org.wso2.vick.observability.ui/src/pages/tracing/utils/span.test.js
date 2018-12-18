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

/* eslint max-lines: ["off"] */

import Constants from "../../common/constants";
import Span from "./span";

describe("Span", () => {
    describe("isSiblingOf()", () => {
        it("should return true if the sibling span is provided", () => {
            const clientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const serverSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10010,
                duration: 50,
                tags: "{}"
            });

            expect(clientSpan.isSiblingOf(serverSpan)).toBe(true);
            expect(serverSpan.isSiblingOf(clientSpan)).toBe(true);
        });

        it("should return false if a span from another trace is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const traceIdMismatchesSpan = new Span({
                traceId: "trace-b-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });

            expect(span.isSiblingOf(traceIdMismatchesSpan)).toBe(false);
        });

        it("should return false if a non equal span ID is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const spanIdMismatchedSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10010,
                duration: 50,
                tags: "{}"
            });

            expect(span.isSiblingOf(spanIdMismatchedSpan)).toBe(false);
        });

        it(`should return false if this span's or the provided span's kind is ${Constants.Span.Kind.CONSUMER}, `
                + `${Constants.Span.Kind.PRODUCER} or empty`, () => {
            const clientKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const serverKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const producerKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "producer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.PRODUCER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const consumerKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CONSUMER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const emptyKindSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "producer-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
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

        it("should return false if this span and the provided span is of the same type; "
                + `${Constants.Span.Kind.CLIENT}/${Constants.Span.Kind.SERVER}`, () => {
            const clientKindSpanA = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const clientKindSpanB = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const serverKindSpanA = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const serverKindSpanB = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });

            // Same kind siblings of type client or server
            expect(clientKindSpanA.isSiblingOf(clientKindSpanB)).toBe(false);
            expect(serverKindSpanA.isSiblingOf(serverKindSpanB)).toBe(false);
        });

        it("should return false if null/undefine is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
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
                parentId: "span-b-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.PRODUCER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const parentSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentId: "span-c-id",
                serviceName: "parent-service",
                operationName: "get-resource",
                startTime: 10050,
                duration: 145,
                tags: "{}"
            });

            expect(parentSpan.isParentOf(childSpan)).toBe(true);
        });

        it("should return true if the sibling server span is provided", () => {
            const siblingServerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "sibling-span-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const siblingClientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "sibling-span-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10050,
                duration: 153,
                tags: "{}"
            });

            expect(siblingClientSpan.isParentOf(siblingServerSpan)).toBe(true);
        });

        it("should return false if the sibling client span is provided", () => {
            const siblingServerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "sibling-span-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const siblingClientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "sibling-span-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10050,
                duration: 153,
                tags: "{}"
            });

            expect(siblingServerSpan.isParentOf(siblingClientSpan)).toBe(false);
        });

        it("should return false if a span from another trace is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-b-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const mismatchedTraceIdSpan = new Span({
                traceId: "trace-b-id",
                spanId: "span-b-id",
                parentId: "span-c-id",
                serviceName: "parent-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });

            expect(mismatchedTraceIdSpan.isParentOf(span)).toBe(false);
        });

        it("should return false if a non-related span is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-b-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const nonRelatedSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-c-id",
                parentId: "span-d-id",
                serviceName: "non-parent-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });

            expect(nonRelatedSpan.isParentOf(span)).toBe(false);
        });

        it("should return false if null/undefined is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });

            expect(span.isParentOf(null)).toBe(false);
            expect(span.isSiblingOf(undefined)).toBe(false);
        });
    });

    describe("hasSibling()", () => {
        it("should return true if a client or server span is provided", () => {
            const clientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const serverSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10010,
                duration: 50,
                tags: "{}"
            });

            expect(clientSpan.hasSibling()).toBe(true);
            expect(serverSpan.hasSibling()).toBe(true);
        });

        it("should return false if a span of kind not equal to client or server is provided", () => {
            const producerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.PRODUCER,
                startTime: 10010,
                duration: 50,
                tags: "{}"
            });
            const consumerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CONSUMER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const unknownTypeSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                startTime: 10010,
                duration: 50,
                tags: "{}"
            });

            expect(producerSpan.hasSibling()).toBe(false);
            expect(consumerSpan.hasSibling()).toBe(false);
            expect(unknownTypeSpan.hasSibling()).toBe(false);
        });
    });

    describe("addSpanReference()", () => {
        it("should add as child and return true if the child is provided", () => {
            const span = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CONSUMER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const childSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-x-id",
                parentId: "span-a-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
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
                parentId: "span-a-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.PRODUCER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const parentSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
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
                parentId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const siblingServerSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-b-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
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
                parentId: "span-b-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const siblingClientSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-a-id",
                parentId: "span-b-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
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
                parentId: "span-x-id",
                serviceName: "consumer-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const nonRelatedSpan = new Span({
                traceId: "trace-a-id",
                spanId: "span-b-id",
                parentId: "span-y-id",
                serviceName: "vick-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
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
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
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

    describe("resetSpanReferences()", () => {
        it("should clear all the span references", () => {
            const span = new Span({
                traceId: "trace-id",
                spanId: "span-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const siblingSpan = new Span({
                traceId: "trace-id",
                spanId: "span-server-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const parentSpan = new Span({
                traceId: "trace-id",
                spanId: "span-parent-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const childSpanA = new Span({
                traceId: "trace-id",
                spanId: "span-child-a-id",
                parentId: "span-id",
                serviceName: "worker",
                operationName: "work",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const childSpanB = new Span({
                traceId: "trace-id",
                spanId: "span-child-b-id",
                parentId: "span-id",
                serviceName: "worker",
                operationName: "work",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });

            span.addSpanReference(siblingSpan);
            span.addSpanReference(parentSpan);
            span.addSpanReference(childSpanA);
            span.addSpanReference(childSpanB);
            span.resetSpanReferences();

            expect(span.children.size).toBe(0);
            expect(span.parent).toBeNull();
            expect(span.sibling).toBeNull();
            expect(span.treeDepth).toBe(0);
        });
    });

    describe("walk()", () => {
        it("should walk down the tree", () => {
            const parentSpan = new Span({
                traceId: "trace-id",
                spanId: "span-parent-id",
                parentId: "span-parent-parent-id",
                serviceName: "parent-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 1000,
                tags: "{}"
            });
            const siblingSpan = new Span({
                traceId: "trace-id",
                spanId: "span-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10020,
                duration: 980,
                tags: "{}"
            });
            const span = new Span({
                traceId: "trace-id",
                spanId: "span-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10010,
                duration: 960,
                tags: "{}"
            });
            const childSpanA = new Span({
                traceId: "trace-id",
                spanId: "span-child--a-id",
                parentId: "span-id",
                serviceName: "worker-A",
                operationName: "work",
                startTime: 10030,
                duration: 400,
                tags: "{}"
            });
            const childSpanB = new Span({
                traceId: "trace-id",
                spanId: "span-child-b-id",
                parentId: "span-id",
                serviceName: "worker-B",
                operationName: "work",
                startTime: 10530,
                duration: 300,
                tags: "{}"
            });

            const spansList = [span, parentSpan, siblingSpan, childSpanA, childSpanB];
            for (let i = 0; i < spansList.length; i++) {
                for (let j = 0; j < spansList.length; j++) {
                    if (i !== j) {
                        spansList[i].addSpanReference(spansList[j]);
                    }
                }
            }

            const initialData = {count: 0};
            const preWalkNodes = [];
            const postWalkNodes = [];
            const walkData = [initialData];
            parentSpan.walk((node, data) => {
                expect(preWalkNodes).not.toContain(node);
                expect(postWalkNodes).not.toContain(node);
                expect(walkData).toContain(data);

                const newData = {
                    id: node.getUniqueId(),
                    count: data.count + 1
                };
                preWalkNodes.push(node);
                walkData.push(newData);
                return newData;
            }, initialData, (node) => {
                expect(preWalkNodes).toContain(node);
                expect(postWalkNodes).not.toContain(node);

                postWalkNodes.push(node);
            });

            expect(preWalkNodes[0]).toBe(parentSpan);
            expect(preWalkNodes[1]).toBe(siblingSpan);
            expect(preWalkNodes[2]).toBe(span);
            expect(preWalkNodes[3]).toBe(childSpanA);
            expect(preWalkNodes[4]).toBe(childSpanB);

            expect(postWalkNodes[0]).toBe(childSpanA);
            expect(postWalkNodes[1]).toBe(childSpanB);
            expect(postWalkNodes[2]).toBe(span);
            expect(postWalkNodes[3]).toBe(siblingSpan);
            expect(postWalkNodes[4]).toBe(parentSpan);

            expect(walkData[0].id).toBeUndefined();
            expect(walkData[0].count).toBe(0);
            expect(walkData[1].id).toBe(parentSpan.getUniqueId());
            expect(walkData[1].count).toBe(1);
            expect(walkData[2].id).toBe(siblingSpan.getUniqueId());
            expect(walkData[2].count).toBe(2);
            expect(walkData[3].id).toBe(span.getUniqueId());
            expect(walkData[3].count).toBe(3);
            expect(walkData[4].id).toBe(childSpanA.getUniqueId());
            expect(walkData[4].count).toBe(4);
            expect(walkData[5].id).toBe(childSpanB.getUniqueId());
            expect(walkData[5].count).toBe(4);
        });
    });

    describe("getUniqueId()", () => {
        it("should return a unique ID across traces", () => {
            const span = new Span({
                traceId: "trace-id",
                spanId: "span-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const siblingSpan = new Span({
                traceId: "trace-id",
                spanId: "span-server-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const parentSpan = new Span({
                traceId: "trace-id",
                spanId: "span-parent-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const differentTraceSpan = new Span({
                traceId: "trace-different-id",
                spanId: "span-id",
                parentId: "span-parent-id",
                serviceName: "test-service",
                operationName: "set-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 200000,
                duration: 1203,
                tags: "{}"
            });

            expect(span.getUniqueId()).not.toBe(siblingSpan.getUniqueId());
            expect(span.getUniqueId()).not.toBe(parentSpan.getUniqueId());
            expect(span.getUniqueId()).not.toBe(differentTraceSpan.getUniqueId());
        });
    });

    const globalGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-a-id",
        parentId: "trace-x-id",
        serviceName: "global-gateway",
        operationName: "get-hr-info",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10000000,
        duration: 3160000,
        tags: "{}"
    });
    const hrCellGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-b-id",
        parentId: "span-a-id",
        cell: "hr",
        serviceName: "hr-cell-gateway",
        operationName: "call-hr-cell",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10020000,
        duration: 3090000,
        tags: "{}"
    });
    const employeeServiceServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-c-id",
        parentId: "span-b-id",
        cell: "hr",
        serviceName: "hr--employee-service",
        operationName: "get-employee-data",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10040000,
        duration: 3040000,
        tags: "{}"
    });
    const istioMixerServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-d-id",
        parentId: "span-c-id",
        serviceName: "istio-mixer",
        operationName: "is-authorized",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10060000,
        duration: 940000,
        tags: "{}"
    });

    describe("isFromCellGateway()", () => {
        it("should return true if the span is from a Cell Gateway", () => {
            expect(hrCellGatewayServerSpan.isFromCellGateway()).toBe(true);
        });

        it("should return false if the span is not from a Cell Gateway", () => {
            expect(globalGatewayServerSpan.isFromCellGateway()).toBe(false);
            expect(istioMixerServerSpan.isFromCellGateway()).toBe(false);
            expect(employeeServiceServerSpan.isFromCellGateway()).toBe(false);
        });
    });

    describe("isFromVICKSystemComponent()", () => {
        it("should return true if the span is from Global Gateway", () => {
            expect(globalGatewayServerSpan.isFromVICKSystemComponent()).toBe(true);
        });

        it("should return true if the span is from a Cell Gateway", () => {
            expect(hrCellGatewayServerSpan.isFromVICKSystemComponent()).toBe(true);
        });

        it("should return true if the span is from Istio Mixer", () => {
            expect(istioMixerServerSpan.isFromVICKSystemComponent()).toBe(false);
        });

        it("should return false if the span is from a custom service", () => {
            expect(employeeServiceServerSpan.isFromVICKSystemComponent()).toBe(false);
        });
    });

    describe("isFromIstioSystemComponent()", () => {
        it("should return true if the span is from Global Gateway", () => {
            expect(globalGatewayServerSpan.isFromIstioSystemComponent()).toBe(false);
        });

        it("should return true if the span is from a Cell Gateway", () => {
            expect(hrCellGatewayServerSpan.isFromIstioSystemComponent()).toBe(false);
        });

        it("should return true if the span is from Istio Mixer", () => {
            expect(istioMixerServerSpan.isFromIstioSystemComponent()).toBe(true);
        });

        it("should return false if the span is from a custom service", () => {
            expect(employeeServiceServerSpan.isFromIstioSystemComponent()).toBe(false);
        });
    });

    describe("shallowClone()", () => {
        it("should return a copy of the initial span without span references", () => {
            const span = new Span({
                traceId: "trace-id",
                spanId: "span-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.CLIENT,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const siblingSpan = new Span({
                traceId: "trace-id",
                spanId: "span-server-id",
                parentId: "span-parent-id",
                serviceName: "server-service",
                operationName: "get-resource",
                kind: Constants.Span.Kind.SERVER,
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const parentSpan = new Span({
                traceId: "trace-id",
                spanId: "span-parent-id",
                parentId: "span-parent-id",
                serviceName: "client-service",
                operationName: "get-resource",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const childSpanA = new Span({
                traceId: "trace-id",
                spanId: "span-child--aid",
                parentId: "span-id",
                serviceName: "worker",
                operationName: "work",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });
            const childSpanB = new Span({
                traceId: "trace-id",
                spanId: "span-child-b-id",
                parentId: "span-id",
                serviceName: "worker",
                operationName: "work",
                startTime: 10000,
                duration: 100,
                tags: "{}"
            });

            span.addSpanReference(siblingSpan);
            span.addSpanReference(parentSpan);
            span.addSpanReference(childSpanA);
            span.addSpanReference(childSpanB);
            span.componentType = Constants.ComponentType.VICK;
            span.cell = {
                name: "cell-a"
            };

            const clone = span.shallowClone();

            expect(clone).not.toBe(span);
            expect(clone.traceId).toBe(span.traceId);
            expect(clone.spanId).toBe(span.spanId);
            expect(clone.parentId).toBe(span.parentId);
            expect(clone.serviceName).toBe(span.serviceName);
            expect(clone.operationName).toBe(span.operationName);
            expect(clone.kind).toBe(span.kind);
            expect(clone.startTime).toBe(span.startTime);
            expect(clone.duration).toBe(span.duration);
            expect(clone.tags).not.toBe(span.tags);
            for (const tagKey in clone.tags) {
                if (clone.tags.hasOwnProperty(tagKey)) {
                    expect(span.tags[tagKey]).not.toBeNull();
                    expect(clone.tags[tagKey]).toBe(span.tags[tagKey]);
                }
            }
            for (const tagKey in span.tags) {
                if (span.tags.hasOwnProperty(tagKey)) {
                    expect(clone.tags[tagKey]).not.toBeNull();
                    expect(span.tags[tagKey]).toBe(clone.tags[tagKey]);
                }
            }
            expect(clone.componentType).toBe(span.componentType);
            expect(clone.cell).not.toBeNull();
            expect(clone.cell).not.toBe(span.cell);
            expect(clone.cell.name).toBe(span.cell.name);

            expect(clone.parent).toBeNull();
            expect(clone.sibling).toBeNull();
            expect(clone.children.size).toBe(0);
        });
    });
});
