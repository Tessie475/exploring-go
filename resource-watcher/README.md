# Resource Watcher

A standalone Kubernetes controller built with Go that uses `client-go` informers to watch Pods in real-time. It detects and logs lifecycle events without requiring Kubebuilder scaffolding.

## What It Does

When running, the watcher connects to a Kubernetes cluster and monitors Pods. It looks for specific patterns or issues, such as:
- **Crash Detection:** Detects when a pod's container restarts and logs the new restart count.
- **Phase Changes:** Triggers alerts when a pod moves between states (e.g., from `Pending` to `Running` or `Failed`).
- **Antipattern Warnings:** Identifies and warns when a new pod specifies an image with the `:latest` tag (or defaults to it).
- **CrashLoopBackOff:** Specifically alerts when a pod enters a crash loop.

## Usage

### Prerequisites
- Go 1.21+
- A running Kubernetes cluster (like `kind` or `minikube`)
- Local `~/.kube/config` set up (or run strictly inside the cluster)

### Run Locally

You can run the watcher directly against your active cluster context:

```bash
go run main.go
```

To restrict the watcher to a specific namespace:

```bash
go run main.go --namespace=my-namespace
```

### Test Actions

While the watcher is running, create or modify pods in another terminal to see events handled in real-time.

```bash
# Triggers an antipattern warning for the :latest tag
kubectl run test-nginx --image=nginx:latest

# Trigger a phase change and deletion event
kubectl delete pod test-nginx

# Trigger crash detection / CrashLoopBackOff warnings
kubectl run crasher --image=busybox -- /bin/false
```

## How It Works

This tool is built directly on top of `client-go`, specifically utilizing the `informers` and `cache` packages.

1. Sets up a `kubernetes.Clientset` using out-of-cluster (`~/.kube/config`) or in-cluster configuration.
2. Initializes a `SharedInformerFactory` to efficiently query and cache resources.
3. Retrieves the Pod informer and registers event handlers (`AddFunc`, `UpdateFunc`, `DeleteFunc`).
4. Handles `os/signal` events to ensure a graceful shutdown instead of abruptly killing the connection.

## Built With

- [Go](https://go.dev/)
- [client-go](https://github.com/kubernetes/client-go) (Kubernetes client library for Go)
