package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang/glog"
	"io/ioutil"
	vickingressv1 "github.com/wso2/product-vick/system/ingress-controller/pkg/apis/vick/ingress/v1alpha1"
	"net/http"
	"strconv"
	"strings"
	"github.com/wso2/product-vick/system/ingress-controller/config"
)

const (
	clientRegistrationPath  = "/client-registration/"
	publisherPath           = "/api/am/publisher/"
	clientRegisterPath      = "/register"
	apisPath                = "/apis"
	lifecylePath            = "/apis/change-lifecycle"
	tokenScopes             = "apim:api_view apim:api_create apim:api_publish apim:label_manage"
	httpProtocol            = "http"
	httpsProtocol           = "https"
	defaultVisibility       = "PUBLIC"
	defaultTier             = "Unlimited"
	defaultAuthType         = "Application & Application User"
	defaultSwaggerVersion   = "2.0"
	defaultGatewayEnvs      = "Production and Sandbox"
	defaultIsDefaultVersion = "true"
)

type registrationResponse struct {
	ClientId string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope string `json:"scope"`
	Error string `json:"error"`
	ErrorDesc string `json:"error_description"`
}

type createApiSuccessResponse struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Context string `json:"context"`
	Version string `json:"version"`
	Status string `json:"status"`
}

type createApiErrorResponse struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Description string `json:"description"`
}

type publisheApiErrorResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Description string `json:"description"`
}

type updateApiErrorResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Description string `json:"description"`
}

type searchApiResult struct {
	List []version `json:"list"`
}

type version struct {
	Version string `json:"version"`
	Id string `json:"id"`
}

type IngressRuleCreator struct {
	IngressConfig *config.IngressConfig
	HttpClient *http.Client
}

// TODO bring an abstraction for the HTTP client so that unit test writing is simplified
//type HttpCaller interface {
//	DoHttpCall () (*http.Response, error)
//}
//
//type DefaultHttpCaller struct {
//	client *http.Client
//	headers *map[string]string
//	req *http.Request
//}
//
//func (caller *DefaultHttpCaller) DoHttpCall () (*http.Response, error)  {
//	for name, val := range *caller.headers {
//		caller.req.Header.Set(name, val)
//	}
//	return caller.client.Do(caller.req)
//}
//
//func NewDefaultHttpCallerWithHttpClient () *DefaultHttpCaller {
//	transport := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//	client := &http.Client{Transport: transport}
//	return &DefaultHttpCaller{client:client}
//}

func NewIngressRuleCreator (ingressConfig *config.IngressConfig) *IngressRuleCreator {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	return &IngressRuleCreator{IngressConfig:ingressConfig, HttpClient:client}
}

//func RegisterClientHelper (ingRuleCreator *IngressRuleCreator) (*map[string]string, *http.Request, error) {
//	// create headers
//	authHeader := getBasicAuthHeader(ingRuleCreator.IngressConfig.Username, ingRuleCreator.IngressConfig.Password)
//	var headers map[string]string
//	headers["Authorization"] = authHeader
//	headers["Content-Type"] = "application/json"
//	headers["Accept"] = "application/json"
//	// create request
//	requestBody, err := json.Marshal(ingRuleCreator.IngressConfig.RegisterPayload)
//	if err != nil {
//		return nil, nil, err
//	}
//	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + clientRegistrationPath +
//		ingRuleCreator.IngressConfig.ApiVersion + clientRegisterPath, bytes.NewBuffer(requestBody))
//	if err != nil {
//		return nil, nil, err
//	}
//	return &headers, request, nil
//}

func (ingRuleCreator *IngressRuleCreator) RegisterClient() (string, string, error)  {
	// create basic auth header
	authHeader := getBasicAuthHeader(ingRuleCreator.IngressConfig.Username, ingRuleCreator.IngressConfig.Password)
	requestBody, err := json.Marshal(ingRuleCreator.IngressConfig.RegisterPayload)
	if err != nil {
		return "", "", err
	}
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + clientRegistrationPath +
		ingRuleCreator.IngressConfig.ApiVersion + clientRegisterPath, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", err
	}
	request.Header.Set("Authorization", "Basic " + authHeader)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// perform the http request
	response, err := ingRuleCreator.HttpClient.Do(request)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()
	// read response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}
	var regResponse registrationResponse
	err = json.Unmarshal(body, &regResponse)
	if err != nil {
		return "", "", err
	}
	if len(regResponse.ClientId) == 0 || len(regResponse.ClientSecret) == 0 {
		return "", "", errors.New("Client registration response does not contain id/secret")
	}
	return regResponse.ClientId, regResponse.ClientSecret, nil
}

func (ingRuleCreator *IngressRuleCreator) GenerateAccessToken(clientId string, clientSecret string) (string, error) {
	// get basic auth header for registered client, with clientId and clientSecret
	authHeader := getBasicAuthHeader(clientId, clientSecret)
	payload := "grant_type=password&username=" + ingRuleCreator.IngressConfig.Username + "&password=" +
		ingRuleCreator.IngressConfig.Password + "&scope=" + tokenScopes
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.TokenEp,
		bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Basic " + authHeader)
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	// perform the http request
	response, err := ingRuleCreator.HttpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	// read response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var tokResponse tokenResponse
	err = json.Unmarshal(body, &tokResponse)
	if err != nil {
		return "", err
	}
	// check if there is an error
	if len(tokResponse.Error) != 0 {
		return "", errors.New("Error in retrieving token: " + tokResponse.Error + ", "+ tokResponse.ErrorDesc)
	}
	if len(tokResponse.AccessToken) == 0  {
		return "", errors.New("Token response does not contain access token")
	}
	// the token scopes should match the expected scope set
	glog.Infof("Scopes of the access token: %s", tokenScopes)
	return tokResponse.AccessToken, nil
}

func getBasicAuthHeader (username string, password string) (string) {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func (ingRuleCreator *IngressRuleCreator) CreateApi (token string, ingSpec *vickingressv1.IngressSpec) error {
	err := validateAndSetDefaults(ingSpec)
	if err != nil {
		return err
	}
	// create
	apiId, err := createApi(ingressRuleCreator, token, ingSpec)
	if err != nil {
		return err
	}
	// publish
	err = publishApi(ingressRuleCreator, apiId, token)
	if err != nil {
		return err
	}
	return nil
}

func createApi (ingRuleCreator *IngressRuleCreator, token string, ingSpec *vickingressv1.IngressSpec) (string, error) {
	createApiPayload, err := createApiCreationPayload(ingSpec)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath, bytes.NewBuffer([]byte(createApiPayload)))
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// perform the http request
	response, err := ingRuleCreator.HttpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	// read response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	glog.Infof("Api creation response body: %+v\n", body)

	if response.StatusCode == 201 {
		// api created successfully
		glog.Infof("Api with name: %s, context: %s, version: %s created successfully", ingSpec.Name,
			ingSpec.Context, ingSpec.Version)
		var creatApiRes createApiSuccessResponse
		err = json.Unmarshal(body, &creatApiRes)
		if err != nil {
			return "", err
		}
		return creatApiRes.Id, nil

	} else {
		glog.Infof("Api creation failed for api: %s, context: %s, version: %s", ingSpec.Name, ingSpec.Context,
			ingSpec.Version)
		// marshal the error response to get the error message
		var createApiErrorRes createApiErrorResponse
		err = json.Unmarshal(body, &createApiErrorRes)
		if err != nil {
			return "", err
		}
		glog.Info(createApiErrorRes)
		return "", errors.New("Api creation failed for api: " + ingSpec.Name + ", context: " + ingSpec.Context +
			", version: " + ingSpec.Version + ", error code: " + createApiErrorRes.Code + ", error msg: " +
			createApiErrorRes.Message + ", error desc: " + createApiErrorRes.Description)
	}
}

func (ingRuleCreator *IngressRuleCreator) UpdateApi (token string, ingSpec *vickingressv1.IngressSpec) error {
	err, apiId := getApiId(ingressRuleCreator, ingSpec.Context, ingSpec.Version, token)
	if err != nil {
		return err
	}
	if isEmpty(apiId) {
		glog.Infof("API with context %s and version %s" , ingSpec.Context, ingSpec.Version +
			" not found, unable to update")
		return nil
	}

	err = validateAndSetDefaults(ingSpec)
	if err != nil {
		return err
	}
	// update
	apiId, err = updateApi(ingressRuleCreator, token, ingSpec, apiId)
	if err != nil {
		return err
	}
	// publish
	err = publishApi(ingressRuleCreator, apiId, token)
	if err != nil {
		return err
	}
	return nil
}

func updateApi (ingRuleCreator *IngressRuleCreator, token string, ingSpec *vickingressv1.IngressSpec, apiId string) (string, error) {
	createApiPayload, err := createApiCreationPayload(ingSpec)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("PUT", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "/" + apiId, bytes.NewBuffer([]byte(createApiPayload)))
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// perform the http request
	response, err := ingRuleCreator.HttpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	// read response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode == 200 {
		// api created successfully
		glog.Infof("Api with name: %s, context: %s, version: %s updated successfully", ingSpec.Name,
			ingSpec.Context, ingSpec.Version)
		var updateApiSuccessRes createApiSuccessResponse
		err = json.Unmarshal(body, &updateApiSuccessRes)
		if err != nil {
			return "", err
		}
		return updateApiSuccessRes.Id, nil
	} else {
		glog.Infof("Api update failed for api: %s, context: %s, version: %s", ingSpec.Name,
			ingSpec.Context, ingSpec.Version)
		// marshal the error response to get the error message
		var publishApiErrorRes updateApiErrorResponse
		err = json.Unmarshal(body, &publishApiErrorRes)
		if err != nil {
			return "", err
		}
		glog.Info(publishApiErrorRes)
		return "", errors.New("Api update failed for api: " + ingSpec.Name + ", context: " + ingSpec.Context +
			", version: " + ingSpec.Version + ", error code: " + strconv.Itoa(publishApiErrorRes.Code) + ", error msg: " +
			publishApiErrorRes.Message + ", error desc: " + publishApiErrorRes.Description)

	}
}

func getApiId (ingRuleCreator *IngressRuleCreator, context string, version string, token string) (error,string) {

	searchApiRes, err := searchApis(ingRuleCreator, context, token)
	if err != nil {
		return err, ""
	}

	for _, result := range searchApiRes.List {
		if version == result.Version {
			// api with this particular context and version exists
			glog.Infof("API found with id %s, context %s, and version %s , " + "starting to delete",
				result.Id, context, version)
			return nil, result.Id
		}
	}
	glog.Infof("API with context %s and version %s" , context, version + " not found, unable to delete")
	return nil, ""
}

func (ingRuleCreator *IngressRuleCreator) DeleteApi (token string, ingSpec *vickingressv1.IngressSpec) error {
	err, apiId := getApiId(ingressRuleCreator, ingSpec.Context, ingSpec.Version, token)
	if err != nil {
		return err
	}
	if isEmpty(apiId) {
		glog.Infof("API with context %s and version %s" , ingSpec.Context, ingSpec.Version +
			" not found, unable to delete")
		return nil
	}
	return deleteApi(ingressRuleCreator, apiId, token)
}

func deleteApi (ingRuleCreator *IngressRuleCreator, apiId string, token string) error {
	request, err := http.NewRequest("DELETE", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "/" + apiId, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// perform the http request
	response, err := ingRuleCreator.HttpClient.Do(request)
	if err != nil {
		return err
	}
	response.Body.Close()
	glog.Infof("Api with id %s deleted", apiId)
	return nil
}

func searchApis (ingRuleCreator *IngressRuleCreator, context string, token string) (searchApiResult, error) {

	var searchApiRes searchApiResult
	request, err := http.NewRequest("GET", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + apisPath + "?query=context:" + context, nil)
	if err != nil {
		return searchApiRes, err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// perform the http request
	response, err := ingRuleCreator.HttpClient.Do(request)
	if err != nil {
		return searchApiRes, err
	}
	defer response.Body.Close()
	// read response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return searchApiRes, err
	}
	err = json.Unmarshal(body, &searchApiRes)
	if err != nil {
		return searchApiRes, err
	}

	return searchApiRes, nil
}

func publishApi (ingRuleCreator *IngressRuleCreator, apiId string, token string) error {


	glog.Infof("Going to publish the API: %s", apiId)
	request, err := http.NewRequest("POST", ingRuleCreator.IngressConfig.BaseUrl + publisherPath +
		ingRuleCreator.IngressConfig.ApiVersion + lifecylePath + "?apiId=" + apiId + "&action=Publish", nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer " + token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	// perform the http request
	response, err := ingRuleCreator.HttpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		glog.Infof("Api with id %s published successfully", apiId)
	} else {
		// read response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		// marshal the error response to get the error message
		var publishApiErrorRes publisheApiErrorResponse
		err = json.Unmarshal(body, &publishApiErrorRes)
		if err != nil {
			return err
		}
		return errors.New("Api update failed for api " + apiId + ", error code: " + strconv.Itoa(publishApiErrorRes.Code) +
			", error msg: " + publishApiErrorRes.Message + ", error desc: " + publishApiErrorRes.Description)
	}
	return nil
}

func validateAndSetDefaults (ingSpec *vickingressv1.IngressSpec) error {
	// validate mandatory params
	if isEmpty(ingSpec.Name) {
		return handleParamEmpty("api name")
	}
	if isEmpty(ingSpec.Context) {
		return handleParamEmpty("api context")
	}
	if isEmpty(ingSpec.Version) {
		handleParamEmpty("api version")
	}
	if ingSpec.Paths == nil || len(ingSpec.Paths) == 0 {
		return handleParamNotSpecified("api path(s)");
	}
	// either endpointConfig or explicit endpoints should be available
	if isEmpty(ingSpec.EndpointConfig) && ingSpec.Endpoints == nil {
		return handleParamNotSpecified("endpointConfig or endpoints")
	}
	// set defaults for non-mandatory params
	if isEmpty(ingSpec.IsDefaultVersion) || isValidBoolean(ingSpec.IsDefaultVersion) {
		ingSpec.IsDefaultVersion = defaultIsDefaultVersion
	}
	if ingSpec.Transports == nil || len(ingSpec.Transports) == 0 {
		ingSpec.Transports = []string {
			httpProtocol, httpsProtocol,
		}
	}
	if ingSpec.Tiers == nil || len(ingSpec.Tiers) == 0 {
		ingSpec.Tiers = []string {
			defaultTier,
		}
	}
	if isEmpty(ingSpec.Visibility) {
		ingSpec.Visibility = defaultVisibility
	}
	if isEmpty(ingSpec.AuthType) {
		ingSpec.AuthType = defaultAuthType
	}
	if isEmpty(ingSpec.SwaggerVersion) {
		ingSpec.SwaggerVersion = defaultSwaggerVersion
	}
	if isEmpty(ingSpec.GatewayEnvironments) {
		ingSpec.GatewayEnvironments = defaultGatewayEnvs
	}
	// validate endpoints
	for _, endpoint := range ingSpec.Endpoints {
		if isEmpty(endpoint.Type) {
			return handleParamEmpty("endpoint.type")
		}
		if isEmpty(endpoint.ServiceName) {
			return handleParamEmpty("endpoint.serviceName")
		}
		if endpoint.Port <= 0 {
			return handleIntParamEmpty("endpoint.port")
		}
		if isEmpty(endpoint.Protocol) {
			return handleParamEmpty("endpoint.protocol")
		}
		// protocol should be either http or https
		if endpoint.Protocol != "http" && endpoint.Protocol != "https" {
			return handleConfigError("Unsupported protocol '" + endpoint.Protocol + "' for endpoint")
		}
	}

	// validate paths
	for _, path := range ingSpec.Paths {
		handleParamEmpty("path.operation")
		handleParamEmpty("path.context")
		// add '/' at the begining of the context if not exists
		prefixSlash(&path)
		// change each context to include '/*' at the end
		addWildCardResourceMatcher(&path)
		// if no tiers are set, use the api level tiers
		if path.Tiers == nil || len(path.Tiers) == 0 {
			path.Tiers = ingSpec.Tiers
		}
		// if no auth Type is set, use api level auth type
		if isEmpty(path.AuthType) {
			path.AuthType = ingSpec.AuthType
		}
	}
	return nil
}

func isValidBoolean (strVal string) bool {
	_, err := strconv.ParseBool(strVal)
	if err != nil {
		return false
	}
	return true
}

func prefixSlash (path *vickingressv1.Path) {
	if !strings.HasPrefix(path.Context, "/") {
		path.Context = "/" + path.Context
	}
}

func addWildCardResourceMatcher (path *vickingressv1.Path) {
	if strings.HasSuffix(path.Context, "/") {
		path.Context = path.Context + "*"
	} else {
		path.Context = path.Context + "/*"
	}
}

func createApiCreationPayload (ingSpec *vickingressv1.IngressSpec) (string, error) {
	defaultVersionBool, _ := strconv.ParseBool(ingSpec.IsDefaultVersion)
	apiPayload := map[string]interface{}{
		"name": ingSpec.Name,
		"context": ingSpec.Context,
		"version": ingSpec.Version,
		"isDefaultVersion": defaultVersionBool,
		"transport": ingSpec.Transports,
		"tiers": ingSpec.Tiers,
		"gatewayEnvironments": ingSpec.GatewayEnvironments,
		"visibility": ingSpec.Visibility,
	}
	if ingSpec.Labels != nil && len(ingSpec.Labels) > 0 {
		apiPayload["labels"] = ingSpec.Labels
	}
	var err error
	apiPayload["apiDefinition"], err = createApiDefinitionPayload(ingSpec)
	if err != nil {
		return "", err
	}
	apiPayload["endpointConfig"], err = createEndpointConfigPayload(ingSpec)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(&apiPayload)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func createApiDefinitionPayload(ingSpec *vickingressv1.IngressSpec) (string, error) {
	// TODO: if the api definition is given inline in the ingress rule, use that
	apiDefinitionPayload := make(map[string]interface{})
	apiDefinitionPayload["swagger"] = ingSpec.SwaggerVersion
	apiDefPayloadPaths := make(map[string]interface{})
	apiDefinitionPayload["paths"] = apiDefPayloadPaths
	for _, path := range ingSpec.Paths {
		switch path.Operation {
		case "POST", "post", "PUT", "put" :
				apiDefPayloadPaths[path.Context] = map[string]interface{}{
					path.Operation: map[string]interface{}{
						"x-auth-type":       path.AuthType,
						"x-throttling-Tier": path.Tiers,
						"parameters": []interface{}{
							map[string]string{
								"name":     "body",
								"required": "true",
								"in":       "body",
							},
						},
					},
				}
		default:
			apiDefPayloadPaths[path.Context] = map[string]interface{}{
					path.Operation: map[string]interface{}{
						"x-auth-type":       path.AuthType,
						"x-throttling-defaultTier": path.Tiers,
					},
				}
		}
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&apiDefinitionPayload)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func createEndpointConfigPayload(ingSpec *vickingressv1.IngressSpec) (string, error) {
	endpointConfigPayload := make(map[string]interface{})
	endpointConfigPayload["endpoint_type"] = "http"
	for _, endpoint := range ingSpec.Endpoints {
		if strings.Contains(endpoint.Type, "production") {
			// production endpoint
			endpointConfigPayload["production_endpoints"] = map[string]interface{} {
				"url": endpoint.Protocol + "://" + endpoint.ServiceName + ":" + strconv.Itoa(endpoint.Port),
			}
		} else if strings.Contains(endpoint.Type, "sandbox") {
			// sandbox endpoint
			endpointConfigPayload["sandbox_endpoints"] = map[string]interface{} {
				"url": endpoint.Protocol + "://" + endpoint.ServiceName + ":" + strconv.Itoa(endpoint.Port),
			}
		} else {
			// other endpoint
			// sandbox endpoint
			endpointConfigPayload[endpoint.Type] = map[string]interface{} {
				"url": endpoint.Protocol + "://" + endpoint.ServiceName + ":" + string(endpoint.Port),
			}
		}
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&endpointConfigPayload)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func handleIntParamEmpty (parameter string) error {
	return handleParamNotSpecified(parameter);
}

func handleParamEmpty(parameter string) error  {
	return handleParamNotSpecified(parameter)
}

func isEmpty(obj string) bool {
	return len(obj) == 0
}

func handleParamNotSpecified (param string) error {
	return errors.New("Mandatory parameter '" + param + "' not provided")
}

func handleConfigError (msg string) error {
	return errors.New(msg)
}

