# ARK Server Commander

> ⚠️ **Development Stage Notice**: This project is currently in development stage and features may be incomplete or have stability issues. It is recommended for testing environments only and should not be used in production.

[English](README.md) | [中文](README-zh.md)

- ARK Survival Evolved server management tool for Linux.
- ARK servers come with ArkApi plugin system built-in.

## 🎮 Features

### ✅ Implemented Features
- 🐳 Each ARK server runs in an independent Docker container
- 🔌 Servers come with ArkApi pre-installed
- 🔄 Server containers support automatic restart on crash
- ⬆️ Automatic server files and mod updates on first startup
- 💾 Automatic creation and management of Docker volumes for game data storage
- 🖥️ Add and manage multiple ARK servers
- ⚙️ Configure server settings and configuration parameters
- ▶️ One-click server start/stop

### 🚧 Planned Features
- 🎮 RCON command execution
- 📊 Server running status monitoring
- 🎨 Mod management integration with Steam Workshop
- 🔧 ArkApi plugin management
- 📋 Server log viewing
- 💾 Server save and configuration backup
- 🔍 Tool version update checking
- 🖼️ Server image updates
- ⚡ Optional server files and mod updates
- 🌐 i18n internationalization support
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
version: '3.8'

services:
  # ARK Server Management System (Frontend + Backend)
  ark-commander:
    image: tbro98/arkservercommander:latest
    container_name: ark-commander
    ports:
      # You can modify the port mapping
      - "8080:8080"
    environment:
      # Will be optimized later, keep default for now
      - JWT_SECRET=your-secret-key-here
      - DB_PATH=/data/ark_server.db
      - SERVER_PORT=8080
    volumes:
      # Database storage
      - ./data:/data
      # Docker socket (for managing Docker containers)
      - /var/run/docker.sock:/var/run/docker.sock
    restart: unless-stopped
    # System requires privileged mode to operate host Docker
    privileged: true

```

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
A: Server data backup operations are not yet implemented. Server data is stored in Docker volumes ark-server-<server_number>, you can backup manually.

### ❓ Q: How to view ARK server logs?
A: The server program currently cannot output logs directly in docker logs, you need to view the server log files, and we'll see how to optimize this later.

### 🖼️ ARK Server Image
- This system uses the `tbro98/ase-server:latest` image to run ARK servers
- Image source: [ASE-Server-Docker](https://github.com/tbro199803/ASE-Server-Docker)

## 📸 Interface Screenshots
![](./docs/imgs/img_servers.png)
![](./docs/imgs/ima_base.png)
![](./docs/imgs/img_GameUserSettings.png)
![](./docs/imgs/img_GameIni.png)
![](./docs/imgs/img_args.png) 