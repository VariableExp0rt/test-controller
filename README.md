# Testing Controllers (and having fun)

This repository is to document some of the progress I have made into having a deployment controller monitor annotations and labels, and have a reconcile function to do *something* when they are absent.

From the project directory in your $GOPATH

`go run *.go -kubeconfig=<path/to/.kube/config`

You may also build this from a Dockerfile and deploy on the cluster, you shouldn't need the supply the *-kubeconfig* flag that way. Subsequently, use the `kubectl logs <pod-controller-is-running-on>` to identify events controller is watching.
