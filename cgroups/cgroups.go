package cgroups

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	procCgroupIdxSubsystems = 1
	procCgroupIdxPath       = 2
)

// Cgroups represents a structure a cgroups across supported subsystems
type Cgroups struct {
	Cpuset  cpuset
	Cpuacct cpuacct
	Memory  memory
}

func readFile(root string, filename string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(root, filename))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadCgroups loads a structure containing all cgroups for a given specification file
func LoadCgroups(specPath string, cgroupsRootPath string) (Cgroups, error) {
	var cgroups Cgroups
	// Find and open the cgroup file for the process
	cgroupsPath := filepath.Clean(specPath)
	cgroupsFile, err := os.Open(cgroupsPath)
	if err != nil {
		return cgroups, err
	}
	defer cgroupsFile.Close()
	// Load the cgroup file as CSV data
	csvReader := csv.NewReader(cgroupsFile)
	csvReader.Comma = ':'
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		return cgroups, err
	}
	// Structure the CSV data into a map
	for _, csvLine := range csvLines {
		subsystems := strings.Split(csvLine[procCgroupIdxSubsystems], ",")
		for _, subsystem := range subsystems {
			// Empty subsystem names are possible for some reason, so skip over those
			if len(subsystem) < 1 {
				log.Debug("skipping empty subsystem")
				continue
			}
			cgroupAbsolutePath := filepath.Join(cgroupsRootPath, strings.TrimPrefix(subsystem, "name="), csvLine[procCgroupIdxPath])
			if _, err := os.Stat(cgroupAbsolutePath); os.IsNotExist(err) {
				return cgroups, fmt.Errorf("cgroup path doesn't exist: %s", cgroupAbsolutePath)
			}
			switch subsystem {
			case "cpuset":
				cgroups.Cpuset = cpuset(cgroupAbsolutePath)
			case "cpuacct":
				cgroups.Cpuacct = cpuacct(cgroupAbsolutePath)
			case "memory":
				cgroups.Memory = memory(cgroupAbsolutePath)
			default:
				log.Debugf("skipping unimplemented subsystem: %v", subsystem)
			}
		}
	}
	return cgroups, nil
}

// LoadProcessCgroups loads a structure containing all cgroups for a given process
func LoadProcessCgroups(pid int, cgroupsRootPath string) (Cgroups, error) {
	// Find and open the cgroup file for the process
	cgroupsPath := filepath.Join("/proc", strconv.Itoa(pid), "cgroup")
	return LoadCgroups(cgroupsPath, cgroupsRootPath)
}
