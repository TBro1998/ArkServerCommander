# Kilo Code - 项目开发规范

本文档定义了本项目的开发规范、编码风格和最佳实践。所有贡献者都应遵循这些规则，以确保代码质量、可读性和可维护性。

## 目录

- [Kilo Code - 项目开发规范](#kilo-code---项目开发规范)
  - [目录](#目录)
  - [总体原则](#总体原则)
  - [后端 (Go)](#后端-go)
    - [代码风格](#代码风格)
    - [项目结构](#项目结构)
    - [API 设计](#api-设计)
    - [错误处理](#错误处理)
    - [测试](#测试)
  - [前端 (Next.js)](#前端-nextjs)
    - [代码风格](#代码风格-1)
    - [组件化](#组件化)
    - [状态管理](#状态管理)
    - [API 请求](#api-请求)
    - [国际化](#国际化)
  - [Git 工作流](#git-工作流)
    - [分支模型](#分支模型)
    - [提交信息](#提交信息)
  - [文档](#文档)

## 总体原则

- **清晰性**: 代码应该易于理解。优先考虑可读性，而不是过度优化。
- **一致性**: 遵循项目中已有的编码风格和模式。
- **单一职责**: 每个函数、模块或组件都应该只做一件事。

## 后端 (Go)

### 代码风格

- 遵循官方的 Go 编码规范 (`gofmt`, `golint`)。
- 命名应清晰、简洁，并能反映其功能。
- 使用 `go mod` 进行依赖管理。

### 项目结构

后端服务遵循以下结构，旨在实现关注点分离：

- `/controllers`: 包含处理 HTTP 请求的控制器。
- `/models`: 定义数据结构。
- `/routes`: 定义 API 路由和中间件。
- `/service`: 包含核心业务逻辑。
- `/database`: 处理数据库连接和操作。
- `/middleware`: 存放认证等中间件。
- `/utils`: 存放辅助工具函数。

### API 设计

- API 应遵循 RESTful 原则。
- 使用 JSON 作为主要的数据交换格式。
- 所有 API 路由都以 `/api` 为前缀。目前未在 URL 中使用版本控制 (例如, `/api/v1`)。

### 错误处理

- 错误应在函数调用链中显式向上传递。
- 不要在底层函数中记录错误日志然后向上传递。日志记录应在顶层处理程序中进行。
- 使用自定义错误类型来传递更多上下文信息。

### 测试

- 所有核心功能都必须有单元测试。
- 测试文件应与源文件放在同一个包中，并以 `_test.go` 结尾。

## 前端 (Next.js)

### 代码风格

- 使用 Prettier 和 ESLint 来保证代码风格一致。
- 使用 TypeScript 进行类型检查。
- 遵循 React 社区的最佳实践。

### 组件化

- 将 UI 拆分为可重用的小组件，存放在 `/src/components` 目录下。
- 项目使用 `shadcn/ui`，它基于 Radix UI 和 Tailwind CSS，提供了一套可定制、可重用的 UI 组件。
- 自定义组件应遵循与 `shadcn/ui` 类似的风格和结构。

### 状态管理

- 全局状态管理使用 `Zustand`。
- Store 应按功能模块进行划分，例如 `useAuthStore` (`/src/stores/auth.ts`) 专门处理认证逻辑。
- Store 中应包含状态 (`state`) 和操作状态的方法 (`actions`)。

### API 请求

- 使用 `axios` 进行网络请求。全局的 `axios` 实例在 `/src/lib/axios.ts` 中配置。
- API 调用逻辑封装在 `Zustand` store 的 `actions` 中 (例如 `useAuthStore` 中的 `login`, `getProfile` 等)。
- 在 `actions` 中统一处理请求的加载（`isLoading`）和错误状态。

### 国际化

- 使用 `next-intl` 进行国际化。
- 所有面向用户的字符串都必须通过翻译函数处理。
- 语言文件存放在 `/messages` 目录中，例如 `en.ts` 和 `zh.ts`。

## Git 工作流

### 分支模型

- `main`: 稳定的主分支，用于生产发布。
- `develop`: 开发分支，集成所有已完成的功能。
- `feature/xxx`: 功能开发分支，从 `develop` 创建。
- `fix/xxx`: Bug 修复分支。

### 提交信息

遵循 Conventional Commits 规范:

```
<type>(<scope>): <subject>

<body>

<footer>
```

- **type**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- **scope**: 可选，表示修改的范围 (例如, `api`, `ui`, `db`)
- **subject**: 简要描述。

## 文档

- API 文档应使用 Swagger/OpenAPI 自动生成。
- 复杂的功能或架构决策应在 `/docs` 目录中进行记录。
- README 文件应保持最新，包含项目设置和启动说明。