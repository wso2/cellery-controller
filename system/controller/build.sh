#!/usr/bin/env bash

export CGO_ENABLED=0

go build -o vick-controller -x ./cmd/controller

docker build -t wso2vick/vick-controller .

docker push wso2vick/vick-controller
