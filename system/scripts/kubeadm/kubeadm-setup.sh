#!/bin/bash
#This script provides steps to install K8s using kubeadmin tool on Ubuntu 18.04 LTS.

UBUNTU_VERSION=16.04
K8S_VERSION=1.11.3-00
node_type=master

function vickHelp() {
    echo "Welcome to the Virtual Integration Computer for Kubernetes (VICK) setup."
    echo
    echo "Usage: "
    echo "  ./kubeadm-setup.sh -t [node-type] "
    echo
    echo -en "  -t\t"
    echo "[OPTIONAL] Kubernetes node type (master or worker). If not specified, the default node type will be master."

    echo
    echo "Run: "$'\e[1m'"./kubeadm-setup.sh -t worker "$'\e[0m'" to Install VICK with Kubernetes worker node"
    echo
    exit 1
}

while getopts :t FLAG; do
  case $FLAG in
    t)
      node_type=$OPTARG
      ;;
    \?)
      vickHelp
      ;;
  esac
done

if [ $# -eq 0 ]; then
    continue
elif [[ $1 != "-t" ]]; then
    vickHelp
fi

if grep -Fq $UBUNTU_VERSION /etc/os-release
then
	echo "Installing K8s version $K8S_VERSION on Ubuntu $UBUNTU_VERSION"
else
	echo "You need Ubuntu version $UBUNTU_VERSION to run this script"
	exit 0
fi

    #Update all installed packages.
    #apt-get update
    #apt-get upgrade

    #if you get an error similar to
    #'[ERROR Swap]: running with swap on is not supported. Please disable swap', disable swap:
    sudo swapoff -a

    # install some utils
    sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common

    #Install Docker
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

    if [ $UBUNTU_VERSION == "16.04" ]; then
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu xenial stable"
    elif [ $UBUNTU_VERSION == "18.04" ]; then
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
    else
        #default tested version
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu xenial stable"
    fi
    sudo apt-get update
    sudo apt-get install -y docker.io

    #Install NFS client
    sudo apt-get install -y nfs-common

    #Enable docker service
    sudo systemctl enable docker.service

    #Update the apt source list
    curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] http://apt.kubernetes.io/ kubernetes-xenial main"

    #Install K8s components
    sudo apt-get update
    sudo apt-get install -y kubelet=$K8S_VERSION kubeadm=$K8S_VERSION kubectl=$K8S_VERSION

    sudo apt-mark hold kubelet kubeadm kubectl

if [ $node_type == "master" ]; then
    #Initialize the k8s cluster
    sudo kubeadm init --pod-network-cidr=10.244.0.0/16

    sleep 60

    #Create .kube file if it does not exists
    mkdir -p $HOME/.kube

    #Move Kubernetes config file if it exists
    if [ -f $HOME/.kube/config ]; then
        mv $HOME/.kube/config $HOME/.kube/config.back
    fi

    sudo cp -f /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config

    #if you are using a single node which acts as both a master and a worker
    #untaint the node so that pods will get scheduled:
    kubectl taint nodes --all node-role.kubernetes.io/master-

    #Install Flannel network
    kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/v0.10.0/Documentation/kube-flannel.yml

    #Install admission plugins
    echo "Installing K8s admission plugins"
    sudo sed -i 's/--enable-admission-plugins=NodeRestriction/--enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota/' /etc/kubernetes/manifests/kube-apiserver.yaml
	
    #Wait to restart the K8s with new admission plugins
    sleep 60
  
   # Install Istio and ingress-enginx
   ./vick-init-kubeadm.sh 

   #Start to deploy VICK control plane 
    cd ../../control-plane/global/

    kubectl apply -f vick-ns-init.yaml
    HOST_NAME=$(hostname | tr '[:upper:]' '[:lower:]')

    kubectl label nodes $HOST_NAME disk=local

    #Create folders required by the mysql PVC
    sudo mkdir -p /mnt/mysql
    sudo chown 999:999 /mnt/mysql

    #Setup MySQL Server
    ./mysql-server-deploy.sh

    #Setup NFS Server
    #./nfs-server-deploy.sh
    
    sleep 60

    #Create folders required by the APIM PVC
    sudo mkdir -p /mnt/apim_repository_deployment_server
    sudo chown 802:802 /mnt/apim_repository_deployment_server

    #Install VICK control plane
    ./vick-control-plane-deploy.sh

    #Install VICK crds
    kubectl apply -f ../../../build/target/vick.yaml

    echo "K8s Master node installation is finished"

elif [ $node_type == "worker" ]; then
    read -p "Enter the Master node IP and the Token [master_node_ip token discovery_token_ca_cert_hash]:" master_node_ip token discovery_token_ca_cert_hash
    if [ -n "$master_node_ip" ] && [ -n "$token" ] && [ -n "$discovery_token_ca_cert_hash" ]; then
        echo $master_node_ip $token $discovery_token_ca_cert_hash
        #Add more worker nodes.
        sudo kubeadm join $master_node_ip:6443 --token $token --discovery-token-ca-cert-hash $discovery_token_ca_cert_hash
    else
        echo " Enter all three argument"
    fi
else
    echo "Enter correct arguments"
fi

