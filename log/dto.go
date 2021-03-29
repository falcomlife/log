package log

import (
	"time"
)

type Node struct {
	Name          string            `json:"name"`          // 名称
	Cpu           map[string]Cpu    `json:"cpu"`           // cpu
	Memory        map[string]Memory `json:"memory"`        // 内存
	CpuSumMax     float64           `json:"cpuSumMax"`     // cpu各个核心总和最大值
	CpuSumMaxTime time.Time         `json:"cpuSumMaxTime"` // cpu各个核心总和最大值发生时间
	CpuSumMin     float64           `json:"cpuSumMin"`     // cpu各个核心总和最小值
	CpuSumMinTime time.Time         `json:"cpuSumMinTime"` // cpu各个核心总和最小值发生时间
	CpuSumAvg     float64           `json:"cpuSumAvg"`     // cpu各个核心平均值
	CpuVolatility float64           `json:"cpuVolatility"` // cpu震动幅度
	CpuMaxRatio   float64           `json:"cpuMaxRatio"`   // cpu值斜率最大
	CpuLaster     float64           `json:"cpuLaster"`     // 上次cpu记录值
	MemMax        float64           `json:"memMax"`        // 内存最大值
	MemMaxTime    time.Time         `json:"memMaxTime"`    // 内存最大值发生时间
	MemMin        float64           `json:"memMin"`        // 内存最小值
	MemMinTime    time.Time         `json:"memMinTime"`    // 内存最小值发生时间
	MemAvg        float64           `json:"memAvg"`        // 内存平均值
	MemVolatility float64           `json:"memVolatility"` // 内存震动幅度
	MemMaxRatio   float64           `json:"memMaxRatio"`   // 内存斜率最大值
	MemLaster     float64           `json:"memLaster"`     // 上次内存记录值
	DiskUsed      float64           `json:"diskUsed"`      // 磁盘使用值
	DiskUsedRatio float64           `json:"diskUsedRatio"` // 磁盘使用率
	DiskTotal     float64           `json:"diskTotal"`     // 磁盘总值
	DiskLeftTime  float64           `json:"diskLeftTime"`  // 磁盘剩余可用时间
	Amplitude     float64           `json:"amplitude"`     // 波动
	Allocatable   Allocatable       `json:"allocatable"`   // node的资源可分配值
}

type Pod struct {
	Name          string      `json:"name"`          // 名称
	Namespace     string      `json:"namespace"`     // 命名空间
	Node          string      `json:"node"`          // 节点名
	CpuSumMax     float64     `json:"cpuSumMax"`     // cpu各个核心总和最大值
	CpuSumMaxTime time.Time   `json:"cpuSumMaxTime"` // cpu各个核心总和最大值发生时间
	CpuSumMin     float64     `json:"cpuSumMin"`     // cpu各个核心总和最小值
	CpuSumMinTime time.Time   `json:"cpuSumMinTime"` // cpu各个核心总和最小值发生时间
	CpuSumAvg     float64     `json:"cpuSumAvg"`     // cpu各个核心平均值
	CpuVolatility float64     `json:"cpuVolatility"` // cpu震动幅度
	CpuMaxRatio   float64     `json:"cpuMaxRatio"`   // cpu值斜率最大
	CpuLaster     float64     `json:"cpuLaster"`     // 上次cpu记录值
	MemMax        float64     `json:"memMax"`        // 内存最大值
	MemMaxTime    time.Time   `json:"memMaxTime"`    // 内存最大值发生时间
	MemMin        float64     `json:"memMin"`        // 内存最小值
	MemMinTime    time.Time   `json:"memMinTime"`    // 内存最小值发生时间
	MemAvg        float64     `json:"memAvg"`        // 内存平均值
	MemVolatility float64     `json:"memVolatility"` // 内存震动幅度
	MemMaxRatio   float64     `json:"memMaxRatio"`   // 内存斜率最大值
	MemLaster     float64     `json:"memLaster"`     // 上次内存记录值
	Amplitude     float64     `json:"amplitude"`     // 波动
	Allocatable   Allocatable `json:"allocatable"`   // node的资源可分配值
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

type Deployment struct {
	Name        string `json:"name"`
	NameSpace   string `json:"nameSpace"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	GitUrl      string `json:"gitUrl"`
	GatewatUrl  string `json:"gatewayUrl"`
	InnerUrl    string `json:"innerUrl"`
	Replicas    int32  `json:"replicas"`
	Ready       int32  `json:"ready"`
}
