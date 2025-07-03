# ARK 服务器管理器

一个简单易用的 ARK 生存进化服务器管理工具，帮助您轻松管理和配置 ARK 服务器。

## 🎮 功能特性

### Docker容器化管理
- 🐳 每个ARK服务器运行在独立的Docker容器中
- 📦 自动创建和管理Docker卷存储游戏数据
- 🔄 容器状态实时同步
- 🛡️ 隔离的运行环境，提高安全性和稳定性

### 服务器管理
- 🖥️ 添加和管理多个 ARK 服务器
- ⚙️ 配置服务器设置（端口、地图、玩家数量等）
- 📝 在线编辑 GameUserSettings.ini 和 Game.ini
- 🎛️ RCON 命令执行
- ▶️ 一键启动/停止服务器
- 🗑️ 安全删除服务器配置和数据

### 配置文件管理
- 📄 可视化编辑 GameUserSettings.ini
- 📄 可视化编辑 Game.ini
- ✅ 配置文件格式验证
- 💾 实时保存到 Docker 卷

### 安全管理
- 🔐 用户登录保护
- 🔑 RCON 密码管理
- 👤 个人数据隔离
- 🛡️ JWT 令牌认证

### 用户界面
- 🎨 现代化界面设计
- 📱 移动设备友好
- ⚡ 快速响应操作
- 💬 友好的提示信息
- 📊 实时服务器状态显示

## 🚀 快速开始

### Docker容器化部署（推荐）

1. **克隆项目**
```bash
git clone <repository-url>
cd ArkServerManager
```

2. **使用Docker Compose部署**
```bash
# 复制配置文件
cp docker-compose.example.yml docker-compose.yml

# 启动服务
docker-compose up -d
```

3. **访问管理界面**
- 前端界面：http://localhost:3000
- API文档：http://localhost:8080/swagger/index.html

### 端口分配说明

管理系统端口：
- 前端Web界面：3000
- 后端API：8080

ARK服务器端口（每个服务器独立）：
- 游戏端口：7777, 7778, 7779... (递增分配)
- 查询端口：27015, 27016, 27017... (递增分配)  
- RCON端口：32330, 32331, 32332... (递增分配)

Docker容器内部端口固定：
- 游戏端口：7777/udp
- 查询端口：27015/udp
- RCON端口：32330/tcp

### 传统部署方式

1. **启动后端服务**
```bash
cd server
go run main.go
```

2. **启动前端服务**
```bash
cd web
npm install
npm run dev
```

3. **访问应用**
- 前端：http://localhost:3000
- 后端API：http://localhost:8080

## 📖 使用说明

### 首次使用
1. 访问 http://ip:3000
2. 系统会自动跳转到初始化页面
3. 设置您的管理员账号和密码
4. 初始化完成后登录系统

### 管理服务器
1. 登录后点击"服务器管理"
2. 点击"添加服务器"创建新的服务器配置
3. 填写服务器信息：
   - **服务器名称**：自定义名称
   - **IP地址**：服务器IP地址
   - **游戏端口**：默认 7777
   - **查询端口**：默认 27015
   - **RCON端口**：默认 32330
   - **地图**：选择游戏地图
   - **最大玩家数**：服务器人数上限

### 支持的地图
- The Island (孤岛)
- The Center (中心岛)
- Scorched Earth (焦土)
- Aberration (畸变)
- Extinction (灭绝)
- Valguero (瓦尔盖罗)
- Genesis (创世纪)
- Crystal Isles (水晶岛)
- Genesis 2 (创世纪2)
- Lost Island (失落岛)
- Fjordur (峡湾)

## 🌐 访问地址

启动成功后，您可以通过以下地址访问：

- **管理界面**：http://localhost:3000
- **API文档**：http://localhost:8080/swagger/index.html

## ❓ 常见问题

### Q: 无法访问管理界面？
A: 请确保前后端服务都已启动，并检查端口 3000 和 8080 是否被占用。

### Q: 可以同时管理多少个服务器？
A: 没有数量限制，您可以根据需要添加任意数量的服务器配置。

### Q: Docker容器无法启动？
A: 检查以下几点：
- Docker服务是否正常运行
- 是否有足够的磁盘空间
- 端口是否被占用
- 镜像 `tbro98/ase-server:latest` 是否能正常拉取

### Q: 如何备份ARK服务器数据？
A: 服务器数据存储在Docker卷中，可以使用以下命令备份：
```bash
# 备份卷数据
docker run --rm -v ark-server-1:/source -v $(pwd):/backup alpine tar czf /backup/ark-server-1-backup.tar.gz -C /source .

# 恢复卷数据
docker run --rm -v ark-server-1:/target -v $(pwd):/backup alpine tar xzf /backup/ark-server-1-backup.tar.gz -C /target
```

### Q: 如何查看ARK服务器日志？
A: 使用Docker命令查看容器日志：
```bash
# 查看特定服务器的日志
docker logs ark-server-1

# 实时跟踪日志
docker logs -f ark-server-1
```

## 🔧 系统要求

### Docker容器化部署
- Docker Engine 20.10+
- Docker Compose 2.0+
- 每个ARK服务器 8GB+ 内存 (推荐)
- 每个ARK服务器 10GB+ 磁盘空间
- 确保Docker有足够权限管理容器和卷

### 传统部署
- Go 1.23+
- Node.js 18+
- 每个服务器 8GB+ 内存 (推荐)

### ARK服务器镜像
本系统使用 `tbro98/ase-server:latest` 镜像来运行ARK服务器


## 📄 许可证

本项目基于 MIT 许可证开源。
