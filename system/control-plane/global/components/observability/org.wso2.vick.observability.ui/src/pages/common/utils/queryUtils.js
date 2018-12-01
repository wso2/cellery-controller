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
import moment from "moment";

class QueryUtils {

    /**
     * Parse a time query string.
     *
     * @param {string} query Query string
     * @returns {moment.Moment} The time referred to by the time query
     */
    static parseTime(query) {
        let time = moment(query, Constants.Pattern.DATE_TIME, true);
        if (new RegExp(Constants.Pattern.Query.RELATIVE_TIME, "i").test(query)) {
            const timeRegex = new RegExp(Constants.Pattern.Query.TIME, "gi");

            // Finding the matching times
            const matches = [];
            let matchResult;
            while ((matchResult = timeRegex.exec(query))) {
                matches.push({
                    amount: matchResult[1],
                    unit: matchResult[2]
                });
            }

            // Calculating the proper time based on the query
            time = moment();
            if (matches) {
                for (let i = 0; i < matches.length; i++) {
                    const match = matches[i];
                    const amount = match.amount;
                    const unit = match.unit.toLowerCase();
                    time = time.subtract(amount, (unit.endsWith("s") ? unit : `${unit}s`));
                }
            }
        } else if (time.format() === "Invalid date") {
            throw Error("Invalid time");
        }
        return time;
    }

}

export default QueryUtils;
