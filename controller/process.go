package controller

import (
	"fmt"
	cron "github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/common"
	"k8s.io/log-controller/log"
	logv1alpha1 "k8s.io/log-controller/pkg/apis/logcontroller/v1alpha1"
	"math"
	"os"
	"strconv"
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

	// Launch two workers to process Log resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	go c.runPrometheusWorker()
	go c.runCronTask(c.PrometheusMetricQueue, c.PrometheusPodMetricQueue)
	go c.runCleanCronTask()
	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// run a work for get prometheus
func (c *Controller) runPrometheusWorker() {
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
				nodes1, err := c.prometheusClient.GetNode(log.NodeCpuUsedPercentage)
				if err != nil {
					continue
				}
				c.batchForCpu(nodes1)
				nodes2, err := c.prometheusClient.GetNode(log.NodeMemoryUsed)
				if err != nil {
					continue
				}
				c.batchForMem(nodes2)
				nodes3, err := c.prometheusClient.GetNode(log.NodeDiskUsed)
				if err != nil {
					continue
				}
				nodes4, err := c.prometheusClient.GetNode(log.NodeDiskTotal)
				if err != nil {
					continue
				}
				c.batchForDisk(nodes3, nodes4)
				pod1, err := c.prometheusClient.GetPod(log.PodCpuUsed)
				if err != nil {
					continue
				}
				c.batchForPodCpu(pod1)
				pod2, err := c.prometheusClient.GetPod(log.PodMemoryUsed)
				if err != nil {
					continue
				}
				c.batchForPodMem(pod2)
				c.prometheusClient.SamplingTimes++
			} else {
				continue
			}
			c.warningForCpu()
		}
	}
}

func (c *Controller) warningForCpu() {
	for nameNode, node := range c.PrometheusMetricQueue {
		for nameSample, nodeSample := range c.NodeCpuAnalysis {
			if nameNode == nameSample {
				cpuExtremePointMedian := nodeSample.ExtremePointMedian * 100
				cpuValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuExtremePointMedian), 64)
				if node.CpuLaster > cpuValue {
					cl := strconv.FormatFloat(node.CpuLaster, 'f', -1, 64)
					cv := strconv.FormatFloat(cpuValue, 'f', -1, 64)
					msg := nameNode + "cpu occupy is high now is " + cl + "% and warning value is " + cv + "%"
					accessToken := log.GetAccessToken(log.CorpId, log.Secret)
					if os.Getenv("LOG.ENV") == "PROD" {
						log.SendMessage(accessToken, msg)
					} else {
						fmt.Println(msg)
					}
				}
			}
		}
	}
}

func (c *Controller) batchForCpu(nodes map[string]log.Node) {
	for _, node := range nodes {
		var nodeOld log.Node
		if _, ok := c.PrometheusMetricQueue[node.Name]; ok {
			// there is node in map
			nodeOld = c.PrometheusMetricQueue[node.Name]
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
		cpuValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuSumNew), 64)
		if nodeOld.CpuSumMax <= cpuSumNew {
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.CpuSumMax = cpuValue
			nodeOld.CpuSumMaxTime = time.Now()
		}
		if nodeOld.CpuSumMin >= cpuSumNew {
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.CpuSumMin = cpuValue
			nodeOld.CpuSumMinTime = time.Now()
		}
		if ratio := cpuValue - nodeOld.CpuLaster; ratio > nodeOld.CpuMaxRatio && nodeOld.CpuLaster != 0 {
			r, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", ratio), 64)
			nodeOld.CpuMaxRatio = r
		}
		nodeOld.CpuLaster = cpuValue

		cpuAvgtemp := (nodeOld.CpuSumAvg*float64(c.prometheusClient.SamplingTimes) + cpuSumNew) / float64(c.prometheusClient.SamplingTimes+1)
		cpuAvg, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuAvgtemp), 64)
		if err != nil {
			klog.Warning(err)
		}
		nodeOld.CpuSumAvg = cpuAvg
		c.PrometheusMetricQueue[node.Name] = nodeOld
	}
}

func (c *Controller) batchForMem(nodes map[string]log.Node) {
	for _, node := range nodes {
		var nodeOld log.Node
		if _, ok := c.PrometheusMetricQueue[node.Name]; ok {
			// there is node in map
			nodeOld = c.PrometheusMetricQueue[node.Name]
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
		memValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memUsedNew), 64)
		if nodeOld.MemMax <= memUsedNew {
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.MemMax = memValue
			nodeOld.MemMaxTime = time.Now()
		}
		if nodeOld.MemMin >= memUsedNew {
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.MemMin = memValue
			nodeOld.MemMinTime = time.Now()
		}
		ratio := memValue - nodeOld.MemLaster
		allocatable := c.nodes[nodeOld.Name].Status.Allocatable
		ft2 := fmt.Sprintf("%.2f", 100*ratio/(float64(allocatable.Memory().Value())/math.Pow(2, 30)))
		memratio, _ := strconv.ParseFloat(ft2, 64)
		if memratio > nodeOld.MemMaxRatio && nodeOld.MemLaster != 0 {
			nodeOld.MemMaxRatio = memratio
		}
		nodeOld.MemLaster = memValue

		memAvgtemp := (nodeOld.MemAvg*float64(c.prometheusClient.SamplingTimes) + memUsedNew) / float64(c.prometheusClient.SamplingTimes+1)
		memAvg, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memAvgtemp), 64)
		if err != nil {
			klog.Warning(err)
		}
		nodeOld.MemAvg = memAvg
		c.PrometheusMetricQueue[node.Name] = nodeOld
	}
}

func (c *Controller) batchForDisk(nodes3 map[string]log.Node, nodes4 map[string]log.Node) {
	for _, node := range nodes3 {
		var nodeOld log.Node
		if _, ok := c.PrometheusMetricQueue[node.Name]; ok {
			// there is node in map
			nodeOld = c.PrometheusMetricQueue[node.Name]
		} else {
			// there is not node in map
			nodeOld = log.Node{
				Name:      node.Name,
				Cpu:       make(map[string]log.Cpu),
				CpuSumMin: math.MaxFloat64,
				MemMin:    math.MaxFloat64,
			}
		}
		diskused, err := strconv.ParseFloat(fmt.Sprintf("%.2f", node.DiskUsed/math.Pow(2, 30)), 64)
		if err != nil {
			klog.Error(err)
		}
		nodeOld.DiskUsed = diskused
		c.PrometheusMetricQueue[node.Name] = nodeOld
	}
	for _, node := range nodes4 {
		var nodeOld log.Node
		if _, ok := c.PrometheusMetricQueue[node.Name]; ok {
			// there is node in map
			nodeOld = c.PrometheusMetricQueue[node.Name]
		} else {
			// there is not node in map
			nodeOld = log.Node{
				Name:      node.Name,
				Cpu:       make(map[string]log.Cpu),
				CpuSumMin: math.MaxFloat64,
				MemMin:    math.MaxFloat64,
			}
		}
		disktotal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", node.DiskTotal/math.Pow(2, 30)), 64)
		if err != nil {
			klog.Error(err)
		}
		nodeOld.DiskTotal = disktotal
		c.PrometheusMetricQueue[node.Name] = nodeOld
	}
}

func (c *Controller) batchForPodCpu(pods map[string]log.Pod) {
	for _, pod := range pods {
		var podOld log.Pod
		if _, ok := c.PrometheusPodMetricQueue[pod.Name]; ok {
			// there is node in map
			podOld = c.PrometheusPodMetricQueue[pod.Name]
		} else {
			// there is not node in map
			podOld = log.Pod{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				CpuSumMin: math.MaxFloat64,
				MemMin:    math.MaxFloat64,
			}
		}
		value, err := strconv.ParseFloat(fmt.Sprintf("%.2f", pod.CpuSumMax*1000), 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		if podOld.CpuSumMax <= value {
			podOld.CpuSumMax = value
			podOld.CpuSumMaxTime = pod.CpuSumMaxTime
		}
		if podOld.CpuSumMin >= value {
			podOld.CpuSumMin = value
			podOld.CpuSumMinTime = pod.CpuSumMinTime
		}
		if ratio := value - podOld.CpuLaster; ratio > podOld.CpuMaxRatio && podOld.CpuLaster != 0 {
			r, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", ratio), 64)
			podOld.CpuMaxRatio = r
		}
		podOld.CpuLaster = value

		cpuAvgtemp := (podOld.CpuSumAvg*float64(c.prometheusClient.SamplingTimes) + value) / float64(c.prometheusClient.SamplingTimes+1)
		cpuAvg, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuAvgtemp), 64)
		if err != nil {
			klog.Warning(err)
		}
		podOld.CpuSumAvg = cpuAvg
		c.PrometheusPodMetricQueue[pod.Name] = podOld
	}
}

func (c *Controller) batchForPodMem(pods map[string]log.Pod) {
	for _, pod := range pods {
		var podOld log.Pod
		if _, ok := c.PrometheusPodMetricQueue[pod.Name]; ok {
			// there is pod in map
			podOld = c.PrometheusPodMetricQueue[pod.Name]
		} else {
			// there is not pod in map
			podOld = log.Pod{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				CpuSumMin: math.MaxFloat64,
				MemMin:    math.MaxFloat64,
			}
		}
		memUsedNew := pod.MemMax / math.Pow(2, 20)
		memValue, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memUsedNew), 64)
		if podOld.MemMax <= memValue {
			if err != nil {
				klog.Warning(err)
				continue
			}
			podOld.MemMax = memValue
			podOld.MemMaxTime = pod.MemMaxTime
		}
		if podOld.MemMin >= memValue {
			if err != nil {
				klog.Warning(err)
				continue
			}
			podOld.MemMin = memValue
			podOld.MemMinTime = pod.MemMinTime
		}
		ratio, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", memValue-podOld.MemLaster), 64)
		if ratio > podOld.MemMaxRatio && podOld.MemLaster != 0 {
			podOld.MemMaxRatio = ratio
		}
		podOld.MemLaster = memValue

		memAvgtemp := (podOld.MemAvg*float64(c.prometheusClient.SamplingTimes) + memUsedNew) / float64(c.prometheusClient.SamplingTimes+1)
		memAvg, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memAvgtemp), 64)
		if err != nil {
			klog.Warning(err)
		}
		podOld.MemAvg = memAvg
		c.PrometheusPodMetricQueue[pod.Name] = podOld
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

func (c *Controller) runCronTask(nodes map[string]log.Node, pods map[string]log.Pod) {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		accessToken := log.GetAccessToken(log.CorpId, log.Secret)
		batchNodes(nodes, c.nodes)
		batchPods(pods)
		analysis(c)
		//log.Draw(c.NodeCpuAnalysis)
		c.prometheusClient.SamplingTimes = 0
		msg := log.Messages("", "", log.TagId, log.AgentId, "今日日报已生成，请访问"+WebUrl+"查看")
		if os.Getenv("LOG.ENV") == "PROD" {
			log.SendMessage(accessToken, msg)
		} else {
			fmt.Println(accessToken, msg)
		}
		obj, shutdown := c.workqueue.Get()
		if shutdown {
			return
		}
		if key, ok := obj.(string); !ok {
			return
		} else {
			namespace, name, err := cache.SplitMetaNamespaceKey(key)
			log, err := c.logsLister.Logs(namespace).Get(name)
			if err != nil {
				klog.Error(err)
				return
			}
			c.recorder.Event(log, corev1.EventTypeNormal, SuccessSended, MessageResourceSended)
		}
		defer c.workqueue.Done(obj)
	}
	if os.Getenv("LOG.ENV") == "PROD" {
		crontab.AddFunc("0 0 20 * * ?", task)

	} else {
		crontab.AddFunc("0 */1 * * * *", task)
	}
	crontab.Start()
	defer crontab.Stop()
	select {}
}

func analysis(c *Controller) error {
	cpu, err := log.AnalysisCpu()
	mem, err := log.AnalysisMemory()
	if err != nil {
		return err
	}
	for name, node := range cpu {
		node.ExtremePointMedian = common.Median(node.GetMaximumPoint())
		cpu[name] = node
	}
	for name, node := range mem {
		node.ExtremePointMedian = common.Median(node.GetMaximumPoint())
		mem[name] = node
	}
	c.NodeCpuAnalysis = cpu
	c.NodeMemoryAnalysis = mem
	return nil
}

func (c *Controller) runCleanCronTask() {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		c.PrometheusMetricQueue = make(map[string]log.Node)
	}
	crontab.AddFunc("0 0 23 * * ?", task)
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
		ft1 := fmt.Sprintf("%.2f", value.CpuSumMax-value.CpuSumMin)
		cpuVolatility, err := strconv.ParseFloat(ft1, 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		value.CpuVolatility = cpuVolatility
		ft2 := fmt.Sprintf("%.2f", 100*(value.MemMax-value.MemMin)/(value.Allocatable.Memory/math.Pow(2, 30)))
		memVolatility, err := strconv.ParseFloat(ft2, 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		value.MemVolatility = memVolatility
		diskratio, err := strconv.ParseFloat(fmt.Sprintf("%.2f", 100*(value.DiskUsed/value.DiskTotal)), 64)
		value.DiskUsedRatio = diskratio
		nodes[key] = value
	}
}

func batchPods(pods map[string]log.Pod) {
	for key, value := range pods {
		cpu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value.CpuSumMax-value.CpuSumMin), 64)
		value.CpuVolatility = cpu
		mem, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value.MemMax-value.MemMin), 64)
		value.MemVolatility = mem
		pods[key] = value
	}
}
