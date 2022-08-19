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
	simplejson "github.com/bitly/go-simplejson"
	"github.com/go-basic/uuid"
	"html/template"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog-controller/common"
	"k8s.io/klog-controller/controller"
	"k8s.io/klog-controller/log"
	"k8s.io/utils/integer"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DeploymentController struct {
	beego.Controller
	Ctl *controller.Controller
}
type WarningController struct {
	beego.Controller
	Ctl *controller.Controller
}
type NodeController struct {
	beego.Controller
	Ctl *controller.Controller
}
type PodController struct {
	beego.Controller
	Ctl *controller.Controller
}
type ChartController struct {
	beego.Controller
	Ctl *controller.Controller
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
		logs.Error(err.Error())
	}
	pod_replicaSet_map := make(map[types.UID]map[string]interface{})
	if pods, err := this.Ctl.Kubeclientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{}); err != nil {
		logs.Error(err.Error())
		return err.Error()
	} else {
		if replicaSets, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets("").List(context.Background(), metav1.ListOptions{}); err != nil {
			logs.Error(err.Error())
			return err.Error()
		} else {
			for _, pod := range pods.Items {
				for _, ownerRc := range pod.OwnerReferences {
					for _, replicaSet := range replicaSets.Items {
						if ownerRc.UID == replicaSet.UID {
							_, ok := pod_replicaSet_map[replicaSet.UID]
							if !ok {
								pod_replicaSet_map[replicaSet.UID] = make(map[string]interface{})
								pod_replicaSet_map[replicaSet.UID]["replicaSet"] = replicaSet
								pod_replicaSet_map[replicaSet.UID]["pods"] = make([]corev1.Pod, 0)
								pod_replicaSet_map[replicaSet.UID]["pods"] = append((pod_replicaSet_map[replicaSet.UID]["pods"]).([]corev1.Pod), pod)
							} else {
								pod_replicaSet_map[replicaSet.UID]["pods"] = append((pod_replicaSet_map[replicaSet.UID]["pods"]).([]corev1.Pod), pod)
							}
						}
					}
				}
			}
		}
	}
	this.Ctl.DeploymentQueue.Range(func(k, v interface{}) bool {
		dep, _ := v.(*v1.Deployment)
		if businessNamespace[dep.Namespace] {
			podNames, startTime, duration, restartTimes, podStatus := this.getPodsByDeployment(dep.UID, dep.Namespace, pod_replicaSet_map)
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
			return n1.Namespace < n2.Namespace
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

	if deployment, err := this.Ctl.Kubeclientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		logs.Error(err.Error())
		return err.Error()
	} else {
		dname := deployment.Name
		var replicasetSlice []v1.ReplicaSet
		if replicasets, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{}); err != nil {
			logs.Error(err.Error())
			return err.Error()
		} else {
			for _, rs := range replicasets.Items {
				for _, replicasetOwner := range rs.OwnerReferences {
					if replicasetOwner.UID == deployment.UID {
						replicasetSlice = append(replicasetSlice, rs)
					}
				}
			}
			pname := make([]string, 0)
			if podAll, err := this.Ctl.Kubeclientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{}); err != nil {
				logs.Error(err.Error())
				return err.Error()
			} else {
				for _, replicaset := range replicasetSlice {
					for _, pod := range podAll.Items {
						for _, owner := range pod.OwnerReferences {
							if owner.UID == replicaset.UID {
								pname = append(pname, pod.Name)
							}
						}
					}
				}
				if eventlist, err := this.Ctl.Kubeclientset.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{}); err != nil {
					logs.Error(err.Error())
					return err.Error()
				} else {
					for _, replicaset := range replicasetSlice {
						rname := replicaset.Name
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
					}
				}
			}
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
		Time:    event.CreationTimestamp.String(),
		Kind:    event.InvolvedObject.Kind,
	}
}

func (this *DeploymentController) getPodsByDeployment(deploymentUid types.UID, namespace string, pod_replicaSet_map map[types.UID]map[string]interface{}) ([]string, []string, []string, []int, []string) {

	podNames := make([]string, 0)
	podStatus := make([]string, 0)
	startTime := make([]string, 0)
	duration := make([]string, 0)
	restartTimes := make([]int, 0)

	for _, v := range pod_replicaSet_map {
		var replicaSet v1.ReplicaSet = v["replicaSet"].(v1.ReplicaSet)
		var pods []corev1.Pod = v["pods"].([]corev1.Pod)
		for _, ownerDep := range replicaSet.OwnerReferences {
			if ownerDep.UID == deploymentUid {
				for _, pod := range pods {
					podNames = append(podNames, pod.Name)
					podStatus = append(podStatus, string(pod.Status.Phase))
					if pod.Status.StartTime != nil {
						startTime = append(startTime, pod.Status.StartTime.String())
					} else {
						startTime = append(startTime, "xxxx-xx-xx xx:xx:xx +0800 CST")
					}
					var durationTimeStamp int64 = 0
					if &pod.Status != nil && pod.Status.StartTime != nil {
						durationTimeStamp = time.Now().Unix() - pod.Status.StartTime.Unix()
					}
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
	if deployment, err := this.Ctl.Kubeclientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
		logs.Error(err.Error())
		return err.Error()
	} else {
		var replicasetSlice []v1.ReplicaSet
		if replicasets, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(namespace).List(context.Background(), metav1.ListOptions{}); err != nil {
			logs.Error(err.Error())
			return err.Error()
		} else {
			for _, rs := range replicasets.Items {
				for _, replicasetOwner := range rs.OwnerReferences {
					if replicasetOwner.UID == deployment.UID && *rs.Spec.Replicas != 0 {
						replicasetSlice = append(replicasetSlice, rs)
					}
				}
			}
			logSlice := make([]log.ServiceLog, 0)
			if podAll, err := this.Ctl.Kubeclientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{}); err != nil {
				logs.Error(err.Error())
				return err.Error()
			} else {
				for _, replicaset := range replicasetSlice {
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
									logSlice = append(logSlice, log.ServiceLog{Name: pod.Name + "/" + container.Name, Content: str})
								}
							}
						}
					}
				}
			}
			b, err := json.Marshal(logSlice)
			if err != nil {
				return err.Error()
			} else {
				return string(b)
			}
		}
	}
}

// restart or rollback servcice
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
		deployment, err := this.Ctl.Kubeclientset.AppsV1().Deployments(input.Namespace).Get(context.Background(), input.Name, metav1.GetOptions{})
		if err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		pods, err := this.Ctl.Kubeclientset.CoreV1().Pods(input.Namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		for _, pod := range pods.Items {
			for _, ownerRc := range pod.OwnerReferences {
				replicaSet, err := this.Ctl.Kubeclientset.AppsV1().ReplicaSets(input.Namespace).Get(context.Background(), ownerRc.Name, metav1.GetOptions{})
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

func (this *DeploymentController) Put() {
	action := this.GetString("action")
	var result string
	switch action {
	case "java":
		result = this.add_java()
	case "javaMul":
		result = this.add_javaMul()
	case "npm":
		result = this.add_npm()
	}
	this.Ctx.WriteString(result)
}

func (this *DeploymentController) add_java() string {
	input := log.Ci{}
	if js, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err != nil {
		logs.Error(err.Error())
		return err.Error()
	} else {
		input.Name = js.Get("name").MustString()
		input.Namespace = js.Get("namespace").MustString()
		input.Describe = js.Get("describe").MustString()
		input.OnlyRefs = js.Get("onlyRefs").MustString()
		input.GitUrl = js.Get("gitUrl").MustString()
		input.Prefix = js.Get("prefix").MustString()
		input.Health = js.Get("health").MustString()
		input.Port = js.Get("port").MustString()
		input.GatewayRealmName = this.Ctl.Log.Spec.Template.GatewayRealmName
		input.GatewayHost = strings.Split(this.Ctl.Log.Spec.Template.GatewayRealmName, "//")[1]
		input.Registry = log.Registry{
			Username: this.Ctl.Log.Spec.Template.Username,
			Password: this.Ctl.Log.Spec.Template.Password,
			Address:  this.Ctl.Log.Spec.Template.Address,
		}
		input.Env = this.Ctl.Log.Spec.Template.Env
		input.Package = log.Package{
			Image:          this.Ctl.Log.Spec.Template.Registry.Java.PackageImage,
			ArtifactsPaths: js.Get("artifactsPaths").MustString(),
		}
		input.Release = log.Release{
			Image: this.Ctl.Log.Spec.Template.Registry.Java.ReleaseImage,
		}
		input.Deploy = log.Deploy{
			Image: this.Ctl.Log.Spec.Template.Registry.Java.DeployImage,
		}
		uid := uuid.New()
		if err := os.MkdirAll(uid+"/.devops/yaml/"+this.Ctl.Log.Spec.Template.Env, os.ModePerm); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := this.generateCiFile("template/java/.gitlab-ci.yml", input, uid); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := this.generateCiFile("template/java/deployment.yaml", input, uid); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := this.generateCiFile("template/java/Dockerfile", input, uid); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := os.Rename(uid+"/Dockerfile", uid+"/.devops/Dockerfile"); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := os.Rename(uid+"/deployment.yaml", uid+"/.devops/yaml/"+this.Ctl.Log.Spec.Template.Env+"/deployment.yaml"); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := common.Zip(uid, "web/ci.zip"); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		logs.Info("success")
		return "success"
	}
}

func (this *DeploymentController) add_javaMul() string {
	input := log.Ci{}
	if js, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err != nil {
		logs.Error(err.Error())
		return err.Error()
	} else {
		//input.Name = js.Get("name").MustString()
		//input.Describe = js.Get("describe").MustString()
		//input.Health = js.Get("health").MustString()
		//input.Port = js.Get("port").MustString()
		input.Modules = make([]log.Modules, 0)
		arr := js.Get("modules").MustArray()
		for _, module := range arr {
			m, _ := module.(map[string]interface{})
			input.Modules = append(input.Modules, log.Modules{
				Name:           m["name"].(string),
				Describe:       m["describe"].(string),
				ArtifactsPaths: m["artifactsPaths"].(string),
				Prefix:         m["prefix"].(string),
				Health:         m["health"].(string),
				Port:           m["port"].(string),
			})
		}
		input.Namespace = js.Get("namespace").MustString()
		input.OnlyRefs = js.Get("onlyRefs").MustString()
		input.GitUrl = js.Get("gitUrl").MustString()
		input.GatewayRealmName = this.Ctl.Log.Spec.Template.GatewayRealmName
		input.GatewayHost = strings.Split(this.Ctl.Log.Spec.Template.GatewayRealmName, "//")[1]
		input.Registry = log.Registry{
			Username: this.Ctl.Log.Spec.Template.Username,
			Password: this.Ctl.Log.Spec.Template.Password,
			Address:  this.Ctl.Log.Spec.Template.Address,
		}
		input.Env = this.Ctl.Log.Spec.Template.Env
		input.Package = log.Package{
			Image:          this.Ctl.Log.Spec.Template.Registry.Java.PackageImage,
			ArtifactsPaths: js.Get("artifactsPaths").MustString(),
		}
		input.Release = log.Release{
			Image: this.Ctl.Log.Spec.Template.Registry.Java.ReleaseImage,
		}
		input.Deploy = log.Deploy{
			Image: this.Ctl.Log.Spec.Template.Registry.Java.DeployImage,
		}
		uid := uuid.New()
		if err := os.MkdirAll(uid+"/.devops/", os.ModePerm); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := this.generateCiFile("template/javaMul/.gitlab-ci.yml", input, uid); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		for _, module := range input.Modules {
			if err := os.MkdirAll(uid+"/.devops/"+module.Name+"/yaml/"+this.Ctl.Log.Spec.Template.Env, os.ModePerm); err != nil {
				logs.Error(err.Error())
				return err.Error()
			}
			input.Name = module.Name
			input.Describe = module.Describe
			input.Port = module.Port
			input.Prefix = module.Prefix
			input.Health = module.Health
			input.Package.ArtifactsPaths = module.ArtifactsPaths
			if err := this.generateCiFile("template/javaMul/deployment.yaml", input, uid); err != nil {
				logs.Error(err.Error())
				return err.Error()
			}
			if err := this.generateCiFile("template/javaMul/Dockerfile", input, uid); err != nil {
				logs.Error(err.Error())
				return err.Error()
			}
			if err := os.Rename(uid+"/Dockerfile", uid+"/.devops/"+module.Name+"/Dockerfile"); err != nil {
				logs.Error(err.Error())
				return err.Error()
			}
			if err := os.Rename(uid+"/deployment.yaml", uid+"/.devops/"+module.Name+"/yaml/"+this.Ctl.Log.Spec.Template.Env+"/deployment.yaml"); err != nil {
				logs.Error(err.Error())
				return err.Error()
			}
		}
		if err := common.Zip(uid, "web/ci.zip"); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		logs.Info("success")
		return "success"
	}
}

func (this *DeploymentController) add_npm() string {
	input := log.Ci{}
	if js, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err != nil {
		logs.Error(err.Error())
		return err.Error()
	} else {
		input.Name = js.Get("name").MustString()
		input.Namespace = js.Get("namespace").MustString()
		input.Describe = js.Get("describe").MustString()
		input.OnlyRefs = js.Get("onlyRefs").MustString()
		input.GitUrl = js.Get("gitUrl").MustString()
		input.Prefix = js.Get("prefix").MustString()
		input.Health = js.Get("health").MustString()
		input.Port = js.Get("port").MustString()
		input.WebRealmName = this.Ctl.Log.Spec.Template.WebRealmName
		input.WebHost = strings.Split(this.Ctl.Log.Spec.Template.WebRealmName, "//")[1]
		input.Registry = log.Registry{
			Username: this.Ctl.Log.Spec.Template.Username,
			Password: this.Ctl.Log.Spec.Template.Password,
			Address:  this.Ctl.Log.Spec.Template.Address,
		}
		input.Env = this.Ctl.Log.Spec.Template.Env
		input.Package = log.Package{
			Image:          this.Ctl.Log.Spec.Template.Registry.Npm.PackageImage,
			ArtifactsPaths: js.Get("artifactsPaths").MustString(),
		}
		input.Release = log.Release{
			Image: this.Ctl.Log.Spec.Template.Registry.Npm.ReleaseImage,
		}
		input.Deploy = log.Deploy{
			Image: this.Ctl.Log.Spec.Template.Registry.Npm.DeployImage,
		}
		uid := uuid.New()
		if err := os.MkdirAll(uid+"/.devops/yaml/"+this.Ctl.Log.Spec.Template.Env, os.ModePerm); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := this.generateCiFile("template/npm/.gitlab-ci.yml", input, uid); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := this.generateCiFile("template/npm/deployment.yaml", input, uid); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := this.generateCiFile("template/npm/Dockerfile", input, uid); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := os.Rename(uid+"/Dockerfile", uid+"/.devops/Dockerfile"); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if source, err := os.Open("template/npm/nginx.conf"); err != nil {
			defer source.Close()
			logs.Error(err.Error())
			return err.Error()
		} else {
			defer source.Close()
			if destination, err := os.Create(uid + "/.devops/nginx.conf"); err != nil {
				defer destination.Close()
				logs.Error(err.Error())
				return err.Error()
			} else {
				defer destination.Close()
				if _, err := io.Copy(destination, source); err != nil {
					logs.Error(err.Error())
					return err.Error()
				}
			}
		}
		if err := os.Rename(uid+"/deployment.yaml", uid+"/.devops/yaml/"+this.Ctl.Log.Spec.Template.Env+"/deployment.yaml"); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		if err := common.Zip(uid, "web/ci.zip"); err != nil {
			logs.Error(err.Error())
			return err.Error()
		}
		logs.Info("success")
		return "success"
	}
}

func (this *DeploymentController) generateCiFile(templateFilePath string, input log.Ci, uid string) error {
	if temp, err := template.ParseFiles(templateFilePath); err != nil {
		logs.Error(err.Error())
		return err
	} else {
		filename := strings.Split(templateFilePath, "/")[2]
		if f, err := os.Create(uid + "/" + filename); err != nil {
			defer f.Close()
			return err

		} else {
			defer f.Close()
			if err := temp.Execute(f, input); err != nil {
				return err
			}
		}
	}
	return nil
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
