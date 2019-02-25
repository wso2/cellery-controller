# Mesh Controller

[![Go Report Card](https://goreportcard.com/badge/github.com/cellery-io/mesh-controller)](https://goreportcard.com/report/github.com/cellery-io/mesh-controller)

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
   
    ```console
    cd $GOPATH/src/github.com/cellery-io/
    git clone https://github.com/cellery-io/mesh-controller.git
    cd mesh-controller
    ```
   
2. Build the controller

    ```console
    # The output binary can be located inside the build directory
    make build.controller 
    ```

3. Run the controller in your local machine (This requires you to provide the kubeconfig file in order to connect to the Kubernetes cluster)

    ```console
    ./build/controller -logtostderr=true --kubeconfig=/home/<username>/.kube/config
    ```
    
## Running the tests

### Unit tests

1. Run the unit tests

    ```console
    make test.controller
    ```
2. Generate and view coverage report

    ```console
    make coverage 
    ```
    
## Deployment

This will guide you on deploying the Mesh Controller inside the Kubernetes cluster
 
### Deploying a released version
 
You can deploy a released version of the Mesh Controller using following steps. 
 
1. Generate Kubernetes YAML’s for the controller
    
    ```console
    # Specify a released version to generate the yaml
    VERSION=v0.1.0 make artifacts
    ```
2. Deploy the controller into Kubernetes cluster
    
    ```console
    kubectl apply -f build/mesh-controller.yaml
    ```
### Deploying a development version
 
If you want to deploy a development version, you can to generate a Kubernetes YAML which uses docker images from your Docker repository.
 
1. Build and push the docker image
    
    ```console
    DOCKER_REPO=<your-repo> VERSION=dev make docker-push.controller artifacts
    ```
    
2. Deploy controller into the Kubernetes cluster
    
    ```console
    kubectl apply -f build/mesh-controller.yaml
    ```

## Running a sample Cell

You can create a sample cell using following command,

```console
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
