# cgroups_exporter
A Prometheus exporter for cgroup-level metrics.

## Compiling from Source

This project is written primarily in Go and requires Go v1.16 or later to
compile.

To build, you can just type `go build`, and Go will handle everything.
Alternatively, if you have Make installed, you can just type `make`. Both of
these methods produce an executable binary called `cgroups_exporter`.

## Usage

To view help, run `./cgroups_exporter -help`.

```
$ ./cgroups_exporter -help
Usage of ./cgroups_exporter:
  -cgroups-root string
        path to the root of the cgroupsv1 hierarchy (default "/sys/fs/cgroup")
  -file string
        path to the cgroup specification file to use if method is file, ignored otherwise (default "/proc/1/cgroup")
  -help
        print usage
  -method string
        one of: file, slurm (default "slurm")
  -port string
        the port to listen on (default "9821")
```

The `-cgroups-root` option allows you to change the default location of the
cgroupsv1 hierarchy if it happens to be mounted somewhere unusual. This can also
be handy when using a container in case you want to mount the hierarchy in a
different location for the sake of the container.

The `-method` option specifies which cgroup hierarchies will be monitored.
Valid options are `file` and `slurm`. If set to `file`, the cgroups
specification file specified by the addtional `-file` option will be read and
used to determine which cgroups to monitor. If set to `slurm`, the program
will monitor the node for any jobs running under the Slurm scheduler and
output labeled stastics for each job over time.

The `-file` option is only used if the `-method` is set to `file` and
specifies the path to the cgroups specification file which indiciates which
cgroups will be used. These cgroup specification files have a structured
format and they're commonly found at `/proc/$$/cgroup` for some process ID
`$$`. The cgroups specified in the specification file will be tracked by this
exporter while all other cgroups will be ignored. If you want to track
cgroups for two or more different processes, you should run two or more
copies of this exporter on different ports.

The `-port` option allows you to change the port that the Prometheus HTTP
server will listen on for requests. The default port is recommended unless
you need to run two or more copies of this exporter.

## Docker

A Docker container is provided for systems where that is more convenient (such
as a Kubernetes cluster). You can build it manually using the provided
`Dockerfile`, or just pull the pre-built copy from Docker Hub. Example usage
follows:

```
docker run -t --rm \
    --mount type=bind,src=/sys/fs/cgroup,dst=/sys/fs/cgroup,readonly \
    phphavok/cgroups_exporter -method file -file /proc/31337/cgroup
```

We specify `-t` so that we're allocated a pseudo-terminal which makes the
logging output look nice and formatted. The `--rm` option automatically
cleans up the container on exit. The first mount command passes through the
cgroupv1 hierarchy (/sys/fs/cgroup) on the parent system to the same location
within the container. By default, Docker will often have some of the cgroup
hierarchy present within the container, but not all of it. This application
will need to see the full hierarchy, so this read-only bind mount takes care
of that. If you run into issues with mounting over the existing hierarchy
within the container, you can change the target to some other location and
then pass the `-cgroup-root` option to the program to accommodate that
change. The entrypoint to the container is the program itself, so you can
just pass its parameters to the run command.

## Singularity

You can use Singularity (e.g., on an HPC cluster) to run the above Docker
container from Docker Hub. If using the `slurm` method, make sure **not** to
have Singularity create a PID namespace for the job (i.e., leave off the `-p`
option), otherwise the container will be unable to properly detect Slurm jobs
running outside the container.

```
singularity run \
    -B /sys/fs/cgroup:/sys/fs/cgroup \
    -B `pwd`:/data \
    docker://phphavok/cgroups_exporter -method file -file /proc/31337/cgroup
```

## Grafana Dashboard

A convenient Grafana dashboard is available [in this
repository](grafana/slurm_cgroups_dashboard.json), and is also published as
[public dashboard 14587](https://grafana.com/grafana/dashboards/14587) for quick
installation. The provided dashboard works for invocations of the
cgroups\_exporter that operate in `-method slurm` mode.
