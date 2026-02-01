package docker_manager

import (
	"testing"
)

// TestRollbackManager 测试回滚管理器
func TestRollbackManager(t *testing.T) {
	rm := NewRollbackManager()

	// 测试添加操作
	executed := false
	rm.AddAction("test", "resource1", "测试操作", func() error {
		executed = true
		return nil
	})

	if rm.Count() != 1 {
		t.Errorf("期望操作数量为1，实际为%d", rm.Count())
	}

	// 测试回滚
	if err := rm.Rollback(); err != nil {
		t.Errorf("回滚失败: %v", err)
	}

	if !executed {
		t.Error("回滚操作未执行")
	}
}

// TestRollbackManagerMultipleActions 测试多个回滚操作
func TestRollbackManagerMultipleActions(t *testing.T) {
	rm := NewRollbackManager()

	var executionOrder []int

	// 添加多个操作
	rm.AddAction("test", "resource1", "操作1", func() error {
		executionOrder = append(executionOrder, 1)
		return nil
	})

	rm.AddAction("test", "resource2", "操作2", func() error {
		executionOrder = append(executionOrder, 2)
		return nil
	})

	rm.AddAction("test", "resource3", "操作3", func() error {
		executionOrder = append(executionOrder, 3)
		return nil
	})

	// 执行回滚
	if err := rm.Rollback(); err != nil {
		t.Errorf("回滚失败: %v", err)
	}

	// 验证执行顺序（应该是逆序）
	if len(executionOrder) != 3 {
		t.Errorf("期望执行3个操作，实际执行%d个", len(executionOrder))
	}

	if executionOrder[0] != 3 || executionOrder[1] != 2 || executionOrder[2] != 1 {
		t.Errorf("回滚顺序错误，期望[3,2,1]，实际%v", executionOrder)
	}
}

// TestRollbackManagerClear 测试清空操作
func TestRollbackManagerClear(t *testing.T) {
	rm := NewRollbackManager()

	rm.AddAction("test", "resource1", "操作1", func() error {
		return nil
	})

	if rm.Count() != 1 {
		t.Errorf("期望操作数量为1，实际为%d", rm.Count())
	}

	rm.Clear()

	if rm.Count() != 0 {
		t.Errorf("清空后期望操作数量为0，实际为%d", rm.Count())
	}
}
