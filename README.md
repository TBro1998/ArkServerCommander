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

- ✅ 前后端分离架构
- ✅ JWT 身份认证
- ✅ 系统初始化设置
- ✅ 用户登录/登出
- ✅ 响应式设计 (Tailwind CSS)
- ✅ SQLite 数据库
- ✅ 单用户系统 (无权限管理)

## 技术栈

### 后端
- **Gin**: Web 框架
- **Gorm**: ORM 数据库操作
- **SQLite**: 数据库
- **JWT**: 身份认证
- **bcrypt**: 密码加密

### 前端
- **Nuxt.js 3**: Vue.js 框架
- **Pinia**: 状态管理
- **Tailwind CSS**: 样式框架
- **TypeScript**: 类型支持

## 快速开始

### 后端启动

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

### 前端启动

1. 进入前端目录：
```bash
cd web
```

2. 安装依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm run dev
```

前端服务将在 `http://localhost:3000` 启动

## 使用说明

1. **初次使用**：
   - 访问 `http://localhost:3000`
   - 系统会自动跳转到初始化页面
   - 设置管理员用户名和密码

2. **登录系统**：
   - 使用刚才设置的用户名和密码登录
   - 登录成功后进入管理界面

3. **系统管理**：
   - 当前版本提供基础的用户认证功能
   - 更多 ARK 服务器管理功能待开发

## API 接口

### 认证接口

- `GET /api/auth/check-init` - 检查系统是否已初始化
- `POST /api/auth/init` - 初始化系统
- `POST /api/auth/login` - 用户登录
- `GET /api/profile` - 获取用户信息 (需要认证)

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

## 开发说明

- 后端使用 SQLite 作为数据库，数据文件会自动创建
- 前端采用 SPA 模式，关闭了 SSR
- 使用 JWT 进行身份认证，Token 存储在 localStorage
- 系统只支持单用户，无复杂权限管理

## 许可证

本项目采用  许可证。
