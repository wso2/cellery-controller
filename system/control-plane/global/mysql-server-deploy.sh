#!/bin/bash

kubectl create configmap mysql-dbscripts --from-file=mysql/dbscripts/ -n vick-system

kubectl apply -f mysql-persistent-volumes-local.yaml -n vick-system

kubectl apply -f mysql-persistent-volume-claim.yaml -n vick-system

kubectl apply -f mysql-deployment.yaml -n vick-system

kubectl apply -f mysql-service.yaml -n vick-system