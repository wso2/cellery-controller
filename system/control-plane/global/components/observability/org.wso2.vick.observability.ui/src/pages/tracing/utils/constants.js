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

/**
 * Constants to be used by Constants UI
 */
const Constants = {
    Span: {
        Kind: {
            CLIENT: "CLIENT",
            SERVER: "SERVER",
            PRODUCER: "PRODUCER",
            CONSUMER: "CONSUMER"
        },
        ComponentType: {
            VICK: "VICK",
            ISTIO: "Istio",
            MICROSERVICE: "Micro-service"
        }
    },
    VICK: {
        Cell: {
            GATEWAY_NAME_PATTERN: /^src:\d+\.\d+\.\d+\.(.+)_(\d+_\d+_\d+)_.+$/,
            MICROSERVICE_NAME_PATTERN: /(.+)--(.+)/
        },
        System: {
            ISTIO_MIXER_NAME: "istio-mixer",
            GLOBAL_GATEWAY_NAME: "global-gateway"
        }
    }
};

export default Constants;
