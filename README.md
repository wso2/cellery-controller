# VICK

[![Go Report Card](https://goreportcard.com/badge/github.com/wso2/product-vick)](https://goreportcard.com/report/github.com/wso2/product-vick)

VICK is a virtual distributed computer, that is optimised for agile development, deployment and operations of [cell-based](https://wso2.com/blogs/thesource/2018/07/introducing-the-wso2-integration-agile-platform/) integration applications, running at scale.

## Installation Guide

VICK installation guide provides you steps to deploy VICK to a production grade K8s cluster and steps to deploy VICK environment for development purposes. 

## Steps to Deploy Single Node VICK Developer Environment

### Minimum System Requirements 

* CPU with minimum 4  cores
* 8 GB RAM
* Ubuntu 16.04 OS with SWAP disabled

### Installing VICK
Run the following command to install VICK and K8s on Ubunut Linux 16.04. This script uses [kubeadm](https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/) tool to deploy K8s.
User has the option to install MySQL Server before installing the VICK or point the VICK installation to an existing MySQL server deployment. 
    
    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-setup-all-in-one.sh | bash -s -- kubeadm
    
### Uninstall VICK

    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-cleanup.sh | bash -s -- kubeadm

Please refer the [script usage document](./system/scripts/all-in-one-installer/README.md) to find the other K8s based platform support.

## Steps to Deploy a Production Quality VICK Environment

### System Requirements 

* K8s deployment done according to K8s [production guidelines](https://kubernetes.io/docs/setup/scratch/)
* MySQL 5.7 server access with database and user creation permissions.
* K8s storage volume compatible network storage (NFS)

### Installing VICK
Run the following command to install VICK in to the existing K8s cluster. Scrip will be prompt the user for the mysql user and password credentials. Also user has to provide 
the network file share server IP (NFS) and the data folder path.
    
    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-setup-all-in-one.sh | bash
    
### Uninstall VICK

    curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-cleanup.sh | bash

Please refer the [script usage document](./system/scripts/all-in-one-installer/README.md) to find the other K8s based platform support.


## Samples

Please see [Employee Portal App](./samples/employee-portal)

<!--
existing cluster installation
* Kubernetes v1.11.3 cluster with Istio 1.0.2 installed.
-->
