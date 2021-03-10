# Kubernetes

All the experiments are running on Minikube

## Installation

You need a vm driver so minikube can run the nodes.\
Can use `docker` if available or `hyperkit`

1. `brew install minikube` (This will also install `kubectl`. Formulae: kubernetes-cli)
2. Install the driver with brew as well
3. Start the cluster (if there is an error use `--v=7 --alsologtostderr` flags to debug)\
`minikube start --driver=docker`

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

Apply:\
`kubectl apply -f service.yml`

Delete:\
`kubectl delete -f service.yml`

For the following you should specify the namespace with `-n namespace`

Get YAML:\
`kubectl get deployment my-deployment -o yaml`

More info:\
`kubectl get pod -o wide`