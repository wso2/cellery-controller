#!/bin/bash
echo "Removing docker kubeadm kubelet kubectl"
sudo kubeadm reset
sudo apt-get purge -y kubelet kubeadm kubectl docker.io docker-ce
