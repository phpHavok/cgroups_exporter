# cgroups_exporter
A Prometheus exporter for cgroup-level metrics.

## Compiling from Source

This project is written primarily in Go and requires Go v1.16 or later to
compile.

To build, you can just type `go build`, and Go will handle everything.
Alternatively, if you have Make installed, you can just type `make`. Both of
these methods produce an executable binary called `cgroups_exporter`.

## Usage

To view help, run `./cgroups_exporter -h`.

```
$ ./cgroups_exporter -h
Usage of ./cgroups_exporter:
  -cgroup-spec string
        path to the cgroup specification file to use (default "/proc/1/cgroup")
  -cgroups-root string
        path to the root of the cgroupsv1 hierarchy (default "/sys/fs/cgroup")
  -help
        print usage
  -port string
        the port to listen on (default "9821")
```

The `-cgroup-spec` option is the most critical and specifies the path to the
cgroups specification file which indiciates which cgroups will be used. These
cgroup specification files have a structured format and they're commonly
found at `/proc/$$/cgroup` for some process ID `$$`. The cgroups specified in
the specification file will be tracked by this exporter while all other
cgroups will be ignored. If you want to track cgroups for two or more
different processes, you should run two or more copies of this exporter on
different ports.

The `-i` option is the most critical and specifies the path to the input file
from which to read the schedule of events. If not specified, the program will
read from standard input.

The `-cgroups-root` option allows you to change the default location of the
cgroupsv1 hierarchy if it happens to be mounted somewhere unusual. This can also
be handy when using a container in case you want to mount the hierarchy in a
different location for the sake of the container.

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
    phphavok/cgroups_exporter
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
container from Docker Hub.

```
singularity run \
    -B /sys/fs/cgroup:/sys/fs/cgroup \
    -B `pwd`:/data \
    docker://phphavok/cgroups_exporter
```