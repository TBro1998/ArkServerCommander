# 结构化日志使用说明

## 环境变量配置

### LOG_LEVEL（日志级别）
- `debug` - 调试信息（最详细）
- `info` - 一般信息（默认）
- `warn` - 警告信息
- `error` - 错误信息
- `fatal` - 致命错误（会导致程序退出）

### LOG_FORMAT（日志格式）
- `json` - JSON格式（生产环境推荐）
  - 便于ELK、Loki等日志收集系统解析
  - 结构化，易于查询和分析
- `console` - 控制台格式（开发环境推荐）
  - 彩色输出，人类可读
  - 更直观的调试体验

## 代码中使用

### 导入
```go
import (
    "ark-server-commander/utils"
    "go.uber.org/zap"
)
```

### 基本使用（结构化字段）

```go
// Info级别
utils.Info("用户登录成功", 
    zap.String("username", "admin"),
    zap.Uint("user_id", 123),
)

// Error级别
utils.Error("数据库连接失败", 
    zap.Error(err),
    zap.String("db_path", dbPath),
)

// Debug级别（只在LOG_LEVEL=debug时输出）
utils.Debug("处理请求", 
    zap.String("method", "POST"),
    zap.String("path", "/api/servers"),
    zap.Duration("latency", time.Since(start)),
)
```

### Printf风格（简便用法）

```go
// 适合简单日志
utils.Infof("服务器启动在端口 %d", port)
utils.Errorf("创建容器失败: %v", err)
utils.Debugf("收到请求: %s %s", method, path)
```

### 常用字段类型

```go
zap.String("key", "value")      // 字符串
zap.Int("key", 123)             // 整数
zap.Uint("key", 123)            // 无符号整数
zap.Bool("key", true)           // 布尔值
zap.Error(err)                  // 错误对象
zap.Duration("key", duration)   // 时间间隔
zap.Time("key", time.Now())     // 时间戳
zap.Any("key", anyValue)        // 任意类型（会序列化）
```

## 输出示例

### JSON格式（生产环境）
```json
{
  "timestamp": "2026-01-29T15:50:00.123Z",
  "level": "info",
  "caller": "main.go:45",
  "message": "用户登录成功",
  "username": "admin",
  "user_id": 123
}
```

### Console格式（开发环境）
```
2026-01-29T15:50:00.123Z	INFO	main.go:45	用户登录成功	{"username": "admin", "user_id": 123}
```

## 日志级别使用建议

- **Debug**: 详细的调试信息，只在开发时启用
  - 函数入口/出口
  - 中间变量值
  - 详细的请求/响应内容

- **Info**: 正常的业务流程记录
  - 服务启动/关闭
  - 用户登录/登出
  - 重要操作完成

- **Warn**: 潜在问题，但不影响运行
  - 配置使用默认值
  - 重试操作
  - 性能降级

- **Error**: 错误情况，需要关注
  - API调用失败
  - 数据库操作失败
  - 业务逻辑错误

- **Fatal**: 致命错误，程序无法继续
  - 配置严重错误
  - 关键依赖不可用
  - 系统资源耗尽

## 最佳实践

1. **使用结构化字段而非字符串拼接**
   ```go
   // ❌ 不推荐
   utils.Infof("用户 %s (ID: %d) 登录成功", username, userID)
   
   // ✅ 推荐
   utils.Info("用户登录成功", 
       zap.String("username", username),
       zap.Uint("user_id", userID),
   )
   ```

2. **错误日志要包含足够的上下文**
   ```go
   utils.Error("创建服务器失败",
       zap.Error(err),
       zap.String("server_name", name),
       zap.Int("port", port),
       zap.String("map", mapName),
   )
   ```

3. **敏感信息要脱敏**
   ```go
   // ❌ 不要直接记录密码
   utils.Debug("收到登录请求", zap.String("password", password))
   
   // ✅ 只记录长度或完全不记录
   utils.Debug("收到登录请求", zap.Int("password_length", len(password)))
   ```

4. **性能敏感的代码路径避免过多日志**
   ```go
   // 高频调用的地方，只在DEBUG级别记录
   utils.Debug("处理请求", zap.String("path", path))
   ```
