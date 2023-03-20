ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:glibc
LABEL maintainer="Ben Kochie <superq@gmail.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/draytek_exporter /bin/draytek_exporter

EXPOSE      9103
USER        nobody
ENTRYPOINT  [ "/bin/draytek_exporter" ]
