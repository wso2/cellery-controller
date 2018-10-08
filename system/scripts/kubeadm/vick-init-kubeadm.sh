#!/bin/bash


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
