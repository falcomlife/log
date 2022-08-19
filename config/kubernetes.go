package config

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

var KubeClient *kubernetes.Clientset

func initKubernetesClient() {
	if KubeClient == nil {
		if client, err := kubernetes.NewForConfig(cfg); err != nil {
			klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
		} else {
			KubeClient = client
		}
	}
}
