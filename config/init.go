package config

func init() {
	analysisFlag()
	createCfg()
	initKubernetesClient()
	initKlogClient()
	initIstioClient()
}
