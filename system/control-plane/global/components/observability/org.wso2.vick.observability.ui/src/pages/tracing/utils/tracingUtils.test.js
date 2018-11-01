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

import Constants from "./constants";
import Span from "./span";
import TracingUtils from "./tracingUtils";

describe("TracingUtils", () => {
    const globalGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-a-id",
        parentSpanId: "trace-x-id",
        serviceName: Constants.VICK.System.GLOBAL_GATEWAY_NAME,
        operationName: "get-hr-info",
        kind: Constants.Span.Kind.SERVER,
        startTime: 100000,
        duration: 10000,
        tags: {}
    });
    const globalGatewayClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-b-id",
        parentSpanId: "span-a-id",
        serviceName: Constants.VICK.System.GLOBAL_GATEWAY_NAME,
        operationName: "call-hr-cell",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10010,
        duration: 9980,
        tags: {}
    });
    const hrCellGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-b-id",
        parentSpanId: "span-a-id",
        serviceName: "src:0.0.0.hr_1_0_0_employee",
        operationName: "call-hr-cell",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10020,
        duration: 9960,
        tags: {}
    });
    const hrCellGatewayClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-c-id",
        parentSpanId: "span-b-id",
        serviceName: "src:0.0.0.hr_1_0_0_employee",
        operationName: "get-employee-data",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10030,
        duration: 9940,
        tags: {}
    });
    const employeeServiceServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-c-id",
        parentSpanId: "span-b-id",
        serviceName: "employee-service",
        operationName: "get-employee-data",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10060,
        duration: 9940,
        tags: {}
    });
    const employeeServiceToIstioMixerClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-d-id",
        parentSpanId: "span-c-id",
        serviceName: "employee-service",
        operationName: "is-authorized",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10060,
        duration: 9940,
        tags: {}
    });
    const istioMixerServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-d-id",
        parentSpanId: "span-c-id",
        serviceName: Constants.VICK.System.ISTIO_MIXER_NAME,
        operationName: "is-authorized",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10040,
        duration: 9920,
        tags: {}
    });
    const istioMixerWorkerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-e-id",
        parentSpanId: "span-d-id",
        serviceName: Constants.VICK.System.ISTIO_MIXER_NAME,
        operationName: "authorization",
        startTime: 10050,
        duration: 9900,
        tags: {}
    });
    const employeeServiceToReviewsServiceClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-f-id",
        parentSpanId: "span-c-id",
        serviceName: "employee-service",
        operationName: "get-reviews",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10060,
        duration: 9940,
        tags: {}
    });
    const reviewsServiceServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-f-id",
        parentSpanId: "span-c-id",
        serviceName: "reviews-service",
        operationName: "get-reviews",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10000,
        duration: 100,
        tags: {}
    });
    const employeeServiceToStockOptionsCellClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-g-id",
        parentSpanId: "span-c-id",
        serviceName: "employee-service",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10060,
        duration: 9940,
        tags: {}
    });
    const stockOptionsCellGatewayServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-g-id",
        parentSpanId: "span-c-id",
        serviceName: "src:0.0.0.stock_options_1_0_0_employee",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10000,
        duration: 100,
        tags: {}
    });
    const stockOptionsCellGatewayClientSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-h-id",
        parentSpanId: "span-g-id",
        serviceName: "src:0.0.0.stock_options_1_0_0_employee",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.CLIENT,
        startTime: 10000,
        duration: 100,
        tags: {}
    });
    const stockOptionsServiceServerSpan = new Span({
        traceId: "trace-x-id",
        spanId: "span-h-id",
        parentSpanId: "span-g-id",
        serviceName: "stock-options-service",
        operationName: "get-employee-stock-options",
        kind: Constants.Span.Kind.SERVER,
        startTime: 10000,
        duration: 100,
        tags: {}
    });

    describe("buildTree()", () => {
        it("should build the tracing tree from the spans list", () => {
            const orderedSpanList = [
                globalGatewayServerSpan, globalGatewayClientSpan, hrCellGatewayServerSpan, hrCellGatewayClientSpan,
                employeeServiceServerSpan, employeeServiceToIstioMixerClientSpan, istioMixerServerSpan,
                istioMixerWorkerSpan, employeeServiceToReviewsServiceClientSpan, reviewsServiceServerSpan,
                employeeServiceToStockOptionsCellClientSpan, stockOptionsCellGatewayServerSpan, stockOptionsCellGatewayClientSpan,
                stockOptionsServiceServerSpan,
            ];

            // Shuffle spans list
            const spansList = orderedSpanList.map(a => [Math.random(), a])
                .sort((a, b) => a[0] - b[0])
                .map(a => a[1]);

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
            expect(istioMixerWorkerSpan.sibling).toBe(null);
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
            expect(employeeServiceToStockOptionsCellClientSpan.children.has(stockOptionsCellGatewayServerSpan)).toBe(true);

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

            expect(globalGatewayServerSpan.cell).toBe(null);
            expect(globalGatewayServerSpan.isSystemComponent).toBe(true);

            expect(globalGatewayClientSpan.cell).toBe(null);
            expect(globalGatewayClientSpan.isSystemComponent).toBe(true);

            expect(hrCellGatewayServerSpan.cell).not.toBeNull();
            expect(hrCellGatewayServerSpan.cell.name).toBe("hr");
            expect(hrCellGatewayServerSpan.cell.version).toBe("1.0.0");
            expect(hrCellGatewayServerSpan.isSystemComponent).toBe(true);

            expect(hrCellGatewayClientSpan.cell).not.toBeNull();
            expect(hrCellGatewayClientSpan.cell.name).toBe("hr");
            expect(hrCellGatewayClientSpan.cell.version).toBe("1.0.0");
            expect(hrCellGatewayClientSpan.isSystemComponent).toBe(true);

            expect(employeeServiceServerSpan.cell).not.toBeNull();
            expect(employeeServiceServerSpan.cell.name).toBe("hr");
            expect(employeeServiceServerSpan.cell.version).toBe("1.0.0");
            expect(employeeServiceServerSpan.isSystemComponent).toBe(false);

            expect(employeeServiceToIstioMixerClientSpan.cell).not.toBeNull();
            expect(employeeServiceToIstioMixerClientSpan.cell.name).toBe("hr");
            expect(employeeServiceToIstioMixerClientSpan.cell.version).toBe("1.0.0");
            expect(employeeServiceToIstioMixerClientSpan.isSystemComponent).toBe(false);

            expect(istioMixerServerSpan.cell).toBeNull();
            expect(istioMixerServerSpan.isSystemComponent).toBe(true);

            expect(istioMixerWorkerSpan.cell).toBeNull();
            expect(istioMixerWorkerSpan.isSystemComponent).toBe(true);

            expect(employeeServiceToReviewsServiceClientSpan.cell).not.toBeNull();
            expect(employeeServiceToReviewsServiceClientSpan.cell.name).toBe("hr");
            expect(employeeServiceToReviewsServiceClientSpan.cell.version).toBe("1.0.0");
            expect(employeeServiceToReviewsServiceClientSpan.isSystemComponent).toBe(false);

            expect(reviewsServiceServerSpan.cell).not.toBeNull();
            expect(reviewsServiceServerSpan.cell.name).toBe("hr");
            expect(reviewsServiceServerSpan.cell.version).toBe("1.0.0");
            expect(reviewsServiceServerSpan.isSystemComponent).toBe(false);

            expect(employeeServiceToStockOptionsCellClientSpan.cell).not.toBeNull();
            expect(employeeServiceToStockOptionsCellClientSpan.cell.name).toBe("hr");
            expect(employeeServiceToStockOptionsCellClientSpan.cell.version).toBe("1.0.0");
            expect(employeeServiceToStockOptionsCellClientSpan.isSystemComponent).toBe(false);

            expect(stockOptionsCellGatewayServerSpan.cell).not.toBeNull();
            expect(stockOptionsCellGatewayServerSpan.cell.name).toBe("stock_options");
            expect(stockOptionsCellGatewayServerSpan.cell.version).toBe("1.0.0");
            expect(stockOptionsCellGatewayServerSpan.isSystemComponent).toBe(true);

            expect(stockOptionsCellGatewayClientSpan.cell).not.toBeNull();
            expect(stockOptionsCellGatewayClientSpan.cell.name).toBe("stock_options");
            expect(stockOptionsCellGatewayClientSpan.cell.version).toBe("1.0.0");
            expect(stockOptionsCellGatewayClientSpan.isSystemComponent).toBe(true);

            expect(stockOptionsServiceServerSpan.cell).not.toBeNull();
            expect(stockOptionsServiceServerSpan.cell.name).toBe("stock_options");
            expect(stockOptionsServiceServerSpan.cell.version).toBe("1.0.0");
            expect(stockOptionsServiceServerSpan.isSystemComponent).toBe(false);
        });
    });

    describe("getCell()", () => {
        it("should return the cell information if the span is from a Cell Gateway", () => {
            const cell = TracingUtils.getCell(hrCellGatewayServerSpan);

            expect(cell.name).toBe("hr");
            expect(cell.version).toBe("1.0.0");
        });

        it("should throw error if the span is not not from a Cell Gateway", () => {
            expect(() => TracingUtils.getCell(globalGatewayServerSpan)).toThrow();
            expect(() => TracingUtils.getCell(istioMixerServerSpan)).toThrow();
            expect(() => TracingUtils.getCell(employeeServiceServerSpan)).toThrow();
        });

        it("should return null if null is provided", () => {
            expect(TracingUtils.getCell(null)).toBeNull();
            expect(TracingUtils.getCell(undefined)).toBeNull();
        });
    });

    describe("isFromCellGateway()", () => {
        it("should return true if the span is from a Cell Gateway", () => {
            expect(TracingUtils.isFromCellGateway(hrCellGatewayServerSpan)).toBe(true);
        });

        it("should return false if the span is not from a Cell Gateway", () => {
            expect(TracingUtils.isFromCellGateway(globalGatewayServerSpan)).toBe(false);
            expect(TracingUtils.isFromCellGateway(istioMixerServerSpan)).toBe(false);
            expect(TracingUtils.isFromCellGateway(employeeServiceServerSpan)).toBe(false);
        });

        it("should return false if null is provided", () => {
            expect(TracingUtils.isFromCellGateway(null)).toBe(false);
            expect(TracingUtils.isFromCellGateway(undefined)).toBe(false);
        });
    });

    describe("isFromSystemComponent()", () => {
        it("should return true if the span is from Global Gateway", () => {
            expect(TracingUtils.isFromSystemComponent(globalGatewayServerSpan)).toBe(true);
        });

        it("should return true if the span is from a Cell Gateway", () => {
            expect(TracingUtils.isFromSystemComponent(hrCellGatewayServerSpan)).toBe(true);
        });

        it("should return true if the span is from Istio Mixer", () => {
            expect(TracingUtils.isFromSystemComponent(istioMixerServerSpan)).toBe(true);
        });

        it("should return false if the span is from a custom service", () => {
            expect(TracingUtils.isFromSystemComponent(employeeServiceServerSpan)).toBe(false);
        });

        it("should return false if null is provided", () => {
            expect(TracingUtils.isFromSystemComponent(null)).toBe(false);
            expect(TracingUtils.isFromSystemComponent(undefined)).toBe(false);
        });
    });

});
