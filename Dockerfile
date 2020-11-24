# This is dockerfile for log project
# Version: 0.0.1
# Author: sora

#FROM golang
FROM  golang:1.14.4

MAINTAINER sora sorawingwind@163.com

WORKDIR /opt

COPY log-controller /opt
COPY message.template /opt

ENV LOG.ENV PROD