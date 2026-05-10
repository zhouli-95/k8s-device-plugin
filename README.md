# K8S Device Plugin

This is a device plugin for kubernetes that enables resource registration through configuration files. It is designed for learning the basic mechanisms of device plugins and for registering virtual resources for validation purposes.

# How to use device plugin

1. Deploy Device Plugin

```bash
kubectl apply -k deployments/kustomize/base
```

Check whether the Pod is in Running state.

```bash
kubectl get po -n device-plugin
```

Example Output:

```
NAME                  READY   STATUS    RESTARTS   AGE
device-plugin-82jcs   1/1     Running   0          8s
device-plugin-gscj2   1/1     Running   0          8s
```

2. Check Node resource

```bash
kubectl describe no <node-name>
```

Example Output:

```
Capacity:
  cpu:                12
  ephemeral-storage:  1055762868Ki
  example.com/dev:    4
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  memory:             16286916Ki
  pods:               110
Allocatable:
  cpu:                12
  ephemeral-storage:  1055762868Ki
  example.com/dev:    4
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  memory:             16286916Ki
  pods:               110
...
```

3. Create test pod which requests 1 device.

```bash
kubectl create -f examples/dev-1-example.yaml
```

Check device in pod.

```bash
kubectl exec -it dev-pod-1 -- ls -l /dev/
```

Example Output:

```
crw-rw-rw- 1 root root 1, 3 May  8 14:53 compute-device001
crw-rw-rw- 1 root root 1, 3 May  8 14:53 control-device
```

4. Uninstall Device Plugin

```bash
kubectl delete -k deployments/kustomize/base
```

# How to develop

Build binary.

```bash
make build
```

Build docker image.

```bash
make docker
```