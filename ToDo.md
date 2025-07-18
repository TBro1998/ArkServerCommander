# 待完成功能

## 镜像更新功能

### 功能需求概述
基于现有的Docker镜像管理系统，重构镜像管理策略并新增镜像更新功能。该功能应遵循完全手动管理的原则，移除系统启动时的自动镜像下载，提供镜像检查、下载、更新和容器重启的完整流程管理。

### 完全手动管理原则
- **移除自动下载**: 取消系统启动时的自动镜像检查和下载 (`main.go:58-66` 中的 `EnsureRequiredImages`)
- **手动触发**: 所有镜像操作（检查、下载、更新）都需要用户主动触发
- **启动验证**: 服务器启动前必须验证必需镜像是否存在，不存在则禁止启动
- **状态监控**: 系统提供状态监控和进度反馈，但不自动执行任何下载或更新操作

### 镜像下载/更新操作策略

#### 镜像不存在时
- **强制手动下载**: 用户必须手动触发下载操作
- **服务器启动阻断**: 如果必需镜像不存在，禁止启动服务器，显示明确的错误提示
- **使用现有API**: 复用 `PullImageWithProgress` 方法进行下载
- **实时反馈**: 提供下载进度和状态的实时反馈

#### 镜像已存在时
- **主动检查更新**: 用户可以主动检查是否有更新版本
- **手动更新**: 如有更新，用户需要手动确认并触发更新操作
- **更新流程**: 更新过程使用 `PullImageWithProgress` 方法

### 容器更新策略
在镜像更新时，需要处理使用该镜像的容器：
- **运行中的容器**: 先停止容器 → 更新镜像 → 重新创建容器 → 重启容器
- **已停止的容器**: 直接更新镜像 → 重新创建容器（下次启动时使用新镜像）

### 技术实现规划

#### 后端API设计

1. **镜像检查更新接口** (`GET /api/servers/images/check-updates`)
   ```go
   // 检查所有镜像是否有更新
   func CheckImageUpdates(c *gin.Context) {
       // 返回: 镜像名称 -> 是否有更新的映射
   }
   ```

2. **镜像手动拉取接口** (`POST /api/servers/images/pull`)
   ```go
   // 请求体: {"image_name": "tbro98/ase-server:latest"}
   func PullImage(c *gin.Context) {
       // 1. 验证镜像名称
       // 2. 执行拉取操作（复用现有的 PullImageWithProgress）
       // 3. 返回拉取状态
   }
   ```

3. **镜像更新接口** (`POST /api/servers/images/update`)
   ```go
   // 请求体: {"image_name": "tbro98/ase-server:latest"}
   func UpdateImage(c *gin.Context) {
       // 1. 检查镜像是否存在
       // 2. 如果存在，检查是否有更新
       // 3. 执行更新拉取操作
       // 4. 处理依赖容器的重建
   }
   ```

4. **容器重建接口** (`POST /api/servers/{id}/recreate`)
   ```go
   func RecreateContainer(c *gin.Context) {
       // 1. 停止容器（如果运行中）
       // 2. 删除容器
       // 3. 重新创建容器
       // 4. 如果原来是运行状态，则重新启动
   }
   ```

#### Docker管理器扩展

基于现有的 `docker_manager` 包，需要新增以下方法：

1. **检查镜像更新** (`docker_image.go`)
   ```go
   // CheckImageUpdate 检查镜像是否有更新
   func (dm *DockerManager) CheckImageUpdate(imageName string) (bool, error)
   
   // GetImageInfo 获取本地镜像详细信息
   func (dm *DockerManager) GetImageInfo(imageName string) (*ImageInfo, error)
   ```

2. **容器重建逻辑** (`docker_manager.go`)
   ```go
   // RecreateContainer 重建容器（保持现有配置）
   func (dm *DockerManager) RecreateContainer(serverID uint) error
   
   // GetContainersByImage 获取使用指定镜像的所有容器
   func (dm *DockerManager) GetContainersByImage(imageName string) ([]ContainerInfo, error)
   ```

3. **镜像删除管理** (`docker_image.go`)
   ```go
   // RemoveOldImage 删除旧版本镜像（在更新完成后）
   func (dm *DockerManager) RemoveOldImage(imageName string, keepLatest bool) error
   ```

#### 前端UI设计

1. **镜像管理页面/模块**
   - 扩展现有的 `ImageStatus.vue` 和 `ImageStatusModal.vue`
   - 添加"检查更新"按钮
   - 添加"更新镜像"按钮（当有更新时显示）
   - 显示镜像版本信息和最后检查时间

2. **更新进度显示**
   - 复用现有的镜像下载进度显示逻辑
   - 添加容器重建进度提示
   - 显示影响的服务器列表

3. **确认对话框**
   ```vue
   // 更新前确认，显示影响的容器列表
   <ImageUpdateConfirmModal 
     :affected-servers="affectedServers"
     :image-name="imageName"
     @confirm="handleUpdateImage"
     @cancel="showConfirm = false"
   />
   ```

#### 服务器服务扩展

在 `server_service.go` 中添加镜像管理相关方法：

```go
// CheckImageUpdates 检查所有管理的镜像更新
func (s *ServerService) CheckImageUpdates() (map[string]bool, error)

// PullImage 手动拉取指定镜像
func (s *ServerService) PullImage(imageName string) error

// UpdateImage 更新指定镜像及相关容器
func (s *ServerService) UpdateImage(imageName string, userID uint) error

// GetAffectedServers 获取使用指定镜像的服务器列表
func (s *ServerService) GetAffectedServers(imageName string, userID uint) ([]ServerInfo, error)

// ValidateRequiredImages 验证启动服务器所需的镜像是否存在
func (s *ServerService) ValidateRequiredImages() (missing []string, err error)
```

### 实现优先级

1. **Phase 0**: 移除自动下载机制 *(最高优先级)*
   - 移除 `main.go` 中的 `EnsureRequiredImages` 自动调用
   - 修改服务器启动逻辑，增加镜像存在性验证
   - 更新前端提示信息，引导用户手动下载镜像

2. **Phase 1**: 手动镜像管理功能
   - 实现镜像手动拉取API (`/api/servers/images/pull`)
   - 实现镜像更新检查API (`/api/servers/images/check-updates`)
   - 扩展前端镜像管理界面，添加手动下载按钮

3. **Phase 2**: 镜像更新和容器重建功能
   - 实现容器重建逻辑
   - 处理服务器状态同步
   - 添加更新确认流程

3. **Phase 3**: 用户体验优化
   - 添加更新历史记录
   - 优化进度显示和错误处理
   - 添加批量更新功能

### 关键变更和影响

#### 移除自动下载的影响
1. **启动行为变化**: 系统启动不再自动下载镜像，`main.go:58-66` 的异步下载逻辑将被移除
2. **服务器启动验证**: `startServerAsync` 方法 (`server_service.go:571`) 需要增强镜像验证逻辑
3. **用户体验变化**: 首次使用需要用户主动下载镜像才能创建/启动服务器

#### 现有代码修改点
1. **main.go**: 移除第58-66行的自动镜像下载逻辑
2. **server_service.go**: 
   - `startServerAsync` 方法增加严格的镜像存在性检查
   - `CreateServer` 方法增加镜像验证
3. **前端组件**: 
   - `ImageStatus.vue` 和 `ImageStatusModal.vue` 需要添加手动下载功能
   - `servers.vue` 需要更新错误提示信息

### 风险和注意事项

1. **向后兼容性**: 移除自动下载可能影响现有用户的使用习惯
2. **数据安全**: 容器重建过程中确保数据卷安全
3. **服务连续性**: 运行中服务器的重启会导致临时服务中断  
4. **网络依赖**: 镜像下载和更新需要稳定的网络连接
5. **版本兼容**: 新镜像版本可能与现有配置不兼容
6. **用户权限**: 确保只有授权用户才能执行镜像操作
7. **启动失败处理**: 需要提供清晰的错误信息指导用户下载缺失的镜像

### 现有代码复用

- `docker_manager/docker_image.go`: 镜像存在检查和拉取逻辑
- `docker_manager/docker_manager.go`: 容器管理基础功能
- `service/server/server_service.go`: 服务器状态管理和容器操作
- `components/dockerimages/`: 前端镜像状态显示组件

此设计充分利用了现有的代码架构，通过扩展的方式实现镜像更新功能，保持代码的一致性和可维护性。
