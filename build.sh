#!/bin/bash

# 读取版本号
version=$(<VERSION)

# 获取当前 Git 提交的哈希
commitHash=$(git rev-parse --short HEAD)

# 获取当前时间并转换为所需格式
currentTime=$(date +"%Y-%m-%d")

# 设置编译平台
export GOOS=windows
export GOARCH=amd64

# 构建项目
wails build -ldflags "-X main.Version=$version -X main.BuildTime=$currentTime -X main.Commit=$commitHash -w -s" --upx
