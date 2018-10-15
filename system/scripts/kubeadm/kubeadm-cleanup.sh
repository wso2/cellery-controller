#!/bin/bash
echo "Removing docker kubeadm kubelet kubectl"
echo
read -p "Do you want to purge kubelet kubectl and docker [Y/n]" deb_remove_option
echo

sudo kubeadm reset

if [ $deb_remove_option == "Y" ]; then
    DEL_LEVEL="purge"
else
    DEL_LEVEL="remove"
fi

sudo apt-get $DEL_LEVEL kubelet kubeadm kubectl docker.io docker-ce

echo "Removing /mnt/mysql and /mnt/apim_repository_deployment_server"
echo

sudo rm -fr /mnt/mysql
sudo rm -fr /mnt/apim_repository_deployment_server

echo "VICK cleanup is finished"
