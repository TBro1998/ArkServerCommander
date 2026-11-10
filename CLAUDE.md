# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

ARK Server Commander 是一个全栈 Web 应用，用于在 Linux 上通过 Web 界面创建、配置和管理运行在 Docker 容器中的多个 ARK: Survival Evolved 服务器。

**技术栈（2026-11 升级后）：**
- 后端：Go 1.27 + Gin 1.12 + GORM 1.31 + SQLite（`glebarez/sqlite`，纯 Go，无需 CGO）+ Docker SDK 28
- 前端：Next.js 16.4（App Router + Turbopack）+ React 19.3 + TypeScript 6 + Tailwind CSS 4 + ESLint 10
- 前端包管理：npm（`ui/package-lock.json`），不再使用 pnpm
- 镜像：构建阶段 `node:26-alpine` 和 `golang:1.27-alpine`，运行阶段 `node:26-alpine`

**版本固定说明：** TypeScript 固定在 6.x。`typescript-eslint` 暂不支持 7.0，装 7.0 会导致 ESLint 无法启动。待其支持后再升级。

## 开发命令

### 后端（工作目录 `server/`）

```bash
go mod download                 # 安装依赖
go run main.go                  # 运行（需要本机 Docker 和 JWT_SECRET）
go build ./... && go vet ./...  # 编译和静态检查
go test ./...                   # 测试
gofmt -l .                      # 列出未格式化的文件，应无输出
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o main .   # 构建二进制
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init            # 重新生成 Swagger（修改控制器注解后）
```

### 前端（工作目录 `ui/`）

```bash
npm ci            # 按锁文件安装依赖
npm run dev       # 开发服务器（Turbopack）
npm run build     # 生产构建（standalone 输出）
npm run lint      # eslint .（Next 16 已移除 next lint）
npx tsc --noEmit  # 类型检查
```

本机 Node 版本只需 >= 20.9；镜像内使用 Node 26。

### Docker

```bash
docker-compose up -d --build   # 构建并启动
docker-compose logs -f         # 查看日志
docker-compose down            # 停止
```

## 环境变量

后端：
- `JWT_SECRET`：必需，至少 32 字符，包含弱口令片段（如 `secret`、`password`）会被拒绝启动
- `DB_PATH`：数据库路径（代码默认 `ark_server.db`，镜像内为 `/data/ark_server.db`）
- `SERVER_PORT`：后端端口（默认 `8080`；容器内由入口脚本固定为 `8081`，`.env` 中的值被忽略）
- `LISTEN_HOST`：监听地址（默认 `127.0.0.1`；容器内固定为 `127.0.0.1`）
- `CORS_ALLOWED_ORIGINS`：跨域来源白名单，逗号分隔；为空表示仅同源
- `GIN_MODE`、`LOG_LEVEL`、`LOG_FORMAT`：Gin 模式和日志配置

前端（Next 服务端）：
- `API_BASE`：后端 API 地址，运行时读取，容器内为 `http://127.0.0.1:8081/api`
- `NEXT_PUBLIC_API_BASE`：开发环境回退值；都未设置时回退到 `http://127.0.0.1:8080/api`
- `COOKIE_SECURE`：通过 HTTPS 访问时设为 `true`，纯 HTTP 时保持 `false`，否则无法登录

示例见 `.env.example`。

## 架构概览

### 部署形态

单个容器内运行两个进程，由 `docker-entrypoint.sh` 启动：

```
浏览器 → Next.js（:8080，对外）→ Go 后端（127.0.0.1:8081，仅容器内）→ Docker Socket → ARK 服务器容器
```

- Go 后端不托管前端静态文件，只提供 `/api`、`/health` 和 `/swagger`
- 任一进程退出，入口脚本会结束另一个并以非零状态退出，交给 Docker 重启策略处理
- Swagger 位于后端的 `/swagger/index.html`，容器内不对外暴露，本地直接运行后端时才能访问
- 容器需要挂载 `/var/run/docker.sock` 并以 `privileged` 运行

### 后端结构（`server/`）

分层：controllers（HTTP）→ service（业务）→ docker_manager（Docker SDK）。

- `main.go`：初始化日志、配置、数据库，检查 Docker，注册 CORS 和路由
- `config/`：环境变量读取和 JWT 密钥校验
- `controllers/{auth,servers,images}/`：请求处理。`servers/errors.go` 的 `respondServiceError` 用 `errors.Is` 把服务层错误映射为 404、400、500
- `service/server/`：服务器生命周期。`errors.go` 定义哨兵错误，不要再用 `err.Error() == ...` 比较字符串
- `service/docker_manager/`：容器、镜像、卷操作和 `RollbackManager`（LIFO 回滚栈）
- `middleware/`：`auth.go`（JWT 校验）、`cors.go`（白名单 CORS）
- `models/`：GORM 模型 User、Server，以及请求和响应结构。`ServerResponse` 不包含管理员密码
- `utils/`：日志、JWT、密码、命名（`asm_info.go`）、配置文件模板（`config_files.go`）
- `docs/`：swag 生成的 Swagger 文件，不要手改

### API 路由（均以 `/api` 开头）

- 公开：`GET /auth/check-init`、`POST /auth/init`、`POST /auth/login`
- 需 JWT：`GET /profile`
- 服务器：`GET|POST /servers`、`GET|PUT|DELETE /servers/:id`、`POST /servers/:id/start|stop|recreate`、`GET /servers/:id/rcon`
- 镜像：`GET /images/status`、`POST /images/pull`、`GET /images/check-updates`、`POST /images/update`、`GET /images/affected`

### 前端结构（`ui/src/`）

- `app/(auth)/`：登录、初始化页面
- `app/(protected)/`：需登录页面（`home`、`servers`、`servers/new`、`servers/[id]/edit`）
- `app/api/`：服务端路由。`auth/{login,init,logout}` 处理令牌 cookie；`[...path]` 是通配代理，转发到后端
- `proxy.ts`：Next 16 中原 `middleware.ts` 的新约定。负责语言 cookie（`NEXT_LOCALE`）和登录重定向
- `lib/backend-server.ts`：`getApiBase`、令牌 cookie 读写、`authProxy`
- `stores/`：Zustand（`auth.ts`、`servers.ts`），HTTP 用 Axios
- `components/ui/`：shadcn/ui 风格基础组件；`components/servers/`：服务器卡片、编辑弹窗、ini 可视化编辑器
- `messages/en.ts`、`messages/zh.ts`：next-intl 翻译（TypeScript 文件，不是 JSON），入口为 `src/i18n.ts`

### 认证流程

1. 首次启动没有用户时，前端进入 `/init` 创建管理员
2. 登录或初始化请求先到 Next 的 `app/api/auth/*`，转发给后端，拿到 JWT 后写入 HttpOnly cookie `auth-token`，响应体里不含令牌
3. 其他 `/api/*` 请求经通配代理：从 cookie 取令牌，转成 `Authorization: Bearer <token>` 发给后端；后端返回 401 时清除 cookie
4. 前端脚本读不到令牌，页面加载时通过 `restoreSession`（请求 `/api/profile`）确认登录状态；退出登录调用 `/api/auth/logout`
5. `proxy.ts` 只检查 cookie 是否存在来做页面重定向，令牌真伪由后端校验

## Docker 集成

命名规范（见 `server/utils/asm_info.go`）：
- 容器：`ase-server-<server_id>`
- 数据卷：`ase-server-<server_id>`
- 插件卷：`ase-server-plugins-<server_id>`
- 镜像：`tbro98/ase-server:latest`（子模块 `ASE-Server-Docker` 是该镜像的源码）

### 服务器生命周期

1. **创建**：事务内写数据库记录并创建两个卷，创建容器；失败时 `defer` 清理已创建的卷（`committed` 标记）
2. **启动 / 停止**：用 `transitionStatus` 原子切换状态（`UPDATE ... WHERE status NOT IN (...)`，看 `RowsAffected`），避免并发请求重复启停；容器操作在后台协程中完成
3. **重建**：必要时先停止，删除旧容器，按当前配置创建新容器
4. **删除**：运行中拒绝删除；软删除数据库记录并删除容器，**不会删除 Docker 卷**，数据会保留在卷中

### 配置文件

GameUserSettings.ini 和 Game.ini 以文本形式保存在数据库中，在创建容器时写入。前端的可视化编辑器（`GameUserSettingsEditor`、`GameIniEditor`）和文本模式互相同步。

## 开发约定

### 添加 API 端点

1. 需要时在 `models/` 增加结构
2. 在 `service/server/` 实现业务，错误用哨兵错误，并在 `controllers/servers/errors.go` 补充映射
3. 在控制器中添加 Swagger 注解，在 `routes/routes.go` 注册路由
4. 运行 swag 重新生成文档
5. 前端调用 `/api/...`，通配代理会自动带上令牌，不需要新增代理文件

### 前端编码注意

- `eslint-config-next` 16 带 `react-hooks` 7 规则：不要在 `useEffect` 里同步调用 `setState`，不要在渲染中调用 `Date.now()` 等不纯函数。现有做法是在事件处理器里直接更新，或在渲染期间根据 props 变化调整状态
- 构建时 Next 会自动维护 `tsconfig.json`（`jsx: react-jsx`、`include` 里的 `.next/dev/types`），不要手动还原
- 管理员密码不随列表返回，卡片上点击时才请求 `/api/servers/:id/rcon`，不缓存

### 提交前检查

```bash
(cd server && gofmt -l . && go build ./... && go vet ./... && go test ./...)
(cd ui && npm run lint && npx tsc --noEmit && npm run build)
```

## 测试

测试覆盖率很低，目前只有 `server/service/docker_manager/rollback_test.go`。启停竞态、卷清理和前端流程都没有自动化测试。

## 已知问题

- 前端新建表单的 RCON 默认端口是 27020，后端模型默认值是 32330，两者不一致
- `proxy.ts` 在缺少 `NEXT_LOCALE` cookie 时会在登录检查之前直接放行
- 登录接口没有限流，`InitUser` 存在并发创建管理员的竞态
- `README.md` 和 `README-zh.md` 仍把卷名写成 `ark-server-<id>`，实际为 `ase-server-<id>`
- `npm audit` 有来自 ESLint 依赖链（`minimatch`）的中等级别告警，仅影响开发期
- 删除服务器不会清理 Docker 卷，需要手动 `docker volume rm`

## 故障排除

- **后端无法启动**：检查 `JWT_SECRET` 是否设置且不少于 32 字符；`docker ps` 确认 Docker 可用；检查数据库目录权限
- **登录后立即回到登录页**：通过 HTTP 访问却设置了 `COOKIE_SECURE=true`
- **页面提示后端服务不可用（502）**：检查 `API_BASE` 是否指向后端，以及后端进程是否在运行
- **前端构建失败**：删除 `ui/.next` 和 `ui/node_modules` 后执行 `npm ci`
- **Docker 操作失败**：确认挂载了 `/var/run/docker.sock` 并以 `privileged` 运行
