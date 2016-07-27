package controller

import (
	"time"

	"github.com/sjenning/kube-led/pin"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/pkg/watch"
)

type LEDController struct {
	client     *unversioned.Client
	controller *framework.Controller
	store      cache.Store
	pin        pin.Pin
	numPods    int
}

func NewLEDController(client *unversioned.Client, resyncPeriod time.Duration, namespace string, pin pin.Pin) (*LEDController, error) {
	c := LEDController{
		client:  client,
		pin:     pin,
		numPods: 0,
	}

	handlers := framework.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.numPods++
		},
		DeleteFunc: func(obj interface{}) {
			c.numPods--
		},
	}
	c.store, c.controller = framework.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options api.ListOptions) (runtime.Object, error) {
				return c.client.Pods(namespace).List(options)
			},
			WatchFunc: func(options api.ListOptions) (watch.Interface, error) {
				return c.client.Pods(namespace).Watch(options)
			},
		},
		&api.Pod{}, resyncPeriod, handlers)

	return &c, nil
}

func (c *LEDController) blinkLED() {
	for i := 0; i < c.numPods; i++ {
		c.pin.On()
		time.Sleep(100 * time.Millisecond)
		c.pin.Off()
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *LEDController) Run(stopCh chan struct{}) {
	go c.controller.Run(stopCh)
	c.pin.Off()
	wait.Until(c.blinkLED, 2*time.Second, stopCh)
}
