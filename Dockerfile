FROM golang:1.16-buster

COPY . /usr/local/cgroups_exporter

WORKDIR /usr/local/cgroups_exporter

RUN make

ENTRYPOINT ["/usr/local/cgroups_exporter/cgroups_exporter"]
