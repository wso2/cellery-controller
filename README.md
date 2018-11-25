# VICK

[![Go Report Card](https://goreportcard.com/badge/github.com/wso2/product-vick)](https://goreportcard.com/report/github.com/wso2/product-vick)

VICK is a virtual distributed computer, that is optimised for agile development, deployment and operations of cell-based integration applications, running at scale.

## Installation Guide

### System Requirements

* CPU with minimum 4  cores
* 8 GB RAM
* Ubuntu 16.04 OS with SWAP disabled

### Installing VICK
Run the following command to install VICK
    
    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/kubeadm/kubeadm-setup.sh | sudo bash
    
### Uninstall VICK

    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/kubeadm/kubeadm-cleanup.sh | sudo bash

## Samples

Please see [Employee Portal App](./samples/employee-portal)

<!--
existing cluster installation
* Kubernetes v1.11.3 cluster with Istio 1.0.2 installed.
-->
