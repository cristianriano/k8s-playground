# Operators

Operators are software extensions to Kubernetes that make use of custom resources to manage applications and their components.
More on the [official doc](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

## Start a new operator

1. Initialize go module `go mod init k8s.example.com/v1`
2. Run kube builder `kubebuilder init --domain <domain>`
   1. Default domain is `my.domain`
3. Create a new API `kubebuilder create api --version v1 --kind <resource>`
   1. Choose to create both resource and controller
   2. Use CamelCase for the resource name like `PdfDocument`
4. Add fields to the Spec of the new resource in `api/v1/<resource>_types.go` in the `<resource>Spec` struct
5. Implement the reconciliation business logic in `internal/controllers/<resource>_controller.go`
6. When updating the Spec of the CRD update the generated yml with `make manifests`
7. Install the crds in the cluster `make install`
8. Last run the operator `make run`
9. Can be testing by creating a resource of the defined spec and checking resulting jobs and containers