FROM golang:alpine AS build-env
ADD . /src
ENV CGO_ENABLED 0
RUN cd /src && apk add git && go build -o exabgp_exporter ./cmd/exabgp_exporter && go build -o exabgp_listener ./cmd/exabgp_listener 

FROM ubuntu:16.04
ENV EXABGP_VERSION 4.0.10
ENV HOME /root
ENV S6_LOGGING 1
EXPOSE 9569
WORKDIR /root
RUN mkdir -p /exabgp/run
RUN mkdir -p /exabgp/etc/exabgp
RUN mkdir -p /gobgp


COPY --from=build-env /src/exabgp_exporter /exabgp/exabgp_exporter
COPY --from=build-env /src/exabgp_listener /exabgp/exabgp_listener
RUN apt-get update && apt-get install -qy --no-install-recommends \
# Python
    python3-pip \
    python3-setuptools \
# Utility
    iproute2 \
    socat \
    curl \
    tmux \
    vim-nox \
 && rm -rf /var/lib/apt/lists/* \
 && pip3 install exabgp==${EXABGP_VERSION}


ADD https://github.com/osrg/gobgp/releases/download/v2.3.0/gobgp_2.3.0_linux_amd64.tar.gz /gobgp.tar.gz
RUN tar -xvf /gobgp.tar.gz -C /gobgp/

ADD https://github.com/just-containers/s6-overlay-builder/releases/download/v1.8.5/s6-overlay-portable-amd64.tar.gz /tmp/
RUN tar xzf /tmp/s6-overlay-portable-amd64.tar.gz -C /

RUN mkfifo /exabgp/run/exabgp.in
RUN mkfifo /exabgp/run/exabgp.out
RUN chmod 666 /exabgp/run/exabgp.*
COPY docker/files/exabgp.conf /exabgp/etc/exabgp/exabgp.conf
COPY docker/files/gobgp.yaml /gobgp/gobgp.yaml
COPY docker/files/rsyslog.conf /etc/rsyslog.conf
COPY docker/*.sh /root/
RUN mkdir /etc/services.d/gobgp
RUN mkdir /etc/services.d/exabgp
COPY docker/exabgp.sh /etc/services.d/exabgp/run
COPY docker/gobgpd.sh /etc/services.d/gobgp/run
ENTRYPOINT [ "/init" ]