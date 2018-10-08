#!/bin/bash
#This script provides steps to install K8s using kubeadmin tool on Ubuntu 18.04 LTS.

#Install Ubuntu 18.04

echo "Installing K8s cluster using Kubeadm"
read -p "Enter the node type [master/worker]:" node_type

if [ $node_type == "master" ]; then
    echo $node_type
    exit 0

    #Update all installed packages.
    apt-get update
    apt-get upgrade

    #Install Docker
    apt-get install -y docker.io

    #Install NFS client
    apt-get install nfs-common

    #Enable docker service
    systemctl enable docker.service

    #Install curl
    apt-get update && apt-get install -y apt-transport-https curl

    #Update the apt source list
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
    echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" >> /etc/apt/sources.list.d/kubernetes.list

    #Install K8s components
    apt-get update
    apt-get install -y kubelet kubeadm kubectl
    apt-mark hold kubelet kubeadm kubectl

    #Initialize the k8s cluster
    kubeadm init --pod-network-cidr=10.244.0.0/16

    #if you are using a single node which acts as both a master and a worker
    #untaint the node so that pods will get scheduled:
    kubectl taint nodes --all node-role.kubernetes.io/master-

    #if you get an error similar to
    #'[ERROR Swap]: running with swap on is not supported. Please disable swap', disable swap:
    swapoff -a

    #Install Flannel network
    kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/v0.10.0/Documentation/kube-flannel.yml

elif [ $node_type == "worker" ]; then
    read -p "Enter the Master node IP and the Token [master_node_ip token discovery_token_ca_cert_hash]:" master_node_ip token discovery_token_ca_cert_hash
    if [ -n "$master_node_ip" ] && [ -n "$token" ] && [ -n "$discovery_token_ca_cert_hash" ]; then
        echo $master_node_ip $token $discovery_token_ca_cert_hash
        exit 0
        #Add more worker nodes.
        kubeadm join $master_node_ip:6443 --token $token --discovery-token-ca-cert-hash $discovery_token_ca_cert_hash
    else
        echo " Enter all three argument"
    fi
else
    echo "Enter correct arguments"
fi