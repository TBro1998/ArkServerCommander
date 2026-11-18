package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"ark-server-commander/utils"
)

// arkfgfRequest 向 arkfgf.com 的接口发 JSON 请求（发布到收录站和联机房共用）。
// body 为 nil 时不带请求体；out 为 nil 时不解析响应。site 是错误信息里的名称。
// 返回 HTTP 状态码、Retry-After、错误码（响应的 code 字段）和错误；状态码不是 200 时返回错误
func arkfgfRequest(client *http.Client, method, url, token string, body, out any, site string) (int, time.Duration, string, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return 0, 0, "", err
		}
	}
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return 0, 0, "", err
	}
	if body != nil || method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "ArkServerCommander/"+utils.Version)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, "", fmt.Errorf("连接%s失败: %w", site, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 256<<10))
	if err != nil {
		return resp.StatusCode, 0, "", fmt.Errorf("读取%s响应失败: %w", site, err)
	}
	if resp.StatusCode != http.StatusOK {
		var e struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}
		_ = json.Unmarshal(data, &e)
		retry := time.Duration(0)
		if secs, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && secs > 0 {
			retry = time.Duration(secs) * time.Second
		}
		msg := e.Error
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return resp.StatusCode, retry, e.Code, fmt.Errorf("%s返回 %d: %s", site, resp.StatusCode, msg)
	}
	if out != nil {
		if err := json.Unmarshal(data, out); err != nil {
			return resp.StatusCode, 0, "", fmt.Errorf("解析%s响应失败: %w", site, err)
		}
	}
	return resp.StatusCode, 0, "", nil
}
