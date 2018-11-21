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

const configuration = require("./config/portal.js");
const path = require("path");
const express = require("express");
const fallback = require("express-history-api-fallback");

const app = express();
const webPortalPort = process.env.PORTAL_PORT || 3000;
const root = path.join(__dirname, "/public");

// REST API for configurations
app.get("/config", function (req, res) {
    res.send(JSON.stringify(configuration));
});

// Serving the React App
app.use(express.static(root));
app.use(fallback("index.html", {
    root: root
}));

const server = app.listen(webPortalPort, () => {
    const host = server.address().address;
    const port = server.address().port;

    console.log("WSO2 VICK Observability Portal listening at http://%s:%s", host, port);
});
