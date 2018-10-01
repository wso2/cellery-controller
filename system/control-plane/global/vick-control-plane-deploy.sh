#!/usr/local/bin/bash

#Setup VICK namespace, create service account and the docker registry credentials
#kubectl apply -f vick-ns-init.yaml

#Create the pub-store config maps
kubectl create configmap apim-conf --from-file=apim-configs/pub-store -n vick-system
kubectl create configmap apim-conf-datasources --from-file=apim-configs/pub-store/datasources/ -n vick-system

#Create the gw config maps
kubectl create configmap gw-conf --from-file=apim-configs/gw -n vick-system
kubectl create configmap gw-conf-datasources --from-file=apim-configs/gw/datasources/ -n vick-system

#Create credentials for docker.wso2.com
kubectl create secret docker-registry wso2creds --docker-server=docker.wso2.com --docker-username=$DOCKER_REG_USER --docker-password=$DOCKER_REG_PASSWD --docker-email=$DOCKER_REG_USER_EMAIL -n vick-system

#Create volumes and volume claims
kubectl apply -f vick-apim-persistent-volumes.yaml -n vick-system
kubectl apply -f vick-apim-persistent-volume-claim.yaml -n vick-system

#Create pub-store deployment and the service
kubectl apply -f vick-apim-pub-store.yaml

#Create pub-store ingress
kubectl apply -f vick-apim-pub-store-ingress.yaml -n vick-system

#Create gateway deployment and the service
kubectl apply -f vick-apim-gw.yaml

#Create gateway ingress
kubectl apply -f vick-apim-gw-ingress.yaml -n vick-system
