#  Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
#
#  WSO2 Inc. licenses this file to you under the Apache License,
#  Version 2.0 (the "License"); you may not use this file except
#  in compliance with the License.
#  You may obtain a copy of the License at
#
#  http:www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.


PROJECT_ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BUILD_ROOT := $(PROJECT_ROOT)/build/target

all: controller global-api-updater cell-sts

controller:
	go build -o ${BUILD_ROOT}/vick-controller -x ./system/controller/cmd/controller/


global-api-updater:
	mvn clean install -f system/control-plane/cell/components/global-api-updater/pom.xml

cell-sts:
	mvn clean install -f system/control-plane/cell/components/cell-sts/pom.xml



