package log

import (
	"fmt"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/common"
	"strconv"
	"time"
)

func (p *PrometheusClient) Get(path string) (map[string]Node, error) {
	httpclient := common.HttpClient{Host: p.Host, Port: p.Port}
	m, err := httpclient.Get(path)
	if err != nil {
		return nil, err
	}
	return emetricMapToNode(m)
}

func emetricMapToNode(m map[string]interface{}) (map[string]Node, error) {
	rm := make(map[string]Node)
	status, _ := m["status"].(string)
	if status != "success" {
		klog.Error("get a error message from server, status is " + status)
		return nil, fmt.Errorf("get a error message from server, status is %v", status)
	}
	data, ok := m["data"].(map[string]interface{})
	if !ok {
		klog.Warning("get data fail")
		return nil, fmt.Errorf("get data fail")
	}
	result, ok := data["result"].([]interface{})
	if !ok {
		klog.Warning("get result fail")
		return nil, fmt.Errorf("get result fail")
	}
	for _, r := range result {
		rr := r.(map[string]interface{})
		metric, ok := rr["metric"].(map[string]interface{})
		if !ok {
			klog.Warning("get metric fail")
			continue
		}
		value, ok := rr["value"].([]interface{})
		if !ok {
			klog.Warning("get value fail")
			continue
		}
		instance := metric["instance"].(string)
		var inst Node
		if _, ok := rm[instance]; ok {
			// instance exist in map
			inst = (rm[instance])
		} else {
			// instance not exist in map
			inst = Node{
				Name: instance,
				Cpu:  make(map[string]Cpu),
			}
		}
		v, ok := (value[1]).(string)
		if !ok {
			klog.Warning("value[1] not a string")
			continue
		}
		cpuValueStr, err := strconv.ParseFloat(v, 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		cpu := Cpu{
			CpuMax:     cpuValueStr,
			CpuMaxTime: time.Now(),
		}
		inst.Cpu[(metric["cpu"].(string))] = cpu
		rm[inst.Name] = inst
	}
	return rm, nil
}