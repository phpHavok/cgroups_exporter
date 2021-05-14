package cgroups

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type memory string

// GetUsageInBytes returns the current memory in use by the cgroup in bytes
func (c memory) GetUsageInBytes() (int, error) {
	data, err := readFile(string(c), "memory.usage_in_bytes")
	if err != nil {
		return 0, err
	}
	usage, err := strconv.Atoi(strings.TrimSpace(data))
	if err != nil {
		log.Errorf("unable to convert memory usage to integer: %v", err)
		return usage, err
	}
	return usage, nil
}
