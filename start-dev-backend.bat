@echo off
echo =========================================
echo 🚀 启动 ARK 服务器管理器后端开发环境
echo =========================================

echo.
echo 📁 切换到后端目录...
cd server

echo.
echo 🔧 检查Go环境...
go version
if %errorlevel% neq 0 (
    echo ❌ Go环境未安装或未配置
    echo 请先安装Go: https://golang.org/dl/
    pause
    exit /b 1
)

echo.
echo 📦 检查依赖...
if not exist "go.sum" (
    echo 正在下载依赖...
    go mod tidy
) else (
    echo 依赖已存在，跳过下载
)

echo.
echo 🐳 检查Docker环境...
docker --version
if %errorlevel% neq 0 (
    echo ⚠️  Docker未安装或未运行
    echo 部分功能可能不可用
) else (
    echo ✅ Docker环境正常
)

echo.
echo 🌐 启动后端服务器...
echo 后端地址: http://localhost:8080
echo API文档: http://localhost:8080/swagger/index.html
echo.
go run main.go

pause 