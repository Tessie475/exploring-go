# Exploring Go

This repository is a collection of small Go programs I've built while learning the language more deeply through hands-on problem solving.

Rather than treating Go as a purely academic exercise, I'm using this space to explore how it can be applied to real engineering tasks such as networking, concurrency, HTTP communication, and command-line tooling.

A lot of these projects are intentionally small, but each one helps me understand an important Go concept by building something practical.

## Why this repository exists

I created this repository to:
- strengthen my understanding of core Go concepts
- build small tools instead of only reading syntax examples
- document my progress in a way that is practical and portfolio-friendly
- explore Go from a systems and infrastructure perspective

## Projects

### Uptime Checker
A concurrent CLI tool that checks whether a list of URLs is reachable, resolves their IP addresses, and reports HTTP status codes.

**Concepts explored:**
- goroutines
- sync.WaitGroup
- HTTP clients
- timeouts
- DNS lookups
- URL parsing

### Config File Validator
A CLI utility that reads and validates YAML/JSON configuration files before deployment, catching syntax errors and invalid values early.

**Concepts explored:**
- file I/O (`os.ReadFile`)
- JSON parsing (`encoding/json`)
- YAML parsing (`gopkg.in/yaml.v3`)
- struct tags for field mapping
- nested structs
- custom validation logic
- error wrapping (`fmt.Errorf` with `%w`)
- working with third-party packages

### System Metrics API
An HTTP server that exposes real-time system metrics (CPU, memory, disk) as JSON endpoints, with a live HTML dashboard that auto-refreshes.

**Concepts explored:**
- HTTP server (`net/http`, `http.HandleFunc`, `http.ListenAndServe`)
- JSON API responses (`json.NewEncoder`)
- HTML templating (`html/template`)
- system metrics collection (`gopsutil`)
- struct composition
- route registration

### AppService Operator
A Kubernetes operator that manages a custom `AppService` resource. When an AppService is created, the operator automatically provisions a Deployment and Service for it.

**Concepts explored:**
- Custom Resource Definitions (CRDs)
- controller reconciliation loop (watch → detect → act)
- Kubebuilder scaffolding and code generation
- `controller-runtime` and `client-go`
- owner references for automatic cleanup
- spec vs status pattern
- RBAC markers
- admission webhooks (mutating and validating)
- TLS certificates with cert-manager
- deploying operators to a cluster

### Resource Watcher
A standalone Kubernetes controller that uses `client-go` informers to watch Pods in real-time, detecting crashes, bad image tags, and state changes.

**Concepts explored:**
- `client-go` directly (without Kubebuilder)
- Kubernetes informers and watch mechanisms
- Event-driven programming with work queues
- Kubeconfig loading and authentication
- Graceful shutdown and signal handling

## What I'm focusing on
Most of the projects in this repository are centered around:
- concurrency
- networking
- CLI applications
- Kubernetes development
- backend/system-style problem solving

## Notes
This is an evolving repository. Some projects may remain intentionally simple, while others may be revisited and improved as my understanding grows.

## Inspiration

Project ideas in this repository were inspired by the excellent 
[Golang for DevOps](https://github.com/techiescamp/golang-for-devops) repository.

I used it as a guide for identifying Go exercises while 
implementing the solutions myself to better understand the language and its 
standard library.