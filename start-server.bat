@echo off
echo 启动 ARK 服务器管理器后端...
cd server
go mod tidy
go run main.go
pause