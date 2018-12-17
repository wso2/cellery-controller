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

import Blue from "@material-ui/core/colors/blue";
import Green from "@material-ui/core/colors/green";
import Red from "@material-ui/core/colors/red";
import {StateHolder} from "../state";
import Yellow from "@material-ui/core/colors/yellow";
import randomColor from "randomcolor";

/**
 * Color Generator.
 */
class ColorGenerator {

    static VICK = "VICK";
    static ISTIO = "Istio";
    static UNKNOWN = "UNKNOWN";
    static SUCCESS = "SUCCESS";
    static WARNING = "WARNING";
    static ERROR = "ERROR";
    static CLIENT_ERROR = "CLIENT_ERROR";

    constructor() {
        this.colorMap = {
            [ColorGenerator.VICK]: "#a53288",
            [ColorGenerator.ISTIO]: "#434da1",
            [ColorGenerator.UNKNOWN]: "#71736f",
            [ColorGenerator.SUCCESS]: Green[500],
            [ColorGenerator.WARNING]: Yellow[800],
            [ColorGenerator.ERROR]: Red[500],
            [ColorGenerator.CLIENT_ERROR]: Blue[500]
        };
    }

    /**
     * Add a list of keys to the current exiting keys.
     *
     * @param {Array.<string>} keys The array of keys to add to the current keys
     */
    addKeys = (keys = []) => {
        const self = this;
        const newKeys = keys.filter((key) => !(key in self.colorMap));
        const colors = ColorGenerator.generateColors(newKeys.length);

        for (let i = 0; i < newKeys.length; i++) {
            self.colorMap[newKeys[i]] = colors[i];
        }
    };

    /**
     * Get the color for a particular key.
     * If the key does not already exist the key will be added and a new color scheme will be generated.
     *
     * @param {string} key The name of the key
     * @returns {string} Hex value for a particular color
     */
    getColor = (key) => {
        if (!(key in this.colorMap)) {
            this.addKeys([key]);
        }
        return this.colorMap[key];
    };

    /**
     * Get the colors for percentage.
     *
     * @param {number} percentage The percentage for which the color is required
     * @param {StateHolder} globalState The global state
     * @returns {string} The color for the percentage
     */
    getColorForPercentage = (percentage, globalState) => {
        let colorKey = ColorGenerator.SUCCESS;
        if (percentage < globalState.get(StateHolder.CONFIG).percentageRangeMinValue.warningThreshold) {
            colorKey = ColorGenerator.WARNING;
        }
        if (percentage < globalState.get(StateHolder.CONFIG).percentageRangeMinValue.errorThreshold) {
            colorKey = ColorGenerator.ERROR;
        }
        if (percentage < 0 || percentage > 1) {
            colorKey = ColorGenerator.UNKNOWN;
        }
        return this.colorMap[colorKey];
    };

    /**
     * Regenerate a new color scheme for the existing keys.
     * This will remove all the previous colors used and generate a new set of colors.
     */
    regenerateNewColorScheme = () => {
        const keyCount = Object.keys(this.colorMap).length;
        const colors = ColorGenerator.generateColors(keyCount);

        let i = 0;
        for (const key in this.colorMap) {
            if (this.colorMap.hasOwnProperty(key)) {
                this.colorMap[key] = colors[i];
                i += 1;
            }
        }
    };

    /**
     * Generate a set of colors.
     *
     * @private
     * @param {number} count The number of colors to generate
     * @returns {Array.<string>} The colors generated
     */
    static generateColors = (count) => randomColor({
        luminosity: "light",
        count: count
    });

}

export default ColorGenerator;
