# ARK Server Commander

> ⚠️ **Development Stage Notice**: This project is currently in development stage and features may be incomplete or have stability issues. It is recommended for testing environments only and should not be used in production.

[English](README.md) | [中文](README-zh.md)

- ARK Survival Evolved server management tool for Linux.
- ARK servers come with ArkApi plugin system built-in.

## 🎮 Features

### ✅ Implemented Features
- 🐳 Each ARK server runs in an independent Docker container
- 🔌 Servers come with ArkApi pre-installed, with the plugin directory mounted separately
- 🔄 Server containers support automatic restart on crash
- ⬆️ Server files and mods are downloaded on first start and kept on their own volume, so recreating a container does not download them again; new servers copy server files from an already installed server
- 💾 Saves, configs, plugins and server files live under `data/servers/<id>/` and can be managed directly on the host; they are kept when a server is deleted unless you choose to delete them
- 🖥️ Add and manage multiple ARK servers, with cluster support (shared cluster directory)
- ⚙️ Edit GameUserSettings.ini, Game.ini and launch arguments visually or as text
- ▶️ One-click server start/stop/restart; containers are recreated automatically on start after launch arguments or mods change
- 📊 Status monitoring: download and mod install progress, joinable state, online players, CPU/memory/disk
- 📋 Live server log viewing (container output and the game's ShooterGame.log)
- 🎮 RCON console: run commands, broadcast, save world; kick/ban online players
- 💾 Save backups: manual backup, download, upload, restore
- ⏰ Scheduled tasks: restart (with advance broadcast), backup, update, command, start/stop
- ⚡ One-click server file and mod updates, or automatic update on every start
- 🎨 Mod management: Workshop names and previews, paste a Workshop link to add a mod, drag to reorder, detect mod updates and update only the changed mods
- 🔧 ArkApi plugin management: install from zip, enable/disable, delete, and edit plugin and ArkApi config.json online
- 🖼️ Docker image management: download, check for updates, update; after an update, servers recreate their containers with the new image on next start, or running servers can be restarted right away
- 🔐 Administrator login (JWT authentication, secret generated automatically by the backend, no configuration needed), with password change
- 🔍 Tool version update checking (settings page)
- 🧹 Clean up cluster data that is no longer used
- 🌐 Chinese / English interface
- 📝 API documentation (Swagger, available when running the backend locally)

### 🚧 Planned Features
See [docs/DEV_PLAN.md](docs/DEV_PLAN.md) for the development plan and progress.
- 👤 Multi-user management
- 🔌 MCP (Model Context Protocol) support

### 🚀 Future Plans
- ☸️ Multi-host management based on K8S
- 🌍 Server listing website, breaking free from poor Steam server search
- 👥 Player user interface

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

### ❓ Q: Where is the login secret? Do I need to configure it?
A: No. On first start the backend generates a strong random secret and stores it in `jwt_secret` in the data directory (the mounted `./data/jwt_secret`). Deleting that file and restarting the container invalidates all logged-in sessions.

### 🖼️ ARK Server Image
- This system uses the `tbro98/ase-server:latest` image to run ARK servers
- Image source: [ASE-Server-Docker](https://github.com/tbro199803/ASE-Server-Docker)

## 📸 Interface Screenshots
![](./docs/zh/images/img_servers.png)
![](./docs/zh/images/ima_base.png)
![](./docs/zh/images/img_GameUserSettings.png)
![](./docs/zh/images/img_GameIni.png)
![](./docs/zh/images/img_args.png) 