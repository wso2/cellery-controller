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

import io.grpc.stub.StreamObserver;
import org.wso2.vick.telemetry.receiver.AttributesBag;
import org.wso2.vick.telemetry.receiver.generated.AttributesOuterClass;
import org.wso2.vick.telemetry.receiver.generated.MixerGrpc;
import org.wso2.vick.telemetry.receiver.generated.Report;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.logging.Logger;

/**
 * Telemetry Service implementation which receives the information.
 */
public class TelemetryServiceImpl extends MixerGrpc.MixerImplBase {
    private static final Logger logger = Logger.getLogger(TelemetryServiceImpl.class.getName());

    @Override
    public void report(Report.ReportRequest request,
                       StreamObserver<Report.ReportResponse> responseObserver) {
        List<AttributesBag> attributesBags = new ArrayList<>();
        AttributedDecoder attributedDecoder = new AttributedDecoder(request);
        List<AttributesOuterClass.CompressedAttributes> attributesList = request.getAttributesList();

        for (AttributesOuterClass.CompressedAttributes attributes : attributesList) {
            attributedDecoder.setCurrentAttributes(attributes);
            AttributesBag attributesBag = new AttributesBag();

            attributes.getStrings().forEach((key, value) -> {
                attributesBag.put(attributedDecoder.getValue(key), attributedDecoder.getValue(value));
            });

            attributes.getStringMaps().forEach((key, stringMap) -> {
                Map<String, String> decodedStringMap = new HashMap<>();
                stringMap.getEntries().forEach((stringMapKey, stringMapValue) -> {
                    decodedStringMap.put(attributedDecoder.getValue(stringMapKey),
                            attributedDecoder.getValue(stringMapValue));
                });
                attributesBag.put(attributedDecoder.getValue(key), decodedStringMap);
            });

            attributes.getBoolsMap().forEach((key, value) -> {
                attributesBag.put(attributedDecoder.getValue(key), value);
            });

            attributes.getInt64SMap().forEach((key, value) -> {
                attributesBag.put(attributedDecoder.getValue(key), value);
            });

            attributes.getDoublesMap().forEach((key, value) -> {
                attributesBag.put(attributedDecoder.getValue(key), value);
            });

            attributes.getBytesMap().forEach((key, value) -> {
                attributesBag.put(attributedDecoder.getValue(key), value);
            });

            attributes.getTimestampsMap().forEach((key, value) -> {
                attributesBag.put(attributedDecoder.getValue(key), value);
            });

            attributes.getDurationsMap().forEach((key, value) -> {
                attributesBag.put(attributedDecoder.getValue(key), value);
            });

            attributesBags.add(attributesBag);
        }
        logger.info(attributesBags.toString());
        Report.ReportResponse reply = Report.ReportResponse.newBuilder().build();
        responseObserver.onNext(reply);
        responseObserver.onCompleted();
    }
}
