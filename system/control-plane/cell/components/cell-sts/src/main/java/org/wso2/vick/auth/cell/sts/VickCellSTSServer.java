package org.wso2.vick.auth.cell.sts;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.sts.service.VickCellAuthorizationService;
import org.wso2.vick.auth.cell.sts.service.VickCellSTSException;

import java.io.IOException;

/**
 * Intercepts outbound calls from micro service proxy.
 */
public class VickCellSTSServer {

    private static final Logger log = LoggerFactory.getLogger(VickCellSTSServer.class);
    private final int port;
    private final Server server;

    private VickCellSTSServer(int port) throws VickCellSTSException {

        this.port = port;
        server = ServerBuilder.forPort(port).addService(new VickCellAuthorizationService()).build();
    }

    /**
     * Start serving requests.
     */
    private void start() throws IOException {

        server.start();
        log.info("Vick Cell STS GRPC Server started, listening on " + port);
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
