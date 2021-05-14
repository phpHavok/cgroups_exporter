package collectors

import (
	"strconv"

	cg "github.com/phpHavok/cgroups_exporter/cgroups"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type cgroupsFileCollector struct {
	cpuacctUsagePerCPUMetric *prometheus.Desc
	memoryUsageInBytesMetric *prometheus.Desc
	cgroupFilePath           string
	cgroupsRootPath          string
}

func NewCgroupsFileCollector(cgroupFilePath string, cgroupsRootPath string) *cgroupsFileCollector {
	return &cgroupsFileCollector{
		cpuacctUsagePerCPUMetric: prometheus.NewDesc("cgroups_file_cpuacct_usage_per_cpu_ns",
			"Per-nanosecond usage of each CPU in a cgroup",
			[]string{"file_path", "cpu_id"}, nil,
		),
		memoryUsageInBytesMetric: prometheus.NewDesc("cgroups_file_memory_usage_in_bytes",
			"Current memory used by the cgroup in bytes",
			[]string{"file_path"}, nil,
		),
		cgroupFilePath:  cgroupFilePath,
		cgroupsRootPath: cgroupsRootPath,
	}
}

func (collector *cgroupsFileCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.cpuacctUsagePerCPUMetric
	ch <- collector.memoryUsageInBytesMetric
}

func (collector *cgroupsFileCollector) Collect(ch chan<- prometheus.Metric) {
	cgroups, err := cg.LoadCgroups(collector.cgroupFilePath, collector.cgroupsRootPath)
	if err != nil {
		log.Fatalf("unable to read cgroups specification file: %v", err)
	}
	// cpuacctUsagePerCPUMetric
	usagePerCPU, err := cgroups.Cpuacct.GetUsagePerCPU()
	if err != nil {
		log.Fatalf("unable to read cpuacct usage per cpu: %v", err)
	}
	for cpuID, cpuUsage := range usagePerCPU {
		ch <- prometheus.MustNewConstMetric(collector.cpuacctUsagePerCPUMetric,
			prometheus.GaugeValue, float64(cpuUsage), collector.cgroupFilePath, strconv.Itoa(cpuID))
	}
	// memoryUsageInBytesMetric
	memoryUsageBytes, err := cgroups.Memory.GetUsageInBytes()
	if err != nil {
		log.Fatalf("unable to read memory usage in bytes: %v", err)
	}
	ch <- prometheus.MustNewConstMetric(collector.memoryUsageInBytesMetric,
		prometheus.GaugeValue, float64(memoryUsageBytes), collector.cgroupFilePath)
}
