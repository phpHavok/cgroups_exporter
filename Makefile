TARGET:=cgroups_exporter
SOURCES:=main.go collector.go cgroups/cgroups.go cgroups/cpuset.go cgroups/cpuacct.go

$(TARGET): $(SOURCES)
	go build -o $@

clean:
	rm -f $(TARGET)

.PHONY: clean
