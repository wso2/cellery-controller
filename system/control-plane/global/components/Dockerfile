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
FROM wso2vick/wso2am:latest

# Copy the jar files
RUN rm -rf /home/wso2carbon/wso2am-2.6.0/repository/components/lib/org.wso2.vick.auth.*.jar
COPY files/*.jar /home/wso2carbon/wso2am-2.6.0/repository/components/lib/

RUN rm -rf /home/wso2carbon/wso2am-2.6.0/repository/deployment/server/webapps/api#identity#vick-auth#*.war
COPY files/*war /home/wso2carbon/wso2am-2.6.0/repository/deployment/server/webapps