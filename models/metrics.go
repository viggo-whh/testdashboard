package models

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type MetricsCategory struct {
	Category  string `json:"category"`
	Nodes     string `json:"nodes,omitempty"`
	PVC       string `json:"pvc,omitempty"`
	Pods      string `json:"pods,omitempty"`
	Selector  string `json:"selector,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func (mc *MetricsCategory) GenerateQuery() *PrometheusQuery {
	switch mc.Category {
	case "cluster":
		return &PrometheusQuery{
			CpuUsage:          fmt.Sprintf(`sum(rate(node_cpu_seconds_total{instance=~"%s",mode=~"user|system"}[1m])) by (component)`, mc.Nodes),
			CpuRequests:       fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s",resource=~"cpu"}) by (component)`, mc.Nodes),
			CpuLimits:         fmt.Sprintf(`sum(kube_pod_container_resource_limitss{node=~"%s",resource=~"cpu"}) by (component)`, mc.Nodes),
			CpuCapacity:       fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s",resource="cpu"}) by (component)`, mc.Nodes),
			CpuAlloCatable:    fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s",resource="cpu"}) by (component)`, mc.Nodes),
			MemUsage:          strings.Replace("sum(node_memory_MemTotal_bytes-(node_memory_MemFree_bytes+node_memory_Cached_bytes+node_memory_Buffers_bytes)) by (instance)", "_bytes", fmt.Sprintf(`_bytes{instance=~"%s"}`, mc.Nodes), -1),
			MemRequests:       fmt.Sprintf(`sum(kube_pod_container_resource_requests{node=~"%s",resource=~"memory"}) by (component)`, mc.Nodes),
			MemLimits:         fmt.Sprintf(`sum(kube_pod_container_resource_limitss{node=~"%s",resource=~"memory"}) by (component)`, mc.Nodes),
			MemCapacity:       fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s",resource="memory"}) by (component)`, mc.Nodes),
			MemoryAlloCatable: fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s",resource="memory"}) by (component)`, mc.Nodes),
			FsUsage:           fmt.Sprintf(`sum(node_filesystem_size_bytes{instance=~"%s",mountpoint=~"/data"}-node_filesystem_avail_bytes{instance=~"%s",mountpoint=~"/data"}) by (instance)`, mc.Nodes, mc.Nodes),
			FsSize:            fmt.Sprintf(`sum(node_filesystem_size_bytes{instance=~"%s",mountpoint=~"/data"}) by (instance)`, mc.Nodes),
			PodUsage:          fmt.Sprintf(`sum(kubelet_running_pod_count{node=~"%s"}) by (component)`, mc.Nodes),
			PodCapacity:       fmt.Sprintf(`sum(kube_node_status_capacity{node=~"%s",resource="pods"}) by (component)`, mc.Nodes),
			PodAlloCatable:    fmt.Sprintf(`sum(kube_node_status_allocatable{node=~"%s",resource="pods"}) by (component)`, mc.Nodes),
		}
	case "nodes":
		return &PrometheusQuery{
			CpuUsage:          `sum(rate(node_cpu_seconds_total{mode=~"user|system"}[1m])) by (instance)`,
			CpuCapacity:       `sum(kube_node_status_capacity{resource="cpu"}) by (node)`,
			CpuAlloCatable:    `sum(kube_node_status_allocatable{resource="cpu"}) by (node)`,
			MemUsage:          `sum(node_memory_Active_bytes-(node_memory_MemFree_bytes+node_memory_Buffers_bytes+node_memory_Cached_bytes)) by (instance)`,
			MemCapacity:       `sum(kube_node_status_capacity{resource="memory"}) by (node)`,
			MemoryAlloCatable: `sum(kube_node_status_allocatable{resource="memory"}) by (node)`,
			FsUsage:           `sum(node_filesystem_size_bytes{mountpoint=~"/data"}-node_filesystem_avail_bytes{mountpoint=~"/data"}) by (instance)`,
			FsSize:            `sum(node_filesystem_size_bytes{mountpoint=~"/data"}) by (instance)`,
		}
	}
	return nil
}

type MetricsQuery struct {
	MemoryUsage       *MetricsCategory `json:"memoryUsage,omitempty"`
	MemoryRequests    *MetricsCategory `json:"memoryRequests,omitempty"`
	MemoryLimits      *MetricsCategory `json:"memoryLimits,omitempty"`
	MemoryCapacity    *MetricsCategory `json:"memoryCapacity,omitempty"`
	MemoryAlloCatable *MetricsCategory `json:"memoryAlloCatable,omitempty"`
	CpuUsages         *MetricsCategory `json:"cpuUsages,omitempty"`
	CpuLimits         *MetricsCategory `json:"cpuLimits,omitempty"`
	CpuRequests       *MetricsCategory `json:"cpuRequests,omitempty"`
	CpuCapacity       *MetricsCategory `json:"cpuCapacity,omitempty"`
	CpuAlloCatable    *MetricsCategory `json:"cpuAlloCatable,omitempty"`
	FsSize            *MetricsCategory `json:"fsSize,omitempty"`
	FsUsage           *MetricsCategory `json:"fsUsage,omitempty"`
	PodUsage          *MetricsCategory `json:"podUsage,omitempty"`
	PodCapacity       *MetricsCategory `json:"podCapacity,omitempty"`
	PodAlloCatable    *MetricsCategory `json:"podAlloCatable,,omitempty"`
}

type PrometheusQuery struct {
	CpuUsage          string
	CpuRequests       string
	CpuLimits         string
	CpuCapacity       string
	CpuAlloCatable    string
	MemUsage          string
	MemRequests       string
	MemLimits         string
	MemCapacity       string
	MemoryAlloCatable string
	FsUsage           string
	FsSize            string
	NetworkReceive    string
	NetworkTransmit   string
	PodUsage          string
	PodCapacity       string
	PodAlloCatable    string
	DiskUsage         string
	DiskCapacity      string
}

func (pq *PrometheusQuery) GetValueByField(field string) string {
	e := reflect.ValueOf(pq).Elem()
	for i := 0; i < e.NumField(); i++ {
		if e.Type().Field(i).Name == field {
			return e.Field(i).Interface().(string)
		}
	}
	return ""
}

type PrometheusQueryResp struct {
	Status string                   `json:"status"`
	Data   *PrometheusQueryRespData `json:"data"`
}

type PrometheusQueryRespData struct {
	ResultType string                      `json:"resultType"`
	Result     []PrometheusQueryRespResult `json:"result"`
}

type PrometheusQueryRespResult struct {
	Metric interface{} `json:"metric"`
	Values interface{} `json:"values"`
}

type PrometheusTracker struct {
	//添加读写锁
	sync.RWMutex
	Metrics map[string]*PrometheusQueryResp
}

func NewPrometheusTracker() *PrometheusTracker {
	return &PrometheusTracker{Metrics: map[string]*PrometheusQueryResp{}}
}

func (pt *PrometheusTracker) Get(key string) (*PrometheusQueryResp, bool) {
	pt.RLock()
	defer pt.RUnlock()
	value, exist := pt.Metrics[key]
	return value, exist
}

func (pt *PrometheusTracker) Set(key string, val *PrometheusQueryResp) {
	pt.Lock()
	defer pt.Unlock()
	pt.Metrics[key] = val
}
