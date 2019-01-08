#!/bin/bash

#kubectl apply -f vick-ns-init.yaml

#Delete the pub-store config maps
kubectl delete configmap apim-conf -n vick-system
kubectl delete configmap apim-conf-datasources -n vick-system

#Delete the pub-store config maps
kubectl delete configmap gw-conf -n vick-system
kubectl delete configmap gw-conf-datasources -n vick-system

#Delete credentials for docker.wso2.com
kubectl delete secret wso2creds -n vick-system

#Delete pub-store deployment and the service
kubectl delete -f vick-apim-pub-store.yaml -n vick-system

#Delete pub-store ingress
kubectl delete -f vick-apim-pub-store-ingress.yaml -n vick-system

#Delete gateway deployment and the service
kubectl delete -f vick-apim-gw.yaml -n vick-system

#Delete gateway ingress
kubectl delete -f vick-apim-gw-ingress.yaml -n vick-system

#Delete volumes and volume claims
kubectl delete -f vick-apim-persistent-volume-claim.yaml -n vick-system
kubectl delete -f vick-apim-persistent-volumes.yaml -n vick-system
