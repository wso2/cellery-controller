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

    /** @private **/
    static LOCAL_STORAGE_ITEM = "colorMap";

    constructor() {
        this.loadColorMap();
        this.initializeMainColors();
        this.listeners = [];
    }

    /**
     * Initialize the main colors in the color scheme.
     */
    initializeMainColors = () => {
        if (!this.colorMap || typeof this.colorMap !== "object") {
            this.colorMap = {};
        }

        this.colorMap[ColorGenerator.VICK] = "#ababab";
        this.colorMap[ColorGenerator.ISTIO] = "#434da1";
        this.colorMap[ColorGenerator.UNKNOWN] = "#71736f";
        this.colorMap[ColorGenerator.SUCCESS] = Green[500];
        this.colorMap[ColorGenerator.WARNING] = Yellow[800];
        this.colorMap[ColorGenerator.ERROR] = Red[500];
        this.colorMap[ColorGenerator.CLIENT_ERROR] = Blue[500];
    };

    /**
     * Add a list of keys to the current exiting keys.
     *
     * @param {Array.<string>} keys The array of keys to add to the current keys
     */
    addKeys = (keys) => {
        const self = this;
        const newKeys = keys.filter((key) => !(key in self.colorMap));
        const colors = this.generateColors(newKeys.length);

        for (let i = 0; i < newKeys.length; i++) {
            self.colorMap[newKeys[i]] = colors[i];
        }
        this.persistColorMap();
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
     * Generate a set of colors.
     *
     * @private
     * @param {number} count The number of colors to generate
     * @returns {Array.<string>} The colors generated
     */
    generateColors = (count) => {
        const newColors = [];
        let colorsLeftCount = count;
        while (colorsLeftCount > 0) {
            const generatedColors = randomColor({
                luminosity: "light",
                count: colorsLeftCount
            });

            // Verifying that the colors are distinct and adding them to the list
            for (const generatedColor of generatedColors) {
                if (!this.colorMap[generatedColor] && !newColors.includes(generatedColor)) {
                    newColors.push(generatedColor);
                    colorsLeftCount -= 1;
                }
            }
        }
        return newColors;
    };

    /**
     * Reset the color scheme in memory and in persistent storage.
     */
    resetColors = () => {
        localStorage.removeItem(ColorGenerator.LOCAL_STORAGE_ITEM);
        this.initializeMainColors();
        this.notify();
    };

    /**
     * Add a listener to listen to color map changes.
     *
     * @param {Function} callback Callback to be called upon color map changes.
     */
    addListener = (callback) => {
        this.listeners.push(callback);
    };

    /**
     * Remove a listener which was added before to listen to color map changes.
     *
     * @param {Function} callback The callback to be removed
     */
    removeListener = (callback) => {
        const removeIndex = this.listeners.indexOf(callback);
        this.listeners.splice(removeIndex, 1);
    };

    /**
     * Notify the listeners about the color map change.
     *
     * @private
     */
    notify = () => {
        this.listeners.forEach((listener) => listener());
    };

    /**
     * Persist the color map in a persistent storage.
     * The Browser local storage is currently used.
     *
     * @private
     */
    persistColorMap = () => {
        localStorage.setItem(ColorGenerator.LOCAL_STORAGE_ITEM, JSON.stringify(this.colorMap));
    };

    /**
     * Load the color man stored in a persistent storage.
     *
     * @private
     */
    loadColorMap = () => {
        this.colorMap = JSON.parse(localStorage.getItem(ColorGenerator.LOCAL_STORAGE_ITEM));
    };

}

export default ColorGenerator;
