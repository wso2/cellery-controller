package org.wso2.vick.telemetry.receiver.generated;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 * <pre>
 * Mixer provides three generated features:
 * 
 * - *Precondition Checking*. Enables callers to verify a number of preconditions
 * before responding to an incoming request from a service consumer.
 * Preconditions can include whether the service consumer is properly
 * authenticated, is on the service’s whitelist, passes ACL checks, and more.
 * - *Quota Management*. Enables services to allocate and free quota on a number
 * of dimensions, Quotas are used as a relatively simple resource management tool
 * to provide some fairness between service consumers when contending for limited
 * resources. Rate limits are examples of quotas.
 * - *Telemetry Reporting*. Enables services to report logging and monitoring.
 * In the future, it will also enable tracing and billing streams intended for
 * both the service operator as well as for service consumers.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.15.0)",
    comments = "Source: mixer/v1/service.proto")
public final class MixerGrpc {

  private MixerGrpc() {}

  public static final String SERVICE_NAME = "istio.mixer.v1.Mixer";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<Check.CheckRequest,
      Check.CheckResponse> getCheckMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Check",
      requestType = Check.CheckRequest.class,
      responseType = Check.CheckResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<Check.CheckRequest,
      Check.CheckResponse> getCheckMethod() {
    io.grpc.MethodDescriptor<Check.CheckRequest, Check.CheckResponse> getCheckMethod;
    if ((getCheckMethod = MixerGrpc.getCheckMethod) == null) {
      synchronized (MixerGrpc.class) {
        if ((getCheckMethod = MixerGrpc.getCheckMethod) == null) {
          MixerGrpc.getCheckMethod = getCheckMethod = 
              io.grpc.MethodDescriptor.<Check.CheckRequest, Check.CheckResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "istio.mixer.v1.Mixer", "Check"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Check.CheckRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Check.CheckResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new MixerMethodDescriptorSupplier("Check"))
                  .build();
          }
        }
     }
     return getCheckMethod;
  }

  private static volatile io.grpc.MethodDescriptor<Report.ReportRequest,
      Report.ReportResponse> getReportMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Report",
      requestType = Report.ReportRequest.class,
      responseType = Report.ReportResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<Report.ReportRequest,
      Report.ReportResponse> getReportMethod() {
    io.grpc.MethodDescriptor<Report.ReportRequest, Report.ReportResponse> getReportMethod;
    if ((getReportMethod = MixerGrpc.getReportMethod) == null) {
      synchronized (MixerGrpc.class) {
        if ((getReportMethod = MixerGrpc.getReportMethod) == null) {
          MixerGrpc.getReportMethod = getReportMethod = 
              io.grpc.MethodDescriptor.<Report.ReportRequest, Report.ReportResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "istio.mixer.v1.Mixer", "Report"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Report.ReportRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  Report.ReportResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new MixerMethodDescriptorSupplier("Report"))
                  .build();
          }
        }
     }
     return getReportMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static MixerStub newStub(io.grpc.Channel channel) {
    return new MixerStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static MixerBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new MixerBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static MixerFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new MixerFutureStub(channel);
  }

  /**
   * <pre>
   * Mixer provides three generated features:
   * 
   * - *Precondition Checking*. Enables callers to verify a number of preconditions
   * before responding to an incoming request from a service consumer.
   * Preconditions can include whether the service consumer is properly
   * authenticated, is on the service’s whitelist, passes ACL checks, and more.
   * - *Quota Management*. Enables services to allocate and free quota on a number
   * of dimensions, Quotas are used as a relatively simple resource management tool
   * to provide some fairness between service consumers when contending for limited
   * resources. Rate limits are examples of quotas.
   * - *Telemetry Reporting*. Enables services to report logging and monitoring.
   * In the future, it will also enable tracing and billing streams intended for
   * both the service operator as well as for service consumers.
   * </pre>
   */
  public static abstract class MixerImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * Checks preconditions and allocate quota before performing an operation.
     * The preconditions enforced depend on the set of supplied attributes and
     * the active configuration.
     * </pre>
     */
    public void check(Check.CheckRequest request,
                      io.grpc.stub.StreamObserver<Check.CheckResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getCheckMethod(), responseObserver);
    }

    /**
     * <pre>
     * Reports telemetry, such as logs and metrics.
     * The reported information depends on the set of supplied attributes and the
     * active configuration.
     * </pre>
     */
    public void report(Report.ReportRequest request,
                       io.grpc.stub.StreamObserver<Report.ReportResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getReportMethod(), responseObserver);
    }

    @Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCheckMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                Check.CheckRequest,
                Check.CheckResponse>(
                  this, METHODID_CHECK)))
          .addMethod(
            getReportMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                Report.ReportRequest,
                Report.ReportResponse>(
                  this, METHODID_REPORT)))
          .build();
    }
  }

  /**
   * <pre>
   * Mixer provides three generated features:
   * 
   * - *Precondition Checking*. Enables callers to verify a number of preconditions
   * before responding to an incoming request from a service consumer.
   * Preconditions can include whether the service consumer is properly
   * authenticated, is on the service’s whitelist, passes ACL checks, and more.
   * - *Quota Management*. Enables services to allocate and free quota on a number
   * of dimensions, Quotas are used as a relatively simple resource management tool
   * to provide some fairness between service consumers when contending for limited
   * resources. Rate limits are examples of quotas.
   * - *Telemetry Reporting*. Enables services to report logging and monitoring.
   * In the future, it will also enable tracing and billing streams intended for
   * both the service operator as well as for service consumers.
   * </pre>
   */
  public static final class MixerStub extends io.grpc.stub.AbstractStub<MixerStub> {
    private MixerStub(io.grpc.Channel channel) {
      super(channel);
    }

    private MixerStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @Override
    protected MixerStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new MixerStub(channel, callOptions);
    }

    /**
     * <pre>
     * Checks preconditions and allocate quota before performing an operation.
     * The preconditions enforced depend on the set of supplied attributes and
     * the active configuration.
     * </pre>
     */
    public void check(Check.CheckRequest request,
                      io.grpc.stub.StreamObserver<Check.CheckResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCheckMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Reports telemetry, such as logs and metrics.
     * The reported information depends on the set of supplied attributes and the
     * active configuration.
     * </pre>
     */
    public void report(Report.ReportRequest request,
                       io.grpc.stub.StreamObserver<Report.ReportResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getReportMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Mixer provides three generated features:
   * 
   * - *Precondition Checking*. Enables callers to verify a number of preconditions
   * before responding to an incoming request from a service consumer.
   * Preconditions can include whether the service consumer is properly
   * authenticated, is on the service’s whitelist, passes ACL checks, and more.
   * - *Quota Management*. Enables services to allocate and free quota on a number
   * of dimensions, Quotas are used as a relatively simple resource management tool
   * to provide some fairness between service consumers when contending for limited
   * resources. Rate limits are examples of quotas.
   * - *Telemetry Reporting*. Enables services to report logging and monitoring.
   * In the future, it will also enable tracing and billing streams intended for
   * both the service operator as well as for service consumers.
   * </pre>
   */
  public static final class MixerBlockingStub extends io.grpc.stub.AbstractStub<MixerBlockingStub> {
    private MixerBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private MixerBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @Override
    protected MixerBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new MixerBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * Checks preconditions and allocate quota before performing an operation.
     * The preconditions enforced depend on the set of supplied attributes and
     * the active configuration.
     * </pre>
     */
    public Check.CheckResponse check(Check.CheckRequest request) {
      return blockingUnaryCall(
          getChannel(), getCheckMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Reports telemetry, such as logs and metrics.
     * The reported information depends on the set of supplied attributes and the
     * active configuration.
     * </pre>
     */
    public Report.ReportResponse report(Report.ReportRequest request) {
      return blockingUnaryCall(
          getChannel(), getReportMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Mixer provides three generated features:
   * 
   * - *Precondition Checking*. Enables callers to verify a number of preconditions
   * before responding to an incoming request from a service consumer.
   * Preconditions can include whether the service consumer is properly
   * authenticated, is on the service’s whitelist, passes ACL checks, and more.
   * - *Quota Management*. Enables services to allocate and free quota on a number
   * of dimensions, Quotas are used as a relatively simple resource management tool
   * to provide some fairness between service consumers when contending for limited
   * resources. Rate limits are examples of quotas.
   * - *Telemetry Reporting*. Enables services to report logging and monitoring.
   * In the future, it will also enable tracing and billing streams intended for
   * both the service operator as well as for service consumers.
   * </pre>
   */
  public static final class MixerFutureStub extends io.grpc.stub.AbstractStub<MixerFutureStub> {
    private MixerFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private MixerFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @Override
    protected MixerFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new MixerFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * Checks preconditions and allocate quota before performing an operation.
     * The preconditions enforced depend on the set of supplied attributes and
     * the active configuration.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<Check.CheckResponse> check(
        Check.CheckRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getCheckMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Reports telemetry, such as logs and metrics.
     * The reported information depends on the set of supplied attributes and the
     * active configuration.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<Report.ReportResponse> report(
        Report.ReportRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getReportMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CHECK = 0;
  private static final int METHODID_REPORT = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final MixerImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(MixerImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @Override
    @SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CHECK:
          serviceImpl.check((Check.CheckRequest) request,
              (io.grpc.stub.StreamObserver<Check.CheckResponse>) responseObserver);
          break;
        case METHODID_REPORT:
          serviceImpl.report((Report.ReportRequest) request,
              (io.grpc.stub.StreamObserver<Report.ReportResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @Override
    @SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class MixerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    MixerBaseDescriptorSupplier() {}

    @Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return Service.getDescriptor();
    }

    @Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Mixer");
    }
  }

  private static final class MixerFileDescriptorSupplier
      extends MixerBaseDescriptorSupplier {
    MixerFileDescriptorSupplier() {}
  }

  private static final class MixerMethodDescriptorSupplier
      extends MixerBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    MixerMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (MixerGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new MixerFileDescriptorSupplier())
              .addMethod(getCheckMethod())
              .addMethod(getReportMethod())
              .build();
        }
      }
    }
    return result;
  }
}
