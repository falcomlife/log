package controller

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"k8s.io/log-controller/log"
	"os"
	"strconv"
	"time"
)

func (c *Controller) runCronTask(nodes map[string]log.Node, pods map[string]log.Pod) {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		batchNodes(nodes, c.nodes)
		batchPods(pods)
		analysis(c)
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

func (c *Controller) runWarningCronTask() {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		defer mutex.Unlock()
		mutex.Lock()
		for nameNode, node := range c.PrometheusMetricQueue {
			for nameSample, nodeSample := range c.NodeCpuAnalysis {
				if nameNode == nameSample {
					cpuExtremePointMedian := nodeSample.ExtremePointMedian * 100
					cpuValue, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuExtremePointMedian), 64)
					if node.CpuLaster > cpuValue {
						cl := strconv.FormatFloat(node.CpuLaster, 'f', -1, 64)
						cv := strconv.FormatFloat(cpuValue, 'f', -1, 64)
						warning := &log.Warning{
							nameNode,
							"Cpu峰值",
							"Cpu占用达到近期高点",
							cpuValue,
							node.CpuLaster,
							time.Now(),
						}
						msg := log.Messages("", "", log.TagId, log.AgentId, nameNode+"，过去一天cpu使用峰值中值为"+cv+",当前使用量为"+cl+"%")
						sendMessageToWechat(msg)
						c.addMessage(nameNode, warning)
					}
				}
			}
		}
	}
	crontab.AddFunc("*/10 * * * * *", task)
	crontab.Start()
	defer crontab.Stop()
	select {}
}

func (c *Controller) runCleanCronTask() {
	crontab := cron.New(cron.WithSeconds())
	task := func() {
		c.PrometheusMetricQueue = make(map[string]log.Node)
		c.Warnings = make([]*log.WarningList, 0)
	}
	crontab.AddFunc("0 0 23 * * ?", task)
	crontab.Start()
	defer crontab.Stop()
	select {}
}
