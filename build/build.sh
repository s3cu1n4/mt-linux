#!/bin/bash

project_path=$(cd `dirname $0`; pwd)
project_name="${project_path##*/}"

cd $project_path


if [ ! -d "../bin" ]; then

mkdir ../bin

fi

# 在Linux下编译 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ../bin/mt-linux ../main.go