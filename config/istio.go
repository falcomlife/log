package config

import (
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/klog/v2"
)

var Istioclient *versionedclient.Clientset

func initIstioClient() {
	if Istioclient == nil {
		if client, err := versionedclient.NewForConfig(cfg); err != nil {
			klog.Fatalf("Error building istio clientset: %s", err.Error())
		} else {
			Istioclient = client
		}
	}
}
