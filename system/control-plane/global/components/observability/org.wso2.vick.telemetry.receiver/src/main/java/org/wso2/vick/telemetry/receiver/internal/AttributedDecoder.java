/*
 *  Copyright (c) 2018 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 *
 */

package org.wso2.vick.telemetry.receiver.internal;

import org.wso2.vick.telemetry.receiver.generated.AttributesOuterClass;
import org.wso2.vick.telemetry.receiver.generated.Report;

/**
 * This class uses the global, request, and attribute dictionary to decode the report request key/value.
 */
public class AttributedDecoder {

    private GlobalDictionary globalDictionary = GlobalDictionary.getInstance();
    private Report.ReportRequest reportRequest;
    private AttributesOuterClass.CompressedAttributes currentAttributes;

    public AttributedDecoder(Report.ReportRequest reportRequest) {
        this.reportRequest = reportRequest;
    }

    public String getValue(int index) {
        if (index >= 0) {
            return globalDictionary.getValue(index);
        } else {
            if (currentAttributes.getWordsCount() == 0) {
                return reportRequest.getDefaultWords((-index) - 1);
            } else {
                return currentAttributes.getWords((-index) - 1);
            }
        }
    }

    public void setCurrentAttributes(AttributesOuterClass.CompressedAttributes compressedAttributes) {
        this.currentAttributes = compressedAttributes;
    }
}
