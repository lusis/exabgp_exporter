FROM golang:alpine AS build-env


FROM golang:alpine
ENV CGO_ENABLED 0
ENV EXABGP_VERSION 4.0.10
ENV HOME /root
ENV S6_LOGGING 1
WORKDIR /root

RUN mkdir -p /exabgp/run
RUN mkdir -p /exabgp/etc/exabgp
RUN mkdir -p /gobgp
RUN mkfifo /exabgp/run/exabgp.in
RUN mkfifo /exabgp/run/exabgp.out
RUN chmod 666 /exabgp/run/exabgp.*
COPY docker/files/exabgp.conf /exabgp/etc/exabgp/exabgp.conf
COPY docker/files/gobgp.yaml /gobgp/gobgp.yaml
COPY docker/files/rsyslog.conf /etc/rsyslog.conf
COPY docker/*.sh /root/

# add support packages
ADD https://github.com/osrg/gobgp/releases/download/v2.3.0/gobgp_2.3.0_linux_amd64.tar.gz /gobgp.tar.gz
ADD https://github.com/just-containers/s6-overlay-builder/releases/download/v1.8.5/s6-overlay-portable-amd64.tar.gz /tmp/

# install support packages
RUN tar xvf /gobgp.tar.gz -C /gobgp/
RUN tar xzf /tmp/s6-overlay-portable-amd64.tar.gz -C /

RUN mkdir /etc/services.d/gobgp
RUN mkdir /etc/services.d/exabgp
RUN mkdir /etc/services.d/exabgp_exporter_good
RUN mkdir /etc/services.d/exabgp_exporter_bad
COPY docker/exabgp.sh /etc/services.d/exabgp/run
COPY docker/gobgpd.sh /etc/services.d/gobgp/run
COPY docker/exporter_valid.sh /etc/services.d/exabgp_exporter_good/run
COPY docker/exporter_invalid.sh /etc/services.d/exabgp_exporter_bad/run

RUN apk add\
    bash \
    py3-pip \
    py3-setuptools \
    socat \
    curl \
    git \
    musl-dev \
    linux-headers \
    python3-dev \
    gcc

# build binary
ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o /exabgp/exabgp_exporter -ldflags "-s" -a -tags netgo ./cmd/exabgp_exporter

ENTRYPOINT [ "/root/install-and-init.sh" ]