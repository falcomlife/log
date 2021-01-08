package log

import (
	"time"
)

type Node struct {
	Name          string            `json:"name"`
	Cpu           map[string]Cpu    `json:"cpu"`
	Memory        map[string]Memory `json:"memory"`
	CpuSumMax     float64           `json:"cpuSumMax"`
	CpuSumMaxTime time.Time         `json:"cpuSumMaxTime"`
	CpuSumMin     float64           `json:"cpuSumMin"`
	CpuSumMinTime time.Time         `json:"cpuSumMinTime"`
	CpuSumAvg     float64           `json:"cpuSumAvg"`
	CpuVolatility float64           `json:"cpuVolatility"`
	CpuMaxRatio   float64           `json:"cpuMaxRatio"`
	CpuLaster     float64           `json:"cpuLaster"`
	MemMax        float64           `json:"memMax"`
	MemMaxTime    time.Time         `json:"memMaxTime"`
	MemMin        float64           `json:"memMin"`
	MemMinTime    time.Time         `json:"memMinTime"`
	MemAvg        float64           `json:"memAvg"`
	MemVolatility float64           `json:"memVolatility"`
	MemMaxRatio   float64           `json:"memMaxRatio"`
	MemLaster     float64           `json:"memLaster"`
	DiskUsed      float64           `json:"diskUsed"`
	DiskUsedRatio float64           `json:"diskUsedRatio"`
	DiskTotal     float64           `json:"diskTotal"`
	Amplitude     float64           `json:"amplitude"`
	Allocatable   Allocatable       `json:"allocatable"`
}

type Pod struct {
	Name          string      `json:"name"`
	Namespace     string      `json:"namespace"`
	Node          string      `json:"node"`
	CpuSumMax     float64     `json:"cpuSumMax"`
	CpuSumMaxTime time.Time   `json:"cpuSumMaxTime"`
	CpuSumMin     float64     `json:"cpuSumMin"`
	CpuSumMinTime time.Time   `json:"cpuSumMinTime"`
	CpuSumAvg     float64     `json:"cpuSumAvg"`
	CpuVolatility float64     `json:"cpuVolatility"`
	CpuMaxRatio   float64     `json:"cpuMaxRatio"`
	CpuLaster     float64     `json:"cpuLaster"`
	MemMax        float64     `json:"memMax"`
	MemMaxTime    time.Time   `json:"memMaxTime"`
	MemMin        float64     `json:"memMin"`
	MemMinTime    time.Time   `json:"memMinTime"`
	MemAvg        float64     `json:"memAvg"`
	MemVolatility float64     `json:"memVolatility"`
	MemMaxRatio   float64     `json:"memMaxRatio"`
	MemLaster     float64     `json:"memLaster"`
	Amplitude     float64     `json:"amplitude"`
	Allocatable   Allocatable `json:"allocatable"`
}

type WarningList struct {
	Name     string     `json:"name"` // 名称
	Children []*Warning `json:"children"`
}

type Warning struct {
	NodeName     string    `json:"nodeName"`     // 主机名
	Type         string    `json:"type"`         // 类型
	Describe     string    `json:"describe"`     // 描述
	WarningValue float64   `json:"warningValue"` // 警戒值
	Actual       float64   `json:"actual"`       // 实际值
	Time         time.Time `json:"time"`         // 时间
	Suspect      []Suspect `json:"suspect"`      // 造成问题的嫌疑人(pod)
}

type Suspect struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Namespace   string  `json:"namespace"`
	ActualValue float64 `json:"actualValue"`
	Volatility  float64 `json:"volatility"`
}
