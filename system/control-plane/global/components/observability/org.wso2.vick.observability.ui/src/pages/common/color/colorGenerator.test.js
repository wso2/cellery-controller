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

describe("ColorGenerator", () => {
    const INITIAL_KEY_COUNT = 7;

    describe("constructor()", () => {
        it("should have VICK and Istio keys by default", () => {
            const colorGenerator = new ColorGenerator();

            expect(Object.keys(colorGenerator.colorMap)).toHaveLength(INITIAL_KEY_COUNT);
            expect(colorGenerator.colorMap[ColorGenerator.VICK]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.ISTIO]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.UNKNOWN]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.SUCCESS]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.WARNING]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.ERROR]).not.toBeUndefined();
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
            expect(colorGenerator.colorMap[ColorGenerator.VICK]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.ISTIO]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.UNKNOWN]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.SUCCESS]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.WARNING]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.ERROR]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_1]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_2]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_3]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_4]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_5]).not.toBeUndefined();
        });

        it("should not add duplicate keys", () => {
            const colorGenerator = new ColorGenerator();
            colorGenerator.addKeys(keyList);
            colorGenerator.addKeys([KEY_2, KEY_5]);

            expect(Object.keys(colorGenerator.colorMap)).toHaveLength(keyList.length + INITIAL_KEY_COUNT);
            expect(colorGenerator.colorMap[ColorGenerator.VICK]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.ISTIO]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.UNKNOWN]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.SUCCESS]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.WARNING]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.ERROR]).not.toBeUndefined();
            expect(colorGenerator.colorMap[ColorGenerator.CLIENT_ERROR]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_1]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_2]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_3]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_4]).not.toBeUndefined();
            expect(colorGenerator.colorMap[KEY_5]).not.toBeUndefined();
        });
    });

    describe("getColor()", () => {
        const keyCount = 200;
        let keyList;

        beforeEach(() => {
            keyList = [];
            keyList.push(ColorGenerator.VICK);
            keyList.push(ColorGenerator.ISTIO);
            keyList.push(ColorGenerator.ERROR);
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
            expect(similarColors).toBeLessThan(5); // There are few overlaps
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

    describe("regenerateNewColorScheme()", () => {
        const keyCount = 200;
        let keyList;

        beforeEach(() => {
            keyList = [];
            keyList.push(ColorGenerator.VICK);
            keyList.push(ColorGenerator.ISTIO);
            keyList.push(ColorGenerator.UNKNOWN);
            keyList.push(ColorGenerator.SUCCESS);
            keyList.push(ColorGenerator.WARNING);
            keyList.push(ColorGenerator.ERROR);
            keyList.push(ColorGenerator.CLIENT_ERROR);
            for (let i = 0; i < keyCount; i++) {
                keyList.push(`key${i}`);
            }
        });

        it("should generate colors for all the existing colors", () => {
            const colorGenerator = new ColorGenerator();
            colorGenerator.addKeys(keyList);
            const spy = jest.spyOn(ColorGenerator, "generateColors");
            colorGenerator.regenerateNewColorScheme();

            expect(spy).toHaveBeenCalledTimes(1);
            expect(spy).toHaveBeenCalledWith(keyList.length);
        });
    });
});
