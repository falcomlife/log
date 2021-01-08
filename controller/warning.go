package controller

import (
	"fmt"
	"k8s.io/log-controller/common"
	"k8s.io/log-controller/log"
	"math"
	"strconv"
	"time"
)

func (c *Controller) warning(nodeName string, t string, describe string, pods []log.Pod, warningValue float64, actualValue float64, time time.Time, message string) {
	suspects := make([]log.Suspect, 0)
	for _, pod := range pods {
		if t == "Cpu快速增长" {
			suspects = append(suspects, log.Suspect{
				"PodCpu",
				pod.Name,
				pod.Namespace,
				common.EllipsisDot(pod.CpuLaster),
				common.EllipsisDot(pod.CpuVolatility),
			})
		} else if t == "内存快速增长" {
			mem1 := pod.MemLaster / math.Pow(2, 30)
			memLaster, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", mem1), 64)
			mem2 := pod.MemVolatility / math.Pow(2, 30)
			memVolatility, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", mem2), 64)
			suspects = append(suspects, log.Suspect{
				"PodMem",
				pod.Name,
				pod.Namespace,
				memLaster,
				memVolatility,
			})
		}
	}
	warning := &log.Warning{
		nodeName,
		t,
		describe,
		warningValue,
		actualValue,
		time,
		suspects,
	}
	msg := log.Messages("", "", log.TagId, log.AgentId, message)
	sendMessageToWechat(msg)
	c.addWarningMessage(nodeName, warning)
}

// Add a new message to warning list
func (c *Controller) addWarningMessage(nameNode string, warning *log.Warning) {
	exsit := false
	warngingtemp := make([]*log.WarningList, len(c.Warnings))
	copy(warngingtemp, c.Warnings)
	for i, nodeWarnings := range warngingtemp {
		if nodeWarnings.Name == nameNode {
			c.Warnings[i].Children = append(c.Warnings[i].Children, warning)
			exsit = true
		}
	}
	if !exsit {
		warninglist := log.WarningList{
			Name:     nameNode,
			Children: make([]*log.Warning, 0),
		}
		warninglist.Children = append(warninglist.Children, warning)
		c.Warnings = append(c.Warnings, &warninglist)
	}
}
