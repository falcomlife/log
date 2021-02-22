/*
	beego api controller
*/
package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	v1 "k8s.io/api/apps/v1"
	controller2 "k8s.io/log-controller/controller"
	"k8s.io/log-controller/log"
	"sort"
)

type DeploymentController struct {
	beego.Controller
	Ctl *controller2.Controller
}
type WarningController struct {
	beego.Controller
	Ctl *controller2.Controller
}
type NodeController struct {
	beego.Controller
	Ctl *controller2.Controller
}
type PodController struct {
	beego.Controller
	Ctl *controller2.Controller
}

// get deployment info
func (this *DeploymentController) Get() {
	var result string
	var list = make([]interface{}, 0)
	this.Ctl.DeploymentQueue.Range(func(k, v interface{}) bool {
		dep, _ := v.(*v1.Deployment)
		if dep.Namespace == "middleplatform" {
			list = append(list, log.Deployment{
				Name:        dep.Name,
				NameSpace:   dep.Namespace,
				Kind:        "Deployment",
				Description: dep.Annotations["description"],
				GitUrl:      dep.Annotations["gitUrl"],
				GatewatUrl:  dep.Annotations["gatewayUrl"],
				InnerUrl:    dep.Annotations["innerUrl"],
				Replicas:    *dep.Spec.Replicas,
				Ready:       dep.Status.ReadyReplicas,
			})
		}
		return true
	})
	if len(list) != 0 {
		sort.SliceStable(list, func(i, j int) bool {
			n1, _ := list[i].(log.Deployment)
			n2, _ := list[j].(log.Deployment)
			return n1.Name < n2.Name
		})
	}
	b, err := json.Marshal(list)
	if err != nil {
		result = err.Error()
	} else {
		result = string(b)
	}
	this.Ctx.WriteString(result)
}

// get warning info
func (this *WarningController) Get() {
	var result string
	b, err := json.Marshal(this.Ctl.Warnings)
	if err != nil {
		result = err.Error()
	} else {
		result = string(b)
	}
	this.Ctx.WriteString(result)
}

// get node info
func (this *NodeController) Get() {
	var result string
	var list = make([]interface{}, 0)
	this.Ctl.PrometheusMetricQueue.Range(func(k, v interface{}) bool {
		list = append(list, v)
		return true
	})
	sort.SliceStable(list, func(i, j int) bool {
		n1, _ := list[i].(log.Node)
		n2, _ := list[j].(log.Node)
		return n1.Name < n2.Name
	})
	var resList = make([]interface{}, 0)
	for _, li := range list {
		resList = append(resList, li)
		resList = append(resList, li)
	}
	b, err := json.Marshal(resList)
	if err != nil {
		result = err.Error()
	} else {
		result = string(b)
	}
	this.Ctx.WriteString(result)
}

// get pod info
func (this *PodController) Get() {
	var result string
	var list = make([]interface{}, 0)
	for _, v := range this.Ctl.PrometheusPodMetricQueue {
		list = append(list, v)
	}
	sort.SliceStable(list, func(i, j int) bool {
		n1, _ := list[i].(log.Pod)
		n2, _ := list[j].(log.Pod)
		return n1.Namespace < n2.Namespace
	})
	b, err := json.Marshal(list)
	if err != nil {
		result = err.Error()
	} else {
		result = string(b)
	}
	this.Ctx.WriteString(result)
}
