#!/bin/bash

# 修复 Go 代码格式的脚本
# 主要问题：多行函数调用后的空格和缩进

cd /root/clawd/ArkServerCommander/server

# 需要修复的文件列表
files=(
    "service/docker_manager/rollback.go"
    "service/docker_manager/container_with_rollback.go"
    "service/docker_manager/rollback_test.go"
    "service/server/server_create_continue.go"
    "service/server/server_create_docker_resources.go"
    "service/server/server_create_with_rollback.go"
    "service/server/server_create_with_transaction.go"
    "service/server/server_start_continue.go"
    "service/server/server_start_with_rollback.go"
)

echo "开始修复 Go 代码格式..."

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "处理: $file"
        # 删除行尾空格
        sed -i 's/[[:space:]]*$//' "$file"
        # 将制表符转换为空格（如果需要）
        # sed -i 's/\t/    /g' "$file"
    else
        echo "文件不存在: $file"
    fi
done

echo "格式修复完成！"
