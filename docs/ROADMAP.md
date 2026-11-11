# 使用体验优化计划表

来源：2026-11-11 从使用者角度对前后端的梳理。状态：`[ ]` 未开始，`[~]` 进行中，`[x]` 完成。
每项完成后在"备注"里写明实现位置和需要注意的地方。

## 一、核心可用性

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 1.1 | 后端：容器日志接口（最近 N 行 + SSE 实时跟随） | [x] | `GET /servers/:id/logs?tail=`、`GET /servers/:id/logs/stream`（SSE，容器停止时发 `end` 事件）。`docker_logs.go` 按 `\r`/`\n` 分行并去掉 ANSI 控制符。Next 通配代理改为流式转发 |
| 1.2 | 前端：日志弹窗（实时跟随、暂停滚动、下载） | [x] | 详情页"日志"标签（`LogViewer.tsx`）：实时、暂停、筛选、上翻时不自动滚动、下载最近 5000 行、断开后手动重连 |
| 1.3 | 后端：细分运行阶段，解析日志 + 查询端口探测 | [x] | `service/server/runtime.go` 每 10 秒采集：从最近 400 行日志判断 downloading（含 steamcmd 进度）/ installing_mods（含 x/y）/ launching，A2S_INFO 有响应即 ready（同时拿到在线人数）。A2S 和 RCON 一样通过 `docker exec` 在容器内执行 |
| 1.4 | 前端：卡片显示运行阶段 | [x] | `StatusBadge.tsx`：下载 x% / 安装 Mod x/y / 加载中 / 可连接；卡片下载阶段显示进度条 |
| 2.1 | 集群：按 cluster_id 共享集群卷 | [x] | 卷 `ase-cluster-<cluster_id>` 挂到 ARK 默认集群目录 `Saved/clusters`（不用 `-ClusterDirOverride`，避免 Proton 下的路径问题）。容器环境变量 `CLUSTER_VOLUME` 用于判断是否重建；cluster_id 格式校验 |
| 2.2 | 游戏本体 / Mod 更新 | [x] | `POST /servers/:id/update-game {game, mods}`：停止 → 删除文件卷里的 `.install_complete` / `.installed_mods` → 启动，入口脚本会重新 `app_update validate` 和下载 Mod。服务器开关 `auto_update`：每次启动前都清除标记。不需要重新构建镜像 |
| 2.3 | 存档备份：手动备份、列表、下载、恢复、删除 | [x] | 存放在 `BACKUP_DIR/<server_id>/`，tar.gz 内为 `Saved/`（跳过 Logs 和 clusters）。运行中先 SaveWorld。恢复要求已停止，先自动做 pre-restore 备份，再清空 SavedArks 和 Config 后解压。支持上传（校验路径和文件类型） |

## 二、运维功能

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 3.1 | RCON 控制台：任意命令、广播、存档 | [x] | `POST /servers/:id/rcon/command`；详情页"控制台"：命令历史（上下键）、常用命令、广播、保存世界。连接被拒时返回 `rcon_unavailable` 友好提示 |
| 3.2 | 在线玩家列表：踢出、封禁 | [x] | `GET /servers/:id/players`（解析 ListPlayers，名称里带逗号也能正确拆分）；踢出 / 封禁带确认 |
| 3.3 | 定时任务 | [x] | 表 `schedules`，`service/server/schedule.go` 每分钟检查（robfig/cron 解析，时区 CONTAINER_TZ）。类型：restart / backup（keep 保留份数）/ update / command / start / stop；重启、停止、更新可设提前 N 分钟广播（在 N、10、5、1 分钟各广播一次）。支持立即执行，记录上次结果和下次时间 |
| 3.4 | 资源监控 | [x] | 监控协程采样 `ContainerStatsOneShot`，与上一次采样相减得到 CPU%；内存扣除 inactive_file。磁盘：`DiskUsage` 取卷大小（慢，缓存 5 分钟，可强制刷新）+ 备份目录大小 |
| 3.5 | 连接信息 | [x] | `PUBLIC_HOST` 或自动探测（ipify / ifconfig.me，缓存 1 小时）。卡片复制 `IP:游戏端口`，详情页另有 `steam://connect/IP:查询端口` |

## 三、状态刷新与反馈

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 4.1 | 服务器列表常规轮询 | [x] | `hooks/use-polling.ts`：页面隐藏时暂停。有启停中或未就绪的服务器 3 秒，否则 10 秒 |
| 4.2 | 首页镜像拉取：刷新后恢复轮询，失败提示 | [x] | 以 `any_pulling` 驱动轮询，刷新页面后自动继续；下载、更新失败都弹出提示 |
| 4.3 | 全局 toast 通知 | [x] | `lib/toast.ts` + `components/ui/toaster.tsx`（挂在根布局），store 外也能调用 |
| 4.4 | 多语言 | [x] | 后端错误响应带 `code`，`useApiError()`：中文界面显示后端详细信息，其他语言按 `apiErrors.<code>` 翻译。布局、状态徽章、store 里写死的中文已移除。补齐了旧组件缺失的翻译键 |

## 四、表单与细节

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 5.1 | 新建：自动填充可用端口、随机生成密码、字段说明 | [x] | `GET /system/suggest-ports`；新建、编辑共用 `ServerForm.tsx`，前端校验端口范围和重复、集群ID、Mod ID |
| 5.2 | 编辑：密码留空提示、未保存离开提醒、保存并重启 | [x] | 有修改时 `beforeunload` 提醒、取消时确认；运行中显示"保存并重启"。自动重启开关可编辑，启动时比对重启策略决定是否重建容器 |
| 5.3 | 删除：弹框关闭、过渡态禁用、可选删除数据 | [x] | `DELETE /servers/:id?delete_data=true` 删除三个卷和备份（集群卷共享，不删）；删除数据需输入服务器名确认。同时删除该服务器的定时任务 |
| 5.4 | Mod：格式校验、显示名称 | [x] | `GET /mods?ids=` 查询 Steam 公开接口（成功缓存 6 小时，失败缓存 5 分钟），无法访问 Steam 时显示"查不到名称"而不是"不存在"。保存时规范化（去空格、支持中文逗号） |
| 5.5 | 首页改为概览面板 | [x] | `GET /system/overview`：服务器 / 运行 / 可连接数、在线人数、CPU、内存、磁盘、最近错误，镜像管理和磁盘明细 |
| 5.6 | 账号：修改登录密码、登录限流、初始化竞态 | [x] | `PUT /auth/password`；登录失败按 IP 10 次 / 用户名 30 次（15 分钟窗口）限流，返回 429；Next 转发 `X-Forwarded-For`。`InitUser` 加互斥锁 |
| 5.7 | 移动端适配 | [x] | 顶部导航在小屏收进侧边菜单；表单、详情页、列表按断点换列 |

## 收尾

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 9.1 | 后端测试补充、Swagger 重新生成 | [x] | 新增：阶段解析、玩家解析、字段规范化、定时任务校验、备份格式校验、CPU 计算、登录限流、禁用状态的定时任务、端口建议、服务器归属 |
| 9.2 | 前后端提交前检查全部通过 | [x] | Linux 容器内 `go build/vet/test` 通过；`npm run lint`、`tsc`、`npm run build` 通过。另在本机用独立数据库和真实 ARK 容器端到端验证 |
| 9.3 | 更新 CLAUDE.md / README | [x] | CLAUDE.md、`.env.example` 已更新 |

## 实现中发现并修复的问题

- GORM 创建记录时会把 bool 零值替换成 `default:true`：`auto_restart=false` 新建服务器、`enabled=false` 新建定时任务都会被存成 true。已在创建后显式写回，并加了测试。
- `proxy.ts` 缺少语言 cookie 时会跳过登录检查（CLAUDE.md 已知问题），已修复；新增的 `/settings` 也受保护。
- 旧的编辑页、新建页引用了不存在的翻译键（`serverUpdateError` 等），界面会显示键名，已补齐。

## 后续可以考虑

- 镜像侧：steamcmd 下载阶段入口脚本还没设置 SIGTERM 转发，停止下载中的服务器要等 `docker stop` 的 30 秒超时（实测）。需要在 `start_server.sh` 开头就设置 trap，并重新构建镜像。
- 修改登录密码后旧 JWT 仍在 24 小时内有效，需要令牌版本号才能立即失效。
- 日志查看只覆盖容器输出；游戏本身的 `ShooterGame.log` 在数据卷的 `Saved/Logs` 下，可以再加一个读取接口。
- 下载服务端时日志被 steamcmd 进度行刷屏，可以加"隐藏进度行"开关。
