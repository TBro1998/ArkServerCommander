# 待实现功能开发计划

来源：2026-11-12 对照 README "已实现 / 待实现功能"和代码逐项核对的结果。状态：`[ ]` 未开始，`[~]` 进行中，`[x]` 完成。
每项完成后在"备注"里写明实现位置和需要注意的地方，并同步更新 README（中英文）的功能列表。

## 核对结论

README "已实现功能"的 12 项都有对应实现，后端 `go build/vet/test`、前端 `tsc`、`lint` 通过。以下描述已过时，已在 README 中改正：

- "容器重建后会重新下载"：服务端文件已改为放在独立的文件卷里，重建容器不会重新下载；新服务器从已安装的服务器复制
- "删除服务器时保留卷中的数据"：默认仍保留，但删除时可以勾选同时删除数据和备份；卷绑定到 `data/servers/<id>/`
- 常见问题里"备份未实现""日志看不到"已不成立

README "待实现功能"中以下几项实际已完成（见 [ROADMAP.md](ROADMAP.md)），已移到"已实现"：

| README 条目 | 实现 |
|---|---|
| RCON 命令执行 | 详情页"控制台"，ROADMAP 3.1、3.2 |
| 服务器运行状态监控 | 运行阶段、在线人数、CPU / 内存 / 磁盘，ROADMAP 1.3、1.4、3.4、5.5 |
| 服务器日志查看 | 详情页"日志"，ROADMAP 1.1、1.2（只有容器输出，见下文 6.1） |
| 服务器存档及配置备份 | 详情页"备份"，ROADMAP 2.3 |
| 可选更新服务端文件和 Mod | "更新游戏/Mod"和 `auto_update`，ROADMAP 2.2 |
| 用户管理中的"修改密码" | 设置页，ROADMAP 5.6（多用户未做，见下文 3.x） |

## 一、Mod 管理对接 Steam 创意工坊

2026-11-13 完成。镜像侧改动（1.1）需要重新构建并推送 `tbro98/ase-server` 后才生效，旧镜像装的 Mod 用 `.mod` 文件的修改时间代替安装时间。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 1.1 | 记录 Mod 安装时间 | [x] | `start_server.sh` 写入 `id 时间戳`，判断已安装改为 `grep -E "^id( \|$)"`，兼容只有 ID 的旧行。后端 `docker_manager/installed_mods.go` 直接读数据目录 `servers/<id>/files/.installed_mods`，旧行用 `ShooterGame/Content/Mods/<id>.mod` 的修改时间 |
| 1.2 | 检测 Mod 更新 | [x] | `service/server/mods.go`：创意工坊 `time_updated` 晚于安装时间即有更新。后台每 30 分钟检查（启动 1 分钟后首次），结果在内存，随服务器列表返回 `mods_outdated`；`GET /servers/:id/mods?refresh=true` 忽略 6 小时缓存。卡片和详情页显示"N 个 Mod 有更新" |
| 1.3 | 只更新指定 Mod | [x] | `update-game` 增加 `mod_ids`（必须在服务器的 Mod 列表里），只从 `.installed_mods` 删除这些行。`mods: true` 时忽略 `mod_ids`。详情页 Mod 卡片默认勾选有更新的 Mod |
| 1.4 | 预览图、更新时间、链接、调整顺序 | [x] | `ModList.tsx`：表单里可以拖动、上下移动、移除；`ServerModsPanel.tsx`：详情页的安装状态和更新。预览图用 `<img>`（Steam CDN 不在 next/image 白名单里） |
| 1.5 | 创意工坊搜索 | [x] 部分 | 实测 `QueryFiles` 不带 Key 返回 403，按计划只做了：粘贴创意工坊链接自动换成 ID（失焦时规范化，重复的去掉），表单里加"在创意工坊浏览 Mod"链接。站内搜索需要用户配置 Steam Web API Key，暂不做 |

## 二、ArkApi 插件管理

2026-11-13 完成。2.6 的版本文件需要重新构建镜像后才有，旧镜像显示"版本未知"。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 2.1 | 列出插件 | [x] | `service/server/plugins.go`，`GET /servers/:id/plugins`：目录名、`PluginInfo.json`（去 BOM，版本可以是数字或字符串）、大小、是否有同名 dll 和 config.json |
| 2.2 | 上传安装、删除 | [x] | `POST /servers/:id/plugins`（multipart `file`，可选 `name`，上限 200 MB，解压后 500 MB、5000 个文件）。每个 `<名称>/<名称>.dll` 所在目录为一个插件，可以带前缀、一次多个；只有一个 dll 时用它的名称。拒绝 `..`、绝对路径、反斜杠、符号链接，按声明大小限制写入量。先解压到插件卷外的临时目录再改名；覆盖安装保留原 config.json 和启用状态。安装、删除、启停只能在已停止时进行（运行中替换 dll 会让游戏崩溃），同一服务器串行 |
| 2.3 | 读写 config.json | [x] | `GET/PUT /servers/:id/plugins/:name/config`、`/servers/:id/arkapi/config`（`files/ShooterGame/Binaries/Win64/config.json`，入口脚本不会覆盖）。保存前校验 JSON，先写临时文件再改名。运行中也能保存 |
| 2.4 | 详情页"插件"标签 | [x] | `PluginsPanel.tsx`：上传、启停开关、删除确认、JSON 编辑弹窗（格式化、校验）。运行中提示可以用 `plugins.unload` / `plugins.load` 重新加载，没有自动执行（未在真实服务器验证） |
| 2.5 | 启用 / 停用 | [x] | 没有用 `plugins/.disabled/`：ArkApi 会遍历 Plugins 下的每个目录，无法确认它会跳过点开头的目录。改为移到不挂载进容器的 `servers/<id>/plugins-disabled/`，删除服务器数据时一起删除 |
| 2.6 | 显示 ArkApi 版本 | [x] | 镜像 `dockerfile` 增加 `ENV ARKAPI_VERSION=3.56`（`COPY ArkApi_${ARKAPI_VERSION}`），入口脚本写入文件卷的 `.arkapi_version`；`version.dll` 存在即视为已安装 |

## 三、用户管理

现状：只有初始化时创建的一个管理员；`servers`、`schedules` 已按 `user_id` 归属，镜像和系统接口是全局的。已支持修改密码和登录限流。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 3.1 | 修改密码后旧令牌立即失效 | [ ] | `User` 增加 `token_version`，写进 JWT，`middleware/auth.go` 校验；修改密码时加一。顺带实现"退出所有设备" |
| 3.2 | 角色：`admin` / `user`，初始化创建的是 admin | [ ] | 镜像管理、系统概览、用户管理只允许 admin |
| 3.3 | 用户管理接口和设置页：新增、删除、重置密码、改角色 | [ ] | 不能删除最后一个 admin；删除用户时其服务器的处理方式要先定（转给 admin 或禁止删除） |
| 3.4 | 服务器共享 / 授权（可选） | [ ] | 先确认是否需要；最简单是 admin 可见全部，普通用户只看自己的。注意 `checkPorts` 已经按全部服务器检查，多用户后不要改成只查自己的 |

## 四、工具版本更新检查

2026-11-13 完成。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 4.1 | 构建时注入版本号 | [x] | `utils.Version`（默认 `dev`），`Dockerfile` 的 `ARG VERSION` 传给 `-ldflags -X`，CI 用 `docker/metadata-action` 的 `version`（tag 去掉 v）。前端不单独注入，从接口取 |
| 4.2 | `GET /system/version` | [x] | `service/server/version.go`：取 GitHub `releases/latest`（不含预发布），成功缓存 6 小时、失败 10 分钟，`?refresh=true` 立即查询；按 semver 比较，`dev` 不提示更新 |
| 4.3 | 设置页显示版本和升级命令 | [x] | `components/settings/VersionCard.tsx`：当前版本、最新版本和发布说明链接，有更新时给出可复制的 `docker compose pull && docker compose up -d`。后端启动日志也打印版本 |

## 五、MCP 支持

让 AI 助手通过 MCP 管理服务器。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 5.1 | API 令牌：设置页生成 / 吊销长期令牌，只保存哈希 | [ ] | MCP 和脚本调用都需要，不能复用 24 小时的登录 JWT |
| 5.2 | MCP 服务端（Streamable HTTP），挂在后端 `/mcp`，经 Next 代理对外 | [ ] | Go 用官方 `modelcontextprotocol/go-sdk`；代理需要支持流式（通配代理已支持） |
| 5.3 | 工具：列出服务器和状态、启动 / 停止 / 重启、执行 RCON、在线玩家、创建备份、读取最近日志 | [ ] | 直接调用 `service/server`，权限与令牌所属用户一致；危险操作（删除、恢复）先不提供 |
| 5.4 | 文档：README 增加 MCP 客户端配置示例 | [ ] | |

## 六、遗留改进

2026-11-13 完成。6.3 需要重新构建镜像后才生效。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 6.1 | 查看游戏日志 | [x] | `GET /servers/:id/logs/game?tail=`（默认 1000，最多 5000）直接读 `servers/<id>/saved/Logs/ShooterGame.log` 末尾 2 MB，按 BOM 和零字节分布识别 UTF-8 / UTF-16LE。日志标签切换"容器输出 / 游戏日志"，游戏日志每 5 秒刷新，停服后也能看。真实 ARK 日志的编码没有验证过，两种都做了测试 |
| 6.2 | 隐藏下载进度行 | [x] | `LogViewer.tsx` 的"隐藏下载进度"勾选框，只对容器输出生效 |
| 6.3 | steamcmd 下载阶段转发 SIGTERM | [x] | `start_server.sh` 开头设置 trap，steamcmd 和 Mod 安装脚本用 `run` 以 `setsid` 放到单独进程组后台运行并 `wait`，收到信号时结束整个进程组并以 143 退出。用模拟的 steamcmd（前台 sleep）测过能立即退出、不留子进程；没有用真实镜像测 |
| 6.4 | 集群清理 | [x] | 首页磁盘明细列出"未使用的集群数据"（有 `ase-cluster-*` 卷或 `data/clusters/<id>` 目录、但没有任何服务器设置该集群 ID），可删除（`DELETE /system/clusters/:clusterId`）。仍有服务器使用，或卷还被旧容器挂载时返回 409 |
| 6.5 | 默认 Game.ini 模板中倍率 key 的位置 | [x] | 按 `docs/ASE_Game.ini.md` 和 `docs/ASE_GameUserSetting.md` 逐个核对：倍率、难度、昼夜速度，以及 `AllowFlyerCarryPvE`、`DisableStructureDecayPvE`、`EnablePvPGamma`（这三个不带 b 前缀）都属于 `[ServerSettings]`，已移到 GameUserSettings.ini 模板；Game.ini 只留它自己的选项，统一放在 `[/script/shootergame.shootergamemode]`。两个可视化编辑器同步调整，Game.ini 编辑器的空文件不再写入全部默认值。注意：模板里原本无效的 `OverrideOfficialDifficulty=5.0`、`AllowFlyerCarryPvE=True` 现在会生效（只影响新建的服务器）。没有用真实服务器验证 |
| 6.6 | 更新 `docs/ISSUES.md` | [x] | 改为复核表：9 个问题中 6 个已解决，2 个部分解决（令牌失效、审计日志），1 个不处理（多实例） |

## 长期计划（暂不排期）

- 多主机管理（可能基于 K8S）
- 服务器收录网站
- 玩家使用界面
