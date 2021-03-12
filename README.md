# Kubernetes

All the experiments are running on Minikube

## Installation

You need a vm driver so minikube can run the nodes.\
Can use `docker` if available or `hyperkit`

1. `brew install minikube` (This will also install `kubectl`. Formulae: kubernetes-cli)
2. Install the driver with brew as well
3. Start the cluster (if there is an error use `--v=7 --alsologtostderr` flags to debug)\
`minikube start --driver=docker`
4. You can specify the max resources used by minikube as well when starting `--cpus 4 --memory 8192`

## Commands

The main way to interact with the cluster is with `kubectl`. But it can get cumbersome for some tasks so we can use [K9 CLI](https://k9scli.io/).\
Install with `brew install derailed/k9s/k9s`

### Contexts
In case you have more than one cluster registered

List:\
`kubectl config get-contexts`

Use:\
`kubectl config use-context minikube`

Rename:\
`kubectl config rename-context old new`

### Get/Update Info

- Apply:
`kubectl apply -f service.yml`

- Delete:
`kubectl delete -f service.yml`

- Specify namespace with the flag `-n namespace`

- Get YAML:
`kubectl get deployment my-deployment -o yaml`

- More info:
`kubectl get pod -o wide`

- Create:
`kubectl create namespace my-namespace `

### Debugging

- When getting an error applying/validating the yaml you can get the description of the expected params like
`kubectl explain ingress.spec.rules.http.paths.backend.service`

## Ingress

First we need to install an ingress controller implementation. We would use K8S Nginx implementation which comes by default in Minikube.
This will create the `IngressController` pod in `kube-system` namespace. For this we need use a driver compatible with ingress.

```
minikube start --driver=virtualbox
minikube addons enable ingress
```

Also make sure the DNS record is pointing to the IP where the Ingress Controller is running.\
For testing edit `/etc/hosts` file with the ip assigned to the Ingress. Check the IP by running\
`kubectl get ingress -n nginx --watch`
