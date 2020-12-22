package log

import (
	"k8s.io/log-controller/common"
	"time"
)

const (
	Period                      = 30
	NodeCpuUsedPercentage       = "/api/v1/query?query=(%0A%20%20(1%20-%20rate(node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D%5B75s%5D))%0A%2F%20ignoring(cpu)%20group_left%0A%20%20count%20without%20(cpu)(%20node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D)%0A)%0A"
	NodeMemoryUsed              = "/api/v1/query?query=(%0A%20%20node_memory_MemTotal_bytes%7Bjob%3D%22node-exporter%22%7D%0A-%0A%20%20node_memory_MemFree_bytes%7Bjob%3D%22node-exporter%22%7D%0A-%0A%20%20node_memory_Buffers_bytes%7Bjob%3D%22node-exporter%22%7D%0A-%0A%20%20node_memory_Cached_bytes%7Bjob%3D%22node-exporter%22%7D%0A)"
	NodeCpuUsedPercentageRange  = "/api/v1/query_range?query=(%0A%20%20(1%20-%20rate(node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D%5B75s%5D))%0A%2F%20ignoring(cpu)%20group_left%0A%20%20count%20without%20(cpu)(%20node_cpu_seconds_total%7Bjob%3D\"node-exporter\"%2C%20mode%3D\"idle\"%7D)%0A)"
	NodeMemoryUsedRange         = "/api/v1/query_range?query=(node_memory_MemTotal_bytes%7bjob%3d%22node-exporter%22%7d-node_memory_MemFree_bytes%7bjob%3d%22node-exporter%22%7d-node_memory_Buffers_bytes%7bjob%3d%22node-exporter%22%7d-node_memory_Cached_bytes%7bjob%3d%22node-exporter%22%7d)"
	NodeCpuUsedPercentageSample = "/api/v1/query_range?query=(sum%20without(cpu)(60-increase(node_cpu_seconds_total%7Bjob=%22node-exporter%22,%20mode=%22idle%22%7D%5B60s%5D)))/(60*count%20without(cpu)%20(node_cpu_seconds_total%7Bjob=%22node-exporter%22,%20mode=%22idle%22%7D))*100"
	NodeDiskUsed                = "/api/v1/query?query=sum%20by%20(instance)(node_filesystem_size_bytes%7Bjob%3D%22node-exporter%22%2C%20fstype!%3D%22%22%2C%20fstype!%3D%22tmpfs%22%2C%20fstype!%3D%22rootfs%22%7D)%0A%20%20-%0Asum%20by%20(instance)(node_filesystem_avail_bytes%7Bjob%3D%22node-exporter%22%2C%20fstype!%3D%22%22%2C%20fstype!%3D%22tmpfs%22%2C%20fstype!%3D%22rootfs%22%7D)%0A"
	NodeDiskTotal               = "/api/v1/query?query=sum%20by(instance)(node_filesystem_size_bytes%7Bjob%3D%22node-exporter%22%2C%20fstype!%3D%22%22%2Cfstype!%3D%22rootfs%22%2Cfstype!%3D%22tmpfs%22%7D)"
	PodCpuUsed                  = "/api/v1/query?query=sum%20by%20(namespace%2Cpod)%20(irate(container_cpu_usage_seconds_total%7Bjob%3D%22kubelet%22%2C%20cluster%3D%22%22%2Cimage!%3D%22%22%2C%20container!%3D%22POD%22%7D%5B4m%5D))"
	PodMemoryUsed               = "/api/v1/query?query=sum%20by(namespace%2Cpod)%20(container_memory_working_set_bytes%7bjob%3d%22kubelet%22%2c+cluster%3d%22%22%2c+container!%3d%22POD%22%2c+container!%3d%22%22%7d)"
	CorpId                      = "wx46e1733df6787273"
	AgentId                     = 1000052
	TagId                       = 1
	Secret                      = "mPPjV5BKhe6xpU414WiUxcMpf4_8P_zj9qsxHk8XNWc"

	Template = "./message.template"
)

type Log struct {
}

type Cpu struct {
	Value float64
	Time  time.Time
}

type Allocatable struct {
	Cpu    float64
	Memory float64
}

type PrometheusClient struct {
	common.HttpClient
	SamplingTimes int64
	Period        int64
}

type NodeSample struct {
	NodeName           string
	Cpu                map[string][]*Cpu
	UpArea             [][2]int
	DownArea           [][2]int
	SampleSumPoint     []float64
	ExtremePointMedian float64 //极值点中值
}
