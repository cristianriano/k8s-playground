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