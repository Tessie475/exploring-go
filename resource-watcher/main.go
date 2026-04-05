package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// CLI flags
	namespace := flag.String("namespace", "", "Namespace to watch (empty = all namespaces)")
	kubeconfig := flag.String("kubeconfig", defaultKubeconfig(), "Path to kubeconfig file")
	flag.Parse()

	// Create the Kubernetes client
	clientset, err := getClientset(*kubeconfig)
	if err != nil {
		fmt.Printf("✗ Failed to create client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│          Resource Watcher Started        │")
	fmt.Println("└─────────────────────────────────────────┘")
	if *namespace == "" {
		fmt.Println("  Watching: all namespaces")
	} else {
		fmt.Printf("  Watching: %s\n", *namespace)
	}
	fmt.Println()

	// Set up the informer factory.
	// The resync period is how often to re-list all resources.
	var factory informers.SharedInformerFactory
	if *namespace == "" {
		factory = informers.NewSharedInformerFactory(clientset, 30*time.Second)
	} else {
		factory = informers.NewSharedInformerFactoryWithOptions(
			clientset, 30*time.Second,
			informers.WithNamespace(*namespace),
		)
	}

	// Get the Pod informer — this watches for pod changes via the K8s API
	podInformer := factory.Core().V1().Pods().Informer()

	// Register event handlers for add, update, and delete events.
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onPodAdd,
		UpdateFunc: onPodUpdate,
		DeleteFunc: onPodDelete,
	})

	// Set up graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start the informer in the background
	factory.Start(stopCh)

	// Wait for the initial cache sync
	fmt.Println("  Syncing cache...")
	if !cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced) {
		fmt.Println("✗ Failed to sync cache")
		os.Exit(1)
	}
	fmt.Println("  ✓ Cache synced. Watching for changes...\n")

	// Block until we get a shutdown signal
	<-sigCh
	fmt.Println("\n  Shutting down...")
	close(stopCh)
}

func onPodAdd(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	fmt.Printf("  [+] Pod Created: %s/%s\n", pod.Namespace, pod.Name)

	// Check for :latest tag
	for _, c := range pod.Spec.Containers {
		if hasLatestTag(c.Image) {
			fmt.Printf("      ⚠ Warning: Container %q uses ':latest' tag (%s)\n", c.Name, c.Image)
		}
	}
}

func onPodUpdate(oldObj, newObj interface{}) {
	oldPod, ok := oldObj.(*corev1.Pod)
	if !ok {
		return
	}
	newPod, ok := newObj.(*corev1.Pod)
	if !ok {
		return
	}

	// Detect restarts
	for i, newStatus := range newPod.Status.ContainerStatuses {
		if i < len(oldPod.Status.ContainerStatuses) {
			oldStatus := oldPod.Status.ContainerStatuses[i]
			if newStatus.RestartCount > oldStatus.RestartCount {
				fmt.Printf("  [!] Restart Detected: %s/%s → container %q (restarts: %d)\n",
					newPod.Namespace, newPod.Name, newStatus.Name, newStatus.RestartCount)
			}
		}
	}

	// Detect phase changes
	if oldPod.Status.Phase != newPod.Status.Phase {
		fmt.Printf("  [~] Phase Changed: %s/%s → %s to %s\n",
			newPod.Namespace, newPod.Name, oldPod.Status.Phase, newPod.Status.Phase)
	}
}

func onPodDelete(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}
	fmt.Printf("  [-] Pod Deleted: %s/%s\n", pod.Namespace, pod.Name)
}

func getClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build config: %w", err)
		}
	}
	return kubernetes.NewForConfig(config)
}

func defaultKubeconfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}

func hasLatestTag(image string) bool {
	// If no colon, it defaults to latest
	lastIdx := -1
	for i := len(image) - 1; i >= 0; i-- {
		if image[i] == ':' {
			lastIdx = i
			break
		}
		if image[i] == '/' {
			break
		}
	}
	if lastIdx == -1 {
		return true
	}
	return image[lastIdx:] == ":latest"
}
