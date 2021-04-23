package cgroups

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type cpuset string

// GetCpus returns an array containing valid CPU numbers for this cgroup sorted
// by CPU index ascending
func (c cpuset) GetCpus() ([]int, error) {
	var cpus []int
	data, err := readFile(string(c), "cpuset.cpus")
	if err != nil {
		return cpus, err
	}

	for _, cpuRange := range strings.Split(data, ",") {
		bounds := strings.Split(cpuRange, "-")
		switch len(bounds) {
		case 1:
			single, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
			if err != nil {
				return cpus, err
			}
			cpus = append(cpus, single)
		case 2:
			lower, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
			if err != nil {
				return cpus, err
			}
			upper, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err != nil {
				return cpus, err
			}
			for i := lower; i <= upper; i++ {
				cpus = append(cpus, i)
			}
		default:
			return cpus, fmt.Errorf("unexpected cpu range bounds: %v", bounds)
		}
	}
	sort.Ints(cpus)
	return cpus, nil
}
