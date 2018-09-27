/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.wso2.vick.apiupdater;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.json.JSONArray;
import org.json.JSONObject;
import org.wso2.vick.apiupdater.beans.controller.API;
import org.wso2.vick.apiupdater.beans.request.ApiCreateRequest;
import org.wso2.vick.apiupdater.beans.controller.ApiDefinition;
import org.wso2.vick.apiupdater.beans.controller.Cell;
import org.wso2.vick.apiupdater.beans.request.Endpoint;
import org.wso2.vick.apiupdater.beans.request.Label;
import org.wso2.vick.apiupdater.beans.request.Method;
import org.wso2.vick.apiupdater.beans.request.Parameter;
import org.wso2.vick.apiupdater.beans.request.PathDefinition;
import org.wso2.vick.apiupdater.beans.request.PathsMapping;
import org.wso2.vick.apiupdater.beans.request.ProductionEndpoint;
import org.wso2.vick.apiupdater.beans.controller.RestConfig;
import org.wso2.vick.apiupdater.exceptions.APIException;
import org.wso2.vick.apiupdater.internals.ConfigManager;
import org.wso2.vick.apiupdater.utils.RequestProcessor;
import org.wso2.vick.apiupdater.utils.Constants;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Class for create APIs in global API Manager.
 */
public class UpdateManager {

    private static String clientId;
    private static String clientSecret;
    private static RestConfig restConfig;
    private static Cell cellConfig;
    private static String apiToken;
    private static final Log log = LogFactory.getLog(UpdateManager.class);

    public static void main(String[] args) {
        try {
            // Encode username password to base64
            restConfig = ConfigManager.getRestConfiguration();
            cellConfig = ConfigManager.getCellConfiguration();
            String username = restConfig.getUsername();
            String password = restConfig.getPassword();
            byte[] message = (username + ":" + password).getBytes(StandardCharsets.UTF_8);
            String userAuth = Base64.getEncoder().encodeToString(message);

            generateClientIDSecret(userAuth);
            generateAccessToken();
            createLabel();
            List apiIds = createAPIs();
            publishAPIs(apiIds);

            log.info("API creation is completed successfully..");
        } catch (APIException e) {
            log.error("Error occurred while creating APIs in Global API manager. " + e.getMessage(), e);
            System.exit(1);
        }
    }

    /**
     * Generate Client ID and Client Secret.
     *
     * @param authHeader Authorization Header
     * @throws APIException Throw an exception if any error occurred.
     */
    private static void generateClientIDSecret(String authHeader) throws APIException {
        if (log.isDebugEnabled()) {
            log.debug("Calling the dynamic client registration endpoint...");
        }
        RequestProcessor requestProcessor = new RequestProcessor();
        String apimBaseURL = restConfig.getApimBaseUrl();
        String applicationResponse = requestProcessor
                .doPost(apimBaseURL + restConfig.getRegisterPath(), Constants.Utils.CONTENT_TYPE_APPLICATION_JSON,
                        Constants.Utils.CONTENT_TYPE_APPLICATION_JSON, Constants.Utils.BASIC + authHeader,
                        restConfig.getRegisterObj().toJSONString());

        if (applicationResponse != null) {
            JSONObject jsonObj = new JSONObject(applicationResponse);
            clientId = jsonObj.getString(Constants.Utils.CLIENT_ID);
            clientSecret = jsonObj.getString(Constants.Utils.CLIENT_SECRET);
        }
    }

    /**
     * Generate access tokens required to invoke RESTful APIs
     *
     * @throws APIException throw API Exception if an error occurred while generating an access token.
     */
    private static void generateAccessToken() throws APIException {
        if (log.isDebugEnabled()) {
            log.debug("Calling token endpoint to generate access tokens...");
        }

        String tokenPayload = Constants.Utils.TOKEN_PAYLOAD.replace("$USER", restConfig.getUsername())
                                                                    .replace("$PASS", restConfig.getPassword());
        apiToken = getToken(tokenPayload);
    }

    /**
     * Invoke Rest API to get token
     *
     * @param tokenPayload Post payload
     * @return access token
     * @throws APIException throw API Exception if an error occurred
     */
    private static String getToken(String tokenPayload) throws APIException {
        RequestProcessor requestProcessor = new RequestProcessor();
        String auth = getBase64EncodedClientIdAndSecret();
        String apiCreateTokenResponse = requestProcessor
                .doPost(restConfig.getTokenEndpoint(), Constants.Utils.CONTENT_TYPE_APPLICATION_URL_ENCODED,
                        Constants.Utils.CONTENT_TYPE_APPLICATION_JSON, Constants.Utils.BASIC + auth, tokenPayload);

        if (apiCreateTokenResponse != null) {
            JSONObject jsonObj = new JSONObject(apiCreateTokenResponse);
            return jsonObj.getString(Constants.Utils.ACCESS_TOKEN);
        } else {
            throw new APIException(
                    "Error while generating the access token from token endpoint: " + restConfig.getTokenEndpoint());
        }
    }

    /**
     * Create microgateway label.
     *
     * @throws APIException throw an API Exception if an error occurred while creating a label
     */
    private static void createLabel() throws APIException {
        RequestProcessor requestProcessor = new RequestProcessor();
        Label label = new Label();
        label.setName(cellConfig.getCell());
        label.setAccessUrls(Collections.singletonList(cellConfig.getCell() + "-gateway"));
        ObjectMapper objectMapper = new ObjectMapper();

        String addLabelTokenResponse;
        try {
            addLabelTokenResponse = requestProcessor
                    .doPost(restConfig.getApimBaseUrl() + restConfig.getAddLabelPath(),
                            Constants.Utils.CONTENT_TYPE_APPLICATION_JSON,
                            Constants.Utils.CONTENT_TYPE_APPLICATION_JSON,
                            Constants.Utils.BEARER + apiToken,
                            objectMapper.writeValueAsString(label));
        } catch (JsonProcessingException e) {
            throw new APIException("Error while serializing the payload: " + label);
        }

        if (addLabelTokenResponse == null) {
            throw new APIException("Error while creating the label. Url: " + restConfig.getApimBaseUrl() +
                                   restConfig.getAddLabelPath());
        }
    }

    /**
     * Create APIs.
     *
     * @return created API Ids.
     * @throws APIException throw API Exception if an error occurred while creating APIs.
     */
    private static List createAPIs() throws APIException {
        if (log.isDebugEnabled()) {
            log.debug("Creating APIs in Global API Manager...");
        }

        JSONArray apiPayloads = createApiPayloads();
        List<String> apiIDs = new ArrayList<>();

        for (int i = 0; i < apiPayloads.length(); i++) {
            RequestProcessor requestProcessor = new RequestProcessor();
            ObjectMapper objectMapper = new ObjectMapper();
            String apiCreateResponse;
            try {
                apiCreateResponse = requestProcessor.doPost(restConfig.getApimBaseUrl() + restConfig.getApiCreatePath(),
                        Constants.Utils.CONTENT_TYPE_APPLICATION_JSON, Constants.Utils.CONTENT_TYPE_APPLICATION_JSON,
                        Constants.Utils.BEARER + apiToken, objectMapper.writeValueAsString(apiPayloads.get(i)));
            } catch (JsonProcessingException e) {
                throw new APIException("Error while serializing the payload: " + apiPayloads.get(i));
            }

            if (apiCreateResponse != null) {
                if (!(apiCreateResponse.contains(Constants.Utils.DUPLICATE_API_ERROR) ||
                      apiCreateResponse.contains(Constants.Utils.DIFFERENT_CONTEXT_ERROR))) {
                    JSONObject jsonObj = new JSONObject(apiCreateResponse);
                    apiIDs.add(jsonObj.getString(Constants.Utils.ID));
                }
            } else {
                throw new APIException("Error while creating the global API from: " + restConfig.getApimBaseUrl() +
                                       restConfig.getApiCreatePath());
            }
        }
        return apiIDs;
    }

    /**
     * Publish APIs in created state.
     *
     * @param ids API Id list
     * @throws APIException Throw API Exception if an error occurred while publishing APIs.
     */
    private static void publishAPIs(List ids) throws APIException {
        if (log.isDebugEnabled()) {
            log.debug("Publishing created APIs in Global API Manager...");
        }

        for (Object id : ids) {
            RequestProcessor requestProcessor = new RequestProcessor();
            String apiPublishResponse;
            apiPublishResponse = requestProcessor
                    .doPost(restConfig.getApimBaseUrl() + restConfig.getApiPublishPath() + "apiId=" + id +
                            "&action=Publish", Constants.Utils.CONTENT_TYPE_APPLICATION_JSON,
                            Constants.Utils.CONTENT_TYPE_APPLICATION_JSON,
                            Constants.Utils.BEARER + apiToken,
                            Constants.Utils.EMPTY_STRING);

            if (apiPublishResponse == null) {
                throw new APIException(
                        "Error while publishing the global API with URL: " + restConfig.getApimBaseUrl() +
                        restConfig.getApiPublishPath() + "apiId=" + id + "&action=Publish");
            }
        }
    }

    private static String getBase64EncodedClientIdAndSecret() {
        byte[] message = (clientId + ":" + clientSecret).getBytes(StandardCharsets.UTF_8);
        return Base64.getEncoder().encodeToString(message);
    }

    private static JSONArray createApiPayloads() throws APIException {
        JSONArray apiPayloadsArray = new JSONArray();
        API[] apis = cellConfig.getApis();

        for (API api : apis) {
            ApiCreateRequest apiCreateRequest = new ApiCreateRequest();
            apiCreateRequest.setName(
                    cellConfig.getCell() + "-" + cellConfig.getVersion() + "-" + api.getContext().replace("/",
                            Constants.Utils.EMPTY_STRING));
            apiCreateRequest.setContext(api.getContext());
            apiCreateRequest.setVersion(cellConfig.getVersion());
            apiCreateRequest.setApiDefinition(getAPIDefinition(api));
            apiCreateRequest.setEndpointConfig(getEndpoint(api));
            Map<String, String> cellLabel = new HashMap<>();
            cellLabel.put("name", cellConfig.getCell());
            apiCreateRequest.setLabels(Collections.singletonList(cellLabel));

            // Add Api to global gateway.
            if (api.isGlobal()) {
                apiCreateRequest.setGatewayEnvironments(Constants.Utils.PRODUCTION_AND_SANDBOX);
            }
            apiPayloadsArray.put(apiCreateRequest);
        }
        return apiPayloadsArray;
    }

    /**
     * Create endpoint_config payload required for API creation payload
     *
     * @param api Api details
     * @return endpoint payload string
     */
    private static String getEndpoint(API api) {
        String response = Constants.Utils.EMPTY_STRING;
        ProductionEndpoint productionEndpoint = new ProductionEndpoint();
        productionEndpoint.setUrl(api.getBackend());

        Endpoint endpoint = new Endpoint();
        endpoint.setProductionEndPoint(productionEndpoint);

        ObjectMapper objectMapper = new ObjectMapper();

        try {
            response = objectMapper.writeValueAsString(endpoint);
        } catch (JsonProcessingException e) {
            log.error("Error occurred while serializing json to string", e);
        }
        return response;
    }

    /**
     * Create api definition payload required for API creation payload
     *
     * @param api Api details
     * @return api definition payload string
     */
    private static String getAPIDefinition(API api) throws APIException {
        PathsMapping apiDefinition = new PathsMapping();
        ApiDefinition[] definitions = api.getDefinitions();

        for (ApiDefinition definition : definitions) {
            PathDefinition pathDefinition = new PathDefinition();
            String methodStr = definition.getMethod();
            Method method = new Method();

            switch (methodStr) {
                case "GET":
                    pathDefinition.setGet(method);
                    break;
                case "POST":
                    Parameter parameter = new Parameter();
                    parameter.setName(Constants.Utils.BODY);
                    parameter.setIn(Constants.Utils.BODY);
                    method.setParameters(Collections.singletonList(parameter));
                    pathDefinition.setPost(method);
                    break;
                default:
                    throw new APIException("Method: " + methodStr + "is not implemented");
            }

            apiDefinition.addPathDefinition(definition.getPath() + Constants.Utils.ALLOW_QUERY_PATTERN, pathDefinition);
        }
        ObjectMapper objectMapper = new ObjectMapper();
        String apiDefinitionStr = Constants.Utils.EMPTY_STRING;

        try {
            apiDefinitionStr = objectMapper.writeValueAsString(apiDefinition);
        } catch (JsonProcessingException e) {
            log.error("Error occurred while serializing json to string", e);
        }
        return apiDefinitionStr;
    }
}
