## Cell Gateway

  
### Running global API Manager in the cluster

If you do not having a running global API Manager, apply following yaml to deploy global APIM. 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: apim-deployment
  labels:
    app: apim
spec:
  replicas: 1
  selector:
    matchLabels:
      app: apim
  template:
    metadata:
      labels:
        app: apim
    spec:
      containers:
      - name: apim
        image: nipunaprashan/wso2am
---        
apiVersion: v1        
kind: Service
apiVersion: v1
metadata:
  name: apim-service
spec:
  selector:
    app: apim
  ports:
  - protocol: TCP
    name: publisher-store
    port: 443
    targetPort: 9443
  type: NodePort
```

### Deploy the sample API in global API Manager

1. Log in to the global API Manager using publicly accessible URL. (https://<apim-dns>/publisher)
2. If you are a first-time user,
   - Close the interactive tutorial that starts automatically.
   - Click the Deploy Sample API button. It deploys a sample API called PizzaShackAPI into the API Manager.
   
   If you have already deployed APIs in API publisher,
      - Create and deploy a sample API
  
### Generate a JWT token to invoke APIs deployed in the cell gateway
  
Once you start the Microgateway, you can use a JWT or an OAuth2 token to invoke the API.
1. Sign in to the API Store (https://\<apim-dns\>/store) with the admin/admin credentials.
2. Select the default application and an available tier, and click Subscribe.
3. Click Applications and choose to edit the Default Application. Set the token type to JWT.
4. Click the Production Keys tab and click regenerate Keys to generate an access token to invoke the API. 
5. copy and save the generated JWT token for future invocation.   


### Group APIs with Labels

1. Log in to the Admin portal ( https://\<apim-dns\>/admin). Use admin as the username and password. 
2. To add a new Microgateway label, click LABELS under MICROGATEWAY, and then click ADD MICROGATEWAY.
3. Create a new label (e.g. cell), add a host (e.g. https://cell-gateway) and click Save.
4. Navigate to the API Publisher (https://\<apim-dns\>/publisher). Sign in using admin as the username and password. 
5. Choose to edit the created sample API (Pizzashack API).
6. Navigate to the Manage section and click Gateway Environments. 
7. Select the newly created label `cell` to attach it to the sample API (Pizzashack API).
8. Click  Save & Publish
9. Similarly, you can select `cell` label for few other Published APIs.


### Deploy a cell gateway

Apply the following yaml to deploy the cell gateway. 

> Note: You can change the arguments in the following yaml file to match with the global gateway configurations.

`args`: [microgateway-project-name, API-label, global-API-manager-username, global-API-manager-password, 
global-api-manager-base-url, microgateway-trust-store-location, microgateway-trust-store-password]
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cell-gateway-deployment
  labels:
    app: cell-gateway
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cell-gateway
  template:
    metadata:
      labels:
        app: cell-gateway
    spec:
      containers:
      - name: cell-gateway
        image: nipunaprashan/microgateway
        volumeMounts:
        - name: targetdir
          mountPath: /target
      initContainers:
      - name: initialize-cell-gateway
        image: nipunaprashan/microgateway_init
        args: ["cell", "cell", "admin", "admin", "https://apim-service/", "lib/platform/bre/security/ballerinaTruststore.p12", "ballerina"]
        volumeMounts:
        - name: targetdir
          mountPath: "/target"
      volumes:
      - name: targetdir
        emptyDir: {}  
```
Or apply following file with above default arguments.

    kubectl apply -f https://raw.githubusercontent.com/wso2/product-vick/master/system/controll-plane/cell/gateway/artifacts/microgateway-deployment.yaml

### Deploy a Service to the cell gateway

```yaml
apiVersion: v1
kind: Service
apiVersion: v1
metadata:
  name: cell-gateway
spec:
  selector:
    app: cell-gateway
  ports:
  - protocol: TCP
    port: 443
    targetPort: 9095
```

Or use

    kubeclt apply -f https://raw.githubusercontent.com/wso2/product-vick/master/system/controll-plane/cell/gateway/artifacts/microgateway-service.yaml
    
### Invoke API deployed in the cell gateway

> If you want to invoke APIs deployed in microgateway directly using microgateway service you have to log into the 
cluster.

Access to one of the pod within the Cluster

    
    kubectl exec -it <pod-name> /bin/bash

then invoke the API
    
    curl -k -X GET "https://cell-gateway/hello/1.0.0/" -H "accept: application/json" -H "Authorization: Bearer <JWT-token>"
    
### Delete the cell gateway and service
    
    kubectl delete -f https://raw.githubusercontent.com/wso2/product-vick/master/system/controll-plane/cell/gateway/artifacts/microgateway-deployment.yaml
    kubeclt delete -f https://raw.githubusercontent.com/wso2/product-vick/master/system/controll-plane/cell/gateway/artifacts/microgateway-service.
    