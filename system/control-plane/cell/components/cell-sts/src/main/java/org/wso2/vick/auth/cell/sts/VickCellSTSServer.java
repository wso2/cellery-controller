package org.wso2.vick.auth.cell.sts;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import org.apache.commons.lang.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.sts.context.store.UserContextStore;
import org.wso2.vick.auth.cell.sts.context.store.UserContextStoreImpl;
import org.wso2.vick.auth.cell.sts.service.VickCellInboundInterceptorService;
import org.wso2.vick.auth.cell.sts.service.VickCellOutboundInterceptorService;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;

import java.io.IOException;

/**
 * Intercepts outbound calls from micro service proxy.
 */
public class VickCellSTSServer {

    private static final Logger log = LoggerFactory.getLogger(VickCellSTSServer.class);
    private final int inboundListeningPort;
    private final Server inboundListener;

    private final int outboundListeningPort;
    private final Server outboundListener;

    private VickCellSTSServer(int inboundListeningPort, int outboundListeningPort) throws VickCellSTSException {

        UserContextStore contextStore = new UserContextStoreImpl();

        this.inboundListeningPort = inboundListeningPort;
        inboundListener = ServerBuilder.forPort(inboundListeningPort)
                .addService(new VickCellInboundInterceptorService(contextStore))
                .build();

        this.outboundListeningPort = outboundListeningPort;
        outboundListener = ServerBuilder.forPort(outboundListeningPort)
                .addService(new VickCellOutboundInterceptorService(contextStore))
                .build();
    }

    /**
     * Start serving requests.
     */
    private void start() throws IOException {

        inboundListener.start();
        outboundListener.start();
        log.info("Vick Cell STS GRPC Server started, listening for inbound traffic on " + inboundListeningPort);
        log.info("Vick Cell STS GRPC Server started, listening for outbound traffic on " + outboundListeningPort);
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            // Use stderr here since the logger may has been reset by its JVM shutdown hook.
            System.err.println("Shutting down Vick Cell STS since JVM is shutting down.");
            VickCellSTSServer.this.stop();
            System.err.println("Vick Cell STS shut down.");
        }));
    }

    /**
     * Stop serving requests and shutdown resources.
     */
    private void stop() {

        if (inboundListener != null) {
            inboundListener.shutdown();
        }

        if (outboundListener != null) {
            outboundListener.shutdown();
        }
    }

    /**
     * Await termination on the main thread since the grpc library uses daemon threads.
     */
    private void blockUntilShutdown() throws InterruptedException {

        if (inboundListener != null) {
            inboundListener.awaitTermination();
        }

        if (outboundListener != null) {
            outboundListener.awaitTermination();
        }
    }

    public static void main(String[] args) {

        VickCellSTSServer server;
        try {
            server = new VickCellSTSServer(8080, 8081);
            server.start();
            server.blockUntilShutdown();
        } catch (Exception e) {
            log.error("Error while starting up the Cell STS.", e);
            // To make the pod go to CrashLoopBackOff state if we encounter any error while starting up
            System.exit(1);
        }
    }

}
