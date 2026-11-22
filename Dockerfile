# 前端构建阶段
FROM node:26-alpine AS frontend-builder

WORKDIR /app/ui

COPY ui/ ./
RUN npm ci

RUN npm run build

# 后端构建阶段
FROM golang:1.27-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# Go 模块代理，可选。访问不了 proxy.golang.org 时传入，例如 --build-arg GOPROXY=https://goproxy.cn,direct；
# 为空时 Go 使用默认代理
ARG GOPROXY

# 复制后端项目文件
COPY server/go.mod server/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码
COPY server/ ./

# 版本号由 CI 按 tag 传入（不带 v，例如 0.0.3），本地构建为 dev，用于界面显示和检查新版本
ARG VERSION=dev

# 整理依赖并构建后端应用
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -trimpath -ldflags="-s -w -X ark-server-commander/utils.Version=${VERSION}" -o main .

# 最终运行阶段：Node 运行 Next.js（对外），Go 后端只在容器内回环地址监听
FROM node:26-alpine

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
RUN mkdir -p /data /server-files

# 服务端文件卷：compose 不用写，Docker 自动创建匿名卷（重建容器时接回原来的卷）。只在 Docker Desktop
# （Windows / macOS，经文件共享访问宿主机磁盘很慢）上并且 /data 绑定了宿主机目录时使用，Linux 上服务端文件仍在 /data，
# 这个卷是空的。见 server/service/docker_manager/files_volume.go
VOLUME ["/server-files"]

# 对外只暴露 Next.js 的 8080；Go 后端监听 127.0.0.1:8081，不对外
EXPOSE 8080

# 设置环境变量
ENV NODE_ENV=production
ENV GIN_MODE=release
ENV DB_PATH=/data/ark_server.db

# 启动应用
CMD ["./docker-entrypoint.sh"]
