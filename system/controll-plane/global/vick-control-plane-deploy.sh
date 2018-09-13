#!/usr/local/bin/bash
#Create ingress-nginx ns and required deployments

#Setup VICK name space, create service account and the docker registry credentials
kubectl apply -f vick-ns-init.yaml

#Create the pub-store config maps
kubectl create configmap apim-conf --from-file=apim-configs/ -n vick-system
kubectl create configmap apim-conf-datasources --from-file=apim-configs/datasources/ -n vick-system

#Create volumes

#Create pub-store deployment and the service
kubectl apply -f vick-apim-pub-store.yaml

#Create pub-store ingress
kubectl apply -f vick-apim-pub-store-ingress.yaml -n vick-system
