# VICK API spec

## YAML Definitions

### Cell

```yaml
apiVersion: vick.wso2.com/v1alpha1
kind: Cell
metadata:
  # Name of the cell which is used by the services to specify the owner
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


### Service

```yaml
apiVersion: vick.wso2.com/v1alpha1
kind: Service
metadata:
  # Name of the service. Other services can use this name to invoke 
  name: server-time-us
spec:
  # Name of the owner cell
  cell: my-cell
  # Number of replicas that this service should have
  replicas: 2
  # Container image of the service 
  image: docker.io/mirage20/time-us
  # Port that container exposes
  containerPort: 8080
  # Port that service exposes so that other service can access.
  servicePort: 80
```
