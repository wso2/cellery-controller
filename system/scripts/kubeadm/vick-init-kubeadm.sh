#!/bin/bash


#Install Istio
echo "Installing Istio"

wget https://github.com/istio/istio/releases/download/1.0.2/istio-1.0.2-linux.tar.gz
tar -xzvf istio-1.0.2-linux.tar.gz

ISTIO_HOME=istio-1.0.2
export PATH=$ISTIO_HOME/bin:$PATH
#kubectl apply -f $ISTIO_HOME/install/kubernetes/istio-demo-auth.yaml
kubectl apply -f $ISTIO_HOME/install/kubernetes/istio-demo.yaml

#Wait for 30 seconds
sleep 30s

#Enabling Istio injection
kubectl label namespace default istio-injection=enabled

#Install Knative
#echo "Installing Knative"
#kubectl apply -f https://github.com/knative/serving/releases/download/v0.1.1/release.yaml

#Install nginx-ingress for control plane ingress
kubectl apply -f mandatory.yaml
kubectl apply -f service-nodeport.yaml

#Init VICK
#echo "Installing VICK"
#kubectl apply -f $VICK_HOME/vic-init.yaml
