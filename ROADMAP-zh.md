# 路线图

[English](ROADMAP.md) | 中文

这里列出 ARK Server Commander 接下来打算做的事。没有固定日期，按优先级分为近期、中期和远期，顺序可能会随反馈调整。每个版本具体改了什么见 [Releases](https://github.com/TBro1998/ArkServerCommander/releases)。

想要某个功能、或者觉得哪项应该先做，欢迎提 [Issue](https://github.com/TBro1998/ArkServerCommander/issues)。

## 最近完成

- 🤖 MCP 接口：Claude Code、Cursor 等 Agent 可以直接管理服务器，令牌分完全控制和只读
- 🌍 一键发布到 [arkfgf.com](https://arkfgf.com) 收录站，自动上报真实在线人数
- 🔗 P2P 联机房：没有公网 IP 也能让玩家通过 [ArkServerLink](ArkServerLink/) 进服
- 🧩 ArkApi 插件目录：从 ark-server-api.com 一键安装和更新插件，config.json 树形编辑
- 🕹️ 可以选择游戏版本（蜗牛海洋正式版 / 通配符原版）
- ⬆️ 设置页一键更新管理端

## 近期：在真实环境中验证

这些功能已经写完，但只在模拟环境或本机 Docker Desktop 上测过，需要在真实的 Linux 服务器上走一遍完整流程：

- [ ] 收录站发布、arkfgf.com 账号绑定、P2P 联机房（包括一个房间多台服务器）与网站联调
- [ ] 原生 Linux 版服务端镜像，以及在同一台服务器上切换 Linux / Proton
- [ ] 通配符原版在 Linux 版服务端和 ArkApi 下是否可用
- [ ] 管理端一键更新在 Linux 上的完整流程（compose 创建的容器、镜像加速器、回滚）
- [ ] 停服时游戏能否在 `docker stop` 的 30 秒内完成存档

同时处理安全方面的遗留问题：

- [ ] 修改登录密码后旧的登录令牌立即失效，支持“退出所有设备”

## 中期

- [ ] 👤 多用户：管理员 / 普通用户角色，在设置页新增、删除用户和重置密码
- [ ] 🔐 服务器授权：普通用户只能看到和管理分配给自己的服务器
- [ ] 📜 操作记录：谁在什么时候对哪台服务器做了什么
- [ ] 🔑 MCP 令牌可以限制只能管理指定的服务器
- [ ] 🔍 创意工坊站内搜索 Mod（需要填写自己的 Steam Web API Key）
- [ ] 💽 减少磁盘占用：现在每台服务器各有一份约 20 GB 的服务端文件，考虑让多台服务器共用（方案待定）

## 远期

- [ ] ☸️ 多主机管理：一个管理端管理多台机器上的服务器
- [ ] 👥 玩家使用界面
