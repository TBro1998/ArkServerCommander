package docker_manager

import (
	"ark-server-commander/utils"
	"fmt"
	"go.uber.org/zap"
)

// RollbackAction 回滚操作类型
type RollbackAction struct {
	Type        string       // 操作类型: "volume", "container", "config"
	ResourceID  string       // 资源ID
	Action      func() error // 回滚函数
	Description string       // 操作描述
}

// RollbackManager 回滚管理器
type RollbackManager struct {
	actions []RollbackAction
}

// NewRollbackManager 创建新的回滚管理器
func NewRollbackManager() *RollbackManager {
	return &RollbackManager{
		actions: make([]RollbackAction, 0),
	}
}

// AddAction 添加回滚操作
func (rm *RollbackManager) AddAction(actionType, resourceID, description string, action func() error) {
	rm.actions = append(rm.actions, RollbackAction{
		Type:        actionType,
		ResourceID:  resourceID,
		Action:      action,
		Description: description,
	})
	utils.Debug("添加回滚操作",
		zap.String("type", actionType),
		zap.String("resource", resourceID),
		zap.String("description", description))
}

// Rollback 执行所有回滚操作（逆序执行）
func (rm *RollbackManager) Rollback() error {
	if len(rm.actions) == 0 {
		utils.Debug("没有需要回滚的操作")
		return nil
	}

	utils.Info("开始执行回滚操作", zap.Int("count", len(rm.actions)))

	var rollbackErrors []error

	// 逆序执行回滚操作
	for i := len(rm.actions) - 1; i >= 0; i-- {
		action := rm.actions[i]
		utils.Info("执行回滚",
			zap.String("type", action.Type),
			zap.String("resource", action.ResourceID),
			zap.String("description", action.Description))

		if err := action.Action(); err != nil {
			utils.Error("回滚操作失败",
				zap.String("type", action.Type),
				zap.String("resource", action.ResourceID),
				zap.Error(err))
			rollbackErrors = append(rollbackErrors, fmt.Errorf("%s 回滚失败: %w", action.Description, err))
		} else {
			utils.Info("回滚操作成功",
				zap.String("type", action.Type),
				zap.String("resource", action.ResourceID))
		}
	}

	if len(rollbackErrors) > 0 {
		return fmt.Errorf("部分回滚操作失败: %v", rollbackErrors)
	}

	utils.Info("所有回滚操作执行完成")
	return nil
}

// Clear 清空回滚操作列表
func (rm *RollbackManager) Clear() {
	rm.actions = make([]RollbackAction, 0)
	utils.Debug("清空回滚操作列表")
}

// Count 返回待回滚操作数量
func (rm *RollbackManager) Count() int {
	return len(rm.actions)
}
