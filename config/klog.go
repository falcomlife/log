package config

import (
	versionclient "k8s.io/klog-controller/pkg/generated/clientset/versioned"
	"k8s.io/klog/v2"
)

var KlogClient *versionclient.Clientset

func initKlogClient() {
	if KlogClient == nil {
		if client, err := versionclient.NewForConfig(cfg); err != nil {
			klog.Fatalf("Error building example clientset: %s", err.Error())
		} else {
			KlogClient = client
		}
	}
}
