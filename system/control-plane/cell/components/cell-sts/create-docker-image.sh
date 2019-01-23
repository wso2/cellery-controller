#!/usr/bin/env bash

mvn clean install
docker build -t wso2vick/cell-sts:test .
docker push wso2vick/cell-sts:test

