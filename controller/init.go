package controller

import (
	"flag"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/klog-controller/config"
	informers "k8s.io/klog-controller/pkg/generated/informers/externalversions"
	"k8s.io/klog-controller/pkg/signals"
	"k8s.io/klog/v2"
	"time"
)

var CrdController *Controller

func RunCrd() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(config.KubeClient, time.Second*10)
	logInformerFactory := informers.NewSharedInformerFactory(config.KlogClient, time.Second*10)
	CrdController = NewController(config.KubeClient,
		config.KlogClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		logInformerFactory.Klogcontroller().V1alpha1().Klogs(),
		kubeInformerFactory.Core().V1().Nodes())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	logInformerFactory.Start(stopCh)

	if err := CrdController.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
