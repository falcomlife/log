go build
docker build -t registry.cn-hangzhou.aliyuncs.com/cotte-internal/log:0.0.1 .
docker push registry.cn-hangzhou.aliyuncs.com/cotte-internal/log:0.0.1