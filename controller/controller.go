/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog-controller/log"
	"k8s.io/klog-controller/pkg/apis/klogcontroller/v1alpha1"
	"k8s.io/klog/v2"
	"sync"

	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	logv1alpha1 "k8s.io/klog-controller/pkg/apis/klogcontroller/v1alpha1"
	clientset "k8s.io/klog-controller/pkg/generated/clientset/versioned"
	logscheme "k8s.io/klog-controller/pkg/generated/clientset/versioned/scheme"
	informers "k8s.io/klog-controller/pkg/generated/informers/externalversions/klogcontroller/v1alpha1"
	listers "k8s.io/klog-controller/pkg/generated/listers/klogcontroller/v1alpha1"
)

const controllerAgentName = "log-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a message is sended to wechat
	SuccessSended = "Sended"
	// ErrResourceExists is used as part of the Event 'reason' when a Log fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"
	// MessageResourceSynced is the message used for an Event fired when a Log
	MessageResourceSended = "Daily report message is sended"
	// web url
	//WebUrl = "https://klog.lync2m.com/#/node"
	WebUrl = "https://klog.ciiplat.com/#/node"
)

// Warning message apply to events when k8s cluster check yaml field
const (
	// field spce.prometheus.name not define
	MessageResourceNoName = "spce.prometheus.name not define"
	// field spce.prometheus.host not define
	MessageResourceNoHost = "spce.prometheus.host not define"
	// field spce.prometheus.protocol not define
	MessageResourceNoProtocol = "spce.prometheus.protocol not define"
	// field spce.prometheus.port not define
	MessageResourceNoPort = "spce.prometheus.port not define"
	// use default value
	MessageUseDefaultWarningSustainedCpuStep = "warning.sustained.cpu.step not define,use default value 15"
	// use default value
	MessageUseDefaultWarningSustainedCpuRange = "warning.sustained.cpu.range not define,use default value 90"
	// use default value
	MessageUseDefaultWarningSustainedCpuWarningValue = "warning.sustained.cpu.warningValue not define,use default value 80"
	// use default value
	MessageUseDefaultWarningSustainedCpuLeftTime = "warning.sustained.cpu.leftTime not define,use default value 60"
	// use default value
	MessageUseDefaultWarningSustainedMemoryStep = "warning.sustained.memory.step not define,use default value 15"
	// use default value
	MessageUseDefaultWarningSustainedMemoryRange = "warning.sustained.memory.range not define,use default value 90"
	// use default value
	MessageUseDefaultWarningSustainedMemoryWarningValue = "warning.sustained.memory.warningValue not define,use default value 80"
	// use default value
	MessageUseDefaultWarningSustainedMemoryLeftTime = "warning.sustained.memory.leftTime not define,use default value 60"
	// use default value
	MessageUseDefaultWarningSustainedDiskRange = "warning.sustained.disk.range not define,use default value 90"
	// use default value
	MessageUseDefaultWarningSustainedDiskLeftTime = "warning.sustained.disk.leftTime not define,use default value 60"
)

const (
	// define value, unit second
	WarningSustainedStepDefault = 15
	// define value, unit second
	WarningSustainedRangeDefault = 90
	// define value, cpu %
	WarningSustainedWarningValueDefault = 80
	// define left time for warning
	WarningSustainedLeftTimeDefault = 60
)

// Controller is the controller implementation for Log resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	Kubeclientset kubernetes.Interface
	// logclientset is a clientset for our own API group
	logclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	// prometheus datasource
	prometheusClient log.PrometheusClient
	// setting from yaml
	warningSetting v1alpha1.Warning
	// queue for the node metrics, those come from prometheus
	PrometheusMetricQueue *sync.Map
	// queue for the deployment
	DeploymentQueue *sync.Map
	// queue for the pod metrics, those come from prometheus
	PrometheusPodMetricQueue map[string]*log.Pod
	//NodeCpuAnalysis          map[string]*log.NodeSample
	//NodeMemoryAnalysis       map[string]*log.NodeSample
	Warnings   []*log.WarningList
	logsLister listers.KlogLister
	logsSynced cache.InformerSynced
	Log        *logv1alpha1.Klog
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder

	nodes map[string]corev1.Node
}

// NewController returns a new log controller
func NewController(
	kubeclientset kubernetes.Interface,
	logclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	logInformer informers.KlogInformer,
	nodeInformer coreinformers.NodeInformer) *Controller {

	// Create event broadcaster
	// Add log-controller types to the default Kubernetes Scheme so Events can be
	// logged for log-controller types.
	utilruntime.Must(logscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		Kubeclientset:            kubeclientset,
		logclientset:             logclientset,
		deploymentsLister:        deploymentInformer.Lister(),
		deploymentsSynced:        deploymentInformer.Informer().HasSynced,
		prometheusClient:         log.PrometheusClient{},
		warningSetting:           v1alpha1.Warning{},
		logsLister:               logInformer.Lister(),
		logsSynced:               logInformer.Informer().HasSynced,
		PrometheusMetricQueue:    &sync.Map{},
		DeploymentQueue:          &sync.Map{},
		PrometheusPodMetricQueue: make(map[string]*log.Pod),
		Warnings:                 make([]*log.WarningList, 0),
		workqueue:                workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Logs"),
		recorder:                 recorder,
		nodes:                    map[string]corev1.Node{},
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Log resources change
	logInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueLog,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueLog(new)
		},
	})

	nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.nodeHandler,
		UpdateFunc: func(old interface{}, new interface{}) {
			controller.nodeHandler(new)
		},
	})

	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleDeploymentAdd,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleDeploymentAdd(new)
		},
		DeleteFunc: controller.handleDeploymentDelete,
	})
	return controller
}

func (c *Controller) nodeHandler(obj interface{}) {
	node := obj.(*corev1.Node)
	c.nodes[node.Name] = *node
}

// enqueueLog takes a Log resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Log.
func (c *Controller) enqueueLog(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

func (c *Controller) handleDeploymentAdd(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	c.DeploymentQueue.Store(deployment.Name, deployment)
}

func (c *Controller) handleDeploymentDelete(obj interface{}) {
	deployment := obj.(*appsv1.Deployment)
	c.DeploymentQueue.Delete(deployment.Name)
}
