#!/bin/sh
# ------------------------------------------------------------------------
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
# ------------------------------------------------------------------------
set -e

#Navigate to the wso2am-micro-gw-2.6.0/bin directory
cd /wso2am-micro-gw-2.6.0/bin

# $1-API project name
# $2-label
# $3-username
# $4-password
# $5-APIM base URL
# $6-Trust store location
# $7-Trust store password

./micro-gw setup $1 -l $2 -u $3 -p $4 -s $5 -t $6 -w $7 --insecure
./micro-gw build $1

# unzip target file
unzip /wso2am-micro-gw-2.6.0/bin/$1/target/micro-gw-$1.zip -d /

# renaming the generated folder to a common name to mount
mv /micro-gw-$1/* /target

# make gateway an executable file
chmod +x /target/bin/gateway
