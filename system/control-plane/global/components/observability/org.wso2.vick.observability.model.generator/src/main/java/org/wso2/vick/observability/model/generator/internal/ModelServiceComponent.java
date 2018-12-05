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

import org.apache.log4j.Logger;
import org.osgi.framework.BundleContext;
import org.osgi.service.component.annotations.Activate;
import org.osgi.service.component.annotations.Component;
import org.osgi.service.component.annotations.Reference;
import org.osgi.service.component.annotations.ReferenceCardinality;
import org.osgi.service.component.annotations.ReferencePolicy;
import org.wso2.carbon.datasource.core.api.DataSourceService;
import org.wso2.vick.observability.model.generator.ModelManager;

import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

/**
 * This class acts as a ServiceComponent which specifies the services that is required by the component.
 */
@Component(
        service = ModelServiceComponent.class,
        immediate = true
)
public class ModelServiceComponent {
    private static final Logger log = Logger.getLogger(ModelServiceComponent.class);

    @Activate
    protected void start(BundleContext bundleContext) throws Exception {
        try {
            ServiceHolder.setGraphStoreManager(new GraphStoreManager());
            ServiceHolder.setModelManager(new ModelManager());
            bundleContext.registerService(ModelManager.class.getName(), ServiceHolder.getModelManager(), null);
            Executors.newScheduledThreadPool(1).scheduleAtFixedRate(new ScheduledDependencyGraphStore(),
                    1, 1, TimeUnit.MINUTES);
        } catch (Throwable throwable) {
            log.error("Error occured while activating the model generation bundle", throwable);
            throw throwable;
        }
    }

    @Reference(
            name = "org.wso2.carbon.datasource.DataSourceService",
            service = DataSourceService.class,
            cardinality = ReferenceCardinality.AT_LEAST_ONE,
            policy = ReferencePolicy.DYNAMIC,
            unbind = "unregisterDataSourceService"
    )
    protected void registerDataSourceService(DataSourceService service) {
        ServiceHolder.setDataSourceService(service);
    }

    protected void unregisterDataSourceService(DataSourceService service) {
        ServiceHolder.setDataSourceService(null);
    }

}
