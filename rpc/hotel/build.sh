#!/bin/bash

# 酒店管理服务构建脚本

echo "开始构建酒店管理服务..."

# 1. 检查 kitex 工具是否安装
if ! command -v kitex &> /dev/null; then
    echo "kitex 工具未安装，正在安装..."
    go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
fi

# 2. 生成 Kitex 代码
echo "正在生成 Kitex 代码..."
cd ../../  # 回到项目根目录
kitex -module example_shop -service hotel_service idl/hotel.thrift

# 3. 整理依赖
echo "正在整理依赖..."
go mod tidy

# 4. 构建服务
echo "正在构建服务..."
cd rpc/hotel/main
go build -o ../../../hotel_service .

echo "构建完成！"
echo "运行服务: ./hotel_service"
