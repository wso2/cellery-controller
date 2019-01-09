# How to install VICK using vick-setup-all-in-one.sh script.

VICK all in one setup script helps to deploy VICK into different providers such as vanilla k8s, GCP and OpenShift. User needs to pass the desired k8s provider to the script as 
given in the example.

If the user selects kubeadm as the k8s provider, installer script will deploy k8s on the Linux. When GCP k8s provider selected, the script creates a k8s cluster, MySQL instance and a NFS 
file share in the GCP platform. 


### Deploy VICK with kubeadm K8s provider.

```
curl <vick-all-in-one script url> | bash -s -- <k8s provider>

```

```
curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-setup-all-in-one.sh | bash -s -- kubeadm

```

### Deploy VICK into GCP K8s provider.

```
curl <vick-all-in-one script url> | bash -s -- <K8s provider> <GCP Project ID> <GCP Compute Zone> 
```

```
curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-setup-all-in-one.sh | bash -s -- GCP proj-vick us-west1-c
```

### Deploy VICK into vanilla K8s provider.
```
curl <vick-all-in-one script url> | bash
```

```
curl https://raw.githubusercontent.com/wso2/product-vick/master/system/scripts/all-in-one-installer/vick-setup-all-in-one.sh | bash
```