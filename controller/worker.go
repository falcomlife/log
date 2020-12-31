package controller

import (
	"fmt"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/common"
	"k8s.io/log-controller/log"
	"math"
	"strconv"
	"sync"
	"time"
)

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
				now, period := common.GetRangeTime(c.warningSetting.Sustained.Cpu.Range)
				stepStr := strconv.FormatInt(c.warningSetting.Sustained.Cpu.Step, 10)
				nodes5, err := c.prometheusClient.GetNode(log.NodeCpuUsedPercentageSample + "&start=" + period + "&end=" + now + "&step=" + stepStr)
				if err != nil {
					continue
				}
				c.batchForNodeCpuSample(nodes5)
				now, period = common.GetRangeTime(c.warningSetting.Sustained.Memory.Range)
				stepStr = strconv.FormatInt(c.warningSetting.Sustained.Memory.Step, 10)
				nodes6, err := c.prometheusClient.GetNode(log.NodeMemUsedSample + "&start=" + period + "&end=" + now + "&step=" + stepStr)
				if err != nil {
					continue
				}
				c.batchForNodeMemSample(nodes6)
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
		}
	}
}

func (c *Controller) batchForCpu(nodes map[string]log.Node) {
	defer mutex.Unlock()
	mutex.Lock()
	for _, node := range nodes {
		var nodeOld = getOldNodes(c.PrometheusMetricQueue, node.Name)
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
			// update max cpu value
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.CpuSumMax = cpuValue
			nodeOld.CpuSumMaxTime = time.Now()
		}
		if nodeOld.CpuSumMin >= cpuSumNew {
			// update min cpu value
			if err != nil {
				klog.Warning(err)
				continue
			}
			nodeOld.CpuSumMin = cpuValue
			nodeOld.CpuSumMinTime = time.Now()
		}
		if ratio := cpuValue - nodeOld.CpuLaster; ratio > nodeOld.CpuMaxRatio && nodeOld.CpuLaster != 0 {
			// update cpu ratio
			r, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", ratio), 64)
			nodeOld.CpuMaxRatio = r
		}
		nodeOld.CpuLaster = cpuValue
		// update cpu avgrage value
		cpuAvgtemp := (nodeOld.CpuSumAvg*float64(c.prometheusClient.SamplingTimes) + cpuSumNew) / float64(c.prometheusClient.SamplingTimes+1)
		cpuAvg, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuAvgtemp), 64)
		if err != nil {
			klog.Warning(err)
		}
		nodeOld.CpuSumAvg = cpuAvg
		// update node in queue
		c.PrometheusMetricQueue.Store(node.Name, nodeOld)
	}
}

func (c *Controller) batchForMem(nodes map[string]log.Node) {
	defer mutex.Unlock()
	mutex.Lock()
	for _, node := range nodes {
		var nodeOld log.Node = getOldNodes(c.PrometheusMetricQueue, node.Name)
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
		c.PrometheusMetricQueue.Store(node.Name, nodeOld)
	}
}

func (c *Controller) batchForDisk(nodes3 map[string]log.Node, nodes4 map[string]log.Node) {
	defer mutex.Unlock()
	mutex.Lock()
	for _, node := range nodes3 {
		var nodeOld log.Node = getOldNodes(c.PrometheusMetricQueue, node.Name)
		diskused, err := strconv.ParseFloat(fmt.Sprintf("%.2f", node.DiskUsed/math.Pow(2, 30)), 64)
		if err != nil {
			klog.Error(err)
		}
		nodeOld.DiskUsed = diskused
		c.PrometheusMetricQueue.Store(node.Name, nodeOld)
	}
	for _, node := range nodes4 {
		var nodeOld log.Node = getOldNodes(c.PrometheusMetricQueue, node.Name)
		disktotal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", node.DiskTotal/math.Pow(2, 30)), 64)
		if err != nil {
			klog.Error(err)
		}
		nodeOld.DiskTotal = disktotal
		c.PrometheusMetricQueue.Store(node.Name, nodeOld)
	}
}

func (c *Controller) batchForNodeCpuSample(nodes5 map[string]log.Node) {
	warningValue := float64(c.warningSetting.Sustained.Cpu.WarningValue)
	warningValueStr := strconv.FormatFloat(warningValue, 'f', -1, 64)
	rangetmp := c.warningSetting.Sustained.Cpu.Range
	rangeStr := strconv.FormatInt(rangetmp, 10)
	for nodeName, node := range nodes5 {
		isHigh := true
		list := make([]log.Cpu, len(node.Cpu))
		for times, cpu := range node.Cpu {
			i, _ := strconv.Atoi(times)
			list[i] = cpu
		}
		for index, cpu := range list {
			if index != 0 && cpu.Value < node.Cpu[strconv.Itoa(index-1)].Value {
				isHigh = false
			}
			if index == len(node.Cpu)-1 {
				sub := cpu.Value - node.Cpu["0"].Value
				if sub <= 0 {
					continue
				}
				leftTime := (warningValue - cpu.Value) / (sub / float64(rangetmp))
				leftTimeStr := fmt.Sprintf("%.1f", leftTime)
				if leftTime <= 0 {
					leftTimeStr = "0"
				}
				if isHigh && leftTime < float64(c.warningSetting.Sustained.Cpu.LeftTime) {
					msg := log.Messages("", "", log.TagId, log.AgentId, nodeName+"主机"+rangeStr+"秒内cpu持续增长，将在"+leftTimeStr+"秒后超过"+warningValueStr+"%cpu使用时间")
					sendMessageToWechat(msg)
					actual, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpu.Value), 64)
					warning := &log.Warning{
						nodeName,
						"Cpu快速增长",
						rangeStr + "秒内cpu持续增长，" + leftTimeStr + "秒后超过" + warningValueStr + "%cpu使用时间",
						warningValue,
						actual,
						time.Now(),
					}
					c.addWarningMessage(nodeName, warning)
				}
			}
		}
	}
}

func (c *Controller) batchForNodeMemSample(nodes6 map[string]log.Node) {
	warningValue := float64(c.warningSetting.Sustained.Memory.WarningValue)
	warningValueStr := strconv.FormatFloat(warningValue, 'f', -1, 64)
	warningValueB := warningValue
	rangetmp := c.warningSetting.Sustained.Memory.Range
	rangeStr := strconv.FormatInt(rangetmp, 10)
	for nodeName, node := range nodes6 {
		isHigh := true
		list := make([]log.Memory, len(node.Memory))
		for times, mem := range node.Memory {
			i, _ := strconv.Atoi(times)
			list[i] = mem
		}
		for index, mem := range list {
			if index != 0 && mem.Value < node.Memory[strconv.Itoa(index-1)].Value {
				isHigh = false
			}
			if index == len(node.Memory)-1 {
				sub := mem.Value - node.Memory["0"].Value
				if sub <= 0 {
					continue
				}
				corev1Node := c.nodes[nodeName]
				warningValueMi := warningValueB / 100 * float64(corev1Node.Status.Allocatable.Memory().Value())
				leftTime := (warningValueMi - mem.Value) / (sub / float64(rangetmp))
				leftTimeStr := fmt.Sprintf("%.1f", leftTime)
				if leftTime <= 0 {
					leftTimeStr = "0"
				}
				if isHigh && (leftTime < 0 || leftTime < float64(c.warningSetting.Sustained.Memory.LeftTime)) {
					msg := log.Messages("", "", log.TagId, log.AgentId, nodeName+"主机"+rangeStr+"秒内内存持续增长，将在"+leftTimeStr+"秒后超过"+warningValueStr+"%内存总量")
					sendMessageToWechat(msg)
					actual, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", mem.Value), 64)
					warning := &log.Warning{
						nodeName,
						"内存快速增长",
						rangeStr + "秒内内存持续增长，" + leftTimeStr + "秒后超过" + warningValueStr + "%内存总量",
						warningValue,
						actual,
						time.Now(),
					}
					c.addWarningMessage(nodeName, warning)
				}
			}
		}
	}
}

func (c *Controller) batchForPodCpu(pods map[string]log.Pod) {
	defer mutex.Unlock()
	mutex.Lock()
	for _, pod := range pods {
		var podOld *log.Pod
		if _, ok := c.PrometheusPodMetricQueue[pod.Name]; ok {
			// there is node in map
			podOld = c.PrometheusPodMetricQueue[pod.Name]
		} else {
			// there is not node in map
			podOld = &log.Pod{
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
	defer mutex.Unlock()
	mutex.Lock()
	for _, pod := range pods {
		var podOld *log.Pod
		if _, ok := c.PrometheusPodMetricQueue[pod.Name]; ok {
			// there is pod in map
			podOld = c.PrometheusPodMetricQueue[pod.Name]
		} else {
			// there is not pod in map
			podOld = &log.Pod{
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

func getOldNodes(syncMap *sync.Map, name string) log.Node {
	var nodeOld log.Node
	if n, ok := syncMap.Load(name); ok {
		// there is node in map
		nodeOld = n.(log.Node)
	} else {
		// there is not node in map
		nodeOld = log.Node{
			Name:      name,
			Cpu:       make(map[string]log.Cpu),
			Memory:    make(map[string]log.Memory),
			CpuSumMin: math.MaxFloat64,
			MemMin:    math.MaxFloat64,
		}
	}
	return nodeOld
}
