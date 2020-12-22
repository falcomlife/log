package log

import (
	"errors"
	"k8s.io/klog/v2"
	"k8s.io/log-controller/common"
	"strconv"
)

func AnalysisCpu(protocol, host, port string) (map[string]*NodeSample, error) {
	now, period := common.GetRangeTime(60 * 60 * 24)
	m, err := getDataByInterface(protocol, host, port, NodeCpuUsedPercentageRange, period, now, "10")
	if err != nil {
		return nil, err

	}
	nodes, _ := batchCpuData(m)
	return nodes, nil
}

func AnalysisMemory(protocol, host, port string) (map[string]*NodeSample, error) {
	now, period := common.GetRangeTime(60 * 60 * 24)
	m, err := getDataByInterface(protocol, host, port, NodeMemoryUsedRange, period, now, "10")
	if err != nil {
		return nil, err
	}
	nodes, _ := batchMemData(m)
	return nodes, nil
}

func batchCpuData(m map[string]interface{}) (map[string]*NodeSample, error) {
	nodes, err := prometheusDataToNodeSample(m, handlerCpu)
	if err != nil {
		return nil, err
	}
	for _, nodeSample := range nodes {
		var valuesTimes []float64
		//var coreCount int = len(nodeSample.Cpu)
		for _, cpu := range nodeSample.Cpu {
			if valuesTimes == nil {
				valuesTimes = make([]float64, len(cpu))
			}
			for index, value := range cpu {
				valuesTimes[index] = valuesTimes[index] + value.Value
			}
		}
		//for index, sumAllCore := range valuesTimes {
		//	valuesTimes[index] = sumAllCore / float64(coreCount)
		//}
		indexArrUp, indexArrDown := getArea(valuesTimes)
		nodeSample.UpArea = indexArrUp
		nodeSample.DownArea = indexArrDown
		nodeSample.SampleSumPoint = valuesTimes
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
		var name string = ""
		if metric["cpu"] != nil {
			name = metric["cpu"].(string)
		}
		f(valuesOrigin, nodeSample, name)
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
	nodeSample.SampleSumPoint = values
	return nil
}

func getDataByInterface(protocol, host, port, path, start, end, step string) (map[string]interface{}, error) {
	p := path + "&start=" + start + "&end=" + end + "&step=" + step
	httpclient := common.HttpClient{Protocol: protocol, Host: host, Port: port}
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
				if index == len(values)-1 && lastMonotonicity == -1 {
					area := [2]int{lastChangeIndex, index}
					indexArrDown = append(indexArrDown, area)
				} else if index == len(values)-1 && lastMonotonicity == 1 {
					areadown := [2]int{index - 1, index}
					indexArrDown = append(indexArrDown, areadown)
				}
				if lastMonotonicity == 1 {
					area := [2]int{lastChangeIndex, index - 1}
					indexArrUp = append(indexArrUp, area)
					lastChangeIndex = index - 1
				}
				lastMonotonicity = -1
			} else if last < v {
				if index == len(values)-1 && lastMonotonicity == 1 {
					area := [2]int{lastChangeIndex, index}
					indexArrUp = append(indexArrUp, area)
				} else if index == len(values)-1 && lastMonotonicity == -1 {
					areaup := [2]int{index - 1, index}
					indexArrUp = append(indexArrUp, areaup)
				}
				if lastMonotonicity == -1 {
					area := [2]int{lastChangeIndex, index - 1}
					indexArrDown = append(indexArrDown, area)
					lastChangeIndex = index - 1
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

func (n *NodeSample) GetMaximumPoint() []float64 {
	maxmumPoint := make([]float64, 0)
	for _, area := range n.UpArea {
		if 0 == area[0] || (len(n.UpArea)-1) == area[1] {
			continue
		}
		maxmumPoint = append(maxmumPoint, n.SampleSumPoint[area[1]])
	}
	return maxmumPoint
}

func (n *NodeSample) GetMinimumPoint() []float64 {
	minmumPoint := make([]float64, 0)
	for _, area := range n.DownArea {
		if 0 == area[0] || (len(n.DownArea)-1) == area[1] {
			continue
		}
		minmumPoint = append(minmumPoint, n.SampleSumPoint[area[1]])
	}
	return minmumPoint
}
