# VICK

[![Go Report Card](https://goreportcard.com/badge/github.com/wso2/product-vick)](https://goreportcard.com/report/github.com/wso2/product-vick)

VICK is a virtual distributed computer, that is optimised for agile development, deployment and operations of cell-based integration applications, running at scale.

## Installation Guide

### System Requirements

* CPU with minimum 4  cores
* 8 GB RAM
* Ubuntu 16.04 OS with SWAP disabled

### Installing VICK
Run the following command to install VICK and K8s on Ubunut Linux. This script uses [kubeadm](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/) tool to deploy K8s.
    
    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-setup-all-in-one.sh | bash -s -- kubeadm
    
### Uninstall VICK

    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-cleanup.sh | bash -s -- kubeadm

Please refer the [script usage document](./system/scripts/all-in-one-installer/README.md) to find the other K8s based platform support.

## Samples

Please see [Employee Portal App](./samples/employee-portal)

<!--
existing cluster installation
* Kubernetes v1.11.3 cluster with Istio 1.0.2 installed.
-->
