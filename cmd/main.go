package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/radovskyb/watcher"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	version      = "0.0.0"            // App version number, set at build time with -ldflags "-X main.version=1.2.3"
	buildInfo    = "No build details" // Build details, set at build time with -ldflags "-X main.buildInfo='Foo bar'"
	isRestarting = false
)

func main() {
	log.Printf("### 🚀 PodKicker %s starting...", version)

	watchFsTarget := os.Getenv("PODKICKER_WATCH")
	if watchFsTarget == "" {
		log.Fatalln("### 💥 Env var PODKICKER_WATCH was not set")
	}
	targetName := os.Getenv("PODKICKER_TARGET_NAME")
	if targetName == "" {
		log.Fatalln("### 💥 Env var PODKICKER_TARGET_NAME was not set")
	}
	targetType := os.Getenv("PODKICKER_TARGET_TYPE")
	if targetType == "" {
		targetType = "deployment"
	}

	// Connect to Kubernetes API
	log.Printf("### 🌀 Attempting to use in cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("### 💻 Connecting to Kubernetes API, using host: %s", config.Host)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Fatal(err)
	}
	namespace := string(namespaceBytes)

	fileWatcher := watcher.New()

	go func() {
		for {
			select {
			case event := <-fileWatcher.Event:
				// Prevent multiple restarts if multiple file events happen at once
				if isRestarting {
					return
				}

				log.Printf("### ⛔ Detected file change: %v", event)

				// This patch is exactly how the `kubectl rollout restart` command works
				patchData := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))

				var err error
				if strings.ToLower(targetType) == "statefulset" {
					_, err = clientset.AppsV1().StatefulSets(namespace).Patch(context.Background(), targetName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
				} else {
					_, err = clientset.AppsV1().Deployments(namespace).Patch(context.Background(), targetName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
				}

				if err != nil {
					log.Printf("### 👎 Warning: Failed to patch %s, restart failed: %v", targetType, err)
				} else {
					isRestarting = true
					log.Printf("### ✅ Target %s, named %s was restarted!", targetType, targetName)
				}
			case err := <-fileWatcher.Error:
				log.Fatalln(err)
			case <-fileWatcher.Closed:
				return
			}
		}
	}()

	// Watch this folder or file for changes
	if err := fileWatcher.AddRecursive(watchFsTarget); err != nil {
		log.Fatalln(err)
	}

	log.Printf("### 👀 Watching '%s' for changes...", watchFsTarget)

	// Start the watching process - it'll check for changes every 200ms.
	if err := fileWatcher.Start(time.Millisecond * 200); err != nil {
		log.Fatalln(err)
	}
}
