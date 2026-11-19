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

`docker-compose.yml` 是正式部署文件，拉取已发布的镜像 `tbro98/arkservercommander:latest`，不读取 `.env`（需要的环境变量直接写进 `environment`）。

```bash
docker compose up -d       # 启动
docker compose logs -f     # 查看日志
docker compose down        # 停止
docker build -t arkservercommander:dev .   # 需要验证镜像时手动构建
```

### 本地开发运行

`bash dev.sh` 在本机同时运行后端（`go run .`，`127.0.0.1:8081`）和前端开发服务器（`next dev`，`http://localhost:3000`），不构建镜像，Ctrl+C 同时停止。默认值：`DB_PATH=data-dev/ark_server.db`（数据、`jwt_secret`、备份和 ARK 服务器文件都在 `data-dev/`，已加入 `.gitignore`）、`GIN_MODE=debug`、`LOG_FORMAT=console`，前端的 `API_BASE` 指向后端端口；都可以在运行前 export 覆盖，端口用 `SERVER_PORT` / `UI_PORT`。脚本不读取 `.env`（`.env.example` 里的 `DB_PATH=/data/...` 是容器内路径）。

Windows 上用 PowerShell 运行 `.\dev.ps1`（执行策略不允许时 `powershell -ExecutionPolicy Bypass -File .\dev.ps1`），默认值和 `dev.sh` 相同，用 `$env:XXX = '...'` 覆盖。需要 Docker Desktop（Linux 容器）：数据目录用 `E:/...` 这样的 Windows 路径作为绑定卷的 `device`，Docker Desktop 能识别，不用设置 `HOST_DATA_DIR`；脚本按当前 `docker context` 设置 `DOCKER_HOST`（后端的 Docker SDK 不读 context）；前端直接用 `node` 运行 next，不经过 `npx.cmd`；退出时用 `taskkill /T` 结束两个进程树。`dev.ps1` 要保存为带 BOM 的 UTF-8，否则 Windows PowerShell 5.1 会把中文读成乱码。ARK 文件在 Windows 磁盘上，经文件共享挂进容器，读写较慢。

注意：ARK 服务器容器名是 `ase-server-<id>`，不要和正式部署在同一个 Docker 上同时使用；ARK 容器写入 `data-dev/servers/` 的文件属主是 root，在 Linux 上以普通用户运行后端时，写配置文件、恢复备份、复制服务端文件和删除服务器数据都可能因权限失败（正式部署的后端容器以 root 运行；Windows 的 Docker Desktop 文件共享没有这个问题）。

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
- `HOST_DATA_DIR`：数据目录在宿主机上的绝对路径，可选。ARK 服务器的卷绑定到数据目录下，Docker 需要宿主机路径；为空时后端检查自身容器的挂载自动探测（目标为数据目录的那一项的 Source），探测失败会拒绝启动并提示设置它；后端直接运行在宿主机上时使用数据目录本身
- `ARKFGF_API`：arkfgf.com 接口地址（默认 `https://arkfgf.com/api/v1`），设为空或 `off` 时关闭发布功能和 P2P 联机房，见下文"发布到收录站""P2P 联机房"
- `GIN_MODE`、`LOG_LEVEL`、`LOG_FORMAT`：Gin 模式和日志配置

前端（Next 服务端）：
- `API_BASE`：后端 API 地址，运行时读取，容器内为 `http://127.0.0.1:8081/api`
- `NEXT_PUBLIC_API_BASE`：开发环境回退值；都未设置时回退到 `http://127.0.0.1:8080/api`
- `COOKIE_SECURE`：通过 HTTPS 访问时设为 `true`，纯 HTTP 时保持 `false`，否则无法登录

示例见 `.env.example`。`docker-compose.yml` 和 `dev.sh` 都不读取 `.env`。

## 架构概览

### 部署形态

单个容器内运行两个进程，由 `docker-entrypoint.sh` 启动：

```
浏览器 → Next.js（:8080，对外）→ Go 后端（127.0.0.1:8081，仅容器内）→ Docker Socket → ARK 服务器容器
Agent（MCP 客户端）→ Next.js /mcp → Go 后端 /api/mcp
```

- Go 后端不托管前端静态文件，只提供 `/api`、`/health` 和 `/swagger`
- 任一进程退出，入口脚本会结束另一个并以非零状态退出，交给 Docker 重启策略处理
- Swagger 位于后端的 `/swagger/index.html`，容器内不对外暴露，本地直接运行后端时才能访问
- 容器需要挂载 `/var/run/docker.sock` 并以 `privileged` 运行

### 后端结构（`server/`）

分层：controllers（HTTP）→ service（业务）→ docker_manager（Docker SDK）。

- `main.go`：初始化日志、配置、数据库，检查 Docker，启动运行信息监控（`StartRuntimeMonitor`）、定时任务调度（`StartScheduler`）、Mod 更新检查（`StartModUpdateChecker`）、收录站发布（`StartPublisher`）和 P2P 联机房（`StartLinkManager`），注册 CORS 和路由
- `config/`：环境变量读取，JWT 密钥加载（环境变量优先，否则读取或生成 `jwt_secret` 文件）
- `controllers/{auth,servers,images,system}/`：请求处理。`servers/errors.go` 的 `RespondServiceError` 用 `errors.Is` 把服务层错误映射为 HTTP 状态码，响应带 `error`（中文详情）和 `code`（错误码，见 `service/server/errors.go` 的 `errorCodes`）。`auth/ratelimit.go` 是登录失败限流
- `service/server/`：服务器生命周期和镜像管理。`errors.go` 定义哨兵错误和错误码，不要再用 `err.Error() == ...` 比较字符串；`managedImages` 是允许拉取和更新的镜像白名单。其他文件：`operations.go`（重启、更新游戏、日志、RCON、玩家）、`runtime.go`（运行信息监控）、`backup.go`、`schedule.go`、`system.go`（概览、公网地址、端口建议、创意工坊查询）、`mods.go`（Mod 安装状态和更新检查，后台每 30 分钟）、`plugins.go`（ArkApi 插件）、`arkapi.go`（ArkApi 版本库和安装）、`images.go`（启动时自动拉取镜像、哪些镜像是必需的）、`gamelog.go`（读取 ShooterGame.log）、`server_config.go`（配置导出、导入、复制到其他服务器）、`config_values.go`（`SetConfigValues`）、`clusters.go`（未使用的集群）、`version.go`（GitHub 最新版本）、`publish.go`（发布到 arkfgf.com 收录站）、`link.go`（P2P 联机房）、`profile.go`（发布和联机房共用的服务器资料 `buildServerProfile`）、`arkfgf.go`（arkfgf.com 的 JSON 请求、错误码和 Retry-After）、`response.go`（`toResponse`）
- `service/docker_manager/`：容器、镜像、卷操作、`RollbackManager`（LIFO 回滚栈），`docker_config.go`（卷内配置文件读写）、`docker_exec.go`（容器内执行命令 `ExecCommand`，返回 stdout / stderr / 退出码；`execPython` 基于它）、`docker_link.go`（联机房节点容器的创建、检查、删除、列出，数据目录 `servers/<id>/link`）、`docker_rcon.go`（RCON）、`docker_a2s.go`（查询端口 A2S_INFO）、`docker_logs.go`（日志）、`docker_stats.go`（CPU/内存、卷大小）、`docker_backup.go`（存档导出导入、清除安装标记）、`docker_seed.go`（检查和复制服务端文件）、`localfs.go`（`copyTree`，相当于 `cp -a`）、`installed_mods.go`（读写 `.installed_mods`）、`clusters.go`（列出和删除集群卷）、`storage.go`（数据目录布局，`ServerFilesDir` 等返回后端本地路径。卷里的所有文件操作——配置文件、备份导出恢复、安装标记、复制服务端文件、插件、Mod 记录、游戏日志——都直接读写这些目录，不再用临时容器，也不需要 alpine 镜像；目录内容由 ARK 容器写入，访问时用 `os.Root`，符号链接不能指到目录外）
- `middleware/`：`auth.go`（JWT 校验）、`cors.go`（白名单 CORS）
- `mcpserver/`：MCP 服务（官方 `github.com/modelcontextprotocol/go-sdk`），见下文"MCP 接口"
- `service/apitoken/`：MCP 用的 API 令牌（创建、列出、撤销、校验），`controllers/auth/tokens.go` 是对应接口
- `models/`：GORM 模型 User、Server、Schedule，以及请求和响应结构。`ServerResponse` 不包含管理员密码，`image_outdated` 表示容器用的不是服务器当前应使用的镜像（镜像更新过或切换了 ArkApi），`arkapi` 决定镜像，`image` 是由它得出的镜像名，`runtime` 是运行中服务器的实时信息（阶段、人数、资源），`mods_outdated` 是创意工坊上有更新的 Mod。注意：带 `default:true` 的 bool 字段（`AutoRestart`、`Schedule.Enabled`）创建时 GORM 会把 false 换成 true，需要创建后显式写回
- `utils/`：日志、JWT、密码、命名和镜像名常量（`asm_info.go`，`ArkServerLinuxImage` / `ArkServerProtonImage`，`ArkServerImageFor(arkApi)` 选镜像，联机房节点容器名 `GetLinkContainerName` 和镜像 `EasyTierImage`）、配置文件模板（`config_files.go`）、版本号（`version.go`，构建时用 `-ldflags -X ark-server-commander/utils.Version=` 注入，默认 `dev`）
- `docs/`：swag 生成的 Swagger 文件，不要手改

### API 路由（均以 `/api` 开头）

- 公开：`GET /auth/check-init`、`POST /auth/init`、`POST /auth/login`
- 需 JWT：`GET /profile`、`PUT /auth/password`、`GET|POST /tokens`、`DELETE /tokens/:id`（MCP 令牌，只能用 JWT 管理）
- MCP：`POST /mcp`（API 令牌认证，不接受 JWT；无状态模式下 GET / DELETE 返回 405）
- 服务器：`GET|POST /servers`、`GET|PUT|DELETE /servers/:id`（`?delete_data=true` 同时删卷和备份）、`POST /servers/:id/start|stop|restart|recreate`、`POST /servers/:id/update-game`（`{game, mods, mod_ids}`）、`GET /servers/:id/rcon`（管理员密码）
- 配置：`GET /servers/:id/config/export`（导出，不含端口和管理员密码）、`POST /servers/:id/config/import`（`{config, parts}`）、`POST /servers/:id/config/copy`（`{target_ids, parts}`，逐台返回结果），见下文"配置导入导出"
- 运维：`GET /servers/:id/logs?tail=`、`GET /servers/:id/logs/stream`（SSE）、`GET /servers/:id/logs/game?tail=`（ShooterGame.log）、`POST /servers/:id/rcon/command`、`GET /servers/:id/players`、`GET /servers/:id/mods?refresh=`（Mod 安装和更新状态）、`GET /servers/:id/publish/stats`（收录站访问统计）
- 联机房：`GET /servers/:id/link/members`、`POST /servers/:id/link/kick`（`{steam_id, ban}`）、`DELETE /servers/:id/link/bans/:steam_id`，见下文"P2P 联机房"
- 插件：`GET|POST /servers/:id/plugins`（multipart `file`，可选 `name`）、`PUT|DELETE /servers/:id/plugins/:name`（`{enabled}`）、`GET|PUT /servers/:id/plugins/:name/config`、`GET|PUT /servers/:id/arkapi/config`（`{content}`）、`POST /servers/:id/arkapi/install`（`{version}`，为空装最新版本）
- ArkApi 版本库：`GET /arkapi/versions`（`?refresh=true` 重新查询 GitHub）、`POST /arkapi/versions`（multipart `file`，可选 `version`）、`DELETE /arkapi/versions/:version`
- 备份：`GET|POST /servers/:id/backups`、`POST /servers/:id/backups/upload`（multipart `file`）、`GET|DELETE /servers/:id/backups/:name`、`POST /servers/:id/backups/:name/restore`
- 定时任务：`GET|POST /servers/:id/schedules`、`PUT|DELETE /servers/:id/schedules/:scheduleId`、`POST /servers/:id/schedules/:scheduleId/run`
- 系统：`GET /system/overview`（`?refresh=true` 重新计算磁盘，含 `unused_clusters`）、`GET /system/info`（公网地址）、`GET /system/version`（`?refresh=true`）、`DELETE /system/clusters/:clusterId`、`GET /system/suggest-ports`、`GET /system/server-defaults`（新建服务器的默认配置文件和启动参数）、`GET /mods?ids=`
- 镜像：`GET /images/status`、`POST /images/pull`、`GET /images/check-updates`、`POST /images/update`（`{image_name, restart_running}`）、`GET /images/affected`

### 前端结构（`ui/src/`）

- `app/(auth)/`：登录、初始化页面
- `app/(protected)/`：需登录页面（`home` 概览、`servers` 列表（卡片 / 列表两种视图，选择存在 localStorage 的 `servers-view`）、`servers/new`、`servers/[id]` 详情（概览 / 控制台 / 日志 / 备份 / 定时任务 / 插件 / 管理，`?tab=` 选择标签）、`servers/[id]/edit`、`settings`）
- `app/api/`：服务端路由。`auth/{login,init,logout}` 处理令牌 cookie；`[...path]` 是通配代理，请求体和响应体都以流转发（支持上传、下载和 SSE），并转发 `X-Forwarded-For`；它拒绝 `/api/mcp`（会带上登录 cookie，令牌无效时还会清 cookie）
- `app/mcp/route.ts`：MCP 入口，原样转发到后端 `/api/mcp`，透传 `Authorization`、`Accept`、`Mcp-*` 头，不读 cookie。`proxy.ts` 的 matcher 排除了 `mcp`
- `proxy.ts`：Next 16 中原 `middleware.ts` 的新约定。负责语言 cookie（`NEXT_LOCALE`）和登录重定向
- `lib/backend-server.ts`：`getApiBase`、令牌 cookie 读写、`authProxy`
- `lib/api-error.ts`：`useApiError()`（中文界面显示后端 `error`，其他语言按 `apiErrors.<code>` 翻译）和旧的 `getErrorMessage`；`lib/toast.ts`：全局通知（`toast.success/error`）；`lib/format.ts`：大小、时间、随机密码、复制；`hooks/use-polling.ts`：页面可见时轮询；`lib/ini.ts`：INI 的按 section 读取（`indexIni` / `getIniValue`）和就地修改（`setIniValue`，值为 undefined 时删除键）；`lib/server-settings.ts`：服务器设置字典（见下文"服务器设置界面"）；`lib/images.ts`：镜像名常量和显示名称
- `stores/`：Zustand（`auth.ts`、`servers.ts`），HTTP 用 Axios
- 视觉风格（方舟游戏内界面风格的深色 HUD）：根节点固定 `className="dark"`，`globals.css` 定义主题变量（深蓝底色、半透明面板 `--card`、青色 `--primary`、琥珀色 `--amber` 表示进行中、终端底色 `--terminal`）。`@theme` 里把 Tailwind 的 `gray/red/green/amber/yellow/blue/sky/purple/orange` 色阶整体**反转**：数字小的是深色背景色调，数字大的是浅色文字，所以 `bg-red-50 text-red-800` 这类写法在深色下自然是"深红底 + 浅红字"。新代码照此约定：背景用 50–200，图标和圆点用 400–600，文字用 600–900；不要用 `bg-white`、`text-gray-200` 这类按浅色主题思路写的类。自定义工具类（用 `@utility` 定义，可以加变体前缀）：`hud`（左上、右下青色角标，元素本身要有定位）、`hud-caps`（英文大写加字距，用于导航、标题和主要按钮，中文不受影响）、`glow-text`、`font-display`（Saira，来自 `@fontsource-variable/saira`，本地打包）、`num`（等宽数字）。不使用游戏的图片、Logo 或字体等素材，只借鉴界面风格。状态徽章三种颜色：绿色可用、琥珀色进行中、石色已停止
- 地图颜色：`lib/maps.ts` 的 `mapColor` 给每张地图一个固定颜色，服务器卡片、详情页顶部（`MapBand`，停止时变灰）、首页名单和地图选择器都用它；`hooks/use-map-name.ts` 取地图显示名
- 公共组件：`PageHeader`（页面标题行）、`AuthShell`（登录 / 初始化外框，含 `FormError`）、`Wordmark`、`components/ui/page-loader.tsx`（统一加载状态）；`TabsList variant="line"` 用于页面级标签（详情页、服务器表单）
- `components/ui/`：shadcn/ui 风格基础组件，另有 `confirm-dialog.tsx`、`toaster.tsx`；`components/servers/`：`ServerCard`（卡片视图，含共用的 `MapBand`、`ServerTags`）、`ServerRow`（列表视图的一行和表头，列宽固定在 `ROW_GRID`）、`ServerActions`（`useServerControls` 管理启停、重启、更新、重建、删除和这些确认弹窗，调用方渲染 `controls.dialogs`；`PowerButtons` 启停和重启，`ServerShortcuts` 控制台 / 日志 / 编辑设置图标入口，`IconAction` 带悬停提示的图标按钮。卡片和列表不再有"更多"菜单，常用操作直接显示，不常用和有风险的操作在详情页"管理"标签）、`ServerManagePanel`（"管理"标签：更新游戏 / Mod、重建容器、配置导出导入复制、删除，每项带说明）、`StatusBadge`（状态 + 运行阶段）、`ServerForm`（新建和编辑共用）、`LogViewer`、`ConsolePanel`、`BackupsPanel`、`SchedulesPanel`、`PluginsPanel`、`ModList`（表单里可拖动排序）、`ServerModsPanel`（详情页 Mod 安装和更新）、`ServerConfigDialogs`（导入配置、复制设置到其他服务器）、`ConnectionInfo`、`settings/`（设置项卡片 `SettingTile` / `SettingControl`、自定义启动参数和配置文件面板）；`components/settings/VersionCard`（版本和升级命令）、`components/settings/McpCard`（MCP 地址和令牌管理）；`components/docker/`：镜像状态和更新确认弹窗（在 `home` 页面使用）
- `messages/en.ts`、`messages/zh.ts`：next-intl 翻译（TypeScript 文件，不是 JSON），入口为 `src/i18n.ts`

### 认证流程

1. 首次启动没有用户时，前端进入 `/init` 创建管理员
2. 登录或初始化请求先到 Next 的 `app/api/auth/*`，转发给后端，拿到 JWT 后写入 HttpOnly cookie `auth-token`，响应体里不含令牌
3. 其他 `/api/*` 请求经通配代理：从 cookie 取令牌，转成 `Authorization: Bearer <token>` 发给后端；后端返回 401 时清除 cookie
4. 前端脚本读不到令牌，页面加载时通过 `restoreSession`（请求 `/api/profile`）确认登录状态；退出登录调用 `/api/auth/logout`
5. `proxy.ts` 只检查 cookie 是否存在来做页面重定向，令牌真伪由后端校验

## MCP 接口

让 Agent 脱离前端管理服务器。`mcpserver/` 用 `mcp.NewStreamableHTTPHandler`（`Stateless` + `JSONResponse`）挂在后端 `/api/mcp`，外层 `auth.RequireBearerToken` 用 `apitoken.Verify` 校验令牌。

- 令牌：`asc_` + 32 字节随机数（base64url），表 `api_tokens` 只存 SHA-256 和前几位 `prefix`，明文只在创建时返回一次；没有过期时间，撤销即失效；令牌所属用户不存在时无效；`last_used_at` 最多每分钟写一次
- 权限：启动时建两个 `mcp.Server`，按令牌的 scope 选择。只读令牌的服务器只注册 `query` 工具；`action`（修改类）和 `secret`（如 `get_admin_password`）只注册到完全控制的服务器，只读令牌看不到也调用不了
- 工具直接调用 service 层（不经过 HTTP），用户 ID 来自 `req.Extra.TokenInfo.UserID`。输入用 MCP 专用结构（带 `omitempty` 和 `jsonschema` 描述，jsonschema-go 把没有 `omitempty` 的字段当必填），再转成 `models` 请求并用 `binding.Validator.ValidateStruct` 套用 REST 的同一套校验。服务层错误返回为 `[错误码] 中文信息`
- 列表结果包在对象里（`{"servers": [...]}`），结构化输出是对象
- 不提供文件上传下载（插件 zip、备份上传 / 下载、ArkApi 压缩包上传）
- `set_config_values` 调用 `ServerService.SetConfigValues`：以暂存配置或卷里的文件为基础用 `utils.SetIniValue` 改键，读取失败直接报错（不能当空文件覆盖）；同名键会合并，需要重复的键（如 `OverrideNamedEngramEntries`）要用 `update_server` 传完整内容
- 新增工具时：在 `mcpserver/tools_*.go` 用 `query` / `secret` / `action(additive|idempotent|destructive)` 注册，描述用英文写清楚异步行为和前置条件；服务器说明在 `mcpserver.go` 的 `instructions`

## Docker 集成

命名规范（见 `server/utils/asm_info.go`）：
- 容器：`ase-server-<server_id>`
- 数据卷：`ase-server-<server_id>`
- 插件卷：`ase-server-plugins-<server_id>`
- 服务端文件卷：名称见 `utils.GetServerFilesVolumeName`
- P2P 联机房节点容器：`ase-link-<server_id>`（镜像 `easytier/easytier:v2.6.4`，见下文"P2P 联机房"）
- 镜像：`tbro98/ase-server:linux`（原生 Linux 版服务端，默认）和 `tbro98/ase-server:proton`（GE-Proton 运行 Windows 版，`arkapi` 为 true 的服务器用），子模块 `ASE-Server-Docker` 的同一个 dockerfile 用 `--target linux|proton` 构建，入口脚本按 `SERVER_PLATFORM` 区分。`latest` 不再更新，停在最后一个自带 ArkApi 的版本，留给旧版管理端

卷是绑定到数据目录的命名卷（local 驱动，`type=none, o=bind`，见 `docker_manager/storage.go`、`docker_volume.go`）：卷名不变，数据实际放在后端容器的 `/data`（宿主机 `./data`）下，宿主机上可以直接管理文件：`servers/<id>/files`（服务端文件）、`servers/<id>/saved`（存档和配置）、`servers/<id>/plugins`（插件）、`clusters/<cluster_id>`（集群共享）。后端启动时 `InitStorage` 确定数据目录的宿主机路径，失败则退出。删除卷不会删除目录里的文件，需要再调用 `RemoveServerDirs`；集群卷在创建容器前由 `EnsureClusterVolume` 建好，否则 Docker 会自动创建一个普通卷。磁盘占用由 `LocalVolumeSizes` 遍历目录统计（Docker 的 DiskUsage 查不到绑定卷的大小），key 仍是卷名。

容器挂载三个卷：服务端文件卷挂到 `/home/steam/arkserver`，数据卷嵌套挂到 `ShooterGame/Saved`，插件卷嵌套挂到 `ArkApi/Plugins`。设置了 `cluster_id` 时再挂共享的集群卷 `ase-cluster-<cluster_id>` 到 ARK 默认集群目录 `Saved/clusters`（容器环境变量 `CLUSTER_VOLUME` 记录挂载的卷，用于判断是否需要重建）。镜像不再预下载游戏本体，入口脚本在首次启动时把服务端文件下载到文件卷；安装标记 `.install_complete` 和已装 Mod 记录 `.installed_mods`（每行 `Mod ID 安装时间`，旧版只有 ID，后端用 `Content/Mods/<id>.mod` 的修改时间代替）也在卷里，所以重建容器不会重新下载，新增 Mod 时只下载缺少的。`.install_complete` 的内容是安装时的平台（`linux` / `windows`，旧版是空文件，视为 windows），切换 ArkApi 换了镜像后平台不同，入口脚本会重新 validate，补下载这个平台的主程序（Mod 两个平台都用 WindowsNoEditor，不用重下）。Linux 版服务端读 `Config/LinuxServer`，入口脚本把它建成指向 `WindowsServer` 的符号链接，管理端只维护 `WindowsServer`；备份导出时跳过符号链接。ArkApi 不随镜像发布（`assets/arkapi/ArkApi_3.56.zip` 内置在后端里），由 `arkapi.go` 解压到文件卷的 `ShooterGame/Binaries/Win64`（保留已有的 `config.json`，跳过压缩包里的 `ArkApi/Plugins`），版本号写入 `.arkapi_version`。版本库：内置版本、数据目录 `arkapi/<版本>.zip`（从 GitHub `Michidu/ARK-Server-API` 下载或上传的），以及 GitHub 上还没下载的版本（安装时下载）。GitHub 上只到 3.54，3.56 来自 gameservershub。启用 ArkApi 的服务器启动时如果还没有 `version.dll`，自动装本地最新版本（不联网）。停用的插件移到不挂载进容器的 `servers/<id>/plugins-disabled/`；插件的安装、删除、启停只能在已停止时进行。镜像侧的改动需要重新构建并推送镜像后才生效。

### 服务器生命周期

1. **创建**：先检查端口冲突（`checkPorts`，含游戏端口+1），再在事务内写数据库记录、创建三个卷（服务端文件卷、数据卷、插件卷）并写入默认配置文件；容器不在创建时生成，首次启动时才创建。事务未提交时 `defer` 清理已创建的卷和数据目录（`committed` 标记；事务回滚后 ID 可能被复用，残留文件会串到下一台服务器）
2. **启动 / 停止**：用 `transitionStatus` 原子切换状态（`UPDATE ... WHERE status NOT IN (...)`，看 `RowsAffected`），避免并发请求重复启停；容器操作在后台协程中完成，结果由 `finishOperation` 写入状态和 `last_error`。启动前先把暂存的配置文件写入数据卷；服务端文件卷里还没有服务端主程序时，`seed.go` 从另一台已安装完成（有 `.install_complete` 和主程序，优先已停止）的服务器复制整个文件卷（含 Mod 和 `.installed_mods`），再删掉 `.install_complete`，入口脚本启动时 validate 只下载差异，复制期间 `runtime.phase` 为 `copying`，复制失败只记日志、照常从 Steam 下载（来源优先有目标平台主程序的；目标已有任一平台主程序时不复制）；服务器要用的 ARK 镜像不存在时先自动拉取（`images.go`，期间 `runtime.phase` 为 `pulling_image`）；开启 `auto_update` 时清除安装标记；如果 `SERVER_ARGS`、`GameModIds`、`CLUSTER_VOLUME`、重启策略（`auto_restart`）与当前配置不同，或容器的镜像 ID 不是 `ArkServerImageFor(arkapi)` 当前指向的镜像（镜像更新过或切换了 ArkApi），会先重建容器；启动超时会停掉容器。停止前先通过容器内 RCON（`docker_rcon.go`，用 `docker exec` 运行 python3）执行 `SaveWorld`，失败只记日志、不阻止停止；入口脚本会把 SIGTERM 转发给游戏进程（需要重新构建镜像才生效）。`GenerateArgs` 对 map 按 key 排序，保证相同配置生成相同参数；参数同时以数组形式放进容器 `Cmd`（入口脚本优先用 `$@`，名称和密码里的空格不会被拆开），`SERVER_ARGS` 只用于重建判断和兼容旧镜像。管理员密码不在启动参数和环境变量里，每次启动前由 `syncAdminPassword` 以数据库里的值写入 `GameUserSettings.ini` 的 `[ServerSettings]`（读取失败时不会覆盖整个文件，只有 `ErrConfigNotFound` 才视为空文件）。`GetServers` 通过 `syncStatusWithDocker` 只校正 running / stopped，不覆盖 starting / stopping；后端启动时 `ReconcileStatuses` 校正被中断的启停操作
3. **运行阶段**：容器 running 不代表能连接。`runtime.go` 每 10 秒对运行中的服务器读取最近 400 行日志判断阶段（`downloading` 含 steamcmd 进度、`installing_mods` 含 x/y、`launching`），再在容器内发 A2S_INFO，有响应即 `ready` 并拿到在线人数；同时采样 CPU/内存。结果缓存在内存，随服务器列表返回 `runtime`。阶段标志行来自 `ASE-Server-Docker/scripts/start_server.sh`，改脚本输出时要同步改 `applyPhase`
4. **重启 / 更新游戏**：`restartWith` 先停止（如果在运行）、执行可选的 `beforeStart`、再启动；更新游戏的 `beforeStart` 是删除文件卷里的 `.install_complete` / `.installed_mods`（`mod_ids` 时只删除这些 Mod 的行），入口脚本会重新 `app_update validate` 并重新下载 Mod
5. **重建**：`/servers/:id/recreate` 先原子占位为 stopping，再在后台执行：运行中的服务器先停止（失败则放弃），写入暂存配置，按当前配置创建新容器，结束时为 stopped，失败原因写入 `last_error`（详情页"管理"标签里有重建按钮）
6. **删除**：只允许删除已停止的服务器（starting / stopping 返回 409）；先删容器，再以 `status = stopped` 为条件软删除数据库记录，同时删除定时任务。默认保留数据目录和备份；`delete_data=true` 时删除三个卷、`servers/<id>/` 目录和备份目录（集群卷和 `clusters/<id>` 共享，不删）

### 发布到收录站

服务器的 `publish` 开关（表单"运维"分组，默认关闭，保存后立即生效，不用重启）。`service/server/publish.go` 的 `Publisher` 每 15 秒检查一次：服务器每次变成 running 且 `runtime.phase` 为 `ready` 时 `POST {ARKFGF_API}/publish` 发布一次（名称、端口、地图、Mod、ArkApi、集群，倍率和 PvE 先取启动参数 `?Key=` 再取配置文件），收录站返回的令牌存在 `servers.publish_token`；之后不再发请求，没有心跳，停服也不通知：收录站每 30 分钟自己查询端口，查不到就显示离线。运行中短暂不是 ready（监控的一次 A2S 超时）2 分钟内不算断开，不会重新发布。关闭发布或删除服务器（软删除，记录还在）后，有令牌的服务器调用 `/publish/unpublish` 删除收录站上的记录，成功（或收录站返回 `invalid_token`）后清空令牌，失败一直重试，管理端重启后继续。失败按 1、2、4…30 分钟退避，429 / 503 按 `Retry-After`，停服再启动不会跳过退避。状态只在内存里，随服务器响应返回 `publish_status`（`waiting` / `published` / `error` / `disabled`，失败时带 `error_code` 和中文 `error`），详情页概览显示。

收录站（外层仓库 `web/backend/controllers/publish.go`）不需要登录：收录的 IP 是请求的来源 IP（`CF-Connecting-IP`），收录站自己向该 IP 的查询端口发 A2S_INFO 验证，名称、地图、人数、是否有密码以 A2S 为准，所以查询端口（UDP）必须能从公网访问。收录站每 30 分钟查询一次所有自动发布的服务器（连续 24 小时查不到的不再查询，等下次发布），另有按 IP / 端口 / 全站的频率限制和每 IP 台数上限，30 天没上线的自动删除。服务器运行中公网 IP 变了的话，要等下次重新可以连接（重启）时才会更新。

收录资料（`service/server/listing.go`）：简介（500 字）、标签（8 个，每个 16 字）、QQ 群号、Discord 邀请链接，存在 `servers.listing_*` 列，接口里是 `listing` 对象（创建、更新、MCP 的 create/update_server 都能传，传了就整体替换），详情页"收录站"卡片里编辑。每次发布都带上这四项（空值表示清空；收录站对不带这些字段的旧版请求保持原值）。运行中修改会把发布状态重置，距上次发布至少 1 分钟后重新发布。同一张卡片显示收录站上的访问统计（`GET /api/servers/:id/publish/stats` 代理收录站的 `GET /publish/stats`，用发布令牌认证，缓存 5 分钟；没有令牌时返回 409 `not_published`）和状态徽章的引用代码（HTML / Markdown / BBCode）。

### P2P 联机房

服务器的 `link` 开关（表单"运维"分组"P2P 联机房"，默认关闭，与 `publish` 独立，保存后立即生效）。没有公网 IP 的服主也能让玩家连进来：网站 arkfgf.com 登记房间并运行只做打洞信令的公共 EasyTier 节点，Commander 为每台开了联机房的服务器运行一个 EasyTier 节点容器，玩家用 ARKFGF 联机助手（子模块 `ArkServerLink/`，单独发布到它自己的 GitHub Release，下载页 arkfgf.com/download）从网站领取一次性凭证加入这个 EasyTier 网络，直连到节点。只直连不中继（`--p2p-only`），打洞失败的玩家连不上。三方接口约定见外层仓库 `docs/p2p-link.md`，改动要和网站、客户端同步。

`service/server/link.go` 的 `LinkManager` 每 15 秒检查一次（`StartLinkManager`，`ARKFGF_API` 为空时只删除残留的节点容器）：
- 登记：没有令牌或资料变化（`buildLinkRoomProfile`，和发布到收录站共用 `profile.go` 的 `buildServerProfile`：名称、地图、人数、是否有加入密码、PvP、ArkApi、集群、Mod、倍率、收录资料，再加游戏端口、查询端口、虚拟 IP `10.126.126.1`，游戏版本取自 A2S 名称末尾的 `(vxxx)`，停服后沿用上次的值）时 `POST {ARKFGF_API}/link/rooms`（带令牌即更新）。返回的令牌、房间 ID、网络名存在 `servers.link_token` / `link_room_id` / `link_network_name`；`peers` 和房间地址只在内存里，管理端重启后第一轮重新登记拿回来。令牌对不上时网站会新建房间，此时删除节点容器和凭证目录（旧房间的凭证不能留在节点里）
- 本地参数：`link_network_secret`（32 字节随机数 hex，只在本机）和 `link_rpc_port`（节点 RPC 端口，只监听宿主机 127.0.0.1）第一次处理时生成。RPC 端口 = 20000 + 服务器 ID，和数据库里已知的端口（各服务器的游戏端口、+1、查询、RCON，其他节点的 RPC 端口）冲突时顺延；管理端在容器里，探测不了宿主机上其他程序占用的端口
- 节点容器 `ase-link-<id>`（镜像 `utils.EasyTierImage` = `easytier/easytier:v2.6.4`，入口就是 `easytier-core`）：服务器 running 且 `runtime.phase` 为 ready 过一次后创建，整个运行期间保留（A2S 偶尔超时不删容器，玩家不掉线），服务器不在运行就删除。宿主机网络、`unless-stopped`、数据目录 `servers/<id>/link/` 挂到 `/data`（凭证文件 `/data/cred.json`）。参数见 `linkNodeArgs`（和约定里服主节点的参数完全一致），标签 `com.arkservercommander.link.hash` 是镜像名 + 参数的 SHA-256，参数（名称、端口、网络名、peers 等）或镜像 ID 变化时重建。镜像不存在时后台自动下载（`ClaimPull`），下载期间节点不启动。其他 `ase-link-*` 容器（关闭了联机房、已删除、没在运行）每轮都会被删除。`ListServerContainers` 只匹配 `ase-server-<id>`，不会把节点当成服务器
- 心跳 `POST /link/rooms/:id/heartbeat`：在线时每 30 秒、离线时每 60 秒；在线状态变化、有凭证要上传、有吊销要确认、踢人之后（`poke`）提前发。在线 = 服务器 ready（运行中查询不到 2 分钟内仍算在线）且节点容器在运行；人数取运行监控的 A2S 结果。响应的 `want_credentials` 由下一轮在节点里签发（`docker exec` 运行 `easytier-cli -p 127.0.0.1:<rpc> -o json credential generate --ttl 2592000 --credential-id <UUID> --reusable false`，每轮最多 20 个，每次心跳最多上传 50 个，`expires_at` 为 30 天后），`revoke` 由下一轮在节点里 `credential revoke <id>`（输出 `{"success": false}` 表示凭证不存在，也算完成）后随心跳回报 `revoked_ack`。签发和吊销失败只记节点错误，心跳照发；没上传成功的凭证和确认留到下一次。网络名、peers 以心跳响应为准。心跳返回 `invalid_token` / 404 时清空令牌、删除节点和凭证目录，退避后重新登记新房间
- 下架：关闭联机房或删除服务器（软删除）且有令牌时 `DELETE /link/rooms/:id`，成功（或 `invalid_token` / 404）后删除节点容器和 `servers/<id>/link/`，清空令牌、房间 ID、网络名和网络密钥，失败一直重试。`DeleteServer` 也会顺手删除节点容器
- 失败按 1、2、4…30 分钟退避，429 / 503 按 `Retry-After`（`arkfgf.go` 的 `arkfgfRequest` 是和发布共用的请求函数）
- 状态只在内存里，随服务器响应返回 `link_status`（`waiting` / `online` / `offline` / `error` / `disabled`，带 `room_url`，失败时带 `error_code` 和中文 `error`；`error_code` 除了网站的错误码还有 `network`、`node_error`、`image_error`、`local_error`）

接口 `GET /api/servers/:id/link/members`、`POST /api/servers/:id/link/kick`（`{steam_id, ban}`）、`DELETE /api/servers/:id/link/bans/:steam_id` 用房间令牌代理到网站的同名接口（没有令牌 409 `link_not_registered`，连不上 502 `link_unavailable`，网站返回 400 / 404 / 409 / 422 时 400 `link_request_failed`，SteamID 不是 17 位数字 400 `invalid_steam_id`）。踢人后网站把凭证放进下一次心跳的 `revoke`，所以玩家最多约 30 秒后断开。前端：详情页概览的"联机房"一行显示状态，开启后显示 `LinkCard`（房间链接、加入方式说明、成员的踢出 / 封禁、封禁列表的解封，都有确认弹窗）。

### 存档备份

备份放在 `BACKUP_DIR/<server_id>/<时间>_<类型>.tar.gz`（类型：manual / scheduled / pre-restore / uploaded，时间按 `CONTAINER_TZ`），归档内是 `Saved/`（跳过 `Logs` 和 `clusters`）。运行中备份先 RCON `SaveWorld`。恢复只能在已停止时进行：先自动做 pre-restore 备份，占位为 stopping，直接在数据目录里清空 `SavedArks` 和 `Config` 再解压（逐项检查路径，只写普通文件和目录），清除暂存配置。上传的归档要通过 `ValidateBackup`（gzip、只允许 `Saved/` 下的普通文件和目录、必须有 `SavedArks`）。同一服务器的备份 / 恢复由 `backupLocks` 串行化

### 定时任务

`schedules` 表，`schedule.go` 每分钟整点检查（`robfig/cron` 5 段表达式或 `@daily` 等，时区 `CONTAINER_TZ`，系统休眠后最多补 5 分钟）。restart / stop / update 可设 `warn_minutes`，在 N、10、5、1 分钟各广播一次（`message` 中的 `{minutes}` 会替换）；backup 的 `keep` 只清理 scheduled 类型的旧备份。定时重启只处理运行中的服务器

### 镜像更新

1. **检查更新**：`CheckImageUpdate` 用 `DistributionInspect` 取镜像仓库中该标签的 manifest digest，与本地镜像的 `RepoDigests` 比较，不下载镜像层。结果记录在内存中，由 `/images/status` 的 `has_update` 返回，拉取成功后清除。本地构建、没有 `RepoDigests` 的镜像会一直显示有更新
2. **更新**：先用 `ClaimPull` 原子占用拉取权（并发请求返回 409），再后台拉取新镜像，失败原因记录在镜像状态的 `error` 字段并显示在首页；拉取成功且镜像 ID 变化时会清理被替换的旧镜像（不强删，仍被容器引用时只记日志），拉取进度流中的 `error` 会被当作失败。已停止的服务器在下次启动时自动重建；`restart_running` 为 true 时，运行中且镜像已过期的服务器会依次经过 stopping → stopped → starting，启动时重建容器
3. 前端：首页"镜像管理"检查更新，点击"更新"先请求 `/images/affected`，再弹出 `ImageUpdateConfirmModal` 选择是否重启运行中的服务器；服务器卡片根据 `image_outdated` 显示"旧镜像"标记。`/images/affected` 和 `restart_running` 只涉及使用该镜像的服务器。`/images/status` 里每个镜像带 `required`：ARK 镜像只有在服务器用到时才算（没有服务器时算 Linux 版），EasyTier 镜像只有在有服务器开启联机房时才算（更新它之后，运行中的联机房节点在下一轮按新镜像重建），`any_not_ready` 只看需要的镜像。创建和启动服务器不依赖任何镜像就绪（ARK 镜像启动时自动拉取），所以前端不再根据镜像状态禁用启动 / 创建按钮。`alpine:latest` 已不再使用，也不在管理的镜像里

### 配置文件

GameUserSettings.ini 和 Game.ini 不存数据库（暂存内容除外），而是直接保存在数据卷的 `Config/WindowsServer/` 下：`docker_config.go` 直接读写后端本地的 `servers/<id>/saved/Config/WindowsServer/`，写入先写临时文件再改名。服务器未停止时保存的配置先存入 `pending_game_user_settings` / `pending_game_ini`（ARK 停服时会覆盖 GameUserSettings.ini），下次启动或重建前写入数据卷，响应里 `config_pending` 为 true。前端每次修改一个设置项，就用 `lib/ini.ts` 的 `setIniValue` 在原文本的对应 section 里改这一行，未知键、注释和自定义 section 原样保留。

### 配置导入导出

`models.ServerConfig`（`format` 固定为 `ark-server-commander/server-config`，`version` 为 1）包含地图、人数、集群、Mod、ArkApi / 自动重启 / 自动更新、两个配置文件和启动参数，不含端口和管理员密码，导出时去掉 GameUserSettings.ini 里的 `ServerAdminPassword`。导入和复制按 `parts` 选择部分：`config_files`、`server_args`、`mods`、`map`、`cluster`、`options`（为空表示全部），最终都经过 `UpdateServer`，未停止的服务器配置文件先暂存。写入配置文件时用 `utils.KeepIniKeys` 保留目标服务器自己的 `SessionName`、`ServerAdminPassword`、`Port`、`QueryPort`、`RCONPort`、`MaxPlayers`（同时导入了 `options` 时人数改成新值），前端 `lib/server-config.ts` 的 `applyConfigToForm` 用 `lib/ini.ts` 的 `keepIniValues` 实现同样的规则，两边要保持一致。前端：详情页"管理"标签里导出（下载 JSON）、导入（`ImportConfigDialog` 的 `apply` 模式，调用 `/config/import` 直接保存）和复制到其他服务器（`CopyConfigDialog`，直接保存）；新建和编辑页底部"导入配置"（`ImportConfigDialog` 的 `form` 模式，只填进表单，保存后才生效）。MCP 工具 `export_server_config`、`import_server_config`、`copy_server_config`。

### 服务器设置界面

新建和编辑页（`ServerForm`）不按文件分区，而是按用途分组（常规、网络、Mod 与集群、游戏规则、倍率、玩家、恐龙、繁殖、建筑、部落与 PvP、运维、高级），左侧导航带修改计数，顶部可以跨分组搜索（名称、键名、说明）。所有设置项定义在 `lib/server-settings.ts` 的 `SETTING_GROUPS`，`source` 记录存放位置：`field`（服务器记录字段，由表单单独渲染）、`gus`（GameUserSettings.ini）、`game`（Game.ini）、`query`（`?Key=Value`）、`arg`（`-Key`）。`readSetting` / `writeSetting` 统一读写，界面每次渲染都从文本重新解析，所以"高级 → 配置文件"里直接改文本也会同步到设置项。`default` 是没有设置时游戏的表现，和它不同的值显示为已修改，恢复默认就是删掉这个键或参数（开关参数关闭也是删除）。新增设置项：在字典里加一行，在 `messages/*.ts` 的 `serverSettings.params`（必要时 `hints`）加中英文名称。新建页先请求 `/system/server-defaults`，以后端的默认模板和推荐启动参数为基础编辑，创建时再把模板里的名称、人数和欢迎消息换成表单里的值（`new/page.tsx` 的 `applyIdentity`）。

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
- 详情页用列表接口 `/api/servers` 轮询，不要轮询 `GET /servers/:id`（它会读取配置文件和 Docker 状态，开销较大）
- 翻译用 ICU 格式：文案里不能直接写 `<xxx>`（会被当成标签）或不传参数的 `{xxx}`；新增键时中英文都要加
- 错误提示用 `useApiError()` + `toast`，不要在 store 里写死中文
- `Button` 会给没有 `size-*` 类的图标强制设为 `size-4`，按钮里的图标要改大小时用 `size-3.5` 这类写法，`h-3 w-3` 不生效
- 按钮禁用时浏览器不显示 `title`，需要说明原因时把 `title` 放在外层 `span` 上（见 `ServerManagePanel`、`IconAction`）

### 发布

推送 `v*` tag 触发 `docker-build.yml`：构建推送镜像后创建 GitHub Release，内容取附注 tag 的正文（第一行概括之后的部分），再接上镜像信息。更新说明由 Claude 根据上一个 tag 以来的提交生成，分 English、中文两栏（English 在前），流程见 `.claude/skills/release/SKILL.md`（打 tag 必须用 `--cleanup=whitespace`，否则 `##` 标题会被当注释删掉）。仓库直接提交到 main，GitHub 自动生成的 "What's Changed" 是空的。

### 提交前检查

```bash
(cd server && gofmt -l . && go build ./... && go vet ./... && go test ./...)
(cd ui && npm run lint && npx tsc --noEmit && npm run build)
```

## 测试

测试覆盖率仍然较低：`server/service/docker_manager/rollback_test.go`、`docker_image_test.go`（`ClaimPull` 并发占用、拉取进度消息解析）、`docker_backup_test.go`（备份格式校验、CPU 计算）、`localfs_test.go`（配置文件读写、备份导出恢复、检查和复制服务端文件）、`server/config/config_test.go`（JWT 密钥生成和复用）、`server/models/server_args_test.go`（启动参数顺序确定）、`server/service/server/server_service_test.go`（端口冲突、保存不覆盖状态、端口建议、禁用的定时任务，使用临时 SQLite）、`parse_test.go`（运行阶段、ListPlayers 解析、字段规范化、定时任务校验）、`server/controllers/auth/ratelimit_test.go`、`server/utils/ini_test.go`（`SetIniValue`、`KeepIniKeys`）、`config_files_test.go`（默认模板的 section）、`installed_mods_test.go`、`plugins_test.go`（zip 识别、覆盖安装、启停删除）、`arkapi_test.go`（版本比较、压缩包根目录识别、安装保留配置、自动安装、版本库上传删除）、`gamelog_test.go`（UTF-8 / UTF-16 尾部读取）、`version_test.go`（版本比较）、`server/service/apitoken/apitoken_test.go`（令牌生成、校验、撤销）、`server/mcpserver/mcpserver_test.go`（用 SDK 客户端走真实 HTTP：令牌认证、只读令牌的工具过滤、错误码）、`server/service/server/publish_test.go`（每次可以连接只发布一次、宽限期、停服不发请求、关闭发布和删除后取消发布、退避和 Retry-After，用 httptest 模拟收录站）、`server/service/server/server_config_test.go`（配置导出去掉密码、复制时保留目标的名称端口密码、按部分应用、运行中暂存）、`server/service/server/link_test.go`（联机房：登记和令牌保存、资料变化重新登记、在线 / 离线心跳间隔、节点参数、参数或镜像变化时重建、停服删除节点、凭证补充数量和每轮上限、吊销和确认、节点出错、关闭和删除后下架、退避和 Retry-After、令牌失效重新登记、成员接口代理，用 httptest 模拟网站、假的 Docker 操作）、`link_docker_test.go`（默认跳过，`ASC_LINK_DOCKER_TEST=1` 时用真实 Docker 按生成的参数启动 EasyTier 节点，签发和吊销凭证后删除）。启停流程、RCON、配置暂存和前端流程都没有自动化测试。Windows 上 `config_test.go` 的 0600 权限检查会失败（Windows 没有这种权限位），以 Linux 结果为准。本机没有 Go 时，可以用 `docker run --rm -v $PWD/server:/src -w /src golang:1.27-alpine go test ./...` 在容器里跑。

## 已知问题

- Linux 版镜像（原生服务端、`Config/LinuxServer` 符号链接、`steamclient.so` 补链接）以及在同一个文件卷上切换平台后 steamcmd 的行为，都没有在真实服务器上验证过；ArkApi 压缩包的目录结构按"找到 `version.dll` 所在目录"处理，只用内置 3.56 和测试数据验证过。镜像需要推送新 tag 构建出 `linux` / `proton` 两个标签后，新版管理端才能启动服务器
- 加入 `arkapi` 列之前创建的服务器在迁移时全部设为启用 ArkApi（保持 Proton 版），新服务器默认关闭
- 登录限流按 IP（`X-Forwarded-For`，由 Next 转发）计数，客户端能伪造该头绕过 IP 限制，用户名维度的上限仍然有效；限流记录在内存，重启后清零
- 修改登录密码后，旧 JWT 在 24 小时有效期内仍然可用
- `npm audit` 有来自 ESLint 依赖链（`minimatch`）的中等级别告警，仅影响开发期
- 入口脚本已转发 SIGTERM，但游戏进程能否在 `docker stop` 的 30 秒内完成存档，没有在真实环境验证过，停服前的 `SaveWorld` 继续保留。steamcmd 下载阶段的 trap 已加（用 `setsid` 单独进程组），只用模拟进程测过，需要重新构建镜像后才生效
- 集群卷和 `data/clusters/<id>` 不会随服务器删除；没有服务器使用后会出现在首页磁盘明细的"未使用的集群数据"里，可以在那里删除。已有服务器改成集群后，原数据卷里 `Saved/clusters` 的旧数据会被集群卷遮住
- 服务器数据在宿主机 `./data/servers/<id>/` 下，文件属主是 root，外部编辑需要 sudo；运行中修改 `GameUserSettings.ini` 会在停服时被 ARK 覆盖，要在停止后再改；`data` 所在磁盘要容纳每台服务器约 20 GB。早期版本创建的普通 Docker 卷不会被识别，也没有做迁移（工具未上线），需要删除旧服务器后重建
- 运行阶段依赖入口脚本的日志文本；`ShooterGame.log` 的编码（UTF-8 / UTF-16LE）是自动识别的，没有用真实服务器的日志验证
- 创意工坊名称查询依赖访问 `api.steampowered.com`，访问不了时只显示 ID
- 管理员密码现在存放在卷里的 `GameUserSettings.ini`，`docker inspect` 看不到；数据库里仍是明文，并且每次启动都以数据库的值为准，在 ini 里直接改 `ServerAdminPassword` 会被覆盖
- 每台服务器各有一份完整的服务端文件卷，磁盘占用大（新服务器从已安装的服务器复制，省的是下载，不省磁盘）；Mod 平时只下载尚未安装的，创意工坊有更新时详情页会提示，需要手动更新所选 Mod、"更新游戏/Mod"或开启 `auto_update` 才会重新下载
- 默认配置模板按 `docs/ASE_Game.ini.md` / `docs/ASE_GameUserSetting.md` 分配：倍率、难度、昼夜速度在 GameUserSettings.ini 的 `[ServerSettings]`，Game.ini 只放自己的选项，都在 `[/script/shootergame.shootergamemode]`。旧服务器的 Game.ini 里可能还留着这些无效的键，设置界面只读写所属 section 里的键，不显示它们。未在真实服务器验证
- 创意工坊搜索接口需要 Steam Web API Key，没有做站内搜索，只支持粘贴链接
- 版本检查依赖访问 `api.github.com`（未认证每小时 60 次，结果缓存 6 小时）
- MCP 令牌长期有效、权限是整个账户（完全控制或只读两档），没有按服务器限制；MCP 和网页共用对外端口，公网暴露时需要 HTTPS 反向代理，否则令牌明文传输
- 发布到收录站时，收录站看到的来源 IP 必须就是玩家连接的 IP：出口和入口 IP 不同（多 IP、NAT、代理）的机器会因 A2S 验证失败发布不了；只支持 IPv4
- P2P 联机房没有在真实服务器上和网站、客户端联调过（只有单元测试和本机 Docker Desktop 上的节点启动、签发 / 吊销凭证）。玩家流量经节点投递到宿主机 `127.0.0.1:<游戏端口>`，再靠 Docker 的端口映射（userland-proxy 或 iptables 对回环地址的 DNAT，需要 `route_localnet` 或 userland-proxy 开启，Docker 默认开启 userland-proxy）到达 ARK 容器，`userland-proxy: false` 的宿主机上可能不通；游戏端口只映射到某个具体 IP 而不是 0.0.0.0 时也不通。节点 RPC 端口 = 20000 + 服务器 ID（冲突时顺延），只避开数据库里已知的端口，被宿主机上其他程序占用时节点起不来（状态显示 `node_error`）。所有玩家的流量在 ARK 看来都来自同一个地址（本机回环 / Docker 网关），按 IP 的封禁和统计不可用。凭证在节点的 `cred.json` 和网站各有一份，管理端重启时内存里还没上传的凭证会丢失（只是浪费，30 天后过期）
- `GetServers` 现在一次 `ContainerList` 取状态，依赖容器名匹配 `ase-server-<id>`；`GetImageStatus` 的 `overall_status` 是状态码（`ready` / `pulling` / `not_ready`）

## 故障排除

- **后端无法启动**：`docker ps` 确认 Docker 可用；检查数据库目录是否可写（需要写入 `jwt_secret`）；如果手动设置了 `JWT_SECRET`，确认不少于 32 字符且不含弱口令片段
- **所有用户被登出**：`jwt_secret` 文件被删除或数据目录没有持久化，重启后生成了新密钥
- **登录后立即回到登录页**：通过 HTTP 访问却设置了 `COOKIE_SECURE=true`
- **页面提示后端服务不可用（502）**：检查 `API_BASE` 是否指向后端，以及后端进程是否在运行
- **前端构建失败**：删除 `ui/.next` 和 `ui/node_modules` 后执行 `npm ci`
- **Docker 操作失败**：确认挂载了 `/var/run/docker.sock` 并以 `privileged` 运行
- **服务器显示"运行中"但连不上**：看卡片上的阶段（下载服务端 / 安装 Mod / 加载中），详情页"日志"标签能看到实时进度；首次启动要下载约 20 GB
- **连接地址不对**：设置 `PUBLIC_HOST`
- **玩家版本不匹配进不了服**：服务器详情页"管理"标签里的"更新游戏/Mod"
