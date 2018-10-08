package org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 * <pre>
 * A generic interface for performing authorization check on incoming
 * requests to a networked service.
 * </pre>
 */
@javax.annotation.Generated(
        value = "by gRPC proto compiler (version 1.15.0)",
        comments = "Source: envoy/service/auth/v2alpha/external_auth.proto")
public final class AuthorizationGrpc {

    private AuthorizationGrpc() {

    }

    public static final String SERVICE_NAME = "envoy.service.auth.v2alpha.Authorization";

    // Static method descriptors that strictly reflect the proto.
    private static volatile io.grpc.MethodDescriptor<ExternalAuth.CheckRequest,
            ExternalAuth.CheckResponse> getCheckMethod;

    @io.grpc.stub.annotations.RpcMethod(
            fullMethodName = SERVICE_NAME + '/' + "Check",
            requestType = ExternalAuth.CheckRequest.class,
            responseType = ExternalAuth.CheckResponse.class,
            methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
    public static io.grpc.MethodDescriptor<ExternalAuth.CheckRequest,
            ExternalAuth.CheckResponse> getCheckMethod() {

        io.grpc.MethodDescriptor<ExternalAuth.CheckRequest, ExternalAuth.CheckResponse> getCheckMethod;
        if ((getCheckMethod = AuthorizationGrpc.getCheckMethod) == null) {
            synchronized (AuthorizationGrpc.class) {
                if ((getCheckMethod = AuthorizationGrpc.getCheckMethod) == null) {
                    AuthorizationGrpc.getCheckMethod = getCheckMethod =
                            io.grpc.MethodDescriptor.<ExternalAuth.CheckRequest, ExternalAuth.CheckResponse>newBuilder()
                                    .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                                    .setFullMethodName(generateFullMethodName(
                                            "envoy.service.auth.v2alpha.Authorization", "Check"))
                                    .setSampledToLocalTracing(true)
                                    .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                                            ExternalAuth.CheckRequest.getDefaultInstance()))
                                    .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                                            ExternalAuth.CheckResponse.getDefaultInstance()))
                                    .setSchemaDescriptor(new AuthorizationMethodDescriptorSupplier("Check"))
                                    .build();
                }
            }
        }
        return getCheckMethod;
    }

    /**
     * Creates a new async stub that supports all call types for the service
     */
    public static AuthorizationStub newStub(io.grpc.Channel channel) {

        return new AuthorizationStub(channel);
    }

    /**
     * Creates a new blocking-style stub that supports unary and streaming output calls on the service
     */
    public static AuthorizationBlockingStub newBlockingStub(
            io.grpc.Channel channel) {

        return new AuthorizationBlockingStub(channel);
    }

    /**
     * Creates a new ListenableFuture-style stub that supports unary calls on the service
     */
    public static AuthorizationFutureStub newFutureStub(
            io.grpc.Channel channel) {

        return new AuthorizationFutureStub(channel);
    }

    /**
     * <pre>
     * A generic interface for performing authorization check on incoming
     * requests to a networked service.
     * </pre>
     */
    public static abstract class AuthorizationImplBase implements io.grpc.BindableService {

        /**
         * <pre>
         * Performs authorization check based on the attributes associated with the
         * incoming request, and returns status `OK` or not `OK`.
         * </pre>
         */
        public void check(ExternalAuth.CheckRequest request,
                          io.grpc.stub.StreamObserver<ExternalAuth.CheckResponse> responseObserver) {

            asyncUnimplementedUnaryCall(getCheckMethod(), responseObserver);
        }

        @Override
        public final io.grpc.ServerServiceDefinition bindService() {

            return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
                    .addMethod(
                            getCheckMethod(),
                            asyncUnaryCall(
                                    new MethodHandlers<
                                            ExternalAuth.CheckRequest,
                                            ExternalAuth.CheckResponse>(
                                            this, METHODID_CHECK)))
                    .build();
        }
    }

    /**
     * <pre>
     * A generic interface for performing authorization check on incoming
     * requests to a networked service.
     * </pre>
     */
    public static final class AuthorizationStub extends io.grpc.stub.AbstractStub<AuthorizationStub> {

        private AuthorizationStub(io.grpc.Channel channel) {

            super(channel);
        }

        private AuthorizationStub(io.grpc.Channel channel,
                                  io.grpc.CallOptions callOptions) {

            super(channel, callOptions);
        }

        @Override
        protected AuthorizationStub build(io.grpc.Channel channel,
                                          io.grpc.CallOptions callOptions) {

            return new AuthorizationStub(channel, callOptions);
        }

        /**
         * <pre>
         * Performs authorization check based on the attributes associated with the
         * incoming request, and returns status `OK` or not `OK`.
         * </pre>
         */
        public void check(ExternalAuth.CheckRequest request,
                          io.grpc.stub.StreamObserver<ExternalAuth.CheckResponse> responseObserver) {

            asyncUnaryCall(
                    getChannel().newCall(getCheckMethod(), getCallOptions()), request, responseObserver);
        }
    }

    /**
     * <pre>
     * A generic interface for performing authorization check on incoming
     * requests to a networked service.
     * </pre>
     */
    public static final class AuthorizationBlockingStub extends io.grpc.stub.AbstractStub<AuthorizationBlockingStub> {

        private AuthorizationBlockingStub(io.grpc.Channel channel) {

            super(channel);
        }

        private AuthorizationBlockingStub(io.grpc.Channel channel,
                                          io.grpc.CallOptions callOptions) {

            super(channel, callOptions);
        }

        @Override
        protected AuthorizationBlockingStub build(io.grpc.Channel channel,
                                                  io.grpc.CallOptions callOptions) {

            return new AuthorizationBlockingStub(channel, callOptions);
        }

        /**
         * <pre>
         * Performs authorization check based on the attributes associated with the
         * incoming request, and returns status `OK` or not `OK`.
         * </pre>
         */
        public ExternalAuth.CheckResponse check(ExternalAuth.CheckRequest request) {

            return blockingUnaryCall(
                    getChannel(), getCheckMethod(), getCallOptions(), request);
        }
    }

    /**
     * <pre>
     * A generic interface for performing authorization check on incoming
     * requests to a networked service.
     * </pre>
     */
    public static final class AuthorizationFutureStub extends io.grpc.stub.AbstractStub<AuthorizationFutureStub> {

        private AuthorizationFutureStub(io.grpc.Channel channel) {

            super(channel);
        }

        private AuthorizationFutureStub(io.grpc.Channel channel,
                                        io.grpc.CallOptions callOptions) {

            super(channel, callOptions);
        }

        @Override
        protected AuthorizationFutureStub build(io.grpc.Channel channel,
                                                io.grpc.CallOptions callOptions) {

            return new AuthorizationFutureStub(channel, callOptions);
        }

        /**
         * <pre>
         * Performs authorization check based on the attributes associated with the
         * incoming request, and returns status `OK` or not `OK`.
         * </pre>
         */
        public com.google.common.util.concurrent.ListenableFuture<ExternalAuth.CheckResponse> check(
                ExternalAuth.CheckRequest request) {

            return futureUnaryCall(
                    getChannel().newCall(getCheckMethod(), getCallOptions()), request);
        }
    }

    private static final int METHODID_CHECK = 0;

    private static final class MethodHandlers<Req, Resp> implements
            io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
            io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
            io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
            io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {

        private final AuthorizationImplBase serviceImpl;
        private final int methodId;

        MethodHandlers(AuthorizationImplBase serviceImpl, int methodId) {

            this.serviceImpl = serviceImpl;
            this.methodId = methodId;
        }

        @Override
        @SuppressWarnings("unchecked")
        public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {

            switch (methodId) {
                case METHODID_CHECK:
                    serviceImpl.check((ExternalAuth.CheckRequest) request,
                            (io.grpc.stub.StreamObserver<ExternalAuth.CheckResponse>) responseObserver);
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

    private static abstract class AuthorizationBaseDescriptorSupplier
            implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {

        AuthorizationBaseDescriptorSupplier() {

        }

        @Override
        public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {

            return ExternalAuth.getDescriptor();
        }

        @Override
        public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {

            return getFileDescriptor().findServiceByName("Authorization");
        }
    }

    private static final class AuthorizationFileDescriptorSupplier
            extends AuthorizationBaseDescriptorSupplier {

        AuthorizationFileDescriptorSupplier() {

        }
    }

    private static final class AuthorizationMethodDescriptorSupplier
            extends AuthorizationBaseDescriptorSupplier
            implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {

        private final String methodName;

        AuthorizationMethodDescriptorSupplier(String methodName) {

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
            synchronized (AuthorizationGrpc.class) {
                result = serviceDescriptor;
                if (result == null) {
                    serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                            .setSchemaDescriptor(new AuthorizationFileDescriptorSupplier())
                            .addMethod(getCheckMethod())
                            .build();
                }
            }
        }
        return result;
    }
}
