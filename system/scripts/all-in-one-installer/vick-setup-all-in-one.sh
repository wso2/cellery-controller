#!/usr/bin/env bash
# ------------------------------------------------------------------------
#
# Copyright 2018 WSO2, Inc. (http://wso2.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License
#
# ------------------------------------------------------------------------

#cat vick-setup-all-in-one.sh |bash -s -- GCP knative-deep

function install_k8s_kubeadm () {
    K8S_VERSION=$1
    UBUNTU_VERSION=$(cat /etc/lsb-release | sed -n '/DISTRIB_RELEASE/p' | awk -F'=' '{print $2}')

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
    sudo apt-get install -y kubelet=$K8S_VERSION kubeadm=$K8S_VERSION kubectl=$K8S_VERSION
    sudo apt-mark hold kubelet kubeadm kubectl
}

function configure_k8s_kubeadm () {
node_type=$1
if [ -z $node_type ]; then
    node_type="master"
fi

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

}

#Create new NFS share
function create_nfs_share_gcp () {
echo "Creating NFS share in GCP"
local nfs_share_location=$1
local nfs_server_ip

gcloud beta filestore instances create nfs-server \
    --project=knative-deep \
    --location=us-west1-c \
    --tier=STANDARD \
    --file-share=name="${nfs_share_location}",capacity=1TB \
    --network=name="default"

sleep 30

nfs_server_ip=$(gcloud beta filestore instances describe nfs-server --project=knative-deep --location=us-west1-c \
--format flattened | awk '/ipAddresses/  {print $2}')

if [ -n $nfs_server_ip ]; then
    nfs_config_params["NFS_SERVER_IP"]=$nfs_server_ip
    nfs_config_params["NFS_SHARE_LOCATION"]="/${nfs_share_location}"
else
    echo "Error in NFS share creation"
    exit 0
fi
}

#Read NFS connection details form the user
function read_nfs_connection () {
echo "Read NFS connection"
    local nfs_server_ip
    local nfs_share_location

    echo "Configuring NFS volume"
    echo
    read -p "NFS server IP: " nfs_server_ip < /dev/tty
    if [[ ! -z "${nfs_server_ip/ //}" ]]; then
        nfs_config_params["NFS_SERVER_IP"]=$nfs_server_ip
    fi
    read -p "NFS share location: " nfs_share_location < /dev/tty
    if [[ ! -z "${nfs_share_location/ //}" ]]; then
        nfs_config_params["NFS_SHARE_LOCATION"]=$nfs_share_location
    fi
    echo
}

#Update APIM PVs
function update_apim_nfs_volumes () {
 local download_location=$1
    for param in "${!nfs_config_params[@]}"
    do
        sed -i '' "s|$param|${nfs_config_params[$param]}|g" ${download_location}/vick-apim-artifacts-persistent-volumes.yaml
    done
}

#Create K8s cluster in GCP
#User need to configure GCP CLI in the machine
function install_k8s_gcp () {
gcp_project=$1
#create GCP project
#gcloud projects create $gcp_project --set-as-default

#Point GCP to the new project
gcloud config set project $gcp_project

#Enable required GCP APIs
gcloud services enable \
  cloudapis.googleapis.com \
  container.googleapis.com \
  containerregistry.googleapis.com

CLUSTER_NAME=vick-knative
CLUSTER_ZONE=us-west1-c

echo "Creating K8s cluster $CLUSTER_NAM in in zone $CLUSTER_ZONE"

#Create K8s cluster
gcloud -q container clusters create $CLUSTER_NAME \
  --zone=$CLUSTER_ZONE \
  --cluster-version=latest \
  --machine-type=n1-standard-4 \
  --enable-autoscaling --min-nodes=1 --max-nodes=10 \
  --enable-autorepair \
  --scopes=service-control,service-management,compute-rw,storage-ro,cloud-platform,logging-write,monitoring-write,pubsub,datastore \
  --num-nodes=1

#Wait for 60 seconds
sleep 60s

#Grant cluster-admin permissions to the current user
K8S_ADMIN=$(gcloud config get-value core/account)

echo "Grant K8s cluster admin persmission to user $K8S_ADMIN"

kubectl create clusterrolebinding cluster-admin-binding \
--clusterrole=cluster-admin \
--user=$(gcloud config get-value core/account)

}

function deploy_mysql_server () {
    local download_location=$1
    #Create folders required by the mysql PVC
    if [ -d /mnt/mysql ]; then
        sudo mv /mnt/mysql "/mnt/mysql.$(date +%s)"
    fi
    sudo mkdir -p /mnt/mysql
    #Change the folder ownership to mysql server user.
    sudo chown 999:999 /mnt/mysql

    kubectl create configmap mysql-dbscripts --from-file=${download_location}/mysql/dbscripts/ -n vick-system
    kubectl apply -f ${download_location}/mysql-persistent-volumes-local.yaml -n vick-system
    kubectl apply -f ${download_location}/mysql-persistent-volume-claim.yaml -n vick-system
    kubectl apply -f ${download_location}/mysql-deployment.yaml -n vick-system
    #Wait till the mysql deployment availability
    kubectl wait deployment/wso2apim-with-analytics-mysql-deployment --for condition=available --timeout=6000s -n vick-system
    kubectl apply -f ${download_location}/mysql-service.yaml -n vick-system
}

function deploy_mysql_server_gcp () {
    local download_location=$1
    local sql_instance_name=$2
    local service_account
    local mysql_server_ip
    gcloud -q sql instances create ${sql_instance_name} --tier=db-n1-standard-1 --gce-zone=us-west1-c
    service_account=$(gcloud beta sql instances describe ${sql_instance_name} --format flattened | awk '/serviceAccountEmailAddress/ {print $2}')
    #if service account is zero exit
    gsutil -q mb --retention 3600s -l us-west1 gs://vickdb
    gsutil cp ${download_location}/mysql/dbscripts/init.sql gs://vickdb/init.sql
    gsutil acl ch -u ${service_account}:R gs://vickdb/init.sql
    gcloud -q sql import sql ${sql_instance_name} gs://vickdb/init.sql
    gcloud -q sql instances patch ${sql_instance_name} --authorized-networks=0.0.0.0/0

    mysql_server_ip=$(gcloud beta sql instances describe ${sql_instance_name} --format flattened | awk '/.ipAddress/ {print $2}')
    config_params["MYSQL_DATABASE_HOST"]=$mysql_server_ip
}

#Configure remote mysql server endpoint in control plane configuration files.
function read_control_plane_datasources_configs () {
    local db_user
    local db_passwd
    local db_hostname="wso2apim-with-analytics-rdbms-service"

    if [ $iaas == "kubeadm" ]; then
    echo "Configuring remote MySQL server"
    #read -p "Database host name: " db_hostname < /dev/tty
        if [[ ! -z "${db_hostname/ //}" ]]; then
                echo "Db Hostname ${db_hostname}"
                config_params["MYSQL_DATABASE_HOST"]=$db_hostname
        fi
    fi
    read -p "Database user name: " db_user < /dev/tty
    if [[ ! -z "${db_user/ //}" ]]; then
            config_params["DATABASE_USERNAME"]=$db_user
    fi
    read -p "Database user password: " db_passwd < /dev/tty
    if [[ ! -z "${db_passwd/ //}" ]]; then
            config_params["DATABASE_PASSWORD"]=$db_passwd
    fi
}

function update_control_plance_sql () {
local download_location=$1
    for param in "${!config_params[@]}"
    do
        sed -i "s/$param/${config_params[$param]}/g" ${download_location}/mysql/dbscripts/init.sql
    done

}

#Update pub-store/gw/sp worker/ sp dashboard datasources
function update_control_plane_datasources () {
    local download_location=$1
    for param in "${!config_params[@]}"
    do
        sed -i "s/$param/${config_params[$param]}/g" ${download_location}/apim-configs/pub-store/datasources/master-datasources.xml
        sed -i "s/$param/${config_params[$param]}/g" ${download_location}/apim-configs/gw/datasources/master-datasources.xml
        sed -i "s/$param/${config_params[$param]}/g" ${download_location}/sp-worker/conf/deployment.yaml
        sed -i "s/$param/${config_params[$param]}/g" ${download_location}/sp-dashboard/conf/deployment.yaml
    done
}

function deploy_global_gw () {

    local download_location=$1
    local iaas=$2

    if [ $iaas == "kubeadm" ]; then
        #Create folders required by the APIM GW PVC
        if [ -d /mnt/apim_repository_deployment_server ]; then
            sudo mv /mnt/apim_repository_deployment_server "/mnt/apim_repository_deployment_server.$(date +%s)"
        fi
        #Create folders required by the APIM PVC
        sudo mkdir -p /mnt/apim_repository_deployment_server
        sudo chown 802:802 /mnt/apim_repository_deployment_server
        #Create apim local volumes and volume claims
        kubectl apply -f ${download_location}/vick-apim-persistent-volumes-local.yaml -n vick-system
        kubectl apply -f ${download_location}/vick-apim-persistent-volume-claim-local.yaml -n vick-system
    elif [ $iaas == "GCP" ]; then
        #Create apim NFS volumes and volume claims
        kubectl apply -f ${download_location}/vick-apim-artifacts-persistent-volumes.yaml -n vick-system
        kubectl apply -f ${download_location}/vick-apim-artifacts-persistent-volume-claim.yaml -n vick-system
    fi

    #Create the gw config maps
    kubectl create configmap gw-conf --from-file=${download_location}/apim-configs/gw -n vick-system
    kubectl create configmap gw-conf-datasources --from-file=${download_location}/apim-configs/gw/datasources/ -n vick-system
    #Create KM config maps
    kubectl create configmap conf-identity --from-file=${download_location}/apim-configs/gw/identity -n vick-system
    kubectl create configmap apim-template --from-file=${download_location}/apim-configs/gw/resources/api_templates -n vick-system
    kubectl create configmap apim-tomcat --from-file=${download_location}/apim-configs/gw/tomcat -n vick-system
    kubectl create configmap apim-security --from-file=${download_location}/apim-configs/gw/security -n vick-system

    #Create gateway deployment and the service
    kubectl apply -f ${download_location}/vick-apim-gw.yaml -n vick-system
     #Wait till the gateway deployment availability
    kubectl wait deployment/gateway --for condition=available --timeout=6000s -n vick-system
    #Create gateway ingress
    kubectl apply -f ${download_location}/vick-apim-gw-ingress.yaml -n vick-system
}

function deploy_sp_dashboard_worker () {
    local download_location=$1
    #Create SP worker configmaps
    kubectl create configmap sp-worker-siddhi --from-file=${download_location}/sp-worker/siddhi -n vick-system
    kubectl create configmap sp-worker-conf --from-file=${download_location}/sp-worker/conf -n vick-system
    kubectl create configmap sp-worker-bin --from-file=${download_location}/sp-worker/bin -n vick-system
    #Create SP worker deployment
    kubectl apply -f ${download_location}/vick-sp-worker-deployment.yaml -n vick-system
    kubectl apply -f ${download_location}/vick-sp-worker-service.yaml -n vick-system
    #Create SP dashboard configmaps
    kubectl create configmap sp-dashboard-conf --from-file=${download_location}/status-dashboard/conf -n vick-system
    #kubectl create configmap sp-worker-bin --from-file=sp-worker/bin -n vick-system
    #Create SP status dashboard deployment
    kubectl apply -f ${download_location}/vick-sp-dashboard-deployment.yaml -n vick-system
    kubectl apply -f ${download_location}/vick-sp-dashboard-service.yaml -n vick-system
    #Create SP dashboard ingress
    kubectl apply -f ${download_location}/vick-sp-dashboard-ingress.yaml -n vick-system
}
function init_control_plane () {
    local download_location=$1
    #Setup VICK namespace, create service account and the docker registry credentials
    kubectl apply -f ${download_location}/vick-ns-init.yaml

    HOST_NAME=$(hostname | tr '[:upper:]' '[:lower:]')
    #label the node
    kubectl label nodes $HOST_NAME disk=local

    #Create credentials for docker.wso2.com
    #kubectl create secret docker-registry wso2creds --docker-server=docker.wso2.com --docker-username=$DOCKER_REG_USER --docker-password=$DOCKER_REG_PASSWD --docker-email=$DOCKER_REG_USER_EMAIL -n vick-system
}

function deploy_istio () {
    local download_location=$1
    istio_version=$2
    ISTIO_HOME=${download_location}/istio-${istio_version}
    wget https://github.com/istio/istio/releases/download/${istio_version}/istio-${istio_version}-linux.tar.gz -P ${download_location}
    tar -xzvf ${download_location}/istio-${istio_version}-linux.tar.gz -C ${download_location}
    export PATH=$ISTIO_HOME/bin:$PATH
    kubectl apply -f $ISTIO_HOME/install/kubernetes/helm/istio/templates/crds.yaml
    #kubectl apply -f $ISTIO_HOME/install/kubernetes/istio-demo.yaml
    #kubectl apply -f $ISTIO_HOME/install/kubernetes/istio-demo-auth.yaml
    kubectl apply -f ${download_location}/istio-demo-vick.yaml
    kubectl wait deployment/istio-pilot --for condition=available --timeout=6000s -n istio-system
    #Enabling Istio injection
    kubectl label namespace default istio-injection=enabled
}

function deploy_vick_crds () {
    local download_location=$1
    #Install VICK crds
    kubectl apply -f ${download_location}/vick.yaml
}

function create_artifact_folder () {
 local tmp_folder=$1
 if [ -d $tmp_folder ]; then
        mv $tmp_folder ${tmp_folder}.$(date +%s)
    fi

    mkdir ${tmp_folder}
}
function download_vick_artifacts () {

    local base_url=$1
    local download_path=$2
    yaml_list=("$@")

    for file_path in "${yaml_list[@]}"
    do
      dir_name=""
      if [[ $file_path =~ / ]]; then
        dir_name=${file_path%/*}
      fi
      wget "$base_url/$file_path" -P "$download_path/$dir_name" -a vick-setup.log
    done
}

function install_nginx_ingress_kubeadm () {
    local download_location=$1
    #Install nginx-ingress for control plane ingress
    kubectl apply -f ${download_location}/mandatory.yaml
    kubectl apply -f ${download_location}/service-nodeport.yaml
}

function install_nginx_ingress_gcp () {
    local download_location=$1
    #Install nginx-ingress for control plane ingress
    kubectl apply -f ${download_location}/mandatory.yaml
}
#-----------------------------------------------------------------------------------------------------------------------
#Get the IaaS type form the user
iaas=$1
gcp_project=$2

#sanity check
#Bash 4 / gcloud tools if GCP

k8s_version="1.11.3-00"
istio_version="1.0.2"
download_path="tmp-wso2"
git_base_url="https://raw.githubusercontent.com/wso2/product-vick/master"

control_plane_base_url="${git_base_url}/system/control-plane/global"
control_plane_yaml=(
    "mysql-deployment.yaml"
    "mysql-persistent-volume-claim.yaml"
    "mysql-persistent-volumes-local.yaml"
    "mysql-persistent-volumes.yaml"
    "mysql-service.yaml"
    "nfs-deployment.yaml"
    "nfs-persistent-volume-claim.yaml"
    "nfs-persistent-volumes-local.yaml"
    "nfs-server-service.yaml"
    "vick-apim-gw-ingress.yaml"
    "vick-apim-gw.yaml"
    "vick-apim-persistent-volume-claim-local.yaml"
    "vick-apim-persistent-volume-claim.yaml"
    "vick-apim-persistent-volumes-local.yaml"
    "vick-apim-persistent-volumes.yaml"
    "vick-apim-pub-store-ingress.yaml"
    "vick-apim-pub-store.yaml"
    "vick-ns-init.yaml"
    "vick-sp-dashboard-deployment.yaml"
    "vick-sp-dashboard-ingress.yaml"
    "vick-sp-dashboard-service.yaml"
    "vick-sp-persistent-volumes.yaml"
    "vick-sp-worker-deployment.yaml"
    "vick-sp-worker-service.yaml"
    "vick-apim-artifacts-persistent-volumes.yaml"
    "vick-apim-artifacts-persistent-volume-claim.yaml"
    "mandatory.yaml"
    "service-nodeport.yaml"
    "apim-configs/gw/datasources/master-datasources.xml"
    "apim-configs/gw/user-mgt.xml"
    "apim-configs/gw/identity/identity.xml"
    "apim-configs/gw/tomcat/catalina-server.xml"
    "apim-configs/gw/carbon.xml"
    "apim-configs/gw/security/Owasp.CsrfGuard.Carbon.properties"
    "apim-configs/gw/registry.xml"
    "apim-configs/gw/resources/api_templates/velocity_template.xml"
    "apim-configs/gw/api-manager.xml"
    "apim-configs/gw/log4j.properties"
    "apim-configs/pub-store/datasources/master-datasources.xml"
    "apim-configs/pub-store/user-mgt.xml"
    "apim-configs/pub-store/identity/identity.xml"
    "apim-configs/pub-store/carbon.xml"
    "apim-configs/pub-store/registry.xml"
    "apim-configs/pub-store/resources/api_templates/velocity_template.xml"
    "apim-configs/pub-store/api-manager.xml"
    "apim-configs/pub-store/log4j.properties"
    "sp-worker/bin/carbon.sh"
    "sp-worker/siddhi/tracer-app.siddhi"
    "sp-worker/siddhi/telemetry-app.siddhi"
    "sp-worker/conf/deployment.yaml"
    "status-dashboard//conf/deployment.yaml"
    "mysql/dbscripts/init.sql"
)

crd_base_url="${git_base_url}/build/target"
crd_yaml=("vick.yaml")

istio_base_url="${git_base_url}/system/scripts/kubeadm"
istio_yaml=("istio-demo-vick.yaml")

declare -A config_params
declare -A nfs_config_params
#-----------------------------------------------------------------------------------------------------------------------
#Create temporary foldr to download vick artifacts
create_artifact_folder $download_path

echo "Downloading vick artifacts"

download_vick_artifacts $control_plane_base_url $download_path "${control_plane_yaml[@]}"

download_vick_artifacts $crd_base_url  $download_path "${crd_yaml[@]}"

download_vick_artifacts $istio_base_url $download_path "${istio_yaml[@]}"

#Install K8s
if [ $iaas == "GCP" ]; then
    echo "GCP selected"
    install_k8s_gcp $gcp_project
elif [ $iaas == "kubeadm" ]; then
     echo "kubeadm selected"
    install_k8s_kubeadm $k8s_version
    #configure master node
    configure_k8s_kubeadm
else
    echo "Installation script only supports GCP and Kubeadm"
fi

#Init control plane
echo "Creating vick-system namespace and the service account"

init_control_plane $download_path

#Deploy/Configure NFS APIM artifact
read -p "Do you want to deploy NFS server [Y/n]: " install_nfs < /dev/tty

if [ $install_nfs == "Y" ]; then
    if [ $iaas == "GCP" ]; then
        create_nfs_share_gcp "data"
    elif [ $iaas == "kubeadm" ]; then
        echo "Kubeadm based setup does not require a NFS server"
    fi
elif [ $install_nfs == "n" ]; then
    if [ $iaas != "kubeadm" ]; then
         read_nfs_connection
    fi
fi
update_apim_nfs_volumes $download_path

#Deploy/configure MySQL / APIM datasources
read -p "Do you want to deploy MySQL server in to vick-system namespace [Y/n]: " install_mysql < /dev/tty

if [ $install_mysql == "Y" ] && [ $iaas == "GCP" ]; then
    #Read db user / passwd
    read_control_plane_datasources_configs
    #Update the sql
    update_control_plance_sql $download_path
    deploy_mysql_server_gcp $download_path "vick-mysql-9"
elif [ $install_mysql == "Y" ] && [ $iaas == "kubeadm" ]; then
    read_control_plane_datasources_configs
    #update the sql file
    update_control_plance_sql $download_path
    deploy_mysql_server $download_path
else
    read_control_plane_datasources_configs
fi
update_control_plane_datasources $download_path

echo "Deploying the control plane API Manager"

deploy_global_gw $download_path $iaas

echo "Deploying SP"

deploy_sp_dashboard_worker $download_path $iaas

echo "Deploying Istio"

deploy_istio $download_path $istio_version

echo "Deploying vick crds"

deploy_vick_crds $download_path

echo "Deploying nginx-ingress"

if [ $iaas == "kubeadm" ]; then
    install_nginx_ingress_kubeadm $download_path
fi
#check GCP ingress

