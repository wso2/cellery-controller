## VICK Guide: Minikube

### Prerequisite

* Kubernetes v1.10.5 cluster with Istio 1.0.1 installed
 
<!--
###Starting a cluster
 
    minikube start --memory=8192 --cpus=4 --kubernetes-version=v1.10.5 --bootstrapper=kubeadm --extra-config=controller-manager.cluster-signing-cert-file="/var/lib/localkube/certs/ca.crt" --extra-config=controller-manager.cluster-signing-key-file="/var/lib/localkube/certs/ca.key" --extra-config=apiserver.admission-control="LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,MutatingAdmissionWebhook"
    
### Installing Istio

    curl -L https://raw.githubusercontent.com/knative/serving/master/third_party/istio-1.0.1/istio.yaml   | sed 's/LoadBalancer/NodePort/'   | kubectl apply -f -
    
* Wait until istio pods get started. Run following for check pod status
    
        kubectl get pods -n istio-system 
-->
        
### Installing VICK

    kubectl apply -f https://raw.githubusercontent.com/Mirage20/vick-test/master/artifacts/vick.yaml
    
* Wait until vick controller get started. Run the following command for watch the pod status

        watch kubectl get pods -n vick-system

### Creating a Cell

Apply the following yaml to create a sample cell

```yaml
apiVersion: vick.wso2.com/v1alpha1
kind: Cell
metadata:
  name: my-cell
spec:
  # Currently not applied
  gateway:
    http:
      port: 80
      routes:
      - match:
        - uri: /foo
        - uri: /bar
        service: server-time-us
```
Or use

    kubectl apply -f https://raw.githubusercontent.com/Mirage20/vick-test/master/artifacts/cell/example-cell.yaml

### Deploying a Service to the Cell

Apply the following yaml to deploy the service within the cell

```yaml
apiVersion: vick.wso2.com/v1alpha1
kind: Service
metadata:
  name: server-time-us
spec:
  cell: my-cell
  replicas: 2
  image: docker.io/mirage20/time-us
  containerPort: 8080
  servicePort: 80
```
Or use

    kubectl apply -f https://raw.githubusercontent.com/Mirage20/vick-test/master/artifacts/service/example-service.yaml

### Check the deployment using kubectl

Get cells

    kubectl get cells

Get services

    kubectl get vsvc
    
Verify the created pods

    kubectl get pods
    
### Access the service within the Cluster

SSH into a different pod using
    
    kubectl exec -it <pod-name> bash

and send a request
    
    curl server-time-us-service
    
### Deleting the cell 

NOTE: This will also delete all the services within that cell

    kubectl delete -f https://raw.githubusercontent.com/Mirage20/vick-test/master/artifacts/cell/example-cell.yaml

### Uninstall VICK

    kubectl delete -f https://raw.githubusercontent.com/Mirage20/vick-test/master/artifacts/vick.yaml
