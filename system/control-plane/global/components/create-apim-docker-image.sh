#!/usr/bin/env bash

# Cleanup previous artifacts
rm -rf files
mkdir files

mvn clean install

# Copy the artifacts to files folder
cp ./org.wso2.vick.auth.extensions/target/org.wso2.vick.auth.extensions-*.jar ./files/
cp ./org.wso2.vick.auth.sts.core/target/org.wso2.vick.auth.sts.core-*.jar ./files/
cp ./org.wso2.vick.auth.sts.endpoint/target/api#identity#vick-auth#*.war ./files/

docker build -t wso2vick/wso2am:latest .
docker push wso2vick/wso2am:latest
