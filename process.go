package main

import (
	"bytes"
	"fmt"
	cron "github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/log"
	logv1alpha1 "k8s.io/log-controller/pkg/apis/logcontroller/v1alpha1"
	"math"
	"strconv"
	"text/template"
	"time"
)

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
	go c.runPrometheusWorder()
	go c.runCronTask(c.prometheusMetricQueue)
	// Launch two workers to process Log resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// run a work for get prometheus
func (c *Controller) runPrometheusWorder() {
	// default period value is 30
	period := c.prometheusClient.Period
	if period == 0 {
		period = log.Period
	}
	for {
		select {
		case <-time.Tick(time.Duration(period) * time.Second):
			name := c.prometheusClient.Name
			host := c.prometheusClient.Host
			port := c.prometheusClient.Port
			if name != "" && host != "" && port != "" {
				nodes1, err := c.prometheusClient.Get(log.NodeCpuUsedPercentage)
				if err != nil {
					continue
				}
				c.batchForCpu(nodes1)
				nodes2, err := c.prometheusClient.Get(log.NodeMemoryUsed)
				if err != nil {
					continue
				}
				c.batchForMem(nodes2)
				c.prometheusClient.SamplingTimes++
			} else {
				continue
			}
		}
	}
}

func (c *Controller) batchForCpu(nodes map[string]log.Node) {
	for _, node := range nodes {
		var nodeOld log.Node
		if _, ok := c.prometheusMetricQueue[node.Name]; ok {
			// there is node in map
			nodeOld = c.prometheusMetricQueue[node.Name]
		} else {
			// there is not node in map
			nodeOld = log.Node{
				Name:      node.Name,
				Cpu:       make(map[string]log.Cpu),
				CpuSumMin: math.MaxFloat64,
				MemMin:    math.MaxFloat64,
			}
		}
		cpuNewvValueSum := 0.0
		for cpuNewKey, cpuNewValue := range node.Cpu {
			var cpuOld log.Cpu
			if _, ok := nodeOld.Cpu[cpuNewKey]; ok {
				// there is cpu in map
				cpuOld = nodeOld.Cpu[cpuNewKey]
			} else {
				// there is not cpu in map
				cpuOld = log.Cpu{}
			}
			cpuOld.Value = cpuNewValue.Value
			cpuOld.Time = cpuNewValue.Time
			nodeOld.Cpu[cpuNewKey] = cpuOld
			cpuNewvValueSum += cpuNewValue.Value
		}
		cpuSumNew := cpuNewvValueSum / float64(len(node.Cpu)) * 100
		if nodeOld.CpuSumMax < cpuSumNew {
			cpuValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuSumNew), 64)
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.CpuSumMax = cpuValue
			nodeOld.CpuSumMaxTime = time.Now()
		}
		if nodeOld.CpuSumMin > cpuSumNew {
			cpuValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuSumNew), 64)
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.CpuSumMin = cpuValue
			nodeOld.CpuSumMinTime = time.Now()
		}
		cpuAvgtemp := (nodeOld.CpuSumAvg*float64(c.prometheusClient.SamplingTimes) + cpuSumNew) / float64(c.prometheusClient.SamplingTimes+1)
		cpuAvg, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuAvgtemp), 64)
		if err != nil {
			klog.Warning(err)
		}
		nodeOld.CpuSumAvg = cpuAvg
		c.prometheusMetricQueue[node.Name] = nodeOld
	}
}

func (c *Controller) batchForMem(nodes map[string]log.Node) {
	for _, node := range nodes {
		var nodeOld log.Node
		if _, ok := c.prometheusMetricQueue[node.Name]; ok {
			// there is node in map
			nodeOld = c.prometheusMetricQueue[node.Name]
		} else {
			// there is not node in map
			nodeOld = log.Node{
				Name:      node.Name,
				Cpu:       make(map[string]log.Cpu),
				CpuSumMin: math.MaxFloat64,
				MemMin:    math.MaxFloat64,
			}
		}
		memUsedNew := node.MemMax / math.Pow(2, 30)
		if nodeOld.MemMax < memUsedNew {
			memValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memUsedNew), 64)
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.MemMax = memValue
			nodeOld.MemMaxTime = time.Now()
		}
		if nodeOld.MemMin > memUsedNew {
			memValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memUsedNew), 64)
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.MemMin = memValue
			nodeOld.MemMinTime = time.Now()
		}
		memAvgtemp := (nodeOld.MemAvg*float64(c.prometheusClient.SamplingTimes) + memUsedNew) / float64(c.prometheusClient.SamplingTimes+1)
		memAvg, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memAvgtemp), 64)
		if err != nil {
			klog.Warning(err)
		}
		nodeOld.MemAvg = memAvg
		c.prometheusMetricQueue[node.Name] = nodeOld
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

	c.recorder.Event(log, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

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
	c.prometheusClient.Name = name
	c.prometheusClient.Host = host
	c.prometheusClient.Protocol = protocol
	c.prometheusClient.Port = port
	c.prometheusClient.Period = period
	return nil
}

func (c *Controller) runCronTask(nodes map[string]log.Node) {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		accessToken := log.GetAccessToken(log.CorpId, log.Secret)
		tmpl, err := template.ParseFiles(log.Template)
		if err != nil {
			klog.Warning("create template failed, err:")
			return
		}
		buf := new(bytes.Buffer)
		batchNodes(nodes, c.nodes)
		tmpl.Execute(buf, nodes)
		c.prometheusClient.SamplingTimes = 0
		content := buf.String()
		msg := log.Messages("", "", log.TagId, log.AgentId, content)
		fmt.Println(accessToken,msg)
		//log.SendMessage(accessToken, msg)
	}
	//crontab.AddFunc("0 0 18 * * ?", task)
	crontab.AddFunc("0 */5 * * * *", task)
	crontab.Start()
	defer crontab.Stop()
	select {}
}

func batchNodes(nodes map[string]log.Node, corev1Nodes map[string]corev1.Node) {
	for key, value := range nodes {
		allocatable := corev1Nodes[key].Status.Allocatable
		cpuAllocatable := allocatable.Cpu().Value()
		memAllocatable := allocatable.Memory().Value()
		value.Allocatable.Memory = float64(memAllocatable)
		value.Allocatable.Cpu = float64(cpuAllocatable)
		ft1 := fmt.Sprintf("%.2f",value.CpuSumMax - value.CpuSumMin)
		cpuVolatility, err := strconv.ParseFloat(ft1, 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		value.CpuVolatility = cpuVolatility
		ft2 := fmt.Sprintf("%.4f", (value.MemMax-value.MemMin)/(value.Allocatable.Memory/math.Pow(2, 30)))
		memVolatility, err := strconv.ParseFloat(ft2, 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		value.MemVolatility = memVolatility*100
		nodes[key] = value
	}
	fmt.Println(nodes)
}
