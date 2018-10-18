# Provisioner

Creates and deletes the volumes on Hetzner.

## Deployment

Create a secret containing the Hetzner project api token. The credentials will be used by the provisioner and also the driver.

**WARNING:**
Credentials will be passed on every node and persisted to disk for every volume.
Future releases will fix this issue.

```console
$ kubectl create secret generic hetzner-token \
--from-literal=api=s0m3l33tt0k3n
```

```console
$ kubectl create -f deploy/deployment.yaml
deployment "hetzner-provisioner" created
```
If you are not using RBAC or OpenShift you can continue to the usage section.

### Authorization

If your cluster has RBAC enabled or you are running OpenShift you must authorize the provisioner.
If you are in a namespace/project other than "hetzner-provisioner" either edit `deploy/auth/clusterrolebinding.yaml` or edit the `oadm policy` command accordingly.

#### RBAC
```console
$ kubectl create -f deploy/auth/serviceaccount.yaml
serviceaccount "hetzner-provisioner" created
$ kubectl create -f deploy/auth/clusterrole.yaml
clusterrole "hetzner-provisioner-runner" created
$ kubectl create -f deploy/auth/clusterrolebinding.yaml
clusterrolebinding "run-hetzner-provisioner" created
$ kubectl patch deployment hetzner-provisioner -p '{"spec":{"template":{"spec":{"serviceAccount":"hetzner-provisioner"}}}}'
```

#### OpenShift
Openshift will be supported after the proof of concept works.

## Usage

First a [`StorageClass`](https://kubernetes.io/docs/user-guide/persistent-volumes/#storageclasses) for claims to ask for needs to be created.

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: hetzner-cloud-storage
provisioner: stevenklar/hetzner-provisioner
```

### Parameters

Once you have finished configuring the class to have the name you chose when deploying the provisioner, create it.

```console
$ kubectl create -f deploy/class.yaml 
storageclass "hetzner-cloud-storage" created
```

When you create a claim that asks for the class, a volume will be automatically created.

```console
$ kubectl create -f deploy/claim.yaml 
persistentvolumeclaim "hetzner" created
$ kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM               STORAGECLASS        REASON    AGE
pvc-25de044e-5c93-11e8-885e-fe778d7189b9   1Mi        RWO            Delete           Bound     default/hetzner   hetzner-cloud-storage             1s
```
