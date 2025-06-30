# 开发指南

## 项目架构

本项目采用前后端分离架构：

- **前端**: Nuxt.js 3 + Pinia + Tailwind CSS
- **后端**: Gin + Gorm + SQLite + JWT

## 开发环境设置

### 1. 环境要求

- Node.js 18+ 
- Go 1.21+
- Git

### 2. 克隆项目

```bash
git clone <repository-url>
cd ArkServerManager
```

### 3. 后端开发

```bash
# 进入后端目录
cd server

# 安装Go依赖
go mod tidy

# 启动开发服务器
go run main.go
```

后端服务默认运行在 `http://localhost:8080`

#### 后端目录结构

- `config/` - 配置文件
- `controllers/` - API控制器
- `database/` - 数据库初始化
- `middleware/` - 中间件
- `models/` - 数据模型
- `routes/` - 路由定义
- `utils/` - 工具函数

#### 主要文件说明

- `main.go` - 应用入口
- `models/user.go` - 用户模型
- `controllers/auth.go` - 认证控制器
- `middleware/auth.go` - JWT认证中间件

### 4. 前端开发

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务默认运行在 `http://localhost:3000`

#### 前端目录结构

- `pages/` - 路由页面
- `layouts/` - 布局组件
- `stores/` - Pinia状态管理
- `middleware/` - 路由中间件

#### 主要文件说明

- `app.vue` - 应用根组件
- `pages/index.vue` - 首页
- `pages/login.vue` - 登录页
- `pages/init.vue` - 初始化页
- `stores/auth.js` - 认证状态管理
- `middleware/auth.js` - 认证路由中间件

## API 接口文档

### 认证接口

#### 检查初始化状态
```
GET /api/auth/check-init
Response: { "initialized": boolean }
```

#### 系统初始化
```
POST /api/auth/init
Body: { "username": string, "password": string }
Response: { "message": string, "token": string, "user": UserResponse }
```

#### 用户登录
```
POST /api/auth/login
Body: { "username": string, "password": string }
Response: { "message": string, "token": string, "user": UserResponse }
```

#### 获取用户信息
```
GET /api/profile
Headers: { "Authorization": "Bearer <token>" }
Response: { "user": UserResponse }
```

## 数据库

项目使用 SQLite 数据库，数据库文件默认为 `ark_server.db`。

### 用户表结构

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);
```

## 开发规范

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方编码规范
- 函数和变量使用驼峰命名
- 包名使用小写

### Vue/JavaScript 规范

- 使用 ESLint 进行代码检查
- 组件名使用 PascalCase
- 文件名使用 kebab-case
- 使用 Composition API

### 提交规范

使用传统提交格式：

```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码样式修改
refactor: 代码重构
test: 测试相关
chore: 其他修改
```

## 部署

### 生产环境构建

#### 后端

```bash
cd server
go build -o ark-server main.go
```

#### 前端

```bash
cd web
npm run build
npm run generate
```

### 环境变量

#### 后端环境变量

```bash
JWT_SECRET=your-production-secret-key
DB_PATH=/path/to/production.db
SERVER_PORT=8080
```

#### 前端环境变量

```bash
API_BASE=https://your-api-domain.com/api
```

## 常见问题

### 1. 跨域问题

后端已配置 CORS，允许前端访问。如果仍有问题，检查 `main.go` 中的 CORS 配置。

### 2. JWT Token 失效

Token 默认24小时过期，可以在 `utils/jwt.go` 中修改过期时间。

### 3. 数据库连接问题

确保有写入权限到数据库文件所在目录。

### 4. 端口冲突

- 后端默认端口: 8080
- 前端默认端口: 3000

可以通过环境变量修改端口。

## 扩展开发

### 添加新的 API 接口

1. 在 `models/` 中定义数据模型
2. 在 `controllers/` 中添加控制器函数
3. 在 `routes/routes.go` 中注册路由
4. 前端在 `stores/` 中添加对应的状态管理

### 添加新的页面

1. 在 `pages/` 中创建 Vue 组件
2. 如需认证，添加 `middleware: 'auth'`
3. 更新导航和路由

### 数据库迁移

项目使用 Gorm 自动迁移，修改模型后重启服务即可自动更新数据库结构。