package docker_manager

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// 并发的拉取请求只能有一个占用成功，否则同一个镜像会被同时拉取多次
func TestClaimPullIsExclusive(t *testing.T) {
	const image = "test/claim-pull:latest"

	var wins atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ClaimPull(image) {
				wins.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := wins.Load(); got != 1 {
		t.Fatalf("期望只有 1 个请求占用成功，实际 %d 个", got)
	}
	if !IsImagePulling(image) {
		t.Fatal("占用后应该处于拉取中")
	}

	// 失败后释放占用，并记录原因
	finishPull(image, errors.New("network down"))
	if IsImagePulling(image) {
		t.Fatal("结束后应该释放占用")
	}
	if imagePullError[image] != "network down" {
		t.Fatalf("失败原因没有记录: %q", imagePullError[image])
	}

	// 释放后可以再次占用，并清除上次的失败原因
	if !ClaimPull(image) {
		t.Fatal("释放后应该可以再次占用")
	}
	if _, ok := imagePullError[image]; ok {
		t.Fatal("新的拉取开始时应该清除上次的失败原因")
	}
	finishPull(image, nil)
}
