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
- `JWT_SECRET`：可选。未设置时，首次启动自动生成 48 字节随机密钥，保存到数据库所在目录的 `jwt_secret`（权限 0600，镜像内为 `/data/jwt_secret`），之后重启复用。手动设置时至少 32 字符，包含弱口令片段（如 `secret`、`password`）会被拒绝启动
- `DB_PATH`：数据库路径（代码默认 `ark_server.db`，镜像内为 `/data/ark_server.db`）
- `SERVER_PORT`：后端端口（默认 `8080`；容器内由入口脚本固定为 `8081`，`.env` 中的值被忽略）
- `LISTEN_HOST`：监听地址（默认 `127.0.0.1`；容器内固定为 `127.0.0.1`）
- `CONTAINER_TZ`：ARK 服务器容器的时区（默认 `Asia/Shanghai`），只对新建或重建的容器生效
- `CORS_ALLOWED_ORIGINS`：跨域来源白名单，逗号分隔；为空表示仅同源
- `PUBLIC_HOST`：玩家连接用的公网地址（IP 或域名）；为空时自动探测公网 IP（ipify / ifconfig.me，缓存 1 小时）
- `BACKUP_DIR`：存档备份目录（默认数据库所在目录下的 `backups`，镜像内为 `/data/backups`）
- `GIN_MODE`、`LOG_LEVEL`、`LOG_FORMAT`：Gin 模式和日志配置

前端（Next 服务端）：
- `API_BASE`：后端 API 地址，运行时读取，容器内为 `http://127.0.0.1:8081/api`
- `NEXT_PUBLIC_API_BASE`：开发环境回退值；都未设置时回退到 `http://127.0.0.1:8080/api`
- `COOKIE_SECURE`：通过 HTTPS 访问时设为 `true`，纯 HTTP 时保持 `false`，否则无法登录

示例见 `.env.example`。`docker-compose.yml` 中 `.env` 是可选的（`required: false`），不建 `.env` 也能直接启动。

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

- `main.go`：初始化日志、配置、数据库，检查 Docker，启动运行信息监控（`StartRuntimeMonitor`）和定时任务调度（`StartScheduler`），注册 CORS 和路由
- `config/`：环境变量读取，JWT 密钥加载（环境变量优先，否则读取或生成 `jwt_secret` 文件）
- `controllers/{auth,servers,images,system}/`：请求处理。`servers/errors.go` 的 `RespondServiceError` 用 `errors.Is` 把服务层错误映射为 HTTP 状态码，响应带 `error`（中文详情）和 `code`（错误码，见 `service/server/errors.go` 的 `errorCodes`）。`auth/ratelimit.go` 是登录失败限流
- `service/server/`：服务器生命周期和镜像管理。`errors.go` 定义哨兵错误和错误码，不要再用 `err.Error() == ...` 比较字符串；`managedImages` 是允许拉取和更新的镜像白名单。其他文件：`operations.go`（重启、更新游戏、日志、RCON、玩家）、`runtime.go`（运行信息监控）、`backup.go`、`schedule.go`、`system.go`（概览、公网地址、端口建议、创意工坊查询）、`response.go`（`toResponse`）
- `service/docker_manager/`：容器、镜像、卷操作、`RollbackManager`（LIFO 回滚栈），`docker_config.go`（卷内配置文件读写）、`docker_exec.go`（容器内执行 python3 脚本、`RunHelper` 跑一次性 alpine 容器）、`docker_rcon.go`（RCON）、`docker_a2s.go`（查询端口 A2S_INFO）、`docker_logs.go`（日志）、`docker_stats.go`（CPU/内存、卷大小）、`docker_backup.go`（存档导出导入、清除安装标记）
- `middleware/`：`auth.go`（JWT 校验）、`cors.go`（白名单 CORS）
- `models/`：GORM 模型 User、Server、Schedule，以及请求和响应结构。`ServerResponse` 不包含管理员密码，`image_outdated` 表示容器仍在用旧镜像，`runtime` 是运行中服务器的实时信息（阶段、人数、资源）。注意：带 `default:true` 的 bool 字段（`AutoRestart`、`Schedule.Enabled`）创建时 GORM 会把 false 换成 true，需要创建后显式写回
- `utils/`：日志、JWT、密码、命名和镜像名常量（`asm_info.go`，`ArkServerImage`）、配置文件模板（`config_files.go`）
- `docs/`：swag 生成的 Swagger 文件，不要手改

### API 路由（均以 `/api` 开头）

- 公开：`GET /auth/check-init`、`POST /auth/init`、`POST /auth/login`
- 需 JWT：`GET /profile`、`PUT /auth/password`
- 服务器：`GET|POST /servers`、`GET|PUT|DELETE /servers/:id`（`?delete_data=true` 同时删卷和备份）、`POST /servers/:id/start|stop|restart|recreate`、`POST /servers/:id/update-game`（`{game, mods}`）、`GET /servers/:id/rcon`（管理员密码）
- 运维：`GET /servers/:id/logs?tail=`、`GET /servers/:id/logs/stream`（SSE）、`POST /servers/:id/rcon/command`、`GET /servers/:id/players`
- 备份：`GET|POST /servers/:id/backups`、`POST /servers/:id/backups/upload`（multipart `file`）、`GET|DELETE /servers/:id/backups/:name`、`POST /servers/:id/backups/:name/restore`
- 定时任务：`GET|POST /servers/:id/schedules`、`PUT|DELETE /servers/:id/schedules/:scheduleId`、`POST /servers/:id/schedules/:scheduleId/run`
- 系统：`GET /system/overview`（`?refresh=true` 重新计算磁盘）、`GET /system/info`（公网地址）、`GET /system/suggest-ports`、`GET /mods?ids=`
- 镜像：`GET /images/status`、`POST /images/pull`、`GET /images/check-updates`、`POST /images/update`（`{image_name, restart_running}`）、`GET /images/affected`

### 前端结构（`ui/src/`）

- `app/(auth)/`：登录、初始化页面
- `app/(protected)/`：需登录页面（`home` 概览、`servers` 列表、`servers/new`、`servers/[id]` 详情（概览 / 控制台 / 日志 / 备份 / 定时任务，`?tab=` 选择标签）、`servers/[id]/edit`、`settings`）
- `app/api/`：服务端路由。`auth/{login,init,logout}` 处理令牌 cookie；`[...path]` 是通配代理，请求体和响应体都以流转发（支持上传、下载和 SSE），并转发 `X-Forwarded-For`
- `proxy.ts`：Next 16 中原 `middleware.ts` 的新约定。负责语言 cookie（`NEXT_LOCALE`）和登录重定向
- `lib/backend-server.ts`：`getApiBase`、令牌 cookie 读写、`authProxy`
- `lib/api-error.ts`：`useApiError()`（中文界面显示后端 `error`，其他语言按 `apiErrors.<code>` 翻译）和旧的 `getErrorMessage`；`lib/toast.ts`：全局通知（`toast.success/error`）；`lib/format.ts`：大小、时间、随机密码、复制；`hooks/use-polling.ts`：页面可见时轮询；`lib/ini.ts`：`mergeIni`，INI 文本的就地合并
- `stores/`：Zustand（`auth.ts`、`servers.ts`），HTTP 用 Axios
- `components/ui/`：shadcn/ui 风格基础组件，另有 `confirm-dialog.tsx`、`toaster.tsx`；`components/servers/`：`ServerCard`、`ServerActions`（启停按钮、更多菜单和全部确认弹窗，卡片和详情页共用）、`StatusBadge`（状态 + 运行阶段）、`ServerForm`（新建和编辑共用）、`LogViewer`、`ConsolePanel`、`BackupsPanel`、`SchedulesPanel`、`ModList`、`ConnectionInfo`、ini 可视化编辑器；`components/docker/`：镜像状态和更新确认弹窗（在 `home` 页面使用）
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
- 服务端文件卷：名称见 `utils.GetServerFilesVolumeName`
- 镜像：`tbro98/ase-server:latest`（子模块 `ASE-Server-Docker` 是该镜像的源码）

容器挂载三个卷：服务端文件卷挂到 `/home/steam/arkserver`，数据卷嵌套挂到 `ShooterGame/Saved`，插件卷嵌套挂到 `ArkApi/Plugins`。设置了 `cluster_id` 时再挂共享的集群卷 `ase-cluster-<cluster_id>` 到 ARK 默认集群目录 `Saved/clusters`（容器环境变量 `CLUSTER_VOLUME` 记录挂载的卷，用于判断是否需要重建）。镜像不再预下载游戏本体，入口脚本在首次启动时把服务端文件下载到文件卷；安装标记 `.install_complete` 和已装 Mod 记录 `.installed_mods` 也在卷里，所以重建容器不会重新下载，新增 Mod 时只下载缺少的。ArkApi 随镜像放在 `/home/steam/arkapi`，每次启动同步到卷里（已存在的 `config.json` 不覆盖）。镜像侧的改动需要重新构建并推送镜像后才生效。

### 服务器生命周期

1. **创建**：先检查端口冲突（`checkPorts`，含游戏端口+1），再在事务内写数据库记录、创建三个卷（服务端文件卷、数据卷、插件卷）并写入默认配置文件；容器不在创建时生成，首次启动时才创建。事务未提交时 `defer` 清理已创建的卷（`committed` 标记）
2. **启动 / 停止**：用 `transitionStatus` 原子切换状态（`UPDATE ... WHERE status NOT IN (...)`，看 `RowsAffected`），避免并发请求重复启停；容器操作在后台协程中完成，结果由 `finishOperation` 写入状态和 `last_error`。启动前先把暂存的配置文件写入数据卷；服务端文件卷里还没有服务端主程序时，`seed.go` 从另一台已安装完成（有 `.install_complete` 和主程序，优先已停止）的服务器复制整个文件卷（含 Mod 和 `.installed_mods`），再删掉 `.install_complete`，入口脚本启动时 validate 只下载差异，复制期间 `runtime.phase` 为 `copying`，复制失败只记日志、照常从 Steam 下载；开启 `auto_update` 时清除安装标记；如果 `SERVER_ARGS`、`GameModIds`、`CLUSTER_VOLUME`、重启策略（`auto_restart`）与当前配置不同，或容器的镜像 ID 不是 `ArkServerImage` 当前指向的镜像，会先重建容器；启动超时会停掉容器。停止前先通过容器内 RCON（`docker_rcon.go`，用 `docker exec` 运行 python3）执行 `SaveWorld`，失败只记日志、不阻止停止；入口脚本会把 SIGTERM 转发给游戏进程（需要重新构建镜像才生效）。`GenerateArgs` 对 map 按 key 排序，保证相同配置生成相同参数；参数同时以数组形式放进容器 `Cmd`（入口脚本优先用 `$@`，名称和密码里的空格不会被拆开），`SERVER_ARGS` 只用于重建判断和兼容旧镜像。管理员密码不在启动参数和环境变量里，每次启动前由 `syncAdminPassword` 以数据库里的值写入 `GameUserSettings.ini` 的 `[ServerSettings]`（读取失败时不会覆盖整个文件，只有 `ErrConfigNotFound` 才视为空文件）。`GetServers` 通过 `syncStatusWithDocker` 只校正 running / stopped，不覆盖 starting / stopping；后端启动时 `ReconcileStatuses` 校正被中断的启停操作
3. **运行阶段**：容器 running 不代表能连接。`runtime.go` 每 10 秒对运行中的服务器读取最近 400 行日志判断阶段（`downloading` 含 steamcmd 进度、`installing_mods` 含 x/y、`launching`），再在容器内发 A2S_INFO，有响应即 `ready` 并拿到在线人数；同时采样 CPU/内存。结果缓存在内存，随服务器列表返回 `runtime`。阶段标志行来自 `ASE-Server-Docker/scripts/start_server.sh`，改脚本输出时要同步改 `applyPhase`
4. **重启 / 更新游戏**：`restartWith` 先停止（如果在运行）、执行可选的 `beforeStart`、再启动；更新游戏的 `beforeStart` 是删除文件卷里的 `.install_complete` / `.installed_mods`，入口脚本会重新 `app_update validate` 并重新下载 Mod
5. **重建**：`/servers/:id/recreate` 先原子占位为 stopping，再在后台执行：运行中的服务器先停止（失败则放弃），写入暂存配置，按当前配置创建新容器，结束时为 stopped，失败原因写入 `last_error`（服务器卡片上有重建按钮）
6. **删除**：只允许删除已停止的服务器（starting / stopping 返回 409）；先删容器，再以 `status = stopped` 为条件软删除数据库记录，同时删除定时任务。默认保留 Docker 卷和备份；`delete_data=true` 时删除三个卷和备份目录（集群卷共享，不删）

### 存档备份

备份放在 `BACKUP_DIR/<server_id>/<时间>_<类型>.tar.gz`（类型：manual / scheduled / pre-restore / uploaded，时间按 `CONTAINER_TZ`），归档内是 `Saved/`（跳过 `Logs` 和 `clusters`）。运行中备份先 RCON `SaveWorld`。恢复只能在已停止时进行：先自动做 pre-restore 备份，占位为 stopping，用 `RunHelper` 清空 `SavedArks` 和 `Config`，再 `CopyToContainer` 解压，清除暂存配置。上传的归档要通过 `ValidateBackup`（gzip、只允许 `Saved/` 下的普通文件和目录、必须有 `SavedArks`）。同一服务器的备份 / 恢复由 `backupLocks` 串行化

### 定时任务

`schedules` 表，`schedule.go` 每分钟整点检查（`robfig/cron` 5 段表达式或 `@daily` 等，时区 `CONTAINER_TZ`，系统休眠后最多补 5 分钟）。restart / stop / update 可设 `warn_minutes`，在 N、10、5、1 分钟各广播一次（`message` 中的 `{minutes}` 会替换）；backup 的 `keep` 只清理 scheduled 类型的旧备份。定时重启只处理运行中的服务器

### 镜像更新

1. **检查更新**：`CheckImageUpdate` 用 `DistributionInspect` 取镜像仓库中该标签的 manifest digest，与本地镜像的 `RepoDigests` 比较，不下载镜像层。结果记录在内存中，由 `/images/status` 的 `has_update` 返回，拉取成功后清除。本地构建、没有 `RepoDigests` 的镜像会一直显示有更新
2. **更新**：先用 `ClaimPull` 原子占用拉取权（并发请求返回 409），再后台拉取新镜像，失败原因记录在镜像状态的 `error` 字段并显示在首页；拉取成功且镜像 ID 变化时会清理被替换的旧镜像（不强删，仍被容器引用时只记日志），拉取进度流中的 `error` 会被当作失败。已停止的服务器在下次启动时自动重建；`restart_running` 为 true 时，运行中且镜像已过期的服务器会依次经过 stopping → stopped → starting，启动时重建容器
3. 前端：首页"镜像管理"检查更新，点击"更新"先请求 `/images/affected`，再弹出 `ImageUpdateConfirmModal` 选择是否重启运行中的服务器；服务器卡片根据 `image_outdated` 显示"旧镜像"标记

### 配置文件

GameUserSettings.ini 和 Game.ini 不存数据库（暂存内容除外），而是直接保存在数据卷的 `Config/WindowsServer/` 下：`docker_config.go` 创建一个挂载数据卷、不启动的临时 `alpine:latest` 容器，读用 `CopyFromContainer`，写用 `CopyToContainer`（tar 写入，不经过 shell）。所以 `alpine:latest` 也是必需镜像。服务器未停止时保存的配置先存入 `pending_game_user_settings` / `pending_game_ini`（ARK 停服时会覆盖 GameUserSettings.ini），下次启动或重建前写入数据卷，响应里 `config_pending` 为 true。前端的可视化编辑器（`GameUserSettingsEditor`、`GameIniEditor`）和文本模式互相同步，可视化修改通过 `lib/ini.ts` 的 `mergeIni` 在原文本上就地合并，只写入改过的键，未知键、注释和自定义 section 原样保留。

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
- 详情页用列表接口 `/api/servers` 轮询，不要轮询 `GET /servers/:id`（它会用临时容器读取配置文件，开销大）
- 翻译用 ICU 格式：文案里不能直接写 `<xxx>`（会被当成标签）或不传参数的 `{xxx}`；新增键时中英文都要加
- 错误提示用 `useApiError()` + `toast`，不要在 store 里写死中文

### 提交前检查

```bash
(cd server && gofmt -l . && go build ./... && go vet ./... && go test ./...)
(cd ui && npm run lint && npx tsc --noEmit && npm run build)
```

## 测试

测试覆盖率仍然较低：`server/service/docker_manager/rollback_test.go`、`docker_image_test.go`（`ClaimPull` 并发占用）、`docker_backup_test.go`（备份格式校验、CPU 计算）、`server/config/config_test.go`（JWT 密钥生成和复用）、`server/models/server_args_test.go`（启动参数顺序确定）、`server/service/server/server_service_test.go`（端口冲突、保存不覆盖状态、端口建议、禁用的定时任务，使用临时 SQLite）、`parse_test.go`（运行阶段、ListPlayers 解析、字段规范化、定时任务校验）、`server/controllers/auth/ratelimit_test.go`、`server/utils/ini_test.go`（`SetIniValue`）。启停流程、RCON、配置暂存、备份恢复和前端流程都没有自动化测试。Windows 上 `config_test.go` 的 0600 权限检查会失败（Windows 没有这种权限位），以 Linux 结果为准。本机没有 Go 时，可以用 `docker run --rm -v $PWD/server:/src -w /src golang:1.27-alpine go test ./...` 在容器里跑。

## 已知问题

- 登录限流按 IP（`X-Forwarded-For`，由 Next 转发）计数，客户端能伪造该头绕过 IP 限制，用户名维度的上限仍然有效；限流记录在内存，重启后清零
- 修改登录密码后，旧 JWT 在 24 小时有效期内仍然可用
- `npm audit` 有来自 ESLint 依赖链（`minimatch`）的中等级别告警，仅影响开发期
- 入口脚本已转发 SIGTERM，但游戏进程能否在 `docker stop` 的 30 秒内完成存档，没有在真实环境验证过，停服前的 `SaveWorld` 继续保留。steamcmd 下载阶段还没设置 trap，停止下载中的服务器要等满 30 秒（已实测）
- 集群卷不会随服务器删除；不再使用的集群卷需要手动 `docker volume rm ase-cluster-<id>`。已有服务器改成集群后，原数据卷里 `Saved/clusters` 的旧数据会被集群卷遮住
- 运行阶段依赖入口脚本的日志文本；游戏本身的 `ShooterGame.log` 在数据卷 `Saved/Logs` 下，界面上看不到
- 创意工坊名称查询依赖访问 `api.steampowered.com`，访问不了时只显示 ID
- 管理员密码现在存放在卷里的 `GameUserSettings.ini`，`docker inspect` 看不到；数据库里仍是明文，并且每次启动都以数据库的值为准，在 ini 里直接改 `ServerAdminPassword` 会被覆盖
- 每台服务器各有一份完整的服务端文件卷，磁盘占用大（新服务器从已安装的服务器复制，省的是下载，不省磁盘）；Mod 平时只下载尚未安装的，已装 Mod 需要用"更新游戏/Mod"或开启 `auto_update` 才会重新下载
- 默认 Game.ini 模板里的倍率类 key（如 `XPMultiplier`、`TamingSpeedMultiplier`）可能应该放在 GameUserSettings.ini 的 `[ServerSettings]` 下，未验证
- `GetServers` 现在一次 `ContainerList` 取状态，依赖容器名匹配 `ase-server-<id>`；`GetImageStatus` 的 `overall_status` 是状态码（`ready` / `pulling` / `not_ready`），`can_create_server` 取决于 alpine 镜像是否就绪

## 故障排除

- **后端无法启动**：`docker ps` 确认 Docker 可用；检查数据库目录是否可写（需要写入 `jwt_secret`）；如果手动设置了 `JWT_SECRET`，确认不少于 32 字符且不含弱口令片段
- **所有用户被登出**：`jwt_secret` 文件被删除或数据目录没有持久化，重启后生成了新密钥
- **登录后立即回到登录页**：通过 HTTP 访问却设置了 `COOKIE_SECURE=true`
- **页面提示后端服务不可用（502）**：检查 `API_BASE` 是否指向后端，以及后端进程是否在运行
- **前端构建失败**：删除 `ui/.next` 和 `ui/node_modules` 后执行 `npm ci`
- **Docker 操作失败**：确认挂载了 `/var/run/docker.sock` 并以 `privileged` 运行
- **服务器显示"运行中"但连不上**：看卡片上的阶段（下载服务端 / 安装 Mod / 加载中），详情页"日志"标签能看到实时进度；首次启动要下载约 20 GB
- **连接地址不对**：设置 `PUBLIC_HOST`
- **玩家版本不匹配进不了服**：服务器菜单里"更新游戏/Mod"
