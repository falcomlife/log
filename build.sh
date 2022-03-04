#!/bin/bash
if [ "$1" = "" ] && [ "$1" != "ciiplat" ] && [ "$1" != "lync2m" ] && [ "$1" != "mdy" ] && [ "$1" != "xy" ] && [ "$1" != "release" ];then
	echo "please input a right env(ciiplat/lync2m/mdy/xy/release)"
	exit 1
else
  go build
  docker build -t swr.cn-north-4.myhuaweicloud.com/cotte-internal/log:0.0.3-$1 .
  docker push swr.cn-north-4.myhuaweicloud.com/cotte-internal/log:0.0.3-$1
	exit 0
fi
