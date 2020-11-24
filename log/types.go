package log

import (
	"k8s.io/log-controller/common"
	"time"
)

const (
	Period                = 30
	NodeCpuUsedPercentage = "/api/v1/query?query=(%0A%20%20(1%20-%20rate(node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D%5B75s%5D))%0A%2F%20ignoring(cpu)%20group_left%0A%20%20count%20without%20(cpu)(%20node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D)%0A)%0A"

	CorpId  = "wx46e1733df6787273"
	AgentId = 1000052
	TagId   = 1
	Secret  = "mPPjV5BKhe6xpU414WiUxcMpf4_8P_zj9qsxHk8XNWc"
)

type Log struct {
}

type Node struct {
	Name       string
	Cpu        map[string]Cpu
	CpuSum     float64
	CpuSumTime time.Time
	MemMax     float64
	MemMaxTime time.Time
	MemAvg     float64
	Amplitude  float64
}

type Cpu struct {
	CpuMax     float64
	CpuMaxTime time.Time
	CpuAvg     float64
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
	Period int64
}
