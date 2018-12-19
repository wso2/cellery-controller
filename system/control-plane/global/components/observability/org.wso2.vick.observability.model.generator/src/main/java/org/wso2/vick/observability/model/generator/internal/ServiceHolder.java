/*
 *  Copyright (c) 2018, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package org.wso2.vick.observability.model.generator.internal;

import org.wso2.carbon.datasource.core.api.DataSourceService;
import org.wso2.vick.observability.model.generator.ModelManager;

/**
 * This class holds the registered services by OSGi, that is required by the entire component.
 */
public class ServiceHolder {
    private static DataSourceService dataSourceService;
    private static ModelStoreManager modelStoreManager;
    private static ModelManager modelManager;
    private static ModelPeriodicProcessor periodicProcessor;

    private ServiceHolder() {
    }

    public static DataSourceService getDataSourceService() {
        return dataSourceService;
    }

    public static void setDataSourceService(DataSourceService dataSourceService) {
        ServiceHolder.dataSourceService = dataSourceService;
    }

    public static ModelStoreManager getModelStoreManager() {
        return modelStoreManager;
    }

    public static void setModelStoreManager(ModelStoreManager modelStoreManager) {
        ServiceHolder.modelStoreManager = modelStoreManager;
    }

    public static ModelManager getModelManager() {
        return modelManager;
    }

    public static void setModelManager(ModelManager modelManager) {
        ServiceHolder.modelManager = modelManager;
    }

    public static ModelPeriodicProcessor getPeriodicProcessor() {
        return periodicProcessor;
    }

    public static void setPeriodicProcessor(ModelPeriodicProcessor periodicProcessor) {
        ServiceHolder.periodicProcessor = periodicProcessor;
    }
}
