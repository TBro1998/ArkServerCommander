#!/usr/bin/env bash
# 开发用：在本机直接运行 Go 后端和 Next 开发服务器，不构建镜像（正式部署用 docker-compose.yml）
# 用法：bash dev.sh，Ctrl+C 同时停止两者
# 需要：Go 1.27、Node >= 20.9、本机可用的 Docker
# 环境变量可在运行前 export 覆盖，例如 PUBLIC_HOST=1.2.3.4 bash dev.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 数据放在 data-dev，不和正式部署的 ./data（容器内创建，属主是 root）混用
export DB_PATH="${DB_PATH:-$ROOT/data-dev/ark_server.db}"
export SERVER_PORT="${SERVER_PORT:-8081}"
export LISTEN_HOST="${LISTEN_HOST:-127.0.0.1}"
export GIN_MODE="${GIN_MODE:-debug}"
export LOG_LEVEL="${LOG_LEVEL:-debug}"
export LOG_FORMAT="${LOG_FORMAT:-console}"
UI_PORT="${UI_PORT:-3000}"
export API_BASE="${API_BASE:-http://127.0.0.1:$SERVER_PORT/api}"

mkdir -p "$(dirname "$DB_PATH")"

if [ ! -d "$ROOT/ui/node_modules" ]; then
  echo "[dev] 安装前端依赖 (npm ci)"
  (cd "$ROOT/ui" && npm ci)
fi

# 任一进程退出或按下 Ctrl+C 时结束整个进程组（含 go run 编译出的子进程和 Next 的工作进程）
trap 'trap - EXIT INT TERM; kill 0 2>/dev/null' EXIT INT TERM

echo "[dev] 后端: http://$LISTEN_HOST:$SERVER_PORT  (Swagger: /swagger/index.html)"
(cd "$ROOT/server" && exec go run .) &

echo "[dev] 前端: http://localhost:$UI_PORT"
(cd "$ROOT/ui" && exec npx next dev --turbopack --port "$UI_PORT") &

wait -n
