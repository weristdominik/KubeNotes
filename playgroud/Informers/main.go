package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Load ~/.kube/config
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Shared informer factory
	factory := informers.NewSharedInformerFactory(
		clientset,
		30*time.Second, // resync period
	)

	// ConfigMap informer
	// Runs:
	// LIST: curl --cert client.crt --key client.key https://127.0.0.1:6443/api/v1/configmaps
	// WATCH: curl --cert client.crt --key client.key https://127.0.0.1:6443/api/v1/configmaps?watch=true
	configMapInformer := factory.Core().V1().ConfigMaps().Informer()

	// Secret informer
	// Runs:
	// LIST: curl --cert client.crt --key client.key https://127.0.0.1:6443/api/v1/secrets
	// WATCH: curl --cert client.crt --key client.key https://127.0.0.1:6443/api/v1/secrets?watch=true
	secretInformer := factory.Core().V1().Secrets().Informer()

	// Subscribe to events
	configMapInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("[ConfigMap ADD] %s/%s\n", cm.Namespace, cm.Name)
		},

		// UpdateFunc: func(oldObj, newObj interface{}) {
		// 	cm := newObj.(*corev1.ConfigMap)
		// 	fmt.Printf("[UPDATE] %s/%s\n", cm.Namespace, cm.Name)
		// },

		DeleteFunc: func(obj interface{}) {
			cm := obj.(*corev1.ConfigMap)
			fmt.Printf("[ConfigMap DELETE] %s/%s\n", cm.Namespace, cm.Name)
		},
	})

	secretInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s := obj.(*corev1.Secret)
			fmt.Printf("[Secret ADD] %s/%s\n", s.Namespace, s.Name)
		},
	})

	stopCh := make(chan struct{})

	// Start all informers in the factory
	factory.Start(stopCh)

	// Wait for cache sync
	if !cache.WaitForCacheSync(stopCh, configMapInformer.HasSynced, secretInformer.HasSynced) {
		panic("failed to sync cache")
	}

	// Wait for Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	fmt.Println("Shutting down...")
	close(stopCh)
}
