package main

import (
	"strconv"

	cg "github.com/phpHavok/cgroups_exporter/cgroups"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type cgroupsCollector struct {
	cpuacctUsagePerCPUMetric *prometheus.Desc
	cgroupSpecPath           string
	cgroupsRootPath          string
}

func newCgroupsCollector(cgroupSpecPath string, cgroupsRootPath string) *cgroupsCollector {
	return &cgroupsCollector{
		cpuacctUsagePerCPUMetric: prometheus.NewDesc("cgroups_cpuacct_usage_per_cpu_ns",
			"Per-nanosecond usage of each CPU in a cgroup",
			[]string{"cpu_id"}, nil,
		),
		cgroupSpecPath:  cgroupSpecPath,
		cgroupsRootPath: cgroupsRootPath,
	}
}

func (collector *cgroupsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.cpuacctUsagePerCPUMetric
}

func (collector *cgroupsCollector) Collect(ch chan<- prometheus.Metric) {
	cgroups, err := cg.LoadCgroups(collector.cgroupSpecPath, collector.cgroupsRootPath)
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
			prometheus.GaugeValue, float64(cpuUsage), strconv.Itoa(cpuID))
	}
}
