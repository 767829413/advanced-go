#!/usr/bin/env bash
# 实现遍历获取指定目录下指定后缀的文件
for file in $(find ./api/ -type f -name "*.proto"); do
    echo $file
    protoc -I . --go_out=. --go-grpc_out=. $file
done
