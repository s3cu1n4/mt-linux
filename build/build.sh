#!/bin/bash

project_path=$(cd `dirname $0`; pwd)
project_name="${project_path##*/}"

cd $project_path
cd ..

# 在Linux下编译 
go build -a -o ./bin/mt-linux