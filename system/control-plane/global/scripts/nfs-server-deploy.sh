#!/bin/bash

kubectl create -f nfs-persistent-volumes-local.yaml -n vick-system

kubectl apply -f nfs-persistent-volume-claim.yaml -n vick-system

kubectl apply -f nfs-deployment.yaml -n vick-system

kubectl apply -f nfs-server-service.yaml -n vick-system