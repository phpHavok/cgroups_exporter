TARGET:=cgroups_exporter
SOURCES:=main.go cgroups/cgroups.go cgroups/cpuset.go cgroups/cpuacct.go cgroups/memory.go collectors/file.go collectors/slurm.go

$(TARGET): $(SOURCES)
	go build -o $@

clean:
	rm -f $(TARGET)

.PHONY: clean
