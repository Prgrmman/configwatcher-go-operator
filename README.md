# Operator
This project show cases how to implement a golang operator. It uses operator-sdk to initialize and build the project. The operator logic is based on a use case of application redeploying automatically when the configmap, which is mounted on the application, is updated. This is a common scenario that many applications encounter whereby a change in configmap does not automatically start an application and some manual intervention is required. This operator automate this common chore.

## Pre-requisites
In order to run deploy and test this operator locally, the following tools are requried:
* openshift local installation
      
  https://access.redhat.com/documentation/en-us/red_hat_openshift_local/2.17/html/getting_started_guide/index]
* Operator-sdk for golang
  
  https://sdk.operatorframework.io/docs/installation/]
  
* Podman is installed and the classpath is configured for it:

  https://podman.io/getting-started/installation]

## Testing it on local OCP or K8 cluster
You need a local OCP cluster in order to deploy and run the operator and test it with an application. These intructions assume that there is an OCP cluster running locally.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `oc cluster-info` shows).

### Running on the cluster
1. Install the CRD (custom resource definition)
   From the root of the project directory:
```sh
  make install
```

2. Build and push your image to the location specified by `IMG`. In order to push to the registry, you need to be logged on the registry with podman

```sh
make docker-build docker-push IMG=<some-registry>/configwatcher-go-operator:tag
```

3. In order to make sure that the operator is deployed to a specific namespace, and that namespace has access to the registry in step 2) above, the following changes need to be made:
   
   - Create a project in the OCP
      ```
        oc new-project <project-name>
      ```
   - Change the **config/default/kustomization.yaml** file to edit the project name to match the name of your project
      ```
        # Adds namespace to all resources.
         namespace: <project-name>

      ```
    - Create a secret in the OCP project to access the registry:

       ```
       oc create secret docker-registry my-secret --docker-server=quay.io --docker-username=<u-name> --docker-password=<password> --docker-email=<email>
       ```
    - Change the rbac/service_account.yaml file with the following:
```
apiVersion: v1
imagePullSecrets:
- name: my-secret
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: serviceaccount
    app.kubernetes.io/instance: controller-manager
    app.kubernetes.io/component: rbac
    app.kubernetes.io/created-by: memcached-operator
    app.kubernetes.io/part-of: memcached-operator
    app.kubernetes.io/managed-by: kustomize
  name: controller-manager
  namespace: system
secrets:
- name: my-secret
```
      
      

4. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/configwatcher-go-operator:tag
```
**note:** ake sure you are on the targeted OCP project created in the previous step

5. Deploy the application on the OCP. This is the sample application having a configmap mounted that would be updated in the later step to test the workings of the operator:

```
oc apply -f extra/web-app.yaml
```

6. Deploy the CR (custom resource) on the cluster:

```
 apply -f config/samples/tutorials_v1_configwatcher.yaml
```
7. Create a route for the application

```
 oc expose svc webapp
```
8. Access the application using the route

```
curl http://$(oc get route -o jsonpath='{.items[0].spec.host}')
```
9. Update the configmap using this command:

```
oc  patch configmap webapp-config -p '{"data":{"message":"Greets from your smooth operator!"}}'
```

10. Verify that your application message is updated by running the step 8)

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy IMG=<the image in the registry>
```



### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster

