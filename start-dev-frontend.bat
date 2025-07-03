@echo off
echo =========================================
echo 🚀 启动 ARK 服务器管理器前端开发环境
echo =========================================

echo.
echo 📁 切换到前端目录...
cd web

echo.
echo 📦 检查依赖...
if not exist "node_modules" (
    echo 正在安装依赖...
    pnpm install
) else (
    echo 依赖已存在，跳过安装
)

echo.
echo 🌐 启动开发服务器...
echo 前端地址: http://localhost:3000
echo API代理: http://localhost:3000/api -> http://localhost:8080/api
echo.
echo 请确保后端服务已在 http://localhost:8080 运行
echo.
pnpm dev

pause 