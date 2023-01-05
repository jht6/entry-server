#!/bin/bash

tag=`date '+%Y%m%d_%H%M%S'`

echo "TAG=${tag}"

# 编译
go mod download
go build -o eserver

# 构建镜像
docker build -t mirrors.xx.com/jht/entry-server:$tag .
docker push mirrors.xx.com/jht/entry-server:$tag

rm eserver