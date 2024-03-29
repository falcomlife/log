# This is dockerfile for log project
# Version: 0.0.1
# Author: sora

#FROM golang
FROM  golang:1.14.4

MAINTAINER sora sorawingwind@163.com

WORKDIR /opt

COPY klog-controller /opt
COPY template /opt/template
COPY web /opt/web

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    LC_CTYPE=en_US.UTF-8

ENV LOG.ENV PROD