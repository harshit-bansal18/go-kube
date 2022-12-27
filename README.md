# go-kube
This repository contains Golang implementation of dynamic config injection in Kubernetes cluster. 

## Requirements

* kubectl
* minikube
* go v1.18+

## Compiling the Code
```
go mod tidy
go build -o ./app ./cmd/kube/
```

## Setup
Start the minikube virtual machine and set the env
```
minikube start
eval $(minikube docker-env)
```
This will create docker images inside the minikube vm.

## Create Docker image

Build the docker image to run inside the container.
```
docker build -t gokube:v1 .
```
Change the tag `v1` when updating the image. Also change the `image` field in `sample-pod.yaml` when deploying updated image. Remember all the images will be lost after deleting the vm.

## Start Kubernetes Cluster

If you have RBAC enabled on your cluster, use the following snippet to create role binding which will grant the default service account view permissions.
```
kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default
```
Use the following command to create configmap from `configmap/sample-configmap.yaml` file.
```
kubectl create -f configmap/sample-configmap.yaml
```
Use the following command to create a pod from `sample-pod.yaml` file present inside the repository.
```
kubectl apply -f sample-pod.yaml
```
Check the status of pod using
```
kubectl get pods
```
## Access the shell
Now get the shell to container to test the config client
```
kubectl exec --stdin --tty configmap-tester -- /bin/bash
```
Run the following from `/` to execute `app`.
```
./app
```
Open new shell with minikube docker env for next steps.

## Update the config
Modify existing fields or add new fields to `sample-configmap.yaml` under the `data` section.
\
Use the following command to upate the configmap object in kubernetes cluster.
```
kubectl replace -f configmap/sample-configmap.yaml
```
Now you should be able to see the new configmap data printed on the shell where `./app` was executed.

## Deleting Resources

```
kubectl delete pods configmap-tester
kubectl delete configmaps sample-config
```
To delete the vm after testing use
```
minikube delete
```

