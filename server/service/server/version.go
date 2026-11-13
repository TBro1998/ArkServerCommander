package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"ark-server-commander/utils"

	"go.uber.org/zap"
)

// VersionInfo 当前版本和最新发布版本
type VersionInfo struct {
	Current     string `json:"current"`
	Latest      string `json:"latest,omitempty"`       // 最新发布版本（不含预发布），查询失败时为空
	HasUpdate   bool   `json:"has_update"`             // 当前是开发版本（dev）时始终为 false
	ReleaseURL  string `json:"release_url,omitempty"`  // 最新版本的发布说明
	PublishedAt string `json:"published_at,omitempty"` // 最新版本的发布时间
	CheckedAt   string `json:"checked_at,omitempty"`
	Error       string `json:"error,omitempty"` // 查询失败的原因
}

var latestRelease struct {
	sync.Mutex
	tag, url, published string
	err                 string
	checked             time.Time
}

// GetVersionInfo 返回当前版本，并查询 GitHub 上的最新发布版本。成功结果缓存 6 小时，失败 10 分钟；refresh 时重新查询
func (s *ServerService) GetVersionInfo(refresh bool) VersionInfo {
	latestRelease.Lock()
	defer latestRelease.Unlock()

	ttl := 6 * time.Hour
	if latestRelease.err != "" {
		ttl = 10 * time.Minute
	}
	if refresh || latestRelease.checked.IsZero() || time.Since(latestRelease.checked) > ttl {
		tag, url, published, err := fetchLatestRelease()
		latestRelease.checked = time.Now()
		if err != nil {
			utils.Warn("检查新版本失败", zap.Error(err))
			latestRelease.err = err.Error()
		} else {
			latestRelease.tag, latestRelease.url, latestRelease.published, latestRelease.err = tag, url, published, ""
		}
	}

	info := VersionInfo{
		Current:     utils.Version,
		Latest:      strings.TrimPrefix(latestRelease.tag, "v"),
		ReleaseURL:  latestRelease.url,
		PublishedAt: latestRelease.published,
		CheckedAt:   latestRelease.checked.Format(time.RFC3339),
		Error:       latestRelease.err,
	}
	if info.Latest != "" {
		info.HasUpdate = newerVersion(info.Latest, info.Current)
	}
	return info
}

func fetchLatestRelease() (tag, url, published string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/"+utils.ReleaseRepo+"/releases/latest", nil)
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "ArkServerCommander/"+utils.Version)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("GitHub 接口返回 %d", resp.StatusCode)
	}
	var body struct {
		TagName     string `json:"tag_name"`
		HTMLURL     string `json:"html_url"`
		PublishedAt string `json:"published_at"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body); err != nil {
		return "", "", "", err
	}
	if body.TagName == "" {
		return "", "", "", fmt.Errorf("GitHub 没有返回版本号")
	}
	return body.TagName, body.HTMLURL, body.PublishedAt, nil
}

// semver 解析 x.y.z 或 x.y.z-pre（可带 v 前缀），解析失败返回 ok=false
func semver(v string) (nums [3]int, pre string, ok bool) {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		if v[i] == '-' {
			pre = v[i+1:]
			if j := strings.IndexByte(pre, '+'); j >= 0 {
				pre = pre[:j]
			}
		}
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return nums, "", false
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return nums, "", false
		}
		nums[i] = n
	}
	return nums, pre, true
}

// newerVersion latest 是否比 current 新。current 不是版本号（例如 dev）时返回 false。
// 预发布版本比同号的正式版本旧；两个预发布版本按字符串比较
func newerVersion(latest, current string) bool {
	l, lpre, ok1 := semver(latest)
	c, cpre, ok2 := semver(current)
	if !ok1 || !ok2 {
		return false
	}
	for i := 0; i < 3; i++ {
		if l[i] != c[i] {
			return l[i] > c[i]
		}
	}
	switch {
	case lpre == cpre:
		return false
	case lpre == "":
		return true
	case cpre == "":
		return false
	}
	return lpre > cpre
}
