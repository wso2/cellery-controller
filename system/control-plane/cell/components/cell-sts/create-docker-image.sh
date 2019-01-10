#!/usr/bin/env bash

mvn clean install
docker build -t wso2vick/cell-sts:envoy-cell-gw .
docker push wso2vick/cell-sts:envoy-cell-gw

