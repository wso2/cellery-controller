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

const Constants = {};

// Pattern Constants
{
    const Pattern = {
        DATE_TIME: "YYYY-MM-DD HH:mm:ss"
    };

    // Query Pattern Constants
    {
        const Query = {
            SECONDS: "second(?:s)?",
            MINUTES: "minute(?:s)?",
            HOURS: "hour(?:s)?",
            DAYS: "day(?:s)?",
            MONTHS: "month(?:s)?",
            YEARS: "year(?:s)?"
        };
        Query.TIME_UNIT
            = `${Query.YEARS}|${Query.MONTHS}|${Query.DAYS}|${Query.HOURS}|${Query.MINUTES}|${Query.SECONDS}`;
        Query.TIME = `([0-9]+)\\s*(${Query.TIME_UNIT})`;
        Query.RELATIVE_TIME = `^\\s*now\\s*(?:-\\s*(?:${Query.TIME}\\s*)+)?$`;

        Pattern.Query = Query;
    }

    Constants.Pattern = Pattern;
}

export default Constants;
