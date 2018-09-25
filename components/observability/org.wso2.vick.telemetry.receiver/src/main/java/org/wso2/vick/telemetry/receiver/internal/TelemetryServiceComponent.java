/*
 *  Copyright (c) 2018 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  WSO2 Inc. licenses this file to you under the Apache License,
 *  Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 *
 */

package org.wso2.vick.telemetry.receiver.internal;

import org.apache.log4j.Logger;
import org.osgi.framework.BundleContext;
import org.osgi.service.component.annotations.Activate;
import org.osgi.service.component.annotations.Component;


/**
 * This is the internal class that is used to activate the telemtry receiver component, and starts the GRPC server.
 */
@Component(
        name = "org.wso2.vick.telemetry.receiver.internal.TelemetryServiceComponent",
        service = TelemetryServiceComponent.class,
        immediate = true
)
public class TelemetryServiceComponent {

    private static final Logger log = Logger.getLogger(TelemetryServiceComponent.class);

    @Activate
    protected void start(BundleContext bundleContext) throws Exception {
       if (log.isDebugEnabled()) {
           log.debug("Telemetry service component started");
       }
    }

}
