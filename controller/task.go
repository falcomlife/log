package controller

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog-controller/db"
	"k8s.io/klog-controller/log"
	"k8s.io/klog/v2"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

// Generate daily
func (c *Controller) runCronTask(nodes *sync.Map, pods map[string]*log.Pod) {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		batchNodes(nodes, c.nodes)
		batchPods(pods)
		//analysis(c)
		//log.Draw(c.NodeCpuAnalysis)
		c.prometheusClient.SamplingTimes = 0
		msg := log.Messages("", "", log.TagId, log.AgentId, "今日日报已生成，请访问"+WebUrl+"查看")
		sendMessageToWechat(msg)
		event(c)
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

// Refresh top one percent value daily
// Deprecated: Use AI instance of this function
//func analysis(c *Controller) error {
//	cpu, err := log.AnalysisCpu(c.prometheusClient.Protocol, c.prometheusClient.Host, c.prometheusClient.Port)
//	mem, err := log.AnalysisMemory(c.prometheusClient.Protocol, c.prometheusClient.Host, c.prometheusClient.Port)
//	if err != nil {
//		return err
//	}
//	for name, node := range cpu {
//		node.ExtremePointMedian = common.Median(node.GetMaximumPoint())
//		cpu[name] = node
//	}
//	for name, node := range mem {
//		node.ExtremePointMedian = common.Median(node.GetMaximumPoint())
//		mem[name] = node
//	}
//	c.NodeCpuAnalysis = cpu
//	c.NodeMemoryAnalysis = mem
//	return nil
//}

// Count the nodes information of the day
func batchNodes(nodes *sync.Map, corev1Nodes map[string]corev1.Node) {
	nodes.Range(func(keyOri, valueOri interface{}) bool {
		key := keyOri.(string)
		value := valueOri.(log.Node)
		allocatable := corev1Nodes[key].Status.Allocatable
		cpuAllocatable := allocatable.Cpu().Value()
		memAllocatable := allocatable.Memory().Value()
		value.Allocatable.Memory = float64(memAllocatable)
		value.Allocatable.Cpu = float64(cpuAllocatable)
		ft1 := fmt.Sprintf("%.2f", value.CpuSumMax-value.CpuSumMin)
		cpuVolatility, err := strconv.ParseFloat(ft1, 64)
		if err != nil {
			klog.Warning(err)
			return false
		}
		value.CpuVolatility = cpuVolatility
		ft2 := fmt.Sprintf("%.2f", 100*(value.MemMax-value.MemMin)/(value.Allocatable.Memory/math.Pow(2, 30)))
		memVolatility, err := strconv.ParseFloat(ft2, 64)
		if err != nil {
			klog.Warning(err)
			return false
		}
		value.MemVolatility = memVolatility
		diskratio, err := strconv.ParseFloat(fmt.Sprintf("%.2f", 100*(value.DiskUsed/value.DiskTotal)), 64)
		value.DiskUsedRatio = diskratio
		nodes.Store(key, value)
		return true
	})
}

// Count the pods information of the day
func batchPods(pods map[string]*log.Pod) {
	for key, value := range pods {
		cpu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value.CpuSumMax-value.CpuSumMin), 64)
		value.CpuVolatility = cpu
		mem, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value.MemMax-value.MemMin), 64)
		value.MemVolatility = mem
		pods[key] = value
	}
}

// Generate record in event
func event(c *Controller) {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return
	}
	if key, ok := obj.(string); !ok {
		return
	} else {
		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		log, err := c.logsLister.Klogs(namespace).Get(name)
		if err != nil {
			klog.Error(err)
			return
		}
		c.recorder.Event(log, corev1.EventTypeNormal, SuccessSended, MessageResourceSended)
	}
	defer c.workqueue.Done(obj)
}

func (c *Controller) runWarningCronTask() {
	crontab := cron.New(cron.WithSeconds())
	//cpu_sendflag := true
	//memory_sendflag := true
	indexLaster := make(map[int64]float64)
	task := func() {
		defer mutex.Unlock()
		mutex.Lock()
		// TODO get AI result from db and add warning to list
		results := db.GetResult(c.nodes)
		for _, result := range results {
			margemap := make(map[string]interface{})
			json.Unmarshal([]byte(result.MergeIndex), &margemap)
			indexes := margemap["index"].([]interface{})
			datas := margemap["data"].([]interface{})
			for i, index := range indexes {
				int64Index := int64(index.(float64))
				if _, exsit := indexLaster[int64Index]; !exsit {
					data := datas[i].([]interface{})
					float64Data := data[0].(float64)
					c.warning(result.Instance, "Cpu峰值", "Cpu占用达到近期高点", make([]log.Pod, 0), float64Data, 0, time.Unix(int64Index/1000, 0), result.Instance+"使用量达到相对峰值")
					indexLaster[int64Index] = float64Data
				}
			}
		}
		//c.PrometheusMetricQueue.Range(func(nameNodeOri, nodeOri interface{}) bool {
		//	nameNode := nameNodeOri.(string)
		//	node := nodeOri.(log.Node)
		//	for nameSample, nodeSample := range c.NodeCpuAnalysis {
		//		if nameNode == nameSample {
		//			cpuExtremePointMedian := nodeSample.ExtremePointMedian * 100
		//			cpuValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuExtremePointMedian), 64)
		//			if node.CpuLaster > cpuValue && node.CpuLaster > float64(c.warningSetting.ExtremePointMedian.Cpu.WarningValue) && cpu_sendflag {
		//				cl := strconv.FormatFloat(node.CpuLaster, 'f', -1, 64)
		//				cv := strconv.FormatFloat(cpuValue, 'f', -1, 64)
		//				c.warning(nameNode, "Cpu峰值", "Cpu占用达到近期高点", make([]log.Pod, 0), cpuValue, node.CpuLaster, time.Now(), nameNode+"，过去一天cpu使用峰值中值为"+cv+"%,当前使用量为"+cl+"%")
		//				cpu_sendflag = false
		//			} else if node.CpuLaster <= cpuValue && !cpu_sendflag {
		//				cpu_sendflag = true
		//			}
		//		}
		//	}
		//	for nameSample, nodeSample := range c.NodeMemoryAnalysis {
		//		if nameNode == nameSample {
		//			extremePointMedian := nodeSample.ExtremePointMedian / math.Pow(2, 30)
		//			allocatable := c.nodes[nameNode].Status.Allocatable
		//			memoryAllocatable := float64(allocatable.Memory().Value()) / math.Pow(2, 30)
		//			memoryWarningValue := (float64(c.warningSetting.ExtremePointMedian.Memory.WarningValue) / 100) * float64(memoryAllocatable)
		//			if node.MemLaster > extremePointMedian && node.MemLaster > memoryWarningValue && memory_sendflag {
		//				cl := strconv.FormatFloat(node.MemLaster, 'f', -1, 64)
		//				cv := fmt.Sprintf("%.2f", extremePointMedian)
		//				c.warning(nameNode, "内存峰值", "内存占用达到近期高点", make([]log.Pod, 0), extremePointMedian, node.CpuLaster, time.Now(), nameNode+"，过去一天内存使用峰值中值为"+cv+"Gi,当前使用量为"+cl+"Gi")
		//				memory_sendflag = false
		//			} else if node.MemLaster <= extremePointMedian && memory_sendflag {
		//				memory_sendflag = true
		//			}
		//		}
		//	}
		//	return true
		//})
	}
	crontab.AddFunc("*/10 * * * * *", task)
	crontab.Start()
	defer crontab.Stop()
	select {}
}

// Clear data of the day, run after *Controller.runCronTask task
func (c *Controller) runCleanCronTask() {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		queueOrigin := *c.PrometheusMetricQueue
		queue := &queueOrigin
		queue.Range(func(key interface{}, value interface{}) bool {
			c.PrometheusMetricQueue.Delete(key)
			return true
		})
		c.Warnings = make([]*log.WarningList, 0)
	}
	if os.Getenv("LOG.ENV") == "PROD" {
		crontab.AddFunc("0 0 23 * * ?", task)
	} else {
		crontab.AddFunc("0 */10 * * * *", task)
	}
	crontab.Start()
	defer crontab.Stop()
	select {}
}

// Send message to wechat
func sendMessageToWechat(msg string) {
	accessToken := log.GetAccessToken(log.CorpId, log.Secret)
	if os.Getenv("LOG.ENV") == "PROD" {
		log.SendMessage(accessToken, msg)
	} else {
		fmt.Println(accessToken, msg)
	}
}
