# This is dockerfile for log project
# Version: 0.0.1
# Author: sora

#FROM golang
FROM  golang:1.16.6

MAINTAINER sora sorawingwind@163.com

WORKDIR /opt

COPY log-controller /opt
COPY message.template /opt
COPY web /opt/web

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ENV LOG.ENV PROD
