# This is dockerfile for micro service cluster
# Version: 1.0
# Author: sora

FROM  {{.Registry.Address}}/cotte-internal/nginx:0.0.1

ARG CONTAINERNAME

WORKDIR /usr/share/nginx/html

COPY nginx.conf /etc/nginx/nginx.conf
ADD web.tar.gz /usr/share/nginx/html/{{.Name}}/