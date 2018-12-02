
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

/* eslint prefer-promise-reject-errors: ["off"] */

import AuthUtils from "./authUtils";
import {ConfigHolder} from "../config";
import HttpUtils from "./httpUtils";
import axios from "axios";

jest.mock("axios");

describe("HttpUtils", () => {
    describe("parseQueryParams()", () => {
        it("should parse and return the query parameters as an object", () => {
            const query = HttpUtils.parseQueryParams("?key1=value1&key2=value2&key3=value3");

            expect(Object.keys(query)).toHaveLength(3);
            expect(query.key1).toBe("value1");
            expect(query.key2).toBe("value2");
            expect(query.key3).toBe("value3");
        });

        it("should parse and return the query if the query string doesn't have the '?' prefix", () => {
            const query = HttpUtils.parseQueryParams("key1=value1&key2=value2&key3=value3");

            expect(Object.keys(query)).toHaveLength(3);
            expect(query.key1).toBe("value1");
            expect(query.key2).toBe("value2");
            expect(query.key3).toBe("value3");
        });

        it("should parse and return the decoded query parameters if the query parameters had been URL encoded", () => {
            const query = HttpUtils.parseQueryParams("?key1=value%201&key2=value%2D2&key%213=value3");

            expect(Object.keys(query)).toHaveLength(3);
            expect(query.key1).toBe("value 1");
            expect(query.key2).toBe("value-2");
            expect(query["key!3"]).toBe("value3");
        });

        it("should true for the parameters which has a key but doesn't have a value", () => {
            const query = HttpUtils.parseQueryParams("?key1=&key2");

            expect(Object.keys(query)).toHaveLength(2);
            expect(query.key1).toBe(true);
            expect(query.key2).toBe(true);
        });

        it("should return object with no keys without failing if empty query is provided", () => {
            expect(Object.keys(HttpUtils.parseQueryParams("?"))).toHaveLength(0);
            expect(Object.keys(HttpUtils.parseQueryParams(""))).toHaveLength(0);
            expect(Object.keys(HttpUtils.parseQueryParams("?="))).toHaveLength(0);
            expect(Object.keys(HttpUtils.parseQueryParams("?&&=&&"))).toHaveLength(0);
            expect(Object.keys(HttpUtils.parseQueryParams(undefined))).toHaveLength(0);
            expect(Object.keys(HttpUtils.parseQueryParams(null))).toHaveLength(0);
        });
    });

    describe("callBackendAPI()", () => {
        const backEndURL = "http://www.example.com";
        let config;

        beforeEach(() => {
            config = new ConfigHolder();
            config.set(ConfigHolder.USER, "user1");
            config.set(ConfigHolder.BACKEND_URL, backEndURL);
        });

        it("should add application/json header", async () => {
            axios.mockImplementation((config) => {
                expect(Object.keys(config)).toHaveLength(3);
                expect(config.method).toBe("GET");
                expect(config.url).toBe(`${backEndURL}/test`);
                expect(Object.keys(config.headers)).toHaveLength(3);
                expect(config.headers["X-Key"]).toBe("value");
                expect(config.headers.Accept).toBe("application/json");
                expect(config.headers["Content-Type"]).toBe("application/json");

                return new Promise((resolve) => {
                    resolve({
                        status: 200,
                        data: [
                            {
                                event: {
                                    value: "testData"
                                }
                            }
                        ]
                    });
                });
            });

            await HttpUtils.callBackendAPI({
                method: "GET",
                url: "/test",
                headers: {
                    "X-Key": "value"
                }
            }, config);
        });

        it("should add application/json Accept and Content-Type headers if no headers are provided", async () => {
            axios.mockImplementation((config) => {
                expect(Object.keys(config)).toHaveLength(3);
                expect(config.method).toBe("GET");
                expect(config.url).toBe(`${backEndURL}/test`);
                expect(Object.keys(config.headers)).toHaveLength(2);
                expect(config.headers.Accept).toBe("application/json");
                expect(config.headers["Content-Type"]).toBe("application/json");

                return new Promise((resolve) => {
                    resolve({
                        status: 200,
                        data: [
                            {
                                event: {
                                    value: "testData"
                                }
                            }
                        ]
                    });
                });
            });

            await HttpUtils.callBackendAPI({
                method: "GET",
                url: "/test"
            }, config);
        });

        it("should not change headers if Accept header is already provided", async () => {
            axios.mockImplementation((config) => {
                expect(Object.keys(config)).toHaveLength(3);
                expect(config.method).toBe("GET");
                expect(config.url).toBe(`${backEndURL}/test`);
                expect(Object.keys(config.headers)).toHaveLength(2);
                expect(config.headers.Accept).toBe("application/xml");
                expect(config.headers["Content-Type"]).toBe("application/json");

                return new Promise((resolve) => {
                    resolve({
                        status: 200,
                        data: [
                            {
                                event: {
                                    value: "testData"
                                }
                            }
                        ]
                    });
                });
            });

            await HttpUtils.callBackendAPI({
                method: "GET",
                url: "/test",
                headers: {
                    Accept: "application/xml"
                }
            }, config);
        });

        it("should add application/json header if Content-Type header is provided", async () => {
            axios.mockImplementation((config) => {
                expect(Object.keys(config)).toHaveLength(3);
                expect(config.method).toBe("GET");
                expect(config.url).toBe(`${backEndURL}/test`);
                expect(Object.keys(config.headers)).toHaveLength(2);
                expect(config.headers.Accept).toBe("application/json");
                expect(config.headers["Content-Type"]).toBe("application/json");

                return new Promise((resolve) => {
                    resolve({
                        status: 200,
                        data: [
                            {
                                event: {
                                    value: "testData"
                                }
                            }
                        ]
                    });
                });
            });

            await HttpUtils.callBackendAPI({
                method: "GET",
                url: "/test"
            }, config);
        });

        it("should not change headers if Content-Type header is already provided", async () => {
            axios.mockImplementation((config) => {
                expect(Object.keys(config)).toHaveLength(3);
                expect(config.method).toBe("GET");
                expect(config.url).toBe(`${backEndURL}/test`);
                expect(Object.keys(config.headers)).toHaveLength(2);
                expect(config.headers.Accept).toBe("application/json");
                expect(config.headers["Content-Type"]).toBe("application/xml");

                return new Promise((resolve) => {
                    resolve({
                        status: 200,
                        data: [
                            {
                                event: {
                                    value: "testData"
                                }
                            }
                        ]
                    });
                });
            });

            await HttpUtils.callBackendAPI({
                method: "GET",
                url: "/test",
                headers: {
                    "Content-Type": "application/xml"
                }
            }, config);
        });

        const ERROR_DATA = "testError";
        const AXIOS_OUTPUT_DATA = [
            {
                event: "testEvent"
            }
        ];
        const SUCCESS_OUTPUT_DATA = [
            "testEvent"
        ];
        const resolveStatusCodes = [
            200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
            300, 301, 302, 303, 304, 305, 306, 307, 308
        ];
        const rejectStatusCodes = [
            400, 401, 402, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 418,
            421, 422, 423, 424, 426, 428, 429, 431, 451,
            500, 501, 502, 503, 504, 505, 506, 507, 508, 510, 511
        ];

        const mockResolve = (statusCode) => {
            axios.mockResolvedValue(new Promise((resolve) => {
                resolve({
                    status: statusCode,
                    data: AXIOS_OUTPUT_DATA
                });
            }));

            return HttpUtils.callBackendAPI({
                method: "GET",
                url: "/test",
                headers: {
                    Accept: "application/xml"
                }
            }, config);
        };
        const mockReject = (statusCode) => {
            axios.mockResolvedValue(new Promise((resolve, reject) => {
                reject({
                    response: {
                        status: statusCode,
                        data: ERROR_DATA
                    }
                });
            }));

            return HttpUtils.callBackendAPI({
                method: "GET",
                url: "/test",
                headers: {
                    Accept: "application/xml"
                }
            }, config);
        };

        resolveStatusCodes.forEach((statusCode) => {
            it(`should resolve with response data when axios resolves with a ${statusCode} status code`, () => {
                expect.assertions(1);
                return expect(mockResolve(statusCode)).resolves.toEqual(SUCCESS_OUTPUT_DATA);
            });
        });

        rejectStatusCodes.forEach((statusCode) => {
            it(`should reject with response data when axios resolves with a ${statusCode} status code`, () => {
                expect.assertions(1);
                return expect(mockResolve(statusCode)).rejects.toEqual(AXIOS_OUTPUT_DATA);
            });
        });

        rejectStatusCodes.filter((statusCode) => statusCode !== 401).forEach((statusCode) => {
            it(`should reject with response data when axios rejects with a ${statusCode} status code`, () => {
                expect.assertions(1);
                return expect(mockReject(statusCode)).rejects.toEqual(new Error(ERROR_DATA));
            });
        });

        it("should sign out and reject with response when axios rejects with a 401 status code", async () => {
            const spy = jest.spyOn(AuthUtils, "signOut");

            await expect(mockReject(401)).rejects.toEqual(new Error(ERROR_DATA));
            expect(spy).toHaveBeenCalledTimes(1);
        });
    });
});
