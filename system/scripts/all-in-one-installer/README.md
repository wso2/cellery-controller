# How to install VICK using vick-setup-all-in-one.sh script.

### Deploy VICK with kubeadm K8s provider.

```
cat vick-setup-all-in-one.sh |bash -s -- <K8s provider>
```

```
cat vick-setup-all-in-one.sh |bash -s -- kubeadm
```

### Deploy VICK with GCP K8s provider.

```
cat vick-setup-all-in-one.sh |bash -s -- <K8s provider> <GCP project name>
```

```
cat vick-setup-all-in-one.sh |bash -s -- GCP proj-vick
```
