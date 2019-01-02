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

import Constants from "../constants";
import QueryUtils from "./queryUtils";
import moment from "moment";

describe("QueryUtils", () => {
    describe("parseTime()", () => {
        const validateTime = (timeA, timeB) => expect(Math.abs(timeA.diff(timeB))).toBeLessThan(1000);

        it("should return the time object when a string formatted time is returned", () => {
            const runTest = (time) => validateTime(QueryUtils.parseTime(time),
                moment(time, Constants.Pattern.DATE_TIME));

            runTest("24 Jul 2018, 02:33 PM");
            runTest("14 Nov 2006, 07:23 AM");
        });

        it("should return the proper time when a properly formatted relative time query is provided", () => {
            validateTime(QueryUtils.parseTime("now"), moment());
            validateTime(QueryUtils.parseTime("now - 2 years"), moment().subtract("2", "years"));
            validateTime(QueryUtils.parseTime("now - 3 year"), moment().subtract("3", "years"));
            validateTime(QueryUtils.parseTime("now - 7 months"), moment().subtract("7", "months"));
            validateTime(QueryUtils.parseTime("now - 10 month"), moment().subtract("10", "months"));
            validateTime(QueryUtils.parseTime("now - 26 days"), moment().subtract("26", "days"));
            validateTime(QueryUtils.parseTime("now - 18 day"), moment().subtract("18", "days"));
            validateTime(QueryUtils.parseTime("now - 21 hours"), moment().subtract("21", "hours"));
            validateTime(QueryUtils.parseTime("now - 14 hour"), moment().subtract("14", "hour"));
            validateTime(QueryUtils.parseTime("now - 36 minutes"), moment().subtract("36", "minutes"));
            validateTime(QueryUtils.parseTime("now - 42 minute"), moment().subtract("42", "minutes"));
            validateTime(QueryUtils.parseTime("now - 51 seconds"), moment().subtract("51", "seconds"));
            validateTime(QueryUtils.parseTime("now - 37 second"), moment().subtract("37", "seconds"));
        });

        it("should return the proper time being case insensitive to the query", () => {
            validateTime(QueryUtils.parseTime("NoW"), moment());
            validateTime(QueryUtils.parseTime("NOW"), moment());
            validateTime(QueryUtils.parseTime("Now - 5 Years"), moment().subtract("5", "years"));
            validateTime(QueryUtils.parseTime("noW - 17 sEconds"), moment().subtract("17", "seconds"));
            validateTime(QueryUtils.parseTime("NOW - 6 MONTHS"), moment().subtract("6", "months"));
        });

        it("should return the proper time when a properly formatted complex relative time query is provided", () => {
            validateTime(QueryUtils.parseTime("now - 14 days 43 minutes"),
                moment().subtract("14", "days").subtract("43", "minutes"));
            validateTime(QueryUtils.parseTime("now - 5 years 3 months 1 day 21 hours 42 minutes 23 seconds"),
                moment().subtract("5", "years").subtract("3", "months").subtract("1", "days")
                    .subtract("21", "hours").subtract("42", "minutes").subtract("23", "seconds"));
        });

        it("should return the proper time when a improperly formatted relative time query is provided", () => {
            validateTime(QueryUtils.parseTime(" now"), moment());
            validateTime(QueryUtils.parseTime("now "), moment());
            validateTime(QueryUtils.parseTime(" now "), moment());
            validateTime(QueryUtils.parseTime("  now"), moment());
            validateTime(QueryUtils.parseTime("now - 2days"), moment().subtract("2", "days"));
            validateTime(QueryUtils.parseTime("now -2 days"), moment().subtract("2", "days"));
            validateTime(QueryUtils.parseTime("now-2days"), moment().subtract("2", "days"));
            validateTime(QueryUtils.parseTime("now- 2 days"), moment().subtract("2", "days"));
            validateTime(QueryUtils.parseTime(" now- 2 days"), moment().subtract("2", "days"));
            validateTime(QueryUtils.parseTime("now- 2 days "), moment().subtract("2", "days"));
            validateTime(QueryUtils.parseTime("now -  2 days "), moment().subtract("2", "days"));
            validateTime(QueryUtils.parseTime("now-2days 21hours"),
                moment().subtract("2", "days").subtract("21", "hours"));
            validateTime(QueryUtils.parseTime("now-14days43minutes"),
                moment().subtract("14", "days").subtract("43", "minutes"));
        });

        it("should return throw an error when the query format is invalid", () => {
            expect(() => QueryUtils.parseTime("2 days")).toThrow();
            expect(() => QueryUtils.parseTime("5 seconds")).toThrow();
            expect(() => QueryUtils.parseTime("26")).toThrow();
            expect(() => QueryUtils.parseTime("now + 4 minutes")).toThrow();
            expect(() => QueryUtils.parseTime("now 10 hours")).toThrow();
            expect(() => QueryUtils.parseTime("invalid date x")).toThrow();
            expect(() => QueryUtils.parseTime("")).toThrow();
            expect(() => QueryUtils.parseTime(undefined)).toThrow();
            expect(() => QueryUtils.parseTime(null)).toThrow();
        });
    });

    describe("getTimeGranularity()", () => {
        const currentTimeMilliseconds = moment().valueOf();

        const currentTime = () => moment(currentTimeMilliseconds);
        const getValidator = (expectedGranularity) => (fromTime, toTime) => {
            const timeGranularity = QueryUtils.getTimeGranularity(fromTime, toTime);
            expect(timeGranularity).toBe(expectedGranularity);
        };

        it("should return years if the time difference is greater than 3 years", () => {
            const validate = getValidator("years");

            validate(currentTime(), currentTime().add(3, "years"));
            validate(currentTime(), currentTime().add(3, "years").add(1, "milliseconds"));
            validate(currentTime(), currentTime().add(10, "years"));
            validate(currentTime(), currentTime().add(100000, "years"));
        });

        it("should return months if the time difference is between 3 months and 3 years (including 3 months)", () => {
            const validate = getValidator("months");

            validate(currentTime(), currentTime().add(3, "months"));
            validate(currentTime(), currentTime().add(3, "months").add(1, "milliseconds"));
            validate(currentTime(), currentTime().add(10, "months"));
            validate(currentTime(), currentTime().add(2, "years"));
            validate(currentTime(), currentTime().add(3, "years").subtract(1, "milliseconds"));
        });

        it("should return days if the time difference is between 3 days and 3 months (including 3 days)", () => {
            const validate = getValidator("days");

            validate(currentTime(), currentTime().add(3, "days"));
            validate(currentTime(), currentTime().add(3, "days").add(1, "milliseconds"));
            validate(currentTime(), currentTime().add(17, "days"));
            validate(currentTime(), currentTime().add(2, "months"));
            validate(currentTime(), currentTime().add(3, "months").subtract(1, "milliseconds"));
        });

        it("should return hours if the time difference is between 3 hours and 3 days (including 3 hours)", () => {
            const validate = getValidator("hours");

            validate(currentTime(), currentTime().add(3, "hours"));
            validate(currentTime(), currentTime().add(3, "hours").add(1, "milliseconds"));
            validate(currentTime(), currentTime().add(13, "hours"));
            validate(currentTime(), currentTime().add(2, "days"));
            validate(currentTime(), currentTime().add(3, "days").subtract(1, "milliseconds"));
        });

        it("should return minutes if the time difference is between 3 mins and 3 hours (including 3 mins)", () => {
            const validate = getValidator("minutes");

            validate(currentTime(), currentTime().add(3, "minutes"));
            validate(currentTime(), currentTime().add(3, "minutes").add(1, "milliseconds"));
            validate(currentTime(), currentTime().add(10, "minutes"));
            validate(currentTime(), currentTime().add(2, "hours"));
            validate(currentTime(), currentTime().add(3, "hours").subtract(1, "milliseconds"));
        });

        it("should return seconds if the time difference is lower than 3 mins", () => {
            const validate = getValidator("seconds");

            validate(currentTime(), currentTime().add(3, "minute").subtract(1, "milliseconds"));
            validate(currentTime(), currentTime().add(2, "minute"));
            validate(currentTime(), currentTime().add(10, "seconds"));
            validate(currentTime(), currentTime().add(148, "milliseconds"));
            validate(currentTime(), currentTime().add(1, "milliseconds"));
        });
    });
});
