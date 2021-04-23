package main

import (
	"flag"
	"net/http"

	cg "github.com/phpHavok/cgroups_exporter/cgroups"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Official port reserved for this project, see:
// https://github.com/prometheus/prometheus/wiki/Default-port-allocations
const officialPort = "9821"

func main() {
	portPtr := flag.String("port", officialPort, "the port to listen on")
	helpPtr := flag.Bool("help", false, "print usage")
	cgroupsRootPathPtr := flag.String("cgroups-root", "/sys/fs/cgroup", "path to the root of the cgroupsv1 hierarchy")
	cgroupSpecPathPtr := flag.String("cgroup-spec", "/proc/1/cgroup", "path to the cgroup specification file to use")
	flag.Parse()
	if *helpPtr {
		flag.Usage()
		return
	}
	// Print some help debug information
	log.Printf("serving cgroups from hierarchy root %s", *cgroupsRootPathPtr)
	log.Printf("reading cgroups from specification %s", *cgroupSpecPathPtr)
	// Trial load to make sure we won't hit problems later in the application
	_, err := cg.LoadCgroups(*cgroupSpecPathPtr, *cgroupsRootPathPtr)
	if err != nil {
		log.Fatalf("unable to read cgroups specification file: %v", err)
	}
	// Create and register our cgroups collector
	cgroupsCollector := newCgroupsCollector(*cgroupSpecPathPtr, *cgroupsRootPathPtr)
	prometheus.MustRegister(cgroupsCollector)
	// Serve Prometheus HTTP requests
	log.Printf("listening on port %s", *portPtr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+(*portPtr), nil))
}
