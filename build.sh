#!/bin/bash
if [ "$1" = "" ] && [ "$1" != "ciiplat" ] && [ "$1" != "lync2m" ];then
	echo "please input a right env(ciiplat/lync2m)"
	exit 1
else
  go build
  docker build -t registry.cn-hangzhou.aliyuncs.com/cotte-internal/log:0.0.2-$1 .
  docker push registry.cn-hangzhou.aliyuncs.com/cotte-internal/log:0.0.2-$1
	exit 0
fi
