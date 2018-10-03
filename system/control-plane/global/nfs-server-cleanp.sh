#!/bin/bash

kubectl delete -f nfs-server-service.yaml -n vick-system

kubectl delete -f nfs-deployment.yaml -n vick-system

kubectl delete -f nfs-persistent-volume-claim.yaml -n vick-system

kubectl delete -f nfs-persistent-volumes-local.yaml -n vick-system
