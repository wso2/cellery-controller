#!/bin/bash

#This script provides steps to install K8s using kubeadmin tool on Ubuntu.

#Install Ubuntu 18.04

#Update all installed packages.
#apt-get update
#apt-get upgrade

#Install Docker
#apt-get install -y docker.io

#Install NFS client
#apt-get install nfs-common

#Enable docker service
#systemctl enable docker.service

#Install curl
#apt-get update && apt-get install -y apt-transport-https curl

#Update the apt source list
#curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
#cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
# deb http://apt.kubernetes.io/ kubernetes-xenial main
#EOF

#Install K8s components
#apt-get update
#apt-get install -y kubelet kubeadm kubectl
#apt-mark hold kubelet kubeadm kubectl

#Initialize the k8s cluster
#kubeadm init --pod-network-cidr=10.244.0.0/16

#Install Flannel network
#kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/v0.10.0/Documentation/kube-flannel.yml

#Add more worker nodes.
#kubeadm join <Master Node IP>:6443 --token <token> --discovery-token-ca-cert-hash <discovery-token-ca-cert-hash>

#Install Istio
echo "Installing Istio"
#ISTIO_HOME=/Users/vickdev/wso2/VICK/vick-init/istio-1.0.2
export PATH=$ISTIO_HOME/bin:$PATH
kubectl apply -f $ISTIO_HOME/install/kubernetes/istio-demo-auth.yaml

#Wait for 30 seconds
sleep 30s

#Enabling Istio injection
kubectl label namespace default istio-injection=enabled

#Install Knative
#echo "Installing Knative"
#kubectl apply -f https://github.com/knative/serving/releases/download/v0.1.1/release.yaml

#KNATIVE_HOME=/Users/vickdev/wso2/VICK/vick-init/knative-0.1.1
#kubectl apply -f $KNATIVE_HOME/release.yaml

#Install nginx-ingress for control plane ingress
kubectl apply -f mandatory.yaml
kubectl apply -f service-nodeport.yaml

#Init VICK
#echo "Installing VICK"
#kubectl apply -f $VICK_HOME/vic-init.yaml
