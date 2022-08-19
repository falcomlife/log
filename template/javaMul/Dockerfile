# This is dockerfile for micro service cluster
# Version: 1.0
# Author: sora

FROM {{.Registry.Address}}/cotte-internal/centos7-tools-jdk8-git:0.0.3
ARG CONTAINERNAME
WORKDIR /opt
#copy config to container
COPY microservice.jar /opt
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ENV LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    LC_CTYPE=en_US.UTF-8