# AppService Operator

A Kubernetes operator built with Go and Kubebuilder that manages a custom `AppService` resource. When you create an AppService, the operator automatically creates and manages a Deployment and Service for it.

## What It Does

```
You apply:                     Operator creates:
┌─────────────────┐           ┌─────────────────┐
│   AppService    │    ───►   │   Deployment    │
│  image: nginx   │           │   (N replicas)  │
│  replicas: 2    │           └─────────────────┘
│  port: 80       │           ┌─────────────────┐
└─────────────────┘    ───►   │    Service      │
                              │   (port 80)     │
                              └─────────────────┘
```

## Prerequisites

- Go 1.21+
- Docker
- kubectl
- A Kubernetes cluster ([kind](https://kind.sigs.k8s.io/), minikube, or Docker Desktop)
- [Kubebuilder](https://book.kubebuilder.io/)

## Usage

### Run Locally (Development)

```bash
# Install the CRD into your cluster
make install

# Run the operator locally
make run
```

### Deploy to Cluster (Production-style)

```bash
# Build the Docker image
make docker-build IMG=appservice-operator:latest

# Load into kind cluster
kind load docker-image appservice-operator:latest

# Deploy to cluster
make deploy IMG=appservice-operator:latest

# Verify
kubectl get pods -n appservice-operator-system
```

### Create an AppService

```yaml
apiVersion: apps.example.com/v1alpha1
kind: AppService
metadata:
  name: my-app
spec:
  image: nginx:latest
  replicas: 2
  port: 80
```

```bash
kubectl apply -f config/samples/apps_v1alpha1_appservice.yaml

# Verify resources were created
kubectl get appservices
kubectl get deployments
kubectl get services
kubectl get pods
```

### Test Scenarios

```bash
# Scale — change replicas and re-apply
kubectl edit appservice my-app

# Delete — Deployment and Service auto-delete via owner references
kubectl delete appservice my-app
```

## AppService Spec

| Field | Type | Description |
|---|---|---|
| `image` | string | Container image to deploy |
| `replicas` | int32 | Number of pod replicas (default: 1) |
| `port` | int32 | Container port to expose |

## AppService Status

| Field | Type | Description |
|---|---|---|
| `availableReplicas` | int32 | Number of pods currently running |
| `conditions` | []Condition | Current state of the resource |

## Admission Webhooks

The operator includes both mutating and validating webhooks that intercept AppService resources before they're saved.

> **Note:** Webhooks only work when deployed to the cluster (not with `make run`), because they require TLS certificates from cert-manager.

### Mutating Webhook (Defaults)

Automatically modifies resources before saving:
- Sets `replicas` to 2 if not specified
- Adds `managed-by: appservice-operator` label if missing

### Validating Webhook (Rules)

| Rule | Result |
|---|---|
| `image` is empty | ✗ Rejected |
| `replicas` outside 1-10 | ✗ Rejected |
| `port` outside 1-65535 | ✗ Rejected |
| Image uses `:latest` tag | ⚠️ Warning (allowed) |

### Testing Webhooks

```bash
# Should be REJECTED (empty image):
kubectl apply -f - <<EOF
apiVersion: apps.example.com/v1alpha1
kind: AppService
metadata:
  name: bad-app
spec:
  image: ""
  replicas: 2
  port: 80
EOF

# Should SUCCEED with defaults applied:
kubectl apply -f - <<EOF
apiVersion: apps.example.com/v1alpha1
kind: AppService
metadata:
  name: good-app
spec:
  image: nginx:latest
  port: 80
EOF

# Check that replicas was defaulted to 2:
kubectl get appservice good-app -o yaml | grep replicas
```

## How It Works

1. User creates an AppService custom resource
2. **Mutating webhook** sets defaults (replicas, labels)
3. **Validating webhook** checks rules (image required, replicas 1-10, valid port)
4. The controller detects the new resource via the reconciliation loop
5. It creates a Deployment with the specified image, replicas, and port
6. It creates a Service that routes traffic to the pods
7. It updates the AppService status with the actual replica count
8. If the AppService is updated, the controller updates the Deployment/Service to match
9. If the AppService is deleted, owner references cause automatic cleanup

## Project Structure

```
appservice-operator/
├── api/v1alpha1/
│   └── appservice_types.go              # CRD type definitions (Spec, Status)
├── internal/
│   ├── controller/
│   │   └── appservice_controller.go     # Reconcile logic
│   └── webhook/v1alpha1/
│       └── appservice_webhook.go        # Mutating + validating webhooks
├── config/
│   ├── crd/                             # Generated CRD manifests
│   ├── webhook/                         # Webhook server configuration
│   ├── certmanager/                     # TLS certificate config
│   └── samples/                         # Sample AppService YAML
├── cmd/
│   └── main.go                          # Entry point
└── Makefile                             # Build, install, deploy commands
```

## Built With

- [Go](https://go.dev/)
- [Kubebuilder](https://book.kubebuilder.io/)
- [controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [cert-manager](https://cert-manager.io/) (TLS for webhooks)
- [kind](https://kind.sigs.k8s.io/) (local Kubernetes)
