# kubectl-get-yaml

This plugin helps you to get rid of annoying `managedFields` field in object's yaml presentation.

Usage is quite simple:

`
kubectl get-yaml [flags] <object type> <object name>
`

Avaliable flags:

`-n/--namespace` - specify namespace

`-c/--context` - specify context

`-k/--kubeconfig` - specify kubeconfig file
