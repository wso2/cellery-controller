#!/usr/bin/env bash

mvn clean install
docker build -t wso2vick/cell-sts:latest .
docker push wso2vick/cell-sts:latest

