# Mesh Controller

[![Build Status](https://travis-ci.org/wso2-cellery/mesh-controller.svg?branch=master)](https://travis-ci.org/wso2-cellery/mesh-controller)
[![codecov.io](https://codecov.io/gh/wso2-cellery/mesh-controller/branch/master/graph/badge.svg)](https://codecov.io/gh/wso2-cellery/mesh-controller/branch/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/cellery-io/mesh-controller)](https://goreportcard.com/report/github.com/cellery-io/mesh-controller)
[![GitHub (pre-)release](https://img.shields.io/github/release/cellery-io/mesh-controller/all.svg)](https://github.com/cellery-io/mesh-controller/releases)
[![GitHub (Pre-)Release Date](https://img.shields.io/github/release-date-pre/cellery-io/mesh-controller.svg)](https://github.com/cellery-io/mesh-controller/releases)
[![GitHub last commit](https://img.shields.io/github/last-commit/cellery-io/mesh-controller.svg)](https://github.com/cellery-io/mesh-controller/commits/master)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
  
**Note**: If you just wish to get started with trying out Cellery, you can start your journey at [wso2-cellery/sdk](https://github.com/wso2-cellery/sdk) repository.  

Mesh controller is a Kubernetes implementation for the Cellery mesh (a.k.a Cell Mesh) builds using Istio. To learn more about about Cellery, check out the [cellery-sdk](https://github.com/cellery-io/sdk) repository.

## Getting Started

Mesh Controller requires a Kubernetes cluster with Istio installed in order to work properly.

### Prerequisites

* Golang v1.11+
* Kubernetes cluster and client v1.11 
* Istio v1.0.2
* GNU Make 4.1+
* Docker

### Building and Running

1. Clone this repository to your $GOPATH
   
    ```bash
    cd $GOPATH/src/github.com/cellery-io/
    git clone https://github.com/cellery-io/mesh-controller.git
    cd mesh-controller
    ```
   
2. Build the controller

    ```bash
    # The output binary can be located inside the build directory
    make build.controller 
    ```

3. Run the controller in your local machine (This requires you to provide the kubeconfig file in order to connect to the Kubernetes cluster)

    ```bash
    ./build/controller -logtostderr=true --kubeconfig=/home/<username>/.kube/config
    ```
    
## Running the tests

### Unit tests

1. Run the unit tests

    ```bash
    make test.controller
    ```
2. Generate and view coverage report

    ```bash
    make coverage 
    ```
    
## Deployment

This will guide you on deploying the Mesh Controller inside the Kubernetes cluster
 
### Deploying a released version
 
You can deploy a released version of the Mesh Controller using following steps. 
 
1. Generate Kubernetes YAML’s for the controller
    
    ```bash
    # Specify a released version to generate the yaml
    VERSION=v0.1.0 make artifacts
    ```
2. Deploy the controller into Kubernetes cluster
    
    ```bash
    kubectl apply -f build/mesh-controller.yaml
    ```
### Deploying a development version
 
If you want to deploy a development version, you can to generate a Kubernetes YAML which uses docker images from your Docker repository.
 
1. Build and push the docker image
    
    ```bash
    DOCKER_REPO=<your-repo> VERSION=dev make docker-push.controller artifacts
    ```
    
2. Deploy controller into the Kubernetes cluster
    
    ```bash
    kubectl apply -f build/mesh-controller.yaml
    ```

## Running a sample Cell

You can create a sample cell using following command,

```bash
cat <<EOF | kubectl create -f -
apiVersion: mesh.cellery.io/v1alpha1
kind: Cell
metadata:
  name: my-cell
spec:
  gatewayTemplate:
    spec:
      http:
      - context: time
        definitions:
        - path: /
          method: GET
        backend: server-time
        global: false
      - context: hello
        definitions:
        - path: /
          method: GET
        backend: node-hello
        global: false
  servicesTemplates:
  - metadata:
      name: time-us
    spec:
      replicas: 1
      container:
        image: docker.io/mirage20/time-us
        ports:
        - containerPort: 8080
      servicePort: 80
  - metadata:
      name: time-uk
    spec:
      replicas: 1
      container:
        image: docker.io/mirage20/time-uk
        ports:
        - containerPort: 8080
      servicePort: 80
  - metadata:
      name: server-time
    spec:
      replicas: 1
      container:
        image: docker.io/mirage20/time
        ports:
        - containerPort: 8080
      servicePort: 80
  - metadata:
      name: debug
    spec:
      replicas: 1
      container:
        image: docker.io/mirage20/k8s-debug-tools
      servicePort: 80
EOF
```

## Contributing

If you are interested in contributing please see [CONTRIBUTING.md](CONTRIBUTING.md)
