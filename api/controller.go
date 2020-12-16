package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	controller2 "k8s.io/log-controller/controller"
	"k8s.io/log-controller/log"
	"sort"
)

type NodeController struct {
	beego.Controller
	Ctl *controller2.Controller
}
type PodController struct {
	beego.Controller
	Ctl *controller2.Controller
}

// get node info
func (this *NodeController) Get() {
	var result string
	var list = make([]interface{}, 0)
	for _, v := range this.Ctl.PrometheusMetricQueue {
		list = append(list, v)
	}
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
