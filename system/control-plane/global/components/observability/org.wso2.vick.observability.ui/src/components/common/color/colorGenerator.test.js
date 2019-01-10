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

import ColorGenerator from "./colorGenerator";
import StateHolder from "../state/stateHolder";

describe("ColorGenerator", () => {
    const INITIAL_KEYS = [
        ColorGenerator.VICK, ColorGenerator.ISTIO, ColorGenerator.UNKNOWN, ColorGenerator.SUCCESS,
        ColorGenerator.WARNING, ColorGenerator.ERROR, ColorGenerator.CLIENT_ERROR
    ];
    const INITIAL_KEY_COUNT = INITIAL_KEYS.length;

    const validateInitialKeys = (colorGenerator) => {
        for (const key of INITIAL_KEYS) {
            expect(colorGenerator.colorMap[key]).not.toBeUndefined();
        }
    };

    describe("constructor()", () => {
        it("should have VICK and Istio keys by default", () => {
            const colorGenerator = new ColorGenerator();

            expect(Object.keys(colorGenerator.colorMap)).toHaveLength(INITIAL_KEY_COUNT);
            validateInitialKeys(colorGenerator);
        });
    });

    describe("addKeys()", () => {
        const KEY_1 = "key1";
        const KEY_2 = "key2";
        const KEY_3 = "key3";
        const KEY_4 = "key4";
        const KEY_5 = "key5";
        const keyList = [KEY_1, KEY_2, KEY_3, KEY_4, KEY_5];

        it("should add the key set provided", () => {
            const colorGenerator = new ColorGenerator();
            colorGenerator.addKeys(keyList);

            expect(Object.keys(colorGenerator.colorMap)).toHaveLength(keyList.length + INITIAL_KEY_COUNT);
            validateInitialKeys(colorGenerator);
            expect(colorGenerator.colorMap[KEY_1]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_2]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_3]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_4]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_5]).not.toBeUndefined();
        });

        it("should not add duplicate keys nor change the color", () => {
            const colorGenerator = new ColorGenerator();
            colorGenerator.addKeys(keyList);

            const key2Color = colorGenerator.colorMap[KEY_2];
            const key5Color = colorGenerator.colorMap[KEY_5];
            colorGenerator.addKeys([KEY_2, KEY_5]);

            expect(Object.keys(colorGenerator.colorMap)).toHaveLength(keyList.length + INITIAL_KEY_COUNT);
            validateInitialKeys(colorGenerator);
            expect(colorGenerator.colorMap[KEY_1]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_3]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_4]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_2]).toBe(key2Color);
            expect(colorGenerator.colorMap[KEY_5]).toBe(key5Color);
        });
    });

    describe("getColor()", () => {
        const keyCount = 200;
        let keyList;

        beforeEach(() => {
            keyList = [];
            for (const key of INITIAL_KEYS) {
                keyList.push(key);
            }
            for (let i = 0; i < keyCount; i++) {
                keyList.push(`key${i}`);
            }
        });

        it("should return colors different to one another", () => {
            const colorGenerator = new ColorGenerator();
            colorGenerator.addKeys(keyList);
            const spy = jest.spyOn(colorGenerator, "addKeys");

            let similarColors = 0;
            for (let i = 0; i < keyList.length; i++) {
                for (let j = 0; j < keyList.length; j++) {
                    if (i !== j) {
                        const colorA = colorGenerator.getColor(keyList[i]);
                        const colorB = colorGenerator.getColor(keyList[j]);
                        if (colorA === colorB) {
                            similarColors += 1;
                        }
                    }
                }
            }
            expect(similarColors).toBe(0);
            expect(spy).not.toHaveBeenCalledTimes(keyCount);
        });

        it("should return the same color for the same key when invoked multiple times", () => {
            const colorGenerator = new ColorGenerator();
            const spy = jest.spyOn(colorGenerator, "addKeys");

            const colors = [];
            for (let i = 0; i < keyList.length; i++) {
                colors.push(colorGenerator.getColor(keyList[i]));
            }
            for (let i = 0; i < keyList.length; i++) {
                const newColor = colorGenerator.getColor(keyList[i]);
                expect(colors[i]).toBe(newColor);
            }
            expect(spy).toHaveBeenCalledTimes(keyCount);
        });

        it("should add the key if an unknown key is provided", () => {
            const colorGenerator = new ColorGenerator();
            const spy = jest.spyOn(colorGenerator, "addKeys");

            for (let i = 0; i < keyList.length; i++) {
                const key = keyList[i];
                const color = colorGenerator.getColor(key);
                expect(color).not.toBeUndefined();
                expect(colorGenerator.colorMap[key]).not.toBeUndefined();
            }
            expect(spy).toHaveBeenCalledTimes(keyCount);
        });
    });

    describe("getColorForPercentage()", () => {
        const globalState = new StateHolder();
        const colorGenerator = new ColorGenerator();
        globalState.set(StateHolder.CONFIG, {
            percentageRangeMinValue: {
                errorThreshold: 0.5,
                warningThreshold: 0.7
            }
        });

        const unknownColor = colorGenerator.colorMap[ColorGenerator.UNKNOWN];
        const errorColor = colorGenerator.colorMap[ColorGenerator.ERROR];
        const warningColor = colorGenerator.colorMap[ColorGenerator.WARNING];
        const successColor = colorGenerator.colorMap[ColorGenerator.SUCCESS];

        it("should return unknown color if a value is less than 0 or greater than 1 is provided", () => {
            expect(colorGenerator.getColorForPercentage(-1, globalState)).toBe(unknownColor);
            expect(colorGenerator.getColorForPercentage(2, globalState)).toBe(unknownColor);
        });

        it("should return error color if a value 0 and 0.5 is provided (including 0)", () => {
            expect(colorGenerator.getColorForPercentage(0, globalState)).toBe(errorColor);
            expect(colorGenerator.getColorForPercentage(0.3, globalState)).toBe(errorColor);
        });

        it("should return warning color if a value between 0.5 and 0.7 is provided (including 0.5)", () => {
            expect(colorGenerator.getColorForPercentage(0.5, globalState)).toBe(warningColor);
            expect(colorGenerator.getColorForPercentage(0.6, globalState)).toBe(warningColor);
        });

        it("should return success color if a value between 0.7 and 1 is provided (including 0.7 and 1)", () => {
            expect(colorGenerator.getColorForPercentage(0.7, globalState)).toBe(successColor);
            expect(colorGenerator.getColorForPercentage(0.8, globalState)).toBe(successColor);
            expect(colorGenerator.getColorForPercentage(1, globalState)).toBe(successColor);
        });
    });
});
