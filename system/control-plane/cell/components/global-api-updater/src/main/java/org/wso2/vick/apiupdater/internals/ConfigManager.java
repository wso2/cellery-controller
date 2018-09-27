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
package org.wso2.vick.apiupdater.internals;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.vick.apiupdater.beans.controller.Cell;
import org.wso2.vick.apiupdater.beans.controller.RestConfig;
import org.wso2.vick.apiupdater.exceptions.APIException;
import org.wso2.vick.apiupdater.utils.Constants;

import java.io.File;
import java.io.IOException;

/**
 * Methods to read the configuration files
 */
public class ConfigManager {
    private static final Log log = LogFactory.getLog(ConfigManager.class);
    private static volatile Cell cell = null;
    private static volatile RestConfig restConfig = null;

    /**
     * Initializes the Cell configuration
     */
    private static void cellInitialize() throws IOException {
        synchronized (ConfigManager.class) {
            if (cell == null) {
                if (log.isDebugEnabled()) {
                    log.debug("Loading Configuration..");
                }
                cell = loadCellConfig();
            }
        }
    }

    /**
     * Initializes the REST configuration
     */
    private static void restConfigInitialize() throws IOException {
        synchronized (ConfigManager.class) {
            if (restConfig == null) {
                if (log.isDebugEnabled()) {
                    log.debug("Loading Configuration..");
                }
                restConfig = loadRESTConfig();
            }
        }
    }

    /**
     * Returns the whole configuration as a {@link Cell} beans
     *
     * @return cell bean
     */
    public static Cell getCellConfiguration() throws APIException {
        try {
            if (cell == null) {
                cellInitialize();
            }
        } catch (IOException e) {
            String errorMessage = "Error occurred while initializing configuration.";
            log.error(errorMessage, e);
            throw new APIException(errorMessage, e);
        }
        return cell;
    }

    /**
     * Returns the whole configuration as a {@link RestConfig} beans
     *
     * @return REST bean
     */
    public static RestConfig getRestConfiguration() throws APIException {
        try {
            if (restConfig == null) {
                restConfigInitialize();
            }
        } catch (IOException e) {
            String errorMessage = "Error occurred while initializing configuration.";
            log.error(errorMessage, e);
            throw new APIException(errorMessage, e);
        }
        return restConfig;
    }

    /**
     * Load Cell configuration
     *
     * @return Cell Config
     */
    private static Cell loadCellConfig() throws IOException {
        ObjectMapper objectMapper = new ObjectMapper();
        return objectMapper.readValue(new File(Constants.Utils.CELL_CONFIGURATION_FILE_PATH), Cell.class);
    }

    /**
     * Load configurations required for REST APIs
     *
     * @return REST Config
     */
    private static RestConfig loadRESTConfig() throws IOException {
        ObjectMapper objectMapper = new ObjectMapper();
        return objectMapper.readValue(new File(Constants.Utils.REST_CONFIGURATION_FILE_PATH), RestConfig.class);
    }
}
