# Kubernetes

All the experiments are running on Minikube

## Installation

You need to install minikube to run this demo and a driver.

1. `brew install minikube` (This will also install `kubectl`. Formulae: kubernetes-cli)

2. Install the desired driver with brew (more available [here](https://minikube.sigs.k8s.io/docs/drivers/)):
  - `docker` (already available but doesn't support Ingress)
  - `virtualbox` (not working atm)
  - `hyperkit` (recommended)

3. Start the cluster with\
`minikube start --driver=<step_2>`

  - If there is an error use `--v=7 --alsologtostderr` flags to debug
  - Is possible to specify the max resources used by minikube as well when starting `--cpus 4 --memory 8192`

4. Enable Ingress to access the cluster from the local machine. Remember to use a compatible driver.
More info in [Ingress section](#ingress)\
`minikube addons enable ingress`

5. Install a basic nginx service with 2 pods `kubectl apply -f nginx-hello.yml`. Includes
  - Namespace nginx
  - Deployment and service
  - Ingress for the service
  - Network policy
  - Secrets for TLS (self signes more [here](#tls))

6. _(Optional)_ Use [Helm](#helm-charts) to install [Prometheus](#prometheus) and [Grafana](#grafana) for basic monitoring.

**NOTE:** The cluster is transient, not persisted. Once you deleted you have to create everything from scratch.

## Commands

The main way to interact with the cluster is with `kubectl`. But it can get cumbersome for some tasks
so is recommended to use **[K9 CLI](https://k9scli.io/)**.\
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
minikube start --driver=hyperkit
minikube addons enable ingress
```

Also make sure the DNS record is pointing to the IP where the Ingress Controller is running.

For testing edit `/etc/hosts` file (with sudo) and point the ip assigned to the Ingress to `nginx-example.com` (you can test it using ping).\
Check the IP by running\
`kubectl get ingress -n nginx --watch`

### TLS

To add TLS support to the Ingress controller first inject the SSL certificate and key, encoded base64 in a secret resource.\
And add the TLS config on the Ingress spec. (check the [manifest file](./nginx-hello.yml) where the secret is already defined)

For testing we can generate our self-signed certificates for the testing domain
(already provided in the repo using `nginx-example.com` domain)
```
openssl req -x509 -newkey rsa:4096 -keyout <domain>.key -out <domain>.cert -days 365 -sha256
base64 <domain>.key | pbcopy
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

Prepare the installation\
`kubectl apply -f prometheus-config.yml`
This will create a namespace for the installation and a network policy to allow prometheus to pull
metrics from the nginx containers

Install it with helm\
`helm install prometheus prometheus-community/prometheus --values prometheus-values.yml -n prometheus`

To check which values can be configured\
`helm show values prometheus-community/prometheus > values.yml`

## Grafana

Webapp for metrics and logs visualizations

Prepare the installation\
`kubectl apply -f grafana-config.yml`
This will create the grafana namespace for the chart and a network policy to allow communication _from_
and _to_ prometheus

Install it with helm\
`helm install grafana grafana/grafana --values grafana-values.yml -n grafana`

The `grafana-values.yml` configure:
- An ingress to access Grafana from `grafana.nginx-example.com` (remember to add it to `/etc/host` file too)
- Configure Prometheus as a datasource
- Install 2 dashboards from Grafana gallery and configure their datasources
- A sidecar container that scans for ConfigMaps with `grafana-dashboard` label to automatically add dashboards (lable can be customized)

Available values for configuration [here](https://artifacthub.io/packages/helm/grafana/grafana#configuration)

Once installed to get `admin` user password:\
`kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo`

## Jsonnet

Create and reuse json as code. To install
```
go install github.com/google/go-jsonnet/cmd/jsonnet
go install github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb@latest
```

To install the formater as well\
`go install github.com/google/go-jsonnet/cmd/jsonnetfmt@latest`

*Note: If using asdf run after installation `asdf reshim golang`*

To generate Grafana dashboards install [Grafonet 7.0](https://grafana.github.io/grafonnet-lib/getting-started/) and run jsonnet
```
jb install https://github.com/grafana/grafonnet-lib/grafonnet-7.0
# The following will generate nginx.json file with the dashboard config
jsonnet -J vendor/ -m . dashboards.jsonnet
kubectl create configmap grafana-dashboard-nginx-json --from-file=nginx.json --dry-run=client -o yaml -n grafana | yq e '.metadata.labels.grafana-dashboard = "true"' - | kubectl apply -f -
```