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
import org.wso2.vick.telemetry.receiver.generated.MixerGrpc;
import org.wso2.vick.telemetry.receiver.generated.Report;

import java.util.logging.Logger;

/**
 * Telemetry Service implementation which receives the information.
 */
public class TelemetryServiceImpl extends MixerGrpc.MixerImplBase {
    private static final Logger logger = Logger.getLogger(TelemetryServiceImpl.class.getName());

    @Override
    public void report(Report.ReportRequest request,
                       StreamObserver<Report.ReportResponse> responseObserver) {
        logger.info(request.toString());
        Report.ReportResponse reply = Report.ReportResponse.newBuilder().build();
        responseObserver.onNext(reply);
        responseObserver.onCompleted();
    }
}
