# ARK Server Commander

> ⚠️ **Development Stage Notice**: This project is currently in development stage and features may be incomplete or have stability issues. It is recommended for testing environments only and should not be used in production.

[English](README.md) | [中文](README-zh.md)

- ARK Survival Evolved server management tool for Linux.
- Runs the native Linux server by default; ArkApi can be enabled per server when you need plugins (the Windows server runs under GE-Proton).

## 🎮 Features

### ✅ Implemented Features
- 🐳 Each ARK server runs in an independent Docker container
- 🔌 ArkApi can be enabled per server: enabled servers run the Windows build on the Proton image, others run the native Linux build; switching takes effect on the next start
- 🔄 Server containers support automatic restart on crash
- ⬆️ Server files and mods are downloaded on first start and kept on their own volume, so recreating a container does not download them again; new servers copy server files from an already installed server
- 💾 Saves, configs, plugins and server files live under `data/servers/<id>/` (except server files on Docker Desktop, see the FAQ) and can be managed directly on the host; they are kept when a server is deleted unless you choose to delete them
- 🖥️ Add and manage multiple ARK servers, with cluster support (shared cluster directory)
- ⚙️ Edit GameUserSettings.ini, Game.ini and launch arguments visually or as text
- 📤 Export a server's config to a file, import it from a file or another server, or copy it to several servers at once (choose game settings, launch arguments, mods and more; names, ports and admin passwords are never overwritten)
- ▶️ One-click server start/stop/restart; containers are recreated automatically on start after launch arguments or mods change
- 📊 Status monitoring: download and mod install progress, joinable state, online players, CPU/memory/disk
- 📋 Live server log viewing (container output and the game's ShooterGame.log)
- 🎮 RCON console: run commands, broadcast, save world; kick/ban online players
- 💾 Save backups: manual backup, download, upload, restore
- ⏰ Scheduled tasks: restart (with advance broadcast), backup, update, command, start/stop
- ⚡ One-click server file and mod updates, or automatic update on every start
- 🎨 Mod management: Workshop names and previews, paste a Workshop link to add a mod, drag to reorder, detect mod updates and update only the changed mods
- 🔧 ArkApi management: 3.56 is built in, other versions can be downloaded from GitHub or uploaded as a zip and installed with one click; plugins can be installed from zip, enabled/disabled, deleted, and their config.json and ArkApi's edited online
- 🖼️ Docker image management: download, check for updates, update; after an update, servers recreate their containers with the new image on next start, or running servers can be restarted right away
- 🔐 Administrator login (JWT authentication, secret generated automatically by the backend, no configuration needed), with password change
- 🔍 Tool version update checking (settings page)
- 🧹 Clean up cluster data that is no longer used
- 🌐 Chinese / English interface
- 📝 API documentation (Swagger, available when running the backend locally)
- 🤖 MCP endpoint for agents (Claude Code, Claude Desktop, Cursor…), with full-control or read-only access tokens

### 🗺️ Roadmap
See [ROADMAP.md](ROADMAP.md) for what's planned next. Highlights:
- 🧪 Verifying recent features on real Linux servers (native Linux image, P2P rooms, one-click update, etc.)
- 👤 Multi-user management and per-server access
- ☸️ Multi-host management

## 🚀 Quick Start

### 🔧 System Requirements

- 8GB+ memory per ARK server (recommended)
- 20GB+ disk space per ARK server

### 📦 Install 1Panel
We recommend installing 1Panel on your server to manage the Docker environment
[1Panel Installation Guide](https://1panel.cn/docs/v2/installation/online_installation/)

### 🐳 Docker Containerized Deployment

Copy the docker-compose.yml, or use the following configuration directly:
```yml
services:
  # ARK Server Management System (Frontend + Backend)
  ark-commander:
    image: tbro98/arkservercommander:latest
    container_name: ark-commander
    ports:
      # You can modify the port mapping
      - "8080:8080"
    volumes:
      # Database and the auto-generated login secret
      - ./data:/data
      # Docker socket (for managing Docker containers)
      - /var/run/docker.sock:/var/run/docker.sock
    restart: unless-stopped
    # System requires privileged mode to operate host Docker
    privileged: true
```

No extra configuration is needed. If you access the system over HTTPS, add the environment variable `COOKIE_SECURE=true`; see [.env.example](.env.example) for other options.

# 🚀 Start the service
```
sudo docker-compose up -d
```

🌐 Access the system interface via ip+port

## 📖 User Guide

### 🆕 First Time Use
1. The system will automatically redirect to the initialization page
2. Set up your administrator account and password
3. After initialization, log into the system

### 🖥️ Managing Servers
1. After logging in, click "Server Management"
2. Click "Add Server" to create a new server configuration

### 🤖 Agent Access (MCP)
The commander exposes an MCP (Model Context Protocol) server so agents can manage servers without the web UI: list status, start/stop/restart, create and edit servers, change individual config keys, export, import and copy server configs, run RCON commands, view logs and players, manage backups, schedules, plugins and images.

1. Open **Settings → Agent access (MCP)**, enter a token name and click **Create** (turn on **Read-only** for an agent that should only look). The token is shown only once.
2. Add the server to your agent. For Claude Code:
   ```bash
   claude mcp add --transport http ark-commander http://<host>:8080/mcp --header "Authorization: Bearer <token>"
   ```
   Other clients use the same URL (Streamable HTTP) and `Authorization` header; the settings page shows a ready-to-copy JSON config.
3. Revoke a token on the same page; it stops working immediately.

**Let your agent install it.** Replace `<host>` and `<token>` below, then paste the whole paragraph to your agent (Claude Code, Cursor, Codex, etc.):

```text
Please connect to my ARK Server Commander via MCP. Add an MCP server named "ark-commander" using the Streamable HTTP transport, URL http://<host>:8080/mcp, with the request header "Authorization: Bearer <token>". If you are Claude Code, run: claude mcp add --transport http --scope user ark-commander http://<host>:8080/mcp --header "Authorization: Bearer <token>". Otherwise add the equivalent entry to your MCP client config file ({"mcpServers": {"ark-commander": {"type": "http", "url": "http://<host>:8080/mcp", "headers": {"Authorization": "Bearer <token>"}}}}), merging it with any existing servers instead of overwriting them. Keep the token out of version control and do not echo it back to me. When done, tell me if a restart or reload is needed to load the server; once its tools are available, call list_servers and summarize each server's name, status, map and online players.
```

The endpoint uses the same port as the web UI. If it is reachable from the internet, put it behind an HTTPS reverse proxy, because the token is sent with every request. Uploading plugins and uploading/downloading backups are only available in the web UI.

### 🗺️ Supported Maps - Can be extended with custom additions later
- The Island
- The Center
- Scorched Earth
- Aberration
- Extinction
- Valguero
- Genesis
- Crystal Isles
- Genesis 2
- Lost Island
- Fjordur

## ❓ FAQ

### ❓ Q: How to backup ARK server data?
A: Use the "Backups" tab on the server detail page to back up, download, upload and restore saves, or add a scheduled backup under "Schedules". Backups are stored in `data/backups/<server_id>/`. All server data lives under `data/servers/<server_id>/` on the host, which you can also back up yourself.

### ❓ Q: How to view ARK server logs?
A: Open the "Logs" tab on the server detail page to follow logs live, pause, filter and download them. Switch to "Game log" to see the game's own `ShooterGame.log` (stored under `data/servers/<server_id>/saved/Logs`).

### ❓ Q: How to update ARK server images?
A: After logging in, click "Check Updates" under image management on the home page, then click "Update" when a new version is available. When confirming you can choose:
- **Update Image Only**: running servers are not affected and recreate their containers with the new image on next start
- **Update Image and Restart Running Servers**: once the download finishes, those servers are stopped, recreated and started again

Servers still on the old image show an "Old Image" badge on their card. Server files, saves and configs live in volumes, so recreating a container loses nothing and does not download them again.

### ❓ Q: How to update the commander itself?
A: In "Settings → About", click "Check for updates", then "Update now" when a new version is available. The commander downloads the new image, and a temporary updater container (`ark-commander-updater`) recreates the commander container with it. The page is unavailable for about a minute and reloads automatically when done. Running ARK servers and the P2P room keep running; if the new version does not start properly, the previous version is restored automatically. The database is copied to `./data/ark_server.db.before-update` first.

One-click update requires running the `tbro98/arkservercommander:latest` image from `docker-compose.yml` with `/var/run/docker.sock` mounted. If you pinned a version or build your own image, run `docker compose pull && docker compose up -d` in the directory containing `docker-compose.yml`.

### ❓ Q: Where is the login secret? Do I need to configure it?
A: No. On first start the backend generates a strong random secret and stores it in `jwt_secret` in the data directory (the mounted `./data/jwt_secret`). Deleting that file and restarting the container invalidates all logged-in sessions.

### ❓ Q: Can I run it on Docker Desktop (Windows / macOS)?
A: Yes, without changing `docker-compose.yml`. Docker Desktop reaches the host disk (`./data`) through file sharing, which makes downloading the ~20 GB of server files very slow. On Docker Desktop the manager therefore puts **the server files of new servers** in a volume built into the image (Docker's own filesystem); saves, configs, plugins and backups stay in `./data` where you can manage them directly. Server files already downloaded into `./data` keep being used and are not downloaded again. These server files are not visible on the host. After `docker compose down` the next start gets a new volume, and the manager copies the server files over automatically (this can take a few minutes). Set the environment variable `SERVER_FILES_VOLUME=off` to turn this off.

### 🖼️ ARK Server Image
- Servers without ArkApi use `tbro98/ase-server:linux` (native Linux server); servers with ArkApi use `tbro98/ase-server:proton` (Windows server under GE-Proton)
- A missing image is downloaded automatically when a server that needs it starts; the images no longer ship ArkApi, the manager installs it into the server files
- `tbro98/ase-server:latest` is no longer updated. It stays at the last version that bundled ArkApi, for older manager versions
- Image source: [ASE-Server-Docker](https://github.com/tbro199803/ASE-Server-Docker)

### 🔗 P2P Rooms (no public IP)
- Turn on "P2P room" in a server's settings and players join through [arkfgf.com/rooms](https://arkfgf.com/rooms) with the open-source client [ArkServerLink](https://github.com/TBro1998/ArkServerLink) (submodule `ArkServerLink/`, download at [arkfgf.com/download](https://arkfgf.com/download))

## 📸 Interface Screenshots
**Home: servers, online players, resource usage, images and disk**

![Home: servers, online players, resource usage, images and disk](./docs/zh/images/home.png)

**Server list: card view with start/stop and shortcuts to console, logs and settings**

![Server list: card view with start/stop and shortcuts to console, logs and settings](./docs/zh/images/servers.png)

**List view: one server per row with status, players, resources and address**

![List view: one server per row with status, players, resources and address](./docs/zh/images/servers-list.png)

**Server details: connection info, configuration and mod updates**

![Server details: connection info, configuration and mod updates](./docs/zh/images/server-detail.png)

**Manage: update game/mods, recreate the container, export, import and copy configs, delete the server**

![Manage: update game/mods, recreate the container, export, import and copy configs, delete the server](./docs/zh/images/server-manage.png)

**Scheduled tasks: restart, backup and update with countdown broadcasts**

![Scheduled tasks: restart, backup and update with countdown broadcasts](./docs/zh/images/server-schedules.png)

**Edit server: settings grouped by purpose, side navigation with counts of changed items**

![Edit server: settings grouped by purpose, side navigation with counts of changed items](./docs/zh/images/form-general.png)

**Mods and clusters: Workshop mod names and load order, cross-server transfers**

![Mods and clusters: Workshop mod names and load order, cross-server transfers](./docs/zh/images/form-mods.png)

**Rates: settings from GameUserSettings.ini, Game.ini and startup arguments in one place**

![Rates: settings from GameUserSettings.ini, Game.ini and startup arguments in one place](./docs/zh/images/form-rates.png)

**Search: find settings by name, key or description across all groups**

![Search: find settings by name, key or description across all groups](./docs/zh/images/form-search.png)
