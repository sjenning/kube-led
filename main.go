package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"

	"github.com/sjenning/kube-led/controller"
	"github.com/sjenning/kube-led/pin/minnowboard"
)

var (
	namespace = flag.String("namespace", "", "Namespace to watch (defaults to all namespaces)")
)

func main() {
	flag.Parse()

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		fmt.Printf("failed loading client config")
		return
	}
	client := unversioned.NewOrDie(config)

	pin, err := minnowboard.NewPin(minnowboard.LED_D2)
	if err != nil {
		fmt.Println(err)
		return
	}

	controller, err := controller.NewLEDController(client, 30*time.Second, *namespace, pin)
	if err != nil {
		fmt.Println("failed to create controller")
		return
	}
	stopCh := make(chan struct{})
	go controller.Run(stopCh)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(stopCh)
	pin.Close()
}
