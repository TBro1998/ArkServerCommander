# 前端构建阶段
FROM node:24-alpine AS frontend-builder

WORKDIR /app/ui

COPY ui/ ./
RUN npm install

RUN npm run build

# 后端构建阶段
FROM golang:1.24-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# 复制后端项目文件
COPY server/go.mod server/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码
COPY server/ ./

# 整理依赖并构建后端应用
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -trimpath -ldflags="-s -w" -o main .

# 最终运行阶段：Node 运行 Next.js（对外），Go 后端只在容器内回环地址监听
FROM node:24-alpine

# 安装必要的包
RUN apk add --no-cache ca-certificates docker-cli sqlite wget

# 设置工作目录
WORKDIR /app

# 从后端构建阶段复制可执行文件
COPY --from=backend-builder /app/main ./main

# 从前端构建阶段复制 Next.js standalone 产物（server.js 在 standalone 根目录）
COPY --from=frontend-builder /app/ui/.next/standalone ./ui
COPY --from=frontend-builder /app/ui/.next/static ./ui/.next/static
COPY --from=frontend-builder /app/ui/public ./ui/public

# 入口脚本：同时启动 Go 后端和 Next.js
COPY docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh

# 创建数据目录
RUN mkdir -p /data

# 对外只暴露 Next.js 的 8080；Go 后端监听 127.0.0.1:8081，不对外
EXPOSE 8080

# 设置环境变量
ENV NODE_ENV=production
ENV GIN_MODE=release
ENV DB_PATH=/data/ark_server.db

# 启动应用
CMD ["./docker-entrypoint.sh"]
