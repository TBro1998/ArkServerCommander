# ARK Server Commander - Agent 开发指南

> 本文档为 AI Agent 提供项目结构、开发规范和关键信息，用于辅助后续开发工作。

## 📋 目录

- [项目概述](#项目概述)
- [技术栈](#技术栈)
- [项目架构](#项目架构)
- [开发规范](#开发规范)
- [关键模块说明](#关键模块说明)
- [开发工作流](#开发工作流)
- [常见任务](#常见任务)

---

## 📖 项目概述

**ARK Server Commander** 是一个用于 Linux 环境的 ARK: Survival Evolved 服务器管理工具。

### 核心特性
- 🐳 基于 Docker 容器化部署，每个 ARK 服务器独立运行
- 🔌 内置 ArkApi 插件系统支持
- 🖥️ Web 界面管理多个 ARK 服务器
- ⚙️ 可视化配置服务器参数和设置
- 💾 自动管理 Docker 卷存储游戏数据

### 开发状态
⚠️ **当前处于开发阶段**，部分功能尚未完成或可能存在稳定性问题。

---

## 🛠️ 技术栈

### 后端 (Backend)
- **语言**: Go 1.23+
- **Web 框架**: Gin
- **数据库**: SQLite (GORM)
- **容器管理**: Docker SDK for Go
- **认证**: JWT (JSON Web Token)
- **API 文档**: Swagger/OpenAPI

### 前端 (Frontend)
- **框架**: Next.js 15+ (App Router)
- **语言**: TypeScript
- **UI 组件**: shadcn/ui (基于 Radix UI + Tailwind CSS)
- **状态管理**: Zustand
- **HTTP 客户端**: Axios
- **国际化**: next-intl

### 部署
- **容器化**: Docker + Docker Compose
- **镜像仓库**: Docker Hub (tbro98/arkservercommander)
- **ARK 服务器镜像**: tbro98/ase-server:latest

---

## 🏗️ 项目架构

### 目录结构

```
ArkServerCommander/
├── server/                 # Go 后端服务
│   ├── config/            # 配置管理
│   ├── controllers/       # HTTP 控制器
│   │   ├── auth/         # 认证相关
│   │   ├── images/       # Docker 镜像管理
│   │   └── servers/      # 服务器管理
│   ├── database/         # 数据库连接和初始化
│   ├── docs/             # Swagger API 文档
│   ├── middleware/       # 中间件（认证等）
│   ├── models/           # 数据模型
│   ├── routes/           # 路由定义
│   ├── service/          # 业务逻辑层
│   │   ├── docker_manager/  # Docker 容器管理
│   │   └── server/          # 服务器服务
│   ├── utils/            # 工具函数
│   ├── main.go           # 入口文件
│   ├── go.mod            # Go 依赖
│   └── go.sum
│
├── ui/                    # Next.js 前端应用
│   ├── src/
│   │   ├── app/          # App Router 页面
│   │   ├── components/   # React 组件
│   │   ├── lib/          # 工具库（axios 配置等）
│   │   └── stores/       # Zustand 状态管理
│   ├── messages/         # 国际化翻译文件
│   ├── public/           # 静态资源
│   ├── package.json
│   └── tsconfig.json
│
├── data/                  # 数据持久化目录
│   └── ark_server.db     # SQLite 数据库
│
├── docs/                  # 项目文档
│   └── zh/               # 中文文档
│
├── docker-compose.yml     # Docker Compose 配置
├── Dockerfile            # 多阶段构建配置
├── .env.example          # 环境变量示例
└── README.md             # 项目说明
```

### 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                     用户浏览器                            │
│                   (Web Interface)                        │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/HTTPS
                     ▼
┌─────────────────────────────────────────────────────────┐
│              ARK Commander Container                     │
│  ┌──────────────────────┐  ┌──────────────────────┐    │
│  │   Next.js Frontend   │  │    Go Backend API    │    │
│  │   (Port 3000)        │◄─┤    (Port 8080)       │    │
│  └──────────────────────┘  └──────────┬───────────┘    │
│                                        │                 │
│                            ┌───────────▼───────────┐    │
│                            │   SQLite Database     │    │
│                            │  (ark_server.db)      │    │
│                            └───────────────────────┘    │
└────────────────────────────┬───────────────────────────┘
                             │ Docker Socket
                             ▼
┌─────────────────────────────────────────────────────────┐
│                    Docker Engine                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ ARK Server 1 │  │ ARK Server 2 │  │ ARK Server N │  │
│  │ (Container)  │  │ (Container)  │  │ (Container)  │  │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  │
│         │                 │                 │           │
│  ┌──────▼───────┐  ┌──────▼───────┐  ┌──────▼───────┐  │
│  │ Docker Volume│  │ Docker Volume│  │ Docker Volume│  │
│  │ (Game Data)  │  │ (Game Data)  │  │ (Game Data)  │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 数据流

1. **用户操作** → Next.js UI
2. **UI 发起请求** → Axios → Go Backend API
3. **API 处理**:
   - 认证中间件验证 JWT
   - 控制器接收请求
   - 服务层执行业务逻辑
   - 数据库层持久化数据
   - Docker SDK 管理容器
4. **响应返回** → UI 更新状态

---

## 📐 开发规范

### 总体原则

- **清晰性**: 代码应该易于理解，优先考虑可读性
- **一致性**: 遵循项目中已有的编码风格和模式
- **单一职责**: 每个函数、模块或组件都应该只做一件事

---

### 后端开发规范 (Go)

#### 代码风格
- 遵循官方 Go 编码规范 (`gofmt`, `golint`)
- 命名清晰简洁，能反映其功能
- 使用 `go mod` 进行依赖管理

#### 项目结构
- `/controllers`: HTTP 请求处理控制器
- `/models`: 数据结构定义
- `/routes`: API 路由和中间件
- `/service`: 核心业务逻辑
- `/database`: 数据库连接和操作
- `/middleware`: 认证等中间件
- `/utils`: 辅助工具函数

#### API 设计
- 遵循 RESTful 原则
- 使用 JSON 作为数据交换格式
- 所有 API 路由以 `/api` 为前缀
- 目前未使用版本控制（如 `/api/v1`）

#### 错误处理
- 错误应在函数调用链中显式向上传递
- 不要在底层函数中记录错误日志然后向上传递
- 日志记录应在顶层处理程序中进行
- 使用自定义错误类型传递更多上下文信息

#### 测试
- 所有核心功能必须有单元测试
- 测试文件与源文件放在同一包中，以 `_test.go` 结尾

---

### 前端开发规范 (Next.js)

#### 代码风格
- 使用 Prettier 和 ESLint 保证代码风格一致
- 使用 TypeScript 进行类型检查
- 遵循 React 社区最佳实践

#### 组件化
- 将 UI 拆分为可重用的小组件，存放在 `/src/components`
- 使用 `shadcn/ui` (基于 Radix UI 和 Tailwind CSS)
- 自定义组件应遵循与 `shadcn/ui` 类似的风格和结构

#### 状态管理
- 全局状态管理使用 `Zustand`
- Store 按功能模块划分（如 `useAuthStore` 在 `/src/stores/auth.ts`）
- Store 包含状态 (`state`) 和操作方法 (`actions`)

#### API 请求
- 使用 `axios` 进行网络请求
- 全局 `axios` 实例在 `/src/lib/axios.ts` 中配置
- API 调用逻辑封装在 `Zustand` store 的 `actions` 中
- 在 `actions` 中统一处理请求的加载（`isLoading`）和错误状态

#### 国际化
- 使用 `next-intl` 进行国际化
- 所有面向用户的字符串必须通过翻译函数处理
- 语言文件存放在 `/messages` 目录（如 `en.ts`, `zh.ts`）

---

### Git 工作流

#### 分支模型
- `main`: 稳定的主分支，用于生产发布
- `develop`: 开发分支，集成所有已完成的功能
- `feature/xxx`: 功能开发分支，从 `develop` 创建
- `fix/xxx`: Bug 修复分支

#### 提交信息
遵循 Conventional Commits 规范:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**type 类型**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

**scope**: 可选，表示修改范围（如 `api`, `ui`, `db`）

**示例**:
```
feat(server): 添加服务器状态监控功能
fix(ui): 修复服务器列表刷新问题
docs: 更新 README 安装说明
```

---

## 🔑 关键模块说明

### 后端核心模块

#### 1. 认证系统 (`/server/controllers/auth`, `/server/middleware/auth.go`)

**功能**:
- 用户注册、登录
- JWT Token 生成和验证
- 密码加密存储（bcrypt）

**关键文件**:
- `controllers/auth/auth.go`: 认证控制器
- `middleware/auth.go`: JWT 验证中间件
- `utils/jwt.go`: JWT 工具函数
- `utils/password.go`: 密码加密工具

**安全要求**:
- JWT_SECRET 必须配置且长度 ≥ 32 字符
- 应用启动时会验证 JWT_SECRET 强度
- 密码使用 bcrypt 加密存储

#### 2. 服务器管理 (`/server/controllers/servers`, `/server/service/server`)

**功能**:
- 创建、启动、停止、删除 ARK 服务器
- 配置服务器参数（地图、端口、密码等）
- 管理 GameUserSettings.ini 和 Game.ini 配置文件

**关键文件**:
- `controllers/servers/server.go`: 服务器控制器
- `service/server/server_service.go`: 服务器业务逻辑
- `models/server.go`: 服务器数据模型
- `models/server_args.go`: 服务器启动参数模型

**支持的地图**:
- TheIsland, TheCenter, ScorchedEarth_P
- Aberration_P, Extinction, Valguero_P
- Genesis, CrystalIsles, Gen2, LostIsland, Fjordur

#### 3. Docker 容器管理 (`/server/service/docker_manager`)

**功能**:
- 创建和管理 ARK 服务器 Docker 容器
- 管理 Docker 卷（游戏数据存储）
- 拉取和管理 Docker 镜像

**关键文件**:
- `service/docker_manager/docker_image.go`: 镜像管理
- `service/docker_manager/docker_container.go`: 容器管理（推测）

**使用的镜像**:
- ARK 服务器: `tbro98/ase-server:latest`
- 管理系统: `tbro98/arkservercommander:latest`

#### 4. 配置文件管理 (`/server/utils/config_files.go`)

**功能**:
- 生成和更新 GameUserSettings.ini
- 生成和更新 Game.ini
- 处理 ARK 服务器配置参数

**配置文件位置**:
- 容器内: `/ark/ShooterGame/Saved/Config/LinuxServer/`
- 通过 Docker 卷持久化

#### 5. 数据库层 (`/server/database`, `/server/models`)

**数据库**: SQLite (使用 GORM)

**主要模型**:
- `User`: 用户账户
- `Server`: ARK 服务器配置
- `ServerArgs`: 服务器启动参数

**数据库文件**: `/data/ark_server.db`

---

### 前端核心模块

#### 1. 认证状态管理 (`/ui/src/stores/auth.ts`)

**功能**:
- 用户登录/登出
- Token 管理
- 用户信息获取
- 认证状态持久化

**关键方法**:
- `login(username, password)`: 用户登录
- `logout()`: 用户登出
- `getProfile()`: 获取用户信息
- `checkAuth()`: 检查认证状态

#### 2. 服务器管理界面 (`/ui/src/app/[locale]/servers`)

**功能**:
- 服务器列表展示
- 添加/编辑服务器
- 启动/停止服务器
- 配置服务器参数

**主要组件**:
- 服务器列表组件
- 服务器表单组件
- 配置编辑器组件

#### 3. Axios 配置 (`/ui/src/lib/axios.ts`)

**功能**:
- 全局 axios 实例配置
- 请求/响应拦截器
- 自动添加 JWT Token
- 统一错误处理

**配置**:
- Base URL: `/api` (相对路径)
- 自动从 localStorage 读取 token
- 401 错误自动跳转登录页

#### 4. 国际化 (`/ui/messages`)

**支持语言**:
- 英文 (`en.ts`)
- 中文 (`zh.ts`)

**使用方式**:
```typescript
import { useTranslations } from 'next-intl';

const t = useTranslations('ServerManagement');
const title = t('title'); // 获取翻译文本
```

---

## 🔄 开发工作流

### 本地开发环境搭建

#### 后端开发

```bash
# 1. 进入后端目录
cd server/

# 2. 安装依赖
go mod download

# 3. 配置环境变量
export JWT_SECRET='your-strong-secret-key-here-min-32-chars'
export DB_PATH='./data/ark_server.db'
export SERVER_PORT='8080'

# 4. 运行后端服务
go run main.go
```

#### 前端开发

```bash
# 1. 进入前端目录
cd ui/

# 2. 安装依赖
pnpm install

# 3. 配置环境变量
# 编辑 .env.development
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.development

# 4. 运行开发服务器
pnpm dev
```

---

### Docker 开发环境

```bash
# 1. 构建镜像
docker-compose build

# 2. 启动服务
docker-compose up -d

# 3. 查看日志
docker-compose logs -f

# 4. 停止服务
docker-compose down
```

---

### 代码检查和格式化

#### 后端 (Go)

```bash
# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 运行测试
go test ./...
```

#### 前端 (TypeScript/Next.js)

```bash
# 代码检查
pnpm lint

# 格式化代码
pnpm format

# 类型检查
pnpm type-check
```

---

### 构建和部署

#### 构建 Docker 镜像

```bash
# 构建镜像
docker build -t tbro98/arkservercommander:latest .

# 推送到 Docker Hub
docker push tbro98/arkservercommander:latest
```

#### 多阶段构建说明

Dockerfile 使用多阶段构建:
1. **Stage 1**: 构建 Go 后端
2. **Stage 2**: 构建 Next.js 前端
3. **Stage 3**: 组合运行时镜像

---

## 📝 常见任务

### 添加新的 API 端点

#### 1. 定义数据模型 (`/server/models`)

```go
// models/example.go
package models

type Example struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    Name      string `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 2. 创建控制器 (`/server/controllers`)

```go
// controllers/example/example.go
package example

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func GetExample(c *gin.Context) {
    // 业务逻辑
    c.JSON(http.StatusOK, gin.H{
        "message": "success",
    })
}
```

#### 3. 注册路由 (`/server/routes/routes.go`)

```go
// 在 SetupRoutes 函数中添加
api := r.Group("/api")
{
    example := api.Group("/example")
    {
        example.GET("/", exampleController.GetExample)
    }
}
```

---

### 添加新的前端页面

#### 1. 创建页面组件 (`/ui/src/app/[locale]`)

```typescript
// app/[locale]/example/page.tsx
export default function ExamplePage() {
  return (
    <div>
      <h1>Example Page</h1>
    </div>
  );
}
```

#### 2. 创建 Zustand Store (`/ui/src/stores`)

```typescript
// stores/example.ts
import { create } from 'zustand';

interface ExampleState {
  data: any[];
  isLoading: boolean;
  fetchData: () => Promise<void>;
}

export const useExampleStore = create<ExampleState>((set) => ({
  data: [],
  isLoading: false,
  fetchData: async () => {
    set({ isLoading: true });
    try {
      const response = await axios.get('/api/example');
      set({ data: response.data });
    } finally {
      set({ isLoading: false });
    }
  },
}));
```

#### 3. 添加国际化文本 (`/ui/messages`)

```typescript
// messages/zh.ts
export default {
  Example: {
    title: '示例页面',
    description: '这是一个示例',
  },
};
```

---

### 数据库迁移

#### 添加新表或字段

```go
// 在 database/database.go 的 InitDB 函数中
func InitDB() {
    // 自动迁移
    db.AutoMigrate(
        &models.User{},
        &models.Server{},
        &models.NewModel{}, // 添加新模型
    )
}
```

GORM 会自动处理表结构变更。

---

### Docker 容器管理

#### 创建 ARK 服务器容器

```go
// 使用 Docker SDK
containerConfig := &container.Config{
    Image: "tbro98/ase-server:latest",
    Env: []string{
        "MAP=TheIsland",
        "SERVER_NAME=My ARK Server",
    },
}

hostConfig := &container.HostConfig{
    PortBindings: nat.PortMap{
        "7777/udp": []nat.PortBinding{{HostPort: "7777"}},
    },
}
```

---

### 调试技巧

#### 后端调试

```bash
# 查看详细日志
go run main.go 2>&1 | tee server.log

# 使用 delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go
```

#### 前端调试

```bash
# 开启详细日志
NEXT_PUBLIC_DEBUG=true pnpm dev

# 查看构建信息
pnpm build --debug
```

#### Docker 调试

```bash
# 查看容器日志
docker logs ark-commander -f

# 进入容器
docker exec -it ark-commander /bin/sh

# 查看 ARK 服务器容器
docker ps | grep ark-server

# 查看 Docker 卷
docker volume ls | grep ark-server
```

---

## ⚠️ 重要注意事项

### 安全相关

1. **JWT Secret 配置**
   - ⚠️ 必须配置强密钥（≥32字符）
   - ❌ 不要使用默认值
   - ❌ 不要提交到版本控制
   - ✅ 使用 `openssl rand -base64 48` 生成

2. **Docker Socket 权限**
   - 容器需要 `privileged: true` 来管理宿主机 Docker
   - 这是必需的，但要注意安全风险
   - 不要在不受信任的环境中运行

3. **数据备份**
   - 定期备份 SQLite 数据库 (`/data/ark_server.db`)
   - 备份 Docker 卷中的游戏数据
   - 使用 `docker volume backup` 或手动复制

### 性能相关

1. **资源要求**
   - 每个 ARK 服务器需要 8GB+ 内存
   - 每个服务器需要 10GB+ 磁盘空间
   - 建议使用 SSD 存储

2. **容器限制**
   - 考虑为 ARK 服务器容器设置资源限制
   - 避免单台主机运行过多服务器

### 开发相关

1. **API 变更**
   - 修改 API 后更新 Swagger 文档
   - 同步更新前端 API 调用
   - 考虑向后兼容性

2. **数据库变更**
   - GORM 自动迁移可能不够用
   - 复杂变更需要手动迁移脚本
   - 测试前先备份数据库

3. **Docker 镜像**
   - 修改后重新构建镜像
   - 测试多阶段构建是否正常
   - 注意镜像大小优化

---

## 🚧 待开发功能清单

### 高优先级

- [ ] 🖼️ 服务器镜像更新功能
- [ ] 🎮 RCON 命令执行
- [ ] 📊 服务器运行状态监控
- [ ] 📋 服务器日志查看

### 中优先级

- [ ] 🎨 Mod 管理（Steam Workshop 集成）
- [ ] 🔧 ArkApi 插件管理
- [ ] 💾 服务器存档和配置备份
- [ ] 🔍 工具版本更新检查
- [ ] ⚡ 可选的服务器文件和 Mod 更新

### 低优先级

- [ ] 🌐 完善 i18n 国际化支持
- [ ] 🔌 MCP (Mod Configuration Protocol) 支持

### 未来计划

- [ ] ☸️ 基于 K8S 的多主机管理
- [ ] 🌍 服务器列表网站（改善 Steam 服务器搜索体验）
- [ ] 👥 玩家用户界面

---

## 📚 快速参考

### 常用命令

```bash
# 后端开发
cd server && go run main.go

# 前端开发
cd ui && pnpm dev

# Docker 部署
docker-compose up -d

# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 停止服务
docker-compose down
```

### 常用路径

```bash
# 数据库
/data/ark_server.db

# Docker Socket
/var/run/docker.sock

# ARK 服务器配置
/ark/ShooterGame/Saved/Config/LinuxServer/

# Docker 卷
docker volume ls | grep ark-server
```

### 环境变量

```bash
# 必需
JWT_SECRET=your-secret-key-here

# 可选
DB_PATH=/data/ark_server.db
SERVER_PORT=8080
```

### API 端点示例

```bash
# 认证
POST /api/auth/login
POST /api/auth/register
GET  /api/auth/profile

# 服务器管理
GET    /api/servers
POST   /api/servers
GET    /api/servers/:id
PUT    /api/servers/:id
DELETE /api/servers/:id
POST   /api/servers/:id/start
POST   /api/servers/:id/stop

# Docker 镜像
GET  /api/images
POST /api/images/pull
```

---

## 🔗 相关资源

### 项目相关

- **GitHub 仓库**: [ArkServerCommander](https://github.com/tbro199803/ArkServerCommander)
- **Docker Hub**: [tbro98/arkservercommander](https://hub.docker.com/r/tbro98/arkservercommander)
- **ARK 服务器镜像**: [tbro98/ase-server](https://hub.docker.com/r/tbro98/ase-server)
- **ASE Server Docker**: [GitHub](https://github.com/tbro199803/ASE-Server-Docker)

### 技术文档

#### Go 相关
- [Go 官方文档](https://go.dev/doc/)
- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM](https://gorm.io/docs/)
- [Docker SDK for Go](https://docs.docker.com/engine/api/sdk/)

#### Next.js 相关
- [Next.js 文档](https://nextjs.org/docs)
- [shadcn/ui](https://ui.shadcn.com/)
- [Zustand](https://docs.pmnd.rs/zustand/getting-started/introduction)
- [next-intl](https://next-intl-docs.vercel.app/)
- [Tailwind CSS](https://tailwindcss.com/docs)

#### Docker 相关
- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [1Panel 文档](https://1panel.cn/docs/)

### ARK 相关

- [ARK: Survival Evolved Wiki](https://ark.fandom.com/wiki/ARK_Survival_Evolved_Wiki)
- [ArkApi 文档](https://arkapi.net/)
- [ARK Server Manager](https://arkservermanager.freeforums.net/)

---

## 📖 文档维护

### 更新此文档

当项目发生以下变化时，请更新此文档：

1. **架构变更**: 添加新的服务、模块或重大重构
2. **技术栈变更**: 更换框架、库或工具
3. **开发规范变更**: 修改编码规范或工作流程
4. **新功能完成**: 从待开发清单移到已实现功能
5. **重要决策**: 记录架构决策和技术选型理由

### 文档版本

- **创建日期**: 2026-02-01
- **最后更新**: 2026-02-01
- **维护者**: AI Agent (呜呼)

---

## 💡 开发建议

### 给 AI Agent 的提示

当你需要开发新功能时：

1. **先阅读相关模块**: 查看现有代码的实现方式
2. **遵循现有模式**: 保持代码风格一致
3. **测试先行**: 编写测试用例，确保功能正常
4. **文档同步**: 更新 API 文档和用户文档
5. **安全第一**: 注意认证、授权和数据验证

### 常见问题排查

**问题**: 容器无法启动
- 检查 JWT_SECRET 是否配置
- 检查 Docker Socket 权限
- 查看容器日志 `docker logs ark-commander`

**问题**: 前端无法连接后端
- 检查 NEXT_PUBLIC_API_URL 配置
- 检查后端是否正常运行
- 检查网络和端口映射

**问题**: ARK 服务器无法启动
- 检查内存是否充足（需要 8GB+）
- 检查磁盘空间（需要 10GB+）
- 查看 ARK 服务器容器日志

---

## 🎯 总结

这个项目是一个**全栈 Web 应用**，用于管理 ARK: Survival Evolved 游戏服务器：

- **后端**: Go + Gin + GORM + Docker SDK
- **前端**: Next.js + TypeScript + shadcn/ui + Zustand
- **部署**: Docker + Docker Compose
- **架构**: 容器化微服务架构

**核心功能**: 通过 Web 界面管理多个 ARK 服务器的创建、配置、启动和停止。

**开发重点**: 
- 保持代码清晰和一致性
- 注重安全性（JWT、密码加密）
- 优化用户体验（国际化、响应式设计）
- 完善待开发功能清单

---

**祝开发顺利！** 🚀

如有问题或需要更新此文档，请随时修改。

