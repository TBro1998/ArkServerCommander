# 多阶段构建 - 前端构建阶段
FROM node:20-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/web

# 复制前端项目文件
COPY web/package*.json ./
COPY web/pnpm-lock.yaml ./

# 安装 pnpm 并下载依赖
RUN npm install -g pnpm && pnpm install

# 复制前端源代码
COPY web/ ./

# 构建前端并生成静态文件
RUN pnpm build && pnpm generate

# 后端构建阶段
FROM golang:1.23-alpine AS backend-builder

# 安装必要的包（纯 Go 构建不需要 git）
# RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制后端项目文件
COPY server/go.mod server/go.sum ./

# 下载依赖
RUN go mod download

# 复制后端源代码
COPY server/ ./

# 整理依赖并构建后端应用
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 最终运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk add --no-cache ca-certificates docker-cli sqlite wget

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从后端构建阶段复制可执行文件
COPY --from=backend-builder /app/main .

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /app/web/.output/public ./static

# 创建数据目录
RUN mkdir -p /data && chown -R appuser:appgroup /data /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release
ENV DB_PATH=/data/ark_server.db
ENV STATIC_PATH=/app/static


# 启动应用
CMD ["./main"] 