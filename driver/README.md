# Driver

Mounts and unmounts the volumes on the hosts itself.

## Installing

The volume driver needs to be installed on all nodes.

### Kubernetes

Kubernetes can install the driver through a *DaemonSet*, which can be created
by running this command:

```console
$ kubectl apply -f deploy/daemon.yaml
```

### OpenShift

Openshift will be supported after the proof of concept works.
