#!/usr/bin/env bash

export CGO_ENABLED=0

go build -o vick-controller -x ./cmd/controller

docker build -t mirage20/vick-controller .

docker push mirage20/vick-controller