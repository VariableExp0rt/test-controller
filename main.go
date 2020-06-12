//Hello from main

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Print("controller version")

	kubeconfig := ""
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")
	flag.Parse()
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	var (
		config *rest.Config
		err    error
	)
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v", err)
		os.Exit(1)
	}
	client := kubernetes.NewForConfigOrDie(config)

	sharedInformers := informers.NewSharedInformerFactory(client, 10*time.Minute)
	Controller := NewController(client, sharedInformers.Apps().V1().Deployments())

	sharedInformers.Start(nil)
	Controller.Run(nil)
}
