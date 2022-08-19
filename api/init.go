package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"k8s.io/klog-controller/config"
	"k8s.io/klog-controller/controller"
	"time"
)

func BeegoInit() {
	//db.Orm()
	for {
		select {
		case <-time.After(3 * time.Second):
			if controller.CrdController != nil {
				logs.SetLogger(logs.AdapterConsole)
				beego.BConfig.WebConfig.AutoRender = false
				beego.BConfig.CopyRequestBody = true
				beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{AllowAllOrigins: true, AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"}, ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"}, AllowCredentials: true}))
				beego.SetStaticPath("/", "web")
				beego.Router("/nodes", &NodeController{Ctl: controller.CrdController})
				beego.Router("/pods", &PodController{Ctl: controller.CrdController})
				beego.Router("/warnings", &WarningController{Ctl: controller.CrdController})
				beego.Router("/services", &DeploymentController{Ctl: controller.CrdController})
				beego.Router("/chart", &ChartController{Ctl: controller.CrdController})
				beego.Router("/whitelist", &EnvoyFilterController{Ctl: controller.CrdController, IstioClient: config.Istioclient})
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
