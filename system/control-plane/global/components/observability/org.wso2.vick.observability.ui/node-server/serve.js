/*
 * Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

const path = require("path");
const fs = require("fs");
const express = require("express");
const fallback = require("express-history-api-fallback");

const app = express();
const webPortalPort = process.env.PORTAL_PORT || 3000;

const configRoot = path.join(__dirname, "/config");
const portalConfigFile = `${configRoot}/portal.json`;
console.log(`Using Portal Configuration from ${portalConfigFile} file`);

let portalConfig;
const loadPortalConfig = () => {
    portalConfig = fs.readFileSync(`${portalConfigFile}`, "utf8");
    console.log("Loaded new Portal Configuration");
};
loadPortalConfig();

// Watching for config changes
fs.watch(configRoot, null, () => {
    loadPortalConfig();
});

// REST API for configurations
app.get("/config", (req, res) => {
    res.set("Content-Type", "application/json");
    res.send(portalConfig);
});

if (process.env.APP_ENV !== "DEV") {
    const appRoot = path.join(__dirname, "/public");
    console.log(`Using App from ${appRoot} directory`);

    // Serving the React App
    app.use(express.static(appRoot));
    app.use(fallback("index.html", {
        root: appRoot
    }));
} else {
    console.log("Serving Only the Observability Portal Configuration");
}

const server = app.listen(webPortalPort, () => {
    const host = server.address().address;
    const port = server.address().port;

    console.log("Cellery Observability Portal listening at http://%s:%s", host, port);
});
