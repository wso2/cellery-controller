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
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.auth.cell.sts.service;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.vick.auth.cell.sts.model.CellStsRequest;
import org.wso2.vick.auth.cell.sts.model.CellStsResponse;

/**
 * Intercepts inbounds calls to pods within the Cell.
 */
public class VickCellInboundInterceptorService extends VickCellInterceptorService {

    private Logger log = LoggerFactory.getLogger(VickCellInboundInterceptorService.class);


    public VickCellInboundInterceptorService(VickCellStsService cellStsService) throws VickCellSTSException {
        super(cellStsService);
    }

    @Override
    protected void handleRequest(CellStsRequest cellStsRequest,
                                 CellStsResponse cellStsResponse) throws VickCellSTSException {

        log.debug("Intercepted Sidecar Inbound Request:\nSource:{}\nDestination:{}\n", cellStsRequest.getSource(),
                cellStsRequest.getDestination());
        cellStsService.handleInboundRequest(cellStsRequest, cellStsResponse);
    }

}
