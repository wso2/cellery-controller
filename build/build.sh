#!/usr/bin/env bash

export CGO_ENABLED=0

go build -o ./target/vick-controller -x ./../system/controller/cmd/controller/
