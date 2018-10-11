package org.wso2.vick.auth.cell.sts;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.vick.auth.cell.sts.service.VickCellOutboundAuthorizationService;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;

import java.io.IOException;

/**
 * Intercepts outbound calls from micro service proxy.
 */
public class VickCellSTSServer {

    private static final Log log = LogFactory.getLog(VickCellSTSServer.class);
    private final int port;
    private final Server server;

    private VickCellSTSServer(int port) throws VickCellSTSException {

        this.port = port;
        server = ServerBuilder.forPort(port).addService(new VickCellOutboundAuthorizationService()).build();
    }

    /**
     * Start serving requests.
     */
    private void start() throws IOException {

        server.start();
        log.info("Vick Cell STS GRPC Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            // Use stderr here since the logger may has been reset by its JVM shutdown hook.
            System.err.println("*** Shutting down gRPC server since JVM is shutting down");
            VickCellSTSServer.this.stop();
            System.err.println("*** Server shut down");
        }));
    }

    /**
     * Stop serving requests and shutdown resources.
     */
    private void stop() {

        if (server != null) {
            server.shutdown();
        }
    }

    /**
     * Await termination on the main thread since the grpc library uses daemon threads.
     */
    private void blockUntilShutdown() throws InterruptedException {

        if (server != null) {
            server.awaitTermination();
        }
    }

    public static void main(String[] args) {

        VickCellSTSServer server = null;
        try {
            server = new VickCellSTSServer(8080);
            server.start();
            server.blockUntilShutdown();
        } catch (Exception e) {
            log.error("Error while starting up the Cell STS.", e);
            // To make the pod go to CrashLoopBackOff state if we encounter any error while starting up
            System.exit(1);
        }
    }

}
