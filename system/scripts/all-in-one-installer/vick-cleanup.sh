#!/usr/bin/env bash

#cat vick-cleanup.sh |bash -s -- kubeadm

iaas=$1

if [ -z $iaas ]; then
    echo "Please provide the K8s provider. [ GCP | kubeadm ]."
    exit 1
fi

if [ $iaas == "kubeadm" ]; then

    echo "Removing docker kubeadm kubelet kubectl."
    echo
    read -p "Do you want to purge kubelet kubectl and docker [y/N]: " deb_remove_option < /dev/tty

    if [ $deb_remove_option == "y" ]; then
        DEL_LEVEL="purge"
    else
        DEL_LEVEL="remove"
    fi

    yes | sudo kubeadm reset

    sudo apt-get $DEL_LEVEL --allow-change-held-packages kubelet kubeadm kubectl docker.io docker-ce

    echo "Removing /mnt/mysql and /mnt/apim_repository_deployment_server."
    echo

    if [ -d /mnt/mysql ]; then
        sudo mv /mnt/mysql "/mnt/mysql.$(date +%s)"
    fi

    if [ -d /mnt/apim_repository_deployment_server ]; then
        sudo mv /mnt/apim_repository_deployment_server "/mnt/apim_repository_deployment_server.$(date +%s)"
    fi
fi

echo "VICK cleanup is finished."
