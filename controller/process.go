package controller

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/log"
	logv1alpha1 "k8s.io/log-controller/pkg/apis/logcontroller/v1alpha1"
	"sync"
	"time"
)

var mutex sync.Mutex

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Log controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.logsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")

	// Launch two workers to process Log resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	go c.runPrometheusWorker()
	go initAnalysis(c)
	go c.runCronTask(c.PrometheusMetricQueue, c.PrometheusPodMetricQueue)
	go c.runWarningCronTask()
	go c.runCleanCronTask()
	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func initAnalysis(c *Controller) {
	for {
		if c.prometheusClient.Protocol != "" && c.prometheusClient.Host != "" && c.prometheusClient.Port != "" {
			analysis(c)
			break
		}
	}
}

// runWorker is a long-running function that will couanyntinually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}
	defer c.workqueue.Done(obj)
	var key string
	var ok bool
	if key, ok = obj.(string); !ok {
		c.workqueue.Forget(obj)
		utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
		return true
	}
	if err := c.syncHandler(key); err != nil {
		// Put the item back on the workqueue to handle any transient errors.
		c.workqueue.AddRateLimited(key)
		fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
	}
	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Log resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Log resource with this namespace/name
	log, err := c.logsLister.Logs(namespace).Get(name)
	if err != nil {
		// The Log resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("log '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	err = c.updateLogStatus(log)
	if err != nil {
		return err
	}
	//c.recorder.Event(log, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// Update log ,whick get from kubernetes cluster
func (c *Controller) updateLogStatus(l *logv1alpha1.Log) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	logCopy := l.DeepCopy()
	name := logCopy.Spec.Prometheus.Name
	protocol := logCopy.Spec.Prometheus.Protocol
	host := logCopy.Spec.Prometheus.Host
	port := logCopy.Spec.Prometheus.Port
	period := logCopy.Spec.Prometheus.Period
	warning := logCopy.Spec.Warning
	if period == 0 {
		// default period value is 30 second
		period = log.Period
	}
	if name == "" {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageResourceNoName)
		return fmt.Errorf("field spce.prometheus.name not define")
	}
	if host == "" {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageResourceNoHost)
		return fmt.Errorf("field spce.prometheus.host not define")
	}
	if protocol == "" {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageResourceNoProtocol)
		return fmt.Errorf("field spce.prometheus.host not define")
	}
	if port == "" {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageResourceNoPort)
		return fmt.Errorf("field spce.prometheus.port not define")
	}
	if warning.Sustained.Cpu.Step == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedCpuStep)
		warning.Sustained.Cpu.Step = WarningSustainedStepDefault
	}
	if warning.Sustained.Cpu.Range == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedCpuRange)
		warning.Sustained.Cpu.Range = WarningSustainedRangeDefault
	}
	if warning.Sustained.Cpu.WarningValue == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedCpuWarningValue)
		warning.Sustained.Cpu.WarningValue = WarningSustainedWarningValueDefault
	}
	if warning.Sustained.Cpu.LeftTime == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedCpuLeftTime)
		warning.Sustained.Cpu.LeftTime = WarningSustainedLeftTimeDefault
	}

	if warning.Sustained.Memory.Step == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedMemoryStep)
		warning.Sustained.Memory.Step = WarningSustainedStepDefault
	}
	if warning.Sustained.Memory.Range == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedMemoryRange)
		warning.Sustained.Memory.Range = WarningSustainedRangeDefault
	}
	if warning.Sustained.Memory.WarningValue == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedMemoryWarningValue)
		warning.Sustained.Memory.WarningValue = WarningSustainedWarningValueDefault
	}
	if warning.Sustained.Memory.LeftTime == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedMemoryLeftTime)
		warning.Sustained.Memory.LeftTime = WarningSustainedLeftTimeDefault
	}

	if warning.Sustained.Disk.Range == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedDiskRange)
		warning.Sustained.Disk.Range = WarningSustainedRangeDefault
	}
	if warning.Sustained.Disk.LeftTime == 0 {
		c.recorder.Event(l, corev1.EventTypeWarning, ErrResourceExists, MessageUseDefaultWarningSustainedDiskLeftTime)
		warning.Sustained.Disk.LeftTime = WarningSustainedLeftTimeDefault
	}
	c.prometheusClient.Name = name
	c.prometheusClient.Host = host
	c.prometheusClient.Protocol = protocol
	c.prometheusClient.Port = port
	c.prometheusClient.Period = period
	c.warningSetting = warning
	return nil
}
