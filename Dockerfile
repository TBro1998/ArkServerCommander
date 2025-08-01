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

# 最终运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk add --no-cache ca-certificates docker-cli sqlite wget

# 设置工作目录
WORKDIR /app

# 从后端构建阶段复制可执行文件
COPY --from=backend-builder /app/main .

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /app/ui/public ./static/public
COPY --from=frontend-builder /app/ui/.next/standalone ./static/.next/standalone
COPY --from=frontend-builder /app/ui/.next/static ./static/.next/static

# 创建数据目录
RUN mkdir -p /data

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release
ENV DB_PATH=/data/ark_server.db

# 启动应用
CMD ["./main"]