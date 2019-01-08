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

#Deploy VICK with kubeadm K8s provider.
#cat vick-setup-all-in-one.sh | bash -s -- kubeadm

#Deploy VICK with GCP K8s provier.
#cat vick-setup-all-in-one.sh | bash -s -- GCP proj-vick

function install_k8s_kubeadm () {
local K8S_VERSION=$1

local UBUNTU_VERSION=$(cat /etc/lsb-release | sed -n '/DISTRIB_RELEASE/p' | awk -F'=' '{print $2}')

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
sudo apt-get update -qq
sudo apt-get install -y docker.io
#Install NFS client
sudo apt-get install -y nfs-common
#Enable docker service
sudo systemctl enable docker.service
#Update the apt source list
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] http://apt.kubernetes.io/ kubernetes-xenial main"
#Install K8s components
echo "⚙️ Installing kubeadm"
sudo apt-get install -y kubelet=$K8S_VERSION kubeadm=$K8S_VERSION kubectl=$K8S_VERSION
    sudo apt-mark hold kubelet kubeadm kubectl
}

function configure_k8s_kubeadm () {
local node_type=$1

if [ -z $node_type ]; then
    node_type="master"
fi

if [ $node_type == "master" ]; then
    #Initialize the k8s cluster
    echo "⚙️ Installing k8s using kubeadm"
    sudo kubeadm init --pod-network-cidr=10.244.0.0/16
    #Wait till k8s cluster initialization.
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
 echo "👣 k8s Master node installation is finished"

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
local nfs_share_location=$1
local gcp_compute_zone=$2
local nfs_server_ip

echo "ℹ️ Creating NFS share in GCP"
gcloud beta filestore instances create nfs-server \
    --location=$gcp_compute_zone \
    --tier=STANDARD \
    --file-share=name="${nfs_share_location}",capacity=1TB \
    --network=name="default"
#Wati till NFS server creation finished.
sleep 30

nfs_server_ip=$(gcloud beta filestore instances describe nfs-server --location=$gcp_compute_zone \
--format flattened | awk '/ipAddresses/  {print $2}')

if [ -n $nfs_server_ip ]; then
    nfs_config_params["NFS_SERVER_IP"]=$nfs_server_ip
    nfs_config_params["NFS_SHARE_LOCATION"]="/${nfs_share_location}"
else
    echo "🙅 Error in NFS share creation"
    exit 0
fi
}

#Read NFS connection details form the user
function read_nfs_connection () {
local nfs_server_ip
local nfs_share_location

echo "ℹ️ Configuring NFS volume"
echo
read -p "⛏️ NFS server IP: " nfs_server_ip < /dev/tty
if [[ ! -z "${nfs_server_ip/ //}" ]]; then
    nfs_config_params["NFS_SERVER_IP"]=$nfs_server_ip
fi
read -p "⛏️ NFS share location: " nfs_share_location < /dev/tty
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
    sed -i "s|$param|${nfs_config_params[$param]}|g" ${download_location}/vick-apim-artifacts-persistent-volumes.yaml
done
}

#Create K8s cluster in GCP
#User need to configure GCP CLI in the machine
function install_k8s_gcp () {
local gcp_project=$1
local gcp_compute_zone=$2

#Point GCP to the new project
gcloud config set project $gcp_project
gcloud config set compute/zone $gcp_compute_zone

#Enable required GCP APIs
gcloud services enable \
  cloudapis.googleapis.com \
  container.googleapis.com \
  containerregistry.googleapis.com

CLUSTER_NAME=cellery-runtime

echo "ℹ️ Creating K8s cluster $CLUSTER_NAM in in zone $gcp_compute_zone"

#Create K8s cluster
gcloud -q --verbosity=error container clusters create $CLUSTER_NAME \
  --cluster-version=latest \
  --machine-type=n1-standard-4 \
  --enable-autoscaling --min-nodes=1 --max-nodes=10 \
  --enable-autorepair \
  --scopes=service-control,service-management,compute-rw,storage-ro,cloud-platform,logging-write,monitoring-write,pubsub,datastore \
  --num-nodes=1

#Wait for 60 seconds to finish the k8s cluster creation.
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
local gcp_compute_zone=$3
local service_account
local mysql_server_ip

gcloud -q sql instances create ${sql_instance_name} --tier=db-n1-standard-1 --gce-zone=$gcp_compute_zone
service_account=$(gcloud beta sql instances describe ${sql_instance_name} --format flattened | awk '/serviceAccountEmailAddress/ {print $2}')
gcloud -q sql instances patch ${sql_instance_name} --authorized-networks=0.0.0.0/0 --database-flags character_set_server=latin1
#Wait till finish the MySQL server restart after the patching
sleep 60
gcloud sql users set-password root --instance=${sql_instance_name} --prompt-for-password --host=%
#Wait till the credential update sync.
sleep 30
cat tmp-wso2/mysql/dbscripts/init.sql | gcloud sql connect ${sql_instance_name} --user=root

mysql_server_ip=$(gcloud beta sql instances describe ${sql_instance_name}  --format flattened | awk '/.ipAddress/ {print $2}')
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
read -p "⛏️ Database user name: " db_user < /dev/tty
if [[ ! -z "${db_user/ //}" ]]; then
        config_params["DATABASE_USERNAME"]=$db_user
fi
read -s -p "⛏️ Database user password: " db_passwd < /dev/tty
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
done
}

function deploy_global_gw () {
local download_location=$1
local iaas=$2

if [ $iaas == "kubeadm" ] || [ $iaas == "k8s" ]; then
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

function deploy_global_pubstore () {
    local download_location=$1

#Create pubstore ingress
#pubstore ingress is pointed to gatway service. In the future pubstore ingress will be pointed to pubstore service.
kubectl apply -f ${download_location}/vick-apim-pub-store-ingress.yaml -n vick-system
}

function deploy_sp_dashboard_worker () {
local download_location=$1

#Create SP worker configmaps
kubectl create configmap sp-worker-siddhi --from-file=${download_location}/sp-worker/siddhi -n vick-system
kubectl create configmap sp-worker-conf --from-file=${download_location}/sp-worker/conf -n vick-system
#kubectl create configmap sp-worker-bin --from-file=${download_location}/sp-worker/bin -n vick-system
#Create SP worker deployment
kubectl apply -f ${download_location}/vick-sp-worker-deployment.yaml -n vick-system
kubectl apply -f ${download_location}/vick-sp-worker-service.yaml -n vick-system
#Create SP dashboard configmaps
#kubectl create configmap sp-dashboard-conf --from-file=${download_location}/status-dashboard/conf -n vick-system
#kubectl create configmap sp-worker-bin --from-file=sp-worker/bin -n vick-system
#Create vick dashboard deployment, service and ingress.
kubectl apply -f ${download_location}/vick-observability-portal.yaml -n vick-system
kubectl apply -f ${download_location}/vick-sp-worker-ingress.yaml -n vick-system
}

function init_control_plane () {
local download_location=$1

#Setup VICK namespace, create service account and the docker registry credentials
kubectl apply -f ${download_location}/vick-ns-init.yaml

HOST_NAME=$(hostname | tr '[:upper:]' '[:lower:]')
#label the node if k8s provider is kubeadm

kubectl label nodes $HOST_NAME disk=local

#Create credentials for docker.wso2.com
#kubectl create secret docker-registry wso2creds --docker-server=docker.wso2.com --docker-username=$DOCKER_REG_USER --docker-password=$DOCKER_REG_PASSWD --docker-email=$DOCKER_REG_USER_EMAIL -n vick-system
}

function deploy_istio () {
local download_location=$1
local istio_version=$2

ISTIO_HOME=${download_location}/istio-${istio_version}
wget https://github.com/istio/istio/releases/download/${istio_version}/istio-${istio_version}-linux.tar.gz -P ${download_location}
tar -xzf ${download_location}/istio-${istio_version}-linux.tar.gz -C ${download_location}
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
local yaml_list=("$@")

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
kubectl apply -f ${download_location}/cloud-generic.yaml
}
#-----------------------------------------------------------------------------------------------------------------------
#Get the IaaS type form the user
iaas=$1
gcp_project=$2
gcp_compute_zone=$3

#sanity check
#Bash 4 / gcloud tools if GCP

k8s_version="1.11.3-00"
istio_version="1.0.2"
download_path="tmp-wso2"
git_base_url="https://raw.githubusercontent.com/wso2/product-vick/master"

control_plane_base_url="${git_base_url}/system/control-plane/global/k8s-artifacts"
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
    "vick-observability-portal.yaml"
    "vick-sp-persistent-volumes.yaml"
    "vick-sp-worker-deployment.yaml"
    "vick-sp-worker-service.yaml"
    "vick-sp-worker-ingress.yaml"
    "vick-apim-artifacts-persistent-volumes.yaml"
    "vick-apim-artifacts-persistent-volume-claim.yaml"
    "mandatory.yaml"
    "service-nodeport.yaml"
    "cloud-generic.yaml"
    "apim-configs/gw/datasources/master-datasources.xml"
    "apim-configs/gw/user-mgt.xml"
    "apim-configs/gw/identity/identity.xml"
    "apim-configs/gw/tomcat/catalina-server.xml"
    "apim-configs/gw/carbon.xml"
    "apim-configs/gw/security/Owasp.CsrfGuard.Carbon.properties"
    "apim-configs/gw/registry.xml"
    "apim-configs/gw/resources/api_templates/velocity_template.xml"
    "apim-configs/gw/api-manager.xml"
    "apim-configs/gw/synapse-handlers.xml"
    "apim-configs/gw/log4j.properties"
    "apim-configs/pub-store/datasources/master-datasources.xml"
    "apim-configs/pub-store/user-mgt.xml"
    "apim-configs/pub-store/identity/identity.xml"
    "apim-configs/pub-store/carbon.xml"
    "apim-configs/pub-store/registry.xml"
    "apim-configs/pub-store/resources/api_templates/velocity_template.xml"
    "apim-configs/pub-store/api-manager.xml"
    "apim-configs/pub-store/log4j.properties"
    "sp-worker/siddhi/tracing-app.siddhi"
    "sp-worker/siddhi/istio-telemetry-app.siddhi"
    "sp-worker/siddhi/telemetry-app.siddhi"
    "sp-worker/conf/deployment.yaml"
    "mysql/dbscripts/init.sql"
)

crd_base_url="${git_base_url}/build/target"
crd_yaml=("vick.yaml")

istio_base_url="${git_base_url}/system/scripts/kubeadm"
istio_yaml=("istio-demo-vick.yaml")

declare -A config_params
declare -A nfs_config_params
#-----------------------------------------------------------------------------------------------------------------------


#Install K8s
if [[ -n ${iaas/[ ]*\n/} ]]; then
    if [ $iaas == "GCP" ]; then
        echo "ℹ️ Selected k8s provider: GCP"
        echo "ℹ️ GCP Project $gcp_project hosted in $gcp_compute_zone"
        if [ -n $gcp_project ]; then
            install_k8s_gcp $gcp_project $gcp_compute_zone
        else
            echo "GCP project name is required"
            exit 0
        fi
    elif [ $iaas == "kubeadm" ]; then
        echo "ℹ️ Selected k8s provider: kubeadm"
        install_k8s_kubeadm $k8s_version
        #configure master node
        configure_k8s_kubeadm
    else
        echo "Installation script supported k8s providers are GCP and Kubeadm."
        exit 0
    fi
else
    echo "Installing VICK into an existing K8s cluster"
fi

#Create temporary foldr to download vick artifacts
create_artifact_folder $download_path

echo "🕷️ Downloading vick artifacts to ${download_path}"

download_vick_artifacts $control_plane_base_url $download_path "${control_plane_yaml[@]}"
download_vick_artifacts $crd_base_url  $download_path "${crd_yaml[@]}"
download_vick_artifacts $istio_base_url $download_path "${istio_yaml[@]}"

#Init control plane
echo "🔧 Creating vick-system namespace and the service account"

init_control_plane $download_path

#Deploy/Configure NFS APIM artifact
 if [ $iaas == "GCP" ]; then
    read -p "⛏️ Do you want to deploy a NFS server [Y/n]: " install_nfs < /dev/tty

    if [ $install_nfs == "n" ]; then
         read_nfs_connection
    else
        create_nfs_share_gcp "data" $gcp_compute_zone
    fi
    update_apim_nfs_volumes $download_path
fi

#Deploy/configure MySQL / APIM datasources
read -p "⛏️ Do you want to deploy MySQL server [y/N]: " install_mysql < /dev/tty

if [ $install_mysql == "y" ]; then

    if [ $iaas == "GCP" ]; then
        #Read db user / passwd
        read_control_plane_datasources_configs
        #Update the sql
        update_control_plance_sql $download_path
        deploy_mysql_server_gcp $download_path "vick-mysql-$((1 + RANDOM % 1000))" $gcp_compute_zone
    elif [ $iaas == "kubeadm" ]; then
        read_control_plane_datasources_configs
        #update the sql file
        update_control_plance_sql $download_path
        deploy_mysql_server $download_path
    else
        echo "🔧 Deploy MySQL server into the existing K8s clusters"
        read_control_plane_datasources_configs
        update_control_plance_sql $download_path
        deploy_mysql_server $download_path
    fi
else
    read_control_plane_datasources_configs
fi

update_control_plane_datasources $download_path

echo "ℹ️ Start to Deploying the VICK control plane"
echo "🔧 Deploying the control plane API Manager"

deploy_global_gw $download_path $iaas
deploy_global_pubstore $download_path

echo "🔧Deploying Stream Processor"

deploy_sp_dashboard_worker $download_path $iaas

echo "🔧 Deploying Istio version $istio_version"

deploy_istio $download_path $istio_version

echo "🔧 Deploying VICK CRDs"

deploy_vick_crds $download_path

echo "🔧 Deploying ingress-nginx"

if [ $iaas == "kubeadm" ]; then
    install_nginx_ingress_kubeadm $download_path
elif [ $iaas == "GCP" ]; then
    install_nginx_ingress_gcp $download_path
fi

echo "ℹ️ VICK installation is finished."
echo "-=🎉=-"

