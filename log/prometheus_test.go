package log

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/common"
	"strconv"
	"testing"
	"time"
)

func Test_prometheus_cpu(t *testing.T) {
	now, period := getTimestamp()
	m, err := getDataByInterface(NodeCpuUsedPercentageRange, strconv.FormatInt(period, 10), strconv.FormatInt(now, 10), "10")
	if err != nil {
		return
	}
	fmt.Println(time.Now())
	nodes, _ := batchCpuData(m)
	fmt.Println(time.Now())
	fmt.Println(nodes)
}

func Test_prometheus_mem(t *testing.T) {
	now, period := getTimestamp()
	m, err := getDataByInterface(NodeMemoryUsedRange, strconv.FormatInt(period, 10), strconv.FormatInt(now, 10), "10")
	if err != nil {
		return
	}
	fmt.Println(time.Now())
	nodes, _ := batchMemData(m)
	fmt.Println(time.Now())
	fmt.Println(nodes)
}

func batchCpuData(m map[string]interface{}) (map[string]*NodeSample, error) {
	nodes, err := prometheusDataToNodeSample(m, handlerCpu)
	if err != nil {
		return nil, err
	}
	for _, nodeSample := range nodes {
		var valuesTimes []float64
		var coreCount int = len(nodeSample.Cpu)
		for _, cpu := range nodeSample.Cpu {
			if valuesTimes == nil {
				valuesTimes = make([]float64, len(cpu))
			}
			for index, value := range cpu {
				valuesTimes[index] = valuesTimes[index] + value.Value
			}
		}
		for index, sumAllCore := range valuesTimes {
			valuesTimes[index] = sumAllCore / float64(coreCount)
		}
		indexArrUp, indexArrDown := getArea(valuesTimes)
		nodeSample.UpArea = indexArrUp
		nodeSample.DownArea = indexArrDown
	}
	return nodes, nil
}

func batchMemData(m map[string]interface{}) (map[string]*NodeSample, error) {
	nodes, err := prometheusDataToNodeSample(m, handlerMem)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func prometheusDataToNodeSample(m map[string]interface{}, f func(val []interface{}, nodeSample *NodeSample, name string) error) (map[string]*NodeSample, error) {
	status := m["status"]
	if status != "success" {
		klog.Error("Failed to get message from prometheus:" + m["error"].(string))
		return nil, errors.New("Failed to get message from prometheus:" + m["error"].(string))
	}
	data, _ := m["data"].(map[string]interface{})
	result, _ := data["result"].([]interface{})
	nodes := make(map[string]*NodeSample)
	for _, r := range result {
		rr := r.(map[string]interface{})
		metric, _ := rr["metric"].(map[string]interface{})
		valuesOrigin, _ := rr["values"].([]interface{})
		nodeName := metric["instance"].(string)
		nodeSample := nodes[nodeName]
		if nodeSample == nil {
			nodeSample = &NodeSample{
				NodeName: metric["instance"].(string),
				Cpu:      make(map[string][]*Cpu),
			}
		}
		f(valuesOrigin, nodeSample, nodeName)
		nodes[nodeName] = nodeSample
	}
	return nodes, nil
}

func handlerCpu(valuesOrigin []interface{}, nodeSample *NodeSample, name string) error {
	for _, val := range valuesOrigin {
		cpuSingleCore := nodeSample.Cpu[name]
		if cpuSingleCore == nil {
			cpuSingleCore = make([]*Cpu, 0)
		}
		valtran := val.([]interface{})
		s := valtran[1].(string)
		cpuSingleCoreValueTimes, err := strconv.ParseFloat(s, 64)
		if err != nil {
			klog.Error(err)
			return err
		}
		cpuSingleCore = append(cpuSingleCore, &Cpu{Value: cpuSingleCoreValueTimes})
		nodeSample.Cpu[name] = cpuSingleCore
	}
	return nil
}

func handlerMem(valuesOrigin []interface{}, nodeSample *NodeSample, name string) error {
	values := make([]float64, 0)
	for _, val := range valuesOrigin {
		valueStr, err := strconv.ParseFloat(val.([]interface{})[1].(string), 64)
		if err != nil {
			klog.Error(err)
			return err
		}
		values = append(values, valueStr)
	}
	indexArrUp, indexArrDown := getArea(values)
	nodeSample.UpArea = indexArrUp
	nodeSample.DownArea = indexArrDown
	return nil
}

func getDataByInterface(path string, start string, end string, step string) (map[string]interface{}, error) {
	p := path + "&start=" + start + "&end=" + end + "&step=" + step
	httpclient := common.HttpClient{Protocol: "http", Host: "172.16.77.154", Port: "39090"}
	m := make(map[string]interface{})
	error := httpclient.Get(p, &m)
	if error != nil {
		klog.Error(error)
		return nil, error
	}
	return m, nil
}

func getArea(values []float64) ([][2]int, [][2]int) {
	var indexArrUp = make([][2]int, 0)
	var indexArrDown = make([][2]int, 0)
	var last float64
	var lastMonotonicity int
	var lastChangeIndex int
	for index, v := range values {
		if index == 0 {
			last = v
		} else {
			if last > v {
				if lastMonotonicity == 1 {
					area := [2]int{lastChangeIndex, index}
					indexArrUp = append(indexArrUp, area)
					lastChangeIndex = index
				}
				if index == len(values)-1 && lastMonotonicity == -1 {
					area := [2]int{lastChangeIndex, index}
					indexArrDown = append(indexArrDown, area)
					break
				}
				lastMonotonicity = -1
			} else if last < v {
				if lastMonotonicity == -1 {
					area := [2]int{lastChangeIndex, index}
					indexArrDown = append(indexArrDown, area)
					lastChangeIndex = index
				}
				if index == len(values)-1 && lastMonotonicity == 1 {
					area := [2]int{lastChangeIndex, index}
					indexArrUp = append(indexArrUp, area)
					break
				}
				lastMonotonicity = 1
			} else {
				continue
			}
		}
		last = v
	}
	return indexArrUp, indexArrDown
}

func getTimestamp() (int64, int64) {
	timeUnix := time.Now().Unix()
	timeperiod := timeUnix - 60*60*24
	return timeUnix, timeperiod
}
