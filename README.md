# ARK Server Commander

> ⚠️ **Development Stage Notice**: This project is currently in development stage and features may be incomplete or have stability issues. It is recommended for testing environments only and should not be used in production.

[English](README.md) | [中文](README-zh.md)

- ARK Survival Evolved server management tool for Linux.
- ARK servers come with ArkApi plugin system built-in.

## 🎮 Features

### ✅ Implemented Features
- 🐳 Each ARK server runs in an independent Docker container
- 🔌 Servers come with ArkApi pre-installed, with the plugin directory on its own Docker volume
- 🔄 Server containers support automatic restart on crash
- ⬆️ Server files are updated and mods downloaded on a container's first start (repeated after the container is recreated)
- 💾 Docker volumes are created automatically for saves, configs and plugins, and are kept when a server is deleted
- 🖥️ Add and manage multiple ARK servers
- ⚙️ Edit GameUserSettings.ini, Game.ini and launch arguments visually or as text
- ▶️ One-click server start/stop; containers are recreated automatically on start after launch arguments or mods change
- 🖼️ Docker image management: download, check for updates, update; after an update, servers recreate their containers with the new image on next start, or running servers can be restarted right away
- 🔐 Administrator login (JWT authentication, secret generated automatically by the backend, no configuration needed)
- 🌐 Chinese / English interface
- 📝 API documentation (Swagger, available when running the backend locally)

### 🚧 Planned Features
- 🎮 RCON command execution
- 📊 Server running status monitoring
- 🎨 Mod management integration with Steam Workshop
- 🔧 ArkApi plugin management
- 📋 Server log viewing
- 💾 Server save and configuration backup
- 👤 User management (password change, multiple users)
- 🔍 Tool version update checking
- ⚡ Optional server files and mod updates
- 🔌 MCP (Mod Configuration Protocol) support

### 🚀 Future Plans
- ☸️ Multi-host management based on K8S
- 🌍 Server listing website, breaking free from poor Steam server search
- 👥 Player user interface

## 🚀 Quick Start

### 🔧 System Requirements

- 8GB+ memory per ARK server (recommended)
- 10GB+ disk space per ARK server

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
A: Server data backup is not implemented yet. Saves and configs are stored in the Docker volume `ase-server-<server_number>` and ArkApi plugins in `ase-server-plugins-<server_number>`; you can back them up manually. Deleting a server does not remove these volumes; run `docker volume rm` when you no longer need them.

### ❓ Q: How to view ARK server logs?
A: The server program currently cannot output logs directly in docker logs, you need to view the server log files, and we'll see how to optimize this later.

### ❓ Q: How to update ARK server images?
A: After logging in, click "Check Updates" under image management on the home page, then click "Update" when a new version is available. When confirming you can choose:
- **Update Image Only**: running servers are not affected and recreate their containers with the new image on next start
- **Update Image and Restart Running Servers**: once the download finishes, those servers are stopped, recreated and started again

Servers still on the old image show an "Old Image" badge on their card. The first start after a container is recreated re-validates server files and re-downloads mods, which takes a while; saves and configs live in volumes and are not lost.

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