#!/usr/bin/env bash

mvn clean install
docker build -t wso2vick/cell-sts .
docker push wso2vick/cell-sts

