# ------------------------------------------------------------------------
#
# Copyright 2018 WSO2, Inc. (http://wso2.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License
#
# ------------------------------------------------------------------------

FROM docker.wso2.com/wso2sp-worker:4.3.0

ARG VICK_TARGET_FILES=./target/files

# Remove connectors with bugs (The connectors with fixes are copied with the "libs" directory)
RUN rm ${WSO2_SERVER_HOME}/lib/siddhi-io-http-1.0.39.jar

# Copy the osgified jars
COPY --chown=wso2carbon:wso2 ${VICK_TARGET_FILES}/lib/*.jar ${WSO2_SERVER_HOME}/lib/

# Expose ports
EXPOSE 9090 9443 9712 9612 7711 7611 9091 9411 9092 9123
