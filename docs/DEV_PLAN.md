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

现状：Mod ID 格式校验、规范化，`GET /mods?ids=` 查询名称、更新时间和预览图（`service/server/system.go`）。没有搜索，也不知道已装 Mod 是否过期。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 1.1 | 记录 Mod 安装时间：入口脚本把 `.installed_mods` 改为 `id 时间戳` 格式（兼容旧格式），后端读取 | [ ] | 镜像侧改动，需要重新构建镜像；后端用 `RunHelper` 或数据目录直接读文件卷 |
| 1.2 | 检测 Mod 更新：对比创意工坊 `time_updated` 和安装时间，详情页 / 卡片标出"有 Mod 更新" | [ ] | 依赖 1.1；只检查运行中或最近启动过的服务器，复用 `/mods` 的缓存 |
| 1.3 | 只更新指定 Mod：`update-game` 增加 `mod_ids`，只从 `.installed_mods` 删除这些 ID | [ ] | 现在 `mods: true` 会重新下载全部 Mod |
| 1.4 | Mod 列表展示预览图、更新时间、链接，支持拖动调整加载顺序 | [ ] | 顺序即 `GameModIds` 顺序，改顺序会触发容器重建 |
| 1.5 | 创意工坊搜索并添加（按名称 / ID） | [ ] | Steam 搜索接口 `IPublishedFileService/QueryFiles` 需要 Web API Key，先确认能否不用 Key；不行则只支持粘贴链接自动解析 ID |

## 二、ArkApi 插件管理

现状：插件目录是独立的卷，绑定到宿主机 `data/servers/<id>/plugins`，只能在宿主机上手动管理。ArkApi 本体随镜像同步。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 2.1 | 后端：列出插件（目录名、`PluginInfo.json` 里的版本和说明、大小） | [ ] | 直接读数据目录，不需要临时容器 |
| 2.2 | 后端：上传 zip 安装插件、删除插件 | [ ] | 校验 zip 路径（禁止 `..`、绝对路径、符号链接），限制大小；只允许在停止时操作，或提示需要重启 |
| 2.3 | 后端：读写插件 `config.json`，以及 ArkApi 本体的 `config.json` | [ ] | 保存前校验 JSON；ArkApi 的 `config.json` 在文件卷 `ShooterGame/Binaries/Win64/config.json` |
| 2.4 | 前端：详情页新增"插件"标签（列表、上传、删除、JSON 编辑） | [ ] | 中英文翻译都要加 |
| 2.5 | 启用 / 禁用插件 | [ ] | ArkApi 没有禁用开关，可以把目录移到 `plugins/.disabled/`，需确认 ArkApi 不会加载点开头的目录 |
| 2.6 | 显示 ArkApi 版本 | [ ] | 从镜像或 `version.dll` 读取，镜像侧可能要写一个版本文件 |

## 三、用户管理

现状：只有初始化时创建的一个管理员；`servers`、`schedules` 已按 `user_id` 归属，镜像和系统接口是全局的。已支持修改密码和登录限流。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 3.1 | 修改密码后旧令牌立即失效 | [ ] | `User` 增加 `token_version`，写进 JWT，`middleware/auth.go` 校验；修改密码时加一。顺带实现"退出所有设备" |
| 3.2 | 角色：`admin` / `user`，初始化创建的是 admin | [ ] | 镜像管理、系统概览、用户管理只允许 admin |
| 3.3 | 用户管理接口和设置页：新增、删除、重置密码、改角色 | [ ] | 不能删除最后一个 admin；删除用户时其服务器的处理方式要先定（转给 admin 或禁止删除） |
| 3.4 | 服务器共享 / 授权（可选） | [ ] | 先确认是否需要；最简单是 admin 可见全部，普通用户只看自己的。注意 `checkPorts` 已经按全部服务器检查，多用户后不要改成只查自己的 |

## 四、工具版本更新检查

现状：管理器自身没有版本号，无法检查更新。ARK 服务器镜像已有检查更新。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 4.1 | 构建时注入版本号（Git tag / commit）：后端 `-ldflags -X`，前端构建参数 | [ ] | 改 `Dockerfile` 和 `.github/workflows/docker-build.yml` |
| 4.2 | `GET /system/version`：当前版本、最新版本、是否有更新 | [ ] | 最新版本取 GitHub Releases，或用 `DistributionInspect` 比较 `tbro98/arkservercommander:latest` 的 digest（与镜像检查同一套逻辑）；结果缓存 |
| 4.3 | 前端：设置页 / 页脚显示版本，有更新时提示升级命令 | [ ] | 管理器容器不自我更新（要替换正在运行的自身），只给出 `docker-compose pull && docker-compose up -d` |

## 五、MCP 支持

让 AI 助手通过 MCP 管理服务器。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 5.1 | API 令牌：设置页生成 / 吊销长期令牌，只保存哈希 | [ ] | MCP 和脚本调用都需要，不能复用 24 小时的登录 JWT |
| 5.2 | MCP 服务端（Streamable HTTP），挂在后端 `/mcp`，经 Next 代理对外 | [ ] | Go 用官方 `modelcontextprotocol/go-sdk`；代理需要支持流式（通配代理已支持） |
| 5.3 | 工具：列出服务器和状态、启动 / 停止 / 重启、执行 RCON、在线玩家、创建备份、读取最近日志 | [ ] | 直接调用 `service/server`，权限与令牌所属用户一致；危险操作（删除、恢复）先不提供 |
| 5.4 | 文档：README 增加 MCP 客户端配置示例 | [ ] | |

## 六、遗留改进

来自 ROADMAP "后续可以考虑"和 CLAUDE.md 已知问题。

| # | 项目 | 状态 | 备注 |
|---|------|------|------|
| 6.1 | 查看游戏日志 `ShooterGame.log`（数据卷 `Saved/Logs`） | [ ] | 直接读数据目录 `servers/<id>/saved/Logs`；日志标签里切换"容器输出 / 游戏日志" |
| 6.2 | 日志查看增加"隐藏下载进度行"开关 | [ ] | `LogViewer.tsx` 已能识别进度行（只是变灰），加过滤即可 |
| 6.3 | steamcmd 下载阶段转发 SIGTERM | [ ] | `start_server.sh` 现在只在启动游戏后设置 trap（第 99 行），下载中停止要等 30 秒超时；需要重新构建镜像 |
| 6.4 | 集群清理：列出不再被任何服务器使用的集群卷，提供删除 | [ ] | 现在要手动 `docker volume rm ase-cluster-<id>` 并删除 `data/clusters/<id>` |
| 6.5 | 验证默认 Game.ini 模板中倍率类 key 的位置 | [ ] | `XPMultiplier`、`TamingSpeedMultiplier` 等可能应在 GameUserSettings.ini 的 `[ServerSettings]`，需真实服务器验证后改 `utils/config_files.go` |
| 6.6 | 更新 `docs/ISSUES.md` | [ ] | 2026-02 的报告，其中端口冲突、登出、登录限流等已修复，内容与现状不符 |

## 长期计划（暂不排期）

- 多主机管理（可能基于 K8S）
- 服务器收录网站
- 玩家使用界面
