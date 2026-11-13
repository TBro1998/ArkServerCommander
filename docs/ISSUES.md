# ARK Server Commander - 问题汇总

> 原报告生成于 2026-02-01，2026-11-13 按当前代码逐项复核。
> 仍未解决的问题同时记录在 CLAUDE.md 的"已知问题"，新的开发计划见 [DEV_PLAN.md](DEV_PLAN.md)。

## 复核结果

| # | 原问题 | 状态 | 现状 |
|---|--------|------|------|
| 1 | Docker 卷删除用字符串拼接插件卷名 | 已解决 | `docker_volume.go` 的 `serverVolumeNames` 从 `serverVolumeSpecs` 取三个卷名，与创建时同一来源；删除数据时还会删除 `servers/<id>/` 目录（`RemoveServerDirs`） |
| 2 | 创建 / 更新服务器时没有端口冲突检测 | 已解决 | `checkPorts` 检查全部服务器（不分用户），包括游戏端口 +1；新建时 `GET /system/suggest-ports` 推荐可用端口 |
| 3 | 首次启动的更新完全依赖镜像，应用层无法感知进度 | 已解决 | 运行信息监控（`runtime.go`）解析日志得到下载进度、Mod 安装进度和是否可连接；日志页实时查看容器输出和游戏日志；可单独更新游戏本体、全部或部分 Mod |
| 4 | GameUserSettings.ini / Game.ini 内容没有格式校验 | 已解决 | 创建和更新时调用 `utils.ValidateINIContent`；可视化编辑器只合并改过的键 |
| 5 | 认证缺少刷新、登出、密码强度、登录限流 | 部分解决 | 已有：登出（清除 HttpOnly cookie）、密码 8–72 字节、按 IP 和用户名的登录限流、修改密码。未做：令牌刷新；修改密码或登出后旧 JWT 在 24 小时有效期内仍可用（需要令牌版本号，见 DEV_PLAN 3.1） |
| 6 | 异步更新状态没有并发保护 | 已解决 | `transitionStatus` 用 `UPDATE ... WHERE status NOT IN (...)` 原子切换，结果由 `finishOperation` 写回；删除以 `status = stopped` 为条件 |
| 7 | 镜像拉取状态是全局变量，多实例部署会不同步 | 不处理 | 项目设计为单实例（SQLite、单容器、本机 Docker），不支持多实例部署。单实例内的并发拉取由 `ClaimPull` 原子占用，重复请求返回 409 |
| 8 | 错误信息不友好，没有统一错误码 | 已解决 | 服务层哨兵错误 + `errorCodes`，响应带 `error`（中文详情）和 `code`，前端 `useApiError()` 按语言显示 |
| 9 | 关键操作缺少日志，没有审计记录 | 部分解决 | 后端使用 zap 结构化日志，启停、镜像、备份等操作有日志。没有"谁在什么时间做了什么"的操作审计记录，多用户之前意义不大 |

## 仍需关注

- 旧 JWT 在修改密码后仍可用 24 小时（问题 5）
- 没有操作审计日志（问题 9）
- 其他已知限制（登录限流可被伪造的 `X-Forwarded-For` 绕过 IP 维度、数据库里的管理员密码是明文等）见 CLAUDE.md "已知问题"
