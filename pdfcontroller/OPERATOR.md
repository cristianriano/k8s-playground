# Operators

Operators are software extensions to Kubernetes that make use of custom resources to manage applications and their components.
More on the [official doc](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

## Start a new operator

1. Initialize go module `go mod init k8s.example.com/v1`
2. Run kube builder `kubebuilder init`
3. Create a new API `kubebuilder create api --version v1 --kind PdfDocument`
   1. Create both resource and controller
4. 