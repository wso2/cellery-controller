#!/usr/local/bin/bash

#create GCP project
#gcloud projects create my-knative-project --set-as-default

#Point GCP tothe new project
#gcloud config set project my-knative-project

#Enable required GCP APIs
#gcloud services enable \
#  cloudapis.googleapis.com \
#  container.googleapis.com \
#  containerregistry.googleapis.com

CLUSTER_NAME=vick-knative
CLUSTER_ZONE=us-west1-c

echo "Creating K8s cluster $CLUSTER_NAM in in zone $CLUSTER_ZONE"

#Create K8s cluster
gcloud container clusters create $CLUSTER_NAME \
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

#Install Istio
echo "Installing Istio"
#ISTIO_HOME=/Users/vickdev/wso2/VICK/vick-init/istio-1.0.2
#export PATH=$ISTIO_HOME/bin:$PATH
#kubectl apply -f $ISTIO_HOME/install/kubernetes/istio-demo-auth.yaml

#Wait for 30 seconds
sleep 30s

#Enabling Istio injection
kubectl label namespace default istio-injection=enabled

#Install Knative
echo "Installing Knative"
#kubectl apply -f https://github.com/knative/serving/releases/download/v0.1.1/release.yaml

#KNATIVE_HOME=/Users/vickdev/wso2/VICK/vick-init/knative-0.1.1
#kubectl apply -f $KNATIVE_HOME/release.yaml

#Install nginx-ingress for control plane ingress
kubectl apply -f mandatory.yaml

#Init VICK
echo "Installing VICK"
#kubectl apply -f $VICK_HOME/vic-init.yaml
