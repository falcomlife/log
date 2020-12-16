/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/api"
	controller2 "k8s.io/log-controller/controller"
	clientset "k8s.io/log-controller/pkg/generated/clientset/versioned"
	informers "k8s.io/log-controller/pkg/generated/informers/externalversions"
	"k8s.io/log-controller/pkg/signals"
	"os"
	"time"
)

var (
	masterURL  string
	kubeconfig string
	controller *controller2.Controller
)

func main() {
	go runCrd()
	beegoInit()
}

func runCrd() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()
	var cfg *rest.Config
	var err error
	if os.Getenv("LOG.ENV") == "PROD" {
		cfg, err = rest.InClusterConfig()
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	}
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	//kubeClient, err := kubernetes.NewForConfig(cfg)
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*10)
	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*10)

	controller = controller2.NewController(kubeClient,
		exampleClient,
		exampleInformerFactory.Logcontroller().V1alpha1().Logs(),
		kubeInformerFactory.Core().V1().Nodes())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "test.kubeconfig", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func beegoInit() {
	for {
		select {
		case <-time.After(3 * time.Second):
			if controller != nil {
				beego.SetStaticPath("/", "web")
				beego.Router("/nodes", &api.NodeController{Ctl: controller})
				beego.Router("/pods", &api.PodController{Ctl: controller})
				beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
					AllowAllOrigins:  true,
					AllowMethods:     []string{"*"},
					AllowHeaders:     []string{"*"},
					ExposeHeaders:    []string{"Content-Length"},
					AllowCredentials: true,
				}))
				beego.Run()
			}
		}
	}
}
