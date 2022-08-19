package config

import "flag"

var kubeConfig string
var masterUrl string
var env string

func analysisFlag() {
	flag.StringVar(&kubeConfig, "kubeconfig", "test.kubeconfig", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterUrl, "masterurl", "", "url to master. Only required if out-of-cluster.")
	flag.StringVar(&env, "env", "test", "program environment")
	flag.Parse()
}
