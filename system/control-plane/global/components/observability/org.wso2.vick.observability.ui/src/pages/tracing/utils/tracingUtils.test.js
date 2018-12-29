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
    let globalGatewayServerSpan;
    let globalGatewayClientSpan;
    let hrCellGatewayServerSpan;
    let hrCellGatewayClientSpan;
    let employeeServiceServerSpan;
    let employeeServiceToIstioMixerClientSpan;
    let istioMixerServerSpan;
    let istioMixerWorkerSpan;
    let employeeServiceToReviewsServiceClientSpan;
    let reviewsServiceServerSpan;
    let employeeServiceToStockOptionsCellClientSpan;
    let stockOptionsCellGatewayServerSpan;
    let stockOptionsCellGatewayClientSpan;
    let stockOptionsServiceServerSpan;
    let orderedSpanList;

    const setup = () => {
        globalGatewayServerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "trace-x-id",
            parentId: "undefined",
            serviceName: "global-gateway",
            operationName: "get-hr-info",
            kind: Constants.Span.Kind.SERVER,
            startTime: 10000000,
            duration: 3160000,
            tags: JSON.stringify({keyA: "valueA", key2: "value2"})
        });
        globalGatewayClientSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-b-id",
            parentId: "trace-x-id",
            serviceName: "global-gateway",
            operationName: "call-hr-cell",
            kind: Constants.Span.Kind.CLIENT,
            startTime: 10010000,
            duration: 3110000,
            tags: JSON.stringify({key1: "value1", key2: "value2"})
        });
        hrCellGatewayServerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-b-id",
            parentId: "trace-x-id",
            cell: "hr",
            serviceName: "cell-gateway",
            operationName: "call-hr-cell",
            kind: Constants.Span.Kind.SERVER,
            startTime: 10020000,
            duration: 3090000,
            tags: JSON.stringify({component: "proxy", key1: "value1", key2: "value2"})
        });
        hrCellGatewayClientSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-c-id",
            parentId: "span-b-id",
            cell: "hr",
            serviceName: "cell-gateway",
            operationName: "get-employee-data",
            kind: Constants.Span.Kind.CLIENT,
            startTime: 10030000,
            duration: 3060000,
            tags: JSON.stringify({component: "proxy", key12: "value12", key2: "value2"})
        });
        employeeServiceServerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-c-id",
            parentId: "span-b-id",
            cell: "hr",
            serviceName: "employee-service",
            operationName: "get-employee-data",
            kind: Constants.Span.Kind.SERVER,
            startTime: 10040000,
            duration: 3040000,
            tags: JSON.stringify({component: "proxy", key1: "value1", key21: "value21"})
        });
        employeeServiceToIstioMixerClientSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-d-id",
            parentId: "span-c-id",
            cell: "hr",
            serviceName: "employee-service",
            operationName: "is-authorized",
            kind: Constants.Span.Kind.CLIENT,
            startTime: 10050000,
            duration: 990000,
            tags: JSON.stringify({component: "proxy", key1: "value1", keyF: "valueF"})
        });
        istioMixerServerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-d-id",
            parentId: "span-c-id",
            serviceName: "istio-mixer",
            operationName: "is-authorized",
            kind: Constants.Span.Kind.SERVER,
            startTime: 10060000,
            duration: 940000,
            tags: JSON.stringify({lc: "istio-mixer"})
        });
        istioMixerWorkerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-e-id",
            parentId: "span-d-id",
            serviceName: "istio-mixer",
            operationName: "authorization",
            startTime: 10070000,
            duration: 890000,
            tags: JSON.stringify({lc: "istio-mixer"})
        });
        employeeServiceToReviewsServiceClientSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-f-id",
            parentId: "span-c-id",
            cell: "hr",
            serviceName: "employee-service",
            operationName: "get-reviews",
            kind: Constants.Span.Kind.CLIENT,
            startTime: 11050000,
            duration: 990000,
            tags: JSON.stringify({component: "proxy", keyX: "valueX", key2: "value2"})
        });
        reviewsServiceServerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-f-id",
            parentId: "span-c-id",
            cell: "hr",
            serviceName: "reviews-service",
            operationName: "get-reviews",
            kind: Constants.Span.Kind.SERVER,
            startTime: 11100000,
            duration: 890000,
            tags: JSON.stringify({component: "proxy", key6: "value6", keyE: "valueE"})
        });
        employeeServiceToStockOptionsCellClientSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-g-id",
            parentId: "span-c-id",
            cell: "hr",
            serviceName: "employee-service",
            operationName: "get-employee-stock-options",
            kind: Constants.Span.Kind.CLIENT,
            startTime: 12060000,
            duration: 990000,
            tags: JSON.stringify({component: "proxy", key14: "value14", key5: "value5"})
        });
        stockOptionsCellGatewayServerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-g-id",
            parentId: "span-c-id",
            cell: "stock-options",
            serviceName: "cell-gateway",
            operationName: "get-employee-stock-options",
            kind: Constants.Span.Kind.SERVER,
            startTime: 12100000,
            duration: 890000,
            tags: JSON.stringify({component: "proxy", key1: "value1", key27: "value27"})
        });
        stockOptionsCellGatewayClientSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-h-id",
            parentId: "span-g-id",
            cell: "stock-options",
            serviceName: "cell-gateway",
            operationName: "get-employee-stock-options",
            kind: Constants.Span.Kind.CLIENT,
            startTime: 12150000,
            duration: 790000,
            tags: JSON.stringify({component: "proxy", keyX1: "valueX1", key2A: "value2A"})
        });
        stockOptionsServiceServerSpan = new Span({
            traceId: "trace-x-id",
            spanId: "span-h-id",
            parentId: "span-g-id",
            cell: "stock-options",
            serviceName: "stock-options-service",
            operationName: "get-employee-stock-options",
            kind: Constants.Span.Kind.SERVER,
            startTime: 12200000,
            duration: 690000,
            tags: JSON.stringify({component: "proxy", keyG1: "valueG1", key8: "value8"})
        });

        orderedSpanList = [
            globalGatewayServerSpan, globalGatewayClientSpan, hrCellGatewayServerSpan, hrCellGatewayClientSpan,
            employeeServiceServerSpan, employeeServiceToIstioMixerClientSpan, istioMixerServerSpan,
            istioMixerWorkerSpan, employeeServiceToReviewsServiceClientSpan, reviewsServiceServerSpan,
            employeeServiceToStockOptionsCellClientSpan, stockOptionsCellGatewayServerSpan,
            stockOptionsCellGatewayClientSpan, stockOptionsServiceServerSpan
        ];
    };

    let initialSpanDataList;
    const storeInitialSpanData = () => {
        initialSpanDataList = [];
        for (const span of orderedSpanList) {
            initialSpanDataList.push({
                traceId: span.traceId,
                spanId: span.spanId,
                parentId: span.parentId
            });
        }
    };
    const validateInitialSpanData = () => {
        expect(initialSpanDataList).toHaveLength(orderedSpanList.length);
        for (let i = 0; i < initialSpanDataList.length; i++) {
            const spanData = initialSpanDataList[i];
            const span = orderedSpanList[i];

            expect(span.traceId).toBe(spanData.traceId);
            expect(span.spanId).toBe(spanData.spanId);
            expect(span.parentId).toBe(spanData.parentId);
        }
    };

    describe("buildTree()", () => {
        beforeEach(setup);

        it("should build the tracing tree from the spans list", () => {
            // Shuffle spans list
            const spansList = orderedSpanList.map((a) => [Math.random(), a])
                .sort((a, b) => a[0] - b[0])
                .map((a) => a[1]);
            storeInitialSpanData();

            expect(TracingUtils.buildTree(spansList)).toBe(globalGatewayServerSpan);

            validateInitialSpanData();

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

        it("should build the tracing tree from the spans list when trace starts from CLIENT span", () => {
            let additionalSpan;

            const setupAdditionalSpans = () => {
                additionalSpan = new Span({
                    traceId: globalGatewayServerSpan.traceId,
                    spanId: globalGatewayServerSpan.traceId,
                    parentId: "undefined",
                    serviceName: "external-service",
                    operationName: "call-vick",
                    kind: Constants.Span.Kind.CLIENT,
                    startTime: globalGatewayServerSpan.startTime - 10,
                    duration: globalGatewayServerSpan.duration + 20,
                    tags: JSON.stringify({keyA: "valueA"})
                });
            };

            const buildAndValidate = () => {
                expect(TracingUtils.buildTree(orderedSpanList)).toBe(additionalSpan);

                // Additional span added
                expect(additionalSpan.parent).toBeNull();
                expect(additionalSpan.sibling).toBe(globalGatewayServerSpan);
                expect(additionalSpan.children.size).toBe(1);
                expect(additionalSpan.children.has(globalGatewayServerSpan)).toBe(true);

                expect(globalGatewayServerSpan.parent).toBe(additionalSpan);
                expect(globalGatewayServerSpan.sibling).toBe(additionalSpan);
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
            };

            setupAdditionalSpans();
            orderedSpanList.push(additionalSpan);
            globalGatewayServerSpan.parentId = "undefined";
            storeInitialSpanData();

            buildAndValidate();
            validateInitialSpanData();

            // Changing the order of spans and testing again
            setup();
            setupAdditionalSpans();
            orderedSpanList.splice(0, 0, additionalSpan);
            globalGatewayServerSpan.parentId = "undefined";
            storeInitialSpanData();

            buildAndValidate();
            validateInitialSpanData();
        });

        it("should build the tracing tree from the spans list when there are internal spans with kinds set", () => {
            const additionalGlobalGatewaySpan1 = new Span({
                traceId: globalGatewayServerSpan.traceId,
                spanId: `${globalGatewayServerSpan.spanId}-1`,
                parentId: globalGatewayServerSpan.spanId,
                serviceName: globalGatewayServerSpan.serviceName,
                operationName: `${globalGatewayServerSpan.serviceName}:additional-operation-1`,
                kind: Constants.Span.Kind.CLIENT,
                startTime: globalGatewayServerSpan.startTime + 1,
                duration: globalGatewayServerSpan.startTime - 2,
                tags: JSON.stringify({key3: "value3"})
            });
            const additionalGlobalGatewaySpan2 = new Span({
                traceId: globalGatewayServerSpan.traceId,
                spanId: `${globalGatewayServerSpan.spanId}-2`,
                parentId: `${globalGatewayServerSpan.spanId}-1`,
                serviceName: globalGatewayServerSpan.serviceName,
                operationName: `${globalGatewayServerSpan.serviceName}:additional-operation-2`,
                startTime: globalGatewayServerSpan.startTime + 2,
                duration: globalGatewayServerSpan.startTime - 4,
                tags: JSON.stringify({key2A: "value2A"})
            });
            orderedSpanList.push(additionalGlobalGatewaySpan1);
            orderedSpanList.push(additionalGlobalGatewaySpan2);
            globalGatewayClientSpan.parentId = additionalGlobalGatewaySpan2.spanId;
            hrCellGatewayServerSpan.parentId = additionalGlobalGatewaySpan2.spanId;

            storeInitialSpanData();

            expect(TracingUtils.buildTree(orderedSpanList)).toBe(globalGatewayServerSpan);

            validateInitialSpanData();

            expect(globalGatewayServerSpan.parent).toBeNull();
            expect(globalGatewayServerSpan.sibling).toBeNull();
            expect(globalGatewayServerSpan.children.size).toBe(1);
            expect(globalGatewayServerSpan.children.has(additionalGlobalGatewaySpan1)).toBe(true);

            // Additional global gateway 1 span added
            expect(additionalGlobalGatewaySpan1.parent).toBe(globalGatewayServerSpan);
            expect(additionalGlobalGatewaySpan1.sibling).toBeNull();
            expect(additionalGlobalGatewaySpan1.children.size).toBe(1);
            expect(additionalGlobalGatewaySpan1.children.has(additionalGlobalGatewaySpan2)).toBe(true);

            // Additional global gateway 2 span added
            expect(additionalGlobalGatewaySpan2.parent).toBe(additionalGlobalGatewaySpan1);
            expect(additionalGlobalGatewaySpan2.sibling).toBeNull();
            expect(additionalGlobalGatewaySpan2.children.size).toBe(1);
            expect(additionalGlobalGatewaySpan2.children.has(globalGatewayClientSpan)).toBe(true);

            expect(globalGatewayClientSpan.parent).toBe(additionalGlobalGatewaySpan2);
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

        it("should build the tracing tree from the spans list with the improperly connected spans", () => {
            const additionalSpan = new Span({
                traceId: hrCellGatewayClientSpan.traceId,
                spanId: `${hrCellGatewayClientSpan.spanId}-1`,
                parentId: hrCellGatewayClientSpan.parentId,
                cell: hrCellGatewayClientSpan.cell.name,
                serviceName: hrCellGatewayClientSpan.serviceName,
                operationName: `${hrCellGatewayClientSpan.serviceName}:additional-operation-1`,
                kind: Constants.Span.Kind.CLIENT,
                startTime: hrCellGatewayClientSpan.startTime - 1,
                duration: hrCellGatewayClientSpan.duration + 2,
                tags: JSON.stringify({component: "proxy", key3: "value3"})
            });
            orderedSpanList.push(additionalSpan);
            hrCellGatewayClientSpan.parentId = additionalSpan.spanId;
            employeeServiceServerSpan.parentId = additionalSpan.spanId;

            storeInitialSpanData();

            expect(TracingUtils.buildTree(orderedSpanList)).toBe(globalGatewayServerSpan);

            validateInitialSpanData();

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
            expect(hrCellGatewayServerSpan.children.has(additionalSpan)).toBe(true);

            // Additional span added
            expect(additionalSpan.parent).toBe(hrCellGatewayServerSpan);
            expect(additionalSpan.sibling).toBeNull();
            expect(additionalSpan.children.size).toBe(1);
            expect(additionalSpan.children.has(hrCellGatewayClientSpan)).toBe(true);

            expect(hrCellGatewayClientSpan.parent).toBe(additionalSpan);
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

        it("should throw an error if the tree provided spans contain more than two root spans", () => {
            stockOptionsCellGatewayServerSpan.parentId = "invalid-span-id";

            storeInitialSpanData();

            expect(() => TracingUtils.buildTree(orderedSpanList)).toThrow();

            validateInitialSpanData();
        });

        it("should throw an error if the tree provided spans contains two spans which are not siblings", () => {
            globalGatewayServerSpan.spanId = "invalid-span-id";

            storeInitialSpanData();

            expect(() => TracingUtils.buildTree(orderedSpanList)).toThrow();

            validateInitialSpanData();
        });
    });


    describe("labelSpanTree()", () => {
        beforeEach(setup);

        it("should label the necessary nodes according to cell, component type", () => {
            storeInitialSpanData();

            const rootSpan = TracingUtils.buildTree(orderedSpanList);
            TracingUtils.labelSpanTree(rootSpan);

            validateInitialSpanData();

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

    describe("getTreeRoot()", () => {
        beforeEach(setup);

        it("should return the root from the tree structure", () => {
            const rootSpan = TracingUtils.buildTree(orderedSpanList);

            storeInitialSpanData();

            expect(TracingUtils.getTreeRoot(orderedSpanList)).toBe(rootSpan);

            validateInitialSpanData();
        });

        it("should thow an error if called before building the tree structure", () => {
            storeInitialSpanData();

            expect(() => TracingUtils.getTreeRoot(orderedSpanList)).toThrow();

            validateInitialSpanData();
        });

        it("should thow an error if an array containing more than one tree root is provided", () => {
            TracingUtils.buildTree(orderedSpanList);
            stockOptionsServiceServerSpan.parent = null;
            stockOptionsCellGatewayClientSpan.children.delete(stockOptionsServiceServerSpan);

            storeInitialSpanData();

            expect(() => TracingUtils.getTreeRoot(orderedSpanList)).toThrow();

            validateInitialSpanData();
        });
    });

    describe("getOrderedList()", () => {
        beforeEach(setup);

        it("should return the nodes ordered by start time and tree structure", () => {
            storeInitialSpanData();

            const rootSpan = TracingUtils.buildTree(orderedSpanList);
            const resultList = TracingUtils.getOrderedList(rootSpan);

            validateInitialSpanData();

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
        beforeEach(setup);

        it("should remove the references in other nodes that point to this to form the tree", () => {
            storeInitialSpanData();

            TracingUtils.buildTree(orderedSpanList);
            TracingUtils.removeSpanFromTree(employeeServiceServerSpan);

            validateInitialSpanData();

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
            expect(hrCellGatewayClientSpan.sibling).toBeNull();
            expect(hrCellGatewayClientSpan.children.size).toBe(3);
            expect(hrCellGatewayClientSpan.children.has(employeeServiceToIstioMixerClientSpan)).toBe(true);
            expect(hrCellGatewayClientSpan.children.has(employeeServiceToReviewsServiceClientSpan)).toBe(true);
            expect(hrCellGatewayClientSpan.children.has(employeeServiceToStockOptionsCellClientSpan)).toBe(true);

            expect(employeeServiceToIstioMixerClientSpan.parent).toBe(hrCellGatewayClientSpan);
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

            expect(employeeServiceToReviewsServiceClientSpan.parent).toBe(hrCellGatewayClientSpan);
            expect(employeeServiceToReviewsServiceClientSpan.sibling).toBe(reviewsServiceServerSpan);
            expect(employeeServiceToReviewsServiceClientSpan.children.size).toBe(1);
            expect(employeeServiceToReviewsServiceClientSpan.children.has(reviewsServiceServerSpan)).toBe(true);

            expect(reviewsServiceServerSpan.parent).toBe(employeeServiceToReviewsServiceClientSpan);
            expect(reviewsServiceServerSpan.sibling).toBe(employeeServiceToReviewsServiceClientSpan);
            expect(reviewsServiceServerSpan.children.size).toBe(0);

            expect(employeeServiceToStockOptionsCellClientSpan.parent).toBe(hrCellGatewayClientSpan);
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

        it("should remove the references in other nodes without failing when the root is to be removed", () => {
            storeInitialSpanData();

            TracingUtils.buildTree(orderedSpanList);
            TracingUtils.removeSpanFromTree(globalGatewayServerSpan);

            validateInitialSpanData();

            expect(globalGatewayClientSpan.parent).toBeNull();
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

        it("should remove the references in other nodes without failing when a leaf is to be removed", () => {
            storeInitialSpanData();

            TracingUtils.buildTree(orderedSpanList);
            TracingUtils.removeSpanFromTree(istioMixerWorkerSpan);

            validateInitialSpanData();

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
            expect(istioMixerServerSpan.children.size).toBe(0);

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

    describe("resetTreeSpanReferences()", () => {
        beforeEach(setup);

        it("should clear all connections to other nodes", () => {
            storeInitialSpanData();

            TracingUtils.buildTree(orderedSpanList);
            TracingUtils.resetTreeSpanReferences(orderedSpanList);

            validateInitialSpanData();

            for (let i = 0; i < orderedSpanList.length; i++) {
                expect(orderedSpanList[i].children.size).toBe(0);
                expect(orderedSpanList[i].parent).toBeNull();
                expect(orderedSpanList[i].sibling).toBeNull();
            }
        });
    });
});
