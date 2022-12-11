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

4. Install a basic nginx service with 2 pods that includes
  - Namespace nginx
  - Deployment and service
  - Ingress for the service
  - Network policy
  - Secrets for TLS (self signes more [here](#tls))
```
kubectl apply -f resources/11-nginx-hello.yml \
              -f resources/11-nginx-ingress.yml
```

5. Enable Ingress to access the cluster from the local machine. Remember to use a compatible driver.
More info in [Ingress section](#ingress)\
`minikube addons enable ingress`

6. _(Optional)_ Using [Traefik](#traefik) as a proxy to the nginx containers

7. _(Optional)_ Use [Helm](#helm-charts) to install [Prometheus](#prometheus) and [Grafana](#grafana) for basic monitoring.

**NOTE:** The cluster is transient, not persisted. Once you deleted you have to create everything from scratch.
If everything is configured you can start the cluster with basic config running the script `setup_example`.

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
And add the TLS config on the Ingress spec. (check the [manifest file](./resources/11-nginx-ingress.yml) where the secret is already defined)

For testing we can generate our self-signed certificates for the testing domain
(already provided in the repo using `nginx-example.com` domain)
```
openssl req -x509 -newkey rsa:4096 -keyout <domain>.key -out <domain>.cert -days 365 -sha256
base64 <domain>.key | pbcopy
```

## Traefik

[Traefik proxy](https://traefik.io/traefik/) is an OpenSource application proxy. Can run as a standalone
container or inside a Kubernetes cluster. Is configuration is dynamic allowing to add or remove services
as the become available or not.

In Kubernetes Traefik is notified when an Ingress resources is created/changed via the Kubernetes API
and updates the configuration.

First create the role for Traefik to have permissions and be notified about new Ingress

`kubectl apply -f resources/00-traefik-role.yml`

Then create the deployment and service (in the default namespace). The service will be of type LoadBalancer which means is up to the cluster
to expose it. Minikube has a [tunnel functionality](https://minikube.sigs.k8s.io/docs/handbook/accessing/#loadbalancer-access) to allow access.

```
kubectl apply -f resources/01-traefik-service.yml
minikube tunnel # start it in another terminal
```

To check the exposed IP run `kubectl get svc --watch`. Use the external-ip to reach the service

```
NAME                        TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)          AGE
kubernetes                  ClusterIP      10.96.0.1        <none>           443/TCP          4m15s
traefik-dashboard-service   LoadBalancer   10.103.204.156   10.103.204.156   8080:31422/TCP   65s
traefik-web-service         LoadBalancer   10.110.30.175    10.110.30.175    80:32648/TCP     65s
```

Access the dashboard in `10.103.204.156:8080/dashboard` and make sure the nginx service appears.
Once it does nginx should be reachable under `10.110.30.175/nginx`

**Remember that for Traefik to detect a new Ingress in another Namespace a NetworkPolicy should be in place for both incoming and outgoing traffic!!**

_For some reason pointing the DNS to Traefik doesn't work._

Configuration in Traefik can refer to two different things:

### Dynamic Configuration

Contains everything that defines how the requests are handled by your system (middleware, routers, etc).
This configuration can change and is seamlessly hot-reloaded, without any request interruption or connection loss.
It's obtain from providers: whether an orchestrator, a service registry, or a plain old configuration file.

It's defined in this [file](./resources/traefik_dynamic_config.yml)

Mount it with

`kubectl create configmap traefik-dynamic-config --from-file=resources/traefik_dynamic_config.yml -n default`

### Static Configuration

Also startup config. Sets up connections to providers and define the entrypoints Traefik will listen to (these elements don't change often).
There are three different, mutually exclusive ways to define static configuration options in Traefik:

1. In a configuration file (/etc/traefik/traefik.yml)
2. In the command-line arguments (available arguments [here](https://doc.traefik.io/traefik/reference/static-configuration/cli/))
3. As environment variables (available ENVS [here](https://doc.traefik.io/traefik/reference/static-configuration/env/))

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
`kubectl apply -f resources/12-prometheus-config.yml`
This will create a namespace for the installation and a network policy to allow prometheus to pull
metrics from the nginx containers

Install it with helm\
`helm install prometheus prometheus-community/prometheus --values prometheus-values.yml -n prometheus`

To check which values can be configured\
`helm show values prometheus-community/prometheus > values.yml`

## Grafana

Webapp for metrics and logs visualizations

Prepare the installation\
`kubectl apply -f resources/13-grafana-config.yml`
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

### Jsonnet

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