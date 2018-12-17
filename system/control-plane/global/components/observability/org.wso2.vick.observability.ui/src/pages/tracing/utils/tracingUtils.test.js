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
import TracingUtils from "./tracingUtils";

describe("TracingUtils", () => {
    const globalGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "trace-x-id",
        parentId: "undefined",
        serviceName: "global-gateway",
        operationName: "get-hr-info",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10000000,
        duration: 3160000,
        tags: "{}"
    });
    const globalGatewayClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-b-id",
        parentId: "trace-x-id",
        serviceName: "global-gateway",
        operationName: "call-hr-cell",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10010000,
        duration: 3110000,
        tags: "{}"
    });
    const hrCellGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-b-id",
        parentId: "trace-x-id",
        cell: "hr",
        serviceName: "hr-cell-gateway",
        operationName: "call-hr-cell",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10020000,
        duration: 3090000,
        tags: "{}"
    });
    const hrCellGatewayClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-c-id",
        parentId: "span-b-id",
        cell: "hr",
        serviceName: "hr-cell-gateway",
        operationName: "get-employee-data",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10030000,
        duration: 3060000,
        tags: "{}"
    });
    const employeeServiceServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-c-id",
        parentId: "span-b-id",
        cell: "hr",
        serviceName: "employee-service",
        operationName: "get-employee-data",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10040000,
        duration: 3040000,
        tags: "{}"
    });
    const employeeServiceToIstioMixerClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-d-id",
        parentId: "span-c-id",
        cell: "hr",
        serviceName: "employee-service",
        operationName: "is-authorized",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10050000,
        duration: 990000,
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
    const istioMixerWorkerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-e-id",
        parentId: "span-d-id",
        serviceName: "istio-mixer",
        operationName: "authorization",
        startTime: 10070000,
        duration: 890000,
        tags: "{}"
    });
    const employeeServiceToReviewsServiceClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-f-id",
        parentId: "span-c-id",
        cell: "hr",
        serviceName: "employee-service",
        operationName: "get-reviews",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 11050000,
        duration: 990000,
        tags: "{}"
    });
    const reviewsServiceServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-f-id",
        parentId: "span-c-id",
        cell: "hr",
        serviceName: "reviews-service",
        operationName: "get-reviews",
        kind: Constants.Span.Kind.SERVER,
        startTime: 11100000,
        duration: 890000,
        tags: "{}"
    });
    const employeeServiceToStockOptionsCellClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-g-id",
        parentId: "span-c-id",
        cell: "hr",
        serviceName: "employee-service",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 12060000,
        duration: 990000,
        tags: "{}"
    });
    const stockOptionsCellGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-g-id",
        parentId: "span-c-id",
        cell: "stock-options",
        serviceName: "stock-options-cell-gateway",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.SERVER,
        startTime: 12100000,
        duration: 890000,
        tags: "{}"
    });
    const stockOptionsCellGatewayClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-h-id",
        parentId: "span-g-id",
        cell: "stock-options",
        serviceName: "stock-options-cell-gateway",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 12150000,
        duration: 790000,
        tags: "{}"
    });
    const stockOptionsServiceServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-h-id",
        parentId: "span-g-id",
        cell: "stock-options",
        serviceName: "stock-options-service",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.SERVER,
        startTime: 12200000,
        duration: 690000,
        tags: "{}"
    });

    const orderedSpanList = [
        globalGatewayServerSpan, globalGatewayClientSpan, hrCellGatewayServerSpan, hrCellGatewayClientSpan,
        employeeServiceServerSpan, employeeServiceToIstioMixerClientSpan, istioMixerServerSpan,
        istioMixerWorkerSpan, employeeServiceToReviewsServiceClientSpan, reviewsServiceServerSpan,
        employeeServiceToStockOptionsCellClientSpan, stockOptionsCellGatewayServerSpan,
        stockOptionsCellGatewayClientSpan, stockOptionsServiceServerSpan
    ];

    describe("buildTree()", () => {
        it("should build the tracing tree from the spans list", () => {
            // Shuffle spans list
            const spansList = orderedSpanList.map((a) => [Math.random(), a])
                .sort((a, b) => a[0] - b[0])
                .map((a) => a[1]);

            expect(TracingUtils.buildTree(spansList)).toBe(globalGatewayServerSpan);

            expect(globalGatewayServerSpan.parent).toBeNull();
            expect(globalGatewayServerSpan.sibling).toBeNull();
            expect(globalGatewayServerSpan.children.size).toBe(1);
            expect(globalGatewayServerSpan.children.has(globalGatewayClientSpan)).toBe(true);

            expect(globalGatewayClientSpan.parent).toBe(globalGatewayServerSpan);
            expect(globalGatewayClientSpan.sibling).toBe(hrCellGatewayServerSpan);
            expect(globalGatewayClientSpan.children.size).toBe(1);
            expect(globalGatewayClientSpan.children.has(hrCellGatewayServerSpan)).toBe(true);

            expect(hrCellGatewayServerSpan.parent).toBe(globalGatewayClientSpan);
            expect(hrCellGatewayServerSpan.sibling).toBe(globalGatewayClientSpan);
            expect(hrCellGatewayServerSpan.children.size).toBe(1);
            expect(hrCellGatewayServerSpan.children.has(hrCellGatewayClientSpan)).toBe(true);

            expect(hrCellGatewayClientSpan.parent).toBe(hrCellGatewayServerSpan);
            expect(hrCellGatewayClientSpan.sibling).toBe(employeeServiceServerSpan);
            expect(hrCellGatewayClientSpan.children.size).toBe(1);
            expect(hrCellGatewayClientSpan.children.has(employeeServiceServerSpan)).toBe(true);

            expect(employeeServiceServerSpan.parent).toBe(hrCellGatewayClientSpan);
            expect(employeeServiceServerSpan.sibling).toBe(hrCellGatewayClientSpan);
            expect(employeeServiceServerSpan.children.size).toBe(3);
            expect(employeeServiceServerSpan.children.has(employeeServiceToIstioMixerClientSpan)).toBe(true);
            expect(employeeServiceServerSpan.children.has(employeeServiceToReviewsServiceClientSpan)).toBe(true);
            expect(employeeServiceServerSpan.children.has(employeeServiceToStockOptionsCellClientSpan)).toBe(true);

            expect(employeeServiceToIstioMixerClientSpan.parent).toBe(employeeServiceServerSpan);
            expect(employeeServiceToIstioMixerClientSpan.sibling).toBe(istioMixerServerSpan);
            expect(employeeServiceToIstioMixerClientSpan.children.size).toBe(1);
            expect(employeeServiceToIstioMixerClientSpan.children.has(istioMixerServerSpan)).toBe(true);

            expect(istioMixerServerSpan.parent).toBe(employeeServiceToIstioMixerClientSpan);
            expect(istioMixerServerSpan.sibling).toBe(employeeServiceToIstioMixerClientSpan);
            expect(istioMixerServerSpan.children.size).toBe(1);
            expect(istioMixerServerSpan.children.has(istioMixerWorkerSpan)).toBe(true);

            expect(istioMixerWorkerSpan.parent).toBe(istioMixerServerSpan);
            expect(istioMixerWorkerSpan.sibling).toBeNull();
            expect(istioMixerWorkerSpan.children.size).toBe(0);

            expect(employeeServiceToReviewsServiceClientSpan.parent).toBe(employeeServiceServerSpan);
            expect(employeeServiceToReviewsServiceClientSpan.sibling).toBe(reviewsServiceServerSpan);
            expect(employeeServiceToReviewsServiceClientSpan.children.size).toBe(1);
            expect(employeeServiceToReviewsServiceClientSpan.children.has(reviewsServiceServerSpan)).toBe(true);

            expect(reviewsServiceServerSpan.parent).toBe(employeeServiceToReviewsServiceClientSpan);
            expect(reviewsServiceServerSpan.sibling).toBe(employeeServiceToReviewsServiceClientSpan);
            expect(reviewsServiceServerSpan.children.size).toBe(0);

            expect(employeeServiceToStockOptionsCellClientSpan.parent).toBe(employeeServiceServerSpan);
            expect(employeeServiceToStockOptionsCellClientSpan.sibling).toBe(stockOptionsCellGatewayServerSpan);
            expect(employeeServiceToStockOptionsCellClientSpan.children.size).toBe(1);
            expect(employeeServiceToStockOptionsCellClientSpan.children.has(stockOptionsCellGatewayServerSpan))
                .toBe(true);

            expect(stockOptionsCellGatewayServerSpan.parent).toBe(employeeServiceToStockOptionsCellClientSpan);
            expect(stockOptionsCellGatewayServerSpan.sibling).toBe(employeeServiceToStockOptionsCellClientSpan);
            expect(stockOptionsCellGatewayServerSpan.children.size).toBe(1);
            expect(stockOptionsCellGatewayServerSpan.children.has(stockOptionsCellGatewayClientSpan)).toBe(true);

            expect(stockOptionsCellGatewayClientSpan.parent).toBe(stockOptionsCellGatewayServerSpan);
            expect(stockOptionsCellGatewayClientSpan.sibling).toBe(stockOptionsServiceServerSpan);
            expect(stockOptionsCellGatewayClientSpan.children.size).toBe(1);
            expect(stockOptionsCellGatewayClientSpan.children.has(stockOptionsServiceServerSpan)).toBe(true);

            expect(stockOptionsServiceServerSpan.parent).toBe(stockOptionsCellGatewayClientSpan);
            expect(stockOptionsServiceServerSpan.sibling).toBe(stockOptionsCellGatewayClientSpan);
            expect(stockOptionsServiceServerSpan.children.size).toBe(0);
        });
    });


    describe("labelSpanTree()", () => {
        it("should label the necessary nodes according to cell, component type", () => {
            TracingUtils.labelSpanTree(globalGatewayServerSpan);

            expect(globalGatewayServerSpan.cell).toBeNull();
            expect(globalGatewayServerSpan.componentType).toBe(Constants.ComponentType.VICK);
            expect(globalGatewayServerSpan.treeDepth).toBe(0);

            expect(globalGatewayClientSpan.cell).toBeNull();
            expect(globalGatewayClientSpan.componentType).toBe(Constants.ComponentType.VICK);
            expect(globalGatewayClientSpan.treeDepth).toBe(1);

            expect(hrCellGatewayServerSpan.cell).not.toBeNull();
            expect(hrCellGatewayServerSpan.cell.name).toBe("hr");
            expect(hrCellGatewayServerSpan.componentType).toBe(Constants.ComponentType.VICK);
            expect(hrCellGatewayServerSpan.treeDepth).toBe(2);

            expect(hrCellGatewayClientSpan.cell).not.toBeNull();
            expect(hrCellGatewayClientSpan.cell.name).toBe("hr");
            expect(hrCellGatewayClientSpan.componentType).toBe(Constants.ComponentType.VICK);
            expect(hrCellGatewayClientSpan.treeDepth).toBe(3);

            expect(employeeServiceServerSpan.cell).not.toBeNull();
            expect(employeeServiceServerSpan.cell.name).toBe("hr");
            expect(employeeServiceServerSpan.componentType).toBe(Constants.ComponentType.MICROSERVICE);
            expect(employeeServiceServerSpan.treeDepth).toBe(4);

            expect(employeeServiceToIstioMixerClientSpan.cell).not.toBeNull();
            expect(employeeServiceToIstioMixerClientSpan.cell.name).toBe("hr");
            expect(employeeServiceToIstioMixerClientSpan.componentType).toBe(Constants.ComponentType.MICROSERVICE);
            expect(employeeServiceToIstioMixerClientSpan.treeDepth).toBe(5);

            expect(istioMixerServerSpan.cell).toBeNull();
            expect(istioMixerServerSpan.componentType).toBe(Constants.ComponentType.ISTIO);
            expect(istioMixerServerSpan.treeDepth).toBe(6);

            expect(istioMixerWorkerSpan.cell).toBeNull();
            expect(istioMixerWorkerSpan.componentType).toBe(Constants.ComponentType.ISTIO);
            expect(istioMixerWorkerSpan.treeDepth).toBe(7);

            expect(employeeServiceToReviewsServiceClientSpan.cell).not.toBeNull();
            expect(employeeServiceToReviewsServiceClientSpan.cell.name).toBe("hr");
            expect(employeeServiceToReviewsServiceClientSpan.componentType)
                .toBe(Constants.ComponentType.MICROSERVICE);
            expect(employeeServiceToReviewsServiceClientSpan.treeDepth).toBe(5);

            expect(reviewsServiceServerSpan.cell).not.toBeNull();
            expect(reviewsServiceServerSpan.cell.name).toBe("hr");
            expect(reviewsServiceServerSpan.componentType).toBe(Constants.ComponentType.MICROSERVICE);
            expect(reviewsServiceServerSpan.treeDepth).toBe(6);

            expect(employeeServiceToStockOptionsCellClientSpan.cell).not.toBeNull();
            expect(employeeServiceToStockOptionsCellClientSpan.cell.name).toBe("hr");
            expect(employeeServiceToStockOptionsCellClientSpan.componentType)
                .toBe(Constants.ComponentType.MICROSERVICE);
            expect(employeeServiceToStockOptionsCellClientSpan.treeDepth).toBe(5);

            expect(stockOptionsCellGatewayServerSpan.cell).not.toBeNull();
            expect(stockOptionsCellGatewayServerSpan.cell.name).toBe("stock-options");
            expect(stockOptionsCellGatewayServerSpan.componentType).toBe(Constants.ComponentType.VICK);
            expect(stockOptionsCellGatewayServerSpan.treeDepth).toBe(6);

            expect(stockOptionsCellGatewayClientSpan.cell).not.toBeNull();
            expect(stockOptionsCellGatewayClientSpan.cell.name).toBe("stock-options");
            expect(stockOptionsCellGatewayClientSpan.componentType).toBe(Constants.ComponentType.VICK);
            expect(stockOptionsCellGatewayClientSpan.treeDepth).toBe(7);

            expect(stockOptionsServiceServerSpan.cell).not.toBeNull();
            expect(stockOptionsServiceServerSpan.cell.name).toBe("stock-options");
            expect(stockOptionsServiceServerSpan.componentType).toBe(Constants.ComponentType.MICROSERVICE);
            expect(stockOptionsServiceServerSpan.treeDepth).toBe(8);
        });
    });

    describe("getOrderedList()", () => {
        it("should return the nodes ordered by start time and tree structure", () => {
            const resultList = TracingUtils.getOrderedList(globalGatewayServerSpan);

            expect(resultList).toHaveLength(14);
            expect(resultList[0]).toBe(globalGatewayServerSpan);
            expect(resultList[1]).toBe(globalGatewayClientSpan);
            expect(resultList[2]).toBe(hrCellGatewayServerSpan);
            expect(resultList[3]).toBe(hrCellGatewayClientSpan);
            expect(resultList[4]).toBe(employeeServiceServerSpan);
            expect(resultList[5]).toBe(employeeServiceToIstioMixerClientSpan);
            expect(resultList[6]).toBe(istioMixerServerSpan);
            expect(resultList[7]).toBe(istioMixerWorkerSpan);
            expect(resultList[8]).toBe(employeeServiceToReviewsServiceClientSpan);
            expect(resultList[9]).toBe(reviewsServiceServerSpan);
            expect(resultList[10]).toBe(employeeServiceToStockOptionsCellClientSpan);
            expect(resultList[11]).toBe(stockOptionsCellGatewayServerSpan);
            expect(resultList[12]).toBe(stockOptionsCellGatewayClientSpan);
            expect(resultList[13]).toBe(stockOptionsServiceServerSpan);
        });
    });

    describe("removeSpanFromTree()", () => {
        it("should remove the references in other nodes that point to this to form the tree", () => {
            TracingUtils.removeSpanFromTree(employeeServiceServerSpan);

            expect(employeeServiceToIstioMixerClientSpan.parent).toBe(hrCellGatewayClientSpan);
            expect(employeeServiceToReviewsServiceClientSpan.parent).toBe(hrCellGatewayClientSpan);
            expect(employeeServiceToStockOptionsCellClientSpan.parent).toBe(hrCellGatewayClientSpan);

            expect(hrCellGatewayClientSpan.children.size).toBe(3);
            expect(hrCellGatewayClientSpan.children.has(employeeServiceToIstioMixerClientSpan)).toBe(true);
            expect(hrCellGatewayClientSpan.children.has(employeeServiceToReviewsServiceClientSpan)).toBe(true);
            expect(hrCellGatewayClientSpan.children.has(employeeServiceToStockOptionsCellClientSpan)).toBe(true);
            expect(hrCellGatewayClientSpan.sibling).toBeNull();
        });
    });

    describe("resetTreeSpanReferences()", () => {
        it("should clear all connections to other nodes", () => {
            TracingUtils.resetTreeSpanReferences(orderedSpanList);

            for (let i = 0; i < orderedSpanList.length; i++) {
                expect(orderedSpanList[i].children.size).toBe(0);
                expect(orderedSpanList[i].parent).toBeNull();
                expect(orderedSpanList[i].sibling).toBeNull();
            }
        });
    });
});
