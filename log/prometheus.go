package log

import (
	"fmt"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/common"
	"strconv"
	"time"
)

func (p *PrometheusClient) GetNode(path string) (map[string]Node, error) {
	httpclient := common.HttpClient{Protocol: p.Protocol, Host: p.Host, Port: p.Port}
	m := make(map[string]interface{})
	err := httpclient.Get(path, &m)
	if err != nil {
		return nil, err
	}
	return emetricMapToNode(m, path)
}

func emetricMapToNode(m map[string]interface{}, path string) (map[string]Node, error) {
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
		valueFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		if path == NodeCpuUsedPercentage {
			cpu := Cpu{
				Value: valueFloat,
				Time:  time.Now(),
			}
			inst.Cpu[(metric["cpu"].(string))] = cpu
			rm[inst.Name] = inst
		} else if path == NodeMemoryUsed {
			inst.MemMax = valueFloat
			inst.MemMaxTime = time.Now()
			inst.MemMin = valueFloat
			inst.MemMinTime = time.Now()
			inst.MemAvg = valueFloat
			rm[inst.Name] = inst
		} else if path == NodeDiskUsed {
			inst.DiskTotal = valueFloat
			inst.DiskUsed = valueFloat
			rm[inst.Name] = inst
		} else if path == NodeDiskTotal {
			inst.DiskTotal = valueFloat
			inst.DiskUsed = valueFloat
			rm[inst.Name] = inst
		}
	}
	return rm, nil
}

func (p *PrometheusClient) GetPod(path string) (map[string]Pod, error) {
	httpclient := common.HttpClient{Protocol: p.Protocol, Host: p.Host, Port: p.Port}
	m := make(map[string]interface{})
	err := httpclient.Get(path, &m)
	if err != nil {
		return nil, err
	}
	return emetricMapToPod(m, path)
}

func emetricMapToPod(m map[string]interface{}, path string) (map[string]Pod, error) {
	rm := make(map[string]Pod)
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
		if len(metric) == 0 {
			continue
		}
		instance := metric["pod"].(string)
		ns := metric["namespace"].(string)
		var inst Pod
		if _, ok := rm[instance]; ok {
			// instance exist in map
			inst = (rm[instance])
		} else {
			// instance not exist in map
			inst = Pod{
				Name:      instance,
				Namespace: ns,
			}
		}
		v, ok := (value[1]).(string)
		if !ok {
			klog.Warning("value[1] not a string")
			continue
		}
		valueFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			klog.Warning(err)
			continue
		}
		inst.CpuSumMax = valueFloat
		inst.CpuSumMaxTime = time.Now()
		inst.CpuSumMin = valueFloat
		inst.CpuSumMinTime = time.Now()
		inst.CpuSumAvg = valueFloat
		rm[inst.Name] = inst
		inst.MemMax = valueFloat
		inst.MemMaxTime = time.Now()
		inst.MemMin = valueFloat
		inst.MemMinTime = time.Now()
		inst.MemAvg = valueFloat
		rm[inst.Name] = inst
	}
	return rm, nil
}
