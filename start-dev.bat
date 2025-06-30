@echo off
echo ===== 启动 ARK 服务器管理器开发环境 =====
echo.

echo [1/2] 启动后端服务器...
cd server
start "ARK Server Backend" cmd /k "go run main.go"
cd ..

echo [2/2] 等待3秒后启动前端...
timeout /t 3 /nobreak > nul

cd web
start "ARK Server Frontend" cmd /k "pnpm dev"
cd ..

echo.
echo ===== 开发环境启动完成 =====
echo 后端地址: http://localhost:8080
echo 前端地址: http://localhost:3000
echo Swagger API文档: http://localhost:8080/swagger/index.html
echo.
echo 按任意键关闭此窗口...
pause > nul 