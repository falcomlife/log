/*
	beego api controller
*/
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/log-controller/common"
	controller2 "k8s.io/log-controller/controller"
	"k8s.io/log-controller/log"
	"k8s.io/utils/integer"
	"sort"
	"strconv"
	"time"
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
type ChartController struct {
	beego.Controller
	Ctl *controller2.Controller
}

// get deployment info
func (this *DeploymentController) Get() {
	action := this.GetString("action")
	name := this.GetString("name")
	namespace := this.GetString("namespace")
	var result string
	switch action {
	case "list":
		result = deployment_list(this)
	case "desc":
		result = deployment_desc(this, name, namespace)
	case "log":
		result = deployment_log(this, name, namespace)
	}

	this.Ctx.WriteString(result)
}

func deployment_list(this *DeploymentController) string {
	var list = make([]interface{}, 0)
	namespaces, err := this.Ctl.Kubeclientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	businessNamespace := make(map[string]bool)
	for _, namespace := range namespaces.Items {
		value, ok := namespace.Labels["role"]
		if ok && value == "business" {
			businessNamespace[namespace.Name] = true
		}
	}
	if err != nil {
		print(err)
	}
	this.Ctl.DeploymentQueue.Range(func(k, v interface{}) bool {
		dep, _ := v.(*v1.Deployment)
		if businessNamespace[dep.Namespace] {
			podNames, startTime, duration, restartTimes, podStatus := this.getPodsByDeployment(dep.UID, dep.Namespace)
			list = append(list, log.Deployment{
				Name:         dep.Name,
				Namespace:    dep.Namespace,
				Kind:         "Deployment",
				Description:  dep.Annotations["description"],
				GitUrl:       dep.Annotations["gitUrl"],
				GatewatUrl:   dep.Annotations["gatewayUrl"],
				PodNames:     podNames,
				PodStatus:    podStatus,
				StartTime:    startTime,
				Duration:     duration,
				RestartTimes: restartTimes,
				InnerUrl:     dep.Annotations["innerUrl"],
				Replicas:     *dep.Spec.Replicas,
				Ready:        dep.Status.ReadyReplicas,
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
		return err.Error()
	} else {
		return string(b)
	}
}

func deployment_desc(this *DeploymentController, name string, namespace string) string {
	m := make(map[string]*[]log.ServiceInfo, 0)
	podInfos := make([]log.ServiceInfo, 0)
	replicasetInfos := make([]log.ServiceInfo, 0)
	deploymentInfos := make([]log.ServiceInfo, 0)
	m["pod"] = &podInfos
	m["replicaset"] = &replicasetInfos
	m["deployment"] = &deploymentInfos

	deployment, err := this.Ctl.Kubeclientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	dname := deployment.Name
	if err != nil {
		return err.Error()
	}
	var replicaset v1.ReplicaSet
	replicasets, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{})
	for _, rs := range replicasets.Items {
		for _, replicasetOwner := range rs.OwnerReferences {
			if replicasetOwner.UID == deployment.UID && *rs.Spec.Replicas != 0 {
				replicaset = rs
				break
			}
		}
	}
	rname := replicaset.Name
	pname := make([]string, 0)
	podAll, err := this.Ctl.Kubeclientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err.Error()
	}
	for _, pod := range podAll.Items {
		for _, owner := range pod.OwnerReferences {
			if owner.UID == replicaset.UID {
				pname = append(pname, pod.Name)
			}
		}
	}

	eventlist, err := this.Ctl.Kubeclientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err.Error()
	}
	for _, event := range eventlist.Items {
		var serviceInfo log.ServiceInfo
		if event.InvolvedObject.Kind == "Pod" {
			for _, pn := range pname {
				if pn == event.InvolvedObject.Name {
					serviceInfo = this.packageForServiceInfo(event)
					podInfos = append(podInfos, serviceInfo)
				}
			}
		} else if event.InvolvedObject.Kind == "ReplicaSet" && event.InvolvedObject.Name == rname {
			serviceInfo = this.packageForServiceInfo(event)
			replicasetInfos = append(replicasetInfos, serviceInfo)
		} else if event.InvolvedObject.Kind == "Deployment" && event.InvolvedObject.Name == dname {
			serviceInfo = this.packageForServiceInfo(event)
			deploymentInfos = append(deploymentInfos, serviceInfo)
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

func (this *DeploymentController) packageForServiceInfo(event corev1.Event) log.ServiceInfo {
	return log.ServiceInfo{
		Name:    event.InvolvedObject.Name,
		Reason:  event.Reason,
		Message: event.Message,
		Type:    event.Type,
		Time:    event.LastTimestamp.String(),
		Kind:    event.InvolvedObject.Kind,
	}
}

func (this *DeploymentController) getPodsByDeployment(deploymentUid types.UID, namespace string) ([]string, []string, []string, []int, []string) {
	pods, err := this.Ctl.Kubeclientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		print(err)
	}
	podNames := make([]string, 0)
	podStatus := make([]string, 0)
	startTime := make([]string, 0)
	duration := make([]string, 0)
	restartTimes := make([]int, 0)

	for _, pod := range pods.Items {
		for _, ownerRc := range pod.OwnerReferences {
			replicaSet, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(namespace).Get(context.Background(), ownerRc.Name, metav1.GetOptions{})
			if err != nil {
				print(err)
				return nil, nil, nil, nil, nil
			}
			for _, ownerDep := range replicaSet.OwnerReferences {
				if ownerDep.UID == deploymentUid {
					podNames = append(podNames, pod.Name)
					podStatus = append(podStatus, string(pod.Status.Phase))
					if pod.Status.StartTime != nil {
						startTime = append(startTime, pod.Status.StartTime.String())
					} else {
						startTime = append(startTime, "xxxx-xx-xx xx:xx:xx +0800 CST")
					}
					durationTimeStamp := time.Now().Unix() - pod.Status.StartTime.Unix()
					duration = append(duration, strconv.Itoa(int(durationTimeStamp)))
					restartTimes = append(restartTimes, maxContainerRestarts(&pod))
				}
			}
		}
	}
	return podNames, startTime, duration, restartTimes, podStatus
}

func maxContainerRestarts(pod *corev1.Pod) int {
	maxRestarts := 0
	for _, c := range pod.Status.ContainerStatuses {
		maxRestarts = integer.IntMax(maxRestarts, int(c.RestartCount))
	}
	return maxRestarts
}

func deployment_log(this *DeploymentController, name string, namespace string) string {
	deployment, err := this.Ctl.Kubeclientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err.Error()
	}
	var replicaset v1.ReplicaSet
	replicasets, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{})
	for _, rs := range replicasets.Items {
		for _, replicasetOwner := range rs.OwnerReferences {
			if replicasetOwner.UID == deployment.UID && *rs.Spec.Replicas != 0 {
				replicaset = rs
				break
			}
		}
	}
	logs := make([]log.ServiceLog, 0)
	podAll, err := this.Ctl.Kubeclientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	for _, pod := range podAll.Items {
		for _, owner := range pod.OwnerReferences {
			if owner.UID == replicaset.UID {
				for _, container := range pod.Spec.Containers {
					req := this.Ctl.Kubeclientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{Follow: false, Container: container.Name})
					podLogs, err := req.Stream(context.Background())
					if err != nil {
						return "error in opening stream"
					}
					defer podLogs.Close()

					buf := new(bytes.Buffer)
					_, err = io.Copy(buf, podLogs)
					if err != nil {
						return "error in copy information from podLogs to buf"
					}
					str := buf.String()
					logs = append(logs, log.ServiceLog{Name: pod.Name + "/" + container.Name, Content: str})
				}
			}
		}
	}
	b, err := json.Marshal(logs)
	if err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

func (this *DeploymentController) Post() {
	input := log.DeploymentBody{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &input); err != nil {
		this.Data["json"] = err.Error()
	}
	var result string = ""
	switch input.Action {
	case "restart":
		result = this.restart(input)
	case "rollback":
		result = this.rollback(input)
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *DeploymentController) restart(input log.DeploymentBody) string {
	if input.Kind == "Deployment" {
		logs.Info("debug 1")
		deployment, err := this.Ctl.Kubeclientset.AppsV1().Deployments(input.Namespace).Get(context.Background(), input.Name, metav1.GetOptions{})
		if err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		logs.Info("debug 2")
		pods, err := this.Ctl.Kubeclientset.CoreV1().Pods(input.Namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		logs.Info("debug 3")
		for _, pod := range pods.Items {
			for _, ownerRc := range pod.OwnerReferences {
				replicaSet, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(input.Namespace).Get(context.Background(), ownerRc.Name, metav1.GetOptions{})
				logs.Info("debug 4")
				if err != nil {
					logs.Error(err.Error())
					return err.Error()
				}
				for _, ownerDep := range replicaSet.OwnerReferences {
					if ownerDep.UID == deployment.UID {
						err := this.Ctl.Kubeclientset.CoreV1().Pods(input.Namespace).Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
						if err != nil {
							logs.Error(err.Error())
							return err.Error()
						}
						logs.Info("debug 5")
					}
				}
			}
		}
	}
	return "success"
}

func (this *DeploymentController) rollback(input log.DeploymentBody) string {
	if input.Kind == "Deployment" {
		deployment, err := this.Ctl.Kubeclientset.AppsV1().Deployments(input.Namespace).Get(context.Background(), input.Name, metav1.GetOptions{})
		if err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		replicasetMatch := make([]v1.ReplicaSet, 0)
		replicasets, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(input.Namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		for _, replicaset := range replicasets.Items {
			for _, replicasetOwner := range replicaset.OwnerReferences {
				if replicasetOwner.UID == deployment.UID {
					replicasetMatch = append(replicasetMatch, replicaset)
				}
			}
		}
		sort.Sort(common.ReplicaSetSlice(replicasetMatch))
		deployment.Spec.Template.Spec.Containers = replicasetMatch[1].Spec.Template.Spec.Containers
		_, err = this.Ctl.Kubeclientset.AppsV1().Deployments(input.Namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
		if err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
	}
	return "success"
}

// get warning info
// Deprecated
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
// Deprecated
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
// Deprecated
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

// Deprecated
func (this *ChartController) Get() {
	var list = make([]string, 0)
	value := this.Input()
	date := value.Get("date")
	files, _ := ioutil.ReadDir("web/ai/result/" + date)
	for _, f := range files {
		list = append(list, f.Name())
	}
	result, _ := json.Marshal(list)
	this.Ctx.WriteString(string(result))
}

//func (this *ChartController) Get() {
//	r := make([]datasource.Result, 0)
//	datasource.Exec(&r, "SELECT * FROM analy.result where time >= ? and time < ? order by time desc limit ?;", "2021-07-18", "2021-07-19", 6)
//	result := make([]Chart, 0)
//	for _, res := range r {
//		originMap := make(map[string]interface{})
//		err := json.Unmarshal([]byte(res.OriginIndex), &originMap)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		originIndex, _ := originMap["index"].([]interface{})
//		originData, _ := originMap["data"].([]interface{})
//		fmt.Println(originIndex,originData)
//		var chart = Chart{}
//		origin(res, &chart)
//		auto(res, &chart)
//		//iforest(res, chart)
//		//merge(res, chart)
//		result = append(result, chart)
//	}
//	var resultStr string
//	resultByte, err := json.Marshal(result)
//	if err != nil {
//		resultStr = err.Error()
//	} else {
//		resultStr = string(resultByte)
//	}
//	this.Ctx.WriteString(resultStr)
//}
//
//func origin(res datasource.Result, chart *Chart) {
//	originMap := make(map[string]interface{})
//	err := json.Unmarshal([]byte(res.OriginIndex), &originMap)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	originIndex, _ := originMap["index"].([]interface{})
//	originData, _ := originMap["data"].([]interface{})
//	xAxis := XAxis{
//		Data: make([]string, 0),
//	}
//	series := Series{
//		Data: make([]string, 0),
//	}
//	for _, oi := range originIndex {
//		oiflo, _ := oi.(float64)
//		xAxis.Data = append(xAxis.Data, (time.Unix(int64(oiflo)/1000, 0).Format("15:04:05")))
//	}
//	for _, od := range originData {
//		o, _ := od.([]interface{})
//		odflo, _ := o[0].(float64)
//		series.Data = append(series.Data, strconv.FormatFloat(odflo, 'f', -1, 64))
//	}
//	chart.XAxis = xAxis
//	chart.Series = append(chart.Series, series)
//}
//
//func auto(res datasource.Result, chart *Chart) {
//	originMap := make(map[string]interface{})
//	autoMap := make(map[string]interface{})
//	json.Unmarshal([]byte(res.OriginIndex), &originMap)
//	json.Unmarshal([]byte(res.AutoIndex), &autoMap)
//
//	indexs, datas := mergedata(originMap, autoMap)
//
//	xAxis := XAxis{
//		Data: make([]string, 0),
//	}
//	series := Series{
//		Data: make([]string, 0),
//	}
//	for _, oi := range indexs {
//		value, _ := oi.(float64)
//		xAxis.Data = append(xAxis.Data, (time.Unix(int64(value)/1000, 0).Format("15:04:05")))
//	}
//	for _, od := range datas {
//		series.Data = append(series.Data, od)
//	}
//	chart.XAxis = xAxis
//	chart.Series = append(chart.Series, series)
//}

//func iforest(res datasource.Result, chart Chart) {
//	originMap := make(map[string]interface{})
//	iforestMap := make(map[string]interface{})
//
//	json.Unmarshal([]byte(res.OriginIndex), originMap)
//	json.Unmarshal([]byte(res.IforestIndex), iforestMap)
//
//	indexs, datas := mergedata(originMap, iforestMap)
//
//	xAxis := XAxis{
//		Data: make([]string, 0),
//	}
//	series := Series{
//		Data: make([]string, 0),
//	}
//	for _, oi := range indexs {
//		value, ok := oi.(int64)
//		if !ok {
//			fmt.Println(ok)
//		}
//		xAxis.Data = append(xAxis.Data, (time.Unix(value, 0).Format("15:04:05")))
//	}
//	for _, od := range datas {
//		series.Data = append(series.Data, string(od))
//	}
//	chart.XAxis = xAxis
//	chart.Series = append(chart.Series, series)
//}

//
//func merge(res datasource.Result, chart Chart) {
//	originMap := make(map[string]string)
//	json.Unmarshal([]byte(res.IforestIndex), originMap)
//	originIndex := originMap["index"]
//	originData := originMap["data"]
//	xAxis := XAxis{
//		Data: make([]string, 0),
//	}
//	series := Series{
//		Data: make([]string, 0),
//	}
//	for _, oi := range originIndex {
//
//		xAxis.Data = append(xAxis.Data, (time.Unix(int64(oi), 0).Format("15:04:05")))
//	}
//	for _, od := range originData {
//		series.Data = append(series.Data, string(od))
//	}
//	chart.XAxis = xAxis
//	chart.Series = append(chart.Series, series)
//}
//}
