package org.wso2.vick.telemetry.receiver;

import org.wso2.vick.telemetry.receiver.generated.MixerGrpc;
import org.wso2.vick.telemetry.receiver.generated.Report;

/**
 * Telemetry Service implementation which receives the information.
 */
public class TelemetryServiceImpl extends MixerGrpc.MixerImplBase {
    @Override
    public void report(Report.ReportRequest request,
                       io.grpc.stub.StreamObserver<Report.ReportResponse> responseObserver) {
        Report.ReportResponse reply = Report.ReportResponse.newBuilder().build();
        responseObserver.onNext(reply);
        responseObserver.onCompleted();
    }
}
