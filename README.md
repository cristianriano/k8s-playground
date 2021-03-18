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

Change Namespace for subsequent commands:\
`kubectl config set-context --current --namespace=prometheus`

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
- By default services can reach each other among namespaces, which for testing is great. Nevertheless there
are the valid NetworkPolicies defined
- Kube-DNS sets a url for each service by default like `http://service.namespace.svc.cluster.local`

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

### TLS

To add TLS support to the Ingress controller first inject the SSL certificate and key encoded base64 in a secret resource.\
And add the TLS config on the Ingress spec.\

For testing we can generate our self-signed certificates
```
openssl req -x509 -newkey rsa:4096 -keyout domain.key -out domain.cert -days 365 -sha256
base64 domain.key | pbcopy
```

## Helm Charts

Helm provides an interface to define, install, update and share pre-configured components for a K8S cluster as repositories.\
To get charts use the [Artifact Hub](https://artifacthub.io/)

1. Install helm locally `brew install helm`
2. Add the desired repo, search the chart you need and install (it will be installed in the cluster than
kubectl is pointing to)

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts

helm search repo prometheus-community/prometheus

helm install my-prometheus --values prometheus-values.yml prometheus-community/grafana
```

3. When installing provide a yml file to override the default values
4. You can list the current running charts with `helm ls --all-namespaces`
5. Also use helm for uninstalling (specify the namespace) `helm uninstall grafana --namespace grafana`
6. Update `helm upgrade grafana grafana/grafana --namespace grafana --values grafana-values.yml`

## Prometheus

In charge of collecting metrics and trigger alerts. It has 3 components:
- Storage (Time Series Database in Disk)
- Retrieval (pull)
- HTTP Servier

Install it with helm\
`helm install prometheus prometheus-community/prometheus --values prometheus-values.yml -n prometheus-namespace`

To check which values can be configured\
`helm show values prometheus-community/prometheus > values.yml`

## Grafana

Webapp for metrics and logs visualizations

Install it with helm\
`helm install grafana grafana/grafana --values grafana-values.yml -n grafana`

Available values for configuration [here](https://artifacthub.io/packages/helm/grafana/grafana#configuration)

Once installed to get `admin` user password:\
`kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo`

## Jsonnet

Create and reuse json as code. To install\
`go get github.com/google/go-jsonnet/cmd/jsonnet`

To install the formater as well\
`go get github.com/google/go-jsonnet/cmd/jsonnetfmt`

*Note: If using asdf run after installation `asdf reshim golang`*

To generate Grafana dashboards install grafana-builder and run jsonnet
```
jb install
jsonnet -J vendor/ -m . dashboards.jsonnet
```