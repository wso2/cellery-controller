// Copyright (c)  WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
//
// WSO2 Inc. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file   except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

import ballerina/http;
import ballerina/log;
import ballerina/auth;
import ballerina/config;
import ballerina/runtime;
import ballerina/system;
import ballerina/time;
import ballerina/io;
import ballerina/reflect;
import wso2/gateway;


endpoint http:Client spEndpoint {
    url: "http://vick-wso2sp-worker.vick-system.svc.cluster.local:9092"
};

// Extension filter used to send custom error messages and to do customizations.
@Description { value: "Represents the extension filter, which used to customize and extend the request and response
                handling" }
@Field {value:"filterRequest: intercepts the request flow"}
public type ExtensionFilter object {

    @Description {value:"filterRequest: Request filter function"}
    public function filterRequest (http:Listener listener, http:Request request, http:FilterContext context) returns boolean {
    io:println("hit filter request");
	context.attributes["REQUEST_METHOD"] = request.method;
	context.attributes["USER_AGENT"] = request.userAgent;   
        return true;
    }

    public function filterResponse(http:Response response, http:FilterContext context) returns boolean {
    io:println("hit filter response");
	string hostname = system:getEnv("HOSTNAME");
	if(hostname.equalsIgnoreCase("")){
	 	hostname = "my-cell--defaultpod.default-namespace";
	}
	string serverName = "kubernetes://" + hostname;
	string serviceName = context.serviceName;
	string serviceMethod = <string> context.attributes["REQUEST_METHOD"];
	int responseTime = time:currentTime().time;
	int requestTime =  check <int>context.attributes["REQUEST_TIME"];
	float responseDuration = <float>responseTime - requestTime;
	int responseCode =  response.statusCode; 
	string userAgent =  <string> context.attributes["USER_AGENT"];
	string requestIP =  <string> context.attributes["remote_address"];

	http:Request req = new;
	json eventPayload = { timestamp: responseTime, serverName: serverName, serviceName:serviceName, serviceMethod: serviceMethod, 
			responseTime:responseDuration, httpResponseCode: responseCode, userAgent:userAgent, requestIP:requestIP };
	json msgPayload = {event: eventPayload};
	req.setJsonPayload(msgPayload);

    	http:Response|error spResponse = spEndpoint->post("/vick-request", req);
    	match spResponse {
        http:Response resp => {
            var msg = resp.getJsonPayload();
            match msg {
                json jsonPayload => {
                    
                }
                error err => {
                    log:printError(err.message, err = err);
                }
            }
          }
        error err => { log:printError(err.message, err = err); }
    	}

        return true;
    }
};

@Description {value:"This method can be used to send custom error message in an authentication failure"}
function setAuthenticationErrorResponse(http:Response response, http:FilterContext context) {
    //Un comment the following code and set the proper error messages

    //int statusCode = check <int>context.attributes[gateway:HTTP_STATUS_CODE];
    //string errorDescription = <string>context.attributes[gateway:ERROR_DESCRIPTION];
    //string errorMesssage = <string>context.attributes[gateway:ERROR_MESSAGE];
    //int errorCode = check <int>context.attributes[gateway:ERROR_CODE];
    //response.statusCode = statusCode;
    //response.setContentType(gateway:APPLICATION_JSON);
    //json payload = {fault : {
    //    code : errorCode,
    //    message : errorMesssage,
    //    description : errorDescription
    //}};
    //response.setJsonPayload(payload);
}

@Description {value:"This method can be used to send custom error message in an authorization failure"}
function setAuthorizationErrorResponse(http:Response response, http:FilterContext context) {

}

@Description {value:"This method can be used to send custom error message when message throttled out"}
function setThrottleFailureResponse(http:Response response, http:FilterContext context) {

}

@Description {value:"This method can be used to send custom general error message "}
function setGenericErrorResponse(http:Response response, http:FilterContext context) {

}

