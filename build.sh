#!/bin/bash
if [ "$1" = "" ] && [ "$1" != "ciiplat" ] && [ "$1" != "lync2m" ] && [ "$1" != "mdy" ] && [ "$1" != "xy" ];then
	echo "please input a right env(ciiplat/lync2m/mdy/xy)"
	exit 1
else
  go build
  docker build -t registry.cn-hangzhou.aliyuncs.com/cotte-internal/log:0.0.2-$1 .
  docker push registry.cn-hangzhou.aliyuncs.com/cotte-internal/log:0.0.2-$1
	exit 0
fi
