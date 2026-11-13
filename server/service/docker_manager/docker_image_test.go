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

func pullMsg(status, id string, current, total int64) *ImagePullProgress {
	msg := &ImagePullProgress{Status: status, ID: id}
	if total > 0 || current > 0 {
		msg.ProgressDetail = &struct {
			Current int64 `json:"current"`
			Total   int64 `json:"total"`
		}{current, total}
	}
	return msg
}

// 按 Docker 拉取流的真实顺序记录层状态："Pulling from" 不算一层，
// "Download complete" 不能当成已完成，本地已有的层单独标记
func TestApplyPullProgress(t *testing.T) {
	layers := make(map[string]*LayerStatus)
	steps := []*ImagePullProgress{
		pullMsg("Pulling from tbro98/ase-server", "latest", 0, 0),
		pullMsg("Already exists", "aaa", 0, 0),
		pullMsg("Pulling fs layer", "bbb", 0, 0),
		pullMsg("Pulling fs layer", "ccc", 0, 0),
		pullMsg("Waiting", "ccc", 0, 0),
		pullMsg("Downloading", "bbb", 300, 1000),
		pullMsg("Verifying Checksum", "bbb", 0, 0),
		pullMsg("Download complete", "bbb", 0, 0),
	}
	for _, msg := range steps {
		applyPullProgress(layers, msg)
	}

	if _, ok := layers["latest"]; ok {
		t.Fatal("Pulling from 的 ID 是镜像标签，不应记成一层")
	}
	if got := layers["aaa"].Status; got != "exists" {
		t.Fatalf("Already exists 的层状态应为 exists，实际 %s", got)
	}
	if b := layers["bbb"]; b.Status != "downloaded" || b.Size != 1000 || b.Progress != 1000 {
		t.Fatalf("下载完成的层状态不对: %+v", *b)
	}
	if c := layers["ccc"]; c.Status != "pending" || c.Size != 0 {
		t.Fatalf("排队中的层应为 pending 且大小未知: %+v", *c)
	}

	// 解压阶段的进度单独计算，大小保持下载阶段的值
	applyPullProgress(layers, pullMsg("Extracting", "bbb", 500, 1200))
	if b := layers["bbb"]; b.Status != "extracting" || b.Size != 1000 || b.Progress != 500 {
		t.Fatalf("解压中的层状态不对: %+v", *b)
	}
	applyPullProgress(layers, pullMsg("Pull complete", "bbb", 0, 0))
	if b := layers["bbb"]; b.Status != "complete" || b.Progress != b.Size {
		t.Fatalf("完成的层状态不对: %+v", *b)
	}
}
