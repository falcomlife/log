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

package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/log"

	clientset "k8s.io/log-controller/pkg/generated/clientset/versioned"
	logscheme "k8s.io/log-controller/pkg/generated/clientset/versioned/scheme"
	informers "k8s.io/log-controller/pkg/generated/informers/externalversions/logcontroller/v1alpha1"
	listers "k8s.io/log-controller/pkg/generated/listers/logcontroller/v1alpha1"
)

const controllerAgentName = "log-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Log is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Log fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"
	// MessageResourceSynced is the message used for an Event fired when a Log
	// is synced successfully
	MessageResourceSynced = "Log synced successfully"
	// field spce.prometheus.period not define
	MessageResourceNoPeriod = "spce.prometheus.period not define"
	// field spce.prometheus.name not define
	MessageResourceNoName = "spce.prometheus.name not define"
	// field spce.prometheus.host not define
	MessageResourceNoHost = "spce.prometheus.host not define"
	// field spce.prometheus.protocol not define
	MessageResourceNoProtocol = "spce.prometheus.protocol not define"
	// field spce.prometheus.port not define
	MessageResourceNoPort = "spce.prometheus.port not define"
)

// Controller is the controller implementation for Log resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// logclientset is a clientset for our own API group
	logclientset clientset.Interface
	// prometheus datasource
	prometheusClient log.PrometheusClient
	// queue for the metrics, those come from prometheus
	prometheusMetricQueue map[string]log.Node

	logsLister listers.LogLister
	logsSynced cache.InformerSynced
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new log controller
func NewController(
	kubeclientset kubernetes.Interface,
	logclientset clientset.Interface,
	logInformer informers.LogInformer) *Controller {

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
		kubeclientset:         kubeclientset,
		logclientset:          logclientset,
		prometheusClient:      log.PrometheusClient{},
		logsLister:            logInformer.Lister(),
		logsSynced:            logInformer.Informer().HasSynced,
		prometheusMetricQueue: make(map[string]log.Node),
		workqueue:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Logs"),
		recorder:              recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Log resources change
	logInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueLog,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueLog(new)
		},
	})
	return controller
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

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Log resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Log resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Log, we should not do anything more
		// with it.
		if ownerRef.Kind != "Log" {
			return
		}

		log, err := c.logsLister.Logs(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of log '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueLog(log)
		return
	}
}
