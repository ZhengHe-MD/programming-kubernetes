# cnat-client-go

## Prerequisites

It makes use of the generators in [k8s.io/code-generator](https://github.com/kubernetes/code-generator)
to generate a typed client, informers, listers and deep-copy functions. Please make sure you have projects
laid out in the following way:

```base
# the path to code-generator
.../Workspace/github.com/kubernetes/code-generator
# the path to this demo
.../Workspace/github.com/ZhengHe-MD/programming-kubernetes/cnat-client-go
```

## Development

```bash
# build custom controller binary:
$ go build -o cnat-controller .

# launch custom controller locally:
$ ./cnat-controller -kubeconfig=$HOME/.kube/config

# after API/CRD changes:
$ ./hack/update-codegen.sh

# register At custom resource definition:
$ kubectl apply -f artifacts/examples/cnat-crd.yaml

# create an At custom resource:
$ kubectl apply -f artifacts/examples/cnat-example.yaml
```