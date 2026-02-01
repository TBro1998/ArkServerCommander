# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

ARK Server Commander 是一个用于管理 Linux 上 ARK Survival Evolved 游戏服务器的全栈 Web 应用程序。它提供了一个现代化的 Web 界面，用于创建、配置和管理运行在 Docker 容器中的多个 ARK 服务器。

**技术栈：**
- 后端：Go 1.24 + Gin 框架 + GORM + SQLite
- 前端：Next.js 15.4.2 + React 19 + TypeScript + Tailwind CSS 4
- 基础设施：Docker + Docker Compose

## 开发命令

### 后端 (Go)

**工作目录：** `server/`

```bash
# 安装依赖
go mod download

# 整理依赖
go mod tidy

# 运行后端服务器（需要 Docker）
go run main.go

# 构建后端二进制文件
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -trimpath -ldflags="-s -w" -o main .

# 生成 Swagger 文档
swag init

# 运行测试
go test ./...

# 运行特定测试
go test ./service/docker_manager -v -run TestRollback

# 格式化代码
gofmt -w .
go fmt ./...
```

**环境变量：**
- `JWT_SECRET`：JWT 密钥（最少 32 字符，必需）
- `DB_PATH`：数据库文件路径（默认：`/data/ark_server.db`）
- `SERVER_PORT`：服务器端口（默认：`8080`）
- `GIN_MODE`：Gin 模式（`debug` 或 `release`）

### 前端 (Next.js)

**工作目录：** `ui/`

```bash
# 安装依赖
npm install

# 使用 Turbopack 运行开发服务器
npm run dev

# 构建生产版本
npm run build

# 启动生产服务器
npm start

# 运行代码检查
npm run lint
```

### Docker

```bash
# 构建并启动所有服务
docker-compose up -d

# 构建镜像
docker build -t arkservercommander:latest .

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重新构建并重启
docker-compose up -d --build
```

## 架构概览

### 后端结构 (`/server`)

后端遵循分层架构模式：

```
HTTP 层（路由 + 控制器）
    ↓
服务层（业务逻辑）
    ↓
数据层（模型 + 数据库）
```

**关键目录：**

- **`/config`**：配置管理（JWT 密钥验证、环境变量）
- **`/controllers`**：HTTP 请求处理器（认证、服务器、镜像）
- **`/models`**：GORM 数据模型（User、Server、ServerArgs）
- **`/database`**：数据库初始化和连接管理
- **`/middleware`**：JWT 认证中间件
- **`/routes`**：API 路由注册和 SPA 回退路由
- **`/service/docker_manager`**：Docker SDK 操作（容器、镜像、卷、回滚）
- **`/service/server`**：带事务支持的服务器生命周期管理
- **`/utils`**：日志记录（Zap）、JWT 工具、辅助函数
- **`/docs`**：自动生成的 Swagger/OpenAPI 文档

**重要文件：**
- [main.go](server/main.go)：应用程序入口点，初始化日志、配置、数据库、Docker 管理器
- [routes/routes.go](server/routes/routes.go)：API 路由注册和 SPA 路由
- [service/docker_manager/docker_manager.go](server/service/docker_manager/docker_manager.go)：单例 Docker 客户端管理器
- [service/server/server_service.go](server/service/server/server_service.go)：带事务支持的服务器 CRUD 操作

### 前端结构 (`/ui`)

Next.js 15 App Router 使用路由组：

- **`/src/app/(auth)`**：公开认证页面（登录、初始化）
- **`/src/app/(protected)`**：需要认证的受保护页面（服务器、主页）
- **`/src/app/api`**：Next.js API 路由（代理到后端）
- **`/public`**：静态资源
- **`/messages`**：国际化翻译文件（en.json、zh.json）

**状态管理：**
- Zustand 用于全局状态
- Axios 用于 HTTP 客户端（带拦截器）

### Docker 集成

应用程序使用 Docker Go SDK 管理 ARK 服务器容器：

1. **容器管理**：创建、启动、停止、重建 ARK 服务器容器
2. **卷管理**：自动创建持久化卷（`ark-server-<id>`）
3. **镜像管理**：拉取和更新 `tbro98/ase-server:latest` 镜像
4. **回滚支持**：基于事务的操作，失败时自动回滚

**容器命名规范：** `ark-server-<server_id>`
**卷命名规范：** `ark-server-<server_id>`

## 核心概念

### 服务器生命周期

1. **创建**：用户创建服务器配置 → 创建数据库记录
2. **启动**：服务器服务创建 Docker 容器 → 如需要则拉取镜像 → 创建卷 → 启动容器
3. **停止**：停止 Docker 容器（容器保留）
4. **重建**：停止并删除旧容器 → 使用更新的配置创建新容器 → 失败时回滚
5. **删除**：停止容器 → 删除容器 → 删除卷 → 删除数据库记录

### 事务与回滚系统

服务器服务使用带回滚支持的事务模式：

- 操作被跟踪在回滚栈中
- 失败时，操作按 LIFO 顺序反向执行
- 示例：如果在创建卷后容器创建失败，卷会自动删除

参见 [service/docker_manager/rollback.go](server/service/docker_manager/rollback.go) 了解实现细节。

### 认证流程

1. 首次启动 → 系统检查是否存在管理员用户
2. 如果没有用户 → 重定向到初始化页面 → 创建管理员用户
3. 登录 → 生成 JWT 令牌 → 令牌存储在 cookie 中
4. 受保护路由 → 中间件验证 JWT → 提取用户 ID
5. API 请求 → 令牌在 Authorization 头中发送（`Bearer <token>`）

### 配置文件

ARK 服务器使用两个主要配置文件，在数据库中以 JSON 字符串存储：

- **GameUserSettings.ini**：服务器设置（会话名称、密码、最大玩家数等）
- **Game.ini**：游戏规则和机制

这些配置在运行时通过 Docker 卷注入到容器中。

## 常见开发模式

### 添加新的 API 端点

1. 如需要，在 `/models` 中定义模型
2. 在 `/controllers` 中创建控制器函数
3. 为控制器添加 Swagger 注解
4. 在 `/routes/routes.go` 中注册路由
5. 重新生成 Swagger 文档：`swag init`

### 添加新的 Docker 操作

1. 在 `/service/docker_manager/docker_manager.go` 中为 `DockerManager` 添加方法
2. 如果操作需要回滚，在 `/service/docker_manager/rollback.go` 中添加回滚函数
3. 在服务器服务中使用该操作，并提供事务支持

### 前端 API 集成

1. 在 `/ui/src/app/api` 中创建 API 路由（代理到后端）
2. 在组件中使用 Axios 调用 Next.js API 路由
3. 通过基于 cookie 的 JWT 处理认证

## 重要说明

### 安全性

- **JWT 密钥**：必须至少 32 字符。应用程序拒绝使用弱密钥启动。
- **Docker Socket**：应用程序需要访问 `/var/run/docker.sock`（特权模式）
- **CORS**：当前允许所有来源（开发模式）- 生产环境应限制

### 数据库

- 使用 GORM 的 SQLite 数据库
- 启动时自动迁移
- User 和 Server 模型启用软删除
- 数据库文件位置：`/data/ark_server.db`（可通过 `DB_PATH` 配置）

### Docker 要求

- 必须安装并运行 Docker
- 应用程序在启动时检查 Docker 状态
- 需要特权模式来管理主机 Docker 容器

### 日志记录

- 使用 Uber Zap 的结构化日志
- 日志级别：Debug、Info、Warn、Error、Fatal
- 日志输出到 stdout（由 Docker 捕获）

## API 文档

Swagger 文档自动生成，可通过以下地址访问：
- **本地**：http://localhost:8080/swagger/index.html
- **注解**：在控制器函数中使用 Swag 注解
- **重新生成**：在 `/server` 目录运行 `swag init`

## 测试

目前测试覆盖率较低。存在一个测试文件：
- [service/docker_manager/rollback_test.go](server/service/docker_manager/rollback_test.go)

运行测试：`go test ./...`

## 部署

应用程序使用多阶段 Docker 构建：

1. **前端构建阶段**：Node 24 Alpine → npm install → npm build
2. **后端构建阶段**：Go 1.24 Alpine → go build
3. **运行阶段**：Alpine + Docker CLI + SQLite → 复制二进制文件和静态文件

**生产环境检查清单：**
- 设置强 `JWT_SECRET`（最少 32 字符）
- 配置 `DB_PATH` 用于持久化存储
- 挂载 `/var/run/docker.sock` 以访问 Docker
- 设置 `GIN_MODE=release`
- 在生产环境中限制 CORS 来源

## 故障排除

### 后端无法启动
- 检查 JWT_SECRET 是否已设置且 >= 32 字符
- 验证 Docker 是否运行：`docker ps`
- 检查数据库文件权限

### 前端构建失败
- 清理 `.next` 目录：`rm -rf .next`
- 删除 `node_modules` 并重新安装：`rm -rf node_modules && npm install`

### Docker 操作失败
- 验证 Docker socket 是否已挂载：`-v /var/run/docker.sock:/var/run/docker.sock`
- 检查容器是否具有特权模式：`privileged: true`
- 确保 Docker 守护进程可访问

## 项目状态

此项目正在积极开发中。功能状态详见 [README.md](README.md)：
- ✅ 已实现：服务器 CRUD、Docker 管理、认证
- 🚧 计划中：RCON、监控、模组管理、备份、国际化
- 🚀 未来计划：K8S 支持、服务器列表网站
