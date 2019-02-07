#  Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
#
#  WSO2 Inc. licenses this file to you under the Apache License,
#  Version 2.0 (the "License"); you may not use this file except
#  in compliance with the License.
#  You may obtain a copy of the License at
#
#  http:www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.

PROJECT_ROOT := $(realpath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
PROJECT_PKG := github.com/cellery-io/mesh-controller
BUILD_DIRECTORY := build
BUILD_ROOT := $(PROJECT_ROOT)/$(BUILD_DIRECTORY)
GOFILES		= $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./pkg/client/*")

MAIN_PACKAGES := controller
BUILD_TARGETS := $(addprefix build., $(MAIN_PACKAGES))
TEST_TARGETS := $(addprefix test., $(MAIN_PACKAGES))
DOCKER_TARGETS := $(addprefix docker., $(MAIN_PACKAGES))
DOCKER_PUSH_TARGETS := $(addprefix docker-push., $(MAIN_PACKAGES))
DOCKER_REPO := celleryio
DOCKER_IMAGE_PREFIX := mesh
DOCKER_IMAGE_TAG := latest

all: build

.PHONY: $(BUILD_TARGETS)
$(BUILD_TARGETS):
	$(eval TARGET=$(patsubst build.%,%,$@))
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_ROOT)/$(TARGET) -x $(PROJECT_ROOT)/cmd/$(TARGET)

.PHONY: build
build: $(BUILD_TARGETS)

.PHONY: $(TEST_TARGETS)
$(TEST_TARGETS):
	$(eval TARGET=$(patsubst test.%,%,$@))
	go test -covermode=count -coverprofile=$(BUILD_ROOT)/coverage.out ./pkg/$(TARGET)/...

.PHONY: test
test: $(TEST_TARGETS)

coverage: test
	go tool cover -html=$(BUILD_ROOT)/coverage.out

.PHONY: $(DOCKER_TARGETS)
$(DOCKER_TARGETS): docker.% : build.%
	$(eval TARGET=$(patsubst docker.%,%,$@))
	docker build -f $(PROJECT_ROOT)/docker/$(TARGET)/Dockerfile $(BUILD_ROOT) -t $(DOCKER_REPO)/$(DOCKER_IMAGE_PREFIX)-$(TARGET):$(DOCKER_IMAGE_TAG)

.PHONY: docker
docker: $(DOCKER_TARGETS)

.PHONY: $(DOCKER_PUSH_TARGETS)
$(DOCKER_PUSH_TARGETS): docker-push.% : docker.%
	$(eval TARGET=$(patsubst docker-push.%,%,$@))
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE_PREFIX)-$(TARGET):$(DOCKER_IMAGE_TAG)

.PHONY: docker-push
docker-push: $(DOCKER_PUSH_TARGETS)

.PHONY: clean
clean:
	rm -rf $(BUILD_ROOT)

.PHONY: code.format
code.format: tools.goimports
	@goimports -local $(PROJECT_PKG) -w -l $(GOFILES)

.PHONY: code.format-check
code.format-check: tools.goimports
	@goimports -local $(PROJECT_PKG) -l $(GOFILES)

.PHONY: tools tools.goimports

tools: tools.goimports

tools.goimports:
	@command -v goimports >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "goimports not found. Running 'go get golang.org/x/tools/cmd/goimports'"; \
		go get golang.org/x/tools/cmd/goimports; \
	fi
