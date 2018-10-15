package org.wso2.vick.auth.cell.sts.service;

import com.google.rpc.Code;
import com.google.rpc.Status;
import com.mashape.unirest.http.HttpResponse;
import com.mashape.unirest.http.JsonNode;
import com.mashape.unirest.http.Unirest;
import com.mashape.unirest.http.exceptions.UnirestException;
import io.grpc.stub.StreamObserver;
import org.apache.commons.lang.StringUtils;
import org.apache.http.conn.ssl.NoopHostnameVerifier;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.ssl.SSLContextBuilder;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.wso2.vick.auth.cell.sts.generated.envoy.core.Base;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.AuthorizationGrpc;
import org.wso2.vick.auth.cell.sts.generated.envoy.service.auth.v2alpha.ExternalAuth;
import org.wso2.vick.sts.core.VickSTSConstants;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;

/**
 * Intercept outbound HTTP calls from the cell and injects STS token required for authorization. Interception is engaged
 * using an EnvoyFilter
 * <p>
 * refer {@}
 */
public class VickCellOutboundAuthorizationService extends AuthorizationGrpc.AuthorizationImplBase {

    private static final Logger log = LoggerFactory.getLogger(VickCellOutboundAuthorizationService.class);
    private static final String AUTHORIZATION_HEADER_NAME = "authorization";
    private static final String STS_RESPONSE_TOKEN_PARAM = "token";
    private static final String CELL_NAME_ENV_VARIABLE = "CELL_NAME";
    private static final String STS_CONFIG_PATH_ENV_VARIABLE = "CONF_PATH";
    private static final String CONFIG_FILE_PATH = "/etc/config/sts.json";

    private static final String CONFIG_STS_ENDPOINT = "endpoint";
    private static final String CONFIG_AUTH_USERNAME = "username";
    private static final String CONFIG_AUTH_PASSWORD = "password";
    private static final String BEARER_HEADER_VALUE_PREFIX = "Bearer ";

    private static final String REQUEST_ID = "request.id";
    private static final String REQUEST_ID_HEADER = "x-request-id";
    private static final String DESTINATION_HEADER = ":authority";

    private String stsEndpointUrl;
    private String userName;
    private String password;
    private String cellName;

    public VickCellOutboundAuthorizationService() throws VickCellSTSException {

        setUpConfigurationParams();
        setHttpClientProperties();
    }

    private void setUpConfigurationParams() throws VickCellSTSException {

        try {
            String configFilePath = getConfigFilePath();
            String content = new String(Files.readAllBytes(Paths.get(configFilePath)));
            JSONObject config = (JSONObject) new JSONParser().parse(content);
            stsEndpointUrl = (String) config.get(CONFIG_STS_ENDPOINT);
            userName = (String) config.get(CONFIG_AUTH_USERNAME);
            password = (String) config.get(CONFIG_AUTH_PASSWORD);
            cellName = getCellName();

            log.info("Global STS Endpoint: " + stsEndpointUrl);
            log.info("Cell Name: " + cellName);
        } catch (ParseException | IOException e) {
            throw new VickCellSTSException("Error while setting up STS configurations", e);
        }
    }

    private String getConfigFilePath() {

        String configPath = System.getenv(STS_CONFIG_PATH_ENV_VARIABLE);
        return StringUtils.isNotBlank(configPath) ? configPath : CONFIG_FILE_PATH;
    }

    @Override
    public void check(ExternalAuth.CheckRequest request, StreamObserver<ExternalAuth.CheckResponse> responseObserver) {

        try {
            // Add request ID for log correlation.
            MDC.put(REQUEST_ID, getRequestId(request));

            String destination = getDestination(request);
            log.debug("Intercepting Sidecar call to '{}'", destination);
            log.debug("Request from Istio-Proxy:\n{}", request);

            String authzHeaderInRequest = getAuthorizationHeaderValue(request);
            ExternalAuth.CheckResponse response;

            if (StringUtils.isEmpty(authzHeaderInRequest)) {
                log.info("Authorization Header is missing in the outbound call. Injecting a JWT from STS.");

                String stsToken = getTokenToInject(request);
                if (StringUtils.isEmpty(stsToken)) {
                    log.error("No JWT token received from the STS endpoint: " + stsEndpointUrl);
                }
                response = ExternalAuth.CheckResponse.newBuilder()
                        .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
                        .setOkResponse(buildOkHttpResponse(stsToken))
                        .build();
            } else {
                log.info("Authorization Header is present in the request. Continuing without injecting a new JWT.");
                response = ExternalAuth.CheckResponse.newBuilder()
                        .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
                        .build();
            }

            log.debug("Response to istio-proxy:\n{}", response);

            responseObserver.onNext(response);
            responseObserver.onCompleted();
        } catch (VickCellSTSException e) {
            log.error("Error while handling request from istio-proxy to '{}'", getDestination(request), e);
        } finally {
            MDC.clear();
        }
    }

    private ExternalAuth.OkHttpResponse buildOkHttpResponse(String stsToken) {

        ExternalAuth.OkHttpResponse.Builder builder = ExternalAuth.OkHttpResponse.newBuilder();
        if (StringUtils.isNotEmpty(stsToken)) {
            builder.addHeaders(buildHeader(AUTHORIZATION_HEADER_NAME, BEARER_HEADER_VALUE_PREFIX + stsToken));
        }
        return builder.build();
    }

    private String getTokenToInject(ExternalAuth.CheckRequest request) {

        try {
            HttpResponse<JsonNode> apiResponse =
                    Unirest.post(stsEndpointUrl)
                            .basicAuth(userName, password)
                            .field(VickSTSConstants.VickSTSRequest.SUBJECT, cellName)
                            .asJson();

            log.debug("Response from the STS:\nstatus:{}\nbody:{}",
                    apiResponse.getStatus(), apiResponse.getBody().toString());

            if (apiResponse.getStatus() == 200) {
                Object stsTokenValue = apiResponse.getBody().getObject().get(STS_RESPONSE_TOKEN_PARAM);
                return stsTokenValue != null ? stsTokenValue.toString() : null;
            } else {
                log.error("Error from STS endpoint. statusCode= " + apiResponse.getStatus() + ", " +
                        "statusMessage=" + apiResponse.getStatusText());
            }
        } catch (UnirestException e) {
            log.error("Error while obtaining the STS token.", e);
        }

        return null;
    }

    private String getCellName() throws VickCellSTSException {

        // For now we pick the cell name from the environment variable. In future we need to figure out a way to derive
        // values from the authz request.
        String cellName = System.getenv(CELL_NAME_ENV_VARIABLE);
        if (StringUtils.isBlank(cellName)) {
            throw new VickCellSTSException("Environment variable '" + CELL_NAME_ENV_VARIABLE + "' is empty.");
        }
        return cellName;
    }

    private Base.HeaderValueOption buildHeader(String headerName, String headerValue) {

        return Base.HeaderValueOption.newBuilder()
                .setHeader(Base.HeaderValue.newBuilder().setKey(headerName).setValue(headerValue))
                .build();
    }

    private String getAuthorizationHeaderValue(ExternalAuth.CheckRequest request) {

        return getHeader(request, AUTHORIZATION_HEADER_NAME);
    }

    private String getHeader(ExternalAuth.CheckRequest request, String headerKey) {

        return request.getAttributes().getRequest().getHttp().getHeaders().get(headerKey);
    }

    private static void setHttpClientProperties() throws VickCellSTSException {

        try {
            Unirest.setHttpClient(HttpClients.custom()
                    .setSSLContext(new SSLContextBuilder().loadTrustMaterial(null, (x509Certificates, s) -> true).build())
                    .setSSLHostnameVerifier(NoopHostnameVerifier.INSTANCE)
                    .disableRedirectHandling()
                    .build());
        } catch (NoSuchAlgorithmException | KeyManagementException | KeyStoreException e) {
            throw new VickCellSTSException("Error initializing the http client.", e);
        }
    }

    private String getRequestId(ExternalAuth.CheckRequest request) throws VickCellSTSException {

        String id = request.getAttributes().getRequest().getHttp().getHeaders().get(REQUEST_ID_HEADER);
        if (StringUtils.isBlank(id)) {
            throw new VickCellSTSException("Request Id cannot be found in the header: " + REQUEST_ID_HEADER);
        }
        return id;
    }

    private String getDestination(ExternalAuth.CheckRequest request) {

        String destination = request.getAttributes().getRequest().getHttp().getHeaders().get(DESTINATION_HEADER);
        if (StringUtils.isBlank(destination)) {
            destination = getHost(request);
            log.debug("Destination is picked from host value in the request.");
        }
        return destination;
    }

    private String getHost(ExternalAuth.CheckRequest request) {

        return request.getAttributes().getRequest().getHttp().getHost();
    }

}
