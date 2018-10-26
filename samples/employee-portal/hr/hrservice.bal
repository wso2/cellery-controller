import ballerina/http;
import ballerina/log;
import ballerina/internal;
import ballerina/config;
import ballerinax/docker;
import ballerina/io;

@final
string AUTH_HEADER = "Authorization";
@final
string BEARER_PREFIX = "Bearer";
@final
string JWT_SUB_TOKEN = "sub";
@final
string EMPLOYEE_NAME_HEADER = "x-emp-name";

endpoint http:Client employeeDetailsEp {
    url: "http://" + config:getAsString("employeegw.url")
};

endpoint http:Client stockOptionsEp {
    url: "http://" + config:getAsString("stockgw.url")
};

@http:ServiceConfig {
    basePath:"/"
}

@docker:Config {
    registry:"wso2vick",
    name:"sampleapp-hr",
    tag:"v1.0"
}
service<http:Service> hr bind { port: 8080 } {

    @http:ResourceConfig {
        methods:["GET"],
        path:"/"
    }
    getHrDetails (endpoint caller, http:Request req) {
        http:Response res = new;

        string[] headers = req.getHeaderNames();
        foreach header in headers {
            io:println(header + ": " + req.getHeader(untaint header));
        }

        // get JWT header
        if (req.hasHeader(AUTH_HEADER)) {
            match extractJwtTokenFromHeader(req.getHeader(AUTH_HEADER)) {
                string jwtToken => {
                    string[] jwtParts = jwtToken.split("\\.");
                    if (lengthof jwtParts < 3) {
                        // wrong jwt format, incorrect request
                        log:printInfo("JWT token does not have header, payload and signature separated by '.'");
                        res.statusCode = 401;
                    } else {
                        log:printInfo("Request recieved: ");
                        io:println(req);
                        match internal:parseJson(check jwtParts[1].base64Decode()) {
                            json payload => {
                                string sub = check <string> payload[JWT_SUB_TOKEN];
                                string employeeName = sub.split("@")[0];
                                io:println("employee name: " + employeeName);
                                // set the username as a header in the request
                                req.setHeader(EMPLOYEE_NAME_HEADER, employeeName);
                                json employeeDetails;
                                match getEmployeeDetails(req) {
                                    http:Response response => {
                                        match response.getJsonPayload() {
                                            json jsonEmpDetails => {
                                                employeeDetails = jsonEmpDetails;
                                                io:println("employee details: ");
                                                io:println(employeeDetails);
                                            }
                                            error e => {
                                                log:printError("Error in extracting response from Employee service", err = e);
                                                res.statusCode = 500;
                                                caller->respond(res) but { error e1 => log:printError("Error sending response", err = e1) };
                                                done;
                                            }
                                        }
                                    }
                                    error e => {
                                        res.statusCode = 500;
                                        caller->respond(res) but { error e1 => log:printError("Error sending response", err = e1) };
                                        done;
                                    }
                                }
                                json stockOptionDetails;
                                match getStockOptionData(req) {
                                    http:Response response => {
                                        match response.getJsonPayload() {
                                            json jsonStocks => {
                                                stockOptionDetails = jsonStocks;
                                                io:println("stock details: ");
                                                io:println(stockOptionDetails);
                                            }
                                            error e => {
                                                log:printError("Error in extracting response from Stock service", err = e);
                                                res.statusCode = 500;
                                                caller->respond(res) but { error e1 => log:printError("Error sending response", err = e1) };
                                                done;
                                            }
                                        }
                                    }
                                    error e => {
                                        res.statusCode = 500;
                                        caller->respond(res) but { error e1 => log:printError("Error sending response", err = e1) };
                                        done;
                                    }
                                }
                                json resp = buildResponse(employeeDetails, stockOptionDetails);
                                io:println("response: ");
                                io:println(resp);
                                res.setJsonPayload(untaint resp);
                            }
                            error e =>  {
                                log:printError("Error parsing Json", err = e);
                            }
                        }
                    }
                }
                error e => {
                    // jwt token extraction error
                    log:printError("Unabe to extract JWT token", err = e);
                    res.statusCode = 401;
                }
            }
        } else {
            // no authorization header, incorrect request
            log:printInfo("No auth header sent");
            res.statusCode = 401;
        }
        caller->respond(res) but { error e => log:printError("Error sending response", err = e) };
    }
}

function buildResponse(json employeeDetails, json stockOptions) returns json  {
    json response = { employee: { details: employeeDetails, stocks: stockOptions } };
    return response;
}

function getEmployeeDetails(http:Request clientRequest) returns http:Response|error {
    var response = employeeDetailsEp->get("/employee/details", message = untaint clientRequest);
    match response {
        http:Response httpResponse => {
            return httpResponse;
        }
        error e => {
            log:printError("Error in invoking EmployeeDetails service", err = e);
            return e;
        }
    }
}

function getStockOptionData(http:Request clientRequest) returns http:Response|error {
    var response = stockOptionsEp->get("/stock/options", message = untaint clientRequest);
    match response {
        http:Response httpResponse => {
            return httpResponse;
        }
        error e => {
            log:printError("Error in invoking StockOptions service", err = e);
            return e;
        }
    }
}

function extractJwtTokenFromHeader(string authHeader) returns (string)|error {
    // extract auth header
    if (!authHeader.hasPrefix(BEARER_PREFIX)) {
        error err = {message: "Auth failure"};
        return err;
    }
    try {
        return authHeader.substring(6, authHeader.length()).trim();
    } catch (error err) {
        return err;
    }
}

function decodeJwtPayload (string jwtPayload) returns string|error {
    return jwtPayload.base64Decode();
}
