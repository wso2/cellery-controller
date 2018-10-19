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

from docker.wso2.com/wso2sp-dashboard:4.3.0

ARG VICK_PROFILE_FILES=./files
ARG VICK_TARGET_FILES=./target/files

# copy entrypoint bash script to user home
COPY --chown=wso2carbon:wso2  ${VICK_TARGET_FILES}/widgets ${WSO2_SERVER_HOME}/wso2/dashboard/deployment/web-ui-apps/portal/extensions/widgets
COPY --chown=wso2carbon:wso2 ${VICK_PROFILE_FILES}/dashboard/* ${WSO2_SERVER_HOME}/wso2/dashboard/resources/dashboards/
COPY --chown=wso2carbon:wso2 ${VICK_PROFILE_FILES}/patches/widgets/OpenTracingList/widgetConf.json ${WSO2_SERVER_HOME}/wso2/dashboard/deployment/web-ui-apps/portal/extensions/widgets/OpenTracingList/

# expose ports
EXPOSE 9714 9643 9614 7713 7613