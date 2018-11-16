#!/bin/bash

#cat vick-cleanup.sh |bash -s -- GCP

iaas=$1

if [ -z $iaas ]; then
    echo "Please provide the K8s provider. [ GCP | kubeadm ]"
    exit 0
fi

if [ $iaas == "kubeadm" ]; then

    echo "Removing docker kubeadm kubelet kubectl"
    echo
    read -p "Do you want to purge kubelet kubectl and docker [y/N]" deb_remove_option
    echo

    sudo kubeadm reset

    if [ $deb_remove_option == "y" ]; then
        DEL_LEVEL="purge"
    else
        DEL_LEVEL="remove"
    fi

    sudo apt-get $DEL_LEVEL kubelet kubeadm kubectl docker.io docker-ce

    echo "Removing /mnt/mysql and /mnt/apim_repository_deployment_server"
    echo

    sudo rm -fr /mnt/mysql
    sudo rm -fr /mnt/apim_repository_deployment_server
fi

echo "VICK cleanup is finished"
