#!/bin/sh
# 同一个容器里运行两个进程：
#   Next.js (对外 :8080)  ->  Go 后端 (仅容器内回环 127.0.0.1:8081)
# 端口在这里写死，避免 .env 里的 SERVER_PORT 等变量把两者配到同一个端口上。

GO_PORT=8081
WEB_PORT=8080

SERVER_PORT=$GO_PORT LISTEN_HOST=127.0.0.1 ./main &
GO_PID=$!

cd /app/ui || exit 1
API_BASE="http://127.0.0.1:${GO_PORT}/api" PORT=$WEB_PORT HOSTNAME=0.0.0.0 node server.js &
WEB_PID=$!

trap 'kill "$GO_PID" "$WEB_PID" 2>/dev/null; exit 0' TERM INT

# 任一进程退出，就结束另一个并以非零状态退出，让 Docker 的重启策略接管
while kill -0 "$GO_PID" 2>/dev/null && kill -0 "$WEB_PID" 2>/dev/null; do
  sleep 1
done

kill "$GO_PID" "$WEB_PID" 2>/dev/null
exit 1
