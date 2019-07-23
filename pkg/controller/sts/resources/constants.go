/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package resources

const (
	tokenServiceServiceInboundPort  = 8080
	tokenServiceServiceOutboundPort = 8081

	tokenServiceServicePortInboundName  = "grpc-inbound"
	tokenServiceServicePortOutboundName = "grpc-outbound"
	tokenServiceServicePortJWKSName     = "http-jwks"

	tokenServiceContainerInboundPort  = 8080
	tokenServiceContainerOutboundPort = 8081
	tokenServiceContainerJWKSPort     = 8090
	opaServicePort                    = 8181

	unsecuredPathsConfigKey  = "unsecured-paths"
	unsecuredPathsConfigFile = "unsecured-paths.json"

	tokenServiceConfigKey  = "sts-config"
	tokenServiceConfigFile = "sts.json"

	configVolumeName = "config-volume"
	configMountPath  = "/etc/config"

	policyVolumeName      = "cell-policy"
	pocliyConfigMountPath = "/policies"

	keyPairVolumeName = "cell-keys"
	keyPairMountPath  = "/etc/certs"

	caCertsVolumeName = "ca-certs"
	caCertsMountPath  = "/etc/certs/trusted-certs"

	policyConfigKey  = "policy.rego"
	policyConfigFile = "sample.rego"

	envCellNameKey = "CELL_NAME"

	// Envoy filter
	filterInsertPositionFirst  = "FIRST"
	filterInsertPositionLast   = "LAST"
	filterInsertPositionBefore = "BEFORE"
	filterInsertPositionAfter  = "AFTER"

	filterMixer = "mixer"

	filterListenerTypeInbound  = "SIDECAR_INBOUND"
	filterListenerTypeOutbound = "SIDECAR_OUTBOUND"

	HTTPProtocol   = "HTTP"
	baseFilterName = "envoy.ext_authz"
	statPrefix     = "ext_authz"
	filterTimeout  = "10s"
)
