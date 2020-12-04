package log

import (
	"k8s.io/log-controller/common"
	"time"
)

const (
	Period                = 30
	NodeCpuUsedPercentage = "/api/v1/query?query=(%0A%20%20(1%20-%20rate(node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D%5B75s%5D))%0A%2F%20ignoring(cpu)%20group_left%0A%20%20count%20without%20(cpu)(%20node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D)%0A)%0A"
	NodeMemoryUsed        = "/api/v1/query?query=(%0A%20%20node_memory_MemTotal_bytes%7Bjob%3D%22node-exporter%22%7D%0A-%0A%20%20node_memory_MemFree_bytes%7Bjob%3D%22node-exporter%22%7D%0A-%0A%20%20node_memory_Buffers_bytes%7Bjob%3D%22node-exporter%22%7D%0A-%0A%20%20node_memory_Cached_bytes%7Bjob%3D%22node-exporter%22%7D%0A)"
	CorpId                = "wx46e1733df6787273"
	AgentId               = 1000052
	TagId                 = 1
	Secret                = "mPPjV5BKhe6xpU414WiUxcMpf4_8P_zj9qsxHk8XNWc"

	Template = "./message.template"
)

type Log struct {
}

type Node struct {
	Name          string
	Cpu           map[string]Cpu
	CpuSumMax     float64
	CpuSumMaxTime time.Time
	CpuSumMin     float64
	CpuSumMinTime time.Time
	CpuSumAvg     float64
	CpuVolatility float64
	CpuMaxRatio   float64
	CpuLaster     float64
	MemMax        float64
	MemMaxTime    time.Time
	MemMin        float64
	MemMinTime    time.Time
	MemAvg        float64
	MemVolatility float64
	MemMaxRatio   float64
	MemLaster     float64
	Amplitude     float64
	Allocatable   Allocatable
}

type Cpu struct {
	Value float64
	Time  time.Time
}

type Allocatable struct {
	Cpu    float64
	Memory float64
}

type Pod struct {
	Name       string
	Namespace  string
	CpuMax     float64
	CpuMaxTime time.Time
	MemMax     float64
	MemMaxTime time.Time
	CpuAvg     float64
	MemAvg     float64
	Amplitude  float64
}

type PrometheusClient struct {
	common.HttpClient
	SamplingTimes int64
	Period        int64
}
