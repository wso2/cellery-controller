#!/bin/bash

kubectl delete configmap mysql-dbscripts -n vick-system

kubectl delete -f mysql-service.yaml -n vick-system

kubectl delete -f mysql-deployment.yaml -n vick-system

kubectl delete -f mysql-persistent-volume-claim.yaml -n vick-system

kubectl delete -f mysql-persistent-volumes-local.yaml -n vick-system
