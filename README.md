# ARK 服务器管理器

一个基于 Nuxt.js + Gin + Gorm 的前后端分离 ARK 服务器管理系统。

## 项目结构

```
ArkServerManager/
├── server/          # 后端 (Gin + Gorm)
│   ├── config/      # 配置文件
│   ├── controllers/ # 控制器
│   ├── database/    # 数据库相关
│   ├── middleware/  # 中间件
│   ├── models/      # 数据模型
│   ├── routes/      # 路由
│   ├── utils/       # 工具函数
│   ├── main.go      # 主入口
│   └── go.mod       # Go模块
└── web/             # 前端 (Nuxt.js)
    ├── layouts/     # 布局
    ├── middleware/  # 中间件
    ├── pages/       # 页面
    ├── stores/      # 状态管理
    ├── app.vue      # 应用入口
    ├── nuxt.config.ts
    └── package.json
```

## 功能特性

### 🔐 认证系统
- ✅ JWT 身份认证
- ✅ 系统初始化设置
- ✅ 用户登录/登出
- ✅ 密码加密存储 (bcrypt)

### 🖥️ 服务器管理
- ✅ 服务器列表查看
- ✅ 新增 ARK 服务器配置
- ✅ 编辑服务器配置
- ✅ 删除服务器配置
- ✅ 多地图支持（The Island、Scorched Earth 等）
- ✅ 端口配置管理（游戏端口、查询端口、RCON端口）

### 🎨 用户界面
- ✅ 响应式设计 (Tailwind CSS)
- ✅ 现代化UI界面
- ✅ 错误和成功消息提示
- ✅ 加载状态指示
- ✅ 确认对话框

### 🔧 技术特性
- ✅ 前后端分离架构
- ✅ RESTful API 设计
- ✅ Swagger API 文档
- ✅ SQLite 数据库
- ✅ CORS 跨域支持
- ✅ 用户级数据隔离

## 技术栈

### 后端
- **Gin**: Web 框架
- **Gorm**: ORM 数据库操作
- **SQLite**: 数据库
- **JWT**: 身份认证
- **bcrypt**: 密码加密
- **Swagger**: API 文档生成
- **CORS**: 跨域资源共享

### 前端
- **Nuxt.js 3**: Vue.js 框架
- **Pinia**: 状态管理
- **Tailwind CSS**: 样式框架
- **TypeScript**: 类型支持
- **Composition API**: Vue 3 响应式 API

## 快速开始

### 🚀 一键启动 (推荐)

双击运行 `start-dev.bat` 脚本，自动启动前后端服务：

```bash
start-dev.bat
```

### 📋 手动启动

#### 后端启动

1. 进入后端目录：
```bash
cd server
```

2. 安装依赖：
```bash
go mod tidy
```

3. 启动服务器：
```bash
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

#### 前端启动

1. 进入前端目录：
```bash
cd web
```

2. 安装依赖：
```bash
npm install
# 或使用 pnpm
pnpm install
```

3. 启动开发服务器：
```bash
npm run dev
# 或使用 pnpm
pnpm dev
```

前端服务将在 `http://localhost:3000` 启动

### 🌐 访问地址

- **前端界面**: http://localhost:3000
- **后端API**: http://localhost:8080
- **Swagger API文档**: http://localhost:8080/swagger/index.html

## 使用说明

1. **初次使用**：
   - 访问 `http://localhost:3000`
   - 系统会自动跳转到初始化页面
   - 设置管理员用户名和密码

2. **登录系统**：
   - 使用刚才设置的用户名和密码登录
   - 登录成功后进入管理界面

3. **服务器管理**：
   - 点击导航栏的"服务器管理"进入管理界面
   - 可以新增、编辑、删除 ARK 服务器配置
   - 支持配置服务器名称、IP、端口、地图等信息

4. **API 测试**：
   - 点击导航栏的"API文档"查看 Swagger 接口文档
   - 可以在文档中直接测试 API 接口

## API 接口

### 📚 接口文档

本项目使用 Swagger 自动生成 API 文档，启动后端服务后访问：
**http://localhost:8080/swagger/index.html**

### 🔐 认证接口

- `GET /api/auth/check-init` - 检查系统是否已初始化
- `POST /api/auth/init` - 初始化系统
- `POST /api/auth/login` - 用户登录
- `GET /api/profile` - 获取用户信息 (需要认证)

### 🖥️ 服务器管理接口

- `GET /api/servers` - 获取服务器列表 (需要认证)
- `POST /api/servers` - 创建新服务器 (需要认证)
- `GET /api/servers/{id}` - 获取服务器详情 (需要认证)
- `PUT /api/servers/{id}` - 更新服务器配置 (需要认证)
- `DELETE /api/servers/{id}` - 删除服务器 (需要认证)

### 🔑 认证说明

除了公开的认证接口外，其他接口都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <your-jwt-token>
```

### 📋 数据模型

#### 服务器配置 (ServerRequest)
```json
{
  "name": "我的ARK服务器",
  "description": "服务器描述",
  "ip": "127.0.0.1",
  "port": 7777,
  "query_port": 27015,
  "rcon_port": 27020,
  "rcon_password": "管理员密码",
  "map": "TheIsland",
  "max_players": 70
}
```

#### 支持的地图列表
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

## 环境变量

### 后端环境变量

```bash
JWT_SECRET=your-jwt-secret-key              # JWT 密钥
DB_PATH=ark_server.db                       # 数据库文件路径
SERVER_PORT=8080                            # 服务器端口
TRUSTED_PROXIES=127.0.0.1,::1              # 可信任的代理地址 (用逗号分隔)
```

### 前端环境变量

```bash
API_BASE=http://localhost:8080/api  # 后端 API 地址
```

## 🔧 开发说明

### 后端开发
- 后端使用 SQLite 作为数据库，数据文件会自动创建
- 使用 Gorm 进行数据库操作，支持自动迁移
- 所有控制器方法都包含 Swagger 注释，自动生成 API 文档
- 使用 JWT 进行身份认证，密码使用 bcrypt 加密存储
- 支持 CORS 跨域请求，配置为允许前端地址访问

### 前端开发
- 前端采用 SPA 模式，关闭了 SSR
- 使用 Pinia 进行状态管理，包含 auth 和 servers 两个 store
- 使用 Composition API 和 TypeScript 进行开发
- JWT Token 存储在 localStorage 和 Cookie 中
- 所有 API 调用使用配置的基础 URL，支持环境变量配置

### 数据安全
- 每个用户只能管理自己创建的服务器（用户级数据隔离）
- RCON 密码加密存储，API 响应中不包含敏感信息
- 前后端都有完整的输入验证和错误处理

### API 测试
- 使用 Swagger UI 进行接口测试，支持直接在浏览器中测试 API
- 生成的文档包含完整的请求/响应示例
- 支持 JWT 认证，可以在 Swagger 界面中输入 Token

### 开发工具
- 提供一键启动脚本 `start-dev.bat`
- Swagger 文档自动生成，修改代码后需要重新运行 `swag init`
- 前端支持热重载，修改代码后自动刷新

## 许可证

本项目采用  许可证。
