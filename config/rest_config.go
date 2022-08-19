package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var cfg *rest.Config

func createCfg() {

	var err error
	if env == "prod" {
		cfg, err = rest.InClusterConfig()
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags(masterUrl, kubeConfig)
	}
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
}
